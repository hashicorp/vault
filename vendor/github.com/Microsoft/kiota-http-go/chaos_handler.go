package nethttplibrary

import (
	"errors"
	"io"
	"math/rand"
	nethttp "net/http"
	"regexp"
	"strings"

	abstractions "github.com/microsoft/kiota-abstractions-go"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

type ChaosStrategy int

const (
	Manual ChaosStrategy = iota
	Random
)

// ChaosHandlerOptions is a configuration struct holding behavior defined options for a chaos handler
//
// BaseUrl represent the host url for in
// ChaosStrategy Specifies the strategy used for the Testing Handler -> RANDOM/MANUAL
// StatusCode Status code to be returned as part of the error response
// StatusMessage Message to be returned as part of the error response
// ChaosPercentage The percentage of randomness/chaos in the handler
// ResponseBody The response body to be returned as part of the error response
// Headers The response headers to be returned as part of the error response
// StatusMap The Map passed by user containing url-statusCode info
type ChaosHandlerOptions struct {
	BaseUrl         string
	ChaosStrategy   ChaosStrategy
	StatusCode      int
	StatusMessage   string
	ChaosPercentage int
	ResponseBody    *nethttp.Response
	Headers         map[string][]string
	StatusMap       map[string]map[string]int
}

type chaosHandlerOptionsInt interface {
	abstractions.RequestOption
	GetBaseUrl() string
	GetChaosStrategy() ChaosStrategy
	GetStatusCode() int
	GetStatusMessage() string
	GetChaosPercentage() int
	GetResponseBody() *nethttp.Response
	GetHeaders() map[string][]string
	GetStatusMap() map[string]map[string]int
}

func (handlerOptions *ChaosHandlerOptions) GetBaseUrl() string {
	return handlerOptions.BaseUrl
}

func (handlerOptions *ChaosHandlerOptions) GetChaosStrategy() ChaosStrategy {
	return handlerOptions.ChaosStrategy
}

func (handlerOptions *ChaosHandlerOptions) GetStatusCode() int {
	return handlerOptions.StatusCode
}

func (handlerOptions *ChaosHandlerOptions) GetStatusMessage() string {
	return handlerOptions.StatusMessage
}

func (handlerOptions *ChaosHandlerOptions) GetChaosPercentage() int {
	return handlerOptions.ChaosPercentage
}

func (handlerOptions *ChaosHandlerOptions) GetResponseBody() *nethttp.Response {
	return handlerOptions.ResponseBody
}

func (handlerOptions *ChaosHandlerOptions) GetHeaders() map[string][]string {
	return handlerOptions.Headers
}

func (handlerOptions *ChaosHandlerOptions) GetStatusMap() map[string]map[string]int {
	return handlerOptions.StatusMap
}

type ChaosHandler struct {
	options *ChaosHandlerOptions
}

var chaosHandlerKey = abstractions.RequestOptionKey{Key: "ChaosHandler"}

// GetKey returns ChaosHandlerOptions unique name in context object
func (handlerOptions *ChaosHandlerOptions) GetKey() abstractions.RequestOptionKey {
	return chaosHandlerKey
}

// NewChaosHandlerWithOptions creates a new ChaosHandler with the configured options
func NewChaosHandlerWithOptions(handlerOptions *ChaosHandlerOptions) (*ChaosHandler, error) {
	if handlerOptions == nil {
		return nil, errors.New("unexpected argument ChaosHandlerOptions as nil")
	}

	if handlerOptions.ChaosPercentage < 0 || handlerOptions.ChaosPercentage > 100 {
		return nil, errors.New("ChaosPercentage must be between 0 and 100")
	}
	if handlerOptions.ChaosStrategy == Manual {
		if handlerOptions.StatusCode == 0 {
			return nil, errors.New("invalid status code for manual strategy")
		}
	}

	return &ChaosHandler{options: handlerOptions}, nil
}

// NewChaosHandler creates a new ChaosHandler with default configuration options of Random errors at 10%
func NewChaosHandler() *ChaosHandler {
	return &ChaosHandler{
		options: &ChaosHandlerOptions{
			ChaosPercentage: 10,
			ChaosStrategy:   Random,
			StatusMessage:   "A random error message",
		},
	}
}

var methodStatusCode = map[string][]int{
	"GET":    {429, 500, 502, 503, 504},
	"POST":   {429, 500, 502, 503, 504, 507},
	"PUT":    {429, 500, 502, 503, 504, 507},
	"PATCH":  {429, 500, 502, 503, 504},
	"DELETE": {429, 500, 502, 503, 504, 507},
}

var httpStatusCode = map[int]string{
	100: "Continue",
	101: "Switching Protocols",
	102: "Processing",
	103: "Early Hints",
	200: "OK",
	201: "Created",
	202: "Accepted",
	203: "Non-Authoritative Information",
	204: "No Content",
	205: "Reset Content",
	206: "Partial Content",
	207: "Multi-Status",
	208: "Already Reported",
	226: "IM Used",
	300: "Multiple Choices",
	301: "Moved Permanently",
	302: "Found",
	303: "See Other",
	304: "Not Modified",
	305: "Use Proxy",
	307: "Temporary Redirect",
	308: "Permanent Redirect",
	400: "Bad Request",
	401: "Unauthorized",
	402: "Payment Required",
	403: "Forbidden",
	404: "Not Found",
	405: "Method Not Allowed",
	406: "Not Acceptable",
	407: "Proxy Authentication Required",
	408: "Request Timeout",
	409: "Conflict",
	410: "Gone",
	411: "Length Required",
	412: "Precondition Failed",
	413: "Payload Too Large",
	414: "URI Too Long",
	415: "Unsupported Media Type",
	416: "Range Not Satisfiable",
	417: "Expectation Failed",
	421: "Misdirected Request",
	422: "Unprocessable Entity",
	423: "Locked",
	424: "Failed Dependency",
	425: "Too Early",
	426: "Upgrade Required",
	428: "Precondition Required",
	429: "Too Many Requests",
	431: "Request Header Fields Too Large",
	451: "Unavailable For Legal Reasons",
	500: "Internal Server Error",
	501: "Not Implemented",
	502: "Bad Gateway",
	503: "Service Unavailable",
	504: "Gateway Timeout",
	505: "HTTP Version Not Supported",
	506: "Variant Also Negotiates",
	507: "Insufficient Storage",
	508: "Loop Detected",
	510: "Not Extended",
	511: "Network Authentication Required",
}

func generateRandomStatusCode(request *nethttp.Request) int {
	statusCodeArray := methodStatusCode[request.Method]
	return statusCodeArray[rand.Intn(len(statusCodeArray))]
}

func getRelativeURL(handlerOptions chaosHandlerOptionsInt, url string) string {
	baseUrl := handlerOptions.GetBaseUrl()
	if baseUrl != "" {
		return strings.Replace(url, baseUrl, "", 1)
	} else {
		return url
	}
}

func getStatusCode(handlerOptions chaosHandlerOptionsInt, req *nethttp.Request) int {
	requestMethod := req.Method
	statusMap := handlerOptions.GetStatusMap()
	requestURL := req.RequestURI

	if handlerOptions.GetChaosStrategy() == Manual {
		return handlerOptions.GetStatusCode()
	}

	if handlerOptions.GetChaosStrategy() == Random {
		if handlerOptions.GetStatusCode() > 0 {
			return handlerOptions.GetStatusCode()
		} else {
			relativeUrl := getRelativeURL(handlerOptions, requestURL)
			if definedResponses, ok := statusMap[relativeUrl]; ok {
				if mapCode, mapCodeOk := definedResponses[requestMethod]; mapCodeOk {
					return mapCode
				}
			} else {
				for key := range statusMap {
					match, _ := regexp.MatchString(key+"$", relativeUrl)
					if match {
						responseCode := statusMap[key][requestMethod]
						if responseCode != 0 {
							return responseCode
						}
					}
				}
			}
		}
	}

	return generateRandomStatusCode(req)
}

func createResponseBody(handlerOptions chaosHandlerOptionsInt, statusCode int) *nethttp.Response {
	if handlerOptions.GetResponseBody() != nil {
		return handlerOptions.GetResponseBody()
	}

	var stringReader *strings.Reader
	if statusCode > 400 {
		codeMessage := httpStatusCode[statusCode]
		errMessage := handlerOptions.GetStatusMessage()
		stringReader = strings.NewReader("error : { code :  " + codeMessage + " , message : " + errMessage + " }")
	} else {
		stringReader = strings.NewReader("{}")
	}

	return &nethttp.Response{
		StatusCode: statusCode,
		Status:     handlerOptions.GetStatusMessage(),
		Body:       io.NopCloser(stringReader),
		Header:     handlerOptions.GetHeaders(),
	}
}

func createChaosResponse(handler chaosHandlerOptionsInt, req *nethttp.Request) (*nethttp.Response, error) {
	statusCode := getStatusCode(handler, req)
	responseBody := createResponseBody(handler, statusCode)
	return responseBody, nil
}

// ChaosHandlerTriggeredEventKey is the key used for the open telemetry event
const ChaosHandlerTriggeredEventKey = "com.microsoft.kiota.chaos_handler_triggered"

func (middleware ChaosHandler) Intercept(pipeline Pipeline, middlewareIndex int, req *nethttp.Request) (*nethttp.Response, error) {
	reqOption, ok := req.Context().Value(chaosHandlerKey).(chaosHandlerOptionsInt)
	if !ok {
		reqOption = middleware.options
	}

	obsOptions := GetObservabilityOptionsFromRequest(req)
	ctx := req.Context()
	var span trace.Span
	if obsOptions != nil {
		ctx, span = otel.GetTracerProvider().Tracer(obsOptions.GetTracerInstrumentationName()).Start(ctx, "ChaosHandler_Intercept")
		span.SetAttributes(attribute.Bool("com.microsoft.kiota.handler.chaos.enable", true))
		req = req.WithContext(ctx)
		defer span.End()
	}

	if rand.Intn(100) < reqOption.GetChaosPercentage() {
		if span != nil {
			span.AddEvent(ChaosHandlerTriggeredEventKey)
		}
		return createChaosResponse(reqOption, req)
	}

	return pipeline.Next(req, middlewareIndex)
}
