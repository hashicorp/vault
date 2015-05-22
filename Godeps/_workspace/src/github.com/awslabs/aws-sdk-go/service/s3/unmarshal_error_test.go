package s3_test

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/awslabs/aws-sdk-go/aws"
	"github.com/awslabs/aws-sdk-go/internal/test/unit"
	"github.com/awslabs/aws-sdk-go/service/s3"
	"github.com/stretchr/testify/assert"
	"github.com/awslabs/aws-sdk-go/aws/awserr"
)

var _ = unit.Imported

var s3StatusCodeErrorTests = []struct {
	scode   int
	status  string
	body    string
	code    string
	message string
}{
	{301, "Moved Permanently", "", "MovedPermanently", "Moved Permanently"},
	{403, "Forbidden", "", "Forbidden", "Forbidden"},
	{400, "Bad Request", "", "BadRequest", "Bad Request"},
	{404, "Not Found", "", "NotFound", "Not Found"},
	{500, "Internal Error", "", "InternalError", "Internal Error"},
}

func TestStatusCodeError(t *testing.T) {
	for _, test := range s3StatusCodeErrorTests {
		s := s3.New(nil)
		s.Handlers.Send.Clear()
		s.Handlers.Send.PushBack(func(r *aws.Request) {
			body := ioutil.NopCloser(bytes.NewReader([]byte(test.body)))
			r.HTTPResponse = &http.Response{
				ContentLength: int64(len(test.body)),
				StatusCode:    test.scode,
				Status:        test.status,
				Body:          body,
			}
		})
		_, err := s.PutBucketACL(&s3.PutBucketACLInput{
			Bucket: aws.String("bucket"), ACL: aws.String("public-read"),
		})

		assert.Error(t, err)
		assert.Equal(t, test.code, err.(awserr.Error).Code())
		assert.Equal(t, test.message, err.(awserr.Error).Message())
	}
}
