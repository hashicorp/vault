package plugin

import (
	"encoding/json"
	"fmt"
	"strconv"
	"time"
)

// logEntry is the JSON payload that gets sent to Stderr from the plugin to the host
type logEntry struct {
	Message   string        `json:"@message"`
	Level     string        `json:"@level"`
	Timestamp time.Time     `json:"timestamp"`
	KVPairs   []*logEntryKV `json:"kv_pairs"`
}

// logEntryKV is a key value pair within the Output payload
type logEntryKV struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

// parseKVPairs transforms string inputs into []*logEntryKV
func parseKVPairs(kvs ...interface{}) ([]*logEntryKV, error) {
	var result []*logEntryKV
	if len(kvs)%2 != 0 {
		return nil, fmt.Errorf("kv slice needs to be even number, got %d", len(kvs))
	}
	for i := 0; i < len(kvs); i = i + 2 {
		var val string

		switch st := kvs[i+1].(type) {
		case string:
			val = st
		case int:
			val = strconv.FormatInt(int64(st), 10)
		case int64:
			val = strconv.FormatInt(int64(st), 10)
		case int32:
			val = strconv.FormatInt(int64(st), 10)
		case int16:
			val = strconv.FormatInt(int64(st), 10)
		case int8:
			val = strconv.FormatInt(int64(st), 10)
		case uint:
			val = strconv.FormatUint(uint64(st), 10)
		case uint64:
			val = strconv.FormatUint(uint64(st), 10)
		case uint32:
			val = strconv.FormatUint(uint64(st), 10)
		case uint16:
			val = strconv.FormatUint(uint64(st), 10)
		case uint8:
			val = strconv.FormatUint(uint64(st), 10)
		default:
			val = fmt.Sprintf("%v", st)
		}

		result = append(result, &logEntryKV{
			Key:   kvs[i].(string),
			Value: val,
		})
	}

	return result, nil
}

// flattenKVPairs is used to flatten KVPair slice into []interface{}
// for hclog consumption.
func flattenKVPairs(kvs []*logEntryKV) []interface{} {
	var result []interface{}
	for _, kv := range kvs {
		result = append(result, kv.Key)
		result = append(result, kv.Value)
	}

	return result
}

// parseJSON handles parsing JSON output
func parseJSON(input string) (*logEntry, error) {
	var raw map[string]interface{}
	entry := &logEntry{}

	err := json.Unmarshal([]byte(input), &raw)
	if err != nil {
		return nil, err
	}

	// Parse hclog-specific objects
	if v, ok := raw["@message"]; ok {
		entry.Message = v.(string)
		delete(raw, "@message")
	}

	if v, ok := raw["@level"]; ok {
		entry.Level = v.(string)
		delete(raw, "@level")
	}

	if v, ok := raw["@timestamp"]; ok {
		t, err := time.Parse("2006-01-02T15:04:05.000000Z07:00", v.(string))
		if err != nil {
			return nil, err
		}
		entry.Timestamp = t
		delete(raw, "@timestamp")
	}

	// Parse dynamic KV args from the hclog payload.
	kvs := []interface{}{}
	for k, v := range raw {
		kvs = append(kvs, k)
		kvs = append(kvs, v.(string))
	}
	pairs, err := parseKVPairs(kvs...)
	if err != nil {
		return nil, err
	}
	entry.KVPairs = pairs

	return entry, nil
}
