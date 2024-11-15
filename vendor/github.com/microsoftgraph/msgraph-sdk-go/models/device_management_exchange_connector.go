package models

import (
    i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e "time"
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

// DeviceManagementExchangeConnector entity which represents a connection to an Exchange environment.
type DeviceManagementExchangeConnector struct {
    Entity
}
// NewDeviceManagementExchangeConnector instantiates a new DeviceManagementExchangeConnector and sets the default values.
func NewDeviceManagementExchangeConnector()(*DeviceManagementExchangeConnector) {
    m := &DeviceManagementExchangeConnector{
        Entity: *NewEntity(),
    }
    return m
}
// CreateDeviceManagementExchangeConnectorFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
// returns a Parsable when successful
func CreateDeviceManagementExchangeConnectorFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewDeviceManagementExchangeConnector(), nil
}
// GetConnectorServerName gets the connectorServerName property value. The name of the server hosting the Exchange Connector.
// returns a *string when successful
func (m *DeviceManagementExchangeConnector) GetConnectorServerName()(*string) {
    val, err := m.GetBackingStore().Get("connectorServerName")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetExchangeAlias gets the exchangeAlias property value. An alias assigned to the Exchange server
// returns a *string when successful
func (m *DeviceManagementExchangeConnector) GetExchangeAlias()(*string) {
    val, err := m.GetBackingStore().Get("exchangeAlias")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetExchangeConnectorType gets the exchangeConnectorType property value. The type of Exchange Connector.
// returns a *DeviceManagementExchangeConnectorType when successful
func (m *DeviceManagementExchangeConnector) GetExchangeConnectorType()(*DeviceManagementExchangeConnectorType) {
    val, err := m.GetBackingStore().Get("exchangeConnectorType")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*DeviceManagementExchangeConnectorType)
    }
    return nil
}
// GetExchangeOrganization gets the exchangeOrganization property value. Exchange Organization to the Exchange server
// returns a *string when successful
func (m *DeviceManagementExchangeConnector) GetExchangeOrganization()(*string) {
    val, err := m.GetBackingStore().Get("exchangeOrganization")
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
func (m *DeviceManagementExchangeConnector) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := m.Entity.GetFieldDeserializers()
    res["connectorServerName"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetConnectorServerName(val)
        }
        return nil
    }
    res["exchangeAlias"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetExchangeAlias(val)
        }
        return nil
    }
    res["exchangeConnectorType"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetEnumValue(ParseDeviceManagementExchangeConnectorType)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetExchangeConnectorType(val.(*DeviceManagementExchangeConnectorType))
        }
        return nil
    }
    res["exchangeOrganization"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetExchangeOrganization(val)
        }
        return nil
    }
    res["lastSyncDateTime"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetTimeValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetLastSyncDateTime(val)
        }
        return nil
    }
    res["primarySmtpAddress"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetPrimarySmtpAddress(val)
        }
        return nil
    }
    res["serverName"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetServerName(val)
        }
        return nil
    }
    res["status"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetEnumValue(ParseDeviceManagementExchangeConnectorStatus)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetStatus(val.(*DeviceManagementExchangeConnectorStatus))
        }
        return nil
    }
    res["version"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetVersion(val)
        }
        return nil
    }
    return res
}
// GetLastSyncDateTime gets the lastSyncDateTime property value. Last sync time for the Exchange Connector
// returns a *Time when successful
func (m *DeviceManagementExchangeConnector) GetLastSyncDateTime()(*i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time) {
    val, err := m.GetBackingStore().Get("lastSyncDateTime")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)
    }
    return nil
}
// GetPrimarySmtpAddress gets the primarySmtpAddress property value. Email address used to configure the Service To Service Exchange Connector.
// returns a *string when successful
func (m *DeviceManagementExchangeConnector) GetPrimarySmtpAddress()(*string) {
    val, err := m.GetBackingStore().Get("primarySmtpAddress")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetServerName gets the serverName property value. The name of the Exchange server.
// returns a *string when successful
func (m *DeviceManagementExchangeConnector) GetServerName()(*string) {
    val, err := m.GetBackingStore().Get("serverName")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetStatus gets the status property value. The current status of the Exchange Connector.
// returns a *DeviceManagementExchangeConnectorStatus when successful
func (m *DeviceManagementExchangeConnector) GetStatus()(*DeviceManagementExchangeConnectorStatus) {
    val, err := m.GetBackingStore().Get("status")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*DeviceManagementExchangeConnectorStatus)
    }
    return nil
}
// GetVersion gets the version property value. The version of the ExchangeConnectorAgent
// returns a *string when successful
func (m *DeviceManagementExchangeConnector) GetVersion()(*string) {
    val, err := m.GetBackingStore().Get("version")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// Serialize serializes information the current object
func (m *DeviceManagementExchangeConnector) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    err := m.Entity.Serialize(writer)
    if err != nil {
        return err
    }
    {
        err = writer.WriteStringValue("connectorServerName", m.GetConnectorServerName())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteStringValue("exchangeAlias", m.GetExchangeAlias())
        if err != nil {
            return err
        }
    }
    if m.GetExchangeConnectorType() != nil {
        cast := (*m.GetExchangeConnectorType()).String()
        err = writer.WriteStringValue("exchangeConnectorType", &cast)
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteStringValue("exchangeOrganization", m.GetExchangeOrganization())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteTimeValue("lastSyncDateTime", m.GetLastSyncDateTime())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteStringValue("primarySmtpAddress", m.GetPrimarySmtpAddress())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteStringValue("serverName", m.GetServerName())
        if err != nil {
            return err
        }
    }
    if m.GetStatus() != nil {
        cast := (*m.GetStatus()).String()
        err = writer.WriteStringValue("status", &cast)
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteStringValue("version", m.GetVersion())
        if err != nil {
            return err
        }
    }
    return nil
}
// SetConnectorServerName sets the connectorServerName property value. The name of the server hosting the Exchange Connector.
func (m *DeviceManagementExchangeConnector) SetConnectorServerName(value *string)() {
    err := m.GetBackingStore().Set("connectorServerName", value)
    if err != nil {
        panic(err)
    }
}
// SetExchangeAlias sets the exchangeAlias property value. An alias assigned to the Exchange server
func (m *DeviceManagementExchangeConnector) SetExchangeAlias(value *string)() {
    err := m.GetBackingStore().Set("exchangeAlias", value)
    if err != nil {
        panic(err)
    }
}
// SetExchangeConnectorType sets the exchangeConnectorType property value. The type of Exchange Connector.
func (m *DeviceManagementExchangeConnector) SetExchangeConnectorType(value *DeviceManagementExchangeConnectorType)() {
    err := m.GetBackingStore().Set("exchangeConnectorType", value)
    if err != nil {
        panic(err)
    }
}
// SetExchangeOrganization sets the exchangeOrganization property value. Exchange Organization to the Exchange server
func (m *DeviceManagementExchangeConnector) SetExchangeOrganization(value *string)() {
    err := m.GetBackingStore().Set("exchangeOrganization", value)
    if err != nil {
        panic(err)
    }
}
// SetLastSyncDateTime sets the lastSyncDateTime property value. Last sync time for the Exchange Connector
func (m *DeviceManagementExchangeConnector) SetLastSyncDateTime(value *i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)() {
    err := m.GetBackingStore().Set("lastSyncDateTime", value)
    if err != nil {
        panic(err)
    }
}
// SetPrimarySmtpAddress sets the primarySmtpAddress property value. Email address used to configure the Service To Service Exchange Connector.
func (m *DeviceManagementExchangeConnector) SetPrimarySmtpAddress(value *string)() {
    err := m.GetBackingStore().Set("primarySmtpAddress", value)
    if err != nil {
        panic(err)
    }
}
// SetServerName sets the serverName property value. The name of the Exchange server.
func (m *DeviceManagementExchangeConnector) SetServerName(value *string)() {
    err := m.GetBackingStore().Set("serverName", value)
    if err != nil {
        panic(err)
    }
}
// SetStatus sets the status property value. The current status of the Exchange Connector.
func (m *DeviceManagementExchangeConnector) SetStatus(value *DeviceManagementExchangeConnectorStatus)() {
    err := m.GetBackingStore().Set("status", value)
    if err != nil {
        panic(err)
    }
}
// SetVersion sets the version property value. The version of the ExchangeConnectorAgent
func (m *DeviceManagementExchangeConnector) SetVersion(value *string)() {
    err := m.GetBackingStore().Set("version", value)
    if err != nil {
        panic(err)
    }
}
type DeviceManagementExchangeConnectorable interface {
    Entityable
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable
    GetConnectorServerName()(*string)
    GetExchangeAlias()(*string)
    GetExchangeConnectorType()(*DeviceManagementExchangeConnectorType)
    GetExchangeOrganization()(*string)
    GetLastSyncDateTime()(*i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)
    GetPrimarySmtpAddress()(*string)
    GetServerName()(*string)
    GetStatus()(*DeviceManagementExchangeConnectorStatus)
    GetVersion()(*string)
    SetConnectorServerName(value *string)()
    SetExchangeAlias(value *string)()
    SetExchangeConnectorType(value *DeviceManagementExchangeConnectorType)()
    SetExchangeOrganization(value *string)()
    SetLastSyncDateTime(value *i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)()
    SetPrimarySmtpAddress(value *string)()
    SetServerName(value *string)()
    SetStatus(value *DeviceManagementExchangeConnectorStatus)()
    SetVersion(value *string)()
}
