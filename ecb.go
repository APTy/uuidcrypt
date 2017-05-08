package main

import "crypto/cipher"

type ECB struct {
	block cipher.Block
}

func NewECB(block cipher.Block) *ECB {
	return &ECB{
		block: block,
	}
}

func (c *ECB) Encrypt(plaintext []byte) []byte {
	bs := c.block.BlockSize()
	if len(plaintext)%bs != 0 {
		panic("invalid plaintext size")
	}
	ciphertext := make([]byte, len(plaintext))
	c.block.Encrypt(ciphertext, plaintext)
	return ciphertext
}

func (c *ECB) Decrypt(ciphertext []byte) []byte {
	bs := c.block.BlockSize()
	if len(ciphertext)%bs != 0 {
		panic("invalid ciphertext size")
	}
	plaintext := make([]byte, len(ciphertext))
	c.block.Decrypt(plaintext, ciphertext)
	return plaintext
}
