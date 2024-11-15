package models
// Contains properties for upload request states.
type MobileAppContentFileUploadState int

const (
    SUCCESS_MOBILEAPPCONTENTFILEUPLOADSTATE MobileAppContentFileUploadState = iota
    TRANSIENTERROR_MOBILEAPPCONTENTFILEUPLOADSTATE
    ERROR_MOBILEAPPCONTENTFILEUPLOADSTATE
    UNKNOWN_MOBILEAPPCONTENTFILEUPLOADSTATE
    AZURESTORAGEURIREQUESTSUCCESS_MOBILEAPPCONTENTFILEUPLOADSTATE
    AZURESTORAGEURIREQUESTPENDING_MOBILEAPPCONTENTFILEUPLOADSTATE
    AZURESTORAGEURIREQUESTFAILED_MOBILEAPPCONTENTFILEUPLOADSTATE
    AZURESTORAGEURIREQUESTTIMEDOUT_MOBILEAPPCONTENTFILEUPLOADSTATE
    AZURESTORAGEURIRENEWALSUCCESS_MOBILEAPPCONTENTFILEUPLOADSTATE
    AZURESTORAGEURIRENEWALPENDING_MOBILEAPPCONTENTFILEUPLOADSTATE
    AZURESTORAGEURIRENEWALFAILED_MOBILEAPPCONTENTFILEUPLOADSTATE
    AZURESTORAGEURIRENEWALTIMEDOUT_MOBILEAPPCONTENTFILEUPLOADSTATE
    COMMITFILESUCCESS_MOBILEAPPCONTENTFILEUPLOADSTATE
    COMMITFILEPENDING_MOBILEAPPCONTENTFILEUPLOADSTATE
    COMMITFILEFAILED_MOBILEAPPCONTENTFILEUPLOADSTATE
    COMMITFILETIMEDOUT_MOBILEAPPCONTENTFILEUPLOADSTATE
)

func (i MobileAppContentFileUploadState) String() string {
    return []string{"success", "transientError", "error", "unknown", "azureStorageUriRequestSuccess", "azureStorageUriRequestPending", "azureStorageUriRequestFailed", "azureStorageUriRequestTimedOut", "azureStorageUriRenewalSuccess", "azureStorageUriRenewalPending", "azureStorageUriRenewalFailed", "azureStorageUriRenewalTimedOut", "commitFileSuccess", "commitFilePending", "commitFileFailed", "commitFileTimedOut"}[i]
}
func ParseMobileAppContentFileUploadState(v string) (any, error) {
    result := SUCCESS_MOBILEAPPCONTENTFILEUPLOADSTATE
    switch v {
        case "success":
            result = SUCCESS_MOBILEAPPCONTENTFILEUPLOADSTATE
        case "transientError":
            result = TRANSIENTERROR_MOBILEAPPCONTENTFILEUPLOADSTATE
        case "error":
            result = ERROR_MOBILEAPPCONTENTFILEUPLOADSTATE
        case "unknown":
            result = UNKNOWN_MOBILEAPPCONTENTFILEUPLOADSTATE
        case "azureStorageUriRequestSuccess":
            result = AZURESTORAGEURIREQUESTSUCCESS_MOBILEAPPCONTENTFILEUPLOADSTATE
        case "azureStorageUriRequestPending":
            result = AZURESTORAGEURIREQUESTPENDING_MOBILEAPPCONTENTFILEUPLOADSTATE
        case "azureStorageUriRequestFailed":
            result = AZURESTORAGEURIREQUESTFAILED_MOBILEAPPCONTENTFILEUPLOADSTATE
        case "azureStorageUriRequestTimedOut":
            result = AZURESTORAGEURIREQUESTTIMEDOUT_MOBILEAPPCONTENTFILEUPLOADSTATE
        case "azureStorageUriRenewalSuccess":
            result = AZURESTORAGEURIRENEWALSUCCESS_MOBILEAPPCONTENTFILEUPLOADSTATE
        case "azureStorageUriRenewalPending":
            result = AZURESTORAGEURIRENEWALPENDING_MOBILEAPPCONTENTFILEUPLOADSTATE
        case "azureStorageUriRenewalFailed":
            result = AZURESTORAGEURIRENEWALFAILED_MOBILEAPPCONTENTFILEUPLOADSTATE
        case "azureStorageUriRenewalTimedOut":
            result = AZURESTORAGEURIRENEWALTIMEDOUT_MOBILEAPPCONTENTFILEUPLOADSTATE
        case "commitFileSuccess":
            result = COMMITFILESUCCESS_MOBILEAPPCONTENTFILEUPLOADSTATE
        case "commitFilePending":
            result = COMMITFILEPENDING_MOBILEAPPCONTENTFILEUPLOADSTATE
        case "commitFileFailed":
            result = COMMITFILEFAILED_MOBILEAPPCONTENTFILEUPLOADSTATE
        case "commitFileTimedOut":
            result = COMMITFILETIMEDOUT_MOBILEAPPCONTENTFILEUPLOADSTATE
        default:
            return nil, nil
    }
    return &result, nil
}
func SerializeMobileAppContentFileUploadState(values []MobileAppContentFileUploadState) []string {
    result := make([]string, len(values))
    for i, v := range values {
        result[i] = v.String()
    }
    return result
}
func (i MobileAppContentFileUploadState) isMultiValue() bool {
    return false
}
