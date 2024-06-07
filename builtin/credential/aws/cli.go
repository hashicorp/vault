// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package awsauth

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/go-secure-stdlib/awsutil"
	"github.com/hashicorp/vault/api"
)

type CLIHandler struct{}

func (h *CLIHandler) Auth(c *api.Client, m map[string]string) (*api.Secret, error) {
	mount, ok := m["mount"]
	if !ok {
		mount = "aws"
	}

	role, ok := m["role"]
	if !ok {
		role = ""
	}

	headerValue, ok := m["header_value"]
	if !ok {
		headerValue = ""
	}

	logVal, ok := m["log_level"]
	if !ok {
		logVal = "info"
	}
	level := hclog.LevelFromString(logVal)
	if level == hclog.NoLevel {
		return nil, fmt.Errorf("failed to parse 'log_level' value: %q", logVal)
	}
	hlogger := hclog.Default()
	hlogger.SetLevel(level)

	creds, err := awsutil.RetrieveCreds(m["aws_access_key_id"], m["aws_secret_access_key"], m["aws_security_token"], hlogger)
	if err != nil {
		return nil, err
	}

	region := m["region"]
	switch region {
	case "":
		// The CLI has always defaulted to "us-east-1" if a region is not provided.
		region = awsutil.DefaultRegion
	case "auto":
		// Beginning in 1.10 we also accept the "auto" value, which uses the region detection logic in
		// awsutil.GetRegion() to determine the region. That behavior is triggered when region = "".
		region = ""
	}

	loginData, err := awsutil.GenerateLoginData(creds, headerValue, region, hlogger)
	if err != nil {
		return nil, err
	}
	if loginData == nil {
		return nil, fmt.Errorf("got nil response from GenerateLoginData")
	}
	loginData["role"] = role
	path := fmt.Sprintf("auth/%s/login", mount)
	secret, err := c.Logical().Write(path, loginData)
	if err != nil {
		return nil, err
	}
	if secret == nil {
		return nil, fmt.Errorf("empty response from credential provider")
	}

	return secret, nil
}

func (h *CLIHandler) Help() string {
	help := `
Usage: vault login -method=aws [CONFIG K=V...]

  The AWS auth method allows users to authenticate with AWS IAM
  credentials. The AWS IAM credentials, and optionally the AWS region, may be 
  specified in a number of ways, listed in order of precedence below:

    1. Explicitly via the command line (not recommended)

    2. Via the standard AWS environment variables (AWS_ACCESS_KEY, etc.)

    3. Via the ~/.aws/credentials file

    4. Via EC2 instance profile

  Authenticate using locally stored credentials:

      $ vault login -method=aws

  Authenticate by passing keys:

      $ vault login -method=aws aws_access_key_id=... aws_secret_access_key=...

Configuration:

  aws_access_key_id=<string>
      Explicit AWS access key ID

  aws_secret_access_key=<string>
      Explicit AWS secret access key

  aws_security_token=<string>
      Explicit AWS security token for temporary credentials

  header_value=<string>
      Value for the x-vault-aws-iam-server-id header in requests

  mount=<string>
      Path where the AWS credential method is mounted. This is usually provided
      via the -path flag in the "vault login" command, but it can be specified
      here as well. If specified here, it takes precedence over the value for
      -path. The default value is "aws".

  region=<string>
      Explicit AWS region to reach out to for authentication request signing. A value
      of "auto" enables auto-detection of region based on the precedence described above.
      Defaults to "us-east-1" if not specified.

  role=<string>
      Name of the role to request a token against

  log_level=<string>
      Set logging level during AWS credential acquisition. Valid levels are
      trace, debug, info, warn, error. Defaults to info.
`

	return strings.TrimSpace(help)
}
