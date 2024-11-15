package security
import (
    "math"
    "strings"
)
type DataSourceScopes int

const (
    NONE_DATASOURCESCOPES = 1
    ALLTENANTMAILBOXES_DATASOURCESCOPES = 2
    ALLTENANTSITES_DATASOURCESCOPES = 4
    ALLCASECUSTODIANS_DATASOURCESCOPES = 8
    ALLCASENONCUSTODIALDATASOURCES_DATASOURCESCOPES = 16
    UNKNOWNFUTUREVALUE_DATASOURCESCOPES = 32
)

func (i DataSourceScopes) String() string {
    var values []string
    options := []string{"none", "allTenantMailboxes", "allTenantSites", "allCaseCustodians", "allCaseNoncustodialDataSources", "unknownFutureValue"}
    for p := 0; p < 6; p++ {
        mantis := DataSourceScopes(int(math.Pow(2, float64(p))))
        if i&mantis == mantis {
            values = append(values, options[p])
        }
    }
    return strings.Join(values, ",")
}
func ParseDataSourceScopes(v string) (any, error) {
    var result DataSourceScopes
    values := strings.Split(v, ",")
    for _, str := range values {
        switch str {
            case "none":
                result |= NONE_DATASOURCESCOPES
            case "allTenantMailboxes":
                result |= ALLTENANTMAILBOXES_DATASOURCESCOPES
            case "allTenantSites":
                result |= ALLTENANTSITES_DATASOURCESCOPES
            case "allCaseCustodians":
                result |= ALLCASECUSTODIANS_DATASOURCESCOPES
            case "allCaseNoncustodialDataSources":
                result |= ALLCASENONCUSTODIALDATASOURCES_DATASOURCESCOPES
            case "unknownFutureValue":
                result |= UNKNOWNFUTUREVALUE_DATASOURCESCOPES
            default:
                return nil, nil
        }
    }
    return &result, nil
}
func SerializeDataSourceScopes(values []DataSourceScopes) []string {
    result := make([]string, len(values))
    for i, v := range values {
        result[i] = v.String()
    }
    return result
}
func (i DataSourceScopes) isMultiValue() bool {
    return true
}
