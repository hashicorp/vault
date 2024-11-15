package cfclient

//go:generate go run gen_error.go

import (
	"fmt"
)

type CloudFoundryError struct {
	Code        int    `json:"code"`
	ErrorCode   string `json:"error_code"`
	Description string `json:"description"`
}

type CloudFoundryErrorsV3 struct {
	Errors []CloudFoundryErrorV3 `json:"errors"`
}

type CloudFoundryErrorV3 struct {
	Code   int    `json:"code"`
	Title  string `json:"title"`
	Detail string `json:"detail"`
}

// CF APIs v3 can return multiple errors, we take the first one and convert it into a V2 model
func NewCloudFoundryErrorFromV3Errors(cfErrorsV3 CloudFoundryErrorsV3) CloudFoundryError {
	if len(cfErrorsV3.Errors) == 0 {
		return CloudFoundryError{
			0,
			"GO-Client-No-Errors",
			"No Errors in response from V3",
		}
	}

	return CloudFoundryError{
		cfErrorsV3.Errors[0].Code,
		cfErrorsV3.Errors[0].Title,
		cfErrorsV3.Errors[0].Detail,
	}
}

func (cfErr CloudFoundryError) Error() string {
	return fmt.Sprintf("cfclient error (%s|%d): %s", cfErr.ErrorCode, cfErr.Code, cfErr.Description)
}

type CloudFoundryHTTPError struct {
	StatusCode int
	Status     string
	Body       []byte
}

func (e CloudFoundryHTTPError) Error() string {
	return fmt.Sprintf("cfclient: HTTP error (%d): %s", e.StatusCode, e.Status)
}
