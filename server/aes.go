package main

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"io"
	"bytes"
)

// Encrypt encrypts plaintext using AES-256-CBC and returns a base64-encoded ciphertext.
func Encrypt(key []byte, plaintext string) (string, error) {
	if len(key) != 32 {
		return "", errors.New("key must be 32 bytes for AES-256")
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	// Generate a random IV
	iv := make([]byte, aes.BlockSize)
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return "", err
	}

	// PKCS7 padding
	padding := aes.BlockSize - len(plaintext)%aes.BlockSize
	padtext := append([]byte(plaintext), bytes.Repeat([]byte{byte(padding)}, padding)...)

	ciphertext := make([]byte, len(padtext))
	mode := cipher.NewCBCEncrypter(block, iv)
	mode.CryptBlocks(ciphertext, padtext)

	// Prepend IV to ciphertext
	final := append(iv, ciphertext...)
	return base64.StdEncoding.EncodeToString(final), nil
}

// Decrypt decrypts a base64-encoded AES-256-CBC ciphertext and returns the plaintext.
func Decrypt(key []byte, ciphertextBase64 string) (string, error) {
	if len(key) != 32 {
		return "", errors.New("key must be 32 bytes for AES-256")
	}

	data, err := base64.StdEncoding.DecodeString(ciphertextBase64)
	if err != nil {
		return "", err
	}

	if len(data) < aes.BlockSize {
		return "", errors.New("ciphertext too short")
	}

	iv := data[:aes.BlockSize]
	ciphertext := data[aes.BlockSize:]

	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	if len(ciphertext)%aes.BlockSize != 0 {
		return "", errors.New("ciphertext is not a multiple of the block size")
	}

	mode := cipher.NewCBCDecrypter(block, iv)
	plaintext := make([]byte, len(ciphertext))
	mode.CryptBlocks(plaintext, ciphertext)

	// Remove PKCS7 padding
	length := len(plaintext)
	padLen := int(plaintext[length-1])
	if padLen > aes.BlockSize || padLen == 0 {
		return "", errors.New("invalid padding")
	}
	for i := length - padLen; i < length; i++ {
		if plaintext[i] != byte(padLen) {
			return "", errors.New("invalid padding")
		}
	}
	return string(plaintext[:length-padLen]), nil
}
