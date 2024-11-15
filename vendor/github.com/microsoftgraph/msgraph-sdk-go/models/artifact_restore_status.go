package models
type ArtifactRestoreStatus int

const (
    ADDED_ARTIFACTRESTORESTATUS ArtifactRestoreStatus = iota
    SCHEDULING_ARTIFACTRESTORESTATUS
    SCHEDULED_ARTIFACTRESTORESTATUS
    INPROGRESS_ARTIFACTRESTORESTATUS
    SUCCEEDED_ARTIFACTRESTORESTATUS
    FAILED_ARTIFACTRESTORESTATUS
    UNKNOWNFUTUREVALUE_ARTIFACTRESTORESTATUS
)

func (i ArtifactRestoreStatus) String() string {
    return []string{"added", "scheduling", "scheduled", "inProgress", "succeeded", "failed", "unknownFutureValue"}[i]
}
func ParseArtifactRestoreStatus(v string) (any, error) {
    result := ADDED_ARTIFACTRESTORESTATUS
    switch v {
        case "added":
            result = ADDED_ARTIFACTRESTORESTATUS
        case "scheduling":
            result = SCHEDULING_ARTIFACTRESTORESTATUS
        case "scheduled":
            result = SCHEDULED_ARTIFACTRESTORESTATUS
        case "inProgress":
            result = INPROGRESS_ARTIFACTRESTORESTATUS
        case "succeeded":
            result = SUCCEEDED_ARTIFACTRESTORESTATUS
        case "failed":
            result = FAILED_ARTIFACTRESTORESTATUS
        case "unknownFutureValue":
            result = UNKNOWNFUTUREVALUE_ARTIFACTRESTORESTATUS
        default:
            return nil, nil
    }
    return &result, nil
}
func SerializeArtifactRestoreStatus(values []ArtifactRestoreStatus) []string {
    result := make([]string, len(values))
    for i, v := range values {
        result[i] = v.String()
    }
    return result
}
func (i ArtifactRestoreStatus) isMultiValue() bool {
    return false
}
