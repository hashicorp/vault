package models

import (
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
    ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e "github.com/microsoft/kiota-abstractions-go/store"
)

type ParseExpressionResponse struct {
    // Stores model information.
    backingStore ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackingStore
}
// NewParseExpressionResponse instantiates a new ParseExpressionResponse and sets the default values.
func NewParseExpressionResponse()(*ParseExpressionResponse) {
    m := &ParseExpressionResponse{
    }
    m.backingStore = ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackingStoreFactoryInstance();
    m.SetAdditionalData(make(map[string]any))
    return m
}
// CreateParseExpressionResponseFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
// returns a Parsable when successful
func CreateParseExpressionResponseFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewParseExpressionResponse(), nil
}
// GetAdditionalData gets the AdditionalData property value. Stores additional data not described in the OpenAPI description found when deserializing. Can be used for serialization as well.
// returns a map[string]any when successful
func (m *ParseExpressionResponse) GetAdditionalData()(map[string]any) {
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
func (m *ParseExpressionResponse) GetBackingStore()(ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackingStore) {
    return m.backingStore
}
// GetError gets the error property value. Error details, if expression evaluation resulted in an error.
// returns a PublicErrorable when successful
func (m *ParseExpressionResponse) GetError()(PublicErrorable) {
    val, err := m.GetBackingStore().Get("error")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(PublicErrorable)
    }
    return nil
}
// GetEvaluationResult gets the evaluationResult property value. A collection of values produced by the evaluation of the expression.
// returns a []string when successful
func (m *ParseExpressionResponse) GetEvaluationResult()([]string) {
    val, err := m.GetBackingStore().Get("evaluationResult")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]string)
    }
    return nil
}
// GetEvaluationSucceeded gets the evaluationSucceeded property value. true if the evaluation was successful.
// returns a *bool when successful
func (m *ParseExpressionResponse) GetEvaluationSucceeded()(*bool) {
    val, err := m.GetBackingStore().Get("evaluationSucceeded")
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
func (m *ParseExpressionResponse) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := make(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error))
    res["error"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetObjectValue(CreatePublicErrorFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetError(val.(PublicErrorable))
        }
        return nil
    }
    res["evaluationResult"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
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
            m.SetEvaluationResult(res)
        }
        return nil
    }
    res["evaluationSucceeded"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetBoolValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetEvaluationSucceeded(val)
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
    res["parsedExpression"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetObjectValue(CreateAttributeMappingSourceFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetParsedExpression(val.(AttributeMappingSourceable))
        }
        return nil
    }
    res["parsingSucceeded"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetBoolValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetParsingSucceeded(val)
        }
        return nil
    }
    return res
}
// GetOdataType gets the @odata.type property value. The OdataType property
// returns a *string when successful
func (m *ParseExpressionResponse) GetOdataType()(*string) {
    val, err := m.GetBackingStore().Get("odataType")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetParsedExpression gets the parsedExpression property value. An attributeMappingSource object representing the parsed expression.
// returns a AttributeMappingSourceable when successful
func (m *ParseExpressionResponse) GetParsedExpression()(AttributeMappingSourceable) {
    val, err := m.GetBackingStore().Get("parsedExpression")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(AttributeMappingSourceable)
    }
    return nil
}
// GetParsingSucceeded gets the parsingSucceeded property value. true if the expression was parsed successfully.
// returns a *bool when successful
func (m *ParseExpressionResponse) GetParsingSucceeded()(*bool) {
    val, err := m.GetBackingStore().Get("parsingSucceeded")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*bool)
    }
    return nil
}
// Serialize serializes information the current object
func (m *ParseExpressionResponse) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    {
        err := writer.WriteObjectValue("error", m.GetError())
        if err != nil {
            return err
        }
    }
    if m.GetEvaluationResult() != nil {
        err := writer.WriteCollectionOfStringValues("evaluationResult", m.GetEvaluationResult())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteBoolValue("evaluationSucceeded", m.GetEvaluationSucceeded())
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
        err := writer.WriteObjectValue("parsedExpression", m.GetParsedExpression())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteBoolValue("parsingSucceeded", m.GetParsingSucceeded())
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
func (m *ParseExpressionResponse) SetAdditionalData(value map[string]any)() {
    err := m.GetBackingStore().Set("additionalData", value)
    if err != nil {
        panic(err)
    }
}
// SetBackingStore sets the BackingStore property value. Stores model information.
func (m *ParseExpressionResponse) SetBackingStore(value ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackingStore)() {
    m.backingStore = value
}
// SetError sets the error property value. Error details, if expression evaluation resulted in an error.
func (m *ParseExpressionResponse) SetError(value PublicErrorable)() {
    err := m.GetBackingStore().Set("error", value)
    if err != nil {
        panic(err)
    }
}
// SetEvaluationResult sets the evaluationResult property value. A collection of values produced by the evaluation of the expression.
func (m *ParseExpressionResponse) SetEvaluationResult(value []string)() {
    err := m.GetBackingStore().Set("evaluationResult", value)
    if err != nil {
        panic(err)
    }
}
// SetEvaluationSucceeded sets the evaluationSucceeded property value. true if the evaluation was successful.
func (m *ParseExpressionResponse) SetEvaluationSucceeded(value *bool)() {
    err := m.GetBackingStore().Set("evaluationSucceeded", value)
    if err != nil {
        panic(err)
    }
}
// SetOdataType sets the @odata.type property value. The OdataType property
func (m *ParseExpressionResponse) SetOdataType(value *string)() {
    err := m.GetBackingStore().Set("odataType", value)
    if err != nil {
        panic(err)
    }
}
// SetParsedExpression sets the parsedExpression property value. An attributeMappingSource object representing the parsed expression.
func (m *ParseExpressionResponse) SetParsedExpression(value AttributeMappingSourceable)() {
    err := m.GetBackingStore().Set("parsedExpression", value)
    if err != nil {
        panic(err)
    }
}
// SetParsingSucceeded sets the parsingSucceeded property value. true if the expression was parsed successfully.
func (m *ParseExpressionResponse) SetParsingSucceeded(value *bool)() {
    err := m.GetBackingStore().Set("parsingSucceeded", value)
    if err != nil {
        panic(err)
    }
}
type ParseExpressionResponseable interface {
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.AdditionalDataHolder
    ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackedModel
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable
    GetBackingStore()(ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackingStore)
    GetError()(PublicErrorable)
    GetEvaluationResult()([]string)
    GetEvaluationSucceeded()(*bool)
    GetOdataType()(*string)
    GetParsedExpression()(AttributeMappingSourceable)
    GetParsingSucceeded()(*bool)
    SetBackingStore(value ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackingStore)()
    SetError(value PublicErrorable)()
    SetEvaluationResult(value []string)()
    SetEvaluationSucceeded(value *bool)()
    SetOdataType(value *string)()
    SetParsedExpression(value AttributeMappingSourceable)()
    SetParsingSucceeded(value *bool)()
}
