package encoding

import (
	"bytes"
	"crypto/aes"
	"errors"
)

func AESEncrypt(data, key []byte) ([]byte, error) {
	plaintext := PKCS7Pad(data)
	cipher, err := aes.NewCipher(key[:aes.BlockSize])
	if err != nil {
		return nil, err
	}

	if len(plaintext)%aes.BlockSize != 0 {
		return nil, errors.New("Need a multiple of the blocksize 16")
	}

	ciphertext := make([]byte, 0)
	text := make([]byte, 16)
	for len(plaintext) > 0 {
		// 每次运算一个block
		cipher.Encrypt(text, plaintext)
		plaintext = plaintext[aes.BlockSize:]
		ciphertext = append(ciphertext, text...)
	}
	return ciphertext, nil
}

// 加密
func AESEncryptNew(data []byte, key []byte) []byte {
	c, _ := aes.NewCipher(key)

	data = PKCS7Pad(data)
	out := make([]byte, len(data))

	c.Encrypt(out, data)

	return out
}

// 解密
func AESDecrypt(ciphertext, key []byte) ([]byte, error) {
	cipher, err := aes.NewCipher(key[:aes.BlockSize])
	if err != nil {
		return nil, err
	}

	if len(ciphertext)%aes.BlockSize != 0 {
		return nil, errors.New("Need a multiple of the blocksize 16")
	}

	plaintext := make([]byte, 0)
	text := make([]byte, 16)
	for len(ciphertext) > 0 {
		cipher.Decrypt(text, ciphertext)
		ciphertext = ciphertext[aes.BlockSize:]
		plaintext = append(plaintext, text...)
	}

	return PKCS7UPad(plaintext), nil
}

// Padding补全
func PKCS7Pad(data []byte) []byte {
	padding := aes.BlockSize - len(data)%aes.BlockSize
	padtext := bytes.Repeat([]byte{byte(0)}, padding)
	padtext[len(padtext)-1] = byte(padding)
	return append(data, padtext...)
}

func PKCS7UPad(data []byte) []byte {
	padLength := int(data[len(data)-1])
	if padLength >= len(data) {
		return nil
	}
	return data[:len(data)-padLength]
}
