// Copyright (c) 2018 Yandex LLC. All rights reserved.
// Author: Vladimir Skipor <skipor@yandex-team.ru>

package iam

import (
	"github.com/yandex-cloud/go-sdk/gen/iam/awscompatibility"
)

func (i *IAM) AWSCompatibility() *awscompatibility.AWSCompatibility {
	return awscompatibility.NewAWSCompatibility(i.getConn)
}
