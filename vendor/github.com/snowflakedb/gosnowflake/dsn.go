// Copyright (c) 2017-2022 Snowflake Computing Inc. All rights reserved.

package gosnowflake

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"errors"
	"fmt"
	"net"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"
)

const (
	defaultClientTimeout          = 900 * time.Second // Timeout for network round trip + read out http response
	defaultJWTClientTimeout       = 10 * time.Second  // Timeout for network round trip + read out http response but used for JWT auth
	defaultLoginTimeout           = 300 * time.Second // Timeout for retry for login EXCLUDING clientTimeout
	defaultRequestTimeout         = 0 * time.Second   // Timeout for retry for request EXCLUDING clientTimeout
	defaultJWTTimeout             = 60 * time.Second
	defaultExternalBrowserTimeout = 120 * time.Second // Timeout for external browser login
	defaultMaxRetryCount          = 7                 // specifies maximum number of subsequent retries
	defaultDomain                 = ".snowflakecomputing.com"
	cnDomain                      = ".snowflakecomputing.cn"
	topLevelDomainPrefix          = ".snowflakecomputing." // used to extract the domain from host
)

// ConfigBool is a type to represent true or false in the Config
type ConfigBool uint8

const (
	configBoolNotSet ConfigBool = iota // Reserved for unset to let default value fall into this category
	// ConfigBoolTrue represents true for the config field
	ConfigBoolTrue
	// ConfigBoolFalse represents false for the config field
	ConfigBoolFalse
)

// Config is a set of configuration parameters
type Config struct {
	Account   string // Account name
	User      string // Username
	Password  string // Password (requires User)
	Database  string // Database name
	Schema    string // Schema
	Warehouse string // Warehouse
	Role      string // Role
	Region    string // Region

	// ValidateDefaultParameters disable the validation checks for Database, Schema, Warehouse and Role
	// at the time a connection is established
	ValidateDefaultParameters ConfigBool

	Params map[string]*string // other connection parameters

	ClientIP net.IP // IP address for network check
	Protocol string // http or https (optional)
	Host     string // hostname (optional)
	Port     int    // port (optional)

	Authenticator AuthType // The authenticator type

	Passcode           string
	PasscodeInPassword bool

	OktaURL *url.URL

	LoginTimeout           time.Duration // Login retry timeout EXCLUDING network roundtrip and read out http response
	RequestTimeout         time.Duration // request retry timeout EXCLUDING network roundtrip and read out http response
	JWTExpireTimeout       time.Duration // JWT expire after timeout
	ClientTimeout          time.Duration // Timeout for network round trip + read out http response
	JWTClientTimeout       time.Duration // Timeout for network round trip + read out http response used when JWT token auth is taking place
	ExternalBrowserTimeout time.Duration // Timeout for external browser login
	MaxRetryCount          int           // Specifies how many times non-periodic HTTP request can be retried

	Application  string           // application name.
	InsecureMode bool             // driver doesn't check certificate revocation status
	OCSPFailOpen OCSPFailOpenMode // OCSP Fail Open

	Token            string        // Token to use for OAuth other forms of token based auth
	TokenAccessor    TokenAccessor // Optional token accessor to use
	KeepSessionAlive bool          // Enables the session to persist even after the connection is closed

	PrivateKey *rsa.PrivateKey // Private key used to sign JWT

	Transporter http.RoundTripper // RoundTripper to intercept HTTP requests and responses

	DisableTelemetry bool // indicates whether to disable telemetry

	Tracing string // sets logging level

	TmpDirPath string // sets temporary directory used by a driver for operations like encrypting, compressing etc

	MfaToken                       string     // Internally used to cache the MFA token
	IDToken                        string     // Internally used to cache the Id Token for external browser
	ClientRequestMfaToken          ConfigBool // When true the MFA token is cached in the credential manager. True by default in Windows/OSX. False for Linux.
	ClientStoreTemporaryCredential ConfigBool // When true the ID token is cached in the credential manager. True by default in Windows/OSX. False for Linux.

	DisableQueryContextCache bool // Should HTAP query context cache be disabled

	IncludeRetryReason ConfigBool // Should retried request contain retry reason

	ClientConfigFile string // File path to the client configuration json file

	DisableConsoleLogin ConfigBool // Indicates whether console login should be disabled

	DisableSamlURLCheck ConfigBool // Indicates whether the SAML URL check should be disabled
}

// Validate enables testing if config is correct.
// A driver client may call it manually, but it is also called during opening first connection.
func (c *Config) Validate() error {
	if c.TmpDirPath != "" {
		if _, err := os.Stat(c.TmpDirPath); err != nil {
			return err
		}
	}
	return nil
}

// ocspMode returns the OCSP mode in string INSECURE, FAIL_OPEN, FAIL_CLOSED
func (c *Config) ocspMode() string {
	if c.InsecureMode {
		return ocspModeInsecure
	} else if c.OCSPFailOpen == ocspFailOpenNotSet || c.OCSPFailOpen == OCSPFailOpenTrue {
		// by default or set to true
		return ocspModeFailOpen
	}
	return ocspModeFailClosed
}

// DSN constructs a DSN for Snowflake db.
func DSN(cfg *Config) (dsn string, err error) {
	if cfg.Region == "us-west-2" {
		cfg.Region = ""
	}
	// in case account includes region
	region, posDot := extractRegionFromAccount(cfg.Account)
	if region != "" {
		if cfg.Region != "" {
			return "", errRegionConflict()
		}
		cfg.Region = region
		cfg.Account = cfg.Account[:posDot]
	}
	hasHost := true
	if cfg.Host == "" {
		hasHost = false
		if cfg.Region == "" {
			cfg.Host = cfg.Account + defaultDomain
		} else {
			cfg.Host = buildHostFromAccountAndRegion(cfg.Account, cfg.Region)
		}
	}
	err = fillMissingConfigParameters(cfg)
	if err != nil {
		return "", err
	}
	params := &url.Values{}
	if hasHost && cfg.Account != "" {
		// account may not be included in a Host string
		params.Add("account", cfg.Account)
	}
	if cfg.Database != "" {
		params.Add("database", cfg.Database)
	}
	if cfg.Schema != "" {
		params.Add("schema", cfg.Schema)
	}
	if cfg.Warehouse != "" {
		params.Add("warehouse", cfg.Warehouse)
	}
	if cfg.Role != "" {
		params.Add("role", cfg.Role)
	}
	if cfg.Region != "" {
		params.Add("region", cfg.Region)
	}
	if cfg.Authenticator != AuthTypeSnowflake {
		if cfg.Authenticator == AuthTypeOkta {
			params.Add("authenticator", strings.ToLower(cfg.OktaURL.String()))
		} else {
			params.Add("authenticator", strings.ToLower(cfg.Authenticator.String()))
		}
	}
	if cfg.Passcode != "" {
		params.Add("passcode", cfg.Passcode)
	}
	if cfg.PasscodeInPassword {
		params.Add("passcodeInPassword", strconv.FormatBool(cfg.PasscodeInPassword))
	}
	if cfg.ClientTimeout != defaultClientTimeout {
		params.Add("clientTimeout", strconv.FormatInt(int64(cfg.ClientTimeout/time.Second), 10))
	}
	if cfg.JWTClientTimeout != defaultJWTClientTimeout {
		params.Add("jwtClientTimeout", strconv.FormatInt(int64(cfg.JWTClientTimeout/time.Second), 10))
	}
	if cfg.LoginTimeout != defaultLoginTimeout {
		params.Add("loginTimeout", strconv.FormatInt(int64(cfg.LoginTimeout/time.Second), 10))
	}
	if cfg.RequestTimeout != defaultRequestTimeout {
		params.Add("requestTimeout", strconv.FormatInt(int64(cfg.RequestTimeout/time.Second), 10))
	}
	if cfg.JWTExpireTimeout != defaultJWTTimeout {
		params.Add("jwtTimeout", strconv.FormatInt(int64(cfg.JWTExpireTimeout/time.Second), 10))
	}
	if cfg.ExternalBrowserTimeout != defaultExternalBrowserTimeout {
		params.Add("externalBrowserTimeout", strconv.FormatInt(int64(cfg.ExternalBrowserTimeout/time.Second), 10))
	}
	if cfg.MaxRetryCount != defaultMaxRetryCount {
		params.Add("maxRetryCount", strconv.Itoa(cfg.MaxRetryCount))
	}
	if cfg.Application != clientType {
		params.Add("application", cfg.Application)
	}
	if cfg.Protocol != "" && cfg.Protocol != "https" {
		params.Add("protocol", cfg.Protocol)
	}
	if cfg.Token != "" {
		params.Add("token", cfg.Token)
	}
	if cfg.Params != nil {
		for k, v := range cfg.Params {
			params.Add(k, *v)
		}
	}
	if cfg.PrivateKey != nil {
		privateKeyInBytes, err := marshalPKCS8PrivateKey(cfg.PrivateKey)
		if err != nil {
			return "", err
		}
		keyBase64 := base64.URLEncoding.EncodeToString(privateKeyInBytes)
		params.Add("privateKey", keyBase64)
	}
	if cfg.InsecureMode {
		params.Add("insecureMode", strconv.FormatBool(cfg.InsecureMode))
	}
	if cfg.Tracing != "" {
		params.Add("tracing", cfg.Tracing)
	}
	if cfg.TmpDirPath != "" {
		params.Add("tmpDirPath", cfg.TmpDirPath)
	}
	if cfg.DisableQueryContextCache {
		params.Add("disableQueryContextCache", "true")
	}
	if cfg.IncludeRetryReason == ConfigBoolFalse {
		params.Add("includeRetryReason", "false")
	}

	params.Add("ocspFailOpen", strconv.FormatBool(cfg.OCSPFailOpen != OCSPFailOpenFalse))

	params.Add("validateDefaultParameters", strconv.FormatBool(cfg.ValidateDefaultParameters != ConfigBoolFalse))

	if cfg.ClientRequestMfaToken != configBoolNotSet {
		params.Add("clientRequestMfaToken", strconv.FormatBool(cfg.ClientRequestMfaToken != ConfigBoolFalse))
	}

	if cfg.ClientStoreTemporaryCredential != configBoolNotSet {
		params.Add("clientStoreTemporaryCredential", strconv.FormatBool(cfg.ClientStoreTemporaryCredential != ConfigBoolFalse))
	}
	if cfg.ClientConfigFile != "" {
		params.Add("clientConfigFile", cfg.ClientConfigFile)
	}
	if cfg.DisableConsoleLogin != configBoolNotSet {
		params.Add("disableConsoleLogin", strconv.FormatBool(cfg.DisableConsoleLogin != ConfigBoolFalse))
	}
	if cfg.DisableSamlURLCheck != configBoolNotSet {
		params.Add("disableSamlURLCheck", strconv.FormatBool(cfg.DisableSamlURLCheck != ConfigBoolFalse))
	}

	dsn = fmt.Sprintf("%v:%v@%v:%v", url.QueryEscape(cfg.User), url.QueryEscape(cfg.Password), cfg.Host, cfg.Port)
	if params.Encode() != "" {
		dsn += "?" + params.Encode()
	}
	return
}

// ParseDSN parses the DSN string to a Config.
func ParseDSN(dsn string) (cfg *Config, err error) {
	// New config with some default values
	cfg = &Config{
		Params:        make(map[string]*string),
		Authenticator: AuthTypeSnowflake, // Default to snowflake
	}

	// user[:password]@account/database/schema[?param1=value1&paramN=valueN]
	// or
	// user[:password]@account/database[?param1=value1&paramN=valueN]
	// or
	// user[:password]@host:port/database/schema?account=user_account[?param1=value1&paramN=valueN]
	// or
	// host:port/database/schema?account=user_account[?param1=value1&paramN=valueN]

	foundSlash := false
	secondSlash := false
	done := false
	var i int
	posQuestion := len(dsn)
	for i = len(dsn) - 1; i >= 0; i-- {
		switch {
		case dsn[i] == '/':
			foundSlash = true

			// left part is empty if i <= 0
			var j int
			posSecondSlash := i
			if i > 0 {
				for j = i - 1; j >= 0; j-- {
					switch {
					case dsn[j] == '/':
						// second slash
						secondSlash = true
						posSecondSlash = j
					case dsn[j] == '@':
						// username[:password]@...
						cfg.User, cfg.Password = parseUserPassword(j, dsn)
					}
					if dsn[j] == '@' {
						break
					}
				}

				// account or host:port
				err = parseAccountHostPort(cfg, j, posSecondSlash, dsn)
				if err != nil {
					return nil, err
				}
			}
			// [?param1=value1&...&paramN=valueN]
			// Find the first '?' in dsn[i+1:]
			err = parseParams(cfg, i, dsn)
			if err != nil {
				return
			}
			if secondSlash {
				cfg.Database = dsn[posSecondSlash+1 : i]
				cfg.Schema = dsn[i+1 : posQuestion]
			} else {
				cfg.Database = dsn[posSecondSlash+1 : posQuestion]
			}
			done = true
		case dsn[i] == '?':
			posQuestion = i
		}
		if done {
			break
		}
	}
	if !foundSlash {
		// no db or schema is specified
		var j int
		for j = len(dsn) - 1; j >= 0; j-- {
			switch {
			case dsn[j] == '@':
				cfg.User, cfg.Password = parseUserPassword(j, dsn)
			case dsn[j] == '?':
				posQuestion = j
			}
			if dsn[j] == '@' {
				break
			}
		}
		err = parseAccountHostPort(cfg, j, posQuestion, dsn)
		if err != nil {
			return nil, err
		}
		err = parseParams(cfg, posQuestion-1, dsn)
		if err != nil {
			return
		}
	}
	if cfg.Account == "" && hostIncludesTopLevelDomain(cfg.Host) {
		posDot := strings.Index(cfg.Host, ".")
		if posDot > 0 {
			cfg.Account = cfg.Host[:posDot]
		}
	}
	posDot := strings.Index(cfg.Account, ".")
	if posDot >= 0 {
		cfg.Account = cfg.Account[:posDot]
	}

	err = fillMissingConfigParameters(cfg)
	if err != nil {
		return nil, err
	}

	// unescape parameters
	var s string
	s, err = url.QueryUnescape(cfg.User)
	if err != nil {
		return nil, err
	}
	cfg.User = s
	s, err = url.QueryUnescape(cfg.Password)
	if err != nil {
		return nil, err
	}
	cfg.Password = s
	s, err = url.QueryUnescape(cfg.Database)
	if err != nil {
		return nil, err
	}
	cfg.Database = s
	s, err = url.QueryUnescape(cfg.Schema)
	if err != nil {
		return nil, err
	}
	cfg.Schema = s
	s, err = url.QueryUnescape(cfg.Role)
	if err != nil {
		return nil, err
	}
	cfg.Role = s
	s, err = url.QueryUnescape(cfg.Warehouse)
	if err != nil {
		return nil, err
	}
	cfg.Warehouse = s
	return cfg, nil
}

func fillMissingConfigParameters(cfg *Config) error {
	posDash := strings.LastIndex(cfg.Account, "-")
	if posDash > 0 {
		if strings.Contains(strings.ToLower(cfg.Host), ".global.") {
			cfg.Account = cfg.Account[:posDash]
		}
	}
	if strings.Trim(cfg.Account, " ") == "" {
		return errEmptyAccount()
	}

	if authRequiresUser(cfg) && strings.TrimSpace(cfg.User) == "" {
		return errEmptyUsername()
	}

	if authRequiresPassword(cfg) && strings.TrimSpace(cfg.Password) == "" {
		return errEmptyPassword()
	}
	if strings.Trim(cfg.Protocol, " ") == "" {
		cfg.Protocol = "https"
	}
	if cfg.Port == 0 {
		cfg.Port = 443
	}

	cfg.Region = strings.Trim(cfg.Region, " ")
	if cfg.Region != "" {
		// region is specified but not included in Host
		domain, i := extractDomainFromHost(cfg.Host)
		if i >= 1 {
			hostPrefix := cfg.Host[0:i]
			if !strings.HasSuffix(hostPrefix, cfg.Region) {
				cfg.Host = fmt.Sprintf("%v.%v%v", hostPrefix, cfg.Region, domain)
			}
		}
	}
	if cfg.Host == "" {
		if cfg.Region != "" {
			cfg.Host = cfg.Account + "." + cfg.Region + getDomainBasedOnRegion(cfg.Region)
		} else {
			region, _ := extractRegionFromAccount(cfg.Account)
			if region != "" {
				cfg.Host = cfg.Account + getDomainBasedOnRegion(region)
			} else {
				cfg.Host = cfg.Account + defaultDomain
			}
		}
	}
	if cfg.LoginTimeout == 0 {
		cfg.LoginTimeout = defaultLoginTimeout
	}
	if cfg.RequestTimeout == 0 {
		cfg.RequestTimeout = defaultRequestTimeout
	}
	if cfg.JWTExpireTimeout == 0 {
		cfg.JWTExpireTimeout = defaultJWTTimeout
	}
	if cfg.ClientTimeout == 0 {
		cfg.ClientTimeout = defaultClientTimeout
	}
	if cfg.JWTClientTimeout == 0 {
		cfg.JWTClientTimeout = defaultJWTClientTimeout
	}
	if cfg.ExternalBrowserTimeout == 0 {
		cfg.ExternalBrowserTimeout = defaultExternalBrowserTimeout
	}
	if cfg.MaxRetryCount == 0 {
		cfg.MaxRetryCount = defaultMaxRetryCount
	}
	if strings.Trim(cfg.Application, " ") == "" {
		cfg.Application = clientType
	}

	if cfg.OCSPFailOpen == ocspFailOpenNotSet {
		cfg.OCSPFailOpen = OCSPFailOpenTrue
	}

	if cfg.ValidateDefaultParameters == configBoolNotSet {
		cfg.ValidateDefaultParameters = ConfigBoolTrue
	}

	if cfg.IncludeRetryReason == configBoolNotSet {
		cfg.IncludeRetryReason = ConfigBoolTrue
	}

	domain, _ := extractDomainFromHost(cfg.Host)
	if len(cfg.Host) == len(domain) {
		return &SnowflakeError{
			Number:      ErrCodeFailedToParseHost,
			Message:     errMsgFailedToParseHost,
			MessageArgs: []interface{}{cfg.Host},
		}
	}
	return nil
}

func extractDomainFromHost(host string) (domain string, index int) {
	i := strings.LastIndex(strings.ToLower(host), topLevelDomainPrefix)
	if i >= 1 {
		domain = host[i:]
		return domain, i
	}
	return "", i
}

func getDomainBasedOnRegion(region string) string {
	if strings.HasPrefix(strings.ToLower(region), "cn-") {
		return cnDomain
	}
	return defaultDomain
}

func extractRegionFromAccount(account string) (region string, posDot int) {
	posDot = strings.Index(strings.ToLower(account), ".")
	if posDot > 0 {
		return account[posDot+1:], posDot
	}
	return "", posDot
}

func hostIncludesTopLevelDomain(host string) bool {
	return strings.Contains(strings.ToLower(host), topLevelDomainPrefix)
}

func buildHostFromAccountAndRegion(account, region string) string {
	return account + "." + region + getDomainBasedOnRegion(region)
}

func authRequiresUser(cfg *Config) bool {
	return cfg.Authenticator != AuthTypeOAuth &&
		cfg.Authenticator != AuthTypeTokenAccessor &&
		cfg.Authenticator != AuthTypeExternalBrowser
}

func authRequiresPassword(cfg *Config) bool {
	return cfg.Authenticator != AuthTypeOAuth &&
		cfg.Authenticator != AuthTypeTokenAccessor &&
		cfg.Authenticator != AuthTypeExternalBrowser &&
		cfg.Authenticator != AuthTypeJwt
}

// transformAccountToHost transforms account to host
func transformAccountToHost(cfg *Config) (err error) {
	if cfg.Port == 0 && cfg.Host != "" && !hostIncludesTopLevelDomain(cfg.Host) {
		// account name is specified instead of host:port
		cfg.Account = cfg.Host
		region, posDot := extractRegionFromAccount(cfg.Account)
		if region != "" {
			cfg.Region = region
			cfg.Account = cfg.Account[:posDot]
			cfg.Host = buildHostFromAccountAndRegion(cfg.Account, cfg.Region)
		} else {
			cfg.Host = cfg.Account + defaultDomain
		}
		cfg.Port = 443
	}
	return nil
}

// parseAccountHostPort parses the DSN string to attempt to get account or host and port.
func parseAccountHostPort(cfg *Config, posAt, posSlash int, dsn string) (err error) {
	// account or host:port
	var k int
	for k = posAt + 1; k < posSlash; k++ {
		if dsn[k] == ':' {
			cfg.Port, err = strconv.Atoi(dsn[k+1 : posSlash])
			if err != nil {
				err = &SnowflakeError{
					Number:      ErrCodeFailedToParsePort,
					Message:     errMsgFailedToParsePort,
					MessageArgs: []interface{}{dsn[k+1 : posSlash]},
				}
				return
			}
			break
		}
	}
	cfg.Host = dsn[posAt+1 : k]
	return transformAccountToHost(cfg)
}

// parseUserPassword parses the DSN string for username and password
func parseUserPassword(posAt int, dsn string) (user, password string) {
	var k int
	for k = 0; k < posAt; k++ {
		if dsn[k] == ':' {
			password = dsn[k+1 : posAt]
			break
		}
	}
	user = dsn[:k]
	return
}

// parseParams parse parameters
func parseParams(cfg *Config, posQuestion int, dsn string) (err error) {
	for j := posQuestion + 1; j < len(dsn); j++ {
		if dsn[j] == '?' {
			if err = parseDSNParams(cfg, dsn[j+1:]); err != nil {
				return
			}
			break
		}
	}
	return
}

// parseDSNParams parses the DSN "query string". Values must be url.QueryEscape'ed
func parseDSNParams(cfg *Config, params string) (err error) {
	logger.Infof("Query String: %v\n", params)
	for _, v := range strings.Split(params, "&") {
		param := strings.SplitN(v, "=", 2)
		if len(param) != 2 {
			continue
		}
		var value string
		value, err = url.QueryUnescape(param[1])
		if err != nil {
			return err
		}
		switch param[0] {
		// Disable INFILE whitelist / enable all files
		case "account":
			cfg.Account = value
		case "warehouse":
			cfg.Warehouse = value
		case "database":
			cfg.Database = value
		case "schema":
			cfg.Schema = value
		case "role":
			cfg.Role = value
		case "region":
			cfg.Region = value
		case "protocol":
			cfg.Protocol = value
		case "passcode":
			cfg.Passcode = value
		case "passcodeInPassword":
			var vv bool
			vv, err = strconv.ParseBool(value)
			if err != nil {
				return
			}
			cfg.PasscodeInPassword = vv
		case "clientTimeout":
			cfg.ClientTimeout, err = parseTimeout(value)
			if err != nil {
				return
			}
		case "jwtClientTimeout":
			cfg.JWTClientTimeout, err = parseTimeout(value)
			if err != nil {
				return
			}
		case "loginTimeout":
			cfg.LoginTimeout, err = parseTimeout(value)
			if err != nil {
				return
			}
		case "requestTimeout":
			cfg.RequestTimeout, err = parseTimeout(value)
			if err != nil {
				return
			}
		case "jwtTimeout":
			cfg.JWTExpireTimeout, err = parseTimeout(value)
			if err != nil {
				return err
			}
		case "externalBrowserTimeout":
			cfg.ExternalBrowserTimeout, err = parseTimeout(value)
			if err != nil {
				return err
			}
		case "maxRetryCount":
			cfg.MaxRetryCount, err = strconv.Atoi(value)
			if err != nil {
				return err
			}
		case "application":
			cfg.Application = value
		case "authenticator":
			err := determineAuthenticatorType(cfg, value)
			if err != nil {
				return err
			}
		case "insecureMode":
			var vv bool
			vv, err = strconv.ParseBool(value)
			if err != nil {
				return
			}
			cfg.InsecureMode = vv
		case "ocspFailOpen":
			var vv bool
			vv, err = strconv.ParseBool(value)
			if err != nil {
				return
			}
			if vv {
				cfg.OCSPFailOpen = OCSPFailOpenTrue
			} else {
				cfg.OCSPFailOpen = OCSPFailOpenFalse
			}

		case "token":
			cfg.Token = value
		case "privateKey":
			var decodeErr error
			block, decodeErr := base64.URLEncoding.DecodeString(value)
			if decodeErr != nil {
				err = &SnowflakeError{
					Number:  ErrCodePrivateKeyParseError,
					Message: "Base64 decode failed",
				}
				return
			}
			cfg.PrivateKey, err = parsePKCS8PrivateKey(block)
			if err != nil {
				return err
			}
		case "validateDefaultParameters":
			var vv bool
			vv, err = strconv.ParseBool(value)
			if err != nil {
				return
			}
			if vv {
				cfg.ValidateDefaultParameters = ConfigBoolTrue
			} else {
				cfg.ValidateDefaultParameters = ConfigBoolFalse
			}
		case "clientRequestMfaToken":
			var vv bool
			vv, err = strconv.ParseBool(value)
			if err != nil {
				return
			}
			if vv {
				cfg.ClientRequestMfaToken = ConfigBoolTrue
			} else {
				cfg.ClientRequestMfaToken = ConfigBoolFalse
			}
		case "clientStoreTemporaryCredential":
			var vv bool
			vv, err = strconv.ParseBool(value)
			if err != nil {
				return
			}
			if vv {
				cfg.ClientStoreTemporaryCredential = ConfigBoolTrue
			} else {
				cfg.ClientStoreTemporaryCredential = ConfigBoolFalse
			}
		case "tracing":
			cfg.Tracing = value
		case "tmpDirPath":
			cfg.TmpDirPath = value
		case "disableQueryContextCache":
			var b bool
			b, err = strconv.ParseBool(value)
			if err != nil {
				return
			}
			cfg.DisableQueryContextCache = b
		case "includeRetryReason":
			var vv bool
			vv, err = strconv.ParseBool(value)
			if err != nil {
				return
			}
			if vv {
				cfg.IncludeRetryReason = ConfigBoolTrue
			} else {
				cfg.IncludeRetryReason = ConfigBoolFalse
			}
		case "clientConfigFile":
			cfg.ClientConfigFile = value
		case "disableConsoleLogin":
			var vv bool
			vv, err = strconv.ParseBool(value)
			if err != nil {
				return
			}
			if vv {
				cfg.DisableConsoleLogin = ConfigBoolTrue
			} else {
				cfg.DisableConsoleLogin = ConfigBoolFalse
			}
		case "disableSamlURLCheck":
			var vv bool
			vv, err = strconv.ParseBool(value)
			if err != nil {
				return
			}
			if vv {
				cfg.DisableSamlURLCheck = ConfigBoolTrue
			} else {
				cfg.DisableSamlURLCheck = ConfigBoolFalse
			}
		default:
			if cfg.Params == nil {
				cfg.Params = make(map[string]*string)
			}
			// handle session variables $variable=value
			cfg.Params[urlDecodeIfNeeded(param[0])] = &value
		}
	}
	return
}

func parseTimeout(value string) (time.Duration, error) {
	var vv int64
	var err error
	vv, err = strconv.ParseInt(value, 10, 64)
	if err != nil {
		return time.Duration(0), err
	}
	return time.Duration(vv * int64(time.Second)), nil
}

// ConfigParam is used to bind the name of the Config field with the environment variable and set the requirement for it
type ConfigParam struct {
	Name          string
	EnvName       string
	FailOnMissing bool
}

// GetConfigFromEnv is used to parse the environment variable values to specific fields of the Config
func GetConfigFromEnv(properties []*ConfigParam) (*Config, error) {
	var account, user, password, role, host, portStr, protocol, warehouse, database, schema, region, passcode, application string
	var privateKey *rsa.PrivateKey
	var err error
	if len(properties) == 0 || properties == nil {
		return nil, errors.New("missing configuration parameters for the connection")
	}
	for _, prop := range properties {
		value, err := GetFromEnv(prop.EnvName, prop.FailOnMissing)
		if err != nil {
			return nil, err
		}
		switch prop.Name {
		case "Account":
			account = value
		case "User":
			user = value
		case "Password":
			password = value
		case "Role":
			role = value
		case "Host":
			host = value
		case "Port":
			portStr = value
		case "Protocol":
			protocol = value
		case "Warehouse":
			warehouse = value
		case "Database":
			database = value
		case "Region":
			region = value
		case "Passcode":
			passcode = value
		case "Schema":
			schema = value
		case "Application":
			application = value
		case "PrivateKey":
			privateKey, err = parsePrivateKeyFromFile(value)
			if err != nil {
				return nil, err
			}
		}
	}

	port := 443 // snowflake default port
	if len(portStr) > 0 {
		port, err = strconv.Atoi(portStr)
		if err != nil {
			return nil, err
		}
	}

	cfg := &Config{
		Account:     account,
		User:        user,
		Password:    password,
		Role:        role,
		Host:        host,
		Port:        port,
		Protocol:    protocol,
		Warehouse:   warehouse,
		Database:    database,
		Schema:      schema,
		PrivateKey:  privateKey,
		Region:      region,
		Passcode:    passcode,
		Application: application,
	}
	return cfg, nil
}

func parsePrivateKeyFromFile(path string) (*rsa.PrivateKey, error) {
	bytes, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	block, _ := pem.Decode(bytes)
	if block == nil {
		return nil, errors.New("failed to parse PEM block containing the private key")
	}
	privateKey, err := x509.ParsePKCS8PrivateKey(block.Bytes)
	if err != nil {
		return nil, err
	}
	pk, ok := privateKey.(*rsa.PrivateKey)
	if !ok {
		return nil, fmt.Errorf("interface convertion. expected type *rsa.PrivateKey, but got %T", privateKey)
	}
	return pk, nil
}

func extractAccountName(rawAccount string) string {
	posDot := strings.Index(rawAccount, ".")
	if posDot > 0 {
		return strings.ToUpper(rawAccount[:posDot])
	}
	return strings.ToUpper(rawAccount)
}

func urlDecodeIfNeeded(param string) (decodedParam string) {
	unescaped, err := url.QueryUnescape(param)
	if err != nil {
		return param
	}
	return unescaped
}
