// Copyright 2023 MongoDB Inc
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

const (
	clusterOutageSimulationPath = "api/atlas/v1.0/groups/%s/clusters/%s/outageSimulation"
)

// ClusterOutageSimulationService is an interface for interfacing with the cluster outage endpoints of the MongoDB Atlas API.
//
// See more: https://www.mongodb.com/docs/atlas/reference/api-resources-spec/#tag/Cluster-Outage-Simulation
type ClusterOutageSimulationService interface {
	EndOutageSimulation(context.Context, string, string) (*ClusterOutageSimulation, *Response, error)
	GetOutageSimulation(context.Context, string, string) (*ClusterOutageSimulation, *Response, error)
	StartOutageSimulation(context.Context, string, string, *ClusterOutageSimulationRequest) (*ClusterOutageSimulation, *Response, error)
}

// ClusterOutageSimulationServiceOp handles communication with the ClusterOutageSimulationService related methods of the
// MongoDB Atlas API.
type ClusterOutageSimulationServiceOp service

var _ ClusterOutageSimulationService = &ClusterOutageSimulationServiceOp{}

type ClusterOutageSimulation struct {
	// Human-readable label that identifies the cluster that undergoes outage simulation.
	ClusterName *string `json:"clusterName,omitempty"`
	// Unique 24-hexadecimal character string that identifies the project that contains the cluster to undergo outage simulation.
	GroupID *string `json:"groupId,omitempty"`
	// Unique 24-hexadecimal character string that identifies the outage simulation.
	ID *string `json:"id,omitempty"`
	// List of settings that specify the type of cluster outage simulation.
	OutageFilters []ClusterOutageSimulationOutageFilter `json:"outageFilters,omitempty"`
	// Date and time when MongoDB Cloud started the regional outage simulation.
	StartRequestDate *string `json:"startRequestDate,omitempty"`
	// Phase of the outage simulation.  | State       | Indication | |-------------|------------| | `START_REQUESTED`    | User has requested cluster outage simulation.| | `STARTING`           | MongoDB Cloud is starting cluster outage simulation.| | `SIMULATING`         | MongoDB Cloud is simulating cluster outage.| | `RECOVERY_REQUESTED` | User has requested recovery from the simulated outage.| | `RECOVERING`         | MongoDB Cloud is recovering the cluster from the simulated outage.| | `COMPLETE`           | MongoDB Cloud has completed the cluster outage simulation.|
	State *string `json:"state,omitempty"`
}

type ClusterOutageSimulationRequest struct {
	OutageFilters []ClusterOutageSimulationOutageFilter `json:"outageFilters,omitempty"`
}

type ClusterOutageSimulationOutageFilter struct {
	// The cloud provider of the region that undergoes the outage simulation.
	CloudProvider *string `json:"cloudProvider,omitempty"`
	// The name of the region to undergo an outage simulation.
	RegionName *string `json:"regionName,omitempty"`
	// The type of cluster outage to simulate.  | Type       | Description | |------------|-------------| | `REGION`   | Simulates a cluster outage for a region.|
	Type *string `json:"type,omitempty"`
}

// EndOutageSimulation ends a cluster outage simulation.
//
// See more: https://www.mongodb.com/docs/atlas/reference/api-resources-spec/#tag/Cluster-Outage-Simulation/operation/endOutageSimulation
func (s ClusterOutageSimulationServiceOp) EndOutageSimulation(ctx context.Context, groupID, clusterName string) (*ClusterOutageSimulation, *Response, error) {
	if groupID == "" {
		return nil, nil, NewArgError("groupId", "must be set")
	}
	if clusterName == "" {
		return nil, nil, NewArgError("clusterName", "must be set")
	}

	path := fmt.Sprintf(clusterOutageSimulationPath, groupID, clusterName)

	req, err := s.Client.NewRequest(ctx, http.MethodDelete, path, nil)
	if err != nil {
		return nil, nil, err
	}

	root := new(ClusterOutageSimulation)
	resp, err := s.Client.Do(ctx, req, root)
	if err != nil {
		return nil, resp, err
	}

	return root, resp, err
}

// GetOutageSimulation returns one outage simulation for one cluster.
//
// See more: https://www.mongodb.com/docs/atlas/reference/api-resources-spec/#tag/Cluster-Outage-Simulation/operation/getOutageSimulation
func (s ClusterOutageSimulationServiceOp) GetOutageSimulation(ctx context.Context, groupID, clusterName string) (*ClusterOutageSimulation, *Response, error) {
	if groupID == "" {
		return nil, nil, NewArgError("groupId", "must be set")
	}
	if clusterName == "" {
		return nil, nil, NewArgError("clusterName", "must be set")
	}

	path := fmt.Sprintf(clusterOutageSimulationPath, groupID, clusterName)

	req, err := s.Client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, nil, err
	}

	root := new(ClusterOutageSimulation)
	resp, err := s.Client.Do(ctx, req, root)
	if err != nil {
		return nil, resp, err
	}

	return root, resp, err
}

// StartOutageSimulation starts a cluster outage simulation.
//
// See more:https://www.mongodb.com/docs/atlas/reference/api-resources-spec/#tag/Cluster-Outage-Simulation/operation/startOutageSimulation
func (s ClusterOutageSimulationServiceOp) StartOutageSimulation(ctx context.Context, groupID, clusterName string, request *ClusterOutageSimulationRequest) (*ClusterOutageSimulation, *Response, error) {
	if groupID == "" {
		return nil, nil, NewArgError("groupId", "must be set")
	}
	if clusterName == "" {
		return nil, nil, NewArgError("clusterName", "must be set")
	}
	if request == nil {
		return nil, nil, NewArgError("request", "cannot be nil")
	}

	path := fmt.Sprintf(clusterOutageSimulationPath, groupID, clusterName)
	fmt.Println(path)
	req, err := s.Client.NewRequest(ctx, http.MethodPost, path, request)
	if err != nil {
		return nil, nil, err
	}

	root := new(ClusterOutageSimulation)
	resp, err := s.Client.Do(ctx, req, root)
	if err != nil {
		return nil, resp, err
	}

	return root, resp, err
}
