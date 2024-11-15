package msgraphgocore

import (
	"github.com/microsoft/kiota-abstractions-go/serialization"
)

type batchResponse struct {
	responses     []BatchItem
	indexResponse map[string]BatchItem
	isIndexed     bool
}

func NewBatchResponse() BatchResponse {
	return &batchResponse{
		indexResponse: make(map[string]BatchItem),
		isIndexed:     false,
	}
}

// GetResponses returns a slice of BatchItem to the user
func (br *batchResponse) GetResponses() []BatchItem {
	return br.responses
}

// SetResponses adds a slice of BatchItem to the response
func (br *batchResponse) SetResponses(responses []BatchItem) {
	br.responses = responses
}

// AddResponses adds elements to existing response
func (br *batchResponse) AddResponses(responses []BatchItem) {
	for _, v := range responses {
		br.responses = append(br.responses, v)
	}
}

// GetResponseById returns a response payload as a batch item
func (br *batchResponse) GetResponseById(itemId string) BatchItem {
	if !br.isIndexed {

		for _, resp := range br.GetResponses() {
			br.indexResponse[*(resp.GetId())] = resp
		}

		br.isIndexed = true
	}

	return br.indexResponse[itemId]
}

// CreateBatchResponseDiscriminator creates a new instance of the appropriate class based on discriminator value
func CreateBatchResponseDiscriminator(serialization.ParseNode) (serialization.Parsable, error) {
	return NewBatchResponse(), nil
}

// BatchResponse instance of batch request result payload
type BatchResponse interface {
	serialization.Parsable
	GetResponses() []BatchItem
	SetResponses(responses []BatchItem)
	AddResponses(responses []BatchItem)
	GetResponseById(itemId string) BatchItem
	GetFailedResponses() map[string]int32
	GetStatusCodes() map[string]int32
}

// Serialize serializes information the current object
func (br *batchResponse) Serialize(serialization.SerializationWriter) error {
	panic("batch responses are not serializable")
}

// GetFieldDeserializers the deserialization information for the current model
func (br *batchResponse) GetFieldDeserializers() map[string]func(serialization.ParseNode) error {
	res := make(map[string]func(serialization.ParseNode) error)
	res["responses"] = func(n serialization.ParseNode) error {
		val, err := n.GetCollectionOfObjectValues(CreateBatchRequestItemDiscriminator)
		if err != nil {
			return err
		}
		if val != nil {
			res := make([]BatchItem, len(val))
			for i, v := range val {
				res[i] = v.(BatchItem)
			}
			br.SetResponses(res)
		}
		return nil
	}
	return res
}

// GetFailedResponses returns a map of responses that failed
func (br *batchResponse) GetFailedResponses() map[string]int32 {
	statuses := make(map[string]int32)
	for _, response := range br.GetResponses() {
		if *response.GetStatus() > 399 && *response.GetStatus() < 600 {
			statuses[*response.GetId()] = *response.GetStatus()
		}
	}
	return statuses
}

// GetStatusCodes returns a map of responses statuses and the status codes
func (br *batchResponse) GetStatusCodes() map[string]int32 {
	statuses := make(map[string]int32)
	for _, response := range br.GetResponses() {
		statuses[*response.GetId()] = *response.GetStatus()
	}
	return statuses
}
