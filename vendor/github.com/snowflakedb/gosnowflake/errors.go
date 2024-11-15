// Copyright (c) 2017-2022 Snowflake Computing Inc. All rights reserved.

package gosnowflake

import (
	"fmt"
	"runtime/debug"
	"strconv"
	"time"
)

// SnowflakeError is a error type including various Snowflake specific information.
type SnowflakeError struct {
	Number         int
	SQLState       string
	QueryID        string
	Message        string
	MessageArgs    []interface{}
	IncludeQueryID bool // TODO: populate this in connection
}

func (se *SnowflakeError) Error() string {
	message := se.Message
	if len(se.MessageArgs) > 0 {
		message = fmt.Sprintf(se.Message, se.MessageArgs...)
	}
	if se.SQLState != "" {
		if se.IncludeQueryID {
			return fmt.Sprintf("%06d (%s): %s: %s", se.Number, se.SQLState, se.QueryID, message)
		}
		return fmt.Sprintf("%06d (%s): %s", se.Number, se.SQLState, message)
	}
	if se.IncludeQueryID {
		return fmt.Sprintf("%06d: %s: %s", se.Number, se.QueryID, message)
	}
	return fmt.Sprintf("%06d: %s", se.Number, message)
}

func (se *SnowflakeError) generateTelemetryExceptionData() *telemetryData {
	data := &telemetryData{
		Message: map[string]string{
			typeKey:          sqlException,
			sourceKey:        telemetrySource,
			driverTypeKey:    "Go",
			driverVersionKey: SnowflakeGoDriverVersion,
			stacktraceKey:    maskSecrets(string(debug.Stack())),
		},
		Timestamp: time.Now().UnixNano() / int64(time.Millisecond),
	}
	if se.QueryID != "" {
		data.Message[queryIDKey] = se.QueryID
	}
	if se.SQLState != "" {
		data.Message[sqlStateKey] = se.SQLState
	}
	if se.Message != "" {
		data.Message[reasonKey] = se.Message
	}
	if len(se.MessageArgs) > 0 {
		data.Message[reasonKey] = fmt.Sprintf(se.Message, se.MessageArgs...)
	}
	if se.Number != 0 {
		data.Message[errorNumberKey] = strconv.Itoa(se.Number)
	}
	return data
}

func (se *SnowflakeError) sendExceptionTelemetry(sc *snowflakeConn, data *telemetryData) error {
	if sc != nil && sc.telemetry != nil {
		return sc.telemetry.addLog(data)
	}
	return nil // TODO oob telemetry
}

func (se *SnowflakeError) exceptionTelemetry(sc *snowflakeConn) *SnowflakeError {
	data := se.generateTelemetryExceptionData()
	if err := se.sendExceptionTelemetry(sc, data); err != nil {
		logger.WithContext(sc.ctx).Debugf("failed to log to telemetry: %v", data)
	}
	return se
}

// return populated error fields replacing the default response
func populateErrorFields(code int, data *execResponse) *SnowflakeError {
	err := errUnknownError()
	if code != -1 {
		err.Number = code
	}
	if data.Data.SQLState != "" {
		err.SQLState = data.Data.SQLState
	}
	if data.Message != "" {
		err.Message = data.Message
	}
	if data.Data.QueryID != "" {
		err.QueryID = data.Data.QueryID
	}
	return err
}

const (
	/* connection */

	// ErrCodeEmptyAccountCode is an error code for the case where a DNS doesn't include account parameter
	ErrCodeEmptyAccountCode = 260000
	// ErrCodeEmptyUsernameCode is an error code for the case where a DNS doesn't include user parameter
	ErrCodeEmptyUsernameCode = 260001
	// ErrCodeEmptyPasswordCode is an error code for the case where a DNS doesn't include password parameter
	ErrCodeEmptyPasswordCode = 260002
	// ErrCodeFailedToParseHost is an error code for the case where a DNS includes an invalid host name
	ErrCodeFailedToParseHost = 260003
	// ErrCodeFailedToParsePort is an error code for the case where a DNS includes an invalid port number
	ErrCodeFailedToParsePort = 260004
	// ErrCodeIdpConnectionError is an error code for the case where a IDP connection failed
	ErrCodeIdpConnectionError = 260005
	// ErrCodeSSOURLNotMatch is an error code for the case where a SSO URL doesn't match
	ErrCodeSSOURLNotMatch = 260006
	// ErrCodeServiceUnavailable is an error code for the case where service is unavailable.
	ErrCodeServiceUnavailable = 260007
	// ErrCodeFailedToConnect is an error code for the case where a DB connection failed due to wrong account name
	ErrCodeFailedToConnect = 260008
	// ErrCodeRegionOverlap is an error code for the case where a region is specified despite an account region present
	ErrCodeRegionOverlap = 260009
	// ErrCodePrivateKeyParseError is an error code for the case where the private key is not parsed correctly
	ErrCodePrivateKeyParseError = 260010
	// ErrCodeFailedToParseAuthenticator is an error code for the case where a DNS includes an invalid authenticator
	ErrCodeFailedToParseAuthenticator = 260011
	// ErrCodeClientConfigFailed is an error code for the case where clientConfigFile is invalid or applying client configuration fails
	ErrCodeClientConfigFailed = 260012
	// ErrCodeTomlFileParsingFailed is an error code for the case where parsing the toml file is failed because of invalid value.
	ErrCodeTomlFileParsingFailed = 260013
	// ErrCodeFailedToFindDSNInToml is an error code for the case where the DSN does not exist in the toml file.
	ErrCodeFailedToFindDSNInToml = 260014
	// ErrCodeInvalidFilePermission is an error code for the case where the user does not have 0600 permission to the toml file .
	ErrCodeInvalidFilePermission = 260015

	/* network */

	// ErrFailedToPostQuery is an error code for the case where HTTP POST failed.
	ErrFailedToPostQuery = 261000
	// ErrFailedToRenewSession is an error code for the case where session renewal failed.
	ErrFailedToRenewSession = 261001
	// ErrFailedToCancelQuery is an error code for the case where cancel query failed.
	ErrFailedToCancelQuery = 261002
	// ErrFailedToCloseSession is an error code for the case where close session failed.
	ErrFailedToCloseSession = 261003
	// ErrFailedToAuth is an error code for the case where authentication failed for unknown reason.
	ErrFailedToAuth = 261004
	// ErrFailedToAuthSAML is an error code for the case where authentication via SAML failed for unknown reason.
	ErrFailedToAuthSAML = 261005
	// ErrFailedToAuthOKTA is an error code for the case where authentication via OKTA failed for unknown reason.
	ErrFailedToAuthOKTA = 261006
	// ErrFailedToGetSSO is an error code for the case where authentication via OKTA failed for unknown reason.
	ErrFailedToGetSSO = 261007
	// ErrFailedToParseResponse is an error code for when we cannot parse an external browser response from Snowflake.
	ErrFailedToParseResponse = 261008
	// ErrFailedToGetExternalBrowserResponse is an error code for when there's an error reading from the open socket.
	ErrFailedToGetExternalBrowserResponse = 261009
	// ErrFailedToHeartbeat is an error code when a heartbeat fails.
	ErrFailedToHeartbeat = 261010

	/* rows */

	// ErrFailedToGetChunk is an error code for the case where it failed to get chunk of result set
	ErrFailedToGetChunk = 262000

	/* transaction*/

	// ErrNoReadOnlyTransaction is an error code for the case where readonly mode is specified.
	ErrNoReadOnlyTransaction = 263000
	// ErrNoDefaultTransactionIsolationLevel is an error code for the case where non default isolation level is specified.
	ErrNoDefaultTransactionIsolationLevel = 263001

	/* file transfer */

	// ErrInvalidStageFs is an error code denoting an invalid stage in the file system
	ErrInvalidStageFs = 264001
	// ErrFailedToDownloadFromStage is an error code denoting the failure to download a file from the stage
	ErrFailedToDownloadFromStage = 264002
	// ErrFailedToUploadToStage is an error code denoting the failure to upload a file to the stage
	ErrFailedToUploadToStage = 264003
	// ErrInvalidStageLocation is an error code denoting an invalid stage location
	ErrInvalidStageLocation = 264004
	// ErrLocalPathNotDirectory is an error code denoting a local path that is not a directory
	ErrLocalPathNotDirectory = 264005
	// ErrFileNotExists is an error code denoting the file to be transferred does not exist
	ErrFileNotExists = 264006
	// ErrCompressionNotSupported is an error code denoting the user specified compression type is not supported
	ErrCompressionNotSupported = 264007
	// ErrInternalNotMatchEncryptMaterial is an error code denoting the encryption material specified does not match
	ErrInternalNotMatchEncryptMaterial = 264008
	// ErrCommandNotRecognized is an error code denoting the PUT/GET command was not recognized
	ErrCommandNotRecognized = 264009
	// ErrFailedToConvertToS3Client is an error code denoting the failure of an interface to s3.Client conversion
	ErrFailedToConvertToS3Client = 264010
	// ErrNotImplemented is an error code denoting the file transfer feature is not implemented
	ErrNotImplemented = 264011
	// ErrInvalidPadding is an error code denoting the invalid padding of decryption key
	ErrInvalidPadding = 264012

	/* binding */

	// ErrBindSerialization is an error code for a failed serialization of bind variables
	ErrBindSerialization = 265001
	// ErrBindUpload is an error code for the uploading process of bind elements to the stage
	ErrBindUpload = 265002

	/* async */

	// ErrAsync is an error code for an unknown async error
	ErrAsync = 266001

	/* multi-statement */

	// ErrNoResultIDs is an error code for empty result IDs for multi statement queries
	ErrNoResultIDs = 267001

	/* converter */

	// ErrInvalidTimestampTz is an error code for the case where a returned TIMESTAMP_TZ internal value is invalid
	ErrInvalidTimestampTz = 268000
	// ErrInvalidOffsetStr is an error code for the case where an offset string is invalid. The input string must
	// consist of sHHMI where one sign character '+'/'-' followed by zero filled hours and minutes
	ErrInvalidOffsetStr = 268001
	// ErrInvalidBinaryHexForm is an error code for the case where a binary data in hex form is invalid.
	ErrInvalidBinaryHexForm = 268002
	// ErrTooHighTimestampPrecision is an error code for the case where cannot convert Snowflake timestamp to arrow.Timestamp
	ErrTooHighTimestampPrecision = 268003
	// ErrNullValueInArray is an error code for the case where there are null values in an array without arrayValuesNullable set to true
	ErrNullValueInArray = 268004
	// ErrNullValueInMap is an error code for the case where there are null values in a map without mapValuesNullable set to true
	ErrNullValueInMap = 268005

	/* OCSP */

	// ErrOCSPStatusRevoked is an error code for the case where the certificate is revoked.
	ErrOCSPStatusRevoked = 269001
	// ErrOCSPStatusUnknown is an error code for the case where the certificate revocation status is unknown.
	ErrOCSPStatusUnknown = 269002
	// ErrOCSPInvalidValidity is an error code for the case where the OCSP response validity is invalid.
	ErrOCSPInvalidValidity = 269003
	// ErrOCSPNoOCSPResponderURL is an error code for the case where the OCSP responder URL is not attached.
	ErrOCSPNoOCSPResponderURL = 269004

	/* query Status*/

	// ErrQueryStatus when check the status of a query, receive error or no status
	ErrQueryStatus = 279001
	// ErrQueryIDFormat the query ID given to fetch its result is not valid
	ErrQueryIDFormat = 279101
	// ErrQueryReportedError server side reports the query failed with error
	ErrQueryReportedError = 279201
	// ErrQueryIsRunning the query is still running
	ErrQueryIsRunning = 279301

	/* GS error code */

	// ErrSessionGone is an GS error code for the case that session is already closed
	ErrSessionGone = 390111
	// ErrRoleNotExist is a GS error code for the case that the role specified does not exist
	ErrRoleNotExist = 390189
	// ErrObjectNotExistOrAuthorized is a GS error code for the case that the server-side object specified does not exist
	ErrObjectNotExistOrAuthorized = 390201
)

const (
	errMsgFailedToParseHost                  = "failed to parse a host name. host: %v"
	errMsgFailedToParsePort                  = "failed to parse a port number. port: %v"
	errMsgFailedToParseAuthenticator         = "failed to parse an authenticator: %v"
	errMsgInvalidOffsetStr                   = "offset must be a string consist of sHHMI where one sign character '+'/'-' followed by zero filled hours and minutes: %v"
	errMsgInvalidByteArray                   = "invalid byte array: %v"
	errMsgIdpConnectionError                 = "failed to verify URLs. authenticator: %v, token URL:%v, SSO URL:%v"
	errMsgSSOURLNotMatch                     = "SSO URL didn't match. expected: %v, got: %v"
	errMsgFailedToGetChunk                   = "failed to get a chunk of result sets. idx: %v"
	errMsgFailedToPostQuery                  = "failed to POST. HTTP: %v, URL: %v"
	errMsgFailedToRenew                      = "failed to renew session. HTTP: %v, URL: %v"
	errMsgFailedToCancelQuery                = "failed to cancel query. HTTP: %v, URL: %v"
	errMsgFailedToCloseSession               = "failed to close session. HTTP: %v, URL: %v"
	errMsgFailedToAuth                       = "failed to auth for unknown reason. HTTP: %v, URL: %v"
	errMsgFailedToAuthSAML                   = "failed to auth via SAML for unknown reason. HTTP: %v, URL: %v"
	errMsgFailedToAuthOKTA                   = "failed to auth via OKTA for unknown reason. HTTP: %v, URL: %v"
	errMsgFailedToGetSSO                     = "failed to auth via OKTA for unknown reason. HTTP: %v, URL: %v"
	errMsgFailedToParseResponse              = "failed to parse a response from Snowflake. Response: %v"
	errMsgFailedToGetExternalBrowserResponse = "failed to get an external browser response from Snowflake, err: %s"
	errMsgNoReadOnlyTransaction              = "no readonly mode is supported"
	errMsgNoDefaultTransactionIsolationLevel = "no default isolation transaction level is supported"
	errMsgServiceUnavailable                 = "service is unavailable. check your connectivity. you may need a proxy server. HTTP: %v, URL: %v"
	errMsgFailedToConnect                    = "failed to connect to db. verify account name is correct. HTTP: %v, URL: %v"
	errMsgOCSPStatusRevoked                  = "OCSP revoked: reason:%v, at:%v"
	errMsgOCSPStatusUnknown                  = "OCSP unknown"
	errMsgOCSPInvalidValidity                = "invalid validity: producedAt: %v, thisUpdate: %v, nextUpdate: %v"
	errMsgOCSPNoOCSPResponderURL             = "no OCSP server is attached to the certificate. %v"
	errMsgBindColumnMismatch                 = "column %v has a different number of binds (%v) than column 1 (%v)"
	errMsgNotImplemented                     = "not implemented"
	errMsgFeatureNotSupported                = "feature is not supported: %v"
	errMsgCommandNotRecognized               = "%v command not recognized"
	errMsgLocalPathNotDirectory              = "the local path is not a directory: %v"
	errMsgFileNotExists                      = "file does not exist: %v"
	errMsgInvalidStageFs                     = "destination location type is not valid: %v"
	errMsgInternalNotMatchEncryptMaterial    = "number of downloading files doesn't match the encryption materials. files=%v, encmat=%v"
	errMsgFailedToConvertToS3Client          = "failed to convert interface to s3 client"
	errMsgNoResultIDs                        = "no result IDs returned with the multi-statement query"
	errMsgQueryStatus                        = "server ErrorCode=%s, ErrorMessage=%s"
	errMsgInvalidPadding                     = "invalid padding on input"
	errMsgClientConfigFailed                 = "client configuration failed: %v"
	errMsgNullValueInArray                   = "for handling null values in arrays use WithArrayValuesNullable(ctx)"
	errMsgNullValueInMap                     = "for handling null values in maps use WithMapValuesNullable(ctx)"
	errMsgFailedToParseTomlFile              = "failed to parse toml file. the params %v occurred error with value %v"
	errMsgFailedToFindDSNInTomlFile          = "failed to find DSN in toml file."
	errMsgInvalidPermissionToTomlFile        = "file permissions different than read/write for user. Your Permission: %v"
)

// Returned if a DNS doesn't include account parameter.
func errEmptyAccount() *SnowflakeError {
	return &SnowflakeError{
		Number:  ErrCodeEmptyAccountCode,
		Message: "account is empty",
	}
}

// Returned if a DNS doesn't include user parameter.
func errEmptyUsername() *SnowflakeError {
	return &SnowflakeError{
		Number:  ErrCodeEmptyUsernameCode,
		Message: "user is empty",
	}
}

// Returned if a DNS doesn't include password parameter.
func errEmptyPassword() *SnowflakeError {
	return &SnowflakeError{
		Number:  ErrCodeEmptyPasswordCode,
		Message: "password is empty",
	}
}

// Returned if a DSN's implicit region from account parameter and explicit region parameter conflict.
func errRegionConflict() *SnowflakeError {
	return &SnowflakeError{
		Number:  ErrCodeRegionOverlap,
		Message: "two regions specified",
	}
}

// Returned if a DSN includes an invalid authenticator.
func errFailedToParseAuthenticator() *SnowflakeError {
	return &SnowflakeError{
		Number:  ErrCodeFailedToParseAuthenticator,
		Message: "failed to parse an authenticator",
	}
}

// Returned if the server side returns an error without meaningful message.
func errUnknownError() *SnowflakeError {
	return &SnowflakeError{
		Number:   -1,
		SQLState: "-1",
		Message:  "an unknown server side error occurred",
		QueryID:  "-1",
	}
}

func errNullValueInArray() *SnowflakeError {
	return &SnowflakeError{
		Number:  ErrNullValueInArray,
		Message: errMsgNullValueInArray,
	}
}

func errNullValueInMap() *SnowflakeError {
	return &SnowflakeError{
		Number:  ErrNullValueInMap,
		Message: errMsgNullValueInMap,
	}
}
