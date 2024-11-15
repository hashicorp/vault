package models
type CourseStatus int

const (
    NOTSTARTED_COURSESTATUS CourseStatus = iota
    INPROGRESS_COURSESTATUS
    COMPLETED_COURSESTATUS
    UNKNOWNFUTUREVALUE_COURSESTATUS
)

func (i CourseStatus) String() string {
    return []string{"notStarted", "inProgress", "completed", "unknownFutureValue"}[i]
}
func ParseCourseStatus(v string) (any, error) {
    result := NOTSTARTED_COURSESTATUS
    switch v {
        case "notStarted":
            result = NOTSTARTED_COURSESTATUS
        case "inProgress":
            result = INPROGRESS_COURSESTATUS
        case "completed":
            result = COMPLETED_COURSESTATUS
        case "unknownFutureValue":
            result = UNKNOWNFUTUREVALUE_COURSESTATUS
        default:
            return nil, nil
    }
    return &result, nil
}
func SerializeCourseStatus(values []CourseStatus) []string {
    result := make([]string, len(values))
    for i, v := range values {
        result[i] = v.String()
    }
    return result
}
func (i CourseStatus) isMultiValue() bool {
    return false
}
