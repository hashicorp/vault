package models
type PrintJobStateDetail int

const (
    UPLOADPENDING_PRINTJOBSTATEDETAIL PrintJobStateDetail = iota
    TRANSFORMING_PRINTJOBSTATEDETAIL
    COMPLETEDSUCCESSFULLY_PRINTJOBSTATEDETAIL
    COMPLETEDWITHWARNINGS_PRINTJOBSTATEDETAIL
    COMPLETEDWITHERRORS_PRINTJOBSTATEDETAIL
    RELEASEWAIT_PRINTJOBSTATEDETAIL
    INTERPRETING_PRINTJOBSTATEDETAIL
    UNKNOWNFUTUREVALUE_PRINTJOBSTATEDETAIL
)

func (i PrintJobStateDetail) String() string {
    return []string{"uploadPending", "transforming", "completedSuccessfully", "completedWithWarnings", "completedWithErrors", "releaseWait", "interpreting", "unknownFutureValue"}[i]
}
func ParsePrintJobStateDetail(v string) (any, error) {
    result := UPLOADPENDING_PRINTJOBSTATEDETAIL
    switch v {
        case "uploadPending":
            result = UPLOADPENDING_PRINTJOBSTATEDETAIL
        case "transforming":
            result = TRANSFORMING_PRINTJOBSTATEDETAIL
        case "completedSuccessfully":
            result = COMPLETEDSUCCESSFULLY_PRINTJOBSTATEDETAIL
        case "completedWithWarnings":
            result = COMPLETEDWITHWARNINGS_PRINTJOBSTATEDETAIL
        case "completedWithErrors":
            result = COMPLETEDWITHERRORS_PRINTJOBSTATEDETAIL
        case "releaseWait":
            result = RELEASEWAIT_PRINTJOBSTATEDETAIL
        case "interpreting":
            result = INTERPRETING_PRINTJOBSTATEDETAIL
        case "unknownFutureValue":
            result = UNKNOWNFUTUREVALUE_PRINTJOBSTATEDETAIL
        default:
            return nil, nil
    }
    return &result, nil
}
func SerializePrintJobStateDetail(values []PrintJobStateDetail) []string {
    result := make([]string, len(values))
    for i, v := range values {
        result[i] = v.String()
    }
    return result
}
func (i PrintJobStateDetail) isMultiValue() bool {
    return false
}
