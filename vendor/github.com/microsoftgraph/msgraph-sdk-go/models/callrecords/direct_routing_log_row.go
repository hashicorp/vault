package callrecords

import (
    i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e "time"
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
    ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e "github.com/microsoft/kiota-abstractions-go/store"
)

type DirectRoutingLogRow struct {
    // Stores model information.
    backingStore ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackingStore
}
// NewDirectRoutingLogRow instantiates a new DirectRoutingLogRow and sets the default values.
func NewDirectRoutingLogRow()(*DirectRoutingLogRow) {
    m := &DirectRoutingLogRow{
    }
    m.backingStore = ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackingStoreFactoryInstance();
    m.SetAdditionalData(make(map[string]any))
    return m
}
// CreateDirectRoutingLogRowFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
// returns a Parsable when successful
func CreateDirectRoutingLogRowFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewDirectRoutingLogRow(), nil
}
// GetAdditionalData gets the AdditionalData property value. Stores additional data not described in the OpenAPI description found when deserializing. Can be used for serialization as well.
// returns a map[string]any when successful
func (m *DirectRoutingLogRow) GetAdditionalData()(map[string]any) {
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
func (m *DirectRoutingLogRow) GetBackingStore()(ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackingStore) {
    return m.backingStore
}
// GetCalleeNumber gets the calleeNumber property value. Number of the user or bot who received the call. E.164 format, but might include other data.
// returns a *string when successful
func (m *DirectRoutingLogRow) GetCalleeNumber()(*string) {
    val, err := m.GetBackingStore().Get("calleeNumber")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetCallEndSubReason gets the callEndSubReason property value. In addition to the SIP codes, Microsoft has subcodes that indicate the specific issue.
// returns a *int32 when successful
func (m *DirectRoutingLogRow) GetCallEndSubReason()(*int32) {
    val, err := m.GetBackingStore().Get("callEndSubReason")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*int32)
    }
    return nil
}
// GetCallerNumber gets the callerNumber property value. Number of the user or bot who made the call. E.164 format, but might include other data.
// returns a *string when successful
func (m *DirectRoutingLogRow) GetCallerNumber()(*string) {
    val, err := m.GetBackingStore().Get("callerNumber")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetCallType gets the callType property value. Call type and direction.
// returns a *string when successful
func (m *DirectRoutingLogRow) GetCallType()(*string) {
    val, err := m.GetBackingStore().Get("callType")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetCorrelationId gets the correlationId property value. Identifier for the call that you can use when calling Microsoft Support. GUID.
// returns a *string when successful
func (m *DirectRoutingLogRow) GetCorrelationId()(*string) {
    val, err := m.GetBackingStore().Get("correlationId")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetDuration gets the duration property value. Duration of the call in seconds.
// returns a *int32 when successful
func (m *DirectRoutingLogRow) GetDuration()(*int32) {
    val, err := m.GetBackingStore().Get("duration")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*int32)
    }
    return nil
}
// GetEndDateTime gets the endDateTime property value. Only exists for successful (fully established) calls. Time when call ended.
// returns a *Time when successful
func (m *DirectRoutingLogRow) GetEndDateTime()(*i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time) {
    val, err := m.GetBackingStore().Get("endDateTime")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)
    }
    return nil
}
// GetFailureDateTime gets the failureDateTime property value. Only exists for failed (not fully established) calls.
// returns a *Time when successful
func (m *DirectRoutingLogRow) GetFailureDateTime()(*i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time) {
    val, err := m.GetBackingStore().Get("failureDateTime")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)
    }
    return nil
}
// GetFieldDeserializers the deserialization information for the current model
// returns a map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error) when successful
func (m *DirectRoutingLogRow) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := make(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error))
    res["calleeNumber"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetCalleeNumber(val)
        }
        return nil
    }
    res["callEndSubReason"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetInt32Value()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetCallEndSubReason(val)
        }
        return nil
    }
    res["callerNumber"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetCallerNumber(val)
        }
        return nil
    }
    res["callType"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetCallType(val)
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
    res["duration"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetInt32Value()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetDuration(val)
        }
        return nil
    }
    res["endDateTime"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetTimeValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetEndDateTime(val)
        }
        return nil
    }
    res["failureDateTime"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetTimeValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetFailureDateTime(val)
        }
        return nil
    }
    res["finalSipCode"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetInt32Value()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetFinalSipCode(val)
        }
        return nil
    }
    res["finalSipCodePhrase"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetFinalSipCodePhrase(val)
        }
        return nil
    }
    res["id"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetId(val)
        }
        return nil
    }
    res["inviteDateTime"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetTimeValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetInviteDateTime(val)
        }
        return nil
    }
    res["mediaBypassEnabled"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetBoolValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetMediaBypassEnabled(val)
        }
        return nil
    }
    res["mediaPathLocation"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetMediaPathLocation(val)
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
    res["signalingLocation"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetSignalingLocation(val)
        }
        return nil
    }
    res["startDateTime"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetTimeValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetStartDateTime(val)
        }
        return nil
    }
    res["successfulCall"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetBoolValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetSuccessfulCall(val)
        }
        return nil
    }
    res["trunkFullyQualifiedDomainName"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetTrunkFullyQualifiedDomainName(val)
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
// GetFinalSipCode gets the finalSipCode property value. The final response code with which the call ended. For more information, see RFC 3261.
// returns a *int32 when successful
func (m *DirectRoutingLogRow) GetFinalSipCode()(*int32) {
    val, err := m.GetBackingStore().Get("finalSipCode")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*int32)
    }
    return nil
}
// GetFinalSipCodePhrase gets the finalSipCodePhrase property value. Description of the SIP code and Microsoft subcode.
// returns a *string when successful
func (m *DirectRoutingLogRow) GetFinalSipCodePhrase()(*string) {
    val, err := m.GetBackingStore().Get("finalSipCodePhrase")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetId gets the id property value. Unique call identifier. GUID.
// returns a *string when successful
func (m *DirectRoutingLogRow) GetId()(*string) {
    val, err := m.GetBackingStore().Get("id")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetInviteDateTime gets the inviteDateTime property value. The date and time when the initial invite was sent.
// returns a *Time when successful
func (m *DirectRoutingLogRow) GetInviteDateTime()(*i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time) {
    val, err := m.GetBackingStore().Get("inviteDateTime")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)
    }
    return nil
}
// GetMediaBypassEnabled gets the mediaBypassEnabled property value. Indicates whether the trunk was enabled for media bypass.
// returns a *bool when successful
func (m *DirectRoutingLogRow) GetMediaBypassEnabled()(*bool) {
    val, err := m.GetBackingStore().Get("mediaBypassEnabled")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*bool)
    }
    return nil
}
// GetMediaPathLocation gets the mediaPathLocation property value. The datacenter used for media path in a nonbypass call.
// returns a *string when successful
func (m *DirectRoutingLogRow) GetMediaPathLocation()(*string) {
    val, err := m.GetBackingStore().Get("mediaPathLocation")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetOdataType gets the @odata.type property value. The OdataType property
// returns a *string when successful
func (m *DirectRoutingLogRow) GetOdataType()(*string) {
    val, err := m.GetBackingStore().Get("odataType")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetSignalingLocation gets the signalingLocation property value. The datacenter used for signaling for both bypass and nonbypass calls.
// returns a *string when successful
func (m *DirectRoutingLogRow) GetSignalingLocation()(*string) {
    val, err := m.GetBackingStore().Get("signalingLocation")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetStartDateTime gets the startDateTime property value. Call start time.For failed and unanswered calls, this value can be equal to the invite or failure time.
// returns a *Time when successful
func (m *DirectRoutingLogRow) GetStartDateTime()(*i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time) {
    val, err := m.GetBackingStore().Get("startDateTime")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)
    }
    return nil
}
// GetSuccessfulCall gets the successfulCall property value. Success or attempt.
// returns a *bool when successful
func (m *DirectRoutingLogRow) GetSuccessfulCall()(*bool) {
    val, err := m.GetBackingStore().Get("successfulCall")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*bool)
    }
    return nil
}
// GetTrunkFullyQualifiedDomainName gets the trunkFullyQualifiedDomainName property value. Fully qualified domain name of the session border controller.
// returns a *string when successful
func (m *DirectRoutingLogRow) GetTrunkFullyQualifiedDomainName()(*string) {
    val, err := m.GetBackingStore().Get("trunkFullyQualifiedDomainName")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetUserDisplayName gets the userDisplayName property value. Display name of the user.
// returns a *string when successful
func (m *DirectRoutingLogRow) GetUserDisplayName()(*string) {
    val, err := m.GetBackingStore().Get("userDisplayName")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetUserId gets the userId property value. Calling user's ID in Microsoft Graph. This and other user information is null/empty for bot call types. GUID.
// returns a *string when successful
func (m *DirectRoutingLogRow) GetUserId()(*string) {
    val, err := m.GetBackingStore().Get("userId")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetUserPrincipalName gets the userPrincipalName property value. UserPrincipalName (sign-in name) in Microsoft Entra ID. This value is usually the same as the user's SIP Address, and can be the same as the user's email address.
// returns a *string when successful
func (m *DirectRoutingLogRow) GetUserPrincipalName()(*string) {
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
func (m *DirectRoutingLogRow) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    {
        err := writer.WriteStringValue("calleeNumber", m.GetCalleeNumber())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteInt32Value("callEndSubReason", m.GetCallEndSubReason())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteStringValue("callerNumber", m.GetCallerNumber())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteStringValue("callType", m.GetCallType())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteStringValue("correlationId", m.GetCorrelationId())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteInt32Value("duration", m.GetDuration())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteTimeValue("endDateTime", m.GetEndDateTime())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteTimeValue("failureDateTime", m.GetFailureDateTime())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteInt32Value("finalSipCode", m.GetFinalSipCode())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteStringValue("finalSipCodePhrase", m.GetFinalSipCodePhrase())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteStringValue("id", m.GetId())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteTimeValue("inviteDateTime", m.GetInviteDateTime())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteBoolValue("mediaBypassEnabled", m.GetMediaBypassEnabled())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteStringValue("mediaPathLocation", m.GetMediaPathLocation())
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
        err := writer.WriteStringValue("signalingLocation", m.GetSignalingLocation())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteTimeValue("startDateTime", m.GetStartDateTime())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteBoolValue("successfulCall", m.GetSuccessfulCall())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteStringValue("trunkFullyQualifiedDomainName", m.GetTrunkFullyQualifiedDomainName())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteStringValue("userDisplayName", m.GetUserDisplayName())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteStringValue("userId", m.GetUserId())
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
// SetAdditionalData sets the AdditionalData property value. Stores additional data not described in the OpenAPI description found when deserializing. Can be used for serialization as well.
func (m *DirectRoutingLogRow) SetAdditionalData(value map[string]any)() {
    err := m.GetBackingStore().Set("additionalData", value)
    if err != nil {
        panic(err)
    }
}
// SetBackingStore sets the BackingStore property value. Stores model information.
func (m *DirectRoutingLogRow) SetBackingStore(value ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackingStore)() {
    m.backingStore = value
}
// SetCalleeNumber sets the calleeNumber property value. Number of the user or bot who received the call. E.164 format, but might include other data.
func (m *DirectRoutingLogRow) SetCalleeNumber(value *string)() {
    err := m.GetBackingStore().Set("calleeNumber", value)
    if err != nil {
        panic(err)
    }
}
// SetCallEndSubReason sets the callEndSubReason property value. In addition to the SIP codes, Microsoft has subcodes that indicate the specific issue.
func (m *DirectRoutingLogRow) SetCallEndSubReason(value *int32)() {
    err := m.GetBackingStore().Set("callEndSubReason", value)
    if err != nil {
        panic(err)
    }
}
// SetCallerNumber sets the callerNumber property value. Number of the user or bot who made the call. E.164 format, but might include other data.
func (m *DirectRoutingLogRow) SetCallerNumber(value *string)() {
    err := m.GetBackingStore().Set("callerNumber", value)
    if err != nil {
        panic(err)
    }
}
// SetCallType sets the callType property value. Call type and direction.
func (m *DirectRoutingLogRow) SetCallType(value *string)() {
    err := m.GetBackingStore().Set("callType", value)
    if err != nil {
        panic(err)
    }
}
// SetCorrelationId sets the correlationId property value. Identifier for the call that you can use when calling Microsoft Support. GUID.
func (m *DirectRoutingLogRow) SetCorrelationId(value *string)() {
    err := m.GetBackingStore().Set("correlationId", value)
    if err != nil {
        panic(err)
    }
}
// SetDuration sets the duration property value. Duration of the call in seconds.
func (m *DirectRoutingLogRow) SetDuration(value *int32)() {
    err := m.GetBackingStore().Set("duration", value)
    if err != nil {
        panic(err)
    }
}
// SetEndDateTime sets the endDateTime property value. Only exists for successful (fully established) calls. Time when call ended.
func (m *DirectRoutingLogRow) SetEndDateTime(value *i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)() {
    err := m.GetBackingStore().Set("endDateTime", value)
    if err != nil {
        panic(err)
    }
}
// SetFailureDateTime sets the failureDateTime property value. Only exists for failed (not fully established) calls.
func (m *DirectRoutingLogRow) SetFailureDateTime(value *i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)() {
    err := m.GetBackingStore().Set("failureDateTime", value)
    if err != nil {
        panic(err)
    }
}
// SetFinalSipCode sets the finalSipCode property value. The final response code with which the call ended. For more information, see RFC 3261.
func (m *DirectRoutingLogRow) SetFinalSipCode(value *int32)() {
    err := m.GetBackingStore().Set("finalSipCode", value)
    if err != nil {
        panic(err)
    }
}
// SetFinalSipCodePhrase sets the finalSipCodePhrase property value. Description of the SIP code and Microsoft subcode.
func (m *DirectRoutingLogRow) SetFinalSipCodePhrase(value *string)() {
    err := m.GetBackingStore().Set("finalSipCodePhrase", value)
    if err != nil {
        panic(err)
    }
}
// SetId sets the id property value. Unique call identifier. GUID.
func (m *DirectRoutingLogRow) SetId(value *string)() {
    err := m.GetBackingStore().Set("id", value)
    if err != nil {
        panic(err)
    }
}
// SetInviteDateTime sets the inviteDateTime property value. The date and time when the initial invite was sent.
func (m *DirectRoutingLogRow) SetInviteDateTime(value *i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)() {
    err := m.GetBackingStore().Set("inviteDateTime", value)
    if err != nil {
        panic(err)
    }
}
// SetMediaBypassEnabled sets the mediaBypassEnabled property value. Indicates whether the trunk was enabled for media bypass.
func (m *DirectRoutingLogRow) SetMediaBypassEnabled(value *bool)() {
    err := m.GetBackingStore().Set("mediaBypassEnabled", value)
    if err != nil {
        panic(err)
    }
}
// SetMediaPathLocation sets the mediaPathLocation property value. The datacenter used for media path in a nonbypass call.
func (m *DirectRoutingLogRow) SetMediaPathLocation(value *string)() {
    err := m.GetBackingStore().Set("mediaPathLocation", value)
    if err != nil {
        panic(err)
    }
}
// SetOdataType sets the @odata.type property value. The OdataType property
func (m *DirectRoutingLogRow) SetOdataType(value *string)() {
    err := m.GetBackingStore().Set("odataType", value)
    if err != nil {
        panic(err)
    }
}
// SetSignalingLocation sets the signalingLocation property value. The datacenter used for signaling for both bypass and nonbypass calls.
func (m *DirectRoutingLogRow) SetSignalingLocation(value *string)() {
    err := m.GetBackingStore().Set("signalingLocation", value)
    if err != nil {
        panic(err)
    }
}
// SetStartDateTime sets the startDateTime property value. Call start time.For failed and unanswered calls, this value can be equal to the invite or failure time.
func (m *DirectRoutingLogRow) SetStartDateTime(value *i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)() {
    err := m.GetBackingStore().Set("startDateTime", value)
    if err != nil {
        panic(err)
    }
}
// SetSuccessfulCall sets the successfulCall property value. Success or attempt.
func (m *DirectRoutingLogRow) SetSuccessfulCall(value *bool)() {
    err := m.GetBackingStore().Set("successfulCall", value)
    if err != nil {
        panic(err)
    }
}
// SetTrunkFullyQualifiedDomainName sets the trunkFullyQualifiedDomainName property value. Fully qualified domain name of the session border controller.
func (m *DirectRoutingLogRow) SetTrunkFullyQualifiedDomainName(value *string)() {
    err := m.GetBackingStore().Set("trunkFullyQualifiedDomainName", value)
    if err != nil {
        panic(err)
    }
}
// SetUserDisplayName sets the userDisplayName property value. Display name of the user.
func (m *DirectRoutingLogRow) SetUserDisplayName(value *string)() {
    err := m.GetBackingStore().Set("userDisplayName", value)
    if err != nil {
        panic(err)
    }
}
// SetUserId sets the userId property value. Calling user's ID in Microsoft Graph. This and other user information is null/empty for bot call types. GUID.
func (m *DirectRoutingLogRow) SetUserId(value *string)() {
    err := m.GetBackingStore().Set("userId", value)
    if err != nil {
        panic(err)
    }
}
// SetUserPrincipalName sets the userPrincipalName property value. UserPrincipalName (sign-in name) in Microsoft Entra ID. This value is usually the same as the user's SIP Address, and can be the same as the user's email address.
func (m *DirectRoutingLogRow) SetUserPrincipalName(value *string)() {
    err := m.GetBackingStore().Set("userPrincipalName", value)
    if err != nil {
        panic(err)
    }
}
type DirectRoutingLogRowable interface {
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.AdditionalDataHolder
    ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackedModel
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable
    GetBackingStore()(ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackingStore)
    GetCalleeNumber()(*string)
    GetCallEndSubReason()(*int32)
    GetCallerNumber()(*string)
    GetCallType()(*string)
    GetCorrelationId()(*string)
    GetDuration()(*int32)
    GetEndDateTime()(*i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)
    GetFailureDateTime()(*i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)
    GetFinalSipCode()(*int32)
    GetFinalSipCodePhrase()(*string)
    GetId()(*string)
    GetInviteDateTime()(*i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)
    GetMediaBypassEnabled()(*bool)
    GetMediaPathLocation()(*string)
    GetOdataType()(*string)
    GetSignalingLocation()(*string)
    GetStartDateTime()(*i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)
    GetSuccessfulCall()(*bool)
    GetTrunkFullyQualifiedDomainName()(*string)
    GetUserDisplayName()(*string)
    GetUserId()(*string)
    GetUserPrincipalName()(*string)
    SetBackingStore(value ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackingStore)()
    SetCalleeNumber(value *string)()
    SetCallEndSubReason(value *int32)()
    SetCallerNumber(value *string)()
    SetCallType(value *string)()
    SetCorrelationId(value *string)()
    SetDuration(value *int32)()
    SetEndDateTime(value *i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)()
    SetFailureDateTime(value *i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)()
    SetFinalSipCode(value *int32)()
    SetFinalSipCodePhrase(value *string)()
    SetId(value *string)()
    SetInviteDateTime(value *i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)()
    SetMediaBypassEnabled(value *bool)()
    SetMediaPathLocation(value *string)()
    SetOdataType(value *string)()
    SetSignalingLocation(value *string)()
    SetStartDateTime(value *i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)()
    SetSuccessfulCall(value *bool)()
    SetTrunkFullyQualifiedDomainName(value *string)()
    SetUserDisplayName(value *string)()
    SetUserId(value *string)()
    SetUserPrincipalName(value *string)()
}
