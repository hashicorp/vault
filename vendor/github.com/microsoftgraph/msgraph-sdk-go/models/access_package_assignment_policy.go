package models

import (
    i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e "time"
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

type AccessPackageAssignmentPolicy struct {
    Entity
}
// NewAccessPackageAssignmentPolicy instantiates a new AccessPackageAssignmentPolicy and sets the default values.
func NewAccessPackageAssignmentPolicy()(*AccessPackageAssignmentPolicy) {
    m := &AccessPackageAssignmentPolicy{
        Entity: *NewEntity(),
    }
    return m
}
// CreateAccessPackageAssignmentPolicyFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
// returns a Parsable when successful
func CreateAccessPackageAssignmentPolicyFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewAccessPackageAssignmentPolicy(), nil
}
// GetAccessPackage gets the accessPackage property value. Access package containing this policy. Read-only.  Supports $expand.
// returns a AccessPackageable when successful
func (m *AccessPackageAssignmentPolicy) GetAccessPackage()(AccessPackageable) {
    val, err := m.GetBackingStore().Get("accessPackage")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(AccessPackageable)
    }
    return nil
}
// GetAllowedTargetScope gets the allowedTargetScope property value. Principals that can be assigned the access package through this policy. The possible values are: notSpecified, specificDirectoryUsers, specificConnectedOrganizationUsers, specificDirectoryServicePrincipals, allMemberUsers, allDirectoryUsers, allDirectoryServicePrincipals, allConfiguredConnectedOrganizationUsers, allExternalUsers, unknownFutureValue.
// returns a *AllowedTargetScope when successful
func (m *AccessPackageAssignmentPolicy) GetAllowedTargetScope()(*AllowedTargetScope) {
    val, err := m.GetBackingStore().Get("allowedTargetScope")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*AllowedTargetScope)
    }
    return nil
}
// GetAutomaticRequestSettings gets the automaticRequestSettings property value. This property is only present for an auto assignment policy; if absent, this is a request-based policy.
// returns a AccessPackageAutomaticRequestSettingsable when successful
func (m *AccessPackageAssignmentPolicy) GetAutomaticRequestSettings()(AccessPackageAutomaticRequestSettingsable) {
    val, err := m.GetBackingStore().Get("automaticRequestSettings")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(AccessPackageAutomaticRequestSettingsable)
    }
    return nil
}
// GetCatalog gets the catalog property value. Catalog of the access package containing this policy. Read-only.
// returns a AccessPackageCatalogable when successful
func (m *AccessPackageAssignmentPolicy) GetCatalog()(AccessPackageCatalogable) {
    val, err := m.GetBackingStore().Get("catalog")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(AccessPackageCatalogable)
    }
    return nil
}
// GetCreatedDateTime gets the createdDateTime property value. The Timestamp type represents date and time information using ISO 8601 format and is always in UTC time. For example, midnight UTC on Jan 1, 2014 is 2014-01-01T00:00:00Z.
// returns a *Time when successful
func (m *AccessPackageAssignmentPolicy) GetCreatedDateTime()(*i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time) {
    val, err := m.GetBackingStore().Get("createdDateTime")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)
    }
    return nil
}
// GetCustomExtensionStageSettings gets the customExtensionStageSettings property value. The collection of stages when to execute one or more custom access package workflow extensions. Supports $expand.
// returns a []CustomExtensionStageSettingable when successful
func (m *AccessPackageAssignmentPolicy) GetCustomExtensionStageSettings()([]CustomExtensionStageSettingable) {
    val, err := m.GetBackingStore().Get("customExtensionStageSettings")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]CustomExtensionStageSettingable)
    }
    return nil
}
// GetDescription gets the description property value. The description of the policy.
// returns a *string when successful
func (m *AccessPackageAssignmentPolicy) GetDescription()(*string) {
    val, err := m.GetBackingStore().Get("description")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetDisplayName gets the displayName property value. The display name of the policy.
// returns a *string when successful
func (m *AccessPackageAssignmentPolicy) GetDisplayName()(*string) {
    val, err := m.GetBackingStore().Get("displayName")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetExpiration gets the expiration property value. The expiration date for assignments created in this policy.
// returns a ExpirationPatternable when successful
func (m *AccessPackageAssignmentPolicy) GetExpiration()(ExpirationPatternable) {
    val, err := m.GetBackingStore().Get("expiration")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(ExpirationPatternable)
    }
    return nil
}
// GetFieldDeserializers the deserialization information for the current model
// returns a map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error) when successful
func (m *AccessPackageAssignmentPolicy) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := m.Entity.GetFieldDeserializers()
    res["accessPackage"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetObjectValue(CreateAccessPackageFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetAccessPackage(val.(AccessPackageable))
        }
        return nil
    }
    res["allowedTargetScope"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetEnumValue(ParseAllowedTargetScope)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetAllowedTargetScope(val.(*AllowedTargetScope))
        }
        return nil
    }
    res["automaticRequestSettings"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetObjectValue(CreateAccessPackageAutomaticRequestSettingsFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetAutomaticRequestSettings(val.(AccessPackageAutomaticRequestSettingsable))
        }
        return nil
    }
    res["catalog"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetObjectValue(CreateAccessPackageCatalogFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetCatalog(val.(AccessPackageCatalogable))
        }
        return nil
    }
    res["createdDateTime"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetTimeValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetCreatedDateTime(val)
        }
        return nil
    }
    res["customExtensionStageSettings"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetCollectionOfObjectValues(CreateCustomExtensionStageSettingFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            res := make([]CustomExtensionStageSettingable, len(val))
            for i, v := range val {
                if v != nil {
                    res[i] = v.(CustomExtensionStageSettingable)
                }
            }
            m.SetCustomExtensionStageSettings(res)
        }
        return nil
    }
    res["description"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetDescription(val)
        }
        return nil
    }
    res["displayName"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetDisplayName(val)
        }
        return nil
    }
    res["expiration"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetObjectValue(CreateExpirationPatternFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetExpiration(val.(ExpirationPatternable))
        }
        return nil
    }
    res["modifiedDateTime"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetTimeValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetModifiedDateTime(val)
        }
        return nil
    }
    res["questions"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetCollectionOfObjectValues(CreateAccessPackageQuestionFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            res := make([]AccessPackageQuestionable, len(val))
            for i, v := range val {
                if v != nil {
                    res[i] = v.(AccessPackageQuestionable)
                }
            }
            m.SetQuestions(res)
        }
        return nil
    }
    res["requestApprovalSettings"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetObjectValue(CreateAccessPackageAssignmentApprovalSettingsFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetRequestApprovalSettings(val.(AccessPackageAssignmentApprovalSettingsable))
        }
        return nil
    }
    res["requestorSettings"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetObjectValue(CreateAccessPackageAssignmentRequestorSettingsFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetRequestorSettings(val.(AccessPackageAssignmentRequestorSettingsable))
        }
        return nil
    }
    res["reviewSettings"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetObjectValue(CreateAccessPackageAssignmentReviewSettingsFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetReviewSettings(val.(AccessPackageAssignmentReviewSettingsable))
        }
        return nil
    }
    res["specificAllowedTargets"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetCollectionOfObjectValues(CreateSubjectSetFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            res := make([]SubjectSetable, len(val))
            for i, v := range val {
                if v != nil {
                    res[i] = v.(SubjectSetable)
                }
            }
            m.SetSpecificAllowedTargets(res)
        }
        return nil
    }
    return res
}
// GetModifiedDateTime gets the modifiedDateTime property value. The Timestamp type represents date and time information using ISO 8601 format and is always in UTC time. For example, midnight UTC on Jan 1, 2014 is 2014-01-01T00:00:00Z.
// returns a *Time when successful
func (m *AccessPackageAssignmentPolicy) GetModifiedDateTime()(*i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time) {
    val, err := m.GetBackingStore().Get("modifiedDateTime")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)
    }
    return nil
}
// GetQuestions gets the questions property value. Questions that are posed to the  requestor.
// returns a []AccessPackageQuestionable when successful
func (m *AccessPackageAssignmentPolicy) GetQuestions()([]AccessPackageQuestionable) {
    val, err := m.GetBackingStore().Get("questions")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]AccessPackageQuestionable)
    }
    return nil
}
// GetRequestApprovalSettings gets the requestApprovalSettings property value. Specifies the settings for approval of requests for an access package assignment through this policy. For example, if approval is required for new requests.
// returns a AccessPackageAssignmentApprovalSettingsable when successful
func (m *AccessPackageAssignmentPolicy) GetRequestApprovalSettings()(AccessPackageAssignmentApprovalSettingsable) {
    val, err := m.GetBackingStore().Get("requestApprovalSettings")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(AccessPackageAssignmentApprovalSettingsable)
    }
    return nil
}
// GetRequestorSettings gets the requestorSettings property value. Provides additional settings to select who can create a request for an access package assignment through this policy, and what they can include in their request.
// returns a AccessPackageAssignmentRequestorSettingsable when successful
func (m *AccessPackageAssignmentPolicy) GetRequestorSettings()(AccessPackageAssignmentRequestorSettingsable) {
    val, err := m.GetBackingStore().Get("requestorSettings")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(AccessPackageAssignmentRequestorSettingsable)
    }
    return nil
}
// GetReviewSettings gets the reviewSettings property value. Settings for access reviews of assignments through this policy.
// returns a AccessPackageAssignmentReviewSettingsable when successful
func (m *AccessPackageAssignmentPolicy) GetReviewSettings()(AccessPackageAssignmentReviewSettingsable) {
    val, err := m.GetBackingStore().Get("reviewSettings")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(AccessPackageAssignmentReviewSettingsable)
    }
    return nil
}
// GetSpecificAllowedTargets gets the specificAllowedTargets property value. The principals that can be assigned access from an access package through this policy.
// returns a []SubjectSetable when successful
func (m *AccessPackageAssignmentPolicy) GetSpecificAllowedTargets()([]SubjectSetable) {
    val, err := m.GetBackingStore().Get("specificAllowedTargets")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]SubjectSetable)
    }
    return nil
}
// Serialize serializes information the current object
func (m *AccessPackageAssignmentPolicy) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    err := m.Entity.Serialize(writer)
    if err != nil {
        return err
    }
    {
        err = writer.WriteObjectValue("accessPackage", m.GetAccessPackage())
        if err != nil {
            return err
        }
    }
    if m.GetAllowedTargetScope() != nil {
        cast := (*m.GetAllowedTargetScope()).String()
        err = writer.WriteStringValue("allowedTargetScope", &cast)
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteObjectValue("automaticRequestSettings", m.GetAutomaticRequestSettings())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteObjectValue("catalog", m.GetCatalog())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteTimeValue("createdDateTime", m.GetCreatedDateTime())
        if err != nil {
            return err
        }
    }
    if m.GetCustomExtensionStageSettings() != nil {
        cast := make([]i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, len(m.GetCustomExtensionStageSettings()))
        for i, v := range m.GetCustomExtensionStageSettings() {
            if v != nil {
                cast[i] = v.(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable)
            }
        }
        err = writer.WriteCollectionOfObjectValues("customExtensionStageSettings", cast)
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteStringValue("description", m.GetDescription())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteStringValue("displayName", m.GetDisplayName())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteObjectValue("expiration", m.GetExpiration())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteTimeValue("modifiedDateTime", m.GetModifiedDateTime())
        if err != nil {
            return err
        }
    }
    if m.GetQuestions() != nil {
        cast := make([]i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, len(m.GetQuestions()))
        for i, v := range m.GetQuestions() {
            if v != nil {
                cast[i] = v.(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable)
            }
        }
        err = writer.WriteCollectionOfObjectValues("questions", cast)
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteObjectValue("requestApprovalSettings", m.GetRequestApprovalSettings())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteObjectValue("requestorSettings", m.GetRequestorSettings())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteObjectValue("reviewSettings", m.GetReviewSettings())
        if err != nil {
            return err
        }
    }
    if m.GetSpecificAllowedTargets() != nil {
        cast := make([]i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, len(m.GetSpecificAllowedTargets()))
        for i, v := range m.GetSpecificAllowedTargets() {
            if v != nil {
                cast[i] = v.(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable)
            }
        }
        err = writer.WriteCollectionOfObjectValues("specificAllowedTargets", cast)
        if err != nil {
            return err
        }
    }
    return nil
}
// SetAccessPackage sets the accessPackage property value. Access package containing this policy. Read-only.  Supports $expand.
func (m *AccessPackageAssignmentPolicy) SetAccessPackage(value AccessPackageable)() {
    err := m.GetBackingStore().Set("accessPackage", value)
    if err != nil {
        panic(err)
    }
}
// SetAllowedTargetScope sets the allowedTargetScope property value. Principals that can be assigned the access package through this policy. The possible values are: notSpecified, specificDirectoryUsers, specificConnectedOrganizationUsers, specificDirectoryServicePrincipals, allMemberUsers, allDirectoryUsers, allDirectoryServicePrincipals, allConfiguredConnectedOrganizationUsers, allExternalUsers, unknownFutureValue.
func (m *AccessPackageAssignmentPolicy) SetAllowedTargetScope(value *AllowedTargetScope)() {
    err := m.GetBackingStore().Set("allowedTargetScope", value)
    if err != nil {
        panic(err)
    }
}
// SetAutomaticRequestSettings sets the automaticRequestSettings property value. This property is only present for an auto assignment policy; if absent, this is a request-based policy.
func (m *AccessPackageAssignmentPolicy) SetAutomaticRequestSettings(value AccessPackageAutomaticRequestSettingsable)() {
    err := m.GetBackingStore().Set("automaticRequestSettings", value)
    if err != nil {
        panic(err)
    }
}
// SetCatalog sets the catalog property value. Catalog of the access package containing this policy. Read-only.
func (m *AccessPackageAssignmentPolicy) SetCatalog(value AccessPackageCatalogable)() {
    err := m.GetBackingStore().Set("catalog", value)
    if err != nil {
        panic(err)
    }
}
// SetCreatedDateTime sets the createdDateTime property value. The Timestamp type represents date and time information using ISO 8601 format and is always in UTC time. For example, midnight UTC on Jan 1, 2014 is 2014-01-01T00:00:00Z.
func (m *AccessPackageAssignmentPolicy) SetCreatedDateTime(value *i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)() {
    err := m.GetBackingStore().Set("createdDateTime", value)
    if err != nil {
        panic(err)
    }
}
// SetCustomExtensionStageSettings sets the customExtensionStageSettings property value. The collection of stages when to execute one or more custom access package workflow extensions. Supports $expand.
func (m *AccessPackageAssignmentPolicy) SetCustomExtensionStageSettings(value []CustomExtensionStageSettingable)() {
    err := m.GetBackingStore().Set("customExtensionStageSettings", value)
    if err != nil {
        panic(err)
    }
}
// SetDescription sets the description property value. The description of the policy.
func (m *AccessPackageAssignmentPolicy) SetDescription(value *string)() {
    err := m.GetBackingStore().Set("description", value)
    if err != nil {
        panic(err)
    }
}
// SetDisplayName sets the displayName property value. The display name of the policy.
func (m *AccessPackageAssignmentPolicy) SetDisplayName(value *string)() {
    err := m.GetBackingStore().Set("displayName", value)
    if err != nil {
        panic(err)
    }
}
// SetExpiration sets the expiration property value. The expiration date for assignments created in this policy.
func (m *AccessPackageAssignmentPolicy) SetExpiration(value ExpirationPatternable)() {
    err := m.GetBackingStore().Set("expiration", value)
    if err != nil {
        panic(err)
    }
}
// SetModifiedDateTime sets the modifiedDateTime property value. The Timestamp type represents date and time information using ISO 8601 format and is always in UTC time. For example, midnight UTC on Jan 1, 2014 is 2014-01-01T00:00:00Z.
func (m *AccessPackageAssignmentPolicy) SetModifiedDateTime(value *i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)() {
    err := m.GetBackingStore().Set("modifiedDateTime", value)
    if err != nil {
        panic(err)
    }
}
// SetQuestions sets the questions property value. Questions that are posed to the  requestor.
func (m *AccessPackageAssignmentPolicy) SetQuestions(value []AccessPackageQuestionable)() {
    err := m.GetBackingStore().Set("questions", value)
    if err != nil {
        panic(err)
    }
}
// SetRequestApprovalSettings sets the requestApprovalSettings property value. Specifies the settings for approval of requests for an access package assignment through this policy. For example, if approval is required for new requests.
func (m *AccessPackageAssignmentPolicy) SetRequestApprovalSettings(value AccessPackageAssignmentApprovalSettingsable)() {
    err := m.GetBackingStore().Set("requestApprovalSettings", value)
    if err != nil {
        panic(err)
    }
}
// SetRequestorSettings sets the requestorSettings property value. Provides additional settings to select who can create a request for an access package assignment through this policy, and what they can include in their request.
func (m *AccessPackageAssignmentPolicy) SetRequestorSettings(value AccessPackageAssignmentRequestorSettingsable)() {
    err := m.GetBackingStore().Set("requestorSettings", value)
    if err != nil {
        panic(err)
    }
}
// SetReviewSettings sets the reviewSettings property value. Settings for access reviews of assignments through this policy.
func (m *AccessPackageAssignmentPolicy) SetReviewSettings(value AccessPackageAssignmentReviewSettingsable)() {
    err := m.GetBackingStore().Set("reviewSettings", value)
    if err != nil {
        panic(err)
    }
}
// SetSpecificAllowedTargets sets the specificAllowedTargets property value. The principals that can be assigned access from an access package through this policy.
func (m *AccessPackageAssignmentPolicy) SetSpecificAllowedTargets(value []SubjectSetable)() {
    err := m.GetBackingStore().Set("specificAllowedTargets", value)
    if err != nil {
        panic(err)
    }
}
type AccessPackageAssignmentPolicyable interface {
    Entityable
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable
    GetAccessPackage()(AccessPackageable)
    GetAllowedTargetScope()(*AllowedTargetScope)
    GetAutomaticRequestSettings()(AccessPackageAutomaticRequestSettingsable)
    GetCatalog()(AccessPackageCatalogable)
    GetCreatedDateTime()(*i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)
    GetCustomExtensionStageSettings()([]CustomExtensionStageSettingable)
    GetDescription()(*string)
    GetDisplayName()(*string)
    GetExpiration()(ExpirationPatternable)
    GetModifiedDateTime()(*i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)
    GetQuestions()([]AccessPackageQuestionable)
    GetRequestApprovalSettings()(AccessPackageAssignmentApprovalSettingsable)
    GetRequestorSettings()(AccessPackageAssignmentRequestorSettingsable)
    GetReviewSettings()(AccessPackageAssignmentReviewSettingsable)
    GetSpecificAllowedTargets()([]SubjectSetable)
    SetAccessPackage(value AccessPackageable)()
    SetAllowedTargetScope(value *AllowedTargetScope)()
    SetAutomaticRequestSettings(value AccessPackageAutomaticRequestSettingsable)()
    SetCatalog(value AccessPackageCatalogable)()
    SetCreatedDateTime(value *i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)()
    SetCustomExtensionStageSettings(value []CustomExtensionStageSettingable)()
    SetDescription(value *string)()
    SetDisplayName(value *string)()
    SetExpiration(value ExpirationPatternable)()
    SetModifiedDateTime(value *i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)()
    SetQuestions(value []AccessPackageQuestionable)()
    SetRequestApprovalSettings(value AccessPackageAssignmentApprovalSettingsable)()
    SetRequestorSettings(value AccessPackageAssignmentRequestorSettingsable)()
    SetReviewSettings(value AccessPackageAssignmentReviewSettingsable)()
    SetSpecificAllowedTargets(value []SubjectSetable)()
}
