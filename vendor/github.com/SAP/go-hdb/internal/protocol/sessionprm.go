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

import "fmt"

type SessionPrm struct {
	Host, Username, Password       string
	Locale                         string
	BufferSize, FetchSize, Timeout int
}

func (p *SessionPrm) String() string {
	return fmt.Sprintf("session parameters: bufferSize %d fetchSize %d timeout %d locale %s", p.BufferSize, p.FetchSize, p.Timeout, p.Locale)
}
