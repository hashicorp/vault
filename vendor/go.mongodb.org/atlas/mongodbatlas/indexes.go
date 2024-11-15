// Copyright 2021 MongoDB Inc
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package mongodbatlas

import (
	"context"
	"fmt"
	"net/http"
)

const indexesPath = "api/atlas/v1.0/groups/%s/clusters/%s/index"

// IndexesService is an interface for interfacing with the clusters indexes
// endpoints of the MongoDB Atlas API.
//
// See more: https://docs.atlas.mongodb.com/reference/api/indexes/
type IndexesService interface {
	Create(context.Context, string, string, *IndexConfiguration) (*Response, error)
}

// IndexesServiceOp handles communication with the Cluster related methods
// of the MongoDB Atlas API.
type IndexesServiceOp service

var _ IndexesService = &IndexesServiceOp{}

// IndexConfiguration represents a new index requests for a given database and collection.
type IndexConfiguration struct {
	DB         string              `json:"db"`                  // DB the database of the index
	Collection string              `json:"collection"`          // Collection the collection of the index
	Keys       []map[string]string `json:"keys"`                // Keys array of keys to index and their type, sorting of keys is important for an index
	Options    *IndexOptions       `json:"options,omitempty"`   // Options MongoDB index options
	Collation  *CollationOptions   `json:"collation,omitempty"` // Collation Mongo collation index options
}

// IndexOptions represents mongodb index options.
//
// See: https://docs.mongodb.com/manual/reference/method/db.collection.createIndex/#options
type IndexOptions struct {
	Background              bool                    `json:"background,omitempty"`
	PartialFilterExpression *map[string]interface{} `json:"partialFilterExpression,omitempty"`
	StorageEngine           *map[string]interface{} `json:"storageEngine,omitempty"`
	Weights                 *map[string]int         `json:"weights,omitempty"`
	DefaultLanguage         string                  `json:"default_language,omitempty"`
	LanguageOverride        string                  `json:"language_override,omitempty"`
	TextIndexVersion        int                     `json:"textIndexVersion,omitempty"`
	TwodsphereIndexVersion  int                     `json:"2dsphereIndexVersion,omitempty"`
	Bits                    int                     `json:"bits,omitempty"`
	Unique                  bool                    `json:"unique,omitempty"`
	Sparse                  bool                    `json:"sparse,omitempty"`
	GeoMin                  int                     `json:"min,omitempty"`
	GeoMax                  int                     `json:"max,omitempty"`
	BucketSize              int                     `json:"bucketSize,omitempty"`
	Name                    string                  `json:"name,omitempty"`
	ExpireAfterSeconds      int                     `json:"expireAfterSeconds,omitempty"`
}

// CollationOptions represents options for collation indexes.
//
// See: https://docs.mongodb.com/manual/reference/method/db.collection.createIndex/#option-for-collation
type CollationOptions struct {
	Locale          string `json:"locale,omitempty"`
	CaseLevel       bool   `json:"caseLevel,omitempty"`
	CaseFirst       string `json:"caseFirst,omitempty"`
	Strength        int    `json:"strength,omitempty"`
	NumericOrdering bool   `json:"numericOrdering,omitempty"`
	Alternate       string `json:"alternate,omitempty"`
	MaxVariable     string `json:"maxVariable,omitempty"`
	Normalization   bool   `json:"normalization,omitempty"`
	Backwards       bool   `json:"backwards,omitempty"`
}

// Create creates a request for a rolling index creation for the project associated to {GROUP-ID} and the {CLUSTER-NAME}.
//
// See more: https://docs.atlas.mongodb.com/reference/api/rolling-index-create-one/
func (s *IndexesServiceOp) Create(ctx context.Context, groupID, clusterName string, createReq *IndexConfiguration) (*Response, error) {
	if groupID == "" {
		return nil, NewArgError("groupID", "must be set")
	}
	if clusterName == "" {
		return nil, NewArgError("clusterName", "must be set")
	}
	if createReq == nil {
		return nil, NewArgError("createReq", "must be set")
	}

	path := fmt.Sprintf(indexesPath, groupID, clusterName)

	req, err := s.Client.NewRequest(ctx, http.MethodPost, path, createReq)
	if err != nil {
		return nil, err
	}

	resp, err := s.Client.Do(ctx, req, nil)

	return resp, err
}
