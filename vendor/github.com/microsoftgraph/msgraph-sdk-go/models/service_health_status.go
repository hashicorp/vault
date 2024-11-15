package models
type ServiceHealthStatus int

const (
    SERVICEOPERATIONAL_SERVICEHEALTHSTATUS ServiceHealthStatus = iota
    INVESTIGATING_SERVICEHEALTHSTATUS
    RESTORINGSERVICE_SERVICEHEALTHSTATUS
    VERIFYINGSERVICE_SERVICEHEALTHSTATUS
    SERVICERESTORED_SERVICEHEALTHSTATUS
    POSTINCIDENTREVIEWPUBLISHED_SERVICEHEALTHSTATUS
    SERVICEDEGRADATION_SERVICEHEALTHSTATUS
    SERVICEINTERRUPTION_SERVICEHEALTHSTATUS
    EXTENDEDRECOVERY_SERVICEHEALTHSTATUS
    FALSEPOSITIVE_SERVICEHEALTHSTATUS
    INVESTIGATIONSUSPENDED_SERVICEHEALTHSTATUS
    RESOLVED_SERVICEHEALTHSTATUS
    MITIGATEDEXTERNAL_SERVICEHEALTHSTATUS
    MITIGATED_SERVICEHEALTHSTATUS
    RESOLVEDEXTERNAL_SERVICEHEALTHSTATUS
    CONFIRMED_SERVICEHEALTHSTATUS
    REPORTED_SERVICEHEALTHSTATUS
    UNKNOWNFUTUREVALUE_SERVICEHEALTHSTATUS
)

func (i ServiceHealthStatus) String() string {
    return []string{"serviceOperational", "investigating", "restoringService", "verifyingService", "serviceRestored", "postIncidentReviewPublished", "serviceDegradation", "serviceInterruption", "extendedRecovery", "falsePositive", "investigationSuspended", "resolved", "mitigatedExternal", "mitigated", "resolvedExternal", "confirmed", "reported", "unknownFutureValue"}[i]
}
func ParseServiceHealthStatus(v string) (any, error) {
    result := SERVICEOPERATIONAL_SERVICEHEALTHSTATUS
    switch v {
        case "serviceOperational":
            result = SERVICEOPERATIONAL_SERVICEHEALTHSTATUS
        case "investigating":
            result = INVESTIGATING_SERVICEHEALTHSTATUS
        case "restoringService":
            result = RESTORINGSERVICE_SERVICEHEALTHSTATUS
        case "verifyingService":
            result = VERIFYINGSERVICE_SERVICEHEALTHSTATUS
        case "serviceRestored":
            result = SERVICERESTORED_SERVICEHEALTHSTATUS
        case "postIncidentReviewPublished":
            result = POSTINCIDENTREVIEWPUBLISHED_SERVICEHEALTHSTATUS
        case "serviceDegradation":
            result = SERVICEDEGRADATION_SERVICEHEALTHSTATUS
        case "serviceInterruption":
            result = SERVICEINTERRUPTION_SERVICEHEALTHSTATUS
        case "extendedRecovery":
            result = EXTENDEDRECOVERY_SERVICEHEALTHSTATUS
        case "falsePositive":
            result = FALSEPOSITIVE_SERVICEHEALTHSTATUS
        case "investigationSuspended":
            result = INVESTIGATIONSUSPENDED_SERVICEHEALTHSTATUS
        case "resolved":
            result = RESOLVED_SERVICEHEALTHSTATUS
        case "mitigatedExternal":
            result = MITIGATEDEXTERNAL_SERVICEHEALTHSTATUS
        case "mitigated":
            result = MITIGATED_SERVICEHEALTHSTATUS
        case "resolvedExternal":
            result = RESOLVEDEXTERNAL_SERVICEHEALTHSTATUS
        case "confirmed":
            result = CONFIRMED_SERVICEHEALTHSTATUS
        case "reported":
            result = REPORTED_SERVICEHEALTHSTATUS
        case "unknownFutureValue":
            result = UNKNOWNFUTUREVALUE_SERVICEHEALTHSTATUS
        default:
            return nil, nil
    }
    return &result, nil
}
func SerializeServiceHealthStatus(values []ServiceHealthStatus) []string {
    result := make([]string, len(values))
    for i, v := range values {
        result[i] = v.String()
    }
    return result
}
func (i ServiceHealthStatus) isMultiValue() bool {
    return false
}
