/*
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 * http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package endpoints

import (
	"fmt"
	"strings"
	"sync"
)

const keyFormatter = "%s::%s"

type EndpointMapping struct {
	sync.RWMutex
	endpoint map[string]string
}

var endpointMapping = EndpointMapping{endpoint: make(map[string]string)}

// AddEndpointMapping use productId and regionId as key to store the endpoint into inner map
// when using the same productId and regionId as key, the endpoint will be covered.
func AddEndpointMapping(regionId, productId, endpoint string) (err error) {
	key := fmt.Sprintf(keyFormatter, strings.ToLower(regionId), strings.ToLower(productId))
	endpointMapping.Lock()
	endpointMapping.endpoint[key] = endpoint
	endpointMapping.Unlock()
	return nil
}

// GetEndpointFromMap use Product and RegionId as key to find endpoint from inner map
func GetEndpointFromMap(regionId, productId string) string {
	key := fmt.Sprintf(keyFormatter, strings.ToLower(regionId), strings.ToLower(productId))
	endpointMapping.RLock()
	endpoint := endpointMapping.endpoint[key]
	endpointMapping.RUnlock()
	return endpoint
}
