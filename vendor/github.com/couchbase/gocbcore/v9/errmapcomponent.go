package gocbcore

import (
	"encoding/json"

	"github.com/couchbase/gocbcore/v9/memd"
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
	}

	if resp != nil {
		enhErr.StatusCode = resp.Status
		enhErr.Opaque = resp.Opaque

		errMapData := errMgr.getKvErrMapData(enhErr.StatusCode)
		if errMapData != nil {
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
		return errTemporaryFailure
	case ErrMemdKeyExists:
		if req.Command == memd.CmdReplace || (req.Command == memd.CmdDelete && req.Cas != 0) {
			return errCasMismatch
		}
		return errDocumentExists
	case ErrMemdCollectionNotFound:
		return errCollectionNotFound
	case ErrMemdUnknownCommand:
		return errUnsupportedOperation
	case ErrMemdNotSupported:
		return errUnsupportedOperation

	case ErrMemdKeyNotFound:
		return errDocumentNotFound
	case ErrMemdLocked:
		// BUGFIX(brett19): This resolves a bug in the server processing of the LOCKED
		// operation where the server will respond with LOCKED rather than a CAS mismatch.
		if req.Command == memd.CmdUnlockKey {
			return errCasMismatch
		}
		return errDocumentLocked
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
	}

	return err
}
