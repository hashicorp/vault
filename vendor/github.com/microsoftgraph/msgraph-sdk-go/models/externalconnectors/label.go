package externalconnectors
type Label int

const (
    TITLE_LABEL Label = iota
    URL_LABEL
    CREATEDBY_LABEL
    LASTMODIFIEDBY_LABEL
    AUTHORS_LABEL
    CREATEDDATETIME_LABEL
    LASTMODIFIEDDATETIME_LABEL
    FILENAME_LABEL
    FILEEXTENSION_LABEL
    UNKNOWNFUTUREVALUE_LABEL
    ICONURL_LABEL
)

func (i Label) String() string {
    return []string{"title", "url", "createdBy", "lastModifiedBy", "authors", "createdDateTime", "lastModifiedDateTime", "fileName", "fileExtension", "unknownFutureValue", "iconUrl"}[i]
}
func ParseLabel(v string) (any, error) {
    result := TITLE_LABEL
    switch v {
        case "title":
            result = TITLE_LABEL
        case "url":
            result = URL_LABEL
        case "createdBy":
            result = CREATEDBY_LABEL
        case "lastModifiedBy":
            result = LASTMODIFIEDBY_LABEL
        case "authors":
            result = AUTHORS_LABEL
        case "createdDateTime":
            result = CREATEDDATETIME_LABEL
        case "lastModifiedDateTime":
            result = LASTMODIFIEDDATETIME_LABEL
        case "fileName":
            result = FILENAME_LABEL
        case "fileExtension":
            result = FILEEXTENSION_LABEL
        case "unknownFutureValue":
            result = UNKNOWNFUTUREVALUE_LABEL
        case "iconUrl":
            result = ICONURL_LABEL
        default:
            return nil, nil
    }
    return &result, nil
}
func SerializeLabel(values []Label) []string {
    result := make([]string, len(values))
    for i, v := range values {
        result[i] = v.String()
    }
    return result
}
func (i Label) isMultiValue() bool {
    return false
}
