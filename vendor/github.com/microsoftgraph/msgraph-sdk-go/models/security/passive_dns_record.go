package security

import (
    i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e "time"
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

type PassiveDnsRecord struct {
    Artifact
}
// NewPassiveDnsRecord instantiates a new PassiveDnsRecord and sets the default values.
func NewPassiveDnsRecord()(*PassiveDnsRecord) {
    m := &PassiveDnsRecord{
        Artifact: *NewArtifact(),
    }
    odataTypeValue := "#microsoft.graph.security.passiveDnsRecord"
    m.SetOdataType(&odataTypeValue)
    return m
}
// CreatePassiveDnsRecordFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
// returns a Parsable when successful
func CreatePassiveDnsRecordFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewPassiveDnsRecord(), nil
}
// GetArtifact gets the artifact property value. The artifact property
// returns a Artifactable when successful
func (m *PassiveDnsRecord) GetArtifact()(Artifactable) {
    val, err := m.GetBackingStore().Get("artifact")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(Artifactable)
    }
    return nil
}
// GetCollectedDateTime gets the collectedDateTime property value. The date and time that this passiveDnsRecord entry was collected by Microsoft. The Timestamp type represents date and time information using ISO 8601 format and is always in UTC time. For example, midnight UTC on Jan 1, 2014 is 2014-01-01T00:00:00Z.
// returns a *Time when successful
func (m *PassiveDnsRecord) GetCollectedDateTime()(*i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time) {
    val, err := m.GetBackingStore().Get("collectedDateTime")
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
func (m *PassiveDnsRecord) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := m.Artifact.GetFieldDeserializers()
    res["artifact"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetObjectValue(CreateArtifactFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetArtifact(val.(Artifactable))
        }
        return nil
    }
    res["collectedDateTime"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetTimeValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetCollectedDateTime(val)
        }
        return nil
    }
    res["firstSeenDateTime"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetTimeValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetFirstSeenDateTime(val)
        }
        return nil
    }
    res["lastSeenDateTime"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetTimeValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetLastSeenDateTime(val)
        }
        return nil
    }
    res["parentHost"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetObjectValue(CreateHostFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetParentHost(val.(Hostable))
        }
        return nil
    }
    res["recordType"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetRecordType(val)
        }
        return nil
    }
    return res
}
// GetFirstSeenDateTime gets the firstSeenDateTime property value. The date and time when this passiveDnsRecord entry was first seen. The Timestamp type represents date and time information using ISO 8601 format and is always in UTC time. For example, midnight UTC on Jan 1, 2014 is 2014-01-01T00:00:00Z.
// returns a *Time when successful
func (m *PassiveDnsRecord) GetFirstSeenDateTime()(*i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time) {
    val, err := m.GetBackingStore().Get("firstSeenDateTime")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)
    }
    return nil
}
// GetLastSeenDateTime gets the lastSeenDateTime property value. The date and time when this passiveDnsRecord entry was most recently seen. The Timestamp type represents date and time information using ISO 8601 format and is always in UTC time. For example, midnight UTC on Jan 1, 2014 is 2014-01-01T00:00:00Z.
// returns a *Time when successful
func (m *PassiveDnsRecord) GetLastSeenDateTime()(*i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time) {
    val, err := m.GetBackingStore().Get("lastSeenDateTime")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)
    }
    return nil
}
// GetParentHost gets the parentHost property value. The parentHost property
// returns a Hostable when successful
func (m *PassiveDnsRecord) GetParentHost()(Hostable) {
    val, err := m.GetBackingStore().Get("parentHost")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(Hostable)
    }
    return nil
}
// GetRecordType gets the recordType property value. The DNS record type for this passiveDnsRecord entry.
// returns a *string when successful
func (m *PassiveDnsRecord) GetRecordType()(*string) {
    val, err := m.GetBackingStore().Get("recordType")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// Serialize serializes information the current object
func (m *PassiveDnsRecord) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    err := m.Artifact.Serialize(writer)
    if err != nil {
        return err
    }
    {
        err = writer.WriteObjectValue("artifact", m.GetArtifact())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteTimeValue("collectedDateTime", m.GetCollectedDateTime())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteTimeValue("firstSeenDateTime", m.GetFirstSeenDateTime())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteTimeValue("lastSeenDateTime", m.GetLastSeenDateTime())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteObjectValue("parentHost", m.GetParentHost())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteStringValue("recordType", m.GetRecordType())
        if err != nil {
            return err
        }
    }
    return nil
}
// SetArtifact sets the artifact property value. The artifact property
func (m *PassiveDnsRecord) SetArtifact(value Artifactable)() {
    err := m.GetBackingStore().Set("artifact", value)
    if err != nil {
        panic(err)
    }
}
// SetCollectedDateTime sets the collectedDateTime property value. The date and time that this passiveDnsRecord entry was collected by Microsoft. The Timestamp type represents date and time information using ISO 8601 format and is always in UTC time. For example, midnight UTC on Jan 1, 2014 is 2014-01-01T00:00:00Z.
func (m *PassiveDnsRecord) SetCollectedDateTime(value *i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)() {
    err := m.GetBackingStore().Set("collectedDateTime", value)
    if err != nil {
        panic(err)
    }
}
// SetFirstSeenDateTime sets the firstSeenDateTime property value. The date and time when this passiveDnsRecord entry was first seen. The Timestamp type represents date and time information using ISO 8601 format and is always in UTC time. For example, midnight UTC on Jan 1, 2014 is 2014-01-01T00:00:00Z.
func (m *PassiveDnsRecord) SetFirstSeenDateTime(value *i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)() {
    err := m.GetBackingStore().Set("firstSeenDateTime", value)
    if err != nil {
        panic(err)
    }
}
// SetLastSeenDateTime sets the lastSeenDateTime property value. The date and time when this passiveDnsRecord entry was most recently seen. The Timestamp type represents date and time information using ISO 8601 format and is always in UTC time. For example, midnight UTC on Jan 1, 2014 is 2014-01-01T00:00:00Z.
func (m *PassiveDnsRecord) SetLastSeenDateTime(value *i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)() {
    err := m.GetBackingStore().Set("lastSeenDateTime", value)
    if err != nil {
        panic(err)
    }
}
// SetParentHost sets the parentHost property value. The parentHost property
func (m *PassiveDnsRecord) SetParentHost(value Hostable)() {
    err := m.GetBackingStore().Set("parentHost", value)
    if err != nil {
        panic(err)
    }
}
// SetRecordType sets the recordType property value. The DNS record type for this passiveDnsRecord entry.
func (m *PassiveDnsRecord) SetRecordType(value *string)() {
    err := m.GetBackingStore().Set("recordType", value)
    if err != nil {
        panic(err)
    }
}
type PassiveDnsRecordable interface {
    Artifactable
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable
    GetArtifact()(Artifactable)
    GetCollectedDateTime()(*i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)
    GetFirstSeenDateTime()(*i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)
    GetLastSeenDateTime()(*i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)
    GetParentHost()(Hostable)
    GetRecordType()(*string)
    SetArtifact(value Artifactable)()
    SetCollectedDateTime(value *i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)()
    SetFirstSeenDateTime(value *i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)()
    SetLastSeenDateTime(value *i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)()
    SetParentHost(value Hostable)()
    SetRecordType(value *string)()
}
