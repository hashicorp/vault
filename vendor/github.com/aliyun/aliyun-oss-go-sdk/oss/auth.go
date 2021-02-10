package oss

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha1"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"hash"
	"io"
	"net/http"
	"sort"
	"strconv"
	"strings"
)

// headerSorter defines the key-value structure for storing the sorted data in signHeader.
type headerSorter struct {
	Keys []string
	Vals []string
}

// getAdditionalHeaderKeys get exist key in http header
func (conn Conn) getAdditionalHeaderKeys(req *http.Request) ([]string, map[string]string) {
	var keysList []string
	keysMap := make(map[string]string)
	srcKeys := make(map[string]string)

	for k := range req.Header {
		srcKeys[strings.ToLower(k)] = ""
	}

	for _, v := range conn.config.AdditionalHeaders {
		if _, ok := srcKeys[strings.ToLower(v)]; ok {
			keysMap[strings.ToLower(v)] = ""
		}
	}

	for k := range keysMap {
		keysList = append(keysList, k)
	}
	sort.Strings(keysList)
	return keysList, keysMap
}

// signHeader signs the header and sets it as the authorization header.
func (conn Conn) signHeader(req *http.Request, canonicalizedResource string) {
	akIf := conn.config.GetCredentials()
	authorizationStr := ""
	if conn.config.AuthVersion == AuthV2 {
		additionalList, _ := conn.getAdditionalHeaderKeys(req)
		if len(additionalList) > 0 {
			authorizationFmt := "OSS2 AccessKeyId:%v,AdditionalHeaders:%v,Signature:%v"
			additionnalHeadersStr := strings.Join(additionalList, ";")
			authorizationStr = fmt.Sprintf(authorizationFmt, akIf.GetAccessKeyID(), additionnalHeadersStr, conn.getSignedStr(req, canonicalizedResource, akIf.GetAccessKeySecret()))
		} else {
			authorizationFmt := "OSS2 AccessKeyId:%v,Signature:%v"
			authorizationStr = fmt.Sprintf(authorizationFmt, akIf.GetAccessKeyID(), conn.getSignedStr(req, canonicalizedResource, akIf.GetAccessKeySecret()))
		}
	} else {
		// Get the final authorization string
		authorizationStr = "OSS " + akIf.GetAccessKeyID() + ":" + conn.getSignedStr(req, canonicalizedResource, akIf.GetAccessKeySecret())
	}

	// Give the parameter "Authorization" value
	req.Header.Set(HTTPHeaderAuthorization, authorizationStr)
}

func (conn Conn) getSignedStr(req *http.Request, canonicalizedResource string, keySecret string) string {
	// Find out the "x-oss-"'s address in header of the request
	ossHeadersMap := make(map[string]string)
	additionalList, additionalMap := conn.getAdditionalHeaderKeys(req)
	for k, v := range req.Header {
		if strings.HasPrefix(strings.ToLower(k), "x-oss-") {
			ossHeadersMap[strings.ToLower(k)] = v[0]
		} else if conn.config.AuthVersion == AuthV2 {
			if _, ok := additionalMap[strings.ToLower(k)]; ok {
				ossHeadersMap[strings.ToLower(k)] = v[0]
			}
		}
	}
	hs := newHeaderSorter(ossHeadersMap)

	// Sort the ossHeadersMap by the ascending order
	hs.Sort()

	// Get the canonicalizedOSSHeaders
	canonicalizedOSSHeaders := ""
	for i := range hs.Keys {
		canonicalizedOSSHeaders += hs.Keys[i] + ":" + hs.Vals[i] + "\n"
	}

	// Give other parameters values
	// when sign URL, date is expires
	date := req.Header.Get(HTTPHeaderDate)
	contentType := req.Header.Get(HTTPHeaderContentType)
	contentMd5 := req.Header.Get(HTTPHeaderContentMD5)

	// default is v1 signature
	signStr := req.Method + "\n" + contentMd5 + "\n" + contentType + "\n" + date + "\n" + canonicalizedOSSHeaders + canonicalizedResource
	h := hmac.New(func() hash.Hash { return sha1.New() }, []byte(keySecret))

	// v2 signature
	if conn.config.AuthVersion == AuthV2 {
		signStr = req.Method + "\n" + contentMd5 + "\n" + contentType + "\n" + date + "\n" + canonicalizedOSSHeaders + strings.Join(additionalList, ";") + "\n" + canonicalizedResource
		h = hmac.New(func() hash.Hash { return sha256.New() }, []byte(keySecret))
	}

	// convert sign to log for easy to view
	if conn.config.LogLevel >= Debug {
		var signBuf bytes.Buffer
		for i := 0; i < len(signStr); i++ {
			if signStr[i] != '\n' {
				signBuf.WriteByte(signStr[i])
			} else {
				signBuf.WriteString("\\n")
			}
		}
		conn.config.WriteLog(Debug, "[Req:%p]signStr:%s\n", req, signBuf.String())
	}

	io.WriteString(h, signStr)
	signedStr := base64.StdEncoding.EncodeToString(h.Sum(nil))

	return signedStr
}

func (conn Conn) getRtmpSignedStr(bucketName, channelName, playlistName string, expiration int64, keySecret string, params map[string]interface{}) string {
	if params[HTTPParamAccessKeyID] == nil {
		return ""
	}

	canonResource := fmt.Sprintf("/%s/%s", bucketName, channelName)
	canonParamsKeys := []string{}
	for key := range params {
		if key != HTTPParamAccessKeyID && key != HTTPParamSignature && key != HTTPParamExpires && key != HTTPParamSecurityToken {
			canonParamsKeys = append(canonParamsKeys, key)
		}
	}

	sort.Strings(canonParamsKeys)
	canonParamsStr := ""
	for _, key := range canonParamsKeys {
		canonParamsStr = fmt.Sprintf("%s%s:%s\n", canonParamsStr, key, params[key].(string))
	}

	expireStr := strconv.FormatInt(expiration, 10)
	signStr := expireStr + "\n" + canonParamsStr + canonResource

	h := hmac.New(func() hash.Hash { return sha1.New() }, []byte(keySecret))
	io.WriteString(h, signStr)
	signedStr := base64.StdEncoding.EncodeToString(h.Sum(nil))
	return signedStr
}

// newHeaderSorter is an additional function for function SignHeader.
func newHeaderSorter(m map[string]string) *headerSorter {
	hs := &headerSorter{
		Keys: make([]string, 0, len(m)),
		Vals: make([]string, 0, len(m)),
	}

	for k, v := range m {
		hs.Keys = append(hs.Keys, k)
		hs.Vals = append(hs.Vals, v)
	}
	return hs
}

// Sort is an additional function for function SignHeader.
func (hs *headerSorter) Sort() {
	sort.Sort(hs)
}

// Len is an additional function for function SignHeader.
func (hs *headerSorter) Len() int {
	return len(hs.Vals)
}

// Less is an additional function for function SignHeader.
func (hs *headerSorter) Less(i, j int) bool {
	return bytes.Compare([]byte(hs.Keys[i]), []byte(hs.Keys[j])) < 0
}

// Swap is an additional function for function SignHeader.
func (hs *headerSorter) Swap(i, j int) {
	hs.Vals[i], hs.Vals[j] = hs.Vals[j], hs.Vals[i]
	hs.Keys[i], hs.Keys[j] = hs.Keys[j], hs.Keys[i]
}
