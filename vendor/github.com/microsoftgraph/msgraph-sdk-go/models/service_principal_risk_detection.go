package models

import (
    i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e "time"
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

type ServicePrincipalRiskDetection struct {
    Entity
}
// NewServicePrincipalRiskDetection instantiates a new ServicePrincipalRiskDetection and sets the default values.
func NewServicePrincipalRiskDetection()(*ServicePrincipalRiskDetection) {
    m := &ServicePrincipalRiskDetection{
        Entity: *NewEntity(),
    }
    return m
}
// CreateServicePrincipalRiskDetectionFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
// returns a Parsable when successful
func CreateServicePrincipalRiskDetectionFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewServicePrincipalRiskDetection(), nil
}
// GetActivity gets the activity property value. Indicates the activity type the detected risk is linked to.  The possible values are: signin, servicePrincipal. Note that you must use the Prefer: include-unknown-enum-members request header to get the following value(s) in this evolvable enum: servicePrincipal.
// returns a *ActivityType when successful
func (m *ServicePrincipalRiskDetection) GetActivity()(*ActivityType) {
    val, err := m.GetBackingStore().Get("activity")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*ActivityType)
    }
    return nil
}
// GetActivityDateTime gets the activityDateTime property value. Date and time when the risky activity occurred. The DateTimeOffset type represents date and time information using ISO 8601 format and is always in UTC time. For example, midnight UTC on Jan 1, 2014 is 2014-01-01T00:00:00Z
// returns a *Time when successful
func (m *ServicePrincipalRiskDetection) GetActivityDateTime()(*i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time) {
    val, err := m.GetBackingStore().Get("activityDateTime")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)
    }
    return nil
}
// GetAdditionalInfo gets the additionalInfo property value. Additional information associated with the risk detection. This string value is represented as a JSON object with the quotations escaped.
// returns a *string when successful
func (m *ServicePrincipalRiskDetection) GetAdditionalInfo()(*string) {
    val, err := m.GetBackingStore().Get("additionalInfo")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetAppId gets the appId property value. The unique identifier for the associated application.
// returns a *string when successful
func (m *ServicePrincipalRiskDetection) GetAppId()(*string) {
    val, err := m.GetBackingStore().Get("appId")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetCorrelationId gets the correlationId property value. Correlation ID of the sign-in activity associated with the risk detection. This property is null if the risk detection is not associated with a sign-in activity.
// returns a *string when successful
func (m *ServicePrincipalRiskDetection) GetCorrelationId()(*string) {
    val, err := m.GetBackingStore().Get("correlationId")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetDetectedDateTime gets the detectedDateTime property value. Date and time when the risk was detected. The DateTimeOffset type represents date and time information using ISO 8601 format and is always in UTC time. For example, midnight UTC on Jan 1, 2014 is 2014-01-01T00:00:00Z.
// returns a *Time when successful
func (m *ServicePrincipalRiskDetection) GetDetectedDateTime()(*i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time) {
    val, err := m.GetBackingStore().Get("detectedDateTime")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)
    }
    return nil
}
// GetDetectionTimingType gets the detectionTimingType property value. Timing of the detected risk , whether real-time or offline. The possible values are: notDefined, realtime, nearRealtime, offline, unknownFutureValue.
// returns a *RiskDetectionTimingType when successful
func (m *ServicePrincipalRiskDetection) GetDetectionTimingType()(*RiskDetectionTimingType) {
    val, err := m.GetBackingStore().Get("detectionTimingType")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*RiskDetectionTimingType)
    }
    return nil
}
// GetFieldDeserializers the deserialization information for the current model
// returns a map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error) when successful
func (m *ServicePrincipalRiskDetection) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := m.Entity.GetFieldDeserializers()
    res["activity"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetEnumValue(ParseActivityType)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetActivity(val.(*ActivityType))
        }
        return nil
    }
    res["activityDateTime"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetTimeValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetActivityDateTime(val)
        }
        return nil
    }
    res["additionalInfo"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetAdditionalInfo(val)
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
    res["detectedDateTime"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetTimeValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetDetectedDateTime(val)
        }
        return nil
    }
    res["detectionTimingType"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetEnumValue(ParseRiskDetectionTimingType)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetDetectionTimingType(val.(*RiskDetectionTimingType))
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
    res["keyIds"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
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
            m.SetKeyIds(res)
        }
        return nil
    }
    res["lastUpdatedDateTime"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetTimeValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetLastUpdatedDateTime(val)
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
    res["requestId"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetRequestId(val)
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
    res["riskEventType"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetRiskEventType(val)
        }
        return nil
    }
    res["riskLevel"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetEnumValue(ParseRiskLevel)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetRiskLevel(val.(*RiskLevel))
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
    res["servicePrincipalDisplayName"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetServicePrincipalDisplayName(val)
        }
        return nil
    }
    res["servicePrincipalId"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetServicePrincipalId(val)
        }
        return nil
    }
    res["source"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetSource(val)
        }
        return nil
    }
    res["tokenIssuerType"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetEnumValue(ParseTokenIssuerType)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetTokenIssuerType(val.(*TokenIssuerType))
        }
        return nil
    }
    return res
}
// GetIpAddress gets the ipAddress property value. Provides the IP address of the client from where the risk occurred.
// returns a *string when successful
func (m *ServicePrincipalRiskDetection) GetIpAddress()(*string) {
    val, err := m.GetBackingStore().Get("ipAddress")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetKeyIds gets the keyIds property value. The unique identifier for the key credential associated with the risk detection.
// returns a []string when successful
func (m *ServicePrincipalRiskDetection) GetKeyIds()([]string) {
    val, err := m.GetBackingStore().Get("keyIds")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]string)
    }
    return nil
}
// GetLastUpdatedDateTime gets the lastUpdatedDateTime property value. Date and time when the risk detection was last updated.
// returns a *Time when successful
func (m *ServicePrincipalRiskDetection) GetLastUpdatedDateTime()(*i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time) {
    val, err := m.GetBackingStore().Get("lastUpdatedDateTime")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)
    }
    return nil
}
// GetLocation gets the location property value. Location from where the sign-in was initiated.
// returns a SignInLocationable when successful
func (m *ServicePrincipalRiskDetection) GetLocation()(SignInLocationable) {
    val, err := m.GetBackingStore().Get("location")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(SignInLocationable)
    }
    return nil
}
// GetRequestId gets the requestId property value. Request identifier of the sign-in activity associated with the risk detection. This property is null if the risk detection is not associated with a sign-in activity. Supports $filter (eq).
// returns a *string when successful
func (m *ServicePrincipalRiskDetection) GetRequestId()(*string) {
    val, err := m.GetBackingStore().Get("requestId")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetRiskDetail gets the riskDetail property value. Details of the detected risk. Note: Details for this property are only available for Workload Identities Premium customers. Events in tenants without this license will be returned hidden. The possible values are: none, hidden, adminConfirmedServicePrincipalCompromised, adminDismissedAllRiskForServicePrincipal. Note that you must use the Prefer: include-unknown-enum-members request header to get the following value(s) in this evolvable enum: adminConfirmedServicePrincipalCompromised , adminDismissedAllRiskForServicePrincipal.
// returns a *RiskDetail when successful
func (m *ServicePrincipalRiskDetection) GetRiskDetail()(*RiskDetail) {
    val, err := m.GetBackingStore().Get("riskDetail")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*RiskDetail)
    }
    return nil
}
// GetRiskEventType gets the riskEventType property value. The type of risk event detected. The possible values are: investigationsThreatIntelligence, generic, adminConfirmedServicePrincipalCompromised, suspiciousSignins, leakedCredentials, anomalousServicePrincipalActivity, maliciousApplication, suspiciousApplication.
// returns a *string when successful
func (m *ServicePrincipalRiskDetection) GetRiskEventType()(*string) {
    val, err := m.GetBackingStore().Get("riskEventType")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetRiskLevel gets the riskLevel property value. Level of the detected risk. Note: Details for this property are only available for Workload Identities Premium customers. Events in tenants without this license will be returned hidden. The possible values are: low, medium, high, hidden, none.
// returns a *RiskLevel when successful
func (m *ServicePrincipalRiskDetection) GetRiskLevel()(*RiskLevel) {
    val, err := m.GetBackingStore().Get("riskLevel")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*RiskLevel)
    }
    return nil
}
// GetRiskState gets the riskState property value. The state of a detected risky service principal or sign-in activity. The possible values are: none, dismissed, atRisk, confirmedCompromised.
// returns a *RiskState when successful
func (m *ServicePrincipalRiskDetection) GetRiskState()(*RiskState) {
    val, err := m.GetBackingStore().Get("riskState")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*RiskState)
    }
    return nil
}
// GetServicePrincipalDisplayName gets the servicePrincipalDisplayName property value. The display name for the service principal.
// returns a *string when successful
func (m *ServicePrincipalRiskDetection) GetServicePrincipalDisplayName()(*string) {
    val, err := m.GetBackingStore().Get("servicePrincipalDisplayName")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetServicePrincipalId gets the servicePrincipalId property value. The unique identifier for the service principal. Supports $filter (eq).
// returns a *string when successful
func (m *ServicePrincipalRiskDetection) GetServicePrincipalId()(*string) {
    val, err := m.GetBackingStore().Get("servicePrincipalId")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetSource gets the source property value. Source of the risk detection. For example, identityProtection.
// returns a *string when successful
func (m *ServicePrincipalRiskDetection) GetSource()(*string) {
    val, err := m.GetBackingStore().Get("source")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetTokenIssuerType gets the tokenIssuerType property value. Indicates the type of token issuer for the detected sign-in risk. The possible values are: AzureAD.
// returns a *TokenIssuerType when successful
func (m *ServicePrincipalRiskDetection) GetTokenIssuerType()(*TokenIssuerType) {
    val, err := m.GetBackingStore().Get("tokenIssuerType")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*TokenIssuerType)
    }
    return nil
}
// Serialize serializes information the current object
func (m *ServicePrincipalRiskDetection) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    err := m.Entity.Serialize(writer)
    if err != nil {
        return err
    }
    if m.GetActivity() != nil {
        cast := (*m.GetActivity()).String()
        err = writer.WriteStringValue("activity", &cast)
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteTimeValue("activityDateTime", m.GetActivityDateTime())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteStringValue("additionalInfo", m.GetAdditionalInfo())
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
        err = writer.WriteStringValue("correlationId", m.GetCorrelationId())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteTimeValue("detectedDateTime", m.GetDetectedDateTime())
        if err != nil {
            return err
        }
    }
    if m.GetDetectionTimingType() != nil {
        cast := (*m.GetDetectionTimingType()).String()
        err = writer.WriteStringValue("detectionTimingType", &cast)
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
    if m.GetKeyIds() != nil {
        err = writer.WriteCollectionOfStringValues("keyIds", m.GetKeyIds())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteTimeValue("lastUpdatedDateTime", m.GetLastUpdatedDateTime())
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
        err = writer.WriteStringValue("requestId", m.GetRequestId())
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
    {
        err = writer.WriteStringValue("riskEventType", m.GetRiskEventType())
        if err != nil {
            return err
        }
    }
    if m.GetRiskLevel() != nil {
        cast := (*m.GetRiskLevel()).String()
        err = writer.WriteStringValue("riskLevel", &cast)
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
        err = writer.WriteStringValue("servicePrincipalDisplayName", m.GetServicePrincipalDisplayName())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteStringValue("servicePrincipalId", m.GetServicePrincipalId())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteStringValue("source", m.GetSource())
        if err != nil {
            return err
        }
    }
    if m.GetTokenIssuerType() != nil {
        cast := (*m.GetTokenIssuerType()).String()
        err = writer.WriteStringValue("tokenIssuerType", &cast)
        if err != nil {
            return err
        }
    }
    return nil
}
// SetActivity sets the activity property value. Indicates the activity type the detected risk is linked to.  The possible values are: signin, servicePrincipal. Note that you must use the Prefer: include-unknown-enum-members request header to get the following value(s) in this evolvable enum: servicePrincipal.
func (m *ServicePrincipalRiskDetection) SetActivity(value *ActivityType)() {
    err := m.GetBackingStore().Set("activity", value)
    if err != nil {
        panic(err)
    }
}
// SetActivityDateTime sets the activityDateTime property value. Date and time when the risky activity occurred. The DateTimeOffset type represents date and time information using ISO 8601 format and is always in UTC time. For example, midnight UTC on Jan 1, 2014 is 2014-01-01T00:00:00Z
func (m *ServicePrincipalRiskDetection) SetActivityDateTime(value *i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)() {
    err := m.GetBackingStore().Set("activityDateTime", value)
    if err != nil {
        panic(err)
    }
}
// SetAdditionalInfo sets the additionalInfo property value. Additional information associated with the risk detection. This string value is represented as a JSON object with the quotations escaped.
func (m *ServicePrincipalRiskDetection) SetAdditionalInfo(value *string)() {
    err := m.GetBackingStore().Set("additionalInfo", value)
    if err != nil {
        panic(err)
    }
}
// SetAppId sets the appId property value. The unique identifier for the associated application.
func (m *ServicePrincipalRiskDetection) SetAppId(value *string)() {
    err := m.GetBackingStore().Set("appId", value)
    if err != nil {
        panic(err)
    }
}
// SetCorrelationId sets the correlationId property value. Correlation ID of the sign-in activity associated with the risk detection. This property is null if the risk detection is not associated with a sign-in activity.
func (m *ServicePrincipalRiskDetection) SetCorrelationId(value *string)() {
    err := m.GetBackingStore().Set("correlationId", value)
    if err != nil {
        panic(err)
    }
}
// SetDetectedDateTime sets the detectedDateTime property value. Date and time when the risk was detected. The DateTimeOffset type represents date and time information using ISO 8601 format and is always in UTC time. For example, midnight UTC on Jan 1, 2014 is 2014-01-01T00:00:00Z.
func (m *ServicePrincipalRiskDetection) SetDetectedDateTime(value *i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)() {
    err := m.GetBackingStore().Set("detectedDateTime", value)
    if err != nil {
        panic(err)
    }
}
// SetDetectionTimingType sets the detectionTimingType property value. Timing of the detected risk , whether real-time or offline. The possible values are: notDefined, realtime, nearRealtime, offline, unknownFutureValue.
func (m *ServicePrincipalRiskDetection) SetDetectionTimingType(value *RiskDetectionTimingType)() {
    err := m.GetBackingStore().Set("detectionTimingType", value)
    if err != nil {
        panic(err)
    }
}
// SetIpAddress sets the ipAddress property value. Provides the IP address of the client from where the risk occurred.
func (m *ServicePrincipalRiskDetection) SetIpAddress(value *string)() {
    err := m.GetBackingStore().Set("ipAddress", value)
    if err != nil {
        panic(err)
    }
}
// SetKeyIds sets the keyIds property value. The unique identifier for the key credential associated with the risk detection.
func (m *ServicePrincipalRiskDetection) SetKeyIds(value []string)() {
    err := m.GetBackingStore().Set("keyIds", value)
    if err != nil {
        panic(err)
    }
}
// SetLastUpdatedDateTime sets the lastUpdatedDateTime property value. Date and time when the risk detection was last updated.
func (m *ServicePrincipalRiskDetection) SetLastUpdatedDateTime(value *i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)() {
    err := m.GetBackingStore().Set("lastUpdatedDateTime", value)
    if err != nil {
        panic(err)
    }
}
// SetLocation sets the location property value. Location from where the sign-in was initiated.
func (m *ServicePrincipalRiskDetection) SetLocation(value SignInLocationable)() {
    err := m.GetBackingStore().Set("location", value)
    if err != nil {
        panic(err)
    }
}
// SetRequestId sets the requestId property value. Request identifier of the sign-in activity associated with the risk detection. This property is null if the risk detection is not associated with a sign-in activity. Supports $filter (eq).
func (m *ServicePrincipalRiskDetection) SetRequestId(value *string)() {
    err := m.GetBackingStore().Set("requestId", value)
    if err != nil {
        panic(err)
    }
}
// SetRiskDetail sets the riskDetail property value. Details of the detected risk. Note: Details for this property are only available for Workload Identities Premium customers. Events in tenants without this license will be returned hidden. The possible values are: none, hidden, adminConfirmedServicePrincipalCompromised, adminDismissedAllRiskForServicePrincipal. Note that you must use the Prefer: include-unknown-enum-members request header to get the following value(s) in this evolvable enum: adminConfirmedServicePrincipalCompromised , adminDismissedAllRiskForServicePrincipal.
func (m *ServicePrincipalRiskDetection) SetRiskDetail(value *RiskDetail)() {
    err := m.GetBackingStore().Set("riskDetail", value)
    if err != nil {
        panic(err)
    }
}
// SetRiskEventType sets the riskEventType property value. The type of risk event detected. The possible values are: investigationsThreatIntelligence, generic, adminConfirmedServicePrincipalCompromised, suspiciousSignins, leakedCredentials, anomalousServicePrincipalActivity, maliciousApplication, suspiciousApplication.
func (m *ServicePrincipalRiskDetection) SetRiskEventType(value *string)() {
    err := m.GetBackingStore().Set("riskEventType", value)
    if err != nil {
        panic(err)
    }
}
// SetRiskLevel sets the riskLevel property value. Level of the detected risk. Note: Details for this property are only available for Workload Identities Premium customers. Events in tenants without this license will be returned hidden. The possible values are: low, medium, high, hidden, none.
func (m *ServicePrincipalRiskDetection) SetRiskLevel(value *RiskLevel)() {
    err := m.GetBackingStore().Set("riskLevel", value)
    if err != nil {
        panic(err)
    }
}
// SetRiskState sets the riskState property value. The state of a detected risky service principal or sign-in activity. The possible values are: none, dismissed, atRisk, confirmedCompromised.
func (m *ServicePrincipalRiskDetection) SetRiskState(value *RiskState)() {
    err := m.GetBackingStore().Set("riskState", value)
    if err != nil {
        panic(err)
    }
}
// SetServicePrincipalDisplayName sets the servicePrincipalDisplayName property value. The display name for the service principal.
func (m *ServicePrincipalRiskDetection) SetServicePrincipalDisplayName(value *string)() {
    err := m.GetBackingStore().Set("servicePrincipalDisplayName", value)
    if err != nil {
        panic(err)
    }
}
// SetServicePrincipalId sets the servicePrincipalId property value. The unique identifier for the service principal. Supports $filter (eq).
func (m *ServicePrincipalRiskDetection) SetServicePrincipalId(value *string)() {
    err := m.GetBackingStore().Set("servicePrincipalId", value)
    if err != nil {
        panic(err)
    }
}
// SetSource sets the source property value. Source of the risk detection. For example, identityProtection.
func (m *ServicePrincipalRiskDetection) SetSource(value *string)() {
    err := m.GetBackingStore().Set("source", value)
    if err != nil {
        panic(err)
    }
}
// SetTokenIssuerType sets the tokenIssuerType property value. Indicates the type of token issuer for the detected sign-in risk. The possible values are: AzureAD.
func (m *ServicePrincipalRiskDetection) SetTokenIssuerType(value *TokenIssuerType)() {
    err := m.GetBackingStore().Set("tokenIssuerType", value)
    if err != nil {
        panic(err)
    }
}
type ServicePrincipalRiskDetectionable interface {
    Entityable
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable
    GetActivity()(*ActivityType)
    GetActivityDateTime()(*i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)
    GetAdditionalInfo()(*string)
    GetAppId()(*string)
    GetCorrelationId()(*string)
    GetDetectedDateTime()(*i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)
    GetDetectionTimingType()(*RiskDetectionTimingType)
    GetIpAddress()(*string)
    GetKeyIds()([]string)
    GetLastUpdatedDateTime()(*i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)
    GetLocation()(SignInLocationable)
    GetRequestId()(*string)
    GetRiskDetail()(*RiskDetail)
    GetRiskEventType()(*string)
    GetRiskLevel()(*RiskLevel)
    GetRiskState()(*RiskState)
    GetServicePrincipalDisplayName()(*string)
    GetServicePrincipalId()(*string)
    GetSource()(*string)
    GetTokenIssuerType()(*TokenIssuerType)
    SetActivity(value *ActivityType)()
    SetActivityDateTime(value *i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)()
    SetAdditionalInfo(value *string)()
    SetAppId(value *string)()
    SetCorrelationId(value *string)()
    SetDetectedDateTime(value *i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)()
    SetDetectionTimingType(value *RiskDetectionTimingType)()
    SetIpAddress(value *string)()
    SetKeyIds(value []string)()
    SetLastUpdatedDateTime(value *i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)()
    SetLocation(value SignInLocationable)()
    SetRequestId(value *string)()
    SetRiskDetail(value *RiskDetail)()
    SetRiskEventType(value *string)()
    SetRiskLevel(value *RiskLevel)()
    SetRiskState(value *RiskState)()
    SetServicePrincipalDisplayName(value *string)()
    SetServicePrincipalId(value *string)()
    SetSource(value *string)()
    SetTokenIssuerType(value *TokenIssuerType)()
}
