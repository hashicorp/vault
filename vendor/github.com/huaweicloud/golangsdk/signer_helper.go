package golangsdk

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/textproto"
	"net/url"
	"regexp"
	"sort"
	"strings"
	"sync"
	"time"
)

// MemoryCache presents a thread safe memory cache
type MemoryCache struct {
	sync.Mutex                    // handling r/w for cache
	cacheHolder map[string]string // cache holder
	cacheKeys   []string          // cache keys
	MaxCount    int               // max cache entry count
}

// NewCache inits an new MemoryCache
func NewCache(maxCount int) *MemoryCache {
	return &MemoryCache{
		cacheHolder: make(map[string]string, maxCount),
		MaxCount:    maxCount,
	}
}

// Add an new cache item
func (cache *MemoryCache) Add(cacheKey string, cacheData string) {
	cache.Lock()
	defer cache.Unlock()

	if len(cache.cacheKeys) >= cache.MaxCount && len(cache.cacheKeys) > 1 {
		delete(cache.cacheHolder, cache.cacheKeys[0]) // delete first item
		cache.cacheKeys = append(cache.cacheKeys[1:]) // pop first one
	}

	cache.cacheHolder[cacheKey] = cacheData
	cache.cacheKeys = append(cache.cacheKeys, cacheKey)
}

// Get a cache item by its key
func (cache *MemoryCache) Get(cacheKey string) string {
	cache.Lock()
	defer cache.Unlock()

	return cache.cacheHolder[cacheKey]
}

//caseInsencitiveStringArray represents string case insensitive sorting operations
type caseInsencitiveStringArray []string

// noEscape specifies whether the character should be encoded or not
var noEscape [256]bool

func init() {
	// refer to https://docs.oracle.com/javase/7/docs/api/java/net/URLEncoder.html
	for i := 0; i < len(noEscape); i++ {
		noEscape[i] = (i >= 'A' && i <= 'Z') ||
			(i >= 'a' && i <= 'z') ||
			(i >= '0' && i <= '9') ||
			i == '.' ||
			i == '-' ||
			i == '_' ||
			i == '~' //java-sdk-core 3.0.1 HttpUtils.urlEncode
	}
}

// SignOptions represents the options during signing http request, it is concurency safely
type SignOptions struct {
	AccessKey           string //Access Key
	SecretKey           string //Secret key
	RegionName          string // Region name
	ServiceName         string // Service Name
	EnableCacheSignKey  bool   // Cache sign key for one day or not cache, cache is disabled by default
	encodeUrl           bool   //internal use
	SignAlgorithm       string //The algorithm used for sign, the default value is "SDK-HMAC-SHA256" if you don't set its value
	TimeOffsetInseconds int64  // TimeOffsetInseconds is used for adjust x-sdk-date if set its value
}

// StringBuilder wraps bytes.Buffer to implement a high performance string builder
type StringBuilder struct {
	builder bytes.Buffer //string storage
}

// reqSignParams represents the option values used for signing http request
type reqSignParams struct {
	SignOptions
	RequestTime time.Time
	Req         *http.Request
}

// signKeyCacheEntry represents the cache entry of sign key
type signKeyCacheEntry struct {
	Key                    []byte // sign key
	NumberOfDaysSinceEpoch int64  // number of days since epoch
}

// The default sign algorithm
const SignAlgorithmHMACSHA256 = "SDK-HMAC-SHA256"

// The header key of content hash value
const ContentSha256HeaderKey = "x-sdk-content-sha256"

//A regular for searching empty string
var spaceRegexp = regexp.MustCompile(`\s+`)

// cache sign key
var cache = NewCache(300)

//Sign manipulates the http.Request instance with some required authentication headers for SK/SK auth
func Sign(req *http.Request, signOptions SignOptions) {
	signOptions.AccessKey = strings.TrimSpace(signOptions.AccessKey)
	signOptions.SecretKey = strings.TrimSpace(signOptions.SecretKey)
	signOptions.encodeUrl = true

	signParams := reqSignParams{
		SignOptions: signOptions,
		RequestTime: time.Now(),
		Req:         req,
	}

	//t, _ := time.Parse(time.RFC3339, "2018-04-15T04:28:22+00:00")
	//signParams.RequestTime = t

	if signParams.SignAlgorithm == "" {
		signParams.SignAlgorithm = SignAlgorithmHMACSHA256
	}

	addRequiredHeaders(req, signParams.getFormattedSigningDateTime())
	contentSha256 := ""

	if v, ok := req.Header[textproto.CanonicalMIMEHeaderKey(ContentSha256HeaderKey)]; !ok {
		contentSha256 = calculateContentHash(req)
	} else {
		contentSha256 = v[0]
	}

	canonicalRequest := createCanonicalRequest(signParams, contentSha256)

	/*fmt.Println("canonicalRequest: " + canonicalRequest)
	fmt.Println("*****")*/

	strToSign := createStringToSign(canonicalRequest, signParams)
	signKey := deriveSigningKey(signParams)
	signature := computeSignature(strToSign, signKey, signParams.SignAlgorithm)

	req.Header.Set("Authorization", buildAuthorizationHeader(signParams, signature))
}

//ReSign manipulates the http.Request instance with some required authentication headers for SK/SK auth
func ReSign(req *http.Request, signOptions SignOptions) {
	signOptions.AccessKey = strings.TrimSpace(signOptions.AccessKey)
	signOptions.SecretKey = strings.TrimSpace(signOptions.SecretKey)
	signOptions.encodeUrl = true

	signParams := reqSignParams{
		SignOptions: signOptions,
		RequestTime: time.Now(),
		Req:         req,
	}

	if signParams.SignAlgorithm == "" {
		signParams.SignAlgorithm = SignAlgorithmHMACSHA256
	}

	setRequiredHeaders(req, signParams.getFormattedSigningDateTime())
	contentSha256 := ""

	if v, ok := req.Header[textproto.CanonicalMIMEHeaderKey(ContentSha256HeaderKey)]; !ok {
		contentSha256 = calculateContentHash(req)
	} else {
		contentSha256 = v[0]
	}

	canonicalRequest := createCanonicalRequest(signParams, contentSha256)

	strToSign := createStringToSign(canonicalRequest, signParams)
	signKey := deriveSigningKey(signParams)
	signature := computeSignature(strToSign, signKey, signParams.SignAlgorithm)

	req.Header.Set("Authorization", buildAuthorizationHeader(signParams, signature))
}

// deriveSigningKey returns a sign key from cache, or build it and insert it into cache
func deriveSigningKey(signParam reqSignParams) []byte {
	if signParam.EnableCacheSignKey {
		cacheKey := strings.Join([]string{signParam.SecretKey,
			signParam.RegionName,
			signParam.ServiceName,
		}, "-")

		cacheData := cache.Get(cacheKey)

		if cacheData != "" {
			var signKey signKeyCacheEntry
			json.Unmarshal([]byte(cacheData), &signKey)

			if signKey.NumberOfDaysSinceEpoch == signParam.getDaysSinceEpon() {
				return signKey.Key
			}
		}

		signKey := buildSignKey(signParam)
		signKeyStr, _ := json.Marshal(signKeyCacheEntry{
			Key:                    signKey,
			NumberOfDaysSinceEpoch: signParam.getDaysSinceEpon(),
		})
		cache.Add(cacheKey, string(signKeyStr))
		return signKey
	} else {
		return buildSignKey(signParam)
	}
}

func buildSignKey(signParam reqSignParams) []byte {
	var kSecret StringBuilder
	kSecret.Write("SDK").Write(signParam.SecretKey)

	kDate := computeSignature(signParam.getFormattedSigningDate(), kSecret.GetBytes(), signParam.SignAlgorithm)
	kRegion := computeSignature(signParam.RegionName, kDate, signParam.SignAlgorithm)
	kService := computeSignature(signParam.ServiceName, kRegion, signParam.SignAlgorithm)
	return computeSignature("sdk_request", kService, signParam.SignAlgorithm)
}

//HmacSha256 implements the  Keyed-Hash Message Authentication Code computation
func HmacSha256(data string, key []byte) []byte {
	mac := hmac.New(sha256.New, key)
	mac.Write([]byte(data))
	return mac.Sum(nil)
}

// HashSha256 is a wrapper for sha256 implementation
func HashSha256(msg []byte) []byte {
	sh256 := sha256.New()
	sh256.Write(msg)

	return sh256.Sum(nil)
}

// buildAuthorizationHeader builds the authentication header value
func buildAuthorizationHeader(signParam reqSignParams, signature []byte) string {
	var signingCredentials StringBuilder
	signingCredentials.Write(signParam.AccessKey).Write("/").Write(signParam.getScope())

	credential := "Credential=" + signingCredentials.ToString()
	signerHeaders := "SignedHeaders=" + getSignedHeadersString(signParam.Req)
	signatureHeader := "Signature=" + hex.EncodeToString(signature)

	return signParam.SignAlgorithm + " " + strings.Join([]string{
		credential,
		signerHeaders,
		signatureHeader,
	}, ", ")
}

// computeSignature computers the signature with the specified algorithm
// and it only supports SDK-HMAC-SHA256 in currently
func computeSignature(signData string, key []byte, algorithm string) []byte {
	if algorithm == SignAlgorithmHMACSHA256 {
		return HmacSha256(signData, key)
	} else {
		log.Fatalf("Unsupported algorithm %s, please use %s and try again", algorithm, SignAlgorithmHMACSHA256)
		return nil
	}
}

// createStringToSign build the need to be signed string
func createStringToSign(canonicalRequest string, signParams reqSignParams) string {
	return strings.Join([]string{signParams.SignAlgorithm,
		signParams.getFormattedSigningDateTime(),
		signParams.getScope(),
		hex.EncodeToString(HashSha256([]byte(canonicalRequest))),
	}, "\n")
}

// getCanonicalizedResourcePath builds the valid url path for signing
func getCanonicalizedResourcePath(signParas reqSignParams) string {
	urlStr := signParas.Req.URL.Path
	if !strings.HasPrefix(urlStr, "/") {
		urlStr = "/" + urlStr
	}

	if !strings.HasSuffix(urlStr, "/") {
		urlStr = urlStr + "/"
	}

	if signParas.encodeUrl {
		urlStr = urlEncode(urlStr, true)
	}

	if urlStr == "" {
		urlStr = "/"
	}

	return urlStr
}

// urlEncode encodes url path and url querystring according to the following rules:
// The alphanumeric characters "a" through "z", "A" through "Z" and "0" through "9" remain the same.
//The special characters ".", "-", "*", and "_" remain the same.
//The space character " " is converted into a plus sign "%20".
//All other characters are unsafe and are first converted into one or more bytes using some encoding scheme.
func urlEncode(url string, urlPath bool) string {
	var buf bytes.Buffer
	for i := 0; i < len(url); i++ {
		c := url[i]
		if noEscape[c] || (c == '/' && urlPath) {
			buf.WriteByte(c)
		} else {
			fmt.Fprintf(&buf, "%%%02X", c)
		}
	}

	return buf.String()
}

// encodeQueryString build and encode querystring to a string for signing
func encodeQueryString(queryValues url.Values) string {
	var encodedVals = make(map[string]string, len(queryValues))
	var keys = make([]string, len(queryValues))

	i := 0

	for k, _ := range queryValues {
		keys[i] = urlEncode(k, false)
		encodedVals[keys[i]] = k
		i++
	}

	caseInsensitiveSort(keys)

	var queryStr StringBuilder
	for i, k := range keys {
		if i > 0 {
			queryStr.Write("&")
		}

		queryStr.Write(k).Write("=").Write(urlEncode(queryValues.Get(encodedVals[k]), false))
	}

	return queryStr.ToString()
}

// getCanonicalizedQueryString return empty string if in POST method and content is nil, otherwise returns sorted,encoded querystring
func getCanonicalizedQueryString(signParas reqSignParams) string {
	if usePayloadForQueryParameters(signParas.Req) {
		return ""
	} else {
		return encodeQueryString(signParas.Req.URL.Query())
	}
}

// createCanonicalRequest builds canonical string depends the official document  for signing
func createCanonicalRequest(signParas reqSignParams, contentSha256 string) string {
	return strings.Join([]string{signParas.Req.Method,
		getCanonicalizedResourcePath(signParas),
		getCanonicalizedQueryString(signParas),
		getCanonicalizedHeaderString(signParas.Req),
		getSignedHeadersString(signParas.Req),
		contentSha256,
	}, "\n")
}

// calculateContentHash computes the content hash value
func calculateContentHash(req *http.Request) string {
	encodeParas := ""

	//post and content is null use queryString as content -- according to document
	if usePayloadForQueryParameters(req) {
		encodeParas = req.URL.Query().Encode()
	} else {
		if req.Body == nil {
			encodeParas = ""
		} else {
			readBody, _ := ioutil.ReadAll(req.Body)
			req.Body = ioutil.NopCloser(bytes.NewBuffer(readBody))
			encodeParas = string(readBody)
		}
	}

	return hex.EncodeToString(HashSha256([]byte(encodeParas)))
}

// usePayloadForQueryParameters specifies use querystring or not as content for compute content hash
func usePayloadForQueryParameters(req *http.Request) bool {
	if strings.ToLower(req.Method) != "post" {
		return false
	}

	return req.Body == nil
}

// getCanonicalizedHeaderString converts header map to a string for signing
func getCanonicalizedHeaderString(req *http.Request) string {
	var headers StringBuilder

	keys := make([]string, 0)
	for k, _ := range req.Header {
		keys = append(keys, strings.TrimSpace(k))
	}

	caseInsensitiveSort(keys)

	for _, k := range keys {
		k = strings.ToLower(k)
		newKey := spaceRegexp.ReplaceAllString(k, " ")
		headers.Write(newKey)
		headers.Write(":")

		val := req.Header.Get(k)
		val = spaceRegexp.ReplaceAllString(val, " ")
		headers.Write(val)

		headers.Write("\n")
	}

	return headers.ToString()
}

// getSignedHeadersString builds the string for AuthorizationHeader and signing
func getSignedHeadersString(req *http.Request) string {
	var headers StringBuilder

	keys := make([]string, 0)
	for k, _ := range req.Header {
		keys = append(keys, strings.TrimSpace(k))
	}

	caseInsensitiveSort(keys)

	for idx, k := range keys {

		if idx > 0 {
			headers.Write(";")
		}

		headers.Write(strings.ToLower(k))
	}

	return headers.ToString()
}

// addRequiredHeaders adds the required heads to http.request instance
func addRequiredHeaders(req *http.Request, timeStr string) {
	// golang handls port by default
	req.Header.Add("Host", req.URL.Host)
	req.Header.Add("X-Sdk-Date", timeStr)
}

// setRequiredHeaders sets the required heads to http.request for redirection
func setRequiredHeaders(req *http.Request, timeStr string) {
	req.Header.Set("X-Sdk-Date", timeStr)
	req.Header.Del("Authorization")
}

func (s caseInsencitiveStringArray) Len() int {
	return len(s)
}
func (s caseInsencitiveStringArray) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}
func (s caseInsencitiveStringArray) Less(i, j int) bool {
	return strings.ToLower(s[i]) < strings.ToLower(s[j])
}

func caseInsensitiveSort(strSlice []string) {
	sort.Sort(caseInsencitiveStringArray(strSlice))
}

func (signParas *reqSignParams) getSigningDateTimeMilli() int64 {
	return (signParas.RequestTime.UTC().Unix() - signParas.TimeOffsetInseconds) * 1000
}

func (signParas *reqSignParams) getSigningDateTime() time.Time {
	return time.Unix(signParas.getSigningDateTimeMilli()/1000, 0)
}

func (signParas *reqSignParams) getDaysSinceEpon() int64 {
	return signParas.getSigningDateTimeMilli() / 1000 / 3600 / 24
}

func (signParas *reqSignParams) getFormattedSigningDate() string {
	return signParas.getSigningDateTime().UTC().Format("20060102")
}
func (signParas *reqSignParams) getFormattedSigningDateTime() string {
	return signParas.getSigningDateTime().UTC().Format("20060102T150405Z")
}

func (signParas *reqSignParams) getScope() string {
	return strings.Join([]string{signParas.getFormattedSigningDate(),
		signParas.RegionName,
		signParas.ServiceName,
		"sdk_request",
	}, "/")
}

func (buff *StringBuilder) Write(s string) *StringBuilder {
	buff.builder.WriteString((s))
	return buff
}

func (buff *StringBuilder) ToString() string {
	return buff.builder.String()
}

func (buff *StringBuilder) GetBytes() []byte {
	return []byte(buff.ToString())
}
