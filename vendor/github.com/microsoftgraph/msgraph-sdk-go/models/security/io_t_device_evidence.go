package security

import (
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

type IoTDeviceEvidence struct {
    AlertEvidence
}
// NewIoTDeviceEvidence instantiates a new IoTDeviceEvidence and sets the default values.
func NewIoTDeviceEvidence()(*IoTDeviceEvidence) {
    m := &IoTDeviceEvidence{
        AlertEvidence: *NewAlertEvidence(),
    }
    odataTypeValue := "#microsoft.graph.security.ioTDeviceEvidence"
    m.SetOdataType(&odataTypeValue)
    return m
}
// CreateIoTDeviceEvidenceFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
// returns a Parsable when successful
func CreateIoTDeviceEvidenceFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewIoTDeviceEvidence(), nil
}
// GetDeviceId gets the deviceId property value. The deviceId property
// returns a *string when successful
func (m *IoTDeviceEvidence) GetDeviceId()(*string) {
    val, err := m.GetBackingStore().Get("deviceId")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetDeviceName gets the deviceName property value. The deviceName property
// returns a *string when successful
func (m *IoTDeviceEvidence) GetDeviceName()(*string) {
    val, err := m.GetBackingStore().Get("deviceName")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetDevicePageLink gets the devicePageLink property value. The devicePageLink property
// returns a *string when successful
func (m *IoTDeviceEvidence) GetDevicePageLink()(*string) {
    val, err := m.GetBackingStore().Get("devicePageLink")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetDeviceSubType gets the deviceSubType property value. The deviceSubType property
// returns a *string when successful
func (m *IoTDeviceEvidence) GetDeviceSubType()(*string) {
    val, err := m.GetBackingStore().Get("deviceSubType")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetDeviceType gets the deviceType property value. The deviceType property
// returns a *string when successful
func (m *IoTDeviceEvidence) GetDeviceType()(*string) {
    val, err := m.GetBackingStore().Get("deviceType")
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
func (m *IoTDeviceEvidence) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := m.AlertEvidence.GetFieldDeserializers()
    res["deviceId"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetDeviceId(val)
        }
        return nil
    }
    res["deviceName"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetDeviceName(val)
        }
        return nil
    }
    res["devicePageLink"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetDevicePageLink(val)
        }
        return nil
    }
    res["deviceSubType"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetDeviceSubType(val)
        }
        return nil
    }
    res["deviceType"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetDeviceType(val)
        }
        return nil
    }
    res["importance"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetEnumValue(ParseIoTDeviceImportanceType)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetImportance(val.(*IoTDeviceImportanceType))
        }
        return nil
    }
    res["ioTHub"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetObjectValue(CreateAzureResourceEvidenceFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetIoTHub(val.(AzureResourceEvidenceable))
        }
        return nil
    }
    res["ioTSecurityAgentId"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetIoTSecurityAgentId(val)
        }
        return nil
    }
    res["ipAddress"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetObjectValue(CreateIpEvidenceFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetIpAddress(val.(IpEvidenceable))
        }
        return nil
    }
    res["isAuthorized"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetBoolValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetIsAuthorized(val)
        }
        return nil
    }
    res["isProgramming"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetBoolValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetIsProgramming(val)
        }
        return nil
    }
    res["isScanner"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetBoolValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetIsScanner(val)
        }
        return nil
    }
    res["macAddress"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetMacAddress(val)
        }
        return nil
    }
    res["manufacturer"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetManufacturer(val)
        }
        return nil
    }
    res["model"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetModel(val)
        }
        return nil
    }
    res["nics"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetCollectionOfObjectValues(CreateNicEvidenceFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            res := make([]NicEvidenceable, len(val))
            for i, v := range val {
                if v != nil {
                    res[i] = v.(NicEvidenceable)
                }
            }
            m.SetNics(res)
        }
        return nil
    }
    res["operatingSystem"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetOperatingSystem(val)
        }
        return nil
    }
    res["owners"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
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
            m.SetOwners(res)
        }
        return nil
    }
    res["protocols"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
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
            m.SetProtocols(res)
        }
        return nil
    }
    res["purdueLayer"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetPurdueLayer(val)
        }
        return nil
    }
    res["sensor"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetSensor(val)
        }
        return nil
    }
    res["serialNumber"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetSerialNumber(val)
        }
        return nil
    }
    res["site"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetSite(val)
        }
        return nil
    }
    res["source"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetSource(val)
        }
        return nil
    }
    res["sourceRef"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetObjectValue(CreateUrlEvidenceFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetSourceRef(val.(UrlEvidenceable))
        }
        return nil
    }
    res["zone"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetZone(val)
        }
        return nil
    }
    return res
}
// GetImportance gets the importance property value. The importance property
// returns a *IoTDeviceImportanceType when successful
func (m *IoTDeviceEvidence) GetImportance()(*IoTDeviceImportanceType) {
    val, err := m.GetBackingStore().Get("importance")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*IoTDeviceImportanceType)
    }
    return nil
}
// GetIoTHub gets the ioTHub property value. The ioTHub property
// returns a AzureResourceEvidenceable when successful
func (m *IoTDeviceEvidence) GetIoTHub()(AzureResourceEvidenceable) {
    val, err := m.GetBackingStore().Get("ioTHub")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(AzureResourceEvidenceable)
    }
    return nil
}
// GetIoTSecurityAgentId gets the ioTSecurityAgentId property value. The ioTSecurityAgentId property
// returns a *string when successful
func (m *IoTDeviceEvidence) GetIoTSecurityAgentId()(*string) {
    val, err := m.GetBackingStore().Get("ioTSecurityAgentId")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetIpAddress gets the ipAddress property value. The ipAddress property
// returns a IpEvidenceable when successful
func (m *IoTDeviceEvidence) GetIpAddress()(IpEvidenceable) {
    val, err := m.GetBackingStore().Get("ipAddress")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(IpEvidenceable)
    }
    return nil
}
// GetIsAuthorized gets the isAuthorized property value. The isAuthorized property
// returns a *bool when successful
func (m *IoTDeviceEvidence) GetIsAuthorized()(*bool) {
    val, err := m.GetBackingStore().Get("isAuthorized")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*bool)
    }
    return nil
}
// GetIsProgramming gets the isProgramming property value. The isProgramming property
// returns a *bool when successful
func (m *IoTDeviceEvidence) GetIsProgramming()(*bool) {
    val, err := m.GetBackingStore().Get("isProgramming")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*bool)
    }
    return nil
}
// GetIsScanner gets the isScanner property value. The isScanner property
// returns a *bool when successful
func (m *IoTDeviceEvidence) GetIsScanner()(*bool) {
    val, err := m.GetBackingStore().Get("isScanner")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*bool)
    }
    return nil
}
// GetMacAddress gets the macAddress property value. The macAddress property
// returns a *string when successful
func (m *IoTDeviceEvidence) GetMacAddress()(*string) {
    val, err := m.GetBackingStore().Get("macAddress")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetManufacturer gets the manufacturer property value. The manufacturer property
// returns a *string when successful
func (m *IoTDeviceEvidence) GetManufacturer()(*string) {
    val, err := m.GetBackingStore().Get("manufacturer")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetModel gets the model property value. The model property
// returns a *string when successful
func (m *IoTDeviceEvidence) GetModel()(*string) {
    val, err := m.GetBackingStore().Get("model")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetNics gets the nics property value. The nics property
// returns a []NicEvidenceable when successful
func (m *IoTDeviceEvidence) GetNics()([]NicEvidenceable) {
    val, err := m.GetBackingStore().Get("nics")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]NicEvidenceable)
    }
    return nil
}
// GetOperatingSystem gets the operatingSystem property value. The operatingSystem property
// returns a *string when successful
func (m *IoTDeviceEvidence) GetOperatingSystem()(*string) {
    val, err := m.GetBackingStore().Get("operatingSystem")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetOwners gets the owners property value. The owners property
// returns a []string when successful
func (m *IoTDeviceEvidence) GetOwners()([]string) {
    val, err := m.GetBackingStore().Get("owners")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]string)
    }
    return nil
}
// GetProtocols gets the protocols property value. The protocols property
// returns a []string when successful
func (m *IoTDeviceEvidence) GetProtocols()([]string) {
    val, err := m.GetBackingStore().Get("protocols")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]string)
    }
    return nil
}
// GetPurdueLayer gets the purdueLayer property value. The purdueLayer property
// returns a *string when successful
func (m *IoTDeviceEvidence) GetPurdueLayer()(*string) {
    val, err := m.GetBackingStore().Get("purdueLayer")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetSensor gets the sensor property value. The sensor property
// returns a *string when successful
func (m *IoTDeviceEvidence) GetSensor()(*string) {
    val, err := m.GetBackingStore().Get("sensor")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetSerialNumber gets the serialNumber property value. The serialNumber property
// returns a *string when successful
func (m *IoTDeviceEvidence) GetSerialNumber()(*string) {
    val, err := m.GetBackingStore().Get("serialNumber")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetSite gets the site property value. The site property
// returns a *string when successful
func (m *IoTDeviceEvidence) GetSite()(*string) {
    val, err := m.GetBackingStore().Get("site")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetSource gets the source property value. The source property
// returns a *string when successful
func (m *IoTDeviceEvidence) GetSource()(*string) {
    val, err := m.GetBackingStore().Get("source")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetSourceRef gets the sourceRef property value. The sourceRef property
// returns a UrlEvidenceable when successful
func (m *IoTDeviceEvidence) GetSourceRef()(UrlEvidenceable) {
    val, err := m.GetBackingStore().Get("sourceRef")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(UrlEvidenceable)
    }
    return nil
}
// GetZone gets the zone property value. The zone property
// returns a *string when successful
func (m *IoTDeviceEvidence) GetZone()(*string) {
    val, err := m.GetBackingStore().Get("zone")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// Serialize serializes information the current object
func (m *IoTDeviceEvidence) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    err := m.AlertEvidence.Serialize(writer)
    if err != nil {
        return err
    }
    {
        err = writer.WriteStringValue("deviceId", m.GetDeviceId())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteStringValue("deviceName", m.GetDeviceName())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteStringValue("devicePageLink", m.GetDevicePageLink())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteStringValue("deviceSubType", m.GetDeviceSubType())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteStringValue("deviceType", m.GetDeviceType())
        if err != nil {
            return err
        }
    }
    if m.GetImportance() != nil {
        cast := (*m.GetImportance()).String()
        err = writer.WriteStringValue("importance", &cast)
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteObjectValue("ioTHub", m.GetIoTHub())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteStringValue("ioTSecurityAgentId", m.GetIoTSecurityAgentId())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteObjectValue("ipAddress", m.GetIpAddress())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteBoolValue("isAuthorized", m.GetIsAuthorized())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteBoolValue("isProgramming", m.GetIsProgramming())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteBoolValue("isScanner", m.GetIsScanner())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteStringValue("macAddress", m.GetMacAddress())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteStringValue("manufacturer", m.GetManufacturer())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteStringValue("model", m.GetModel())
        if err != nil {
            return err
        }
    }
    if m.GetNics() != nil {
        cast := make([]i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, len(m.GetNics()))
        for i, v := range m.GetNics() {
            if v != nil {
                cast[i] = v.(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable)
            }
        }
        err = writer.WriteCollectionOfObjectValues("nics", cast)
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteStringValue("operatingSystem", m.GetOperatingSystem())
        if err != nil {
            return err
        }
    }
    if m.GetOwners() != nil {
        err = writer.WriteCollectionOfStringValues("owners", m.GetOwners())
        if err != nil {
            return err
        }
    }
    if m.GetProtocols() != nil {
        err = writer.WriteCollectionOfStringValues("protocols", m.GetProtocols())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteStringValue("purdueLayer", m.GetPurdueLayer())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteStringValue("sensor", m.GetSensor())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteStringValue("serialNumber", m.GetSerialNumber())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteStringValue("site", m.GetSite())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteStringValue("source", m.GetSource())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteObjectValue("sourceRef", m.GetSourceRef())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteStringValue("zone", m.GetZone())
        if err != nil {
            return err
        }
    }
    return nil
}
// SetDeviceId sets the deviceId property value. The deviceId property
func (m *IoTDeviceEvidence) SetDeviceId(value *string)() {
    err := m.GetBackingStore().Set("deviceId", value)
    if err != nil {
        panic(err)
    }
}
// SetDeviceName sets the deviceName property value. The deviceName property
func (m *IoTDeviceEvidence) SetDeviceName(value *string)() {
    err := m.GetBackingStore().Set("deviceName", value)
    if err != nil {
        panic(err)
    }
}
// SetDevicePageLink sets the devicePageLink property value. The devicePageLink property
func (m *IoTDeviceEvidence) SetDevicePageLink(value *string)() {
    err := m.GetBackingStore().Set("devicePageLink", value)
    if err != nil {
        panic(err)
    }
}
// SetDeviceSubType sets the deviceSubType property value. The deviceSubType property
func (m *IoTDeviceEvidence) SetDeviceSubType(value *string)() {
    err := m.GetBackingStore().Set("deviceSubType", value)
    if err != nil {
        panic(err)
    }
}
// SetDeviceType sets the deviceType property value. The deviceType property
func (m *IoTDeviceEvidence) SetDeviceType(value *string)() {
    err := m.GetBackingStore().Set("deviceType", value)
    if err != nil {
        panic(err)
    }
}
// SetImportance sets the importance property value. The importance property
func (m *IoTDeviceEvidence) SetImportance(value *IoTDeviceImportanceType)() {
    err := m.GetBackingStore().Set("importance", value)
    if err != nil {
        panic(err)
    }
}
// SetIoTHub sets the ioTHub property value. The ioTHub property
func (m *IoTDeviceEvidence) SetIoTHub(value AzureResourceEvidenceable)() {
    err := m.GetBackingStore().Set("ioTHub", value)
    if err != nil {
        panic(err)
    }
}
// SetIoTSecurityAgentId sets the ioTSecurityAgentId property value. The ioTSecurityAgentId property
func (m *IoTDeviceEvidence) SetIoTSecurityAgentId(value *string)() {
    err := m.GetBackingStore().Set("ioTSecurityAgentId", value)
    if err != nil {
        panic(err)
    }
}
// SetIpAddress sets the ipAddress property value. The ipAddress property
func (m *IoTDeviceEvidence) SetIpAddress(value IpEvidenceable)() {
    err := m.GetBackingStore().Set("ipAddress", value)
    if err != nil {
        panic(err)
    }
}
// SetIsAuthorized sets the isAuthorized property value. The isAuthorized property
func (m *IoTDeviceEvidence) SetIsAuthorized(value *bool)() {
    err := m.GetBackingStore().Set("isAuthorized", value)
    if err != nil {
        panic(err)
    }
}
// SetIsProgramming sets the isProgramming property value. The isProgramming property
func (m *IoTDeviceEvidence) SetIsProgramming(value *bool)() {
    err := m.GetBackingStore().Set("isProgramming", value)
    if err != nil {
        panic(err)
    }
}
// SetIsScanner sets the isScanner property value. The isScanner property
func (m *IoTDeviceEvidence) SetIsScanner(value *bool)() {
    err := m.GetBackingStore().Set("isScanner", value)
    if err != nil {
        panic(err)
    }
}
// SetMacAddress sets the macAddress property value. The macAddress property
func (m *IoTDeviceEvidence) SetMacAddress(value *string)() {
    err := m.GetBackingStore().Set("macAddress", value)
    if err != nil {
        panic(err)
    }
}
// SetManufacturer sets the manufacturer property value. The manufacturer property
func (m *IoTDeviceEvidence) SetManufacturer(value *string)() {
    err := m.GetBackingStore().Set("manufacturer", value)
    if err != nil {
        panic(err)
    }
}
// SetModel sets the model property value. The model property
func (m *IoTDeviceEvidence) SetModel(value *string)() {
    err := m.GetBackingStore().Set("model", value)
    if err != nil {
        panic(err)
    }
}
// SetNics sets the nics property value. The nics property
func (m *IoTDeviceEvidence) SetNics(value []NicEvidenceable)() {
    err := m.GetBackingStore().Set("nics", value)
    if err != nil {
        panic(err)
    }
}
// SetOperatingSystem sets the operatingSystem property value. The operatingSystem property
func (m *IoTDeviceEvidence) SetOperatingSystem(value *string)() {
    err := m.GetBackingStore().Set("operatingSystem", value)
    if err != nil {
        panic(err)
    }
}
// SetOwners sets the owners property value. The owners property
func (m *IoTDeviceEvidence) SetOwners(value []string)() {
    err := m.GetBackingStore().Set("owners", value)
    if err != nil {
        panic(err)
    }
}
// SetProtocols sets the protocols property value. The protocols property
func (m *IoTDeviceEvidence) SetProtocols(value []string)() {
    err := m.GetBackingStore().Set("protocols", value)
    if err != nil {
        panic(err)
    }
}
// SetPurdueLayer sets the purdueLayer property value. The purdueLayer property
func (m *IoTDeviceEvidence) SetPurdueLayer(value *string)() {
    err := m.GetBackingStore().Set("purdueLayer", value)
    if err != nil {
        panic(err)
    }
}
// SetSensor sets the sensor property value. The sensor property
func (m *IoTDeviceEvidence) SetSensor(value *string)() {
    err := m.GetBackingStore().Set("sensor", value)
    if err != nil {
        panic(err)
    }
}
// SetSerialNumber sets the serialNumber property value. The serialNumber property
func (m *IoTDeviceEvidence) SetSerialNumber(value *string)() {
    err := m.GetBackingStore().Set("serialNumber", value)
    if err != nil {
        panic(err)
    }
}
// SetSite sets the site property value. The site property
func (m *IoTDeviceEvidence) SetSite(value *string)() {
    err := m.GetBackingStore().Set("site", value)
    if err != nil {
        panic(err)
    }
}
// SetSource sets the source property value. The source property
func (m *IoTDeviceEvidence) SetSource(value *string)() {
    err := m.GetBackingStore().Set("source", value)
    if err != nil {
        panic(err)
    }
}
// SetSourceRef sets the sourceRef property value. The sourceRef property
func (m *IoTDeviceEvidence) SetSourceRef(value UrlEvidenceable)() {
    err := m.GetBackingStore().Set("sourceRef", value)
    if err != nil {
        panic(err)
    }
}
// SetZone sets the zone property value. The zone property
func (m *IoTDeviceEvidence) SetZone(value *string)() {
    err := m.GetBackingStore().Set("zone", value)
    if err != nil {
        panic(err)
    }
}
type IoTDeviceEvidenceable interface {
    AlertEvidenceable
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable
    GetDeviceId()(*string)
    GetDeviceName()(*string)
    GetDevicePageLink()(*string)
    GetDeviceSubType()(*string)
    GetDeviceType()(*string)
    GetImportance()(*IoTDeviceImportanceType)
    GetIoTHub()(AzureResourceEvidenceable)
    GetIoTSecurityAgentId()(*string)
    GetIpAddress()(IpEvidenceable)
    GetIsAuthorized()(*bool)
    GetIsProgramming()(*bool)
    GetIsScanner()(*bool)
    GetMacAddress()(*string)
    GetManufacturer()(*string)
    GetModel()(*string)
    GetNics()([]NicEvidenceable)
    GetOperatingSystem()(*string)
    GetOwners()([]string)
    GetProtocols()([]string)
    GetPurdueLayer()(*string)
    GetSensor()(*string)
    GetSerialNumber()(*string)
    GetSite()(*string)
    GetSource()(*string)
    GetSourceRef()(UrlEvidenceable)
    GetZone()(*string)
    SetDeviceId(value *string)()
    SetDeviceName(value *string)()
    SetDevicePageLink(value *string)()
    SetDeviceSubType(value *string)()
    SetDeviceType(value *string)()
    SetImportance(value *IoTDeviceImportanceType)()
    SetIoTHub(value AzureResourceEvidenceable)()
    SetIoTSecurityAgentId(value *string)()
    SetIpAddress(value IpEvidenceable)()
    SetIsAuthorized(value *bool)()
    SetIsProgramming(value *bool)()
    SetIsScanner(value *bool)()
    SetMacAddress(value *string)()
    SetManufacturer(value *string)()
    SetModel(value *string)()
    SetNics(value []NicEvidenceable)()
    SetOperatingSystem(value *string)()
    SetOwners(value []string)()
    SetProtocols(value []string)()
    SetPurdueLayer(value *string)()
    SetSensor(value *string)()
    SetSerialNumber(value *string)()
    SetSite(value *string)()
    SetSource(value *string)()
    SetSourceRef(value UrlEvidenceable)()
    SetZone(value *string)()
}
