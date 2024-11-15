package security
type MailboxConfigurationType int

const (
    MAILFORWARDINGRULE_MAILBOXCONFIGURATIONTYPE MailboxConfigurationType = iota
    OWASETTINGS_MAILBOXCONFIGURATIONTYPE
    EWSSETTINGS_MAILBOXCONFIGURATIONTYPE
    MAILDELEGATION_MAILBOXCONFIGURATIONTYPE
    USERINBOXRULE_MAILBOXCONFIGURATIONTYPE
    UNKNOWNFUTUREVALUE_MAILBOXCONFIGURATIONTYPE
)

func (i MailboxConfigurationType) String() string {
    return []string{"mailForwardingRule", "owaSettings", "ewsSettings", "mailDelegation", "userInboxRule", "unknownFutureValue"}[i]
}
func ParseMailboxConfigurationType(v string) (any, error) {
    result := MAILFORWARDINGRULE_MAILBOXCONFIGURATIONTYPE
    switch v {
        case "mailForwardingRule":
            result = MAILFORWARDINGRULE_MAILBOXCONFIGURATIONTYPE
        case "owaSettings":
            result = OWASETTINGS_MAILBOXCONFIGURATIONTYPE
        case "ewsSettings":
            result = EWSSETTINGS_MAILBOXCONFIGURATIONTYPE
        case "mailDelegation":
            result = MAILDELEGATION_MAILBOXCONFIGURATIONTYPE
        case "userInboxRule":
            result = USERINBOXRULE_MAILBOXCONFIGURATIONTYPE
        case "unknownFutureValue":
            result = UNKNOWNFUTUREVALUE_MAILBOXCONFIGURATIONTYPE
        default:
            return nil, nil
    }
    return &result, nil
}
func SerializeMailboxConfigurationType(values []MailboxConfigurationType) []string {
    result := make([]string, len(values))
    for i, v := range values {
        result[i] = v.String()
    }
    return result
}
func (i MailboxConfigurationType) isMultiValue() bool {
    return false
}
