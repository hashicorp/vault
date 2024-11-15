// Copyright 2021 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     https://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package spanner

import "google.golang.org/api/option"

// Returns the default client options used by the generated Spanner client.
//
// This function is only intended for use by the client library, and may be
// removed at any time without any warning.
func DefaultClientOptions() []option.ClientOption {
	return defaultGRPCClientOptions()
}
