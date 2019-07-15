// Copyright (c) 2016, 2018, Oracle and/or its affiliates. All rights reserved.
//
// Helper methods for Oracle Cloud Infrastructure Go SDK Samples
//

package helpers

import (
	"crypto/sha256"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"math"
	"math/rand"
	"reflect"
	"strings"
	"time"

	"github.com/oracle/oci-go-sdk/common"
)

// FatalIfError is equivalent to Println() followed by a call to os.Exit(1) if error is not nil
func FatalIfError(err error) {
	if err != nil {
		log.Fatalln(err.Error())
	}
}

// RetryUntilTrueOrError retries a function until the predicate is true or it reaches a timeout.
// The operation is retried at the give frequency
func RetryUntilTrueOrError(operation func() (interface{}, error), predicate func(interface{}) (bool, error), frequency, timeout <-chan time.Time) error {
	for {
		select {
		case <-timeout:
			return fmt.Errorf("timeout reached")
		case <-frequency:
			result, err := operation()
			if err != nil {
				return err
			}

			isTrue, err := predicate(result)
			if err != nil {
				return err
			}

			if isTrue {
				return nil
			}
		}
	}
}

// FindLifecycleFieldValue finds lifecycle value inside the struct based on reflection
func FindLifecycleFieldValue(request interface{}) (string, error) {
	val := reflect.ValueOf(request)
	if val.Kind() == reflect.Ptr {
		if val.IsNil() {
			return "", fmt.Errorf("can not unmarshal to response a pointer to nil structure")
		}
		val = val.Elem()
	}

	var err error
	typ := val.Type()
	for i := 0; i < typ.NumField(); i++ {
		if err != nil {
			return "", err
		}

		sf := typ.Field(i)

		//unexported
		if sf.PkgPath != "" {
			continue
		}

		sv := val.Field(i)

		if sv.Kind() == reflect.Struct {
			lif, err := FindLifecycleFieldValue(sv.Interface())
			if err == nil {
				return lif, nil
			}
		}
		if !strings.Contains(strings.ToLower(sf.Name), "lifecyclestate") {
			continue
		}
		return sv.String(), nil
	}
	return "", fmt.Errorf("request does not have a lifecycle field")
}

// CheckLifecycleState returns a function that checks for that a struct has the given lifecycle
func CheckLifecycleState(lifecycleState string) func(interface{}) (bool, error) {
	return func(request interface{}) (bool, error) {
		fieldLifecycle, err := FindLifecycleFieldValue(request)
		if err != nil {
			return false, err
		}
		isEqual := fieldLifecycle == lifecycleState
		log.Printf("Current lifecycle state is: %s, waiting for it becomes to: %s", fieldLifecycle, lifecycleState)
		return isEqual, nil
	}
}

// GetRequestMetadataWithDefaultRetryPolicy returns a requestMetadata with default retry policy
// which will do retry for non-200 status code return back from service
// Notes: not all non-200 status code should do retry, this should be based on specific operation
// such as delete operation followed with get operation will retrun 404 if resource already been
// deleted
func GetRequestMetadataWithDefaultRetryPolicy() common.RequestMetadata {
	return common.RequestMetadata{
		RetryPolicy: getDefaultRetryPolicy(),
	}
}

// GetRequestMetadataWithCustomizedRetryPolicy returns a requestMetadata which will do the retry based on
// input function (retry until the function return false)
func GetRequestMetadataWithCustomizedRetryPolicy(fn func(r common.OCIOperationResponse) bool) common.RequestMetadata {
	return common.RequestMetadata{
		RetryPolicy: getExponentialBackoffRetryPolicy(uint(20), fn),
	}
}

func getDefaultRetryPolicy() *common.RetryPolicy {
	// how many times to do the retry
	attempts := uint(10)

	// retry for all non-200 status code
	retryOnAllNon200ResponseCodes := func(r common.OCIOperationResponse) bool {
		return !(r.Error == nil && 199 < r.Response.HTTPResponse().StatusCode && r.Response.HTTPResponse().StatusCode < 300)
	}
	return getExponentialBackoffRetryPolicy(attempts, retryOnAllNon200ResponseCodes)
}

func getExponentialBackoffRetryPolicy(n uint, fn func(r common.OCIOperationResponse) bool) *common.RetryPolicy {
	// the duration between each retry operation, you might want to waite longer each time the retry fails
	exponentialBackoff := func(r common.OCIOperationResponse) time.Duration {
		return time.Duration(math.Pow(float64(2), float64(r.AttemptNumber-1))) * time.Second
	}
	policy := common.NewRetryPolicy(n, fn, exponentialBackoff)
	return &policy
}

// GetRandomString returns a random string with length equals to n
func GetRandomString(n int) string {
	letters := []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

	b := make([]rune, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

// WriteTempFileOfSize output random content to a file
func WriteTempFileOfSize(filesize int64) (fileName string, fileSize int64) {
	hash := sha256.New()
	f, _ := ioutil.TempFile("", "OCIGOSDKSampleFile")
	ra := rand.New(rand.NewSource(time.Now().UnixNano()))
	defer f.Close()
	writer := io.MultiWriter(f, hash)
	written, _ := io.CopyN(writer, ra, filesize)
	fileName = f.Name()
	fileSize = written
	return
}
