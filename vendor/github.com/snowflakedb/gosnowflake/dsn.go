// Copyright (c) 2017-2021 Snowflake Computing Inc. All right reserved.

package gosnowflake

import (
	"crypto/rsa"
	"encoding/base64"
	"fmt"
	"net"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

const (
	defaultClientTimeout  = 900 * time.Second // Timeout for network round trip + read out http response
	defaultLoginTimeout   = 60 * time.Second  // Timeout for retry for login EXCLUDING clientTimeout
	defaultRequestTimeout = 0 * time.Second   // Timeout for retry for request EXCLUDING clientTimeout
	defaultJWTTimeout     = 60 * time.Second
	defaultDomain         = ".snowflakecomputing.com"
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

	LoginTimeout     time.Duration // Login retry timeout EXCLUDING network roundtrip and read out http response
	RequestTimeout   time.Duration // request retry timeout EXCLUDING network roundtrip and read out http response
	JWTExpireTimeout time.Duration // JWT expire after timeout
	ClientTimeout    time.Duration // Timeout for network round trip + read out http response

	Application  string           // application name.
	InsecureMode bool             // driver doesn't check certificate revocation status
	OCSPFailOpen OCSPFailOpenMode // OCSP Fail Open

	Token            string        // Token to use for OAuth other forms of token based auth
	TokenAccessor    TokenAccessor // Optional token accessor to use
	KeepSessionAlive bool          // Enables the session to persist even after the connection is closed

	PrivateKey *rsa.PrivateKey // Private key used to sign JWT

	Transporter http.RoundTripper // RoundTripper to intercept HTTP requests and responses

	DisableTelemetry bool // indicates whether to disable telemetry
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
	hasHost := true
	if cfg.Host == "" {
		hasHost = false
		if cfg.Region == "us-west-2" {
			cfg.Region = ""
		}
		if cfg.Region == "" {
			cfg.Host = cfg.Account + defaultDomain
		} else {
			cfg.Host = cfg.Account + "." + cfg.Region + defaultDomain
		}
	}
	// in case account includes region
	posDot := strings.Index(cfg.Account, ".")
	if posDot > 0 {
		if cfg.Region != "" {
			return "", ErrInvalidRegion
		}
		cfg.Region = cfg.Account[posDot+1:]
		cfg.Account = cfg.Account[:posDot]
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
	if cfg.LoginTimeout != defaultLoginTimeout {
		params.Add("loginTimeout", strconv.FormatInt(int64(cfg.LoginTimeout/time.Second), 10))
	}
	if cfg.RequestTimeout != defaultRequestTimeout {
		params.Add("requestTimeout", strconv.FormatInt(int64(cfg.RequestTimeout/time.Second), 10))
	}
	if cfg.JWTExpireTimeout != defaultJWTTimeout {
		params.Add("jwtTimeout", strconv.FormatInt(int64(cfg.JWTExpireTimeout/time.Second), 10))
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

	params.Add("ocspFailOpen", strconv.FormatBool(cfg.OCSPFailOpen != OCSPFailOpenFalse))

	params.Add("validateDefaultParameters", strconv.FormatBool(cfg.ValidateDefaultParameters != ConfigBoolFalse))

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
	if cfg.Account == "" && strings.HasSuffix(cfg.Host, defaultDomain) {
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
		if strings.Contains(cfg.Host, ".global.") {
			cfg.Account = cfg.Account[:posDash]
		}
	}
	if strings.Trim(cfg.Account, " ") == "" {
		return ErrEmptyAccount
	}

	if cfg.Authenticator != AuthTypeOAuth && strings.Trim(cfg.User, " ") == "" {
		// oauth does not require a username
		return ErrEmptyUsername
	}

	if cfg.Authenticator != AuthTypeExternalBrowser &&
		cfg.Authenticator != AuthTypeOAuth &&
		cfg.Authenticator != AuthTypeJwt &&
		strings.Trim(cfg.Password, " ") == "" {
		// no password parameter is required for EXTERNALBROWSER, OAUTH or JWT.
		return ErrEmptyPassword
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
		i := strings.Index(cfg.Host, defaultDomain)
		if i >= 1 {
			hostPrefix := cfg.Host[0:i]
			if !strings.HasSuffix(hostPrefix, cfg.Region) {
				cfg.Host = hostPrefix + "." + cfg.Region + defaultDomain
			}
		}
	}
	if cfg.Host == "" {
		if cfg.Region != "" {
			cfg.Host = cfg.Account + "." + cfg.Region + defaultDomain
		} else {
			cfg.Host = cfg.Account + defaultDomain
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
	if strings.Trim(cfg.Application, " ") == "" {
		cfg.Application = clientType
	}

	if cfg.OCSPFailOpen == ocspFailOpenNotSet {
		cfg.OCSPFailOpen = OCSPFailOpenTrue
	}

	if cfg.ValidateDefaultParameters == configBoolNotSet {
		cfg.ValidateDefaultParameters = ConfigBoolTrue
	}

	if strings.HasSuffix(cfg.Host, defaultDomain) && len(cfg.Host) == len(defaultDomain) {
		return &SnowflakeError{
			Number:      ErrCodeFailedToParseHost,
			Message:     errMsgFailedToParseHost,
			MessageArgs: []interface{}{cfg.Host},
		}
	}
	return nil
}

// transformAccountToHost transforms host to account name
func transformAccountToHost(cfg *Config) (err error) {
	if cfg.Port == 0 && !strings.HasSuffix(cfg.Host, defaultDomain) && cfg.Host != "" {
		// account name is specified instead of host:port
		cfg.Account = cfg.Host
		cfg.Host = cfg.Account + defaultDomain
		cfg.Port = 443
		posDot := strings.Index(cfg.Account, ".")
		if posDot > 0 {
			cfg.Region = cfg.Account[posDot+1:]
			cfg.Account = cfg.Account[:posDot]
		}
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
		default:
			if cfg.Params == nil {
				cfg.Params = make(map[string]*string)
			}
			cfg.Params[param[0]] = &value
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
