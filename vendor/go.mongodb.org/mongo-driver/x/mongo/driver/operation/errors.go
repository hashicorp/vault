// Copyright (C) MongoDB, Inc. 2017-present.
//
// Licensed under the Apache License, Version 2.0 (the "License"); you may
// not use this file except in compliance with the License. You may obtain
// a copy of the License at http://www.apache.org/licenses/LICENSE-2.0

package operation

import "errors"

var (
	errUnacknowledgedHint = errors.New("the 'hint' command parameter cannot be used with unacknowledged writes")
)
