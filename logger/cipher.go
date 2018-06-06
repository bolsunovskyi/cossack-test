package main

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/sha256"
	"errors"
	"io"
	"os"
)

type Cipher struct {
	writer     cipher.StreamWriter
	reader     cipher.StreamReader
	targetFile ReadWriteSyncer
}

func (c *Cipher) Encrypt(message string) (n int, err error) {
	return c.writer.Write([]byte(message))
}

func (c *Cipher) Write(p []byte) (n int, err error) {
	return c.writer.Write(p)
}

func (c *Cipher) Decrypt(out io.ReadWriter) (n int64, err error) {
	return io.Copy(out, c.reader)
}

func (c *Cipher) Sync() error {
	return c.targetFile.Sync()
}

func (c *Cipher) Close() error {
	return c.targetFile.Close()
}

func MakeCipher(key string, targetReaderWriter *os.File) (*Cipher, error) {
	if key == "" {
		return nil, errors.New("key is empty")
	}

	cph := Cipher{}

	hash := sha256.New()
	hash.Write([]byte(key))

	keyHash := hash.Sum(nil)
	block, err := aes.NewCipher(keyHash)
	if err != nil {
		return nil, err
	}

	iv := keyHash[:aes.BlockSize]

	enc := cipher.NewCFBEncrypter(block, iv)
	dec := cipher.NewCFBDecrypter(block, iv)

	cph.writer = cipher.StreamWriter{S: enc, W: targetReaderWriter}
	cph.reader = cipher.StreamReader{S: dec, R: targetReaderWriter}

	cph.targetFile = targetReaderWriter

	return &cph, nil
}
