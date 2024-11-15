package models

import (
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

type AddressBookAccountTargetContent struct {
    AccountTargetContent
}
// NewAddressBookAccountTargetContent instantiates a new AddressBookAccountTargetContent and sets the default values.
func NewAddressBookAccountTargetContent()(*AddressBookAccountTargetContent) {
    m := &AddressBookAccountTargetContent{
        AccountTargetContent: *NewAccountTargetContent(),
    }
    odataTypeValue := "#microsoft.graph.addressBookAccountTargetContent"
    m.SetOdataType(&odataTypeValue)
    return m
}
// CreateAddressBookAccountTargetContentFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
// returns a Parsable when successful
func CreateAddressBookAccountTargetContentFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewAddressBookAccountTargetContent(), nil
}
// GetAccountTargetEmails gets the accountTargetEmails property value. List of user emails targeted for an attack simulation training campaign.
// returns a []string when successful
func (m *AddressBookAccountTargetContent) GetAccountTargetEmails()([]string) {
    val, err := m.GetBackingStore().Get("accountTargetEmails")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]string)
    }
    return nil
}
// GetFieldDeserializers the deserialization information for the current model
// returns a map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error) when successful
func (m *AddressBookAccountTargetContent) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := m.AccountTargetContent.GetFieldDeserializers()
    res["accountTargetEmails"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
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
            m.SetAccountTargetEmails(res)
        }
        return nil
    }
    return res
}
// Serialize serializes information the current object
func (m *AddressBookAccountTargetContent) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    err := m.AccountTargetContent.Serialize(writer)
    if err != nil {
        return err
    }
    if m.GetAccountTargetEmails() != nil {
        err = writer.WriteCollectionOfStringValues("accountTargetEmails", m.GetAccountTargetEmails())
        if err != nil {
            return err
        }
    }
    return nil
}
// SetAccountTargetEmails sets the accountTargetEmails property value. List of user emails targeted for an attack simulation training campaign.
func (m *AddressBookAccountTargetContent) SetAccountTargetEmails(value []string)() {
    err := m.GetBackingStore().Set("accountTargetEmails", value)
    if err != nil {
        panic(err)
    }
}
type AddressBookAccountTargetContentable interface {
    AccountTargetContentable
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable
    GetAccountTargetEmails()([]string)
    SetAccountTargetEmails(value []string)()
}
