package models
// Possible values for applicationGuardBlockClipboardSharingType
type ApplicationGuardBlockClipboardSharingType int

const (
    // Not Configured
    NOTCONFIGURED_APPLICATIONGUARDBLOCKCLIPBOARDSHARINGTYPE ApplicationGuardBlockClipboardSharingType = iota
    // Block clipboard to share data both from Host to Container and from Container to Host
    BLOCKBOTH_APPLICATIONGUARDBLOCKCLIPBOARDSHARINGTYPE
    // Block clipboard to share data from Host to Container
    BLOCKHOSTTOCONTAINER_APPLICATIONGUARDBLOCKCLIPBOARDSHARINGTYPE
    // Block clipboard to share data from Container to Host
    BLOCKCONTAINERTOHOST_APPLICATIONGUARDBLOCKCLIPBOARDSHARINGTYPE
    // Block clipboard to share data neither from Host to Container nor from Container to Host
    BLOCKNONE_APPLICATIONGUARDBLOCKCLIPBOARDSHARINGTYPE
)

func (i ApplicationGuardBlockClipboardSharingType) String() string {
    return []string{"notConfigured", "blockBoth", "blockHostToContainer", "blockContainerToHost", "blockNone"}[i]
}
func ParseApplicationGuardBlockClipboardSharingType(v string) (any, error) {
    result := NOTCONFIGURED_APPLICATIONGUARDBLOCKCLIPBOARDSHARINGTYPE
    switch v {
        case "notConfigured":
            result = NOTCONFIGURED_APPLICATIONGUARDBLOCKCLIPBOARDSHARINGTYPE
        case "blockBoth":
            result = BLOCKBOTH_APPLICATIONGUARDBLOCKCLIPBOARDSHARINGTYPE
        case "blockHostToContainer":
            result = BLOCKHOSTTOCONTAINER_APPLICATIONGUARDBLOCKCLIPBOARDSHARINGTYPE
        case "blockContainerToHost":
            result = BLOCKCONTAINERTOHOST_APPLICATIONGUARDBLOCKCLIPBOARDSHARINGTYPE
        case "blockNone":
            result = BLOCKNONE_APPLICATIONGUARDBLOCKCLIPBOARDSHARINGTYPE
        default:
            return nil, nil
    }
    return &result, nil
}
func SerializeApplicationGuardBlockClipboardSharingType(values []ApplicationGuardBlockClipboardSharingType) []string {
    result := make([]string, len(values))
    for i, v := range values {
        result[i] = v.String()
    }
    return result
}
func (i ApplicationGuardBlockClipboardSharingType) isMultiValue() bool {
    return false
}
