/*
 *
 * Copyright 2019 gRPC authors.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 *
 */

/*
Package grpcgcp provides grpc supports for Google Cloud APIs.
For now it provides connection management with affinity support.

Note: "channel" is analagous to "connection" in our context.

Usage:

1. First, initialize the api configuration. There are two ways:

	1a. Create a json file defining the configuration and read it.

		// Create some_api_config.json
		{
			"channelPool": {
				"maxSize": 4,
				"maxConcurrentStreamsLowWatermark": 50
			},
			"method": [
				{
					"name": [ "/some.api.v1/Method1" ],
					"affinity": {
						"command": "BIND",
						"affinityKey": "key1"
					}
				},
				{
					"name": [ "/some.api.v1/Method2" ],
					"affinity": {
						"command": "BOUND",
						"affinityKey": "key2"
					}
				},
				{
					"name": [ "/some.api.v1/Method3" ],
					"affinity": {
						"command": "UNBIND",
						"affinityKey": "key3"
					}
				}
			]
		}

		jsonFile, err := ioutil.ReadFile("some_api_config.json")
		if err != nil {
			t.Fatalf("Failed to read config file: %v", err)
		}
		jsonCfg := string(jsonFile)

	1b. Create apiConfig directly and convert it to json.

		// import (
		// 	configpb "github.com/GoogleCloudPlatform/grpc-gcp-go/grpcgcp/grpc_gcp"
		// )

		apiConfig := &configpb.ApiConfig{
			ChannelPool: &configpb.ChannelPoolConfig{
				MaxSize:                          4,
				MaxConcurrentStreamsLowWatermark: 50,
			},
			Method: []*configpb.MethodConfig{
				&configpb.MethodConfig{
					Name: []string{"/some.api.v1/Method1"},
					Affinity: &configpb.AffinityConfig{
						Command:     configpb.AffinityConfig_BIND,
						AffinityKey: "key1",
					},
				},
				&configpb.MethodConfig{
					Name: []string{"/some.api.v1/Method2"},
					Affinity: &configpb.AffinityConfig{
						Command:     configpb.AffinityConfig_BOUND,
						AffinityKey: "key2",
					},
				},
				&configpb.MethodConfig{
					Name: []string{"/some.api.v1/Method3"},
					Affinity: &configpb.AffinityConfig{
						Command:     configpb.AffinityConfig_UNBIND,
						AffinityKey: "key3",
					},
				},
			},
		}

		c, err := protojson.Marshal(apiConfig)
		if err != nil {
			t.Fatalf("cannot json encode config: %v", err)
		}
		jsonCfg := string(c)

2. Make ClientConn with specific DialOptions to enable grpc_gcp load balancer
with provided configuration. And specify gRPC-GCP interceptors.

	conn, err := grpc.Dial(
		target,
		// Register and specify the grpc-gcp load balancer.
		grpc.WithDisableServiceConfig(),
		grpc.WithDefaultServiceConfig(
			fmt.Sprintf(
				`{"loadBalancingConfig": [{"%s":%s}]}`,
				grpcgcp.Name,
				jsonCfg,
			),
		),
		// Set grpcgcp interceptors.
		grpc.WithUnaryInterceptor(grpcgcp.GCPUnaryClientInterceptor),
		grpc.WithStreamInterceptor(grpcgcp.GCPStreamClientInterceptor),
	)
*/
package grpcgcp // import "github.com/GoogleCloudPlatform/grpc-gcp-go/grpcgcp"
