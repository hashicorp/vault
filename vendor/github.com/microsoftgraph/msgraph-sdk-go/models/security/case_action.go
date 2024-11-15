package security
type CaseAction int

const (
    CONTENTEXPORT_CASEACTION CaseAction = iota
    APPLYTAGS_CASEACTION
    CONVERTTOPDF_CASEACTION
    INDEX_CASEACTION
    ESTIMATESTATISTICS_CASEACTION
    ADDTOREVIEWSET_CASEACTION
    HOLDUPDATE_CASEACTION
    UNKNOWNFUTUREVALUE_CASEACTION
    PURGEDATA_CASEACTION
    EXPORTREPORT_CASEACTION
    EXPORTRESULT_CASEACTION
)

func (i CaseAction) String() string {
    return []string{"contentExport", "applyTags", "convertToPdf", "index", "estimateStatistics", "addToReviewSet", "holdUpdate", "unknownFutureValue", "purgeData", "exportReport", "exportResult"}[i]
}
func ParseCaseAction(v string) (any, error) {
    result := CONTENTEXPORT_CASEACTION
    switch v {
        case "contentExport":
            result = CONTENTEXPORT_CASEACTION
        case "applyTags":
            result = APPLYTAGS_CASEACTION
        case "convertToPdf":
            result = CONVERTTOPDF_CASEACTION
        case "index":
            result = INDEX_CASEACTION
        case "estimateStatistics":
            result = ESTIMATESTATISTICS_CASEACTION
        case "addToReviewSet":
            result = ADDTOREVIEWSET_CASEACTION
        case "holdUpdate":
            result = HOLDUPDATE_CASEACTION
        case "unknownFutureValue":
            result = UNKNOWNFUTUREVALUE_CASEACTION
        case "purgeData":
            result = PURGEDATA_CASEACTION
        case "exportReport":
            result = EXPORTREPORT_CASEACTION
        case "exportResult":
            result = EXPORTRESULT_CASEACTION
        default:
            return nil, nil
    }
    return &result, nil
}
func SerializeCaseAction(values []CaseAction) []string {
    result := make([]string, len(values))
    for i, v := range values {
        result[i] = v.String()
    }
    return result
}
func (i CaseAction) isMultiValue() bool {
    return false
}
