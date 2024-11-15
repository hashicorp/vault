package models
type UserPurpose int

const (
    USER_USERPURPOSE UserPurpose = iota
    LINKED_USERPURPOSE
    SHARED_USERPURPOSE
    ROOM_USERPURPOSE
    EQUIPMENT_USERPURPOSE
    OTHERS_USERPURPOSE
    UNKNOWNFUTUREVALUE_USERPURPOSE
)

func (i UserPurpose) String() string {
    return []string{"user", "linked", "shared", "room", "equipment", "others", "unknownFutureValue"}[i]
}
func ParseUserPurpose(v string) (any, error) {
    result := USER_USERPURPOSE
    switch v {
        case "user":
            result = USER_USERPURPOSE
        case "linked":
            result = LINKED_USERPURPOSE
        case "shared":
            result = SHARED_USERPURPOSE
        case "room":
            result = ROOM_USERPURPOSE
        case "equipment":
            result = EQUIPMENT_USERPURPOSE
        case "others":
            result = OTHERS_USERPURPOSE
        case "unknownFutureValue":
            result = UNKNOWNFUTUREVALUE_USERPURPOSE
        default:
            return nil, nil
    }
    return &result, nil
}
func SerializeUserPurpose(values []UserPurpose) []string {
    result := make([]string, len(values))
    for i, v := range values {
        result[i] = v.String()
    }
    return result
}
func (i UserPurpose) isMultiValue() bool {
    return false
}
