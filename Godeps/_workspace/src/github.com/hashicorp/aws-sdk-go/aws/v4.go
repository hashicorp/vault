package aws

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"sort"
	"strconv"
	"strings"
	"time"
)

const (
	authHeaderPrefix = "AWS4-HMAC-SHA256"
	timeFormat       = "20060102T150405Z"
	shortTimeFormat  = "20060102"
)

func (c *Context) sign(r *http.Request) error {
	creds, err := c.Credentials.Credentials()
	if err != nil {
		return err
	}

	date := r.Header.Get("Date")
	t := currentTime().UTC()
	if date != "" {
		var err error
		t, err = time.Parse(http.TimeFormat, date)
		if err != nil {
			return err
		}
	}

	s := signer{
		Request:         r,
		Time:            t,
		Body:            r.Body,
		ServiceName:     c.Service,
		Region:          c.Region,
		AccessKeyID:     creds.AccessKeyID,
		SecretAccessKey: creds.SecretAccessKey,
		SessionToken:    creds.SecurityToken,
		Debug:           0,
	}
	s.sign()
	return nil
}

type signer struct {
	Request         *http.Request
	Time            time.Time
	ServiceName     string
	Region          string
	AccessKeyID     string
	SecretAccessKey string
	SessionToken    string
	Body            io.Reader
	Debug           uint

	formattedTime      string
	formattedShortTime string

	signedHeaders    string
	canonicalHeaders string
	canonicalString  string
	credentialString string
	stringToSign     string
	signature        string
	authorization    string
}

func (v4 *signer) sign() {
	formatted := v4.Time.UTC().Format(timeFormat)

	// remove the old headers
	v4.Request.Header.Del("Date")
	v4.Request.Header.Del("Authorization")

	if v4.SessionToken != "" {
		v4.Request.Header.Set("X-Amz-Security-Token", v4.SessionToken)
	}

	v4.build()

	//v4.Debug = true
	if v4.Debug > 0 {
		fmt.Printf("---[ CANONICAL STRING  ]-----------------------------\n")
		fmt.Printf("%s\n", v4.canonicalString)
		fmt.Printf("-----------------------------------------------------\n\n")
		fmt.Printf("---[ STRING TO SIGN ]--------------------------------\n")
		fmt.Printf("%s\n", v4.stringToSign)
		fmt.Printf("-----------------------------------------------------\n")
	}

	// add the new ones
	v4.Request.Header.Set("Date", formatted)
	v4.Request.Header.Set("Authorization", v4.authorization)
}

func (v4 *signer) build() {
	v4.buildTime()
	v4.buildCanonicalHeaders()
	v4.buildCredentialString()
	v4.buildCanonicalString()
	v4.buildStringToSign()
	v4.buildSignature()
	v4.buildAuthorization()
}

func (v4 *signer) buildTime() {
	v4.formattedTime = v4.Time.UTC().Format(timeFormat)
	v4.formattedShortTime = v4.Time.UTC().Format(shortTimeFormat)
}

func (v4 *signer) buildAuthorization() {
	v4.authorization = strings.Join([]string{
		authHeaderPrefix + " Credential=" + v4.AccessKeyID + "/" + v4.credentialString,
		"SignedHeaders=" + v4.signedHeaders,
		"Signature=" + v4.signature,
	}, ",")
}

func (v4 *signer) buildCredentialString() {
	v4.credentialString = strings.Join([]string{
		v4.formattedShortTime,
		v4.Region,
		v4.ServiceName,
		"aws4_request",
	}, "/")
}

func (v4 *signer) buildCanonicalHeaders() {
	headers := make([]string, 0)
	headers = append(headers, "host")
	for k, _ := range v4.Request.Header {
		if http.CanonicalHeaderKey(k) == "Content-Length" {
			continue // never sign content-length
		}
		headers = append(headers, strings.ToLower(k))
	}
	sort.Strings(headers)

	headerValues := make([]string, len(headers))
	for i, k := range headers {
		if k == "host" {
			headerValues[i] = "host:" + v4.Request.URL.Host
		} else {
			headerValues[i] = k + ":" +
				strings.Join(v4.Request.Header[http.CanonicalHeaderKey(k)], ",")
		}
	}

	v4.signedHeaders = strings.Join(headers, ";")
	v4.canonicalHeaders = strings.Join(headerValues, "\n")
}

func (v4 *signer) buildCanonicalString() {
	v4.canonicalString = strings.Join([]string{
		v4.Request.Method,
		v4.Request.URL.Path,
		v4.Request.URL.Query().Encode(),
		v4.canonicalHeaders + "\n",
		v4.signedHeaders,
		v4.bodyDigest(),
	}, "\n")
}

func (v4 *signer) buildStringToSign() {
	v4.stringToSign = strings.Join([]string{
		authHeaderPrefix,
		v4.formattedTime,
		v4.credentialString,
		hexDigest(makeSha256([]byte(v4.canonicalString))),
	}, "\n")
}

func (v4 *signer) buildSignature() {
	secret := v4.SecretAccessKey
	date := makeHmac([]byte("AWS4"+secret), []byte(v4.formattedShortTime))
	region := makeHmac(date, []byte(v4.Region))
	service := makeHmac(region, []byte(v4.ServiceName))
	credentials := makeHmac(service, []byte("aws4_request"))
	signature := makeHmac(credentials, []byte(v4.stringToSign))
	v4.signature = hexDigest(signature)
}

func (v4 *signer) bodyDigest() string {
	hash := v4.Request.Header.Get("X-Amz-Content-Sha256")
	if hash == "" {
		if v4.Body == nil {
			hash = hexDigest(makeSha256([]byte{}))
		} else {
			// TODO refactor body to support seeking body payloads
			b, _ := ioutil.ReadAll(v4.Body)
			hash = hexDigest(makeSha256(b))
			v4.Request.Body = ioutil.NopCloser(bytes.NewReader(b))
		}
		v4.Request.Header.Add("X-Amz-Content-Sha256", hash)
	}
	return hash
}

func makeHmac(key []byte, data []byte) []byte {
	hash := hmac.New(sha256.New, key)
	hash.Write(data)
	return hash.Sum(nil)
}

func makeSha256(data []byte) []byte {
	hash := sha256.New()
	hash.Write(data)
	return hash.Sum(nil)
}

func makeSha256Reader(reader io.Reader) []byte {
	packet := make([]byte, 4096)
	hash := sha256.New()

	//reader.Seek(0, 0)
	for {
		n, err := reader.Read(packet)
		if n > 0 {
			hash.Write(packet[0:n])
		}
		if err == io.EOF || n == 0 {
			break
		}
	}
	//reader.Seek(0, 0)

	return hash.Sum(nil)
}

func hexDigest(data []byte) string {
	var buffer bytes.Buffer
	for i := range data {
		str := strconv.FormatUint(uint64(data[i]), 16)
		if len(str) < 2 {
			buffer.WriteString("0")
		}
		buffer.WriteString(str)
	}
	return buffer.String()
}
