package ycsdk

import "github.com/yandex-cloud/go-sdk/gen/compute/instancegroup"

const InstancegroupServiceID Endpoint = "compute"

// InstanceGroup returns InstanceGroup object that is used to operate on Yandex Compute InstanceGroup
func (sdk *SDK) InstanceGroup() *instancegroup.InstanceGroup {
	return instancegroup.NewInstanceGroup(sdk.getConn(InstancegroupServiceID))
}
