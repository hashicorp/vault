package models
// Possible values for prompting user for samples submission.
type DefenderPromptForSampleSubmission int

const (
    // User Defined, default value, no intent.
    USERDEFINED_DEFENDERPROMPTFORSAMPLESUBMISSION DefenderPromptForSampleSubmission = iota
    // Always prompt.
    ALWAYSPROMPT_DEFENDERPROMPTFORSAMPLESUBMISSION
    // Send safe samples automatically.
    PROMPTBEFORESENDINGPERSONALDATA_DEFENDERPROMPTFORSAMPLESUBMISSION
    // Never send data.
    NEVERSENDDATA_DEFENDERPROMPTFORSAMPLESUBMISSION
    // Send all data without prompting.
    SENDALLDATAWITHOUTPROMPTING_DEFENDERPROMPTFORSAMPLESUBMISSION
)

func (i DefenderPromptForSampleSubmission) String() string {
    return []string{"userDefined", "alwaysPrompt", "promptBeforeSendingPersonalData", "neverSendData", "sendAllDataWithoutPrompting"}[i]
}
func ParseDefenderPromptForSampleSubmission(v string) (any, error) {
    result := USERDEFINED_DEFENDERPROMPTFORSAMPLESUBMISSION
    switch v {
        case "userDefined":
            result = USERDEFINED_DEFENDERPROMPTFORSAMPLESUBMISSION
        case "alwaysPrompt":
            result = ALWAYSPROMPT_DEFENDERPROMPTFORSAMPLESUBMISSION
        case "promptBeforeSendingPersonalData":
            result = PROMPTBEFORESENDINGPERSONALDATA_DEFENDERPROMPTFORSAMPLESUBMISSION
        case "neverSendData":
            result = NEVERSENDDATA_DEFENDERPROMPTFORSAMPLESUBMISSION
        case "sendAllDataWithoutPrompting":
            result = SENDALLDATAWITHOUTPROMPTING_DEFENDERPROMPTFORSAMPLESUBMISSION
        default:
            return nil, nil
    }
    return &result, nil
}
func SerializeDefenderPromptForSampleSubmission(values []DefenderPromptForSampleSubmission) []string {
    result := make([]string, len(values))
    for i, v := range values {
        result[i] = v.String()
    }
    return result
}
func (i DefenderPromptForSampleSubmission) isMultiValue() bool {
    return false
}
