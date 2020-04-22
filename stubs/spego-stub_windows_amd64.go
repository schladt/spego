package main

/*
#cgo CFLAGS: -IMemoryModule
#cgo LDFLAGS: MemoryModule/build/amd64/MemoryModule.a
#include "MemoryModule/MemoryModule.h"
*/
import "C"

import (
	"flag"
	"fmt"
	"os"
	"unsafe"
)

func launchExe(payload []byte) {

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
