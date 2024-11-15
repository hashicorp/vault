package oss

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
)

// ServiceError contains fields of the error response from Oss Service REST API.
type ServiceError struct {
	XMLName    xml.Name `xml:"Error"`
	Code       string   `xml:"Code"`      // The error code returned from OSS to the caller
	Message    string   `xml:"Message"`   // The detail error message from OSS
	RequestID  string   `xml:"RequestId"` // The UUID used to uniquely identify the request
	HostID     string   `xml:"HostId"`    // The OSS server cluster's Id
	Endpoint   string   `xml:"Endpoint"`
	Ec         string   `xml:"EC"`
	RawMessage string   // The raw messages from OSS
	StatusCode int      // HTTP status code

}

// Error implements interface error
func (e ServiceError) Error() string {
	errorStr := fmt.Sprintf("oss: service returned error: StatusCode=%d, ErrorCode=%s, ErrorMessage=\"%s\", RequestId=%s", e.StatusCode, e.Code, e.Message, e.RequestID)
	if len(e.Endpoint) > 0 {
		errorStr = fmt.Sprintf("%s, Endpoint=%s", errorStr, e.Endpoint)
	}
	if len(e.Ec) > 0 {
		errorStr = fmt.Sprintf("%s, Ec=%s", errorStr, e.Ec)
	}
	return errorStr
}

// UnexpectedStatusCodeError is returned when a storage service responds with neither an error
// nor with an HTTP status code indicating success.
type UnexpectedStatusCodeError struct {
	allowed []int // The expected HTTP stats code returned from OSS
	got     int   // The actual HTTP status code from OSS
}

// Error implements interface error
func (e UnexpectedStatusCodeError) Error() string {
	s := func(i int) string { return fmt.Sprintf("%d %s", i, http.StatusText(i)) }

	got := s(e.got)
	expected := []string{}
	for _, v := range e.allowed {
		expected = append(expected, s(v))
	}
	return fmt.Sprintf("oss: status code from service response is %s; was expecting %s",
		got, strings.Join(expected, " or "))
}

// Got is the actual status code returned by oss.
func (e UnexpectedStatusCodeError) Got() int {
	return e.got
}

// CheckRespCode returns UnexpectedStatusError if the given response code is not
// one of the allowed status codes; otherwise nil.
func CheckRespCode(respCode int, allowed []int) error {
	for _, v := range allowed {
		if respCode == v {
			return nil
		}
	}
	return UnexpectedStatusCodeError{allowed, respCode}
}

// CheckCallbackResp return error if the given response code is not 200
func CheckCallbackResp(resp *Response) error {
	var err error
	contentLengthStr := resp.Headers.Get("Content-Length")
	contentLength, _ := strconv.Atoi(contentLengthStr)
	var bodyBytes []byte
	if contentLength > 0 {
		bodyBytes, _ = ioutil.ReadAll(resp.Body)
	}
	if len(bodyBytes) > 0 {
		srvErr, errIn := serviceErrFromXML(bodyBytes, resp.StatusCode,
			resp.Headers.Get(HTTPHeaderOssRequestID))
		if errIn != nil {
			if len(resp.Headers.Get(HTTPHeaderOssEc)) > 0 {
				err = fmt.Errorf("unknown response body, status code = %d, RequestId = %s, ec = %s", resp.StatusCode, resp.Headers.Get(HTTPHeaderOssRequestID), resp.Headers.Get(HTTPHeaderOssEc))
			} else {
				err = fmt.Errorf("unknown response body, status code= %d, RequestId = %s", resp.StatusCode, resp.Headers.Get(HTTPHeaderOssRequestID))
			}
		} else {
			err = srvErr
		}
	}
	return err
}

func tryConvertServiceError(data []byte, resp *Response, def error) (err error) {
	err = def
	if len(data) > 0 {
		srvErr, errIn := serviceErrFromXML(data, resp.StatusCode, resp.Headers.Get(HTTPHeaderOssRequestID))
		if errIn == nil {
			err = srvErr
		}
	}
	return err
}

// CRCCheckError is returned when crc check is inconsistent between client and server
type CRCCheckError struct {
	clientCRC uint64 // Calculated CRC64 in client
	serverCRC uint64 // Calculated CRC64 in server
	operation string // Upload operations such as PutObject/AppendObject/UploadPart, etc
	requestID string // The request id of this operation
}

// Error implements interface error
func (e CRCCheckError) Error() string {
	return fmt.Sprintf("oss: the crc of %s is inconsistent, client %d but server %d; request id is %s",
		e.operation, e.clientCRC, e.serverCRC, e.requestID)
}

func CheckDownloadCRC(clientCRC, serverCRC uint64) error {
	if clientCRC == serverCRC {
		return nil
	}
	return CRCCheckError{clientCRC, serverCRC, "DownloadFile", ""}
}

func CheckCRC(resp *Response, operation string) error {
	if resp.Headers.Get(HTTPHeaderOssCRC64) == "" || resp.ClientCRC == resp.ServerCRC {
		return nil
	}
	return CRCCheckError{resp.ClientCRC, resp.ServerCRC, operation, resp.Headers.Get(HTTPHeaderOssRequestID)}
}
