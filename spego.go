package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v2"
)

// Config - struct to store configuration
type Config struct {
	PayloadPath string `yaml:"PayloadPath"`
	OutputPath  string `yaml:"OutputPath"`
	StubPath    string `yaml:"StubPath"`
	Magic       string `yaml:"Magic"`
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
	log.Printf("Writing delimiter to %s", config.OutputPath)
	if _, err := f.Write([]byte(config.Magic)); err != nil {
		log.Println("Unable to write delimiter.")
		log.Fatalln(err)
	}

	// write payload
	log.Printf("Writing payload to %s", config.OutputPath)
	if _, err := f.Write(payload); err != nil {
		log.Println("Unable to write payload.")
		log.Fatalln(err)
	}

	// finished
	log.Printf("Complete. Output written to: %s\n", config.OutputPath)

}
