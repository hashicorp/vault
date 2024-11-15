package models
type IdentityUserFlowAttributeInputType int

const (
    TEXTBOX_IDENTITYUSERFLOWATTRIBUTEINPUTTYPE IdentityUserFlowAttributeInputType = iota
    DATETIMEDROPDOWN_IDENTITYUSERFLOWATTRIBUTEINPUTTYPE
    RADIOSINGLESELECT_IDENTITYUSERFLOWATTRIBUTEINPUTTYPE
    DROPDOWNSINGLESELECT_IDENTITYUSERFLOWATTRIBUTEINPUTTYPE
    EMAILBOX_IDENTITYUSERFLOWATTRIBUTEINPUTTYPE
    CHECKBOXMULTISELECT_IDENTITYUSERFLOWATTRIBUTEINPUTTYPE
)

func (i IdentityUserFlowAttributeInputType) String() string {
    return []string{"textBox", "dateTimeDropdown", "radioSingleSelect", "dropdownSingleSelect", "emailBox", "checkboxMultiSelect"}[i]
}
func ParseIdentityUserFlowAttributeInputType(v string) (any, error) {
    result := TEXTBOX_IDENTITYUSERFLOWATTRIBUTEINPUTTYPE
    switch v {
        case "textBox":
            result = TEXTBOX_IDENTITYUSERFLOWATTRIBUTEINPUTTYPE
        case "dateTimeDropdown":
            result = DATETIMEDROPDOWN_IDENTITYUSERFLOWATTRIBUTEINPUTTYPE
        case "radioSingleSelect":
            result = RADIOSINGLESELECT_IDENTITYUSERFLOWATTRIBUTEINPUTTYPE
        case "dropdownSingleSelect":
            result = DROPDOWNSINGLESELECT_IDENTITYUSERFLOWATTRIBUTEINPUTTYPE
        case "emailBox":
            result = EMAILBOX_IDENTITYUSERFLOWATTRIBUTEINPUTTYPE
        case "checkboxMultiSelect":
            result = CHECKBOXMULTISELECT_IDENTITYUSERFLOWATTRIBUTEINPUTTYPE
        default:
            return nil, nil
    }
    return &result, nil
}
func SerializeIdentityUserFlowAttributeInputType(values []IdentityUserFlowAttributeInputType) []string {
    result := make([]string, len(values))
    for i, v := range values {
        result[i] = v.String()
    }
    return result
}
func (i IdentityUserFlowAttributeInputType) isMultiValue() bool {
    return false
}
