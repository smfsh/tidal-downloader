package main

import (
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"io"
	"os"
)

const MasterKey string = "UIlTTEMmmLfGowo/UC60x2H45W6MdGgTRfo/umg4754="

func decryptToken(token string) ([]byte, []byte) {
	dMaster, err := base64.StdEncoding.DecodeString(MasterKey)
	if err != nil {
		panic(err)
	}
	// Token decodes to a 48 byte array.
	dToken, err := base64.StdEncoding.DecodeString(token)
	if err != nil {
		panic(err)
	}

	// Initialization Vector for AES CBC encrypted cipher text.
	iv := dToken[:16]
	// Encrypted cipher text containing block encryption key and nonce to decrypt file.
	ct := dToken[16:]

	// Initialize new cipher block utilizing Tidal master key.
	block, err := aes.NewCipher(dMaster)
	if err != nil {
		panic(err)
	}

	// Initialize decrypter utilizing init vector from encrypted token (cipher text).
	// https://en.wikipedia.org/wiki/Block_cipher_mode_of_operation#Cipher_Block_Chaining_(CBC)
	cbc := cipher.NewCBCDecrypter(block, iv)
	// Decrypt and overwrite the encrypted object with a decrypted one.
	cbc.CryptBlocks(ct, ct)

	// Cipher text is 32 bytes long, the remaining 8 bytes is thrown away.
	key := ct[:16]
	ctriv := ct[16:24]

	return key, ctriv

}

func decryptFile(encFile string, decFile string, key []byte, iv []byte) {
	inFile, err := os.Open(encFile)
	if err != nil {
		panic(err)
	}
	defer inFile.Close()

	outFile, err := os.Create(decFile)
	if err != nil {
		panic(err)
	}
	defer outFile.Close()

	// Initialize a new cipher block utilizing decrypted key from token.
	block, err := aes.NewCipher(key)
	if err != nil {
		panic(err)
	}

	// Initialization vector for CTR decryption must be exactly the same
	// length as the block size being decrypted. Check for this difference
	// and add padding to the end of the byte array if needed.
	if block.BlockSize()-len(iv) > 0 {
		pad := make([]byte, block.BlockSize()-len(iv))
		iv = append(iv, pad...)
	}

	// Initialize decrypter utilizing decrypted nonce from token.
	// https://en.wikipedia.org/wiki/Block_cipher_mode_of_operation#Counter_(CTR)
	ctr := cipher.NewCTR(block, iv)

	// Create a 4k byte buffer to write out the decrypted file in chunks.
	buf := make([]byte, 4096)
	for {
		n, err := inFile.Read(buf)
		if err != nil && err != io.EOF {
			panic(err)
		}

		outBuf := make([]byte, n)
		// Perform decryption using CTR rotation on the current chunk.
		ctr.XORKeyStream(outBuf, buf[:n])
		outFile.Write(outBuf)

		if err == io.EOF {
			break
		}
	}
}
