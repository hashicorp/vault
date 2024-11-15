package models
type AllowInvitesFrom int

const (
    NONE_ALLOWINVITESFROM AllowInvitesFrom = iota
    ADMINSANDGUESTINVITERS_ALLOWINVITESFROM
    ADMINSGUESTINVITERSANDALLMEMBERS_ALLOWINVITESFROM
    EVERYONE_ALLOWINVITESFROM
    UNKNOWNFUTUREVALUE_ALLOWINVITESFROM
)

func (i AllowInvitesFrom) String() string {
    return []string{"none", "adminsAndGuestInviters", "adminsGuestInvitersAndAllMembers", "everyone", "unknownFutureValue"}[i]
}
func ParseAllowInvitesFrom(v string) (any, error) {
    result := NONE_ALLOWINVITESFROM
    switch v {
        case "none":
            result = NONE_ALLOWINVITESFROM
        case "adminsAndGuestInviters":
            result = ADMINSANDGUESTINVITERS_ALLOWINVITESFROM
        case "adminsGuestInvitersAndAllMembers":
            result = ADMINSGUESTINVITERSANDALLMEMBERS_ALLOWINVITESFROM
        case "everyone":
            result = EVERYONE_ALLOWINVITESFROM
        case "unknownFutureValue":
            result = UNKNOWNFUTUREVALUE_ALLOWINVITESFROM
        default:
            return nil, nil
    }
    return &result, nil
}
func SerializeAllowInvitesFrom(values []AllowInvitesFrom) []string {
    result := make([]string, len(values))
    for i, v := range values {
        result[i] = v.String()
    }
    return result
}
func (i AllowInvitesFrom) isMultiValue() bool {
    return false
}
