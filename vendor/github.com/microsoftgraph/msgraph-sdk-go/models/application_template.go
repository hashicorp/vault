package models

import (
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

type ApplicationTemplate struct {
    Entity
}
// NewApplicationTemplate instantiates a new ApplicationTemplate and sets the default values.
func NewApplicationTemplate()(*ApplicationTemplate) {
    m := &ApplicationTemplate{
        Entity: *NewEntity(),
    }
    return m
}
// CreateApplicationTemplateFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
// returns a Parsable when successful
func CreateApplicationTemplateFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewApplicationTemplate(), nil
}
// GetCategories gets the categories property value. The list of categories for the application. Supported values can be: Collaboration, Business Management, Consumer, Content management, CRM, Data services, Developer services, E-commerce, Education, ERP, Finance, Health, Human resources, IT infrastructure, Mail, Management, Marketing, Media, Productivity, Project management, Telecommunications, Tools, Travel, and Web design & hosting.
// returns a []string when successful
func (m *ApplicationTemplate) GetCategories()([]string) {
    val, err := m.GetBackingStore().Get("categories")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]string)
    }
    return nil
}
// GetDescription gets the description property value. A description of the application.
// returns a *string when successful
func (m *ApplicationTemplate) GetDescription()(*string) {
    val, err := m.GetBackingStore().Get("description")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetDisplayName gets the displayName property value. The name of the application.
// returns a *string when successful
func (m *ApplicationTemplate) GetDisplayName()(*string) {
    val, err := m.GetBackingStore().Get("displayName")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetFieldDeserializers the deserialization information for the current model
// returns a map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error) when successful
func (m *ApplicationTemplate) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := m.Entity.GetFieldDeserializers()
    res["categories"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
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
            m.SetCategories(res)
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
    res["homePageUrl"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetHomePageUrl(val)
        }
        return nil
    }
    res["logoUrl"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetLogoUrl(val)
        }
        return nil
    }
    res["publisher"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetPublisher(val)
        }
        return nil
    }
    res["supportedProvisioningTypes"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
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
            m.SetSupportedProvisioningTypes(res)
        }
        return nil
    }
    res["supportedSingleSignOnModes"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
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
            m.SetSupportedSingleSignOnModes(res)
        }
        return nil
    }
    return res
}
// GetHomePageUrl gets the homePageUrl property value. The home page URL of the application.
// returns a *string when successful
func (m *ApplicationTemplate) GetHomePageUrl()(*string) {
    val, err := m.GetBackingStore().Get("homePageUrl")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetLogoUrl gets the logoUrl property value. The URL to get the logo for this application.
// returns a *string when successful
func (m *ApplicationTemplate) GetLogoUrl()(*string) {
    val, err := m.GetBackingStore().Get("logoUrl")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetPublisher gets the publisher property value. The name of the publisher for this application.
// returns a *string when successful
func (m *ApplicationTemplate) GetPublisher()(*string) {
    val, err := m.GetBackingStore().Get("publisher")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetSupportedProvisioningTypes gets the supportedProvisioningTypes property value. The list of provisioning modes supported by this application. The only valid value is sync.
// returns a []string when successful
func (m *ApplicationTemplate) GetSupportedProvisioningTypes()([]string) {
    val, err := m.GetBackingStore().Get("supportedProvisioningTypes")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]string)
    }
    return nil
}
// GetSupportedSingleSignOnModes gets the supportedSingleSignOnModes property value. The list of single sign-on modes supported by this application. The supported values are oidc, password, saml, and notSupported.
// returns a []string when successful
func (m *ApplicationTemplate) GetSupportedSingleSignOnModes()([]string) {
    val, err := m.GetBackingStore().Get("supportedSingleSignOnModes")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]string)
    }
    return nil
}
// Serialize serializes information the current object
func (m *ApplicationTemplate) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    err := m.Entity.Serialize(writer)
    if err != nil {
        return err
    }
    if m.GetCategories() != nil {
        err = writer.WriteCollectionOfStringValues("categories", m.GetCategories())
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
        err = writer.WriteStringValue("homePageUrl", m.GetHomePageUrl())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteStringValue("logoUrl", m.GetLogoUrl())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteStringValue("publisher", m.GetPublisher())
        if err != nil {
            return err
        }
    }
    if m.GetSupportedProvisioningTypes() != nil {
        err = writer.WriteCollectionOfStringValues("supportedProvisioningTypes", m.GetSupportedProvisioningTypes())
        if err != nil {
            return err
        }
    }
    if m.GetSupportedSingleSignOnModes() != nil {
        err = writer.WriteCollectionOfStringValues("supportedSingleSignOnModes", m.GetSupportedSingleSignOnModes())
        if err != nil {
            return err
        }
    }
    return nil
}
// SetCategories sets the categories property value. The list of categories for the application. Supported values can be: Collaboration, Business Management, Consumer, Content management, CRM, Data services, Developer services, E-commerce, Education, ERP, Finance, Health, Human resources, IT infrastructure, Mail, Management, Marketing, Media, Productivity, Project management, Telecommunications, Tools, Travel, and Web design & hosting.
func (m *ApplicationTemplate) SetCategories(value []string)() {
    err := m.GetBackingStore().Set("categories", value)
    if err != nil {
        panic(err)
    }
}
// SetDescription sets the description property value. A description of the application.
func (m *ApplicationTemplate) SetDescription(value *string)() {
    err := m.GetBackingStore().Set("description", value)
    if err != nil {
        panic(err)
    }
}
// SetDisplayName sets the displayName property value. The name of the application.
func (m *ApplicationTemplate) SetDisplayName(value *string)() {
    err := m.GetBackingStore().Set("displayName", value)
    if err != nil {
        panic(err)
    }
}
// SetHomePageUrl sets the homePageUrl property value. The home page URL of the application.
func (m *ApplicationTemplate) SetHomePageUrl(value *string)() {
    err := m.GetBackingStore().Set("homePageUrl", value)
    if err != nil {
        panic(err)
    }
}
// SetLogoUrl sets the logoUrl property value. The URL to get the logo for this application.
func (m *ApplicationTemplate) SetLogoUrl(value *string)() {
    err := m.GetBackingStore().Set("logoUrl", value)
    if err != nil {
        panic(err)
    }
}
// SetPublisher sets the publisher property value. The name of the publisher for this application.
func (m *ApplicationTemplate) SetPublisher(value *string)() {
    err := m.GetBackingStore().Set("publisher", value)
    if err != nil {
        panic(err)
    }
}
// SetSupportedProvisioningTypes sets the supportedProvisioningTypes property value. The list of provisioning modes supported by this application. The only valid value is sync.
func (m *ApplicationTemplate) SetSupportedProvisioningTypes(value []string)() {
    err := m.GetBackingStore().Set("supportedProvisioningTypes", value)
    if err != nil {
        panic(err)
    }
}
// SetSupportedSingleSignOnModes sets the supportedSingleSignOnModes property value. The list of single sign-on modes supported by this application. The supported values are oidc, password, saml, and notSupported.
func (m *ApplicationTemplate) SetSupportedSingleSignOnModes(value []string)() {
    err := m.GetBackingStore().Set("supportedSingleSignOnModes", value)
    if err != nil {
        panic(err)
    }
}
type ApplicationTemplateable interface {
    Entityable
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable
    GetCategories()([]string)
    GetDescription()(*string)
    GetDisplayName()(*string)
    GetHomePageUrl()(*string)
    GetLogoUrl()(*string)
    GetPublisher()(*string)
    GetSupportedProvisioningTypes()([]string)
    GetSupportedSingleSignOnModes()([]string)
    SetCategories(value []string)()
    SetDescription(value *string)()
    SetDisplayName(value *string)()
    SetHomePageUrl(value *string)()
    SetLogoUrl(value *string)()
    SetPublisher(value *string)()
    SetSupportedProvisioningTypes(value []string)()
    SetSupportedSingleSignOnModes(value []string)()
}
