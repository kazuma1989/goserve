# Run `make` to do cross compile.
all: bin/goserve bin/goserve-mac bin/goserve.exe

# For 64-bit OSs only.
bin/goserve:
	GOOS=linux GOARCH=amd64 go build -o bin/goserve
bin/goserve-mac:
	GOOS=darwin GOARCH=amd64 go build -o bin/goserve-mac
bin/goserve.exe:
	GOOS=windows GOARCH=amd64 go build -o bin/goserve.exe

# Run `make clean` to update.
clean:
	rm -rf bin/*
