// Copyright (c) 2021 Snowflake Computing Inc. All right reserved.

package gosnowflake

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"strconv"
)

type snowflakeFileEncryption struct {
	QueryStageMasterKey string `json:"queryStageMasterKey,omitempty"`
	QueryID             string `json:"queryId,omitempty"`
	SMKID               int64  `json:"smkId,omitempty"`
}

// PUT requests return a single encryptionMaterial object whereas GET requests
// return a slice (array) of encryptionMaterial objects, both under the field
// 'encryptionMaterial'
type encryptionWrapper struct {
	snowflakeFileEncryption
	EncryptionMaterials []snowflakeFileEncryption
}

// override default behavior for wrapper
func (ew *encryptionWrapper) UnmarshalJSON(data []byte) error {
	// if GET, unmarshal slice of encryptionMaterial
	if err := json.Unmarshal(data, &ew.EncryptionMaterials); err == nil {
		return err
	}
	// else (if PUT), unmarshal the encryptionMaterial itself
	return json.Unmarshal(data, &ew.snowflakeFileEncryption)
}

type encryptMetadata struct {
	key     string
	iv      string
	matdesc string
}

// encryptStream encrypts a stream buffer using AES128 block cipher in CBC mode
// with PKCS5 padding
func encryptStream(
	sfe *snowflakeFileEncryption,
	src io.Reader,
	out io.Writer,
	chunkSize int) (*encryptMetadata, error) {
	if chunkSize == 0 {
		chunkSize = aes.BlockSize * 4 * 1024
	}
	decodedKey, _ := base64.StdEncoding.DecodeString(sfe.QueryStageMasterKey)
	keySize := len(decodedKey)

	fileKey := getSecureRandom(keySize)
	block, _ := aes.NewCipher(fileKey)
	ivData := getSecureRandom(block.BlockSize())

	mode := cipher.NewCBCEncrypter(block, ivData)
	cipherText := make([]byte, chunkSize)

	// encrypt file with CBC
	var err error
	for {
		chunk := make([]byte, chunkSize)
		n, err := src.Read(chunk)
		if n == 0 || err != nil {
			break
		} else if n%aes.BlockSize != 0 || n != chunkSize {
			chunk = padBytesLength(chunk[:n], aes.BlockSize)
		}
		mode.CryptBlocks(cipherText, chunk)
		out.Write(cipherText[:len(chunk)])

	}
	if err != nil {
		return nil, err
	}

	// encrypt key with ECB
	fileKey = padBytesLength(fileKey, block.BlockSize())
	encryptedFileKey := make([]byte, len(fileKey))
	if err = encryptECB(encryptedFileKey, fileKey, decodedKey); err != nil {
		return nil, err
	}

	matDesc := materialDescriptor{
		strconv.Itoa(int(sfe.SMKID)),
		sfe.QueryID,
		strconv.Itoa(keySize * 8),
	}

	return &encryptMetadata{
		base64.StdEncoding.EncodeToString(encryptedFileKey),
		base64.StdEncoding.EncodeToString(ivData),
		matdescToUnicode(matDesc),
	}, nil
}

func encryptECB(encrypted []byte, fileKey []byte, decodedKey []byte) error {
	block, _ := aes.NewCipher(decodedKey)
	if len(fileKey)%block.BlockSize() != 0 {
		return fmt.Errorf("input not full of blocks")
	}
	if len(encrypted) < len(fileKey) {
		return fmt.Errorf("output length is smaller than input length")
	}
	for len(fileKey) > 0 {
		block.Encrypt(encrypted, fileKey[:block.BlockSize()])
		encrypted = encrypted[block.BlockSize():]
		fileKey = fileKey[block.BlockSize():]
	}
	return nil
}

func decryptECB(decrypted []byte, keyBytes []byte, decodedKey []byte) error {
	block, _ := aes.NewCipher(decodedKey)
	if len(keyBytes)%block.BlockSize() != 0 {
		return fmt.Errorf("input not full of blocks")
	}
	if len(decrypted) < len(keyBytes) {
		return fmt.Errorf("output length is smaller than input length")
	}
	for len(keyBytes) > 0 {
		block.Decrypt(decrypted, keyBytes[:block.BlockSize()])
		keyBytes = keyBytes[block.BlockSize():]
		decrypted = decrypted[block.BlockSize():]
	}
	return nil
}

func encryptFile(
	sfe *snowflakeFileEncryption,
	filename string,
	chunkSize int,
	tmpDir string) (
	*encryptMetadata, string, error) {
	if chunkSize == 0 {
		chunkSize = aes.BlockSize * 4 * 1024
	}
	tmpOutputFile, _ := ioutil.TempFile(tmpDir, baseName(filename)+"#")
	infile, err := os.OpenFile(filename, os.O_CREATE|os.O_RDONLY, os.ModePerm)
	if err != nil {
		return nil, "", err
	}
	meta, err := encryptStream(sfe, infile, tmpOutputFile, chunkSize)
	if err != nil {
		return nil, "", err
	}
	return meta, tmpOutputFile.Name(), nil
}

func decryptFile(
	metadata *encryptMetadata,
	sfe *snowflakeFileEncryption,
	filename string,
	chunkSize int,
	tmpDir string) (
	string, error) {
	if chunkSize == 0 {
		chunkSize = aes.BlockSize * 4 * 1024
	}
	decodedKey, _ := base64.StdEncoding.DecodeString(sfe.QueryStageMasterKey)
	keyBytes, _ := base64.StdEncoding.DecodeString(metadata.key) // encrypted file key
	ivBytes, _ := base64.StdEncoding.DecodeString(metadata.iv)

	// decrypt file key
	decryptedKey := make([]byte, len(keyBytes))
	if err := decryptECB(decryptedKey, keyBytes, decodedKey); err != nil {
		return "", err
	}
	decryptedKey = paddingTrim(decryptedKey)

	// decrypt file
	block, _ := aes.NewCipher(decryptedKey)
	mode := cipher.NewCBCDecrypter(block, ivBytes)

	tmpOutputFile, err := ioutil.TempFile(tmpDir, baseName(filename)+"#")
	if err != nil {
		return "", err
	}
	defer tmpOutputFile.Close()
	infile, err := os.OpenFile(filename, os.O_RDONLY, os.ModePerm)
	if err != nil {
		return "", err
	}
	defer infile.Close()
	var totalFileSize int
	var prevChunk []byte
	for {
		chunk := make([]byte, chunkSize)
		n, err := infile.Read(chunk)
		if n == 0 || err != nil {
			break
		}
		totalFileSize += n
		chunk = chunk[:n]
		mode.CryptBlocks(chunk, chunk)
		tmpOutputFile.Write(chunk)
		prevChunk = chunk
	}
	if err != nil {
		return "", err
	}
	if prevChunk != nil {
		totalFileSize -= paddingOffset(prevChunk)
	}
	tmpOutputFile.Truncate(int64(totalFileSize))
	return tmpOutputFile.Name(), nil
}

type materialDescriptor struct {
	SmkID   string `json:"smkId"`
	QueryID string `json:"queryId"`
	KeySize string `json:"keySize"`
}

func matdescToUnicode(matdesc materialDescriptor) string {
	s, _ := json.Marshal(&matdesc)
	return string(s)
}

func getSecureRandom(byteLength int) []byte {
	token := make([]byte, byteLength)
	rand.Read(token)
	return token
}

func padBytesLength(src []byte, blockSize int) []byte {
	padLength := blockSize - len(src)%blockSize
	padText := bytes.Repeat([]byte{byte(padLength)}, padLength)
	return append(src, padText...)
}

func paddingTrim(src []byte) []byte {
	unpadding := src[len(src)-1]
	return src[:len(src)-int(unpadding)]
}

func paddingOffset(src []byte) int {
	length := len(src)
	return int(src[length-1])
}

type contentKey struct {
	KeyID         string `json:"KeyId,omitempty"`
	EncryptionKey string `json:"EncryptedKey,omitempty"`
	Algorithm     string `json:"Algorithm,omitempty"`
}

type encryptionAgent struct {
	Protocol            string `json:"Protocol,omitempty"`
	EncryptionAlgorithm string `json:"EncryptionAlgorithm,omitempty"`
}

type keyMetadata struct {
	EncryptionLibrary string `json:"EncryptionLibrary,omitempty"`
}

type encryptionData struct {
	EncryptionMode      string          `json:"EncryptionMode,omitempty"`
	WrappedContentKey   contentKey      `json:"WrappedContentKey,omitempty"`
	EncryptionAgent     encryptionAgent `json:"EncryptionAgent,omitempty"`
	ContentEncryptionIV string          `json:"ContentEncryptionIV,omitempty"`
	KeyWrappingMetadata keyMetadata     `json:"KeyWrappingMetadata,omitempty"`
}
