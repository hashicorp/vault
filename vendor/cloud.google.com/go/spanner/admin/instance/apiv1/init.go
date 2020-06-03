/*
Copyright 2020 Google LLC

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package instance

import (
	"context"
	"os"

	"google.golang.org/api/option"
	"google.golang.org/grpc"
)

func init() {
	newInstanceAdminClientHook = func(ctx context.Context, p clientHookParams) ([]option.ClientOption, error) {
		if emulator := os.Getenv("SPANNER_EMULATOR_HOST"); emulator != "" {
			return []option.ClientOption{
				option.WithEndpoint(emulator),
				option.WithGRPCDialOption(grpc.WithInsecure()),
				option.WithoutAuthentication(),
			}, nil
		}

		return nil, nil
	}
}
