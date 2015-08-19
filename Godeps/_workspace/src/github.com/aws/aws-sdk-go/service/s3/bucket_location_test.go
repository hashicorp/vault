package s3_test

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awsutil"
	"github.com/aws/aws-sdk-go/aws/service"
	"github.com/aws/aws-sdk-go/internal/test/unit"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/stretchr/testify/assert"
)

var _ = unit.Imported
var s3LocationTests = []struct {
	body string
	loc  string
}{
	{`<?xml version="1.0" encoding="UTF-8"?><LocationConstraint xmlns="http://s3.amazonaws.com/doc/2006-03-01/"/>`, ``},
	{`<?xml version="1.0" encoding="UTF-8"?><LocationConstraint xmlns="http://s3.amazonaws.com/doc/2006-03-01/">EU</LocationConstraint>`, `EU`},
}

func TestGetBucketLocation(t *testing.T) {
	for _, test := range s3LocationTests {
		s := s3.New(nil)
		s.Handlers.Send.Clear()
		s.Handlers.Send.PushBack(func(r *service.Request) {
			reader := ioutil.NopCloser(bytes.NewReader([]byte(test.body)))
			r.HTTPResponse = &http.Response{StatusCode: 200, Body: reader}
		})

		resp, err := s.GetBucketLocation(&s3.GetBucketLocationInput{Bucket: aws.String("bucket")})
		assert.NoError(t, err)
		if test.loc == "" {
			assert.Nil(t, resp.LocationConstraint)
		} else {
			assert.Equal(t, test.loc, *resp.LocationConstraint)
		}
	}
}

func TestPopulateLocationConstraint(t *testing.T) {
	s := s3.New(nil)
	in := &s3.CreateBucketInput{
		Bucket: aws.String("bucket"),
	}
	req, _ := s.CreateBucketRequest(in)
	err := req.Build()
	assert.NoError(t, err)
	assert.Equal(t, "mock-region", awsutil.ValuesAtPath(req.Params, "CreateBucketConfiguration.LocationConstraint")[0])
	assert.Nil(t, in.CreateBucketConfiguration) // don't modify original params
}

func TestNoPopulateLocationConstraintIfProvided(t *testing.T) {
	s := s3.New(nil)
	req, _ := s.CreateBucketRequest(&s3.CreateBucketInput{
		Bucket: aws.String("bucket"),
		CreateBucketConfiguration: &s3.CreateBucketConfiguration{},
	})
	err := req.Build()
	assert.NoError(t, err)
	assert.Equal(t, 0, len(awsutil.ValuesAtPath(req.Params, "CreateBucketConfiguration.LocationConstraint")))
}

func TestNoPopulateLocationConstraintIfClassic(t *testing.T) {
	s := s3.New(&aws.Config{Region: aws.String("us-east-1")})
	req, _ := s.CreateBucketRequest(&s3.CreateBucketInput{
		Bucket: aws.String("bucket"),
	})
	err := req.Build()
	assert.NoError(t, err)
	assert.Equal(t, 0, len(awsutil.ValuesAtPath(req.Params, "CreateBucketConfiguration.LocationConstraint")))
}
