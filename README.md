# goserve

**goserve** is a distributable HTTP server for your HTMLs.

It is just a small single file with no configuration but gives your HTMLs the power of `http:` protocol instead of `file:` protocol.



## Basic usage

Place a [**goserve** binary](bin/goserve-mac) into your directory and run to show your HTMLs under the directory.

BEFORE (without **goserve**):

```diff
 example/
 ├── index.html
 └── stub.json
```

Run `open index.html`, you will see:

![](doc/before.png)

AFTER (with **goserve** for macOS):

```diff
 example/
+├── goserve-mac
 ├── index.html
 └── stub.json
```

Run `./goserve-mac`, you will see:

![](doc/after.png)



## Routing (mapping between a URL and a file path)

Add `routes.json` in the directory where **goserve** exists.

```diff
 example/
 ├── goserve-mac
 ├── index.html
+├── routes.json
 └── stub.json
```

```json
{
  "/api/message": "stub.json"
}
```

You can see the content of `stub.json` at `http://localhost:8080/api/message` in addition to `http://localhost:8080/stub.json`.




# Development

## Requirements

To build **goserve** from the source, you need:

- [The Go Programming Language](https://golang.org/)
- [Make - GNU Project - Free Software Foundation](https://www.gnu.org/software/make/)
- [UPX: the Ultimate Packer for eXecutables - Homepage](https://upx.github.io/)
