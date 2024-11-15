package models
type FileHashType int

const (
    UNKNOWN_FILEHASHTYPE FileHashType = iota
    SHA1_FILEHASHTYPE
    SHA256_FILEHASHTYPE
    MD5_FILEHASHTYPE
    AUTHENTICODEHASH256_FILEHASHTYPE
    LSHASH_FILEHASHTYPE
    CTPH_FILEHASHTYPE
    UNKNOWNFUTUREVALUE_FILEHASHTYPE
)

func (i FileHashType) String() string {
    return []string{"unknown", "sha1", "sha256", "md5", "authenticodeHash256", "lsHash", "ctph", "unknownFutureValue"}[i]
}
func ParseFileHashType(v string) (any, error) {
    result := UNKNOWN_FILEHASHTYPE
    switch v {
        case "unknown":
            result = UNKNOWN_FILEHASHTYPE
        case "sha1":
            result = SHA1_FILEHASHTYPE
        case "sha256":
            result = SHA256_FILEHASHTYPE
        case "md5":
            result = MD5_FILEHASHTYPE
        case "authenticodeHash256":
            result = AUTHENTICODEHASH256_FILEHASHTYPE
        case "lsHash":
            result = LSHASH_FILEHASHTYPE
        case "ctph":
            result = CTPH_FILEHASHTYPE
        case "unknownFutureValue":
            result = UNKNOWNFUTUREVALUE_FILEHASHTYPE
        default:
            return nil, nil
    }
    return &result, nil
}
func SerializeFileHashType(values []FileHashType) []string {
    result := make([]string, len(values))
    for i, v := range values {
        result[i] = v.String()
    }
    return result
}
func (i FileHashType) isMultiValue() bool {
    return false
}
