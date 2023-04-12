//go:build !enterprise

package eventbus

// EnabledPlugins are plugins that are always allowed to send events.
// This is only visible for testing purposes.
var EnabledPlugins = map[string]struct{}{}
