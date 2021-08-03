package agent

import (
	"context"
	"crypto/rand"
	"crypto/subtle"
	"encoding/base64"
	"fmt"
	"github.com/hashicorp/vault/api"
	"github.com/hashicorp/vault/command/agent/cache"
	"github.com/hashicorp/vault/helper/namespace"
	http2 "github.com/hashicorp/vault/http"
	"github.com/hashicorp/vault/sdk/framework"
	"github.com/hashicorp/vault/sdk/helper/consts"
	"github.com/hashicorp/vault/sdk/logical"
	"golang.org/x/crypto/argon2"
	"io"
	"net/http"
)

type PasswordStore struct {
	KVPath string `hcl:"kv_path"`
}

type passwordStoreHandler struct {
	framework.Backend
	config  *PasswordStore
	proxier cache.Proxier
	client  *api.Client
}

func NewPasswordStore(config *PasswordStore, proxier cache.Proxier, cli *api.Client) *passwordStoreHandler {
	pstore := &passwordStoreHandler{
		config:  config,
		proxier: proxier,
		client:  cli,
	}
	pstore.Backend =
		framework.Backend{
			Paths: []*framework.Path{
				{
					Pattern: "agent/password_store/" + framework.GenericNameRegex("name") + "/test",
					Fields: map[string]*framework.FieldSchema{
						"password": {
							Type:     framework.TypeString,
							Required: true,
						},
						"name": {
							Type:     framework.TypeString,
							Required: true,
						},
					},
					Operations: map[logical.Operation]framework.OperationHandler{
						logical.CreateOperation: &framework.PathOperation{
							Callback: pstore.handleTest,
						},
						logical.UpdateOperation: &framework.PathOperation{
							Callback: pstore.handleTest,
						},
					},
				},
				{
					Pattern: "agent/password_store/" + framework.GenericNameRegex("name"),
					Fields: map[string]*framework.FieldSchema{
						"password": {
							Type:     framework.TypeString,
							Required: true,
						},
						"name": {
							Type:     framework.TypeString,
							Required: true,
						},
					},
					Operations: map[logical.Operation]framework.OperationHandler{
						logical.CreateOperation: &framework.PathOperation{
							Callback: pstore.handleStore,
						},
						logical.UpdateOperation: &framework.PathOperation{
							Callback: pstore.handleStore,
						},
					},
				},
			},
		}
	return pstore
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

func (p *passwordStoreHandler) ServeHTTP(writer http.ResponseWriter, req *http.Request) {
	req = req.WithContext(namespace.ContextWithNamespace(req.Context(), namespace.RootNamespace))
	lreq, _, _, err := http2.BuildLogicalRequestNoAuth(false, writer, req)
	if err == nil {
		resp, err := p.HandleRequest(context.Background(), lreq)
		if resp.IsError() {
			switch err {
			case logical.ErrInvalidRequest:
				writer.WriteHeader(http.StatusBadRequest)
			default:
				writer.WriteHeader(http.StatusInternalServerError)
			}
		}
	}

}

func (p *passwordStoreHandler) handleTest(_ context.Context, request *logical.Request, fd *framework.FieldData) (*logical.Response, error) {
	key, ok := fd.GetOk("name")
	if !ok {
		return logical.ErrorResponse("missing key name"), logical.ErrInvalidRequest
	}
	pw, ok := fd.GetOk("password")
	if !ok {
		return logical.ErrorResponse("missing password"), logical.ErrInvalidRequest
	}
	token := request.Headers[consts.AuthHeaderName]
	req := p.client.NewRequest(http.MethodGet, "/v1/"+p.config.KVPath+"/"+key.(string))
	req.Headers[consts.AuthHeaderName] = []string{token[0]}
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
	hashed := argonHash(salt, []byte(pw.(string)))
	stored, err := base64.StdEncoding.DecodeString(val.Data.Value)
	if err != nil {
		fmt.Printf("error making request: %e", err)
	}
	if subtle.ConstantTimeCompare(hashed, stored) != 1 {
		return logical.ErrorResponse("invalid password"), logical.ErrInvalidRequest
	} else {
		return nil, nil
	}
}

func (p *passwordStoreHandler) handleStore(_ context.Context, request *logical.Request, fd *framework.FieldData) (*logical.Response, error) {
	key, ok := fd.GetOk("name")
	if !ok {
		return logical.ErrorResponse("missing key name"), logical.ErrInvalidRequest
	}
	pw, ok := fd.GetOk("password")
	if !ok {
		return logical.ErrorResponse("missing password"), logical.ErrInvalidRequest
	}
	token := request.Headers[consts.AuthHeaderName]
	salt := make([]byte, 16)
	_, err := io.ReadFull(rand.Reader, salt)
	if err != nil {
		fmt.Printf("error making request: %e\n", err)
	}
	hashed := argonHash(salt, []byte(pw.(string)))
	saltStr := base64.StdEncoding.EncodeToString(salt)
	hashStr := base64.StdEncoding.EncodeToString(hashed)
	req := p.client.NewRequest(http.MethodPost, "/v1/"+p.config.KVPath+"/"+key.(string))
	req.Headers[consts.AuthHeaderName] = []string{token[0]}
	req.SetJSONBody(map[string]interface{}{"value": hashStr, "salt": saltStr})
	resp, err := p.client.RawRequest(req)
	if err != nil {
		fmt.Printf("error making request: %e\n", err)
	}
	if resp.StatusCode > 299 {
		return logical.ErrorResponse("error storing password"), err
	}
	return nil, nil
}

func argonHash(salt, password []byte) []byte {
	return argon2.IDKey([]byte(password), salt, 1, 64*1024, 4, 28)
}
