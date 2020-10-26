/*
 * Copyright 2018 - Present Okta, Inc.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *      http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package cache

import "net/http"

type NoOpCache struct {
}

func NewNoOpCache() Cache {
	return NoOpCache{}
}

func (c NoOpCache) Get(key string) *http.Response {
	return nil
}

func (c NoOpCache) Set(key string, value *http.Response) {

}

func (c NoOpCache) GetString(key string) string {
	return ""
}

func (c NoOpCache) SetString(key string, value string) {

}

func (c NoOpCache) Delete(key string) {

}

func (c NoOpCache) Clear() {

}

func (c NoOpCache) Has(key string) bool {
	return false
}
