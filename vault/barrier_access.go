package vault

// BarrierEncryptorAccess is a wrapper around BarrierEncryptor that allows Core
// to expose its barrier encrypt/decrypt operations through BarrierEncryptorAccess()
// while restricting the ability to modify Core.barrier itself.
type BarrierEncryptorAccess struct {
	barrierEncryptor BarrierEncryptor
}

var _ BarrierEncryptor = (*BarrierEncryptorAccess)(nil)

func NewBarrierEncryptorAccess(barrierEncryptor BarrierEncryptor) *BarrierEncryptorAccess {
	return &BarrierEncryptorAccess{barrierEncryptor: barrierEncryptor}
}

func (b *BarrierEncryptorAccess) Encrypt(key string, plaintext []byte) ([]byte, error) {
	return b.barrierEncryptor.Encrypt(key, plaintext)
}

func (b *BarrierEncryptorAccess) Decrypt(key string, ciphertext []byte) ([]byte, error) {
	return b.barrierEncryptor.Decrypt(key, ciphertext)
}
