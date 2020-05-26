package gocbcore

import (
	"encoding/json"
	"fmt"
	"sort"
	"sync"
	"time"
)

type zombieLogEntry struct {
	connectionID string
	operationID  string
	endpoint     string
	duration     time.Duration
	serviceType  string
}

type zombieLogItem struct {
	ConnectionID     string `json:"c"`
	OperationID      string `json:"i"`
	Endpoint         string `json:"r"`
	ServerDurationUs uint64 `json:"d"`
	ServiceType      string `json:"s"`
}

type zombieLogService struct {
	Service string          `json:"service"`
	Count   int             `json:"count"`
	Top     []zombieLogItem `json:"top"`
}

type zombieLoggerComponent struct {
	zombieLock sync.RWMutex
	zombieOps  []*zombieLogEntry
	interval   time.Duration
	sampleSize int
	stopSig    chan struct{}
}

func newZombieLoggerComponent(interval time.Duration, sampleSize int) *zombieLoggerComponent {
	return &zombieLoggerComponent{
		// zombieOps must have a static capacity for its lifetime, the capacity should
		// never be altered so that it is consistent across the zombieLogger and
		// recordZombieResponse.
		zombieOps:  make([]*zombieLogEntry, 0, sampleSize),
		interval:   interval,
		sampleSize: sampleSize,
		stopSig:    make(chan struct{}),
	}
}

func (zlc *zombieLoggerComponent) Start() {
	lastTick := time.Now()

	for {
		select {
		case <-zlc.stopSig:
			return
		case <-time.After(zlc.interval):
		}

		lastTick = lastTick.Add(zlc.interval)

		// Preallocate space to copy the ops into...
		oldOps := make([]*zombieLogEntry, zlc.sampleSize)

		zlc.zombieLock.Lock()
		// Escape early if we have no ops to log...
		if len(zlc.zombieOps) == 0 {
			zlc.zombieLock.Unlock()
			continue
		}

		// Copy out our ops so we can cheaply print them out without blocking
		// our ops from actually being recorded in other goroutines (which would
		// effectively slow down the op pipeline for logging).
		oldOps = oldOps[0:len(zlc.zombieOps)]
		copy(oldOps, zlc.zombieOps)
		zlc.zombieOps = zlc.zombieOps[:0]

		zlc.zombieLock.Unlock()

		jsonData := zombieLogService{
			Service: "kv",
		}

		for i := len(oldOps) - 1; i >= 0; i-- {
			op := oldOps[i]

			jsonData.Top = append(jsonData.Top, zombieLogItem{
				OperationID:      op.operationID,
				ConnectionID:     op.connectionID,
				Endpoint:         op.endpoint,
				ServerDurationUs: uint64(op.duration / time.Microsecond),
				ServiceType:      op.serviceType,
			})
		}

		jsonData.Count = len(jsonData.Top)

		jsonBytes, err := json.Marshal(jsonData)
		if err != nil {
			logDebugf("Failed to generate zombie logging JSON: %s", err)
		}

		logWarnf("Orphaned responses observed:\n %s", jsonBytes)
	}
}

func (zlc *zombieLoggerComponent) Stop() {
	close(zlc.stopSig)
}

func (zlc *zombieLoggerComponent) RecordZombieResponse(resp *memdQResponse, connID, address string) {
	entry := &zombieLogEntry{
		connectionID: connID,
		operationID:  fmt.Sprintf("0x%x", resp.Opaque),
		endpoint:     address,
		duration:     0,
		serviceType:  fmt.Sprintf("kv:%s", resp.Command.Name()),
	}

	if resp.Packet.ServerDurationFrame != nil {
		entry.duration = resp.Packet.ServerDurationFrame.ServerDuration
	}

	zlc.zombieLock.RLock()

	if cap(zlc.zombieOps) == 0 || (len(zlc.zombieOps) == cap(zlc.zombieOps) &&
		entry.duration < zlc.zombieOps[0].duration) {
		// we are at capacity and we are faster than the fastest slow op or somehow in a state where capacity is 0.
		zlc.zombieLock.RUnlock()
		return
	}
	zlc.zombieLock.RUnlock()

	zlc.zombieLock.Lock()
	if cap(zlc.zombieOps) == 0 || (len(zlc.zombieOps) == cap(zlc.zombieOps) &&
		entry.duration < zlc.zombieOps[0].duration) {
		// we are at capacity and we are faster than the fastest slow op or somehow in a state where capacity is 0.
		zlc.zombieLock.Unlock()
		return
	}

	l := len(zlc.zombieOps)
	i := sort.Search(l, func(i int) bool { return entry.duration < zlc.zombieOps[i].duration })

	// i represents the slot where it should be inserted

	if len(zlc.zombieOps) < cap(zlc.zombieOps) {
		if i == l {
			zlc.zombieOps = append(zlc.zombieOps, entry)
		} else {
			zlc.zombieOps = append(zlc.zombieOps, nil)
			copy(zlc.zombieOps[i+1:], zlc.zombieOps[i:])
			zlc.zombieOps[i] = entry
		}
	} else {
		if i == 0 {
			zlc.zombieOps[i] = entry
		} else {
			copy(zlc.zombieOps[0:i-1], zlc.zombieOps[1:i])
			zlc.zombieOps[i-1] = entry
		}
	}

	zlc.zombieLock.Unlock()
}
