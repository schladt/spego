# SEPGO
Simple Payload Encrypter in GO

## Features
- Portable, single command payload packaging and encryption
- Builds encrypted binaries from executable and raw shellcode payloads
- Uses environment variable to build encryption key. (Keeps samples from executing outside of scoped targets)
- No re-building or compiler toolchain required for payload build
- Cross platform support (builds Windows payloads from any environment)

## Quickstart:
- Download binary from latest release. Choose binary based on system running spego, not target system. All binarys can produce cross platform encrypted executables
- Edit example-config.yml for custom payload, environment variable key values, and other options. (Example is well commented)
- Run spego.exe
```
Usage:
        .\spego-windows-amd64.exe config.yaml
```
- Enjoy!

## Special Notes
- Executable payloads must be compiled with the relocation table intact. Most Microsoft compilers do this by default (MS Visual Studio). However, special flags are needed for payloads using GNU tools (see Makefile in the sample payload directories). Additionally, Go only recently added support for position independent executables (buildmode=pie) in their master branch (expected next release). 
- SPEGO includes an option for a commandline execution password with the -password flag. However, this flag is passed on to the child process so it should not be used for payloads that expect their own command line arguments. In these cases, it's recommended to set a custom environment variable (such as SPEGOPASS in the example-config).

## Build Process

Building is NOT required as all options are configured in the yaml config (unless of course you don't trust random executables on GitHub :) ). Building the main program requires the Go toolchain. Building the stubs requires the MinGW toolchain (including make) with both x86_64-w64-mingw32-gcc and i686-w64-mingw32-gcc set in the path.

- Three step process made easy with Makefiles
  1. Build program stubs
  ```
  cd stubs
  make
  cd ..
  ```
  2. Update bindata
  ```
  go-bindata stubs/bin/spego-stub-windows-386.exe stubs/bin/spego-stub-windows-amd64.exe
  ```
  3. Build spego main encrypter
  ```
  make
  ```
- Enjoy!

## How it works
- Demo coming soon

## References and Ack
This project is built upon and extends the following work:
- MemoryModule - https://github.com/fancycode/MemoryModule
- Ebowla - https://github.com/Genetic-Malware/Ebowla
- Go-Mimikatz https://github.com/vyrus001/go-mimikatz
- Go-Shellcode https://github.com/brimstone/go-shellcode
