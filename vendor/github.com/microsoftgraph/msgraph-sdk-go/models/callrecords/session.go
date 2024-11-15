package callrecords

import (
    i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e "time"
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
    iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242 "github.com/microsoftgraph/msgraph-sdk-go/models"
)

type Session struct {
    iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.Entity
}
// NewSession instantiates a new Session and sets the default values.
func NewSession()(*Session) {
    m := &Session{
        Entity: *iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.NewEntity(),
    }
    return m
}
// CreateSessionFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
// returns a Parsable when successful
func CreateSessionFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewSession(), nil
}
// GetCallee gets the callee property value. Endpoint that answered the session.
// returns a Endpointable when successful
func (m *Session) GetCallee()(Endpointable) {
    val, err := m.GetBackingStore().Get("callee")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(Endpointable)
    }
    return nil
}
// GetCaller gets the caller property value. Endpoint that initiated the session.
// returns a Endpointable when successful
func (m *Session) GetCaller()(Endpointable) {
    val, err := m.GetBackingStore().Get("caller")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(Endpointable)
    }
    return nil
}
// GetEndDateTime gets the endDateTime property value. UTC time when the last user left the session. The DateTimeOffset type represents date and time information using ISO 8601 format and is always in UTC time. For example, midnight UTC on Jan 1, 2014 is 2014-01-01T00:00:00Z
// returns a *Time when successful
func (m *Session) GetEndDateTime()(*i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time) {
    val, err := m.GetBackingStore().Get("endDateTime")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)
    }
    return nil
}
// GetFailureInfo gets the failureInfo property value. Failure information associated with the session if the session failed.
// returns a FailureInfoable when successful
func (m *Session) GetFailureInfo()(FailureInfoable) {
    val, err := m.GetBackingStore().Get("failureInfo")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(FailureInfoable)
    }
    return nil
}
// GetFieldDeserializers the deserialization information for the current model
// returns a map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error) when successful
func (m *Session) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := m.Entity.GetFieldDeserializers()
    res["callee"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetObjectValue(CreateEndpointFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetCallee(val.(Endpointable))
        }
        return nil
    }
    res["caller"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetObjectValue(CreateEndpointFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetCaller(val.(Endpointable))
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
    res["failureInfo"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetObjectValue(CreateFailureInfoFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetFailureInfo(val.(FailureInfoable))
        }
        return nil
    }
    res["isTest"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetBoolValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetIsTest(val)
        }
        return nil
    }
    res["modalities"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
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
            m.SetModalities(res)
        }
        return nil
    }
    res["segments"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetCollectionOfObjectValues(CreateSegmentFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            res := make([]Segmentable, len(val))
            for i, v := range val {
                if v != nil {
                    res[i] = v.(Segmentable)
                }
            }
            m.SetSegments(res)
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
    return res
}
// GetIsTest gets the isTest property value. Specifies whether the session is a test.
// returns a *bool when successful
func (m *Session) GetIsTest()(*bool) {
    val, err := m.GetBackingStore().Get("isTest")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*bool)
    }
    return nil
}
// GetModalities gets the modalities property value. List of modalities present in the session. Possible values are: unknown, audio, video, videoBasedScreenSharing, data, screenSharing, unknownFutureValue.
// returns a []Modality when successful
func (m *Session) GetModalities()([]Modality) {
    val, err := m.GetBackingStore().Get("modalities")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]Modality)
    }
    return nil
}
// GetSegments gets the segments property value. The list of segments involved in the session. Read-only. Nullable.
// returns a []Segmentable when successful
func (m *Session) GetSegments()([]Segmentable) {
    val, err := m.GetBackingStore().Get("segments")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]Segmentable)
    }
    return nil
}
// GetStartDateTime gets the startDateTime property value. UTC time when the first user joined the session. The DateTimeOffset type represents date and time information using ISO 8601 format and is always in UTC time. For example, midnight UTC on Jan 1, 2014 is 2014-01-01T00:00:00Z
// returns a *Time when successful
func (m *Session) GetStartDateTime()(*i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time) {
    val, err := m.GetBackingStore().Get("startDateTime")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)
    }
    return nil
}
// Serialize serializes information the current object
func (m *Session) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    err := m.Entity.Serialize(writer)
    if err != nil {
        return err
    }
    {
        err = writer.WriteObjectValue("callee", m.GetCallee())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteObjectValue("caller", m.GetCaller())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteTimeValue("endDateTime", m.GetEndDateTime())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteObjectValue("failureInfo", m.GetFailureInfo())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteBoolValue("isTest", m.GetIsTest())
        if err != nil {
            return err
        }
    }
    if m.GetModalities() != nil {
        err = writer.WriteCollectionOfStringValues("modalities", SerializeModality(m.GetModalities()))
        if err != nil {
            return err
        }
    }
    if m.GetSegments() != nil {
        cast := make([]i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, len(m.GetSegments()))
        for i, v := range m.GetSegments() {
            if v != nil {
                cast[i] = v.(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable)
            }
        }
        err = writer.WriteCollectionOfObjectValues("segments", cast)
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteTimeValue("startDateTime", m.GetStartDateTime())
        if err != nil {
            return err
        }
    }
    return nil
}
// SetCallee sets the callee property value. Endpoint that answered the session.
func (m *Session) SetCallee(value Endpointable)() {
    err := m.GetBackingStore().Set("callee", value)
    if err != nil {
        panic(err)
    }
}
// SetCaller sets the caller property value. Endpoint that initiated the session.
func (m *Session) SetCaller(value Endpointable)() {
    err := m.GetBackingStore().Set("caller", value)
    if err != nil {
        panic(err)
    }
}
// SetEndDateTime sets the endDateTime property value. UTC time when the last user left the session. The DateTimeOffset type represents date and time information using ISO 8601 format and is always in UTC time. For example, midnight UTC on Jan 1, 2014 is 2014-01-01T00:00:00Z
func (m *Session) SetEndDateTime(value *i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)() {
    err := m.GetBackingStore().Set("endDateTime", value)
    if err != nil {
        panic(err)
    }
}
// SetFailureInfo sets the failureInfo property value. Failure information associated with the session if the session failed.
func (m *Session) SetFailureInfo(value FailureInfoable)() {
    err := m.GetBackingStore().Set("failureInfo", value)
    if err != nil {
        panic(err)
    }
}
// SetIsTest sets the isTest property value. Specifies whether the session is a test.
func (m *Session) SetIsTest(value *bool)() {
    err := m.GetBackingStore().Set("isTest", value)
    if err != nil {
        panic(err)
    }
}
// SetModalities sets the modalities property value. List of modalities present in the session. Possible values are: unknown, audio, video, videoBasedScreenSharing, data, screenSharing, unknownFutureValue.
func (m *Session) SetModalities(value []Modality)() {
    err := m.GetBackingStore().Set("modalities", value)
    if err != nil {
        panic(err)
    }
}
// SetSegments sets the segments property value. The list of segments involved in the session. Read-only. Nullable.
func (m *Session) SetSegments(value []Segmentable)() {
    err := m.GetBackingStore().Set("segments", value)
    if err != nil {
        panic(err)
    }
}
// SetStartDateTime sets the startDateTime property value. UTC time when the first user joined the session. The DateTimeOffset type represents date and time information using ISO 8601 format and is always in UTC time. For example, midnight UTC on Jan 1, 2014 is 2014-01-01T00:00:00Z
func (m *Session) SetStartDateTime(value *i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)() {
    err := m.GetBackingStore().Set("startDateTime", value)
    if err != nil {
        panic(err)
    }
}
type Sessionable interface {
    iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.Entityable
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable
    GetCallee()(Endpointable)
    GetCaller()(Endpointable)
    GetEndDateTime()(*i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)
    GetFailureInfo()(FailureInfoable)
    GetIsTest()(*bool)
    GetModalities()([]Modality)
    GetSegments()([]Segmentable)
    GetStartDateTime()(*i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)
    SetCallee(value Endpointable)()
    SetCaller(value Endpointable)()
    SetEndDateTime(value *i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)()
    SetFailureInfo(value FailureInfoable)()
    SetIsTest(value *bool)()
    SetModalities(value []Modality)()
    SetSegments(value []Segmentable)()
    SetStartDateTime(value *i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)()
}
