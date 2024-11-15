package models
// Types of communityPrivacy.
type CommunityPrivacy int

const (
    // Any user from the tenant can join and participate in the community.
    PUBLIC_COMMUNITYPRIVACY CommunityPrivacy = iota
    // A community administrator must add tenant users to the community before they can participate.
    PRIVATE_COMMUNITYPRIVACY
    // A marker value for members added after the release of this API.
    UNKNOWNFUTUREVALUE_COMMUNITYPRIVACY
)

func (i CommunityPrivacy) String() string {
    return []string{"public", "private", "unknownFutureValue"}[i]
}
func ParseCommunityPrivacy(v string) (any, error) {
    result := PUBLIC_COMMUNITYPRIVACY
    switch v {
        case "public":
            result = PUBLIC_COMMUNITYPRIVACY
        case "private":
            result = PRIVATE_COMMUNITYPRIVACY
        case "unknownFutureValue":
            result = UNKNOWNFUTUREVALUE_COMMUNITYPRIVACY
        default:
            return nil, nil
    }
    return &result, nil
}
func SerializeCommunityPrivacy(values []CommunityPrivacy) []string {
    result := make([]string, len(values))
    for i, v := range values {
        result[i] = v.String()
    }
    return result
}
func (i CommunityPrivacy) isMultiValue() bool {
    return false
}
