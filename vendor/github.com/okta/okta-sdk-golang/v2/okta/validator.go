package okta

import (
	"errors"
	"strings"
)

func validateConfig(c *config) (*config, error) {
	var err error

	err = validateOktaDomain(c)
	if err != nil {
		return nil, err
	}

	if c.Okta.Client.AuthorizationMode == "SSWS" {
		err = validateApiToken(c)
		if err != nil {
			return nil, err
		}
	}

	err = validateAuthorization(c)
	if err != nil {
		return nil, err
	}

	return c, nil
}

func validateOktaDomain(c *config) error {
	if c.Okta.Client.OrgUrl == "" {
		return errors.New("your Okta URL is missing. You can copy your domain from the Okta Developer Console. Follow these instructions to find it: https://bit.ly/finding-okta-domain")
	}

	if strings.Contains(c.Okta.Client.OrgUrl, "{yourOktaDomain}") {
		return errors.New("replace {yourOktaDomain} with your Okta domain. You can copy your domain from the Okta Developer Console. Follow these instructions to find it: https://bit.ly/finding-okta-domain")
	}

	if strings.Contains(c.Okta.Client.OrgUrl, "-admin.okta.com") ||
		strings.Contains(c.Okta.Client.OrgUrl, "-admin.oktapreview.com") ||
		strings.Contains(c.Okta.Client.OrgUrl, "-admin.okta-emea.com") {
		return errors.New("your Okta domain should not contain -admin. Current value: " + c.Okta.Client.OrgUrl + ". You can copy your domain from the Okta Developer Console. Follow these instructions to find it: https://bit.ly/finding-okta-domain")
	}

	if strings.HasSuffix(c.Okta.Client.OrgUrl, ".com.com") {
		return errors.New("it looks like there's a typo in your Okta domain. Current value: " + c.Okta.Client.OrgUrl + ". You can copy your domain from the Okta Developer Console. Follow these instructions to find it: https://bit.ly/finding-okta-domain")
	}

	if c.Okta.Testing.DisableHttpsCheck == false {
		if strings.HasPrefix(c.Okta.Client.OrgUrl, "https://") != true {
			return errors.New("your Okta URL must start with https. Current value: " + c.Okta.Client.OrgUrl + ". You can copy your domain from the Okta Developer Console. Follow these instructions to find it: https://bit.ly/finding-okta-domain")
		}
	}
	return nil
}

func validateApiToken(c *config) error {
	if c.Okta.Client.Token == "" {
		return errors.New("your Okta API token is missing. You can generate one in the Okta Developer Console. Follow these instructions: https://bit.ly/get-okta-api-token")
	}

	if strings.Contains(c.Okta.Client.Token, "{apiToken}") {
		return errors.New("replace {apiToken} with your Okta API token. You can generate one in the Okta Developer Console. Follow these instructions: https://bit.ly/get-okta-api-token")
	}
	return nil
}

func validateAuthorization(c *config) error {
	if c.Okta.Client.AuthorizationMode != "SSWS" &&
		c.Okta.Client.AuthorizationMode != "PrivateKey" {
		return errors.New("the AuthorizaitonMode config option must be one of [SSWS, PrivateKey]. You provided the SDK with " + c.Okta.Client.AuthorizationMode)
	}

	if c.Okta.Client.AuthorizationMode == "PrivateKey" &&
		(c.Okta.Client.ClientId == "" ||
			c.Okta.Client.Scopes == nil ||
			c.Okta.Client.PrivateKey == "") {
		return errors.New("when using AuthorizationMode 'PrivateKey', you must supply 'ClientId', 'Scopes', and 'PrivateKey'")
	}

	return nil
}
