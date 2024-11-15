package models
type AccessPackageRequestType int

const (
    NOTSPECIFIED_ACCESSPACKAGEREQUESTTYPE AccessPackageRequestType = iota
    USERADD_ACCESSPACKAGEREQUESTTYPE
    USERUPDATE_ACCESSPACKAGEREQUESTTYPE
    USERREMOVE_ACCESSPACKAGEREQUESTTYPE
    ADMINADD_ACCESSPACKAGEREQUESTTYPE
    ADMINUPDATE_ACCESSPACKAGEREQUESTTYPE
    ADMINREMOVE_ACCESSPACKAGEREQUESTTYPE
    SYSTEMADD_ACCESSPACKAGEREQUESTTYPE
    SYSTEMUPDATE_ACCESSPACKAGEREQUESTTYPE
    SYSTEMREMOVE_ACCESSPACKAGEREQUESTTYPE
    ONBEHALFADD_ACCESSPACKAGEREQUESTTYPE
    UNKNOWNFUTUREVALUE_ACCESSPACKAGEREQUESTTYPE
)

func (i AccessPackageRequestType) String() string {
    return []string{"notSpecified", "userAdd", "userUpdate", "userRemove", "adminAdd", "adminUpdate", "adminRemove", "systemAdd", "systemUpdate", "systemRemove", "onBehalfAdd", "unknownFutureValue"}[i]
}
func ParseAccessPackageRequestType(v string) (any, error) {
    result := NOTSPECIFIED_ACCESSPACKAGEREQUESTTYPE
    switch v {
        case "notSpecified":
            result = NOTSPECIFIED_ACCESSPACKAGEREQUESTTYPE
        case "userAdd":
            result = USERADD_ACCESSPACKAGEREQUESTTYPE
        case "userUpdate":
            result = USERUPDATE_ACCESSPACKAGEREQUESTTYPE
        case "userRemove":
            result = USERREMOVE_ACCESSPACKAGEREQUESTTYPE
        case "adminAdd":
            result = ADMINADD_ACCESSPACKAGEREQUESTTYPE
        case "adminUpdate":
            result = ADMINUPDATE_ACCESSPACKAGEREQUESTTYPE
        case "adminRemove":
            result = ADMINREMOVE_ACCESSPACKAGEREQUESTTYPE
        case "systemAdd":
            result = SYSTEMADD_ACCESSPACKAGEREQUESTTYPE
        case "systemUpdate":
            result = SYSTEMUPDATE_ACCESSPACKAGEREQUESTTYPE
        case "systemRemove":
            result = SYSTEMREMOVE_ACCESSPACKAGEREQUESTTYPE
        case "onBehalfAdd":
            result = ONBEHALFADD_ACCESSPACKAGEREQUESTTYPE
        case "unknownFutureValue":
            result = UNKNOWNFUTUREVALUE_ACCESSPACKAGEREQUESTTYPE
        default:
            return nil, nil
    }
    return &result, nil
}
func SerializeAccessPackageRequestType(values []AccessPackageRequestType) []string {
    result := make([]string, len(values))
    for i, v := range values {
        result[i] = v.String()
    }
    return result
}
func (i AccessPackageRequestType) isMultiValue() bool {
    return false
}
