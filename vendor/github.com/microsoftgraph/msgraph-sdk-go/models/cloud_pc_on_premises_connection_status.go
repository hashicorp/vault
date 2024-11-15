package models
type CloudPcOnPremisesConnectionStatus int

const (
    PENDING_CLOUDPCONPREMISESCONNECTIONSTATUS CloudPcOnPremisesConnectionStatus = iota
    RUNNING_CLOUDPCONPREMISESCONNECTIONSTATUS
    PASSED_CLOUDPCONPREMISESCONNECTIONSTATUS
    FAILED_CLOUDPCONPREMISESCONNECTIONSTATUS
    WARNING_CLOUDPCONPREMISESCONNECTIONSTATUS
    INFORMATIONAL_CLOUDPCONPREMISESCONNECTIONSTATUS
    UNKNOWNFUTUREVALUE_CLOUDPCONPREMISESCONNECTIONSTATUS
)

func (i CloudPcOnPremisesConnectionStatus) String() string {
    return []string{"pending", "running", "passed", "failed", "warning", "informational", "unknownFutureValue"}[i]
}
func ParseCloudPcOnPremisesConnectionStatus(v string) (any, error) {
    result := PENDING_CLOUDPCONPREMISESCONNECTIONSTATUS
    switch v {
        case "pending":
            result = PENDING_CLOUDPCONPREMISESCONNECTIONSTATUS
        case "running":
            result = RUNNING_CLOUDPCONPREMISESCONNECTIONSTATUS
        case "passed":
            result = PASSED_CLOUDPCONPREMISESCONNECTIONSTATUS
        case "failed":
            result = FAILED_CLOUDPCONPREMISESCONNECTIONSTATUS
        case "warning":
            result = WARNING_CLOUDPCONPREMISESCONNECTIONSTATUS
        case "informational":
            result = INFORMATIONAL_CLOUDPCONPREMISESCONNECTIONSTATUS
        case "unknownFutureValue":
            result = UNKNOWNFUTUREVALUE_CLOUDPCONPREMISESCONNECTIONSTATUS
        default:
            return nil, nil
    }
    return &result, nil
}
func SerializeCloudPcOnPremisesConnectionStatus(values []CloudPcOnPremisesConnectionStatus) []string {
    result := make([]string, len(values))
    for i, v := range values {
        result[i] = v.String()
    }
    return result
}
func (i CloudPcOnPremisesConnectionStatus) isMultiValue() bool {
    return false
}
