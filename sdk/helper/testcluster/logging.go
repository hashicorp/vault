// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package testcluster

import (
	"encoding/json"
	"strings"

	"github.com/hashicorp/go-hclog"
)

func JSONLogNoTimestamp(outlog hclog.Logger, text string) {
	d := json.NewDecoder(strings.NewReader(text))
	m := map[string]interface{}{}
	if err := d.Decode(&m); err != nil {
		outlog.Error("failed to decode json output from dev vault", "error", err, "input", text)
		return
	}
	JSONLogNoTimestampFromMap(outlog, "", m)
}

func JSONLogNoTimestampFromMap(outlog hclog.Logger, min string, m map[string]any) string {
	timeStamp := ""
	if ts, ok := m["@timestamp"]; ok {
		timeStamp = ts.(string)
		// This works because the time format is ISO8601.
		if timeStamp < min {
			return min
		}
	}

	delete(m, "@timestamp")
	message := m["@message"].(string)
	delete(m, "@message")
	level := m["@level"].(string)
	delete(m, "@level")
	if module, ok := m["@module"]; ok {
		delete(m, "@module")
		outlog = outlog.Named(module.(string))
	}

	var pairs []interface{}
	for k, v := range m {
		pairs = append(pairs, k, v)
	}

	outlog.Log(hclog.LevelFromString(level), message, pairs...)
	return timeStamp
}
