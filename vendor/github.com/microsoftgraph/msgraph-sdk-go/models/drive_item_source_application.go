package models
type DriveItemSourceApplication int

const (
    TEAMS_DRIVEITEMSOURCEAPPLICATION DriveItemSourceApplication = iota
    YAMMER_DRIVEITEMSOURCEAPPLICATION
    SHAREPOINT_DRIVEITEMSOURCEAPPLICATION
    ONEDRIVE_DRIVEITEMSOURCEAPPLICATION
    STREAM_DRIVEITEMSOURCEAPPLICATION
    POWERPOINT_DRIVEITEMSOURCEAPPLICATION
    OFFICE_DRIVEITEMSOURCEAPPLICATION
    LOKI_DRIVEITEMSOURCEAPPLICATION
    LOOP_DRIVEITEMSOURCEAPPLICATION
    OTHER_DRIVEITEMSOURCEAPPLICATION
    UNKNOWNFUTUREVALUE_DRIVEITEMSOURCEAPPLICATION
)

func (i DriveItemSourceApplication) String() string {
    return []string{"teams", "yammer", "sharePoint", "oneDrive", "stream", "powerPoint", "office", "loki", "loop", "other", "unknownFutureValue"}[i]
}
func ParseDriveItemSourceApplication(v string) (any, error) {
    result := TEAMS_DRIVEITEMSOURCEAPPLICATION
    switch v {
        case "teams":
            result = TEAMS_DRIVEITEMSOURCEAPPLICATION
        case "yammer":
            result = YAMMER_DRIVEITEMSOURCEAPPLICATION
        case "sharePoint":
            result = SHAREPOINT_DRIVEITEMSOURCEAPPLICATION
        case "oneDrive":
            result = ONEDRIVE_DRIVEITEMSOURCEAPPLICATION
        case "stream":
            result = STREAM_DRIVEITEMSOURCEAPPLICATION
        case "powerPoint":
            result = POWERPOINT_DRIVEITEMSOURCEAPPLICATION
        case "office":
            result = OFFICE_DRIVEITEMSOURCEAPPLICATION
        case "loki":
            result = LOKI_DRIVEITEMSOURCEAPPLICATION
        case "loop":
            result = LOOP_DRIVEITEMSOURCEAPPLICATION
        case "other":
            result = OTHER_DRIVEITEMSOURCEAPPLICATION
        case "unknownFutureValue":
            result = UNKNOWNFUTUREVALUE_DRIVEITEMSOURCEAPPLICATION
        default:
            return nil, nil
    }
    return &result, nil
}
func SerializeDriveItemSourceApplication(values []DriveItemSourceApplication) []string {
    result := make([]string, len(values))
    for i, v := range values {
        result[i] = v.String()
    }
    return result
}
func (i DriveItemSourceApplication) isMultiValue() bool {
    return false
}
