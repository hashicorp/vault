package linodego

import (
	"fmt"
	"os"
	"strings"

	"gopkg.in/ini.v1"
)

const (
	DefaultConfigProfile = "default"
)

var DefaultConfigPaths = []string{
	"%s/.config/linode",
	"%s/.config/linode-cli",
}

type ConfigProfile struct {
	APIToken   string `ini:"token"`
	APIVersion string `ini:"api_version"`
	APIURL     string `ini:"api_url"`
}

type LoadConfigOptions struct {
	Path            string
	Profile         string
	SkipLoadProfile bool
}

// LoadConfig loads a Linode config according to the option's argument.
// If no options are specified, the following defaults will be used:
// Path: ~/.config/linode
// Profile: default
func (c *Client) LoadConfig(options *LoadConfigOptions) error {
	path, err := resolveValidConfigPath()
	if err != nil {
		return err
	}

	profileOption := DefaultConfigProfile

	if options != nil {
		if options.Path != "" {
			path = options.Path
		}

		if options.Profile != "" {
			profileOption = options.Profile
		}
	}

	cfg, err := ini.Load(path)
	if err != nil {
		return err
	}

	defaultConfig := ConfigProfile{
		APIToken:   "",
		APIURL:     APIHost,
		APIVersion: APIVersion,
	}

	if cfg.HasSection("default") {
		err := cfg.Section("default").MapTo(&defaultConfig)
		if err != nil {
			return fmt.Errorf("failed to map default profile: %w", err)
		}
	}

	result := make(map[string]ConfigProfile)

	for _, profile := range cfg.Sections() {
		name := strings.ToLower(profile.Name())

		f := defaultConfig
		if err := profile.MapTo(&f); err != nil {
			return fmt.Errorf("failed to map values: %w", err)
		}

		result[name] = f
	}

	c.configProfiles = result

	if !options.SkipLoadProfile {
		if err := c.UseProfile(profileOption); err != nil {
			return fmt.Errorf("unable to use profile %s: %w", profileOption, err)
		}
	}

	return nil
}

// UseProfile switches client to use the specified profile.
// The specified profile must be already be loaded using client.LoadConfig(...)
func (c *Client) UseProfile(name string) error {
	name = strings.ToLower(name)

	profile, ok := c.configProfiles[name]
	if !ok {
		return fmt.Errorf("profile %s does not exist", name)
	}

	if profile.APIToken == "" {
		return fmt.Errorf("unable to resolve linode_token for profile %s", name)
	}

	if profile.APIURL == "" {
		return fmt.Errorf("unable to resolve linode_api_url for profile %s", name)
	}

	if profile.APIVersion == "" {
		return fmt.Errorf("unable to resolve linode_api_version for profile %s", name)
	}

	c.SetToken(profile.APIToken)
	c.SetBaseURL(profile.APIURL)
	c.SetAPIVersion(profile.APIVersion)
	c.selectedProfile = name
	c.loadedProfile = name

	return nil
}

func FormatConfigPath(path string) (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}

	return fmt.Sprintf(path, homeDir), nil
}

func resolveValidConfigPath() (string, error) {
	for _, cfg := range DefaultConfigPaths {
		p, err := FormatConfigPath(cfg)
		if err != nil {
			return "", err
		}

		if _, err := os.Stat(p); err != nil {
			continue
		}

		return p, err
	}

	// An empty result may not be an error
	return "", nil
}
