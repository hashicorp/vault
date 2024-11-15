package models

import (
    i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e "time"
    i561e97a8befe7661a44c8f54600992b4207a3a0cf6770e5559949bc276de2e22 "github.com/google/uuid"
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

type Application struct {
    DirectoryObject
}
// NewApplication instantiates a new Application and sets the default values.
func NewApplication()(*Application) {
    m := &Application{
        DirectoryObject: *NewDirectoryObject(),
    }
    odataTypeValue := "#microsoft.graph.application"
    m.SetOdataType(&odataTypeValue)
    return m
}
// CreateApplicationFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
// returns a Parsable when successful
func CreateApplicationFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewApplication(), nil
}
// GetAddIns gets the addIns property value. Defines custom behavior that a consuming service can use to call an app in specific contexts. For example, applications that can render file streams can set the addIns property for its 'FileHandler' functionality. This lets services like Microsoft 365 call the application in the context of a document the user is working on.
// returns a []AddInable when successful
func (m *Application) GetAddIns()([]AddInable) {
    val, err := m.GetBackingStore().Get("addIns")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]AddInable)
    }
    return nil
}
// GetApi gets the api property value. Specifies settings for an application that implements a web API.
// returns a ApiApplicationable when successful
func (m *Application) GetApi()(ApiApplicationable) {
    val, err := m.GetBackingStore().Get("api")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(ApiApplicationable)
    }
    return nil
}
// GetAppId gets the appId property value. The unique identifier for the application that is assigned to an application by Microsoft Entra ID. Not nullable. Read-only. Alternate key. Supports $filter (eq).
// returns a *string when successful
func (m *Application) GetAppId()(*string) {
    val, err := m.GetBackingStore().Get("appId")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetApplicationTemplateId gets the applicationTemplateId property value. Unique identifier of the applicationTemplate. Supports $filter (eq, not, ne). Read-only. null if the app wasn't created from an application template.
// returns a *string when successful
func (m *Application) GetApplicationTemplateId()(*string) {
    val, err := m.GetBackingStore().Get("applicationTemplateId")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetAppManagementPolicies gets the appManagementPolicies property value. The appManagementPolicy applied to this application.
// returns a []AppManagementPolicyable when successful
func (m *Application) GetAppManagementPolicies()([]AppManagementPolicyable) {
    val, err := m.GetBackingStore().Get("appManagementPolicies")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]AppManagementPolicyable)
    }
    return nil
}
// GetAppRoles gets the appRoles property value. The collection of roles defined for the application. With app role assignments, these roles can be assigned to users, groups, or service principals associated with other applications. Not nullable.
// returns a []AppRoleable when successful
func (m *Application) GetAppRoles()([]AppRoleable) {
    val, err := m.GetBackingStore().Get("appRoles")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]AppRoleable)
    }
    return nil
}
// GetCertification gets the certification property value. Specifies the certification status of the application.
// returns a Certificationable when successful
func (m *Application) GetCertification()(Certificationable) {
    val, err := m.GetBackingStore().Get("certification")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(Certificationable)
    }
    return nil
}
// GetCreatedDateTime gets the createdDateTime property value. The date and time the application was registered. The DateTimeOffset type represents date and time information using ISO 8601 format and is always in UTC time. For example, midnight UTC on Jan 1, 2014 is 2014-01-01T00:00:00Z. Read-only.  Supports $filter (eq, ne, not, ge, le, in, and eq on null values) and $orderby.
// returns a *Time when successful
func (m *Application) GetCreatedDateTime()(*i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time) {
    val, err := m.GetBackingStore().Get("createdDateTime")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)
    }
    return nil
}
// GetCreatedOnBehalfOf gets the createdOnBehalfOf property value. Supports $filter (/$count eq 0, /$count ne 0). Read-only.
// returns a DirectoryObjectable when successful
func (m *Application) GetCreatedOnBehalfOf()(DirectoryObjectable) {
    val, err := m.GetBackingStore().Get("createdOnBehalfOf")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(DirectoryObjectable)
    }
    return nil
}
// GetDefaultRedirectUri gets the defaultRedirectUri property value. The defaultRedirectUri property
// returns a *string when successful
func (m *Application) GetDefaultRedirectUri()(*string) {
    val, err := m.GetBackingStore().Get("defaultRedirectUri")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetDescription gets the description property value. Free text field to provide a description of the application object to end users. The maximum allowed size is 1,024 characters. Supports $filter (eq, ne, not, ge, le, startsWith) and $search.
// returns a *string when successful
func (m *Application) GetDescription()(*string) {
    val, err := m.GetBackingStore().Get("description")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetDisabledByMicrosoftStatus gets the disabledByMicrosoftStatus property value. Specifies whether Microsoft has disabled the registered application. Possible values are: null (default value), NotDisabled, and DisabledDueToViolationOfServicesAgreement (reasons include suspicious, abusive, or malicious activity, or a violation of the Microsoft Services Agreement).  Supports $filter (eq, ne, not).
// returns a *string when successful
func (m *Application) GetDisabledByMicrosoftStatus()(*string) {
    val, err := m.GetBackingStore().Get("disabledByMicrosoftStatus")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetDisplayName gets the displayName property value. The display name for the application. Supports $filter (eq, ne, not, ge, le, in, startsWith, and eq on null values), $search, and $orderby.
// returns a *string when successful
func (m *Application) GetDisplayName()(*string) {
    val, err := m.GetBackingStore().Get("displayName")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetExtensionProperties gets the extensionProperties property value. Read-only. Nullable. Supports $expand and $filter (/$count eq 0, /$count ne 0).
// returns a []ExtensionPropertyable when successful
func (m *Application) GetExtensionProperties()([]ExtensionPropertyable) {
    val, err := m.GetBackingStore().Get("extensionProperties")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]ExtensionPropertyable)
    }
    return nil
}
// GetFederatedIdentityCredentials gets the federatedIdentityCredentials property value. Federated identities for applications. Supports $expand and $filter (startsWith, /$count eq 0, /$count ne 0).
// returns a []FederatedIdentityCredentialable when successful
func (m *Application) GetFederatedIdentityCredentials()([]FederatedIdentityCredentialable) {
    val, err := m.GetBackingStore().Get("federatedIdentityCredentials")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]FederatedIdentityCredentialable)
    }
    return nil
}
// GetFieldDeserializers the deserialization information for the current model
// returns a map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error) when successful
func (m *Application) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := m.DirectoryObject.GetFieldDeserializers()
    res["addIns"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetCollectionOfObjectValues(CreateAddInFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            res := make([]AddInable, len(val))
            for i, v := range val {
                if v != nil {
                    res[i] = v.(AddInable)
                }
            }
            m.SetAddIns(res)
        }
        return nil
    }
    res["api"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetObjectValue(CreateApiApplicationFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetApi(val.(ApiApplicationable))
        }
        return nil
    }
    res["appId"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetAppId(val)
        }
        return nil
    }
    res["applicationTemplateId"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetApplicationTemplateId(val)
        }
        return nil
    }
    res["appManagementPolicies"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetCollectionOfObjectValues(CreateAppManagementPolicyFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            res := make([]AppManagementPolicyable, len(val))
            for i, v := range val {
                if v != nil {
                    res[i] = v.(AppManagementPolicyable)
                }
            }
            m.SetAppManagementPolicies(res)
        }
        return nil
    }
    res["appRoles"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetCollectionOfObjectValues(CreateAppRoleFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            res := make([]AppRoleable, len(val))
            for i, v := range val {
                if v != nil {
                    res[i] = v.(AppRoleable)
                }
            }
            m.SetAppRoles(res)
        }
        return nil
    }
    res["certification"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetObjectValue(CreateCertificationFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetCertification(val.(Certificationable))
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
    res["createdOnBehalfOf"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetObjectValue(CreateDirectoryObjectFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetCreatedOnBehalfOf(val.(DirectoryObjectable))
        }
        return nil
    }
    res["defaultRedirectUri"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetDefaultRedirectUri(val)
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
    res["disabledByMicrosoftStatus"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetDisabledByMicrosoftStatus(val)
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
    res["extensionProperties"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetCollectionOfObjectValues(CreateExtensionPropertyFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            res := make([]ExtensionPropertyable, len(val))
            for i, v := range val {
                if v != nil {
                    res[i] = v.(ExtensionPropertyable)
                }
            }
            m.SetExtensionProperties(res)
        }
        return nil
    }
    res["federatedIdentityCredentials"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetCollectionOfObjectValues(CreateFederatedIdentityCredentialFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            res := make([]FederatedIdentityCredentialable, len(val))
            for i, v := range val {
                if v != nil {
                    res[i] = v.(FederatedIdentityCredentialable)
                }
            }
            m.SetFederatedIdentityCredentials(res)
        }
        return nil
    }
    res["groupMembershipClaims"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetGroupMembershipClaims(val)
        }
        return nil
    }
    res["homeRealmDiscoveryPolicies"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetCollectionOfObjectValues(CreateHomeRealmDiscoveryPolicyFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            res := make([]HomeRealmDiscoveryPolicyable, len(val))
            for i, v := range val {
                if v != nil {
                    res[i] = v.(HomeRealmDiscoveryPolicyable)
                }
            }
            m.SetHomeRealmDiscoveryPolicies(res)
        }
        return nil
    }
    res["identifierUris"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetCollectionOfPrimitiveValues("string")
        if err != nil {
            return err
        }
        if val != nil {
            res := make([]string, len(val))
            for i, v := range val {
                if v != nil {
                    res[i] = *(v.(*string))
                }
            }
            m.SetIdentifierUris(res)
        }
        return nil
    }
    res["info"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetObjectValue(CreateInformationalUrlFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetInfo(val.(InformationalUrlable))
        }
        return nil
    }
    res["isDeviceOnlyAuthSupported"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetBoolValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetIsDeviceOnlyAuthSupported(val)
        }
        return nil
    }
    res["isFallbackPublicClient"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetBoolValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetIsFallbackPublicClient(val)
        }
        return nil
    }
    res["keyCredentials"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetCollectionOfObjectValues(CreateKeyCredentialFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            res := make([]KeyCredentialable, len(val))
            for i, v := range val {
                if v != nil {
                    res[i] = v.(KeyCredentialable)
                }
            }
            m.SetKeyCredentials(res)
        }
        return nil
    }
    res["logo"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetByteArrayValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetLogo(val)
        }
        return nil
    }
    res["nativeAuthenticationApisEnabled"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetEnumValue(ParseNativeAuthenticationApisEnabled)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetNativeAuthenticationApisEnabled(val.(*NativeAuthenticationApisEnabled))
        }
        return nil
    }
    res["notes"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetNotes(val)
        }
        return nil
    }
    res["oauth2RequirePostResponse"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetBoolValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetOauth2RequirePostResponse(val)
        }
        return nil
    }
    res["optionalClaims"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetObjectValue(CreateOptionalClaimsFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetOptionalClaims(val.(OptionalClaimsable))
        }
        return nil
    }
    res["owners"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetCollectionOfObjectValues(CreateDirectoryObjectFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            res := make([]DirectoryObjectable, len(val))
            for i, v := range val {
                if v != nil {
                    res[i] = v.(DirectoryObjectable)
                }
            }
            m.SetOwners(res)
        }
        return nil
    }
    res["parentalControlSettings"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetObjectValue(CreateParentalControlSettingsFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetParentalControlSettings(val.(ParentalControlSettingsable))
        }
        return nil
    }
    res["passwordCredentials"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetCollectionOfObjectValues(CreatePasswordCredentialFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            res := make([]PasswordCredentialable, len(val))
            for i, v := range val {
                if v != nil {
                    res[i] = v.(PasswordCredentialable)
                }
            }
            m.SetPasswordCredentials(res)
        }
        return nil
    }
    res["publicClient"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetObjectValue(CreatePublicClientApplicationFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetPublicClient(val.(PublicClientApplicationable))
        }
        return nil
    }
    res["publisherDomain"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetPublisherDomain(val)
        }
        return nil
    }
    res["requestSignatureVerification"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetObjectValue(CreateRequestSignatureVerificationFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetRequestSignatureVerification(val.(RequestSignatureVerificationable))
        }
        return nil
    }
    res["requiredResourceAccess"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetCollectionOfObjectValues(CreateRequiredResourceAccessFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            res := make([]RequiredResourceAccessable, len(val))
            for i, v := range val {
                if v != nil {
                    res[i] = v.(RequiredResourceAccessable)
                }
            }
            m.SetRequiredResourceAccess(res)
        }
        return nil
    }
    res["samlMetadataUrl"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetSamlMetadataUrl(val)
        }
        return nil
    }
    res["serviceManagementReference"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetServiceManagementReference(val)
        }
        return nil
    }
    res["servicePrincipalLockConfiguration"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetObjectValue(CreateServicePrincipalLockConfigurationFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetServicePrincipalLockConfiguration(val.(ServicePrincipalLockConfigurationable))
        }
        return nil
    }
    res["signInAudience"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetSignInAudience(val)
        }
        return nil
    }
    res["spa"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetObjectValue(CreateSpaApplicationFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetSpa(val.(SpaApplicationable))
        }
        return nil
    }
    res["synchronization"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetObjectValue(CreateSynchronizationFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetSynchronization(val.(Synchronizationable))
        }
        return nil
    }
    res["tags"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetCollectionOfPrimitiveValues("string")
        if err != nil {
            return err
        }
        if val != nil {
            res := make([]string, len(val))
            for i, v := range val {
                if v != nil {
                    res[i] = *(v.(*string))
                }
            }
            m.SetTags(res)
        }
        return nil
    }
    res["tokenEncryptionKeyId"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetUUIDValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetTokenEncryptionKeyId(val)
        }
        return nil
    }
    res["tokenIssuancePolicies"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetCollectionOfObjectValues(CreateTokenIssuancePolicyFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            res := make([]TokenIssuancePolicyable, len(val))
            for i, v := range val {
                if v != nil {
                    res[i] = v.(TokenIssuancePolicyable)
                }
            }
            m.SetTokenIssuancePolicies(res)
        }
        return nil
    }
    res["tokenLifetimePolicies"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetCollectionOfObjectValues(CreateTokenLifetimePolicyFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            res := make([]TokenLifetimePolicyable, len(val))
            for i, v := range val {
                if v != nil {
                    res[i] = v.(TokenLifetimePolicyable)
                }
            }
            m.SetTokenLifetimePolicies(res)
        }
        return nil
    }
    res["uniqueName"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetUniqueName(val)
        }
        return nil
    }
    res["verifiedPublisher"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetObjectValue(CreateVerifiedPublisherFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetVerifiedPublisher(val.(VerifiedPublisherable))
        }
        return nil
    }
    res["web"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetObjectValue(CreateWebApplicationFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetWeb(val.(WebApplicationable))
        }
        return nil
    }
    return res
}
// GetGroupMembershipClaims gets the groupMembershipClaims property value. Configures the groups claim issued in a user or OAuth 2.0 access token that the application expects. To set this attribute, use one of the following valid string values: None, SecurityGroup (for security groups and Microsoft Entra roles), All (this gets all of the security groups, distribution groups, and Microsoft Entra directory roles that the signed-in user is a member of).
// returns a *string when successful
func (m *Application) GetGroupMembershipClaims()(*string) {
    val, err := m.GetBackingStore().Get("groupMembershipClaims")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetHomeRealmDiscoveryPolicies gets the homeRealmDiscoveryPolicies property value. The homeRealmDiscoveryPolicies property
// returns a []HomeRealmDiscoveryPolicyable when successful
func (m *Application) GetHomeRealmDiscoveryPolicies()([]HomeRealmDiscoveryPolicyable) {
    val, err := m.GetBackingStore().Get("homeRealmDiscoveryPolicies")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]HomeRealmDiscoveryPolicyable)
    }
    return nil
}
// GetIdentifierUris gets the identifierUris property value. Also known as App ID URI, this value is set when an application is used as a resource app. The identifierUris acts as the prefix for the scopes you reference in your API's code, and it must be globally unique. You can use the default value provided, which is in the form api://<appId>, or specify a more readable URI like https://contoso.com/api. For more information on valid identifierUris patterns and best practices, see Microsoft Entra application registration security best practices. Not nullable. Supports $filter (eq, ne, ge, le, startsWith).
// returns a []string when successful
func (m *Application) GetIdentifierUris()([]string) {
    val, err := m.GetBackingStore().Get("identifierUris")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]string)
    }
    return nil
}
// GetInfo gets the info property value. Basic profile information of the application such as  app's marketing, support, terms of service and privacy statement URLs. The terms of service and privacy statement are surfaced to users through the user consent experience. For more info, see How to: Add Terms of service and privacy statement for registered Microsoft Entra apps. Supports $filter (eq, ne, not, ge, le, and eq on null values).
// returns a InformationalUrlable when successful
func (m *Application) GetInfo()(InformationalUrlable) {
    val, err := m.GetBackingStore().Get("info")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(InformationalUrlable)
    }
    return nil
}
// GetIsDeviceOnlyAuthSupported gets the isDeviceOnlyAuthSupported property value. Specifies whether this application supports device authentication without a user. The default is false.
// returns a *bool when successful
func (m *Application) GetIsDeviceOnlyAuthSupported()(*bool) {
    val, err := m.GetBackingStore().Get("isDeviceOnlyAuthSupported")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*bool)
    }
    return nil
}
// GetIsFallbackPublicClient gets the isFallbackPublicClient property value. Specifies the fallback application type as public client, such as an installed application running on a mobile device. The default value is false, which means the fallback application type is confidential client such as a web app. There are certain scenarios where Microsoft Entra ID can't determine the client application type. For example, the ROPC flow where it's configured without specifying a redirect URI. In those cases, Microsoft Entra ID interprets the application type based on the value of this property.
// returns a *bool when successful
func (m *Application) GetIsFallbackPublicClient()(*bool) {
    val, err := m.GetBackingStore().Get("isFallbackPublicClient")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*bool)
    }
    return nil
}
// GetKeyCredentials gets the keyCredentials property value. The collection of key credentials associated with the application. Not nullable. Supports $filter (eq, not, ge, le).
// returns a []KeyCredentialable when successful
func (m *Application) GetKeyCredentials()([]KeyCredentialable) {
    val, err := m.GetBackingStore().Get("keyCredentials")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]KeyCredentialable)
    }
    return nil
}
// GetLogo gets the logo property value. The main logo for the application. Not nullable.
// returns a []byte when successful
func (m *Application) GetLogo()([]byte) {
    val, err := m.GetBackingStore().Get("logo")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]byte)
    }
    return nil
}
// GetNativeAuthenticationApisEnabled gets the nativeAuthenticationApisEnabled property value. Specifies whether the Native Authentication APIs are enabled for the application. The possible values are: none and all. Default is none. For more information, see Native Authentication.
// returns a *NativeAuthenticationApisEnabled when successful
func (m *Application) GetNativeAuthenticationApisEnabled()(*NativeAuthenticationApisEnabled) {
    val, err := m.GetBackingStore().Get("nativeAuthenticationApisEnabled")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*NativeAuthenticationApisEnabled)
    }
    return nil
}
// GetNotes gets the notes property value. Notes relevant for the management of the application.
// returns a *string when successful
func (m *Application) GetNotes()(*string) {
    val, err := m.GetBackingStore().Get("notes")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetOauth2RequirePostResponse gets the oauth2RequirePostResponse property value. The oauth2RequirePostResponse property
// returns a *bool when successful
func (m *Application) GetOauth2RequirePostResponse()(*bool) {
    val, err := m.GetBackingStore().Get("oauth2RequirePostResponse")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*bool)
    }
    return nil
}
// GetOptionalClaims gets the optionalClaims property value. Application developers can configure optional claims in their Microsoft Entra applications to specify the claims that are sent to their application by the Microsoft security token service. For more information, see How to: Provide optional claims to your app.
// returns a OptionalClaimsable when successful
func (m *Application) GetOptionalClaims()(OptionalClaimsable) {
    val, err := m.GetBackingStore().Get("optionalClaims")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(OptionalClaimsable)
    }
    return nil
}
// GetOwners gets the owners property value. Directory objects that are owners of the application. Read-only. Nullable. Supports $expand, $filter (/$count eq 0, /$count ne 0, /$count eq 1, /$count ne 1), and $select nested in $expand.
// returns a []DirectoryObjectable when successful
func (m *Application) GetOwners()([]DirectoryObjectable) {
    val, err := m.GetBackingStore().Get("owners")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]DirectoryObjectable)
    }
    return nil
}
// GetParentalControlSettings gets the parentalControlSettings property value. Specifies parental control settings for an application.
// returns a ParentalControlSettingsable when successful
func (m *Application) GetParentalControlSettings()(ParentalControlSettingsable) {
    val, err := m.GetBackingStore().Get("parentalControlSettings")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(ParentalControlSettingsable)
    }
    return nil
}
// GetPasswordCredentials gets the passwordCredentials property value. The collection of password credentials associated with the application. Not nullable.
// returns a []PasswordCredentialable when successful
func (m *Application) GetPasswordCredentials()([]PasswordCredentialable) {
    val, err := m.GetBackingStore().Get("passwordCredentials")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]PasswordCredentialable)
    }
    return nil
}
// GetPublicClient gets the publicClient property value. Specifies settings for installed clients such as desktop or mobile devices.
// returns a PublicClientApplicationable when successful
func (m *Application) GetPublicClient()(PublicClientApplicationable) {
    val, err := m.GetBackingStore().Get("publicClient")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(PublicClientApplicationable)
    }
    return nil
}
// GetPublisherDomain gets the publisherDomain property value. The verified publisher domain for the application. Read-only. For more information, see How to: Configure an application's publisher domain. Supports $filter (eq, ne, ge, le, startsWith).
// returns a *string when successful
func (m *Application) GetPublisherDomain()(*string) {
    val, err := m.GetBackingStore().Get("publisherDomain")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetRequestSignatureVerification gets the requestSignatureVerification property value. Specifies whether this application requires Microsoft Entra ID to verify the signed authentication requests.
// returns a RequestSignatureVerificationable when successful
func (m *Application) GetRequestSignatureVerification()(RequestSignatureVerificationable) {
    val, err := m.GetBackingStore().Get("requestSignatureVerification")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(RequestSignatureVerificationable)
    }
    return nil
}
// GetRequiredResourceAccess gets the requiredResourceAccess property value. Specifies the resources that the application needs to access. This property also specifies the set of delegated permissions and application roles that it needs for each of those resources. This configuration of access to the required resources drives the consent experience. No more than 50 resource services (APIs) can be configured. Beginning mid-October 2021, the total number of required permissions must not exceed 400. For more information, see Limits on requested permissions per app. Not nullable. Supports $filter (eq, not, ge, le).
// returns a []RequiredResourceAccessable when successful
func (m *Application) GetRequiredResourceAccess()([]RequiredResourceAccessable) {
    val, err := m.GetBackingStore().Get("requiredResourceAccess")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]RequiredResourceAccessable)
    }
    return nil
}
// GetSamlMetadataUrl gets the samlMetadataUrl property value. The URL where the service exposes SAML metadata for federation. This property is valid only for single-tenant applications. Nullable.
// returns a *string when successful
func (m *Application) GetSamlMetadataUrl()(*string) {
    val, err := m.GetBackingStore().Get("samlMetadataUrl")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetServiceManagementReference gets the serviceManagementReference property value. References application or service contact information from a Service or Asset Management database. Nullable.
// returns a *string when successful
func (m *Application) GetServiceManagementReference()(*string) {
    val, err := m.GetBackingStore().Get("serviceManagementReference")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetServicePrincipalLockConfiguration gets the servicePrincipalLockConfiguration property value. Specifies whether sensitive properties of a multitenant application should be locked for editing after the application is provisioned in a tenant. Nullable. null by default.
// returns a ServicePrincipalLockConfigurationable when successful
func (m *Application) GetServicePrincipalLockConfiguration()(ServicePrincipalLockConfigurationable) {
    val, err := m.GetBackingStore().Get("servicePrincipalLockConfiguration")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(ServicePrincipalLockConfigurationable)
    }
    return nil
}
// GetSignInAudience gets the signInAudience property value. Specifies the Microsoft accounts that are supported for the current application. The possible values are: AzureADMyOrg (default), AzureADMultipleOrgs, AzureADandPersonalMicrosoftAccount, and PersonalMicrosoftAccount. See more in the table. The value of this object also limits the number of permissions an app can request. For more information, see Limits on requested permissions per app. The value for this property has implications on other app object properties. As a result, if you change this property, you might need to change other properties first. For more information, see Validation differences for signInAudience.Supports $filter (eq, ne, not).
// returns a *string when successful
func (m *Application) GetSignInAudience()(*string) {
    val, err := m.GetBackingStore().Get("signInAudience")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetSpa gets the spa property value. Specifies settings for a single-page application, including sign out URLs and redirect URIs for authorization codes and access tokens.
// returns a SpaApplicationable when successful
func (m *Application) GetSpa()(SpaApplicationable) {
    val, err := m.GetBackingStore().Get("spa")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(SpaApplicationable)
    }
    return nil
}
// GetSynchronization gets the synchronization property value. Represents the capability for Microsoft Entra identity synchronization through the Microsoft Graph API.
// returns a Synchronizationable when successful
func (m *Application) GetSynchronization()(Synchronizationable) {
    val, err := m.GetBackingStore().Get("synchronization")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(Synchronizationable)
    }
    return nil
}
// GetTags gets the tags property value. Custom strings that can be used to categorize and identify the application. Not nullable. Strings added here will also appear in the tags property of any associated service principals.Supports $filter (eq, not, ge, le, startsWith) and $search.
// returns a []string when successful
func (m *Application) GetTags()([]string) {
    val, err := m.GetBackingStore().Get("tags")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]string)
    }
    return nil
}
// GetTokenEncryptionKeyId gets the tokenEncryptionKeyId property value. Specifies the keyId of a public key from the keyCredentials collection. When configured, Microsoft Entra ID encrypts all the tokens it emits by using the key this property points to. The application code that receives the encrypted token must use the matching private key to decrypt the token before it can be used for the signed-in user.
// returns a *UUID when successful
func (m *Application) GetTokenEncryptionKeyId()(*i561e97a8befe7661a44c8f54600992b4207a3a0cf6770e5559949bc276de2e22.UUID) {
    val, err := m.GetBackingStore().Get("tokenEncryptionKeyId")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*i561e97a8befe7661a44c8f54600992b4207a3a0cf6770e5559949bc276de2e22.UUID)
    }
    return nil
}
// GetTokenIssuancePolicies gets the tokenIssuancePolicies property value. The tokenIssuancePolicies property
// returns a []TokenIssuancePolicyable when successful
func (m *Application) GetTokenIssuancePolicies()([]TokenIssuancePolicyable) {
    val, err := m.GetBackingStore().Get("tokenIssuancePolicies")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]TokenIssuancePolicyable)
    }
    return nil
}
// GetTokenLifetimePolicies gets the tokenLifetimePolicies property value. The tokenLifetimePolicies property
// returns a []TokenLifetimePolicyable when successful
func (m *Application) GetTokenLifetimePolicies()([]TokenLifetimePolicyable) {
    val, err := m.GetBackingStore().Get("tokenLifetimePolicies")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]TokenLifetimePolicyable)
    }
    return nil
}
// GetUniqueName gets the uniqueName property value. The unique identifier that can be assigned to an application and used as an alternate key. Immutable. Read-only.
// returns a *string when successful
func (m *Application) GetUniqueName()(*string) {
    val, err := m.GetBackingStore().Get("uniqueName")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetVerifiedPublisher gets the verifiedPublisher property value. Specifies the verified publisher of the application. For more information about how publisher verification helps support application security, trustworthiness, and compliance, see Publisher verification.
// returns a VerifiedPublisherable when successful
func (m *Application) GetVerifiedPublisher()(VerifiedPublisherable) {
    val, err := m.GetBackingStore().Get("verifiedPublisher")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(VerifiedPublisherable)
    }
    return nil
}
// GetWeb gets the web property value. Specifies settings for a web application.
// returns a WebApplicationable when successful
func (m *Application) GetWeb()(WebApplicationable) {
    val, err := m.GetBackingStore().Get("web")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(WebApplicationable)
    }
    return nil
}
// Serialize serializes information the current object
func (m *Application) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    err := m.DirectoryObject.Serialize(writer)
    if err != nil {
        return err
    }
    if m.GetAddIns() != nil {
        cast := make([]i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, len(m.GetAddIns()))
        for i, v := range m.GetAddIns() {
            if v != nil {
                cast[i] = v.(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable)
            }
        }
        err = writer.WriteCollectionOfObjectValues("addIns", cast)
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteObjectValue("api", m.GetApi())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteStringValue("appId", m.GetAppId())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteStringValue("applicationTemplateId", m.GetApplicationTemplateId())
        if err != nil {
            return err
        }
    }
    if m.GetAppManagementPolicies() != nil {
        cast := make([]i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, len(m.GetAppManagementPolicies()))
        for i, v := range m.GetAppManagementPolicies() {
            if v != nil {
                cast[i] = v.(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable)
            }
        }
        err = writer.WriteCollectionOfObjectValues("appManagementPolicies", cast)
        if err != nil {
            return err
        }
    }
    if m.GetAppRoles() != nil {
        cast := make([]i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, len(m.GetAppRoles()))
        for i, v := range m.GetAppRoles() {
            if v != nil {
                cast[i] = v.(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable)
            }
        }
        err = writer.WriteCollectionOfObjectValues("appRoles", cast)
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteObjectValue("certification", m.GetCertification())
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
    {
        err = writer.WriteObjectValue("createdOnBehalfOf", m.GetCreatedOnBehalfOf())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteStringValue("defaultRedirectUri", m.GetDefaultRedirectUri())
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
        err = writer.WriteStringValue("disabledByMicrosoftStatus", m.GetDisabledByMicrosoftStatus())
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
    if m.GetExtensionProperties() != nil {
        cast := make([]i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, len(m.GetExtensionProperties()))
        for i, v := range m.GetExtensionProperties() {
            if v != nil {
                cast[i] = v.(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable)
            }
        }
        err = writer.WriteCollectionOfObjectValues("extensionProperties", cast)
        if err != nil {
            return err
        }
    }
    if m.GetFederatedIdentityCredentials() != nil {
        cast := make([]i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, len(m.GetFederatedIdentityCredentials()))
        for i, v := range m.GetFederatedIdentityCredentials() {
            if v != nil {
                cast[i] = v.(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable)
            }
        }
        err = writer.WriteCollectionOfObjectValues("federatedIdentityCredentials", cast)
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteStringValue("groupMembershipClaims", m.GetGroupMembershipClaims())
        if err != nil {
            return err
        }
    }
    if m.GetHomeRealmDiscoveryPolicies() != nil {
        cast := make([]i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, len(m.GetHomeRealmDiscoveryPolicies()))
        for i, v := range m.GetHomeRealmDiscoveryPolicies() {
            if v != nil {
                cast[i] = v.(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable)
            }
        }
        err = writer.WriteCollectionOfObjectValues("homeRealmDiscoveryPolicies", cast)
        if err != nil {
            return err
        }
    }
    if m.GetIdentifierUris() != nil {
        err = writer.WriteCollectionOfStringValues("identifierUris", m.GetIdentifierUris())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteObjectValue("info", m.GetInfo())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteBoolValue("isDeviceOnlyAuthSupported", m.GetIsDeviceOnlyAuthSupported())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteBoolValue("isFallbackPublicClient", m.GetIsFallbackPublicClient())
        if err != nil {
            return err
        }
    }
    if m.GetKeyCredentials() != nil {
        cast := make([]i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, len(m.GetKeyCredentials()))
        for i, v := range m.GetKeyCredentials() {
            if v != nil {
                cast[i] = v.(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable)
            }
        }
        err = writer.WriteCollectionOfObjectValues("keyCredentials", cast)
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteByteArrayValue("logo", m.GetLogo())
        if err != nil {
            return err
        }
    }
    if m.GetNativeAuthenticationApisEnabled() != nil {
        cast := (*m.GetNativeAuthenticationApisEnabled()).String()
        err = writer.WriteStringValue("nativeAuthenticationApisEnabled", &cast)
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteStringValue("notes", m.GetNotes())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteBoolValue("oauth2RequirePostResponse", m.GetOauth2RequirePostResponse())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteObjectValue("optionalClaims", m.GetOptionalClaims())
        if err != nil {
            return err
        }
    }
    if m.GetOwners() != nil {
        cast := make([]i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, len(m.GetOwners()))
        for i, v := range m.GetOwners() {
            if v != nil {
                cast[i] = v.(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable)
            }
        }
        err = writer.WriteCollectionOfObjectValues("owners", cast)
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteObjectValue("parentalControlSettings", m.GetParentalControlSettings())
        if err != nil {
            return err
        }
    }
    if m.GetPasswordCredentials() != nil {
        cast := make([]i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, len(m.GetPasswordCredentials()))
        for i, v := range m.GetPasswordCredentials() {
            if v != nil {
                cast[i] = v.(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable)
            }
        }
        err = writer.WriteCollectionOfObjectValues("passwordCredentials", cast)
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteObjectValue("publicClient", m.GetPublicClient())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteStringValue("publisherDomain", m.GetPublisherDomain())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteObjectValue("requestSignatureVerification", m.GetRequestSignatureVerification())
        if err != nil {
            return err
        }
    }
    if m.GetRequiredResourceAccess() != nil {
        cast := make([]i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, len(m.GetRequiredResourceAccess()))
        for i, v := range m.GetRequiredResourceAccess() {
            if v != nil {
                cast[i] = v.(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable)
            }
        }
        err = writer.WriteCollectionOfObjectValues("requiredResourceAccess", cast)
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteStringValue("samlMetadataUrl", m.GetSamlMetadataUrl())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteStringValue("serviceManagementReference", m.GetServiceManagementReference())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteObjectValue("servicePrincipalLockConfiguration", m.GetServicePrincipalLockConfiguration())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteStringValue("signInAudience", m.GetSignInAudience())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteObjectValue("spa", m.GetSpa())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteObjectValue("synchronization", m.GetSynchronization())
        if err != nil {
            return err
        }
    }
    if m.GetTags() != nil {
        err = writer.WriteCollectionOfStringValues("tags", m.GetTags())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteUUIDValue("tokenEncryptionKeyId", m.GetTokenEncryptionKeyId())
        if err != nil {
            return err
        }
    }
    if m.GetTokenIssuancePolicies() != nil {
        cast := make([]i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, len(m.GetTokenIssuancePolicies()))
        for i, v := range m.GetTokenIssuancePolicies() {
            if v != nil {
                cast[i] = v.(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable)
            }
        }
        err = writer.WriteCollectionOfObjectValues("tokenIssuancePolicies", cast)
        if err != nil {
            return err
        }
    }
    if m.GetTokenLifetimePolicies() != nil {
        cast := make([]i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, len(m.GetTokenLifetimePolicies()))
        for i, v := range m.GetTokenLifetimePolicies() {
            if v != nil {
                cast[i] = v.(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable)
            }
        }
        err = writer.WriteCollectionOfObjectValues("tokenLifetimePolicies", cast)
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteStringValue("uniqueName", m.GetUniqueName())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteObjectValue("verifiedPublisher", m.GetVerifiedPublisher())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteObjectValue("web", m.GetWeb())
        if err != nil {
            return err
        }
    }
    return nil
}
// SetAddIns sets the addIns property value. Defines custom behavior that a consuming service can use to call an app in specific contexts. For example, applications that can render file streams can set the addIns property for its 'FileHandler' functionality. This lets services like Microsoft 365 call the application in the context of a document the user is working on.
func (m *Application) SetAddIns(value []AddInable)() {
    err := m.GetBackingStore().Set("addIns", value)
    if err != nil {
        panic(err)
    }
}
// SetApi sets the api property value. Specifies settings for an application that implements a web API.
func (m *Application) SetApi(value ApiApplicationable)() {
    err := m.GetBackingStore().Set("api", value)
    if err != nil {
        panic(err)
    }
}
// SetAppId sets the appId property value. The unique identifier for the application that is assigned to an application by Microsoft Entra ID. Not nullable. Read-only. Alternate key. Supports $filter (eq).
func (m *Application) SetAppId(value *string)() {
    err := m.GetBackingStore().Set("appId", value)
    if err != nil {
        panic(err)
    }
}
// SetApplicationTemplateId sets the applicationTemplateId property value. Unique identifier of the applicationTemplate. Supports $filter (eq, not, ne). Read-only. null if the app wasn't created from an application template.
func (m *Application) SetApplicationTemplateId(value *string)() {
    err := m.GetBackingStore().Set("applicationTemplateId", value)
    if err != nil {
        panic(err)
    }
}
// SetAppManagementPolicies sets the appManagementPolicies property value. The appManagementPolicy applied to this application.
func (m *Application) SetAppManagementPolicies(value []AppManagementPolicyable)() {
    err := m.GetBackingStore().Set("appManagementPolicies", value)
    if err != nil {
        panic(err)
    }
}
// SetAppRoles sets the appRoles property value. The collection of roles defined for the application. With app role assignments, these roles can be assigned to users, groups, or service principals associated with other applications. Not nullable.
func (m *Application) SetAppRoles(value []AppRoleable)() {
    err := m.GetBackingStore().Set("appRoles", value)
    if err != nil {
        panic(err)
    }
}
// SetCertification sets the certification property value. Specifies the certification status of the application.
func (m *Application) SetCertification(value Certificationable)() {
    err := m.GetBackingStore().Set("certification", value)
    if err != nil {
        panic(err)
    }
}
// SetCreatedDateTime sets the createdDateTime property value. The date and time the application was registered. The DateTimeOffset type represents date and time information using ISO 8601 format and is always in UTC time. For example, midnight UTC on Jan 1, 2014 is 2014-01-01T00:00:00Z. Read-only.  Supports $filter (eq, ne, not, ge, le, in, and eq on null values) and $orderby.
func (m *Application) SetCreatedDateTime(value *i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)() {
    err := m.GetBackingStore().Set("createdDateTime", value)
    if err != nil {
        panic(err)
    }
}
// SetCreatedOnBehalfOf sets the createdOnBehalfOf property value. Supports $filter (/$count eq 0, /$count ne 0). Read-only.
func (m *Application) SetCreatedOnBehalfOf(value DirectoryObjectable)() {
    err := m.GetBackingStore().Set("createdOnBehalfOf", value)
    if err != nil {
        panic(err)
    }
}
// SetDefaultRedirectUri sets the defaultRedirectUri property value. The defaultRedirectUri property
func (m *Application) SetDefaultRedirectUri(value *string)() {
    err := m.GetBackingStore().Set("defaultRedirectUri", value)
    if err != nil {
        panic(err)
    }
}
// SetDescription sets the description property value. Free text field to provide a description of the application object to end users. The maximum allowed size is 1,024 characters. Supports $filter (eq, ne, not, ge, le, startsWith) and $search.
func (m *Application) SetDescription(value *string)() {
    err := m.GetBackingStore().Set("description", value)
    if err != nil {
        panic(err)
    }
}
// SetDisabledByMicrosoftStatus sets the disabledByMicrosoftStatus property value. Specifies whether Microsoft has disabled the registered application. Possible values are: null (default value), NotDisabled, and DisabledDueToViolationOfServicesAgreement (reasons include suspicious, abusive, or malicious activity, or a violation of the Microsoft Services Agreement).  Supports $filter (eq, ne, not).
func (m *Application) SetDisabledByMicrosoftStatus(value *string)() {
    err := m.GetBackingStore().Set("disabledByMicrosoftStatus", value)
    if err != nil {
        panic(err)
    }
}
// SetDisplayName sets the displayName property value. The display name for the application. Supports $filter (eq, ne, not, ge, le, in, startsWith, and eq on null values), $search, and $orderby.
func (m *Application) SetDisplayName(value *string)() {
    err := m.GetBackingStore().Set("displayName", value)
    if err != nil {
        panic(err)
    }
}
// SetExtensionProperties sets the extensionProperties property value. Read-only. Nullable. Supports $expand and $filter (/$count eq 0, /$count ne 0).
func (m *Application) SetExtensionProperties(value []ExtensionPropertyable)() {
    err := m.GetBackingStore().Set("extensionProperties", value)
    if err != nil {
        panic(err)
    }
}
// SetFederatedIdentityCredentials sets the federatedIdentityCredentials property value. Federated identities for applications. Supports $expand and $filter (startsWith, /$count eq 0, /$count ne 0).
func (m *Application) SetFederatedIdentityCredentials(value []FederatedIdentityCredentialable)() {
    err := m.GetBackingStore().Set("federatedIdentityCredentials", value)
    if err != nil {
        panic(err)
    }
}
// SetGroupMembershipClaims sets the groupMembershipClaims property value. Configures the groups claim issued in a user or OAuth 2.0 access token that the application expects. To set this attribute, use one of the following valid string values: None, SecurityGroup (for security groups and Microsoft Entra roles), All (this gets all of the security groups, distribution groups, and Microsoft Entra directory roles that the signed-in user is a member of).
func (m *Application) SetGroupMembershipClaims(value *string)() {
    err := m.GetBackingStore().Set("groupMembershipClaims", value)
    if err != nil {
        panic(err)
    }
}
// SetHomeRealmDiscoveryPolicies sets the homeRealmDiscoveryPolicies property value. The homeRealmDiscoveryPolicies property
func (m *Application) SetHomeRealmDiscoveryPolicies(value []HomeRealmDiscoveryPolicyable)() {
    err := m.GetBackingStore().Set("homeRealmDiscoveryPolicies", value)
    if err != nil {
        panic(err)
    }
}
// SetIdentifierUris sets the identifierUris property value. Also known as App ID URI, this value is set when an application is used as a resource app. The identifierUris acts as the prefix for the scopes you reference in your API's code, and it must be globally unique. You can use the default value provided, which is in the form api://<appId>, or specify a more readable URI like https://contoso.com/api. For more information on valid identifierUris patterns and best practices, see Microsoft Entra application registration security best practices. Not nullable. Supports $filter (eq, ne, ge, le, startsWith).
func (m *Application) SetIdentifierUris(value []string)() {
    err := m.GetBackingStore().Set("identifierUris", value)
    if err != nil {
        panic(err)
    }
}
// SetInfo sets the info property value. Basic profile information of the application such as  app's marketing, support, terms of service and privacy statement URLs. The terms of service and privacy statement are surfaced to users through the user consent experience. For more info, see How to: Add Terms of service and privacy statement for registered Microsoft Entra apps. Supports $filter (eq, ne, not, ge, le, and eq on null values).
func (m *Application) SetInfo(value InformationalUrlable)() {
    err := m.GetBackingStore().Set("info", value)
    if err != nil {
        panic(err)
    }
}
// SetIsDeviceOnlyAuthSupported sets the isDeviceOnlyAuthSupported property value. Specifies whether this application supports device authentication without a user. The default is false.
func (m *Application) SetIsDeviceOnlyAuthSupported(value *bool)() {
    err := m.GetBackingStore().Set("isDeviceOnlyAuthSupported", value)
    if err != nil {
        panic(err)
    }
}
// SetIsFallbackPublicClient sets the isFallbackPublicClient property value. Specifies the fallback application type as public client, such as an installed application running on a mobile device. The default value is false, which means the fallback application type is confidential client such as a web app. There are certain scenarios where Microsoft Entra ID can't determine the client application type. For example, the ROPC flow where it's configured without specifying a redirect URI. In those cases, Microsoft Entra ID interprets the application type based on the value of this property.
func (m *Application) SetIsFallbackPublicClient(value *bool)() {
    err := m.GetBackingStore().Set("isFallbackPublicClient", value)
    if err != nil {
        panic(err)
    }
}
// SetKeyCredentials sets the keyCredentials property value. The collection of key credentials associated with the application. Not nullable. Supports $filter (eq, not, ge, le).
func (m *Application) SetKeyCredentials(value []KeyCredentialable)() {
    err := m.GetBackingStore().Set("keyCredentials", value)
    if err != nil {
        panic(err)
    }
}
// SetLogo sets the logo property value. The main logo for the application. Not nullable.
func (m *Application) SetLogo(value []byte)() {
    err := m.GetBackingStore().Set("logo", value)
    if err != nil {
        panic(err)
    }
}
// SetNativeAuthenticationApisEnabled sets the nativeAuthenticationApisEnabled property value. Specifies whether the Native Authentication APIs are enabled for the application. The possible values are: none and all. Default is none. For more information, see Native Authentication.
func (m *Application) SetNativeAuthenticationApisEnabled(value *NativeAuthenticationApisEnabled)() {
    err := m.GetBackingStore().Set("nativeAuthenticationApisEnabled", value)
    if err != nil {
        panic(err)
    }
}
// SetNotes sets the notes property value. Notes relevant for the management of the application.
func (m *Application) SetNotes(value *string)() {
    err := m.GetBackingStore().Set("notes", value)
    if err != nil {
        panic(err)
    }
}
// SetOauth2RequirePostResponse sets the oauth2RequirePostResponse property value. The oauth2RequirePostResponse property
func (m *Application) SetOauth2RequirePostResponse(value *bool)() {
    err := m.GetBackingStore().Set("oauth2RequirePostResponse", value)
    if err != nil {
        panic(err)
    }
}
// SetOptionalClaims sets the optionalClaims property value. Application developers can configure optional claims in their Microsoft Entra applications to specify the claims that are sent to their application by the Microsoft security token service. For more information, see How to: Provide optional claims to your app.
func (m *Application) SetOptionalClaims(value OptionalClaimsable)() {
    err := m.GetBackingStore().Set("optionalClaims", value)
    if err != nil {
        panic(err)
    }
}
// SetOwners sets the owners property value. Directory objects that are owners of the application. Read-only. Nullable. Supports $expand, $filter (/$count eq 0, /$count ne 0, /$count eq 1, /$count ne 1), and $select nested in $expand.
func (m *Application) SetOwners(value []DirectoryObjectable)() {
    err := m.GetBackingStore().Set("owners", value)
    if err != nil {
        panic(err)
    }
}
// SetParentalControlSettings sets the parentalControlSettings property value. Specifies parental control settings for an application.
func (m *Application) SetParentalControlSettings(value ParentalControlSettingsable)() {
    err := m.GetBackingStore().Set("parentalControlSettings", value)
    if err != nil {
        panic(err)
    }
}
// SetPasswordCredentials sets the passwordCredentials property value. The collection of password credentials associated with the application. Not nullable.
func (m *Application) SetPasswordCredentials(value []PasswordCredentialable)() {
    err := m.GetBackingStore().Set("passwordCredentials", value)
    if err != nil {
        panic(err)
    }
}
// SetPublicClient sets the publicClient property value. Specifies settings for installed clients such as desktop or mobile devices.
func (m *Application) SetPublicClient(value PublicClientApplicationable)() {
    err := m.GetBackingStore().Set("publicClient", value)
    if err != nil {
        panic(err)
    }
}
// SetPublisherDomain sets the publisherDomain property value. The verified publisher domain for the application. Read-only. For more information, see How to: Configure an application's publisher domain. Supports $filter (eq, ne, ge, le, startsWith).
func (m *Application) SetPublisherDomain(value *string)() {
    err := m.GetBackingStore().Set("publisherDomain", value)
    if err != nil {
        panic(err)
    }
}
// SetRequestSignatureVerification sets the requestSignatureVerification property value. Specifies whether this application requires Microsoft Entra ID to verify the signed authentication requests.
func (m *Application) SetRequestSignatureVerification(value RequestSignatureVerificationable)() {
    err := m.GetBackingStore().Set("requestSignatureVerification", value)
    if err != nil {
        panic(err)
    }
}
// SetRequiredResourceAccess sets the requiredResourceAccess property value. Specifies the resources that the application needs to access. This property also specifies the set of delegated permissions and application roles that it needs for each of those resources. This configuration of access to the required resources drives the consent experience. No more than 50 resource services (APIs) can be configured. Beginning mid-October 2021, the total number of required permissions must not exceed 400. For more information, see Limits on requested permissions per app. Not nullable. Supports $filter (eq, not, ge, le).
func (m *Application) SetRequiredResourceAccess(value []RequiredResourceAccessable)() {
    err := m.GetBackingStore().Set("requiredResourceAccess", value)
    if err != nil {
        panic(err)
    }
}
// SetSamlMetadataUrl sets the samlMetadataUrl property value. The URL where the service exposes SAML metadata for federation. This property is valid only for single-tenant applications. Nullable.
func (m *Application) SetSamlMetadataUrl(value *string)() {
    err := m.GetBackingStore().Set("samlMetadataUrl", value)
    if err != nil {
        panic(err)
    }
}
// SetServiceManagementReference sets the serviceManagementReference property value. References application or service contact information from a Service or Asset Management database. Nullable.
func (m *Application) SetServiceManagementReference(value *string)() {
    err := m.GetBackingStore().Set("serviceManagementReference", value)
    if err != nil {
        panic(err)
    }
}
// SetServicePrincipalLockConfiguration sets the servicePrincipalLockConfiguration property value. Specifies whether sensitive properties of a multitenant application should be locked for editing after the application is provisioned in a tenant. Nullable. null by default.
func (m *Application) SetServicePrincipalLockConfiguration(value ServicePrincipalLockConfigurationable)() {
    err := m.GetBackingStore().Set("servicePrincipalLockConfiguration", value)
    if err != nil {
        panic(err)
    }
}
// SetSignInAudience sets the signInAudience property value. Specifies the Microsoft accounts that are supported for the current application. The possible values are: AzureADMyOrg (default), AzureADMultipleOrgs, AzureADandPersonalMicrosoftAccount, and PersonalMicrosoftAccount. See more in the table. The value of this object also limits the number of permissions an app can request. For more information, see Limits on requested permissions per app. The value for this property has implications on other app object properties. As a result, if you change this property, you might need to change other properties first. For more information, see Validation differences for signInAudience.Supports $filter (eq, ne, not).
func (m *Application) SetSignInAudience(value *string)() {
    err := m.GetBackingStore().Set("signInAudience", value)
    if err != nil {
        panic(err)
    }
}
// SetSpa sets the spa property value. Specifies settings for a single-page application, including sign out URLs and redirect URIs for authorization codes and access tokens.
func (m *Application) SetSpa(value SpaApplicationable)() {
    err := m.GetBackingStore().Set("spa", value)
    if err != nil {
        panic(err)
    }
}
// SetSynchronization sets the synchronization property value. Represents the capability for Microsoft Entra identity synchronization through the Microsoft Graph API.
func (m *Application) SetSynchronization(value Synchronizationable)() {
    err := m.GetBackingStore().Set("synchronization", value)
    if err != nil {
        panic(err)
    }
}
// SetTags sets the tags property value. Custom strings that can be used to categorize and identify the application. Not nullable. Strings added here will also appear in the tags property of any associated service principals.Supports $filter (eq, not, ge, le, startsWith) and $search.
func (m *Application) SetTags(value []string)() {
    err := m.GetBackingStore().Set("tags", value)
    if err != nil {
        panic(err)
    }
}
// SetTokenEncryptionKeyId sets the tokenEncryptionKeyId property value. Specifies the keyId of a public key from the keyCredentials collection. When configured, Microsoft Entra ID encrypts all the tokens it emits by using the key this property points to. The application code that receives the encrypted token must use the matching private key to decrypt the token before it can be used for the signed-in user.
func (m *Application) SetTokenEncryptionKeyId(value *i561e97a8befe7661a44c8f54600992b4207a3a0cf6770e5559949bc276de2e22.UUID)() {
    err := m.GetBackingStore().Set("tokenEncryptionKeyId", value)
    if err != nil {
        panic(err)
    }
}
// SetTokenIssuancePolicies sets the tokenIssuancePolicies property value. The tokenIssuancePolicies property
func (m *Application) SetTokenIssuancePolicies(value []TokenIssuancePolicyable)() {
    err := m.GetBackingStore().Set("tokenIssuancePolicies", value)
    if err != nil {
        panic(err)
    }
}
// SetTokenLifetimePolicies sets the tokenLifetimePolicies property value. The tokenLifetimePolicies property
func (m *Application) SetTokenLifetimePolicies(value []TokenLifetimePolicyable)() {
    err := m.GetBackingStore().Set("tokenLifetimePolicies", value)
    if err != nil {
        panic(err)
    }
}
// SetUniqueName sets the uniqueName property value. The unique identifier that can be assigned to an application and used as an alternate key. Immutable. Read-only.
func (m *Application) SetUniqueName(value *string)() {
    err := m.GetBackingStore().Set("uniqueName", value)
    if err != nil {
        panic(err)
    }
}
// SetVerifiedPublisher sets the verifiedPublisher property value. Specifies the verified publisher of the application. For more information about how publisher verification helps support application security, trustworthiness, and compliance, see Publisher verification.
func (m *Application) SetVerifiedPublisher(value VerifiedPublisherable)() {
    err := m.GetBackingStore().Set("verifiedPublisher", value)
    if err != nil {
        panic(err)
    }
}
// SetWeb sets the web property value. Specifies settings for a web application.
func (m *Application) SetWeb(value WebApplicationable)() {
    err := m.GetBackingStore().Set("web", value)
    if err != nil {
        panic(err)
    }
}
type Applicationable interface {
    DirectoryObjectable
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable
    GetAddIns()([]AddInable)
    GetApi()(ApiApplicationable)
    GetAppId()(*string)
    GetApplicationTemplateId()(*string)
    GetAppManagementPolicies()([]AppManagementPolicyable)
    GetAppRoles()([]AppRoleable)
    GetCertification()(Certificationable)
    GetCreatedDateTime()(*i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)
    GetCreatedOnBehalfOf()(DirectoryObjectable)
    GetDefaultRedirectUri()(*string)
    GetDescription()(*string)
    GetDisabledByMicrosoftStatus()(*string)
    GetDisplayName()(*string)
    GetExtensionProperties()([]ExtensionPropertyable)
    GetFederatedIdentityCredentials()([]FederatedIdentityCredentialable)
    GetGroupMembershipClaims()(*string)
    GetHomeRealmDiscoveryPolicies()([]HomeRealmDiscoveryPolicyable)
    GetIdentifierUris()([]string)
    GetInfo()(InformationalUrlable)
    GetIsDeviceOnlyAuthSupported()(*bool)
    GetIsFallbackPublicClient()(*bool)
    GetKeyCredentials()([]KeyCredentialable)
    GetLogo()([]byte)
    GetNativeAuthenticationApisEnabled()(*NativeAuthenticationApisEnabled)
    GetNotes()(*string)
    GetOauth2RequirePostResponse()(*bool)
    GetOptionalClaims()(OptionalClaimsable)
    GetOwners()([]DirectoryObjectable)
    GetParentalControlSettings()(ParentalControlSettingsable)
    GetPasswordCredentials()([]PasswordCredentialable)
    GetPublicClient()(PublicClientApplicationable)
    GetPublisherDomain()(*string)
    GetRequestSignatureVerification()(RequestSignatureVerificationable)
    GetRequiredResourceAccess()([]RequiredResourceAccessable)
    GetSamlMetadataUrl()(*string)
    GetServiceManagementReference()(*string)
    GetServicePrincipalLockConfiguration()(ServicePrincipalLockConfigurationable)
    GetSignInAudience()(*string)
    GetSpa()(SpaApplicationable)
    GetSynchronization()(Synchronizationable)
    GetTags()([]string)
    GetTokenEncryptionKeyId()(*i561e97a8befe7661a44c8f54600992b4207a3a0cf6770e5559949bc276de2e22.UUID)
    GetTokenIssuancePolicies()([]TokenIssuancePolicyable)
    GetTokenLifetimePolicies()([]TokenLifetimePolicyable)
    GetUniqueName()(*string)
    GetVerifiedPublisher()(VerifiedPublisherable)
    GetWeb()(WebApplicationable)
    SetAddIns(value []AddInable)()
    SetApi(value ApiApplicationable)()
    SetAppId(value *string)()
    SetApplicationTemplateId(value *string)()
    SetAppManagementPolicies(value []AppManagementPolicyable)()
    SetAppRoles(value []AppRoleable)()
    SetCertification(value Certificationable)()
    SetCreatedDateTime(value *i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)()
    SetCreatedOnBehalfOf(value DirectoryObjectable)()
    SetDefaultRedirectUri(value *string)()
    SetDescription(value *string)()
    SetDisabledByMicrosoftStatus(value *string)()
    SetDisplayName(value *string)()
    SetExtensionProperties(value []ExtensionPropertyable)()
    SetFederatedIdentityCredentials(value []FederatedIdentityCredentialable)()
    SetGroupMembershipClaims(value *string)()
    SetHomeRealmDiscoveryPolicies(value []HomeRealmDiscoveryPolicyable)()
    SetIdentifierUris(value []string)()
    SetInfo(value InformationalUrlable)()
    SetIsDeviceOnlyAuthSupported(value *bool)()
    SetIsFallbackPublicClient(value *bool)()
    SetKeyCredentials(value []KeyCredentialable)()
    SetLogo(value []byte)()
    SetNativeAuthenticationApisEnabled(value *NativeAuthenticationApisEnabled)()
    SetNotes(value *string)()
    SetOauth2RequirePostResponse(value *bool)()
    SetOptionalClaims(value OptionalClaimsable)()
    SetOwners(value []DirectoryObjectable)()
    SetParentalControlSettings(value ParentalControlSettingsable)()
    SetPasswordCredentials(value []PasswordCredentialable)()
    SetPublicClient(value PublicClientApplicationable)()
    SetPublisherDomain(value *string)()
    SetRequestSignatureVerification(value RequestSignatureVerificationable)()
    SetRequiredResourceAccess(value []RequiredResourceAccessable)()
    SetSamlMetadataUrl(value *string)()
    SetServiceManagementReference(value *string)()
    SetServicePrincipalLockConfiguration(value ServicePrincipalLockConfigurationable)()
    SetSignInAudience(value *string)()
    SetSpa(value SpaApplicationable)()
    SetSynchronization(value Synchronizationable)()
    SetTags(value []string)()
    SetTokenEncryptionKeyId(value *i561e97a8befe7661a44c8f54600992b4207a3a0cf6770e5559949bc276de2e22.UUID)()
    SetTokenIssuancePolicies(value []TokenIssuancePolicyable)()
    SetTokenLifetimePolicies(value []TokenLifetimePolicyable)()
    SetUniqueName(value *string)()
    SetVerifiedPublisher(value VerifiedPublisherable)()
    SetWeb(value WebApplicationable)()
}
