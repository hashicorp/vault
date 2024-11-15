package billing

import (
    i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e "time"
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
    iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242 "github.com/microsoftgraph/msgraph-sdk-go/models"
)

type Manifest struct {
    iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.Entity
}
// NewManifest instantiates a new Manifest and sets the default values.
func NewManifest()(*Manifest) {
    m := &Manifest{
        Entity: *iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.NewEntity(),
    }
    return m
}
// CreateManifestFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
// returns a Parsable when successful
func CreateManifestFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewManifest(), nil
}
// GetBlobCount gets the blobCount property value. The total file count for this partner tenant ID.
// returns a *int32 when successful
func (m *Manifest) GetBlobCount()(*int32) {
    val, err := m.GetBackingStore().Get("blobCount")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*int32)
    }
    return nil
}
// GetBlobs gets the blobs property value. A collection of blob objects that contain details of all the files for the partner tenant ID.
// returns a []Blobable when successful
func (m *Manifest) GetBlobs()([]Blobable) {
    val, err := m.GetBackingStore().Get("blobs")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]Blobable)
    }
    return nil
}
// GetCreatedDateTime gets the createdDateTime property value. The date and time when a manifest resource was created. The timestamp type represents date and time information using ISO 8601 format and is always in UTC. For example, midnight UTC on Jan 1, 2014 is 2014-01-01T00:00:00Z.
// returns a *Time when successful
func (m *Manifest) GetCreatedDateTime()(*i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time) {
    val, err := m.GetBackingStore().Get("createdDateTime")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)
    }
    return nil
}
// GetDataFormat gets the dataFormat property value. The billing data file format. The possible value is: compressedJSONLines. Each blob is a compressed file and data in the file is in JSON lines format. Decompress the file to access the data.
// returns a *string when successful
func (m *Manifest) GetDataFormat()(*string) {
    val, err := m.GetBackingStore().Get("dataFormat")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetETag gets the eTag property value. Version of data represented by the manifest. Any change in eTag indicates a new data version.
// returns a *string when successful
func (m *Manifest) GetETag()(*string) {
    val, err := m.GetBackingStore().Get("eTag")
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
func (m *Manifest) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := m.Entity.GetFieldDeserializers()
    res["blobCount"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetInt32Value()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetBlobCount(val)
        }
        return nil
    }
    res["blobs"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetCollectionOfObjectValues(CreateBlobFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            res := make([]Blobable, len(val))
            for i, v := range val {
                if v != nil {
                    res[i] = v.(Blobable)
                }
            }
            m.SetBlobs(res)
        }
        return nil
    }
    res["createdDateTime"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetTimeValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetCreatedDateTime(val)
        }
        return nil
    }
    res["dataFormat"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetDataFormat(val)
        }
        return nil
    }
    res["eTag"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetETag(val)
        }
        return nil
    }
    res["partitionType"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetPartitionType(val)
        }
        return nil
    }
    res["partnerTenantId"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetPartnerTenantId(val)
        }
        return nil
    }
    res["rootDirectory"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetRootDirectory(val)
        }
        return nil
    }
    res["sasToken"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetSasToken(val)
        }
        return nil
    }
    res["schemaVersion"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetSchemaVersion(val)
        }
        return nil
    }
    return res
}
// GetPartitionType gets the partitionType property value. Indicates the division of data. If a given partition has more than the supported number, the data is split into multiple files, each file representing a specific partitionValue. By default, the data in the file is partitioned by the number of line items.
// returns a *string when successful
func (m *Manifest) GetPartitionType()(*string) {
    val, err := m.GetBackingStore().Get("partitionType")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetPartnerTenantId gets the partnerTenantId property value. The Microsoft Entra tenant ID of the partner.
// returns a *string when successful
func (m *Manifest) GetPartnerTenantId()(*string) {
    val, err := m.GetBackingStore().Get("partnerTenantId")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetRootDirectory gets the rootDirectory property value. The root directory that contains all the files.
// returns a *string when successful
func (m *Manifest) GetRootDirectory()(*string) {
    val, err := m.GetBackingStore().Get("rootDirectory")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetSasToken gets the sasToken property value. The SAS token for accessing the directory or an individual file in the directory.
// returns a *string when successful
func (m *Manifest) GetSasToken()(*string) {
    val, err := m.GetBackingStore().Get("sasToken")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetSchemaVersion gets the schemaVersion property value. The version of the manifest schema.
// returns a *string when successful
func (m *Manifest) GetSchemaVersion()(*string) {
    val, err := m.GetBackingStore().Get("schemaVersion")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// Serialize serializes information the current object
func (m *Manifest) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    err := m.Entity.Serialize(writer)
    if err != nil {
        return err
    }
    {
        err = writer.WriteInt32Value("blobCount", m.GetBlobCount())
        if err != nil {
            return err
        }
    }
    if m.GetBlobs() != nil {
        cast := make([]i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, len(m.GetBlobs()))
        for i, v := range m.GetBlobs() {
            if v != nil {
                cast[i] = v.(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable)
            }
        }
        err = writer.WriteCollectionOfObjectValues("blobs", cast)
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteTimeValue("createdDateTime", m.GetCreatedDateTime())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteStringValue("dataFormat", m.GetDataFormat())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteStringValue("eTag", m.GetETag())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteStringValue("partitionType", m.GetPartitionType())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteStringValue("partnerTenantId", m.GetPartnerTenantId())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteStringValue("rootDirectory", m.GetRootDirectory())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteStringValue("sasToken", m.GetSasToken())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteStringValue("schemaVersion", m.GetSchemaVersion())
        if err != nil {
            return err
        }
    }
    return nil
}
// SetBlobCount sets the blobCount property value. The total file count for this partner tenant ID.
func (m *Manifest) SetBlobCount(value *int32)() {
    err := m.GetBackingStore().Set("blobCount", value)
    if err != nil {
        panic(err)
    }
}
// SetBlobs sets the blobs property value. A collection of blob objects that contain details of all the files for the partner tenant ID.
func (m *Manifest) SetBlobs(value []Blobable)() {
    err := m.GetBackingStore().Set("blobs", value)
    if err != nil {
        panic(err)
    }
}
// SetCreatedDateTime sets the createdDateTime property value. The date and time when a manifest resource was created. The timestamp type represents date and time information using ISO 8601 format and is always in UTC. For example, midnight UTC on Jan 1, 2014 is 2014-01-01T00:00:00Z.
func (m *Manifest) SetCreatedDateTime(value *i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)() {
    err := m.GetBackingStore().Set("createdDateTime", value)
    if err != nil {
        panic(err)
    }
}
// SetDataFormat sets the dataFormat property value. The billing data file format. The possible value is: compressedJSONLines. Each blob is a compressed file and data in the file is in JSON lines format. Decompress the file to access the data.
func (m *Manifest) SetDataFormat(value *string)() {
    err := m.GetBackingStore().Set("dataFormat", value)
    if err != nil {
        panic(err)
    }
}
// SetETag sets the eTag property value. Version of data represented by the manifest. Any change in eTag indicates a new data version.
func (m *Manifest) SetETag(value *string)() {
    err := m.GetBackingStore().Set("eTag", value)
    if err != nil {
        panic(err)
    }
}
// SetPartitionType sets the partitionType property value. Indicates the division of data. If a given partition has more than the supported number, the data is split into multiple files, each file representing a specific partitionValue. By default, the data in the file is partitioned by the number of line items.
func (m *Manifest) SetPartitionType(value *string)() {
    err := m.GetBackingStore().Set("partitionType", value)
    if err != nil {
        panic(err)
    }
}
// SetPartnerTenantId sets the partnerTenantId property value. The Microsoft Entra tenant ID of the partner.
func (m *Manifest) SetPartnerTenantId(value *string)() {
    err := m.GetBackingStore().Set("partnerTenantId", value)
    if err != nil {
        panic(err)
    }
}
// SetRootDirectory sets the rootDirectory property value. The root directory that contains all the files.
func (m *Manifest) SetRootDirectory(value *string)() {
    err := m.GetBackingStore().Set("rootDirectory", value)
    if err != nil {
        panic(err)
    }
}
// SetSasToken sets the sasToken property value. The SAS token for accessing the directory or an individual file in the directory.
func (m *Manifest) SetSasToken(value *string)() {
    err := m.GetBackingStore().Set("sasToken", value)
    if err != nil {
        panic(err)
    }
}
// SetSchemaVersion sets the schemaVersion property value. The version of the manifest schema.
func (m *Manifest) SetSchemaVersion(value *string)() {
    err := m.GetBackingStore().Set("schemaVersion", value)
    if err != nil {
        panic(err)
    }
}
type Manifestable interface {
    iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.Entityable
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable
    GetBlobCount()(*int32)
    GetBlobs()([]Blobable)
    GetCreatedDateTime()(*i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)
    GetDataFormat()(*string)
    GetETag()(*string)
    GetPartitionType()(*string)
    GetPartnerTenantId()(*string)
    GetRootDirectory()(*string)
    GetSasToken()(*string)
    GetSchemaVersion()(*string)
    SetBlobCount(value *int32)()
    SetBlobs(value []Blobable)()
    SetCreatedDateTime(value *i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)()
    SetDataFormat(value *string)()
    SetETag(value *string)()
    SetPartitionType(value *string)()
    SetPartnerTenantId(value *string)()
    SetRootDirectory(value *string)()
    SetSasToken(value *string)()
    SetSchemaVersion(value *string)()
}
