package pencrypt

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
)

type AesEncryptor struct {
	plaintext	[]byte
	ciphertext	[]byte
	key 		[]byte
	iv 			[]byte

	// true: encrypt, false: decrypt
	encrypt 	bool
}

func NewAesEncryptor(text []byte, encrypt bool) *AesEncryptor {
	aes := &AesEncryptor{
		encrypt: encrypt,
	}
	if encrypt {
		aes.plaintext = text
	} else {
		aes.ciphertext = text
	}
	return aes
}

// SetKey
// Note that key length must be 16, 24 or 32 bytes to select AES-128, AES-192, or AES-256
// Note that AES block size is 16 bytes
func (p *AesEncryptor) SetKey(key []byte) {
	p.key = key
}

// SetCbcIV
// The length of iv must be the same as the key size
func (p *AesEncryptor) SetCbcIV(iv []byte) {
	p.iv = iv
}

func (p *AesEncryptor) PKCS7Padding(text []byte, blockSize int) []byte {
	padding := blockSize - len(text) % blockSize
	padText := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(text, padText...)
}

func (p *AesEncryptor) PKCS7UnPadding(text []byte) []byte {
	length := len(text)
	unPadding := int(text[length-1])
	return text[:(length - unPadding)]
}

// Encrypt encrypts data with AES algorithm in CBC mode
func (p *AesEncryptor) Encrypt() ([]byte, error) {
	block, err := aes.NewCipher(p.key)
	if err != nil {
		return nil, err
	}

	blockSize := block.BlockSize()
	plaintext := p.PKCS7Padding(p.plaintext, blockSize)
	p.ciphertext = make([]byte, len(plaintext))
	iv := p.key[:blockSize]
	if len(p.iv) > 0 { iv = p.iv }
	blockMode := cipher.NewCBCEncrypter(block, iv)
	blockMode.CryptBlocks(p.ciphertext, plaintext)
	return p.ciphertext, nil
}

// Encrypt2Base64
// Encrypt to Base64
func (p *AesEncryptor) Encrypt2Base64() (string, error) {
	if _, err := p.Encrypt(); err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(p.ciphertext), nil
}

// Decrypt decrypts cipher text with AES algorithm in CBC mode
func (p *AesEncryptor) Decrypt() ([]byte, error) {
	block, err := aes.NewCipher(p.key)
	if err != nil {
		return nil, err
	}

	plaintext := make([]byte, len(p.ciphertext))
	iv := p.key[:block.BlockSize()]
	if len(p.iv) > 0 { iv = p.iv }
	blockMode := cipher.NewCBCDecrypter(block, iv)
	blockMode.CryptBlocks(plaintext, p.ciphertext)
	p.plaintext = p.PKCS7UnPadding(plaintext)
	return p.plaintext, nil
}

// Decrypt5Base64
// Decrypt from Base64
func (p *AesEncryptor) Decrypt5Base64() ([]byte, error) {
	text, err := base64.StdEncoding.DecodeString(string(p.ciphertext))
	if err != nil {
		return nil, err
	}
	p.ciphertext = text
	return p.Decrypt()
}
