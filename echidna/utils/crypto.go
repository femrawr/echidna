package utils

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"echidna/store"
	"io"
	"math/big"
)

var chars = []rune("qwertyuiopasdfghjklzxcvbnm")

func GetRandomString(length int) string {
	buffer := make([]rune, length)
	for i := range buffer {
		num, _ := rand.Int(rand.Reader, big.NewInt(int64(len(chars))))
		buffer[i] = chars[num.Int64()]
	}

	return string(buffer)
}

func EncryptData(data string) []byte {
	if !store.OUTPUT_ENCRYPTED {
		return []byte(data)
	}

	key := sha256.Sum256(store.ENCRYPTION_KEY)

	block, err := aes.NewCipher(key[:])
	if err != nil {
		return nil
	}

	result := make([]byte, aes.BlockSize+len(data))

	iv := result[:aes.BlockSize]
	_, err = io.ReadFull(rand.Reader, iv)
	if err != nil {
		return nil
	}

	stream := cipher.NewCTR(block, iv)
	stream.XORKeyStream(result[aes.BlockSize:], []byte(data))

	return result
}
