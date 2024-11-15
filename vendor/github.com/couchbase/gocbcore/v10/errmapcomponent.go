package gocbcore

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/couchbase/gocbcore/v10/memd"
)

type errMapComponent struct {
	kvErrorMap kvErrorMapPtr
	bucketName string
}

func newErrMapManager(bucketName string) *errMapComponent {
	return &errMapComponent{
		bucketName: bucketName,
	}
}

func (errMgr *errMapComponent) getKvErrMapData(code memd.StatusCode) *kvErrorMapError {
	errMap := errMgr.kvErrorMap.Get()
	if errMap != nil {
		if errData, ok := errMap.Errors[uint16(code)]; ok {
			return &errData
		}
	}
	return nil
}

func (errMgr *errMapComponent) StoreErrorMap(mapBytes []byte) {
	errMap, err := parseKvErrorMap(mapBytes)
	if err != nil {
		logDebugf("Failed to parse kv error map (%s)", err)
		return
	}

	logDebugf("Fetched error map: %+v", errMap)

	// Check if we need to switch the agent itself to a better
	//  error map revision.
	for {
		origMap := errMgr.kvErrorMap.Get()
		if origMap != nil && errMap.Revision < origMap.Revision {
			break
		}

		if errMgr.kvErrorMap.Update(origMap, errMap) {
			break
		}
	}
}

func (errMgr *errMapComponent) ShouldRetry(status memd.StatusCode) bool {
	kvErrData := errMgr.getKvErrMapData(status)
	if kvErrData != nil {
		for _, attr := range kvErrData.Attributes {
			if attr == "auto-retry" || attr == "retry-now" || attr == "retry-later" {
				return true
			}
		}
	}

	return false
}

func (errMgr *errMapComponent) EnhanceKvError(err error, resp *memdQResponse, req *memdQRequest) error {
	enhErr := &KeyValueError{
		InnerError: err,
	}

	if req != nil {
		enhErr.DocumentKey = string(req.Key)
		enhErr.BucketName = errMgr.bucketName
		enhErr.ScopeName = req.ScopeName
		enhErr.CollectionName = req.CollectionName
		enhErr.CollectionID = req.CollectionID

		retryCount, reasons := req.Retries()
		enhErr.RetryReasons = reasons
		enhErr.RetryAttempts = retryCount

		connInfo := req.ConnectionInfo()
		enhErr.LastDispatchedTo = connInfo.lastDispatchedTo
		enhErr.LastDispatchedFrom = connInfo.lastDispatchedFrom
		enhErr.LastConnectionID = connInfo.lastConnectionID
		enhErr.Internal.ResourceUnits = req.ResourceUnits()
	}

	if resp != nil {
		enhErr.StatusCode = resp.Status
		enhErr.Opaque = resp.Opaque

		errMapData := errMgr.getKvErrMapData(enhErr.StatusCode)
		if errMapData != nil {
			var unknownStatusErr *unknownKvStatusCodeError
			if errors.As(err, &unknownStatusErr) {
				enhErr.InnerError = fmt.Errorf("%s (0x%02x)", errMapData.Description, int(resp.Status))
			}

			enhErr.ErrorName = errMapData.Name
			enhErr.ErrorDescription = errMapData.Description
		}

		if memd.DatatypeFlag(resp.Datatype)&memd.DatatypeFlagJSON != 0 {
			var enhancedData struct {
				Error struct {
					Context string `json:"context"`
					Ref     string `json:"ref"`
				} `json:"error"`
			}
			if parseErr := json.Unmarshal(resp.Value, &enhancedData); parseErr == nil {
				enhErr.Context = enhancedData.Error.Context
				enhErr.Ref = enhancedData.Error.Ref
			}
		}
	}

	return enhErr
}

func translateMemdError(err error, req *memdQRequest) error {
	switch err {
	case ErrMemdInvalidArgs:
		return errInvalidArgument
	case ErrMemdInternalError:
		return errInternalServerFailure
	case ErrMemdAccessError:
		return errAuthenticationFailure
	case ErrMemdAuthError:
		return errAuthenticationFailure
	case ErrMemdTmpFail:
		return errTemporaryFailure
	case ErrMemdBusy:
		if req.Command == memd.CmdRangeScanCreate {
			// Range scan create is special cased. We can't change the behaviour for all errors as that
			// would be a breaking change but we need to be able to differentiate busy and tmp fail for create.
			return errBusy
		}
		return errTemporaryFailure
	case ErrMemdKeyExists:
		if req.Command == memd.CmdReplace || (req.Command == memd.CmdDelete && req.Cas != 0) ||
			(req.Command == memd.CmdSubDocMultiMutation && req.Cas != 0) {
			return errCasMismatch
		}
		return errDocumentExists
	case ErrMemdNotStored:
		// GOCBC-1356: memcached does not currently return a NOT_STORED response when inserting a doc, but this was originally
		// the plan, so for safety handle this path.
		if req.Command == memd.CmdAdd {
			return errDocumentExists
		}
		return errNotStored
	case ErrMemdCollectionNotFound:
		return errCollectionNotFound
	case ErrMemdScopeNotFound:
		return errScopeNotFound
	case ErrMemdUnknownCommand:
		return errUnsupportedOperation
	case ErrMemdNotSupported:
		return errUnsupportedOperation
	case ErrMemdDCPStreamIDInvalid:
		return errDCPStreamIDInvalid

	case ErrMemdKeyNotFound:
		return errDocumentNotFound
	case ErrMemdLocked:
		// BUGFIX(brett19): This resolves a bug in the server processing of the LOCKED
		// operation where the server will respond with LOCKED rather than a CAS mismatch.
		if req.Command == memd.CmdUnlockKey {
			return errCasMismatch
		}
		return errDocumentLocked
	case ErrMemdNotLocked:
		return errDocumentNotLocked
	case ErrMemdTooBig:
		return errValueTooLarge
	case ErrMemdSubDocNotJSON:
		return errValueNotJSON
	case ErrMemdDurabilityInvalidLevel:
		return errDurabilityLevelNotAvailable
	case ErrMemdDurabilityImpossible:
		return errDurabilityImpossible
	case ErrMemdSyncWriteAmbiguous:
		return errDurabilityAmbiguous
	case ErrMemdSyncWriteInProgess:
		return errDurableWriteInProgress
	case ErrMemdSyncWriteReCommitInProgress:
		return errDurableWriteReCommitInProgress
	case ErrMemdSubDocPathNotFound:
		return errPathNotFound
	case ErrMemdSubDocPathInvalid:
		return errPathInvalid
	case ErrMemdSubDocPathTooBig:
		return errPathTooBig
	case ErrMemdSubDocDocTooDeep:
		return errPathTooDeep
	case ErrMemdSubDocValueTooDeep:
		return errValueTooDeep
	case ErrMemdSubDocCantInsert:
		return errValueInvalid
	case ErrMemdSubDocNotJSON:
		return errDocumentNotJSON
	case ErrMemdSubDocBadRange:
		return errNumberTooBig
	case ErrMemdSubDocPathMismatch:
		return errPathMismatch
	case ErrMemdBadDelta:
		return errDeltaInvalid
	case ErrMemdSubDocBadDelta:
		return errDeltaInvalid
	case ErrMemdSubDocPathExists:
		return errPathExists
	case ErrXattrUnknownMacro:
		return errXattrUnknownMacro
	case ErrXattrInvalidFlagCombo:
		return errXattrInvalidFlagCombo
	case ErrXattrInvalidKeyCombo:
		return errXattrInvalidKeyCombo
	case ErrMemdSubDocXattrUnknownVAttr:
		return errXattrUnknownVirtualAttribute
	case ErrMemdSubDocXattrCannotModifyVAttr:
		return errXattrCannotModifyVirtualAttribute
	case ErrXattrInvalidOrder:
		return errXattrInvalidOrder
	case ErrMemdNotMyVBucket:
		return errNotMyVBucket
	case ErrMemdRateLimitedNetworkIngress:
		return errRateLimitedFailure
	case ErrMemdRateLimitedNetworkEgress:
		return errRateLimitedFailure
	case ErrMemdRateLimitedMaxConnections:
		return errRateLimitedFailure
	case ErrMemdRateLimitedMaxCommands:
		return errRateLimitedFailure
	case ErrMemdRateLimitedScopeSizeLimitExceeded:
		return errQuotaLimitedFailure
	case ErrMemdRangeScanCancelled:
		return errRangeScanCancelled
	case ErrMemdRangeScanMore:
		return errRangeScanMore
	case ErrMemdRangeScanComplete:
		return errRangeScanComplete
	case ErrMemdRangeScanVbUUIDNotEqual:
		return errRangeScanVbUUIDNotEqual
	}

	return err
}
