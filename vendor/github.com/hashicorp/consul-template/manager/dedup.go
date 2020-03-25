package manager

import (
	"bytes"
	"compress/lzw"
	"encoding/gob"
	"fmt"
	"log"
	"path"
	"sync"
	"time"

	"github.com/mitchellh/hashstructure"

	"github.com/hashicorp/consul-template/config"
	dep "github.com/hashicorp/consul-template/dependency"
	"github.com/hashicorp/consul-template/template"
	"github.com/hashicorp/consul-template/version"
	consulapi "github.com/hashicorp/consul/api"
)

var (
	// sessionCreateRetry is the amount of time we wait
	// to recreate a session when lost.
	sessionCreateRetry = 15 * time.Second

	// lockRetry is the interval on which we try to re-acquire locks
	lockRetry = 10 * time.Second

	// listRetry is the interval on which we retry listing a data path
	listRetry = 10 * time.Second

	// timeout passed through to consul api client Lock
	// here to override in testing (see ./dedup_test.go)
	lockWaitTime = 15 * time.Second
)

const (
	templateNoDataStr = "__NO_DATA__"
)

// templateData is GOB encoded share the dependency values
type templateData struct {
	// Version is the version of Consul Template which created this template data.
	// This is important because users may be running multiple versions of CT
	// with the same templates. This provides a nicer upgrade path.
	Version string

	// Data is the actual template data.
	Data map[string]interface{}
}

func templateNoData() []byte {
	return []byte(templateNoDataStr)
}

// DedupManager is used to de-duplicate which instance of Consul-Template
// is handling each template. For each template, a lock path is determined
// using the MD5 of the template. This path is used to elect a "leader"
// instance.
//
// The leader instance operations like usual, but any time a template is
// rendered, any of the data required for rendering is stored in the
// Consul KV store under the lock path.
//
// The follower instances depend on the leader to do the primary watching
// and rendering, and instead only watch the aggregated data in the KV.
// Followers wait for updates and re-render the template.
//
// If a template depends on 50 views, and is running on 50 machines, that
// would normally require 2500 blocking queries. Using deduplication, one
// instance has 50 view queries, plus 50 additional queries on the lock
// path for a total of 100.
//
type DedupManager struct {
	// config is the deduplicate configuration
	config *config.DedupConfig

	// clients is used to access the underlying clients
	clients *dep.ClientSet

	// Brain is where we inject updates
	brain *template.Brain

	// templates is the set of templates we are trying to dedup
	templates []*template.Template

	// leader tracks if we are currently the leader
	leader     map[*template.Template]<-chan struct{}
	leaderLock sync.RWMutex

	// lastWrite tracks the hash of the data paths
	lastWrite     map[*template.Template]uint64
	lastWriteLock sync.RWMutex

	// updateCh is used to indicate an update watched data
	updateCh chan struct{}

	// wg is used to wait for a clean shutdown
	wg sync.WaitGroup

	stop     bool
	stopCh   chan struct{}
	stopLock sync.Mutex
}

// NewDedupManager creates a new Dedup manager
func NewDedupManager(config *config.DedupConfig, clients *dep.ClientSet, brain *template.Brain, templates []*template.Template) (*DedupManager, error) {
	d := &DedupManager{
		config:    config,
		clients:   clients,
		brain:     brain,
		templates: templates,
		leader:    make(map[*template.Template]<-chan struct{}),
		lastWrite: make(map[*template.Template]uint64),
		updateCh:  make(chan struct{}, 1),
		stopCh:    make(chan struct{}),
	}
	return d, nil
}

// Start is used to start the de-duplication manager
func (d *DedupManager) Start() error {
	log.Printf("[INFO] (dedup) starting de-duplication manager")

	client := d.clients.Consul()
	go d.createSession(client)

	// Start to watch each template
	for _, t := range d.templates {
		go d.watchTemplate(client, t)
	}
	return nil
}

// Stop is used to stop the de-duplication manager
func (d *DedupManager) Stop() error {
	d.stopLock.Lock()
	defer d.stopLock.Unlock()
	if d.stop {
		return nil
	}

	log.Printf("[INFO] (dedup) stopping de-duplication manager")
	d.stop = true
	close(d.stopCh)
	d.wg.Wait()
	return nil
}

// createSession is used to create and maintain a session to Consul
func (d *DedupManager) createSession(client *consulapi.Client) {
START:
	log.Printf("[INFO] (dedup) attempting to create session")
	session := client.Session()
	sessionCh := make(chan struct{})
	ttl := fmt.Sprintf("%.6fs", float64(*d.config.TTL)/float64(time.Second))
	se := &consulapi.SessionEntry{
		Name:      "Consul-Template de-duplication",
		Behavior:  "delete",
		TTL:       ttl,
		LockDelay: 1 * time.Millisecond,
	}
	id, _, err := session.Create(se, nil)
	if err != nil {
		log.Printf("[ERR] (dedup) failed to create session: %v", err)
		goto WAIT
	}
	log.Printf("[INFO] (dedup) created session %s", id)

	// Attempt to lock each template
	for _, t := range d.templates {
		d.wg.Add(1)
		go d.attemptLock(client, id, sessionCh, t)
	}

	// Renew our session periodically
	if err := session.RenewPeriodic("15s", id, nil, d.stopCh); err != nil {
		log.Printf("[ERR] (dedup) failed to renew session: %v", err)
	}
	close(sessionCh)
	d.wg.Wait()

WAIT:
	select {
	case <-time.After(sessionCreateRetry):
		goto START
	case <-d.stopCh:
		return
	}
}

// IsLeader checks if we are currently the leader instance
func (d *DedupManager) IsLeader(tmpl *template.Template) bool {
	d.leaderLock.RLock()
	defer d.leaderLock.RUnlock()

	lockCh, ok := d.leader[tmpl]
	if !ok {
		return false
	}
	select {
	case <-lockCh:
		return false
	default:
		return true
	}
}

// UpdateDeps is used to update the values of the dependencies for a template
func (d *DedupManager) UpdateDeps(t *template.Template, deps []dep.Dependency) error {
	// Calculate the path to write updates to
	dataPath := path.Join(*d.config.Prefix, t.ID(), "data")

	// Package up the dependency data
	td := templateData{
		Version: version.Version,
		Data:    make(map[string]interface{}),
	}
	for _, dp := range deps {
		// Skip any dependencies that can't be shared
		if !dp.CanShare() {
			continue
		}

		// Pull the current value from the brain
		val, ok := d.brain.Recall(dp)
		if ok {
			td.Data[dp.String()] = val
		}
	}

	// Compute stable hash of the data. Note we don't compute this over the actual
	// encoded value since gob encoding does not guarantee stable ordering for
	// maps so spuriously returns a different hash most times. See
	// https://github.com/hashicorp/consul-template/issues/1099.
	hash, err := hashstructure.Hash(td, nil)
	if err != nil {
		return fmt.Errorf("calculating hash failed: %v", err)
	}
	d.lastWriteLock.RLock()
	existing, ok := d.lastWrite[t]
	d.lastWriteLock.RUnlock()
	if ok && existing == hash {
		log.Printf("[INFO] (dedup) de-duplicate data '%s' already current",
			dataPath)
		return nil
	}

	// Encode via GOB and LZW compress
	var buf bytes.Buffer
	compress := lzw.NewWriter(&buf, lzw.LSB, 8)
	enc := gob.NewEncoder(compress)
	if err := enc.Encode(&td); err != nil {
		return fmt.Errorf("encode failed: %v", err)
	}
	compress.Close()

	// Write the KV update
	kvPair := consulapi.KVPair{
		Key:   dataPath,
		Value: buf.Bytes(),
		Flags: consulapi.LockFlagValue,
	}
	client := d.clients.Consul()
	if _, err := client.KV().Put(&kvPair, nil); err != nil {
		return fmt.Errorf("failed to write '%s': %v", dataPath, err)
	}
	log.Printf("[INFO] (dedup) updated de-duplicate data '%s'", dataPath)
	d.lastWriteLock.Lock()
	d.lastWrite[t] = hash
	d.lastWriteLock.Unlock()
	return nil
}

// UpdateCh returns a channel to watch for dependency updates
func (d *DedupManager) UpdateCh() <-chan struct{} {
	return d.updateCh
}

// setLeader sets if we are currently the leader instance
func (d *DedupManager) setLeader(tmpl *template.Template, lockCh <-chan struct{}) {
	// Update the lock state
	d.leaderLock.Lock()
	if lockCh != nil {
		d.leader[tmpl] = lockCh
	} else {
		delete(d.leader, tmpl)
	}
	d.leaderLock.Unlock()

	// Clear the lastWrite hash if we've lost leadership
	if lockCh == nil {
		d.lastWriteLock.Lock()
		delete(d.lastWrite, tmpl)
		d.lastWriteLock.Unlock()
	}

	// Do an async notify of an update
	select {
	case d.updateCh <- struct{}{}:
	default:
	}
}

func (d *DedupManager) watchTemplate(client *consulapi.Client, t *template.Template) {
	log.Printf("[INFO] (dedup) starting watch for template hash %s", t.ID())
	path := path.Join(*d.config.Prefix, t.ID(), "data")

	// Determine if stale queries are allowed
	var allowStale bool
	if *d.config.MaxStale != 0 {
		allowStale = true
	}

	// Setup our query options
	opts := &consulapi.QueryOptions{
		AllowStale: allowStale,
		WaitTime:   60 * time.Second,
	}

	var lastData []byte
	var lastIndex uint64

START:
	// Stop listening if we're stopped
	select {
	case <-d.stopCh:
		return
	default:
	}

	// If we are current the leader, wait for leadership lost
	d.leaderLock.RLock()
	lockCh, ok := d.leader[t]
	d.leaderLock.RUnlock()
	if ok {
		select {
		case <-lockCh:
			goto START
		case <-d.stopCh:
			return
		}
	}

	// Block for updates on the data key
	log.Printf("[INFO] (dedup) listing data for template hash %s", t.ID())
	pair, meta, err := client.KV().Get(path, opts)
	if err != nil {
		log.Printf("[ERR] (dedup) failed to get '%s': %v", path, err)
		select {
		case <-time.After(listRetry):
			goto START
		case <-d.stopCh:
			return
		}
	}
	opts.WaitIndex = meta.LastIndex

	// Stop listening if we're stopped
	select {
	case <-d.stopCh:
		return
	default:
	}

	// If we've exceeded the maximum staleness, retry without stale
	if allowStale && meta.LastContact > *d.config.MaxStale {
		allowStale = false
		log.Printf("[DEBUG] (dedup) %s stale data (last contact exceeded max_stale)", path)
		goto START
	}

	// Re-enable stale queries if allowed
	if *d.config.MaxStale > 0 {
		allowStale = true
	}

	if meta.LastIndex == lastIndex {
		log.Printf("[TRACE] (dedup) %s no new data (index was the same)", path)
		goto START
	}

	if meta.LastIndex < lastIndex {
		log.Printf("[TRACE] (dedup) %s had a lower index, resetting", path)
		lastIndex = 0
		goto START
	}
	lastIndex = meta.LastIndex

	var data []byte
	if pair != nil {
		data = pair.Value
	}
	if bytes.Equal(lastData, data) {
		log.Printf("[TRACE] (dedup) %s no new data (contents were the same)", path)
		goto START
	}
	lastData = data

	// If we are current the leader, wait for leadership lost
	d.leaderLock.RLock()
	lockCh, ok = d.leader[t]
	d.leaderLock.RUnlock()
	if ok {
		select {
		case <-lockCh:
			goto START
		case <-d.stopCh:
			return
		}
	}

	// Parse the data file
	if pair != nil && pair.Flags == consulapi.LockFlagValue && !bytes.Equal(pair.Value, templateNoData()) {
		d.parseData(pair.Key, pair.Value)
	}
	goto START
}

// parseData is used to update brain from a KV data pair
func (d *DedupManager) parseData(path string, raw []byte) {
	// Setup the decompression and decoders
	r := bytes.NewReader(raw)
	decompress := lzw.NewReader(r, lzw.LSB, 8)
	defer decompress.Close()
	dec := gob.NewDecoder(decompress)

	// Decode the data
	var td templateData
	if err := dec.Decode(&td); err != nil {
		log.Printf("[ERR] (dedup) failed to decode '%s': %v",
			path, err)
		return
	}
	if td.Version != version.Version {
		log.Printf("[WARN] (dedup) created with different version (%s vs %s)",
			td.Version, version.Version)
		return
	}
	log.Printf("[INFO] (dedup) loading %d dependencies from '%s'",
		len(td.Data), path)

	// Update the data in the brain
	for hashCode, value := range td.Data {
		d.brain.ForceSet(hashCode, value)
	}

	// Trigger the updateCh
	select {
	case d.updateCh <- struct{}{}:
	default:
	}
}

func (d *DedupManager) attemptLock(client *consulapi.Client, session string, sessionCh chan struct{}, t *template.Template) {
	defer d.wg.Done()
	for {
		log.Printf("[INFO] (dedup) attempting lock for template hash %s", t.ID())
		basePath := path.Join(*d.config.Prefix, t.ID())
		lopts := &consulapi.LockOptions{
			Key:              path.Join(basePath, "data"),
			Value:            templateNoData(),
			Session:          session,
			MonitorRetries:   3,
			MonitorRetryTime: 3 * time.Second,
			LockWaitTime:     lockWaitTime,
		}
		lock, err := client.LockOpts(lopts)
		if err != nil {
			log.Printf("[ERR] (dedup) failed to create lock '%s': %v",
				lopts.Key, err)
			return
		}

		var retryCh <-chan time.Time
		leaderCh, err := lock.Lock(sessionCh)
		if err != nil {
			log.Printf("[ERR] (dedup) failed to acquire lock '%s': %v",
				lopts.Key, err)
			retryCh = time.After(lockRetry)
		} else {
			log.Printf("[INFO] (dedup) acquired lock '%s'", lopts.Key)
			d.setLeader(t, leaderCh)
		}

		select {
		case <-retryCh:
			retryCh = nil
			continue
		case <-leaderCh:
			log.Printf("[WARN] (dedup) lost lock ownership '%s'", lopts.Key)
			d.setLeader(t, nil)
			continue
		case <-sessionCh:
			log.Printf("[INFO] (dedup) releasing session '%s'", lopts.Key)
			d.setLeader(t, nil)
			_, err = client.Session().Destroy(session, nil)
			if err != nil {
				log.Printf("[ERROR] (dedup) failed destroying session '%s', %s", session, err)
			}
			return
		case <-d.stopCh:
			log.Printf("[INFO] (dedup) releasing lock '%s'", lopts.Key)
			_, err = client.Session().Destroy(session, nil)
			if err != nil {
				log.Printf("[ERROR] (dedup) failed destroying session '%s', %s", session, err)
			}
			return
		}
	}
}
