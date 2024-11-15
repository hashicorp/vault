package security
type ExportFormat int

const (
    PST_EXPORTFORMAT ExportFormat = iota
    MSG_EXPORTFORMAT
    EML_EXPORTFORMAT
    UNKNOWNFUTUREVALUE_EXPORTFORMAT
)

func (i ExportFormat) String() string {
    return []string{"pst", "msg", "eml", "unknownFutureValue"}[i]
}
func ParseExportFormat(v string) (any, error) {
    result := PST_EXPORTFORMAT
    switch v {
        case "pst":
            result = PST_EXPORTFORMAT
        case "msg":
            result = MSG_EXPORTFORMAT
        case "eml":
            result = EML_EXPORTFORMAT
        case "unknownFutureValue":
            result = UNKNOWNFUTUREVALUE_EXPORTFORMAT
        default:
            return nil, nil
    }
    return &result, nil
}
func SerializeExportFormat(values []ExportFormat) []string {
    result := make([]string, len(values))
    for i, v := range values {
        result[i] = v.String()
    }
    return result
}
func (i ExportFormat) isMultiValue() bool {
    return false
}
