// Copyright (c) 2017-2021 Snowflake Computing Inc. All right reserved.

package gosnowflake

import (
	"bufio"
	"context"
	"crypto"
	"crypto/tls"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/asn1"
	"encoding/base64"
	"encoding/json"
	"encoding/pem"
	"errors"
	"fmt"
	"golang.org/x/crypto/ocsp"
	"io"
	"io/ioutil"
	"math/big"
	"net"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
	"time"
)

var (
	// caRoot includes the CA certificates.
	caRoot map[string]*x509.Certificate
	// certPOol includes the CA certificates.
	certPool *x509.CertPool
	// cacheDir is the location of OCSP response cache file
	cacheDir = ""
	// cacheFileName is the file name of OCSP response cache file
	cacheFileName = ""
	// cacheUpdated is true if the memory cache is updated
	cacheUpdated = true
)

// OCSPFailOpenMode is OCSP fail open mode. OCSPFailOpenTrue by default and may
// set to ocspModeFailClosed for fail closed mode
type OCSPFailOpenMode uint32

const (
	ocspFailOpenNotSet OCSPFailOpenMode = iota
	// OCSPFailOpenTrue represents OCSP fail open mode.
	OCSPFailOpenTrue
	// OCSPFailOpenFalse represents OCSP fail closed mode.
	OCSPFailOpenFalse
)
const (
	ocspModeFailOpen   = "FAIL_OPEN"
	ocspModeFailClosed = "FAIL_CLOSED"
	ocspModeInsecure   = "INSECURE"
)

// OCSP fail open mode
var ocspFailOpen = OCSPFailOpenTrue

const (
	// defaultOCSPCacheServerTimeout is the total timeout for OCSP cache server.
	defaultOCSPCacheServerTimeout = 5 * time.Second

	// defaultOCSPResponderTimeout is the total timeout for OCSP responder.
	defaultOCSPResponderTimeout = 10 * time.Second
)

const (
	cacheFileBaseName = "ocsp_response_cache.json"
	// cacheExpire specifies cache data expiration time in seconds.
	cacheExpire           = float64(24 * 60 * 60)
	cacheServerURL        = "http://ocsp.snowflakecomputing.com"
	cacheServerEnabledEnv = "SF_OCSP_RESPONSE_CACHE_SERVER_ENABLED"
	cacheServerURLEnv     = "SF_OCSP_RESPONSE_CACHE_SERVER_URL"
	cacheDirEnv           = "SF_OCSP_RESPONSE_CACHE_DIR"
	ocspRetryURLEnv       = "SF_OCSP_RESPONSE_RETRY_URL"
)

const (
	ocspTestInjectValidityErrorEnv        = "SF_OCSP_TEST_INJECT_VALIDITY_ERROR"
	ocspTestInjectUnknownStatusEnv        = "SF_OCSP_TEST_INJECT_UNKNOWN_STATUS"
	ocspTestResponseCacheServerTimeoutEnv = "SF_OCSP_TEST_OCSP_RESPONSE_CACHE_SERVER_TIMEOUT"
	ocspTestResponderTimeoutEnv           = "SF_OCSP_TEST_OCSP_RESPONDER_TIMEOUT"
	ocspTestResponderURLEnv               = "SF_OCSP_TEST_RESPONDER_URL"
	ocspTestNoOCSPURLEnv                  = "SF_OCSP_TEST_NO_OCSP_RESPONDER_URL"
)

const (
	tolerableValidityRatio = 100               // buffer for certificate revocation update time
	maxClockSkew           = 900 * time.Second // buffer for clock skew
)

type ocspStatusCode int

type ocspStatus struct {
	code ocspStatusCode
	err  error
}

const (
	ocspSuccess                ocspStatusCode = 0
	ocspStatusGood             ocspStatusCode = -1
	ocspStatusRevoked          ocspStatusCode = -2
	ocspStatusUnknown          ocspStatusCode = -3
	ocspStatusOthers           ocspStatusCode = -4
	ocspNoServer               ocspStatusCode = -5
	ocspFailedParseOCSPHost    ocspStatusCode = -6
	ocspFailedComposeRequest   ocspStatusCode = -7
	ocspFailedDecomposeRequest ocspStatusCode = -8
	ocspFailedSubmit           ocspStatusCode = -9
	ocspFailedResponse         ocspStatusCode = -10
	ocspFailedExtractResponse  ocspStatusCode = -11
	ocspFailedParseResponse    ocspStatusCode = -12
	ocspInvalidValidity        ocspStatusCode = -13
	ocspMissedCache            ocspStatusCode = -14
	ocspCacheExpired           ocspStatusCode = -15
	ocspFailedDecodeResponse   ocspStatusCode = -16
)

// copied from crypto/ocsp.go
type certID struct {
	HashAlgorithm pkix.AlgorithmIdentifier
	NameHash      []byte
	IssuerKeyHash []byte
	SerialNumber  *big.Int
}

// cache key
type certIDKey struct {
	HashAlgorithm crypto.Hash
	NameHash      string
	IssuerKeyHash string
	SerialNumber  string
}

var (
	ocspResponseCache     map[certIDKey][]interface{}
	ocspResponseCacheLock *sync.RWMutex
)

// copied from crypto/ocsp
var hashOIDs = map[crypto.Hash]asn1.ObjectIdentifier{
	crypto.SHA1:   asn1.ObjectIdentifier([]int{1, 3, 14, 3, 2, 26}),
	crypto.SHA256: asn1.ObjectIdentifier([]int{2, 16, 840, 1, 101, 3, 4, 2, 1}),
	crypto.SHA384: asn1.ObjectIdentifier([]int{2, 16, 840, 1, 101, 3, 4, 2, 2}),
	crypto.SHA512: asn1.ObjectIdentifier([]int{2, 16, 840, 1, 101, 3, 4, 2, 3}),
}

// copied from crypto/ocsp
func getOIDFromHashAlgorithm(target crypto.Hash) asn1.ObjectIdentifier {
	for hash, oid := range hashOIDs {
		if hash == target {
			return oid
		}
	}
	logger.Errorf("no valid OID is found for the hash algorithm. %#v", target)
	return nil
}

func getHashAlgorithmFromOID(target pkix.AlgorithmIdentifier) crypto.Hash {
	for hash, oid := range hashOIDs {
		if oid.Equal(target.Algorithm) {
			return hash
		}
	}
	logger.Errorf("no valid hash algorithm is found for the oid. Falling back to SHA1: %#v", target)
	return crypto.SHA1
}

// calcTolerableValidity returns the maximum validity buffer
func calcTolerableValidity(thisUpdate, nextUpdate time.Time) time.Duration {
	return durationMax(time.Duration(nextUpdate.Sub(thisUpdate)/tolerableValidityRatio), maxClockSkew)
}

// isInValidityRange checks the validity
func isInValidityRange(currTime, thisUpdate, nextUpdate time.Time) bool {
	if currTime.Sub(thisUpdate.Add(-maxClockSkew)) < 0 {
		return false
	}
	if nextUpdate.Add(calcTolerableValidity(thisUpdate, nextUpdate)).Sub(currTime) < 0 {
		return false
	}
	return true
}

func isTestInvalidValidity() bool {
	return strings.EqualFold(os.Getenv(ocspTestInjectValidityErrorEnv), "true")
}

func extractCertIDKeyFromRequest(ocspReq []byte) (*certIDKey, *ocspStatus) {
	r, err := ocsp.ParseRequest(ocspReq)
	if err != nil {
		return nil, &ocspStatus{
			code: ocspFailedDecomposeRequest,
			err:  err,
		}
	}

	// encode CertID, used as a key in the cache
	encodedCertID := &certIDKey{
		r.HashAlgorithm,
		base64.StdEncoding.EncodeToString(r.IssuerNameHash),
		base64.StdEncoding.EncodeToString(r.IssuerKeyHash),
		r.SerialNumber.String(),
	}
	return encodedCertID, &ocspStatus{
		code: ocspSuccess,
	}
}

func encodeCertIDKey(certIDKeyBase64 string) *certIDKey {
	r, err := base64.StdEncoding.DecodeString(certIDKeyBase64)
	if err != nil {
		return nil
	}
	var c certID
	rest, err := asn1.Unmarshal(r, &c)
	if err != nil {
		// error in parsing
		return nil
	}
	if len(rest) > 0 {
		// extra bytes to the end
		return nil
	}
	return &certIDKey{
		getHashAlgorithmFromOID(c.HashAlgorithm),
		base64.StdEncoding.EncodeToString(c.NameHash),
		base64.StdEncoding.EncodeToString(c.IssuerKeyHash),
		c.SerialNumber.String(),
	}
}

func decodeCertIDKey(k *certIDKey) string {
	serialNumber := new(big.Int)
	serialNumber.SetString(k.SerialNumber, 10)
	nameHash, err := base64.StdEncoding.DecodeString(k.NameHash)
	if err != nil {
		return ""
	}
	issuerKeyHash, err := base64.StdEncoding.DecodeString(k.IssuerKeyHash)
	if err != nil {
		return ""
	}
	encodedCertID, err := asn1.Marshal(certID{
		pkix.AlgorithmIdentifier{
			Algorithm:  getOIDFromHashAlgorithm(k.HashAlgorithm),
			Parameters: asn1.RawValue{Tag: 5 /* ASN.1 NULL */},
		},
		nameHash,
		issuerKeyHash,
		serialNumber,
	})
	if err != nil {
		return ""
	}
	return base64.StdEncoding.EncodeToString(encodedCertID)
}

func checkOCSPResponseCache(encodedCertID *certIDKey, subject, issuer *x509.Certificate) *ocspStatus {
	if strings.EqualFold(os.Getenv(cacheServerEnabledEnv), "false") {
		return &ocspStatus{code: ocspNoServer}
	}
	ocspResponseCacheLock.RLock()
	gotValueFromCache := ocspResponseCache[*encodedCertID]
	ocspResponseCacheLock.RUnlock()

	status := extractOCSPCacheResponseValue(gotValueFromCache, subject, issuer)
	if !isValidOCSPStatus(status.code) {
		deleteOCSPCache(encodedCertID)
	}
	return status
}

func deleteOCSPCache(encodedCertID *certIDKey) {
	ocspResponseCacheLock.Lock()
	delete(ocspResponseCache, *encodedCertID)
	cacheUpdated = true
	ocspResponseCacheLock.Unlock()
}

// deleteOCSPCacheAll deletes all entries in the OCSP response cache on memory
func deleteOCSPCacheAll() {
	ocspResponseCacheLock.Lock()
	defer ocspResponseCacheLock.Unlock()
	ocspResponseCache = make(map[certIDKey][]interface{})
}

func validateOCSP(ocspRes *ocsp.Response) *ocspStatus {
	curTime := time.Now()

	if ocspRes == nil {
		return &ocspStatus{
			code: ocspFailedDecomposeRequest,
			err:  errors.New("OCSP Response is nil"),
		}
	}
	if isTestInvalidValidity() || !isInValidityRange(curTime, ocspRes.ThisUpdate, ocspRes.NextUpdate) {
		return &ocspStatus{
			code: ocspInvalidValidity,
			err: &SnowflakeError{
				Number:      ErrOCSPInvalidValidity,
				Message:     errMsgOCSPInvalidValidity,
				MessageArgs: []interface{}{ocspRes.ProducedAt, ocspRes.ThisUpdate, ocspRes.NextUpdate},
			},
		}
	}
	if isTestUnknownStatus() {
		ocspRes.Status = ocsp.Unknown
	}
	return returnOCSPStatus(ocspRes)
}

func returnOCSPStatus(ocspRes *ocsp.Response) *ocspStatus {
	switch ocspRes.Status {
	case ocsp.Good:
		return &ocspStatus{
			code: ocspStatusGood,
			err:  nil,
		}
	case ocsp.Revoked:
		return &ocspStatus{
			code: ocspStatusRevoked,
			err: &SnowflakeError{
				Number:      ErrOCSPStatusRevoked,
				Message:     errMsgOCSPStatusRevoked,
				MessageArgs: []interface{}{ocspRes.RevocationReason, ocspRes.RevokedAt},
			},
		}
	case ocsp.Unknown:
		return &ocspStatus{
			code: ocspStatusUnknown,
			err: &SnowflakeError{
				Number:  ErrOCSPStatusUnknown,
				Message: errMsgOCSPStatusUnknown,
			},
		}
	default:
		return &ocspStatus{
			code: ocspStatusOthers,
			err:  fmt.Errorf("OCSP others. %v", ocspRes.Status),
		}
	}
}

func isTestUnknownStatus() bool {
	return strings.EqualFold(os.Getenv(ocspTestInjectUnknownStatusEnv), "true")
}

func checkOCSPCacheServer(
	ctx context.Context,
	client clientInterface,
	req requestFunc,
	ocspServerHost *url.URL,
	totalTimeout time.Duration) (
	cacheContent *map[string][]interface{},
	ocspS *ocspStatus) {
	var respd map[string][]interface{}
	headers := make(map[string]string)
	res, err := newRetryHTTP(ctx, client, req, ocspServerHost, headers, totalTimeout).execute()
	if err != nil {
		logger.Errorf("failed to get OCSP cache from OCSP Cache Server. %v\n", err)
		return nil, &ocspStatus{
			code: ocspFailedSubmit,
			err:  err,
		}
	}
	defer res.Body.Close()
	logger.Debugf("StatusCode from OCSP Cache Server: %v\n", res.StatusCode)
	if res.StatusCode != http.StatusOK {
		return nil, &ocspStatus{
			code: ocspFailedResponse,
			err:  fmt.Errorf("HTTP code is not OK. %v: %v", res.StatusCode, res.Status),
		}
	}
	logger.Debugf("reading contents")

	dec := json.NewDecoder(res.Body)
	for {
		if err := dec.Decode(&respd); err == io.EOF {
			break
		} else if err != nil {
			logger.Errorf("failed to decode OCSP cache. %v\n", err)
			return nil, &ocspStatus{
				code: ocspFailedExtractResponse,
				err:  err,
			}
		}
	}
	return &respd, &ocspStatus{
		code: ocspSuccess,
	}
}

// retryOCSP is the second level of retry method if the returned contents are corrupted. It often happens with OCSP
// serer and retry helps.
func retryOCSP(
	ctx context.Context,
	client clientInterface,
	req requestFunc,
	ocspHost *url.URL,
	headers map[string]string,
	reqBody []byte,
	issuer *x509.Certificate,
	totalTimeout time.Duration) (
	ocspRes *ocsp.Response,
	ocspResBytes []byte,
	ocspS *ocspStatus) {
	multiplier := 1
	if atomic.LoadUint32((*uint32)(&ocspFailOpen)) == (uint32)(OCSPFailOpenFalse) {
		multiplier = 3 // up to 3 times for Fail Close mode
	}
	res, err := newRetryHTTP(
		ctx, client, req, ocspHost, headers,
		totalTimeout*time.Duration(multiplier)).doPost().setBody(reqBody).execute()
	if err != nil {
		return ocspRes, ocspResBytes, &ocspStatus{
			code: ocspFailedSubmit,
			err:  err,
		}
	}
	defer res.Body.Close()
	logger.Debugf("StatusCode from OCSP Server: %v\n", res.StatusCode)
	if res.StatusCode != http.StatusOK {
		return ocspRes, ocspResBytes, &ocspStatus{
			code: ocspFailedResponse,
			err:  fmt.Errorf("HTTP code is not OK. %v: %v", res.StatusCode, res.Status),
		}
	}
	logger.Debug("reading contents")
	ocspResBytes, err = ioutil.ReadAll(res.Body)
	if err != nil {
		return ocspRes, ocspResBytes, &ocspStatus{
			code: ocspFailedExtractResponse,
			err:  err,
		}
	}
	logger.Debug("parsing OCSP response")
	ocspRes, err = ocsp.ParseResponse(ocspResBytes, issuer)
	if err != nil {
		return ocspRes, ocspResBytes, &ocspStatus{
			code: ocspFailedParseResponse,
			err:  err,
		}
	}

	return ocspRes, ocspResBytes, &ocspStatus{
		code: ocspSuccess,
	}
}

// getRevocationStatus checks the certificate revocation status for subject using issuer certificate.
func getRevocationStatus(ctx context.Context, subject, issuer *x509.Certificate) *ocspStatus {
	logger.Infof("Subject: %v, Issuer: %v\n", subject.Subject, issuer.Subject)

	status, ocspReq, encodedCertID := validateWithCache(subject, issuer)
	if isValidOCSPStatus(status.code) {
		return status
	}
	if ocspReq == nil || encodedCertID == nil {
		return status
	}
	logger.Infof("cache missed\n")
	logger.Infof("OCSP Server: %v\n", subject.OCSPServer)
	if len(subject.OCSPServer) == 0 || isTestNoOCSPURL() {
		return &ocspStatus{
			code: ocspNoServer,
			err: &SnowflakeError{
				Number:      ErrOCSPNoOCSPResponderURL,
				Message:     errMsgOCSPNoOCSPResponderURL,
				MessageArgs: []interface{}{subject.Subject},
			},
		}
	}
	ocspHost := subject.OCSPServer[0]
	u, err := url.Parse(ocspHost)
	if err != nil {
		return &ocspStatus{
			code: ocspFailedParseOCSPHost,
			err:  fmt.Errorf("failed to parse OCSP server host. %v", ocspHost),
		}
	}
	hostnameStr := os.Getenv(ocspTestResponderURLEnv)
	var hostname string
	if retryURL := os.Getenv(ocspRetryURLEnv); retryURL != "" {
		hostname = fmt.Sprintf(retryURL, u.Hostname(), base64.StdEncoding.EncodeToString(ocspReq))
	} else {
		hostname = u.Hostname()
	}
	if hostnameStr != "" {
		u0, err := url.Parse(hostnameStr)
		if err == nil {
			hostname = u0.Hostname()
			u = u0
		}
	}
	headers := make(map[string]string)
	headers[httpHeaderContentType] = "application/ocsp-request"
	headers[httpHeaderAccept] = "application/ocsp-response"
	headers[httpHeaderContentLength] = strconv.Itoa(len(ocspReq))
	headers[httpHeaderHost] = hostname
	timeoutStr := os.Getenv(ocspTestResponderTimeoutEnv)
	timeout := defaultOCSPResponderTimeout
	if timeoutStr != "" {
		var timeoutInMilliseconds int
		timeoutInMilliseconds, err = strconv.Atoi(timeoutStr)
		if err == nil {
			timeout = time.Duration(timeoutInMilliseconds) * time.Millisecond
		}
	}
	ocspClient := &http.Client{
		Timeout:   timeout,
		Transport: snowflakeInsecureTransport,
	}
	ocspRes, ocspResBytes, ocspS := retryOCSP(
		ctx, ocspClient, http.NewRequest, u, headers, ocspReq, issuer, timeout)
	if ocspS.code != ocspSuccess {
		return ocspS
	}

	ret := validateOCSP(ocspRes)
	if !isValidOCSPStatus(ret.code) {
		return ret // return invalid
	}
	v := []interface{}{float64(time.Now().UTC().Unix()), base64.StdEncoding.EncodeToString(ocspResBytes)}
	ocspResponseCacheLock.Lock()
	ocspResponseCache[*encodedCertID] = v
	cacheUpdated = true
	ocspResponseCacheLock.Unlock()
	return ret
}

func isTestNoOCSPURL() bool {
	return strings.EqualFold(os.Getenv(ocspTestNoOCSPURLEnv), "true")
}

func isValidOCSPStatus(status ocspStatusCode) bool {
	return status == ocspStatusGood || status == ocspStatusRevoked || status == ocspStatusUnknown
}

// verifyPeerCertificate verifies all of certificate revocation status
func verifyPeerCertificate(ctx context.Context, verifiedChains [][]*x509.Certificate) (err error) {
	for i := 0; i < len(verifiedChains); i++ {
		// Certificate signed by Root CA. This should be one before the last in the Certificate Chain
		numberOfNoneRootCerts := len(verifiedChains[i]) - 1
		if !verifiedChains[i][numberOfNoneRootCerts].IsCA || string(verifiedChains[i][numberOfNoneRootCerts].RawIssuer) != string(verifiedChains[i][numberOfNoneRootCerts].RawSubject) {
			// Check if the last Non Root Cert is also a CA or is self signed.
			// if the last certificate is not, add it to the list
			rca := caRoot[string(verifiedChains[i][numberOfNoneRootCerts].RawIssuer)]
			if rca == nil {
				return fmt.Errorf("failed to find root CA. pkix.name: %v", verifiedChains[i][numberOfNoneRootCerts].Issuer)
			}
			verifiedChains[i] = append(verifiedChains[i], rca)
			numberOfNoneRootCerts++
		}
		results := getAllRevocationStatus(ctx, verifiedChains[i])
		if r := canEarlyExitForOCSP(results, numberOfNoneRootCerts); r != nil {
			return r.err
		}
	}

	ocspResponseCacheLock.Lock()
	if cacheUpdated {
		writeOCSPCacheFile()
	}
	cacheUpdated = false
	ocspResponseCacheLock.Unlock()
	return nil
}

func canEarlyExitForOCSP(results []*ocspStatus, chainSize int) *ocspStatus {
	msg := ""
	if atomic.LoadUint32((*uint32)(&ocspFailOpen)) == (uint32)(OCSPFailOpenFalse) {
		// Fail closed. any error is returned to stop connection
		for _, r := range results {
			if r.err != nil {
				return r
			}
		}
	} else {
		// Fail open and all results are valid.
		allValid := len(results) == chainSize
		for _, r := range results {
			if !isValidOCSPStatus(r.code) {
				allValid = false
				break
			}
		}
		for _, r := range results {
			if allValid && r.code == ocspStatusRevoked {
				return r
			}
			if r != nil && r.code != ocspStatusGood && r.err != nil {
				msg += "\n" + r.err.Error()
			}
		}
	}
	if len(msg) > 0 {
		logger.Warnf(
			"WARNING!!! Using fail-open to connect. Driver is connecting to an "+
				"HTTPS endpoint without OCSP based Certificate Revocation checking "+
				"as it could not obtain a valid OCSP Response to use from the CA OCSP "+
				"responder. Detail: %v", msg[1:])
	}
	return nil
}

func validateWithCacheForAllCertificates(verifiedChains []*x509.Certificate) bool {
	n := len(verifiedChains) - 1
	for j := 0; j < n; j++ {
		subject := verifiedChains[j]
		issuer := verifiedChains[j+1]
		status, _, _ := validateWithCache(subject, issuer)
		if !isValidOCSPStatus(status.code) {
			return false
		}
	}
	return true
}

func validateWithCache(subject, issuer *x509.Certificate) (*ocspStatus, []byte, *certIDKey) {
	ocspReq, err := ocsp.CreateRequest(subject, issuer, &ocsp.RequestOptions{})
	if err != nil {
		logger.Errorf("failed to create OCSP request from the certificates.\n")
		return &ocspStatus{
			code: ocspFailedComposeRequest,
			err:  errors.New("failed to create a OCSP request"),
		}, nil, nil
	}
	encodedCertID, ocspS := extractCertIDKeyFromRequest(ocspReq)
	if ocspS.code != ocspSuccess {
		logger.Errorf("failed to extract CertID from OCSP Request.\n")
		return &ocspStatus{
			code: ocspFailedComposeRequest,
			err:  errors.New("failed to extract cert ID Key"),
		}, ocspReq, nil
	}
	status := checkOCSPResponseCache(encodedCertID, subject, issuer)
	return status, ocspReq, encodedCertID
}

func downloadOCSPCacheServer() {
	if strings.EqualFold(os.Getenv(cacheServerEnabledEnv), "false") {
		return
	}
	ocspCacheServerURL := os.Getenv(cacheServerURLEnv)
	if ocspCacheServerURL == "" {
		ocspCacheServerURL = fmt.Sprintf("%v/%v", cacheServerURL, cacheFileBaseName)
	}
	u, err := url.Parse(ocspCacheServerURL)
	if err != nil {
		return
	}
	logger.Infof("downloading OCSP Cache from server %v", ocspCacheServerURL)
	timeoutStr := os.Getenv(ocspTestResponseCacheServerTimeoutEnv)
	timeout := defaultOCSPCacheServerTimeout
	if timeoutStr != "" {
		var timeoutInMilliseconds int
		timeoutInMilliseconds, err = strconv.Atoi(timeoutStr)
		if err == nil {
			timeout = time.Duration(timeoutInMilliseconds) * time.Millisecond
		}
	}
	ocspClient := &http.Client{
		Timeout:   timeout,
		Transport: snowflakeInsecureTransport,
	}
	ret, ocspStatus := checkOCSPCacheServer(context.TODO(), ocspClient, http.NewRequest, u, timeout)
	if ocspStatus.code != ocspSuccess {
		return
	}

	ocspResponseCacheLock.Lock()
	for k, cacheValue := range *ret {
		status := extractOCSPCacheResponseValueWithoutSubject(cacheValue)
		if !isValidOCSPStatus(status.code) {
			continue
		}
		cacheKey := encodeCertIDKey(k)
		ocspResponseCache[*cacheKey] = cacheValue
	}
	cacheUpdated = true
	ocspResponseCacheLock.Unlock()
}

func getAllRevocationStatus(ctx context.Context, verifiedChains []*x509.Certificate) []*ocspStatus {
	cached := validateWithCacheForAllCertificates(verifiedChains)
	if !cached {
		downloadOCSPCacheServer()
	}
	n := len(verifiedChains) - 1
	results := make([]*ocspStatus, n)
	for j := 0; j < n; j++ {
		results[j] = getRevocationStatus(ctx, verifiedChains[j], verifiedChains[j+1])
		if !isValidOCSPStatus(results[j].code) {
			return results
		}
	}
	return results
}

// verifyPeerCertificateSerial verifies the certificate revocation status in serial.
func verifyPeerCertificateSerial(_ [][]byte, verifiedChains [][]*x509.Certificate) (err error) {
	overrideCacheDir()
	return verifyPeerCertificate(context.TODO(), verifiedChains)
}

func overrideCacheDir() {
	if os.Getenv(cacheDirEnv) != "" {
		ocspResponseCacheLock.Lock()
		defer ocspResponseCacheLock.Unlock()
		createOCSPCacheDir()
	}
}

// initOCSPCache initializes OCSP Response cache file.
func initOCSPCache() {
	if strings.EqualFold(os.Getenv(cacheServerEnabledEnv), "false") {
		return
	}
	ocspResponseCache = make(map[certIDKey][]interface{})
	ocspResponseCacheLock = &sync.RWMutex{}

	logger.Infof("reading OCSP Response cache file. %v\n", cacheFileName)
	f, err := os.OpenFile(cacheFileName, os.O_CREATE|os.O_RDONLY, os.ModePerm)
	if err != nil {
		logger.Debugf("failed to open. Ignored. %v\n", err)
		return
	}
	defer f.Close()

	buf := make(map[string][]interface{})

	r := bufio.NewReader(f)
	dec := json.NewDecoder(r)
	for {
		if err = dec.Decode(&buf); err == io.EOF {
			break
		} else if err != nil {
			logger.Debugf("failed to read. Ignored. %v\n", err)
			return
		}
	}
	for k, cacheValue := range buf {
		status := extractOCSPCacheResponseValueWithoutSubject(cacheValue)
		if !isValidOCSPStatus(status.code) {
			continue
		}
		cacheKey := encodeCertIDKey(k)
		ocspResponseCache[*cacheKey] = cacheValue

	}
	cacheUpdated = false
}
func extractOCSPCacheResponseValueWithoutSubject(cacheValue []interface{}) *ocspStatus {
	return extractOCSPCacheResponseValue(cacheValue, nil, nil)
}

func extractOCSPCacheResponseValue(cacheValue []interface{}, subject, issuer *x509.Certificate) *ocspStatus {
	subjectName := "Unknown"
	if subject != nil {
		subjectName = subject.Subject.CommonName
	}

	curTime := time.Now()
	if len(cacheValue) != 2 {
		return &ocspStatus{
			code: ocspMissedCache,
			err:  fmt.Errorf("miss cache data. subject: %v", subjectName),
		}
	}
	if ts, ok := cacheValue[0].(float64); ok {
		currentTime := float64(curTime.UTC().Unix())
		if currentTime-ts >= cacheExpire {
			return &ocspStatus{
				code: ocspCacheExpired,
				err: fmt.Errorf("cache expired. current: %v, cache: %v",
					time.Unix(int64(currentTime), 0).UTC(), time.Unix(int64(ts), 0).UTC()),
			}
		}
	} else {
		return &ocspStatus{
			code: ocspFailedDecodeResponse,
			err:  errors.New("the first cache element is not float64"),
		}
	}
	var err error
	var r *ocsp.Response
	if s, ok := cacheValue[1].(string); ok {
		var b []byte
		b, err = base64.StdEncoding.DecodeString(s)
		if err != nil {
			return &ocspStatus{
				code: ocspFailedDecodeResponse,
				err:  fmt.Errorf("failed to decode OCSP Response value in a cache. subject: %v, err: %v", subjectName, err),
			}
		}
		// check the revocation status here
		r, err = ocsp.ParseResponse(b, issuer)
		if err != nil {
			logger.Warnf("the second cache element is not a valid OCSP Response. Ignored. subject: %v\n", subjectName)
			return &ocspStatus{
				code: ocspFailedParseResponse,
				err:  fmt.Errorf("failed to parse OCSP Respose. subject: %v, err: %v", subjectName, err),
			}
		}
	} else {
		return &ocspStatus{
			code: ocspFailedDecodeResponse,
			err:  errors.New("the second cache element is not string"),
		}

	}
	return validateOCSP(r)
}

// writeOCSPCacheFile writes a OCSP Response cache file. This is called if all revocation status is success.
// lock file is used to mitigate race condition with other process.
func writeOCSPCacheFile() {
	if strings.EqualFold(os.Getenv(cacheServerEnabledEnv), "false") {
		return
	}
	logger.Infof("writing OCSP Response cache file. %v\n", cacheFileName)
	cacheLockFileName := cacheFileName + ".lck"
	err := os.Mkdir(cacheLockFileName, 0600)
	switch {
	case os.IsExist(err):
		statinfo, err := os.Stat(cacheLockFileName)
		if err != nil {
			logger.Debugf("failed to write OCSP response cache file. file: %v, err: %v. ignored.\n", cacheFileName, err)
			return
		}
		if time.Since(statinfo.ModTime()) < 15*time.Minute {
			logger.Debugf("other process locks the cache file. %v. ignored.\n", cacheFileName)
			return
		}
		if err = os.Remove(cacheLockFileName); err != nil {
			logger.Debugf("failed to delete lock file. file: %v, err: %v. ignored.\n", cacheLockFileName, err)
			return
		}
		if err = os.Mkdir(cacheLockFileName, 0600); err != nil {
			logger.Debugf("failed to delete lock file. file: %v, err: %v. ignored.\n", cacheLockFileName, err)
			return
		}
	}
	defer os.RemoveAll(cacheLockFileName)

	buf := make(map[string][]interface{})
	for k, v := range ocspResponseCache {
		cacheKeyInBase64 := decodeCertIDKey(&k)
		buf[cacheKeyInBase64] = v
	}

	j, err := json.Marshal(buf)
	if err != nil {
		logger.Debugf("failed to convert OCSP Response cache to JSON. ignored.")
		return
	}
	if err = ioutil.WriteFile(cacheFileName, j, 0644); err != nil {
		logger.Debugf("failed to write OCSP Response cache. err: %v. ignored.\n", err)
	}
}

// readCACerts read a set of root CAs
func readCACerts() {
	raw := []byte(caRootPEM)
	certPool = x509.NewCertPool()
	caRoot = make(map[string]*x509.Certificate)
	var p *pem.Block
	for {
		p, raw = pem.Decode(raw)
		if p == nil {
			break
		}
		if p.Type != "CERTIFICATE" {
			continue
		}
		c, err := x509.ParseCertificate(p.Bytes)
		if err != nil {
			panic("failed to parse CA certificate.")
		}
		certPool.AddCert(c)
		caRoot[string(c.RawSubject)] = c
	}
}

// createOCSPCacheDir creates OCSP response cache directory and set the cache file name.
func createOCSPCacheDir() {
	if strings.EqualFold(os.Getenv(cacheServerEnabledEnv), "false") {
		logger.Info(`OCSP Cache Server disabled. All further access and use of
			OCSP Cache will be disabled for this OCSP Status Query`)
		return
	}
	cacheDir = os.Getenv(cacheDirEnv)
	if cacheDir == "" {
		cacheDir = os.Getenv("SNOWFLAKE_TEST_WORKSPACE")
	}
	if cacheDir == "" {
		switch runtime.GOOS {
		case "windows":
			cacheDir = filepath.Join(os.Getenv("USERPROFILE"), "AppData", "Local", "Snowflake", "Caches")
		case "darwin":
			home := os.Getenv("HOME")
			if home == "" {
				logger.Info("HOME is blank.")
			}
			cacheDir = filepath.Join(home, "Library", "Caches", "Snowflake")
		default:
			home := os.Getenv("HOME")
			if home == "" {
				logger.Info("HOME is blank")
			}
			cacheDir = filepath.Join(home, ".cache", "snowflake")
		}
	}

	if _, err := os.Stat(cacheDir); os.IsNotExist(err) {
		if err = os.MkdirAll(cacheDir, os.ModePerm); err != nil {
			logger.Debugf("failed to create cache directory. %v, err: %v. ignored\n", cacheDir, err)
		}
	}
	cacheFileName = filepath.Join(cacheDir, cacheFileBaseName)
	logger.Infof("reset OCSP cache file. %v", cacheFileName)
}

// deleteOCSPCacheFile deletes the OCSP response cache file
func deleteOCSPCacheFile() {
	os.Remove(cacheFileName)
}

func init() {
	readCACerts()
	createOCSPCacheDir()
	initOCSPCache()
}

// snowflakeInsecureTransport is the transport object that doesn't do certificate revocation check.
var snowflakeInsecureTransport = &http.Transport{
	MaxIdleConns:    10,
	IdleConnTimeout: 30 * time.Minute,
	Proxy:           http.ProxyFromEnvironment,
	DialContext: (&net.Dialer{
		Timeout:   30 * time.Second,
		KeepAlive: 30 * time.Second,
	}).DialContext,
}

// SnowflakeTransport includes the certificate revocation check with OCSP in sequential. By default, the driver uses
// this transport object.
var SnowflakeTransport = &http.Transport{
	TLSClientConfig: &tls.Config{
		RootCAs:               certPool,
		VerifyPeerCertificate: verifyPeerCertificateSerial,
	},
	MaxIdleConns:    10,
	IdleConnTimeout: 30 * time.Minute,
	Proxy:           http.ProxyFromEnvironment,
	DialContext: (&net.Dialer{
		Timeout:   30 * time.Second,
		KeepAlive: 30 * time.Second,
	}).DialContext,
}

// SnowflakeTransportTest includes the certificate revocation check in parallel
var SnowflakeTransportTest = SnowflakeTransport
