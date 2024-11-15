package oss

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha1"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"hash"
	"io"
	"net/http"
	"sort"
	"strconv"
	"strings"
	"time"
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

// getAdditionalHeaderKeysV4 get exist key in http header
func (conn Conn) getAdditionalHeaderKeysV4(req *http.Request) ([]string, map[string]string) {
	var keysList []string
	keysMap := make(map[string]string)
	srcKeys := make(map[string]string)

	for k := range req.Header {
		srcKeys[strings.ToLower(k)] = ""
	}

	for _, v := range conn.config.AdditionalHeaders {
		if _, ok := srcKeys[strings.ToLower(v)]; ok {
			if !strings.EqualFold(v, HTTPHeaderContentMD5) && !strings.EqualFold(v, HTTPHeaderContentType) {
				keysMap[strings.ToLower(v)] = ""
			}
		}
	}

	for k := range keysMap {
		keysList = append(keysList, k)
	}
	sort.Strings(keysList)
	return keysList, keysMap
}

// signHeader signs the header and sets it as the authorization header.
func (conn Conn) signHeader(req *http.Request, canonicalizedResource string, credentials Credentials) {
	akIf := credentials
	authorizationStr := ""
	if conn.config.AuthVersion == AuthV4 {
		strDay := ""
		strDate := req.Header.Get(HttpHeaderOssDate)
		if strDate == "" {
			strDate = req.Header.Get(HTTPHeaderDate)
			t, _ := time.Parse(http.TimeFormat, strDate)
			strDay = t.Format("20060102")
		} else {
			t, _ := time.Parse(timeFormatV4, strDate)
			strDay = t.Format("20060102")
		}
		signHeaderProduct := conn.config.GetSignProduct()
		signHeaderRegion := conn.config.GetSignRegion()

		additionalList, _ := conn.getAdditionalHeaderKeysV4(req)
		if len(additionalList) > 0 {
			authorizationFmt := "OSS4-HMAC-SHA256 Credential=%v/%v/%v/" + signHeaderProduct + "/aliyun_v4_request,AdditionalHeaders=%v,Signature=%v"
			additionnalHeadersStr := strings.Join(additionalList, ";")
			authorizationStr = fmt.Sprintf(authorizationFmt, akIf.GetAccessKeyID(), strDay, signHeaderRegion, additionnalHeadersStr, conn.getSignedStrV4(req, canonicalizedResource, akIf.GetAccessKeySecret(), nil))
		} else {
			authorizationFmt := "OSS4-HMAC-SHA256 Credential=%v/%v/%v/" + signHeaderProduct + "/aliyun_v4_request,Signature=%v"
			authorizationStr = fmt.Sprintf(authorizationFmt, akIf.GetAccessKeyID(), strDay, signHeaderRegion, conn.getSignedStrV4(req, canonicalizedResource, akIf.GetAccessKeySecret(), nil))
		}
	} else if conn.config.AuthVersion == AuthV2 {
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

	if conn.config.LogLevel >= Debug {
		conn.config.WriteLog(Debug, "[Req:%p]signStr:%s\n", req, EscapeLFString(signStr))
	}

	io.WriteString(h, signStr)
	signedStr := base64.StdEncoding.EncodeToString(h.Sum(nil))

	return signedStr
}

func (conn Conn) getSignedStrV4(req *http.Request, canonicalizedResource string, keySecret string, signingTime *time.Time) string {
	// Find out the "x-oss-"'s address in header of the request
	ossHeadersMap := make(map[string]string)
	additionalList, additionalMap := conn.getAdditionalHeaderKeysV4(req)
	for k, v := range req.Header {
		lowKey := strings.ToLower(k)
		if strings.EqualFold(lowKey, HTTPHeaderContentMD5) ||
			strings.EqualFold(lowKey, HTTPHeaderContentType) ||
			strings.HasPrefix(lowKey, "x-oss-") {
			ossHeadersMap[lowKey] = strings.Trim(v[0], " ")
		} else {
			if _, ok := additionalMap[lowKey]; ok {
				ossHeadersMap[lowKey] = strings.Trim(v[0], " ")
			}
		}
	}

	// get day,eg 20210914
	//signingTime
	signDate := ""
	strDay := ""
	if signingTime != nil {
		signDate = signingTime.Format(timeFormatV4)
		strDay = signingTime.Format(shortTimeFormatV4)
	} else {
		var t time.Time
		// Required parameters
		if date := req.Header.Get(HTTPHeaderDate); date != "" {
			signDate = date
			t, _ = time.Parse(http.TimeFormat, date)
		}

		if ossDate := req.Header.Get(HttpHeaderOssDate); ossDate != "" {
			signDate = ossDate
			t, _ = time.Parse(timeFormatV4, ossDate)
		}

		strDay = t.Format("20060102")
	}

	hs := newHeaderSorter(ossHeadersMap)

	// Sort the ossHeadersMap by the ascending order
	hs.Sort()

	// Get the canonicalizedOSSHeaders
	canonicalizedOSSHeaders := ""
	for i := range hs.Keys {
		canonicalizedOSSHeaders += hs.Keys[i] + ":" + hs.Vals[i] + "\n"
	}

	signStr := ""

	// v4 signature
	hashedPayload := DefaultContentSha256
	if val := req.Header.Get(HttpHeaderOssContentSha256); val != "" {
		hashedPayload = val
	}

	// subResource
	resource := canonicalizedResource
	subResource := ""
	subPos := strings.LastIndex(canonicalizedResource, "?")
	if subPos != -1 {
		subResource = canonicalizedResource[subPos+1:]
		resource = canonicalizedResource[0:subPos]
	}

	// get canonical request
	canonicalReuqest := req.Method + "\n" + resource + "\n" + subResource + "\n" + canonicalizedOSSHeaders + "\n" + strings.Join(additionalList, ";") + "\n" + hashedPayload
	rh := sha256.New()
	io.WriteString(rh, canonicalReuqest)
	hashedRequest := hex.EncodeToString(rh.Sum(nil))

	if conn.config.LogLevel >= Debug {
		conn.config.WriteLog(Debug, "[Req:%p]CanonicalRequest:%s\n", req, EscapeLFString(canonicalReuqest))
	}

	// Product & Region
	signedStrV4Product := conn.config.GetSignProduct()
	signedStrV4Region := conn.config.GetSignRegion()

	signStr = "OSS4-HMAC-SHA256" + "\n" + signDate + "\n" + strDay + "/" + signedStrV4Region + "/" + signedStrV4Product + "/aliyun_v4_request" + "\n" + hashedRequest
	if conn.config.LogLevel >= Debug {
		conn.config.WriteLog(Debug, "[Req:%p]signStr:%s\n", req, EscapeLFString(signStr))
	}

	h1 := hmac.New(func() hash.Hash { return sha256.New() }, []byte("aliyun_v4"+keySecret))
	io.WriteString(h1, strDay)
	h1Key := h1.Sum(nil)

	h2 := hmac.New(func() hash.Hash { return sha256.New() }, h1Key)
	io.WriteString(h2, signedStrV4Region)
	h2Key := h2.Sum(nil)

	h3 := hmac.New(func() hash.Hash { return sha256.New() }, h2Key)
	io.WriteString(h3, signedStrV4Product)
	h3Key := h3.Sum(nil)

	h4 := hmac.New(func() hash.Hash { return sha256.New() }, h3Key)
	io.WriteString(h4, "aliyun_v4_request")
	h4Key := h4.Sum(nil)

	h := hmac.New(func() hash.Hash { return sha256.New() }, h4Key)
	io.WriteString(h, signStr)
	return fmt.Sprintf("%x", h.Sum(nil))
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
