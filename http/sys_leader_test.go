// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package http

import (
	"encoding/json"
	"net/http"
	"reflect"
	"testing"
	"time"

	"github.com/hashicorp/vault/vault"
)

func TestSysLeader_get(t *testing.T) {
	core, _, _ := vault.TestCoreUnsealed(t)
	ln, addr := TestServer(t, core)
	defer ln.Close()

	resp, err := http.Get(addr + "/v1/sys/leader")
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	var actual map[string]interface{}
	expected := map[string]interface{}{
		"ha_enabled":                          false,
		"is_self":                             false,
		"leader_address":                      "",
		"leader_cluster_address":              "",
		"performance_standby":                 false,
		"performance_standby_last_remote_wal": json.Number("0"),
		"active_time":                         time.Time{}.UTC().Format(time.RFC3339),
	}
	testResponseStatus(t, resp, 200)
	testResponseBody(t, resp, &actual)
	if !reflect.DeepEqual(actual, expected) {
		t.Fatalf("bad: %#v \n%#v", actual, expected)
	}
}
