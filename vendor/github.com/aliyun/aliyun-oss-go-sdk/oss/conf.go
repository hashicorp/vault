package oss

import (
	"bytes"
	"fmt"
	"log"
	"net"
	"os"
	"time"
)

// Define the level of the output log
const (
	LogOff = iota
	Error
	Warn
	Info
	Debug
)

// LogTag Tag for each level of log
var LogTag = []string{"[error]", "[warn]", "[info]", "[debug]"}

// HTTPTimeout defines HTTP timeout.
type HTTPTimeout struct {
	ConnectTimeout   time.Duration
	ReadWriteTimeout time.Duration
	HeaderTimeout    time.Duration
	LongTimeout      time.Duration
	IdleConnTimeout  time.Duration
}

// HTTPMaxConns defines max idle connections and max idle connections per host
type HTTPMaxConns struct {
	MaxIdleConns        int
	MaxIdleConnsPerHost int
	MaxConnsPerHost     int
}

// Credentials is interface for get AccessKeyID,AccessKeySecret,SecurityToken
type Credentials interface {
	GetAccessKeyID() string
	GetAccessKeySecret() string
	GetSecurityToken() string
}

// CredentialsProvider is interface for get Credential Info
type CredentialsProvider interface {
	GetCredentials() Credentials
}

type CredentialsProviderE interface {
	CredentialsProvider
	GetCredentialsE() (Credentials, error)
}

type defaultCredentials struct {
	config *Config
}

func (defCre *defaultCredentials) GetAccessKeyID() string {
	return defCre.config.AccessKeyID
}

func (defCre *defaultCredentials) GetAccessKeySecret() string {
	return defCre.config.AccessKeySecret
}

func (defCre *defaultCredentials) GetSecurityToken() string {
	return defCre.config.SecurityToken
}

type defaultCredentialsProvider struct {
	config *Config
}

func (defBuild *defaultCredentialsProvider) GetCredentials() Credentials {
	return &defaultCredentials{config: defBuild.config}
}

type envCredentials struct {
	AccessKeyId     string
	AccessKeySecret string
	SecurityToken   string
}

type EnvironmentVariableCredentialsProvider struct {
	cred Credentials
}

func (credentials *envCredentials) GetAccessKeyID() string {
	return credentials.AccessKeyId
}

func (credentials *envCredentials) GetAccessKeySecret() string {
	return credentials.AccessKeySecret
}

func (credentials *envCredentials) GetSecurityToken() string {
	return credentials.SecurityToken
}

func (defBuild *EnvironmentVariableCredentialsProvider) GetCredentials() Credentials {
	var accessID, accessKey, token string
	if defBuild.cred == nil {
		accessID = os.Getenv("OSS_ACCESS_KEY_ID")
		accessKey = os.Getenv("OSS_ACCESS_KEY_SECRET")
		token = os.Getenv("OSS_SESSION_TOKEN")
	} else {
		accessID = defBuild.cred.GetAccessKeyID()
		accessKey = defBuild.cred.GetAccessKeySecret()
		token = defBuild.cred.GetSecurityToken()
	}

	return &envCredentials{
		AccessKeyId:     accessID,
		AccessKeySecret: accessKey,
		SecurityToken:   token,
	}
}

func NewEnvironmentVariableCredentialsProvider() (EnvironmentVariableCredentialsProvider, error) {
	var provider EnvironmentVariableCredentialsProvider
	accessID := os.Getenv("OSS_ACCESS_KEY_ID")
	if accessID == "" {
		return provider, fmt.Errorf("access key id is empty!")
	}
	accessKey := os.Getenv("OSS_ACCESS_KEY_SECRET")
	if accessKey == "" {
		return provider, fmt.Errorf("access key secret is empty!")
	}
	token := os.Getenv("OSS_SESSION_TOKEN")
	envCredential := &envCredentials{
		AccessKeyId:     accessID,
		AccessKeySecret: accessKey,
		SecurityToken:   token,
	}
	return EnvironmentVariableCredentialsProvider{
		cred: envCredential,
	}, nil
}

// Config defines oss configuration
type Config struct {
	Endpoint            string              // OSS endpoint
	AccessKeyID         string              // AccessId
	AccessKeySecret     string              // AccessKey
	RetryTimes          uint                // Retry count by default it's 5.
	UserAgent           string              // SDK name/version/system information
	IsDebug             bool                // Enable debug mode. Default is false.
	Timeout             uint                // Timeout in seconds. By default it's 60.
	SecurityToken       string              // STS Token
	IsCname             bool                // If cname is in the endpoint.
	IsPathStyle         bool                // If Path Style is in the endpoint.
	HTTPTimeout         HTTPTimeout         // HTTP timeout
	HTTPMaxConns        HTTPMaxConns        // Http max connections
	IsUseProxy          bool                // Flag of using proxy.
	ProxyHost           string              // Flag of using proxy host.
	IsAuthProxy         bool                // Flag of needing authentication.
	ProxyUser           string              // Proxy user
	ProxyPassword       string              // Proxy password
	IsEnableMD5         bool                // Flag of enabling MD5 for upload.
	MD5Threshold        int64               // Memory footprint threshold for each MD5 computation (16MB is the default), in byte. When the data is more than that, temp file is used.
	IsEnableCRC         bool                // Flag of enabling CRC for upload.
	LogLevel            int                 // Log level
	Logger              *log.Logger         // For write log
	UploadLimitSpeed    int                 // Upload limit speed:KB/s, 0 is unlimited
	UploadLimiter       *OssLimiter         // Bandwidth limit reader for upload
	DownloadLimitSpeed  int                 // Download limit speed:KB/s, 0 is unlimited
	DownloadLimiter     *OssLimiter         // Bandwidth limit reader for download
	CredentialsProvider CredentialsProvider // User provides interface to get AccessKeyID, AccessKeySecret, SecurityToken
	LocalAddr           net.Addr            // local client host info
	UserSetUa           bool                // UserAgent is set by user or not
	AuthVersion         AuthVersionType     //  v1 or v2, v4 signature,default is v1
	AdditionalHeaders   []string            //  special http headers needed to be sign
	RedirectEnabled     bool                //  only effective from go1.7 onward, enable http redirect or not
	InsecureSkipVerify  bool                //  for https, Whether to skip verifying the server certificate file
	Region              string              //  such as cn-hangzhou
	CloudBoxId          string              //
	Product             string              //  oss or oss-cloudbox, default is oss
	VerifyObjectStrict  bool                //  a flag of verifying object name strictly. Default is enable.
}

// LimitUploadSpeed uploadSpeed:KB/s, 0 is unlimited,default is 0
func (config *Config) LimitUploadSpeed(uploadSpeed int) error {
	if uploadSpeed < 0 {
		return fmt.Errorf("invalid argument, the value of uploadSpeed is less than 0")
	} else if uploadSpeed == 0 {
		config.UploadLimitSpeed = 0
		config.UploadLimiter = nil
		return nil
	}

	var err error
	config.UploadLimiter, err = GetOssLimiter(uploadSpeed)
	if err == nil {
		config.UploadLimitSpeed = uploadSpeed
	}
	return err
}

// LimitDownLoadSpeed downloadSpeed:KB/s, 0 is unlimited,default is 0
func (config *Config) LimitDownloadSpeed(downloadSpeed int) error {
	if downloadSpeed < 0 {
		return fmt.Errorf("invalid argument, the value of downloadSpeed is less than 0")
	} else if downloadSpeed == 0 {
		config.DownloadLimitSpeed = 0
		config.DownloadLimiter = nil
		return nil
	}

	var err error
	config.DownloadLimiter, err = GetOssLimiter(downloadSpeed)
	if err == nil {
		config.DownloadLimitSpeed = downloadSpeed
	}
	return err
}

// WriteLog output log function
func (config *Config) WriteLog(LogLevel int, format string, a ...interface{}) {
	if config.LogLevel < LogLevel || config.Logger == nil {
		return
	}

	var logBuffer bytes.Buffer
	logBuffer.WriteString(LogTag[LogLevel-1])
	logBuffer.WriteString(fmt.Sprintf(format, a...))
	config.Logger.Printf("%s", logBuffer.String())
}

// for get Credentials
func (config *Config) GetCredentials() Credentials {
	return config.CredentialsProvider.GetCredentials()
}

// for get Sign Product
func (config *Config) GetSignProduct() string {
	if config.CloudBoxId != "" {
		return "oss-cloudbox"
	}
	return "oss"
}

// for get Sign Region
func (config *Config) GetSignRegion() string {
	if config.CloudBoxId != "" {
		return config.CloudBoxId
	}
	return config.Region
}

// getDefaultOssConfig gets the default configuration.
func getDefaultOssConfig() *Config {
	config := Config{}

	config.Endpoint = ""
	config.AccessKeyID = ""
	config.AccessKeySecret = ""
	config.RetryTimes = 5
	config.IsDebug = false
	config.UserAgent = userAgent()
	config.Timeout = 60 // Seconds
	config.SecurityToken = ""
	config.IsCname = false
	config.IsPathStyle = false

	config.HTTPTimeout.ConnectTimeout = time.Second * 30   // 30s
	config.HTTPTimeout.ReadWriteTimeout = time.Second * 60 // 60s
	config.HTTPTimeout.HeaderTimeout = time.Second * 60    // 60s
	config.HTTPTimeout.LongTimeout = time.Second * 300     // 300s
	config.HTTPTimeout.IdleConnTimeout = time.Second * 50  // 50s
	config.HTTPMaxConns.MaxIdleConns = 100
	config.HTTPMaxConns.MaxIdleConnsPerHost = 100

	config.IsUseProxy = false
	config.ProxyHost = ""
	config.IsAuthProxy = false
	config.ProxyUser = ""
	config.ProxyPassword = ""

	config.MD5Threshold = 16 * 1024 * 1024 // 16MB
	config.IsEnableMD5 = false
	config.IsEnableCRC = true

	config.LogLevel = LogOff
	config.Logger = log.New(os.Stdout, "", log.LstdFlags)

	provider := &defaultCredentialsProvider{config: &config}
	config.CredentialsProvider = provider

	config.AuthVersion = AuthV1
	config.RedirectEnabled = true
	config.InsecureSkipVerify = false

	config.Product = "oss"

	config.VerifyObjectStrict = true

	return &config
}
