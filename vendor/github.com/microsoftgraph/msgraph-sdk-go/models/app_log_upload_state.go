package models
// AppLogUploadStatus
type AppLogUploadState int

const (
    // Default. Indicates that request is waiting to be processed or under processing.
    PENDING_APPLOGUPLOADSTATE AppLogUploadState = iota
    // Indicates that request is completed with file uploaded to Azure blob for download.
    COMPLETED_APPLOGUPLOADSTATE
    // Indicates that request is completed with file uploaded to Azure blob for download.
    FAILED_APPLOGUPLOADSTATE
    // Evolvable enumeration sentinel value. Do not use.
    UNKNOWNFUTUREVALUE_APPLOGUPLOADSTATE
)

func (i AppLogUploadState) String() string {
    return []string{"pending", "completed", "failed", "unknownFutureValue"}[i]
}
func ParseAppLogUploadState(v string) (any, error) {
    result := PENDING_APPLOGUPLOADSTATE
    switch v {
        case "pending":
            result = PENDING_APPLOGUPLOADSTATE
        case "completed":
            result = COMPLETED_APPLOGUPLOADSTATE
        case "failed":
            result = FAILED_APPLOGUPLOADSTATE
        case "unknownFutureValue":
            result = UNKNOWNFUTUREVALUE_APPLOGUPLOADSTATE
        default:
            return nil, nil
    }
    return &result, nil
}
func SerializeAppLogUploadState(values []AppLogUploadState) []string {
    result := make([]string, len(values))
    for i, v := range values {
        result[i] = v.String()
    }
    return result
}
func (i AppLogUploadState) isMultiValue() bool {
    return false
}
