// Copyright 2024 Google LLC
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

package debug

import "context"

// Logger is the interface used for debug logging. By default, it is unused.
//
// Prefer ContextLogger instead.
type Logger interface {
	// Debugf is for reporting information about internal operations.
	Debugf(format string, args ...interface{})
}

// ContextLogger is the interface used for debug logging. By default, it is unused.
type ContextLogger interface {
	// Debugf is for reporting information about internal operations.
	Debugf(ctx context.Context, format string, args ...interface{})
}
