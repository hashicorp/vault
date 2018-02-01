package lib

import (
	"github.com/hashicorp/serf/serf"
)

// SerfDefaultConfig returns a Consul-flavored Serf default configuration,
// suitable as a basis for a LAN, WAN, segment, or area.
func SerfDefaultConfig() *serf.Config {
	base := serf.DefaultConfig()

	// This effectively disables the annoying queue depth warnings.
	base.QueueDepthWarning = 1000000

	// This enables dynamic sizing of the message queue depth based on the
	// cluster size.
	base.MinQueueDepth = 4096

	return base
}
