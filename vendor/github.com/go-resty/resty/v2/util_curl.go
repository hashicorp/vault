package resty

import (
	"bytes"
	"io"
	"net/http"
	"net/http/cookiejar"

	"net/url"
	"strings"

	"github.com/go-resty/resty/v2/shellescape"
)

func buildCurlRequest(req *http.Request, httpCookiejar http.CookieJar) (curl string) {
	// 1. Generate curl raw headers

	curl = "curl -X " + req.Method + " "
	// req.Host + req.URL.Path + "?" + req.URL.RawQuery + " " + req.Proto + " "
	headers := dumpCurlHeaders(req)
	for _, kv := range *headers {
		curl += `-H ` + shellescape.Quote(kv[0]+": "+kv[1]) + ` `
	}

	// 2. Generate curl cookies
	// TODO validate this block of code, I think its not required since cookie captured via Headers
	if cookieJar, ok := httpCookiejar.(*cookiejar.Jar); ok {
		cookies := cookieJar.Cookies(req.URL)
		if len(cookies) > 0 {
			curl += `-H ` + shellescape.Quote(dumpCurlCookies(cookies)) + " "
		}
	}

	// 3. Generate curl body
	if req.Body != nil {
		buf, _ := io.ReadAll(req.Body)
		req.Body = io.NopCloser(bytes.NewBuffer(buf)) // important!!
		curl += `-d ` + shellescape.Quote(string(buf)) + " "
	}

	urlString := shellescape.Quote(req.URL.String())
	if urlString == "''" {
		urlString = "'http://unexecuted-request'"
	}
	curl += urlString
	return curl
}

// dumpCurlCookies dumps cookies to curl format
func dumpCurlCookies(cookies []*http.Cookie) string {
	sb := strings.Builder{}
	sb.WriteString("Cookie: ")
	for _, cookie := range cookies {
		sb.WriteString(cookie.Name + "=" + url.QueryEscape(cookie.Value) + "&")
	}
	return strings.TrimRight(sb.String(), "&")
}

// dumpCurlHeaders dumps headers to curl format
func dumpCurlHeaders(req *http.Request) *[][2]string {
	headers := [][2]string{}
	for k, vs := range req.Header {
		for _, v := range vs {
			headers = append(headers, [2]string{k, v})
		}
	}
	n := len(headers)
	for i := 0; i < n; i++ {
		for j := n - 1; j > i; j-- {
			jj := j - 1
			h1, h2 := headers[j], headers[jj]
			if h1[0] < h2[0] {
				headers[jj], headers[j] = headers[j], headers[jj]
			}
		}
	}
	return &headers
}
