// Copyright (c) 2021-2022 Snowflake Computing Inc. All rights reserved.

package gosnowflake

import "regexp"

const (
	awsKeyPattern          = `(?i)(aws_key_id|aws_secret_key|access_key_id|secret_access_key)\s*=\s*'([^']+)'`
	awsTokenPattern        = `(?i)(accessToken|tempToken|keySecret)"\s*:\s*"([a-z0-9/+]{32,}={0,2})"`
	sasTokenPattern        = `(?i)(sig|signature|AWSAccessKeyId|password|passcode)=(?P<secret>[a-z0-9%/+]{16,})`
	privateKeyPattern      = `(?im)-----BEGIN PRIVATE KEY-----\\n([a-z0-9/+=\\n]{32,})\\n-----END PRIVATE KEY-----`
	privateKeyDataPattern  = `(?i)"privateKeyData": "([a-z0-9/+=\\n]{10,})"`
	connectionTokenPattern = `(?i)(token|assertion content)([\'\"\s:=]+)([a-z0-9=/_\-\+]{8,})`
	passwordPattern        = `(?i)(password|pwd)([\'\"\s:=]+)([a-z0-9!\"#\$%&\\\'\(\)\*\+\,-\./:;<=>\?\@\[\]\^_\{\|\}~]{8,})`
)

var (
	awsKeyRegexp          = regexp.MustCompile(awsKeyPattern)
	awsTokenRegexp        = regexp.MustCompile(awsTokenPattern)
	sasTokenRegexp        = regexp.MustCompile(sasTokenPattern)
	privateKeyRegexp      = regexp.MustCompile(privateKeyPattern)
	privateKeyDataRegexp  = regexp.MustCompile(privateKeyDataPattern)
	connectionTokenRegexp = regexp.MustCompile(connectionTokenPattern)
	passwordRegexp        = regexp.MustCompile(passwordPattern)
)

func maskConnectionToken(text string) string {
	return connectionTokenRegexp.ReplaceAllString(text, "$1${2}****")
}

func maskPassword(text string) string {
	return passwordRegexp.ReplaceAllString(text, "$1${2}****")
}

func maskAwsKey(text string) string {
	return awsKeyRegexp.ReplaceAllString(text, "${1}****$2")
}

func maskAwsToken(text string) string {
	return awsTokenRegexp.ReplaceAllString(text, "${1}XXXX$2")
}

func maskSasToken(text string) string {
	return sasTokenRegexp.ReplaceAllString(text, "${1}****$2")
}

func maskPrivateKey(text string) string {
	return privateKeyRegexp.ReplaceAllString(text, "-----BEGIN PRIVATE KEY-----\\\\\\\\nXXXX\\\\\\\\n-----END PRIVATE KEY-----")
}

func maskPrivateKeyData(text string) string {
	return privateKeyDataRegexp.ReplaceAllString(text, `"privateKeyData": "XXXX"`)
}

func maskSecrets(text string) string {
	return maskConnectionToken(
		maskPassword(
			maskPrivateKeyData(
				maskPrivateKey(
					maskAwsToken(
						maskSasToken(
							maskAwsKey(text)))))))
}
