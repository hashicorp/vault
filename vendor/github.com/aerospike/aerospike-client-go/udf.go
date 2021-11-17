package aerospike

// UDF carries information about UDFs on the server
type UDF struct {
	// Filename of the UDF
	Filename string
	// Hash digest of the UDF
	Hash string
	// Language of UDF
	Language Language
}
