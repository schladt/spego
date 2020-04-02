package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
)

func main() {
	outPath := "bin/out.exe"
	payloadPath := "Z:\\projects\\spego\\payloads\\go-example\\go-example.exe"
	stubPath := "Z:\\projects\\spego\\stubs\\bin\\spego-stub-windows-amd64.exe"

	// Read payload
	fmt.Println("Reading payload...")
	payload, err := ioutil.ReadFile(payloadPath)
	if err != nil {
		log.Fatalln(err)
	}

	// Read stub
	fmt.Println("Reading stub...")
	stub, err := ioutil.ReadFile(stubPath)
	if err != nil {
		log.Fatalln(err)
	}

	// magic string
	magic := []byte("9BC5440033354F2EBEEED2E6083903390A39B15489892A199F86A17CFA0F55B8")

	// Open output file
	f, err := os.Create(outPath)
	if err != nil {
		log.Fatalln(err)
	}
	defer f.Close()

	// write stub
	fmt.Println("Writing stub...")
	if _, err := f.Write(stub); err != nil {
		log.Fatalln(err)
	}

	// write magic string
	fmt.Println("Writing delimiter...")
	if _, err := f.Write(magic); err != nil {
		log.Fatalln(err)
	}

	// write payload
	fmt.Println("Writing payload...")
	if _, err := f.Write(payload); err != nil {
		log.Fatalln(err)
	}

	// finished
	fmt.Printf("Complete. Output written to: %s\n", outPath)

}
