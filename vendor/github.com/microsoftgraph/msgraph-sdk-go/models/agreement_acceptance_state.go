package models
type AgreementAcceptanceState int

const (
    ACCEPTED_AGREEMENTACCEPTANCESTATE AgreementAcceptanceState = iota
    DECLINED_AGREEMENTACCEPTANCESTATE
    UNKNOWNFUTUREVALUE_AGREEMENTACCEPTANCESTATE
)

func (i AgreementAcceptanceState) String() string {
    return []string{"accepted", "declined", "unknownFutureValue"}[i]
}
func ParseAgreementAcceptanceState(v string) (any, error) {
    result := ACCEPTED_AGREEMENTACCEPTANCESTATE
    switch v {
        case "accepted":
            result = ACCEPTED_AGREEMENTACCEPTANCESTATE
        case "declined":
            result = DECLINED_AGREEMENTACCEPTANCESTATE
        case "unknownFutureValue":
            result = UNKNOWNFUTUREVALUE_AGREEMENTACCEPTANCESTATE
        default:
            return nil, nil
    }
    return &result, nil
}
func SerializeAgreementAcceptanceState(values []AgreementAcceptanceState) []string {
    result := make([]string, len(values))
    for i, v := range values {
        result[i] = v.String()
    }
    return result
}
func (i AgreementAcceptanceState) isMultiValue() bool {
    return false
}
