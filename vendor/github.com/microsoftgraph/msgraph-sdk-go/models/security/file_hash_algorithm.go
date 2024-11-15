package security
type FileHashAlgorithm int

const (
    UNKNOWN_FILEHASHALGORITHM FileHashAlgorithm = iota
    MD5_FILEHASHALGORITHM
    SHA1_FILEHASHALGORITHM
    SHA256_FILEHASHALGORITHM
    SHA256AC_FILEHASHALGORITHM
    UNKNOWNFUTUREVALUE_FILEHASHALGORITHM
)

func (i FileHashAlgorithm) String() string {
    return []string{"unknown", "md5", "sha1", "sha256", "sha256ac", "unknownFutureValue"}[i]
}
func ParseFileHashAlgorithm(v string) (any, error) {
    result := UNKNOWN_FILEHASHALGORITHM
    switch v {
        case "unknown":
            result = UNKNOWN_FILEHASHALGORITHM
        case "md5":
            result = MD5_FILEHASHALGORITHM
        case "sha1":
            result = SHA1_FILEHASHALGORITHM
        case "sha256":
            result = SHA256_FILEHASHALGORITHM
        case "sha256ac":
            result = SHA256AC_FILEHASHALGORITHM
        case "unknownFutureValue":
            result = UNKNOWNFUTUREVALUE_FILEHASHALGORITHM
        default:
            return nil, nil
    }
    return &result, nil
}
func SerializeFileHashAlgorithm(values []FileHashAlgorithm) []string {
    result := make([]string, len(values))
    for i, v := range values {
        result[i] = v.String()
    }
    return result
}
func (i FileHashAlgorithm) isMultiValue() bool {
    return false
}
