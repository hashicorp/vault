package models
type VirtualEventStatus int

const (
    DRAFT_VIRTUALEVENTSTATUS VirtualEventStatus = iota
    PUBLISHED_VIRTUALEVENTSTATUS
    CANCELED_VIRTUALEVENTSTATUS
    UNKNOWNFUTUREVALUE_VIRTUALEVENTSTATUS
)

func (i VirtualEventStatus) String() string {
    return []string{"draft", "published", "canceled", "unknownFutureValue"}[i]
}
func ParseVirtualEventStatus(v string) (any, error) {
    result := DRAFT_VIRTUALEVENTSTATUS
    switch v {
        case "draft":
            result = DRAFT_VIRTUALEVENTSTATUS
        case "published":
            result = PUBLISHED_VIRTUALEVENTSTATUS
        case "canceled":
            result = CANCELED_VIRTUALEVENTSTATUS
        case "unknownFutureValue":
            result = UNKNOWNFUTUREVALUE_VIRTUALEVENTSTATUS
        default:
            return nil, nil
    }
    return &result, nil
}
func SerializeVirtualEventStatus(values []VirtualEventStatus) []string {
    result := make([]string, len(values))
    for i, v := range values {
        result[i] = v.String()
    }
    return result
}
func (i VirtualEventStatus) isMultiValue() bool {
    return false
}
