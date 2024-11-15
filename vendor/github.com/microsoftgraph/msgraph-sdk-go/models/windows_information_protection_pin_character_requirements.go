package models
// Pin Character Requirements
type WindowsInformationProtectionPinCharacterRequirements int

const (
    // Not allow
    NOTALLOW_WINDOWSINFORMATIONPROTECTIONPINCHARACTERREQUIREMENTS WindowsInformationProtectionPinCharacterRequirements = iota
    // Require atleast one
    REQUIREATLEASTONE_WINDOWSINFORMATIONPROTECTIONPINCHARACTERREQUIREMENTS
    // Allow any number
    ALLOW_WINDOWSINFORMATIONPROTECTIONPINCHARACTERREQUIREMENTS
)

func (i WindowsInformationProtectionPinCharacterRequirements) String() string {
    return []string{"notAllow", "requireAtLeastOne", "allow"}[i]
}
func ParseWindowsInformationProtectionPinCharacterRequirements(v string) (any, error) {
    result := NOTALLOW_WINDOWSINFORMATIONPROTECTIONPINCHARACTERREQUIREMENTS
    switch v {
        case "notAllow":
            result = NOTALLOW_WINDOWSINFORMATIONPROTECTIONPINCHARACTERREQUIREMENTS
        case "requireAtLeastOne":
            result = REQUIREATLEASTONE_WINDOWSINFORMATIONPROTECTIONPINCHARACTERREQUIREMENTS
        case "allow":
            result = ALLOW_WINDOWSINFORMATIONPROTECTIONPINCHARACTERREQUIREMENTS
        default:
            return nil, nil
    }
    return &result, nil
}
func SerializeWindowsInformationProtectionPinCharacterRequirements(values []WindowsInformationProtectionPinCharacterRequirements) []string {
    result := make([]string, len(values))
    for i, v := range values {
        result[i] = v.String()
    }
    return result
}
func (i WindowsInformationProtectionPinCharacterRequirements) isMultiValue() bool {
    return false
}
