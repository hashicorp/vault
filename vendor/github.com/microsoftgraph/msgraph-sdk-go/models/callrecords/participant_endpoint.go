package callrecords

import (
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
    iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242 "github.com/microsoftgraph/msgraph-sdk-go/models"
)

type ParticipantEndpoint struct {
    Endpoint
}
// NewParticipantEndpoint instantiates a new ParticipantEndpoint and sets the default values.
func NewParticipantEndpoint()(*ParticipantEndpoint) {
    m := &ParticipantEndpoint{
        Endpoint: *NewEndpoint(),
    }
    odataTypeValue := "#microsoft.graph.callRecords.participantEndpoint"
    m.SetOdataType(&odataTypeValue)
    return m
}
// CreateParticipantEndpointFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
// returns a Parsable when successful
func CreateParticipantEndpointFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewParticipantEndpoint(), nil
}
// GetAssociatedIdentity gets the associatedIdentity property value. Identity associated with the endpoint.
// returns a Identityable when successful
func (m *ParticipantEndpoint) GetAssociatedIdentity()(iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.Identityable) {
    val, err := m.GetBackingStore().Get("associatedIdentity")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.Identityable)
    }
    return nil
}
// GetCpuCoresCount gets the cpuCoresCount property value. CPU number of cores used by the media endpoint.
// returns a *int32 when successful
func (m *ParticipantEndpoint) GetCpuCoresCount()(*int32) {
    val, err := m.GetBackingStore().Get("cpuCoresCount")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*int32)
    }
    return nil
}
// GetCpuName gets the cpuName property value. CPU name used by the media endpoint.
// returns a *string when successful
func (m *ParticipantEndpoint) GetCpuName()(*string) {
    val, err := m.GetBackingStore().Get("cpuName")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetCpuProcessorSpeedInMhz gets the cpuProcessorSpeedInMhz property value. CPU processor speed used by the media endpoint.
// returns a *int32 when successful
func (m *ParticipantEndpoint) GetCpuProcessorSpeedInMhz()(*int32) {
    val, err := m.GetBackingStore().Get("cpuProcessorSpeedInMhz")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*int32)
    }
    return nil
}
// GetFeedback gets the feedback property value. The feedback provided by the user of this endpoint about the quality of the session.
// returns a UserFeedbackable when successful
func (m *ParticipantEndpoint) GetFeedback()(UserFeedbackable) {
    val, err := m.GetBackingStore().Get("feedback")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(UserFeedbackable)
    }
    return nil
}
// GetFieldDeserializers the deserialization information for the current model
// returns a map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error) when successful
func (m *ParticipantEndpoint) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := m.Endpoint.GetFieldDeserializers()
    res["associatedIdentity"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetObjectValue(iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.CreateIdentityFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetAssociatedIdentity(val.(iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.Identityable))
        }
        return nil
    }
    res["cpuCoresCount"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetInt32Value()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetCpuCoresCount(val)
        }
        return nil
    }
    res["cpuName"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetCpuName(val)
        }
        return nil
    }
    res["cpuProcessorSpeedInMhz"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetInt32Value()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetCpuProcessorSpeedInMhz(val)
        }
        return nil
    }
    res["feedback"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetObjectValue(CreateUserFeedbackFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetFeedback(val.(UserFeedbackable))
        }
        return nil
    }
    res["identity"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetObjectValue(iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.CreateIdentitySetFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetIdentity(val.(iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.IdentitySetable))
        }
        return nil
    }
    res["name"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetName(val)
        }
        return nil
    }
    return res
}
// GetIdentity gets the identity property value. Identity associated with the endpoint. The identity property is deprecated and will stop returning data on June 30, 2026. Going forward, use the associatedIdentity property.
// returns a IdentitySetable when successful
func (m *ParticipantEndpoint) GetIdentity()(iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.IdentitySetable) {
    val, err := m.GetBackingStore().Get("identity")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.IdentitySetable)
    }
    return nil
}
// GetName gets the name property value. Name of the device used by the media endpoint.
// returns a *string when successful
func (m *ParticipantEndpoint) GetName()(*string) {
    val, err := m.GetBackingStore().Get("name")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// Serialize serializes information the current object
func (m *ParticipantEndpoint) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    err := m.Endpoint.Serialize(writer)
    if err != nil {
        return err
    }
    {
        err = writer.WriteObjectValue("associatedIdentity", m.GetAssociatedIdentity())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteInt32Value("cpuCoresCount", m.GetCpuCoresCount())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteStringValue("cpuName", m.GetCpuName())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteInt32Value("cpuProcessorSpeedInMhz", m.GetCpuProcessorSpeedInMhz())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteObjectValue("feedback", m.GetFeedback())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteObjectValue("identity", m.GetIdentity())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteStringValue("name", m.GetName())
        if err != nil {
            return err
        }
    }
    return nil
}
// SetAssociatedIdentity sets the associatedIdentity property value. Identity associated with the endpoint.
func (m *ParticipantEndpoint) SetAssociatedIdentity(value iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.Identityable)() {
    err := m.GetBackingStore().Set("associatedIdentity", value)
    if err != nil {
        panic(err)
    }
}
// SetCpuCoresCount sets the cpuCoresCount property value. CPU number of cores used by the media endpoint.
func (m *ParticipantEndpoint) SetCpuCoresCount(value *int32)() {
    err := m.GetBackingStore().Set("cpuCoresCount", value)
    if err != nil {
        panic(err)
    }
}
// SetCpuName sets the cpuName property value. CPU name used by the media endpoint.
func (m *ParticipantEndpoint) SetCpuName(value *string)() {
    err := m.GetBackingStore().Set("cpuName", value)
    if err != nil {
        panic(err)
    }
}
// SetCpuProcessorSpeedInMhz sets the cpuProcessorSpeedInMhz property value. CPU processor speed used by the media endpoint.
func (m *ParticipantEndpoint) SetCpuProcessorSpeedInMhz(value *int32)() {
    err := m.GetBackingStore().Set("cpuProcessorSpeedInMhz", value)
    if err != nil {
        panic(err)
    }
}
// SetFeedback sets the feedback property value. The feedback provided by the user of this endpoint about the quality of the session.
func (m *ParticipantEndpoint) SetFeedback(value UserFeedbackable)() {
    err := m.GetBackingStore().Set("feedback", value)
    if err != nil {
        panic(err)
    }
}
// SetIdentity sets the identity property value. Identity associated with the endpoint. The identity property is deprecated and will stop returning data on June 30, 2026. Going forward, use the associatedIdentity property.
func (m *ParticipantEndpoint) SetIdentity(value iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.IdentitySetable)() {
    err := m.GetBackingStore().Set("identity", value)
    if err != nil {
        panic(err)
    }
}
// SetName sets the name property value. Name of the device used by the media endpoint.
func (m *ParticipantEndpoint) SetName(value *string)() {
    err := m.GetBackingStore().Set("name", value)
    if err != nil {
        panic(err)
    }
}
type ParticipantEndpointable interface {
    Endpointable
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable
    GetAssociatedIdentity()(iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.Identityable)
    GetCpuCoresCount()(*int32)
    GetCpuName()(*string)
    GetCpuProcessorSpeedInMhz()(*int32)
    GetFeedback()(UserFeedbackable)
    GetIdentity()(iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.IdentitySetable)
    GetName()(*string)
    SetAssociatedIdentity(value iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.Identityable)()
    SetCpuCoresCount(value *int32)()
    SetCpuName(value *string)()
    SetCpuProcessorSpeedInMhz(value *int32)()
    SetFeedback(value UserFeedbackable)()
    SetIdentity(value iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.IdentitySetable)()
    SetName(value *string)()
}
