package main

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"time"

	"github.com/cubex/potens-go/adl"
	"github.com/cubex/potens-go/application"
	"github.com/cubex/proto-go/applications"
	zipkin "github.com/openzipkin/zipkin-go-opentracing"
	"github.com/uber-go/zap"
	"golang.org/x/net/context"
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
	applications.RegisterApplicationServer(grpcServer, &server{})

	//Mark service as online and start heart beat
	err = app.Online()
	if err != nil {
		app.Logger.Fatal(err.Error())
	}

	// --- FDL TEST LOOP START ---
	for {
		ent1 := app.ADL("THIS-IS-A-FID")
		err = ent1.Read(adl.Counter("propX"))
		countA := ent1.GetCounter("propX")
		testStrStart := "this is data!"

		if countA <= 0 {
			// Write
			ent1.Write("propX", testStrStart)
		} else {
			str := testStrStart
			for i := 0; i < countA; i++ {
				str += "!"
			}
			ent1.Write("propX", str)
		}

		ent1.WriteMeta("propX", "this is meta")

		ent1.AddSetItem("propX", "test1")
		ent1.AddSetItem("propX", "test2")
		ent1.AddSetItem("propX", "test3")
		//ent1.RemoveSetItem("propX", "test2")

		// Test List Add
		testListName := "TESTLIST"
		ent1.AddListItem(testListName, "1", "ONE")
		ent1.AddListItem(testListName, "2", "TWO")
		ent1.AddListItem(testListName, "3", "THREE")
		ent1.AddListItem(testListName, "4", "FOUR")

		ent1.IncrementCounter("propX")
		ent1.Write("propY", "This is Y data")
		ent1.Commit()

		// TODO: retrieve items with prefix
		ent1.Read(adl.PropertiesWithPrefix("prop"), adl.Meta("propX"), adl.Counter("propX"), adl.Set("propX"), adl.ListRange(testListName, "1", "", 0))
		dataA := ent1.Get("propX")
		countA = ent1.GetCounter("propX")
		data3 := ent1.GetSet("propX")
		data4 := ent1.Get("propY")
		meta1 := ent1.GetMeta("propX")
		list := ent1.GetList(testListName)

		fmt.Printf("\nItem:%s-%s-%s\n", dataA, data4, meta1)
		fmt.Printf("Counter:%d\n", countA)

		for _, d := range data3 {
			fmt.Printf("SET-ITEM:%s\n", d)
		}

		for _, d := range list {
			fmt.Printf("LIST-ITEM:%s - %s\n", d.Key, d.Value)
		}

		if len(list) == 0 {
			fmt.Printf("No list items returned\n")
		}

		t := time.Second * 5
		time.Sleep(t)
	}
	// --- FDL TEST LOOP END ---

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
func (s *server) HandleHTTPRequest(ctx context.Context, in *applications.HTTPRequest) (*applications.HTTPResponse, error) {
	if in.RequestType == applications.HTTPRequest_PAGE_DEFINITION {
		return s.PageDefinition(ctx, in)
	} else if in.RequestType == applications.HTTPRequest_REQUEST {
		return s.HTTPResource(ctx, in)
	}

	return nil, errors.New("Unsupported request type")
}
