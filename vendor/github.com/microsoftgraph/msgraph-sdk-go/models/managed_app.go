package models

import (
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

// ManagedApp abstract class that contains properties and inherited properties for apps that you can manage with an Intune app protection policy.
type ManagedApp struct {
    MobileApp
}
// NewManagedApp instantiates a new ManagedApp and sets the default values.
func NewManagedApp()(*ManagedApp) {
    m := &ManagedApp{
        MobileApp: *NewMobileApp(),
    }
    odataTypeValue := "#microsoft.graph.managedApp"
    m.SetOdataType(&odataTypeValue)
    return m
}
// CreateManagedAppFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
// returns a Parsable when successful
func CreateManagedAppFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
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
                    case "#microsoft.graph.managedAndroidStoreApp":
                        return NewManagedAndroidStoreApp(), nil
                    case "#microsoft.graph.managedIOSLobApp":
                        return NewManagedIOSLobApp(), nil
                    case "#microsoft.graph.managedIOSStoreApp":
                        return NewManagedIOSStoreApp(), nil
                    case "#microsoft.graph.managedMobileLobApp":
                        return NewManagedMobileLobApp(), nil
                }
            }
        }
    }
    return NewManagedApp(), nil
}
// GetAppAvailability gets the appAvailability property value. A managed (MAM) application's availability.
// returns a *ManagedAppAvailability when successful
func (m *ManagedApp) GetAppAvailability()(*ManagedAppAvailability) {
    val, err := m.GetBackingStore().Get("appAvailability")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*ManagedAppAvailability)
    }
    return nil
}
// GetFieldDeserializers the deserialization information for the current model
// returns a map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error) when successful
func (m *ManagedApp) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := m.MobileApp.GetFieldDeserializers()
    res["appAvailability"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetEnumValue(ParseManagedAppAvailability)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetAppAvailability(val.(*ManagedAppAvailability))
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
// GetVersion gets the version property value. The Application's version.
// returns a *string when successful
func (m *ManagedApp) GetVersion()(*string) {
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
func (m *ManagedApp) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    err := m.MobileApp.Serialize(writer)
    if err != nil {
        return err
    }
    if m.GetAppAvailability() != nil {
        cast := (*m.GetAppAvailability()).String()
        err = writer.WriteStringValue("appAvailability", &cast)
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
// SetAppAvailability sets the appAvailability property value. A managed (MAM) application's availability.
func (m *ManagedApp) SetAppAvailability(value *ManagedAppAvailability)() {
    err := m.GetBackingStore().Set("appAvailability", value)
    if err != nil {
        panic(err)
    }
}
// SetVersion sets the version property value. The Application's version.
func (m *ManagedApp) SetVersion(value *string)() {
    err := m.GetBackingStore().Set("version", value)
    if err != nil {
        panic(err)
    }
}
type ManagedAppable interface {
    MobileAppable
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable
    GetAppAvailability()(*ManagedAppAvailability)
    GetVersion()(*string)
    SetAppAvailability(value *ManagedAppAvailability)()
    SetVersion(value *string)()
}
