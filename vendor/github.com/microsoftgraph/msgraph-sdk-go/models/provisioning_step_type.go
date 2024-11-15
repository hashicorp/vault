package models
type ProvisioningStepType int

const (
    IMPORT_PROVISIONINGSTEPTYPE ProvisioningStepType = iota
    SCOPING_PROVISIONINGSTEPTYPE
    MATCHING_PROVISIONINGSTEPTYPE
    PROCESSING_PROVISIONINGSTEPTYPE
    REFERENCERESOLUTION_PROVISIONINGSTEPTYPE
    EXPORT_PROVISIONINGSTEPTYPE
    UNKNOWNFUTUREVALUE_PROVISIONINGSTEPTYPE
)

func (i ProvisioningStepType) String() string {
    return []string{"import", "scoping", "matching", "processing", "referenceResolution", "export", "unknownFutureValue"}[i]
}
func ParseProvisioningStepType(v string) (any, error) {
    result := IMPORT_PROVISIONINGSTEPTYPE
    switch v {
        case "import":
            result = IMPORT_PROVISIONINGSTEPTYPE
        case "scoping":
            result = SCOPING_PROVISIONINGSTEPTYPE
        case "matching":
            result = MATCHING_PROVISIONINGSTEPTYPE
        case "processing":
            result = PROCESSING_PROVISIONINGSTEPTYPE
        case "referenceResolution":
            result = REFERENCERESOLUTION_PROVISIONINGSTEPTYPE
        case "export":
            result = EXPORT_PROVISIONINGSTEPTYPE
        case "unknownFutureValue":
            result = UNKNOWNFUTUREVALUE_PROVISIONINGSTEPTYPE
        default:
            return nil, nil
    }
    return &result, nil
}
func SerializeProvisioningStepType(values []ProvisioningStepType) []string {
    result := make([]string, len(values))
    for i, v := range values {
        result[i] = v.String()
    }
    return result
}
func (i ProvisioningStepType) isMultiValue() bool {
    return false
}
