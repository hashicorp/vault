package models
type CoachmarkLocationType int

const (
    UNKNOWN_COACHMARKLOCATIONTYPE CoachmarkLocationType = iota
    FROMEMAIL_COACHMARKLOCATIONTYPE
    SUBJECT_COACHMARKLOCATIONTYPE
    EXTERNALTAG_COACHMARKLOCATIONTYPE
    DISPLAYNAME_COACHMARKLOCATIONTYPE
    MESSAGEBODY_COACHMARKLOCATIONTYPE
    UNKNOWNFUTUREVALUE_COACHMARKLOCATIONTYPE
)

func (i CoachmarkLocationType) String() string {
    return []string{"unknown", "fromEmail", "subject", "externalTag", "displayName", "messageBody", "unknownFutureValue"}[i]
}
func ParseCoachmarkLocationType(v string) (any, error) {
    result := UNKNOWN_COACHMARKLOCATIONTYPE
    switch v {
        case "unknown":
            result = UNKNOWN_COACHMARKLOCATIONTYPE
        case "fromEmail":
            result = FROMEMAIL_COACHMARKLOCATIONTYPE
        case "subject":
            result = SUBJECT_COACHMARKLOCATIONTYPE
        case "externalTag":
            result = EXTERNALTAG_COACHMARKLOCATIONTYPE
        case "displayName":
            result = DISPLAYNAME_COACHMARKLOCATIONTYPE
        case "messageBody":
            result = MESSAGEBODY_COACHMARKLOCATIONTYPE
        case "unknownFutureValue":
            result = UNKNOWNFUTUREVALUE_COACHMARKLOCATIONTYPE
        default:
            return nil, nil
    }
    return &result, nil
}
func SerializeCoachmarkLocationType(values []CoachmarkLocationType) []string {
    result := make([]string, len(values))
    for i, v := range values {
        result[i] = v.String()
    }
    return result
}
func (i CoachmarkLocationType) isMultiValue() bool {
    return false
}
