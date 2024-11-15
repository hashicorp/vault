package models
type ObjectMappingMetadata int

const (
    ESCROWBEHAVIOR_OBJECTMAPPINGMETADATA ObjectMappingMetadata = iota
    DISABLEMONITORINGFORCHANGES_OBJECTMAPPINGMETADATA
    ORIGINALJOININGPROPERTY_OBJECTMAPPINGMETADATA
    DISPOSITION_OBJECTMAPPINGMETADATA
    ISCUSTOMERDEFINED_OBJECTMAPPINGMETADATA
    EXCLUDEFROMREPORTING_OBJECTMAPPINGMETADATA
    UNSYNCHRONIZED_OBJECTMAPPINGMETADATA
)

func (i ObjectMappingMetadata) String() string {
    return []string{"EscrowBehavior", "DisableMonitoringForChanges", "OriginalJoiningProperty", "Disposition", "IsCustomerDefined", "ExcludeFromReporting", "Unsynchronized"}[i]
}
func ParseObjectMappingMetadata(v string) (any, error) {
    result := ESCROWBEHAVIOR_OBJECTMAPPINGMETADATA
    switch v {
        case "EscrowBehavior":
            result = ESCROWBEHAVIOR_OBJECTMAPPINGMETADATA
        case "DisableMonitoringForChanges":
            result = DISABLEMONITORINGFORCHANGES_OBJECTMAPPINGMETADATA
        case "OriginalJoiningProperty":
            result = ORIGINALJOININGPROPERTY_OBJECTMAPPINGMETADATA
        case "Disposition":
            result = DISPOSITION_OBJECTMAPPINGMETADATA
        case "IsCustomerDefined":
            result = ISCUSTOMERDEFINED_OBJECTMAPPINGMETADATA
        case "ExcludeFromReporting":
            result = EXCLUDEFROMREPORTING_OBJECTMAPPINGMETADATA
        case "Unsynchronized":
            result = UNSYNCHRONIZED_OBJECTMAPPINGMETADATA
        default:
            return nil, nil
    }
    return &result, nil
}
func SerializeObjectMappingMetadata(values []ObjectMappingMetadata) []string {
    result := make([]string, len(values))
    for i, v := range values {
        result[i] = v.String()
    }
    return result
}
func (i ObjectMappingMetadata) isMultiValue() bool {
    return false
}
