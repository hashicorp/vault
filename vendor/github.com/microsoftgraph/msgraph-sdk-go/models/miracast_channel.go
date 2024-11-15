package models
// Possible values for Miracast channel.
type MiracastChannel int

const (
    // User Defined, default value, no intent.
    USERDEFINED_MIRACASTCHANNEL MiracastChannel = iota
    // One.
    ONE_MIRACASTCHANNEL
    // Two.
    TWO_MIRACASTCHANNEL
    // Three.
    THREE_MIRACASTCHANNEL
    // Four.
    FOUR_MIRACASTCHANNEL
    // Five.
    FIVE_MIRACASTCHANNEL
    // Six.
    SIX_MIRACASTCHANNEL
    // Seven.
    SEVEN_MIRACASTCHANNEL
    // Eight.
    EIGHT_MIRACASTCHANNEL
    // Nine.
    NINE_MIRACASTCHANNEL
    // Ten.
    TEN_MIRACASTCHANNEL
    // Eleven.
    ELEVEN_MIRACASTCHANNEL
    // Thirty-Six.
    THIRTYSIX_MIRACASTCHANNEL
    // Forty.
    FORTY_MIRACASTCHANNEL
    // Forty-Four.
    FORTYFOUR_MIRACASTCHANNEL
    // Forty-Eight.
    FORTYEIGHT_MIRACASTCHANNEL
    // OneHundredForty-Nine.
    ONEHUNDREDFORTYNINE_MIRACASTCHANNEL
    // OneHundredFifty-Three.
    ONEHUNDREDFIFTYTHREE_MIRACASTCHANNEL
    // OneHundredFifty-Seven.
    ONEHUNDREDFIFTYSEVEN_MIRACASTCHANNEL
    // OneHundredSixty-One.
    ONEHUNDREDSIXTYONE_MIRACASTCHANNEL
    // OneHundredSixty-Five.
    ONEHUNDREDSIXTYFIVE_MIRACASTCHANNEL
)

func (i MiracastChannel) String() string {
    return []string{"userDefined", "one", "two", "three", "four", "five", "six", "seven", "eight", "nine", "ten", "eleven", "thirtySix", "forty", "fortyFour", "fortyEight", "oneHundredFortyNine", "oneHundredFiftyThree", "oneHundredFiftySeven", "oneHundredSixtyOne", "oneHundredSixtyFive"}[i]
}
func ParseMiracastChannel(v string) (any, error) {
    result := USERDEFINED_MIRACASTCHANNEL
    switch v {
        case "userDefined":
            result = USERDEFINED_MIRACASTCHANNEL
        case "one":
            result = ONE_MIRACASTCHANNEL
        case "two":
            result = TWO_MIRACASTCHANNEL
        case "three":
            result = THREE_MIRACASTCHANNEL
        case "four":
            result = FOUR_MIRACASTCHANNEL
        case "five":
            result = FIVE_MIRACASTCHANNEL
        case "six":
            result = SIX_MIRACASTCHANNEL
        case "seven":
            result = SEVEN_MIRACASTCHANNEL
        case "eight":
            result = EIGHT_MIRACASTCHANNEL
        case "nine":
            result = NINE_MIRACASTCHANNEL
        case "ten":
            result = TEN_MIRACASTCHANNEL
        case "eleven":
            result = ELEVEN_MIRACASTCHANNEL
        case "thirtySix":
            result = THIRTYSIX_MIRACASTCHANNEL
        case "forty":
            result = FORTY_MIRACASTCHANNEL
        case "fortyFour":
            result = FORTYFOUR_MIRACASTCHANNEL
        case "fortyEight":
            result = FORTYEIGHT_MIRACASTCHANNEL
        case "oneHundredFortyNine":
            result = ONEHUNDREDFORTYNINE_MIRACASTCHANNEL
        case "oneHundredFiftyThree":
            result = ONEHUNDREDFIFTYTHREE_MIRACASTCHANNEL
        case "oneHundredFiftySeven":
            result = ONEHUNDREDFIFTYSEVEN_MIRACASTCHANNEL
        case "oneHundredSixtyOne":
            result = ONEHUNDREDSIXTYONE_MIRACASTCHANNEL
        case "oneHundredSixtyFive":
            result = ONEHUNDREDSIXTYFIVE_MIRACASTCHANNEL
        default:
            return nil, nil
    }
    return &result, nil
}
func SerializeMiracastChannel(values []MiracastChannel) []string {
    result := make([]string, len(values))
    for i, v := range values {
        result[i] = v.String()
    }
    return result
}
func (i MiracastChannel) isMultiValue() bool {
    return false
}
