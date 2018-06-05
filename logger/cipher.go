package main

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"io"
)

type Cipher struct {
	stream cipher.Stream
}

func (c *Cipher) pad(src []byte) []byte {
	padding := aes.BlockSize - len(src)%aes.BlockSize
	paddedText := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(src, paddedText...)
}

func (c *Cipher) Encrypt(message string) ([]byte, error) {
	plainBytes := c.pad([]byte(message))
	cipherText := make([]byte, aes.BlockSize+len(plainBytes))
	iv := cipherText[:aes.BlockSize]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return nil, err
	}

	c.stream.XORKeyStream(cipherText[aes.BlockSize:], plainBytes)
	return cipherText, nil
}

func (c *Cipher) Decrypt(cipherText []byte) (string, error) {
	return "", nil
}

func MakeCipher(key string) (*Cipher, error) {
	cph := Cipher{}

	hash := sha256.New()
	hash.Write([]byte(key))

	keyHash := hash.Sum(nil)
	block, err := aes.NewCipher(keyHash)
	if err != nil {
		return nil, err
	}

	iv := make([]byte, aes.BlockSize)
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return nil, err
	}

	cph.stream = cipher.NewCFBEncrypter(block, iv)

	return &cph, nil
}
