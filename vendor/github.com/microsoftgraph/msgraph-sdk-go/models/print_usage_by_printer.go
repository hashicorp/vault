package models

import (
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

type PrintUsageByPrinter struct {
    PrintUsage
}
// NewPrintUsageByPrinter instantiates a new PrintUsageByPrinter and sets the default values.
func NewPrintUsageByPrinter()(*PrintUsageByPrinter) {
    m := &PrintUsageByPrinter{
        PrintUsage: *NewPrintUsage(),
    }
    odataTypeValue := "#microsoft.graph.printUsageByPrinter"
    m.SetOdataType(&odataTypeValue)
    return m
}
// CreatePrintUsageByPrinterFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
// returns a Parsable when successful
func CreatePrintUsageByPrinterFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewPrintUsageByPrinter(), nil
}
// GetFieldDeserializers the deserialization information for the current model
// returns a map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error) when successful
func (m *PrintUsageByPrinter) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := m.PrintUsage.GetFieldDeserializers()
    res["printerId"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetPrinterId(val)
        }
        return nil
    }
    res["printerName"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetPrinterName(val)
        }
        return nil
    }
    return res
}
// GetPrinterId gets the printerId property value. The ID of the printer represented by these statistics.
// returns a *string when successful
func (m *PrintUsageByPrinter) GetPrinterId()(*string) {
    val, err := m.GetBackingStore().Get("printerId")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetPrinterName gets the printerName property value. The name of the printer represented by these statistics.
// returns a *string when successful
func (m *PrintUsageByPrinter) GetPrinterName()(*string) {
    val, err := m.GetBackingStore().Get("printerName")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// Serialize serializes information the current object
func (m *PrintUsageByPrinter) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    err := m.PrintUsage.Serialize(writer)
    if err != nil {
        return err
    }
    {
        err = writer.WriteStringValue("printerId", m.GetPrinterId())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteStringValue("printerName", m.GetPrinterName())
        if err != nil {
            return err
        }
    }
    return nil
}
// SetPrinterId sets the printerId property value. The ID of the printer represented by these statistics.
func (m *PrintUsageByPrinter) SetPrinterId(value *string)() {
    err := m.GetBackingStore().Set("printerId", value)
    if err != nil {
        panic(err)
    }
}
// SetPrinterName sets the printerName property value. The name of the printer represented by these statistics.
func (m *PrintUsageByPrinter) SetPrinterName(value *string)() {
    err := m.GetBackingStore().Set("printerName", value)
    if err != nil {
        panic(err)
    }
}
type PrintUsageByPrinterable interface {
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable
    PrintUsageable
    GetPrinterId()(*string)
    GetPrinterName()(*string)
    SetPrinterId(value *string)()
    SetPrinterName(value *string)()
}
