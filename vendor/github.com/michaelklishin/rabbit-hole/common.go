package rabbithole

import "strconv"

// Extra arguments as a map (on queues, bindings, etc)
type Properties map[string]interface{}

// Port used by RabbitMQ or clients
type Port int

func (p *Port) UnmarshalJSON(b []byte) error {
	stringValue := string(b)
	var parsed int64
	var err error
	if stringValue[0] == '"' && stringValue[len(stringValue)-1] == '"' {
		parsed, err = strconv.ParseInt(stringValue[1:len(stringValue)-1], 10, 32)
	} else {
		parsed, err = strconv.ParseInt(stringValue, 10, 32)
	}
	if err == nil {
		*p = Port(int(parsed))
	}
	return err
}

// RateDetailSample single touple
type RateDetailSample struct {
	Sample    int64 `json:"sample"`
	Timestamp int64 `json:"timestamp"`
}

// Rate of change of a numerical value
type RateDetails struct {
	Rate    float32            `json:"rate"`
	Samples []RateDetailSample `json:"samples"`
}

// RabbitMQ context (Erlang app) running on
// a node
type BrokerContext struct {
	Node        string `json:"node"`
	Description string `json:"description"`
	Path        string `json:"path"`
	Port        Port   `json:"port"`
	Ignore      bool   `json:"ignore_in_use"`
}

// Basic published messages statistics
type MessageStats struct {
	Publish             int64       `json:"publish"`
	PublishDetails      RateDetails `json:"publish_details"`
	Deliver             int64       `json:"deliver"`
	DeliverDetails      RateDetails `json:"deliver_details"`
	DeliverNoAck        int64       `json:"deliver_noack"`
	DeliverNoAckDetails RateDetails `json:"deliver_noack_details"`
	DeliverGet          int64       `json:"deliver_get"`
	DeliverGetDetails   RateDetails `json:"deliver_get_details"`
	Redeliver           int64       `json:"redeliver"`
	RedeliverDetails    RateDetails `json:"redeliver_details"`
	Get                 int64       `json:"get"`
	GetDetails          RateDetails `json:"get_details"`
	GetNoAck            int64       `json:"get_no_ack"`
	GetNoAckDetails     RateDetails `json:"get_no_ack_details"`
	Ack                 int64       `json:"ack"`
	AckDetails          RateDetails `json:"ack_details"`
}
