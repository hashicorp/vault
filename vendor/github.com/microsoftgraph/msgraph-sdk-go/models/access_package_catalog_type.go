package models
type AccessPackageCatalogType int

const (
    USERMANAGED_ACCESSPACKAGECATALOGTYPE AccessPackageCatalogType = iota
    SERVICEDEFAULT_ACCESSPACKAGECATALOGTYPE
    SERVICEMANAGED_ACCESSPACKAGECATALOGTYPE
    UNKNOWNFUTUREVALUE_ACCESSPACKAGECATALOGTYPE
)

func (i AccessPackageCatalogType) String() string {
    return []string{"userManaged", "serviceDefault", "serviceManaged", "unknownFutureValue"}[i]
}
func ParseAccessPackageCatalogType(v string) (any, error) {
    result := USERMANAGED_ACCESSPACKAGECATALOGTYPE
    switch v {
        case "userManaged":
            result = USERMANAGED_ACCESSPACKAGECATALOGTYPE
        case "serviceDefault":
            result = SERVICEDEFAULT_ACCESSPACKAGECATALOGTYPE
        case "serviceManaged":
            result = SERVICEMANAGED_ACCESSPACKAGECATALOGTYPE
        case "unknownFutureValue":
            result = UNKNOWNFUTUREVALUE_ACCESSPACKAGECATALOGTYPE
        default:
            return nil, nil
    }
    return &result, nil
}
func SerializeAccessPackageCatalogType(values []AccessPackageCatalogType) []string {
    result := make([]string, len(values))
    for i, v := range values {
        result[i] = v.String()
    }
    return result
}
func (i AccessPackageCatalogType) isMultiValue() bool {
    return false
}
