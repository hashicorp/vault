package radius // import "layeh.com/radius"

import (
	"bytes"
	"net"
	"testing"
)

func Test_RFC2865_7_1(t *testing.T) {
	// Source: https://tools.ietf.org/html/rfc2865#section-7.1

	secret := []byte("xyzzy5461")

	// Request
	request := []byte{
		0x01, 0x00, 0x00, 0x38, 0x0f, 0x40, 0x3f, 0x94, 0x73, 0x97, 0x80, 0x57, 0xbd, 0x83, 0xd5, 0xcb,
		0x98, 0xf4, 0x22, 0x7a, 0x01, 0x06, 0x6e, 0x65, 0x6d, 0x6f, 0x02, 0x12, 0x0d, 0xbe, 0x70, 0x8d,
		0x93, 0xd4, 0x13, 0xce, 0x31, 0x96, 0xe4, 0x3f, 0x78, 0x2a, 0x0a, 0xee, 0x04, 0x06, 0xc0, 0xa8,
		0x01, 0x10, 0x05, 0x06, 0x00, 0x00, 0x00, 0x03,
	}

	p, err := Parse(request, secret, Builtin)
	if err != nil {
		t.Fatal(err)
	}
	if p.Code != CodeAccessRequest {
		t.Fatal("expecting Code = PacketCodeAccessRequest")
	}
	if p.Identifier != 0 {
		t.Fatal("expecting Identifier = 0")
	}
	if len(p.Attributes) != 4 {
		t.Fatal("expecting 4 attributes")
	}
	if p.String("User-Name") != "nemo" {
		t.Fatal("expecting User-Name = nemo")
	}
	if p.String("User-Password") != "arctangent" {
		t.Fatal("expecting User-Password = arctangent")
	}
	if username, password, ok := p.PAP(); !ok || username != "nemo" || password != "arctangent" {
		t.Fatal("PAP values do not match attributes")
	}
	if ip := p.Value("NAS-IP-Address").(net.IP); !ip.Equal(net.ParseIP("192.168.1.16")) {
		t.Fatal("expecting NAS-IP-Address = 192.168.1.16")
	}
	if p.Value("NAS-Port").(uint32) != uint32(3) {
		t.Fatal("expecting NAS-Port = 3")
	}

	{
		wire, err := p.Encode()
		if err != nil {
			t.Fatal("expecting p.Encode to succeed")
		}
		if !bytes.Equal(wire, request) {
			t.Fatal("expecting p.Encode() and request to equal")
		}
	}

	// Response
	response := []byte{
		0x02, 0x00, 0x00, 0x26, 0x86, 0xfe, 0x22, 0x0e, 0x76, 0x24, 0xba, 0x2a, 0x10, 0x05, 0xf6, 0xbf,
		0x9b, 0x55, 0xe0, 0xb2, 0x06, 0x06, 0x00, 0x00, 0x00, 0x01, 0x0f, 0x06, 0x00, 0x00, 0x00, 0x00,
		0x0e, 0x06, 0xc0, 0xa8, 0x01, 0x03,
	}

	q := Packet{
		Code:          CodeAccessAccept,
		Identifier:    p.Identifier,
		Authenticator: p.Authenticator,
		Secret:        secret,
		Dictionary:    p.Dictionary,
	}
	q.Set("Service-Type", uint32(1))
	q.Set("Login-Service", uint32(0))
	q.Set("Login-IP-Host", net.ParseIP("192.168.1.3"))

	{
		wire, err := q.Encode()
		if err != nil {
			t.Fatal(err)
		}
		if !bytes.Equal(response, wire) {
			t.Fatal("expecing response and q.Encode() to be equal")
		}
	}
}

func Test_RFC2865_7_2(t *testing.T) {
	// Source: https://tools.ietf.org/html/rfc2865#section-7.2

	secret := []byte("xyzzy5461")

	// Request
	request := []byte{
		0x01, 0x01, 0x00, 0x47, 0x2a, 0xee, 0x86, 0xf0, 0x8d, 0x0d, 0x55, 0x96, 0x9c, 0xa5, 0x97, 0x8e,
		0x0d, 0x33, 0x67, 0xa2, 0x01, 0x08, 0x66, 0x6c, 0x6f, 0x70, 0x73, 0x79, 0x03, 0x13, 0x16, 0xe9,
		0x75, 0x57, 0xc3, 0x16, 0x18, 0x58, 0x95, 0xf2, 0x93, 0xff, 0x63, 0x44, 0x07, 0x72, 0x75, 0x04,
		0x06, 0xc0, 0xa8, 0x01, 0x10, 0x05, 0x06, 0x00, 0x00, 0x00, 0x14, 0x06, 0x06, 0x00, 0x00, 0x00,
		0x02, 0x07, 0x06, 0x00, 0x00, 0x00, 0x01,
	}

	p, err := Parse(request, secret, Builtin)
	if err != nil {
		t.Fatal(err)
	}

	if p.Code != CodeAccessRequest {
		t.Fatal("expecting code access request")
	}
	if p.Identifier != 1 {
		t.Fatal("expecting Identifier = 1")
	}
	if p.String("User-Name") != "flopsy" {
		t.Fatal("expecting User-Name = flopsy")
	}
	if ip := p.Value("NAS-IP-Address").(net.IP); !ip.Equal(net.ParseIP("192.168.1.16")) {
		t.Fatal("expecting NAS-IP-Address = 192.168.1.16")
	}
	if p.Value("NAS-Port").(uint32) != uint32(20) {
		t.Fatal("expecting NAS-Port = 20")
	}
	if p.Value("Service-Type").(uint32) != uint32(2) {
		t.Fatal("expecting Service-Type = 2")
	}
	if p.Value("Framed-Protocol").(uint32) != uint32(1) {
		t.Fatal("expecting Framed-Protocol = 1")
	}
}

func TestPasswords(t *testing.T) {
	passwords := []string{
		"",
		"qwerty",
		"helloworld1231231231231233489hegufudhsgdsfygdf8g",
	}

	for _, password := range passwords {
		secret := []byte("xyzzy5461")

		r := New(CodeAccessRequest, secret)
		if r == nil {
			t.Fatal("could not create new RADIUS packet")
		}
		r.Add("User-Password", password)

		b, err := r.Encode()
		if err != nil {
			t.Fatal(err)
		}

		q, err := Parse(b, secret, Builtin)
		if err != nil {
			t.Fatal(err)
		}

		if s := q.String("User-Password"); s != password {
			t.Fatalf("incorrect User-Password (expecting %q, got %q)", password, s)
		}
	}
}
