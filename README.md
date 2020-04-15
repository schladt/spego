# spego
SPEGO - Simple Payload Encrypter in GO

## Features
- Single command payload packaging and encryption
- Uses environment variable to build encryption key. (Keeps samples from executing outside of scoped targets)
- No re-building or compiler toolchain required for payload build
- Cross platform support (builds Windows payloads from any environment)

## Quickstart:
- Download binaries from latest release
- Edit example-config.yml for custom payload, environment variable key values, and other options. (Example is well commented)
- Run spego.exe
```
Usage:
        .\spego-windows-amd64.exe config.yaml
```
- Enjoy!

## Special Notes
- Payloads must be compiled with the relocation table intact. Most MS compilers do this by default (MS Visual Studio). However, special flags are needed for payloads using GNU tools (see Makefile in the sample payload directories). Additionally, Go only recently added support for position independent executables (buildmode=pie) in their master branch (expected next release). 
- SPEGO includes an option for a commandline execution password with the -password flag. However, this flag is pass on to the child process so it should not be used for payloads that expect their own command line arguments. In these cases, it's recommended to set a custom environment variable (such as SPEGOPASS in the example-config).

## Build Process
- Two step process
-- Build program stubs
-- Build spego main encrypter
- Mostly documented in Makefiles
- Msys Toolchain recommended 
- More info coming soon

## How it works
- Coming soon

## References and Ack
This project is built upon and extends the following work:
- MemoryModule - https://github.com/fancycode/MemoryModule
- Ebowla - https://github.com/Genetic-Malware/Ebowla
- Go-Mimikatz https://github.com/vyrus001/go-mimikatz
