// +build !nocgo

/*
Copyright 2016 The Camlistore Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

     http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

// Package netutil identifies the system userid responsible for
// localhost TCP connections.
package netutil // import "camlistore.org/pkg/netutil"

import (
	"os"
	"os/user"
	"strconv"
)

func uidFromUsernameFn(username string) (uid int, err error) {
	if uid := os.Getuid(); uid != 0 && username == os.Getenv("USER") {
		return uid, nil
	}
	u, err := user.Lookup(username)
	if err == nil {
		uid, err := strconv.Atoi(u.Uid)
		return uid, err
	}
	return 0, err
}
