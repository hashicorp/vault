package http

import (
	"bufio"
	"bytes"
	"net/http"
	"reflect"
	"testing"
)

func TestForwardedRequestGenerateParse(t *testing.T) {
	bodBuf := bytes.NewBuffer([]byte(`{ "foo": "bar", "zip": { "argle": "bargle", neet: 0 } }`))
	req, err := http.NewRequest("FOOBAR", "https://pushit.real.good:9281/snicketysnack?furbleburble=bloopetybloop", bodBuf)
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Add(AuthHeaderName, "suppose_i_do")
	req.Header.Add(WrapTTLHeaderName, "43s")

	// We want to get the fields we would expect from an incoming request, so
	// we write it out and then read it again
	buf1 := bytes.NewBuffer(nil)
	err = req.Write(buf1)
	if err != nil {
		t.Fatal(err)
	}

	// Read it back in, parsing like a server
	bufr1 := bufio.NewReader(buf1)
	intreq, err := http.ReadRequest(bufr1)
	if err != nil {
		t.Fatal(err)
	}

	// Generate the request with the forwarded request in the body
	intreq, err = generateForwardedRequest(intreq, "https://bloopety.bloop:8201")
	if err != nil {
		t.Fatal(err)
	}

	// Perform another "round trip"
	buf2 := bytes.NewBuffer(nil)
	err = intreq.Write(buf2)
	if err != nil {
		t.Fatal(err)
	}
	bufr2 := bufio.NewReader(buf2)
	intreq, err = http.ReadRequest(bufr2)
	if err != nil {
		t.Fatal(err)
	}

	// Now extract the forwarded request to generate a final request for processing
	intreq, err = parseForwardedRequest(intreq)
	if err != nil {
		t.Fatal(err)
	}

	switch {
	case req.Method == intreq.Method:
	case req.RemoteAddr == intreq.RemoteAddr:
	case req.Host == intreq.Host:
	case reflect.DeepEqual(req.URL, intreq.URL):
	case reflect.DeepEqual(req.Header, intreq.Header):
	case reflect.DeepEqual(req.Body, intreq.Body):
	case reflect.DeepEqual(req.TLS, intreq.TLS):
	default:
		t.Fatalf("bad:\nreq:\n%#v\nintreq:\n%#v\n", *req, *intreq)
	}
}
