package models
// The current TeamViewer connector status
type RemoteAssistanceOnboardingStatus int

const (
    // The status reported when there is no active TeamViewer connector configured or active
    NOTONBOARDED_REMOTEASSISTANCEONBOARDINGSTATUS RemoteAssistanceOnboardingStatus = iota
    // The status reported when the system has initiated a TeamViewer connection, but the service has not yet completed the confirmation of a connector
    ONBOARDING_REMOTEASSISTANCEONBOARDINGSTATUS
    // The status reported when the system has successfully exchanged account information with TeamViewer and can now initiate remote assistance sessions with clients
    ONBOARDED_REMOTEASSISTANCEONBOARDINGSTATUS
)

func (i RemoteAssistanceOnboardingStatus) String() string {
    return []string{"notOnboarded", "onboarding", "onboarded"}[i]
}
func ParseRemoteAssistanceOnboardingStatus(v string) (any, error) {
    result := NOTONBOARDED_REMOTEASSISTANCEONBOARDINGSTATUS
    switch v {
        case "notOnboarded":
            result = NOTONBOARDED_REMOTEASSISTANCEONBOARDINGSTATUS
        case "onboarding":
            result = ONBOARDING_REMOTEASSISTANCEONBOARDINGSTATUS
        case "onboarded":
            result = ONBOARDED_REMOTEASSISTANCEONBOARDINGSTATUS
        default:
            return nil, nil
    }
    return &result, nil
}
func SerializeRemoteAssistanceOnboardingStatus(values []RemoteAssistanceOnboardingStatus) []string {
    result := make([]string, len(values))
    for i, v := range values {
        result[i] = v.String()
    }
    return result
}
func (i RemoteAssistanceOnboardingStatus) isMultiValue() bool {
    return false
}
