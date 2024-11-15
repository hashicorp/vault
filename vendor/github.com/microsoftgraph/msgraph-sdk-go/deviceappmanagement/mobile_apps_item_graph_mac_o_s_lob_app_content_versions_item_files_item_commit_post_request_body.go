package deviceappmanagement

import (
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
    iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242 "github.com/microsoftgraph/msgraph-sdk-go/models"
    ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e "github.com/microsoft/kiota-abstractions-go/store"
)

type MobileAppsItemGraphMacOSLobAppContentVersionsItemFilesItemCommitPostRequestBody struct {
    // Stores model information.
    backingStore ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackingStore
}
// NewMobileAppsItemGraphMacOSLobAppContentVersionsItemFilesItemCommitPostRequestBody instantiates a new MobileAppsItemGraphMacOSLobAppContentVersionsItemFilesItemCommitPostRequestBody and sets the default values.
func NewMobileAppsItemGraphMacOSLobAppContentVersionsItemFilesItemCommitPostRequestBody()(*MobileAppsItemGraphMacOSLobAppContentVersionsItemFilesItemCommitPostRequestBody) {
    m := &MobileAppsItemGraphMacOSLobAppContentVersionsItemFilesItemCommitPostRequestBody{
    }
    m.backingStore = ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackingStoreFactoryInstance();
    m.SetAdditionalData(make(map[string]any))
    return m
}
// CreateMobileAppsItemGraphMacOSLobAppContentVersionsItemFilesItemCommitPostRequestBodyFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
// returns a Parsable when successful
func CreateMobileAppsItemGraphMacOSLobAppContentVersionsItemFilesItemCommitPostRequestBodyFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewMobileAppsItemGraphMacOSLobAppContentVersionsItemFilesItemCommitPostRequestBody(), nil
}
// GetAdditionalData gets the AdditionalData property value. Stores additional data not described in the OpenAPI description found when deserializing. Can be used for serialization as well.
// returns a map[string]any when successful
func (m *MobileAppsItemGraphMacOSLobAppContentVersionsItemFilesItemCommitPostRequestBody) GetAdditionalData()(map[string]any) {
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
func (m *MobileAppsItemGraphMacOSLobAppContentVersionsItemFilesItemCommitPostRequestBody) GetBackingStore()(ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackingStore) {
    return m.backingStore
}
// GetFieldDeserializers the deserialization information for the current model
// returns a map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error) when successful
func (m *MobileAppsItemGraphMacOSLobAppContentVersionsItemFilesItemCommitPostRequestBody) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := make(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error))
    res["fileEncryptionInfo"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetObjectValue(iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.CreateFileEncryptionInfoFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetFileEncryptionInfo(val.(iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.FileEncryptionInfoable))
        }
        return nil
    }
    return res
}
// GetFileEncryptionInfo gets the fileEncryptionInfo property value. The fileEncryptionInfo property
// returns a FileEncryptionInfoable when successful
func (m *MobileAppsItemGraphMacOSLobAppContentVersionsItemFilesItemCommitPostRequestBody) GetFileEncryptionInfo()(iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.FileEncryptionInfoable) {
    val, err := m.GetBackingStore().Get("fileEncryptionInfo")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.FileEncryptionInfoable)
    }
    return nil
}
// Serialize serializes information the current object
func (m *MobileAppsItemGraphMacOSLobAppContentVersionsItemFilesItemCommitPostRequestBody) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    {
        err := writer.WriteObjectValue("fileEncryptionInfo", m.GetFileEncryptionInfo())
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
func (m *MobileAppsItemGraphMacOSLobAppContentVersionsItemFilesItemCommitPostRequestBody) SetAdditionalData(value map[string]any)() {
    err := m.GetBackingStore().Set("additionalData", value)
    if err != nil {
        panic(err)
    }
}
// SetBackingStore sets the BackingStore property value. Stores model information.
func (m *MobileAppsItemGraphMacOSLobAppContentVersionsItemFilesItemCommitPostRequestBody) SetBackingStore(value ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackingStore)() {
    m.backingStore = value
}
// SetFileEncryptionInfo sets the fileEncryptionInfo property value. The fileEncryptionInfo property
func (m *MobileAppsItemGraphMacOSLobAppContentVersionsItemFilesItemCommitPostRequestBody) SetFileEncryptionInfo(value iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.FileEncryptionInfoable)() {
    err := m.GetBackingStore().Set("fileEncryptionInfo", value)
    if err != nil {
        panic(err)
    }
}
type MobileAppsItemGraphMacOSLobAppContentVersionsItemFilesItemCommitPostRequestBodyable interface {
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.AdditionalDataHolder
    ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackedModel
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable
    GetBackingStore()(ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackingStore)
    GetFileEncryptionInfo()(iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.FileEncryptionInfoable)
    SetBackingStore(value ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackingStore)()
    SetFileEncryptionInfo(value iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.FileEncryptionInfoable)()
}
