// Copyright (C) MongoDB, Inc. 2017-present.
//
// Licensed under the Apache License, Version 2.0 (the "License"); you may
// not use this file except in compliance with the License. You may obtain
// a copy of the License at http://www.apache.org/licenses/LICENSE-2.0

package auth

// Cred is a user's credential.
type Cred struct {
	Source      string
	Username    string
	Password    string
	PasswordSet bool
	Props       map[string]string
}
