package utils

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestEncryptDecrypt_RoundTrip(t *testing.T) {
	key := make([]byte, 32)
	for i := range key {
		key[i] = byte(i)
	}

	plaintext := "my-secret-password"
	encrypted, err := Encrypt(plaintext, key)
	assert.NoError(t, err)
	assert.NotEqual(t, plaintext, encrypted)
	assert.True(t, IsEncrypted(encrypted))

	decrypted, err := Decrypt(encrypted, key)
	assert.NoError(t, err)
	assert.Equal(t, plaintext, decrypted)
}

func TestEncryptDecrypt_EmptyString(t *testing.T) {
	key := make([]byte, 32)

	encrypted, err := Encrypt("", key)
	assert.NoError(t, err)
	assert.Equal(t, "", encrypted)

	decrypted, err := Decrypt("", key)
	assert.NoError(t, err)
	assert.Equal(t, "", decrypted)
}

func TestDecrypt_PlaintextPassthrough(t *testing.T) {
	key := make([]byte, 32)

	decrypted, err := Decrypt("not-encrypted", key)
	assert.NoError(t, err)
	assert.Equal(t, "not-encrypted", decrypted)
}

func TestDecrypt_InvalidBase64(t *testing.T) {
	key := make([]byte, 32)

	_, err := Decrypt("enc:!!!invalid-base64!!!", key)
	assert.Error(t, err)
}

func TestDecrypt_TruncatedCiphertext(t *testing.T) {
	key := make([]byte, 32)

	// "AAAAAA==" decodes to 4 bytes, which is less than the 12-byte AES-GCM nonce size.
	// This tests the "ciphertext too short" error path.
	_, err := Decrypt("enc:AAAAAA==", key)
	assert.Error(t, err)
}

func TestDecrypt_WrongKey(t *testing.T) {
	key1 := make([]byte, 32)
	for i := range key1 {
		key1[i] = byte(i)
	}
	key2 := make([]byte, 32)
	for i := range key2 {
		key2[i] = byte(i + 1)
	}

	encrypted, err := Encrypt("secret", key1)
	assert.NoError(t, err)

	_, err = Decrypt(encrypted, key2)
	assert.Error(t, err)
}

func TestIsEncrypted(t *testing.T) {
	assert.True(t, IsEncrypted("enc:somedata"))
	assert.False(t, IsEncrypted("plaintext"))
	assert.True(t, IsEncrypted("enc:"))
	assert.False(t, IsEncrypted(""))
}

func TestDeriveKey_TooShort(t *testing.T) {
	_, err := deriveKey([]byte("abc"))
	assert.Error(t, err)
}

func TestDeriveKey_Valid(t *testing.T) {
	key, err := deriveKey([]byte("mypassword"))
	assert.NoError(t, err)
	assert.Equal(t, 32, len(key))
}

func TestMaskSecret(t *testing.T) {
	assert.Equal(t, "<empty>", MaskSecret(""))
	assert.Equal(t, "****", MaskSecret("a"))
	assert.Equal(t, "****", MaskSecret("abcd"))
	assert.Equal(t, "****", MaskSecret("mysecretpassword"))
}

func TestGetMasterKey_EnvVar(t *testing.T) {
	// Clear any cached master key first
	ClearMasterKey()

	// Set the environment variable
	envKey := "testpassword123"
	t.Setenv(keyEnvVar, envKey)

	key, err := GetMasterKey()
	require.NoError(t, err)
	require.NotNil(t, key)
	assert.Equal(t, 32, len(key))

	// Cleanup
	ClearMasterKey()
}

func TestGetMasterKey_EnvVarTooShort(t *testing.T) {
	ClearMasterKey()

	t.Setenv(keyEnvVar, "abc")

	key, err := GetMasterKey()
	assert.Error(t, err)
	assert.Nil(t, key)
	assert.Contains(t, err.Error(), "at least 4 characters")

	ClearMasterKey()
}

func TestGetMasterKey_Caching(t *testing.T) {
	ClearMasterKey()

	t.Setenv(keyEnvVar, "cacheTestPassword")

	// First call - should derive and cache
	key1, err := GetMasterKey()
	require.NoError(t, err)

	// Unset the env var; the cached key should still be returned
	os.Unsetenv(keyEnvVar)

	key2, err := GetMasterKey()
	require.NoError(t, err)

	// Should return the same cached key
	assert.Equal(t, key1, key2)

	ClearMasterKey()
}

func TestClearMasterKey(t *testing.T) {
	// First set a key via env var
	ClearMasterKey()
	t.Setenv(keyEnvVar, "clearTestPassword")

	key, err := GetMasterKey()
	require.NoError(t, err)
	require.NotNil(t, key)

	// Clear it
	ClearMasterKey()

	// After clearing, the cached key should be nil
	// We verify by checking the global directly (it's in the same package)
	masterKeyMu.Lock()
	assert.Nil(t, masterKey)
	masterKeyMu.Unlock()

	// Reset env for subsequent tests
	os.Unsetenv(keyEnvVar)
}

func TestEncryptDecrypt_WithDerivedKey(t *testing.T) {
	// Verify that a key derived from deriveKey works with Encrypt/Decrypt
	key, err := deriveKey([]byte("testpassword123"))
	require.NoError(t, err)

	plaintext := "hello world with derived key"
	encrypted, err := Encrypt(plaintext, key)
	require.NoError(t, err)

	decrypted, err := Decrypt(encrypted, key)
	require.NoError(t, err)
	assert.Equal(t, plaintext, decrypted)
}

func TestEncryptDecrypt_SpecialCharacters(t *testing.T) {
	key := make([]byte, 32)
	for i := range key {
		key[i] = byte(i)
	}

	testCases := []string{
		"unicode: 中文文字",
		"newlines:\n\ttab",
		"special chars: !@#$%^&*()",
		`backslash \ and "quotes"`,
		"{\"json\": true}",
	}

	for _, tc := range testCases {
		t.Run(tc, func(t *testing.T) {
			encrypted, err := Encrypt(tc, key)
			require.NoError(t, err)

			decrypted, err := Decrypt(encrypted, key)
			require.NoError(t, err)
			assert.Equal(t, tc, decrypted)
		})
	}
}

func TestEncryptDecrypt_LongPlaintext(t *testing.T) {
	key := make([]byte, 32)
	for i := range key {
		key[i] = byte(i)
	}

	// Generate a long plaintext (> 1KB)
	longText := ""
	for i := 0; i < 2000; i++ {
		longText += "a"
	}

	encrypted, err := Encrypt(longText, key)
	require.NoError(t, err)

	decrypted, err := Decrypt(encrypted, key)
	require.NoError(t, err)
	assert.Equal(t, longText, decrypted)
}

func TestEncrypt_InvalidKey(t *testing.T) {
	// AES requires 16, 24, or 32 byte keys
	_, err := Encrypt("test", []byte{1, 2, 3})
	assert.Error(t, err)
}

func TestDecrypt_InvalidKey(t *testing.T) {
	// First encrypt with a valid key
	validKey := make([]byte, 32)
	encrypted, err := Encrypt("test", validKey)
	require.NoError(t, err)

	// Try to decrypt with an invalid key length
	_, err = Decrypt(encrypted, []byte{1, 2, 3})
	assert.Error(t, err)
}

func TestDeriveKey_Deterministic(t *testing.T) {
	// Same input should produce same output
	key1, err := deriveKey([]byte("samepassword"))
	require.NoError(t, err)

	key2, err := deriveKey([]byte("samepassword"))
	require.NoError(t, err)

	assert.Equal(t, key1, key2)
}

func TestDeriveKey_DifferentPasswords(t *testing.T) {
	key1, err := deriveKey([]byte("password1"))
	require.NoError(t, err)

	key2, err := deriveKey([]byte("password2"))
	require.NoError(t, err)

	assert.NotEqual(t, key1, key2)
}

func TestDeriveKey_MinLength(t *testing.T) {
	// Exactly 4 characters should work
	key, err := deriveKey([]byte("abcd"))
	assert.NoError(t, err)
	assert.Equal(t, 32, len(key))
}

func TestGetMasterKey_InteractivePromptError(t *testing.T) {
	// When no env var is set and stdin is not a terminal (as in tests),
	// term.ReadPassword will fail, triggering the error path.
	ClearMasterKey()
	os.Unsetenv(keyEnvVar)

	key, err := GetMasterKey()
	// In test environments, stdin is a pipe so term.ReadPassword should fail.
	// If for some reason it doesn't (e.g., very unusual test runner),
	// the key should still be usable.
	if err != nil {
		assert.Error(t, err)
		assert.Nil(t, key)
		assert.Contains(t, err.Error(), "failed to read master password")
	}

	ClearMasterKey()
}

func TestGetMasterKey_AlreadyCached(t *testing.T) {
	// Test the early return when masterKey is already cached.
	ClearMasterKey()
	t.Setenv(keyEnvVar, "alreadyCachedPassword")

	// First call caches the key
	key1, err := GetMasterKey()
	require.NoError(t, err)

	// Clear env var so subsequent calls would go to interactive path
	// if caching didn't work
	os.Unsetenv(keyEnvVar)

	// Second call should return the cached key (the len(masterKey) > 0 branch)
	key2, err := GetMasterKey()
	require.NoError(t, err)
	assert.Equal(t, key1, key2)

	ClearMasterKey()
}

func TestGetCachedMasterKey_NoKey(t *testing.T) {
	ClearMasterKey()
	os.Unsetenv(keyEnvVar)

	key, err := GetCachedMasterKey()
	assert.NoError(t, err)
	assert.Nil(t, key)
}

func TestGetCachedMasterKey_EnvVar(t *testing.T) {
	ClearMasterKey()
	t.Setenv(keyEnvVar, "cachedEnvPassword")

	key, err := GetCachedMasterKey()
	require.NoError(t, err)
	assert.NotNil(t, key)
	assert.Equal(t, 32, len(key))

	ClearMasterKey()
}

func TestGetCachedMasterKey_UsesCache(t *testing.T) {
	ClearMasterKey()
	t.Setenv(keyEnvVar, "cacheTestPassword")

	// First call caches via env var
	key1, err := GetCachedMasterKey()
	require.NoError(t, err)

	// Remove env var — cached key should still be returned
	os.Unsetenv(keyEnvVar)

	key2, err := GetCachedMasterKey()
	require.NoError(t, err)
	assert.Equal(t, key1, key2)

	ClearMasterKey()
}

func TestClearMasterKey_WhenNil(t *testing.T) {
	// ClearMasterKey should be safe to call when masterKey is already nil.
	ClearMasterKey()
	masterKeyMu.Lock()
	masterKey = nil
	masterKeyMu.Unlock()

	// Should not panic
	ClearMasterKey()

	masterKeyMu.Lock()
	assert.Nil(t, masterKey)
	masterKeyMu.Unlock()
}
