package main

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v2"
)

// Config - struct to store configuration
type Config struct {
	PayloadPath string            `yaml:"PayloadPath"`
	OutputPath  string            `yaml:"OutputPath"`
	StubPath    string            `yaml:"StubPath"`
	Magic       string            `yaml:"Magic"`
	Password    string            `yaml:"Password"`
	Envs        map[string]string `yaml:"Envs"`
}

func main() {
	if len(os.Args) != 2 {
		fmt.Printf("Usage:\n\t.\\%s config.yaml\n\n", filepath.Base(os.Args[0]))
		log.Fatal("Wrong number of arguments")
	}

	// read in yaml config file
	configBytes, err := ioutil.ReadFile(os.Args[1])
	if err != nil {
		log.Fatalf("Unable to read config file: %v", err)
	}

	config := Config{}
	if err := yaml.Unmarshal(configBytes, &config); err != nil {
		log.Fatalf("Unable to parse config file: %v", err)
	}

	fmt.Printf("Using payload path of %s\n", config.PayloadPath)

	// Read payload
	log.Printf("Reading payload %s", config.PayloadPath)
	payload, err := ioutil.ReadFile(config.PayloadPath)
	if err != nil {
		log.Println("Unable to read payload.")
		log.Fatalln(err)
	}

	// Read stub
	log.Printf("Reading stub %s", config.StubPath)
	stub, err := ioutil.ReadFile(config.StubPath)
	if err != nil {
		log.Println("Unable to read stub file.")
		log.Fatalln(err)
	}

	// create encryption key
	log.Println("Creating encryption key")
	encKey := fmt.Sprintf("password:%s", config.Password)
	envKeys := "" // env vars used to create encryption key
	for k, v := range config.Envs {
		if v != "" {
			encKey = fmt.Sprintf("%s:%s:%s", encKey, strings.ToLower(k), strings.ToLower(v))
			if envKeys == "" {
				envKeys = k
			} else {
				envKeys = fmt.Sprintf("%s:%s", envKeys, k)
			}
		}
	}
	encKeySha256 := sha256.Sum256([]byte(encKey))

	// generate aes cipher black
	c, err := aes.NewCipher(encKeySha256[:])
	if err != nil {
		log.Println("Unable to create AES cipher")
		log.Fatalln(err)
	}

	// create GCM operater
	gcm, err := cipher.NewGCM(c)
	if err != nil {
		log.Println("Unable to create GCM")
		log.Fatalln(err)
	}

	// create nonce
	nonce := make([]byte, gcm.NonceSize())
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		log.Println("Unable to create nonce.")
		log.Fatalln(err)
	}

	// encrypt payload
	log.Println("Encrypting payload")
	encryptedPayload := gcm.Seal(nonce, nonce, payload, nil)

	// open output file
	f, err := os.Create(config.OutputPath)
	if err != nil {
		log.Println("Unable to create output file.")
		log.Fatalln(err)
	}
	defer f.Close()

	// write stub
	log.Printf("Writing stub to %s", config.OutputPath)
	if _, err := f.Write(stub); err != nil {
		log.Println("Unable to write stub.")
		log.Fatalln(err)
	}

	// write magic string
	log.Printf("Writing delimiter and env keys to %s", config.OutputPath)
	if _, err := f.Write([]byte(config.Magic)); err != nil {
		log.Println("Unable to write first delimiter.")
		log.Fatalln(err)
	}

	if _, err := f.Write([]byte(envKeys)); err != nil {
		log.Println("Unable to write env keys.")
		log.Fatalln(err)
	}

	if _, err := f.Write([]byte(config.Magic)); err != nil {
		log.Println("Unable to write second delimiter.")
		log.Fatalln(err)
	}

	// write payload
	log.Printf("Writing payload to %s", config.OutputPath)
	if _, err := f.Write(encryptedPayload); err != nil {
		log.Println("Unable to write payload.")
		log.Fatalln(err)
	}

	// finished
	log.Printf("Complete. Output written to: %s\n", config.OutputPath)

}
