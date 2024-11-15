// Copyright 2021 Couchbase
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     https://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package gocbcore

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

// TransactionsProtocolVersion returns the protocol version that this library supports.
func TransactionsProtocolVersion() string {
	return "2.1"
}

// TransactionsProtocolExtensions returns a list strings representing the various features
// that this specific version of the library supports within its protocol version.
func TransactionsProtocolExtensions() []string {
	return []string{
		"EXT_TRANSACTION_ID",
		"EXT_MEMORY_OPT_UNSTAGING",
		"EXT_BINARY_METADATA",
		"EXT_CUSTOM_METADATA_COLLECTION",
		"EXT_STORE_DURABILITY",
		"EXT_REMOVE_COMPLETED",
		"EXT_ALL_KV_COMBINATIONS",
		"EXT_UNKNOWN_ATR_STATES",
		"BF_CBD_3787",
		"BF_CBD_3705",
		"BF_CBD_3838",
		"BF_CBD_3791",
		"BF_CBD_3794",
		"EXT_QUERY",
		"EXT_SDK_INTEGRATION",
		"EXT_SINGLE_QUERY",
		"EXT_INSERT_EXISTING",
		"EXT_QUERY_CONTEXT",
	}
}

type forwardCompatBehaviour string

// nolint: deadcode,varcheck
const (
	forwardCompatBehaviourRetry forwardCompatBehaviour = "r"
	forwardCompatBehaviourFail  forwardCompatBehaviour = "f"
)

type forwardCompatExtension string

// nolint: deadcode,varcheck
const (
	forwardCompatExtensionTransactionID            forwardCompatExtension = "TI"
	forwardCompatExtensionDeferredCommit           forwardCompatExtension = "DC"
	forwardCompatExtensionTimeOptUnstaging         forwardCompatExtension = "TO"
	forwardCompatExtensionMemoryOptUnstaging       forwardCompatExtension = "MO"
	forwardCompatExtensionCustomMetadataCollection forwardCompatExtension = "CM"
	forwardCompatExtensionBinaryMetadata           forwardCompatExtension = "BM"
	forwardCompatExtensionQuery                    forwardCompatExtension = "QU"
	forwardCompatExtensionStoreDurability          forwardCompatExtension = "SD"
	forwardCompatExtensionRemoveCompleted          forwardCompatExtension = "RC"
	forwardCompatExtensionAllKvCombinations        forwardCompatExtension = "CO"
	forwardCompatExtensionUnknownATRStates         forwardCompatExtension = "UA"
	forwardCompatExtensionBFCBD3787                forwardCompatExtension = "BF3787"
	forwardCompatExtensionBFCBD3705                forwardCompatExtension = "BF3705"
	forwardCompatExtensionBFCBD3838                forwardCompatExtension = "BF3838"
	forwardCompatExtensionBFCBD3791                forwardCompatExtension = "BF3791"
	forwardCompatExtensionBFCBD3794                forwardCompatExtension = "BF3794"
	forwardCompatExtensionSDKIntegration           forwardCompatExtension = "SI"
	forwardCompatExtensionSingleQuery              forwardCompatExtension = "SQ"
	forwardCompatExtensionInsertExisting           forwardCompatExtension = "IX"
	forwardCompatExtensionQueryContext             forwardCompatExtension = "QC"
)

type forwardCompatStage string

// nolint: deadcode,varcheck
const (
	forwardCompatStageWWCReadingATR    forwardCompatStage = "WW_R"
	forwardCompatStageWWCReplacing     forwardCompatStage = "WW_RP"
	forwardCompatStageWWCRemoving      forwardCompatStage = "WW_RM"
	forwardCompatStageWWCInserting     forwardCompatStage = "WW_I"
	forwardCompatStageWWCInsertingGet  forwardCompatStage = "WW_IG"
	forwardCompatStageGets             forwardCompatStage = "G"
	forwardCompatStageGetsReadingATR   forwardCompatStage = "G_A"
	forwardCompatStageGetsCleanupEntry forwardCompatStage = "CL_E"
)

const (
	protocolMajor = 2
	protocolMinor = 0
)

// TransactionForwardCompatibilityEntry represents a forward compatibility entry.
// Internal: This should never be used and is not supported.
type TransactionForwardCompatibilityEntry struct {
	ProtocolVersion   string `json:"p,omitempty"`
	ProtocolExtension string `json:"e,omitempty"`
	Behaviour         string `json:"b,omitempty"`
	RetryInterval     int    `json:"ra,omitempty"`
}

var supportedforwardCompatExtensions = []forwardCompatExtension{
	forwardCompatExtensionTransactionID,
	forwardCompatExtensionMemoryOptUnstaging,
	forwardCompatExtensionCustomMetadataCollection,
	forwardCompatExtensionBinaryMetadata,
	forwardCompatExtensionQuery,
	forwardCompatExtensionStoreDurability,
	forwardCompatExtensionRemoveCompleted,
	forwardCompatExtensionAllKvCombinations,
	forwardCompatExtensionUnknownATRStates,
	forwardCompatExtensionBFCBD3787,
	forwardCompatExtensionBFCBD3705,
	forwardCompatExtensionBFCBD3838,
	forwardCompatExtensionBFCBD3791,
	forwardCompatExtensionBFCBD3794,
	forwardCompatExtensionSDKIntegration,
	forwardCompatExtensionSingleQuery,
	forwardCompatExtensionInsertExisting,
	forwardCompatExtensionQueryContext,
}

func jsonForwardCompatToForwardCompat(fc map[string][]jsonForwardCompatibilityEntry) map[string][]TransactionForwardCompatibilityEntry {
	if fc == nil {
		return nil
	}
	forwardCompat := make(map[string][]TransactionForwardCompatibilityEntry)

	for k, entries := range fc {
		if _, ok := forwardCompat[k]; !ok {
			forwardCompat[k] = make([]TransactionForwardCompatibilityEntry, len(entries))
		}

		for i, entry := range entries {
			forwardCompat[k][i] = TransactionForwardCompatibilityEntry(entry)
		}
	}

	return forwardCompat
}

func checkForwardCompatProtocol(protocolVersion string) (bool, error) {
	if protocolVersion == "" {
		return false, nil
	}

	protocol := strings.Split(protocolVersion, ".")
	if len(protocol) != 2 {
		return false, fmt.Errorf("invalid protocol string: %s", protocolVersion)
	}
	major, err := strconv.Atoi(protocol[0])
	if err != nil {
		return false, wrapError(err, fmt.Sprintf("invalid protocol string: %s", protocolVersion))
	}
	if protocolMajor < major {
		return false, nil
	}
	if protocolMajor == major {
		minor, err := strconv.Atoi(protocol[1])
		if err != nil {
			return false, wrapError(err, fmt.Sprintf("invalid protocol string: %s", protocolVersion))
		}
		if protocolMinor < minor {
			return false, nil
		}
	}

	return true, nil
}

func checkForwardCompatExtension(extension string) bool {
	if extension == "" {
		return false
	}

	for _, supported := range supportedforwardCompatExtensions {
		if string(supported) == extension {
			return true
		}
	}

	return false
}

func checkForwardCompatability(
	stage forwardCompatStage,
	fc map[string][]TransactionForwardCompatibilityEntry,
) (isCompatOut bool, shouldRetryOut bool, retryWaitOut time.Duration, errOut error) {
	if len(fc) == 0 {
		return true, false, 0, nil
	}

	if checks, ok := fc[string(stage)]; ok {
		for _, c := range checks {
			protocolOk, err := checkForwardCompatProtocol(c.ProtocolVersion)
			if err != nil {
				return false, false, 0, err
			}

			if protocolOk {
				continue
			}

			if extensionOk := checkForwardCompatExtension(c.ProtocolExtension); extensionOk {
				continue
			}

			// If we get here then neither protocol or extension are ok.
			switch forwardCompatBehaviour(c.Behaviour) {
			case forwardCompatBehaviourRetry:
				retryWait := time.Duration(c.RetryInterval) * time.Millisecond
				return false, true, retryWait, nil
			default:
				return false, false, 0, nil
			}
		}
	}

	return true, false, 0, nil
}
