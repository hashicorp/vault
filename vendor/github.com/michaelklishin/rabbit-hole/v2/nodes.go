package rabbithole

import (
	"net/url"
)

// OsPid is an operating system process ID.
type OsPid string

// NameDescriptionEnabled represents a named entity with a description.
type NameDescriptionEnabled struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Enabled     bool   `json:"enabled"`
}

// AuthMechanism is a RabbbitMQ authentication and/or authorization mechanism
// available on the node.
type AuthMechanism NameDescriptionEnabled

// ExchangeType is an exchange type available on the node.
type ExchangeType NameDescriptionEnabled

// NameDescriptionVersion represents a named, versioned entity.
type NameDescriptionVersion struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Version     string `json:"version"`
}

// ErlangApp is an Erlang application running on a node.
type ErlangApp NameDescriptionVersion

// ClusterLink is an inter-node communications link entity in clustering
type ClusterLink struct {
	Stats     ClusterLinkStats `json:"stats"`
	Name      string           `json:"name"`
	PeerAddr  string           `json:"peer_addr"`
	PeerPort  uint             `json:"peer_addr"`
	SockAddr  string           `json:"sock_addr"`
	SockPort  uint             `json:"sock_addr"`
	SendBytes uint64           `json:"send_bytes"`
	RecvBytes uint64           `json:"recv_bytes"`
}

// ClusterLinkStats is a stats field in ClusterLink
type ClusterLinkStats struct {
	SendBytes        uint64      `json:"send_bytes"`
	SendBytesDetails RateDetails `json:"send_bytes_details"`
	RecvBytes        uint64      `json:"recv_bytes"`
	RecvBytesDetails RateDetails `json:"recv_bytes_details"`
}

// MetricsGCQueueLength is metrics of gc queuue length
type MetricsGCQueueLength struct {
	ConnectionClosed       int `json:"connection_closed"`
	ChannelClosed          int `json:"channel_closed"`
	ConsumerDeleted        int `json:"consumer_deleted"`
	ExchangeDeleted        int `json:"exchange_deleted"`
	QueueDeleted           int `json:"queue_deleted"`
	VhostDeleted           int `json:"vhost_deleted"`
	NodeNodeDeleted        int `json:"node_node_deleted"`
	ChannelConsumerDeleted int `json:"channel_consumer_deleted"`
}

// NodeInfo describes a RabbitMQ node and its basic metrics (such as resource usage).
type NodeInfo struct {
	Name      string `json:"name"`
	NodeType  string `json:"type"`
	IsRunning bool   `json:"running"`
	OsPid     OsPid  `json:"os_pid"`

	FdUsed                  int         `json:"fd_used"`
	FdUsedDetails           RateDetails `json:"fd_used_details"`
	FdTotal                 int         `json:"fd_total"`
	ProcUsed                int         `json:"proc_used"`
	ProcUsedDetails         RateDetails `json:"proc_used_details"`
	ProcTotal               int         `json:"proc_total"`
	SocketsUsed             int         `json:"sockets_used"`
	SocketsUsedDetails      RateDetails `json:"sockets_used_details"`
	SocketsTotal            int         `json:"sockets_total"`
	MemUsed                 int         `json:"mem_used"`
	MemUsedDetails          RateDetails `json:"mem_used_details"`
	MemLimit                int         `json:"mem_limit"`
	MemAlarm                bool        `json:"mem_alarm"`
	DiskFree                int         `json:"disk_free"`
	DiskFreeDetails         RateDetails `json:"disk_free_details"`
	DiskFreeLimit           int         `json:"disk_free_limit"`
	DiskFreeAlarm           bool        `json:"disk_free_alarm"`
	GCNum                   uint64      `json:"gc_num"`
	GCNumDetails            RateDetails `json:"gc_num_details"`
	GCBytesReclaimed        uint64      `json:"gc_bytes_reclaimed"`
	GCBytesReclaimedDetails RateDetails `json:"gc_bytes_reclaimed_details"`
	ContextSwitches         uint64      `json:"context_switches"`
	ContextSwitchesDetails  RateDetails `json:"context_switches_details"`

	ConnectionCreated        uint64      `json:"connection_created"`
	ConnectionCreatedDetails RateDetails `json:"connection_created_details"`
	ConnectionClosed         uint64      `json:"connection_closed"`
	ConnectionClosedDetails  RateDetails `json:"connection_closed_details"`
	ChannelCreated           uint64      `json:"channel_created"`
	ChannelCreatedDetails    RateDetails `json:"channel_created_details"`
	ChannelClosed            uint64      `json:"channel_closed"`
	ChannelClosedDetails     RateDetails `json:"channel_closed_details"`
	QueueDeclared            uint64      `json:"queue_declared"`
	QueueDeclaredDetails     RateDetails `json:"queue_declared_details"`
	QueueCreated             uint64      `json:"queue_created"`
	QueueCreatedDetails      RateDetails `json:"queue_created_details"`
	QueueDeleted             uint64      `json:"queue_deleted"`
	QueueDeletedDetails      RateDetails `json:"queue_deleted_details"`

	IOReadCount                         uint64      `json:"io_read_count"`
	IOReadCountDetails                  RateDetails `json:"io_read_count_details"`
	IOReadBytes                         uint64      `json:"io_read_bytes"`
	IOReadBytesDetails                  RateDetails `json:"io_read_bytes"`
	IOReadAvgTime                       float64     `json:"io_read_avg_time"`
	IOReadAvgTimeDetails                RateDetails `json:"io_read_avg_time_details"`
	IOWriteCount                        uint64      `json:"io_write_count"`
	IOWriteCountDetails                 RateDetails `json:"io_write_count_details"`
	IOWriteBytes                        uint64      `json:"io_write_bytes"`
	IOWriteBytesDetails                 RateDetails `json:"io_write_bytes_details"`
	IOWriteAvgTime                      float64     `json:"io_write_avg_time"`
	IOWriteAvgTimeDetails               RateDetails `json:"io_write_avg_time_details"`
	IOSyncCount                         uint64      `json:"io_sync_count"`
	IOSyncCountDetails                  RateDetails `json:"io_sync_count_details"`
	IOSyncAvgTime                       float64     `json:"io_sync_avg_time"`
	IOSyncAvgTimeDetails                RateDetails `json:"io_sync_avg_time_details"`
	IOSeekCount                         uint64      `json:"io_seek_count"`
	IOSeekCountDetails                  RateDetails `json:"io_seek_count_details"`
	IOSeekAvgTime                       float64     `json:"io_seek_avg_time"`
	IOSeekAvgTimeDetails                RateDetails `json:"io_seek_avg_time_details"`
	IOReopenCount                       uint64      `json:"io_reopen_count"`
	IOReopenCountDetails                RateDetails `json:"io_reopen_count_details"`
	IOFileHandleOpenAttemptCount        uint64      `json:"io_file_handle_open_attempt_count"`
	IOFileHandleOpenAttemptCountDetails RateDetails `json:"io_file_handle_open_attempt_count_details"`

	MnesiaRAMTxCount uint64 `json:"mnesia_ram_tx_count"`

	MnesiaRAMTxCountDetails            RateDetails `json:"mnesia_ram_tx_count_details"`
	MnesiaDiskTxCount                  uint64      `json:"mnesia_disk_tx_count"`
	MnesiaDiskTxCountDetails           RateDetails `json:"mnesia_disk_tx_count_details"`
	MsgStoreReadCount                  uint64      `json:"msg_store_read_count"`
	MsgStoreReadCountDetails           RateDetails `json:"msg_store_read_count_details"`
	MsgStoreWriteCount                 uint64      `json:"msg_store_write_count"`
	MsgStoreWriteCountDetails          RateDetails `json:"msg_store_write_count_details"`
	QueueIndexJournalWriteCount        uint64      `json:"queue_index_journal_write_count"`
	QueueIndexJournalWriteCountDetails RateDetails `json:"queue_index_journal_write_count_details"`
	QueueIndexWriteCount               uint64      `json:"queue_index_write_count"`
	QueueIndexWriteCountDetails        RateDetails `json:"queue_index_write_count_details"`
	QueueIndexReadCount                uint64      `json:"queue_index_read_count"`
	QueueIndexReadCountDetails         RateDetails `json:"queue_index_read_count_details"`

	// Erlang scheduler run queue length
	RunQueueLength uint32 `json:"run_queue"`
	Processors     uint32 `json:"processors"`
	Uptime         uint64 `json:"uptime"`

	ExchangeTypes  []ExchangeType  `json:"exchange_types"`
	AuthMechanisms []AuthMechanism `json:"auth_mechanisms"`
	ErlangApps     []ErlangApp     `json:"applications"`
	Contexts       []BrokerContext `json:"contexts"`

	Partitions []string `json:"partitions"`

	ClusterLinks []ClusterLink `json:"cluster_links"`

	MetricsGCQueueLength MetricsGCQueueLength `json:"metrics_gc_queue_length"`
}

//
// GET /api/nodes
//

// ListNodes returns a list of cluster nodes.
func (c *Client) ListNodes() (rec []NodeInfo, err error) {
	req, err := newGETRequest(c, "nodes")
	if err != nil {
		return []NodeInfo{}, err
	}

	if err = executeAndParseRequest(c, req, &rec); err != nil {
		return nil, err
	}

	return rec, nil
}

//
// GET /api/nodes/{name}
//

// {
//   "partitions": [],
//   "os_pid": "39292",
//   "fd_used": 35,
//   "fd_total": 256,
//   "sockets_used": 4,
//   "sockets_total": 138,
//   "mem_used": 69964432,
//   "mem_limit": 2960660889,
//   "mem_alarm": false,
//   "disk_free_limit": 50000000,
//   "disk_free": 188362731520,
//   "disk_free_alarm": false,
//   "proc_used": 370,
//   "proc_total": 1048576,
//   "statistics_level": "fine",
//   "uptime": 98355255,
//   "run_queue": 0,
//   "processors": 8,
//   "exchange_types": [
//     {
//       "name": "topic",
//       "description": "AMQP topic exchange, as per the AMQP specification",
//       "enabled": true
//     },
//     {
//       "name": "x-consistent-hash",
//       "description": "Consistent Hashing Exchange",
//       "enabled": true
//     },
//     {
//       "name": "fanout",
//       "description": "AMQP fanout exchange, as per the AMQP specification",
//       "enabled": true
//     },
//     {
//       "name": "direct",
//       "description": "AMQP direct exchange, as per the AMQP specification",
//       "enabled": true
//     },
//     {
//       "name": "headers",
//       "description": "AMQP headers exchange, as per the AMQP specification",
//       "enabled": true
//     }
//   ],
//   "auth_mechanisms": [
//     {
//       "name": "AMQPLAIN",
//       "description": "QPid AMQPLAIN mechanism",
//       "enabled": true
//     },
//     {
//       "name": "PLAIN",
//       "description": "SASL PLAIN authentication mechanism",
//       "enabled": true
//     },
//     {
//       "name": "RABBIT-CR-DEMO",
//       "description": "RabbitMQ Demo challenge-response authentication mechanism",
//       "enabled": false
//     }
//   ],
//   "applications": [
//     {
//       "name": "amqp_client",
//       "description": "RabbitMQ AMQP Client",
//       "version": "3.2.0"
//     },
//     {
//       "name": "asn1",
//       "description": "The Erlang ASN1 compiler version 2.0.3",
//       "version": "2.0.3"
//     },
//     {
//       "name": "cowboy",
//       "description": "Small, fast, modular HTTP server.",
//       "version": "0.5.0-rmq3.2.0-git4b93c2d"
//     },
//     {
//       "name": "crypto",
//       "description": "CRYPTO version 2",
//       "version": "3.1"
//     },
//     {
//       "name": "inets",
//       "description": "INETS  CXC 138 49",
//       "version": "5.9.6"
//     },
//     {
//       "name": "kernel",
//       "description": "ERTS  CXC 138 10",
//       "version": "2.16.3"
//     },
//     {
//       "name": "mnesia",
//       "description": "MNESIA  CXC 138 12",
//       "version": "4.10"
//     },
//     {
//       "name": "mochiweb",
//       "description": "MochiMedia Web Server",
//       "version": "2.7.0-rmq3.2.0-git680dba8"
//     },
//     {
//       "name": "os_mon",
//       "description": "CPO  CXC 138 46",
//       "version": "2.2.13"
//     },
//     {
//       "name": "public_key",
//       "description": "Public key infrastructure",
//       "version": "0.20"
//     },
//     {
//       "name": "rabbit",
//       "description": "RabbitMQ",
//       "version": "3.2.0"
//     },
//     {
//       "name": "rabbitmq_consistent_hash_exchange",
//       "description": "Consistent Hash Exchange Type",
//       "version": "3.2.0"
//     },
//     {
//       "name": "rabbitmq_management",
//       "description": "RabbitMQ Management Console",
//       "version": "3.2.0"
//     },
//     {
//       "name": "rabbitmq_management_agent",
//       "description": "RabbitMQ Management Agent",
//       "version": "3.2.0"
//     },
//     {
//       "name": "rabbitmq_mqtt",
//       "description": "RabbitMQ MQTT Adapter",
//       "version": "3.2.0"
//     },
//     {
//       "name": "rabbitmq_shovel",
//       "description": "Data Shovel for RabbitMQ",
//       "version": "3.2.0"
//     },
//     {
//       "name": "rabbitmq_shovel_management",
//       "description": "Shovel Status",
//       "version": "3.2.0"
//     },
//     {
//       "name": "rabbitmq_stomp",
//       "description": "Embedded Rabbit Stomp Adapter",
//       "version": "3.2.0"
//     },
//     {
//       "name": "rabbitmq_web_dispatch",
//       "description": "RabbitMQ Web Dispatcher",
//       "version": "3.2.0"
//     },
//     {
//       "name": "rabbitmq_web_stomp",
//       "description": "Rabbit WEB-STOMP - WebSockets to Stomp adapter",
//       "version": "3.2.0"
//     },
//     {
//       "name": "sasl",
//       "description": "SASL  CXC 138 11",
//       "version": "2.3.3"
//     },
//     {
//       "name": "sockjs",
//       "description": "SockJS",
//       "version": "0.3.4-rmq3.2.0-git3132eb9"
//     },
//     {
//       "name": "ssl",
//       "description": "Erlang\/OTP SSL application",
//       "version": "5.3.1"
//     },
//     {
//       "name": "stdlib",
//       "description": "ERTS  CXC 138 10",
//       "version": "1.19.3"
//     },
//     {
//       "name": "webmachine",
//       "description": "webmachine",
//       "version": "1.10.3-rmq3.2.0-gite9359c7"
//     },
//     {
//       "name": "xmerl",
//       "description": "XML parser",
//       "version": "1.3.4"
//     }
//   ],
//   "contexts": [
//     {
//       "description": "Redirect to port 15672",
//       "path": "\/",
//       "port": 55672,
//       "ignore_in_use": true
//     },
//     {
//       "description": "RabbitMQ Management",
//       "path": "\/",
//       "port": 15672
//     }
//   ],
//   "name": "rabbit@mercurio",
//   "type": "disc",
//   "running": true
// }

// GetNode return information about a node.
func (c *Client) GetNode(name string) (rec *NodeInfo, err error) {
	req, err := newGETRequest(c, "nodes/"+url.PathEscape(name))
	if err != nil {
		return nil, err
	}

	if err = executeAndParseRequest(c, req, &rec); err != nil {
		return nil, err
	}

	return rec, nil
}
