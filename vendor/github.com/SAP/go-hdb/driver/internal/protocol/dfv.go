package protocol

// Data format version values.
const (
	DfvLevel0 int = 0 // base data format
	DfvLevel1 int = 1 // eval types support all data types
	DfvLevel2 int = 2 // reserved, broken, do not use
	DfvLevel3 int = 3 // additional types Longdate, Secondate, Daydate, Secondtime supported for NGAP
	DfvLevel4 int = 4 // generic support for new date/time types
	DfvLevel5 int = 5 // spatial types in ODBC on request
	DfvLevel6 int = 6 // BINTEXT
	DfvLevel7 int = 7 // with boolean support
	DfvLevel8 int = 8 // with FIXED8/12/16 support
)

var (
	defaultDfv    = DfvLevel8
	supportedDfvs = []int{DfvLevel1, DfvLevel4, DfvLevel6, DfvLevel8}
)

// SupportedDfvs returns a slice of data format versions supported by the driver.
// If parameter defaultOnly is set only the default dfv is returned, otherwise
// all supported dfv values are returned.
func SupportedDfvs(defaultOnly bool) []int {
	if defaultOnly {
		return []int{defaultDfv}
	}
	return supportedDfvs
}

// IsSupportedDfv returns true if the data format version dfv is supported by the driver, false otherwise.
func IsSupportedDfv(dfv int) bool {
	return dfv == DfvLevel1 || dfv == DfvLevel4 || dfv == DfvLevel6 || dfv == DfvLevel8
}
