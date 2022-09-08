// Copyright (c) 2017-2022 Snowflake Computing Inc. All rights reserved.

package ocsp

import (
	"context"
	"crypto"
	"crypto/tls"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/asn1"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/vault/sdk/logical"
	"golang.org/x/crypto/ocsp"
	"io"
	"io/ioutil"
	"math/big"
	"net"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
	"time"
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

type Client struct {
	// caRoot includes the CA certificates.
	caRoot map[string]*x509.Certificate
	// certPOol includes the CA certificates.
	certPool              *x509.CertPool
	ocspResponseCache     map[certIDKey][]interface{}
	ocspResponseCacheLock sync.RWMutex
	// cacheUpdated is true if the memory cache is updated
	cacheUpdated bool
	logFactory   func() hclog.Logger
}

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

// copied from crypto/ocsp
var hashOIDs = map[crypto.Hash]asn1.ObjectIdentifier{
	crypto.SHA1:   asn1.ObjectIdentifier([]int{1, 3, 14, 3, 2, 26}),
	crypto.SHA256: asn1.ObjectIdentifier([]int{2, 16, 840, 1, 101, 3, 4, 2, 1}),
	crypto.SHA384: asn1.ObjectIdentifier([]int{2, 16, 840, 1, 101, 3, 4, 2, 2}),
	crypto.SHA512: asn1.ObjectIdentifier([]int{2, 16, 840, 1, 101, 3, 4, 2, 3}),
}

// copied from crypto/ocsp
func (c *Client) getOIDFromHashAlgorithm(target crypto.Hash) asn1.ObjectIdentifier {
	for hash, oid := range hashOIDs {
		if hash == target {
			return oid
		}
	}
	c.Logger().Error("no valid OID is found for the hash algorithm", "target", target)
	return nil
}

func (c *Client) getHashAlgorithmFromOID(target pkix.AlgorithmIdentifier) crypto.Hash {
	for hash, oid := range hashOIDs {
		if oid.Equal(target.Algorithm) {
			return hash
		}
	}
	c.Logger().Error("no valid hash algorithm is found for the oid. Falling back to SHA1", "target", target)
	return crypto.SHA1
}

// calcTolerableValidity returns the maximum validity buffer
func calcTolerableValidity(thisUpdate, nextUpdate time.Time) time.Duration {
	return durationMax(nextUpdate.Sub(thisUpdate)/tolerableValidityRatio, maxClockSkew)
}

func durationMax(a, b time.Duration) time.Duration {
	if a > b {
		return a
	}
	return b
}

func durationMin(a, b time.Duration) time.Duration {
	return durationMax(b, a)
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

func (c *Client) encodeCertIDKey(certIDKeyBase64 string) *certIDKey {
	r, err := base64.StdEncoding.DecodeString(certIDKeyBase64)
	if err != nil {
		return nil
	}
	var cid certID
	rest, err := asn1.Unmarshal(r, &cid)
	if err != nil {
		// error in parsing
		return nil
	}
	if len(rest) > 0 {
		// extra bytes to the end
		return nil
	}
	return &certIDKey{
		c.getHashAlgorithmFromOID(cid.HashAlgorithm),
		base64.StdEncoding.EncodeToString(cid.NameHash),
		base64.StdEncoding.EncodeToString(cid.IssuerKeyHash),
		cid.SerialNumber.String(),
	}
}

func (c *Client) decodeCertIDKey(k *certIDKey) string {
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
			Algorithm:  c.getOIDFromHashAlgorithm(k.HashAlgorithm),
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

func (c *Client) checkOCSPResponseCache(encodedCertID *certIDKey, subject, issuer *x509.Certificate) (*ocspStatus, error) {
	if strings.EqualFold(os.Getenv(cacheServerEnabledEnv), "false") {
		return &ocspStatus{code: ocspNoServer}, nil
	}
	c.ocspResponseCacheLock.RLock()
	gotValueFromCache := c.ocspResponseCache[*encodedCertID]
	c.ocspResponseCacheLock.RUnlock()

	status, err := c.extractOCSPCacheResponseValue(gotValueFromCache, subject, issuer)
	if err != nil {
		return nil, err
	}
	if !isValidOCSPStatus(status.code) {
		c.deleteOCSPCache(encodedCertID)
	}
	return status, err
}

func (c *Client) deleteOCSPCache(encodedCertID *certIDKey) {
	c.ocspResponseCacheLock.Lock()
	delete(c.ocspResponseCache, *encodedCertID)
	c.cacheUpdated = true
	c.ocspResponseCacheLock.Unlock()
}

func validateOCSP(ocspRes *ocsp.Response) (*ocspStatus, error) {
	curTime := time.Now()

	if ocspRes == nil {
		return nil, errors.New("OCSP Response is nil")
	}
	if isTestInvalidValidity() || !isInValidityRange(curTime, ocspRes.ThisUpdate, ocspRes.NextUpdate) {
		return nil, fmt.Errorf("invalid validity: producedAt: %v, thisUpdate: %v, nextUpdate: %v", ocspRes.ProducedAt, ocspRes.ThisUpdate, ocspRes.NextUpdate)
	}
	if isTestUnknownStatus() {
		ocspRes.Status = ocsp.Unknown
	}
	return returnOCSPStatus(ocspRes), nil
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
		}
	case ocsp.Unknown:
		return &ocspStatus{
			code: ocspStatusUnknown,
			err:  errors.New("OCSP status unknown."),
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

func (c *Client) checkOCSPCacheServer(
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
		c.Logger().Error("failed to get OCSP cache from OCSP Cache Server. ", "err", err)
		return nil, &ocspStatus{
			code: ocspFailedSubmit,
			err:  err,
		}
	}
	defer res.Body.Close()
	c.Logger().Debug("StatusCode from OCSP Cache Server", "statusCode", res.StatusCode)
	if res.StatusCode != http.StatusOK {
		return nil, &ocspStatus{
			code: ocspFailedResponse,
			err:  fmt.Errorf("HTTP code is not OK. %v: %v", res.StatusCode, res.Status),
		}
	}
	c.Logger().Debug("reading contents")

	dec := json.NewDecoder(res.Body)
	for {
		if err := dec.Decode(&respd); err == io.EOF {
			break
		} else if err != nil {
			c.Logger().Error("failed to decode OCSP cache.", "err", err)
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
func (c *Client) retryOCSP(
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
	c.Logger().Debug("StatusCode from OCSP Server:", "statusCode", res.StatusCode)
	if res.StatusCode != http.StatusOK {
		return ocspRes, ocspResBytes, &ocspStatus{
			code: ocspFailedResponse,
			err:  fmt.Errorf("HTTP code is not OK. %v: %v", res.StatusCode, res.Status),
		}
	}
	c.Logger().Debug("reading contents")
	ocspResBytes, err = ioutil.ReadAll(res.Body)
	if err != nil {
		return ocspRes, ocspResBytes, &ocspStatus{
			code: ocspFailedExtractResponse,
			err:  err,
		}
	}
	c.Logger().Debug("parsing OCSP response")
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
func (c *Client) getRevocationStatus(ctx context.Context, subject, issuer *x509.Certificate) (*ocspStatus, error) {
	c.Logger().Info("get-revocation-status", "subject", subject.Subject, "issuer", issuer.Subject)

	status, ocspReq, encodedCertID, err := c.ValidateWithCache(subject, issuer)
	if err != nil {
		return nil, err
	}
	if isValidOCSPStatus(status.code) {
		return status, nil
	}
	if ocspReq == nil || encodedCertID == nil {
		return status, nil
	}
	c.Logger().Info("cache missed")
	c.Logger().Info("OCSP: ", "server", subject.OCSPServer)
	if len(subject.OCSPServer) == 0 || isTestNoOCSPURL() {
		return nil, fmt.Errorf("no OCSP responder URL: subject: %v", subject.Subject)
	}
	ocspHost := subject.OCSPServer[0]
	u, err := url.Parse(ocspHost)
	if err != nil {
		return nil, err
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
	ocspRes, ocspResBytes, ocspS := c.retryOCSP(
		ctx, ocspClient, http.NewRequest, u, headers, ocspReq, issuer, timeout)
	if ocspS.code != ocspSuccess {
		return ocspS, nil
	}

	ret, err := validateOCSP(ocspRes)
	if err != nil {
		return nil, err
	}
	if !isValidOCSPStatus(ret.code) {
		return ret, nil // return invalid
	}
	v := []interface{}{float64(time.Now().UTC().Unix()), base64.StdEncoding.EncodeToString(ocspResBytes)}
	c.ocspResponseCacheLock.Lock()
	c.ocspResponseCache[*encodedCertID] = v
	c.cacheUpdated = true
	c.ocspResponseCacheLock.Unlock()
	return ret, nil
}

func isTestNoOCSPURL() bool {
	return strings.EqualFold(os.Getenv(ocspTestNoOCSPURLEnv), "true")
}

func isValidOCSPStatus(status ocspStatusCode) bool {
	return status == ocspStatusGood || status == ocspStatusRevoked || status == ocspStatusUnknown
}

// VerifyPeerCertificate verifies all of certificate revocation status
func (c *Client) VerifyPeerCertificate(ctx context.Context, verifiedChains [][]*x509.Certificate) (err error) {
	for i := 0; i < len(verifiedChains); i++ {
		// Certificate signed by Root CA. This should be one before the last in the Certificate Chain
		numberOfNoneRootCerts := len(verifiedChains[i]) - 1
		if !verifiedChains[i][numberOfNoneRootCerts].IsCA || string(verifiedChains[i][numberOfNoneRootCerts].RawIssuer) != string(verifiedChains[i][numberOfNoneRootCerts].RawSubject) {
			// Check if the last Non Root Cert is also a CA or is self signed.
			// if the last certificate is not, add it to the list
			rca := c.caRoot[string(verifiedChains[i][numberOfNoneRootCerts].RawIssuer)]
			if rca == nil {
				return fmt.Errorf("failed to find root CA. pkix.name: %v", verifiedChains[i][numberOfNoneRootCerts].Issuer)
			}
			verifiedChains[i] = append(verifiedChains[i], rca)
			numberOfNoneRootCerts++
		}
		results, err := c.getAllRevocationStatus(ctx, verifiedChains[i])
		if err != nil {
			return err
		}
		if r := c.canEarlyExitForOCSP(results, numberOfNoneRootCerts); r != nil {
			return r.err
		}
	}

	return nil
}

func (c *Client) canEarlyExitForOCSP(results []*ocspStatus, chainSize int) *ocspStatus {
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
				msg += "" + r.err.Error()
			}
		}
	}
	if len(msg) > 0 {
		c.Logger().Warn(
			"WARNING!!! Using fail-open to connect. Driver is connecting to an "+
				"HTTPS endpoint without OCSP based Certificate Revocation checking "+
				"as it could not obtain a valid OCSP Response to use from the CA OCSP "+
				"responder", "detail", msg[1:])
	}
	return nil
}

func (c *Client) ValidateWithCacheForAllCertificates(verifiedChains []*x509.Certificate) (bool, error) {
	n := len(verifiedChains) - 1
	for j := 0; j < n; j++ {
		subject := verifiedChains[j]
		issuer := verifiedChains[j+1]
		status, _, _, err := c.ValidateWithCache(subject, issuer)
		if err != nil {
			return false, err
		}
		if !isValidOCSPStatus(status.code) {
			return false, nil
		}
	}
	return true, nil
}

func (c *Client) ValidateWithCache(subject, issuer *x509.Certificate) (*ocspStatus, []byte, *certIDKey, error) {
	ocspReq, err := ocsp.CreateRequest(subject, issuer, &ocsp.RequestOptions{})
	if err != nil {
		return nil, nil, nil, fmt.Errorf("failed to create OCSP request from the certificates: %v", err)
	}
	encodedCertID, ocspS := extractCertIDKeyFromRequest(ocspReq)
	if ocspS.code != ocspSuccess {
		return nil, nil, nil, fmt.Errorf("failed to extract CertID from OCSP Request: %v", err)
	}
	status, err := c.checkOCSPResponseCache(encodedCertID, subject, issuer)
	if err != nil {
		return nil, nil, nil, err
	}
	return status, ocspReq, encodedCertID, nil
}

func (c *Client) getAllRevocationStatus(ctx context.Context, verifiedChains []*x509.Certificate) ([]*ocspStatus, error) {
	_, err := c.ValidateWithCacheForAllCertificates(verifiedChains)
	if err != nil {
		return nil, err
	}
	n := len(verifiedChains) - 1
	results := make([]*ocspStatus, n)
	for j := 0; j < n; j++ {
		results[j], err = c.getRevocationStatus(ctx, verifiedChains[j], verifiedChains[j+1])
		if err != nil {
			return nil, err
		}
		if !isValidOCSPStatus(results[j].code) {
			return results, nil
		}
	}
	return results, nil
}

// verifyPeerCertificateSerial verifies the certificate revocation status in serial.
func (c *Client) verifyPeerCertificateSerial(_ [][]byte, verifiedChains [][]*x509.Certificate) (err error) {
	return c.VerifyPeerCertificate(context.TODO(), verifiedChains)
}

/*
// initOCSPCache initializes OCSP Response cache file.
func (c *Client) initOCSPCache() {
	if strings.EqualFold(os.Getenv(cacheServerEnabledEnv), "false") {
		return
	}
	ocspResponseCache = make(map[certIDKey][]interface{})
	ocspResponseCacheLock = &sync.RWMutex{}

	c.Logger().Info("reading OCSP Response cache file.", "filename", cacheFileName)
	f, err := os.OpenFile(cacheFileName, os.O_CREATE|os.O_RDONLY, os.ModePerm)
	if err != nil {
		c.Logger().Debug("failed to open. Ignored.", "err", err)
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
			c.Logger().Debug("failed to read. Ignored.", "err", err)
			return
		}
	}
	for k, cacheValue := range buf {
		status := c.extractOCSPCacheResponseValueWithoutSubject(cacheValue)
		if !isValidOCSPStatus(status.code) {
			continue
		}
		cacheKey := c.encodeCertIDKey(k)
		ocspResponseCache[*cacheKey] = cacheValue

	}
	cacheUpdated = false
}*/
func (c *Client) extractOCSPCacheResponseValueWithoutSubject(cacheValue []interface{}) (*ocspStatus, error) {
	return c.extractOCSPCacheResponseValue(cacheValue, nil, nil)
}

func (c *Client) extractOCSPCacheResponseValue(cacheValue []interface{}, subject, issuer *x509.Certificate) (*ocspStatus, error) {
	subjectName := "Unknown"
	if subject != nil {
		subjectName = subject.Subject.CommonName
	}

	curTime := time.Now()
	if len(cacheValue) != 2 {
		return &ocspStatus{
			code: ocspMissedCache,
			err:  fmt.Errorf("miss cache data. subject: %v", subjectName),
		}, nil
	}
	if ts, ok := cacheValue[0].(float64); ok {
		currentTime := float64(curTime.UTC().Unix())
		if currentTime-ts >= cacheExpire {
			return &ocspStatus{
				code: ocspCacheExpired,
				err: fmt.Errorf("cache expired. current: %v, cache: %v",
					time.Unix(int64(currentTime), 0).UTC(), time.Unix(int64(ts), 0).UTC()),
			}, nil
		}
	} else {
		return nil, errors.New("the first cache element is not float64")
	}
	var err error
	var r *ocsp.Response
	if s, ok := cacheValue[1].(string); ok {
		var b []byte
		b, err = base64.StdEncoding.DecodeString(s)
		if err != nil {
			return nil, fmt.Errorf("failed to decode OCSP Response value in a cache. subject: %v, err: %v", subjectName, err)

		}
		// check the revocation status here
		r, err = ocsp.ParseResponse(b, issuer)
		if err != nil {
			c.Logger().Warn("the second cache element is not a valid OCSP Response. Ignored.", "subject", subjectName)
			return nil, fmt.Errorf("failed to parse OCSP Respose. subject: %v, err: %v", subjectName, err)
		}
	} else {
		return nil, errors.New("the second cache element is not string")

	}
	return validateOCSP(r)
}

const storageValueKey = "backingStore"

// writeOCSPCacheFile writes a OCSP Response cache file. This is called if all revocation status is success.
// lock file is used to mitigate race condition with other process.
func (c *Client) writeOCSPCacheFile(ctx context.Context, storage logical.Storage) error {
	c.Logger().Debug("writing OCSP Response cache file")

	buf := make(map[string][]interface{})
	for k, v := range c.ocspResponseCache {
		cacheKeyInBase64 := c.decodeCertIDKey(&k)
		buf[cacheKeyInBase64] = v
	}

	j, err := json.Marshal(buf)
	if err != nil {
		return errors.New("failed to convert OCSP Response cache to JSON")
	}

	entry := logical.StorageEntry{
		Key:   "ocsp_cache",
		Value: j,
	}
	return storage.Put(ctx, &entry)
}

/*
// readCACerts read a set of root CAs
func (c *Client) readCACerts() {
	raw := []byte(c.caRootPEM)
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
		c.Logger().Info(`OCSP Cache Server disabled. All further access and use of
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
				c.Logger().Info("HOME is blank.")
			}
			cacheDir = filepath.Join(home, "Library", "Caches", "Snowflake")
		default:
			home := os.Getenv("HOME")
			if home == "" {
				c.Logger().Info("HOME is blank")
			}
			cacheDir = filepath.Join(home, ".cache", "snowflake")
		}
	}

	if _, err := os.Stat(cacheDir); os.IsNotExist(err) {
		if err = os.MkdirAll(cacheDir, os.ModePerm); err != nil {
			c.Logger().Debugf("failed to create cache directory. %v, err: %v. ignored", cacheDir, err)
		}
	}
	cacheFileName = filepath.Join(cacheDir, cacheFileBaseName)
	c.Logger().Infof("reset OCSP cache file. %v", cacheFileName)
}
*/
func New(logFactory func() hclog.Logger) *Client {
	c := Client{
		caRoot:            make(map[string]*x509.Certificate),
		ocspResponseCache: make(map[certIDKey][]interface{}),
		logFactory:        logFactory,
	}

	return &c
}

func (c *Client) Logger() hclog.Logger {
	return c.logFactory()
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
func (c *Client) NewTransport() *http.Transport {
	return &http.Transport{
		TLSClientConfig: &tls.Config{
			RootCAs:               c.certPool,
			VerifyPeerCertificate: c.verifyPeerCertificateSerial,
		},
		MaxIdleConns:    10,
		IdleConnTimeout: 30 * time.Minute,
		Proxy:           http.ProxyFromEnvironment,
		DialContext: (&net.Dialer{
			Timeout:   30 * time.Second,
			KeepAlive: 30 * time.Second,
		}).DialContext,
	}
}

func (c *Client) WriteCache(ctx context.Context, storage logical.Storage) error {
	c.ocspResponseCacheLock.Lock()
	if c.cacheUpdated {
		defer c.ocspResponseCacheLock.Unlock()
		if c.cacheUpdated {
			return c.writeOCSPCacheFile(ctx, storage)
		}
		c.cacheUpdated = false
	}
	return nil
}

/*
                                 Apache License
                           Version 2.0, January 2004
                        http://www.apache.org/licenses/

   TERMS AND CONDITIONS FOR USE, REPRODUCTION, AND DISTRIBUTION

   1. Definitions.

      "License" shall mean the terms and conditions for use, reproduction,
      and distribution as defined by Sections 1 through 9 of this document.

      "Licensor" shall mean the copyright owner or entity authorized by
      the copyright owner that is granting the License.

      "Legal Entity" shall mean the union of the acting entity and all
      other entities that control, are controlled by, or are under common
      control with that entity. For the purposes of this definition,
      "control" means (i) the power, direct or indirect, to cause the
      direction or management of such entity, whether by contract or
      otherwise, or (ii) ownership of fifty percent (50%) or more of the
      outstanding shares, or (iii) beneficial ownership of such entity.

      "You" (or "Your") shall mean an individual or Legal Entity
      exercising permissions granted by this License.

      "Source" form shall mean the preferred form for making modifications,
      including but not limited to software source code, documentation
      source, and configuration files.

      "Object" form shall mean any form resulting from mechanical
      transformation or translation of a Source form, including but
      not limited to compiled object code, generated documentation,
      and conversions to other media types.

      "Work" shall mean the work of authorship, whether in Source or
      Object form, made available under the License, as indicated by a
      copyright notice that is included in or attached to the work
      (an example is provided in the Appendix below).

      "Derivative Works" shall mean any work, whether in Source or Object
      form, that is based on (or derived from) the Work and for which the
      editorial revisions, annotations, elaborations, or other modifications
      represent, as a whole, an original work of authorship. For the purposes
      of this License, Derivative Works shall not include works that remain
      separable from, or merely link (or bind by name) to the interfaces of,
      the Work and Derivative Works thereof.

      "Contribution" shall mean any work of authorship, including
      the original version of the Work and any modifications or additions
      to that Work or Derivative Works thereof, that is intentionally
      submitted to Licensor for inclusion in the Work by the copyright owner
      or by an individual or Legal Entity authorized to submit on behalf of
      the copyright owner. For the purposes of this definition, "submitted"
      means any form of electronic, verbal, or written communication sent
      to the Licensor or its representatives, including but not limited to
      communication on electronic mailing lists, source code control systems,
      and issue tracking systems that are managed by, or on behalf of, the
      Licensor for the purpose of discussing and improving the Work, but
      excluding communication that is conspicuously marked or otherwise
      designated in writing by the copyright owner as "Not a Contribution."

      "Contributor" shall mean Licensor and any individual or Legal Entity
      on behalf of whom a Contribution has been received by Licensor and
      subsequently incorporated within the Work.

   2. Grant of Copyright License. Subject to the terms and conditions of
      this License, each Contributor hereby grants to You a perpetual,
      worldwide, non-exclusive, no-charge, royalty-free, irrevocable
      copyright license to reproduce, prepare Derivative Works of,
      publicly display, publicly perform, sublicense, and distribute the
      Work and such Derivative Works in Source or Object form.

   3. Grant of Patent License. Subject to the terms and conditions of
      this License, each Contributor hereby grants to You a perpetual,
      worldwide, non-exclusive, no-charge, royalty-free, irrevocable
      (except as stated in this section) patent license to make, have made,
      use, offer to sell, sell, import, and otherwise transfer the Work,
      where such license applies only to those patent claims licensable
      by such Contributor that are necessarily infringed by their
      Contribution(s) alone or by combination of their Contribution(s)
      with the Work to which such Contribution(s) was submitted. If You
      institute patent litigation against any entity (including a
      cross-claim or counterclaim in a lawsuit) alleging that the Work
      or a Contribution incorporated within the Work constitutes direct
      or contributory patent infringement, then any patent licenses
      granted to You under this License for that Work shall terminate
      as of the date such litigation is filed.

   4. Redistribution. You may reproduce and distribute copies of the
      Work or Derivative Works thereof in any medium, with or without
      modifications, and in Source or Object form, provided that You
      meet the following conditions:

      (a) You must give any other recipients of the Work or
          Derivative Works a copy of this License; and

      (b) You must cause any modified files to carry prominent notices
          stating that You changed the files; and

      (c) You must retain, in the Source form of any Derivative Works
          that You distribute, all copyright, patent, trademark, and
          attribution notices from the Source form of the Work,
          excluding those notices that do not pertain to any part of
          the Derivative Works; and

      (d) If the Work includes a "NOTICE" text file as part of its
          distribution, then any Derivative Works that You distribute must
          include a readable copy of the attribution notices contained
          within such NOTICE file, excluding those notices that do not
          pertain to any part of the Derivative Works, in at least one
          of the following places: within a NOTICE text file distributed
          as part of the Derivative Works; within the Source form or
          documentation, if provided along with the Derivative Works; or,
          within a display generated by the Derivative Works, if and
          wherever such third-party notices normally appear. The contents
          of the NOTICE file are for informational purposes only and
          do not modify the License. You may add Your own attribution
          notices within Derivative Works that You distribute, alongside
          or as an addendum to the NOTICE text from the Work, provided
          that such additional attribution notices cannot be construed
          as modifying the License.

      You may add Your own copyright statement to Your modifications and
      may provide additional or different license terms and conditions
      for use, reproduction, or distribution of Your modifications, or
      for any such Derivative Works as a whole, provided Your use,
      reproduction, and distribution of the Work otherwise complies with
      the conditions stated in this License.

   5. Submission of Contributions. Unless You explicitly state otherwise,
      any Contribution intentionally submitted for inclusion in the Work
      by You to the Licensor shall be under the terms and conditions of
      this License, without any additional terms or conditions.
      Notwithstanding the above, nothing herein shall supersede or modify
      the terms of any separate license agreement you may have executed
      with Licensor regarding such Contributions.

   6. Trademarks. This License does not grant permission to use the trade
      names, trademarks, service marks, or product names of the Licensor,
      except as required for reasonable and customary use in describing the
      origin of the Work and reproducing the content of the NOTICE file.

   7. Disclaimer of Warranty. Unless required by applicable law or
      agreed to in writing, Licensor provides the Work (and each
      Contributor provides its Contributions) on an "AS IS" BASIS,
      WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or
      implied, including, without limitation, any warranties or conditions
      of TITLE, NON-INFRINGEMENT, MERCHANTABILITY, or FITNESS FOR A
      PARTICULAR PURPOSE. You are solely responsible for determining the
      appropriateness of using or redistributing the Work and assume any
      risks associated with Your exercise of permissions under this License.

   8. Limitation of Liability. In no event and under no legal theory,
      whether in tort (including negligence), contract, or otherwise,
      unless required by applicable law (such as deliberate and grossly
      negligent acts) or agreed to in writing, shall any Contributor be
      liable to You for damages, including any direct, indirect, special,
      incidental, or consequential damages of any character arising as a
      result of this License or out of the use or inability to use the
      Work (including but not limited to damages for loss of goodwill,
      work stoppage, computer failure or malfunction, or any and all
      other commercial damages or losses), even if such Contributor
      has been advised of the possibility of such damages.

   9. Accepting Warranty or Additional Liability. While redistributing
      the Work or Derivative Works thereof, You may choose to offer,
      and charge a fee for, acceptance of support, warranty, indemnity,
      or other liability obligations and/or rights consistent with this
      License. However, in accepting such obligations, You may act only
      on Your own behalf and on Your sole responsibility, not on behalf
      of any other Contributor, and only if You agree to indemnify,
      defend, and hold each Contributor harmless for any liability
      incurred by, or claims asserted against, such Contributor by reason
      of your accepting any such warranty or additional liability.

   END OF TERMS AND CONDITIONS

   APPENDIX: How to apply the Apache License to your work.

      To apply the Apache License to your work, attach the following
      boilerplate notice, with the fields enclosed by brackets "{}"
      replaced with your own identifying information. (Don't include
      the brackets!)  The text should be enclosed in the appropriate
      comment syntax for the file format. We also recommend that a
      file or class name and description of purpose be included on the
      same "printed page" as the copyright notice for easier
      identification within third-party archives.

   Copyright (c) 2017-2022 Snowflake Computing Inc. All rights reserved.

   Licensed under the Apache License, Version 2.0 (the "License");
   you may not use this file except in compliance with the License.
   You may obtain a copy of the License at

       http://www.apache.org/licenses/LICENSE-2.0

   Unless required by applicable law or agreed to in writing, software
   distributed under the License is distributed on an "AS IS" BASIS,
   WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
   See the License for the specific language governing permissions and
   limitations under the License.
*/
