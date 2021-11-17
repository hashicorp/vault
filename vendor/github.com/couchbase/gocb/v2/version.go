package gocb

// Version returns a string representation of the current SDK version.
func Version() string {
	return goCbVersionStr
}

// Identifier returns a string representation of the current SDK identifier.
func Identifier() string {
	return "gocb/" + goCbVersionStr
}
