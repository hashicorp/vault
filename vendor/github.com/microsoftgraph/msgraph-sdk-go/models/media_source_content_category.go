package models
type MediaSourceContentCategory int

const (
    MEETING_MEDIASOURCECONTENTCATEGORY MediaSourceContentCategory = iota
    LIVESTREAM_MEDIASOURCECONTENTCATEGORY
    PRESENTATION_MEDIASOURCECONTENTCATEGORY
    SCREENRECORDING_MEDIASOURCECONTENTCATEGORY
    STORY_MEDIASOURCECONTENTCATEGORY
    PROFILE_MEDIASOURCECONTENTCATEGORY
    CHAT_MEDIASOURCECONTENTCATEGORY
    NOTE_MEDIASOURCECONTENTCATEGORY
    COMMENT_MEDIASOURCECONTENTCATEGORY
    UNKNOWNFUTUREVALUE_MEDIASOURCECONTENTCATEGORY
)

func (i MediaSourceContentCategory) String() string {
    return []string{"meeting", "liveStream", "presentation", "screenRecording", "story", "profile", "chat", "note", "comment", "unknownFutureValue"}[i]
}
func ParseMediaSourceContentCategory(v string) (any, error) {
    result := MEETING_MEDIASOURCECONTENTCATEGORY
    switch v {
        case "meeting":
            result = MEETING_MEDIASOURCECONTENTCATEGORY
        case "liveStream":
            result = LIVESTREAM_MEDIASOURCECONTENTCATEGORY
        case "presentation":
            result = PRESENTATION_MEDIASOURCECONTENTCATEGORY
        case "screenRecording":
            result = SCREENRECORDING_MEDIASOURCECONTENTCATEGORY
        case "story":
            result = STORY_MEDIASOURCECONTENTCATEGORY
        case "profile":
            result = PROFILE_MEDIASOURCECONTENTCATEGORY
        case "chat":
            result = CHAT_MEDIASOURCECONTENTCATEGORY
        case "note":
            result = NOTE_MEDIASOURCECONTENTCATEGORY
        case "comment":
            result = COMMENT_MEDIASOURCECONTENTCATEGORY
        case "unknownFutureValue":
            result = UNKNOWNFUTUREVALUE_MEDIASOURCECONTENTCATEGORY
        default:
            return nil, nil
    }
    return &result, nil
}
func SerializeMediaSourceContentCategory(values []MediaSourceContentCategory) []string {
    result := make([]string, len(values))
    for i, v := range values {
        result[i] = v.String()
    }
    return result
}
func (i MediaSourceContentCategory) isMultiValue() bool {
    return false
}
