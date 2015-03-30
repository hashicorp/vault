package logical

// Auth is the resulting authentication information that is part of
// Response for credential backends.
type Auth struct {
	// Policies is the list of policies that the authenticated user
	// is associated with.
	Policies []string

	// Metadata is used to attach arbitrary string-type metadata to
	// an authenticated user. This metadata will be outputted into the
	// audit log.
	Metadata map[string]string
}
