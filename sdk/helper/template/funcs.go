// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package template

import (
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"strconv"
	"strings"
	"time"

	UUID "github.com/hashicorp/go-uuid"
)

func unixTime() string {
	return strconv.FormatInt(time.Now().Unix(), 10)
}

func unixTimeMillis() string {
	return strconv.FormatInt(time.Now().UnixNano()/int64(time.Millisecond), 10)
}

func timestamp(format string) string {
	return time.Now().Format(format)
}

func truncate(maxLen int, str string) (string, error) {
	if maxLen <= 0 {
		return "", fmt.Errorf("max length must be > 0 but was %d", maxLen)
	}
	if len(str) > maxLen {
		return str[:maxLen], nil
	}
	return str, nil
}

const (
	sha256HashLen = 8
)

func truncateSHA256(maxLen int, str string) (string, error) {
	if maxLen <= 8 {
		return "", fmt.Errorf("max length must be > 8 but was %d", maxLen)
	}

	if len(str) <= maxLen {
		return str, nil
	}

	truncIndex := maxLen - sha256HashLen
	hash := hashSHA256(str[truncIndex:])
	result := fmt.Sprintf("%s%s", str[:truncIndex], hash[:sha256HashLen])
	return result, nil
}

func hashSHA256(str string) string {
	return fmt.Sprintf("%x", sha256.Sum256([]byte(str)))
}

func encodeBase64(str string) string {
	return base64.StdEncoding.EncodeToString([]byte(str))
}

func uppercase(str string) string {
	return strings.ToUpper(str)
}

func lowercase(str string) string {
	return strings.ToLower(str)
}

func replace(find string, replace string, str string) string {
	return strings.ReplaceAll(str, find, replace)
}

func uuid() (string, error) {
	return UUID.GenerateUUID()
}
