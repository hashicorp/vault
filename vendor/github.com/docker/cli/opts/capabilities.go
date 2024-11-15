package opts

import (
	"sort"
	"strings"
)

const (
	// AllCapabilities is a special value to add or drop all capabilities
	AllCapabilities = "ALL"

	// ResetCapabilities is a special value to reset capabilities when updating.
	// This value should only be used when updating, not used on "create".
	ResetCapabilities = "RESET"
)

// NormalizeCapability normalizes a capability by upper-casing, trimming white space
// and adding a CAP_ prefix (if not yet present). This function also accepts the
// "ALL" magic-value, as used by CapAdd/CapDrop.
//
// This function only handles rudimentary formatting; no validation is performed,
// as the list of available capabilities can be updated over time, thus should be
// handled by the daemon.
func NormalizeCapability(capability string) string {
	capability = strings.ToUpper(strings.TrimSpace(capability))
	if capability == AllCapabilities || capability == ResetCapabilities {
		return capability
	}
	if !strings.HasPrefix(capability, "CAP_") {
		capability = "CAP_" + capability
	}
	return capability
}

// CapabilitiesMap normalizes the given capabilities and converts them to a map.
func CapabilitiesMap(caps []string) map[string]bool {
	normalized := make(map[string]bool)
	for _, c := range caps {
		normalized[NormalizeCapability(c)] = true
	}
	return normalized
}

// EffectiveCapAddCapDrop normalizes and sorts capabilities to "add" and "drop",
// and returns the effective capabilities to include in both.
//
// "CapAdd" takes precedence over "CapDrop", so capabilities included in both
// lists are removed from the list of capabilities to drop. The special "ALL"
// capability is also taken into account.
//
// Note that the special "RESET" value is only used when updating an existing
// service, and will be ignored.
//
// Duplicates are removed, and the resulting lists are sorted.
func EffectiveCapAddCapDrop(add, drop []string) (capAdd, capDrop []string) {
	var (
		addCaps  = CapabilitiesMap(add)
		dropCaps = CapabilitiesMap(drop)
	)

	if addCaps[AllCapabilities] {
		// Special case: "ALL capabilities" trumps any other capability added.
		addCaps = map[string]bool{AllCapabilities: true}
	}
	if dropCaps[AllCapabilities] {
		// Special case: "ALL capabilities" trumps any other capability added.
		dropCaps = map[string]bool{AllCapabilities: true}
	}
	for c := range dropCaps {
		if addCaps[c] {
			// Adding a capability takes precedence, so skip dropping
			continue
		}
		if c != ResetCapabilities {
			capDrop = append(capDrop, c)
		}
	}

	for c := range addCaps {
		if c != ResetCapabilities {
			capAdd = append(capAdd, c)
		}
	}

	sort.Strings(capAdd)
	sort.Strings(capDrop)

	return capAdd, capDrop
}
