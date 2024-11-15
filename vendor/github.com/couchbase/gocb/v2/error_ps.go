package gocb

import (
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const (
	preconditionLocked              = "LOCKED"
	preconditionPathMismatch        = "PATH_MISMATCH"
	preconditionDocNotJSON          = "DOC_NOT_JSON"
	preconditionDocTooDeep          = "DOC_TOO_DEEP"
	preconditionValueTooLarge       = "VALUE_TOO_LARGE"
	preconditionValueOutOfRange     = "VALUE_OUT_OF_RANGE"
	preconditionPathValueOutOfRange = "PATH_VALUE_OUT_OF_RANGE"
)

const (
	resourceTypeDocument    = "document"
	resourceTypeIndex       = "queryindex"
	resourceTypeBucket      = "bucket"
	resourceTypeScope       = "scope"
	resourceTypeCollection  = "collection"
	resourceTypePath        = "path"
	resourceTypeSearchIndex = "searchindex"
)

const (
	reasonCasMismatch = "CAS_MISMATCH"
)

func mapPsErrorToGocbError(err error, readOnly bool) *GenericError {
	st, ok := status.FromError(err)
	if !ok {
		return makeGenericError(err, nil)
	}

	return mapPsErrorStatusToGocbError(st, readOnly)
}

func mapPsErrorStatusToGocbError(st *status.Status, readOnly bool) *GenericError {
	context := map[string]interface{}{
		"server":  st.Message(),
		"details": len(st.Details()),
	}

	code := st.Code()
	details := st.Details()
	if len(details) > 0 {
		var baseErr error
		detail := details[0]
		switch d := detail.(type) {
		case *errdetails.PreconditionFailure:
			if len(d.Violations) > 0 {
				context["precondition_violation"] = d.Violations[0].Type

				violation := d.Violations[0]
				switch violation.Type {
				case preconditionLocked:
					baseErr = ErrDocumentLocked
				case preconditionPathMismatch:
					baseErr = ErrPathMismatch
				case preconditionDocNotJSON:
					baseErr = ErrDocumentNotJSON
				case preconditionDocTooDeep:
					baseErr = ErrDocumentTooDeep
				case preconditionValueTooLarge:
					baseErr = ErrValueTooLarge
				case preconditionValueOutOfRange:
					baseErr = ErrValueInvalid
				case preconditionPathValueOutOfRange:
					baseErr = ErrNumberTooBig
				}
			}
		case *errdetails.ResourceInfo:
			if d.ResourceName != "" {
				context["resource_name"] = d.ResourceName
			}
			context["resource_type"] = d.ResourceType

			if code == codes.NotFound {
				switch d.ResourceType {
				case resourceTypeDocument:
					baseErr = ErrDocumentNotFound
				case resourceTypeIndex:
					baseErr = ErrIndexNotFound
				case resourceTypeSearchIndex:
					baseErr = ErrIndexNotFound
				case resourceTypeBucket:
					baseErr = ErrBucketNotFound
				case resourceTypeScope:
					baseErr = ErrScopeNotFound
				case resourceTypeCollection:
					baseErr = ErrCollectionNotFound
				case resourceTypePath:
					baseErr = ErrPathNotFound
				}
			} else if code == codes.AlreadyExists {
				switch d.ResourceType {
				case resourceTypeDocument:
					baseErr = ErrDocumentExists
				case resourceTypeIndex:
					baseErr = ErrIndexExists
				case resourceTypeSearchIndex:
					baseErr = ErrIndexExists
				case resourceTypeBucket:
					baseErr = ErrBucketExists
				case resourceTypeScope:
					baseErr = ErrScopeExists
				case resourceTypeCollection:
					baseErr = ErrCollectionExists
				case resourceTypePath:
					baseErr = ErrPathExists
				}
			}
		case *errdetails.ErrorInfo:
			context["reason"] = d.Reason
			switch d.Reason {
			case reasonCasMismatch:
				baseErr = ErrCasMismatch
			}
		}
		if baseErr != nil {
			return makeGenericError(baseErr, context)
		}
	}

	var baseErr error
	switch st.Code() {
	case codes.Canceled:
		baseErr = ErrRequestCanceled
	case codes.Internal:
		baseErr = ErrInternalServerFailure
	case codes.InvalidArgument:
		baseErr = ErrInvalidArgument
	case codes.DeadlineExceeded:
		if readOnly {
			baseErr = ErrUnambiguousTimeout
		} else {
			baseErr = ErrAmbiguousTimeout
		}
	case codes.Unauthenticated:
		baseErr = wrapError(ErrAuthenticationFailure, "server reported that permission to the resource was denied")
	case codes.PermissionDenied:
		baseErr = wrapError(ErrAuthenticationFailure, "server reported that permission to the resource was denied")
	case codes.Unimplemented:
		baseErr = wrapError(ErrFeatureNotAvailable, st.Message())
	case codes.Unavailable:
		baseErr = ErrServiceNotAvailable
	default:
		baseErr = st.Err()
	}

	return makeGenericError(baseErr, context)
}
