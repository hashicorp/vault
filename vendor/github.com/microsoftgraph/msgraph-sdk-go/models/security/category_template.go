package security

import (
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

type CategoryTemplate struct {
    FilePlanDescriptorTemplate
}
// NewCategoryTemplate instantiates a new CategoryTemplate and sets the default values.
func NewCategoryTemplate()(*CategoryTemplate) {
    m := &CategoryTemplate{
        FilePlanDescriptorTemplate: *NewFilePlanDescriptorTemplate(),
    }
    return m
}
// CreateCategoryTemplateFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
// returns a Parsable when successful
func CreateCategoryTemplateFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewCategoryTemplate(), nil
}
// GetFieldDeserializers the deserialization information for the current model
// returns a map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error) when successful
func (m *CategoryTemplate) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := m.FilePlanDescriptorTemplate.GetFieldDeserializers()
    res["subcategories"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetCollectionOfObjectValues(CreateSubcategoryTemplateFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            res := make([]SubcategoryTemplateable, len(val))
            for i, v := range val {
                if v != nil {
                    res[i] = v.(SubcategoryTemplateable)
                }
            }
            m.SetSubcategories(res)
        }
        return nil
    }
    return res
}
// GetSubcategories gets the subcategories property value. Represents all subcategories under a particular category.
// returns a []SubcategoryTemplateable when successful
func (m *CategoryTemplate) GetSubcategories()([]SubcategoryTemplateable) {
    val, err := m.GetBackingStore().Get("subcategories")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]SubcategoryTemplateable)
    }
    return nil
}
// Serialize serializes information the current object
func (m *CategoryTemplate) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    err := m.FilePlanDescriptorTemplate.Serialize(writer)
    if err != nil {
        return err
    }
    if m.GetSubcategories() != nil {
        cast := make([]i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, len(m.GetSubcategories()))
        for i, v := range m.GetSubcategories() {
            if v != nil {
                cast[i] = v.(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable)
            }
        }
        err = writer.WriteCollectionOfObjectValues("subcategories", cast)
        if err != nil {
            return err
        }
    }
    return nil
}
// SetSubcategories sets the subcategories property value. Represents all subcategories under a particular category.
func (m *CategoryTemplate) SetSubcategories(value []SubcategoryTemplateable)() {
    err := m.GetBackingStore().Set("subcategories", value)
    if err != nil {
        panic(err)
    }
}
type CategoryTemplateable interface {
    FilePlanDescriptorTemplateable
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable
    GetSubcategories()([]SubcategoryTemplateable)
    SetSubcategories(value []SubcategoryTemplateable)()
}
