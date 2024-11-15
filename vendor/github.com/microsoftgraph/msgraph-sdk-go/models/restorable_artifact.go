package models
type RestorableArtifact int

const (
    MESSAGE_RESTORABLEARTIFACT RestorableArtifact = iota
    UNKNOWNFUTUREVALUE_RESTORABLEARTIFACT
)

func (i RestorableArtifact) String() string {
    return []string{"message", "unknownFutureValue"}[i]
}
func ParseRestorableArtifact(v string) (any, error) {
    result := MESSAGE_RESTORABLEARTIFACT
    switch v {
        case "message":
            result = MESSAGE_RESTORABLEARTIFACT
        case "unknownFutureValue":
            result = UNKNOWNFUTUREVALUE_RESTORABLEARTIFACT
        default:
            return nil, nil
    }
    return &result, nil
}
func SerializeRestorableArtifact(values []RestorableArtifact) []string {
    result := make([]string, len(values))
    for i, v := range values {
        result[i] = v.String()
    }
    return result
}
func (i RestorableArtifact) isMultiValue() bool {
    return false
}
