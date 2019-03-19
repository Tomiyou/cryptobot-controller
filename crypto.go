package main

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/md5"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"io"
	"io/ioutil"
	"os"

	"github.com/spf13/cobra"
)

var encryptCmd = &cobra.Command{
	Use:   "encrypt",
	Short: "Encrypt crypto-arbitrage secrets.",
	RunE:  encryptSecretsCmd,
}

var decryptCmd = &cobra.Command{
	Use:   "decrypt",
	Short: "Decrypt crypto-arbitrage secrets.",
	RunE:  decryptSecretsCmd,
}

func encryptSecretsCmd(cmd *cobra.Command, args []string) (err error) {
	// parse key
	if len(args) == 0 {
		return fmt.Errorf("A key is required to run this command.")
	}
	key := args[0]

	// create the tar file to be encrypted
	tarFile, err := createTarFile("keys")
	if err != nil {
		return
	}
	defer tarFile.Close()

	// get bytes from tar file
	data, err := ioutil.ReadAll(tarFile)

	// encrypt the data
	encrypted, err := encrypt(data, key)
	if err != nil {
		return
	}

	// create target file
	target, err := os.Create(botConfig.EncryptedSecretsPath)
	if err != nil {
		return
	}
	defer target.Close()

	// write encrypted bytes to target file
	_, err = target.Write(encrypted)

	return
}

func decryptSecretsCmd(cmd *cobra.Command, args []string) (err error) {
	// parse key
	if len(args) == 0 {
		return fmt.Errorf("A key is required to run this command.")
	}
	key := args[0]

	// read the encrypted file
	data, err := ioutil.ReadFile(botConfig.EncryptedSecretsPath)
	if err != nil {
		return
	}

	// decrypt the data
	decrypted, err := decrypt(data, key)
	if err != nil {
		return
	}

	// now write it to temporary tar file
	tmp, err := os.Create(botConfig.TemporaryTarPath)
	if err != nil {
		return
	}
	defer tmp.Close()
	tmp.Write(decrypted)

	// extract the tar file to the current folder
	err = extractTarFile(".")

	return
}

func createHash(key string) string {
	hasher := md5.New()
	hasher.Write([]byte(key))
	return hex.EncodeToString(hasher.Sum(nil))
}

func encrypt(data []byte, passphrase string) ([]byte, error) {
	block, _ := aes.NewCipher([]byte(createHash(passphrase)))
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}
	nonce := make([]byte, gcm.NonceSize())
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, err
	}
	ciphertext := gcm.Seal(nonce, nonce, data, nil)
	return ciphertext, nil
}

func decrypt(data []byte, passphrase string) ([]byte, error) {
	key := []byte(createHash(passphrase))
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}
	nonceSize := gcm.NonceSize()
	nonce, ciphertext := data[:nonceSize], data[nonceSize:]
	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return nil, err
	}
	return plaintext, nil
}
