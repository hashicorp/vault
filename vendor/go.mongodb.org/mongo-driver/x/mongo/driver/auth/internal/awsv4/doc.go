// Copyright (C) MongoDB, Inc. 2017-present.
//
// Licensed under the Apache License, Version 2.0 (the "License"); you may
// not use this file except in compliance with the License. You may obtain
// a copy of the License at http://www.apache.org/licenses/LICENSE-2.0
//
// Based on github.com/aws/aws-sdk-go v1.34.28 by Amazon.com, Inc.
// See THIRD-PARTY-NOTICES for original license terms

// Package awsv4 implements signing for AWS V4 signer with static credentials,
// and is based on and modified from code in the package aws-sdk-go. The
// modifications remove non-static credentials, support for non-sts services,
// and the options for v4.Signer. They also reduce the number of non-Go
// library dependencies.
package awsv4
