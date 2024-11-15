package models
type Tone int

const (
    TONE0_TONE Tone = iota
    TONE1_TONE
    TONE2_TONE
    TONE3_TONE
    TONE4_TONE
    TONE5_TONE
    TONE6_TONE
    TONE7_TONE
    TONE8_TONE
    TONE9_TONE
    STAR_TONE
    POUND_TONE
    A_TONE
    B_TONE
    C_TONE
    D_TONE
    FLASH_TONE
)

func (i Tone) String() string {
    return []string{"tone0", "tone1", "tone2", "tone3", "tone4", "tone5", "tone6", "tone7", "tone8", "tone9", "star", "pound", "a", "b", "c", "d", "flash"}[i]
}
func ParseTone(v string) (any, error) {
    result := TONE0_TONE
    switch v {
        case "tone0":
            result = TONE0_TONE
        case "tone1":
            result = TONE1_TONE
        case "tone2":
            result = TONE2_TONE
        case "tone3":
            result = TONE3_TONE
        case "tone4":
            result = TONE4_TONE
        case "tone5":
            result = TONE5_TONE
        case "tone6":
            result = TONE6_TONE
        case "tone7":
            result = TONE7_TONE
        case "tone8":
            result = TONE8_TONE
        case "tone9":
            result = TONE9_TONE
        case "star":
            result = STAR_TONE
        case "pound":
            result = POUND_TONE
        case "a":
            result = A_TONE
        case "b":
            result = B_TONE
        case "c":
            result = C_TONE
        case "d":
            result = D_TONE
        case "flash":
            result = FLASH_TONE
        default:
            return nil, nil
    }
    return &result, nil
}
func SerializeTone(values []Tone) []string {
    result := make([]string, len(values))
    for i, v := range values {
        result[i] = v.String()
    }
    return result
}
func (i Tone) isMultiValue() bool {
    return false
}
