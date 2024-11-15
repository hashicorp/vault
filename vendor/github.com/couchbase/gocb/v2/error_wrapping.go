package gocb

import (
	"errors"

	gocbcore "github.com/couchbase/gocbcore/v10"
)

func maybeEnhanceCoreErr(err error) error {
	if kvErr, ok := err.(*gocbcore.KeyValueError); ok {
		return &KeyValueError{
			InnerError:         kvErr.InnerError,
			StatusCode:         kvErr.StatusCode,
			DocumentID:         kvErr.DocumentKey,
			BucketName:         kvErr.BucketName,
			ScopeName:          kvErr.ScopeName,
			CollectionName:     kvErr.CollectionName,
			CollectionID:       kvErr.CollectionID,
			ErrorName:          kvErr.ErrorName,
			ErrorDescription:   kvErr.ErrorDescription,
			Opaque:             kvErr.Opaque,
			Context:            kvErr.Context,
			Ref:                kvErr.Ref,
			RetryReasons:       translateCoreRetryReasons(kvErr.RetryReasons),
			RetryAttempts:      kvErr.RetryAttempts,
			LastDispatchedTo:   kvErr.LastDispatchedTo,
			LastDispatchedFrom: kvErr.LastDispatchedFrom,
			LastConnectionID:   kvErr.LastConnectionID,
		}
	}
	if viewErr, ok := err.(*gocbcore.ViewError); ok {
		return &ViewError{
			InnerError:         viewErr.InnerError,
			DesignDocumentName: viewErr.DesignDocumentName,
			ViewName:           viewErr.ViewName,
			Errors:             translateCoreViewErrorDesc(viewErr.Errors),
			Endpoint:           viewErr.Endpoint,
			RetryReasons:       translateCoreRetryReasons(viewErr.RetryReasons),
			RetryAttempts:      viewErr.RetryAttempts,
			ErrorText:          viewErr.ErrorText,
			HTTPStatusCode:     viewErr.HTTPResponseCode,
		}
	}
	if queryErr, ok := err.(*gocbcore.N1QLError); ok {
		inner := queryErr.InnerError

		if errors.Is(inner, ErrFeatureNotAvailable) {
			if len(queryErr.Errors) > 0 {
				desc := queryErr.Errors[0]
				// We replace the gocbcore wrapped inner feature not available error with our own to provide gocb
				// specific context for the user.
				if desc.Code == 1197 {
					inner = wrapError(ErrFeatureNotAvailable, "this server requires that scope.Query() is used rather than "+
						"cluster.Query(), if this is a transaction then pass a Scope within TransactionQueryOptions")
				}
			}
		}

		return &QueryError{
			InnerError:      inner,
			Statement:       queryErr.Statement,
			ClientContextID: queryErr.ClientContextID,
			Errors:          translateCoreQueryErrorDesc(queryErr.Errors),
			Endpoint:        queryErr.Endpoint,
			RetryReasons:    translateCoreRetryReasons(queryErr.RetryReasons),
			RetryAttempts:   queryErr.RetryAttempts,
			ErrorText:       queryErr.ErrorText,
			HTTPStatusCode:  queryErr.HTTPResponseCode,
		}
	}
	if analyticsErr, ok := err.(*gocbcore.AnalyticsError); ok {
		return &AnalyticsError{
			InnerError:      analyticsErr.InnerError,
			Statement:       analyticsErr.Statement,
			ClientContextID: analyticsErr.ClientContextID,
			Errors:          translateCoreAnalyticsErrorDesc(analyticsErr.Errors),
			Endpoint:        analyticsErr.Endpoint,
			RetryReasons:    translateCoreRetryReasons(analyticsErr.RetryReasons),
			RetryAttempts:   analyticsErr.RetryAttempts,
			ErrorText:       analyticsErr.ErrorText,
			HTTPStatusCode:  analyticsErr.HTTPResponseCode,
		}
	}
	if searchErr, ok := err.(*gocbcore.SearchError); ok {
		return &SearchError{
			InnerError:     searchErr.InnerError,
			Query:          searchErr.Query,
			Endpoint:       searchErr.Endpoint,
			RetryReasons:   translateCoreRetryReasons(searchErr.RetryReasons),
			RetryAttempts:  searchErr.RetryAttempts,
			ErrorText:      searchErr.ErrorText,
			IndexName:      searchErr.IndexName,
			HTTPStatusCode: searchErr.HTTPResponseCode,
		}
	}
	if httpErr, ok := err.(*gocbcore.HTTPError); ok {
		return &HTTPError{
			InnerError:    httpErr.InnerError,
			UniqueID:      httpErr.UniqueID,
			Endpoint:      httpErr.Endpoint,
			RetryReasons:  translateCoreRetryReasons(httpErr.RetryReasons),
			RetryAttempts: httpErr.RetryAttempts,
		}
	}

	if timeoutErr, ok := err.(*gocbcore.TimeoutError); ok {
		return &TimeoutError{
			InnerError:         timeoutErr.InnerError,
			OperationID:        timeoutErr.OperationID,
			Opaque:             timeoutErr.Opaque,
			TimeObserved:       timeoutErr.TimeObserved,
			RetryReasons:       translateCoreRetryReasons(timeoutErr.RetryReasons),
			RetryAttempts:      timeoutErr.RetryAttempts,
			LastDispatchedTo:   timeoutErr.LastDispatchedTo,
			LastDispatchedFrom: timeoutErr.LastDispatchedFrom,
			LastConnectionID:   timeoutErr.LastConnectionID,
		}
	}
	return err
}

func maybeEnhanceKVErr(err error, bucketName, scopeName, collName, docKey string) error {
	return maybeEnhanceCoreErr(err)
}

func maybeEnhanceCollKVErr(err error, coll *Collection, docKey string) error {
	return maybeEnhanceKVErr(err, coll.bucketName(), coll.Name(), coll.ScopeName(), docKey)
}

func maybeEnhanceViewError(err error) error {
	return maybeEnhanceCoreErr(err)
}

func maybeEnhanceCoreQueryError(err error) error {
	return maybeEnhanceCoreErr(err)
}

func maybeEnhanceAnalyticsError(err error) error {
	return maybeEnhanceCoreErr(err)
}

func maybeEnhanceSearchError(err error) error {
	return maybeEnhanceCoreErr(err)
}
