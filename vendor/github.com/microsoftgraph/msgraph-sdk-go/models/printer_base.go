package models

import (
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

type PrinterBase struct {
    Entity
}
// NewPrinterBase instantiates a new PrinterBase and sets the default values.
func NewPrinterBase()(*PrinterBase) {
    m := &PrinterBase{
        Entity: *NewEntity(),
    }
    return m
}
// CreatePrinterBaseFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
// returns a Parsable when successful
func CreatePrinterBaseFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
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
                    case "#microsoft.graph.printer":
                        return NewPrinter(), nil
                    case "#microsoft.graph.printerShare":
                        return NewPrinterShare(), nil
                }
            }
        }
    }
    return NewPrinterBase(), nil
}
// GetCapabilities gets the capabilities property value. The capabilities of the printer/printerShare.
// returns a PrinterCapabilitiesable when successful
func (m *PrinterBase) GetCapabilities()(PrinterCapabilitiesable) {
    val, err := m.GetBackingStore().Get("capabilities")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(PrinterCapabilitiesable)
    }
    return nil
}
// GetDefaults gets the defaults property value. The default print settings of printer/printerShare.
// returns a PrinterDefaultsable when successful
func (m *PrinterBase) GetDefaults()(PrinterDefaultsable) {
    val, err := m.GetBackingStore().Get("defaults")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(PrinterDefaultsable)
    }
    return nil
}
// GetDisplayName gets the displayName property value. The name of the printer/printerShare.
// returns a *string when successful
func (m *PrinterBase) GetDisplayName()(*string) {
    val, err := m.GetBackingStore().Get("displayName")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetFieldDeserializers the deserialization information for the current model
// returns a map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error) when successful
func (m *PrinterBase) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := m.Entity.GetFieldDeserializers()
    res["capabilities"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetObjectValue(CreatePrinterCapabilitiesFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetCapabilities(val.(PrinterCapabilitiesable))
        }
        return nil
    }
    res["defaults"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetObjectValue(CreatePrinterDefaultsFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetDefaults(val.(PrinterDefaultsable))
        }
        return nil
    }
    res["displayName"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetDisplayName(val)
        }
        return nil
    }
    res["isAcceptingJobs"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetBoolValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetIsAcceptingJobs(val)
        }
        return nil
    }
    res["jobs"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetCollectionOfObjectValues(CreatePrintJobFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            res := make([]PrintJobable, len(val))
            for i, v := range val {
                if v != nil {
                    res[i] = v.(PrintJobable)
                }
            }
            m.SetJobs(res)
        }
        return nil
    }
    res["location"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetObjectValue(CreatePrinterLocationFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetLocation(val.(PrinterLocationable))
        }
        return nil
    }
    res["manufacturer"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetManufacturer(val)
        }
        return nil
    }
    res["model"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetModel(val)
        }
        return nil
    }
    res["status"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetObjectValue(CreatePrinterStatusFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetStatus(val.(PrinterStatusable))
        }
        return nil
    }
    return res
}
// GetIsAcceptingJobs gets the isAcceptingJobs property value. Specifies whether the printer/printerShare is currently accepting new print jobs.
// returns a *bool when successful
func (m *PrinterBase) GetIsAcceptingJobs()(*bool) {
    val, err := m.GetBackingStore().Get("isAcceptingJobs")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*bool)
    }
    return nil
}
// GetJobs gets the jobs property value. The list of jobs that are queued for printing by the printer/printerShare.
// returns a []PrintJobable when successful
func (m *PrinterBase) GetJobs()([]PrintJobable) {
    val, err := m.GetBackingStore().Get("jobs")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]PrintJobable)
    }
    return nil
}
// GetLocation gets the location property value. The physical and/or organizational location of the printer/printerShare.
// returns a PrinterLocationable when successful
func (m *PrinterBase) GetLocation()(PrinterLocationable) {
    val, err := m.GetBackingStore().Get("location")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(PrinterLocationable)
    }
    return nil
}
// GetManufacturer gets the manufacturer property value. The manufacturer of the printer/printerShare.
// returns a *string when successful
func (m *PrinterBase) GetManufacturer()(*string) {
    val, err := m.GetBackingStore().Get("manufacturer")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetModel gets the model property value. The model name of the printer/printerShare.
// returns a *string when successful
func (m *PrinterBase) GetModel()(*string) {
    val, err := m.GetBackingStore().Get("model")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetStatus gets the status property value. The status property
// returns a PrinterStatusable when successful
func (m *PrinterBase) GetStatus()(PrinterStatusable) {
    val, err := m.GetBackingStore().Get("status")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(PrinterStatusable)
    }
    return nil
}
// Serialize serializes information the current object
func (m *PrinterBase) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    err := m.Entity.Serialize(writer)
    if err != nil {
        return err
    }
    {
        err = writer.WriteObjectValue("capabilities", m.GetCapabilities())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteObjectValue("defaults", m.GetDefaults())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteStringValue("displayName", m.GetDisplayName())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteBoolValue("isAcceptingJobs", m.GetIsAcceptingJobs())
        if err != nil {
            return err
        }
    }
    if m.GetJobs() != nil {
        cast := make([]i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, len(m.GetJobs()))
        for i, v := range m.GetJobs() {
            if v != nil {
                cast[i] = v.(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable)
            }
        }
        err = writer.WriteCollectionOfObjectValues("jobs", cast)
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
        err = writer.WriteStringValue("manufacturer", m.GetManufacturer())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteStringValue("model", m.GetModel())
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
    return nil
}
// SetCapabilities sets the capabilities property value. The capabilities of the printer/printerShare.
func (m *PrinterBase) SetCapabilities(value PrinterCapabilitiesable)() {
    err := m.GetBackingStore().Set("capabilities", value)
    if err != nil {
        panic(err)
    }
}
// SetDefaults sets the defaults property value. The default print settings of printer/printerShare.
func (m *PrinterBase) SetDefaults(value PrinterDefaultsable)() {
    err := m.GetBackingStore().Set("defaults", value)
    if err != nil {
        panic(err)
    }
}
// SetDisplayName sets the displayName property value. The name of the printer/printerShare.
func (m *PrinterBase) SetDisplayName(value *string)() {
    err := m.GetBackingStore().Set("displayName", value)
    if err != nil {
        panic(err)
    }
}
// SetIsAcceptingJobs sets the isAcceptingJobs property value. Specifies whether the printer/printerShare is currently accepting new print jobs.
func (m *PrinterBase) SetIsAcceptingJobs(value *bool)() {
    err := m.GetBackingStore().Set("isAcceptingJobs", value)
    if err != nil {
        panic(err)
    }
}
// SetJobs sets the jobs property value. The list of jobs that are queued for printing by the printer/printerShare.
func (m *PrinterBase) SetJobs(value []PrintJobable)() {
    err := m.GetBackingStore().Set("jobs", value)
    if err != nil {
        panic(err)
    }
}
// SetLocation sets the location property value. The physical and/or organizational location of the printer/printerShare.
func (m *PrinterBase) SetLocation(value PrinterLocationable)() {
    err := m.GetBackingStore().Set("location", value)
    if err != nil {
        panic(err)
    }
}
// SetManufacturer sets the manufacturer property value. The manufacturer of the printer/printerShare.
func (m *PrinterBase) SetManufacturer(value *string)() {
    err := m.GetBackingStore().Set("manufacturer", value)
    if err != nil {
        panic(err)
    }
}
// SetModel sets the model property value. The model name of the printer/printerShare.
func (m *PrinterBase) SetModel(value *string)() {
    err := m.GetBackingStore().Set("model", value)
    if err != nil {
        panic(err)
    }
}
// SetStatus sets the status property value. The status property
func (m *PrinterBase) SetStatus(value PrinterStatusable)() {
    err := m.GetBackingStore().Set("status", value)
    if err != nil {
        panic(err)
    }
}
type PrinterBaseable interface {
    Entityable
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable
    GetCapabilities()(PrinterCapabilitiesable)
    GetDefaults()(PrinterDefaultsable)
    GetDisplayName()(*string)
    GetIsAcceptingJobs()(*bool)
    GetJobs()([]PrintJobable)
    GetLocation()(PrinterLocationable)
    GetManufacturer()(*string)
    GetModel()(*string)
    GetStatus()(PrinterStatusable)
    SetCapabilities(value PrinterCapabilitiesable)()
    SetDefaults(value PrinterDefaultsable)()
    SetDisplayName(value *string)()
    SetIsAcceptingJobs(value *bool)()
    SetJobs(value []PrintJobable)()
    SetLocation(value PrinterLocationable)()
    SetManufacturer(value *string)()
    SetModel(value *string)()
    SetStatus(value PrinterStatusable)()
}
