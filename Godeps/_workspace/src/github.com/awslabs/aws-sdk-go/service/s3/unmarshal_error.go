package s3

import (
	"encoding/xml"
	"io"
	"strings"

	"github.com/awslabs/aws-sdk-go/aws"
)

type xmlErrorResponse struct {
	XMLName xml.Name `xml:"Error"`
	Code    string   `xml:"Code"`
	Message string   `xml:"Message"`
}

func unmarshalError(r *aws.Request) {
	defer r.HTTPResponse.Body.Close()

	if r.HTTPResponse.ContentLength == int64(0) {
		// No body, use status code to generate an APIError
		r.Error = aws.APIError{
			StatusCode: r.HTTPResponse.StatusCode,
			Code:       strings.Replace(r.HTTPResponse.Status, " ", "", -1),
			Message:    r.HTTPResponse.Status,
		}
		return
	}

	resp := &xmlErrorResponse{}
	err := xml.NewDecoder(r.HTTPResponse.Body).Decode(resp)
	if err != nil && err != io.EOF {
		r.Error = err
	} else {
		r.Error = aws.APIError{
			StatusCode: r.HTTPResponse.StatusCode,
			Code:       resp.Code,
			Message:    resp.Message,
		}
	}
}
