# https://qiita.com/harukasan/items/37698ec799678c12e71d
VERSION := $(shell git describe --tags HEAD 2>/dev/null || git rev-parse --short HEAD 2>/dev/null)

# https://qiita.com/dtan4/items/8c417b629b6b2033a541#%E5%A4%89%E6%95%B0%E5%AE%9A%E7%BE%A9
LDFLAGS := -ldflags="-s -w -X \"main.Version=$(VERSION)\" -extldflags \"-static\""

# Run `make` to do cross compile.
all: bin/goserve bin/goserve-mac bin/goserve.exe bin/goserve-32.exe

# For Linux 64-bit, macOS 64-bit, Windows 32/64-bit.
bin/goserve:
	GOOS=linux   GOARCH=amd64 go build $(LDFLAGS) -o bin/goserve
bin/goserve-mac:
	GOOS=darwin  GOARCH=amd64 go build $(LDFLAGS) -o bin/goserve-mac
bin/goserve.exe:
	GOOS=windows GOARCH=amd64 go build $(LDFLAGS) -o bin/goserve.exe
bin/goserve-32.exe:
	GOOS=windows GOARCH=386   go build $(LDFLAGS) -o bin/goserve-32.exe

# Run `make clean` to force updating artifacts.
clean:
	rm -rf bin/*
