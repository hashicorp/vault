package models
type AccessPackageCatalogState int

const (
    UNPUBLISHED_ACCESSPACKAGECATALOGSTATE AccessPackageCatalogState = iota
    PUBLISHED_ACCESSPACKAGECATALOGSTATE
    UNKNOWNFUTUREVALUE_ACCESSPACKAGECATALOGSTATE
)

func (i AccessPackageCatalogState) String() string {
    return []string{"unpublished", "published", "unknownFutureValue"}[i]
}
func ParseAccessPackageCatalogState(v string) (any, error) {
    result := UNPUBLISHED_ACCESSPACKAGECATALOGSTATE
    switch v {
        case "unpublished":
            result = UNPUBLISHED_ACCESSPACKAGECATALOGSTATE
        case "published":
            result = PUBLISHED_ACCESSPACKAGECATALOGSTATE
        case "unknownFutureValue":
            result = UNKNOWNFUTUREVALUE_ACCESSPACKAGECATALOGSTATE
        default:
            return nil, nil
    }
    return &result, nil
}
func SerializeAccessPackageCatalogState(values []AccessPackageCatalogState) []string {
    result := make([]string, len(values))
    for i, v := range values {
        result[i] = v.String()
    }
    return result
}
func (i AccessPackageCatalogState) isMultiValue() bool {
    return false
}
