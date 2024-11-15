package models

import (
    i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e "time"
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

type DirectoryDefinition struct {
    Entity
}
// NewDirectoryDefinition instantiates a new DirectoryDefinition and sets the default values.
func NewDirectoryDefinition()(*DirectoryDefinition) {
    m := &DirectoryDefinition{
        Entity: *NewEntity(),
    }
    return m
}
// CreateDirectoryDefinitionFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
// returns a Parsable when successful
func CreateDirectoryDefinitionFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewDirectoryDefinition(), nil
}
// GetDiscoverabilities gets the discoverabilities property value. The discoverabilities property
// returns a *DirectoryDefinitionDiscoverabilities when successful
func (m *DirectoryDefinition) GetDiscoverabilities()(*DirectoryDefinitionDiscoverabilities) {
    val, err := m.GetBackingStore().Get("discoverabilities")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*DirectoryDefinitionDiscoverabilities)
    }
    return nil
}
// GetDiscoveryDateTime gets the discoveryDateTime property value. Represents the discovery date and time using ISO 8601 format and is always in UTC time. For example, midnight UTC on Jan 1, 2014 is 2014-01-01T00:00:00Z.
// returns a *Time when successful
func (m *DirectoryDefinition) GetDiscoveryDateTime()(*i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time) {
    val, err := m.GetBackingStore().Get("discoveryDateTime")
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
func (m *DirectoryDefinition) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := m.Entity.GetFieldDeserializers()
    res["discoverabilities"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetEnumValue(ParseDirectoryDefinitionDiscoverabilities)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetDiscoverabilities(val.(*DirectoryDefinitionDiscoverabilities))
        }
        return nil
    }
    res["discoveryDateTime"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetTimeValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetDiscoveryDateTime(val)
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
    res["objects"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetCollectionOfObjectValues(CreateObjectDefinitionFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            res := make([]ObjectDefinitionable, len(val))
            for i, v := range val {
                if v != nil {
                    res[i] = v.(ObjectDefinitionable)
                }
            }
            m.SetObjects(res)
        }
        return nil
    }
    res["readOnly"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetBoolValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetReadOnly(val)
        }
        return nil
    }
    res["version"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetVersion(val)
        }
        return nil
    }
    return res
}
// GetName gets the name property value. Name of the directory. Must be unique within the synchronization schema. Not nullable.
// returns a *string when successful
func (m *DirectoryDefinition) GetName()(*string) {
    val, err := m.GetBackingStore().Get("name")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetObjects gets the objects property value. Collection of objects supported by the directory.
// returns a []ObjectDefinitionable when successful
func (m *DirectoryDefinition) GetObjects()([]ObjectDefinitionable) {
    val, err := m.GetBackingStore().Get("objects")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]ObjectDefinitionable)
    }
    return nil
}
// GetReadOnly gets the readOnly property value. Whether this object is read-only.
// returns a *bool when successful
func (m *DirectoryDefinition) GetReadOnly()(*bool) {
    val, err := m.GetBackingStore().Get("readOnly")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*bool)
    }
    return nil
}
// GetVersion gets the version property value. Read only value that indicates version discovered. null if discovery hasn't yet occurred.
// returns a *string when successful
func (m *DirectoryDefinition) GetVersion()(*string) {
    val, err := m.GetBackingStore().Get("version")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// Serialize serializes information the current object
func (m *DirectoryDefinition) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    err := m.Entity.Serialize(writer)
    if err != nil {
        return err
    }
    if m.GetDiscoverabilities() != nil {
        cast := (*m.GetDiscoverabilities()).String()
        err = writer.WriteStringValue("discoverabilities", &cast)
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteTimeValue("discoveryDateTime", m.GetDiscoveryDateTime())
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
    if m.GetObjects() != nil {
        cast := make([]i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, len(m.GetObjects()))
        for i, v := range m.GetObjects() {
            if v != nil {
                cast[i] = v.(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable)
            }
        }
        err = writer.WriteCollectionOfObjectValues("objects", cast)
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteBoolValue("readOnly", m.GetReadOnly())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteStringValue("version", m.GetVersion())
        if err != nil {
            return err
        }
    }
    return nil
}
// SetDiscoverabilities sets the discoverabilities property value. The discoverabilities property
func (m *DirectoryDefinition) SetDiscoverabilities(value *DirectoryDefinitionDiscoverabilities)() {
    err := m.GetBackingStore().Set("discoverabilities", value)
    if err != nil {
        panic(err)
    }
}
// SetDiscoveryDateTime sets the discoveryDateTime property value. Represents the discovery date and time using ISO 8601 format and is always in UTC time. For example, midnight UTC on Jan 1, 2014 is 2014-01-01T00:00:00Z.
func (m *DirectoryDefinition) SetDiscoveryDateTime(value *i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)() {
    err := m.GetBackingStore().Set("discoveryDateTime", value)
    if err != nil {
        panic(err)
    }
}
// SetName sets the name property value. Name of the directory. Must be unique within the synchronization schema. Not nullable.
func (m *DirectoryDefinition) SetName(value *string)() {
    err := m.GetBackingStore().Set("name", value)
    if err != nil {
        panic(err)
    }
}
// SetObjects sets the objects property value. Collection of objects supported by the directory.
func (m *DirectoryDefinition) SetObjects(value []ObjectDefinitionable)() {
    err := m.GetBackingStore().Set("objects", value)
    if err != nil {
        panic(err)
    }
}
// SetReadOnly sets the readOnly property value. Whether this object is read-only.
func (m *DirectoryDefinition) SetReadOnly(value *bool)() {
    err := m.GetBackingStore().Set("readOnly", value)
    if err != nil {
        panic(err)
    }
}
// SetVersion sets the version property value. Read only value that indicates version discovered. null if discovery hasn't yet occurred.
func (m *DirectoryDefinition) SetVersion(value *string)() {
    err := m.GetBackingStore().Set("version", value)
    if err != nil {
        panic(err)
    }
}
type DirectoryDefinitionable interface {
    Entityable
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable
    GetDiscoverabilities()(*DirectoryDefinitionDiscoverabilities)
    GetDiscoveryDateTime()(*i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)
    GetName()(*string)
    GetObjects()([]ObjectDefinitionable)
    GetReadOnly()(*bool)
    GetVersion()(*string)
    SetDiscoverabilities(value *DirectoryDefinitionDiscoverabilities)()
    SetDiscoveryDateTime(value *i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)()
    SetName(value *string)()
    SetObjects(value []ObjectDefinitionable)()
    SetReadOnly(value *bool)()
    SetVersion(value *string)()
}
