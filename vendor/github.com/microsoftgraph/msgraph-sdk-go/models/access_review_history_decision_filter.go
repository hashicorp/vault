package models
type AccessReviewHistoryDecisionFilter int

const (
    APPROVE_ACCESSREVIEWHISTORYDECISIONFILTER AccessReviewHistoryDecisionFilter = iota
    DENY_ACCESSREVIEWHISTORYDECISIONFILTER
    NOTREVIEWED_ACCESSREVIEWHISTORYDECISIONFILTER
    DONTKNOW_ACCESSREVIEWHISTORYDECISIONFILTER
    NOTNOTIFIED_ACCESSREVIEWHISTORYDECISIONFILTER
    UNKNOWNFUTUREVALUE_ACCESSREVIEWHISTORYDECISIONFILTER
)

func (i AccessReviewHistoryDecisionFilter) String() string {
    return []string{"approve", "deny", "notReviewed", "dontKnow", "notNotified", "unknownFutureValue"}[i]
}
func ParseAccessReviewHistoryDecisionFilter(v string) (any, error) {
    result := APPROVE_ACCESSREVIEWHISTORYDECISIONFILTER
    switch v {
        case "approve":
            result = APPROVE_ACCESSREVIEWHISTORYDECISIONFILTER
        case "deny":
            result = DENY_ACCESSREVIEWHISTORYDECISIONFILTER
        case "notReviewed":
            result = NOTREVIEWED_ACCESSREVIEWHISTORYDECISIONFILTER
        case "dontKnow":
            result = DONTKNOW_ACCESSREVIEWHISTORYDECISIONFILTER
        case "notNotified":
            result = NOTNOTIFIED_ACCESSREVIEWHISTORYDECISIONFILTER
        case "unknownFutureValue":
            result = UNKNOWNFUTUREVALUE_ACCESSREVIEWHISTORYDECISIONFILTER
        default:
            return nil, nil
    }
    return &result, nil
}
func SerializeAccessReviewHistoryDecisionFilter(values []AccessReviewHistoryDecisionFilter) []string {
    result := make([]string, len(values))
    for i, v := range values {
        result[i] = v.String()
    }
    return result
}
func (i AccessReviewHistoryDecisionFilter) isMultiValue() bool {
    return false
}
