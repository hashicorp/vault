package models

import (
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

type FileStorageContainerCustomPropertyDictionary struct {
    Dictionary
}
// NewFileStorageContainerCustomPropertyDictionary instantiates a new FileStorageContainerCustomPropertyDictionary and sets the default values.
func NewFileStorageContainerCustomPropertyDictionary()(*FileStorageContainerCustomPropertyDictionary) {
    m := &FileStorageContainerCustomPropertyDictionary{
        Dictionary: *NewDictionary(),
    }
    return m
}
// CreateFileStorageContainerCustomPropertyDictionaryFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
// returns a Parsable when successful
func CreateFileStorageContainerCustomPropertyDictionaryFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewFileStorageContainerCustomPropertyDictionary(), nil
}
// GetFieldDeserializers the deserialization information for the current model
// returns a map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error) when successful
func (m *FileStorageContainerCustomPropertyDictionary) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := m.Dictionary.GetFieldDeserializers()
    return res
}
// Serialize serializes information the current object
func (m *FileStorageContainerCustomPropertyDictionary) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    err := m.Dictionary.Serialize(writer)
    if err != nil {
        return err
    }
    return nil
}
type FileStorageContainerCustomPropertyDictionaryable interface {
    Dictionaryable
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable
}
