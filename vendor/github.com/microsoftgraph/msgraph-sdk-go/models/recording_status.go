package models
type RecordingStatus int

const (
    UNKNOWN_RECORDINGSTATUS RecordingStatus = iota
    NOTRECORDING_RECORDINGSTATUS
    RECORDING_RECORDINGSTATUS
    FAILED_RECORDINGSTATUS
    UNKNOWNFUTUREVALUE_RECORDINGSTATUS
)

func (i RecordingStatus) String() string {
    return []string{"unknown", "notRecording", "recording", "failed", "unknownFutureValue"}[i]
}
func ParseRecordingStatus(v string) (any, error) {
    result := UNKNOWN_RECORDINGSTATUS
    switch v {
        case "unknown":
            result = UNKNOWN_RECORDINGSTATUS
        case "notRecording":
            result = NOTRECORDING_RECORDINGSTATUS
        case "recording":
            result = RECORDING_RECORDINGSTATUS
        case "failed":
            result = FAILED_RECORDINGSTATUS
        case "unknownFutureValue":
            result = UNKNOWNFUTUREVALUE_RECORDINGSTATUS
        default:
            return nil, nil
    }
    return &result, nil
}
func SerializeRecordingStatus(values []RecordingStatus) []string {
    result := make([]string, len(values))
    for i, v := range values {
        result[i] = v.String()
    }
    return result
}
func (i RecordingStatus) isMultiValue() bool {
    return false
}
