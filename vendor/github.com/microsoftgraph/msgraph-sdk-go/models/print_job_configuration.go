package models

import (
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
    ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e "github.com/microsoft/kiota-abstractions-go/store"
)

type PrintJobConfiguration struct {
    // Stores model information.
    backingStore ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackingStore
}
// NewPrintJobConfiguration instantiates a new PrintJobConfiguration and sets the default values.
func NewPrintJobConfiguration()(*PrintJobConfiguration) {
    m := &PrintJobConfiguration{
    }
    m.backingStore = ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackingStoreFactoryInstance();
    m.SetAdditionalData(make(map[string]any))
    return m
}
// CreatePrintJobConfigurationFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
// returns a Parsable when successful
func CreatePrintJobConfigurationFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewPrintJobConfiguration(), nil
}
// GetAdditionalData gets the AdditionalData property value. Stores additional data not described in the OpenAPI description found when deserializing. Can be used for serialization as well.
// returns a map[string]any when successful
func (m *PrintJobConfiguration) GetAdditionalData()(map[string]any) {
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
func (m *PrintJobConfiguration) GetBackingStore()(ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackingStore) {
    return m.backingStore
}
// GetCollate gets the collate property value. Whether the printer should collate pages wehen printing multiple copies of a multi-page document.
// returns a *bool when successful
func (m *PrintJobConfiguration) GetCollate()(*bool) {
    val, err := m.GetBackingStore().Get("collate")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*bool)
    }
    return nil
}
// GetColorMode gets the colorMode property value. The color mode the printer should use to print the job. Valid values are described in the table below. Read-only.
// returns a *PrintColorMode when successful
func (m *PrintJobConfiguration) GetColorMode()(*PrintColorMode) {
    val, err := m.GetBackingStore().Get("colorMode")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*PrintColorMode)
    }
    return nil
}
// GetCopies gets the copies property value. The number of copies that should be printed. Read-only.
// returns a *int32 when successful
func (m *PrintJobConfiguration) GetCopies()(*int32) {
    val, err := m.GetBackingStore().Get("copies")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*int32)
    }
    return nil
}
// GetDpi gets the dpi property value. The resolution to use when printing the job, expressed in dots per inch (DPI). Read-only.
// returns a *int32 when successful
func (m *PrintJobConfiguration) GetDpi()(*int32) {
    val, err := m.GetBackingStore().Get("dpi")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*int32)
    }
    return nil
}
// GetDuplexMode gets the duplexMode property value. The duplex mode the printer should use when printing the job. Valid values are described in the table below. Read-only.
// returns a *PrintDuplexMode when successful
func (m *PrintJobConfiguration) GetDuplexMode()(*PrintDuplexMode) {
    val, err := m.GetBackingStore().Get("duplexMode")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*PrintDuplexMode)
    }
    return nil
}
// GetFeedOrientation gets the feedOrientation property value. The orientation to use when feeding media into the printer. Valid values are described in the following table. Read-only.
// returns a *PrinterFeedOrientation when successful
func (m *PrintJobConfiguration) GetFeedOrientation()(*PrinterFeedOrientation) {
    val, err := m.GetBackingStore().Get("feedOrientation")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*PrinterFeedOrientation)
    }
    return nil
}
// GetFieldDeserializers the deserialization information for the current model
// returns a map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error) when successful
func (m *PrintJobConfiguration) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := make(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error))
    res["collate"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetBoolValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetCollate(val)
        }
        return nil
    }
    res["colorMode"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetEnumValue(ParsePrintColorMode)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetColorMode(val.(*PrintColorMode))
        }
        return nil
    }
    res["copies"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetInt32Value()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetCopies(val)
        }
        return nil
    }
    res["dpi"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetInt32Value()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetDpi(val)
        }
        return nil
    }
    res["duplexMode"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetEnumValue(ParsePrintDuplexMode)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetDuplexMode(val.(*PrintDuplexMode))
        }
        return nil
    }
    res["feedOrientation"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetEnumValue(ParsePrinterFeedOrientation)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetFeedOrientation(val.(*PrinterFeedOrientation))
        }
        return nil
    }
    res["finishings"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetCollectionOfEnumValues(ParsePrintFinishing)
        if err != nil {
            return err
        }
        if val != nil {
            res := make([]PrintFinishing, len(val))
            for i, v := range val {
                if v != nil {
                    res[i] = *(v.(*PrintFinishing))
                }
            }
            m.SetFinishings(res)
        }
        return nil
    }
    res["fitPdfToPage"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetBoolValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetFitPdfToPage(val)
        }
        return nil
    }
    res["inputBin"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetInputBin(val)
        }
        return nil
    }
    res["margin"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetObjectValue(CreatePrintMarginFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetMargin(val.(PrintMarginable))
        }
        return nil
    }
    res["mediaSize"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetMediaSize(val)
        }
        return nil
    }
    res["mediaType"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetMediaType(val)
        }
        return nil
    }
    res["multipageLayout"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetEnumValue(ParsePrintMultipageLayout)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetMultipageLayout(val.(*PrintMultipageLayout))
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
    res["orientation"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetEnumValue(ParsePrintOrientation)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetOrientation(val.(*PrintOrientation))
        }
        return nil
    }
    res["outputBin"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetOutputBin(val)
        }
        return nil
    }
    res["pageRanges"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetCollectionOfObjectValues(CreateIntegerRangeFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            res := make([]IntegerRangeable, len(val))
            for i, v := range val {
                if v != nil {
                    res[i] = v.(IntegerRangeable)
                }
            }
            m.SetPageRanges(res)
        }
        return nil
    }
    res["pagesPerSheet"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetInt32Value()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetPagesPerSheet(val)
        }
        return nil
    }
    res["quality"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetEnumValue(ParsePrintQuality)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetQuality(val.(*PrintQuality))
        }
        return nil
    }
    res["scaling"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetEnumValue(ParsePrintScaling)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetScaling(val.(*PrintScaling))
        }
        return nil
    }
    return res
}
// GetFinishings gets the finishings property value. Finishing processes to use when printing.
// returns a []PrintFinishing when successful
func (m *PrintJobConfiguration) GetFinishings()([]PrintFinishing) {
    val, err := m.GetBackingStore().Get("finishings")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]PrintFinishing)
    }
    return nil
}
// GetFitPdfToPage gets the fitPdfToPage property value. True to fit each page of a PDF document to a physical sheet of media; false to let the printer decide how to lay out impressions.
// returns a *bool when successful
func (m *PrintJobConfiguration) GetFitPdfToPage()(*bool) {
    val, err := m.GetBackingStore().Get("fitPdfToPage")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*bool)
    }
    return nil
}
// GetInputBin gets the inputBin property value. The input bin (tray) to use when printing. See the printer's capabilities for a list of supported input bins.
// returns a *string when successful
func (m *PrintJobConfiguration) GetInputBin()(*string) {
    val, err := m.GetBackingStore().Get("inputBin")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetMargin gets the margin property value. The margin settings to use when printing.
// returns a PrintMarginable when successful
func (m *PrintJobConfiguration) GetMargin()(PrintMarginable) {
    val, err := m.GetBackingStore().Get("margin")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(PrintMarginable)
    }
    return nil
}
// GetMediaSize gets the mediaSize property value. The media size to use when printing. Supports standard size names for ISO and ANSI media sizes. Valid values listed in the printerCapabilities topic.
// returns a *string when successful
func (m *PrintJobConfiguration) GetMediaSize()(*string) {
    val, err := m.GetBackingStore().Get("mediaSize")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetMediaType gets the mediaType property value. The default media (such as paper) type to print the document on.
// returns a *string when successful
func (m *PrintJobConfiguration) GetMediaType()(*string) {
    val, err := m.GetBackingStore().Get("mediaType")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetMultipageLayout gets the multipageLayout property value. The direction to lay out pages when multiple pages are being printed per sheet. Valid values are described in the following table.
// returns a *PrintMultipageLayout when successful
func (m *PrintJobConfiguration) GetMultipageLayout()(*PrintMultipageLayout) {
    val, err := m.GetBackingStore().Get("multipageLayout")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*PrintMultipageLayout)
    }
    return nil
}
// GetOdataType gets the @odata.type property value. The OdataType property
// returns a *string when successful
func (m *PrintJobConfiguration) GetOdataType()(*string) {
    val, err := m.GetBackingStore().Get("odataType")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetOrientation gets the orientation property value. The orientation setting the printer should use when printing the job. Valid values are described in the following table.
// returns a *PrintOrientation when successful
func (m *PrintJobConfiguration) GetOrientation()(*PrintOrientation) {
    val, err := m.GetBackingStore().Get("orientation")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*PrintOrientation)
    }
    return nil
}
// GetOutputBin gets the outputBin property value. The output bin to place completed prints into. See the printer's capabilities for a list of supported output bins.
// returns a *string when successful
func (m *PrintJobConfiguration) GetOutputBin()(*string) {
    val, err := m.GetBackingStore().Get("outputBin")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetPageRanges gets the pageRanges property value. The page ranges to print. Read-only.
// returns a []IntegerRangeable when successful
func (m *PrintJobConfiguration) GetPageRanges()([]IntegerRangeable) {
    val, err := m.GetBackingStore().Get("pageRanges")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]IntegerRangeable)
    }
    return nil
}
// GetPagesPerSheet gets the pagesPerSheet property value. The number of document pages to print on each sheet.
// returns a *int32 when successful
func (m *PrintJobConfiguration) GetPagesPerSheet()(*int32) {
    val, err := m.GetBackingStore().Get("pagesPerSheet")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*int32)
    }
    return nil
}
// GetQuality gets the quality property value. The print quality to use when printing the job. Valid values are described in the table below. Read-only.
// returns a *PrintQuality when successful
func (m *PrintJobConfiguration) GetQuality()(*PrintQuality) {
    val, err := m.GetBackingStore().Get("quality")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*PrintQuality)
    }
    return nil
}
// GetScaling gets the scaling property value. Specifies how the printer should scale the document data to fit the requested media. Valid values are described in the following table.
// returns a *PrintScaling when successful
func (m *PrintJobConfiguration) GetScaling()(*PrintScaling) {
    val, err := m.GetBackingStore().Get("scaling")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*PrintScaling)
    }
    return nil
}
// Serialize serializes information the current object
func (m *PrintJobConfiguration) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    {
        err := writer.WriteBoolValue("collate", m.GetCollate())
        if err != nil {
            return err
        }
    }
    if m.GetColorMode() != nil {
        cast := (*m.GetColorMode()).String()
        err := writer.WriteStringValue("colorMode", &cast)
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteInt32Value("copies", m.GetCopies())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteInt32Value("dpi", m.GetDpi())
        if err != nil {
            return err
        }
    }
    if m.GetDuplexMode() != nil {
        cast := (*m.GetDuplexMode()).String()
        err := writer.WriteStringValue("duplexMode", &cast)
        if err != nil {
            return err
        }
    }
    if m.GetFeedOrientation() != nil {
        cast := (*m.GetFeedOrientation()).String()
        err := writer.WriteStringValue("feedOrientation", &cast)
        if err != nil {
            return err
        }
    }
    if m.GetFinishings() != nil {
        err := writer.WriteCollectionOfStringValues("finishings", SerializePrintFinishing(m.GetFinishings()))
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteBoolValue("fitPdfToPage", m.GetFitPdfToPage())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteStringValue("inputBin", m.GetInputBin())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteObjectValue("margin", m.GetMargin())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteStringValue("mediaSize", m.GetMediaSize())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteStringValue("mediaType", m.GetMediaType())
        if err != nil {
            return err
        }
    }
    if m.GetMultipageLayout() != nil {
        cast := (*m.GetMultipageLayout()).String()
        err := writer.WriteStringValue("multipageLayout", &cast)
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
    if m.GetOrientation() != nil {
        cast := (*m.GetOrientation()).String()
        err := writer.WriteStringValue("orientation", &cast)
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteStringValue("outputBin", m.GetOutputBin())
        if err != nil {
            return err
        }
    }
    if m.GetPageRanges() != nil {
        cast := make([]i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, len(m.GetPageRanges()))
        for i, v := range m.GetPageRanges() {
            if v != nil {
                cast[i] = v.(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable)
            }
        }
        err := writer.WriteCollectionOfObjectValues("pageRanges", cast)
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteInt32Value("pagesPerSheet", m.GetPagesPerSheet())
        if err != nil {
            return err
        }
    }
    if m.GetQuality() != nil {
        cast := (*m.GetQuality()).String()
        err := writer.WriteStringValue("quality", &cast)
        if err != nil {
            return err
        }
    }
    if m.GetScaling() != nil {
        cast := (*m.GetScaling()).String()
        err := writer.WriteStringValue("scaling", &cast)
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
func (m *PrintJobConfiguration) SetAdditionalData(value map[string]any)() {
    err := m.GetBackingStore().Set("additionalData", value)
    if err != nil {
        panic(err)
    }
}
// SetBackingStore sets the BackingStore property value. Stores model information.
func (m *PrintJobConfiguration) SetBackingStore(value ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackingStore)() {
    m.backingStore = value
}
// SetCollate sets the collate property value. Whether the printer should collate pages wehen printing multiple copies of a multi-page document.
func (m *PrintJobConfiguration) SetCollate(value *bool)() {
    err := m.GetBackingStore().Set("collate", value)
    if err != nil {
        panic(err)
    }
}
// SetColorMode sets the colorMode property value. The color mode the printer should use to print the job. Valid values are described in the table below. Read-only.
func (m *PrintJobConfiguration) SetColorMode(value *PrintColorMode)() {
    err := m.GetBackingStore().Set("colorMode", value)
    if err != nil {
        panic(err)
    }
}
// SetCopies sets the copies property value. The number of copies that should be printed. Read-only.
func (m *PrintJobConfiguration) SetCopies(value *int32)() {
    err := m.GetBackingStore().Set("copies", value)
    if err != nil {
        panic(err)
    }
}
// SetDpi sets the dpi property value. The resolution to use when printing the job, expressed in dots per inch (DPI). Read-only.
func (m *PrintJobConfiguration) SetDpi(value *int32)() {
    err := m.GetBackingStore().Set("dpi", value)
    if err != nil {
        panic(err)
    }
}
// SetDuplexMode sets the duplexMode property value. The duplex mode the printer should use when printing the job. Valid values are described in the table below. Read-only.
func (m *PrintJobConfiguration) SetDuplexMode(value *PrintDuplexMode)() {
    err := m.GetBackingStore().Set("duplexMode", value)
    if err != nil {
        panic(err)
    }
}
// SetFeedOrientation sets the feedOrientation property value. The orientation to use when feeding media into the printer. Valid values are described in the following table. Read-only.
func (m *PrintJobConfiguration) SetFeedOrientation(value *PrinterFeedOrientation)() {
    err := m.GetBackingStore().Set("feedOrientation", value)
    if err != nil {
        panic(err)
    }
}
// SetFinishings sets the finishings property value. Finishing processes to use when printing.
func (m *PrintJobConfiguration) SetFinishings(value []PrintFinishing)() {
    err := m.GetBackingStore().Set("finishings", value)
    if err != nil {
        panic(err)
    }
}
// SetFitPdfToPage sets the fitPdfToPage property value. True to fit each page of a PDF document to a physical sheet of media; false to let the printer decide how to lay out impressions.
func (m *PrintJobConfiguration) SetFitPdfToPage(value *bool)() {
    err := m.GetBackingStore().Set("fitPdfToPage", value)
    if err != nil {
        panic(err)
    }
}
// SetInputBin sets the inputBin property value. The input bin (tray) to use when printing. See the printer's capabilities for a list of supported input bins.
func (m *PrintJobConfiguration) SetInputBin(value *string)() {
    err := m.GetBackingStore().Set("inputBin", value)
    if err != nil {
        panic(err)
    }
}
// SetMargin sets the margin property value. The margin settings to use when printing.
func (m *PrintJobConfiguration) SetMargin(value PrintMarginable)() {
    err := m.GetBackingStore().Set("margin", value)
    if err != nil {
        panic(err)
    }
}
// SetMediaSize sets the mediaSize property value. The media size to use when printing. Supports standard size names for ISO and ANSI media sizes. Valid values listed in the printerCapabilities topic.
func (m *PrintJobConfiguration) SetMediaSize(value *string)() {
    err := m.GetBackingStore().Set("mediaSize", value)
    if err != nil {
        panic(err)
    }
}
// SetMediaType sets the mediaType property value. The default media (such as paper) type to print the document on.
func (m *PrintJobConfiguration) SetMediaType(value *string)() {
    err := m.GetBackingStore().Set("mediaType", value)
    if err != nil {
        panic(err)
    }
}
// SetMultipageLayout sets the multipageLayout property value. The direction to lay out pages when multiple pages are being printed per sheet. Valid values are described in the following table.
func (m *PrintJobConfiguration) SetMultipageLayout(value *PrintMultipageLayout)() {
    err := m.GetBackingStore().Set("multipageLayout", value)
    if err != nil {
        panic(err)
    }
}
// SetOdataType sets the @odata.type property value. The OdataType property
func (m *PrintJobConfiguration) SetOdataType(value *string)() {
    err := m.GetBackingStore().Set("odataType", value)
    if err != nil {
        panic(err)
    }
}
// SetOrientation sets the orientation property value. The orientation setting the printer should use when printing the job. Valid values are described in the following table.
func (m *PrintJobConfiguration) SetOrientation(value *PrintOrientation)() {
    err := m.GetBackingStore().Set("orientation", value)
    if err != nil {
        panic(err)
    }
}
// SetOutputBin sets the outputBin property value. The output bin to place completed prints into. See the printer's capabilities for a list of supported output bins.
func (m *PrintJobConfiguration) SetOutputBin(value *string)() {
    err := m.GetBackingStore().Set("outputBin", value)
    if err != nil {
        panic(err)
    }
}
// SetPageRanges sets the pageRanges property value. The page ranges to print. Read-only.
func (m *PrintJobConfiguration) SetPageRanges(value []IntegerRangeable)() {
    err := m.GetBackingStore().Set("pageRanges", value)
    if err != nil {
        panic(err)
    }
}
// SetPagesPerSheet sets the pagesPerSheet property value. The number of document pages to print on each sheet.
func (m *PrintJobConfiguration) SetPagesPerSheet(value *int32)() {
    err := m.GetBackingStore().Set("pagesPerSheet", value)
    if err != nil {
        panic(err)
    }
}
// SetQuality sets the quality property value. The print quality to use when printing the job. Valid values are described in the table below. Read-only.
func (m *PrintJobConfiguration) SetQuality(value *PrintQuality)() {
    err := m.GetBackingStore().Set("quality", value)
    if err != nil {
        panic(err)
    }
}
// SetScaling sets the scaling property value. Specifies how the printer should scale the document data to fit the requested media. Valid values are described in the following table.
func (m *PrintJobConfiguration) SetScaling(value *PrintScaling)() {
    err := m.GetBackingStore().Set("scaling", value)
    if err != nil {
        panic(err)
    }
}
type PrintJobConfigurationable interface {
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.AdditionalDataHolder
    ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackedModel
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable
    GetBackingStore()(ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackingStore)
    GetCollate()(*bool)
    GetColorMode()(*PrintColorMode)
    GetCopies()(*int32)
    GetDpi()(*int32)
    GetDuplexMode()(*PrintDuplexMode)
    GetFeedOrientation()(*PrinterFeedOrientation)
    GetFinishings()([]PrintFinishing)
    GetFitPdfToPage()(*bool)
    GetInputBin()(*string)
    GetMargin()(PrintMarginable)
    GetMediaSize()(*string)
    GetMediaType()(*string)
    GetMultipageLayout()(*PrintMultipageLayout)
    GetOdataType()(*string)
    GetOrientation()(*PrintOrientation)
    GetOutputBin()(*string)
    GetPageRanges()([]IntegerRangeable)
    GetPagesPerSheet()(*int32)
    GetQuality()(*PrintQuality)
    GetScaling()(*PrintScaling)
    SetBackingStore(value ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackingStore)()
    SetCollate(value *bool)()
    SetColorMode(value *PrintColorMode)()
    SetCopies(value *int32)()
    SetDpi(value *int32)()
    SetDuplexMode(value *PrintDuplexMode)()
    SetFeedOrientation(value *PrinterFeedOrientation)()
    SetFinishings(value []PrintFinishing)()
    SetFitPdfToPage(value *bool)()
    SetInputBin(value *string)()
    SetMargin(value PrintMarginable)()
    SetMediaSize(value *string)()
    SetMediaType(value *string)()
    SetMultipageLayout(value *PrintMultipageLayout)()
    SetOdataType(value *string)()
    SetOrientation(value *PrintOrientation)()
    SetOutputBin(value *string)()
    SetPageRanges(value []IntegerRangeable)()
    SetPagesPerSheet(value *int32)()
    SetQuality(value *PrintQuality)()
    SetScaling(value *PrintScaling)()
}
