package main

import "crypto/cipher"

// ECB performs encryption and decryption in the Electronic Codebook
// operation mode. It is not fit for operating on many blocks of data
// but can be used on fixed-size, single-block length data.
//
// As such, it is intended to only be used to encrypt the byte
// representation of UUIDs.
type ECB struct {
	block cipher.Block
}

// NewECB returns a new ECB object using the provided cipher block.
func NewECB(block cipher.Block) *ECB {
	return &ECB{
		block: block,
	}
}

// Encrypt encrypts a block-sized length plaintext into ciphertext.
func (c *ECB) Encrypt(plaintext []byte) []byte {
	bs := c.block.BlockSize()
	if len(plaintext)%bs != 0 {
		panic("invalid plaintext size")
	}
	ciphertext := make([]byte, len(plaintext))
	c.block.Encrypt(ciphertext, plaintext)
	return ciphertext
}

// Decrypt decrypts a block-sized length ciphertext into plaintext.
func (c *ECB) Decrypt(ciphertext []byte) []byte {
	bs := c.block.BlockSize()
	if len(ciphertext)%bs != 0 {
		panic("invalid ciphertext size")
	}
	plaintext := make([]byte, len(ciphertext))
	c.block.Decrypt(plaintext, ciphertext)
	return plaintext
}
