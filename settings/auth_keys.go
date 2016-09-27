package settings

import (
	"bufio"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"log"
	"os"
)

func GetPrivateKey() *rsa.PrivateKey {

	fileData := getKeyData("private_key")

	privateKey, err := x509.ParsePKCS1PrivateKey(fileData.Bytes)
	if err != nil {
		log.Fatal(err)
	}

	return privateKey
}

func GetPublicKey() *rsa.PublicKey {

	fileData := getKeyData("public_key.pub")

	publicKey, err := x509.ParsePKIXPublicKey(fileData.Bytes)
	if err != nil {
		log.Fatal(err)
	}

	rsaPublicKey, ok := publicKey.(*rsa.PublicKey)
	if !ok {
		log.Fatalf("Failed to cast decoded public key to type *rsa.Public")
	}

	return rsaPublicKey
}

func getKeyData(fileName string) *pem.Block {

	file, err := os.Open("/home/dipen/go/src/github.com/patelndipen/AP1/settings/" + os.Getenv("GO_ENV") + "/" + fileName)
	if err != nil {
		log.Fatal(err)
	}

	defer file.Close()

	fileInfo, err := file.Stat()
	if err != nil {
		log.Fatal(err)
	}

	fileBytes := make([]byte, fileInfo.Size())

	_, err = bufio.NewReader(file).Read(fileBytes)
	if err != nil {
		log.Fatal(err)
	}

	fileData, _ := pem.Decode([]byte(fileBytes))
	if fileData == nil {
		log.Fatalf("Failed to decode PEM format")
	}

	return fileData

}
