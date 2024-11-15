package models
type AppliedConditionalAccessPolicyResult int

const (
    SUCCESS_APPLIEDCONDITIONALACCESSPOLICYRESULT AppliedConditionalAccessPolicyResult = iota
    FAILURE_APPLIEDCONDITIONALACCESSPOLICYRESULT
    NOTAPPLIED_APPLIEDCONDITIONALACCESSPOLICYRESULT
    NOTENABLED_APPLIEDCONDITIONALACCESSPOLICYRESULT
    UNKNOWN_APPLIEDCONDITIONALACCESSPOLICYRESULT
    UNKNOWNFUTUREVALUE_APPLIEDCONDITIONALACCESSPOLICYRESULT
    REPORTONLYSUCCESS_APPLIEDCONDITIONALACCESSPOLICYRESULT
    REPORTONLYFAILURE_APPLIEDCONDITIONALACCESSPOLICYRESULT
    REPORTONLYNOTAPPLIED_APPLIEDCONDITIONALACCESSPOLICYRESULT
    REPORTONLYINTERRUPTED_APPLIEDCONDITIONALACCESSPOLICYRESULT
)

func (i AppliedConditionalAccessPolicyResult) String() string {
    return []string{"success", "failure", "notApplied", "notEnabled", "unknown", "unknownFutureValue", "reportOnlySuccess", "reportOnlyFailure", "reportOnlyNotApplied", "reportOnlyInterrupted"}[i]
}
func ParseAppliedConditionalAccessPolicyResult(v string) (any, error) {
    result := SUCCESS_APPLIEDCONDITIONALACCESSPOLICYRESULT
    switch v {
        case "success":
            result = SUCCESS_APPLIEDCONDITIONALACCESSPOLICYRESULT
        case "failure":
            result = FAILURE_APPLIEDCONDITIONALACCESSPOLICYRESULT
        case "notApplied":
            result = NOTAPPLIED_APPLIEDCONDITIONALACCESSPOLICYRESULT
        case "notEnabled":
            result = NOTENABLED_APPLIEDCONDITIONALACCESSPOLICYRESULT
        case "unknown":
            result = UNKNOWN_APPLIEDCONDITIONALACCESSPOLICYRESULT
        case "unknownFutureValue":
            result = UNKNOWNFUTUREVALUE_APPLIEDCONDITIONALACCESSPOLICYRESULT
        case "reportOnlySuccess":
            result = REPORTONLYSUCCESS_APPLIEDCONDITIONALACCESSPOLICYRESULT
        case "reportOnlyFailure":
            result = REPORTONLYFAILURE_APPLIEDCONDITIONALACCESSPOLICYRESULT
        case "reportOnlyNotApplied":
            result = REPORTONLYNOTAPPLIED_APPLIEDCONDITIONALACCESSPOLICYRESULT
        case "reportOnlyInterrupted":
            result = REPORTONLYINTERRUPTED_APPLIEDCONDITIONALACCESSPOLICYRESULT
        default:
            return nil, nil
    }
    return &result, nil
}
func SerializeAppliedConditionalAccessPolicyResult(values []AppliedConditionalAccessPolicyResult) []string {
    result := make([]string, len(values))
    for i, v := range values {
        result[i] = v.String()
    }
    return result
}
func (i AppliedConditionalAccessPolicyResult) isMultiValue() bool {
    return false
}
