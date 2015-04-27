package aws_test

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"

	"github.com/hashicorp/aws-sdk-go/aws"
)

func TestRestRequest(t *testing.T) {
	var m sync.Mutex
	var httpReq *http.Request

	server := httptest.NewServer(http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			m.Lock()
			defer m.Unlock()

			httpReq = r

			fmt.Fprintln(w, `woo`)
		},
	))
	defer server.Close()

	client := aws.RestClient{
		Context: aws.Context{
			Service: "animals",
			Region:  "us-west-2",
			Credentials: aws.Creds(
				"accessKeyID",
				"secretAccessKey",
				"securityToken",
			),
		},
		Client: http.DefaultClient,
	}

	req, err := http.NewRequest("GET", server.URL+"/yay", nil)
	if err != nil {
		t.Fatal(err)
	}

	resp, err := client.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Fatal(err)
	}

	if v, want := string(body), "woo\n"; v != want {
		t.Errorf("Response entity was %q, but expected %q", v, want)
	}

	m.Lock()
	defer m.Unlock()

	if v, want := httpReq.Method, "GET"; v != want {
		t.Errorf("Method was %v but expected %v", v, want)
	}

	if httpReq.Header.Get("Authorization") == "" {
		t.Error("Authorization header is missing")
	}

	if v, want := httpReq.Header.Get("User-Agent"), "aws-go"; v != want {
		t.Errorf("User-Agent was %v but expected %v", v, want)
	}

	if v, want := httpReq.URL.String(), "/yay"; v != want {
		t.Errorf("URL was %v but expected %v", v, want)
	}
}

func TestRestRequestXMLError(t *testing.T) {
	var m sync.Mutex
	var httpReq *http.Request

	server := httptest.NewServer(http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			m.Lock()
			defer m.Unlock()

			httpReq = r

			w.Header().Set("Content-Type", "application/xml")
			w.WriteHeader(500)
			fmt.Fprintln(w, `<Error>
<Code>bonus</Code>
<BucketName>bingo</BucketName>
<Message>the bad thing</Message>
<RequestId>woo woo</RequestId>
<HostId>woo woo</HostId>
</Error>`)
		},
	))
	defer server.Close()

	client := aws.RestClient{
		Context: aws.Context{
			Service: "animals",
			Region:  "us-west-2",
			Credentials: aws.Creds(
				"accessKeyID",
				"secretAccessKey",
				"securityToken",
			),
		},
		Client: http.DefaultClient,
	}

	req, err := http.NewRequest("GET", server.URL+"/yay", nil)
	if err != nil {
		t.Fatal(err)
	}

	_, err = client.Do(req)
	if err == nil {
		t.Fatal("Expected an error but none was returned")
	}

	if err, ok := err.(aws.APIError); ok {
		if v, want := err.Code, "bonus"; v != want {
			t.Errorf("Error code was %v, but expected %v", v, want)
		}

		if v, want := err.Message, "the bad thing"; v != want {
			t.Errorf("Error message was %v, but expected %v", v, want)
		}
	} else {
		t.Errorf("Unknown error returned: %#v", err)
	}
}

func TestRestRequestJSONError(t *testing.T) {
	var m sync.Mutex
	var httpReq *http.Request

	server := httptest.NewServer(http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			m.Lock()
			defer m.Unlock()

			httpReq = r

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(500)
			fmt.Fprintln(w, `{"Code":"bonus", "Message":"the bad thing"}`)
		},
	))
	defer server.Close()

	client := aws.RestClient{
		Context: aws.Context{
			Service: "animals",
			Region:  "us-west-2",
			Credentials: aws.Creds(
				"accessKeyID",
				"secretAccessKey",
				"securityToken",
			),
		},
		Client: http.DefaultClient,
	}

	req, err := http.NewRequest("GET", server.URL+"/yay", nil)
	if err != nil {
		t.Fatal(err)
	}

	_, err = client.Do(req)
	if err == nil {
		t.Fatal("Expected an error but none was returned")
	}

	if err, ok := err.(aws.APIError); ok {
		if v, want := err.Code, "bonus"; v != want {
			t.Errorf("Error code was %v, but expected %v", v, want)
		}

		if v, want := err.Message, "the bad thing"; v != want {
			t.Errorf("Error message was %v, but expected %v", v, want)
		}
	} else {
		t.Errorf("Unknown error returned: %#v", err)
	}
}

func TestEscapePath(t *testing.T) {
	for _, x := range []struct {
		in   string
		want string
	}{
		{"", ""},
		{"ABCDEFGHIJKLMNOPQRTSUVWXYZ", "ABCDEFGHIJKLMNOPQRTSUVWXYZ"},
		{"abcdefghijklmnopqrtsuvwxyz", "abcdefghijklmnopqrtsuvwxyz"},
		{"0123456789", "0123456789"},
		{"_-~./:", "_-~./:"},
		{"test? file", "test%3F%20file"},
		{`hello? sausage/êé/Hello, 世界/ " ' @ < > & ?/z.txt`, "hello%3F%20sausage/%C3%AA%C3%A9/Hello%2C%20%E4%B8%96%E7%95%8C/%20%22%20%27%20%40%20%3C%20%3E%20%26%20%3F/z.txt"},
	} {
		got := aws.EscapePath(x.in)
		if got != x.want {
			t.Errorf("EscapePath(%q) got %q, want %v", x.in, got, x.want)
		}
	}
}
