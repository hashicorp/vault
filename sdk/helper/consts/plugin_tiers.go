// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package consts

import (
	"encoding/json"
	"fmt"
)

var PluginTiers = []PluginTier{
	PluginTierUnknown,
	PluginTierCommunity,
	PluginTierPartner,
	PluginTierOfficial,
}

type PluginTier uint32

const (
	// PluginTierUnknown defines unknown plugin tier
	// DO NOT change the order of the enum as it
	// could cause the wrong plugin tier to be read
	// from storage for a given underlying number
	PluginTierUnknown PluginTier = iota
	// PluginTierCommunity defines community plugin tier
	// DO NOT change the order of the enum as it
	// could cause the wrong plugin tier to be read
	// from storage for a given underlying number
	PluginTierCommunity
	// PluginTierPartner defines partner plugin tier
	// DO NOT change the order of the enum as it
	// could cause the wrong plugin tier to be read
	// from storage for a given underlying number
	PluginTierPartner
	// PluginTierOfficial defines enterprise plugin tier
	// DO NOT change the order of the enum as it
	// could cause the wrong plugin tier to be read
	// from storage for a given underlying number
	PluginTierOfficial
)

func (p PluginTier) String() string {
	switch p {
	case PluginTierUnknown:
		return "unknown"
	case PluginTierCommunity:
		return "community"
	case PluginTierPartner:
		return "partner"
	case PluginTierOfficial:
		return "official"
	default:
		return "unsupported"
	}
}

func ParsePluginTier(pluginTier string) (PluginTier, error) {
	switch pluginTier {
	case "unknown", "":
		return PluginTierUnknown, nil
	case "community":
		return PluginTierCommunity, nil
	case "partner":
		return PluginTierPartner, nil
	case "official":
		return PluginTierOfficial, nil
	default:
		return PluginTierUnknown, fmt.Errorf("%q is not a supported plugin tier", pluginTier)
	}
}

// UnmarshalJSON implements json.Unmarshaler. It supports unmarshaling either a
// string or a uint32. All new serialization will be as a string, but we
// previously serialized as a uint32 so we need to support that for backwards
// compatibility.
func (p *PluginTier) UnmarshalJSON(data []byte) error {
	var asString string
	err := json.Unmarshal(data, &asString)
	if err == nil {
		*p, err = ParsePluginTier(asString)
		return err
	}

	var asUint32 uint32
	err = json.Unmarshal(data, &asUint32)
	if err != nil {
		return err
	}
	*p = PluginTier(asUint32)
	switch *p {
	case PluginTierUnknown, PluginTierCommunity, PluginTierPartner, PluginTierOfficial:
		return nil
	default:
		return fmt.Errorf("%d is not a supported plugin tier", asUint32)
	}
}

// MarshalJSON implements json.Marshaler.
func (p PluginTier) MarshalJSON() ([]byte, error) {
	return json.Marshal(p.String())
}
