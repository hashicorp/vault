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
		Fields: map[string]*framework.FieldSchema{
			"signature_version": &framework.FieldSchema{
				Type:        framework.TypeString,
				Description: "X-Ops-Sign",
			},
			"client_name": &framework.FieldSchema{
				Type:        framework.TypeString,
				Description: "X-Ops-UserId",
			},
			"timestamp": &framework.FieldSchema{
				Type:        framework.TypeString,
				Description: "X-Ops-Timestamp",
			},
			"signature": &framework.FieldSchema{
				Type:        framework.TypeString,
				Description: "X-Ops-Authorization-* concatinated together",
			},
		},
		Callbacks: map[logical.Operation]framework.OperationFunc{
			logical.UpdateOperation: b.pathLogin,
		},
	}
}

func (b *backend) pathLogin(req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	client := data.Get("client_name").(string)
	ts := data.Get("timestamp").(string)
	sig := data.Get("signature").(string)
	sigVer := data.Get("signature_version").(string)

	key, err := b.retrievePubKey(req, client)
	if err != nil {
		return nil, err
	}

	auth := authenticate(client, ts, sig, sigVer, key, req.MountPoint+req.Path)
	if !auth {
		return logical.ErrorResponse("Couldn't authenticate client"), nil
	}

	allowedSkew := time.Minute * 5
	now := time.Now().UTC()
	headerTime, err := time.Parse(time.RFC3339, data.Get("timestamp").(string))
	if err != nil {
		return nil, err
	}

	if math.Abs(float64(now.Sub(headerTime))) > float64(allowedSkew) {
		return nil, fmt.Errorf("clock skew is too great for request")
	}

	policies, err := b.getNodePolicies(req, client)
	if err != nil {
		return nil, err
	}

	return &logical.Response{
		Auth: &logical.Auth{
			Policies:    policies,
			DisplayName: client,
			LeaseOptions: logical.LeaseOptions{
				Renewable: true,
			},
			InternalData: map[string]interface{}{
				"request_path":      req.MountPoint + req.Path,
				"signature_version": data.Get("signature_version"),
				"signature":         data.Get("signature"),
				"client_name":       data.Get("client_name"),
				"timestamp":         data.Get("timestamp"),
			},
		},
	}, nil
}

func (b *backend) pathLoginRenew(req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
	if req.Auth == nil {
		return nil, fmt.Errorf("request auth was nil")
	}

	reqPath := req.Auth.InternalData["request_path"].(string)
	sig := req.Auth.InternalData["signature"].(string)
	sigVer := req.Auth.InternalData["signature_version"].(string)
	client := req.Auth.InternalData["client_name"].(string)
	ts := req.Auth.InternalData["timestamp"].(string)

	key, err := b.retrievePubKey(req, client)
	if err != nil {
		return nil, err
	}

	auth := authenticate(client, ts, sig, sigVer, key, reqPath)
	if !auth {
		return nil, fmt.Errorf("couldn't authenticate renew request")
	}

	policies, err := b.getNodePolicies(req, client)
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

func (b *backend) getNodePolicies(req *logical.Request, node string) ([]string, error) {
	nodeInfo, err := b.retrieveNodeInfo(req, node)
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

func (b *backend) retrieveNodeInfo(req *logical.Request, targetName string) (*nodeResponse, error) {
	config, err := b.Config(req.Storage)
	if err != nil {
		return nil, err
	}

	nodeURL, err := url.Parse(config.BaseURL + "/nodes/" + targetName)
	if err != nil {
		return nil, err
	}

	headers, err := authHeaders(config, nodeURL, "GET", true)
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

func (b *backend) retrievePubKey(req *logical.Request, targetName string) (*rsa.PublicKey, error) {
	config, err := b.Config(req.Storage)
	if err != nil {
		return nil, err
	}

	keyURL, err := url.Parse(config.BaseURL + "/clients/" + targetName)
	if err != nil {
		return nil, err
	}

	headers, err := authHeaders(config, keyURL, "GET", true)
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

func authenticate(client string, ts string, sig string, sigVer string, key *rsa.PublicKey, path string) bool {
	bodyHash := sha1.Sum([]byte(""))
	hashedPath := sha1.Sum([]byte(path))
	headers := []string{
		"Method:POST",
		"Hashed Path:" + base64.StdEncoding.EncodeToString(hashedPath[:]),
		"X-Ops-Content-Hash:" + base64.StdEncoding.EncodeToString(bodyHash[:]),
		"X-Ops-Timestamp:" + ts,
		"X-Ops-UserId:" + client,
	}
	headerString := strings.Join(headers, "\n")
	decSig, err := base64.StdEncoding.DecodeString(sig)
	if err != nil {
		return false
	}
	err = rsa.VerifyPKCS1v15(key, crypto.Hash(0), []byte(headerString), decSig)
	if err != nil {
		return false
	}
	return true
}

func authHeaders(conf *config, url *url.URL, method string, split bool) (http.Header, error) {
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
	ret := make(http.Header)
	if split {
		splitSig := splitOn60(base64.StdEncoding.EncodeToString(sig))
		for i := range splitSig {
			ret.Set(fmt.Sprintf("X-Ops-Authorization-%d", i+1), splitSig[i])
		}
	} else {
		ret.Set("X-Ops-Authorization", base64.StdEncoding.EncodeToString(sig))
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
