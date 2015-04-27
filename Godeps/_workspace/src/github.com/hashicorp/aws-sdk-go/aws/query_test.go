package aws_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"reflect"
	"sync"
	"testing"
	"time"

	"github.com/hashicorp/aws-sdk-go/aws"
)

func TestQueryRequest(t *testing.T) {
	var m sync.Mutex
	var httpReq *http.Request
	var form url.Values

	server := httptest.NewServer(http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			m.Lock()
			defer m.Unlock()

			httpReq = r

			if err := r.ParseForm(); err != nil {
				t.Fatal(err)
			}
			form = r.Form

			fmt.Fprintln(w, `<Thing><IpAddress>woo</IpAddress></Thing>`)
		},
	))
	defer server.Close()

	client := aws.QueryClient{
		Context: aws.Context{
			Service: "animals",
			Region:  "us-west-2",
			Credentials: aws.Creds(
				"accessKeyID",
				"secretAccessKey",
				"securityToken",
			),
		},
		Client:     http.DefaultClient,
		Endpoint:   server.URL,
		APIVersion: "1.1",
	}

	req := fakeQueryRequest{
		PresentString:  aws.String("string"),
		PresentBoolean: aws.True(),
		PresentInteger: aws.Integer(1),
		PresentLong:    aws.Long(2),
		PresentDouble:  aws.Double(1.2),
		PresentFloat:   aws.Float(2.3),
		PresentTime:    time.Date(2001, 1, 1, 2, 1, 1, 0, time.FixedZone("UTC+1", 3600)),
		PresentSlice:   []string{"one", "two"},
		PresentStruct:  &EmbeddedStruct{Value: aws.String("v")},
		PresentStructSlice: []EmbeddedStruct{
			{Value: aws.String("p")},
			{Value: aws.String("q")},
		},
		PresentMap: map[string]EmbeddedStruct{
			"aa": EmbeddedStruct{Value: aws.String("AA")},
			"bb": EmbeddedStruct{Value: aws.String("BB")},
		},
	}
	var resp fakeQueryResponse
	if err := client.Do("GetIP", "POST", "/", &req, &resp); err != nil {
		t.Fatal(err)
	}

	m.Lock()
	defer m.Unlock()

	if v, want := httpReq.Method, "POST"; v != want {
		t.Errorf("Method was %v but expected %v", v, want)
	}

	if httpReq.Header.Get("Authorization") == "" {
		t.Error("Authorization header is missing")
	}

	if v, want := httpReq.Header.Get("Content-Type"), "application/x-www-form-urlencoded"; v != want {
		t.Errorf("Content-Type was %v but expected %v", v, want)
	}

	if v, want := httpReq.Header.Get("User-Agent"), "aws-go"; v != want {
		t.Errorf("User-Agent was %v but expected %v", v, want)
	}

	if err := httpReq.ParseForm(); err != nil {
		t.Fatal(err)
	}

	expectedForm := url.Values{
		"Action":                     []string{"GetIP"},
		"Version":                    []string{"1.1"},
		"PresentString":              []string{"string"},
		"PresentBoolean":             []string{"true"},
		"PresentInteger":             []string{"1"},
		"PresentLong":                []string{"2"},
		"PresentDouble":              []string{"1.2"},
		"PresentFloat":               []string{"2.3"},
		"PresentTime":                []string{"2001-01-01T01:01:01Z"},
		"PresentSlice.1":             []string{"one"},
		"PresentSlice.2":             []string{"two"},
		"PresentStruct.Value":        []string{"v"},
		"PresentStructSlice.1.Value": []string{"p"},
		"PresentStructSlice.2.Value": []string{"q"},
		"PresentMap.1.Name":          []string{"aa"},
		"PresentMap.1.Value.Value":   []string{"AA"},
		"PresentMap.2.Name":          []string{"bb"},
		"PresentMap.2.Value.Value":   []string{"BB"},
	}

	if !reflect.DeepEqual(form, expectedForm) {
		t.Errorf("Post body was \n%s\n but expected \n%s", form.Encode(), expectedForm.Encode())
	}

	if want := (fakeQueryResponse{IPAddress: "woo"}); want != resp {
		t.Errorf("Response was %#v, but expected %#v", resp, want)
	}
}

func TestQueryRequestError(t *testing.T) {
	var m sync.Mutex
	var httpReq *http.Request
	var form url.Values

	server := httptest.NewServer(http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			m.Lock()
			defer m.Unlock()

			httpReq = r

			if err := r.ParseForm(); err != nil {
				t.Fatal(err)
			}
			form = r.Form

			w.WriteHeader(400)
			fmt.Fprintln(w, `<ErrorResponse>
<RequestId>woo</RequestId>
<Error>
<Type>Problem</Type>
<Code>Uh Oh</Code>
<Message>You done did it</Message>
</Error>
</ErrorResponse>`)
		},
	))
	defer server.Close()

	client := aws.QueryClient{
		Context: aws.Context{
			Service: "animals",
			Region:  "us-west-2",
			Credentials: aws.Creds(
				"accessKeyID",
				"secretAccessKey",
				"securityToken",
			),
		},
		Client:     http.DefaultClient,
		Endpoint:   server.URL,
		APIVersion: "1.1",
	}

	req := fakeQueryRequest{}
	var resp fakeQueryResponse
	err := client.Do("GetIP", "POST", "/", &req, &resp)
	if err == nil {
		t.Fatal("Expected an error but none was returned")
	}

	if err, ok := err.(aws.APIError); ok {
		if v, want := err.Type, "Problem"; v != want {
			t.Errorf("Error type was %v, but expected %v", v, want)
		}

		if v, want := err.Code, "Uh Oh"; v != want {
			t.Errorf("Error type was %v, but expected %v", v, want)
		}

		if v, want := err.Message, "You done did it"; v != want {
			t.Errorf("Error message was %v, but expected %v", v, want)
		}
	} else {
		t.Errorf("Unknown error returned: %#v", err)
	}
}

type fakeQueryRequest struct {
	PresentString aws.StringValue `query:"PresentString"`
	MissingString aws.StringValue `query:"MissingString"`

	PresentInteger aws.IntegerValue `query:"PresentInteger"`
	MissingInteger aws.IntegerValue `query:"MissingInteger"`

	PresentLong aws.LongValue `query:"PresentLong"`
	MissingLong aws.LongValue `query:"MissingLong"`

	PresentDouble aws.DoubleValue `query:"PresentDouble"`
	MissingDouble aws.DoubleValue `query:"MissingDouble"`

	PresentFloat aws.FloatValue `query:"PresentFloat"`
	MissingFloat aws.FloatValue `query:"MissingFloat"`

	PresentBoolean aws.BooleanValue `query:"PresentBoolean"`
	MissingBoolean aws.BooleanValue `query:"MissingBoolean"`

	PresentTime time.Time `query:"PresentTime"`
	MissingTime time.Time `query:"MissingTime"`

	PresentSlice []string `query:"PresentSlice"`
	MissingSlice []string `query:"MissingSlice"`

	PresentStructSlice []EmbeddedStruct `query:"PresentStructSlice"`
	MissingStructSlice []EmbeddedStruct `query:"MissingStructSlice"`

	PresentMap map[string]EmbeddedStruct `query:"PresentMap"`
	MissingMap map[string]EmbeddedStruct `query:"MissingMap"`

	PresentStruct *EmbeddedStruct `query:"PresentStruct"`
	MissingStruct *EmbeddedStruct `query:"MissingStruct"`
}

type EmbeddedStruct struct {
	Value aws.StringValue
}

type fakeQueryResponse struct {
	IPAddress string `xml:"IpAddress"`
}
