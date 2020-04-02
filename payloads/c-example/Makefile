UNAME := $(shell uname)

ifeq ($(UNAME), Linux)
ifndef PLATFORM
PLATFORM = x86_64
endif
CC = $(PLATFORM)-w64-mingw32-gcc
CXX = $(PLATFORM)-w64-mingw32-g++
LINK = $(PLATFORM)-w64-mingw32-ld
else
CC = gcc
CXX = g++
LINK = ld
endif

CFLAGS  = -Wall -g
LDFLAGS = -static

all:
	$(CC) $(LDFLAGS) -Wl,--dynamicbase,--export-all-symbols c-example.c -o c-example.exe