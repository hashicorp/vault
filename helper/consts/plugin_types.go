package consts

import "fmt"

var PluginTypes = []PluginType{PluginTypeCredential, PluginTypeSecrets, PluginTypeDatabase}

type PluginType int

const (
	PluginTypeUnknown PluginType = iota
	PluginTypeCredential
	PluginTypeSecrets
	PluginTypeDatabase
)

func (p PluginType) String() string {
	switch p {
	case PluginTypeSecrets:
		return "secret"
	case PluginTypeCredential:
		return "auth"
	case PluginTypeDatabase:
		return "database"
	default:
		return "unknown"
	}
}

func ParsePluginType(pluginType string) (PluginType, error) {
	switch pluginType {
	case "secret":
		return PluginTypeSecrets, nil
	case "auth":
		return PluginTypeCredential, nil
	case "database":
		return PluginTypeDatabase, nil
	default:
		return PluginTypeUnknown, fmt.Errorf("%s is not a supported plugin type", pluginType)
	}
}
