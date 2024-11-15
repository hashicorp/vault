package gosnowflake

import (
	"sort"
	"strconv"
	"sync"
)

const (
	queryContextCacheSizeParamName = "QUERY_CONTEXT_CACHE_SIZE"
	defaultQueryContextCacheSize   = 5
)

type queryContext struct {
	Entries []queryContextEntry `json:"entries,omitempty"`
}

type queryContextEntry struct {
	ID        int    `json:"id"`
	Timestamp int64  `json:"timestamp"`
	Priority  int    `json:"priority"`
	Context   string `json:"context,omitempty"`
}

type queryContextCache struct {
	mutex   *sync.Mutex
	entries []queryContextEntry
}

func (qcc *queryContextCache) init() *queryContextCache {
	qcc.mutex = &sync.Mutex{}
	return qcc
}

func (qcc *queryContextCache) add(sc *snowflakeConn, qces ...queryContextEntry) {
	qcc.mutex.Lock()
	defer qcc.mutex.Unlock()
	if len(qces) == 0 {
		qcc.prune(0)
	} else {
		for _, newQce := range qces {
			logger.Debugf("adding query context: %v", newQce)
			newQceProcessed := false
			for existingQceIdx, existingQce := range qcc.entries {
				if newQce.ID == existingQce.ID {
					newQceProcessed = true
					if newQce.Timestamp > existingQce.Timestamp {
						qcc.entries[existingQceIdx] = newQce
					} else if newQce.Timestamp == existingQce.Timestamp {
						if newQce.Priority != existingQce.Priority {
							qcc.entries[existingQceIdx] = newQce
						}
					}
				}
			}
			if !newQceProcessed {
				for existingQceIdx, existingQce := range qcc.entries {
					if newQce.Priority == existingQce.Priority {
						qcc.entries[existingQceIdx] = newQce
						newQceProcessed = true
					}
				}
			}
			if !newQceProcessed {
				qcc.entries = append(qcc.entries, newQce)
			}
		}
		sort.Slice(qcc.entries, func(idx1, idx2 int) bool {
			return qcc.entries[idx1].Priority < qcc.entries[idx2].Priority
		})
		qcc.prune(qcc.getQueryContextCacheSize(sc))
	}
}

func (qcc *queryContextCache) prune(size int) {
	if len(qcc.entries) > size {
		qcc.entries = qcc.entries[0:size]
	}
}

func (qcc *queryContextCache) getQueryContextCacheSize(sc *snowflakeConn) int {
	paramsMutex.Lock()
	sizeStr, ok := sc.cfg.Params[queryContextCacheSizeParamName]
	paramsMutex.Unlock()
	if ok {
		size, err := strconv.Atoi(*sizeStr)
		if err != nil {
			logger.Warnf("cannot parse %v as int as query context cache size: %v", sizeStr, err)
		} else {
			return size
		}
	}
	return defaultQueryContextCacheSize
}
