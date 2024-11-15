package models
type MailDestinationRoutingReason int

const (
    NONE_MAILDESTINATIONROUTINGREASON MailDestinationRoutingReason = iota
    MAILFLOWRULE_MAILDESTINATIONROUTINGREASON
    SAFESENDER_MAILDESTINATIONROUTINGREASON
    BLOCKEDSENDER_MAILDESTINATIONROUTINGREASON
    ADVANCEDSPAMFILTERING_MAILDESTINATIONROUTINGREASON
    DOMAINALLOWLIST_MAILDESTINATIONROUTINGREASON
    DOMAINBLOCKLIST_MAILDESTINATIONROUTINGREASON
    NOTINADDRESSBOOK_MAILDESTINATIONROUTINGREASON
    FIRSTTIMESENDER_MAILDESTINATIONROUTINGREASON
    AUTOPURGETOINBOX_MAILDESTINATIONROUTINGREASON
    AUTOPURGETOJUNK_MAILDESTINATIONROUTINGREASON
    AUTOPURGETODELETED_MAILDESTINATIONROUTINGREASON
    OUTBOUND_MAILDESTINATIONROUTINGREASON
    NOTJUNK_MAILDESTINATIONROUTINGREASON
    JUNK_MAILDESTINATIONROUTINGREASON
    UNKNOWNFUTUREVALUE_MAILDESTINATIONROUTINGREASON
)

func (i MailDestinationRoutingReason) String() string {
    return []string{"none", "mailFlowRule", "safeSender", "blockedSender", "advancedSpamFiltering", "domainAllowList", "domainBlockList", "notInAddressBook", "firstTimeSender", "autoPurgeToInbox", "autoPurgeToJunk", "autoPurgeToDeleted", "outbound", "notJunk", "junk", "unknownFutureValue"}[i]
}
func ParseMailDestinationRoutingReason(v string) (any, error) {
    result := NONE_MAILDESTINATIONROUTINGREASON
    switch v {
        case "none":
            result = NONE_MAILDESTINATIONROUTINGREASON
        case "mailFlowRule":
            result = MAILFLOWRULE_MAILDESTINATIONROUTINGREASON
        case "safeSender":
            result = SAFESENDER_MAILDESTINATIONROUTINGREASON
        case "blockedSender":
            result = BLOCKEDSENDER_MAILDESTINATIONROUTINGREASON
        case "advancedSpamFiltering":
            result = ADVANCEDSPAMFILTERING_MAILDESTINATIONROUTINGREASON
        case "domainAllowList":
            result = DOMAINALLOWLIST_MAILDESTINATIONROUTINGREASON
        case "domainBlockList":
            result = DOMAINBLOCKLIST_MAILDESTINATIONROUTINGREASON
        case "notInAddressBook":
            result = NOTINADDRESSBOOK_MAILDESTINATIONROUTINGREASON
        case "firstTimeSender":
            result = FIRSTTIMESENDER_MAILDESTINATIONROUTINGREASON
        case "autoPurgeToInbox":
            result = AUTOPURGETOINBOX_MAILDESTINATIONROUTINGREASON
        case "autoPurgeToJunk":
            result = AUTOPURGETOJUNK_MAILDESTINATIONROUTINGREASON
        case "autoPurgeToDeleted":
            result = AUTOPURGETODELETED_MAILDESTINATIONROUTINGREASON
        case "outbound":
            result = OUTBOUND_MAILDESTINATIONROUTINGREASON
        case "notJunk":
            result = NOTJUNK_MAILDESTINATIONROUTINGREASON
        case "junk":
            result = JUNK_MAILDESTINATIONROUTINGREASON
        case "unknownFutureValue":
            result = UNKNOWNFUTUREVALUE_MAILDESTINATIONROUTINGREASON
        default:
            return nil, nil
    }
    return &result, nil
}
func SerializeMailDestinationRoutingReason(values []MailDestinationRoutingReason) []string {
    result := make([]string, len(values))
    for i, v := range values {
        result[i] = v.String()
    }
    return result
}
func (i MailDestinationRoutingReason) isMultiValue() bool {
    return false
}
