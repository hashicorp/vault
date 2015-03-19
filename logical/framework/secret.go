package framework

// Secret is a type of secret that can be returned from a backend.
type Secret struct {
	// Type is the name of this secret type. This is used to setup the
	// vault ID and to look up the proper secret structure when revocation/
	// renewal happens. Once this is set this should not be changed.
	Type string

	// Fields is the mapping of data fields and schema that comprise
	// the structure of this secret.
	Fields map[string]*FieldSchema
}
