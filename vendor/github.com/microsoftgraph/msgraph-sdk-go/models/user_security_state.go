package models

import (
    i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e "time"
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
    ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e "github.com/microsoft/kiota-abstractions-go/store"
)

type UserSecurityState struct {
    // Stores model information.
    backingStore ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackingStore
}
// NewUserSecurityState instantiates a new UserSecurityState and sets the default values.
func NewUserSecurityState()(*UserSecurityState) {
    m := &UserSecurityState{
    }
    m.backingStore = ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackingStoreFactoryInstance();
    m.SetAdditionalData(make(map[string]any))
    return m
}
// CreateUserSecurityStateFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
// returns a Parsable when successful
func CreateUserSecurityStateFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewUserSecurityState(), nil
}
// GetAadUserId gets the aadUserId property value. AAD User object identifier (GUID) - represents the physical/multi-account user entity.
// returns a *string when successful
func (m *UserSecurityState) GetAadUserId()(*string) {
    val, err := m.GetBackingStore().Get("aadUserId")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetAccountName gets the accountName property value. Account name of user account (without Active Directory domain or DNS domain) - (also called mailNickName).
// returns a *string when successful
func (m *UserSecurityState) GetAccountName()(*string) {
    val, err := m.GetBackingStore().Get("accountName")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetAdditionalData gets the AdditionalData property value. Stores additional data not described in the OpenAPI description found when deserializing. Can be used for serialization as well.
// returns a map[string]any when successful
func (m *UserSecurityState) GetAdditionalData()(map[string]any) {
    val , err :=  m.backingStore.Get("additionalData")
    if err != nil {
        panic(err)
    }
    if val == nil {
        var value = make(map[string]any);
        m.SetAdditionalData(value);
    }
    return val.(map[string]any)
}
// GetBackingStore gets the BackingStore property value. Stores model information.
// returns a BackingStore when successful
func (m *UserSecurityState) GetBackingStore()(ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackingStore) {
    return m.backingStore
}
// GetDomainName gets the domainName property value. NetBIOS/Active Directory domain of user account (that is, domain/account format).
// returns a *string when successful
func (m *UserSecurityState) GetDomainName()(*string) {
    val, err := m.GetBackingStore().Get("domainName")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetEmailRole gets the emailRole property value. For email-related alerts - user account's email 'role'. Possible values are: unknown, sender, recipient.
// returns a *EmailRole when successful
func (m *UserSecurityState) GetEmailRole()(*EmailRole) {
    val, err := m.GetBackingStore().Get("emailRole")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*EmailRole)
    }
    return nil
}
// GetFieldDeserializers the deserialization information for the current model
// returns a map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error) when successful
func (m *UserSecurityState) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := make(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error))
    res["aadUserId"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetAadUserId(val)
        }
        return nil
    }
    res["accountName"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetAccountName(val)
        }
        return nil
    }
    res["domainName"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetDomainName(val)
        }
        return nil
    }
    res["emailRole"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetEnumValue(ParseEmailRole)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetEmailRole(val.(*EmailRole))
        }
        return nil
    }
    res["isVpn"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetBoolValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetIsVpn(val)
        }
        return nil
    }
    res["logonDateTime"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetTimeValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetLogonDateTime(val)
        }
        return nil
    }
    res["logonId"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetLogonId(val)
        }
        return nil
    }
    res["logonIp"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetLogonIp(val)
        }
        return nil
    }
    res["logonLocation"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetLogonLocation(val)
        }
        return nil
    }
    res["logonType"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetEnumValue(ParseLogonType)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetLogonType(val.(*LogonType))
        }
        return nil
    }
    res["@odata.type"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetOdataType(val)
        }
        return nil
    }
    res["onPremisesSecurityIdentifier"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetOnPremisesSecurityIdentifier(val)
        }
        return nil
    }
    res["riskScore"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetRiskScore(val)
        }
        return nil
    }
    res["userAccountType"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetEnumValue(ParseUserAccountSecurityType)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetUserAccountType(val.(*UserAccountSecurityType))
        }
        return nil
    }
    res["userPrincipalName"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetUserPrincipalName(val)
        }
        return nil
    }
    return res
}
// GetIsVpn gets the isVpn property value. Indicates whether the user logged on through a VPN.
// returns a *bool when successful
func (m *UserSecurityState) GetIsVpn()(*bool) {
    val, err := m.GetBackingStore().Get("isVpn")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*bool)
    }
    return nil
}
// GetLogonDateTime gets the logonDateTime property value. Time at which the sign-in occurred. The Timestamp type represents date and time information using ISO 8601 format and is always in UTC time. For example, midnight UTC on Jan 1, 2014 is 2014-01-01T00:00:00Z.
// returns a *Time when successful
func (m *UserSecurityState) GetLogonDateTime()(*i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time) {
    val, err := m.GetBackingStore().Get("logonDateTime")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)
    }
    return nil
}
// GetLogonId gets the logonId property value. User sign-in ID.
// returns a *string when successful
func (m *UserSecurityState) GetLogonId()(*string) {
    val, err := m.GetBackingStore().Get("logonId")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetLogonIp gets the logonIp property value. IP Address the sign-in request originated from.
// returns a *string when successful
func (m *UserSecurityState) GetLogonIp()(*string) {
    val, err := m.GetBackingStore().Get("logonIp")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetLogonLocation gets the logonLocation property value. Location (by IP address mapping) associated with a user sign-in event by this user.
// returns a *string when successful
func (m *UserSecurityState) GetLogonLocation()(*string) {
    val, err := m.GetBackingStore().Get("logonLocation")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetLogonType gets the logonType property value. Method of user sign in. Possible values are: unknown, interactive, remoteInteractive, network, batch, service.
// returns a *LogonType when successful
func (m *UserSecurityState) GetLogonType()(*LogonType) {
    val, err := m.GetBackingStore().Get("logonType")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*LogonType)
    }
    return nil
}
// GetOdataType gets the @odata.type property value. The OdataType property
// returns a *string when successful
func (m *UserSecurityState) GetOdataType()(*string) {
    val, err := m.GetBackingStore().Get("odataType")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetOnPremisesSecurityIdentifier gets the onPremisesSecurityIdentifier property value. Active Directory (on-premises) Security Identifier (SID) of the user.
// returns a *string when successful
func (m *UserSecurityState) GetOnPremisesSecurityIdentifier()(*string) {
    val, err := m.GetBackingStore().Get("onPremisesSecurityIdentifier")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetRiskScore gets the riskScore property value. Provider-generated/calculated risk score of the user account. Recommended value range of 0-1, which equates to a percentage.
// returns a *string when successful
func (m *UserSecurityState) GetRiskScore()(*string) {
    val, err := m.GetBackingStore().Get("riskScore")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetUserAccountType gets the userAccountType property value. User account type (group membership), per Windows definition. Possible values are: unknown, standard, power, administrator.
// returns a *UserAccountSecurityType when successful
func (m *UserSecurityState) GetUserAccountType()(*UserAccountSecurityType) {
    val, err := m.GetBackingStore().Get("userAccountType")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*UserAccountSecurityType)
    }
    return nil
}
// GetUserPrincipalName gets the userPrincipalName property value. User sign-in name - internet format: (user account name)@(user account DNS domain name).
// returns a *string when successful
func (m *UserSecurityState) GetUserPrincipalName()(*string) {
    val, err := m.GetBackingStore().Get("userPrincipalName")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// Serialize serializes information the current object
func (m *UserSecurityState) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    {
        err := writer.WriteStringValue("aadUserId", m.GetAadUserId())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteStringValue("accountName", m.GetAccountName())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteStringValue("domainName", m.GetDomainName())
        if err != nil {
            return err
        }
    }
    if m.GetEmailRole() != nil {
        cast := (*m.GetEmailRole()).String()
        err := writer.WriteStringValue("emailRole", &cast)
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteBoolValue("isVpn", m.GetIsVpn())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteTimeValue("logonDateTime", m.GetLogonDateTime())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteStringValue("logonId", m.GetLogonId())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteStringValue("logonIp", m.GetLogonIp())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteStringValue("logonLocation", m.GetLogonLocation())
        if err != nil {
            return err
        }
    }
    if m.GetLogonType() != nil {
        cast := (*m.GetLogonType()).String()
        err := writer.WriteStringValue("logonType", &cast)
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteStringValue("@odata.type", m.GetOdataType())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteStringValue("onPremisesSecurityIdentifier", m.GetOnPremisesSecurityIdentifier())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteStringValue("riskScore", m.GetRiskScore())
        if err != nil {
            return err
        }
    }
    if m.GetUserAccountType() != nil {
        cast := (*m.GetUserAccountType()).String()
        err := writer.WriteStringValue("userAccountType", &cast)
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteStringValue("userPrincipalName", m.GetUserPrincipalName())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteAdditionalData(m.GetAdditionalData())
        if err != nil {
            return err
        }
    }
    return nil
}
// SetAadUserId sets the aadUserId property value. AAD User object identifier (GUID) - represents the physical/multi-account user entity.
func (m *UserSecurityState) SetAadUserId(value *string)() {
    err := m.GetBackingStore().Set("aadUserId", value)
    if err != nil {
        panic(err)
    }
}
// SetAccountName sets the accountName property value. Account name of user account (without Active Directory domain or DNS domain) - (also called mailNickName).
func (m *UserSecurityState) SetAccountName(value *string)() {
    err := m.GetBackingStore().Set("accountName", value)
    if err != nil {
        panic(err)
    }
}
// SetAdditionalData sets the AdditionalData property value. Stores additional data not described in the OpenAPI description found when deserializing. Can be used for serialization as well.
func (m *UserSecurityState) SetAdditionalData(value map[string]any)() {
    err := m.GetBackingStore().Set("additionalData", value)
    if err != nil {
        panic(err)
    }
}
// SetBackingStore sets the BackingStore property value. Stores model information.
func (m *UserSecurityState) SetBackingStore(value ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackingStore)() {
    m.backingStore = value
}
// SetDomainName sets the domainName property value. NetBIOS/Active Directory domain of user account (that is, domain/account format).
func (m *UserSecurityState) SetDomainName(value *string)() {
    err := m.GetBackingStore().Set("domainName", value)
    if err != nil {
        panic(err)
    }
}
// SetEmailRole sets the emailRole property value. For email-related alerts - user account's email 'role'. Possible values are: unknown, sender, recipient.
func (m *UserSecurityState) SetEmailRole(value *EmailRole)() {
    err := m.GetBackingStore().Set("emailRole", value)
    if err != nil {
        panic(err)
    }
}
// SetIsVpn sets the isVpn property value. Indicates whether the user logged on through a VPN.
func (m *UserSecurityState) SetIsVpn(value *bool)() {
    err := m.GetBackingStore().Set("isVpn", value)
    if err != nil {
        panic(err)
    }
}
// SetLogonDateTime sets the logonDateTime property value. Time at which the sign-in occurred. The Timestamp type represents date and time information using ISO 8601 format and is always in UTC time. For example, midnight UTC on Jan 1, 2014 is 2014-01-01T00:00:00Z.
func (m *UserSecurityState) SetLogonDateTime(value *i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)() {
    err := m.GetBackingStore().Set("logonDateTime", value)
    if err != nil {
        panic(err)
    }
}
// SetLogonId sets the logonId property value. User sign-in ID.
func (m *UserSecurityState) SetLogonId(value *string)() {
    err := m.GetBackingStore().Set("logonId", value)
    if err != nil {
        panic(err)
    }
}
// SetLogonIp sets the logonIp property value. IP Address the sign-in request originated from.
func (m *UserSecurityState) SetLogonIp(value *string)() {
    err := m.GetBackingStore().Set("logonIp", value)
    if err != nil {
        panic(err)
    }
}
// SetLogonLocation sets the logonLocation property value. Location (by IP address mapping) associated with a user sign-in event by this user.
func (m *UserSecurityState) SetLogonLocation(value *string)() {
    err := m.GetBackingStore().Set("logonLocation", value)
    if err != nil {
        panic(err)
    }
}
// SetLogonType sets the logonType property value. Method of user sign in. Possible values are: unknown, interactive, remoteInteractive, network, batch, service.
func (m *UserSecurityState) SetLogonType(value *LogonType)() {
    err := m.GetBackingStore().Set("logonType", value)
    if err != nil {
        panic(err)
    }
}
// SetOdataType sets the @odata.type property value. The OdataType property
func (m *UserSecurityState) SetOdataType(value *string)() {
    err := m.GetBackingStore().Set("odataType", value)
    if err != nil {
        panic(err)
    }
}
// SetOnPremisesSecurityIdentifier sets the onPremisesSecurityIdentifier property value. Active Directory (on-premises) Security Identifier (SID) of the user.
func (m *UserSecurityState) SetOnPremisesSecurityIdentifier(value *string)() {
    err := m.GetBackingStore().Set("onPremisesSecurityIdentifier", value)
    if err != nil {
        panic(err)
    }
}
// SetRiskScore sets the riskScore property value. Provider-generated/calculated risk score of the user account. Recommended value range of 0-1, which equates to a percentage.
func (m *UserSecurityState) SetRiskScore(value *string)() {
    err := m.GetBackingStore().Set("riskScore", value)
    if err != nil {
        panic(err)
    }
}
// SetUserAccountType sets the userAccountType property value. User account type (group membership), per Windows definition. Possible values are: unknown, standard, power, administrator.
func (m *UserSecurityState) SetUserAccountType(value *UserAccountSecurityType)() {
    err := m.GetBackingStore().Set("userAccountType", value)
    if err != nil {
        panic(err)
    }
}
// SetUserPrincipalName sets the userPrincipalName property value. User sign-in name - internet format: (user account name)@(user account DNS domain name).
func (m *UserSecurityState) SetUserPrincipalName(value *string)() {
    err := m.GetBackingStore().Set("userPrincipalName", value)
    if err != nil {
        panic(err)
    }
}
type UserSecurityStateable interface {
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.AdditionalDataHolder
    ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackedModel
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable
    GetAadUserId()(*string)
    GetAccountName()(*string)
    GetBackingStore()(ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackingStore)
    GetDomainName()(*string)
    GetEmailRole()(*EmailRole)
    GetIsVpn()(*bool)
    GetLogonDateTime()(*i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)
    GetLogonId()(*string)
    GetLogonIp()(*string)
    GetLogonLocation()(*string)
    GetLogonType()(*LogonType)
    GetOdataType()(*string)
    GetOnPremisesSecurityIdentifier()(*string)
    GetRiskScore()(*string)
    GetUserAccountType()(*UserAccountSecurityType)
    GetUserPrincipalName()(*string)
    SetAadUserId(value *string)()
    SetAccountName(value *string)()
    SetBackingStore(value ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackingStore)()
    SetDomainName(value *string)()
    SetEmailRole(value *EmailRole)()
    SetIsVpn(value *bool)()
    SetLogonDateTime(value *i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)()
    SetLogonId(value *string)()
    SetLogonIp(value *string)()
    SetLogonLocation(value *string)()
    SetLogonType(value *LogonType)()
    SetOdataType(value *string)()
    SetOnPremisesSecurityIdentifier(value *string)()
    SetRiskScore(value *string)()
    SetUserAccountType(value *UserAccountSecurityType)()
    SetUserPrincipalName(value *string)()
}
