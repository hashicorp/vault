// Copyright (C) MongoDB, Inc. 2017-present.
//
// Licensed under the Apache License, Version 2.0 (the "License"); you may
// not use this file except in compliance with the License. You may obtain
// a copy of the License at http://www.apache.org/licenses/LICENSE-2.0

// Package auth is not for public use.
//
// The API for packages in the 'private' directory have no stability
// guarantee.
//
// The packages within the 'private' directory would normally be put into an
// 'internal' directory to prohibit their use outside the 'mongo' directory.
// However, some MongoDB tools require very low-level access to the building
// blocks of a driver, so we have placed them under 'private' to allow these
// packages to be imported by projects that need them.
//
// These package APIs may be modified in backwards-incompatible ways at any
// time.
//
// You are strongly discouraged from directly using any packages
// under 'private'.
package auth
