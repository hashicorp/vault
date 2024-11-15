package identitygovernance
import (
    "math"
    "strings"
)
type LifecycleTaskCategory int

const (
    JOINER_LIFECYCLETASKCATEGORY = 1
    LEAVER_LIFECYCLETASKCATEGORY = 2
    UNKNOWNFUTUREVALUE_LIFECYCLETASKCATEGORY = 4
    MOVER_LIFECYCLETASKCATEGORY = 8
)

func (i LifecycleTaskCategory) String() string {
    var values []string
    options := []string{"joiner", "leaver", "unknownFutureValue", "mover"}
    for p := 0; p < 4; p++ {
        mantis := LifecycleTaskCategory(int(math.Pow(2, float64(p))))
        if i&mantis == mantis {
            values = append(values, options[p])
        }
    }
    return strings.Join(values, ",")
}
func ParseLifecycleTaskCategory(v string) (any, error) {
    var result LifecycleTaskCategory
    values := strings.Split(v, ",")
    for _, str := range values {
        switch str {
            case "joiner":
                result |= JOINER_LIFECYCLETASKCATEGORY
            case "leaver":
                result |= LEAVER_LIFECYCLETASKCATEGORY
            case "unknownFutureValue":
                result |= UNKNOWNFUTUREVALUE_LIFECYCLETASKCATEGORY
            case "mover":
                result |= MOVER_LIFECYCLETASKCATEGORY
            default:
                return nil, nil
        }
    }
    return &result, nil
}
func SerializeLifecycleTaskCategory(values []LifecycleTaskCategory) []string {
    result := make([]string, len(values))
    for i, v := range values {
        result[i] = v.String()
    }
    return result
}
func (i LifecycleTaskCategory) isMultiValue() bool {
    return true
}
