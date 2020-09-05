package sarama

import (
	"encoding/binary"
	"encoding/hex"

	"gopkg.in/jcmturner/gokrb5.v7/credentials"
	"gopkg.in/jcmturner/gokrb5.v7/gssapi"
	"gopkg.in/jcmturner/gokrb5.v7/iana/keyusage"
	"gopkg.in/jcmturner/gokrb5.v7/messages"
	"gopkg.in/jcmturner/gokrb5.v7/types"
)

type KafkaGSSAPIHandler struct {
	client         *MockKerberosClient
	badResponse    bool
	badKeyChecksum bool
}

func (h *KafkaGSSAPIHandler) MockKafkaGSSAPI(buffer []byte) []byte {
	// Default payload used for verify
	err := h.client.Login() // Mock client construct keys when login
	if err != nil {
		return nil
	}
	if h.badResponse { // Returns trash
		return []byte{0x00, 0x00, 0x00, 0x01, 0xAD}
	}

	var pack = gssapi.WrapToken{
		Flags:     KRB5_USER_AUTH,
		EC:        12,
		RRC:       0,
		SndSeqNum: 3398292281,
		Payload:   []byte{0x11, 0x00}, // 1100
	}
	// Compute checksum
	if h.badKeyChecksum {
		pack.CheckSum = []byte{0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF}
	} else {
		err = pack.SetCheckSum(h.client.ASRep.DecryptedEncPart.Key, keyusage.GSSAPI_ACCEPTOR_SEAL)
		if err != nil {
			return nil
		}
	}

	packBytes, err := pack.Marshal()
	if err != nil {
		return nil
	}
	lenBytes := len(packBytes)
	response := make([]byte, lenBytes+4)
	copy(response[4:], packBytes)
	binary.BigEndian.PutUint32(response, uint32(lenBytes))
	return response
}

type MockKerberosClient struct {
	asRepBytes  string
	ASRep       messages.ASRep
	credentials *credentials.Credentials
	mockError   error
	errorStage  string
}

func (c *MockKerberosClient) Login() error {
	if c.errorStage == "login" && c.mockError != nil {
		return c.mockError
	}
	c.asRepBytes = "6b8202e9308202e5a003020105a10302010ba22b30293027a103020113a220041e301c301aa003020112a1131b114" +
		"558414d504c452e434f4d636c69656e74a30d1b0b4558414d504c452e434f4da4133011a003020101a10a30081b06636c69656e7" +
		"4a5820156618201523082014ea003020105a10d1b0b4558414d504c452e434f4da220301ea003020102a11730151b066b7262746" +
		"7741b0b4558414d504c452e434f4da382011430820110a003020112a103020101a28201020481ffdb9891175d106818e61008c51" +
		"d0b3462bca92f3bf9d4cfa82de4c4d7aff9994ec87c573e3a3d54dcb2bb79618c76f2bf4a3d006f90d5bdbd049bc18f48be39203" +
		"549ca02acaf63f292b12404f9b74c34b83687119d8f56552ccc0c50ebee2a53bb114c1b4619bb1d5d31f0f49b4d40a08a9b4c046" +
		"2e1398d0b648be1c0e50c552ad16e1d8d8e74263dd0bf0ec591e4797dfd40a9a1be4ae830d03a306e053fd7586fef84ffc5e4a83" +
		"7c3122bf3e6a40fe87e84019f6283634461b955712b44a5f7386c278bff94ec2c2dc0403247e29c2450e853471ceababf9b8911f" +
		"997f2e3010b046d2c49eb438afb0f4c210821e80d4ffa4c9521eb895dcd68610b3feaa682012c30820128a003020112a282011f0" +
		"482011bce73cbce3f1dd17661c412005f0f2257c756fe8e98ff97e6ec24b7bab66e5fd3a3827aeeae4757af0c6e892948122d8b2" +
		"03c8df48df0ef5d142d0e416d688f11daa0fcd63d96bdd431d02b8e951c664eeff286a2be62383d274a04016d5f0e141da58cb86" +
		"331de64063062f4f885e8e9ce5b181ca2fdc67897c5995e0ae1ae0c171a64493ff7bd91bc6d89cd4fce1e2b3ea0a10e34b0d5eda" +
		"aa38ee727b50c5632ed1d2f2b457908e616178d0d80b72af209fb8ac9dbaa1768fa45931392b36b6d8c12400f8ded2efaa0654d0" +
		"da1db966e8b5aab4706c800f95d559664646041fdb38b411c62fc0fbe0d25083a28562b0e1c8df16e62e9d5626b0addee489835f" +
		"eedb0f26c05baa596b69b17f47920aa64b29dc77cfcc97ba47885"
	apRepBytes, err := hex.DecodeString(c.asRepBytes)
	if err != nil {
		return err
	}
	err = c.ASRep.Unmarshal(apRepBytes)
	if err != nil {
		return err
	}
	c.credentials = credentials.New("client", "EXAMPLE.COM").WithPassword("qwerty")
	_, err = c.ASRep.DecryptEncPart(c.credentials)
	if err != nil {
		return err
	}
	return nil
}

func (c *MockKerberosClient) GetServiceTicket(spn string) (messages.Ticket, types.EncryptionKey, error) {
	if c.errorStage == "service_ticket" && c.mockError != nil {
		return messages.Ticket{}, types.EncryptionKey{}, c.mockError
	}
	return c.ASRep.Ticket, c.ASRep.DecryptedEncPart.Key, nil
}

func (c *MockKerberosClient) Domain() string {
	return "EXAMPLE.COM"
}
func (c *MockKerberosClient) CName() types.PrincipalName {
	var p = types.PrincipalName{
		NameType: KRB5_USER_AUTH,
		NameString: []string{
			"kafka",
			"kafka",
		},
	}
	return p
}
func (c *MockKerberosClient) Destroy() {
	// Do nothing.
}
