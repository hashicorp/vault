package models

import (
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
    ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e "github.com/microsoft/kiota-abstractions-go/store"
)

type PrinterCapabilities struct {
    // Stores model information.
    backingStore ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackingStore
}
// NewPrinterCapabilities instantiates a new PrinterCapabilities and sets the default values.
func NewPrinterCapabilities()(*PrinterCapabilities) {
    m := &PrinterCapabilities{
    }
    m.backingStore = ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackingStoreFactoryInstance();
    m.SetAdditionalData(make(map[string]any))
    return m
}
// CreatePrinterCapabilitiesFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
// returns a Parsable when successful
func CreatePrinterCapabilitiesFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewPrinterCapabilities(), nil
}
// GetAdditionalData gets the AdditionalData property value. Stores additional data not described in the OpenAPI description found when deserializing. Can be used for serialization as well.
// returns a map[string]any when successful
func (m *PrinterCapabilities) GetAdditionalData()(map[string]any) {
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
func (m *PrinterCapabilities) GetBackingStore()(ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackingStore) {
    return m.backingStore
}
// GetBottomMargins gets the bottomMargins property value. A list of supported bottom margins(in microns) for the printer.
// returns a []int32 when successful
func (m *PrinterCapabilities) GetBottomMargins()([]int32) {
    val, err := m.GetBackingStore().Get("bottomMargins")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]int32)
    }
    return nil
}
// GetCollation gets the collation property value. True if the printer supports collating when printing muliple copies of a multi-page document; false otherwise.
// returns a *bool when successful
func (m *PrinterCapabilities) GetCollation()(*bool) {
    val, err := m.GetBackingStore().Get("collation")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*bool)
    }
    return nil
}
// GetColorModes gets the colorModes property value. The color modes supported by the printer. Valid values are described in the following table.
// returns a []PrintColorMode when successful
func (m *PrinterCapabilities) GetColorModes()([]PrintColorMode) {
    val, err := m.GetBackingStore().Get("colorModes")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]PrintColorMode)
    }
    return nil
}
// GetContentTypes gets the contentTypes property value. A list of supported content (MIME) types that the printer supports. It is not guaranteed that the Universal Print service supports printing all of these MIME types.
// returns a []string when successful
func (m *PrinterCapabilities) GetContentTypes()([]string) {
    val, err := m.GetBackingStore().Get("contentTypes")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]string)
    }
    return nil
}
// GetCopiesPerJob gets the copiesPerJob property value. The range of copies per job supported by the printer.
// returns a IntegerRangeable when successful
func (m *PrinterCapabilities) GetCopiesPerJob()(IntegerRangeable) {
    val, err := m.GetBackingStore().Get("copiesPerJob")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(IntegerRangeable)
    }
    return nil
}
// GetDpis gets the dpis property value. The list of print resolutions in DPI that are supported by the printer.
// returns a []int32 when successful
func (m *PrinterCapabilities) GetDpis()([]int32) {
    val, err := m.GetBackingStore().Get("dpis")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]int32)
    }
    return nil
}
// GetDuplexModes gets the duplexModes property value. The list of duplex modes that are supported by the printer. Valid values are described in the following table.
// returns a []PrintDuplexMode when successful
func (m *PrinterCapabilities) GetDuplexModes()([]PrintDuplexMode) {
    val, err := m.GetBackingStore().Get("duplexModes")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]PrintDuplexMode)
    }
    return nil
}
// GetFeedOrientations gets the feedOrientations property value. The list of feed orientations that are supported by the printer.
// returns a []PrinterFeedOrientation when successful
func (m *PrinterCapabilities) GetFeedOrientations()([]PrinterFeedOrientation) {
    val, err := m.GetBackingStore().Get("feedOrientations")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]PrinterFeedOrientation)
    }
    return nil
}
// GetFieldDeserializers the deserialization information for the current model
// returns a map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error) when successful
func (m *PrinterCapabilities) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := make(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error))
    res["bottomMargins"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetCollectionOfPrimitiveValues("int32")
        if err != nil {
            return err
        }
        if val != nil {
            res := make([]int32, len(val))
            for i, v := range val {
                if v != nil {
                    res[i] = *(v.(*int32))
                }
            }
            m.SetBottomMargins(res)
        }
        return nil
    }
    res["collation"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetBoolValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetCollation(val)
        }
        return nil
    }
    res["colorModes"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetCollectionOfEnumValues(ParsePrintColorMode)
        if err != nil {
            return err
        }
        if val != nil {
            res := make([]PrintColorMode, len(val))
            for i, v := range val {
                if v != nil {
                    res[i] = *(v.(*PrintColorMode))
                }
            }
            m.SetColorModes(res)
        }
        return nil
    }
    res["contentTypes"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetCollectionOfPrimitiveValues("string")
        if err != nil {
            return err
        }
        if val != nil {
            res := make([]string, len(val))
            for i, v := range val {
                if v != nil {
                    res[i] = *(v.(*string))
                }
            }
            m.SetContentTypes(res)
        }
        return nil
    }
    res["copiesPerJob"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetObjectValue(CreateIntegerRangeFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetCopiesPerJob(val.(IntegerRangeable))
        }
        return nil
    }
    res["dpis"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetCollectionOfPrimitiveValues("int32")
        if err != nil {
            return err
        }
        if val != nil {
            res := make([]int32, len(val))
            for i, v := range val {
                if v != nil {
                    res[i] = *(v.(*int32))
                }
            }
            m.SetDpis(res)
        }
        return nil
    }
    res["duplexModes"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetCollectionOfEnumValues(ParsePrintDuplexMode)
        if err != nil {
            return err
        }
        if val != nil {
            res := make([]PrintDuplexMode, len(val))
            for i, v := range val {
                if v != nil {
                    res[i] = *(v.(*PrintDuplexMode))
                }
            }
            m.SetDuplexModes(res)
        }
        return nil
    }
    res["feedOrientations"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetCollectionOfEnumValues(ParsePrinterFeedOrientation)
        if err != nil {
            return err
        }
        if val != nil {
            res := make([]PrinterFeedOrientation, len(val))
            for i, v := range val {
                if v != nil {
                    res[i] = *(v.(*PrinterFeedOrientation))
                }
            }
            m.SetFeedOrientations(res)
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
    res["inputBins"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetCollectionOfPrimitiveValues("string")
        if err != nil {
            return err
        }
        if val != nil {
            res := make([]string, len(val))
            for i, v := range val {
                if v != nil {
                    res[i] = *(v.(*string))
                }
            }
            m.SetInputBins(res)
        }
        return nil
    }
    res["isColorPrintingSupported"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetBoolValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetIsColorPrintingSupported(val)
        }
        return nil
    }
    res["isPageRangeSupported"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetBoolValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetIsPageRangeSupported(val)
        }
        return nil
    }
    res["leftMargins"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetCollectionOfPrimitiveValues("int32")
        if err != nil {
            return err
        }
        if val != nil {
            res := make([]int32, len(val))
            for i, v := range val {
                if v != nil {
                    res[i] = *(v.(*int32))
                }
            }
            m.SetLeftMargins(res)
        }
        return nil
    }
    res["mediaColors"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetCollectionOfPrimitiveValues("string")
        if err != nil {
            return err
        }
        if val != nil {
            res := make([]string, len(val))
            for i, v := range val {
                if v != nil {
                    res[i] = *(v.(*string))
                }
            }
            m.SetMediaColors(res)
        }
        return nil
    }
    res["mediaSizes"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetCollectionOfPrimitiveValues("string")
        if err != nil {
            return err
        }
        if val != nil {
            res := make([]string, len(val))
            for i, v := range val {
                if v != nil {
                    res[i] = *(v.(*string))
                }
            }
            m.SetMediaSizes(res)
        }
        return nil
    }
    res["mediaTypes"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetCollectionOfPrimitiveValues("string")
        if err != nil {
            return err
        }
        if val != nil {
            res := make([]string, len(val))
            for i, v := range val {
                if v != nil {
                    res[i] = *(v.(*string))
                }
            }
            m.SetMediaTypes(res)
        }
        return nil
    }
    res["multipageLayouts"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetCollectionOfEnumValues(ParsePrintMultipageLayout)
        if err != nil {
            return err
        }
        if val != nil {
            res := make([]PrintMultipageLayout, len(val))
            for i, v := range val {
                if v != nil {
                    res[i] = *(v.(*PrintMultipageLayout))
                }
            }
            m.SetMultipageLayouts(res)
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
    res["orientations"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetCollectionOfEnumValues(ParsePrintOrientation)
        if err != nil {
            return err
        }
        if val != nil {
            res := make([]PrintOrientation, len(val))
            for i, v := range val {
                if v != nil {
                    res[i] = *(v.(*PrintOrientation))
                }
            }
            m.SetOrientations(res)
        }
        return nil
    }
    res["outputBins"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetCollectionOfPrimitiveValues("string")
        if err != nil {
            return err
        }
        if val != nil {
            res := make([]string, len(val))
            for i, v := range val {
                if v != nil {
                    res[i] = *(v.(*string))
                }
            }
            m.SetOutputBins(res)
        }
        return nil
    }
    res["pagesPerSheet"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetCollectionOfPrimitiveValues("int32")
        if err != nil {
            return err
        }
        if val != nil {
            res := make([]int32, len(val))
            for i, v := range val {
                if v != nil {
                    res[i] = *(v.(*int32))
                }
            }
            m.SetPagesPerSheet(res)
        }
        return nil
    }
    res["qualities"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetCollectionOfEnumValues(ParsePrintQuality)
        if err != nil {
            return err
        }
        if val != nil {
            res := make([]PrintQuality, len(val))
            for i, v := range val {
                if v != nil {
                    res[i] = *(v.(*PrintQuality))
                }
            }
            m.SetQualities(res)
        }
        return nil
    }
    res["rightMargins"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetCollectionOfPrimitiveValues("int32")
        if err != nil {
            return err
        }
        if val != nil {
            res := make([]int32, len(val))
            for i, v := range val {
                if v != nil {
                    res[i] = *(v.(*int32))
                }
            }
            m.SetRightMargins(res)
        }
        return nil
    }
    res["scalings"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetCollectionOfEnumValues(ParsePrintScaling)
        if err != nil {
            return err
        }
        if val != nil {
            res := make([]PrintScaling, len(val))
            for i, v := range val {
                if v != nil {
                    res[i] = *(v.(*PrintScaling))
                }
            }
            m.SetScalings(res)
        }
        return nil
    }
    res["supportsFitPdfToPage"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetBoolValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetSupportsFitPdfToPage(val)
        }
        return nil
    }
    res["topMargins"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetCollectionOfPrimitiveValues("int32")
        if err != nil {
            return err
        }
        if val != nil {
            res := make([]int32, len(val))
            for i, v := range val {
                if v != nil {
                    res[i] = *(v.(*int32))
                }
            }
            m.SetTopMargins(res)
        }
        return nil
    }
    return res
}
// GetFinishings gets the finishings property value. Finishing processes the printer supports for a printed document.
// returns a []PrintFinishing when successful
func (m *PrinterCapabilities) GetFinishings()([]PrintFinishing) {
    val, err := m.GetBackingStore().Get("finishings")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]PrintFinishing)
    }
    return nil
}
// GetInputBins gets the inputBins property value. Supported input bins for the printer.
// returns a []string when successful
func (m *PrinterCapabilities) GetInputBins()([]string) {
    val, err := m.GetBackingStore().Get("inputBins")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]string)
    }
    return nil
}
// GetIsColorPrintingSupported gets the isColorPrintingSupported property value. True if color printing is supported by the printer; false otherwise. Read-only.
// returns a *bool when successful
func (m *PrinterCapabilities) GetIsColorPrintingSupported()(*bool) {
    val, err := m.GetBackingStore().Get("isColorPrintingSupported")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*bool)
    }
    return nil
}
// GetIsPageRangeSupported gets the isPageRangeSupported property value. True if the printer supports printing by page ranges; false otherwise.
// returns a *bool when successful
func (m *PrinterCapabilities) GetIsPageRangeSupported()(*bool) {
    val, err := m.GetBackingStore().Get("isPageRangeSupported")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*bool)
    }
    return nil
}
// GetLeftMargins gets the leftMargins property value. A list of supported left margins(in microns) for the printer.
// returns a []int32 when successful
func (m *PrinterCapabilities) GetLeftMargins()([]int32) {
    val, err := m.GetBackingStore().Get("leftMargins")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]int32)
    }
    return nil
}
// GetMediaColors gets the mediaColors property value. The media (i.e., paper) colors supported by the printer.
// returns a []string when successful
func (m *PrinterCapabilities) GetMediaColors()([]string) {
    val, err := m.GetBackingStore().Get("mediaColors")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]string)
    }
    return nil
}
// GetMediaSizes gets the mediaSizes property value. The media sizes supported by the printer. Supports standard size names for ISO and ANSI media sizes. Valid values are in the following table.
// returns a []string when successful
func (m *PrinterCapabilities) GetMediaSizes()([]string) {
    val, err := m.GetBackingStore().Get("mediaSizes")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]string)
    }
    return nil
}
// GetMediaTypes gets the mediaTypes property value. The media types supported by the printer.
// returns a []string when successful
func (m *PrinterCapabilities) GetMediaTypes()([]string) {
    val, err := m.GetBackingStore().Get("mediaTypes")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]string)
    }
    return nil
}
// GetMultipageLayouts gets the multipageLayouts property value. The presentation directions supported by the printer. Supported values are described in the following table.
// returns a []PrintMultipageLayout when successful
func (m *PrinterCapabilities) GetMultipageLayouts()([]PrintMultipageLayout) {
    val, err := m.GetBackingStore().Get("multipageLayouts")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]PrintMultipageLayout)
    }
    return nil
}
// GetOdataType gets the @odata.type property value. The OdataType property
// returns a *string when successful
func (m *PrinterCapabilities) GetOdataType()(*string) {
    val, err := m.GetBackingStore().Get("odataType")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetOrientations gets the orientations property value. The print orientations supported by the printer. Valid values are described in the following table.
// returns a []PrintOrientation when successful
func (m *PrinterCapabilities) GetOrientations()([]PrintOrientation) {
    val, err := m.GetBackingStore().Get("orientations")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]PrintOrientation)
    }
    return nil
}
// GetOutputBins gets the outputBins property value. The printer's supported output bins (trays).
// returns a []string when successful
func (m *PrinterCapabilities) GetOutputBins()([]string) {
    val, err := m.GetBackingStore().Get("outputBins")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]string)
    }
    return nil
}
// GetPagesPerSheet gets the pagesPerSheet property value. Supported number of Input Pages to impose upon a single Impression.
// returns a []int32 when successful
func (m *PrinterCapabilities) GetPagesPerSheet()([]int32) {
    val, err := m.GetBackingStore().Get("pagesPerSheet")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]int32)
    }
    return nil
}
// GetQualities gets the qualities property value. The print qualities supported by the printer.
// returns a []PrintQuality when successful
func (m *PrinterCapabilities) GetQualities()([]PrintQuality) {
    val, err := m.GetBackingStore().Get("qualities")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]PrintQuality)
    }
    return nil
}
// GetRightMargins gets the rightMargins property value. A list of supported right margins(in microns) for the printer.
// returns a []int32 when successful
func (m *PrinterCapabilities) GetRightMargins()([]int32) {
    val, err := m.GetBackingStore().Get("rightMargins")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]int32)
    }
    return nil
}
// GetScalings gets the scalings property value. Supported print scalings.
// returns a []PrintScaling when successful
func (m *PrinterCapabilities) GetScalings()([]PrintScaling) {
    val, err := m.GetBackingStore().Get("scalings")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]PrintScaling)
    }
    return nil
}
// GetSupportsFitPdfToPage gets the supportsFitPdfToPage property value. True if the printer supports scaling PDF pages to match the print media size; false otherwise.
// returns a *bool when successful
func (m *PrinterCapabilities) GetSupportsFitPdfToPage()(*bool) {
    val, err := m.GetBackingStore().Get("supportsFitPdfToPage")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*bool)
    }
    return nil
}
// GetTopMargins gets the topMargins property value. A list of supported top margins(in microns) for the printer.
// returns a []int32 when successful
func (m *PrinterCapabilities) GetTopMargins()([]int32) {
    val, err := m.GetBackingStore().Get("topMargins")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]int32)
    }
    return nil
}
// Serialize serializes information the current object
func (m *PrinterCapabilities) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    if m.GetBottomMargins() != nil {
        err := writer.WriteCollectionOfInt32Values("bottomMargins", m.GetBottomMargins())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteBoolValue("collation", m.GetCollation())
        if err != nil {
            return err
        }
    }
    if m.GetColorModes() != nil {
        err := writer.WriteCollectionOfStringValues("colorModes", SerializePrintColorMode(m.GetColorModes()))
        if err != nil {
            return err
        }
    }
    if m.GetContentTypes() != nil {
        err := writer.WriteCollectionOfStringValues("contentTypes", m.GetContentTypes())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteObjectValue("copiesPerJob", m.GetCopiesPerJob())
        if err != nil {
            return err
        }
    }
    if m.GetDpis() != nil {
        err := writer.WriteCollectionOfInt32Values("dpis", m.GetDpis())
        if err != nil {
            return err
        }
    }
    if m.GetDuplexModes() != nil {
        err := writer.WriteCollectionOfStringValues("duplexModes", SerializePrintDuplexMode(m.GetDuplexModes()))
        if err != nil {
            return err
        }
    }
    if m.GetFeedOrientations() != nil {
        err := writer.WriteCollectionOfStringValues("feedOrientations", SerializePrinterFeedOrientation(m.GetFeedOrientations()))
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
    if m.GetInputBins() != nil {
        err := writer.WriteCollectionOfStringValues("inputBins", m.GetInputBins())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteBoolValue("isColorPrintingSupported", m.GetIsColorPrintingSupported())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteBoolValue("isPageRangeSupported", m.GetIsPageRangeSupported())
        if err != nil {
            return err
        }
    }
    if m.GetLeftMargins() != nil {
        err := writer.WriteCollectionOfInt32Values("leftMargins", m.GetLeftMargins())
        if err != nil {
            return err
        }
    }
    if m.GetMediaColors() != nil {
        err := writer.WriteCollectionOfStringValues("mediaColors", m.GetMediaColors())
        if err != nil {
            return err
        }
    }
    if m.GetMediaSizes() != nil {
        err := writer.WriteCollectionOfStringValues("mediaSizes", m.GetMediaSizes())
        if err != nil {
            return err
        }
    }
    if m.GetMediaTypes() != nil {
        err := writer.WriteCollectionOfStringValues("mediaTypes", m.GetMediaTypes())
        if err != nil {
            return err
        }
    }
    if m.GetMultipageLayouts() != nil {
        err := writer.WriteCollectionOfStringValues("multipageLayouts", SerializePrintMultipageLayout(m.GetMultipageLayouts()))
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
    if m.GetOrientations() != nil {
        err := writer.WriteCollectionOfStringValues("orientations", SerializePrintOrientation(m.GetOrientations()))
        if err != nil {
            return err
        }
    }
    if m.GetOutputBins() != nil {
        err := writer.WriteCollectionOfStringValues("outputBins", m.GetOutputBins())
        if err != nil {
            return err
        }
    }
    if m.GetPagesPerSheet() != nil {
        err := writer.WriteCollectionOfInt32Values("pagesPerSheet", m.GetPagesPerSheet())
        if err != nil {
            return err
        }
    }
    if m.GetQualities() != nil {
        err := writer.WriteCollectionOfStringValues("qualities", SerializePrintQuality(m.GetQualities()))
        if err != nil {
            return err
        }
    }
    if m.GetRightMargins() != nil {
        err := writer.WriteCollectionOfInt32Values("rightMargins", m.GetRightMargins())
        if err != nil {
            return err
        }
    }
    if m.GetScalings() != nil {
        err := writer.WriteCollectionOfStringValues("scalings", SerializePrintScaling(m.GetScalings()))
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteBoolValue("supportsFitPdfToPage", m.GetSupportsFitPdfToPage())
        if err != nil {
            return err
        }
    }
    if m.GetTopMargins() != nil {
        err := writer.WriteCollectionOfInt32Values("topMargins", m.GetTopMargins())
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
func (m *PrinterCapabilities) SetAdditionalData(value map[string]any)() {
    err := m.GetBackingStore().Set("additionalData", value)
    if err != nil {
        panic(err)
    }
}
// SetBackingStore sets the BackingStore property value. Stores model information.
func (m *PrinterCapabilities) SetBackingStore(value ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackingStore)() {
    m.backingStore = value
}
// SetBottomMargins sets the bottomMargins property value. A list of supported bottom margins(in microns) for the printer.
func (m *PrinterCapabilities) SetBottomMargins(value []int32)() {
    err := m.GetBackingStore().Set("bottomMargins", value)
    if err != nil {
        panic(err)
    }
}
// SetCollation sets the collation property value. True if the printer supports collating when printing muliple copies of a multi-page document; false otherwise.
func (m *PrinterCapabilities) SetCollation(value *bool)() {
    err := m.GetBackingStore().Set("collation", value)
    if err != nil {
        panic(err)
    }
}
// SetColorModes sets the colorModes property value. The color modes supported by the printer. Valid values are described in the following table.
func (m *PrinterCapabilities) SetColorModes(value []PrintColorMode)() {
    err := m.GetBackingStore().Set("colorModes", value)
    if err != nil {
        panic(err)
    }
}
// SetContentTypes sets the contentTypes property value. A list of supported content (MIME) types that the printer supports. It is not guaranteed that the Universal Print service supports printing all of these MIME types.
func (m *PrinterCapabilities) SetContentTypes(value []string)() {
    err := m.GetBackingStore().Set("contentTypes", value)
    if err != nil {
        panic(err)
    }
}
// SetCopiesPerJob sets the copiesPerJob property value. The range of copies per job supported by the printer.
func (m *PrinterCapabilities) SetCopiesPerJob(value IntegerRangeable)() {
    err := m.GetBackingStore().Set("copiesPerJob", value)
    if err != nil {
        panic(err)
    }
}
// SetDpis sets the dpis property value. The list of print resolutions in DPI that are supported by the printer.
func (m *PrinterCapabilities) SetDpis(value []int32)() {
    err := m.GetBackingStore().Set("dpis", value)
    if err != nil {
        panic(err)
    }
}
// SetDuplexModes sets the duplexModes property value. The list of duplex modes that are supported by the printer. Valid values are described in the following table.
func (m *PrinterCapabilities) SetDuplexModes(value []PrintDuplexMode)() {
    err := m.GetBackingStore().Set("duplexModes", value)
    if err != nil {
        panic(err)
    }
}
// SetFeedOrientations sets the feedOrientations property value. The list of feed orientations that are supported by the printer.
func (m *PrinterCapabilities) SetFeedOrientations(value []PrinterFeedOrientation)() {
    err := m.GetBackingStore().Set("feedOrientations", value)
    if err != nil {
        panic(err)
    }
}
// SetFinishings sets the finishings property value. Finishing processes the printer supports for a printed document.
func (m *PrinterCapabilities) SetFinishings(value []PrintFinishing)() {
    err := m.GetBackingStore().Set("finishings", value)
    if err != nil {
        panic(err)
    }
}
// SetInputBins sets the inputBins property value. Supported input bins for the printer.
func (m *PrinterCapabilities) SetInputBins(value []string)() {
    err := m.GetBackingStore().Set("inputBins", value)
    if err != nil {
        panic(err)
    }
}
// SetIsColorPrintingSupported sets the isColorPrintingSupported property value. True if color printing is supported by the printer; false otherwise. Read-only.
func (m *PrinterCapabilities) SetIsColorPrintingSupported(value *bool)() {
    err := m.GetBackingStore().Set("isColorPrintingSupported", value)
    if err != nil {
        panic(err)
    }
}
// SetIsPageRangeSupported sets the isPageRangeSupported property value. True if the printer supports printing by page ranges; false otherwise.
func (m *PrinterCapabilities) SetIsPageRangeSupported(value *bool)() {
    err := m.GetBackingStore().Set("isPageRangeSupported", value)
    if err != nil {
        panic(err)
    }
}
// SetLeftMargins sets the leftMargins property value. A list of supported left margins(in microns) for the printer.
func (m *PrinterCapabilities) SetLeftMargins(value []int32)() {
    err := m.GetBackingStore().Set("leftMargins", value)
    if err != nil {
        panic(err)
    }
}
// SetMediaColors sets the mediaColors property value. The media (i.e., paper) colors supported by the printer.
func (m *PrinterCapabilities) SetMediaColors(value []string)() {
    err := m.GetBackingStore().Set("mediaColors", value)
    if err != nil {
        panic(err)
    }
}
// SetMediaSizes sets the mediaSizes property value. The media sizes supported by the printer. Supports standard size names for ISO and ANSI media sizes. Valid values are in the following table.
func (m *PrinterCapabilities) SetMediaSizes(value []string)() {
    err := m.GetBackingStore().Set("mediaSizes", value)
    if err != nil {
        panic(err)
    }
}
// SetMediaTypes sets the mediaTypes property value. The media types supported by the printer.
func (m *PrinterCapabilities) SetMediaTypes(value []string)() {
    err := m.GetBackingStore().Set("mediaTypes", value)
    if err != nil {
        panic(err)
    }
}
// SetMultipageLayouts sets the multipageLayouts property value. The presentation directions supported by the printer. Supported values are described in the following table.
func (m *PrinterCapabilities) SetMultipageLayouts(value []PrintMultipageLayout)() {
    err := m.GetBackingStore().Set("multipageLayouts", value)
    if err != nil {
        panic(err)
    }
}
// SetOdataType sets the @odata.type property value. The OdataType property
func (m *PrinterCapabilities) SetOdataType(value *string)() {
    err := m.GetBackingStore().Set("odataType", value)
    if err != nil {
        panic(err)
    }
}
// SetOrientations sets the orientations property value. The print orientations supported by the printer. Valid values are described in the following table.
func (m *PrinterCapabilities) SetOrientations(value []PrintOrientation)() {
    err := m.GetBackingStore().Set("orientations", value)
    if err != nil {
        panic(err)
    }
}
// SetOutputBins sets the outputBins property value. The printer's supported output bins (trays).
func (m *PrinterCapabilities) SetOutputBins(value []string)() {
    err := m.GetBackingStore().Set("outputBins", value)
    if err != nil {
        panic(err)
    }
}
// SetPagesPerSheet sets the pagesPerSheet property value. Supported number of Input Pages to impose upon a single Impression.
func (m *PrinterCapabilities) SetPagesPerSheet(value []int32)() {
    err := m.GetBackingStore().Set("pagesPerSheet", value)
    if err != nil {
        panic(err)
    }
}
// SetQualities sets the qualities property value. The print qualities supported by the printer.
func (m *PrinterCapabilities) SetQualities(value []PrintQuality)() {
    err := m.GetBackingStore().Set("qualities", value)
    if err != nil {
        panic(err)
    }
}
// SetRightMargins sets the rightMargins property value. A list of supported right margins(in microns) for the printer.
func (m *PrinterCapabilities) SetRightMargins(value []int32)() {
    err := m.GetBackingStore().Set("rightMargins", value)
    if err != nil {
        panic(err)
    }
}
// SetScalings sets the scalings property value. Supported print scalings.
func (m *PrinterCapabilities) SetScalings(value []PrintScaling)() {
    err := m.GetBackingStore().Set("scalings", value)
    if err != nil {
        panic(err)
    }
}
// SetSupportsFitPdfToPage sets the supportsFitPdfToPage property value. True if the printer supports scaling PDF pages to match the print media size; false otherwise.
func (m *PrinterCapabilities) SetSupportsFitPdfToPage(value *bool)() {
    err := m.GetBackingStore().Set("supportsFitPdfToPage", value)
    if err != nil {
        panic(err)
    }
}
// SetTopMargins sets the topMargins property value. A list of supported top margins(in microns) for the printer.
func (m *PrinterCapabilities) SetTopMargins(value []int32)() {
    err := m.GetBackingStore().Set("topMargins", value)
    if err != nil {
        panic(err)
    }
}
type PrinterCapabilitiesable interface {
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.AdditionalDataHolder
    ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackedModel
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable
    GetBackingStore()(ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackingStore)
    GetBottomMargins()([]int32)
    GetCollation()(*bool)
    GetColorModes()([]PrintColorMode)
    GetContentTypes()([]string)
    GetCopiesPerJob()(IntegerRangeable)
    GetDpis()([]int32)
    GetDuplexModes()([]PrintDuplexMode)
    GetFeedOrientations()([]PrinterFeedOrientation)
    GetFinishings()([]PrintFinishing)
    GetInputBins()([]string)
    GetIsColorPrintingSupported()(*bool)
    GetIsPageRangeSupported()(*bool)
    GetLeftMargins()([]int32)
    GetMediaColors()([]string)
    GetMediaSizes()([]string)
    GetMediaTypes()([]string)
    GetMultipageLayouts()([]PrintMultipageLayout)
    GetOdataType()(*string)
    GetOrientations()([]PrintOrientation)
    GetOutputBins()([]string)
    GetPagesPerSheet()([]int32)
    GetQualities()([]PrintQuality)
    GetRightMargins()([]int32)
    GetScalings()([]PrintScaling)
    GetSupportsFitPdfToPage()(*bool)
    GetTopMargins()([]int32)
    SetBackingStore(value ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackingStore)()
    SetBottomMargins(value []int32)()
    SetCollation(value *bool)()
    SetColorModes(value []PrintColorMode)()
    SetContentTypes(value []string)()
    SetCopiesPerJob(value IntegerRangeable)()
    SetDpis(value []int32)()
    SetDuplexModes(value []PrintDuplexMode)()
    SetFeedOrientations(value []PrinterFeedOrientation)()
    SetFinishings(value []PrintFinishing)()
    SetInputBins(value []string)()
    SetIsColorPrintingSupported(value *bool)()
    SetIsPageRangeSupported(value *bool)()
    SetLeftMargins(value []int32)()
    SetMediaColors(value []string)()
    SetMediaSizes(value []string)()
    SetMediaTypes(value []string)()
    SetMultipageLayouts(value []PrintMultipageLayout)()
    SetOdataType(value *string)()
    SetOrientations(value []PrintOrientation)()
    SetOutputBins(value []string)()
    SetPagesPerSheet(value []int32)()
    SetQualities(value []PrintQuality)()
    SetRightMargins(value []int32)()
    SetScalings(value []PrintScaling)()
    SetSupportsFitPdfToPage(value *bool)()
    SetTopMargins(value []int32)()
}
