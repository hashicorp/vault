// Copyright (c) 2016, 2018, 2019, Oracle and/or its affiliates. All rights reserved.
// Code generated. DO NOT EDIT.

// Audit API
//
// API for the Audit Service. You can use this API for queries, but not bulk-export operations.
//

package audit

import (
	"github.com/oracle/oci-go-sdk/common"
)

// AuditEvent The representation of AuditEvent
type AuditEvent struct {

	// The OCID of the tenant.
	TenantId *string `mandatory:"false" json:"tenantId"`

	// The OCID of the compartment.
	CompartmentId *string `mandatory:"false" json:"compartmentId"`

	// The name of the compartment. This value is the friendly name associated with compartmentId.
	// This value can change, but the service logs the value that appeared at the time of the audit event.
	CompartmentName *string `mandatory:"false" json:"compartmentName"`

	// The GUID of the event.
	EventId *string `mandatory:"false" json:"eventId"`

	// The name of the event.
	// Example: `LaunchInstance`
	EventName *string `mandatory:"false" json:"eventName"`

	// The source of the event.
	EventSource *string `mandatory:"false" json:"eventSource"`

	// The type of the event.
	EventType *string `mandatory:"false" json:"eventType"`

	// The time the event occurred, expressed in RFC 3339 (https://tools.ietf.org/html/rfc3339) timestamp format.
	EventTime *common.SDKTime `mandatory:"false" json:"eventTime"`

	// The OCID of the user whose action triggered the event.
	PrincipalId *string `mandatory:"false" json:"principalId"`

	// The credential ID of the user. This value is extracted from the HTTP 'Authorization' request header. It consists of the tenantId, userId, and user fingerprint, all delimited by a slash (/).
	CredentialId *string `mandatory:"false" json:"credentialId"`

	// The HTTP method of the request.
	RequestAction *string `mandatory:"false" json:"requestAction"`

	// The opc-request-id of the request.
	RequestId *string `mandatory:"false" json:"requestId"`

	// The user agent of the client that made the request.
	RequestAgent *string `mandatory:"false" json:"requestAgent"`

	// The HTTP header fields and values in the request.
	RequestHeaders map[string][]string `mandatory:"false" json:"requestHeaders"`

	// The IP address of the source of the request.
	RequestOrigin *string `mandatory:"false" json:"requestOrigin"`

	// The query parameter fields and values for the request.
	RequestParameters map[string][]string `mandatory:"false" json:"requestParameters"`

	// The resource targeted by the request.
	RequestResource *string `mandatory:"false" json:"requestResource"`

	// The headers of the response.
	ResponseHeaders map[string][]string `mandatory:"false" json:"responseHeaders"`

	// The status code of the response.
	ResponseStatus *string `mandatory:"false" json:"responseStatus"`

	// The time of the response to the audited request, expressed in RFC 3339 (https://tools.ietf.org/html/rfc3339) timestamp format.
	ResponseTime *common.SDKTime `mandatory:"false" json:"responseTime"`

	// Metadata of interest from the response payload. For example, the OCID of a resource.
	ResponsePayload map[string]interface{} `mandatory:"false" json:"responsePayload"`

	// The name of the user or service. This value is the friendly name associated with principalId.
	UserName *string `mandatory:"false" json:"userName"`
}

func (m AuditEvent) String() string {
	return common.PointerString(m)
}
