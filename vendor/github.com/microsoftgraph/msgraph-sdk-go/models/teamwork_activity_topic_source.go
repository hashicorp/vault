package models
type TeamworkActivityTopicSource int

const (
    ENTITYURL_TEAMWORKACTIVITYTOPICSOURCE TeamworkActivityTopicSource = iota
    TEXT_TEAMWORKACTIVITYTOPICSOURCE
)

func (i TeamworkActivityTopicSource) String() string {
    return []string{"entityUrl", "text"}[i]
}
func ParseTeamworkActivityTopicSource(v string) (any, error) {
    result := ENTITYURL_TEAMWORKACTIVITYTOPICSOURCE
    switch v {
        case "entityUrl":
            result = ENTITYURL_TEAMWORKACTIVITYTOPICSOURCE
        case "text":
            result = TEXT_TEAMWORKACTIVITYTOPICSOURCE
        default:
            return nil, nil
    }
    return &result, nil
}
func SerializeTeamworkActivityTopicSource(values []TeamworkActivityTopicSource) []string {
    result := make([]string, len(values))
    for i, v := range values {
        result[i] = v.String()
    }
    return result
}
func (i TeamworkActivityTopicSource) isMultiValue() bool {
    return false
}
