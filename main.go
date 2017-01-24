package main

import (
	"golang.org/x/net/context"

	"github.com/cubex/proto-go/platform"

	"flag"

	"os"
	"os/signal"

	"errors"

	zipkin "github.com/openzipkin/zipkin-go-opentracing"
	"github.com/uber-go/zap"
	"github.com/cubex/potens-go/application"
)

type server struct{}

var (
	app        application.CubexApplication
	zipkinHttp = flag.String("zipkin-http", "http://localhost:9411/api/v1/spans", "Zipkin HTTP Collector URL")
)

func main() {

	//Get a new service
	var err error
	app, err = application.New(nil, nil)
	if err != nil {
		app.Logger.Fatal(err.Error())
	}

	//Retrieve TLS Cert and Register with the discovery service
	if *zipkinHttp != "" {
		collector, err := zipkin.NewHTTPCollector(*zipkinHttp)
		defer collector.Close()
		if err != nil {
			app.Logger.Fatal("unable to create Zipkin HTTP collector", zap.Error(err))
		}
		err = app.Start(collector)
	} else {
		err = app.Start(nil)
	}

	if err != nil {
		app.Logger.Fatal(err.Error())
	}

	//Create your gRPC Server with the correct certificates
	lis, grpcServer, err := app.CreateServer()
	if err != nil {
		app.Logger.Fatal(err.Error())
	}

	/*
	 * Start up your service here
	 */
	platform.RegisterCubexPlatformServer(grpcServer, &server{})

	//Mark service as online and start heart beat
	err = app.Online()
	if err != nil {
		app.Logger.Fatal(err.Error())
	}

	//When interrupted, take the app offline with discovery first
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		for range c {
			app.Close()
			os.Exit(1)
		}
	}()

	//Listen for connections
	err = grpcServer.Serve(lis)

	if err != nil {
		//Mark the gservice instance as offline
		app.Close()
		app.Logger.Fatal(err.Error())
	}
}

// HandleHTTPRequest handles requests from HTTP sources
func (s *server) HandleHTTPRequest(ctx context.Context, in *platform.HTTPRequest) (*platform.HTTPResponse, error) {
	if in.RequestType == platform.HTTPRequest_PAGE_DEFINITION {
		return s.PageDefinition(ctx, in)
	} else if in.RequestType == platform.HTTPRequest_REQUEST {
		return s.HTTPResource(ctx, in)
	}

	return nil, errors.New("Unsupported request type")
}
