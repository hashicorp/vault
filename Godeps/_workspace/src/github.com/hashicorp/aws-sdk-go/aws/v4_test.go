package aws

import (
	"net/http"
	"strings"
	"testing"
	"time"
)

func buildSigner(serviceName string, region string, signTime time.Time, body string) signer {
	endpoint := "https://" + serviceName + "." + region + ".amazonaws.com"
	reader := strings.NewReader(body)
	req, _ := http.NewRequest("POST", endpoint, reader)
	req.Header.Add("X-Amz-Target", "prefix.Operation")
	req.Header.Add("Content-Type", "application/x-amz-json-1.0")
	req.Header.Add("Content-Length", string(len(body)))

	return signer{
		Request:         req,
		Time:            signTime,
		Body:            reader,
		ServiceName:     serviceName,
		Region:          region,
		AccessKeyID:     "AKID",
		SecretAccessKey: "SECRET",
		SessionToken:    "SESSION",
	}
}

func removeWS(text string) string {
	text = strings.Replace(text, " ", "", -1)
	text = strings.Replace(text, "\n", "", -1)
	text = strings.Replace(text, "\t", "", -1)
	return text
}

func assertEqual(t *testing.T, expected, given string) {
	if removeWS(expected) != removeWS(given) {
		t.Errorf("\nExpected: %s\nGiven:    %s", expected, given)
	}
}

func TestSignRequest(t *testing.T) {
	signer := buildSigner("dynamodb", "us-east-1", time.Unix(0, 0), "{}")
	signer.sign()

	expectedDate := "19700101T000000Z"
	expectedAuth := `
    AWS4-HMAC-SHA256
    Credential=AKID/19700101/us-east-1/dynamodb/aws4_request,
    SignedHeaders=content-type;host;x-amz-security-token;x-amz-target,
    Signature=4662104789134800e088b6a2bf3a1153ca7d38ecfc07a69bff2859f04900b67f
  `

	assertEqual(t, expectedAuth, signer.Request.Header.Get("Authorization"))
	assertEqual(t, expectedDate, signer.Request.Header.Get("Date"))
}

func BenchmarkSignRequest(b *testing.B) {
	signer := buildSigner("dynamodb", "us-east-1", time.Now(), "{}")
	for i := 0; i < b.N; i++ {
		signer.sign()
	}
}
