package configutil

import (
	"bytes"
	"context"
	"encoding/base64"
	"errors"
	"fmt"
	"regexp"

	wrapping "github.com/hashicorp/go-kms-wrapping"
	"google.golang.org/protobuf/proto"
)

var (
	encryptRegex = regexp.MustCompile(`{{encrypt\(.*\)}}`)
	decryptRegex = regexp.MustCompile(`{{decrypt\(.*\)}}`)
)

func EncryptDecrypt(rawStr string, decrypt, strip bool, wrapper wrapping.Wrapper) (string, error) {
	var locs [][]int
	raw := []byte(rawStr)
	searchVal := "{{encrypt("
	replaceVal := "{{decrypt("
	suffixVal := ")}}"
	if decrypt {
		searchVal = "{{decrypt("
		replaceVal = "{{encrypt("
		locs = decryptRegex.FindAllIndex(raw, -1)
	} else {
		locs = encryptRegex.FindAllIndex(raw, -1)
	}
	if strip {
		replaceVal = ""
		suffixVal = ""
	}

	out := make([]byte, 0, len(rawStr)*2)
	var prevMaxLoc int
	for _, match := range locs {
		if len(match) != 2 {
			return "", fmt.Errorf("expected two values for match, got %d", len(match))
		}

		// Append everything from the end of the last match to the beginning of this one
		out = append(out, raw[prevMaxLoc:match[0]]...)

		// Transform. First pull off the suffix/prefix
		matchBytes := raw[match[0]:match[1]]
		matchBytes = bytes.TrimSuffix(bytes.TrimPrefix(matchBytes, []byte(searchVal)), []byte(")}}"))
		var finalVal string

		// Now encrypt or decrypt
		switch decrypt {
		case false:
			outBlob, err := wrapper.Encrypt(context.Background(), matchBytes, nil)
			if err != nil {
				return "", fmt.Errorf("error encrypting parameter: %w", err)
			}
			if outBlob == nil {
				return "", errors.New("nil value returned from encrypting parameter")
			}
			outMsg, err := proto.Marshal(outBlob)
			if err != nil {
				return "", fmt.Errorf("error marshaling encrypted parameter: %w", err)
			}
			finalVal = base64.RawURLEncoding.EncodeToString(outMsg)

		default:
			inMsg, err := base64.RawURLEncoding.DecodeString(string(matchBytes))
			if err != nil {
				return "", fmt.Errorf("error decoding encrypted parameter: %w", err)
			}
			inBlob := new(wrapping.EncryptedBlobInfo)
			if err := proto.Unmarshal(inMsg, inBlob); err != nil {
				return "", fmt.Errorf("error unmarshaling encrypted parameter: %w", err)
			}
			dec, err := wrapper.Decrypt(context.Background(), inBlob, nil)
			if err != nil {
				return "", fmt.Errorf("error decrypting encrypted parameter: %w", err)
			}
			finalVal = string(dec)
		}

		// Append new value
		out = append(out, []byte(fmt.Sprintf("%s%s%s", replaceVal, finalVal, suffixVal))...)
		prevMaxLoc = match[1]
	}
	// At the end, append the rest
	out = append(out, raw[prevMaxLoc:]...)
	return string(out), nil
}
