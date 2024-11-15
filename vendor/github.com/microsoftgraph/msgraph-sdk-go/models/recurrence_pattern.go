package models

import (
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
    ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e "github.com/microsoft/kiota-abstractions-go/store"
)

type RecurrencePattern struct {
    // Stores model information.
    backingStore ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackingStore
}
// NewRecurrencePattern instantiates a new RecurrencePattern and sets the default values.
func NewRecurrencePattern()(*RecurrencePattern) {
    m := &RecurrencePattern{
    }
    m.backingStore = ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackingStoreFactoryInstance();
    m.SetAdditionalData(make(map[string]any))
    return m
}
// CreateRecurrencePatternFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
// returns a Parsable when successful
func CreateRecurrencePatternFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewRecurrencePattern(), nil
}
// GetAdditionalData gets the AdditionalData property value. Stores additional data not described in the OpenAPI description found when deserializing. Can be used for serialization as well.
// returns a map[string]any when successful
func (m *RecurrencePattern) GetAdditionalData()(map[string]any) {
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
func (m *RecurrencePattern) GetBackingStore()(ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackingStore) {
    return m.backingStore
}
// GetDayOfMonth gets the dayOfMonth property value. The day of the month on which the event occurs. Required if type is absoluteMonthly or absoluteYearly.
// returns a *int32 when successful
func (m *RecurrencePattern) GetDayOfMonth()(*int32) {
    val, err := m.GetBackingStore().Get("dayOfMonth")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*int32)
    }
    return nil
}
// GetDaysOfWeek gets the daysOfWeek property value. A collection of the days of the week on which the event occurs. The possible values are: sunday, monday, tuesday, wednesday, thursday, friday, saturday. If type is relativeMonthly or relativeYearly, and daysOfWeek specifies more than one day, the event falls on the first day that satisfies the pattern.  Required if type is weekly, relativeMonthly, or relativeYearly.
// returns a []DayOfWeek when successful
func (m *RecurrencePattern) GetDaysOfWeek()([]DayOfWeek) {
    val, err := m.GetBackingStore().Get("daysOfWeek")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]DayOfWeek)
    }
    return nil
}
// GetFieldDeserializers the deserialization information for the current model
// returns a map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error) when successful
func (m *RecurrencePattern) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := make(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error))
    res["dayOfMonth"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetInt32Value()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetDayOfMonth(val)
        }
        return nil
    }
    res["daysOfWeek"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetCollectionOfEnumValues(ParseDayOfWeek)
        if err != nil {
            return err
        }
        if val != nil {
            res := make([]DayOfWeek, len(val))
            for i, v := range val {
                if v != nil {
                    res[i] = *(v.(*DayOfWeek))
                }
            }
            m.SetDaysOfWeek(res)
        }
        return nil
    }
    res["firstDayOfWeek"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetEnumValue(ParseDayOfWeek)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetFirstDayOfWeek(val.(*DayOfWeek))
        }
        return nil
    }
    res["index"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetEnumValue(ParseWeekIndex)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetIndex(val.(*WeekIndex))
        }
        return nil
    }
    res["interval"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetInt32Value()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetInterval(val)
        }
        return nil
    }
    res["month"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetInt32Value()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetMonth(val)
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
    res["type"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetEnumValue(ParseRecurrencePatternType)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetTypeEscaped(val.(*RecurrencePatternType))
        }
        return nil
    }
    return res
}
// GetFirstDayOfWeek gets the firstDayOfWeek property value. The first day of the week. The possible values are: sunday, monday, tuesday, wednesday, thursday, friday, saturday. Default is sunday. Required if type is weekly.
// returns a *DayOfWeek when successful
func (m *RecurrencePattern) GetFirstDayOfWeek()(*DayOfWeek) {
    val, err := m.GetBackingStore().Get("firstDayOfWeek")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*DayOfWeek)
    }
    return nil
}
// GetIndex gets the index property value. Specifies on which instance of the allowed days specified in daysOfWeek the event occurs, counted from the first instance in the month. The possible values are: first, second, third, fourth, last. Default is first. Optional and used if type is relativeMonthly or relativeYearly.
// returns a *WeekIndex when successful
func (m *RecurrencePattern) GetIndex()(*WeekIndex) {
    val, err := m.GetBackingStore().Get("index")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*WeekIndex)
    }
    return nil
}
// GetInterval gets the interval property value. The number of units between occurrences, where units can be in days, weeks, months, or years, depending on the type. Required.
// returns a *int32 when successful
func (m *RecurrencePattern) GetInterval()(*int32) {
    val, err := m.GetBackingStore().Get("interval")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*int32)
    }
    return nil
}
// GetMonth gets the month property value. The month in which the event occurs.  This is a number from 1 to 12.
// returns a *int32 when successful
func (m *RecurrencePattern) GetMonth()(*int32) {
    val, err := m.GetBackingStore().Get("month")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*int32)
    }
    return nil
}
// GetOdataType gets the @odata.type property value. The OdataType property
// returns a *string when successful
func (m *RecurrencePattern) GetOdataType()(*string) {
    val, err := m.GetBackingStore().Get("odataType")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetTypeEscaped gets the type property value. The recurrence pattern type: daily, weekly, absoluteMonthly, relativeMonthly, absoluteYearly, relativeYearly. Required. For more information, see values of type property.
// returns a *RecurrencePatternType when successful
func (m *RecurrencePattern) GetTypeEscaped()(*RecurrencePatternType) {
    val, err := m.GetBackingStore().Get("typeEscaped")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*RecurrencePatternType)
    }
    return nil
}
// Serialize serializes information the current object
func (m *RecurrencePattern) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    {
        err := writer.WriteInt32Value("dayOfMonth", m.GetDayOfMonth())
        if err != nil {
            return err
        }
    }
    if m.GetDaysOfWeek() != nil {
        err := writer.WriteCollectionOfStringValues("daysOfWeek", SerializeDayOfWeek(m.GetDaysOfWeek()))
        if err != nil {
            return err
        }
    }
    if m.GetFirstDayOfWeek() != nil {
        cast := (*m.GetFirstDayOfWeek()).String()
        err := writer.WriteStringValue("firstDayOfWeek", &cast)
        if err != nil {
            return err
        }
    }
    if m.GetIndex() != nil {
        cast := (*m.GetIndex()).String()
        err := writer.WriteStringValue("index", &cast)
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteInt32Value("interval", m.GetInterval())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteInt32Value("month", m.GetMonth())
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
    if m.GetTypeEscaped() != nil {
        cast := (*m.GetTypeEscaped()).String()
        err := writer.WriteStringValue("type", &cast)
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
func (m *RecurrencePattern) SetAdditionalData(value map[string]any)() {
    err := m.GetBackingStore().Set("additionalData", value)
    if err != nil {
        panic(err)
    }
}
// SetBackingStore sets the BackingStore property value. Stores model information.
func (m *RecurrencePattern) SetBackingStore(value ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackingStore)() {
    m.backingStore = value
}
// SetDayOfMonth sets the dayOfMonth property value. The day of the month on which the event occurs. Required if type is absoluteMonthly or absoluteYearly.
func (m *RecurrencePattern) SetDayOfMonth(value *int32)() {
    err := m.GetBackingStore().Set("dayOfMonth", value)
    if err != nil {
        panic(err)
    }
}
// SetDaysOfWeek sets the daysOfWeek property value. A collection of the days of the week on which the event occurs. The possible values are: sunday, monday, tuesday, wednesday, thursday, friday, saturday. If type is relativeMonthly or relativeYearly, and daysOfWeek specifies more than one day, the event falls on the first day that satisfies the pattern.  Required if type is weekly, relativeMonthly, or relativeYearly.
func (m *RecurrencePattern) SetDaysOfWeek(value []DayOfWeek)() {
    err := m.GetBackingStore().Set("daysOfWeek", value)
    if err != nil {
        panic(err)
    }
}
// SetFirstDayOfWeek sets the firstDayOfWeek property value. The first day of the week. The possible values are: sunday, monday, tuesday, wednesday, thursday, friday, saturday. Default is sunday. Required if type is weekly.
func (m *RecurrencePattern) SetFirstDayOfWeek(value *DayOfWeek)() {
    err := m.GetBackingStore().Set("firstDayOfWeek", value)
    if err != nil {
        panic(err)
    }
}
// SetIndex sets the index property value. Specifies on which instance of the allowed days specified in daysOfWeek the event occurs, counted from the first instance in the month. The possible values are: first, second, third, fourth, last. Default is first. Optional and used if type is relativeMonthly or relativeYearly.
func (m *RecurrencePattern) SetIndex(value *WeekIndex)() {
    err := m.GetBackingStore().Set("index", value)
    if err != nil {
        panic(err)
    }
}
// SetInterval sets the interval property value. The number of units between occurrences, where units can be in days, weeks, months, or years, depending on the type. Required.
func (m *RecurrencePattern) SetInterval(value *int32)() {
    err := m.GetBackingStore().Set("interval", value)
    if err != nil {
        panic(err)
    }
}
// SetMonth sets the month property value. The month in which the event occurs.  This is a number from 1 to 12.
func (m *RecurrencePattern) SetMonth(value *int32)() {
    err := m.GetBackingStore().Set("month", value)
    if err != nil {
        panic(err)
    }
}
// SetOdataType sets the @odata.type property value. The OdataType property
func (m *RecurrencePattern) SetOdataType(value *string)() {
    err := m.GetBackingStore().Set("odataType", value)
    if err != nil {
        panic(err)
    }
}
// SetTypeEscaped sets the type property value. The recurrence pattern type: daily, weekly, absoluteMonthly, relativeMonthly, absoluteYearly, relativeYearly. Required. For more information, see values of type property.
func (m *RecurrencePattern) SetTypeEscaped(value *RecurrencePatternType)() {
    err := m.GetBackingStore().Set("typeEscaped", value)
    if err != nil {
        panic(err)
    }
}
type RecurrencePatternable interface {
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.AdditionalDataHolder
    ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackedModel
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable
    GetBackingStore()(ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackingStore)
    GetDayOfMonth()(*int32)
    GetDaysOfWeek()([]DayOfWeek)
    GetFirstDayOfWeek()(*DayOfWeek)
    GetIndex()(*WeekIndex)
    GetInterval()(*int32)
    GetMonth()(*int32)
    GetOdataType()(*string)
    GetTypeEscaped()(*RecurrencePatternType)
    SetBackingStore(value ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackingStore)()
    SetDayOfMonth(value *int32)()
    SetDaysOfWeek(value []DayOfWeek)()
    SetFirstDayOfWeek(value *DayOfWeek)()
    SetIndex(value *WeekIndex)()
    SetInterval(value *int32)()
    SetMonth(value *int32)()
    SetOdataType(value *string)()
    SetTypeEscaped(value *RecurrencePatternType)()
}
