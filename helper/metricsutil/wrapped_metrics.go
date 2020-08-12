package metricsutil

import (
	"strings"

	metrics "github.com/armon/go-metrics"
	"github.com/hashicorp/vault/helper/namespace"
)

// NamespaceLabel creates a metrics label for the given
// Namespace: root is "root"; others are path with the
// final '/' removed.
func NamespaceLabel(ns *namespace.Namespace) metrics.Label {
	switch {
	case ns == nil:
		return metrics.Label{"namespace", "root"}
	case ns.ID == namespace.RootNamespaceID:
		return metrics.Label{"namespace", "root"}
	default:
		return metrics.Label{"namespace",
			strings.Trim(ns.Path, "/")}
	}
}
