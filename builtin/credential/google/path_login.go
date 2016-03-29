package google

import (
	"fmt"
	"reflect"
	"sort"
	"errors"

	"github.com/hashicorp/vault/logical"
	"github.com/hashicorp/vault/logical/framework"

	"golang.org/x/oauth2"
	"golang.org/x/net/context"
	goauth "google.golang.org/api/oauth2/v2"
	"time"
)

const loginPath = "login"
const googleAuthCodeParameterName = "code"

func pathLogin(b *backend) *framework.Path {
	return &framework.Path{
		Pattern: loginPath,
		Fields: map[string]*framework.FieldSchema{
			googleAuthCodeParameterName: &framework.FieldSchema{
				Type:        framework.TypeString,
				Description: "Google authentication code",
			},
		},

		Callbacks: map[logical.Operation]framework.OperationFunc{
			logical.UpdateOperation: b.pathLogin,
		},
	}
}

const refreshToken = "refreshToken"

func (b *backend) pathLogin(
	req *logical.Request, data *framework.FieldData) (*logical.Response, error) {

	code := data.Get(googleAuthCodeParameterName).(string)

	var verifyResp *verifyCredentialsResp
	if verifyResponse, resp, err := b.verifyCredentials(req, code, nil); err != nil {
		return nil, err
	} else if resp != nil {
		return resp, nil
	} else {
		verifyResp = verifyResponse
	}

	config, err := b.Config(req.Storage)
	if err != nil {
		return nil, err
	}

	ttl, _, err := b.SanitizeTTL(config.TTL.String(), config.MaxTTL.String())
	if err != nil {
		return logical.ErrorResponse(fmt.Sprintf("[ERR]:%s", err)), nil
	}

	internalData := map[string]interface{}{
		refreshToken: verifyResp.RefreshToken,
	}

	return &logical.Response{
		Auth: &logical.Auth{
			InternalData: internalData,
			Policies: verifyResp.Policies,
			Metadata: map[string]string{
				"username": 	verifyResp.User,
				"domain":      	verifyResp.Domain,
			},
			DisplayName: 		verifyResp.Name,
			LeaseOptions: logical.LeaseOptions{
				TTL:       ttl,
				Renewable: true,
			},
		},
	}, nil
}

func (b *backend) pathLoginRenew(
	req *logical.Request, d *framework.FieldData) (*logical.Response, error) {

	previousTokenObject := req.Auth.InternalData[refreshToken]
	if (previousTokenObject == nil) {
		return nil, errors.New("no refresh token from previous login")
	}
	previousTokenMap := previousTokenObject.(map[string]interface{})
	expiryString := previousTokenMap["expiry"].(string)
	expiry, err := time.Parse(time.RFC3339Nano, expiryString)
	if (err != nil) {
		return nil, fmt.Errorf("could not parse time (%s) from persisted token", expiryString)
	}
	refreshToken := &oauth2.Token{
		AccessToken: previousTokenMap["access_token"].(string),
		TokenType: previousTokenMap["token_type"].(string),
		RefreshToken: previousTokenMap["refresh_token"].(string),
		Expiry: expiry,
	}

	var verifyResp *verifyCredentialsResp

	if verifyResponse, resp, err := b.verifyCredentials(req, "", refreshToken); err != nil {
		return nil, err
	} else if resp != nil {
		return resp, nil
	} else {
		verifyResp = verifyResponse
	}
	sort.Strings(req.Auth.Policies)
	if !reflect.DeepEqual(sliceToMap(verifyResp.Policies), sliceToMap(req.Auth.Policies)) {
		return logical.ErrorResponse(fmt.Sprintf("policies do not match. new policies: %s. old policies: %s.", verifyResp.Policies, req.Auth.Policies)), nil
	}

	config, err := b.Config(req.Storage)
	if err != nil {
		return nil, err
	}
	return framework.LeaseExtend(config.TTL, config.MaxTTL, b.System())(req, d)
}

func (b *backend) verifyCredentials(req *logical.Request, code string, tok *oauth2.Token) (*verifyCredentialsResp, *logical.Response, error) {

	config, err := b.Config(req.Storage)

	if err != nil {
		return nil, nil, err
	}

	if config.ApplicationID == "" {
		return nil, logical.ErrorResponse(writeConfigPathHelp), nil
	}

	if config.ApplicationSecret == "" {
		return nil, logical.ErrorResponse(writeConfigPathHelp), nil
	}

	googleConfig := applicationOauth2Config(config.ApplicationID, config.ApplicationSecret)

	if (tok == nil && code != "") {
		tok, err = googleConfig.Exchange(oauth2.NoContext, code)
		if (err != nil) {
			return nil, nil, err
		}
	}

	httpClient := googleConfig.Client(context.Background(), tok)
	service, err := goauth.New(httpClient)
	if (err != nil) {
		return nil, nil, err
	}

	me := goauth.NewUserinfoV2MeService(service)
	info, err := me.Get().Do()
	if (err != nil) {
		return nil, nil, err
	}

	user := info.Email
	domain := info.Hd


	if domain != config.Domain {
		return nil, logical.ErrorResponse(fmt.Sprintf("user %s is of domain %s, not part of required domain %s", user, domain, config.Domain)), nil
	}

	userId := localPartFromEmail(user)

	policiesList, err := b.Map.Policies(req.Storage, userId)
	//be compatible with core, see issue https://github.com/hashicorp/vault/issues/1256
	if strListContains(policiesList, "root") {
		policiesList = []string{"root"}
	} else {
		policiesList = append(policiesList, "default")
	}

	if err != nil {
		return nil, nil, err
	}
	return &verifyCredentialsResp{
		User:     user,
		Domain:      domain,
		Policies: policiesList,
		RefreshToken: tok,
		Name: info.Name,
	}, nil, nil
}

type verifyCredentialsResp struct {
	User    string
	Domain  string
	Name	string
	Policies []string
	RefreshToken *oauth2.Token
}
