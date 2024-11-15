package models

import (
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
    ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e "github.com/microsoft/kiota-abstractions-go/store"
)

type TitleArea struct {
    // Stores model information.
    backingStore ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackingStore
}
// NewTitleArea instantiates a new TitleArea and sets the default values.
func NewTitleArea()(*TitleArea) {
    m := &TitleArea{
    }
    m.backingStore = ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackingStoreFactoryInstance();
    m.SetAdditionalData(make(map[string]any))
    return m
}
// CreateTitleAreaFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
// returns a Parsable when successful
func CreateTitleAreaFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewTitleArea(), nil
}
// GetAdditionalData gets the AdditionalData property value. Stores additional data not described in the OpenAPI description found when deserializing. Can be used for serialization as well.
// returns a map[string]any when successful
func (m *TitleArea) GetAdditionalData()(map[string]any) {
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
// GetAlternativeText gets the alternativeText property value. Alternative text on the title area.
// returns a *string when successful
func (m *TitleArea) GetAlternativeText()(*string) {
    val, err := m.GetBackingStore().Get("alternativeText")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetBackingStore gets the BackingStore property value. Stores model information.
// returns a BackingStore when successful
func (m *TitleArea) GetBackingStore()(ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackingStore) {
    return m.backingStore
}
// GetEnableGradientEffect gets the enableGradientEffect property value. Indicates whether the title area has a gradient effect enabled.
// returns a *bool when successful
func (m *TitleArea) GetEnableGradientEffect()(*bool) {
    val, err := m.GetBackingStore().Get("enableGradientEffect")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*bool)
    }
    return nil
}
// GetFieldDeserializers the deserialization information for the current model
// returns a map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error) when successful
func (m *TitleArea) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := make(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error))
    res["alternativeText"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetAlternativeText(val)
        }
        return nil
    }
    res["enableGradientEffect"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetBoolValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetEnableGradientEffect(val)
        }
        return nil
    }
    res["imageWebUrl"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetImageWebUrl(val)
        }
        return nil
    }
    res["layout"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetEnumValue(ParseTitleAreaLayoutType)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetLayout(val.(*TitleAreaLayoutType))
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
    res["serverProcessedContent"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetObjectValue(CreateServerProcessedContentFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetServerProcessedContent(val.(ServerProcessedContentable))
        }
        return nil
    }
    res["showAuthor"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetBoolValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetShowAuthor(val)
        }
        return nil
    }
    res["showPublishedDate"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetBoolValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetShowPublishedDate(val)
        }
        return nil
    }
    res["showTextBlockAboveTitle"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetBoolValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetShowTextBlockAboveTitle(val)
        }
        return nil
    }
    res["textAboveTitle"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetTextAboveTitle(val)
        }
        return nil
    }
    res["textAlignment"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetEnumValue(ParseTitleAreaTextAlignmentType)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetTextAlignment(val.(*TitleAreaTextAlignmentType))
        }
        return nil
    }
    return res
}
// GetImageWebUrl gets the imageWebUrl property value. URL of the image in the title area.
// returns a *string when successful
func (m *TitleArea) GetImageWebUrl()(*string) {
    val, err := m.GetBackingStore().Get("imageWebUrl")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetLayout gets the layout property value. Enumeration value that indicates the layout of the title area. The possible values are: imageAndTitle, plain, colorBlock, overlap, unknownFutureValue.
// returns a *TitleAreaLayoutType when successful
func (m *TitleArea) GetLayout()(*TitleAreaLayoutType) {
    val, err := m.GetBackingStore().Get("layout")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*TitleAreaLayoutType)
    }
    return nil
}
// GetOdataType gets the @odata.type property value. The OdataType property
// returns a *string when successful
func (m *TitleArea) GetOdataType()(*string) {
    val, err := m.GetBackingStore().Get("odataType")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetServerProcessedContent gets the serverProcessedContent property value. Contains collections of data that can be processed by server side services like search index and link fixup.
// returns a ServerProcessedContentable when successful
func (m *TitleArea) GetServerProcessedContent()(ServerProcessedContentable) {
    val, err := m.GetBackingStore().Get("serverProcessedContent")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(ServerProcessedContentable)
    }
    return nil
}
// GetShowAuthor gets the showAuthor property value. Indicates whether the author should be shown in title area.
// returns a *bool when successful
func (m *TitleArea) GetShowAuthor()(*bool) {
    val, err := m.GetBackingStore().Get("showAuthor")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*bool)
    }
    return nil
}
// GetShowPublishedDate gets the showPublishedDate property value. Indicates whether the published date should be shown in title area.
// returns a *bool when successful
func (m *TitleArea) GetShowPublishedDate()(*bool) {
    val, err := m.GetBackingStore().Get("showPublishedDate")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*bool)
    }
    return nil
}
// GetShowTextBlockAboveTitle gets the showTextBlockAboveTitle property value. Indicates whether the text block above title should be shown in title area.
// returns a *bool when successful
func (m *TitleArea) GetShowTextBlockAboveTitle()(*bool) {
    val, err := m.GetBackingStore().Get("showTextBlockAboveTitle")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*bool)
    }
    return nil
}
// GetTextAboveTitle gets the textAboveTitle property value. The text above title line.
// returns a *string when successful
func (m *TitleArea) GetTextAboveTitle()(*string) {
    val, err := m.GetBackingStore().Get("textAboveTitle")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetTextAlignment gets the textAlignment property value. Enumeration value that indicates the text alignment of the title area. The possible values are: left, center, unknownFutureValue.
// returns a *TitleAreaTextAlignmentType when successful
func (m *TitleArea) GetTextAlignment()(*TitleAreaTextAlignmentType) {
    val, err := m.GetBackingStore().Get("textAlignment")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*TitleAreaTextAlignmentType)
    }
    return nil
}
// Serialize serializes information the current object
func (m *TitleArea) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    {
        err := writer.WriteStringValue("alternativeText", m.GetAlternativeText())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteBoolValue("enableGradientEffect", m.GetEnableGradientEffect())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteStringValue("imageWebUrl", m.GetImageWebUrl())
        if err != nil {
            return err
        }
    }
    if m.GetLayout() != nil {
        cast := (*m.GetLayout()).String()
        err := writer.WriteStringValue("layout", &cast)
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
    {
        err := writer.WriteObjectValue("serverProcessedContent", m.GetServerProcessedContent())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteBoolValue("showAuthor", m.GetShowAuthor())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteBoolValue("showPublishedDate", m.GetShowPublishedDate())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteBoolValue("showTextBlockAboveTitle", m.GetShowTextBlockAboveTitle())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteStringValue("textAboveTitle", m.GetTextAboveTitle())
        if err != nil {
            return err
        }
    }
    if m.GetTextAlignment() != nil {
        cast := (*m.GetTextAlignment()).String()
        err := writer.WriteStringValue("textAlignment", &cast)
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
func (m *TitleArea) SetAdditionalData(value map[string]any)() {
    err := m.GetBackingStore().Set("additionalData", value)
    if err != nil {
        panic(err)
    }
}
// SetAlternativeText sets the alternativeText property value. Alternative text on the title area.
func (m *TitleArea) SetAlternativeText(value *string)() {
    err := m.GetBackingStore().Set("alternativeText", value)
    if err != nil {
        panic(err)
    }
}
// SetBackingStore sets the BackingStore property value. Stores model information.
func (m *TitleArea) SetBackingStore(value ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackingStore)() {
    m.backingStore = value
}
// SetEnableGradientEffect sets the enableGradientEffect property value. Indicates whether the title area has a gradient effect enabled.
func (m *TitleArea) SetEnableGradientEffect(value *bool)() {
    err := m.GetBackingStore().Set("enableGradientEffect", value)
    if err != nil {
        panic(err)
    }
}
// SetImageWebUrl sets the imageWebUrl property value. URL of the image in the title area.
func (m *TitleArea) SetImageWebUrl(value *string)() {
    err := m.GetBackingStore().Set("imageWebUrl", value)
    if err != nil {
        panic(err)
    }
}
// SetLayout sets the layout property value. Enumeration value that indicates the layout of the title area. The possible values are: imageAndTitle, plain, colorBlock, overlap, unknownFutureValue.
func (m *TitleArea) SetLayout(value *TitleAreaLayoutType)() {
    err := m.GetBackingStore().Set("layout", value)
    if err != nil {
        panic(err)
    }
}
// SetOdataType sets the @odata.type property value. The OdataType property
func (m *TitleArea) SetOdataType(value *string)() {
    err := m.GetBackingStore().Set("odataType", value)
    if err != nil {
        panic(err)
    }
}
// SetServerProcessedContent sets the serverProcessedContent property value. Contains collections of data that can be processed by server side services like search index and link fixup.
func (m *TitleArea) SetServerProcessedContent(value ServerProcessedContentable)() {
    err := m.GetBackingStore().Set("serverProcessedContent", value)
    if err != nil {
        panic(err)
    }
}
// SetShowAuthor sets the showAuthor property value. Indicates whether the author should be shown in title area.
func (m *TitleArea) SetShowAuthor(value *bool)() {
    err := m.GetBackingStore().Set("showAuthor", value)
    if err != nil {
        panic(err)
    }
}
// SetShowPublishedDate sets the showPublishedDate property value. Indicates whether the published date should be shown in title area.
func (m *TitleArea) SetShowPublishedDate(value *bool)() {
    err := m.GetBackingStore().Set("showPublishedDate", value)
    if err != nil {
        panic(err)
    }
}
// SetShowTextBlockAboveTitle sets the showTextBlockAboveTitle property value. Indicates whether the text block above title should be shown in title area.
func (m *TitleArea) SetShowTextBlockAboveTitle(value *bool)() {
    err := m.GetBackingStore().Set("showTextBlockAboveTitle", value)
    if err != nil {
        panic(err)
    }
}
// SetTextAboveTitle sets the textAboveTitle property value. The text above title line.
func (m *TitleArea) SetTextAboveTitle(value *string)() {
    err := m.GetBackingStore().Set("textAboveTitle", value)
    if err != nil {
        panic(err)
    }
}
// SetTextAlignment sets the textAlignment property value. Enumeration value that indicates the text alignment of the title area. The possible values are: left, center, unknownFutureValue.
func (m *TitleArea) SetTextAlignment(value *TitleAreaTextAlignmentType)() {
    err := m.GetBackingStore().Set("textAlignment", value)
    if err != nil {
        panic(err)
    }
}
type TitleAreaable interface {
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.AdditionalDataHolder
    ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackedModel
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable
    GetAlternativeText()(*string)
    GetBackingStore()(ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackingStore)
    GetEnableGradientEffect()(*bool)
    GetImageWebUrl()(*string)
    GetLayout()(*TitleAreaLayoutType)
    GetOdataType()(*string)
    GetServerProcessedContent()(ServerProcessedContentable)
    GetShowAuthor()(*bool)
    GetShowPublishedDate()(*bool)
    GetShowTextBlockAboveTitle()(*bool)
    GetTextAboveTitle()(*string)
    GetTextAlignment()(*TitleAreaTextAlignmentType)
    SetAlternativeText(value *string)()
    SetBackingStore(value ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackingStore)()
    SetEnableGradientEffect(value *bool)()
    SetImageWebUrl(value *string)()
    SetLayout(value *TitleAreaLayoutType)()
    SetOdataType(value *string)()
    SetServerProcessedContent(value ServerProcessedContentable)()
    SetShowAuthor(value *bool)()
    SetShowPublishedDate(value *bool)()
    SetShowTextBlockAboveTitle(value *bool)()
    SetTextAboveTitle(value *string)()
    SetTextAlignment(value *TitleAreaTextAlignmentType)()
}
