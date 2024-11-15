package models
type LobbyBypassScope int

const (
    ORGANIZER_LOBBYBYPASSSCOPE LobbyBypassScope = iota
    ORGANIZATION_LOBBYBYPASSSCOPE
    ORGANIZATIONANDFEDERATED_LOBBYBYPASSSCOPE
    EVERYONE_LOBBYBYPASSSCOPE
    UNKNOWNFUTUREVALUE_LOBBYBYPASSSCOPE
    INVITED_LOBBYBYPASSSCOPE
    ORGANIZATIONEXCLUDINGGUESTS_LOBBYBYPASSSCOPE
)

func (i LobbyBypassScope) String() string {
    return []string{"organizer", "organization", "organizationAndFederated", "everyone", "unknownFutureValue", "invited", "organizationExcludingGuests"}[i]
}
func ParseLobbyBypassScope(v string) (any, error) {
    result := ORGANIZER_LOBBYBYPASSSCOPE
    switch v {
        case "organizer":
            result = ORGANIZER_LOBBYBYPASSSCOPE
        case "organization":
            result = ORGANIZATION_LOBBYBYPASSSCOPE
        case "organizationAndFederated":
            result = ORGANIZATIONANDFEDERATED_LOBBYBYPASSSCOPE
        case "everyone":
            result = EVERYONE_LOBBYBYPASSSCOPE
        case "unknownFutureValue":
            result = UNKNOWNFUTUREVALUE_LOBBYBYPASSSCOPE
        case "invited":
            result = INVITED_LOBBYBYPASSSCOPE
        case "organizationExcludingGuests":
            result = ORGANIZATIONEXCLUDINGGUESTS_LOBBYBYPASSSCOPE
        default:
            return nil, nil
    }
    return &result, nil
}
func SerializeLobbyBypassScope(values []LobbyBypassScope) []string {
    result := make([]string, len(values))
    for i, v := range values {
        result[i] = v.String()
    }
    return result
}
func (i LobbyBypassScope) isMultiValue() bool {
    return false
}
