package main

import (
	"bytes"
	"testing"
)

func TestCipher_Encrypt(t *testing.T) {
	fileBuffer := &bytes.Buffer{}

	cph, err := MakeCipher("test", fileBuffer)
	if err != nil {
		t.Fatal(err)
	}

	_, err = cph.Encrypt("hello\n")
	if err != nil {
		t.Fatal(err)
	}

	_, err = cph.Encrypt("world\n")
	if err != nil {
		t.Fatal(err)
	}

	cph, err = MakeCipher("test", fileBuffer)
	if err != nil {
		t.Fatal(err)
	}

	outBuffer := &bytes.Buffer{}

	if _, err := cph.Decrypt(outBuffer); err != nil {
		t.Fatal(err)
	}

	if outBuffer.String() != "hello\nworld\n" {
		t.Error("wrong decrypted message")
		t.Fatal(outBuffer.String())
	}
}
