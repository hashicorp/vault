package radius

// NonAuthenticResponseError is returned when a client was expecting
// a valid response but did not receive one.
type NonAuthenticResponseError struct {
}

func (e *NonAuthenticResponseError) Error() string {
	return `radius: non-authentic response`
}
