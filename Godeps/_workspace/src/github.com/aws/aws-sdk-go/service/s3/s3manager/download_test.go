package s3manager_test

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"regexp"
	"strconv"
	"sync"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/internal/test/unit"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/stretchr/testify/assert"
)

var _ = unit.Imported

func dlLoggingSvc(data []byte) (*s3.S3, *[]string, *[]string) {
	var m sync.Mutex
	names := []string{}
	ranges := []string{}

	svc := s3.New(nil)
	svc.Handlers.Send.Clear()
	svc.Handlers.Send.PushBack(func(r *aws.Request) {
		m.Lock()
		defer m.Unlock()

		names = append(names, r.Operation.Name)
		ranges = append(ranges, *r.Params.(*s3.GetObjectInput).Range)

		rerng := regexp.MustCompile(`bytes=(\d+)-(\d+)`)
		rng := rerng.FindStringSubmatch(r.HTTPRequest.Header.Get("Range"))
		start, _ := strconv.ParseInt(rng[1], 10, 64)
		fin, _ := strconv.ParseInt(rng[2], 10, 64)
		fin++

		if fin > int64(len(data)) {
			fin = int64(len(data))
		}

		r.HTTPResponse = &http.Response{
			StatusCode: 200,
			Body:       ioutil.NopCloser(bytes.NewReader(data[start:fin])),
			Header:     http.Header{},
		}
		r.HTTPResponse.Header.Set("Content-Range", fmt.Sprintf("bytes %d-%d/%d",
			start, fin, len(data)))
	})

	return svc, &names, &ranges
}

type dlwriter struct {
	buf []byte
}

func newDLWriter(size int) *dlwriter {
	return &dlwriter{buf: make([]byte, size)}
}

func (d dlwriter) WriteAt(p []byte, pos int64) (n int, err error) {
	if pos > int64(len(d.buf)) {
		return 0, io.EOF
	}

	written := 0
	for i, b := range p {
		if i >= len(d.buf) {
			break
		}
		d.buf[pos+int64(i)] = b
		written++
	}
	return written, nil
}

func TestDownloadOrder(t *testing.T) {
	s, names, ranges := dlLoggingSvc(buf12MB)

	opts := &s3manager.DownloadOptions{S3: s, Concurrency: 1}
	d := s3manager.NewDownloader(opts)
	w := newDLWriter(len(buf12MB))
	n, err := d.Download(w, &s3.GetObjectInput{
		Bucket: aws.String("bucket"),
		Key:    aws.String("key"),
	})

	assert.Nil(t, err)
	assert.Equal(t, int64(len(buf12MB)), n)
	assert.Equal(t, []string{"GetObject", "GetObject", "GetObject"}, *names)
	assert.Equal(t, []string{"bytes=0-5242879", "bytes=5242880-10485759", "bytes=10485760-15728639"}, *ranges)

	count := 0
	for _, b := range w.buf {
		count += int(b)
	}
	assert.Equal(t, 0, count)
}

func TestDownloadZero(t *testing.T) {
	s, names, ranges := dlLoggingSvc([]byte{})

	opts := &s3manager.DownloadOptions{S3: s}
	d := s3manager.NewDownloader(opts)
	w := newDLWriter(0)
	n, err := d.Download(w, &s3.GetObjectInput{
		Bucket: aws.String("bucket"),
		Key:    aws.String("key"),
	})

	assert.Nil(t, err)
	assert.Equal(t, int64(0), n)
	assert.Equal(t, []string{"GetObject"}, *names)
	assert.Equal(t, []string{"bytes=0-5242879"}, *ranges)
}

func TestDownloadSetPartSize(t *testing.T) {
	s, names, ranges := dlLoggingSvc([]byte{1, 2, 3})

	opts := &s3manager.DownloadOptions{S3: s, PartSize: 1, Concurrency: 1}
	d := s3manager.NewDownloader(opts)
	w := newDLWriter(3)
	n, err := d.Download(w, &s3.GetObjectInput{
		Bucket: aws.String("bucket"),
		Key:    aws.String("key"),
	})

	assert.Nil(t, err)
	assert.Equal(t, int64(3), n)
	assert.Equal(t, []string{"GetObject", "GetObject", "GetObject"}, *names)
	assert.Equal(t, []string{"bytes=0-0", "bytes=1-1", "bytes=2-2"}, *ranges)
	assert.Equal(t, []byte{1, 2, 3}, w.buf)
}

func TestDownloadError(t *testing.T) {
	s, names, _ := dlLoggingSvc([]byte{1, 2, 3})
	opts := &s3manager.DownloadOptions{S3: s, PartSize: 1, Concurrency: 1}

	num := 0
	s.Handlers.Send.PushBack(func(r *aws.Request) {
		num++
		if num > 1 {
			r.HTTPResponse.StatusCode = 400
			r.HTTPResponse.Body = ioutil.NopCloser(bytes.NewReader([]byte{}))
		}
	})

	d := s3manager.NewDownloader(opts)
	w := newDLWriter(3)
	n, err := d.Download(w, &s3.GetObjectInput{
		Bucket: aws.String("bucket"),
		Key:    aws.String("key"),
	})

	assert.NotNil(t, err)
	assert.Equal(t, int64(1), n)
	assert.Equal(t, []string{"GetObject", "GetObject"}, *names)
	assert.Equal(t, []byte{1, 0, 0}, w.buf)
}
