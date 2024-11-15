package models

import (
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
    ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e "github.com/microsoft/kiota-abstractions-go/store"
)

type Print struct {
    // Stores model information.
    backingStore ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackingStore
}
// NewPrint instantiates a new Print and sets the default values.
func NewPrint()(*Print) {
    m := &Print{
    }
    m.backingStore = ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackingStoreFactoryInstance();
    m.SetAdditionalData(make(map[string]any))
    return m
}
// CreatePrintFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
// returns a Parsable when successful
func CreatePrintFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewPrint(), nil
}
// GetAdditionalData gets the AdditionalData property value. Stores additional data not described in the OpenAPI description found when deserializing. Can be used for serialization as well.
// returns a map[string]any when successful
func (m *Print) GetAdditionalData()(map[string]any) {
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
func (m *Print) GetBackingStore()(ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackingStore) {
    return m.backingStore
}
// GetConnectors gets the connectors property value. The list of available print connectors.
// returns a []PrintConnectorable when successful
func (m *Print) GetConnectors()([]PrintConnectorable) {
    val, err := m.GetBackingStore().Get("connectors")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]PrintConnectorable)
    }
    return nil
}
// GetFieldDeserializers the deserialization information for the current model
// returns a map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error) when successful
func (m *Print) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := make(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error))
    res["connectors"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetCollectionOfObjectValues(CreatePrintConnectorFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            res := make([]PrintConnectorable, len(val))
            for i, v := range val {
                if v != nil {
                    res[i] = v.(PrintConnectorable)
                }
            }
            m.SetConnectors(res)
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
    res["operations"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetCollectionOfObjectValues(CreatePrintOperationFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            res := make([]PrintOperationable, len(val))
            for i, v := range val {
                if v != nil {
                    res[i] = v.(PrintOperationable)
                }
            }
            m.SetOperations(res)
        }
        return nil
    }
    res["printers"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetCollectionOfObjectValues(CreatePrinterFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            res := make([]Printerable, len(val))
            for i, v := range val {
                if v != nil {
                    res[i] = v.(Printerable)
                }
            }
            m.SetPrinters(res)
        }
        return nil
    }
    res["services"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetCollectionOfObjectValues(CreatePrintServiceFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            res := make([]PrintServiceable, len(val))
            for i, v := range val {
                if v != nil {
                    res[i] = v.(PrintServiceable)
                }
            }
            m.SetServices(res)
        }
        return nil
    }
    res["settings"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetObjectValue(CreatePrintSettingsFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetSettings(val.(PrintSettingsable))
        }
        return nil
    }
    res["shares"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetCollectionOfObjectValues(CreatePrinterShareFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            res := make([]PrinterShareable, len(val))
            for i, v := range val {
                if v != nil {
                    res[i] = v.(PrinterShareable)
                }
            }
            m.SetShares(res)
        }
        return nil
    }
    res["taskDefinitions"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetCollectionOfObjectValues(CreatePrintTaskDefinitionFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            res := make([]PrintTaskDefinitionable, len(val))
            for i, v := range val {
                if v != nil {
                    res[i] = v.(PrintTaskDefinitionable)
                }
            }
            m.SetTaskDefinitions(res)
        }
        return nil
    }
    return res
}
// GetOdataType gets the @odata.type property value. The OdataType property
// returns a *string when successful
func (m *Print) GetOdataType()(*string) {
    val, err := m.GetBackingStore().Get("odataType")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetOperations gets the operations property value. The list of print long running operations.
// returns a []PrintOperationable when successful
func (m *Print) GetOperations()([]PrintOperationable) {
    val, err := m.GetBackingStore().Get("operations")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]PrintOperationable)
    }
    return nil
}
// GetPrinters gets the printers property value. The list of printers registered in the tenant.
// returns a []Printerable when successful
func (m *Print) GetPrinters()([]Printerable) {
    val, err := m.GetBackingStore().Get("printers")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]Printerable)
    }
    return nil
}
// GetServices gets the services property value. The list of available Universal Print service endpoints.
// returns a []PrintServiceable when successful
func (m *Print) GetServices()([]PrintServiceable) {
    val, err := m.GetBackingStore().Get("services")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]PrintServiceable)
    }
    return nil
}
// GetSettings gets the settings property value. Tenant-wide settings for the Universal Print service.
// returns a PrintSettingsable when successful
func (m *Print) GetSettings()(PrintSettingsable) {
    val, err := m.GetBackingStore().Get("settings")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(PrintSettingsable)
    }
    return nil
}
// GetShares gets the shares property value. The list of printer shares registered in the tenant.
// returns a []PrinterShareable when successful
func (m *Print) GetShares()([]PrinterShareable) {
    val, err := m.GetBackingStore().Get("shares")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]PrinterShareable)
    }
    return nil
}
// GetTaskDefinitions gets the taskDefinitions property value. List of abstract definition for a task that can be triggered when various events occur within Universal Print.
// returns a []PrintTaskDefinitionable when successful
func (m *Print) GetTaskDefinitions()([]PrintTaskDefinitionable) {
    val, err := m.GetBackingStore().Get("taskDefinitions")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]PrintTaskDefinitionable)
    }
    return nil
}
// Serialize serializes information the current object
func (m *Print) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    if m.GetConnectors() != nil {
        cast := make([]i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, len(m.GetConnectors()))
        for i, v := range m.GetConnectors() {
            if v != nil {
                cast[i] = v.(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable)
            }
        }
        err := writer.WriteCollectionOfObjectValues("connectors", cast)
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
    if m.GetOperations() != nil {
        cast := make([]i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, len(m.GetOperations()))
        for i, v := range m.GetOperations() {
            if v != nil {
                cast[i] = v.(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable)
            }
        }
        err := writer.WriteCollectionOfObjectValues("operations", cast)
        if err != nil {
            return err
        }
    }
    if m.GetPrinters() != nil {
        cast := make([]i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, len(m.GetPrinters()))
        for i, v := range m.GetPrinters() {
            if v != nil {
                cast[i] = v.(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable)
            }
        }
        err := writer.WriteCollectionOfObjectValues("printers", cast)
        if err != nil {
            return err
        }
    }
    if m.GetServices() != nil {
        cast := make([]i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, len(m.GetServices()))
        for i, v := range m.GetServices() {
            if v != nil {
                cast[i] = v.(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable)
            }
        }
        err := writer.WriteCollectionOfObjectValues("services", cast)
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteObjectValue("settings", m.GetSettings())
        if err != nil {
            return err
        }
    }
    if m.GetShares() != nil {
        cast := make([]i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, len(m.GetShares()))
        for i, v := range m.GetShares() {
            if v != nil {
                cast[i] = v.(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable)
            }
        }
        err := writer.WriteCollectionOfObjectValues("shares", cast)
        if err != nil {
            return err
        }
    }
    if m.GetTaskDefinitions() != nil {
        cast := make([]i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, len(m.GetTaskDefinitions()))
        for i, v := range m.GetTaskDefinitions() {
            if v != nil {
                cast[i] = v.(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable)
            }
        }
        err := writer.WriteCollectionOfObjectValues("taskDefinitions", cast)
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
func (m *Print) SetAdditionalData(value map[string]any)() {
    err := m.GetBackingStore().Set("additionalData", value)
    if err != nil {
        panic(err)
    }
}
// SetBackingStore sets the BackingStore property value. Stores model information.
func (m *Print) SetBackingStore(value ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackingStore)() {
    m.backingStore = value
}
// SetConnectors sets the connectors property value. The list of available print connectors.
func (m *Print) SetConnectors(value []PrintConnectorable)() {
    err := m.GetBackingStore().Set("connectors", value)
    if err != nil {
        panic(err)
    }
}
// SetOdataType sets the @odata.type property value. The OdataType property
func (m *Print) SetOdataType(value *string)() {
    err := m.GetBackingStore().Set("odataType", value)
    if err != nil {
        panic(err)
    }
}
// SetOperations sets the operations property value. The list of print long running operations.
func (m *Print) SetOperations(value []PrintOperationable)() {
    err := m.GetBackingStore().Set("operations", value)
    if err != nil {
        panic(err)
    }
}
// SetPrinters sets the printers property value. The list of printers registered in the tenant.
func (m *Print) SetPrinters(value []Printerable)() {
    err := m.GetBackingStore().Set("printers", value)
    if err != nil {
        panic(err)
    }
}
// SetServices sets the services property value. The list of available Universal Print service endpoints.
func (m *Print) SetServices(value []PrintServiceable)() {
    err := m.GetBackingStore().Set("services", value)
    if err != nil {
        panic(err)
    }
}
// SetSettings sets the settings property value. Tenant-wide settings for the Universal Print service.
func (m *Print) SetSettings(value PrintSettingsable)() {
    err := m.GetBackingStore().Set("settings", value)
    if err != nil {
        panic(err)
    }
}
// SetShares sets the shares property value. The list of printer shares registered in the tenant.
func (m *Print) SetShares(value []PrinterShareable)() {
    err := m.GetBackingStore().Set("shares", value)
    if err != nil {
        panic(err)
    }
}
// SetTaskDefinitions sets the taskDefinitions property value. List of abstract definition for a task that can be triggered when various events occur within Universal Print.
func (m *Print) SetTaskDefinitions(value []PrintTaskDefinitionable)() {
    err := m.GetBackingStore().Set("taskDefinitions", value)
    if err != nil {
        panic(err)
    }
}
type Printable interface {
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.AdditionalDataHolder
    ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackedModel
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable
    GetBackingStore()(ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackingStore)
    GetConnectors()([]PrintConnectorable)
    GetOdataType()(*string)
    GetOperations()([]PrintOperationable)
    GetPrinters()([]Printerable)
    GetServices()([]PrintServiceable)
    GetSettings()(PrintSettingsable)
    GetShares()([]PrinterShareable)
    GetTaskDefinitions()([]PrintTaskDefinitionable)
    SetBackingStore(value ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackingStore)()
    SetConnectors(value []PrintConnectorable)()
    SetOdataType(value *string)()
    SetOperations(value []PrintOperationable)()
    SetPrinters(value []Printerable)()
    SetServices(value []PrintServiceable)()
    SetSettings(value PrintSettingsable)()
    SetShares(value []PrinterShareable)()
    SetTaskDefinitions(value []PrintTaskDefinitionable)()
}
