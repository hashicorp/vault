// Copyright (c) 2017-2021 Snowflake Computing Inc. All rights reserved.

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
			driverTypeKey:    "Go",
			driverVersionKey: SnowflakeGoDriverVersion,
			stacktraceKey:    maskSecrets(string(debug.Stack())),
		},
		Timestamp: time.Now().UnixNano(),
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
	if sc != nil {
		return sc.telemetry.addLog(data)
	}
	return nil // TODO oob telemetry
}

func (se *SnowflakeError) exceptionTelemetry(sc *snowflakeConn) *SnowflakeError {
	data := se.generateTelemetryExceptionData()
	if err := se.sendExceptionTelemetry(sc, data); err != nil {
		logger.Debugf("failed to log to telemetry: %v", data)
	}
	return se
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

	/* binding */

	// ErrBindSerialization is an error code for a failed serialization of bind variables
	ErrBindSerialization = 265001
	// ErrBindUpload is an error code for the uploading process of bind elements to the stage
	ErrBindUpload = 265002

	/* converter */

	// ErrInvalidTimestampTz is an error code for the case where a returned TIMESTAMP_TZ internal value is invalid
	ErrInvalidTimestampTz = 268000
	// ErrInvalidOffsetStr is an error code for the case where a offset string is invalid. The input string must
	// consist of sHHMI where one sign character '+'/'-' followed by zero filled hours and minutes
	ErrInvalidOffsetStr = 268001
	// ErrInvalidBinaryHexForm is an error code for the case where a binary data in hex form is invalid.
	ErrInvalidBinaryHexForm = 268002

	/* OCSP */

	// ErrOCSPStatusRevoked is an error code for the case where the certificate is revoked.
	ErrOCSPStatusRevoked = 269001
	// ErrOCSPStatusUnknown is an error code for the case where the certificate revocation status is unknown.
	ErrOCSPStatusUnknown = 269002
	// ErrOCSPInvalidValidity is an error code for the case where the OCSP response validity is invalid.
	ErrOCSPInvalidValidity = 269003
	// ErrOCSPNoOCSPResponderURL is an error code for the case where the OCSP responder URL is not attached.
	ErrOCSPNoOCSPResponderURL = 269004

	/* Query Status*/

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
)

var (
	// ErrEmptyAccount is returned if a DNS doesn't include account parameter.
	ErrEmptyAccount = &SnowflakeError{
		Number:  ErrCodeEmptyAccountCode,
		Message: "account is empty",
	}
	// ErrEmptyUsername is returned if a DNS doesn't include user parameter.
	ErrEmptyUsername = &SnowflakeError{
		Number:  ErrCodeEmptyUsernameCode,
		Message: "user is empty",
	}
	// ErrEmptyPassword is returned if a DNS doesn't include password parameter.
	ErrEmptyPassword = &SnowflakeError{
		Number:  ErrCodeEmptyPasswordCode,
		Message: "password is empty"}

	// ErrInvalidRegion is returned if a DSN's implicit region from account parameter and explicit region parameter conflict.
	ErrInvalidRegion = &SnowflakeError{
		Number:  ErrCodeRegionOverlap,
		Message: "two regions specified"}
)
