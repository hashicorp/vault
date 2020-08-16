// Copyright (c) 2018 Yandex LLC. All rights reserved.
// Author: Maxim Kolganov <manykey@yandex-team.ru>

package ycsdk

import (
	"github.com/yandex-cloud/go-sdk/gen/containerregistry"
)

const (
	ContainerRegistryServiceID Endpoint = "container-registry"
)

// ContainerRegistry returns ContainerRegistry object that is used to operate on Yandex Container Registry
func (sdk *SDK) ContainerRegistry() *containerregistry.ContainerRegistry {
	return containerregistry.NewContainerRegistry(sdk.getConn(ContainerRegistryServiceID))
}
