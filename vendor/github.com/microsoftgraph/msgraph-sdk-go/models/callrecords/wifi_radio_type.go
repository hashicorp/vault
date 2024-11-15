package callrecords
type WifiRadioType int

const (
    UNKNOWN_WIFIRADIOTYPE WifiRadioType = iota
    WIFI80211A_WIFIRADIOTYPE
    WIFI80211B_WIFIRADIOTYPE
    WIFI80211G_WIFIRADIOTYPE
    WIFI80211N_WIFIRADIOTYPE
    WIFI80211AC_WIFIRADIOTYPE
    WIFI80211AX_WIFIRADIOTYPE
    UNKNOWNFUTUREVALUE_WIFIRADIOTYPE
)

func (i WifiRadioType) String() string {
    return []string{"unknown", "wifi80211a", "wifi80211b", "wifi80211g", "wifi80211n", "wifi80211ac", "wifi80211ax", "unknownFutureValue"}[i]
}
func ParseWifiRadioType(v string) (any, error) {
    result := UNKNOWN_WIFIRADIOTYPE
    switch v {
        case "unknown":
            result = UNKNOWN_WIFIRADIOTYPE
        case "wifi80211a":
            result = WIFI80211A_WIFIRADIOTYPE
        case "wifi80211b":
            result = WIFI80211B_WIFIRADIOTYPE
        case "wifi80211g":
            result = WIFI80211G_WIFIRADIOTYPE
        case "wifi80211n":
            result = WIFI80211N_WIFIRADIOTYPE
        case "wifi80211ac":
            result = WIFI80211AC_WIFIRADIOTYPE
        case "wifi80211ax":
            result = WIFI80211AX_WIFIRADIOTYPE
        case "unknownFutureValue":
            result = UNKNOWNFUTUREVALUE_WIFIRADIOTYPE
        default:
            return nil, nil
    }
    return &result, nil
}
func SerializeWifiRadioType(values []WifiRadioType) []string {
    result := make([]string, len(values))
    for i, v := range values {
        result[i] = v.String()
    }
    return result
}
func (i WifiRadioType) isMultiValue() bool {
    return false
}
