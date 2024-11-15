package models
type RegistryValueType int

const (
    UNKNOWN_REGISTRYVALUETYPE RegistryValueType = iota
    BINARY_REGISTRYVALUETYPE
    DWORD_REGISTRYVALUETYPE
    DWORDLITTLEENDIAN_REGISTRYVALUETYPE
    DWORDBIGENDIAN_REGISTRYVALUETYPE
    EXPANDSZ_REGISTRYVALUETYPE
    LINK_REGISTRYVALUETYPE
    MULTISZ_REGISTRYVALUETYPE
    NONE_REGISTRYVALUETYPE
    QWORD_REGISTRYVALUETYPE
    QWORDLITTLEENDIAN_REGISTRYVALUETYPE
    SZ_REGISTRYVALUETYPE
    UNKNOWNFUTUREVALUE_REGISTRYVALUETYPE
)

func (i RegistryValueType) String() string {
    return []string{"unknown", "binary", "dword", "dwordLittleEndian", "dwordBigEndian", "expandSz", "link", "multiSz", "none", "qword", "qwordlittleEndian", "sz", "unknownFutureValue"}[i]
}
func ParseRegistryValueType(v string) (any, error) {
    result := UNKNOWN_REGISTRYVALUETYPE
    switch v {
        case "unknown":
            result = UNKNOWN_REGISTRYVALUETYPE
        case "binary":
            result = BINARY_REGISTRYVALUETYPE
        case "dword":
            result = DWORD_REGISTRYVALUETYPE
        case "dwordLittleEndian":
            result = DWORDLITTLEENDIAN_REGISTRYVALUETYPE
        case "dwordBigEndian":
            result = DWORDBIGENDIAN_REGISTRYVALUETYPE
        case "expandSz":
            result = EXPANDSZ_REGISTRYVALUETYPE
        case "link":
            result = LINK_REGISTRYVALUETYPE
        case "multiSz":
            result = MULTISZ_REGISTRYVALUETYPE
        case "none":
            result = NONE_REGISTRYVALUETYPE
        case "qword":
            result = QWORD_REGISTRYVALUETYPE
        case "qwordlittleEndian":
            result = QWORDLITTLEENDIAN_REGISTRYVALUETYPE
        case "sz":
            result = SZ_REGISTRYVALUETYPE
        case "unknownFutureValue":
            result = UNKNOWNFUTUREVALUE_REGISTRYVALUETYPE
        default:
            return nil, nil
    }
    return &result, nil
}
func SerializeRegistryValueType(values []RegistryValueType) []string {
    result := make([]string, len(values))
    for i, v := range values {
        result[i] = v.String()
    }
    return result
}
func (i RegistryValueType) isMultiValue() bool {
    return false
}
