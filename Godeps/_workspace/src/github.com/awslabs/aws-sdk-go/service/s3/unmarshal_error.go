package s3

import (
	"encoding/xml"
	"io"
	"strings"

	"github.com/awslabs/aws-sdk-go/aws"
	"github.com/awslabs/aws-sdk-go/internal/apierr"
)

type xmlErrorResponse struct {
	XMLName xml.Name `xml:"Error"`
	Code    string   `xml:"Code"`
	Message string   `xml:"Message"`
}

func unmarshalError(r *aws.Request) {
	defer r.HTTPResponse.Body.Close()

	if r.HTTPResponse.ContentLength == int64(0) {
		// No body, use status code to generate an awserr.Error
		r.Error = apierr.NewRequestError(
			apierr.New(strings.Replace(r.HTTPResponse.Status, " ", "", -1), r.HTTPResponse.Status, nil),
			r.HTTPResponse.StatusCode,
			"",
		)
		return
	}

	resp := &xmlErrorResponse{}
	err := xml.NewDecoder(r.HTTPResponse.Body).Decode(resp)
	if err != nil && err != io.EOF {
		r.Error = apierr.New("Unmarshal", "failed to decode S3 XML error response", nil)
	} else {
		r.Error = apierr.NewRequestError(
			apierr.New(resp.Code, resp.Message, nil),
			r.HTTPResponse.StatusCode,
			"",
		)
	}
}
