package main

import (
	"bytes"
	"errors"
	"html/template"
	"os"
	"regexp"
	"strings"

	"github.com/cubex/portcullis-go"
	"github.com/cubex/potens-go/webui"
	"github.com/cubex/potens-go/webui/breadcrumb"
	"github.com/cubex/proto-go/applications"
	"github.com/uber-go/zap"
	"golang.org/x/net/context"
)

func (s *server) PageDefinition(ctx context.Context, in *applications.HTTPRequest) (*applications.HTTPResponse, error) {
	//Retrieve auth data
	response := webui.CreateResponse()
	authData := portcullis.FromContext(ctx)
	app.Logger.Debug("Received HTTP Request", zap.String("userID", authData.UserID), zap.String("projectID", authData.ProjectID))

	if in.Path == "/" {
		webui.SetPageTitle(response, "List Of Things")
	} else {

		// Load the hero template for the page if one exists
		if _, err := os.Stat("heroTemplates" + in.Path + ".html"); err == nil {
			defaultTemplate := template.Must(template.ParseFiles("heroTemplates" + in.Path + ".html"))
			buf := new(bytes.Buffer)
			defaultTemplate.Execute(buf, "")
			response.Body = buf.String()
		}

		pageKey := in.Path[1:]
		r, _ := regexp.Compile("(\\w)(\\w*)")
		pageTitle := r.ReplaceAllStringFunc(strings.Replace(pageKey, "-", " ", -1), func(str string) string {
			return strings.ToUpper(str[0:1]) + strings.ToLower(str[1:])
		})

		bread := breadcrumb.Breadcrumb{}
		bread.AddItem(breadcrumb.BreadcrumbItem{Url: in.Path, Title: pageTitle})
		webui.SetBackPath(response, "/")
		webui.SetBreadcrumb(response, bread)
		webui.SetPageTitle(response, pageTitle)
		webui.SetPageFID(response, pageKey)
	}

	return response, nil
}

// HandleHTTPRequest handles requests from HTTP sources
func (s *server) HTTPResource(ctx context.Context, in *applications.HTTPRequest) (*applications.HTTPResponse, error) {
	//Retrieve auth data
	response := webui.CreateResponse()

	defaultTemplate := template.New("")
	if in.Path == "/" {
		defaultTemplate = template.Must(template.ParseFiles("templates/default.html"))
	} else if strings.HasSuffix(in.Path, "/details") {
		defaultTemplate = template.Must(template.ParseFiles("templates/example/details.html"))
	} else {
		defaultTemplate = template.Must(template.ParseFiles("templates/example/example.html"))
	}

	buf := new(bytes.Buffer)
	defaultTemplate.Execute(buf, "")
	response.Body = buf.String()

	return response, nil
}

func (s *server) HandleSocketAction(ctx context.Context, in *applications.SocketRequest) (*applications.HTTPResponse, error) {
	return nil, errors.New("Sockets not supported")
}

func (s *server) ModifyRelationship(context.Context, *applications.ProjectModifyRequest) (*applications.ProjectModifyResponse, error) {
	return nil, errors.New("Sockets not supported")
}
