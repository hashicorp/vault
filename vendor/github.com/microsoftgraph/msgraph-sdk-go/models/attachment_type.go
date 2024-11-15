package models
type AttachmentType int

const (
    FILE_ATTACHMENTTYPE AttachmentType = iota
    ITEM_ATTACHMENTTYPE
    REFERENCE_ATTACHMENTTYPE
)

func (i AttachmentType) String() string {
    return []string{"file", "item", "reference"}[i]
}
func ParseAttachmentType(v string) (any, error) {
    result := FILE_ATTACHMENTTYPE
    switch v {
        case "file":
            result = FILE_ATTACHMENTTYPE
        case "item":
            result = ITEM_ATTACHMENTTYPE
        case "reference":
            result = REFERENCE_ATTACHMENTTYPE
        default:
            return nil, nil
    }
    return &result, nil
}
func SerializeAttachmentType(values []AttachmentType) []string {
    result := make([]string, len(values))
    for i, v := range values {
        result[i] = v.String()
    }
    return result
}
func (i AttachmentType) isMultiValue() bool {
    return false
}
