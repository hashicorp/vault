package rabbithole

import (
	"encoding/json"
	"errors"
	"strconv"
)

// Properties are extra arguments as a map (on queues, bindings, etc)
type Properties map[string]interface{}

// Port used by RabbitMQ or clients
type Port int

// UnmarshalJSON deserialises a port that can be an integer or string
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

// RateDetails fields represent rate of change of a numerical value
type RateDetails struct {
	Rate    float32            `json:"rate"`
	Samples []RateDetailSample `json:"samples"`
}

// BrokerContext represents a context (Erlang application) running on a node
// a node
type BrokerContext struct {
	Node        string `json:"node"`
	Description string `json:"description"`
	Path        string `json:"path"`
	Port        Port   `json:"port"`
	Ignore      bool   `json:"ignore_in_use"`
}

// MessageStats fields repsent a number of metrics related to published messages
type MessageStats struct {
	Publish                 int64       `json:"publish"`
	PublishDetails          RateDetails `json:"publish_details"`
	Deliver                 int64       `json:"deliver"`
	DeliverDetails          RateDetails `json:"deliver_details"`
	DeliverNoAck            int64       `json:"deliver_noack"`
	DeliverNoAckDetails     RateDetails `json:"deliver_noack_details"`
	DeliverGet              int64       `json:"deliver_get"`
	DeliverGetDetails       RateDetails `json:"deliver_get_details"`
	Redeliver               int64       `json:"redeliver"`
	RedeliverDetails        RateDetails `json:"redeliver_details"`
	Get                     int64       `json:"get"`
	GetDetails              RateDetails `json:"get_details"`
	GetNoAck                int64       `json:"get_no_ack"`
	GetNoAckDetails         RateDetails `json:"get_no_ack_details"`
	Ack                     int64       `json:"ack"`
	AckDetails              RateDetails `json:"ack_details"`
	ReturnUnroutable        int64       `json:"return_unroutable"`
	ReturnUnroutableDetails RateDetails `json:"return_unroutable_details"`
	DropUnroutable          int64       `json:"drop_unroutable"`
	DropUnroutableDetails   RateDetails `json:"drop_unroutable_details"`
}

// URISet represents a set of URIs used by Shovel, Federation, and so on.
// The URIs from this set are tried until one of them succeeds
// (a shovel or federation link successfully connects and authenticates with it)
type URISet []string

// UnmarshalJSON can unmarshal a single URI string or a list of
// URI strings
func (s *URISet) UnmarshalJSON(b []byte) error {
	// the value is a single URI, a string
	if b[0] == '"' {
		var uri string
		if err := json.Unmarshal(b, &uri); err != nil {
			return err
		}
		*s = []string{uri}
		return nil
	}

	// the value is a list
	var uris []string
	if err := json.Unmarshal(b, &uris); err != nil {
		return err
	}
	*s = uris
	return nil
}

// AutoDelete is a boolean but RabbitMQ may return the string "undefined"
type AutoDelete bool

// UnmarshalJSON can unmarshal a string or a boolean
func (d *AutoDelete) UnmarshalJSON(b []byte) error {
	switch string(b) {
	case "\"undefined\"":
		// auto_delete is "undefined", map it to true
		*d = AutoDelete(true)
	case "true":
		*d = AutoDelete(true)
	case "false":
		*d = AutoDelete(false)
	default:
		return errors.New("Unknown value of auto_delete")
	}

	return nil
}
