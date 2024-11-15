// Copyright 2013 M-Lab
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// +build appengine

package digest

import (
	"appengine"
	"appengine/urlfetch"
)

// GAETransport returns an implementation of http.RoundTripper that uses the GAE
// urlfetch.Transport and is capable of HTTP digest authentication.
func GAETransport(c appengine.Context, username, password string) *Transport {
	t := &Transport{
		Username: username,
		Password: password,
	}
	t.Transport = &urlfetch.Transport{
		Context: c,
	}
	return t
}
