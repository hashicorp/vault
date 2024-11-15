package models
type Modality int

const (
    AUDIO_MODALITY Modality = iota
    VIDEO_MODALITY
    VIDEOBASEDSCREENSHARING_MODALITY
    DATA_MODALITY
    UNKNOWNFUTUREVALUE_MODALITY
)

func (i Modality) String() string {
    return []string{"audio", "video", "videoBasedScreenSharing", "data", "unknownFutureValue"}[i]
}
func ParseModality(v string) (any, error) {
    result := AUDIO_MODALITY
    switch v {
        case "audio":
            result = AUDIO_MODALITY
        case "video":
            result = VIDEO_MODALITY
        case "videoBasedScreenSharing":
            result = VIDEOBASEDSCREENSHARING_MODALITY
        case "data":
            result = DATA_MODALITY
        case "unknownFutureValue":
            result = UNKNOWNFUTUREVALUE_MODALITY
        default:
            return nil, nil
    }
    return &result, nil
}
func SerializeModality(values []Modality) []string {
    result := make([]string, len(values))
    for i, v := range values {
        result[i] = v.String()
    }
    return result
}
func (i Modality) isMultiValue() bool {
    return false
}
