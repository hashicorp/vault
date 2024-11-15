package models

import (
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

type Domain struct {
    Entity
}
// NewDomain instantiates a new Domain and sets the default values.
func NewDomain()(*Domain) {
    m := &Domain{
        Entity: *NewEntity(),
    }
    return m
}
// CreateDomainFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
// returns a Parsable when successful
func CreateDomainFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewDomain(), nil
}
// GetAuthenticationType gets the authenticationType property value. Indicates the configured authentication type for the domain. The value is either Managed or Federated. Managed indicates a cloud managed domain where Microsoft Entra ID performs user authentication. Federated indicates authentication is federated with an identity provider such as the tenant's on-premises Active Directory via Active Directory Federation Services. Not nullable.  To update this property in delegated scenarios, the calling app must be assigned the Directory.AccessAsUser.All delegated permission.
// returns a *string when successful
func (m *Domain) GetAuthenticationType()(*string) {
    val, err := m.GetBackingStore().Get("authenticationType")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetAvailabilityStatus gets the availabilityStatus property value. This property is always null except when the verify action is used. When the verify action is used, a domain entity is returned in the response. The availabilityStatus property of the domain entity in the response is either AvailableImmediately or EmailVerifiedDomainTakeoverScheduled.
// returns a *string when successful
func (m *Domain) GetAvailabilityStatus()(*string) {
    val, err := m.GetBackingStore().Get("availabilityStatus")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetDomainNameReferences gets the domainNameReferences property value. The objects such as users and groups that reference the domain ID. Read-only, Nullable. Doesn't support $expand. Supports $filter by the OData type of objects returned. For example, /domains/{domainId}/domainNameReferences/microsoft.graph.user and /domains/{domainId}/domainNameReferences/microsoft.graph.group.
// returns a []DirectoryObjectable when successful
func (m *Domain) GetDomainNameReferences()([]DirectoryObjectable) {
    val, err := m.GetBackingStore().Get("domainNameReferences")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]DirectoryObjectable)
    }
    return nil
}
// GetFederationConfiguration gets the federationConfiguration property value. Domain settings configured by a customer when federated with Microsoft Entra ID. Doesn't support $expand.
// returns a []InternalDomainFederationable when successful
func (m *Domain) GetFederationConfiguration()([]InternalDomainFederationable) {
    val, err := m.GetBackingStore().Get("federationConfiguration")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]InternalDomainFederationable)
    }
    return nil
}
// GetFieldDeserializers the deserialization information for the current model
// returns a map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error) when successful
func (m *Domain) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := m.Entity.GetFieldDeserializers()
    res["authenticationType"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetAuthenticationType(val)
        }
        return nil
    }
    res["availabilityStatus"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetAvailabilityStatus(val)
        }
        return nil
    }
    res["domainNameReferences"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
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
            m.SetDomainNameReferences(res)
        }
        return nil
    }
    res["federationConfiguration"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetCollectionOfObjectValues(CreateInternalDomainFederationFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            res := make([]InternalDomainFederationable, len(val))
            for i, v := range val {
                if v != nil {
                    res[i] = v.(InternalDomainFederationable)
                }
            }
            m.SetFederationConfiguration(res)
        }
        return nil
    }
    res["isAdminManaged"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetBoolValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetIsAdminManaged(val)
        }
        return nil
    }
    res["isDefault"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetBoolValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetIsDefault(val)
        }
        return nil
    }
    res["isInitial"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetBoolValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetIsInitial(val)
        }
        return nil
    }
    res["isRoot"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetBoolValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetIsRoot(val)
        }
        return nil
    }
    res["isVerified"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetBoolValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetIsVerified(val)
        }
        return nil
    }
    res["manufacturer"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetManufacturer(val)
        }
        return nil
    }
    res["model"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetModel(val)
        }
        return nil
    }
    res["passwordNotificationWindowInDays"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetInt32Value()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetPasswordNotificationWindowInDays(val)
        }
        return nil
    }
    res["passwordValidityPeriodInDays"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetInt32Value()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetPasswordValidityPeriodInDays(val)
        }
        return nil
    }
    res["rootDomain"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetObjectValue(CreateDomainFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetRootDomain(val.(Domainable))
        }
        return nil
    }
    res["serviceConfigurationRecords"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetCollectionOfObjectValues(CreateDomainDnsRecordFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            res := make([]DomainDnsRecordable, len(val))
            for i, v := range val {
                if v != nil {
                    res[i] = v.(DomainDnsRecordable)
                }
            }
            m.SetServiceConfigurationRecords(res)
        }
        return nil
    }
    res["state"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetObjectValue(CreateDomainStateFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetState(val.(DomainStateable))
        }
        return nil
    }
    res["supportedServices"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
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
            m.SetSupportedServices(res)
        }
        return nil
    }
    res["verificationDnsRecords"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetCollectionOfObjectValues(CreateDomainDnsRecordFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            res := make([]DomainDnsRecordable, len(val))
            for i, v := range val {
                if v != nil {
                    res[i] = v.(DomainDnsRecordable)
                }
            }
            m.SetVerificationDnsRecords(res)
        }
        return nil
    }
    return res
}
// GetIsAdminManaged gets the isAdminManaged property value. The value of the property is false if the DNS record management of the domain is delegated to Microsoft 365. Otherwise, the value is true. Not nullable
// returns a *bool when successful
func (m *Domain) GetIsAdminManaged()(*bool) {
    val, err := m.GetBackingStore().Get("isAdminManaged")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*bool)
    }
    return nil
}
// GetIsDefault gets the isDefault property value. true if this is the default domain that is used for user creation. There's only one default domain per company. Not nullable.
// returns a *bool when successful
func (m *Domain) GetIsDefault()(*bool) {
    val, err := m.GetBackingStore().Get("isDefault")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*bool)
    }
    return nil
}
// GetIsInitial gets the isInitial property value. true if this is the initial domain created by Microsoft Online Services (contoso.com). There's only one initial domain per company. Not nullable
// returns a *bool when successful
func (m *Domain) GetIsInitial()(*bool) {
    val, err := m.GetBackingStore().Get("isInitial")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*bool)
    }
    return nil
}
// GetIsRoot gets the isRoot property value. true if the domain is a verified root domain. Otherwise, false if the domain is a subdomain or unverified. Not nullable.
// returns a *bool when successful
func (m *Domain) GetIsRoot()(*bool) {
    val, err := m.GetBackingStore().Get("isRoot")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*bool)
    }
    return nil
}
// GetIsVerified gets the isVerified property value. true if the domain completed domain ownership verification. Not nullable.
// returns a *bool when successful
func (m *Domain) GetIsVerified()(*bool) {
    val, err := m.GetBackingStore().Get("isVerified")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*bool)
    }
    return nil
}
// GetManufacturer gets the manufacturer property value. The manufacturer property
// returns a *string when successful
func (m *Domain) GetManufacturer()(*string) {
    val, err := m.GetBackingStore().Get("manufacturer")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetModel gets the model property value. The model property
// returns a *string when successful
func (m *Domain) GetModel()(*string) {
    val, err := m.GetBackingStore().Get("model")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetPasswordNotificationWindowInDays gets the passwordNotificationWindowInDays property value. Specifies the number of days before a user receives notification that their password expires. If the property isn't set, a default value of 14 days is used.
// returns a *int32 when successful
func (m *Domain) GetPasswordNotificationWindowInDays()(*int32) {
    val, err := m.GetBackingStore().Get("passwordNotificationWindowInDays")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*int32)
    }
    return nil
}
// GetPasswordValidityPeriodInDays gets the passwordValidityPeriodInDays property value. Specifies the length of time that a password is valid before it must be changed. If the property isn't set, a default value of 90 days is used.
// returns a *int32 when successful
func (m *Domain) GetPasswordValidityPeriodInDays()(*int32) {
    val, err := m.GetBackingStore().Get("passwordValidityPeriodInDays")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*int32)
    }
    return nil
}
// GetRootDomain gets the rootDomain property value. Root domain of a subdomain. Read-only, Nullable. Supports $expand.
// returns a Domainable when successful
func (m *Domain) GetRootDomain()(Domainable) {
    val, err := m.GetBackingStore().Get("rootDomain")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(Domainable)
    }
    return nil
}
// GetServiceConfigurationRecords gets the serviceConfigurationRecords property value. DNS records the customer adds to the DNS zone file of the domain before the domain can be used by Microsoft Online services. Read-only, Nullable. Doesn't support $expand.
// returns a []DomainDnsRecordable when successful
func (m *Domain) GetServiceConfigurationRecords()([]DomainDnsRecordable) {
    val, err := m.GetBackingStore().Get("serviceConfigurationRecords")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]DomainDnsRecordable)
    }
    return nil
}
// GetState gets the state property value. Status of asynchronous operations scheduled for the domain.
// returns a DomainStateable when successful
func (m *Domain) GetState()(DomainStateable) {
    val, err := m.GetBackingStore().Get("state")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(DomainStateable)
    }
    return nil
}
// GetSupportedServices gets the supportedServices property value. The capabilities assigned to the domain. Can include 0, 1 or more of following values: Email, Sharepoint, EmailInternalRelayOnly, OfficeCommunicationsOnline, SharePointDefaultDomain, FullRedelegation, SharePointPublic, OrgIdAuthentication, Yammer, Intune. The values that you can add or remove using the API include: Email, OfficeCommunicationsOnline, Yammer. Not nullable.
// returns a []string when successful
func (m *Domain) GetSupportedServices()([]string) {
    val, err := m.GetBackingStore().Get("supportedServices")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]string)
    }
    return nil
}
// GetVerificationDnsRecords gets the verificationDnsRecords property value. DNS records that the customer adds to the DNS zone file of the domain before the customer can complete domain ownership verification with Microsoft Entra ID. Read-only, Nullable. Doesn't support $expand.
// returns a []DomainDnsRecordable when successful
func (m *Domain) GetVerificationDnsRecords()([]DomainDnsRecordable) {
    val, err := m.GetBackingStore().Get("verificationDnsRecords")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]DomainDnsRecordable)
    }
    return nil
}
// Serialize serializes information the current object
func (m *Domain) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    err := m.Entity.Serialize(writer)
    if err != nil {
        return err
    }
    {
        err = writer.WriteStringValue("authenticationType", m.GetAuthenticationType())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteStringValue("availabilityStatus", m.GetAvailabilityStatus())
        if err != nil {
            return err
        }
    }
    if m.GetDomainNameReferences() != nil {
        cast := make([]i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, len(m.GetDomainNameReferences()))
        for i, v := range m.GetDomainNameReferences() {
            if v != nil {
                cast[i] = v.(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable)
            }
        }
        err = writer.WriteCollectionOfObjectValues("domainNameReferences", cast)
        if err != nil {
            return err
        }
    }
    if m.GetFederationConfiguration() != nil {
        cast := make([]i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, len(m.GetFederationConfiguration()))
        for i, v := range m.GetFederationConfiguration() {
            if v != nil {
                cast[i] = v.(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable)
            }
        }
        err = writer.WriteCollectionOfObjectValues("federationConfiguration", cast)
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteBoolValue("isAdminManaged", m.GetIsAdminManaged())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteBoolValue("isDefault", m.GetIsDefault())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteBoolValue("isInitial", m.GetIsInitial())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteBoolValue("isRoot", m.GetIsRoot())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteBoolValue("isVerified", m.GetIsVerified())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteStringValue("manufacturer", m.GetManufacturer())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteStringValue("model", m.GetModel())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteInt32Value("passwordNotificationWindowInDays", m.GetPasswordNotificationWindowInDays())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteInt32Value("passwordValidityPeriodInDays", m.GetPasswordValidityPeriodInDays())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteObjectValue("rootDomain", m.GetRootDomain())
        if err != nil {
            return err
        }
    }
    if m.GetServiceConfigurationRecords() != nil {
        cast := make([]i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, len(m.GetServiceConfigurationRecords()))
        for i, v := range m.GetServiceConfigurationRecords() {
            if v != nil {
                cast[i] = v.(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable)
            }
        }
        err = writer.WriteCollectionOfObjectValues("serviceConfigurationRecords", cast)
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteObjectValue("state", m.GetState())
        if err != nil {
            return err
        }
    }
    if m.GetSupportedServices() != nil {
        err = writer.WriteCollectionOfStringValues("supportedServices", m.GetSupportedServices())
        if err != nil {
            return err
        }
    }
    if m.GetVerificationDnsRecords() != nil {
        cast := make([]i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, len(m.GetVerificationDnsRecords()))
        for i, v := range m.GetVerificationDnsRecords() {
            if v != nil {
                cast[i] = v.(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable)
            }
        }
        err = writer.WriteCollectionOfObjectValues("verificationDnsRecords", cast)
        if err != nil {
            return err
        }
    }
    return nil
}
// SetAuthenticationType sets the authenticationType property value. Indicates the configured authentication type for the domain. The value is either Managed or Federated. Managed indicates a cloud managed domain where Microsoft Entra ID performs user authentication. Federated indicates authentication is federated with an identity provider such as the tenant's on-premises Active Directory via Active Directory Federation Services. Not nullable.  To update this property in delegated scenarios, the calling app must be assigned the Directory.AccessAsUser.All delegated permission.
func (m *Domain) SetAuthenticationType(value *string)() {
    err := m.GetBackingStore().Set("authenticationType", value)
    if err != nil {
        panic(err)
    }
}
// SetAvailabilityStatus sets the availabilityStatus property value. This property is always null except when the verify action is used. When the verify action is used, a domain entity is returned in the response. The availabilityStatus property of the domain entity in the response is either AvailableImmediately or EmailVerifiedDomainTakeoverScheduled.
func (m *Domain) SetAvailabilityStatus(value *string)() {
    err := m.GetBackingStore().Set("availabilityStatus", value)
    if err != nil {
        panic(err)
    }
}
// SetDomainNameReferences sets the domainNameReferences property value. The objects such as users and groups that reference the domain ID. Read-only, Nullable. Doesn't support $expand. Supports $filter by the OData type of objects returned. For example, /domains/{domainId}/domainNameReferences/microsoft.graph.user and /domains/{domainId}/domainNameReferences/microsoft.graph.group.
func (m *Domain) SetDomainNameReferences(value []DirectoryObjectable)() {
    err := m.GetBackingStore().Set("domainNameReferences", value)
    if err != nil {
        panic(err)
    }
}
// SetFederationConfiguration sets the federationConfiguration property value. Domain settings configured by a customer when federated with Microsoft Entra ID. Doesn't support $expand.
func (m *Domain) SetFederationConfiguration(value []InternalDomainFederationable)() {
    err := m.GetBackingStore().Set("federationConfiguration", value)
    if err != nil {
        panic(err)
    }
}
// SetIsAdminManaged sets the isAdminManaged property value. The value of the property is false if the DNS record management of the domain is delegated to Microsoft 365. Otherwise, the value is true. Not nullable
func (m *Domain) SetIsAdminManaged(value *bool)() {
    err := m.GetBackingStore().Set("isAdminManaged", value)
    if err != nil {
        panic(err)
    }
}
// SetIsDefault sets the isDefault property value. true if this is the default domain that is used for user creation. There's only one default domain per company. Not nullable.
func (m *Domain) SetIsDefault(value *bool)() {
    err := m.GetBackingStore().Set("isDefault", value)
    if err != nil {
        panic(err)
    }
}
// SetIsInitial sets the isInitial property value. true if this is the initial domain created by Microsoft Online Services (contoso.com). There's only one initial domain per company. Not nullable
func (m *Domain) SetIsInitial(value *bool)() {
    err := m.GetBackingStore().Set("isInitial", value)
    if err != nil {
        panic(err)
    }
}
// SetIsRoot sets the isRoot property value. true if the domain is a verified root domain. Otherwise, false if the domain is a subdomain or unverified. Not nullable.
func (m *Domain) SetIsRoot(value *bool)() {
    err := m.GetBackingStore().Set("isRoot", value)
    if err != nil {
        panic(err)
    }
}
// SetIsVerified sets the isVerified property value. true if the domain completed domain ownership verification. Not nullable.
func (m *Domain) SetIsVerified(value *bool)() {
    err := m.GetBackingStore().Set("isVerified", value)
    if err != nil {
        panic(err)
    }
}
// SetManufacturer sets the manufacturer property value. The manufacturer property
func (m *Domain) SetManufacturer(value *string)() {
    err := m.GetBackingStore().Set("manufacturer", value)
    if err != nil {
        panic(err)
    }
}
// SetModel sets the model property value. The model property
func (m *Domain) SetModel(value *string)() {
    err := m.GetBackingStore().Set("model", value)
    if err != nil {
        panic(err)
    }
}
// SetPasswordNotificationWindowInDays sets the passwordNotificationWindowInDays property value. Specifies the number of days before a user receives notification that their password expires. If the property isn't set, a default value of 14 days is used.
func (m *Domain) SetPasswordNotificationWindowInDays(value *int32)() {
    err := m.GetBackingStore().Set("passwordNotificationWindowInDays", value)
    if err != nil {
        panic(err)
    }
}
// SetPasswordValidityPeriodInDays sets the passwordValidityPeriodInDays property value. Specifies the length of time that a password is valid before it must be changed. If the property isn't set, a default value of 90 days is used.
func (m *Domain) SetPasswordValidityPeriodInDays(value *int32)() {
    err := m.GetBackingStore().Set("passwordValidityPeriodInDays", value)
    if err != nil {
        panic(err)
    }
}
// SetRootDomain sets the rootDomain property value. Root domain of a subdomain. Read-only, Nullable. Supports $expand.
func (m *Domain) SetRootDomain(value Domainable)() {
    err := m.GetBackingStore().Set("rootDomain", value)
    if err != nil {
        panic(err)
    }
}
// SetServiceConfigurationRecords sets the serviceConfigurationRecords property value. DNS records the customer adds to the DNS zone file of the domain before the domain can be used by Microsoft Online services. Read-only, Nullable. Doesn't support $expand.
func (m *Domain) SetServiceConfigurationRecords(value []DomainDnsRecordable)() {
    err := m.GetBackingStore().Set("serviceConfigurationRecords", value)
    if err != nil {
        panic(err)
    }
}
// SetState sets the state property value. Status of asynchronous operations scheduled for the domain.
func (m *Domain) SetState(value DomainStateable)() {
    err := m.GetBackingStore().Set("state", value)
    if err != nil {
        panic(err)
    }
}
// SetSupportedServices sets the supportedServices property value. The capabilities assigned to the domain. Can include 0, 1 or more of following values: Email, Sharepoint, EmailInternalRelayOnly, OfficeCommunicationsOnline, SharePointDefaultDomain, FullRedelegation, SharePointPublic, OrgIdAuthentication, Yammer, Intune. The values that you can add or remove using the API include: Email, OfficeCommunicationsOnline, Yammer. Not nullable.
func (m *Domain) SetSupportedServices(value []string)() {
    err := m.GetBackingStore().Set("supportedServices", value)
    if err != nil {
        panic(err)
    }
}
// SetVerificationDnsRecords sets the verificationDnsRecords property value. DNS records that the customer adds to the DNS zone file of the domain before the customer can complete domain ownership verification with Microsoft Entra ID. Read-only, Nullable. Doesn't support $expand.
func (m *Domain) SetVerificationDnsRecords(value []DomainDnsRecordable)() {
    err := m.GetBackingStore().Set("verificationDnsRecords", value)
    if err != nil {
        panic(err)
    }
}
type Domainable interface {
    Entityable
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable
    GetAuthenticationType()(*string)
    GetAvailabilityStatus()(*string)
    GetDomainNameReferences()([]DirectoryObjectable)
    GetFederationConfiguration()([]InternalDomainFederationable)
    GetIsAdminManaged()(*bool)
    GetIsDefault()(*bool)
    GetIsInitial()(*bool)
    GetIsRoot()(*bool)
    GetIsVerified()(*bool)
    GetManufacturer()(*string)
    GetModel()(*string)
    GetPasswordNotificationWindowInDays()(*int32)
    GetPasswordValidityPeriodInDays()(*int32)
    GetRootDomain()(Domainable)
    GetServiceConfigurationRecords()([]DomainDnsRecordable)
    GetState()(DomainStateable)
    GetSupportedServices()([]string)
    GetVerificationDnsRecords()([]DomainDnsRecordable)
    SetAuthenticationType(value *string)()
    SetAvailabilityStatus(value *string)()
    SetDomainNameReferences(value []DirectoryObjectable)()
    SetFederationConfiguration(value []InternalDomainFederationable)()
    SetIsAdminManaged(value *bool)()
    SetIsDefault(value *bool)()
    SetIsInitial(value *bool)()
    SetIsRoot(value *bool)()
    SetIsVerified(value *bool)()
    SetManufacturer(value *string)()
    SetModel(value *string)()
    SetPasswordNotificationWindowInDays(value *int32)()
    SetPasswordValidityPeriodInDays(value *int32)()
    SetRootDomain(value Domainable)()
    SetServiceConfigurationRecords(value []DomainDnsRecordable)()
    SetState(value DomainStateable)()
    SetSupportedServices(value []string)()
    SetVerificationDnsRecords(value []DomainDnsRecordable)()
}
