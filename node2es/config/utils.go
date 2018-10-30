package config

import (
	"encoding/base64"

	"github.com/hongxincn/promexp/node2es/common"
)

func GetDecryptedPassword(encryptedPassword string) (string, error) {
	if encryptedPassword == "" {
		return "", nil
	}
	aesEnc := common.AesEncrypt{Key: ENCRYPT_KEY}
	rawBytes, err := base64.StdEncoding.DecodeString(encryptedPassword)
	if err != nil {
		return "", err
	}

	decryptedPassword, err := aesEnc.Decrypt(rawBytes)
	if err != nil {
		return "", err
	}
	return decryptedPassword, nil
}
