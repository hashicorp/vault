package models

import (
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

type PrintUsage struct {
    Entity
}
// NewPrintUsage instantiates a new PrintUsage and sets the default values.
func NewPrintUsage()(*PrintUsage) {
    m := &PrintUsage{
        Entity: *NewEntity(),
    }
    return m
}
// CreatePrintUsageFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
// returns a Parsable when successful
func CreatePrintUsageFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
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
                    case "#microsoft.graph.printUsageByPrinter":
                        return NewPrintUsageByPrinter(), nil
                    case "#microsoft.graph.printUsageByUser":
                        return NewPrintUsageByUser(), nil
                }
            }
        }
    }
    return NewPrintUsage(), nil
}
// GetBlackAndWhitePageCount gets the blackAndWhitePageCount property value. The blackAndWhitePageCount property
// returns a *int64 when successful
func (m *PrintUsage) GetBlackAndWhitePageCount()(*int64) {
    val, err := m.GetBackingStore().Get("blackAndWhitePageCount")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*int64)
    }
    return nil
}
// GetColorPageCount gets the colorPageCount property value. The colorPageCount property
// returns a *int64 when successful
func (m *PrintUsage) GetColorPageCount()(*int64) {
    val, err := m.GetBackingStore().Get("colorPageCount")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*int64)
    }
    return nil
}
// GetCompletedBlackAndWhiteJobCount gets the completedBlackAndWhiteJobCount property value. The completedBlackAndWhiteJobCount property
// returns a *int64 when successful
func (m *PrintUsage) GetCompletedBlackAndWhiteJobCount()(*int64) {
    val, err := m.GetBackingStore().Get("completedBlackAndWhiteJobCount")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*int64)
    }
    return nil
}
// GetCompletedColorJobCount gets the completedColorJobCount property value. The completedColorJobCount property
// returns a *int64 when successful
func (m *PrintUsage) GetCompletedColorJobCount()(*int64) {
    val, err := m.GetBackingStore().Get("completedColorJobCount")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*int64)
    }
    return nil
}
// GetCompletedJobCount gets the completedJobCount property value. The completedJobCount property
// returns a *int64 when successful
func (m *PrintUsage) GetCompletedJobCount()(*int64) {
    val, err := m.GetBackingStore().Get("completedJobCount")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*int64)
    }
    return nil
}
// GetDoubleSidedSheetCount gets the doubleSidedSheetCount property value. The doubleSidedSheetCount property
// returns a *int64 when successful
func (m *PrintUsage) GetDoubleSidedSheetCount()(*int64) {
    val, err := m.GetBackingStore().Get("doubleSidedSheetCount")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*int64)
    }
    return nil
}
// GetFieldDeserializers the deserialization information for the current model
// returns a map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error) when successful
func (m *PrintUsage) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := m.Entity.GetFieldDeserializers()
    res["blackAndWhitePageCount"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetInt64Value()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetBlackAndWhitePageCount(val)
        }
        return nil
    }
    res["colorPageCount"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetInt64Value()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetColorPageCount(val)
        }
        return nil
    }
    res["completedBlackAndWhiteJobCount"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetInt64Value()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetCompletedBlackAndWhiteJobCount(val)
        }
        return nil
    }
    res["completedColorJobCount"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetInt64Value()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetCompletedColorJobCount(val)
        }
        return nil
    }
    res["completedJobCount"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetInt64Value()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetCompletedJobCount(val)
        }
        return nil
    }
    res["doubleSidedSheetCount"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetInt64Value()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetDoubleSidedSheetCount(val)
        }
        return nil
    }
    res["incompleteJobCount"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetInt64Value()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetIncompleteJobCount(val)
        }
        return nil
    }
    res["mediaSheetCount"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetInt64Value()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetMediaSheetCount(val)
        }
        return nil
    }
    res["pageCount"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetInt64Value()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetPageCount(val)
        }
        return nil
    }
    res["singleSidedSheetCount"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetInt64Value()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetSingleSidedSheetCount(val)
        }
        return nil
    }
    res["usageDate"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetDateOnlyValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetUsageDate(val)
        }
        return nil
    }
    return res
}
// GetIncompleteJobCount gets the incompleteJobCount property value. The incompleteJobCount property
// returns a *int64 when successful
func (m *PrintUsage) GetIncompleteJobCount()(*int64) {
    val, err := m.GetBackingStore().Get("incompleteJobCount")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*int64)
    }
    return nil
}
// GetMediaSheetCount gets the mediaSheetCount property value. The mediaSheetCount property
// returns a *int64 when successful
func (m *PrintUsage) GetMediaSheetCount()(*int64) {
    val, err := m.GetBackingStore().Get("mediaSheetCount")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*int64)
    }
    return nil
}
// GetPageCount gets the pageCount property value. The pageCount property
// returns a *int64 when successful
func (m *PrintUsage) GetPageCount()(*int64) {
    val, err := m.GetBackingStore().Get("pageCount")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*int64)
    }
    return nil
}
// GetSingleSidedSheetCount gets the singleSidedSheetCount property value. The singleSidedSheetCount property
// returns a *int64 when successful
func (m *PrintUsage) GetSingleSidedSheetCount()(*int64) {
    val, err := m.GetBackingStore().Get("singleSidedSheetCount")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*int64)
    }
    return nil
}
// GetUsageDate gets the usageDate property value. The usageDate property
// returns a *DateOnly when successful
func (m *PrintUsage) GetUsageDate()(*i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.DateOnly) {
    val, err := m.GetBackingStore().Get("usageDate")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.DateOnly)
    }
    return nil
}
// Serialize serializes information the current object
func (m *PrintUsage) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    err := m.Entity.Serialize(writer)
    if err != nil {
        return err
    }
    {
        err = writer.WriteInt64Value("blackAndWhitePageCount", m.GetBlackAndWhitePageCount())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteInt64Value("colorPageCount", m.GetColorPageCount())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteInt64Value("completedBlackAndWhiteJobCount", m.GetCompletedBlackAndWhiteJobCount())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteInt64Value("completedColorJobCount", m.GetCompletedColorJobCount())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteInt64Value("completedJobCount", m.GetCompletedJobCount())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteInt64Value("doubleSidedSheetCount", m.GetDoubleSidedSheetCount())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteInt64Value("incompleteJobCount", m.GetIncompleteJobCount())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteInt64Value("mediaSheetCount", m.GetMediaSheetCount())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteInt64Value("pageCount", m.GetPageCount())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteInt64Value("singleSidedSheetCount", m.GetSingleSidedSheetCount())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteDateOnlyValue("usageDate", m.GetUsageDate())
        if err != nil {
            return err
        }
    }
    return nil
}
// SetBlackAndWhitePageCount sets the blackAndWhitePageCount property value. The blackAndWhitePageCount property
func (m *PrintUsage) SetBlackAndWhitePageCount(value *int64)() {
    err := m.GetBackingStore().Set("blackAndWhitePageCount", value)
    if err != nil {
        panic(err)
    }
}
// SetColorPageCount sets the colorPageCount property value. The colorPageCount property
func (m *PrintUsage) SetColorPageCount(value *int64)() {
    err := m.GetBackingStore().Set("colorPageCount", value)
    if err != nil {
        panic(err)
    }
}
// SetCompletedBlackAndWhiteJobCount sets the completedBlackAndWhiteJobCount property value. The completedBlackAndWhiteJobCount property
func (m *PrintUsage) SetCompletedBlackAndWhiteJobCount(value *int64)() {
    err := m.GetBackingStore().Set("completedBlackAndWhiteJobCount", value)
    if err != nil {
        panic(err)
    }
}
// SetCompletedColorJobCount sets the completedColorJobCount property value. The completedColorJobCount property
func (m *PrintUsage) SetCompletedColorJobCount(value *int64)() {
    err := m.GetBackingStore().Set("completedColorJobCount", value)
    if err != nil {
        panic(err)
    }
}
// SetCompletedJobCount sets the completedJobCount property value. The completedJobCount property
func (m *PrintUsage) SetCompletedJobCount(value *int64)() {
    err := m.GetBackingStore().Set("completedJobCount", value)
    if err != nil {
        panic(err)
    }
}
// SetDoubleSidedSheetCount sets the doubleSidedSheetCount property value. The doubleSidedSheetCount property
func (m *PrintUsage) SetDoubleSidedSheetCount(value *int64)() {
    err := m.GetBackingStore().Set("doubleSidedSheetCount", value)
    if err != nil {
        panic(err)
    }
}
// SetIncompleteJobCount sets the incompleteJobCount property value. The incompleteJobCount property
func (m *PrintUsage) SetIncompleteJobCount(value *int64)() {
    err := m.GetBackingStore().Set("incompleteJobCount", value)
    if err != nil {
        panic(err)
    }
}
// SetMediaSheetCount sets the mediaSheetCount property value. The mediaSheetCount property
func (m *PrintUsage) SetMediaSheetCount(value *int64)() {
    err := m.GetBackingStore().Set("mediaSheetCount", value)
    if err != nil {
        panic(err)
    }
}
// SetPageCount sets the pageCount property value. The pageCount property
func (m *PrintUsage) SetPageCount(value *int64)() {
    err := m.GetBackingStore().Set("pageCount", value)
    if err != nil {
        panic(err)
    }
}
// SetSingleSidedSheetCount sets the singleSidedSheetCount property value. The singleSidedSheetCount property
func (m *PrintUsage) SetSingleSidedSheetCount(value *int64)() {
    err := m.GetBackingStore().Set("singleSidedSheetCount", value)
    if err != nil {
        panic(err)
    }
}
// SetUsageDate sets the usageDate property value. The usageDate property
func (m *PrintUsage) SetUsageDate(value *i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.DateOnly)() {
    err := m.GetBackingStore().Set("usageDate", value)
    if err != nil {
        panic(err)
    }
}
type PrintUsageable interface {
    Entityable
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable
    GetBlackAndWhitePageCount()(*int64)
    GetColorPageCount()(*int64)
    GetCompletedBlackAndWhiteJobCount()(*int64)
    GetCompletedColorJobCount()(*int64)
    GetCompletedJobCount()(*int64)
    GetDoubleSidedSheetCount()(*int64)
    GetIncompleteJobCount()(*int64)
    GetMediaSheetCount()(*int64)
    GetPageCount()(*int64)
    GetSingleSidedSheetCount()(*int64)
    GetUsageDate()(*i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.DateOnly)
    SetBlackAndWhitePageCount(value *int64)()
    SetColorPageCount(value *int64)()
    SetCompletedBlackAndWhiteJobCount(value *int64)()
    SetCompletedColorJobCount(value *int64)()
    SetCompletedJobCount(value *int64)()
    SetDoubleSidedSheetCount(value *int64)()
    SetIncompleteJobCount(value *int64)()
    SetMediaSheetCount(value *int64)()
    SetPageCount(value *int64)()
    SetSingleSidedSheetCount(value *int64)()
    SetUsageDate(value *i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.DateOnly)()
}
