package models
type RiskDetail int

const (
    NONE_RISKDETAIL RiskDetail = iota
    ADMINGENERATEDTEMPORARYPASSWORD_RISKDETAIL
    USERPERFORMEDSECUREDPASSWORDCHANGE_RISKDETAIL
    USERPERFORMEDSECUREDPASSWORDRESET_RISKDETAIL
    ADMINCONFIRMEDSIGNINSAFE_RISKDETAIL
    AICONFIRMEDSIGNINSAFE_RISKDETAIL
    USERPASSEDMFADRIVENBYRISKBASEDPOLICY_RISKDETAIL
    ADMINDISMISSEDALLRISKFORUSER_RISKDETAIL
    ADMINCONFIRMEDSIGNINCOMPROMISED_RISKDETAIL
    HIDDEN_RISKDETAIL
    ADMINCONFIRMEDUSERCOMPROMISED_RISKDETAIL
    UNKNOWNFUTUREVALUE_RISKDETAIL
    M365DADMINDISMISSEDDETECTION_RISKDETAIL
    ADMINCONFIRMEDSERVICEPRINCIPALCOMPROMISED_RISKDETAIL
    ADMINDISMISSEDALLRISKFORSERVICEPRINCIPAL_RISKDETAIL
    USERCHANGEDPASSWORDONPREMISES_RISKDETAIL
    ADMINDISMISSEDRISKFORSIGNIN_RISKDETAIL
    ADMINCONFIRMEDACCOUNTSAFE_RISKDETAIL
)

func (i RiskDetail) String() string {
    return []string{"none", "adminGeneratedTemporaryPassword", "userPerformedSecuredPasswordChange", "userPerformedSecuredPasswordReset", "adminConfirmedSigninSafe", "aiConfirmedSigninSafe", "userPassedMFADrivenByRiskBasedPolicy", "adminDismissedAllRiskForUser", "adminConfirmedSigninCompromised", "hidden", "adminConfirmedUserCompromised", "unknownFutureValue", "m365DAdminDismissedDetection", "adminConfirmedServicePrincipalCompromised", "adminDismissedAllRiskForServicePrincipal", "userChangedPasswordOnPremises", "adminDismissedRiskForSignIn", "adminConfirmedAccountSafe"}[i]
}
func ParseRiskDetail(v string) (any, error) {
    result := NONE_RISKDETAIL
    switch v {
        case "none":
            result = NONE_RISKDETAIL
        case "adminGeneratedTemporaryPassword":
            result = ADMINGENERATEDTEMPORARYPASSWORD_RISKDETAIL
        case "userPerformedSecuredPasswordChange":
            result = USERPERFORMEDSECUREDPASSWORDCHANGE_RISKDETAIL
        case "userPerformedSecuredPasswordReset":
            result = USERPERFORMEDSECUREDPASSWORDRESET_RISKDETAIL
        case "adminConfirmedSigninSafe":
            result = ADMINCONFIRMEDSIGNINSAFE_RISKDETAIL
        case "aiConfirmedSigninSafe":
            result = AICONFIRMEDSIGNINSAFE_RISKDETAIL
        case "userPassedMFADrivenByRiskBasedPolicy":
            result = USERPASSEDMFADRIVENBYRISKBASEDPOLICY_RISKDETAIL
        case "adminDismissedAllRiskForUser":
            result = ADMINDISMISSEDALLRISKFORUSER_RISKDETAIL
        case "adminConfirmedSigninCompromised":
            result = ADMINCONFIRMEDSIGNINCOMPROMISED_RISKDETAIL
        case "hidden":
            result = HIDDEN_RISKDETAIL
        case "adminConfirmedUserCompromised":
            result = ADMINCONFIRMEDUSERCOMPROMISED_RISKDETAIL
        case "unknownFutureValue":
            result = UNKNOWNFUTUREVALUE_RISKDETAIL
        case "m365DAdminDismissedDetection":
            result = M365DADMINDISMISSEDDETECTION_RISKDETAIL
        case "adminConfirmedServicePrincipalCompromised":
            result = ADMINCONFIRMEDSERVICEPRINCIPALCOMPROMISED_RISKDETAIL
        case "adminDismissedAllRiskForServicePrincipal":
            result = ADMINDISMISSEDALLRISKFORSERVICEPRINCIPAL_RISKDETAIL
        case "userChangedPasswordOnPremises":
            result = USERCHANGEDPASSWORDONPREMISES_RISKDETAIL
        case "adminDismissedRiskForSignIn":
            result = ADMINDISMISSEDRISKFORSIGNIN_RISKDETAIL
        case "adminConfirmedAccountSafe":
            result = ADMINCONFIRMEDACCOUNTSAFE_RISKDETAIL
        default:
            return nil, nil
    }
    return &result, nil
}
func SerializeRiskDetail(values []RiskDetail) []string {
    result := make([]string, len(values))
    for i, v := range values {
        result[i] = v.String()
    }
    return result
}
func (i RiskDetail) isMultiValue() bool {
    return false
}
