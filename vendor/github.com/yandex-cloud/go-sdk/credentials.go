// Copyright (c) 2018 Yandex LLC. All rights reserved.
// Author: Vladimir Skipor <skipor@yandex-team.ru>

package ycsdk

import (
	"context"
	"crypto/rsa"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"net/http/httputil"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/golang/protobuf/ptypes"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/grpclog"
	"google.golang.org/grpc/status"

	iampb "github.com/yandex-cloud/go-genproto/yandex/cloud/iam/v1"
	"github.com/yandex-cloud/go-sdk/iamkey"
	"github.com/yandex-cloud/go-sdk/pkg/sdkerrors"
)

// Credentials is an abstraction of API authorization credentials.
// See https://cloud.yandex.ru/docs/iam/concepts/authorization/authorization for details.
// Note that functions that return Credentials may return different Credentials implementation
// in next SDK version, and this is not considered breaking change.
type Credentials interface {
	// YandexCloudAPICredentials is a marker method. All compatible Credentials implementations have it
	YandexCloudAPICredentials()
}

// ExchangeableCredentials can be exchanged for IAM Token in IAM Token Service, that can be used
// to authorize API calls.
// See https://cloud.yandex.ru/docs/iam/concepts/authorization/iam-token for details.
type ExchangeableCredentials interface {
	Credentials
	// IAMTokenRequest returns request for fresh IAM token or error.
	IAMTokenRequest() (*iampb.CreateIamTokenRequest, error)
}

// NonExchangeableCredentials allows to get IAM Token without calling IAM Token Service.
type NonExchangeableCredentials interface {
	Credentials
	// IAMToken returns IAM Token.
	IAMToken(ctx context.Context) (*iampb.CreateIamTokenResponse, error)
}

// OAuthToken returns API credentials for user Yandex Passport OAuth token, that can be received
// on page https://oauth.yandex.ru/authorize?response_type=token&client_id=1a6990aa636648e9b2ef855fa7bec2fb
// See https://cloud.yandex.ru/docs/iam/concepts/authorization/oauth-token for details.
func OAuthToken(token string) Credentials {
	return exchangeableCredentialsFunc(func() (*iampb.CreateIamTokenRequest, error) {
		return &iampb.CreateIamTokenRequest{
			Identity: &iampb.CreateIamTokenRequest_YandexPassportOauthToken{
				YandexPassportOauthToken: token,
			},
		}, nil
	})
}

// ServiceAccountKey returns credentials for the given IAM Key. The key is used to sign JWT tokens.
// JWT tokens are exchanged for IAM Tokens used to authorize API calls.
// This authorization method is not supported for IAM Keys issued for User Accounts.
func ServiceAccountKey(key *iamkey.Key) (Credentials, error) {
	jwtBuilder, err := newServiceAccountJWTBuilder(key)
	if err != nil {
		return nil, err
	}
	return exchangeableCredentialsFunc(func() (*iampb.CreateIamTokenRequest, error) {
		signedJWT, err := jwtBuilder.SignedToken()
		if err != nil {
			return nil, sdkerrors.WithMessage(err, "JWT sign failed")
		}
		return &iampb.CreateIamTokenRequest{
			Identity: &iampb.CreateIamTokenRequest_Jwt{
				Jwt: signedJWT,
			},
		}, nil
	}), nil
}

// InstanceServiceAccount returns credentials for Compute Instance Service Account.
// That is, for SDK build with InstanceServiceAccount credentials and used on Compute Instance
// created with yandex.cloud.compute.v1.CreateInstanceRequest.service_account_id, API calls
// will be authenticated with this ServiceAccount ID.
// TODO(skipor): doc link
func InstanceServiceAccount() NonExchangeableCredentials {
	return newInstanceServiceAccountCredentials(InstanceMetadataAddr)
}

func newServiceAccountJWTBuilder(key *iamkey.Key) (*serviceAccountJWTBuilder, error) {
	err := validateServiceAccountKey(key)
	if err != nil {
		return nil, sdkerrors.WithMessage(err, "key validation failed")
	}
	rsaPrivateKey, err := jwt.ParseRSAPrivateKeyFromPEM([]byte(key.PrivateKey))
	if err != nil {
		return nil, sdkerrors.WithMessage(err, "private key parsing failed")
	}
	return &serviceAccountJWTBuilder{
		key:           key,
		rsaPrivateKey: rsaPrivateKey,
	}, nil
}

func validateServiceAccountKey(key *iamkey.Key) error {
	if key.Id == "" {
		return errors.New("key id is missing")
	}
	if key.GetServiceAccountId() == "" {
		return fmt.Errorf("key should de issued for service account, but subject is %#v", key.Subject)
	}
	return nil
}

type serviceAccountJWTBuilder struct {
	key           *iamkey.Key
	rsaPrivateKey *rsa.PrivateKey
}

func (b *serviceAccountJWTBuilder) SignedToken() (string, error) {
	return b.issueToken().SignedString(b.rsaPrivateKey)
}

func (b *serviceAccountJWTBuilder) issueToken() *jwt.Token {
	issuedAt := time.Now()
	token := jwt.NewWithClaims(jwtSigningMethodPS256WithSaltLengthEqualsHash, jwt.StandardClaims{
		Issuer:    b.key.GetServiceAccountId(),
		IssuedAt:  issuedAt.Unix(),
		ExpiresAt: issuedAt.Add(time.Hour).Unix(),
		Audience:  "https://iam.api.cloud.yandex.net/iam/v1/tokens",
	})
	token.Header["kid"] = b.key.Id
	return token
}

// NOTE(skipor): by default, Go RSA PSS uses PSSSaltLengthAuto, which is not accepted by jwt.io and some python libraries.
// Should be removed after https://github.com/dgrijalva/jwt-go/issues/285 fix.
var jwtSigningMethodPS256WithSaltLengthEqualsHash = &jwt.SigningMethodRSAPSS{
	SigningMethodRSA: jwt.SigningMethodPS256.SigningMethodRSA,
	Options: &rsa.PSSOptions{
		SaltLength: rsa.PSSSaltLengthEqualsHash,
	},
}

type exchangeableCredentialsFunc func() (iamTokenReq *iampb.CreateIamTokenRequest, err error)

var _ ExchangeableCredentials = (exchangeableCredentialsFunc)(nil)

func (exchangeableCredentialsFunc) YandexCloudAPICredentials() {}
func (f exchangeableCredentialsFunc) IAMTokenRequest() (iamTokenReq *iampb.CreateIamTokenRequest, err error) {
	return f()
}

func newInstanceServiceAccountCredentials(metadataServiceAddr string) NonExchangeableCredentials {
	return &instanceServiceAccountCredentials{
		metadataServiceAddr: metadataServiceAddr,
		client: http.Client{
			Transport: &http.Transport{
				DialContext: (&net.Dialer{
					Timeout:   time.Second, // One second should be enough for localhost connection.
					KeepAlive: -1,          // No keep alive. Near token per hour requested.
				}).DialContext,
			},
		},
	}
}

type instanceServiceAccountCredentials struct {
	metadataServiceAddr string
	client              http.Client
}

func (c *instanceServiceAccountCredentials) YandexCloudAPICredentials() {}

func (c *instanceServiceAccountCredentials) IAMToken(ctx context.Context) (*iampb.CreateIamTokenResponse, error) {
	token, err := c.iamToken(ctx)
	if err != nil {
		return nil, sdkerrors.WithMessagef(err, "failed to get compute instance service account token from instance metadata service: GET %s", c.url())
	}
	return token, nil
}

func (c *instanceServiceAccountCredentials) iamToken(ctx context.Context) (*iampb.CreateIamTokenResponse, error) {
	req, err := http.NewRequest("GET", c.url(), nil)
	if err != nil {
		return nil, sdkerrors.WithMessage(err, "request make failed")
	}
	req.Header.Set("Metadata-Flavor", "Google")
	reqDump, _ := httputil.DumpRequestOut(req, false)
	grpclog.Infof("Going to request instance SA token in metadata service:\n%s", reqDump)
	resp, err := c.client.Do(req.WithContext(ctx))
	if err != nil {
		return nil, fmt.Errorf("%s.\n"+
			"Are you inside compute instance?",
			err)
	}
	defer resp.Body.Close()
	respDump, _ := httputil.DumpResponse(resp, false)
	grpclog.Infof("Metadata service instance SA token response (without body, because contains sensitive token):\n%s", respDump)

	if resp.StatusCode == http.StatusNotFound {
		return nil, fmt.Errorf("%s.\n"+
			"Is this compute instance running using Service Account? That is, Instance.service_account_id should not be empty.",
			resp.Status)
	}
	body, err := ioutil.ReadAll(resp.Body)

	if resp.StatusCode != http.StatusOK {
		if err != nil {
			body = []byte(fmt.Sprintf("Failed response body read failed: %s", err.Error()))
		}
		grpclog.Errorf("Metadata service instance SA token get failed: %s. Body:\n%s", resp.Status, body)
		return nil, fmt.Errorf("%s", resp.Status)
	}
	if err != nil {
		return nil, fmt.Errorf("reponse read failed: %s", err)
	}

	var tokenResponse struct {
		AccessToken string `json:"access_token"`
		ExpiresIn   int64  `json:"expires_in"`
		TokenType   string `json:"token_type"`
	}

	err = json.Unmarshal(body, &tokenResponse)
	if err != nil {
		grpclog.Errorf("Failed to unmarshal instance metadata service SA token response body.\nError: %s\nBody:\n%s", err, body)
		return nil, sdkerrors.WithMessage(err, "body unmarshal failed")
	}
	expiresAt := ptypes.TimestampNow()
	expiresAt.Seconds += tokenResponse.ExpiresIn - 1
	expiresAt.Nanos = 0 // Truncate is for readability.
	return &iampb.CreateIamTokenResponse{
		IamToken:  tokenResponse.AccessToken,
		ExpiresAt: expiresAt,
	}, nil
}

func (c *instanceServiceAccountCredentials) url() string {
	return fmt.Sprintf("http://%s/computeMetadata/v1/instance/service-accounts/default/token", c.metadataServiceAddr)
}

// NoCredentials implements Credentials, it allows to create unauthenticated connections
type NoCredentials struct{}

func (creds NoCredentials) YandexCloudAPICredentials() {}

// IAMToken always returns gRPC error with status UNAUTHENTICATED
func (creds NoCredentials) IAMToken(ctx context.Context) (*iampb.CreateIamTokenResponse, error) {
	return nil, status.Error(codes.Unauthenticated, "unauthenticated connection")
}

// IAMTokenCredentials implements Credentials with IAM token as-is
type IAMTokenCredentials struct {
	iamToken string
}

func (creds IAMTokenCredentials) YandexCloudAPICredentials() {}

func (creds IAMTokenCredentials) IAMToken(ctx context.Context) (*iampb.CreateIamTokenResponse, error) {
	return &iampb.CreateIamTokenResponse{
		IamToken: creds.iamToken,
	}, nil
}

func NewIAMTokenCredentials(iamToken string) NonExchangeableCredentials {
	return &IAMTokenCredentials{
		iamToken: iamToken,
	}
}
