package main

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/hmac"
	"crypto/md5"
)

type CryptType int

const (
	EncryptType CryptType = 1 + iota
	DecryptType
)

type Processor interface {
	Process([]byte) []byte
}

func NewCrypterProcessor(secret, namespace []byte, cryptType CryptType) Processor {
	key := keyGen(secret, namespace)
	block, err := aes.NewCipher(key)
	if err != nil {
		panic(err)
	}
	cipher, err := cipher.NewGCM(block)
	if err != nil {
		panic(err)
	}
	return &crypterProcessor{
		key:       key,
		cipher:    cipher,
		cryptType: cryptType,
	}
}

type crypterProcessor struct {
	key       []byte
	cipher    cipher.AEAD
	cryptType CryptType
	nonce     []byte
}

func keyGen(secret, namespace []byte) []byte {
	mac := hmac.New(md5.New, secret)
	mac.Write(namespace)
	return mac.Sum(nil)
}

func (p *crypterProcessor) Process(in []byte) []byte {
	switch p.cryptType {
	case EncryptType:
		return p.encrypt(in)
	case DecryptType:
		return p.decrypt(in)
	default:
	}
	return in
}

func (p *crypterProcessor) encrypt(in []byte) []byte {
	out := make([]byte, len(in))
	for i := range in {
		out[i] = in[i] ^ p.key[i]
	}
	return out
}

func (p *crypterProcessor) decrypt(in []byte) []byte {
	out := make([]byte, len(in))
	for i := range in {
		out[i] = in[i] ^ p.key[i]
	}
	return out
}
