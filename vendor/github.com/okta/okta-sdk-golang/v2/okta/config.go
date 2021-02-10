/*
 * Copyright 2018 - Present Okta, Inc.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *      http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package okta

import (
	"errors"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"syscall"

	"github.com/okta/okta-sdk-golang/v2/okta/cache"
)

type config struct {
	Okta struct {
		Client struct {
			Cache struct {
				Enabled    bool  `yaml:"enabled" envconfig:"OKTA_CLIENT_CACHE_ENABLED"`
				DefaultTtl int32 `yaml:"defaultTtl" envconfig:"OKTA_CLIENT_CACHE_DEFAULT_TTL"`
				DefaultTti int32 `yaml:"defaultTti" envconfig:"OKTA_CLIENT_CACHE_DEFAULT_TTI"`
			} `yaml:"cache"`
			Proxy struct {
				Port     int32  `yaml:"port" envconfig:"OKTA_CLIENT_PROXY_PORT"`
				Host     string `yaml:"host" envconfig:"OKTA_CLIENT_PROXY_HOST"`
				Username string `yaml:"username" envconfig:"OKTA_CLIENT_PROXY_USERNAME"`
				Password string `yaml:"password" envconfig:"OKTA_CLIENT_PROXY_PASSWORD"`
			} `yaml:"proxy"`
			ConnectionTimeout int64 `yaml:"connectionTimeout" envconfig:"OKTA_CLIENT_CONNECTION_TIMEOUT"`
			RequestTimeout    int64 `yaml:"requestTimeout" envconfig:"OKTA_CLIENT_REQUEST_TIMEOUT"`
			RateLimit         struct {
				MaxRetries int32 `yaml:"maxRetries" envconfig:"OKTA_CLIENT_RATE_LIMIT_MAX_RETRIES"`
				MaxBackoff int64 `yaml:"maxBackoff" envconfig:"OKTA_CLIENT_RATE_LIMIT_MAX_BACKOFF"`
			} `yaml:"rateLimit"`
			OrgUrl            string   `yaml:"orgUrl" envconfig:"OKTA_CLIENT_ORGURL"`
			Token             string   `yaml:"token" envconfig:"OKTA_CLIENT_TOKEN"`
			AuthorizationMode string   `yaml:"authorizationMode" envconfig:"OKTA_CLIENT_AUTHORIZATIONMODE"`
			ClientId          string   `yaml:"clientId" envconfig:"OKTA_CLIENT_CLIENTID"`
			Scopes            []string `yaml:"scopes" envconfig:"OKTA_CLIENT_SCOPES"`
			PrivateKey        string   `yaml:"privateKey" envconfig:"OKTA_CLIENT_PRIVATEKEY"`
		} `yaml:"client"`
		Testing struct {
			DisableHttpsCheck bool `yaml:"disableHttpsCheck" envconfig:"OKTA_TESTING_DISABLE_HTTPS_CHECK"`
		} `yaml:"testing"`
	} `yaml:"okta"`
	UserAgentExtra string
	HttpClient     http.Client
	CacheManager   cache.Cache
}

type ConfigSetter func(*config)

func WithCache(cache bool) ConfigSetter {
	return func(c *config) {
		c.Okta.Client.Cache.Enabled = cache
	}
}

func WithCacheManager(cacheManager cache.Cache) ConfigSetter {
	return func(c *config) {
		c.CacheManager = cacheManager
	}
}

func WithCacheTtl(i int32) ConfigSetter {
	return func(c *config) {
		c.Okta.Client.Cache.DefaultTtl = i
	}
}

func WithCacheTti(i int32) ConfigSetter {
	return func(c *config) {
		c.Okta.Client.Cache.DefaultTti = i
	}
}

func WithConnectionTimeout(i int64) ConfigSetter {
	return func(c *config) {
		c.Okta.Client.ConnectionTimeout = i
	}
}

func WithProxyPort(i int32) ConfigSetter {
	return func(c *config) {
		c.Okta.Client.Proxy.Port = i
	}
}

func WithProxyHost(host string) ConfigSetter {
	return func(c *config) {
		c.Okta.Client.Proxy.Host = host
	}
}

func WithProxyUsername(username string) ConfigSetter {
	return func(c *config) {
		c.Okta.Client.Proxy.Username = username
	}
}

func WithProxyPassword(pass string) ConfigSetter {
	return func(c *config) {
		c.Okta.Client.Proxy.Password = pass
	}
}

func WithOrgUrl(url string) ConfigSetter {
	return func(c *config) {
		c.Okta.Client.OrgUrl = url
	}
}

func WithToken(token string) ConfigSetter {
	return func(c *config) {
		c.Okta.Client.Token = token
	}
}

func WithUserAgentExtra(userAgent string) ConfigSetter {
	return func(c *config) {
		c.UserAgentExtra = userAgent
	}
}

func WithHttpClient(httpClient http.Client) ConfigSetter {
	return func(c *config) {
		c.HttpClient = httpClient
	}
}

func WithTestingDisableHttpsCheck(httpsCheck bool) ConfigSetter {
	return func(c *config) {
		c.Okta.Testing.DisableHttpsCheck = httpsCheck
	}
}

func WithRequestTimeout(requestTimeout int64) ConfigSetter {
	return func(c *config) {
		c.Okta.Client.RequestTimeout = requestTimeout
	}
}

func WithRateLimitMaxRetries(maxRetries int32) ConfigSetter {
	return func(c *config) {
		c.Okta.Client.RateLimit.MaxRetries = maxRetries
	}
}

func WithRateLimitMaxBackOff(maxBackoff int64) ConfigSetter {
	return func(c *config) {
		c.Okta.Client.RateLimit.MaxBackoff = maxBackoff
	}
}

func WithAuthorizationMode(authzMode string) ConfigSetter {
	return func(c *config) {
		c.Okta.Client.AuthorizationMode = authzMode
	}
}

func WithClientId(clientId string) ConfigSetter {
	return func(c *config) {
		c.Okta.Client.ClientId = clientId
	}
}

func WithScopes(scopes []string) ConfigSetter {
	return func(c *config) {
		c.Okta.Client.Scopes = scopes
	}
}

// WithPrivateKey sets private key key. Can be either a path to a private key or private key itself.
func WithPrivateKey(privateKey string) ConfigSetter {
	return func(c *config) {
		if fileExists(privateKey) {
			content, err := ioutil.ReadFile(privateKey)
			if err != nil {
				log.Fatalf("failed to read from provided private key file path: %v", err)
			}
			c.Okta.Client.PrivateKey = string(content)
		} else {
			c.Okta.Client.PrivateKey = privateKey
		}
	}
}

func fileExists(filename string) bool {
	info, err := os.Stat(filename)
	if err != nil {
		if os.IsNotExist(err) || errors.Is(err, syscall.ENAMETOOLONG) {
			return false
		}
		log.Println("can not get information about the file containing private key, using provided value as the key itself")
		return false
	}
	return !info.IsDir()
}
