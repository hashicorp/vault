package gocbcore

import (
	"bytes"
	"fmt"
)

type routeConfig struct {
	revID          int64
	revEpoch       int64
	uuid           string
	name           string
	bktType        bucketType
	kvServerList   []string
	capiEpList     []string
	mgmtEpList     []string
	n1qlEpList     []string
	ftsEpList      []string
	cbasEpList     []string
	eventingEpList []string
	gsiEpList      []string
	backupEpList   []string
	vbMap          *vbucketMap
	ketamaMap      *ketamaContinuum

	clusterCapabilitiesVer []int
	clusterCapabilities    map[string][]string

	bucketCapabilities    []string
	bucketCapabilitiesVer string
}

func (config *routeConfig) DebugString() string {
	var buffer bytes.Buffer

	buffer.WriteString(fmt.Sprintf("Revision ID: %d\n", config.revID))
	buffer.WriteString(fmt.Sprintf("Revision Epoch: %d\n", config.revEpoch))

	if config.name != "" {
		fmt.Fprintf(&buffer, "Bucket: %s\n", config.name)
	}

	addEps := func(title string, eps []string) {
		fmt.Fprintf(&buffer, "%s Eps:\n", title)
		for _, ep := range eps {
			fmt.Fprintf(&buffer, "  - %s\n", ep)
		}
	}

	addEps("Capi", config.capiEpList)
	addEps("Mgmt", config.mgmtEpList)
	addEps("N1ql", config.n1qlEpList)
	addEps("FTS", config.ftsEpList)
	addEps("CBAS", config.cbasEpList)
	addEps("Eventing", config.eventingEpList)
	addEps("GSI", config.gsiEpList)
	addEps("Backup", config.backupEpList)

	if config.vbMap != nil {
		fmt.Fprintln(&buffer, "VBMap:")
		fmt.Fprintf(&buffer, "%+v\n", config.vbMap)
	} else {
		fmt.Fprintln(&buffer, "VBMap: not-used")
	}

	if config.ketamaMap != nil {
		fmt.Fprintln(&buffer, "KetamaMap:")
		fmt.Fprintf(&buffer, "%+v\n", config.ketamaMap)
	} else {
		fmt.Fprintln(&buffer, "KetamaMap: not-used")
	}

	return buffer.String()
}

func (config *routeConfig) IsValid() bool {
	if len(config.kvServerList) == 0 || len(config.mgmtEpList) == 0 {
		return false
	}
	switch config.bktType {
	case bktTypeCouchbase:
		return config.vbMap != nil && config.vbMap.IsValid()
	case bktTypeMemcached:
		return config.ketamaMap != nil && config.ketamaMap.IsValid()
	case bktTypeNone:
		return true
	default:
		return false
	}
}

func (config *routeConfig) IsGCCCPConfig() bool {
	return config.bktType == bktTypeNone
}

func (config *routeConfig) ContainsClusterCapability(version int, category, capability string) bool {
	caps := config.clusterCapabilities
	capsVer := config.clusterCapabilitiesVer
	if len(capsVer) == 0 || caps == nil {
		return false
	}

	if capsVer[0] == version {
		for cat, catCapabilities := range caps {
			switch cat {
			case category:
				for _, capa := range catCapabilities {
					switch capa {
					case capability:
						return true
					}
				}
			}
		}
	}

	return false
}

func (config *routeConfig) ContainsBucketCapability(needleCap string) bool {
	for _, capa := range config.bucketCapabilities {
		if capa == needleCap {
			return true
		}
	}
	return false
}
