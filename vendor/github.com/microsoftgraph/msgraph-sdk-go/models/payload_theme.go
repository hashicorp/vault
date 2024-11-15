package models
type PayloadTheme int

const (
    UNKNOWN_PAYLOADTHEME PayloadTheme = iota
    OTHER_PAYLOADTHEME
    ACCOUNTACTIVATION_PAYLOADTHEME
    ACCOUNTVERIFICATION_PAYLOADTHEME
    BILLING_PAYLOADTHEME
    CLEANUPMAIL_PAYLOADTHEME
    CONTROVERSIAL_PAYLOADTHEME
    DOCUMENTRECEIVED_PAYLOADTHEME
    EXPENSE_PAYLOADTHEME
    FAX_PAYLOADTHEME
    FINANCEREPORT_PAYLOADTHEME
    INCOMINGMESSAGES_PAYLOADTHEME
    INVOICE_PAYLOADTHEME
    ITEMRECEIVED_PAYLOADTHEME
    LOGINALERT_PAYLOADTHEME
    MAILRECEIVED_PAYLOADTHEME
    PASSWORD_PAYLOADTHEME
    PAYMENT_PAYLOADTHEME
    PAYROLL_PAYLOADTHEME
    PERSONALIZEDOFFER_PAYLOADTHEME
    QUARANTINE_PAYLOADTHEME
    REMOTEWORK_PAYLOADTHEME
    REVIEWMESSAGE_PAYLOADTHEME
    SECURITYUPDATE_PAYLOADTHEME
    SERVICESUSPENDED_PAYLOADTHEME
    SIGNATUREREQUIRED_PAYLOADTHEME
    UPGRADEMAILBOXSTORAGE_PAYLOADTHEME
    VERIFYMAILBOX_PAYLOADTHEME
    VOICEMAIL_PAYLOADTHEME
    ADVERTISEMENT_PAYLOADTHEME
    EMPLOYEEENGAGEMENT_PAYLOADTHEME
    UNKNOWNFUTUREVALUE_PAYLOADTHEME
)

func (i PayloadTheme) String() string {
    return []string{"unknown", "other", "accountActivation", "accountVerification", "billing", "cleanUpMail", "controversial", "documentReceived", "expense", "fax", "financeReport", "incomingMessages", "invoice", "itemReceived", "loginAlert", "mailReceived", "password", "payment", "payroll", "personalizedOffer", "quarantine", "remoteWork", "reviewMessage", "securityUpdate", "serviceSuspended", "signatureRequired", "upgradeMailboxStorage", "verifyMailbox", "voicemail", "advertisement", "employeeEngagement", "unknownFutureValue"}[i]
}
func ParsePayloadTheme(v string) (any, error) {
    result := UNKNOWN_PAYLOADTHEME
    switch v {
        case "unknown":
            result = UNKNOWN_PAYLOADTHEME
        case "other":
            result = OTHER_PAYLOADTHEME
        case "accountActivation":
            result = ACCOUNTACTIVATION_PAYLOADTHEME
        case "accountVerification":
            result = ACCOUNTVERIFICATION_PAYLOADTHEME
        case "billing":
            result = BILLING_PAYLOADTHEME
        case "cleanUpMail":
            result = CLEANUPMAIL_PAYLOADTHEME
        case "controversial":
            result = CONTROVERSIAL_PAYLOADTHEME
        case "documentReceived":
            result = DOCUMENTRECEIVED_PAYLOADTHEME
        case "expense":
            result = EXPENSE_PAYLOADTHEME
        case "fax":
            result = FAX_PAYLOADTHEME
        case "financeReport":
            result = FINANCEREPORT_PAYLOADTHEME
        case "incomingMessages":
            result = INCOMINGMESSAGES_PAYLOADTHEME
        case "invoice":
            result = INVOICE_PAYLOADTHEME
        case "itemReceived":
            result = ITEMRECEIVED_PAYLOADTHEME
        case "loginAlert":
            result = LOGINALERT_PAYLOADTHEME
        case "mailReceived":
            result = MAILRECEIVED_PAYLOADTHEME
        case "password":
            result = PASSWORD_PAYLOADTHEME
        case "payment":
            result = PAYMENT_PAYLOADTHEME
        case "payroll":
            result = PAYROLL_PAYLOADTHEME
        case "personalizedOffer":
            result = PERSONALIZEDOFFER_PAYLOADTHEME
        case "quarantine":
            result = QUARANTINE_PAYLOADTHEME
        case "remoteWork":
            result = REMOTEWORK_PAYLOADTHEME
        case "reviewMessage":
            result = REVIEWMESSAGE_PAYLOADTHEME
        case "securityUpdate":
            result = SECURITYUPDATE_PAYLOADTHEME
        case "serviceSuspended":
            result = SERVICESUSPENDED_PAYLOADTHEME
        case "signatureRequired":
            result = SIGNATUREREQUIRED_PAYLOADTHEME
        case "upgradeMailboxStorage":
            result = UPGRADEMAILBOXSTORAGE_PAYLOADTHEME
        case "verifyMailbox":
            result = VERIFYMAILBOX_PAYLOADTHEME
        case "voicemail":
            result = VOICEMAIL_PAYLOADTHEME
        case "advertisement":
            result = ADVERTISEMENT_PAYLOADTHEME
        case "employeeEngagement":
            result = EMPLOYEEENGAGEMENT_PAYLOADTHEME
        case "unknownFutureValue":
            result = UNKNOWNFUTUREVALUE_PAYLOADTHEME
        default:
            return nil, nil
    }
    return &result, nil
}
func SerializePayloadTheme(values []PayloadTheme) []string {
    result := make([]string, len(values))
    for i, v := range values {
        result[i] = v.String()
    }
    return result
}
func (i PayloadTheme) isMultiValue() bool {
    return false
}
