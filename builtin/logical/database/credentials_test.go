package database

import (
	"context"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"testing"

	"github.com/hashicorp/vault/sdk/helper/base62"
	"github.com/hashicorp/vault/sdk/logical"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func Test_newPasswordGenerator(t *testing.T) {
	type args struct {
		config map[string]interface{}
	}
	tests := []struct {
		name    string
		args    args
		want    passwordGenerator
		wantErr bool
	}{
		{
			name: "newPasswordGenerator with nil config",
			args: args{
				config: nil,
			},
			want: passwordGenerator{
				PasswordPolicy: "",
			},
		},
		{
			name: "newPasswordGenerator without password_policy",
			args: args{
				config: map[string]interface{}{},
			},
			want: passwordGenerator{
				PasswordPolicy: "",
			},
		},
		{
			name: "newPasswordGenerator with password_policy",
			args: args{
				config: map[string]interface{}{
					"password_policy": "test-policy",
				},
			},
			want: passwordGenerator{
				PasswordPolicy: "test-policy",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := newPasswordGenerator(tt.args.config)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			assert.Equal(t, tt.want, got)
		})
	}
}

func Test_newRSAKeyGenerator(t *testing.T) {
	type args struct {
		config map[string]interface{}
	}
	tests := []struct {
		name    string
		args    args
		want    rsaKeyGenerator
		wantErr bool
	}{
		{
			name: "newRSAKeyGenerator with nil config",
			args: args{
				config: nil,
			},
			want: rsaKeyGenerator{
				Format:  "pkcs8",
				KeyBits: 2048,
			},
		},
		{
			name: "newRSAKeyGenerator with empty config",
			args: args{
				config: map[string]interface{}{},
			},
			want: rsaKeyGenerator{
				Format:  "pkcs8",
				KeyBits: 2048,
			},
		},
		{
			name: "newRSAKeyGenerator with zero value format",
			args: args{
				config: map[string]interface{}{
					"format": "",
				},
			},
			want: rsaKeyGenerator{
				Format:  "pkcs8",
				KeyBits: 2048,
			},
		},
		{
			name: "newRSAKeyGenerator with zero value key_bits",
			args: args{
				config: map[string]interface{}{
					"key_bits": "0",
				},
			},
			want: rsaKeyGenerator{
				Format:  "pkcs8",
				KeyBits: 2048,
			},
		},
		{
			name: "newRSAKeyGenerator with format",
			args: args{
				config: map[string]interface{}{
					"format": "pkcs8",
				},
			},
			want: rsaKeyGenerator{
				Format:  "pkcs8",
				KeyBits: 2048,
			},
		},
		{
			name: "newRSAKeyGenerator with format case insensitive",
			args: args{
				config: map[string]interface{}{
					"format": "PKCS8",
				},
			},
			want: rsaKeyGenerator{
				Format:  "PKCS8",
				KeyBits: 2048,
			},
		},
		{
			name: "newRSAKeyGenerator with 3072 key_bits",
			args: args{
				config: map[string]interface{}{
					"key_bits": "3072",
				},
			},
			want: rsaKeyGenerator{
				Format:  "pkcs8",
				KeyBits: 3072,
			},
		},
		{
			name: "newRSAKeyGenerator with 4096 key_bits",
			args: args{
				config: map[string]interface{}{
					"key_bits": "4096",
				},
			},
			want: rsaKeyGenerator{
				Format:  "pkcs8",
				KeyBits: 4096,
			},
		},
		{
			name: "newRSAKeyGenerator with invalid key_bits",
			args: args{
				config: map[string]interface{}{
					"key_bits": "4097",
				},
			},
			wantErr: true,
		},
		{
			name: "newRSAKeyGenerator with invalid format",
			args: args{
				config: map[string]interface{}{
					"format": "pkcs1",
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := newRSAKeyGenerator(tt.args.config)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			assert.Equal(t, tt.want, got)
		})
	}
}

func Test_passwordGenerator_generate(t *testing.T) {
	config := logical.TestBackendConfig()
	b := Backend(config)
	b.Setup(context.Background(), config)

	type args struct {
		config  map[string]interface{}
		mock    func() interface{}
		passGen logical.PasswordGenerator
	}
	tests := []struct {
		name       string
		args       args
		wantRegexp string
		wantErr    bool
	}{
		{
			name: "wrapper missing v4 and v5 interface",
			args: args{
				mock: func() interface{} {
					return nil
				},
			},
			wantErr: true,
		},
		{
			name: "v4: generate password using GenerateCredentials",
			args: args{
				mock: func() interface{} {
					v4Mock := new(mockLegacyDatabase)
					v4Mock.On("GenerateCredentials", mock.Anything).
						Return("v4-generated-password", nil).
						Times(1)
					return v4Mock
				},
			},
			wantRegexp: "^v4-generated-password$",
		},
		{
			name: "v5: generate password without policy",
			args: args{
				mock: func() interface{} {
					return new(mockNewDatabase)
				},
			},
			wantRegexp: "^[a-zA-Z0-9-]{20}$",
		},
		{
			name: "v5: generate password with non-existing policy",
			args: args{
				config: map[string]interface{}{
					"password_policy": "not-created",
				},
				mock: func() interface{} {
					return new(mockNewDatabase)
				},
			},
			wantErr: true,
		},
		{
			name: "v5: generate password with existing policy",
			args: args{
				config: map[string]interface{}{
					"password_policy": "test-policy",
				},
				mock: func() interface{} {
					return new(mockNewDatabase)
				},
				passGen: func() (string, error) {
					return base62.Random(30)
				},
			},
			wantRegexp: "^[a-zA-Z0-9]{30}$",
		},
		{
			name: "v5: generate password with existing policy static",
			args: args{
				config: map[string]interface{}{
					"password_policy": "test-policy",
				},
				mock: func() interface{} {
					return new(mockNewDatabase)
				},
				passGen: func() (string, error) {
					return "policy-generated-password", nil
				},
			},
			wantRegexp: "^policy-generated-password$",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Set up the version wrapper with a mock database implementation
			wrapper := databaseVersionWrapper{}
			switch m := tt.args.mock().(type) {
			case *mockNewDatabase:
				wrapper.v5 = m
			case *mockLegacyDatabase:
				wrapper.v4 = m
			}

			// Set the password policy for the test case
			config.System.(*logical.StaticSystemView).SetPasswordPolicy(
				"test-policy", tt.args.passGen)

			// Generate the password
			pg, err := newPasswordGenerator(tt.args.config)
			got, err := pg.generate(context.Background(), b, wrapper)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			assert.Regexp(t, tt.wantRegexp, got)

			// Assert all expected calls took place on the mock
			if m, ok := wrapper.v5.(*mockNewDatabase); ok {
				m.AssertExpectations(t)
			}
			if m, ok := wrapper.v4.(*mockLegacyDatabase); ok {
				m.AssertExpectations(t)
			}
		})
	}
}

func Test_passwordGenerator_configMap(t *testing.T) {
	type args struct {
		config map[string]interface{}
	}
	tests := []struct {
		name string
		args args
		want map[string]interface{}
	}{
		{
			name: "nil config results in empty map",
			args: args{
				config: nil,
			},
			want: map[string]interface{}{},
		},
		{
			name: "empty config results in empty map",
			args: args{
				config: map[string]interface{}{},
			},
			want: map[string]interface{}{},
		},
		{
			name: "input config is equal to output config",
			args: args{
				config: map[string]interface{}{
					"password_policy": "test-policy",
				},
			},
			want: map[string]interface{}{
				"password_policy": "test-policy",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pg, err := newPasswordGenerator(tt.args.config)
			assert.NoError(t, err)
			cm, err := pg.configMap()
			assert.NoError(t, err)
			assert.Equal(t, tt.want, cm)
		})
	}
}

func Test_rsaKeyGenerator_generate(t *testing.T) {
	type args struct {
		config map[string]interface{}
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "generate RSA key with nil default config",
			args: args{
				config: nil,
			},
		},
		{
			name: "generate RSA key with empty default config",
			args: args{
				config: map[string]interface{}{},
			},
		},
		{
			name: "generate RSA key with 2048 key_bits and format",
			args: args{
				config: map[string]interface{}{
					"key_bits": "2048",
					"format":   "pkcs8",
				},
			},
		},
		{
			name: "generate RSA key with 2048 key_bits",
			args: args{
				config: map[string]interface{}{
					"key_bits": "2048",
				},
			},
		},
		{
			name: "generate RSA key with 3072 key_bits",
			args: args{
				config: map[string]interface{}{
					"key_bits": "3072",
				},
			},
		},
		{
			name: "generate RSA key with 4096 key_bits",
			args: args{
				config: map[string]interface{}{
					"key_bits": "4096",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Generate the RSA key pair
			kg, err := newRSAKeyGenerator(tt.args.config)
			public, private, err := kg.generate(rand.Reader)
			assert.NoError(t, err)
			assert.NotEmpty(t, public)
			assert.NotEmpty(t, private)

			// Decode the public and private key PEMs
			pubBlock, pubRest := pem.Decode(public)
			privBlock, privRest := pem.Decode(private)
			assert.NotNil(t, pubBlock)
			assert.Empty(t, pubRest)
			assert.Equal(t, "PUBLIC KEY", pubBlock.Type)
			assert.NotNil(t, privBlock)
			assert.Empty(t, privRest)
			assert.Equal(t, "PRIVATE KEY", privBlock.Type)

			// Assert that we can parse the public key PEM block
			pub, err := x509.ParsePKIXPublicKey(pubBlock.Bytes)
			assert.NoError(t, err)
			assert.NotNil(t, pub)
			assert.IsType(t, &rsa.PublicKey{}, pub)

			// Assert that we can parse the private key PEM block in
			// the configured format
			switch kg.Format {
			case "pkcs8":
				priv, err := x509.ParsePKCS8PrivateKey(privBlock.Bytes)
				assert.NoError(t, err)
				assert.NotNil(t, priv)
				assert.IsType(t, &rsa.PrivateKey{}, priv)
			default:
				t.Fatal("unknown format")
			}
		})
	}
}

func Test_rsaKeyGenerator_configMap(t *testing.T) {
	type args struct {
		config map[string]interface{}
	}
	tests := []struct {
		name string
		args args
		want map[string]interface{}
	}{
		{
			name: "nil config results in defaults",
			args: args{
				config: nil,
			},
			want: map[string]interface{}{
				"format":   "pkcs8",
				"key_bits": 2048,
			},
		},
		{
			name: "empty config results in defaults",
			args: args{
				config: map[string]interface{}{},
			},
			want: map[string]interface{}{
				"format":   "pkcs8",
				"key_bits": 2048,
			},
		},
		{
			name: "config with format",
			args: args{
				config: map[string]interface{}{
					"format": "pkcs8",
				},
			},
			want: map[string]interface{}{
				"format":   "pkcs8",
				"key_bits": 2048,
			},
		},
		{
			name: "config with key_bits",
			args: args{
				config: map[string]interface{}{
					"key_bits": 4096,
				},
			},
			want: map[string]interface{}{
				"format":   "pkcs8",
				"key_bits": 4096,
			},
		},
		{
			name: "config with format and key_bits",
			args: args{
				config: map[string]interface{}{
					"format":   "pkcs8",
					"key_bits": 3072,
				},
			},
			want: map[string]interface{}{
				"format":   "pkcs8",
				"key_bits": 3072,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			kg, err := newRSAKeyGenerator(tt.args.config)
			assert.NoError(t, err)
			cm, err := kg.configMap()
			assert.NoError(t, err)
			assert.Equal(t, tt.want, cm)
		})
	}
}
