package models
type AttributeDefinitionMetadata int

const (
    BASEATTRIBUTENAME_ATTRIBUTEDEFINITIONMETADATA AttributeDefinitionMetadata = iota
    COMPLEXOBJECTDEFINITION_ATTRIBUTEDEFINITIONMETADATA
    ISCONTAINER_ATTRIBUTEDEFINITIONMETADATA
    ISCUSTOMERDEFINED_ATTRIBUTEDEFINITIONMETADATA
    ISDOMAINQUALIFIED_ATTRIBUTEDEFINITIONMETADATA
    LINKPROPERTYNAMES_ATTRIBUTEDEFINITIONMETADATA
    LINKTYPENAME_ATTRIBUTEDEFINITIONMETADATA
    MAXIMUMLENGTH_ATTRIBUTEDEFINITIONMETADATA
    REFERENCEDPROPERTY_ATTRIBUTEDEFINITIONMETADATA
)

func (i AttributeDefinitionMetadata) String() string {
    return []string{"BaseAttributeName", "ComplexObjectDefinition", "IsContainer", "IsCustomerDefined", "IsDomainQualified", "LinkPropertyNames", "LinkTypeName", "MaximumLength", "ReferencedProperty"}[i]
}
func ParseAttributeDefinitionMetadata(v string) (any, error) {
    result := BASEATTRIBUTENAME_ATTRIBUTEDEFINITIONMETADATA
    switch v {
        case "BaseAttributeName":
            result = BASEATTRIBUTENAME_ATTRIBUTEDEFINITIONMETADATA
        case "ComplexObjectDefinition":
            result = COMPLEXOBJECTDEFINITION_ATTRIBUTEDEFINITIONMETADATA
        case "IsContainer":
            result = ISCONTAINER_ATTRIBUTEDEFINITIONMETADATA
        case "IsCustomerDefined":
            result = ISCUSTOMERDEFINED_ATTRIBUTEDEFINITIONMETADATA
        case "IsDomainQualified":
            result = ISDOMAINQUALIFIED_ATTRIBUTEDEFINITIONMETADATA
        case "LinkPropertyNames":
            result = LINKPROPERTYNAMES_ATTRIBUTEDEFINITIONMETADATA
        case "LinkTypeName":
            result = LINKTYPENAME_ATTRIBUTEDEFINITIONMETADATA
        case "MaximumLength":
            result = MAXIMUMLENGTH_ATTRIBUTEDEFINITIONMETADATA
        case "ReferencedProperty":
            result = REFERENCEDPROPERTY_ATTRIBUTEDEFINITIONMETADATA
        default:
            return nil, nil
    }
    return &result, nil
}
func SerializeAttributeDefinitionMetadata(values []AttributeDefinitionMetadata) []string {
    result := make([]string, len(values))
    for i, v := range values {
        result[i] = v.String()
    }
    return result
}
func (i AttributeDefinitionMetadata) isMultiValue() bool {
    return false
}
