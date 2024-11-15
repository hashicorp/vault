package common

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/opentracing/opentracing-go/ext"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/denverdino/aliyungo/util"
	"github.com/opentracing/opentracing-go"
)

// RemovalPolicy.N add index to array item
// RemovalPolicy=["a", "b"] => RemovalPolicy.1="a" RemovalPolicy.2="b"
type FlattenArray []string

// string contains underline which will be replaced with dot
// SystemDisk_Category => SystemDisk.Category
type UnderlineString string

// A Client represents a client of ECS services
type Client struct {
	AccessKeyId     string //Access Key Id
	AccessKeySecret string //Access Key Secret
	securityToken   string
	debug           bool
	httpClient      *http.Client
	endpoint        string
	version         string
	serviceCode     string
	regionID        Region
	businessInfo    string
	userAgent       string
	disableTrace    bool
	span            opentracing.Span
	logger          *Logger
}

// Initialize properties of a client instance
func (client *Client) Init(endpoint, version, accessKeyId, accessKeySecret string) {
	client.AccessKeyId = accessKeyId
	ak := accessKeySecret
	if !strings.HasSuffix(ak, "&") {
		ak += "&"
	}
	client.AccessKeySecret = ak
	client.InitClient()
	client.endpoint = endpoint
	client.version = version
}

// Initialize properties of a client instance including regionID
func (client *Client) NewInit(endpoint, version, accessKeyId, accessKeySecret, serviceCode string, regionID Region) {
	client.Init(endpoint, version, accessKeyId, accessKeySecret)
	client.serviceCode = serviceCode
	client.regionID = regionID
}

// Initialize properties of a client instance including regionID
//only for hz regional Domain
func (client *Client) NewInit4RegionalDomain(endpoint, version, accessKeyId, accessKeySecret, serviceCode string, regionID Region) {
	client.Init(endpoint, version, accessKeyId, accessKeySecret)
	client.serviceCode = serviceCode
	client.regionID = regionID

	client.setEndpoint4RegionalDomain(client.regionID, client.serviceCode, client.AccessKeyId, client.AccessKeySecret, client.securityToken)
}

// Intialize client object when all properties are ready
func (client *Client) InitClient() *Client {
	client.debug = false

	// create DefaultTransport manully, because transport doesn't has clone method in go 1.10
	t := &http.Transport{
		Proxy: http.ProxyFromEnvironment,
		DialContext: (&net.Dialer{
			Timeout:   30 * time.Second,
			KeepAlive: 30 * time.Second,
			DualStack: true,
		}).DialContext,
		MaxIdleConns:          100,
		IdleConnTimeout:       90 * time.Second,
		TLSHandshakeTimeout:   10 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
	}

	handshakeTimeoutStr, ok := os.LookupEnv("TLSHandshakeTimeout")
	if ok {
		handshakeTimeout, err := strconv.Atoi(handshakeTimeoutStr)
		if err != nil {
			log.Printf("Get TLSHandshakeTimeout from env error: %v.", err)
		} else {
			t.TLSHandshakeTimeout = time.Duration(handshakeTimeout) * time.Second
		}
	}

	responseHeaderTimeoutStr, ok := os.LookupEnv("ResponseHeaderTimeout")
	if ok {
		responseHeaderTimeout, err := strconv.Atoi(responseHeaderTimeoutStr)
		if err != nil {
			log.Printf("Get ResponseHeaderTimeout from env error: %v.", err)
		} else {
			t.ResponseHeaderTimeout = time.Duration(responseHeaderTimeout) * time.Second
		}
	}

	expectContinueTimeoutStr, ok := os.LookupEnv("ExpectContinueTimeout")
	if ok {
		expectContinueTimeout, err := strconv.Atoi(expectContinueTimeoutStr)
		if err != nil {
			log.Printf("Get ExpectContinueTimeout from env error: %v.", err)
		} else {
			t.ExpectContinueTimeout = time.Duration(expectContinueTimeout) * time.Second
		}
	}

	idleConnTimeoutStr, ok := os.LookupEnv("IdleConnTimeout")
	if ok {
		idleConnTimeout, err := strconv.Atoi(idleConnTimeoutStr)
		if err != nil {
			log.Printf("Get IdleConnTimeout from env error: %v.", err)
		} else {
			t.IdleConnTimeout = time.Duration(idleConnTimeout) * time.Second
		}
	}

	client.httpClient = &http.Client{
		Transport: t,
	}

	httpTimeoutStr, ok := os.LookupEnv("HttpTimeout")
	if ok {
		httpTimeout, err := strconv.Atoi(httpTimeoutStr)
		if err != nil {
			log.Printf("Get HttpTimeout from env error: %v.", err)
		} else {
			client.httpClient.Timeout = time.Duration(httpTimeout) * time.Second
		}
	}

	return client
}

// Intialize client object when all properties are ready
//only for regional domain hz
func (client *Client) InitClient4RegionalDomain() *Client {
	client.InitClient()
	//set endpoint
	client.setEndpoint4RegionalDomain(client.regionID, client.serviceCode, client.AccessKeyId, client.AccessKeySecret, client.securityToken)
	return client
}

func (client *Client) NewInitForAssumeRole(endpoint, version, accessKeyId, accessKeySecret, serviceCode string, regionID Region, securityToken string) {
	client.NewInit(endpoint, version, accessKeyId, accessKeySecret, serviceCode, regionID)
	client.securityToken = securityToken
}

//getLocationEndpoint
func (client *Client) getEndpointByLocation() string {
	locationClient := NewLocationClient(client.AccessKeyId, client.AccessKeySecret, client.securityToken)
	locationClient.SetDebug(true)
	return locationClient.DescribeOpenAPIEndpoint(client.regionID, client.serviceCode)
}

//NewClient using location service
func (client *Client) setEndpointByLocation(region Region, serviceCode, accessKeyId, accessKeySecret, securityToken string) {
	locationClient := NewLocationClient(accessKeyId, accessKeySecret, securityToken)
	locationClient.SetDebug(true)
	ep := locationClient.DescribeOpenAPIEndpoint(region, serviceCode)

	if ep != "" {
		client.endpoint = ep
	}
}

// Get openapi endpoint accessed by ecs instance.
// For some UnitRegions, the endpoint pattern is https://[product].[regionid].aliyuncs.com
// For some CentralRegions, the endpoint pattern is  https://[product].vpc-proxy.aliyuncs.com
// The other region, the endpoint pattern is https://[product]-vpc.[regionid].aliyuncs.com
func (client *Client) setEndpoint4RegionalDomain(region Region, serviceCode, accessKeyId, accessKeySecret, securityToken string) {
	if endpoint, ok := CentralDomainServices[serviceCode]; ok {
		client.endpoint = fmt.Sprintf("https://%s", endpoint)
		return
	}
	for _, service := range RegionalDomainServices {
		if service == serviceCode {
			if ep, ok := UnitRegions[region]; ok {
				client.endpoint = fmt.Sprintf("https://%s.%s.aliyuncs.com", serviceCode, ep)
				return
			}

			client.endpoint = fmt.Sprintf("https://%s%s.%s.aliyuncs.com", serviceCode, "-vpc", region)
			return
		}
	}
	locationClient := NewLocationClient(accessKeyId, accessKeySecret, securityToken)
	locationClient.SetDebug(true)
	ep := locationClient.DescribeOpenAPIEndpoint(region, serviceCode)

	if ep != "" {
		client.endpoint = ep
	}
}

// Ensure all necessary properties are valid
func (client *Client) ensureProperties() error {
	var msg string

	if client.endpoint == "" {
		msg = fmt.Sprintf("endpoint cannot be empty!")
	} else if client.version == "" {
		msg = fmt.Sprintf("version cannot be empty!")
	} else if client.AccessKeyId == "" {
		msg = fmt.Sprintf("AccessKeyId cannot be empty!")
	} else if client.AccessKeySecret == "" {
		msg = fmt.Sprintf("AccessKeySecret cannot be empty!")
	}

	if msg != "" {
		return errors.New(msg)
	}

	return nil
}

// ----------------------------------------------------
// WithXXX methods
// ----------------------------------------------------

// WithEndpoint sets custom endpoint
func (client *Client) WithEndpoint(endpoint string) *Client {
	client.SetEndpoint(endpoint)
	return client
}

// WithVersion sets custom version
func (client *Client) WithVersion(version string) *Client {
	client.SetVersion(version)
	return client
}

// WithRegionID sets Region ID
func (client *Client) WithRegionID(regionID Region) *Client {
	client.SetRegionID(regionID)
	return client
}

//WithServiceCode sets serviceCode
func (client *Client) WithServiceCode(serviceCode string) *Client {
	client.SetServiceCode(serviceCode)
	return client
}

// WithAccessKeyId sets new AccessKeyId
func (client *Client) WithAccessKeyId(id string) *Client {
	client.SetAccessKeyId(id)
	return client
}

// WithAccessKeySecret sets new AccessKeySecret
func (client *Client) WithAccessKeySecret(secret string) *Client {
	client.SetAccessKeySecret(secret)
	return client
}

// WithSecurityToken sets securityToken
func (client *Client) WithSecurityToken(securityToken string) *Client {
	client.SetSecurityToken(securityToken)
	return client
}

// WithDebug sets debug mode to log the request/response message
func (client *Client) WithDebug(debug bool) *Client {
	client.SetDebug(debug)
	return client
}

// WithBusinessInfo sets business info to log the request/response message
func (client *Client) WithBusinessInfo(businessInfo string) *Client {
	client.SetBusinessInfo(businessInfo)
	return client
}

// WithUserAgent sets user agent to the request/response message
func (client *Client) WithUserAgent(userAgent string) *Client {
	client.SetUserAgent(userAgent)
	return client
}

// WithUserAgent sets user agent to the request/response message
func (client *Client) WithDisableTrace(disableTrace bool) *Client {
	client.SetDisableTrace(disableTrace)
	return client
}

// WithUserAgent sets user agent to the request/response message
func (client *Client) WithSpan(span opentracing.Span) *Client {
	client.SetSpan(span)
	return client
}

// ----------------------------------------------------
// SetXXX methods
// ----------------------------------------------------

// SetEndpoint sets custom endpoint
func (client *Client) SetEndpoint(endpoint string) {
	client.endpoint = endpoint
}

func (client *Client) GetEndpoint() string {
	return client.endpoint
}

// SetEndpoint sets custom version
func (client *Client) SetVersion(version string) {
	client.version = version
}

// SetEndpoint sets Region ID
func (client *Client) SetRegionID(regionID Region) {
	client.regionID = regionID
}

//SetServiceCode sets serviceCode
func (client *Client) SetServiceCode(serviceCode string) {
	client.serviceCode = serviceCode
}

// SetAccessKeyId sets new AccessKeyId
func (client *Client) SetAccessKeyId(id string) {
	client.AccessKeyId = id
}

// SetAccessKeySecret sets new AccessKeySecret
func (client *Client) SetAccessKeySecret(secret string) {
	client.AccessKeySecret = secret + "&"
}

// SetDebug sets debug mode to log the request/response message
func (client *Client) SetDebug(debug bool) {
	client.debug = debug
}

// SetBusinessInfo sets business info to log the request/response message
func (client *Client) SetBusinessInfo(businessInfo string) {
	if strings.HasPrefix(businessInfo, "/") {
		client.businessInfo = businessInfo
	} else if businessInfo != "" {
		client.businessInfo = "/" + businessInfo
	}
}

// SetUserAgent sets user agent to the request/response message
func (client *Client) SetUserAgent(userAgent string) {
	client.userAgent = userAgent
}

//set SecurityToken
func (client *Client) SetSecurityToken(securityToken string) {
	client.securityToken = securityToken
}

// SetTransport sets transport to the http client
func (client *Client) SetTransport(transport http.RoundTripper) {
	if client.httpClient == nil {
		client.httpClient = &http.Client{}
	}
	client.httpClient.Transport = transport
}

// SetDisableTrace close trace mode
func (client *Client) SetDisableTrace(disableTrace bool) {
	client.disableTrace = disableTrace
}

// SetSpan set the parent span
func (client *Client) SetSpan(span opentracing.Span) {
	client.span = span
}

func (client *Client) initEndpoint() error {
	// if set any value to "CUSTOMIZED_ENDPOINT" could skip location service.
	// example: export CUSTOMIZED_ENDPOINT=true
	if os.Getenv("CUSTOMIZED_ENDPOINT") != "" {
		return nil
	}

	if client.endpoint != "" {
		return nil
	}

	if client.serviceCode != "" && client.regionID != "" {
		endpoint := client.getEndpointByLocation()
		if endpoint == "" {
			return GetCustomError("InvalidEndpoint", "endpoint is empty,pls check")
		}
		client.endpoint = endpoint
	}
	return nil
}

// Invoke sends the raw HTTP request for ECS services
func (client *Client) Invoke(action string, args interface{}, response interface{}) (err error) {
	if err := client.ensureProperties(); err != nil {
		return err
	}

	// log request
	fieldMap := make(map[string]string)
	initLogMsg(fieldMap)
	defer func() {
		client.printLog(fieldMap, err)
	}()

	request := Request{}
	request.init(client.version, action, client.AccessKeyId, client.securityToken, client.regionID)

	query := util.ConvertToQueryValues(request)
	util.SetQueryValues(args, &query)

	// Sign request
	signature := util.CreateSignatureForRequest(ECSRequestMethod, &query, client.AccessKeySecret)

	// Generate the request URL
	requestURL := client.endpoint + "?" + query.Encode() + "&Signature=" + url.QueryEscape(signature)

	httpReq, err := http.NewRequest(ECSRequestMethod, requestURL, nil)

	if err != nil {
		return GetClientError(err)
	}

	// TODO move to util and add build val flag
	httpReq.Header.Set("X-SDK-Client", `AliyunGO/`+Version+client.businessInfo)
	httpReq.Header.Set("User-Agent", httpReq.UserAgent()+" "+client.userAgent)

	// Set tracer
	var span opentracing.Span
	if ok := opentracing.IsGlobalTracerRegistered(); ok && !client.disableTrace {
		tracer := opentracing.GlobalTracer()
		var rootCtx opentracing.SpanContext

		if client.span != nil {
			rootCtx = client.span.Context()
		}

		span = tracer.StartSpan(
			"AliyunGO-"+request.Action,
			opentracing.ChildOf(rootCtx),
			opentracing.Tag{string(ext.Component), "AliyunGO"},
			opentracing.Tag{"ActionName", request.Action})

		defer span.Finish()
		tracer.Inject(
			span.Context(),
			opentracing.HTTPHeaders,
			opentracing.HTTPHeadersCarrier(httpReq.Header))
	}

	putMsgToMap(fieldMap, httpReq)
	t0 := time.Now()
	fieldMap["{start_time}"] = t0.Format("2006-01-02 15:04:05")
	httpResp, err := client.httpClient.Do(httpReq)
	t1 := time.Now()
	fieldMap["{cost}"] = t1.Sub(t0).String()
	if err != nil {
		if span != nil {
			ext.LogError(span, err)
		}
		return GetClientError(err)
	}
	fieldMap["{code}"] = strconv.Itoa(httpResp.StatusCode)
	fieldMap["{res_headers}"] = TransToString(httpResp.Header)
	statusCode := httpResp.StatusCode

	if client.debug {
		log.Printf("Invoke %s %s %d (%v)", ECSRequestMethod, requestURL, statusCode, t1.Sub(t0))
	}

	if span != nil {
		ext.HTTPStatusCode.Set(span, uint16(httpResp.StatusCode))
	}

	defer httpResp.Body.Close()
	body, err := ioutil.ReadAll(httpResp.Body)
	fieldMap["{res_body}"] = string(body)

	if err != nil {
		return GetClientError(err)
	}

	if client.debug {
		var prettyJSON bytes.Buffer
		err = json.Indent(&prettyJSON, body, "", "    ")
		if err != nil {
			log.Printf("Failed in json.Indent: %v\n", err)
		} else {
			log.Printf("JSON body: %s\n", prettyJSON.String())
		}
	}

	if statusCode >= 400 && statusCode <= 599 {
		errorResponse := ErrorResponse{}
		err = json.Unmarshal(body, &errorResponse)
		if err != nil {
			log.Printf("Failed in json.Unmarshal: %v\n", err)
		}
		ecsError := &Error{
			ErrorResponse: errorResponse,
			StatusCode:    statusCode,
		}
		return ecsError
	}

	err = json.Unmarshal(body, response)
	//log.Printf("%++v", response)
	if err != nil {
		return GetClientError(err)
	}

	return nil
}

// Invoke sends the raw HTTP request for ECS services
func (client *Client) InvokeByFlattenMethod(action string, args interface{}, response interface{}) (err error) {
	if err := client.ensureProperties(); err != nil {
		return err
	}

	// log request
	fieldMap := make(map[string]string)
	initLogMsg(fieldMap)
	defer func() {
		client.printLog(fieldMap, err)
	}()

	//init endpoint
	if err := client.initEndpoint(); err != nil {
		return err
	}

	request := Request{}
	request.init(client.version, action, client.AccessKeyId, client.securityToken, client.regionID)

	query := util.ConvertToQueryValues(request)

	util.SetQueryValueByFlattenMethod(args, &query)

	// Sign request
	signature := util.CreateSignatureForRequest(ECSRequestMethod, &query, client.AccessKeySecret)

	// Generate the request URL
	requestURL := client.endpoint + "?" + query.Encode() + "&Signature=" + url.QueryEscape(signature)

	httpReq, err := http.NewRequest(ECSRequestMethod, requestURL, nil)

	if err != nil {
		return GetClientError(err)
	}

	// TODO move to util and add build val flag
	httpReq.Header.Set("X-SDK-Client", `AliyunGO/`+Version+client.businessInfo)
	httpReq.Header.Set("User-Agent", httpReq.UserAgent()+" "+client.userAgent)

	// Set tracer
	var span opentracing.Span
	if ok := opentracing.IsGlobalTracerRegistered(); ok && !client.disableTrace {
		tracer := opentracing.GlobalTracer()
		var rootCtx opentracing.SpanContext

		if client.span != nil {
			rootCtx = client.span.Context()
		}

		span = tracer.StartSpan(
			"AliyunGO-"+request.Action,
			opentracing.ChildOf(rootCtx),
			opentracing.Tag{string(ext.Component), "AliyunGO"},
			opentracing.Tag{"ActionName", request.Action})

		defer span.Finish()
		tracer.Inject(
			span.Context(),
			opentracing.HTTPHeaders,
			opentracing.HTTPHeadersCarrier(httpReq.Header))
	}

	putMsgToMap(fieldMap, httpReq)
	t0 := time.Now()
	fieldMap["{start_time}"] = t0.Format("2006-01-02 15:04:05")
	httpResp, err := client.httpClient.Do(httpReq)
	t1 := time.Now()
	fieldMap["{cost}"] = t1.Sub(t0).String()
	if err != nil {
		if span != nil {
			ext.LogError(span, err)
		}
		return GetClientError(err)
	}
	fieldMap["{code}"] = strconv.Itoa(httpResp.StatusCode)
	fieldMap["{res_headers}"] = TransToString(httpResp.Header)
	statusCode := httpResp.StatusCode

	if client.debug {
		log.Printf("Invoke %s %s %d (%v)", ECSRequestMethod, requestURL, statusCode, t1.Sub(t0))
	}

	if span != nil {
		ext.HTTPStatusCode.Set(span, uint16(httpResp.StatusCode))
	}

	defer httpResp.Body.Close()
	body, err := ioutil.ReadAll(httpResp.Body)
	fieldMap["{res_body}"] = string(body)

	if err != nil {
		return GetClientError(err)
	}

	if client.debug {
		var prettyJSON bytes.Buffer
		err = json.Indent(&prettyJSON, body, "", "    ")
		if err != nil {
			log.Printf("Failed in json.Indent: %v\n", err)
		}
		log.Println(prettyJSON.String())
	}

	if statusCode >= 400 && statusCode <= 599 {
		errorResponse := ErrorResponse{}
		err = json.Unmarshal(body, &errorResponse)
		if err != nil {
			log.Printf("Failed in json.Unmarshal: %v\n", err)
		}
		ecsError := &Error{
			ErrorResponse: errorResponse,
			StatusCode:    statusCode,
		}
		return ecsError
	}

	err = json.Unmarshal(body, response)
	//log.Printf("%++v", response)
	if err != nil {
		return GetClientError(err)
	}

	return nil
}

// Invoke sends the raw HTTP request for ECS services
//改进了一下上面那个方法，可以使用各种Http方法
//2017.1.30 增加了一个path参数，用来拓展访问的地址
func (client *Client) InvokeByAnyMethod(method, action, path string, args interface{}, response interface{}) (err error) {
	if err := client.ensureProperties(); err != nil {
		return err
	}

	// log request
	fieldMap := make(map[string]string)
	initLogMsg(fieldMap)
	defer func() {
		client.printLog(fieldMap, err)
	}()

	//init endpoint
	//if err := client.initEndpoint(); err != nil {
	//	return err
	//}

	request := Request{}
	request.init(client.version, action, client.AccessKeyId, client.securityToken, client.regionID)
	data := util.ConvertToQueryValues(request)
	util.SetQueryValues(args, &data)

	// Sign request
	signature := util.CreateSignatureForRequest(method, &data, client.AccessKeySecret)

	data.Add("Signature", signature)
	// Generate the request URL
	var (
		httpReq *http.Request
	)
	if method == http.MethodGet {
		requestURL := client.endpoint + path + "?" + data.Encode()
		//fmt.Println(requestURL)
		httpReq, err = http.NewRequest(method, requestURL, nil)
	} else {
		//fmt.Println(client.endpoint + path)
		httpReq, err = http.NewRequest(method, client.endpoint+path, strings.NewReader(data.Encode()))
		httpReq.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	}

	if err != nil {
		return GetClientError(err)
	}

	// TODO move to util and add build val flag
	httpReq.Header.Set("X-SDK-Client", `AliyunGO/`+Version+client.businessInfo)
	httpReq.Header.Set("User-Agent", httpReq.Header.Get("User-Agent")+" "+client.userAgent)

	// Set tracer
	var span opentracing.Span
	if ok := opentracing.IsGlobalTracerRegistered(); ok && !client.disableTrace {
		tracer := opentracing.GlobalTracer()
		var rootCtx opentracing.SpanContext

		if client.span != nil {
			rootCtx = client.span.Context()
		}

		span = tracer.StartSpan(
			"AliyunGO-"+request.Action,
			opentracing.ChildOf(rootCtx),
			opentracing.Tag{string(ext.Component), "AliyunGO"},
			opentracing.Tag{"ActionName", request.Action})

		defer span.Finish()
		tracer.Inject(
			span.Context(),
			opentracing.HTTPHeaders,
			opentracing.HTTPHeadersCarrier(httpReq.Header))
	}

	putMsgToMap(fieldMap, httpReq)
	t0 := time.Now()
	fieldMap["{start_time}"] = t0.Format("2006-01-02 15:04:05")
	httpResp, err := client.httpClient.Do(httpReq)
	t1 := time.Now()
	fieldMap["{cost}"] = t1.Sub(t0).String()
	if err != nil {
		if span != nil {
			ext.LogError(span, err)
		}
		return GetClientError(err)
	}
	fieldMap["{code}"] = strconv.Itoa(httpResp.StatusCode)
	fieldMap["{res_headers}"] = TransToString(httpResp.Header)
	statusCode := httpResp.StatusCode

	if client.debug {
		log.Printf("Invoke %s %s %d (%v) %v", ECSRequestMethod, client.endpoint, statusCode, t1.Sub(t0), data.Encode())
	}

	if span != nil {
		ext.HTTPStatusCode.Set(span, uint16(httpResp.StatusCode))
	}

	defer httpResp.Body.Close()
	body, err := ioutil.ReadAll(httpResp.Body)
	fieldMap["{res_body}"] = string(body)

	if err != nil {
		return GetClientError(err)
	}

	if client.debug {
		var prettyJSON bytes.Buffer
		err = json.Indent(&prettyJSON, body, "", "    ")
		log.Println(prettyJSON.String())
	}

	if statusCode >= 400 && statusCode <= 599 {
		errorResponse := ErrorResponse{}
		err = json.Unmarshal(body, &errorResponse)
		ecsError := &Error{
			ErrorResponse: errorResponse,
			StatusCode:    statusCode,
		}
		return ecsError
	}

	err = json.Unmarshal(body, response)
	//log.Printf("%++v", response)
	if err != nil {
		return GetClientError(err)
	}

	return nil
}

// GenerateClientToken generates the Client Token with random string
func (client *Client) GenerateClientToken() string {
	return util.CreateRandomString()
}

func GetClientErrorFromString(str string) error {
	return &Error{
		ErrorResponse: ErrorResponse{
			Code:    "AliyunGoClientFailure",
			Message: str,
		},
		StatusCode: -1,
	}
}

func GetClientError(err error) error {
	return GetClientErrorFromString(err.Error())
}

func GetCustomError(code, message string) error {
	return &Error{
		ErrorResponse: ErrorResponse{
			Code:    code,
			Message: message,
		},
		StatusCode: 400,
	}
}

func putMsgToMap(fieldMap map[string]string, request *http.Request) {
	fieldMap["{host}"] = request.Host
	fieldMap["{method}"] = request.Method
	fieldMap["{uri}"] = request.URL.RequestURI()
	fieldMap["{pid}"] = strconv.Itoa(os.Getpid())
	fieldMap["{version}"] = strings.Split(request.Proto, "/")[1]
	hostname, _ := os.Hostname()
	fieldMap["{hostname}"] = hostname
	fieldMap["{req_headers}"] = TransToString(request.Header)
	fieldMap["{target}"] = request.URL.Path + request.URL.RawQuery
}
