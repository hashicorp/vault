package mstypes

// LPWSTR implements https://msdn.microsoft.com/en-us/library/cc230355.aspx
type LPWSTR struct {
	Value string `ndr:"pointer,conformant,varying"`
}

func (s *LPWSTR) String() string {
	return s.Value
}
