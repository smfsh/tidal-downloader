package main

import (
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"fmt"
	"io"
	"os"
)

const MASTER_KEY string = "UIlTTEMmmLfGowo/UC60x2H45W6MdGgTRfo/umg4754="

func decryptToken(token string) ([]byte, []byte) {
	dMaster, err := base64.StdEncoding.DecodeString(MASTER_KEY)
	if err != nil {
		panic(err)
	}
	fmt.Printf("%08b", dMaster)
	dToken, err := base64.StdEncoding.DecodeString(token)
	if err != nil {
		panic(err)
	}

	// Initialization Vector.
	iv := dToken[:16]
	// Cipher Text.
	ct := dToken[16:]

	block, err := aes.NewCipher(dMaster)
	if err != nil {
		panic(err)
	}

	cbc := cipher.NewCBCDecrypter(block, iv)
	cbc.CryptBlocks(ct, ct)

	println(len(ct))
	fmt.Println(string(ct))
	key := ct[:16]
	nonce := ct[16:24]
	fmt.Println(string(key))
	fmt.Println(string(nonce))

	return key, nonce

}

func decryptFile(encFile string, decFile string, key []byte, nonce []byte) {
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

	block, err := aes.NewCipher(key)
	if err != nil {
		panic(err)
	}

	ctr := cipher.NewCTR(block, nonce)

	buf := make([]byte, 4096)
	for {
		n, err := inFile.Read(buf)
		if err != nil && err != io.EOF {
			panic(err)
		}

		outBuf := make([]byte, n)
		ctr.XORKeyStream(outBuf, buf[:n])
		outFile.Write(outBuf)

		if err == io.EOF {
			break
		}
	}
}
