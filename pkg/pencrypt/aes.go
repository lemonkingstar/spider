package pencrypt

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"errors"
)

type AesEncryptor struct {
	// Note that key length must be 16, 24 or 32 bytes to select AES-128, AES-192, or AES-256
	// Note that AES block size is 16 bytes
	key []byte
	// Note The length of iv must be the same as the block size
	iv []byte
}

func NewAesEncryptor(key, iv []byte) (*AesEncryptor, error) {
	switch len(key) {
	case 16, 24, 32:
	default:
		return nil, errors.New("invalid key length")
	}

	if iv != nil && len(iv) != aes.BlockSize {
		return nil, errors.New("invalid iv length")
	}
	return &AesEncryptor{key: key, iv: iv}, nil
}

func (p *AesEncryptor) pkcs7Padding(text []byte, blockSize int) []byte {
	padding := blockSize - len(text)%blockSize
	padText := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(text, padText...)
}

func (p *AesEncryptor) pkcs7UnPadding(text []byte) ([]byte, error) {
	length := len(text)
	if length == 0 {
		return nil, errors.New("padding text is nil")
	}
	unPadding := int(text[length-1])
	return text[:(length - unPadding)], nil
}

// Encrypt encrypts data with AES algorithm in CBC mode
func (p *AesEncryptor) Encrypt(plaintext []byte) ([]byte, error) {
	block, err := aes.NewCipher(p.key)
	if err != nil {
		return nil, err
	}

	blockSize := block.BlockSize()
	plaintext2 := p.pkcs7Padding(plaintext, blockSize)
	ciphertext := make([]byte, len(plaintext2))
	iv := p.key[:blockSize]
	if len(p.iv) > 0 {
		iv = p.iv
	}
	blockMode := cipher.NewCBCEncrypter(block, iv)
	blockMode.CryptBlocks(ciphertext, plaintext2)
	return ciphertext, nil
}

func (p *AesEncryptor) EncryptToBase64(plaintext string) (string, error) {
	ciphertext, err := p.Encrypt([]byte(plaintext))
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(ciphertext), nil
}

// Decrypt decrypts cipher text with AES algorithm in CBC mode
func (p *AesEncryptor) Decrypt(ciphertext []byte) ([]byte, error) {
	block, err := aes.NewCipher(p.key)
	if err != nil {
		return nil, err
	}

	plaintext := make([]byte, len(ciphertext))
	iv := p.key[:block.BlockSize()]
	if len(p.iv) > 0 {
		iv = p.iv
	}
	blockMode := cipher.NewCBCDecrypter(block, iv)
	blockMode.CryptBlocks(plaintext, ciphertext)
	plaintext2, err := p.pkcs7UnPadding(plaintext)
	if err != nil {
		return nil, err
	}
	return plaintext2, nil
}

func (p *AesEncryptor) DecryptFromBase64(ciphertext string) (string, error) {
	plaintext, err := base64.StdEncoding.DecodeString(ciphertext)
	if err != nil {
		return "", err
	}
	plaintext2, err := p.Decrypt(plaintext)
	if err != nil {
		return "", err
	}
	return string(plaintext2), nil
}
