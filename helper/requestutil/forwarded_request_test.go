package requestutil

import (
	"bufio"
	"bytes"
	"crypto/tls"
	"net/http"
	"reflect"
	"testing"
)

func TestForwardedRequestGenerateParse(t *testing.T) {
	bodBuf := bytes.NewReader([]byte(`{ "foo": "bar", "zip": { "argle": "bargle", neet: 0 } }`))
	req, err := http.NewRequest("FOOBAR", "https://pushit.real.good:9281/snicketysnack?furbleburble=bloopetybloop", bodBuf)
	if err != nil {
		t.Fatal(err)
	}
	req.TLS = &tls.ConnectionState{
		Version:           tls.VersionTLS12,
		HandshakeComplete: true,
		ServerName:        "tralala",
	}

	// We want to get the fields we would expect from an incoming request, so
	// we write it out and then read it again
	buf1 := bytes.NewBuffer(nil)
	err = req.Write(buf1)
	if err != nil {
		t.Fatal(err)
	}

	// Read it back in, parsing like a server
	bufr1 := bufio.NewReader(buf1)
	initialReq, err := http.ReadRequest(bufr1)
	if err != nil {
		t.Fatal(err)
	}

	// Generate the request with the forwarded request in the body
	req, err = GenerateForwardedRequest(initialReq, "https://bloopety.bloop:8201")
	if err != nil {
		t.Fatal(err)
	}

	// Perform another "round trip"
	buf2 := bytes.NewBuffer(nil)
	err = req.Write(buf2)
	if err != nil {
		t.Fatal(err)
	}
	bufr2 := bufio.NewReader(buf2)
	intreq, err := http.ReadRequest(bufr2)
	if err != nil {
		t.Fatal(err)
	}

	// Now extract the forwarded request to generate a final request for processing
	finalReq, err := ParseForwardedRequest(intreq)
	if err != nil {
		t.Fatal(err)
	}

	switch {
	case initialReq.Method != finalReq.Method:
		t.Fatalf("bad method:\ninitialReq:\n%#v\nfinalReq:\n%#v\n", *initialReq, *finalReq)
	case initialReq.RemoteAddr != finalReq.RemoteAddr:
		t.Fatalf("bad remoteaddr:\ninitialReq:\n%#v\nfinalReq:\n%#v\n", *initialReq, *finalReq)
	case initialReq.Host != finalReq.Host:
		t.Fatalf("bad host:\ninitialReq:\n%#v\nfinalReq:\n%#v\n", *initialReq, *finalReq)
	case !reflect.DeepEqual(initialReq.URL, finalReq.URL):
		t.Fatalf("bad url:\ninitialReq:\n%#v\nfinalReq:\n%#v\n", *initialReq.URL, *finalReq.URL)
	case !reflect.DeepEqual(initialReq.Header, finalReq.Header):
		t.Fatalf("bad header:\ninitialReq:\n%#v\nfinalReq:\n%#v\n", *initialReq, *finalReq)
	case !reflect.DeepEqual(initialReq.TLS, finalReq.TLS):
		t.Fatalf("bad tls:\ninitialReq:\n%#v\nfinalReq:\n%#v\n", *initialReq, *finalReq)
	default:
		// Compare bodies
		bodBuf.Seek(0, 0)
		initBuf := bytes.NewBuffer(nil)
		_, err = initBuf.ReadFrom(bodBuf)
		if err != nil {
			t.Fatal(err)
		}
		finBuf := bytes.NewBuffer(nil)
		_, err = finBuf.ReadFrom(finalReq.Body)
		if err != nil {
			t.Fatal(err)
		}
		if !reflect.DeepEqual(initBuf.Bytes(), finBuf.Bytes()) {
			t.Fatalf("badbody :\ninitialReq:\n%#v\nfinalReq:\n%#v\n", initBuf.Bytes(), finBuf.Bytes())
		}
	}
}
