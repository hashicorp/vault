package models

import (
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

type ColumnDefinition struct {
    Entity
}
// NewColumnDefinition instantiates a new ColumnDefinition and sets the default values.
func NewColumnDefinition()(*ColumnDefinition) {
    m := &ColumnDefinition{
        Entity: *NewEntity(),
    }
    return m
}
// CreateColumnDefinitionFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
// returns a Parsable when successful
func CreateColumnDefinitionFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewColumnDefinition(), nil
}
// GetBoolean gets the boolean property value. This column stores Boolean values.
// returns a BooleanColumnable when successful
func (m *ColumnDefinition) GetBoolean()(BooleanColumnable) {
    val, err := m.GetBackingStore().Get("boolean")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(BooleanColumnable)
    }
    return nil
}
// GetCalculated gets the calculated property value. This column's data is calculated based on other columns.
// returns a CalculatedColumnable when successful
func (m *ColumnDefinition) GetCalculated()(CalculatedColumnable) {
    val, err := m.GetBackingStore().Get("calculated")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(CalculatedColumnable)
    }
    return nil
}
// GetChoice gets the choice property value. This column stores data from a list of choices.
// returns a ChoiceColumnable when successful
func (m *ColumnDefinition) GetChoice()(ChoiceColumnable) {
    val, err := m.GetBackingStore().Get("choice")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(ChoiceColumnable)
    }
    return nil
}
// GetColumnGroup gets the columnGroup property value. For site columns, the name of the group this column belongs to. Helps organize related columns.
// returns a *string when successful
func (m *ColumnDefinition) GetColumnGroup()(*string) {
    val, err := m.GetBackingStore().Get("columnGroup")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetContentApprovalStatus gets the contentApprovalStatus property value. This column stores content approval status.
// returns a ContentApprovalStatusColumnable when successful
func (m *ColumnDefinition) GetContentApprovalStatus()(ContentApprovalStatusColumnable) {
    val, err := m.GetBackingStore().Get("contentApprovalStatus")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(ContentApprovalStatusColumnable)
    }
    return nil
}
// GetCurrency gets the currency property value. This column stores currency values.
// returns a CurrencyColumnable when successful
func (m *ColumnDefinition) GetCurrency()(CurrencyColumnable) {
    val, err := m.GetBackingStore().Get("currency")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(CurrencyColumnable)
    }
    return nil
}
// GetDateTime gets the dateTime property value. This column stores DateTime values.
// returns a DateTimeColumnable when successful
func (m *ColumnDefinition) GetDateTime()(DateTimeColumnable) {
    val, err := m.GetBackingStore().Get("dateTime")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(DateTimeColumnable)
    }
    return nil
}
// GetDefaultValue gets the defaultValue property value. The default value for this column.
// returns a DefaultColumnValueable when successful
func (m *ColumnDefinition) GetDefaultValue()(DefaultColumnValueable) {
    val, err := m.GetBackingStore().Get("defaultValue")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(DefaultColumnValueable)
    }
    return nil
}
// GetDescription gets the description property value. The user-facing description of the column.
// returns a *string when successful
func (m *ColumnDefinition) GetDescription()(*string) {
    val, err := m.GetBackingStore().Get("description")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetDisplayName gets the displayName property value. The user-facing name of the column.
// returns a *string when successful
func (m *ColumnDefinition) GetDisplayName()(*string) {
    val, err := m.GetBackingStore().Get("displayName")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetEnforceUniqueValues gets the enforceUniqueValues property value. If true, no two list items may have the same value for this column.
// returns a *bool when successful
func (m *ColumnDefinition) GetEnforceUniqueValues()(*bool) {
    val, err := m.GetBackingStore().Get("enforceUniqueValues")
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
func (m *ColumnDefinition) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := m.Entity.GetFieldDeserializers()
    res["boolean"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetObjectValue(CreateBooleanColumnFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetBoolean(val.(BooleanColumnable))
        }
        return nil
    }
    res["calculated"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetObjectValue(CreateCalculatedColumnFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetCalculated(val.(CalculatedColumnable))
        }
        return nil
    }
    res["choice"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetObjectValue(CreateChoiceColumnFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetChoice(val.(ChoiceColumnable))
        }
        return nil
    }
    res["columnGroup"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetColumnGroup(val)
        }
        return nil
    }
    res["contentApprovalStatus"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetObjectValue(CreateContentApprovalStatusColumnFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetContentApprovalStatus(val.(ContentApprovalStatusColumnable))
        }
        return nil
    }
    res["currency"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetObjectValue(CreateCurrencyColumnFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetCurrency(val.(CurrencyColumnable))
        }
        return nil
    }
    res["dateTime"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetObjectValue(CreateDateTimeColumnFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetDateTime(val.(DateTimeColumnable))
        }
        return nil
    }
    res["defaultValue"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetObjectValue(CreateDefaultColumnValueFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetDefaultValue(val.(DefaultColumnValueable))
        }
        return nil
    }
    res["description"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetDescription(val)
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
    res["enforceUniqueValues"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetBoolValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetEnforceUniqueValues(val)
        }
        return nil
    }
    res["geolocation"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetObjectValue(CreateGeolocationColumnFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetGeolocation(val.(GeolocationColumnable))
        }
        return nil
    }
    res["hidden"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetBoolValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetHidden(val)
        }
        return nil
    }
    res["hyperlinkOrPicture"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetObjectValue(CreateHyperlinkOrPictureColumnFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetHyperlinkOrPicture(val.(HyperlinkOrPictureColumnable))
        }
        return nil
    }
    res["indexed"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetBoolValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetIndexed(val)
        }
        return nil
    }
    res["isDeletable"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetBoolValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetIsDeletable(val)
        }
        return nil
    }
    res["isReorderable"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetBoolValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetIsReorderable(val)
        }
        return nil
    }
    res["isSealed"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetBoolValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetIsSealed(val)
        }
        return nil
    }
    res["lookup"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetObjectValue(CreateLookupColumnFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetLookup(val.(LookupColumnable))
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
    res["number"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetObjectValue(CreateNumberColumnFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetNumber(val.(NumberColumnable))
        }
        return nil
    }
    res["personOrGroup"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetObjectValue(CreatePersonOrGroupColumnFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetPersonOrGroup(val.(PersonOrGroupColumnable))
        }
        return nil
    }
    res["propagateChanges"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetBoolValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetPropagateChanges(val)
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
    res["required"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetBoolValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetRequired(val)
        }
        return nil
    }
    res["sourceColumn"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetObjectValue(CreateColumnDefinitionFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetSourceColumn(val.(ColumnDefinitionable))
        }
        return nil
    }
    res["sourceContentType"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetObjectValue(CreateContentTypeInfoFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetSourceContentType(val.(ContentTypeInfoable))
        }
        return nil
    }
    res["term"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetObjectValue(CreateTermColumnFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetTerm(val.(TermColumnable))
        }
        return nil
    }
    res["text"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetObjectValue(CreateTextColumnFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetText(val.(TextColumnable))
        }
        return nil
    }
    res["thumbnail"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetObjectValue(CreateThumbnailColumnFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetThumbnail(val.(ThumbnailColumnable))
        }
        return nil
    }
    res["type"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetEnumValue(ParseColumnTypes)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetTypeEscaped(val.(*ColumnTypes))
        }
        return nil
    }
    res["validation"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetObjectValue(CreateColumnValidationFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetValidation(val.(ColumnValidationable))
        }
        return nil
    }
    return res
}
// GetGeolocation gets the geolocation property value. This column stores a geolocation.
// returns a GeolocationColumnable when successful
func (m *ColumnDefinition) GetGeolocation()(GeolocationColumnable) {
    val, err := m.GetBackingStore().Get("geolocation")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(GeolocationColumnable)
    }
    return nil
}
// GetHidden gets the hidden property value. Specifies whether the column is displayed in the user interface.
// returns a *bool when successful
func (m *ColumnDefinition) GetHidden()(*bool) {
    val, err := m.GetBackingStore().Get("hidden")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*bool)
    }
    return nil
}
// GetHyperlinkOrPicture gets the hyperlinkOrPicture property value. This column stores hyperlink or picture values.
// returns a HyperlinkOrPictureColumnable when successful
func (m *ColumnDefinition) GetHyperlinkOrPicture()(HyperlinkOrPictureColumnable) {
    val, err := m.GetBackingStore().Get("hyperlinkOrPicture")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(HyperlinkOrPictureColumnable)
    }
    return nil
}
// GetIndexed gets the indexed property value. Specifies whether the column values can be used for sorting and searching.
// returns a *bool when successful
func (m *ColumnDefinition) GetIndexed()(*bool) {
    val, err := m.GetBackingStore().Get("indexed")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*bool)
    }
    return nil
}
// GetIsDeletable gets the isDeletable property value. Indicates whether this column can be deleted.
// returns a *bool when successful
func (m *ColumnDefinition) GetIsDeletable()(*bool) {
    val, err := m.GetBackingStore().Get("isDeletable")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*bool)
    }
    return nil
}
// GetIsReorderable gets the isReorderable property value. Indicates whether values in the column can be reordered. Read-only.
// returns a *bool when successful
func (m *ColumnDefinition) GetIsReorderable()(*bool) {
    val, err := m.GetBackingStore().Get("isReorderable")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*bool)
    }
    return nil
}
// GetIsSealed gets the isSealed property value. Specifies whether the column can be changed.
// returns a *bool when successful
func (m *ColumnDefinition) GetIsSealed()(*bool) {
    val, err := m.GetBackingStore().Get("isSealed")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*bool)
    }
    return nil
}
// GetLookup gets the lookup property value. This column's data is looked up from another source in the site.
// returns a LookupColumnable when successful
func (m *ColumnDefinition) GetLookup()(LookupColumnable) {
    val, err := m.GetBackingStore().Get("lookup")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(LookupColumnable)
    }
    return nil
}
// GetName gets the name property value. The API-facing name of the column as it appears in the fields on a listItem. For the user-facing name, see displayName.
// returns a *string when successful
func (m *ColumnDefinition) GetName()(*string) {
    val, err := m.GetBackingStore().Get("name")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetNumber gets the number property value. This column stores number values.
// returns a NumberColumnable when successful
func (m *ColumnDefinition) GetNumber()(NumberColumnable) {
    val, err := m.GetBackingStore().Get("number")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(NumberColumnable)
    }
    return nil
}
// GetPersonOrGroup gets the personOrGroup property value. This column stores Person or Group values.
// returns a PersonOrGroupColumnable when successful
func (m *ColumnDefinition) GetPersonOrGroup()(PersonOrGroupColumnable) {
    val, err := m.GetBackingStore().Get("personOrGroup")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(PersonOrGroupColumnable)
    }
    return nil
}
// GetPropagateChanges gets the propagateChanges property value. If 'true', changes to this column will be propagated to lists that implement the column.
// returns a *bool when successful
func (m *ColumnDefinition) GetPropagateChanges()(*bool) {
    val, err := m.GetBackingStore().Get("propagateChanges")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*bool)
    }
    return nil
}
// GetReadOnly gets the readOnly property value. Specifies whether the column values can be modified.
// returns a *bool when successful
func (m *ColumnDefinition) GetReadOnly()(*bool) {
    val, err := m.GetBackingStore().Get("readOnly")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*bool)
    }
    return nil
}
// GetRequired gets the required property value. Specifies whether the column value isn't optional.
// returns a *bool when successful
func (m *ColumnDefinition) GetRequired()(*bool) {
    val, err := m.GetBackingStore().Get("required")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*bool)
    }
    return nil
}
// GetSourceColumn gets the sourceColumn property value. The source column for the content type column.
// returns a ColumnDefinitionable when successful
func (m *ColumnDefinition) GetSourceColumn()(ColumnDefinitionable) {
    val, err := m.GetBackingStore().Get("sourceColumn")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(ColumnDefinitionable)
    }
    return nil
}
// GetSourceContentType gets the sourceContentType property value. ContentType from which this column is inherited from. Present only in contentTypes columns response. Read-only.
// returns a ContentTypeInfoable when successful
func (m *ColumnDefinition) GetSourceContentType()(ContentTypeInfoable) {
    val, err := m.GetBackingStore().Get("sourceContentType")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(ContentTypeInfoable)
    }
    return nil
}
// GetTerm gets the term property value. This column stores taxonomy terms.
// returns a TermColumnable when successful
func (m *ColumnDefinition) GetTerm()(TermColumnable) {
    val, err := m.GetBackingStore().Get("term")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(TermColumnable)
    }
    return nil
}
// GetText gets the text property value. This column stores text values.
// returns a TextColumnable when successful
func (m *ColumnDefinition) GetText()(TextColumnable) {
    val, err := m.GetBackingStore().Get("text")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(TextColumnable)
    }
    return nil
}
// GetThumbnail gets the thumbnail property value. This column stores thumbnail values.
// returns a ThumbnailColumnable when successful
func (m *ColumnDefinition) GetThumbnail()(ThumbnailColumnable) {
    val, err := m.GetBackingStore().Get("thumbnail")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(ThumbnailColumnable)
    }
    return nil
}
// GetTypeEscaped gets the type property value. For site columns, the type of column. Read-only.
// returns a *ColumnTypes when successful
func (m *ColumnDefinition) GetTypeEscaped()(*ColumnTypes) {
    val, err := m.GetBackingStore().Get("typeEscaped")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*ColumnTypes)
    }
    return nil
}
// GetValidation gets the validation property value. This column stores validation formula and message for the column.
// returns a ColumnValidationable when successful
func (m *ColumnDefinition) GetValidation()(ColumnValidationable) {
    val, err := m.GetBackingStore().Get("validation")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(ColumnValidationable)
    }
    return nil
}
// Serialize serializes information the current object
func (m *ColumnDefinition) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    err := m.Entity.Serialize(writer)
    if err != nil {
        return err
    }
    {
        err = writer.WriteObjectValue("boolean", m.GetBoolean())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteObjectValue("calculated", m.GetCalculated())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteObjectValue("choice", m.GetChoice())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteStringValue("columnGroup", m.GetColumnGroup())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteObjectValue("contentApprovalStatus", m.GetContentApprovalStatus())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteObjectValue("currency", m.GetCurrency())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteObjectValue("dateTime", m.GetDateTime())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteObjectValue("defaultValue", m.GetDefaultValue())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteStringValue("description", m.GetDescription())
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
        err = writer.WriteBoolValue("enforceUniqueValues", m.GetEnforceUniqueValues())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteObjectValue("geolocation", m.GetGeolocation())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteBoolValue("hidden", m.GetHidden())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteObjectValue("hyperlinkOrPicture", m.GetHyperlinkOrPicture())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteBoolValue("indexed", m.GetIndexed())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteBoolValue("isDeletable", m.GetIsDeletable())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteBoolValue("isReorderable", m.GetIsReorderable())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteBoolValue("isSealed", m.GetIsSealed())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteObjectValue("lookup", m.GetLookup())
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
    {
        err = writer.WriteObjectValue("number", m.GetNumber())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteObjectValue("personOrGroup", m.GetPersonOrGroup())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteBoolValue("propagateChanges", m.GetPropagateChanges())
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
        err = writer.WriteBoolValue("required", m.GetRequired())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteObjectValue("sourceColumn", m.GetSourceColumn())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteObjectValue("sourceContentType", m.GetSourceContentType())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteObjectValue("term", m.GetTerm())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteObjectValue("text", m.GetText())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteObjectValue("thumbnail", m.GetThumbnail())
        if err != nil {
            return err
        }
    }
    if m.GetTypeEscaped() != nil {
        cast := (*m.GetTypeEscaped()).String()
        err = writer.WriteStringValue("type", &cast)
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteObjectValue("validation", m.GetValidation())
        if err != nil {
            return err
        }
    }
    return nil
}
// SetBoolean sets the boolean property value. This column stores Boolean values.
func (m *ColumnDefinition) SetBoolean(value BooleanColumnable)() {
    err := m.GetBackingStore().Set("boolean", value)
    if err != nil {
        panic(err)
    }
}
// SetCalculated sets the calculated property value. This column's data is calculated based on other columns.
func (m *ColumnDefinition) SetCalculated(value CalculatedColumnable)() {
    err := m.GetBackingStore().Set("calculated", value)
    if err != nil {
        panic(err)
    }
}
// SetChoice sets the choice property value. This column stores data from a list of choices.
func (m *ColumnDefinition) SetChoice(value ChoiceColumnable)() {
    err := m.GetBackingStore().Set("choice", value)
    if err != nil {
        panic(err)
    }
}
// SetColumnGroup sets the columnGroup property value. For site columns, the name of the group this column belongs to. Helps organize related columns.
func (m *ColumnDefinition) SetColumnGroup(value *string)() {
    err := m.GetBackingStore().Set("columnGroup", value)
    if err != nil {
        panic(err)
    }
}
// SetContentApprovalStatus sets the contentApprovalStatus property value. This column stores content approval status.
func (m *ColumnDefinition) SetContentApprovalStatus(value ContentApprovalStatusColumnable)() {
    err := m.GetBackingStore().Set("contentApprovalStatus", value)
    if err != nil {
        panic(err)
    }
}
// SetCurrency sets the currency property value. This column stores currency values.
func (m *ColumnDefinition) SetCurrency(value CurrencyColumnable)() {
    err := m.GetBackingStore().Set("currency", value)
    if err != nil {
        panic(err)
    }
}
// SetDateTime sets the dateTime property value. This column stores DateTime values.
func (m *ColumnDefinition) SetDateTime(value DateTimeColumnable)() {
    err := m.GetBackingStore().Set("dateTime", value)
    if err != nil {
        panic(err)
    }
}
// SetDefaultValue sets the defaultValue property value. The default value for this column.
func (m *ColumnDefinition) SetDefaultValue(value DefaultColumnValueable)() {
    err := m.GetBackingStore().Set("defaultValue", value)
    if err != nil {
        panic(err)
    }
}
// SetDescription sets the description property value. The user-facing description of the column.
func (m *ColumnDefinition) SetDescription(value *string)() {
    err := m.GetBackingStore().Set("description", value)
    if err != nil {
        panic(err)
    }
}
// SetDisplayName sets the displayName property value. The user-facing name of the column.
func (m *ColumnDefinition) SetDisplayName(value *string)() {
    err := m.GetBackingStore().Set("displayName", value)
    if err != nil {
        panic(err)
    }
}
// SetEnforceUniqueValues sets the enforceUniqueValues property value. If true, no two list items may have the same value for this column.
func (m *ColumnDefinition) SetEnforceUniqueValues(value *bool)() {
    err := m.GetBackingStore().Set("enforceUniqueValues", value)
    if err != nil {
        panic(err)
    }
}
// SetGeolocation sets the geolocation property value. This column stores a geolocation.
func (m *ColumnDefinition) SetGeolocation(value GeolocationColumnable)() {
    err := m.GetBackingStore().Set("geolocation", value)
    if err != nil {
        panic(err)
    }
}
// SetHidden sets the hidden property value. Specifies whether the column is displayed in the user interface.
func (m *ColumnDefinition) SetHidden(value *bool)() {
    err := m.GetBackingStore().Set("hidden", value)
    if err != nil {
        panic(err)
    }
}
// SetHyperlinkOrPicture sets the hyperlinkOrPicture property value. This column stores hyperlink or picture values.
func (m *ColumnDefinition) SetHyperlinkOrPicture(value HyperlinkOrPictureColumnable)() {
    err := m.GetBackingStore().Set("hyperlinkOrPicture", value)
    if err != nil {
        panic(err)
    }
}
// SetIndexed sets the indexed property value. Specifies whether the column values can be used for sorting and searching.
func (m *ColumnDefinition) SetIndexed(value *bool)() {
    err := m.GetBackingStore().Set("indexed", value)
    if err != nil {
        panic(err)
    }
}
// SetIsDeletable sets the isDeletable property value. Indicates whether this column can be deleted.
func (m *ColumnDefinition) SetIsDeletable(value *bool)() {
    err := m.GetBackingStore().Set("isDeletable", value)
    if err != nil {
        panic(err)
    }
}
// SetIsReorderable sets the isReorderable property value. Indicates whether values in the column can be reordered. Read-only.
func (m *ColumnDefinition) SetIsReorderable(value *bool)() {
    err := m.GetBackingStore().Set("isReorderable", value)
    if err != nil {
        panic(err)
    }
}
// SetIsSealed sets the isSealed property value. Specifies whether the column can be changed.
func (m *ColumnDefinition) SetIsSealed(value *bool)() {
    err := m.GetBackingStore().Set("isSealed", value)
    if err != nil {
        panic(err)
    }
}
// SetLookup sets the lookup property value. This column's data is looked up from another source in the site.
func (m *ColumnDefinition) SetLookup(value LookupColumnable)() {
    err := m.GetBackingStore().Set("lookup", value)
    if err != nil {
        panic(err)
    }
}
// SetName sets the name property value. The API-facing name of the column as it appears in the fields on a listItem. For the user-facing name, see displayName.
func (m *ColumnDefinition) SetName(value *string)() {
    err := m.GetBackingStore().Set("name", value)
    if err != nil {
        panic(err)
    }
}
// SetNumber sets the number property value. This column stores number values.
func (m *ColumnDefinition) SetNumber(value NumberColumnable)() {
    err := m.GetBackingStore().Set("number", value)
    if err != nil {
        panic(err)
    }
}
// SetPersonOrGroup sets the personOrGroup property value. This column stores Person or Group values.
func (m *ColumnDefinition) SetPersonOrGroup(value PersonOrGroupColumnable)() {
    err := m.GetBackingStore().Set("personOrGroup", value)
    if err != nil {
        panic(err)
    }
}
// SetPropagateChanges sets the propagateChanges property value. If 'true', changes to this column will be propagated to lists that implement the column.
func (m *ColumnDefinition) SetPropagateChanges(value *bool)() {
    err := m.GetBackingStore().Set("propagateChanges", value)
    if err != nil {
        panic(err)
    }
}
// SetReadOnly sets the readOnly property value. Specifies whether the column values can be modified.
func (m *ColumnDefinition) SetReadOnly(value *bool)() {
    err := m.GetBackingStore().Set("readOnly", value)
    if err != nil {
        panic(err)
    }
}
// SetRequired sets the required property value. Specifies whether the column value isn't optional.
func (m *ColumnDefinition) SetRequired(value *bool)() {
    err := m.GetBackingStore().Set("required", value)
    if err != nil {
        panic(err)
    }
}
// SetSourceColumn sets the sourceColumn property value. The source column for the content type column.
func (m *ColumnDefinition) SetSourceColumn(value ColumnDefinitionable)() {
    err := m.GetBackingStore().Set("sourceColumn", value)
    if err != nil {
        panic(err)
    }
}
// SetSourceContentType sets the sourceContentType property value. ContentType from which this column is inherited from. Present only in contentTypes columns response. Read-only.
func (m *ColumnDefinition) SetSourceContentType(value ContentTypeInfoable)() {
    err := m.GetBackingStore().Set("sourceContentType", value)
    if err != nil {
        panic(err)
    }
}
// SetTerm sets the term property value. This column stores taxonomy terms.
func (m *ColumnDefinition) SetTerm(value TermColumnable)() {
    err := m.GetBackingStore().Set("term", value)
    if err != nil {
        panic(err)
    }
}
// SetText sets the text property value. This column stores text values.
func (m *ColumnDefinition) SetText(value TextColumnable)() {
    err := m.GetBackingStore().Set("text", value)
    if err != nil {
        panic(err)
    }
}
// SetThumbnail sets the thumbnail property value. This column stores thumbnail values.
func (m *ColumnDefinition) SetThumbnail(value ThumbnailColumnable)() {
    err := m.GetBackingStore().Set("thumbnail", value)
    if err != nil {
        panic(err)
    }
}
// SetTypeEscaped sets the type property value. For site columns, the type of column. Read-only.
func (m *ColumnDefinition) SetTypeEscaped(value *ColumnTypes)() {
    err := m.GetBackingStore().Set("typeEscaped", value)
    if err != nil {
        panic(err)
    }
}
// SetValidation sets the validation property value. This column stores validation formula and message for the column.
func (m *ColumnDefinition) SetValidation(value ColumnValidationable)() {
    err := m.GetBackingStore().Set("validation", value)
    if err != nil {
        panic(err)
    }
}
type ColumnDefinitionable interface {
    Entityable
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable
    GetBoolean()(BooleanColumnable)
    GetCalculated()(CalculatedColumnable)
    GetChoice()(ChoiceColumnable)
    GetColumnGroup()(*string)
    GetContentApprovalStatus()(ContentApprovalStatusColumnable)
    GetCurrency()(CurrencyColumnable)
    GetDateTime()(DateTimeColumnable)
    GetDefaultValue()(DefaultColumnValueable)
    GetDescription()(*string)
    GetDisplayName()(*string)
    GetEnforceUniqueValues()(*bool)
    GetGeolocation()(GeolocationColumnable)
    GetHidden()(*bool)
    GetHyperlinkOrPicture()(HyperlinkOrPictureColumnable)
    GetIndexed()(*bool)
    GetIsDeletable()(*bool)
    GetIsReorderable()(*bool)
    GetIsSealed()(*bool)
    GetLookup()(LookupColumnable)
    GetName()(*string)
    GetNumber()(NumberColumnable)
    GetPersonOrGroup()(PersonOrGroupColumnable)
    GetPropagateChanges()(*bool)
    GetReadOnly()(*bool)
    GetRequired()(*bool)
    GetSourceColumn()(ColumnDefinitionable)
    GetSourceContentType()(ContentTypeInfoable)
    GetTerm()(TermColumnable)
    GetText()(TextColumnable)
    GetThumbnail()(ThumbnailColumnable)
    GetTypeEscaped()(*ColumnTypes)
    GetValidation()(ColumnValidationable)
    SetBoolean(value BooleanColumnable)()
    SetCalculated(value CalculatedColumnable)()
    SetChoice(value ChoiceColumnable)()
    SetColumnGroup(value *string)()
    SetContentApprovalStatus(value ContentApprovalStatusColumnable)()
    SetCurrency(value CurrencyColumnable)()
    SetDateTime(value DateTimeColumnable)()
    SetDefaultValue(value DefaultColumnValueable)()
    SetDescription(value *string)()
    SetDisplayName(value *string)()
    SetEnforceUniqueValues(value *bool)()
    SetGeolocation(value GeolocationColumnable)()
    SetHidden(value *bool)()
    SetHyperlinkOrPicture(value HyperlinkOrPictureColumnable)()
    SetIndexed(value *bool)()
    SetIsDeletable(value *bool)()
    SetIsReorderable(value *bool)()
    SetIsSealed(value *bool)()
    SetLookup(value LookupColumnable)()
    SetName(value *string)()
    SetNumber(value NumberColumnable)()
    SetPersonOrGroup(value PersonOrGroupColumnable)()
    SetPropagateChanges(value *bool)()
    SetReadOnly(value *bool)()
    SetRequired(value *bool)()
    SetSourceColumn(value ColumnDefinitionable)()
    SetSourceContentType(value ContentTypeInfoable)()
    SetTerm(value TermColumnable)()
    SetText(value TextColumnable)()
    SetThumbnail(value ThumbnailColumnable)()
    SetTypeEscaped(value *ColumnTypes)()
    SetValidation(value ColumnValidationable)()
}
