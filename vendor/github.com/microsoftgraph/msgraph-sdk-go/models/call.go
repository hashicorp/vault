package models

import (
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

type Call struct {
    Entity
}
// NewCall instantiates a new Call and sets the default values.
func NewCall()(*Call) {
    m := &Call{
        Entity: *NewEntity(),
    }
    return m
}
// CreateCallFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
// returns a Parsable when successful
func CreateCallFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewCall(), nil
}
// GetAudioRoutingGroups gets the audioRoutingGroups property value. The audioRoutingGroups property
// returns a []AudioRoutingGroupable when successful
func (m *Call) GetAudioRoutingGroups()([]AudioRoutingGroupable) {
    val, err := m.GetBackingStore().Get("audioRoutingGroups")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]AudioRoutingGroupable)
    }
    return nil
}
// GetCallbackUri gets the callbackUri property value. The callback URL on which callbacks are delivered. Must be an HTTPS URL.
// returns a *string when successful
func (m *Call) GetCallbackUri()(*string) {
    val, err := m.GetBackingStore().Get("callbackUri")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetCallChainId gets the callChainId property value. A unique identifier for all the participant calls in a conference or a unique identifier for two participant calls in a P2P call.  This identifier must be copied over from Microsoft.Graph.Call.CallChainId.
// returns a *string when successful
func (m *Call) GetCallChainId()(*string) {
    val, err := m.GetBackingStore().Get("callChainId")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetCallOptions gets the callOptions property value. Contains the optional features for the call.
// returns a CallOptionsable when successful
func (m *Call) GetCallOptions()(CallOptionsable) {
    val, err := m.GetBackingStore().Get("callOptions")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(CallOptionsable)
    }
    return nil
}
// GetCallRoutes gets the callRoutes property value. The routing information on how the call was retargeted. Read-only.
// returns a []CallRouteable when successful
func (m *Call) GetCallRoutes()([]CallRouteable) {
    val, err := m.GetBackingStore().Get("callRoutes")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]CallRouteable)
    }
    return nil
}
// GetChatInfo gets the chatInfo property value. The chat information. Required information for joining a meeting.
// returns a ChatInfoable when successful
func (m *Call) GetChatInfo()(ChatInfoable) {
    val, err := m.GetBackingStore().Get("chatInfo")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(ChatInfoable)
    }
    return nil
}
// GetContentSharingSessions gets the contentSharingSessions property value. The contentSharingSessions property
// returns a []ContentSharingSessionable when successful
func (m *Call) GetContentSharingSessions()([]ContentSharingSessionable) {
    val, err := m.GetBackingStore().Get("contentSharingSessions")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]ContentSharingSessionable)
    }
    return nil
}
// GetDirection gets the direction property value. The direction of the call. The possible values are incoming or outgoing. Read-only.
// returns a *CallDirection when successful
func (m *Call) GetDirection()(*CallDirection) {
    val, err := m.GetBackingStore().Get("direction")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*CallDirection)
    }
    return nil
}
// GetFieldDeserializers the deserialization information for the current model
// returns a map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error) when successful
func (m *Call) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := m.Entity.GetFieldDeserializers()
    res["audioRoutingGroups"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetCollectionOfObjectValues(CreateAudioRoutingGroupFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            res := make([]AudioRoutingGroupable, len(val))
            for i, v := range val {
                if v != nil {
                    res[i] = v.(AudioRoutingGroupable)
                }
            }
            m.SetAudioRoutingGroups(res)
        }
        return nil
    }
    res["callbackUri"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetCallbackUri(val)
        }
        return nil
    }
    res["callChainId"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetCallChainId(val)
        }
        return nil
    }
    res["callOptions"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetObjectValue(CreateCallOptionsFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetCallOptions(val.(CallOptionsable))
        }
        return nil
    }
    res["callRoutes"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetCollectionOfObjectValues(CreateCallRouteFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            res := make([]CallRouteable, len(val))
            for i, v := range val {
                if v != nil {
                    res[i] = v.(CallRouteable)
                }
            }
            m.SetCallRoutes(res)
        }
        return nil
    }
    res["chatInfo"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetObjectValue(CreateChatInfoFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetChatInfo(val.(ChatInfoable))
        }
        return nil
    }
    res["contentSharingSessions"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetCollectionOfObjectValues(CreateContentSharingSessionFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            res := make([]ContentSharingSessionable, len(val))
            for i, v := range val {
                if v != nil {
                    res[i] = v.(ContentSharingSessionable)
                }
            }
            m.SetContentSharingSessions(res)
        }
        return nil
    }
    res["direction"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetEnumValue(ParseCallDirection)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetDirection(val.(*CallDirection))
        }
        return nil
    }
    res["incomingContext"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetObjectValue(CreateIncomingContextFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetIncomingContext(val.(IncomingContextable))
        }
        return nil
    }
    res["mediaConfig"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetObjectValue(CreateMediaConfigFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetMediaConfig(val.(MediaConfigable))
        }
        return nil
    }
    res["mediaState"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetObjectValue(CreateCallMediaStateFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetMediaState(val.(CallMediaStateable))
        }
        return nil
    }
    res["meetingInfo"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetObjectValue(CreateMeetingInfoFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetMeetingInfo(val.(MeetingInfoable))
        }
        return nil
    }
    res["myParticipantId"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetMyParticipantId(val)
        }
        return nil
    }
    res["operations"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetCollectionOfObjectValues(CreateCommsOperationFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            res := make([]CommsOperationable, len(val))
            for i, v := range val {
                if v != nil {
                    res[i] = v.(CommsOperationable)
                }
            }
            m.SetOperations(res)
        }
        return nil
    }
    res["participants"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetCollectionOfObjectValues(CreateParticipantFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            res := make([]Participantable, len(val))
            for i, v := range val {
                if v != nil {
                    res[i] = v.(Participantable)
                }
            }
            m.SetParticipants(res)
        }
        return nil
    }
    res["requestedModalities"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetCollectionOfEnumValues(ParseModality)
        if err != nil {
            return err
        }
        if val != nil {
            res := make([]Modality, len(val))
            for i, v := range val {
                if v != nil {
                    res[i] = *(v.(*Modality))
                }
            }
            m.SetRequestedModalities(res)
        }
        return nil
    }
    res["resultInfo"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetObjectValue(CreateResultInfoFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetResultInfo(val.(ResultInfoable))
        }
        return nil
    }
    res["source"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetObjectValue(CreateParticipantInfoFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetSource(val.(ParticipantInfoable))
        }
        return nil
    }
    res["state"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetEnumValue(ParseCallState)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetState(val.(*CallState))
        }
        return nil
    }
    res["subject"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetSubject(val)
        }
        return nil
    }
    res["targets"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetCollectionOfObjectValues(CreateInvitationParticipantInfoFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            res := make([]InvitationParticipantInfoable, len(val))
            for i, v := range val {
                if v != nil {
                    res[i] = v.(InvitationParticipantInfoable)
                }
            }
            m.SetTargets(res)
        }
        return nil
    }
    res["tenantId"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetTenantId(val)
        }
        return nil
    }
    res["toneInfo"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetObjectValue(CreateToneInfoFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetToneInfo(val.(ToneInfoable))
        }
        return nil
    }
    res["transcription"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetObjectValue(CreateCallTranscriptionInfoFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetTranscription(val.(CallTranscriptionInfoable))
        }
        return nil
    }
    return res
}
// GetIncomingContext gets the incomingContext property value. Call context associated with an incoming call.
// returns a IncomingContextable when successful
func (m *Call) GetIncomingContext()(IncomingContextable) {
    val, err := m.GetBackingStore().Get("incomingContext")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(IncomingContextable)
    }
    return nil
}
// GetMediaConfig gets the mediaConfig property value. The media configuration. Required.
// returns a MediaConfigable when successful
func (m *Call) GetMediaConfig()(MediaConfigable) {
    val, err := m.GetBackingStore().Get("mediaConfig")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(MediaConfigable)
    }
    return nil
}
// GetMediaState gets the mediaState property value. Read-only. The call media state.
// returns a CallMediaStateable when successful
func (m *Call) GetMediaState()(CallMediaStateable) {
    val, err := m.GetBackingStore().Get("mediaState")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(CallMediaStateable)
    }
    return nil
}
// GetMeetingInfo gets the meetingInfo property value. The meeting information. Required information for meeting scenarios.
// returns a MeetingInfoable when successful
func (m *Call) GetMeetingInfo()(MeetingInfoable) {
    val, err := m.GetBackingStore().Get("meetingInfo")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(MeetingInfoable)
    }
    return nil
}
// GetMyParticipantId gets the myParticipantId property value. The myParticipantId property
// returns a *string when successful
func (m *Call) GetMyParticipantId()(*string) {
    val, err := m.GetBackingStore().Get("myParticipantId")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetOperations gets the operations property value. The operations property
// returns a []CommsOperationable when successful
func (m *Call) GetOperations()([]CommsOperationable) {
    val, err := m.GetBackingStore().Get("operations")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]CommsOperationable)
    }
    return nil
}
// GetParticipants gets the participants property value. The participants property
// returns a []Participantable when successful
func (m *Call) GetParticipants()([]Participantable) {
    val, err := m.GetBackingStore().Get("participants")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]Participantable)
    }
    return nil
}
// GetRequestedModalities gets the requestedModalities property value. The list of requested modalities. Possible values are: unknown, audio, video, videoBasedScreenSharing, data.
// returns a []Modality when successful
func (m *Call) GetRequestedModalities()([]Modality) {
    val, err := m.GetBackingStore().Get("requestedModalities")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]Modality)
    }
    return nil
}
// GetResultInfo gets the resultInfo property value. The result information. For example, the result can hold termination reason. Read-only.
// returns a ResultInfoable when successful
func (m *Call) GetResultInfo()(ResultInfoable) {
    val, err := m.GetBackingStore().Get("resultInfo")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(ResultInfoable)
    }
    return nil
}
// GetSource gets the source property value. The originator of the call.
// returns a ParticipantInfoable when successful
func (m *Call) GetSource()(ParticipantInfoable) {
    val, err := m.GetBackingStore().Get("source")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(ParticipantInfoable)
    }
    return nil
}
// GetState gets the state property value. The call state. Possible values are: incoming, establishing, ringing, established, hold, transferring, transferAccepted, redirecting, terminating, terminated. Read-only.
// returns a *CallState when successful
func (m *Call) GetState()(*CallState) {
    val, err := m.GetBackingStore().Get("state")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*CallState)
    }
    return nil
}
// GetSubject gets the subject property value. The subject of the conversation.
// returns a *string when successful
func (m *Call) GetSubject()(*string) {
    val, err := m.GetBackingStore().Get("subject")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetTargets gets the targets property value. The targets of the call. Required information for creating peer to peer call.
// returns a []InvitationParticipantInfoable when successful
func (m *Call) GetTargets()([]InvitationParticipantInfoable) {
    val, err := m.GetBackingStore().Get("targets")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]InvitationParticipantInfoable)
    }
    return nil
}
// GetTenantId gets the tenantId property value. The tenantId property
// returns a *string when successful
func (m *Call) GetTenantId()(*string) {
    val, err := m.GetBackingStore().Get("tenantId")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetToneInfo gets the toneInfo property value. The toneInfo property
// returns a ToneInfoable when successful
func (m *Call) GetToneInfo()(ToneInfoable) {
    val, err := m.GetBackingStore().Get("toneInfo")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(ToneInfoable)
    }
    return nil
}
// GetTranscription gets the transcription property value. The transcription information for the call. Read-only.
// returns a CallTranscriptionInfoable when successful
func (m *Call) GetTranscription()(CallTranscriptionInfoable) {
    val, err := m.GetBackingStore().Get("transcription")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(CallTranscriptionInfoable)
    }
    return nil
}
// Serialize serializes information the current object
func (m *Call) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    err := m.Entity.Serialize(writer)
    if err != nil {
        return err
    }
    if m.GetAudioRoutingGroups() != nil {
        cast := make([]i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, len(m.GetAudioRoutingGroups()))
        for i, v := range m.GetAudioRoutingGroups() {
            if v != nil {
                cast[i] = v.(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable)
            }
        }
        err = writer.WriteCollectionOfObjectValues("audioRoutingGroups", cast)
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteStringValue("callbackUri", m.GetCallbackUri())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteStringValue("callChainId", m.GetCallChainId())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteObjectValue("callOptions", m.GetCallOptions())
        if err != nil {
            return err
        }
    }
    if m.GetCallRoutes() != nil {
        cast := make([]i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, len(m.GetCallRoutes()))
        for i, v := range m.GetCallRoutes() {
            if v != nil {
                cast[i] = v.(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable)
            }
        }
        err = writer.WriteCollectionOfObjectValues("callRoutes", cast)
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteObjectValue("chatInfo", m.GetChatInfo())
        if err != nil {
            return err
        }
    }
    if m.GetContentSharingSessions() != nil {
        cast := make([]i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, len(m.GetContentSharingSessions()))
        for i, v := range m.GetContentSharingSessions() {
            if v != nil {
                cast[i] = v.(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable)
            }
        }
        err = writer.WriteCollectionOfObjectValues("contentSharingSessions", cast)
        if err != nil {
            return err
        }
    }
    if m.GetDirection() != nil {
        cast := (*m.GetDirection()).String()
        err = writer.WriteStringValue("direction", &cast)
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteObjectValue("incomingContext", m.GetIncomingContext())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteObjectValue("mediaConfig", m.GetMediaConfig())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteObjectValue("mediaState", m.GetMediaState())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteObjectValue("meetingInfo", m.GetMeetingInfo())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteStringValue("myParticipantId", m.GetMyParticipantId())
        if err != nil {
            return err
        }
    }
    if m.GetOperations() != nil {
        cast := make([]i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, len(m.GetOperations()))
        for i, v := range m.GetOperations() {
            if v != nil {
                cast[i] = v.(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable)
            }
        }
        err = writer.WriteCollectionOfObjectValues("operations", cast)
        if err != nil {
            return err
        }
    }
    if m.GetParticipants() != nil {
        cast := make([]i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, len(m.GetParticipants()))
        for i, v := range m.GetParticipants() {
            if v != nil {
                cast[i] = v.(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable)
            }
        }
        err = writer.WriteCollectionOfObjectValues("participants", cast)
        if err != nil {
            return err
        }
    }
    if m.GetRequestedModalities() != nil {
        err = writer.WriteCollectionOfStringValues("requestedModalities", SerializeModality(m.GetRequestedModalities()))
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteObjectValue("resultInfo", m.GetResultInfo())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteObjectValue("source", m.GetSource())
        if err != nil {
            return err
        }
    }
    if m.GetState() != nil {
        cast := (*m.GetState()).String()
        err = writer.WriteStringValue("state", &cast)
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteStringValue("subject", m.GetSubject())
        if err != nil {
            return err
        }
    }
    if m.GetTargets() != nil {
        cast := make([]i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, len(m.GetTargets()))
        for i, v := range m.GetTargets() {
            if v != nil {
                cast[i] = v.(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable)
            }
        }
        err = writer.WriteCollectionOfObjectValues("targets", cast)
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteStringValue("tenantId", m.GetTenantId())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteObjectValue("toneInfo", m.GetToneInfo())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteObjectValue("transcription", m.GetTranscription())
        if err != nil {
            return err
        }
    }
    return nil
}
// SetAudioRoutingGroups sets the audioRoutingGroups property value. The audioRoutingGroups property
func (m *Call) SetAudioRoutingGroups(value []AudioRoutingGroupable)() {
    err := m.GetBackingStore().Set("audioRoutingGroups", value)
    if err != nil {
        panic(err)
    }
}
// SetCallbackUri sets the callbackUri property value. The callback URL on which callbacks are delivered. Must be an HTTPS URL.
func (m *Call) SetCallbackUri(value *string)() {
    err := m.GetBackingStore().Set("callbackUri", value)
    if err != nil {
        panic(err)
    }
}
// SetCallChainId sets the callChainId property value. A unique identifier for all the participant calls in a conference or a unique identifier for two participant calls in a P2P call.  This identifier must be copied over from Microsoft.Graph.Call.CallChainId.
func (m *Call) SetCallChainId(value *string)() {
    err := m.GetBackingStore().Set("callChainId", value)
    if err != nil {
        panic(err)
    }
}
// SetCallOptions sets the callOptions property value. Contains the optional features for the call.
func (m *Call) SetCallOptions(value CallOptionsable)() {
    err := m.GetBackingStore().Set("callOptions", value)
    if err != nil {
        panic(err)
    }
}
// SetCallRoutes sets the callRoutes property value. The routing information on how the call was retargeted. Read-only.
func (m *Call) SetCallRoutes(value []CallRouteable)() {
    err := m.GetBackingStore().Set("callRoutes", value)
    if err != nil {
        panic(err)
    }
}
// SetChatInfo sets the chatInfo property value. The chat information. Required information for joining a meeting.
func (m *Call) SetChatInfo(value ChatInfoable)() {
    err := m.GetBackingStore().Set("chatInfo", value)
    if err != nil {
        panic(err)
    }
}
// SetContentSharingSessions sets the contentSharingSessions property value. The contentSharingSessions property
func (m *Call) SetContentSharingSessions(value []ContentSharingSessionable)() {
    err := m.GetBackingStore().Set("contentSharingSessions", value)
    if err != nil {
        panic(err)
    }
}
// SetDirection sets the direction property value. The direction of the call. The possible values are incoming or outgoing. Read-only.
func (m *Call) SetDirection(value *CallDirection)() {
    err := m.GetBackingStore().Set("direction", value)
    if err != nil {
        panic(err)
    }
}
// SetIncomingContext sets the incomingContext property value. Call context associated with an incoming call.
func (m *Call) SetIncomingContext(value IncomingContextable)() {
    err := m.GetBackingStore().Set("incomingContext", value)
    if err != nil {
        panic(err)
    }
}
// SetMediaConfig sets the mediaConfig property value. The media configuration. Required.
func (m *Call) SetMediaConfig(value MediaConfigable)() {
    err := m.GetBackingStore().Set("mediaConfig", value)
    if err != nil {
        panic(err)
    }
}
// SetMediaState sets the mediaState property value. Read-only. The call media state.
func (m *Call) SetMediaState(value CallMediaStateable)() {
    err := m.GetBackingStore().Set("mediaState", value)
    if err != nil {
        panic(err)
    }
}
// SetMeetingInfo sets the meetingInfo property value. The meeting information. Required information for meeting scenarios.
func (m *Call) SetMeetingInfo(value MeetingInfoable)() {
    err := m.GetBackingStore().Set("meetingInfo", value)
    if err != nil {
        panic(err)
    }
}
// SetMyParticipantId sets the myParticipantId property value. The myParticipantId property
func (m *Call) SetMyParticipantId(value *string)() {
    err := m.GetBackingStore().Set("myParticipantId", value)
    if err != nil {
        panic(err)
    }
}
// SetOperations sets the operations property value. The operations property
func (m *Call) SetOperations(value []CommsOperationable)() {
    err := m.GetBackingStore().Set("operations", value)
    if err != nil {
        panic(err)
    }
}
// SetParticipants sets the participants property value. The participants property
func (m *Call) SetParticipants(value []Participantable)() {
    err := m.GetBackingStore().Set("participants", value)
    if err != nil {
        panic(err)
    }
}
// SetRequestedModalities sets the requestedModalities property value. The list of requested modalities. Possible values are: unknown, audio, video, videoBasedScreenSharing, data.
func (m *Call) SetRequestedModalities(value []Modality)() {
    err := m.GetBackingStore().Set("requestedModalities", value)
    if err != nil {
        panic(err)
    }
}
// SetResultInfo sets the resultInfo property value. The result information. For example, the result can hold termination reason. Read-only.
func (m *Call) SetResultInfo(value ResultInfoable)() {
    err := m.GetBackingStore().Set("resultInfo", value)
    if err != nil {
        panic(err)
    }
}
// SetSource sets the source property value. The originator of the call.
func (m *Call) SetSource(value ParticipantInfoable)() {
    err := m.GetBackingStore().Set("source", value)
    if err != nil {
        panic(err)
    }
}
// SetState sets the state property value. The call state. Possible values are: incoming, establishing, ringing, established, hold, transferring, transferAccepted, redirecting, terminating, terminated. Read-only.
func (m *Call) SetState(value *CallState)() {
    err := m.GetBackingStore().Set("state", value)
    if err != nil {
        panic(err)
    }
}
// SetSubject sets the subject property value. The subject of the conversation.
func (m *Call) SetSubject(value *string)() {
    err := m.GetBackingStore().Set("subject", value)
    if err != nil {
        panic(err)
    }
}
// SetTargets sets the targets property value. The targets of the call. Required information for creating peer to peer call.
func (m *Call) SetTargets(value []InvitationParticipantInfoable)() {
    err := m.GetBackingStore().Set("targets", value)
    if err != nil {
        panic(err)
    }
}
// SetTenantId sets the tenantId property value. The tenantId property
func (m *Call) SetTenantId(value *string)() {
    err := m.GetBackingStore().Set("tenantId", value)
    if err != nil {
        panic(err)
    }
}
// SetToneInfo sets the toneInfo property value. The toneInfo property
func (m *Call) SetToneInfo(value ToneInfoable)() {
    err := m.GetBackingStore().Set("toneInfo", value)
    if err != nil {
        panic(err)
    }
}
// SetTranscription sets the transcription property value. The transcription information for the call. Read-only.
func (m *Call) SetTranscription(value CallTranscriptionInfoable)() {
    err := m.GetBackingStore().Set("transcription", value)
    if err != nil {
        panic(err)
    }
}
type Callable interface {
    Entityable
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable
    GetAudioRoutingGroups()([]AudioRoutingGroupable)
    GetCallbackUri()(*string)
    GetCallChainId()(*string)
    GetCallOptions()(CallOptionsable)
    GetCallRoutes()([]CallRouteable)
    GetChatInfo()(ChatInfoable)
    GetContentSharingSessions()([]ContentSharingSessionable)
    GetDirection()(*CallDirection)
    GetIncomingContext()(IncomingContextable)
    GetMediaConfig()(MediaConfigable)
    GetMediaState()(CallMediaStateable)
    GetMeetingInfo()(MeetingInfoable)
    GetMyParticipantId()(*string)
    GetOperations()([]CommsOperationable)
    GetParticipants()([]Participantable)
    GetRequestedModalities()([]Modality)
    GetResultInfo()(ResultInfoable)
    GetSource()(ParticipantInfoable)
    GetState()(*CallState)
    GetSubject()(*string)
    GetTargets()([]InvitationParticipantInfoable)
    GetTenantId()(*string)
    GetToneInfo()(ToneInfoable)
    GetTranscription()(CallTranscriptionInfoable)
    SetAudioRoutingGroups(value []AudioRoutingGroupable)()
    SetCallbackUri(value *string)()
    SetCallChainId(value *string)()
    SetCallOptions(value CallOptionsable)()
    SetCallRoutes(value []CallRouteable)()
    SetChatInfo(value ChatInfoable)()
    SetContentSharingSessions(value []ContentSharingSessionable)()
    SetDirection(value *CallDirection)()
    SetIncomingContext(value IncomingContextable)()
    SetMediaConfig(value MediaConfigable)()
    SetMediaState(value CallMediaStateable)()
    SetMeetingInfo(value MeetingInfoable)()
    SetMyParticipantId(value *string)()
    SetOperations(value []CommsOperationable)()
    SetParticipants(value []Participantable)()
    SetRequestedModalities(value []Modality)()
    SetResultInfo(value ResultInfoable)()
    SetSource(value ParticipantInfoable)()
    SetState(value *CallState)()
    SetSubject(value *string)()
    SetTargets(value []InvitationParticipantInfoable)()
    SetTenantId(value *string)()
    SetToneInfo(value ToneInfoable)()
    SetTranscription(value CallTranscriptionInfoable)()
}
