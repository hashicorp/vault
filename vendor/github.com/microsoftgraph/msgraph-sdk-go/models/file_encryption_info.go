package models

import (
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
    ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e "github.com/microsoft/kiota-abstractions-go/store"
)

// FileEncryptionInfo contains properties for file encryption information for the content version of a line of business app.
type FileEncryptionInfo struct {
    // Stores model information.
    backingStore ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackingStore
}
// NewFileEncryptionInfo instantiates a new FileEncryptionInfo and sets the default values.
func NewFileEncryptionInfo()(*FileEncryptionInfo) {
    m := &FileEncryptionInfo{
    }
    m.backingStore = ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackingStoreFactoryInstance();
    m.SetAdditionalData(make(map[string]any))
    return m
}
// CreateFileEncryptionInfoFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
// returns a Parsable when successful
func CreateFileEncryptionInfoFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewFileEncryptionInfo(), nil
}
// GetAdditionalData gets the AdditionalData property value. Stores additional data not described in the OpenAPI description found when deserializing. Can be used for serialization as well.
// returns a map[string]any when successful
func (m *FileEncryptionInfo) GetAdditionalData()(map[string]any) {
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
func (m *FileEncryptionInfo) GetBackingStore()(ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackingStore) {
    return m.backingStore
}
// GetEncryptionKey gets the encryptionKey property value. The key used to encrypt the file content.
// returns a []byte when successful
func (m *FileEncryptionInfo) GetEncryptionKey()([]byte) {
    val, err := m.GetBackingStore().Get("encryptionKey")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]byte)
    }
    return nil
}
// GetFieldDeserializers the deserialization information for the current model
// returns a map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error) when successful
func (m *FileEncryptionInfo) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := make(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error))
    res["encryptionKey"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetByteArrayValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetEncryptionKey(val)
        }
        return nil
    }
    res["fileDigest"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetByteArrayValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetFileDigest(val)
        }
        return nil
    }
    res["fileDigestAlgorithm"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetFileDigestAlgorithm(val)
        }
        return nil
    }
    res["initializationVector"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetByteArrayValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetInitializationVector(val)
        }
        return nil
    }
    res["mac"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetByteArrayValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetMac(val)
        }
        return nil
    }
    res["macKey"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetByteArrayValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetMacKey(val)
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
    res["profileIdentifier"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetProfileIdentifier(val)
        }
        return nil
    }
    return res
}
// GetFileDigest gets the fileDigest property value. The file digest prior to encryption. ProfileVersion1 requires a non-null FileDigest.
// returns a []byte when successful
func (m *FileEncryptionInfo) GetFileDigest()([]byte) {
    val, err := m.GetBackingStore().Get("fileDigest")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]byte)
    }
    return nil
}
// GetFileDigestAlgorithm gets the fileDigestAlgorithm property value. The file digest algorithm. ProfileVersion1 currently only supports SHA256 for the FileDigestAlgorithm.
// returns a *string when successful
func (m *FileEncryptionInfo) GetFileDigestAlgorithm()(*string) {
    val, err := m.GetBackingStore().Get("fileDigestAlgorithm")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetInitializationVector gets the initializationVector property value. The initialization vector (IV) used for the encryption algorithm. Must be 16 bytes.
// returns a []byte when successful
func (m *FileEncryptionInfo) GetInitializationVector()([]byte) {
    val, err := m.GetBackingStore().Get("initializationVector")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]byte)
    }
    return nil
}
// GetMac gets the mac property value. The hash of the concatenation of the IV and encrypted file content. Must be 32 bytes.
// returns a []byte when successful
func (m *FileEncryptionInfo) GetMac()([]byte) {
    val, err := m.GetBackingStore().Get("mac")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]byte)
    }
    return nil
}
// GetMacKey gets the macKey property value. The key used to compute the message authentication code of the concatenation of the IV and encrypted file content. Must be 32 bytes.
// returns a []byte when successful
func (m *FileEncryptionInfo) GetMacKey()([]byte) {
    val, err := m.GetBackingStore().Get("macKey")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]byte)
    }
    return nil
}
// GetOdataType gets the @odata.type property value. The OdataType property
// returns a *string when successful
func (m *FileEncryptionInfo) GetOdataType()(*string) {
    val, err := m.GetBackingStore().Get("odataType")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetProfileIdentifier gets the profileIdentifier property value. The profile identifier. Maps to the strategy used to encrypt the file. Currently, only ProfileVersion1 is supported.
// returns a *string when successful
func (m *FileEncryptionInfo) GetProfileIdentifier()(*string) {
    val, err := m.GetBackingStore().Get("profileIdentifier")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// Serialize serializes information the current object
func (m *FileEncryptionInfo) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    {
        err := writer.WriteByteArrayValue("encryptionKey", m.GetEncryptionKey())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteByteArrayValue("fileDigest", m.GetFileDigest())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteStringValue("fileDigestAlgorithm", m.GetFileDigestAlgorithm())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteByteArrayValue("initializationVector", m.GetInitializationVector())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteByteArrayValue("mac", m.GetMac())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteByteArrayValue("macKey", m.GetMacKey())
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
        err := writer.WriteStringValue("profileIdentifier", m.GetProfileIdentifier())
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
func (m *FileEncryptionInfo) SetAdditionalData(value map[string]any)() {
    err := m.GetBackingStore().Set("additionalData", value)
    if err != nil {
        panic(err)
    }
}
// SetBackingStore sets the BackingStore property value. Stores model information.
func (m *FileEncryptionInfo) SetBackingStore(value ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackingStore)() {
    m.backingStore = value
}
// SetEncryptionKey sets the encryptionKey property value. The key used to encrypt the file content.
func (m *FileEncryptionInfo) SetEncryptionKey(value []byte)() {
    err := m.GetBackingStore().Set("encryptionKey", value)
    if err != nil {
        panic(err)
    }
}
// SetFileDigest sets the fileDigest property value. The file digest prior to encryption. ProfileVersion1 requires a non-null FileDigest.
func (m *FileEncryptionInfo) SetFileDigest(value []byte)() {
    err := m.GetBackingStore().Set("fileDigest", value)
    if err != nil {
        panic(err)
    }
}
// SetFileDigestAlgorithm sets the fileDigestAlgorithm property value. The file digest algorithm. ProfileVersion1 currently only supports SHA256 for the FileDigestAlgorithm.
func (m *FileEncryptionInfo) SetFileDigestAlgorithm(value *string)() {
    err := m.GetBackingStore().Set("fileDigestAlgorithm", value)
    if err != nil {
        panic(err)
    }
}
// SetInitializationVector sets the initializationVector property value. The initialization vector (IV) used for the encryption algorithm. Must be 16 bytes.
func (m *FileEncryptionInfo) SetInitializationVector(value []byte)() {
    err := m.GetBackingStore().Set("initializationVector", value)
    if err != nil {
        panic(err)
    }
}
// SetMac sets the mac property value. The hash of the concatenation of the IV and encrypted file content. Must be 32 bytes.
func (m *FileEncryptionInfo) SetMac(value []byte)() {
    err := m.GetBackingStore().Set("mac", value)
    if err != nil {
        panic(err)
    }
}
// SetMacKey sets the macKey property value. The key used to compute the message authentication code of the concatenation of the IV and encrypted file content. Must be 32 bytes.
func (m *FileEncryptionInfo) SetMacKey(value []byte)() {
    err := m.GetBackingStore().Set("macKey", value)
    if err != nil {
        panic(err)
    }
}
// SetOdataType sets the @odata.type property value. The OdataType property
func (m *FileEncryptionInfo) SetOdataType(value *string)() {
    err := m.GetBackingStore().Set("odataType", value)
    if err != nil {
        panic(err)
    }
}
// SetProfileIdentifier sets the profileIdentifier property value. The profile identifier. Maps to the strategy used to encrypt the file. Currently, only ProfileVersion1 is supported.
func (m *FileEncryptionInfo) SetProfileIdentifier(value *string)() {
    err := m.GetBackingStore().Set("profileIdentifier", value)
    if err != nil {
        panic(err)
    }
}
type FileEncryptionInfoable interface {
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.AdditionalDataHolder
    ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackedModel
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable
    GetBackingStore()(ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackingStore)
    GetEncryptionKey()([]byte)
    GetFileDigest()([]byte)
    GetFileDigestAlgorithm()(*string)
    GetInitializationVector()([]byte)
    GetMac()([]byte)
    GetMacKey()([]byte)
    GetOdataType()(*string)
    GetProfileIdentifier()(*string)
    SetBackingStore(value ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackingStore)()
    SetEncryptionKey(value []byte)()
    SetFileDigest(value []byte)()
    SetFileDigestAlgorithm(value *string)()
    SetInitializationVector(value []byte)()
    SetMac(value []byte)()
    SetMacKey(value []byte)()
    SetOdataType(value *string)()
    SetProfileIdentifier(value *string)()
}
