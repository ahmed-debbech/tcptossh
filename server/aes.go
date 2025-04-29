package main

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"io"
	"strconv"
	"time"
)

const (
	nonceSize = 12 // Standard nonce size for AES-GCM
	timeTolerance = 5 * time.Second // Allow for some clock drift
)

// generateKey creates a random AES key. Keep this key SECRET!
func GenerateKey() ([]byte, error) {
	key := make([]byte, 32) // AES-256 key
	_, err := io.ReadFull(rand.Reader, key)
	if err != nil {
		return nil, err
	}
	return key, nil
}

// encrypt encrypts the data using AES-GCM with a nonce and timestamp.
// The output format is base64(nonce|timestamp|ciphertext).
func Encrypt(key []byte, plaintext string) (string, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	nonce := make([]byte, nonceSize)
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return "", err
	}

	timestamp := time.Now().UTC().UnixNano()
	timestampBytes := []byte(strconv.FormatInt(timestamp, 10))

	// Prepend nonce and timestamp to the plaintext
	data := append(nonce, append(timestampBytes, []byte(plaintext)...)...)

	ciphertext := gcm.Seal(nil, nonce, data, nil)

	// Encode to base64 for easier handling and transmission
	return base64.StdEncoding.EncodeToString(ciphertext), nil
}

// decrypt decrypts the data, verifies the nonce, timestamp, and integrity.
// It expects the input to be in the format base64(nonce|timestamp|ciphertext).
func Decrypt(key []byte, ciphertextBase64 string) (string, error) {
	ciphertext, err := base64.StdEncoding.DecodeString(ciphertextBase64)
	if err != nil {
		return "", err
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	if len(ciphertext) < nonceSize {
		return "", err
	}

	nonce := ciphertext[:nonceSize]
	encryptedData := ciphertext[nonceSize:]

	// Extract timestamp
	var timestampStr string
	var dataWithoutTimestamp []byte
	for i, b := range encryptedData {
		if b < '0' || b > '9' {
			timestampStr = string(encryptedData[:i])
			dataWithoutTimestamp = encryptedData[i:]
			break
		}
		if i == len(encryptedData)-1 {
			return "", err
		}
	}

	timestampInt, err := strconv.ParseInt(timestampStr, 10, 64)
	if err != nil {
		return "", err
	}
	messageTime := time.Unix(0, timestampInt).UTC()
	now := time.Now().UTC()

	// Check for replay attack based on timestamp
	if now.Sub(messageTime) > timeTolerance {
		return "", err
	}
	if messageTime.After(now.Add(timeTolerance)) {
		return "", err
	}

	plaintextBytes, err := gcm.Open(nil, nonce, dataWithoutTimestamp, nil)
	if err != nil {
		return "", err
	}

	return string(plaintextBytes), nil
}