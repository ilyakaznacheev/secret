package handler

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"
	"io"
)

// generateKey generates a random encryption key
func generateKey() string {
	nonce := make([]byte, 12)
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		panic(err.Error())
	}
	return fmt.Sprintf("%x", nonce)
}

// encryptSecret encrypts value with AES using key
func encryptSecret(key, value string) (string, error) {
	rawKey := []byte(key)
	rawText := []byte(value)

	block, err := aes.NewCipher(rawKey)
	if err != nil {
		return "", err
	}

	//IV needs to be unique, but doesn't have to be secure.
	//It's common to put it at the beginning of the ciphertext.
	cipherText := make([]byte, aes.BlockSize+len(rawText))
	iv := cipherText[:aes.BlockSize]
	if _, err = io.ReadFull(rand.Reader, iv); err != nil {
		return "", err
	}

	stream := cipher.NewCFBEncrypter(block, iv)
	stream.XORKeyStream(cipherText[aes.BlockSize:], rawText)

	//returns to base64 encoded string
	encText := base64.StdEncoding.EncodeToString(cipherText)
	return encText, nil
}

// decryptSecret decrypts value with AES using key
func decryptSecret(key, value string) (string, error) {
	rawKey := []byte(key)

	rawText, err := base64.StdEncoding.DecodeString(value)
	if err != nil {
		return "", err
	}

	block, err := aes.NewCipher(rawKey)
	if err != nil {
		return "", err
	}

	if len(rawText) < aes.BlockSize {
		err = errors.New("ciphertext block size is too short")
		return "", err
	}

	//IV needs to be unique, but doesn't have to be secure.
	//It's common to put it at the beginning of the ciphertext.
	iv := rawText[:aes.BlockSize]
	rawText = rawText[aes.BlockSize:]

	stream := cipher.NewCFBDecrypter(block, iv)
	// XORKeyStream can work in-place if the two arguments are the same.
	stream.XORKeyStream(rawText, rawText)

	decText := string(rawText)
	return decText, nil
}
