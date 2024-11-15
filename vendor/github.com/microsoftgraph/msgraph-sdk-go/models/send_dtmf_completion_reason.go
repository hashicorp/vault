package models
type SendDtmfCompletionReason int

const (
    UNKNOWN_SENDDTMFCOMPLETIONREASON SendDtmfCompletionReason = iota
    COMPLETEDSUCCESSFULLY_SENDDTMFCOMPLETIONREASON
    MEDIAOPERATIONCANCELED_SENDDTMFCOMPLETIONREASON
    UNKNOWNFUTUREVALUE_SENDDTMFCOMPLETIONREASON
)

func (i SendDtmfCompletionReason) String() string {
    return []string{"unknown", "completedSuccessfully", "mediaOperationCanceled", "unknownFutureValue"}[i]
}
func ParseSendDtmfCompletionReason(v string) (any, error) {
    result := UNKNOWN_SENDDTMFCOMPLETIONREASON
    switch v {
        case "unknown":
            result = UNKNOWN_SENDDTMFCOMPLETIONREASON
        case "completedSuccessfully":
            result = COMPLETEDSUCCESSFULLY_SENDDTMFCOMPLETIONREASON
        case "mediaOperationCanceled":
            result = MEDIAOPERATIONCANCELED_SENDDTMFCOMPLETIONREASON
        case "unknownFutureValue":
            result = UNKNOWNFUTUREVALUE_SENDDTMFCOMPLETIONREASON
        default:
            return nil, nil
    }
    return &result, nil
}
func SerializeSendDtmfCompletionReason(values []SendDtmfCompletionReason) []string {
    result := make([]string, len(values))
    for i, v := range values {
        result[i] = v.String()
    }
    return result
}
func (i SendDtmfCompletionReason) isMultiValue() bool {
    return false
}
