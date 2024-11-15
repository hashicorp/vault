package models

import (
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

type HorizontalSection struct {
    Entity
}
// NewHorizontalSection instantiates a new HorizontalSection and sets the default values.
func NewHorizontalSection()(*HorizontalSection) {
    m := &HorizontalSection{
        Entity: *NewEntity(),
    }
    return m
}
// CreateHorizontalSectionFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
// returns a Parsable when successful
func CreateHorizontalSectionFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewHorizontalSection(), nil
}
// GetColumns gets the columns property value. The set of vertical columns in this section.
// returns a []HorizontalSectionColumnable when successful
func (m *HorizontalSection) GetColumns()([]HorizontalSectionColumnable) {
    val, err := m.GetBackingStore().Get("columns")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]HorizontalSectionColumnable)
    }
    return nil
}
// GetEmphasis gets the emphasis property value. Enumeration value that indicates the emphasis of the section background. The possible values are: none, netural, soft, strong, unknownFutureValue.
// returns a *SectionEmphasisType when successful
func (m *HorizontalSection) GetEmphasis()(*SectionEmphasisType) {
    val, err := m.GetBackingStore().Get("emphasis")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*SectionEmphasisType)
    }
    return nil
}
// GetFieldDeserializers the deserialization information for the current model
// returns a map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error) when successful
func (m *HorizontalSection) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := m.Entity.GetFieldDeserializers()
    res["columns"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetCollectionOfObjectValues(CreateHorizontalSectionColumnFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            res := make([]HorizontalSectionColumnable, len(val))
            for i, v := range val {
                if v != nil {
                    res[i] = v.(HorizontalSectionColumnable)
                }
            }
            m.SetColumns(res)
        }
        return nil
    }
    res["emphasis"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetEnumValue(ParseSectionEmphasisType)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetEmphasis(val.(*SectionEmphasisType))
        }
        return nil
    }
    res["layout"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetEnumValue(ParseHorizontalSectionLayoutType)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetLayout(val.(*HorizontalSectionLayoutType))
        }
        return nil
    }
    return res
}
// GetLayout gets the layout property value. Layout type of the section. The possible values are: none, oneColumn, twoColumns, threeColumns, oneThirdLeftColumn, oneThirdRightColumn, fullWidth, unknownFutureValue.
// returns a *HorizontalSectionLayoutType when successful
func (m *HorizontalSection) GetLayout()(*HorizontalSectionLayoutType) {
    val, err := m.GetBackingStore().Get("layout")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*HorizontalSectionLayoutType)
    }
    return nil
}
// Serialize serializes information the current object
func (m *HorizontalSection) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    err := m.Entity.Serialize(writer)
    if err != nil {
        return err
    }
    if m.GetColumns() != nil {
        cast := make([]i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, len(m.GetColumns()))
        for i, v := range m.GetColumns() {
            if v != nil {
                cast[i] = v.(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable)
            }
        }
        err = writer.WriteCollectionOfObjectValues("columns", cast)
        if err != nil {
            return err
        }
    }
    if m.GetEmphasis() != nil {
        cast := (*m.GetEmphasis()).String()
        err = writer.WriteStringValue("emphasis", &cast)
        if err != nil {
            return err
        }
    }
    if m.GetLayout() != nil {
        cast := (*m.GetLayout()).String()
        err = writer.WriteStringValue("layout", &cast)
        if err != nil {
            return err
        }
    }
    return nil
}
// SetColumns sets the columns property value. The set of vertical columns in this section.
func (m *HorizontalSection) SetColumns(value []HorizontalSectionColumnable)() {
    err := m.GetBackingStore().Set("columns", value)
    if err != nil {
        panic(err)
    }
}
// SetEmphasis sets the emphasis property value. Enumeration value that indicates the emphasis of the section background. The possible values are: none, netural, soft, strong, unknownFutureValue.
func (m *HorizontalSection) SetEmphasis(value *SectionEmphasisType)() {
    err := m.GetBackingStore().Set("emphasis", value)
    if err != nil {
        panic(err)
    }
}
// SetLayout sets the layout property value. Layout type of the section. The possible values are: none, oneColumn, twoColumns, threeColumns, oneThirdLeftColumn, oneThirdRightColumn, fullWidth, unknownFutureValue.
func (m *HorizontalSection) SetLayout(value *HorizontalSectionLayoutType)() {
    err := m.GetBackingStore().Set("layout", value)
    if err != nil {
        panic(err)
    }
}
type HorizontalSectionable interface {
    Entityable
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable
    GetColumns()([]HorizontalSectionColumnable)
    GetEmphasis()(*SectionEmphasisType)
    GetLayout()(*HorizontalSectionLayoutType)
    SetColumns(value []HorizontalSectionColumnable)()
    SetEmphasis(value *SectionEmphasisType)()
    SetLayout(value *HorizontalSectionLayoutType)()
}
