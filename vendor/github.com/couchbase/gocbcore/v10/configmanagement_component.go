package gocbcore

import (
	"encoding/binary"
	"fmt"
	"math/rand"
	"sync"

	"github.com/couchbase/gocbcore/v10/memd"
)

type configManagementComponent struct {
	useSSL      bool
	networkType string

	seedNodeAddr      string
	localLoopbackAddr *localLoopbackAddress

	currentConfig *routeConfig
	configLock    sync.Mutex

	cfgChangeWatchers []routeConfigWatcher
	watchersLock      sync.Mutex

	srcServers []routeEndpoint

	seenConfig bool

	configFetcher      *cccpConfigFetcher
	configFetchSig     chan struct{}
	configFetchSigLock sync.Mutex

	shutdownSig chan struct{}
}

type configManagerProperties struct {
	UseTLS       bool
	SeedNodeAddr string
	NetworkType  string
	SrcMemdAddrs []routeEndpoint
	SrcHTTPAddrs []routeEndpoint
}

type routeConfigWatcher interface {
	OnNewRouteConfig(cfg *routeConfig)
}

type configManager interface {
	AddConfigWatcher(watcher routeConfigWatcher)
	RemoveConfigWatcher(watcher routeConfigWatcher)
}

func newConfigManager(props configManagerProperties) *configManagementComponent {
	return &configManagementComponent{
		useSSL:       props.UseTLS,
		seedNodeAddr: props.SeedNodeAddr,
		networkType:  props.NetworkType,
		srcServers:   append(props.SrcMemdAddrs, props.SrcHTTPAddrs...),
		currentConfig: &routeConfig{
			revID: -1,
		},
		shutdownSig: make(chan struct{}),
	}
}

// SetConfigFetcher sets the config fetcher for the manager, this must be done before OnNewConfig can be called.
func (cm *configManagementComponent) SetConfigFetcher(fetcher *cccpConfigFetcher) {
	cm.configFetcher = fetcher
}

func (cm *configManagementComponent) UseTLS(use bool) {
	cm.configLock.Lock()
	cm.useSSL = use
	cm.configLock.Unlock()
}

func (cm *configManagementComponent) TLSEnabled() bool {
	cm.configLock.Lock()
	useSSL := cm.useSSL
	cm.configLock.Unlock()

	return useSSL
}

func (cm *configManagementComponent) CurrentRev() (int64, int64) {
	cm.configLock.Lock()
	revID := cm.currentConfig.revID
	revEpoch := cm.currentConfig.revEpoch
	cm.configLock.Unlock()

	return revID, revEpoch
}

func (cm *configManagementComponent) OnNewConfig(cfg *cfgBucket) {
	cm.onNewConfig(cfg)
}

func (cm *configManagementComponent) onNewConfig(cfg *cfgBucket) bool {
	var routeCfg *routeConfig
	cm.configLock.Lock()
	if cm.seenConfig {
		routeCfg = cfg.BuildRouteConfig(cm.useSSL, cm.networkType, false, cm.localLoopbackAddr)
	} else {
		routeCfg = cm.buildFirstRouteConfig(cfg, cm.useSSL)
		if routeCfg == nil {
			cm.configLock.Unlock()
			// If the routeCfg isn't valid then ignore it.
			return false
		}
		logDebugf("Using network type %s for connections", cm.networkType)
	}
	if !routeCfg.IsValid() {
		cm.configLock.Unlock()
		logDebugf("Routing data is not valid, skipping update: \n%s", routeCfg.DebugString())
		return false
	}

	// There's something wrong with this route config so don't send it to the watchers.
	if !cm.canUpdateRouteConfig(routeCfg) {
		cm.configLock.Unlock()
		return false
	}

	cm.currentConfig = routeCfg
	cm.seenConfig = true
	cm.configLock.Unlock()

	logDebugf("Sending out mux routing data (update)...")
	logDebugf("New Routing Data:\n%s", routeCfg.DebugString())

	// We can end up deadlocking if we iterate whilst in the lock and a watcher decides to remove itself.
	cm.watchersLock.Lock()
	watchers := make([]routeConfigWatcher, len(cm.cfgChangeWatchers))
	copy(watchers, cm.cfgChangeWatchers)
	cm.watchersLock.Unlock()

	for _, watcher := range watchers {
		watcher.OnNewRouteConfig(routeCfg)
	}

	return true
}

func (cm *configManagementComponent) RefreshConfig(snapshot *pipelineSnapshot) {
	currentRev, currentEpoch := cm.CurrentRev()
	cm.configFetchSigLock.Lock()
	if cm.configFetchSig != nil {
		// Someone else is already fetching a config so let's bail out.
		cm.configFetchSigLock.Unlock()
		return
	}
	cm.configFetchSig = make(chan struct{})
	cm.configFetchSigLock.Unlock()

	cm.fetchConfig(snapshot, currentRev, currentEpoch)

	cm.configFetchSigLock.Lock()
	close(cm.configFetchSig)
	cm.configFetchSig = nil
	cm.configFetchSigLock.Unlock()
}

func (cm *configManagementComponent) OnNewConfigChangeNotifBrief(snapshot *pipelineSnapshot, notif []byte) {
	if cm.configFetcher == nil {
		// No point in doing anything if we can't fetch a config anyway.
		return
	}
	if len(notif) != 16 {
		logWarnf("Invalid clustermap notification brief data size")
		return
	}
	serverRevEpoch := int64(binary.BigEndian.Uint64(notif[0:]))
	serverRevID := int64(binary.BigEndian.Uint64(notif[8:]))

	var currentRev, currentEpoch int64
	for {
		currentRev, currentEpoch = cm.CurrentRev()

		if serverRevEpoch < currentEpoch {
			logDebugf("Ignoring configuration notification as it has an older revision epoch. Old: %d, new: %d", currentEpoch, serverRevEpoch)
			return
		} else if serverRevEpoch == currentEpoch {
			if serverRevID == 0 {
				logDebugf("Unversioned configuration notification data, switching.")
			} else if serverRevID == currentRev {
				logDebugf("Ignoring configuration notification with identical revision number - %d", serverRevID)
				return
			} else if serverRevID < currentRev {
				logDebugf("Ignoring new configuration notification as it has an older revision id. Old: %d, new: %d", currentRev, serverRevID)
				return
			}
		}

		var waitSig chan struct{}
		cm.configFetchSigLock.Lock()
		if cm.configFetchSig == nil {
			cm.configFetchSig = make(chan struct{})
			cm.configFetchSigLock.Unlock()
			break
		}
		waitSig = cm.configFetchSig
		cm.configFetchSigLock.Unlock()

		<-waitSig
	}

	cm.fetchConfig(snapshot, currentRev, currentEpoch)

	cm.configFetchSigLock.Lock()
	close(cm.configFetchSig)
	cm.configFetchSig = nil
	cm.configFetchSigLock.Unlock()
}

func (cm *configManagementComponent) fetchConfig(snapshot *pipelineSnapshot, currentRev, currentEpoch int64) {
	if cm.configFetcher == nil {
		logDebugf("CfgManager: Cannot fetch config as the configFetcher is unset, likely because the agent is in ns server mode")
		return
	}

	numNodes := snapshot.NumPipelines()
	nodeIdx := rand.Intn(numNodes) // #nosec G404

	// We try to fetch the config from each node once.
	// If we cannot get it from any node then we just return.
	snapshot.Iterate(nodeIdx, func(pipeline *memdPipeline) bool {
		nodeIdx = (nodeIdx + 1) % numNodes
		if !pipeline.SupportsFeature(memd.FeatureClusterMapKnownVersion) {
			// No point in sending a request to a node that doesn't support known versions.
			return false
		}
		cfgBytes, err := cm.configFetcher.GetClusterConfig(pipeline, currentRev, currentEpoch, cm.shutdownSig)
		if err != nil {
			logDebugf("CfgManager: Failed to fetch config: %s", err)
			return false
		}
		if len(cfgBytes) == 0 {
			// The server didn't know about this revision.
			return false
		}

		logDebugf("CfgManager: Got Block: %s", string(cfgBytes))

		hostName, err := hostFromHostPort(pipeline.Address())
		if err != nil {
			logWarnf("CfgManager:Failed to parse source address. %s", err)
			return false
		}

		bk, err := parseConfig(cfgBytes, hostName)
		if err != nil {
			logDebugf("CfgManager:Failed to parse config. %v", err)
			return false
		}

		return cm.onNewConfig(bk)
	})
}

func (cm *configManagementComponent) Close() {
	close(cm.shutdownSig)
}

func (cm *configManagementComponent) Watchers() []routeConfigWatcher {
	cm.watchersLock.Lock()
	watchers := make([]routeConfigWatcher, len(cm.cfgChangeWatchers))
	copy(watchers, cm.cfgChangeWatchers)
	cm.watchersLock.Unlock()

	return watchers
}

func (cm *configManagementComponent) ResetConfig() {
	cm.configLock.Lock()
	cm.currentConfig = &routeConfig{
		revID: -1,
	}
	cm.configLock.Unlock()
}

func (cm *configManagementComponent) AddConfigWatcher(watcher routeConfigWatcher) {
	cm.watchersLock.Lock()
	cm.cfgChangeWatchers = append(cm.cfgChangeWatchers, watcher)
	cm.watchersLock.Unlock()
}

func (cm *configManagementComponent) RemoveConfigWatcher(watcher routeConfigWatcher) {
	var idx int
	var found bool
	cm.watchersLock.Lock()
	for i, w := range cm.cfgChangeWatchers {
		if w == watcher {
			idx = i
			found = true
			break
		}
	}

	if !found {
		cm.watchersLock.Unlock()
		return
	}

	if idx == len(cm.cfgChangeWatchers) {
		cm.cfgChangeWatchers = cm.cfgChangeWatchers[:idx]
	} else {
		cm.cfgChangeWatchers = append(cm.cfgChangeWatchers[:idx], cm.cfgChangeWatchers[idx+1:]...)
	}
	cm.watchersLock.Unlock()
}

// We should never be receiving concurrent updates and nothing should be accessing
// our internal route config so we shouldn't need to lock here.
func (cm *configManagementComponent) canUpdateRouteConfig(cfg *routeConfig) bool {
	oldCfg := cm.currentConfig

	// Check some basic things to ensure consistency!
	// If oldCfg name was empty and the new cfg isn't then we're moving from cluster to bucket connection.
	if cfg.revID > -1 && (oldCfg.name != "" && cfg.name != "") {
		if (cfg.vbMap == nil) != (oldCfg.vbMap == nil) {
			logErrorf("Received a configuration with a different number of vbuckets %s-%s.  Ignoring.", oldCfg.name, cfg.name)
			return false
		}

		if cfg.vbMap != nil && cfg.vbMap.NumVbuckets() != oldCfg.vbMap.NumVbuckets() {
			logErrorf("Received a configuration with a different number of vbuckets %s-%s.  Ignoring.", oldCfg.name, cfg.name)
			return false
		}
	}

	// Check that the new config data is newer than the current one, in the case where we've done a select bucket
	// against an existing connection then the revisions could be the same. In that case the configuration still
	// needs to be applied.
	// In the case where the rev epochs are the same then we need to compare rev IDs. If the new config epoch is lower
	// than the old one then we ignore it, if it's newer then we apply the new config.
	if cfg.bktType != oldCfg.bktType {
		logDebugf("Configuration data changed bucket type, switching.")
	} else if !cfg.IsNewerThan(oldCfg) {
		return false
	}

	return true
}

func (cm *configManagementComponent) buildFirstRouteConfig(config *cfgBucket, useSSL bool) *routeConfig {
	if cm.seedNodeAddr != "" {
		for _, node := range config.NodesExt {
			if node.ThisNode {
				cm.localLoopbackAddr = &localLoopbackAddress{
					LoopbackAddr: cm.seedNodeAddr,
					Identifier:   fmt.Sprintf("%s:%d", node.Hostname, node.Services.Mgmt),
				}
				break
			}
		}

		if cm.localLoopbackAddr == nil {
			logWarnf("Ignoring config, nodesExt entry contained no thisNode node")
			return &routeConfig{}
		}
	}
	if cm.networkType != "" && cm.networkType != "auto" {
		return config.BuildRouteConfig(useSSL, cm.networkType, true, cm.localLoopbackAddr)
	}

	defaultRouteConfig := config.BuildRouteConfig(useSSL, "default", true, cm.localLoopbackAddr)

	var kvServerList []routeEndpoint
	var mgmtEpList []routeEndpoint
	if useSSL {
		kvServerList = defaultRouteConfig.kvServerList.SSLEndpoints
		mgmtEpList = defaultRouteConfig.mgmtEpList.SSLEndpoints
	} else {
		kvServerList = defaultRouteConfig.kvServerList.NonSSLEndpoints
		mgmtEpList = defaultRouteConfig.mgmtEpList.NonSSLEndpoints
	}

	// Iterate over all the source servers and check if any addresses match as default or external network types
	for _, srcServer := range cm.srcServers {
		// First we check if the source server is from the defaults list
		srcInDefaultConfig := false
		for _, endpoint := range kvServerList {
			if trimSchemePrefix(endpoint.Address) == srcServer.Address {
				srcInDefaultConfig = true
			}
		}
		for _, endpoint := range mgmtEpList {
			if endpoint == srcServer {
				srcInDefaultConfig = true
			}
		}
		if srcInDefaultConfig {
			cm.networkType = "default"
			return defaultRouteConfig
		}
	}

	// Next lets see if we have an external config, if so, default to that
	externalRouteCfg := config.BuildRouteConfig(useSSL, "external", true, cm.localLoopbackAddr)
	if externalRouteCfg.IsValid() {
		cm.networkType = "external"
		return externalRouteCfg
	}

	// If all else fails, default to the implicit default config
	cm.networkType = "default"
	return defaultRouteConfig
}

func (cm *configManagementComponent) NetworkType() string {
	return cm.networkType
}
