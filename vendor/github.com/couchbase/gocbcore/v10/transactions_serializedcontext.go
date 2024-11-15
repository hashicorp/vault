// Copyright 2021 Couchbase
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

package gocbcore

type jsonSerializedMutation struct {
	Bucket     string `json:"bkt"`
	Scope      string `json:"scp"`
	Collection string `json:"coll"`
	ID         string `json:"id"`
	Cas        string `json:"cas"`
	Type       string `json:"type"`
}

type jsonSerializedAttempt struct {
	ID struct {
		Transaction string `json:"txn"`
		Attempt     string `json:"atmpt"`
	} `json:"id"`
	ATR struct {
		Bucket     string `json:"bkt"`
		Scope      string `json:"scp"`
		Collection string `json:"coll"`
		ID         string `json:"id"`
	} `json:"atr"`
	Config struct {
		KeyValueTimeoutMs int    `json:"kvTimeoutMs"`
		DurabilityLevel   string `json:"durabilityLevel"`
		NumAtrs           int    `json:"numAtrs"`
	} `json:"config"`
	State struct {
		TimeLeftMs int `json:"timeLeftMs"`
	} `json:"state"`
	Mutations []jsonSerializedMutation `json:"mutations"`
}
