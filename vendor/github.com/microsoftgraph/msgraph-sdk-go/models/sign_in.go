package models

import (
    i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e "time"
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

type SignIn struct {
    Entity
}
// NewSignIn instantiates a new SignIn and sets the default values.
func NewSignIn()(*SignIn) {
    m := &SignIn{
        Entity: *NewEntity(),
    }
    return m
}
// CreateSignInFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
// returns a Parsable when successful
func CreateSignInFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewSignIn(), nil
}
// GetAppDisplayName gets the appDisplayName property value. App name displayed in the Microsoft Entra admin center.  Supports $filter (eq, startsWith).
// returns a *string when successful
func (m *SignIn) GetAppDisplayName()(*string) {
    val, err := m.GetBackingStore().Get("appDisplayName")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetAppId gets the appId property value. Unique GUID that represents the app ID in the Microsoft Entra ID.  Supports $filter (eq).
// returns a *string when successful
func (m *SignIn) GetAppId()(*string) {
    val, err := m.GetBackingStore().Get("appId")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetAppliedConditionalAccessPolicies gets the appliedConditionalAccessPolicies property value. Provides a list of conditional access policies that the corresponding sign-in activity triggers. Apps need more Conditional Access-related privileges to read the details of this property. For more information, see Viewing applied conditional access (CA) policies in sign-ins.
// returns a []AppliedConditionalAccessPolicyable when successful
func (m *SignIn) GetAppliedConditionalAccessPolicies()([]AppliedConditionalAccessPolicyable) {
    val, err := m.GetBackingStore().Get("appliedConditionalAccessPolicies")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]AppliedConditionalAccessPolicyable)
    }
    return nil
}
// GetClientAppUsed gets the clientAppUsed property value. Identifies the client used for the sign-in activity. Modern authentication clients include Browser, modern clients. Legacy authentication clients include Exchange ActiveSync, IMAP, MAPI, SMTP, POP, and other clients.  Supports $filter (eq).
// returns a *string when successful
func (m *SignIn) GetClientAppUsed()(*string) {
    val, err := m.GetBackingStore().Get("clientAppUsed")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetConditionalAccessStatus gets the conditionalAccessStatus property value. Reports status of an activated conditional access policy. Possible values are: success, failure, notApplied, and unknownFutureValue.  Supports $filter (eq).
// returns a *ConditionalAccessStatus when successful
func (m *SignIn) GetConditionalAccessStatus()(*ConditionalAccessStatus) {
    val, err := m.GetBackingStore().Get("conditionalAccessStatus")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*ConditionalAccessStatus)
    }
    return nil
}
// GetCorrelationId gets the correlationId property value. The request ID sent from the client when the sign-in is initiated. Used to troubleshoot sign-in activity.  Supports $filter (eq).
// returns a *string when successful
func (m *SignIn) GetCorrelationId()(*string) {
    val, err := m.GetBackingStore().Get("correlationId")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetCreatedDateTime gets the createdDateTime property value. Date and time (UTC) the sign-in was initiated. Example: midnight on Jan 1, 2014 is reported as 2014-01-01T00:00:00Z.  Supports $orderby, $filter (eq, le, and ge).
// returns a *Time when successful
func (m *SignIn) GetCreatedDateTime()(*i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time) {
    val, err := m.GetBackingStore().Get("createdDateTime")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)
    }
    return nil
}
// GetDeviceDetail gets the deviceDetail property value. Device information from where the sign-in occurred; includes device ID, operating system, and browser.  Supports $filter (eq, startsWith) on browser and operatingSytem properties.
// returns a DeviceDetailable when successful
func (m *SignIn) GetDeviceDetail()(DeviceDetailable) {
    val, err := m.GetBackingStore().Get("deviceDetail")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(DeviceDetailable)
    }
    return nil
}
// GetFieldDeserializers the deserialization information for the current model
// returns a map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error) when successful
func (m *SignIn) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := m.Entity.GetFieldDeserializers()
    res["appDisplayName"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetAppDisplayName(val)
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
    res["appliedConditionalAccessPolicies"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetCollectionOfObjectValues(CreateAppliedConditionalAccessPolicyFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            res := make([]AppliedConditionalAccessPolicyable, len(val))
            for i, v := range val {
                if v != nil {
                    res[i] = v.(AppliedConditionalAccessPolicyable)
                }
            }
            m.SetAppliedConditionalAccessPolicies(res)
        }
        return nil
    }
    res["clientAppUsed"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetClientAppUsed(val)
        }
        return nil
    }
    res["conditionalAccessStatus"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetEnumValue(ParseConditionalAccessStatus)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetConditionalAccessStatus(val.(*ConditionalAccessStatus))
        }
        return nil
    }
    res["correlationId"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetCorrelationId(val)
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
    res["deviceDetail"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetObjectValue(CreateDeviceDetailFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetDeviceDetail(val.(DeviceDetailable))
        }
        return nil
    }
    res["ipAddress"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetIpAddress(val)
        }
        return nil
    }
    res["isInteractive"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetBoolValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetIsInteractive(val)
        }
        return nil
    }
    res["location"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetObjectValue(CreateSignInLocationFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetLocation(val.(SignInLocationable))
        }
        return nil
    }
    res["resourceDisplayName"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetResourceDisplayName(val)
        }
        return nil
    }
    res["resourceId"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetResourceId(val)
        }
        return nil
    }
    res["riskDetail"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetEnumValue(ParseRiskDetail)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetRiskDetail(val.(*RiskDetail))
        }
        return nil
    }
    res["riskEventTypes"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetCollectionOfEnumValues(ParseRiskEventType)
        if err != nil {
            return err
        }
        if val != nil {
            res := make([]RiskEventType, len(val))
            for i, v := range val {
                if v != nil {
                    res[i] = *(v.(*RiskEventType))
                }
            }
            m.SetRiskEventTypes(res)
        }
        return nil
    }
    res["riskEventTypes_v2"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
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
            m.SetRiskEventTypesV2(res)
        }
        return nil
    }
    res["riskLevelAggregated"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetEnumValue(ParseRiskLevel)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetRiskLevelAggregated(val.(*RiskLevel))
        }
        return nil
    }
    res["riskLevelDuringSignIn"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetEnumValue(ParseRiskLevel)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetRiskLevelDuringSignIn(val.(*RiskLevel))
        }
        return nil
    }
    res["riskState"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetEnumValue(ParseRiskState)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetRiskState(val.(*RiskState))
        }
        return nil
    }
    res["status"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetObjectValue(CreateSignInStatusFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetStatus(val.(SignInStatusable))
        }
        return nil
    }
    res["userDisplayName"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetUserDisplayName(val)
        }
        return nil
    }
    res["userId"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetUserId(val)
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
// GetIpAddress gets the ipAddress property value. IP address of the client used to sign in.  Supports $filter (eq, startsWith).
// returns a *string when successful
func (m *SignIn) GetIpAddress()(*string) {
    val, err := m.GetBackingStore().Get("ipAddress")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetIsInteractive gets the isInteractive property value. Indicates whether a sign-in is interactive.
// returns a *bool when successful
func (m *SignIn) GetIsInteractive()(*bool) {
    val, err := m.GetBackingStore().Get("isInteractive")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*bool)
    }
    return nil
}
// GetLocation gets the location property value. Provides the city, state, and country code where the sign-in originated.  Supports $filter (eq, startsWith) on city, state, and countryOrRegion properties.
// returns a SignInLocationable when successful
func (m *SignIn) GetLocation()(SignInLocationable) {
    val, err := m.GetBackingStore().Get("location")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(SignInLocationable)
    }
    return nil
}
// GetResourceDisplayName gets the resourceDisplayName property value. Name of the resource the user signed into.  Supports $filter (eq).
// returns a *string when successful
func (m *SignIn) GetResourceDisplayName()(*string) {
    val, err := m.GetBackingStore().Get("resourceDisplayName")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetResourceId gets the resourceId property value. ID of the resource that the user signed into.  Supports $filter (eq).
// returns a *string when successful
func (m *SignIn) GetResourceId()(*string) {
    val, err := m.GetBackingStore().Get("resourceId")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetRiskDetail gets the riskDetail property value. The reason behind a specific state of a risky user, sign-in, or a risk event. The possible values are none, adminGeneratedTemporaryPassword, userPerformedSecuredPasswordChange, userPerformedSecuredPasswordReset, adminConfirmedSigninSafe, aiConfirmedSigninSafe, userPassedMFADrivenByRiskBasedPolicy, adminDismissedAllRiskForUser, adminConfirmedSigninCompromised, hidden, adminConfirmedUserCompromised, unknownFutureValue, adminConfirmedServicePrincipalCompromised, adminDismissedAllRiskForServicePrincipal, m365DAdminDismissedDetection, userChangedPasswordOnPremises, adminDismissedRiskForSignIn, adminConfirmedAccountSafe. You must use the Prefer: include-unknown-enum-members request header to get the following value or values in this evolvable enum: adminConfirmedServicePrincipalCompromised, adminDismissedAllRiskForServicePrincipal, m365DAdminDismissedDetection, userChangedPasswordOnPremises, adminDismissedRiskForSignIn, adminConfirmedAccountSafe.The value none means that Microsoft Entra risk detection did not flag the user or the sign-in as a risky event so far.  Supports $filter (eq). Note: Details for this property are only available for Microsoft Entra ID P2 customers. All other customers are returned hidden.
// returns a *RiskDetail when successful
func (m *SignIn) GetRiskDetail()(*RiskDetail) {
    val, err := m.GetBackingStore().Get("riskDetail")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*RiskDetail)
    }
    return nil
}
// GetRiskEventTypes gets the riskEventTypes property value. The riskEventTypes property
// returns a []RiskEventType when successful
func (m *SignIn) GetRiskEventTypes()([]RiskEventType) {
    val, err := m.GetBackingStore().Get("riskEventTypes")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]RiskEventType)
    }
    return nil
}
// GetRiskEventTypesV2 gets the riskEventTypes_v2 property value. The list of risk event types associated with the sign-in. Possible values: unlikelyTravel, anonymizedIPAddress, maliciousIPAddress, unfamiliarFeatures, malwareInfectedIPAddress, suspiciousIPAddress, leakedCredentials, investigationsThreatIntelligence, generic, or unknownFutureValue.  Supports $filter (eq, startsWith).
// returns a []string when successful
func (m *SignIn) GetRiskEventTypesV2()([]string) {
    val, err := m.GetBackingStore().Get("riskEventTypes_v2")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]string)
    }
    return nil
}
// GetRiskLevelAggregated gets the riskLevelAggregated property value. Aggregated risk level. The possible values are: none, low, medium, high, hidden, and unknownFutureValue. The value hidden means the user or sign-in wasn't enabled for Microsoft Entra ID Protection.  Supports $filter (eq).  Note: Details for this property are only available for Microsoft Entra ID P2 customers. All other customers are returned hidden.
// returns a *RiskLevel when successful
func (m *SignIn) GetRiskLevelAggregated()(*RiskLevel) {
    val, err := m.GetBackingStore().Get("riskLevelAggregated")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*RiskLevel)
    }
    return nil
}
// GetRiskLevelDuringSignIn gets the riskLevelDuringSignIn property value. Risk level during sign-in. The possible values are: none, low, medium, high, hidden, and unknownFutureValue. The value hidden means the user or sign-in wasn't enabled for Microsoft Entra ID Protection.  Supports $filter (eq). Note: Details for this property are only available for Microsoft Entra ID P2 customers. All other customers are returned hidden.
// returns a *RiskLevel when successful
func (m *SignIn) GetRiskLevelDuringSignIn()(*RiskLevel) {
    val, err := m.GetBackingStore().Get("riskLevelDuringSignIn")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*RiskLevel)
    }
    return nil
}
// GetRiskState gets the riskState property value. Reports status of the risky user, sign-in, or a risk event. The possible values are: none, confirmedSafe, remediated, dismissed, atRisk, confirmedCompromised, unknownFutureValue.  Supports $filter (eq).
// returns a *RiskState when successful
func (m *SignIn) GetRiskState()(*RiskState) {
    val, err := m.GetBackingStore().Get("riskState")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*RiskState)
    }
    return nil
}
// GetStatus gets the status property value. Sign-in status. Includes the error code and description of the error (if a sign-in failure occurs).  Supports $filter (eq) on errorCode property.
// returns a SignInStatusable when successful
func (m *SignIn) GetStatus()(SignInStatusable) {
    val, err := m.GetBackingStore().Get("status")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(SignInStatusable)
    }
    return nil
}
// GetUserDisplayName gets the userDisplayName property value. Display name of the user that initiated the sign-in.  Supports $filter (eq, startsWith).
// returns a *string when successful
func (m *SignIn) GetUserDisplayName()(*string) {
    val, err := m.GetBackingStore().Get("userDisplayName")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetUserId gets the userId property value. ID of the user that initiated the sign-in.  Supports $filter (eq).
// returns a *string when successful
func (m *SignIn) GetUserId()(*string) {
    val, err := m.GetBackingStore().Get("userId")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetUserPrincipalName gets the userPrincipalName property value. User principal name of the user that initiated the sign-in. This value is always in lowercase. For guest users whose values in the user object typically contain #EXT# before the domain part, this property stores the value in both lowercase and the 'true' format. For example, while the user object stores AdeleVance_fabrikam.com#EXT#@contoso.com, the sign-in logs store adelevance@fabrikam.com. Supports $filter (eq, startsWith).
// returns a *string when successful
func (m *SignIn) GetUserPrincipalName()(*string) {
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
func (m *SignIn) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    err := m.Entity.Serialize(writer)
    if err != nil {
        return err
    }
    {
        err = writer.WriteStringValue("appDisplayName", m.GetAppDisplayName())
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
    if m.GetAppliedConditionalAccessPolicies() != nil {
        cast := make([]i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, len(m.GetAppliedConditionalAccessPolicies()))
        for i, v := range m.GetAppliedConditionalAccessPolicies() {
            if v != nil {
                cast[i] = v.(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable)
            }
        }
        err = writer.WriteCollectionOfObjectValues("appliedConditionalAccessPolicies", cast)
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteStringValue("clientAppUsed", m.GetClientAppUsed())
        if err != nil {
            return err
        }
    }
    if m.GetConditionalAccessStatus() != nil {
        cast := (*m.GetConditionalAccessStatus()).String()
        err = writer.WriteStringValue("conditionalAccessStatus", &cast)
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteStringValue("correlationId", m.GetCorrelationId())
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
        err = writer.WriteObjectValue("deviceDetail", m.GetDeviceDetail())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteStringValue("ipAddress", m.GetIpAddress())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteBoolValue("isInteractive", m.GetIsInteractive())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteObjectValue("location", m.GetLocation())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteStringValue("resourceDisplayName", m.GetResourceDisplayName())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteStringValue("resourceId", m.GetResourceId())
        if err != nil {
            return err
        }
    }
    if m.GetRiskDetail() != nil {
        cast := (*m.GetRiskDetail()).String()
        err = writer.WriteStringValue("riskDetail", &cast)
        if err != nil {
            return err
        }
    }
    if m.GetRiskEventTypes() != nil {
        err = writer.WriteCollectionOfStringValues("riskEventTypes", SerializeRiskEventType(m.GetRiskEventTypes()))
        if err != nil {
            return err
        }
    }
    if m.GetRiskEventTypesV2() != nil {
        err = writer.WriteCollectionOfStringValues("riskEventTypes_v2", m.GetRiskEventTypesV2())
        if err != nil {
            return err
        }
    }
    if m.GetRiskLevelAggregated() != nil {
        cast := (*m.GetRiskLevelAggregated()).String()
        err = writer.WriteStringValue("riskLevelAggregated", &cast)
        if err != nil {
            return err
        }
    }
    if m.GetRiskLevelDuringSignIn() != nil {
        cast := (*m.GetRiskLevelDuringSignIn()).String()
        err = writer.WriteStringValue("riskLevelDuringSignIn", &cast)
        if err != nil {
            return err
        }
    }
    if m.GetRiskState() != nil {
        cast := (*m.GetRiskState()).String()
        err = writer.WriteStringValue("riskState", &cast)
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteObjectValue("status", m.GetStatus())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteStringValue("userDisplayName", m.GetUserDisplayName())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteStringValue("userId", m.GetUserId())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteStringValue("userPrincipalName", m.GetUserPrincipalName())
        if err != nil {
            return err
        }
    }
    return nil
}
// SetAppDisplayName sets the appDisplayName property value. App name displayed in the Microsoft Entra admin center.  Supports $filter (eq, startsWith).
func (m *SignIn) SetAppDisplayName(value *string)() {
    err := m.GetBackingStore().Set("appDisplayName", value)
    if err != nil {
        panic(err)
    }
}
// SetAppId sets the appId property value. Unique GUID that represents the app ID in the Microsoft Entra ID.  Supports $filter (eq).
func (m *SignIn) SetAppId(value *string)() {
    err := m.GetBackingStore().Set("appId", value)
    if err != nil {
        panic(err)
    }
}
// SetAppliedConditionalAccessPolicies sets the appliedConditionalAccessPolicies property value. Provides a list of conditional access policies that the corresponding sign-in activity triggers. Apps need more Conditional Access-related privileges to read the details of this property. For more information, see Viewing applied conditional access (CA) policies in sign-ins.
func (m *SignIn) SetAppliedConditionalAccessPolicies(value []AppliedConditionalAccessPolicyable)() {
    err := m.GetBackingStore().Set("appliedConditionalAccessPolicies", value)
    if err != nil {
        panic(err)
    }
}
// SetClientAppUsed sets the clientAppUsed property value. Identifies the client used for the sign-in activity. Modern authentication clients include Browser, modern clients. Legacy authentication clients include Exchange ActiveSync, IMAP, MAPI, SMTP, POP, and other clients.  Supports $filter (eq).
func (m *SignIn) SetClientAppUsed(value *string)() {
    err := m.GetBackingStore().Set("clientAppUsed", value)
    if err != nil {
        panic(err)
    }
}
// SetConditionalAccessStatus sets the conditionalAccessStatus property value. Reports status of an activated conditional access policy. Possible values are: success, failure, notApplied, and unknownFutureValue.  Supports $filter (eq).
func (m *SignIn) SetConditionalAccessStatus(value *ConditionalAccessStatus)() {
    err := m.GetBackingStore().Set("conditionalAccessStatus", value)
    if err != nil {
        panic(err)
    }
}
// SetCorrelationId sets the correlationId property value. The request ID sent from the client when the sign-in is initiated. Used to troubleshoot sign-in activity.  Supports $filter (eq).
func (m *SignIn) SetCorrelationId(value *string)() {
    err := m.GetBackingStore().Set("correlationId", value)
    if err != nil {
        panic(err)
    }
}
// SetCreatedDateTime sets the createdDateTime property value. Date and time (UTC) the sign-in was initiated. Example: midnight on Jan 1, 2014 is reported as 2014-01-01T00:00:00Z.  Supports $orderby, $filter (eq, le, and ge).
func (m *SignIn) SetCreatedDateTime(value *i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)() {
    err := m.GetBackingStore().Set("createdDateTime", value)
    if err != nil {
        panic(err)
    }
}
// SetDeviceDetail sets the deviceDetail property value. Device information from where the sign-in occurred; includes device ID, operating system, and browser.  Supports $filter (eq, startsWith) on browser and operatingSytem properties.
func (m *SignIn) SetDeviceDetail(value DeviceDetailable)() {
    err := m.GetBackingStore().Set("deviceDetail", value)
    if err != nil {
        panic(err)
    }
}
// SetIpAddress sets the ipAddress property value. IP address of the client used to sign in.  Supports $filter (eq, startsWith).
func (m *SignIn) SetIpAddress(value *string)() {
    err := m.GetBackingStore().Set("ipAddress", value)
    if err != nil {
        panic(err)
    }
}
// SetIsInteractive sets the isInteractive property value. Indicates whether a sign-in is interactive.
func (m *SignIn) SetIsInteractive(value *bool)() {
    err := m.GetBackingStore().Set("isInteractive", value)
    if err != nil {
        panic(err)
    }
}
// SetLocation sets the location property value. Provides the city, state, and country code where the sign-in originated.  Supports $filter (eq, startsWith) on city, state, and countryOrRegion properties.
func (m *SignIn) SetLocation(value SignInLocationable)() {
    err := m.GetBackingStore().Set("location", value)
    if err != nil {
        panic(err)
    }
}
// SetResourceDisplayName sets the resourceDisplayName property value. Name of the resource the user signed into.  Supports $filter (eq).
func (m *SignIn) SetResourceDisplayName(value *string)() {
    err := m.GetBackingStore().Set("resourceDisplayName", value)
    if err != nil {
        panic(err)
    }
}
// SetResourceId sets the resourceId property value. ID of the resource that the user signed into.  Supports $filter (eq).
func (m *SignIn) SetResourceId(value *string)() {
    err := m.GetBackingStore().Set("resourceId", value)
    if err != nil {
        panic(err)
    }
}
// SetRiskDetail sets the riskDetail property value. The reason behind a specific state of a risky user, sign-in, or a risk event. The possible values are none, adminGeneratedTemporaryPassword, userPerformedSecuredPasswordChange, userPerformedSecuredPasswordReset, adminConfirmedSigninSafe, aiConfirmedSigninSafe, userPassedMFADrivenByRiskBasedPolicy, adminDismissedAllRiskForUser, adminConfirmedSigninCompromised, hidden, adminConfirmedUserCompromised, unknownFutureValue, adminConfirmedServicePrincipalCompromised, adminDismissedAllRiskForServicePrincipal, m365DAdminDismissedDetection, userChangedPasswordOnPremises, adminDismissedRiskForSignIn, adminConfirmedAccountSafe. You must use the Prefer: include-unknown-enum-members request header to get the following value or values in this evolvable enum: adminConfirmedServicePrincipalCompromised, adminDismissedAllRiskForServicePrincipal, m365DAdminDismissedDetection, userChangedPasswordOnPremises, adminDismissedRiskForSignIn, adminConfirmedAccountSafe.The value none means that Microsoft Entra risk detection did not flag the user or the sign-in as a risky event so far.  Supports $filter (eq). Note: Details for this property are only available for Microsoft Entra ID P2 customers. All other customers are returned hidden.
func (m *SignIn) SetRiskDetail(value *RiskDetail)() {
    err := m.GetBackingStore().Set("riskDetail", value)
    if err != nil {
        panic(err)
    }
}
// SetRiskEventTypes sets the riskEventTypes property value. The riskEventTypes property
func (m *SignIn) SetRiskEventTypes(value []RiskEventType)() {
    err := m.GetBackingStore().Set("riskEventTypes", value)
    if err != nil {
        panic(err)
    }
}
// SetRiskEventTypesV2 sets the riskEventTypes_v2 property value. The list of risk event types associated with the sign-in. Possible values: unlikelyTravel, anonymizedIPAddress, maliciousIPAddress, unfamiliarFeatures, malwareInfectedIPAddress, suspiciousIPAddress, leakedCredentials, investigationsThreatIntelligence, generic, or unknownFutureValue.  Supports $filter (eq, startsWith).
func (m *SignIn) SetRiskEventTypesV2(value []string)() {
    err := m.GetBackingStore().Set("riskEventTypes_v2", value)
    if err != nil {
        panic(err)
    }
}
// SetRiskLevelAggregated sets the riskLevelAggregated property value. Aggregated risk level. The possible values are: none, low, medium, high, hidden, and unknownFutureValue. The value hidden means the user or sign-in wasn't enabled for Microsoft Entra ID Protection.  Supports $filter (eq).  Note: Details for this property are only available for Microsoft Entra ID P2 customers. All other customers are returned hidden.
func (m *SignIn) SetRiskLevelAggregated(value *RiskLevel)() {
    err := m.GetBackingStore().Set("riskLevelAggregated", value)
    if err != nil {
        panic(err)
    }
}
// SetRiskLevelDuringSignIn sets the riskLevelDuringSignIn property value. Risk level during sign-in. The possible values are: none, low, medium, high, hidden, and unknownFutureValue. The value hidden means the user or sign-in wasn't enabled for Microsoft Entra ID Protection.  Supports $filter (eq). Note: Details for this property are only available for Microsoft Entra ID P2 customers. All other customers are returned hidden.
func (m *SignIn) SetRiskLevelDuringSignIn(value *RiskLevel)() {
    err := m.GetBackingStore().Set("riskLevelDuringSignIn", value)
    if err != nil {
        panic(err)
    }
}
// SetRiskState sets the riskState property value. Reports status of the risky user, sign-in, or a risk event. The possible values are: none, confirmedSafe, remediated, dismissed, atRisk, confirmedCompromised, unknownFutureValue.  Supports $filter (eq).
func (m *SignIn) SetRiskState(value *RiskState)() {
    err := m.GetBackingStore().Set("riskState", value)
    if err != nil {
        panic(err)
    }
}
// SetStatus sets the status property value. Sign-in status. Includes the error code and description of the error (if a sign-in failure occurs).  Supports $filter (eq) on errorCode property.
func (m *SignIn) SetStatus(value SignInStatusable)() {
    err := m.GetBackingStore().Set("status", value)
    if err != nil {
        panic(err)
    }
}
// SetUserDisplayName sets the userDisplayName property value. Display name of the user that initiated the sign-in.  Supports $filter (eq, startsWith).
func (m *SignIn) SetUserDisplayName(value *string)() {
    err := m.GetBackingStore().Set("userDisplayName", value)
    if err != nil {
        panic(err)
    }
}
// SetUserId sets the userId property value. ID of the user that initiated the sign-in.  Supports $filter (eq).
func (m *SignIn) SetUserId(value *string)() {
    err := m.GetBackingStore().Set("userId", value)
    if err != nil {
        panic(err)
    }
}
// SetUserPrincipalName sets the userPrincipalName property value. User principal name of the user that initiated the sign-in. This value is always in lowercase. For guest users whose values in the user object typically contain #EXT# before the domain part, this property stores the value in both lowercase and the 'true' format. For example, while the user object stores AdeleVance_fabrikam.com#EXT#@contoso.com, the sign-in logs store adelevance@fabrikam.com. Supports $filter (eq, startsWith).
func (m *SignIn) SetUserPrincipalName(value *string)() {
    err := m.GetBackingStore().Set("userPrincipalName", value)
    if err != nil {
        panic(err)
    }
}
type SignInable interface {
    Entityable
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable
    GetAppDisplayName()(*string)
    GetAppId()(*string)
    GetAppliedConditionalAccessPolicies()([]AppliedConditionalAccessPolicyable)
    GetClientAppUsed()(*string)
    GetConditionalAccessStatus()(*ConditionalAccessStatus)
    GetCorrelationId()(*string)
    GetCreatedDateTime()(*i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)
    GetDeviceDetail()(DeviceDetailable)
    GetIpAddress()(*string)
    GetIsInteractive()(*bool)
    GetLocation()(SignInLocationable)
    GetResourceDisplayName()(*string)
    GetResourceId()(*string)
    GetRiskDetail()(*RiskDetail)
    GetRiskEventTypes()([]RiskEventType)
    GetRiskEventTypesV2()([]string)
    GetRiskLevelAggregated()(*RiskLevel)
    GetRiskLevelDuringSignIn()(*RiskLevel)
    GetRiskState()(*RiskState)
    GetStatus()(SignInStatusable)
    GetUserDisplayName()(*string)
    GetUserId()(*string)
    GetUserPrincipalName()(*string)
    SetAppDisplayName(value *string)()
    SetAppId(value *string)()
    SetAppliedConditionalAccessPolicies(value []AppliedConditionalAccessPolicyable)()
    SetClientAppUsed(value *string)()
    SetConditionalAccessStatus(value *ConditionalAccessStatus)()
    SetCorrelationId(value *string)()
    SetCreatedDateTime(value *i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)()
    SetDeviceDetail(value DeviceDetailable)()
    SetIpAddress(value *string)()
    SetIsInteractive(value *bool)()
    SetLocation(value SignInLocationable)()
    SetResourceDisplayName(value *string)()
    SetResourceId(value *string)()
    SetRiskDetail(value *RiskDetail)()
    SetRiskEventTypes(value []RiskEventType)()
    SetRiskEventTypesV2(value []string)()
    SetRiskLevelAggregated(value *RiskLevel)()
    SetRiskLevelDuringSignIn(value *RiskLevel)()
    SetRiskState(value *RiskState)()
    SetStatus(value SignInStatusable)()
    SetUserDisplayName(value *string)()
    SetUserId(value *string)()
    SetUserPrincipalName(value *string)()
}
