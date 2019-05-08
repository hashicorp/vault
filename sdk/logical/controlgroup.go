package logical

import (
	"time"
)

type ControlGroup struct {
	Authorizations []*Authz  `json:"authorizations"`
	RequestTime    time.Time `json:"request_time"`
	Approved       bool      `json:"approved"`
	NamespaceID    string    `json:"namespace_id"`
}

type Authz struct {
	Token             string    `json:"token"`
	AuthorizationTime time.Time `json:"authorization_time"`
}
