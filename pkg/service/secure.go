package service

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"

	"golang.org/x/crypto/pbkdf2"
)

const (
	saltlen    = 32
	keylen     = 32
	iterations = 100002
)

// returns ciphertext of the following format:
// [32 bit salt][128 bit iv][encrypted plaintext]
func EncryptPrivateKey(plaintext string, password string) string {
	// allocate memory to hold the header of the ciphertext
	header := make([]byte, saltlen+aes.BlockSize)

	// generate salt
	salt := header[:saltlen]
	if _, err := io.ReadFull(rand.Reader, salt); err != nil {
		panic(err)
	}

	// generate initialization vector
	iv := header[saltlen : aes.BlockSize+saltlen]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		panic(err)
	}

	// generate a 32 bit key with the provided password
	key := pbkdf2.Key([]byte(password), salt, iterations, keylen, sha256.New)

	// generate a hmac for the message with the key
	mac := hmac.New(sha256.New, key)
	mac.Write([]byte(plaintext))
	hmac := mac.Sum(nil)

	// append this hmac to the plaintext
	plaintext = string(hmac) + plaintext

	//create the cipher
	block, err := aes.NewCipher(key)
	if err != nil {
		panic(err)
	}

	// allocate space for the ciphertext and write the header to it
	ciphertext := make([]byte, len(header)+len(plaintext))
	copy(ciphertext, header)

	// encrypt
	stream := cipher.NewCFBEncrypter(block, iv)
	stream.XORKeyStream(ciphertext[aes.BlockSize+saltlen:], []byte(plaintext))
	return hex.EncodeToString(ciphertext)
}

func DecryptPrivateKey(encrypted string, password string) string {
	ciphertext, err := hex.DecodeString(encrypted)
	if err != nil {
		panic(err)
	}
	// get the salt from the ciphertext
	salt := ciphertext[:saltlen]
	// get the IV from the ciphertext
	iv := ciphertext[saltlen : aes.BlockSize+saltlen]
	// generate the key with the KDF
	key := pbkdf2.Key([]byte(password), salt, iterations, keylen, sha256.New)

	block, err := aes.NewCipher(key)
	if err != nil {
		panic(err)
	}

	if len(ciphertext) < aes.BlockSize {
		return ""
	}

	decrypted := ciphertext[saltlen+aes.BlockSize:]
	stream := cipher.NewCFBDecrypter(block, iv)
	stream.XORKeyStream(decrypted, decrypted)

	// extract hmac from plaintext
	extractedMac := decrypted[:32]
	plaintext := decrypted[32:]

	// validate the hmac
	mac := hmac.New(sha256.New, key)
	mac.Write(plaintext)
	expectedMac := mac.Sum(nil)
	if !hmac.Equal(extractedMac, expectedMac) {
		return ""
	}

	return string(plaintext)
}

func GenerateSecureToken(length int) string {
	b := make([]byte, length)
	if _, err := rand.Read(b); err != nil {
		return ""
	}
	return hex.EncodeToString(b)
}

func GeneratePasswordHash(password string) string {
	hash := sha256.New()
	hash.Write([]byte(password))
	return fmt.Sprintf("%x", hash.Sum([]byte(salt)))
}
