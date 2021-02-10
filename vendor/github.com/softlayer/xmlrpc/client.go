package xmlrpc

import (
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/cookiejar"
	"net/rpc"
	"net/url"
	"strconv"
	"sync"
	"time"
)

type Client struct {
	*rpc.Client
}

// clientCodec is rpc.ClientCodec interface implementation.
type clientCodec struct {
	// url presents url of xmlrpc service
	url *url.URL

	// httpClient works with HTTP protocol
	httpClient *http.Client

	// cookies stores cookies received on last request
	cookies http.CookieJar

	// responses presents map of active requests. It is required to return request id, that
	// rpc.Client can mark them as done.

	responsesMu sync.RWMutex
	responses   map[uint64]*http.Response
	response Response

	// ready presents channel, that is used to link request and it`s response.
	ready chan uint64
	// close notifies codec is closed.
	close chan uint64

}

func (codec *clientCodec) WriteRequest(request *rpc.Request, args interface{}) (err error) {
	httpRequest, err := NewRequest(codec.url.String(), request.ServiceMethod, args)
	if err != nil {
		return err
	}

	if codec.cookies != nil {
		for _, cookie := range codec.cookies.Cookies(codec.url) {
			httpRequest.AddCookie(cookie)
		}
	}


	httpResponse, err := codec.httpClient.Do(httpRequest)

	if err != nil {
		return err
	}

	if codec.cookies != nil {
		codec.cookies.SetCookies(codec.url, httpResponse.Cookies())
	}

	codec.responsesMu.Lock()
	codec.responses[request.Seq] = httpResponse
	codec.responsesMu.Unlock()
	codec.ready <- request.Seq

	return nil
}

func (codec *clientCodec) ReadResponseHeader(response *rpc.Response) (err error) {
	var seq uint64
	select {
		case seq = <-codec.ready:
		case <-codec.close:
			return errors.New("codec is closed")
	}
	response.Seq = seq

	codec.responsesMu.RLock()
	httpResponse := codec.responses[seq]
	delete(codec.responses, seq)
	codec.responsesMu.RUnlock()

	defer httpResponse.Body.Close()

	contentLength := httpResponse.ContentLength
	if contentLength == -1 {
		if ntcoentLengthHeader, ok := httpResponse.Header["Ntcoent-Length"]; ok {
			ntcoentLength, err := strconv.ParseInt(ntcoentLengthHeader[0], 10, 64)
			if err == nil {
				contentLength = ntcoentLength
			}
		}
	}

	var respData []byte
	if contentLength != -1 {
		respData = make([]byte, contentLength)
		_, err = io.ReadFull(httpResponse.Body, respData)
	} else {
		respData, err = ioutil.ReadAll(httpResponse.Body)
	}

	if err != nil {
		response.Error = err.Error()
		return nil
	}


	resp := NewResponse(respData, httpResponse.StatusCode)

	if resp.Failed() {
		err := resp.Err()
		response.Error = fmt.Sprintf("%v", err)
		return err

	}
	codec.response = *resp

	if httpResponse.StatusCode < 200 || httpResponse.StatusCode >= 300 {
		return &XmlRpcError{HttpStatusCode: httpResponse.StatusCode}
	}

	return nil
}

func (codec *clientCodec) ReadResponseBody(v interface{}) (err error) {
	if v == nil {
		return nil
	}
	return codec.response.Unmarshal(v)
}

func (codec *clientCodec) Close() error {
	if transport, ok := codec.httpClient.Transport.(*http.Transport); ok {
		transport.CloseIdleConnections()
	}

	close(codec.close)
	return nil
}

// NewClient returns instance of rpc.Client object, that is used to send request to xmlrpc service.
func NewClient(requrl string, transport http.RoundTripper, timeout time.Duration) (*Client, error) {
	if transport == nil {
		transport = http.DefaultTransport
	}

	httpClient := &http.Client{
		Transport: transport,
		Timeout:   timeout,
	}

	jar, err := cookiejar.New(nil)

	if err != nil {
		return nil, err
	}

	u, err := url.Parse(requrl)

	if err != nil {
		return nil, err
	}

	codec := clientCodec{
		url:        u,
		httpClient: httpClient,
		ready:      make(chan uint64),
		close:      make(chan uint64),
		responses:  make(map[uint64]*http.Response),
		cookies:    jar,
	}

	return &Client{rpc.NewClientWithCodec(&codec)}, nil
}
