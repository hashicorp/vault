// Copyright (C) MongoDB, Inc. 2017-present.
//
// Licensed under the Apache License, Version 2.0 (the "License"); you may
// not use this file except in compliance with the License. You may obtain
// a copy of the License at http://www.apache.org/licenses/LICENSE-2.0

package internal

import "errors"

var (
	// ErrLoadBalancedWithMultipleHosts is returned when loadBalanced=true is specified in a URI with multiple hosts.
	ErrLoadBalancedWithMultipleHosts = errors.New("loadBalanced cannot be set to true if multiple hosts are specified")
	// ErrLoadBalancedWithReplicaSet is returned when loadBalanced=true is specified in a URI with the replicaSet option.
	ErrLoadBalancedWithReplicaSet = errors.New("loadBalanced cannot be set to true if a replica set name is specified")
	// ErrLoadBalancedWithDirectConnection is returned when loadBalanced=true is specified in a URI with the directConnection option.
	ErrLoadBalancedWithDirectConnection = errors.New("loadBalanced cannot be set to true if the direct connection option is specified")
)
