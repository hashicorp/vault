// Copyright 2018 by David A. Golden. All rights reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License"); you may
// not use this file except in compliance with the License. You may obtain
// a copy of the License at http://www.apache.org/licenses/LICENSE-2.0

// Package scram provides client and server implementations of the Salted
// Challenge Response Authentication Mechanism (SCRAM) described in RFC-5802
// and RFC-7677.
//
// Usage
//
// The scram package provides variables, `SHA1`, `SHA256`, and `SHA512`, that
// are used to construct Client or Server objects.
//
//     clientSHA1,   err := scram.SHA1.NewClient(username, password, authID)
//     clientSHA256, err := scram.SHA256.NewClient(username, password, authID)
//     clientSHA512, err := scram.SHA512.NewClient(username, password, authID)
//
//     serverSHA1,   err := scram.SHA1.NewServer(credentialLookupFcn)
//     serverSHA256, err := scram.SHA256.NewServer(credentialLookupFcn)
//     serverSHA512, err := scram.SHA512.NewServer(credentialLookupFcn)
//
// These objects are used to construct ClientConversation or
// ServerConversation objects that are used to carry out authentication.
package scram
