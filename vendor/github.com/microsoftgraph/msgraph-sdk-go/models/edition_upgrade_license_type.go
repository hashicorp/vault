package models
// Edition Upgrade License type
type EditionUpgradeLicenseType int

const (
    // Product Key Type
    PRODUCTKEY_EDITIONUPGRADELICENSETYPE EditionUpgradeLicenseType = iota
    // License File Type
    LICENSEFILE_EDITIONUPGRADELICENSETYPE
)

func (i EditionUpgradeLicenseType) String() string {
    return []string{"productKey", "licenseFile"}[i]
}
func ParseEditionUpgradeLicenseType(v string) (any, error) {
    result := PRODUCTKEY_EDITIONUPGRADELICENSETYPE
    switch v {
        case "productKey":
            result = PRODUCTKEY_EDITIONUPGRADELICENSETYPE
        case "licenseFile":
            result = LICENSEFILE_EDITIONUPGRADELICENSETYPE
        default:
            return nil, nil
    }
    return &result, nil
}
func SerializeEditionUpgradeLicenseType(values []EditionUpgradeLicenseType) []string {
    result := make([]string, len(values))
    for i, v := range values {
        result[i] = v.String()
    }
    return result
}
func (i EditionUpgradeLicenseType) isMultiValue() bool {
    return false
}
