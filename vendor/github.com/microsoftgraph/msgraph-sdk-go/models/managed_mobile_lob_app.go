package models

import (
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

// ManagedMobileLobApp an abstract base class containing properties for all managed mobile line of business apps.
type ManagedMobileLobApp struct {
    ManagedApp
}
// NewManagedMobileLobApp instantiates a new ManagedMobileLobApp and sets the default values.
func NewManagedMobileLobApp()(*ManagedMobileLobApp) {
    m := &ManagedMobileLobApp{
        ManagedApp: *NewManagedApp(),
    }
    odataTypeValue := "#microsoft.graph.managedMobileLobApp"
    m.SetOdataType(&odataTypeValue)
    return m
}
// CreateManagedMobileLobAppFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
// returns a Parsable when successful
func CreateManagedMobileLobAppFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    if parseNode != nil {
        mappingValueNode, err := parseNode.GetChildNode("@odata.type")
        if err != nil {
            return nil, err
        }
        if mappingValueNode != nil {
            mappingValue, err := mappingValueNode.GetStringValue()
            if err != nil {
                return nil, err
            }
            if mappingValue != nil {
                switch *mappingValue {
                    case "#microsoft.graph.managedAndroidLobApp":
                        return NewManagedAndroidLobApp(), nil
                    case "#microsoft.graph.managedIOSLobApp":
                        return NewManagedIOSLobApp(), nil
                }
            }
        }
    }
    return NewManagedMobileLobApp(), nil
}
// GetCommittedContentVersion gets the committedContentVersion property value. The internal committed content version.
// returns a *string when successful
func (m *ManagedMobileLobApp) GetCommittedContentVersion()(*string) {
    val, err := m.GetBackingStore().Get("committedContentVersion")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetContentVersions gets the contentVersions property value. The list of content versions for this app.
// returns a []MobileAppContentable when successful
func (m *ManagedMobileLobApp) GetContentVersions()([]MobileAppContentable) {
    val, err := m.GetBackingStore().Get("contentVersions")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]MobileAppContentable)
    }
    return nil
}
// GetFieldDeserializers the deserialization information for the current model
// returns a map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error) when successful
func (m *ManagedMobileLobApp) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := m.ManagedApp.GetFieldDeserializers()
    res["committedContentVersion"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetCommittedContentVersion(val)
        }
        return nil
    }
    res["contentVersions"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetCollectionOfObjectValues(CreateMobileAppContentFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            res := make([]MobileAppContentable, len(val))
            for i, v := range val {
                if v != nil {
                    res[i] = v.(MobileAppContentable)
                }
            }
            m.SetContentVersions(res)
        }
        return nil
    }
    res["fileName"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetFileName(val)
        }
        return nil
    }
    res["size"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetInt64Value()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetSize(val)
        }
        return nil
    }
    return res
}
// GetFileName gets the fileName property value. The name of the main Lob application file.
// returns a *string when successful
func (m *ManagedMobileLobApp) GetFileName()(*string) {
    val, err := m.GetBackingStore().Get("fileName")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetSize gets the size property value. The total size, including all uploaded files.
// returns a *int64 when successful
func (m *ManagedMobileLobApp) GetSize()(*int64) {
    val, err := m.GetBackingStore().Get("size")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*int64)
    }
    return nil
}
// Serialize serializes information the current object
func (m *ManagedMobileLobApp) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    err := m.ManagedApp.Serialize(writer)
    if err != nil {
        return err
    }
    {
        err = writer.WriteStringValue("committedContentVersion", m.GetCommittedContentVersion())
        if err != nil {
            return err
        }
    }
    if m.GetContentVersions() != nil {
        cast := make([]i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, len(m.GetContentVersions()))
        for i, v := range m.GetContentVersions() {
            if v != nil {
                cast[i] = v.(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable)
            }
        }
        err = writer.WriteCollectionOfObjectValues("contentVersions", cast)
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteStringValue("fileName", m.GetFileName())
        if err != nil {
            return err
        }
    }
    return nil
}
// SetCommittedContentVersion sets the committedContentVersion property value. The internal committed content version.
func (m *ManagedMobileLobApp) SetCommittedContentVersion(value *string)() {
    err := m.GetBackingStore().Set("committedContentVersion", value)
    if err != nil {
        panic(err)
    }
}
// SetContentVersions sets the contentVersions property value. The list of content versions for this app.
func (m *ManagedMobileLobApp) SetContentVersions(value []MobileAppContentable)() {
    err := m.GetBackingStore().Set("contentVersions", value)
    if err != nil {
        panic(err)
    }
}
// SetFileName sets the fileName property value. The name of the main Lob application file.
func (m *ManagedMobileLobApp) SetFileName(value *string)() {
    err := m.GetBackingStore().Set("fileName", value)
    if err != nil {
        panic(err)
    }
}
// SetSize sets the size property value. The total size, including all uploaded files.
func (m *ManagedMobileLobApp) SetSize(value *int64)() {
    err := m.GetBackingStore().Set("size", value)
    if err != nil {
        panic(err)
    }
}
type ManagedMobileLobAppable interface {
    ManagedAppable
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable
    GetCommittedContentVersion()(*string)
    GetContentVersions()([]MobileAppContentable)
    GetFileName()(*string)
    GetSize()(*int64)
    SetCommittedContentVersion(value *string)()
    SetContentVersions(value []MobileAppContentable)()
    SetFileName(value *string)()
    SetSize(value *int64)()
}
