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

// AUTO-GENERATED!  DO NOT EDIT FILE DIRECTLY

package okta

import (
	"context"
	"fmt"
	"io/ioutil"
	"os/user"

	"github.com/okta/okta-sdk-golang/okta/cache"

	"github.com/go-yaml/yaml"
	"github.com/kelseyhightower/envconfig"
)

const Version = "0.1.0"

type Client struct {
	config *config

	requestExecutor *RequestExecutor

	resource resource

	Application *ApplicationResource
	Group       *GroupResource
	LogEvent    *LogEventResource
	Policy      *PolicyResource
	Session     *SessionResource
	User        *UserResource
	Factor      *FactorResource
}

type resource struct {
	client *Client
}

func NewClient(ctx context.Context, conf ...ConfigSetter) (*Client, error) {
	config := &config{}

	setConfigDefaults(config)
	config = readConfigFromSystem(*config)
	config = readConfigFromApplication(*config)
	config = readConfigFromEnvironment(*config)

	for _, confSetter := range conf {
		confSetter(config)
	}

	var oktaCache cache.Cache
	if !config.Okta.Client.Cache.Enabled {
		oktaCache = cache.NewNoOpCache()
	} else {
		if config.CacheManager == nil {
			oktaCache = cache.NewGoCache(config.Okta.Client.Cache.DefaultTtl,
				config.Okta.Client.Cache.DefaultTti)
		} else {
			oktaCache = config.CacheManager
		}
	}

	config.CacheManager = oktaCache

	config, err := validateConfig(config)
	if err != nil {
		panic(err)
	}

	c := &Client{}
	c.config = config
	c.requestExecutor = NewRequestExecutor(&config.HttpClient, oktaCache, config)

	c.resource.client = c

	c.Application = (*ApplicationResource)(&c.resource)
	c.Group = (*GroupResource)(&c.resource)
	c.LogEvent = (*LogEventResource)(&c.resource)
	c.Policy = (*PolicyResource)(&c.resource)
	c.Session = (*SessionResource)(&c.resource)
	c.User = (*UserResource)(&c.resource)
	c.Factor = (*FactorResource)(&c.resource)
	return c, nil
}

func (c *Client) GetConfig() *config {
	return c.config
}

func (c *Client) GetRequestExecutor() *RequestExecutor {
	return c.requestExecutor
}

func setConfigDefaults(c *config) {
	var conf []ConfigSetter

	conf = append(conf,
		WithConnectionTimeout(30),
		WithCache(true),
		WithCacheTtl(300),
		WithCacheTti(300),
		WithUserAgentExtra(""))
	WithTestingDisableHttpsCheck(false)

	for _, confSetter := range conf {
		confSetter(c)
	}
}

func readConfigFromFile(location string) (*config, error) {
	yamlConfig, err := ioutil.ReadFile(location)

	if err != nil {
		return nil, err
	}

	conf := config{}
	err = yaml.Unmarshal(yamlConfig, &conf)
	if err != nil {
		return nil, err
	}

	return &conf, err
}

func readConfigFromSystem(c config) *config {
	currUser, err := user.Current()
	if err != nil {
		return &c
	}
	if currUser.HomeDir == "" {
		return &c
	}

	conf, err := readConfigFromFile(currUser.HomeDir + "/.okta/okta.yaml")

	if err != nil {
		return &c
	}

	return conf
}

func readConfigFromApplication(c config) *config {
	conf, err := readConfigFromFile(".okta.yaml")

	if err != nil {
		return &c
	}

	return conf
}

func readConfigFromEnvironment(c config) *config {
	err := envconfig.Process("okta", &c)
	if err != nil {
		fmt.Println("error parsing")
		return &c
	}
	return &c
}
