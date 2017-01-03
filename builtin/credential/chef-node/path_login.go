package chefNode

import (
	"bytes"
	"crypto"
	"crypto/rsa"
	"crypto/sha1"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"fmt"
	"math"
	"net/http"
	"sort"
	"strings"

	"net/url"

	"time"

	"io/ioutil"

	"encoding/json"

	"github.com/hashicorp/vault/helper/policyutil"
	"github.com/hashicorp/vault/helper/strutil"
	"github.com/hashicorp/vault/logical"
	"github.com/hashicorp/vault/logical/framework"
)

func pathLogin(b *backend) *framework.Path {
	return &framework.Path{
		Pattern: "login",
		Fields:  map[string]*framework.FieldSchema{},
		Callbacks: map[logical.Operation]framework.OperationFunc{
			logical.UpdateOperation: b.pathLogin,
		},
	}
}

func (b *backend) pathLogin(req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	key, err := b.retrievePubKey(req)
	if err != nil {
		return nil, err
	}

	auth := authenticate(req, key)
	if !auth {
		return logical.ErrorResponse("Couldn't authenticate client"), nil
	}

	allowedSkew := time.Minute * 5
	now := time.Now().UTC()
	headerTime, err := time.Parse(time.RFC3339, req.Connection.Header.Get("X-Ops-Timestamp"))
	if err != nil {
		return nil, err
	}

	if math.Abs(float64(now.Sub(headerTime))) > float64(allowedSkew) {
		return nil, fmt.Errorf("clock skew is too great for request")
	}

	policies, err := b.getNodePolicies(req)
	if err != nil {
		return nil, err
	}

	return &logical.Response{
		Auth: &logical.Auth{
			Policies:    policies,
			DisplayName: req.Connection.Header.Get("X-Ops-Userid"),
			LeaseOptions: logical.LeaseOptions{
				Renewable: true,
			},
			InternalData: map[string]interface{}{
				"headers": req.Connection.Header,
			},
		},
	}, nil
}

func (b *backend) pathLoginRenew(req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
	if req.Auth == nil {
		return nil, fmt.Errorf("request auth was nil")
	}

	headers := req.Auth.InternalData["headers"].(map[string]interface{})
	parsedHeader := make(http.Header)
	for k, v := range headers {
		headerStrings := make([]string, len(v.([]interface{})))
		for j, w := range v.([]interface{}) {
			headerStrings[j] = w.(string)
		}
		parsedHeader[k] = headerStrings
	}

	req.Connection.Header = parsedHeader

	key, err := b.retrievePubKey(req)
	if err != nil {
		return nil, err
	}

	auth := authenticate(req, key)
	if !auth {
		return nil, fmt.Errorf("couldn't authenticate renew request")
	}

	policies, err := b.getNodePolicies(req)
	if err != nil {
		return nil, fmt.Errorf("coulnd't retrieve current policy list")
	}

	if !policyutil.EquivalentPolicies(policies, req.Auth.Policies) {
		return nil, fmt.Errorf("policies have changed, not renewing")
	}

	return framework.LeaseExtend(0, 0, b.System())(req, d)
}

func constructAuthorization(h http.Header) string {
	authHeaders := make(map[string]string)
	var keys []string
	var ret bytes.Buffer

	for k, v := range h {
		if strings.HasPrefix(k, "X-Ops-Authorization-") {
			authHeaders[k] = v[0]
			keys = append(keys, k)
		}
	}
	sort.Strings(keys)
	for _, v := range keys {
		ret.WriteString(authHeaders[v])
	}
	return ret.String()
}

func (b *backend) getNodePolicies(req *logical.Request) ([]string, error) {
	nodeInfo, err := b.retrieveNodeInfo(req)
	if err != nil {
		return nil, err
	}

	var envPols []string
	envEntry, err := b.Environment(req.Storage, nodeInfo.ChefEnv)
	if err != nil {
		return nil, err
	}
	if envEntry != nil {
		envPols = envEntry.Policies
	}

	var rolePols []string
	roles := nodeInfo.AutoAttrs.Roles
	for _, v := range roles {
		roleEntry, err := b.Role(req.Storage, v)
		if err != nil {
			return nil, err
		}
		if roleEntry != nil {
			rolePols = append(rolePols, roleEntry.Policies...)
		}
	}

	var tagPols []string
	tags := nodeInfo.NormalAttrs.Tags
	for _, v := range tags {
		tagEntry, err := b.Tag(req.Storage, v)
		if err != nil {
			return nil, err
		}
		if tagEntry != nil {
			tagPols = append(tagPols, tagEntry.Policies...)
		}
	}

	var allPol []string
	allPol = append(allPol, envPols...)
	allPol = append(allPol, rolePols...)
	allPol = append(allPol, tagPols...)
	allPol = strutil.RemoveDuplicates(allPol)

	return allPol, nil
}

func (b *backend) retrieveNodeInfo(req *logical.Request) (*nodeResponse, error) {
	targetName := req.Connection.Header.Get("X-Ops-Userid")
	if targetName == "" {
		return nil, fmt.Errorf("Couldn't find client id to lookup public key")
	}

	config, err := b.Config(req.Storage)
	if err != nil {
		return nil, err
	}

	nodeURL, err := url.Parse(config.BaseURL + "/nodes/" + targetName)
	if err != nil {
		return nil, err
	}

	headers, err := authHeaders(config, nodeURL, "GET")
	if err != nil {
		return nil, err
	}

	nodeReq, err := http.NewRequest("GET", nodeURL.String(), nil)
	if err != nil {
		return nil, err
	}

	nodeReq.Header = headers
	client := &http.Client{}
	resp, err := client.Do(nodeReq)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var jb nodeResponse
	if err := json.Unmarshal(body, &jb); err != nil {
		return nil, err
	}

	return &jb, nil
}

func (b *backend) retrievePubKey(req *logical.Request) (*rsa.PublicKey, error) {
	targetName := req.Connection.Header.Get("X-Ops-Userid")
	if targetName == "" {
		return nil, fmt.Errorf("Couldn't find client id to lookup public key")
	}
	config, err := b.Config(req.Storage)
	if err != nil {
		return nil, err
	}

	keyURL, err := url.Parse(config.BaseURL + "/clients/" + targetName)
	if err != nil {
		return nil, err
	}

	headers, err := authHeaders(config, keyURL, "GET")
	if err != nil {
		return nil, err
	}

	clientReq, err := http.NewRequest("GET", keyURL.String(), nil)
	if err != nil {
		return nil, err
	}

	clientReq.Header = headers
	client := &http.Client{}
	resp, err := client.Do(clientReq)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	var jb clientResponse
	if err := json.Unmarshal(body, &jb); err != nil {
		return nil, err
	}
	key, err := parsePublicKey(jb.ClientKey)
	if err != nil {
		return nil, err
	}
	return key, nil
}

func authenticate(req *logical.Request, key *rsa.PublicKey) bool {
	h := req.Connection.Header
	sig, _ := base64.StdEncoding.DecodeString(constructAuthorization(h))
	hashedPath := sha1.Sum([]byte(req.MountPoint + req.Path))
	headers := []string{
		"Method:" + h.Get("Method"),
		"Hashed Path:" + base64.StdEncoding.EncodeToString(hashedPath[:]),
		"X-Ops-Content-Hash:" + h.Get("X-Ops-Content-Hash"),
		"X-Ops-Timestamp:" + h.Get("X-Ops-Timestamp"),
		"X-Ops-UserId:" + h.Get("X-Ops-Userid"),
	}
	headerString := strings.Join(headers, "\n")
	err := rsa.VerifyPKCS1v15(key, crypto.Hash(0), []byte(headerString), sig)
	if err != nil {
		return false
	}
	return true
}

func authHeaders(conf *config, url *url.URL, method string) (http.Header, error) {
	hashedPath := sha1.Sum([]byte(url.EscapedPath()))
	// So far nothing we do requires a body
	bodyHash := sha1.Sum([]byte(""))
	ts := time.Now().UTC().Format(time.RFC3339)
	headers := []string{
		"Method:" + method,
		"Hashed Path:" + base64.StdEncoding.EncodeToString(hashedPath[:]),
		"X-Ops-Content-Hash:" + base64.StdEncoding.EncodeToString(bodyHash[:]),
		"X-Ops-Timestamp:" + ts,
		"X-Ops-UserId:" + conf.ClientName,
	}

	headerString := strings.Join(headers, "\n")
	key, err := parsePrivateKey(conf.ClientKey)
	if err != nil {
		return nil, err
	}

	sig, err := rsa.SignPKCS1v15(nil, key, crypto.Hash(0), []byte(headerString))
	if err != nil {
		return nil, err
	}
	splitSig := splitOn60(base64.StdEncoding.EncodeToString(sig))
	ret := make(http.Header)
	for i := range splitSig {
		ret.Set(fmt.Sprintf("X-Ops-Authorization-%d", i+1), splitSig[i])
	}
	ret.Set("X-Ops-Sign", "algorithm=sha1;version=1.0;")
	ret.Set("Method", method)
	ret.Set("X-Ops-Timestamp", ts)
	ret.Set("X-Ops-Content-Hash", base64.StdEncoding.EncodeToString(bodyHash[:]))
	ret.Set("X-Ops-Userid", conf.ClientName)
	ret.Set("Accept", "application/json")
	ret.Set("X-Chef-Version", "12.0.0")
	ret.Set("host", url.Host)

	return ret, nil
}

func splitOn60(toSplit string) []string {
	size := int(math.Ceil(float64(len(toSplit)) / 60.0))
	sl := make([]string, size)
	for i := 0; i < size-1; i++ {
		sl[i] = toSplit[(i * 60) : (i*60)+60]
	}
	sl[size-1] = toSplit[(size-1)*60:]
	return sl
}

func parsePrivateKey(key string) (*rsa.PrivateKey, error) {
	block, _ := pem.Decode([]byte(key))
	if block == nil {
		return nil, fmt.Errorf("Couldn't parse PEM data")
	}
	privkey, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return nil, err
	}
	return privkey, nil
}

func parsePublicKey(key string) (*rsa.PublicKey, error) {
	block, _ := pem.Decode([]byte(key))
	if block == nil {
		return nil, fmt.Errorf("Couldn't parse PEM data")
	}
	pubkey, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, err
	}
	return pubkey.(*rsa.PublicKey), nil
}

type clientResponse struct {
	ClientKey string `json:"public_key"`
}

type nodeResponse struct {
	ChefEnv string `json:"chef_environment"`

	AutoAttrs struct {
		Roles []string
	} `json:"automatic"`

	NormalAttrs struct {
		Tags []string
	} `json:"normal"`
}
