package security
type ExportFileStructure int

const (
    NONE_EXPORTFILESTRUCTURE ExportFileStructure = iota
    DIRECTORY_EXPORTFILESTRUCTURE
    PST_EXPORTFILESTRUCTURE
    UNKNOWNFUTUREVALUE_EXPORTFILESTRUCTURE
)

func (i ExportFileStructure) String() string {
    return []string{"none", "directory", "pst", "unknownFutureValue"}[i]
}
func ParseExportFileStructure(v string) (any, error) {
    result := NONE_EXPORTFILESTRUCTURE
    switch v {
        case "none":
            result = NONE_EXPORTFILESTRUCTURE
        case "directory":
            result = DIRECTORY_EXPORTFILESTRUCTURE
        case "pst":
            result = PST_EXPORTFILESTRUCTURE
        case "unknownFutureValue":
            result = UNKNOWNFUTUREVALUE_EXPORTFILESTRUCTURE
        default:
            return nil, nil
    }
    return &result, nil
}
func SerializeExportFileStructure(values []ExportFileStructure) []string {
    result := make([]string, len(values))
    for i, v := range values {
        result[i] = v.String()
    }
    return result
}
func (i ExportFileStructure) isMultiValue() bool {
    return false
}
