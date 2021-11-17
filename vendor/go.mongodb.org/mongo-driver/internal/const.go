// Copyright (C) MongoDB, Inc. 2017-present.
//
// Licensed under the Apache License, Version 2.0 (the "License"); you may
// not use this file except in compliance with the License. You may obtain
// a copy of the License at http://www.apache.org/licenses/LICENSE-2.0

package internal // import "go.mongodb.org/mongo-driver/internal"

// Version is the current version of the driver.
var Version = "local build"

// SetMockServiceID enables a mode in which the driver mocks server support for returning a "serviceId" field in "hello"
// command responses by using the value of "topologyVersion.processId".  This is used for testing load balancer support
// until an upstream service can support running behind a load balancer.
var SetMockServiceID = false
