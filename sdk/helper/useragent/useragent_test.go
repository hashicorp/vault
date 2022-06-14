package useragent

import (
	"testing"

	"github.com/hashicorp/vault/sdk/logical"
)

func TestUserAgent(t *testing.T) {
	projectURL = "https://vault-test.com"
	rt = "go5.0"
	versionFunc = func() string { return "1.2.3" }

	type args struct {
		comments []string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "User agent",
			args: args{},
			want: "Vault/1.2.3 (+https://vault-test.com; go5.0)",
		},
		{
			name: "User agent with additional comment",
			args: args{
				comments: []string{"pid-abcdefg"},
			},
			want: "Vault/1.2.3 (+https://vault-test.com; go5.0; pid-abcdefg)",
		},
		{
			name: "User agent with additional comments",
			args: args{
				comments: []string{"pid-abcdefg", "cloud-provider"},
			},
			want: "Vault/1.2.3 (+https://vault-test.com; go5.0; pid-abcdefg; cloud-provider)",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := String(tt.args.comments...); got != tt.want {
				t.Errorf("String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUserAgentPlugin(t *testing.T) {
	projectURL = "https://vault-test.com"
	rt = "go5.0"
	env := &logical.PluginEnvironment{
		VaultVersion: "1.2.3",
	}

	type args struct {
		pluginName string
		comments   []string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "Plugin user agent without plugin name",
			args: args{},
			want: "Vault/1.2.3 (+https://vault-test.com; go5.0)",
		},
		{
			name: "Plugin user agent with plugin name",
			args: args{
				pluginName: "azure-auth",
			},
			want: "Vault/1.2.3 (+https://vault-test.com; azure-auth; go5.0)",
		},
		{
			name: "Plugin user agent with plugin name and additional comment",
			args: args{
				pluginName: "azure-auth",
				comments:   []string{"pid-abcdefg"},
			},
			want: "Vault/1.2.3 (+https://vault-test.com; azure-auth; go5.0; pid-abcdefg)",
		},
		{
			name: "Plugin user agent with plugin name and additional comments",
			args: args{
				pluginName: "azure-auth",
				comments:   []string{"pid-abcdefg", "cloud-provider"},
			},
			want: "Vault/1.2.3 (+https://vault-test.com; azure-auth; go5.0; pid-abcdefg; cloud-provider)",
		},
		{
			name: "Plugin user agent with no plugin name and additional comments",
			args: args{
				comments: []string{"pid-abcdefg", "cloud-provider"},
			},
			want: "Vault/1.2.3 (+https://vault-test.com; go5.0; pid-abcdefg; cloud-provider)",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := PluginString(env, tt.args.pluginName, tt.args.comments...); got != tt.want {
				t.Errorf("PluginString() = %v, want %v", got, tt.want)
			}
		})
	}
}
