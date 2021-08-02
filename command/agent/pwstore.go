package agent

import (
	"crypto/rand"
	"crypto/subtle"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/hashicorp/vault/api"
	"github.com/hashicorp/vault/command/agent/cache"
	"github.com/hashicorp/vault/sdk/helper/consts"
	"golang.org/x/crypto/argon2"
	"io"
	"net/http"
	"regexp"
)

type PasswordStore struct {
	KVPath string `hcl:"kv_path"`
}

type passwordStoreHandler struct {
	config  *PasswordStore
	proxier cache.Proxier
	client  *api.Client
}

var pat *regexp.Regexp

func init() {
	var err error
	pat, err = regexp.Compile("/v1/password_store/(?P<name>[^/]+)/(?P<operation>.*)")
	if err != nil {
		fmt.Printf("error compiling regexp: %e", err)
	}
}

func NewPasswordStore(config *PasswordStore, proxier cache.Proxier, cli *api.Client) *passwordStoreHandler {
	return &passwordStoreHandler{
		config:  config,
		proxier: proxier,
		client:  cli,
	}
}

type requestBody struct {
	Password string
}

type kvResponse struct {
	Data struct {
		Value string
		Salt  string
	}
}

func (p *passwordStoreHandler) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	parts := pat.FindStringSubmatch(request.URL.Path)
	if parts == nil {
		writer.WriteHeader(http.StatusBadRequest)
		return
	}
	token := request.Header.Get(consts.AuthHeaderName)
	dec := json.NewDecoder(request.Body)
	var pw requestBody
	err := dec.Decode(&pw)
	if err != nil {
		fmt.Printf("error decoding body: %e", err)
		return
	}

	key := parts[1]
	operation := parts[2]
	switch operation {
	case "check":
		req := p.client.NewRequest(http.MethodGet, "/v1/"+p.config.KVPath+"/"+key)
		req.Headers[consts.AuthHeaderName] = []string{token}
		resp, err := p.client.RawRequest(req)
		if err != nil {
			fmt.Printf("error making request: %e", err)
		}
		var val kvResponse
		err = resp.DecodeJSON(&val)
		if err != nil {
			fmt.Printf("error making request: %e", err)
		}

		salt, err := base64.StdEncoding.DecodeString(val.Data.Salt)
		if err != nil {
			fmt.Printf("error making request: %e", err)
		}
		hashed := argonHash(salt, []byte(pw.Password))
		stored, err := base64.StdEncoding.DecodeString(val.Data.Value)
		if err != nil {
			fmt.Printf("error making request: %e", err)
		}
		if subtle.ConstantTimeCompare(hashed, stored) != 1 {
			writer.WriteHeader(http.StatusBadRequest)
		} else {
			writer.WriteHeader(http.StatusOK)
		}
	case "set":
		salt := make([]byte, 16)
		_, err := io.ReadFull(rand.Reader, salt)
		if err != nil {
			fmt.Printf("error making request: %e\n", err)
		}
		hashed := argonHash(salt, []byte(pw.Password))
		saltStr := base64.StdEncoding.EncodeToString(salt)
		hashStr := base64.StdEncoding.EncodeToString(hashed)
		req := p.client.NewRequest(http.MethodPost, "/v1/"+p.config.KVPath+"/"+key)
		req.Headers[consts.AuthHeaderName] = []string{token}
		req.SetJSONBody(map[string]interface{}{"value": hashStr, "salt": saltStr})
		resp, err := p.client.RawRequest(req)
		if err != nil {
			fmt.Printf("error making request: %e\n", err)
		}
		writer.WriteHeader(resp.StatusCode)
	default:
		writer.WriteHeader(http.StatusBadRequest)
	}
}

func argonHash(salt, password []byte) []byte {
	return argon2.IDKey([]byte(password), salt, 1, 64*1024, 4, 28)
}
