package security
type AlertDetermination int

const (
    UNKNOWN_ALERTDETERMINATION AlertDetermination = iota
    APT_ALERTDETERMINATION
    MALWARE_ALERTDETERMINATION
    SECURITYPERSONNEL_ALERTDETERMINATION
    SECURITYTESTING_ALERTDETERMINATION
    UNWANTEDSOFTWARE_ALERTDETERMINATION
    OTHER_ALERTDETERMINATION
    MULTISTAGEDATTACK_ALERTDETERMINATION
    COMPROMISEDACCOUNT_ALERTDETERMINATION
    PHISHING_ALERTDETERMINATION
    MALICIOUSUSERACTIVITY_ALERTDETERMINATION
    NOTMALICIOUS_ALERTDETERMINATION
    NOTENOUGHDATATOVALIDATE_ALERTDETERMINATION
    CONFIRMEDACTIVITY_ALERTDETERMINATION
    LINEOFBUSINESSAPPLICATION_ALERTDETERMINATION
    UNKNOWNFUTUREVALUE_ALERTDETERMINATION
)

func (i AlertDetermination) String() string {
    return []string{"unknown", "apt", "malware", "securityPersonnel", "securityTesting", "unwantedSoftware", "other", "multiStagedAttack", "compromisedAccount", "phishing", "maliciousUserActivity", "notMalicious", "notEnoughDataToValidate", "confirmedActivity", "lineOfBusinessApplication", "unknownFutureValue"}[i]
}
func ParseAlertDetermination(v string) (any, error) {
    result := UNKNOWN_ALERTDETERMINATION
    switch v {
        case "unknown":
            result = UNKNOWN_ALERTDETERMINATION
        case "apt":
            result = APT_ALERTDETERMINATION
        case "malware":
            result = MALWARE_ALERTDETERMINATION
        case "securityPersonnel":
            result = SECURITYPERSONNEL_ALERTDETERMINATION
        case "securityTesting":
            result = SECURITYTESTING_ALERTDETERMINATION
        case "unwantedSoftware":
            result = UNWANTEDSOFTWARE_ALERTDETERMINATION
        case "other":
            result = OTHER_ALERTDETERMINATION
        case "multiStagedAttack":
            result = MULTISTAGEDATTACK_ALERTDETERMINATION
        case "compromisedAccount":
            result = COMPROMISEDACCOUNT_ALERTDETERMINATION
        case "phishing":
            result = PHISHING_ALERTDETERMINATION
        case "maliciousUserActivity":
            result = MALICIOUSUSERACTIVITY_ALERTDETERMINATION
        case "notMalicious":
            result = NOTMALICIOUS_ALERTDETERMINATION
        case "notEnoughDataToValidate":
            result = NOTENOUGHDATATOVALIDATE_ALERTDETERMINATION
        case "confirmedActivity":
            result = CONFIRMEDACTIVITY_ALERTDETERMINATION
        case "lineOfBusinessApplication":
            result = LINEOFBUSINESSAPPLICATION_ALERTDETERMINATION
        case "unknownFutureValue":
            result = UNKNOWNFUTUREVALUE_ALERTDETERMINATION
        default:
            return nil, nil
    }
    return &result, nil
}
func SerializeAlertDetermination(values []AlertDetermination) []string {
    result := make([]string, len(values))
    for i, v := range values {
        result[i] = v.String()
    }
    return result
}
func (i AlertDetermination) isMultiValue() bool {
    return false
}
