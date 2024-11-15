package models
// The enum to specify the channels for Microsoft Edge apps.
type MicrosoftEdgeChannel int

const (
    // The Dev Channel is intended to help you plan and develop with the latest capabilities of Microsoft Edge.
    DEV_MICROSOFTEDGECHANNEL MicrosoftEdgeChannel = iota
    // The Beta Channel is intended for production deployment to a representative sample set of users. New features ship about every 4 weeks. Security and quality updates ship as needed.
    BETA_MICROSOFTEDGECHANNEL
    // The Stable Channel is intended for broad deployment within organizations, and it's the channel that most users should be on. New features ship about every 4 weeks. Security and quality updates ship as needed.
    STABLE_MICROSOFTEDGECHANNEL
    // Evolvable enumeration sentinel value. Do not use.
    UNKNOWNFUTUREVALUE_MICROSOFTEDGECHANNEL
)

func (i MicrosoftEdgeChannel) String() string {
    return []string{"dev", "beta", "stable", "unknownFutureValue"}[i]
}
func ParseMicrosoftEdgeChannel(v string) (any, error) {
    result := DEV_MICROSOFTEDGECHANNEL
    switch v {
        case "dev":
            result = DEV_MICROSOFTEDGECHANNEL
        case "beta":
            result = BETA_MICROSOFTEDGECHANNEL
        case "stable":
            result = STABLE_MICROSOFTEDGECHANNEL
        case "unknownFutureValue":
            result = UNKNOWNFUTUREVALUE_MICROSOFTEDGECHANNEL
        default:
            return nil, nil
    }
    return &result, nil
}
func SerializeMicrosoftEdgeChannel(values []MicrosoftEdgeChannel) []string {
    result := make([]string, len(values))
    for i, v := range values {
        result[i] = v.String()
    }
    return result
}
func (i MicrosoftEdgeChannel) isMultiValue() bool {
    return false
}
