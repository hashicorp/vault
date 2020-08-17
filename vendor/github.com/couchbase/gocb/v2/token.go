package gocb

import (
	"encoding/json"
	"fmt"
	"strconv"

	gocbcore "github.com/couchbase/gocbcore/v9"
)

// MutationToken holds the mutation state information from an operation.
type MutationToken struct {
	token      gocbcore.MutationToken
	bucketName string
}

type bucketToken struct {
	SeqNo  uint64 `json:"seqno"`
	VbUUID string `json:"vbuuid"`
}

// BucketName returns the name of the bucket that this token belongs to.
func (mt MutationToken) BucketName() string {
	return mt.bucketName
}

// PartitionUUID returns the UUID of the vbucket that this token belongs to.
func (mt MutationToken) PartitionUUID() uint64 {
	return uint64(mt.token.VbUUID)
}

// PartitionID returns the ID of the vbucket that this token belongs to.
func (mt MutationToken) PartitionID() uint64 {
	return uint64(mt.token.VbID)
}

// SequenceNumber returns the sequence number of the vbucket that this token belongs to.
func (mt MutationToken) SequenceNumber() uint64 {
	return uint64(mt.token.SeqNo)
}

func (mt bucketToken) MarshalJSON() ([]byte, error) {
	info := []interface{}{mt.SeqNo, mt.VbUUID}
	return json.Marshal(info)
}

func (mt *bucketToken) UnmarshalJSON(data []byte) error {
	info := []interface{}{&mt.SeqNo, &mt.VbUUID}
	return json.Unmarshal(data, &info)
}

type bucketTokens map[string]*bucketToken
type mutationStateData map[string]*bucketTokens

type searchMutationState map[string]map[string]uint64

// MutationState holds and aggregates MutationToken's across multiple operations.
type MutationState struct {
	tokens []MutationToken
}

// NewMutationState creates a new MutationState for tracking mutation state.
func NewMutationState(tokens ...MutationToken) *MutationState {
	mt := &MutationState{}
	mt.Add(tokens...)
	return mt
}

// Add includes an operation's mutation information in this mutation state.
func (mt *MutationState) Add(tokens ...MutationToken) {
	for _, token := range tokens {
		if token.bucketName != "" {
			mt.tokens = append(mt.tokens, token)
		}
	}
}

// MutationStateInternal specifies internal operations.
// Internal: This should never be used and is not supported.
type MutationStateInternal struct {
	mt *MutationState
}

// Internal return a new MutationStateInternal.
// Internal: This should never be used and is not supported.
func (mt *MutationState) Internal() *MutationStateInternal {
	return &MutationStateInternal{
		mt: mt,
	}
}

// Add includes an operation's mutation information in this mutation state.
func (mti *MutationStateInternal) Add(bucket string, tokens ...gocbcore.MutationToken) {
	for _, token := range tokens {
		mti.mt.Add(MutationToken{
			bucketName: bucket,
			token:      token,
		})
	}
}

// Tokens returns the tokens belonging to the mutation state.
func (mti *MutationStateInternal) Tokens() []MutationToken {
	return mti.mt.tokens
}

// MarshalJSON marshal's this mutation state to JSON.
func (mt *MutationState) MarshalJSON() ([]byte, error) {
	var data mutationStateData
	for _, token := range mt.tokens {
		if data == nil {
			data = make(mutationStateData)
		}

		bucketName := token.bucketName
		if (data)[bucketName] == nil {
			tokens := make(bucketTokens)
			(data)[bucketName] = &tokens
		}

		vbID := fmt.Sprintf("%d", token.token.VbID)
		stateToken := (*(data)[bucketName])[vbID]
		if stateToken == nil {
			stateToken = &bucketToken{}
			(*(data)[bucketName])[vbID] = stateToken
		}

		stateToken.SeqNo = uint64(token.token.SeqNo)
		stateToken.VbUUID = fmt.Sprintf("%d", token.token.VbUUID)

	}

	return json.Marshal(data)
}

// UnmarshalJSON unmarshal's a mutation state from JSON.
func (mt *MutationState) UnmarshalJSON(data []byte) error {
	var stateData mutationStateData
	err := json.Unmarshal(data, &stateData)
	if err != nil {
		return err
	}

	for bucketName, bTokens := range stateData {
		for vbIDStr, stateToken := range *bTokens {
			vbID, err := strconv.Atoi(vbIDStr)
			if err != nil {
				return err
			}
			vbUUID, err := strconv.Atoi(stateToken.VbUUID)
			if err != nil {
				return err
			}
			token := MutationToken{
				bucketName: bucketName,
				token: gocbcore.MutationToken{
					VbID:   uint16(vbID),
					VbUUID: gocbcore.VbUUID(vbUUID),
					SeqNo:  gocbcore.SeqNo(stateToken.SeqNo),
				},
			}

			mt.tokens = append(mt.tokens, token)
		}
	}

	return nil
}

// toSearchMutationState is specific to search, search doesn't accept tokens in the same format as other services.
func (mt *MutationState) toSearchMutationState() searchMutationState {
	data := make(searchMutationState)
	for _, token := range mt.tokens {
		_, ok := data[token.bucketName]
		if !ok {
			data[token.bucketName] = make(map[string]uint64)
		}

		data[token.bucketName][fmt.Sprintf("%d/%d", token.token.VbID, token.token.VbUUID)] = uint64(token.token.SeqNo)
	}

	return data
}
