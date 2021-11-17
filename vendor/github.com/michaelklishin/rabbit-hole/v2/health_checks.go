package rabbithole

import (
	"encoding/json"
	"net/http"
	"strconv"
)

type TimeUnit string

const (
	SECONDS TimeUnit = "seconds"
	DAYS    TimeUnit = "days"
	MONTHS  TimeUnit = "months"
	YEARS   TimeUnit = "years"
)

type Protocol string

const (
	AMQP    Protocol = "amqp"
	AMQPS   Protocol = "amqp/ssl"
	AMQP091 Protocol = "amqp091"
	AMQP10  Protocol = "amqp10"

	MQTT       Protocol = "mqtt"
	STOMP      Protocol = "stomp"
	WebMQTT    Protocol = "web-mqtt"
	WebSTOMP   Protocol = "web-stomp"
	HTTP       Protocol = "http"
	HTTPS      Protocol = "https"
	Prometheus Protocol = "http/prometheus"
	Clustering Protocol = "clustering"
)

// HealthCheck represents a generic health check endpoint response
// Related RabbitMQ doc guide: https://www.rabbitmq.com/monitoring.html
type HealthCheck interface {
	// Returns true if the check is ok, otherwise false
	Ok() bool
}

// HealthCheckStatus represents a generic health check endpoint response
// Related RabbitMQ doc guide: https://www.rabbitmq.com/monitoring.html
type HealthCheckStatus struct {
	HealthCheck
	Status string `json:"status"`
	Reason string `json:"reason,omitempty"`
}

// Ok returns true if the health check succeeded
func (h *HealthCheckStatus) Ok() bool {
	return h.Status == "ok"
}

// AlarmInEffect represents a resource alarm in effect on a node
type AlarmInEffect struct {
	Node     string `json:"node"`
	Resource string `json:"resource"`
}

// ResourceAlarmCheckStatus represents the response from HealthCheckALarms
type ResourceAlarmCheckStatus struct {
	HealthCheck
	Status string          `json:"status"`
	Reason string          `json:"reason,omitempty"`
	Alarms []AlarmInEffect `json:"alarms,omitempty"`
}

// Ok returns true if the health check succeeded
func (h *ResourceAlarmCheckStatus) Ok() bool {
	return h.Status == "ok"
}

// HealthCheckAlarms checks if there are resource alarms in effect in the cluster
// Related RabbitMQ doc guide: https://www.rabbitmq.com/alarms.html
func (c *Client) HealthCheckAlarms() (rec ResourceAlarmCheckStatus, err error) {
	err = c.executeCheck("health/checks/alarms", &rec)
	return rec, err
}

// HealthCheckLocalAlarms checks if there are resource alarms in effect on the target node
// Related RabbitMQ doc guide: https://www.rabbitmq.com/alarms.html
func (c *Client) HealthCheckLocalAlarms() (rec ResourceAlarmCheckStatus, err error) {
	err = c.executeCheck("health/checks/local-alarms", &rec)
	return rec, err
}

// HealthCheckCertificateExpiration checks the expiration date on the certificates for every listener configured to use TLS.
// Valid units: days, weeks, months, years. The value of the within argument is the number of units.
// So, when within is 2 and unit is "months", the expiration period used by the check will be the next two months.
func (c *Client) HealthCheckCertificateExpiration(within uint, unit TimeUnit) (rec HealthCheckStatus, err error) {
	err = c.executeCheck("health/checks/certificate-expiration/"+strconv.Itoa(int(within))+"/"+string(unit), &rec)
	return rec, err
}

// PortListenerCheckStatus represents the response from HealthCheckPortListener
type PortListenerCheckStatus struct {
	HealthCheck
	Status  string `json:"status"`
	Reason  string `json:"reason,omitempty"`
	Port    uint   `json:"port,omitempty"`
	Missing uint   `json:"missing,omitempty"`
	Ports   []uint `json:"ports,omitempty"`
}

// Ok returns true if the health check succeeded
func (h *PortListenerCheckStatus) Ok() bool {
	return h.Status == "ok"
}

// HealthCheckPortListener checks if there is an active listener on the give port.
// Relevant RabbitMQ doc guide: https://www.rabbitmq.com/monitoring.html
func (c *Client) HealthCheckPortListener(port uint) (rec PortListenerCheckStatus, err error) {
	err = c.executeCheck("health/checks/port-listener/"+strconv.Itoa(int(port)), &rec)
	return rec, err
}

// ProtocolListenerCheckStatus represents the response from HealthCheckProtocolListener
type ProtocolListenerCheckStatus struct {
	HealthCheck
	Status    string   `json:"status"`
	Reason    string   `json:"reason,omitempty"`
	Missing   string   `json:"missing,omitempty"`
	Protocols []string `json:"protocols,omitempty"`
}

// Ok returns true if the health check succeeded
func (h *ProtocolListenerCheckStatus) Ok() bool {
	return h.Status == "ok"
}

// HealthCheckProtocolListener checks if there is an active listener for the given protocol
// Valid protocol names are: amqp091, amqp10, mqtt, stomp, web-mqtt, web-stomp, http, https, clustering
// Relevant RabbitMQ doc guide: https://www.rabbitmq.com/monitoring.html
func (c *Client) HealthCheckProtocolListener(protocol Protocol) (rec ProtocolListenerCheckStatus, err error) {
	err = c.executeCheck("health/checks/protocol-listener/"+string(protocol), &rec)
	return rec, err
}

// HealthCheckVirtualHosts checks if all virtual hosts are running on the target node
func (c *Client) HealthCheckVirtualHosts() (rec HealthCheckStatus, err error) {
	err = c.executeCheck("health/checks/virtual-hosts", &rec)
	return rec, err
}

// HealthCheckNodeIsMirrorSyncCritical checks if there are classic mirrored queues without synchronised mirrors online
// (queues that would potentially lose data if the target node is shut down).
func (c *Client) HealthCheckNodeIsMirrorSyncCritical() (rec HealthCheckStatus, err error) {
	err = c.executeCheck("health/checks/node-is-mirror-sync-critical", &rec)
	return rec, err
}

// HealthCheckNodeIsQuorumCritical checks if there are quorum queues with minimum online quorum (queues that would lose
// their quorum and availability if the target node is shut down).
// Relevant RabbitMQ doc guide: https://www.rabbitmq.com/quorum-queues.html
func (c *Client) HealthCheckNodeIsQuorumCritical() (rec HealthCheckStatus, err error) {
	err = c.executeCheck("health/checks/node-is-quorum-critical", &rec)
	return rec, err
}

func (c *Client) executeCheck(path string, rec interface{}) error {
	req, err := newGETRequest(c, path)
	httpc := &http.Client{
		Timeout: c.timeout,
	}
	if c.transport != nil {
		httpc.Transport = c.transport
	}
	resp, err := httpc.Do(req)
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	if resp.StatusCode < http.StatusBadRequest || resp.StatusCode == http.StatusServiceUnavailable {
		if err = json.NewDecoder(resp.Body).Decode(&rec); err != nil {
			return err
		}

		return nil
	}

	if err = parseResponseErrors(resp); err != nil {
		return err
	}

	return nil
}
