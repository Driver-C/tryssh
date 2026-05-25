// Package utils provides common utilities for the tryssh application.
package utils

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"io"
	"os"
	"strings"
	"sync"

	"golang.org/x/term"
)

const (
	encryptedPrefix = "enc:"
	keyEnvVar       = "TRYSSH_MASTER_KEY"
	kdfIter         = 100000
)

var (
	masterKey   []byte
	masterKeyMu sync.Mutex
)

// GetMasterKey returns the cached master key, prompting for it if necessary.
func GetMasterKey() ([]byte, error) {
	masterKeyMu.Lock()
	defer masterKeyMu.Unlock()

	if len(masterKey) > 0 {
		return masterKey, nil
	}

	// Try environment variable first
	if envKey := os.Getenv(keyEnvVar); envKey != "" {
		key, err := deriveKey([]byte(envKey))
		if err != nil {
			return nil, err
		}
		masterKey = key
		return masterKey, nil
	}

	// Prompt interactively
	fmt.Print("Enter master password: ")
	pwdBytes, err := term.ReadPassword(int(os.Stdin.Fd()))
	fmt.Println()
	if err != nil {
		return nil, fmt.Errorf("failed to read master password: %w", err)
	}
	if len(pwdBytes) == 0 {
		return nil, nil
	}

	key, err := deriveKey(pwdBytes)
	for i := range pwdBytes {
		pwdBytes[i] = 0
	}
	if err != nil {
		return nil, err
	}
	masterKey = key
	return masterKey, nil
}

// ClearMasterKey removes the cached master key from memory.
func ClearMasterKey() {
	masterKeyMu.Lock()
	defer masterKeyMu.Unlock()
	for i := range masterKey {
		masterKey[i] = 0
	}
	masterKey = nil
}

// deriveKey derives a 32-byte AES key from a password using iterated HMAC-SHA256.
func deriveKey(password []byte) ([]byte, error) {
	if len(password) < 4 {
		return nil, fmt.Errorf("master password must be at least 4 characters")
	}
	salt := []byte("tryssh-config-v1")
	key := make([]byte, 32)
	h := hmac.New(sha256.New, password)
	h.Write(salt)
	copy(key, h.Sum(nil))
	for i := 0; i < kdfIter; i++ {
		h = hmac.New(sha256.New, password)
		h.Write(key)
		key = h.Sum(key[:0])
	}
	return key, nil
}

// Encrypt encrypts plaintext and returns a prefixed base64 string.
func Encrypt(plaintext string, key []byte) (string, error) {
	if plaintext == "" {
		return "", nil
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	nonce := make([]byte, aesGCM.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return "", err
	}

	ciphertext := aesGCM.Seal(nonce, nonce, []byte(plaintext), nil)
	return encryptedPrefix + base64.StdEncoding.EncodeToString(ciphertext), nil
}

// Decrypt decrypts a prefixed base64 string back to plaintext.
func Decrypt(encrypted string, key []byte) (string, error) {
	if encrypted == "" {
		return "", nil
	}
	if !IsEncrypted(encrypted) {
		return encrypted, nil
	}

	data, err := base64.StdEncoding.DecodeString(encrypted[len(encryptedPrefix):])
	if err != nil {
		return "", fmt.Errorf("failed to decode encrypted value: %w", err)
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	nonceSize := aesGCM.NonceSize()
	if len(data) < nonceSize {
		return "", fmt.Errorf("ciphertext too short")
	}

	nonce, ciphertext := data[:nonceSize], data[nonceSize:]
	plaintext, err := aesGCM.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return "", fmt.Errorf("decryption failed: %w", err)
	}

	return string(plaintext), nil
}

// IsEncrypted checks if a value has the encrypted prefix.
func IsEncrypted(s string) bool {
	return strings.HasPrefix(s, encryptedPrefix)
}
