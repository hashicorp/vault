/*
Copyright 2014 SAP SE

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

package protocol

import (
	"os"
	"strconv"
	"strings"

	"github.com/SAP/go-hdb/internal/bufio"
)

type clientID []byte

func newClientID() clientID {
	if h, err := os.Hostname(); err == nil {
		return clientID(strings.Join([]string{strconv.Itoa(os.Getpid()), h}, "@"))
	}
	return clientID(strconv.Itoa(os.Getpid()))
}

func (id clientID) kind() partKind {
	return partKind(pkClientID)
}

func (id clientID) size() (int, error) {
	return len(id), nil
}

func (id clientID) numArg() int {
	return 1
}

func (id clientID) write(wr *bufio.Writer) error {
	wr.Write(id)

	if trace {
		outLogger.Printf("client id: %s", id)
	}
	return nil
}
