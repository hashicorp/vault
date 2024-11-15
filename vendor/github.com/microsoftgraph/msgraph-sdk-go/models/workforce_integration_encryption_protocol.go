package models
type WorkforceIntegrationEncryptionProtocol int

const (
    SHAREDSECRET_WORKFORCEINTEGRATIONENCRYPTIONPROTOCOL WorkforceIntegrationEncryptionProtocol = iota
    UNKNOWNFUTUREVALUE_WORKFORCEINTEGRATIONENCRYPTIONPROTOCOL
)

func (i WorkforceIntegrationEncryptionProtocol) String() string {
    return []string{"sharedSecret", "unknownFutureValue"}[i]
}
func ParseWorkforceIntegrationEncryptionProtocol(v string) (any, error) {
    result := SHAREDSECRET_WORKFORCEINTEGRATIONENCRYPTIONPROTOCOL
    switch v {
        case "sharedSecret":
            result = SHAREDSECRET_WORKFORCEINTEGRATIONENCRYPTIONPROTOCOL
        case "unknownFutureValue":
            result = UNKNOWNFUTUREVALUE_WORKFORCEINTEGRATIONENCRYPTIONPROTOCOL
        default:
            return nil, nil
    }
    return &result, nil
}
func SerializeWorkforceIntegrationEncryptionProtocol(values []WorkforceIntegrationEncryptionProtocol) []string {
    result := make([]string, len(values))
    for i, v := range values {
        result[i] = v.String()
    }
    return result
}
func (i WorkforceIntegrationEncryptionProtocol) isMultiValue() bool {
    return false
}
