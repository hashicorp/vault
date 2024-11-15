package models
// BitLockerEncryptionMethod types
type BitLockerEncryptionMethod int

const (
    // AES-CBC 128-bit.
    AESCBC128_BITLOCKERENCRYPTIONMETHOD BitLockerEncryptionMethod = iota
    // AES-CBC 256-bit.
    AESCBC256_BITLOCKERENCRYPTIONMETHOD
    // XTS-AES 128-bit.
    XTSAES128_BITLOCKERENCRYPTIONMETHOD
    // XTS-AES 256-bit.
    XTSAES256_BITLOCKERENCRYPTIONMETHOD
)

func (i BitLockerEncryptionMethod) String() string {
    return []string{"aesCbc128", "aesCbc256", "xtsAes128", "xtsAes256"}[i]
}
func ParseBitLockerEncryptionMethod(v string) (any, error) {
    result := AESCBC128_BITLOCKERENCRYPTIONMETHOD
    switch v {
        case "aesCbc128":
            result = AESCBC128_BITLOCKERENCRYPTIONMETHOD
        case "aesCbc256":
            result = AESCBC256_BITLOCKERENCRYPTIONMETHOD
        case "xtsAes128":
            result = XTSAES128_BITLOCKERENCRYPTIONMETHOD
        case "xtsAes256":
            result = XTSAES256_BITLOCKERENCRYPTIONMETHOD
        default:
            return nil, nil
    }
    return &result, nil
}
func SerializeBitLockerEncryptionMethod(values []BitLockerEncryptionMethod) []string {
    result := make([]string, len(values))
    for i, v := range values {
        result[i] = v.String()
    }
    return result
}
func (i BitLockerEncryptionMethod) isMultiValue() bool {
    return false
}
