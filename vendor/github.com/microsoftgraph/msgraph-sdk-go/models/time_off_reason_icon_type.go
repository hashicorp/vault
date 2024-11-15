package models
type TimeOffReasonIconType int

const (
    NONE_TIMEOFFREASONICONTYPE TimeOffReasonIconType = iota
    CAR_TIMEOFFREASONICONTYPE
    CALENDAR_TIMEOFFREASONICONTYPE
    RUNNING_TIMEOFFREASONICONTYPE
    PLANE_TIMEOFFREASONICONTYPE
    FIRSTAID_TIMEOFFREASONICONTYPE
    DOCTOR_TIMEOFFREASONICONTYPE
    NOTWORKING_TIMEOFFREASONICONTYPE
    CLOCK_TIMEOFFREASONICONTYPE
    JURYDUTY_TIMEOFFREASONICONTYPE
    GLOBE_TIMEOFFREASONICONTYPE
    CUP_TIMEOFFREASONICONTYPE
    PHONE_TIMEOFFREASONICONTYPE
    WEATHER_TIMEOFFREASONICONTYPE
    UMBRELLA_TIMEOFFREASONICONTYPE
    PIGGYBANK_TIMEOFFREASONICONTYPE
    DOG_TIMEOFFREASONICONTYPE
    CAKE_TIMEOFFREASONICONTYPE
    TRAFFICCONE_TIMEOFFREASONICONTYPE
    PIN_TIMEOFFREASONICONTYPE
    SUNNY_TIMEOFFREASONICONTYPE
    UNKNOWNFUTUREVALUE_TIMEOFFREASONICONTYPE
)

func (i TimeOffReasonIconType) String() string {
    return []string{"none", "car", "calendar", "running", "plane", "firstAid", "doctor", "notWorking", "clock", "juryDuty", "globe", "cup", "phone", "weather", "umbrella", "piggyBank", "dog", "cake", "trafficCone", "pin", "sunny", "unknownFutureValue"}[i]
}
func ParseTimeOffReasonIconType(v string) (any, error) {
    result := NONE_TIMEOFFREASONICONTYPE
    switch v {
        case "none":
            result = NONE_TIMEOFFREASONICONTYPE
        case "car":
            result = CAR_TIMEOFFREASONICONTYPE
        case "calendar":
            result = CALENDAR_TIMEOFFREASONICONTYPE
        case "running":
            result = RUNNING_TIMEOFFREASONICONTYPE
        case "plane":
            result = PLANE_TIMEOFFREASONICONTYPE
        case "firstAid":
            result = FIRSTAID_TIMEOFFREASONICONTYPE
        case "doctor":
            result = DOCTOR_TIMEOFFREASONICONTYPE
        case "notWorking":
            result = NOTWORKING_TIMEOFFREASONICONTYPE
        case "clock":
            result = CLOCK_TIMEOFFREASONICONTYPE
        case "juryDuty":
            result = JURYDUTY_TIMEOFFREASONICONTYPE
        case "globe":
            result = GLOBE_TIMEOFFREASONICONTYPE
        case "cup":
            result = CUP_TIMEOFFREASONICONTYPE
        case "phone":
            result = PHONE_TIMEOFFREASONICONTYPE
        case "weather":
            result = WEATHER_TIMEOFFREASONICONTYPE
        case "umbrella":
            result = UMBRELLA_TIMEOFFREASONICONTYPE
        case "piggyBank":
            result = PIGGYBANK_TIMEOFFREASONICONTYPE
        case "dog":
            result = DOG_TIMEOFFREASONICONTYPE
        case "cake":
            result = CAKE_TIMEOFFREASONICONTYPE
        case "trafficCone":
            result = TRAFFICCONE_TIMEOFFREASONICONTYPE
        case "pin":
            result = PIN_TIMEOFFREASONICONTYPE
        case "sunny":
            result = SUNNY_TIMEOFFREASONICONTYPE
        case "unknownFutureValue":
            result = UNKNOWNFUTUREVALUE_TIMEOFFREASONICONTYPE
        default:
            return nil, nil
    }
    return &result, nil
}
func SerializeTimeOffReasonIconType(values []TimeOffReasonIconType) []string {
    result := make([]string, len(values))
    for i, v := range values {
        result[i] = v.String()
    }
    return result
}
func (i TimeOffReasonIconType) isMultiValue() bool {
    return false
}
