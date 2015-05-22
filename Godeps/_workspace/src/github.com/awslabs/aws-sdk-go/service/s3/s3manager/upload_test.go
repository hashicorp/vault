package s3manager_test

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"sync"
	"testing"

	"github.com/awslabs/aws-sdk-go/aws"
	"github.com/awslabs/aws-sdk-go/aws/awserr"
	"github.com/awslabs/aws-sdk-go/aws/awsutil"
	"github.com/awslabs/aws-sdk-go/internal/test/unit"
	"github.com/awslabs/aws-sdk-go/service/s3"
	"github.com/awslabs/aws-sdk-go/service/s3/s3manager"
	"github.com/stretchr/testify/assert"
)

var _ = unit.Imported
var buf12MB = make([]byte, 1024*1024*12)
var buf2MB = make([]byte, 1024*1024*2)

func val(i interface{}, s string) interface{} {
	return awsutil.ValuesAtPath(i, s)[0]
}

func loggingSvc() (*s3.S3, *[]string, *[]interface{}) {
	var m sync.Mutex
	partNum := 0
	names := []string{}
	params := []interface{}{}
	svc := s3.New(nil)
	svc.Handlers.Unmarshal.Clear()
	svc.Handlers.UnmarshalMeta.Clear()
	svc.Handlers.UnmarshalError.Clear()
	svc.Handlers.Send.Clear()
	svc.Handlers.Send.PushBack(func(r *aws.Request) {
		m.Lock()
		defer m.Unlock()

		names = append(names, r.Operation.Name)
		params = append(params, r.Params)

		r.HTTPResponse = &http.Response{
			StatusCode: 200,
			Body:       ioutil.NopCloser(bytes.NewReader([]byte{})),
		}

		switch data := r.Data.(type) {
		case *s3.CreateMultipartUploadOutput:
			data.UploadID = aws.String("UPLOAD-ID")
		case *s3.UploadPartOutput:
			partNum++
			data.ETag = aws.String(fmt.Sprintf("ETAG%d", partNum))
		case *s3.CompleteMultipartUploadOutput:
			data.Location = aws.String("https://location")
		}
	})

	return svc, &names, &params
}

func buflen(i interface{}) int {
	r := i.(io.Reader)
	b, _ := ioutil.ReadAll(r)
	return len(b)
}

func TestUploadOrderMulti(t *testing.T) {
	s, ops, args := loggingSvc()
	resp, err := s3manager.Upload(s, &s3manager.UploadInput{
		Bucket: aws.String("Bucket"),
		Key:    aws.String("Key"),
		Body:   bytes.NewReader(buf12MB),
	}, nil)

	assert.NoError(t, err)
	assert.Equal(t, []string{"CreateMultipartUpload", "UploadPart", "UploadPart", "UploadPart", "CompleteMultipartUpload"}, *ops)
	assert.Equal(t, "https://location", resp.Location)
	assert.Equal(t, "UPLOAD-ID", resp.UploadID)

	// Validate input values

	// UploadPart
	assert.Equal(t, "UPLOAD-ID", val((*args)[1], "UploadID"))
	assert.Equal(t, "UPLOAD-ID", val((*args)[2], "UploadID"))
	assert.Equal(t, "UPLOAD-ID", val((*args)[3], "UploadID"))

	// CompleteMultipartUpload
	assert.Equal(t, "UPLOAD-ID", val((*args)[4], "UploadID"))
	assert.Equal(t, int64(1), val((*args)[4], "MultipartUpload.Parts[0].PartNumber"))
	assert.Equal(t, int64(2), val((*args)[4], "MultipartUpload.Parts[1].PartNumber"))
	assert.Equal(t, int64(3), val((*args)[4], "MultipartUpload.Parts[2].PartNumber"))
	assert.Regexp(t, `^ETAG\d+$`, val((*args)[4], "MultipartUpload.Parts[0].ETag"))
	assert.Regexp(t, `^ETAG\d+$`, val((*args)[4], "MultipartUpload.Parts[1].ETag"))
	assert.Regexp(t, `^ETAG\d+$`, val((*args)[4], "MultipartUpload.Parts[2].ETag"))
}

func TestUploadOrderMultiDifferentPartSize(t *testing.T) {
	s, ops, args := loggingSvc()
	_, err := s3manager.Upload(s, &s3manager.UploadInput{
		Bucket: aws.String("Bucket"),
		Key:    aws.String("Key"),
		Body:   bytes.NewReader(buf12MB),
	}, &s3manager.UploadOptions{PartSize: 1024 * 1024 * 7, Concurrency: 1})

	assert.NoError(t, err)
	assert.Equal(t, []string{"CreateMultipartUpload", "UploadPart", "UploadPart", "CompleteMultipartUpload"}, *ops)

	// Part lengths
	assert.Equal(t, 1024*1024*7, buflen(val((*args)[1], "Body")))
	assert.Equal(t, 1024*1024*5, buflen(val((*args)[2], "Body")))
}

func TestUploadOrderSingle(t *testing.T) {
	s, ops, _ := loggingSvc()
	resp, err := s3manager.Upload(s, &s3manager.UploadInput{
		Bucket: aws.String("Bucket"),
		Key:    aws.String("Key"),
		Body:   bytes.NewReader(buf2MB),
	}, nil)

	assert.NoError(t, err)
	assert.Equal(t, []string{"PutObject"}, *ops)
	assert.NotEqual(t, "", resp.Location)
	assert.Equal(t, "", resp.UploadID)
}

func TestUploadOrderSingleFailure(t *testing.T) {
	s, ops, _ := loggingSvc()
	s.Handlers.Send.PushBack(func(r *aws.Request) {
		r.HTTPResponse.StatusCode = 400
	})
	resp, err := s3manager.Upload(s, &s3manager.UploadInput{
		Bucket: aws.String("Bucket"),
		Key:    aws.String("Key"),
		Body:   bytes.NewReader(buf2MB),
	}, nil)

	assert.Error(t, err)
	assert.Equal(t, []string{"PutObject"}, *ops)
	assert.Nil(t, resp)
}

func TestUploadOrderZero(t *testing.T) {
	s, ops, args := loggingSvc()
	resp, err := s3manager.Upload(s, &s3manager.UploadInput{
		Bucket: aws.String("Bucket"),
		Key:    aws.String("Key"),
		Body:   bytes.NewReader(make([]byte, 0)),
	}, nil)

	assert.NoError(t, err)
	assert.Equal(t, []string{"PutObject"}, *ops)
	assert.NotEqual(t, "", resp.Location)
	assert.Equal(t, "", resp.UploadID)
	assert.Equal(t, 0, buflen(val((*args)[0], "Body")))
}

func TestUploadOrderMultiFailure(t *testing.T) {
	s, ops, _ := loggingSvc()
	s.Handlers.Send.PushBack(func(r *aws.Request) {
		switch t := r.Data.(type) {
		case *s3.UploadPartOutput:
			if *t.ETag == "ETAG2" {
				r.HTTPResponse.StatusCode = 400
			}
		}
	})
	_, err := s3manager.Upload(s, &s3manager.UploadInput{
		Bucket: aws.String("Bucket"),
		Key:    aws.String("Key"),
		Body:   bytes.NewReader(buf12MB),
	}, &s3manager.UploadOptions{Concurrency: 1})

	assert.Error(t, err)
	assert.Equal(t, []string{"CreateMultipartUpload", "UploadPart", "UploadPart", "AbortMultipartUpload"}, *ops)
}

func TestUploadOrderMultiFailureOnComplete(t *testing.T) {
	s, ops, _ := loggingSvc()
	s.Handlers.Send.PushBack(func(r *aws.Request) {
		switch r.Data.(type) {
		case *s3.CompleteMultipartUploadOutput:
			r.HTTPResponse.StatusCode = 400
		}
	})
	_, err := s3manager.Upload(s, &s3manager.UploadInput{
		Bucket: aws.String("Bucket"),
		Key:    aws.String("Key"),
		Body:   bytes.NewReader(buf12MB),
	}, nil)

	assert.Error(t, err)
	assert.Equal(t, []string{"CreateMultipartUpload", "UploadPart", "UploadPart",
		"UploadPart", "CompleteMultipartUpload", "AbortMultipartUpload"}, *ops)
}

func TestUploadOrderMultiFailureOnCreate(t *testing.T) {
	s, ops, _ := loggingSvc()
	s.Handlers.Send.PushBack(func(r *aws.Request) {
		switch r.Data.(type) {
		case *s3.CreateMultipartUploadOutput:
			r.HTTPResponse.StatusCode = 400
		}
	})
	_, err := s3manager.Upload(s, &s3manager.UploadInput{
		Bucket: aws.String("Bucket"),
		Key:    aws.String("Key"),
		Body:   bytes.NewReader(make([]byte, 1024*1024*12)),
	}, nil)

	assert.Error(t, err)
	assert.Equal(t, []string{"CreateMultipartUpload"}, *ops)
}

func TestUploadOrderMultiFailureLeaveParts(t *testing.T) {
	s, ops, _ := loggingSvc()
	s.Handlers.Send.PushBack(func(r *aws.Request) {
		switch data := r.Data.(type) {
		case *s3.UploadPartOutput:
			if *data.ETag == "ETAG2" {
				r.HTTPResponse.StatusCode = 400
			}
		}
	})
	_, err := s3manager.Upload(s, &s3manager.UploadInput{
		Bucket: aws.String("Bucket"),
		Key:    aws.String("Key"),
		Body:   bytes.NewReader(make([]byte, 1024*1024*12)),
	}, &s3manager.UploadOptions{Concurrency: 1, LeavePartsOnError: true})

	assert.Error(t, err)
	assert.Equal(t, []string{"CreateMultipartUpload", "UploadPart", "UploadPart"}, *ops)
}

var failreaderCount = 0

type failreader struct{ times int }

func (f failreader) Read(b []byte) (int, error) {
	failreaderCount++
	if failreaderCount >= f.times {
		return 0, fmt.Errorf("random failure")
	}
	return len(b), nil
}

func TestUploadOrderReadFail1(t *testing.T) {
	failreaderCount = 0
	s, ops, _ := loggingSvc()
	_, err := s3manager.Upload(s, &s3manager.UploadInput{
		Bucket: aws.String("Bucket"),
		Key:    aws.String("Key"),
		Body:   failreader{1},
	}, nil)

	assert.Equal(t, "ReadRequestBody", err.(awserr.Error).Code())
	assert.EqualError(t, err.(awserr.Error).OrigErr(), "random failure")
	assert.Equal(t, []string{}, *ops)
}

func TestUploadOrderReadFail2(t *testing.T) {
	failreaderCount = 0
	s, ops, _ := loggingSvc()
	_, err := s3manager.Upload(s, &s3manager.UploadInput{
		Bucket: aws.String("Bucket"),
		Key:    aws.String("Key"),
		Body:   failreader{2},
	}, nil)

	assert.Equal(t, "MultipartUpload", err.(awserr.Error).Code())
	assert.Equal(t, "ReadRequestBody", err.(awserr.Error).OrigErr().(awserr.Error).Code())
	assert.EqualError(t, err.(awserr.Error).OrigErr().(awserr.Error).OrigErr(), "random failure")
	assert.Equal(t, []string{"CreateMultipartUpload", "AbortMultipartUpload"}, *ops)
}
