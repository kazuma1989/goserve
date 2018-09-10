# https://qiita.com/dtan4/items/8c417b629b6b2033a541#%E5%A4%89%E6%95%B0%E5%AE%9A%E7%BE%A9
LDFLAGS := -ldflags="-s -w -extldflags \"-static\""

# Run `make` to do cross compile.
all: bin/goserve bin/goserve-mac bin/goserve.exe

# For 64-bit OSs only.
bin/goserve:
	GOOS=linux   GOARCH=amd64 go build $(LDFLAGS) -o bin/goserve
bin/goserve-mac:
	GOOS=darwin  GOARCH=amd64 go build $(LDFLAGS) -o bin/goserve-mac
bin/goserve.exe:
	GOOS=windows GOARCH=amd64 go build $(LDFLAGS) -o bin/goserve.exe

# Run `make clean` to update.
clean:
	rm -rf bin/*
