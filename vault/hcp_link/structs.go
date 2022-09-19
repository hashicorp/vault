package hcp_link

// SessionStatus is used to express the current status of the SCADA session.
type SessionStatus = string

const (
	// Connected HCP link connection status when it is connected
	Connected = SessionStatus("connected")
	// Disconnected HCP link connection status when it is disconnected
	Disconnected = SessionStatus("disconnected")
)

type WrappedHCPLinkVault struct {
	HCPLinkVaultInterface
}

type HCPLinkVaultInterface interface {
	Shutdown() error
	GetScadaSessionStatus() string
	GetConnectionStatusMessage(string) string
}
