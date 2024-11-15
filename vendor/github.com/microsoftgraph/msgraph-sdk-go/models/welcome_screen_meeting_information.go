package models
// Possible values for welcome screen meeting information.
type WelcomeScreenMeetingInformation int

const (
    // User Defined, default value, no intent.
    USERDEFINED_WELCOMESCREENMEETINGINFORMATION WelcomeScreenMeetingInformation = iota
    // Show organizer and time only.
    SHOWORGANIZERANDTIMEONLY_WELCOMESCREENMEETINGINFORMATION
    // Show organizer, time and subject (subject is hidden for private meetings).
    SHOWORGANIZERANDTIMEANDSUBJECT_WELCOMESCREENMEETINGINFORMATION
)

func (i WelcomeScreenMeetingInformation) String() string {
    return []string{"userDefined", "showOrganizerAndTimeOnly", "showOrganizerAndTimeAndSubject"}[i]
}
func ParseWelcomeScreenMeetingInformation(v string) (any, error) {
    result := USERDEFINED_WELCOMESCREENMEETINGINFORMATION
    switch v {
        case "userDefined":
            result = USERDEFINED_WELCOMESCREENMEETINGINFORMATION
        case "showOrganizerAndTimeOnly":
            result = SHOWORGANIZERANDTIMEONLY_WELCOMESCREENMEETINGINFORMATION
        case "showOrganizerAndTimeAndSubject":
            result = SHOWORGANIZERANDTIMEANDSUBJECT_WELCOMESCREENMEETINGINFORMATION
        default:
            return nil, nil
    }
    return &result, nil
}
func SerializeWelcomeScreenMeetingInformation(values []WelcomeScreenMeetingInformation) []string {
    result := make([]string, len(values))
    for i, v := range values {
        result[i] = v.String()
    }
    return result
}
func (i WelcomeScreenMeetingInformation) isMultiValue() bool {
    return false
}
