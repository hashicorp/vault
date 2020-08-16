// Copyright (c) 2019 YANDEX LLC.

package ycsdk

import "github.com/yandex-cloud/go-sdk/gen/loadbalancer"

const (
	LoadBalancerServiceID Endpoint = "load-balancer"
)

// LoadBalancer returns LoadBalancer object that is used to operate on load balancers
func (sdk *SDK) LoadBalancer() *loadbalancer.LoadBalancer {
	return loadbalancer.NewLoadBalancer(sdk.getConn(LoadBalancerServiceID))
}
