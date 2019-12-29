package pac

import (
	"bytes"
	"fmt"

	"gopkg.in/jcmturner/rpc.v1/mstypes"
	"gopkg.in/jcmturner/rpc.v1/ndr"
)

// Claims reference: https://msdn.microsoft.com/en-us/library/hh553895.aspx

// ClientClaimsInfo implements https://msdn.microsoft.com/en-us/library/hh536365.aspx
type ClientClaimsInfo struct {
	ClaimsSetMetadata mstypes.ClaimsSetMetadata
	ClaimsSet         mstypes.ClaimsSet
}

// Unmarshal bytes into the ClientClaimsInfo struct
func (k *ClientClaimsInfo) Unmarshal(b []byte) (err error) {
	dec := ndr.NewDecoder(bytes.NewReader(b))
	m := new(mstypes.ClaimsSetMetadata)
	err = dec.Decode(m)
	if err != nil {
		err = fmt.Errorf("error unmarshaling ClientClaimsInfo ClaimsSetMetadata: %v", err)
	}
	k.ClaimsSetMetadata = *m
	k.ClaimsSet, err = k.ClaimsSetMetadata.ClaimsSet()
	if err != nil {
		err = fmt.Errorf("error unmarshaling ClientClaimsInfo ClaimsSet: %v", err)
	}
	return
}
