package msgraphgocore

import (
	nethttp "net/http"

	runtime "runtime"

	uuid "github.com/google/uuid"
	khttp "github.com/microsoft/kiota-http-go"
)

// GraphTelemetryHandler is a middleware handler that adds telemetry headers to requests.
type GraphTelemetryHandler struct {
	sdkVersion string
}

// NewGraphTelemetryHandler creates a new GraphTelemetryHandler.
func NewGraphTelemetryHandler(options *GraphClientOptions) *GraphTelemetryHandler {
	serviceVersionPrefix := ""
	if options != nil && options.GraphServiceLibraryVersion != "" {
		serviceVersionPrefix += "graph-go"
		if options.GraphServiceVersion != "" {
			serviceVersionPrefix += "-" + options.GraphServiceVersion
		}
		serviceVersionPrefix += "/" + options.GraphServiceLibraryVersion
		serviceVersionPrefix += ", "
	}
	featuresSuffix := ""
	if runtime.GOOS != "" {
		featuresSuffix += " hostOS=" + runtime.GOOS + ";"
	}
	if runtime.GOARCH != "" {
		featuresSuffix += " hostArch=" + runtime.GOARCH + ";"
	}
	goVersion := runtime.Version()
	if goVersion != "" {
		featuresSuffix += " runtimeEnvironment=" + goVersion + ";"
	}
	if featuresSuffix != "" {
		featuresSuffix = " (" + featuresSuffix[1:] + ")"
	}
	return &GraphTelemetryHandler{
		sdkVersion: serviceVersionPrefix + "graph-go-core/" + CoreVersion + featuresSuffix,
	}
}
func (middleware GraphTelemetryHandler) Intercept(pipeline khttp.Pipeline, middlewareIndex int, req *nethttp.Request) (*nethttp.Response, error) {
	req.Header.Add("SdkVersion", middleware.sdkVersion)
	req.Header.Add("client-request-id", uuid.NewString())
	return pipeline.Next(req, middlewareIndex)
}
