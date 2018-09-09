# goserve

**goserve** is a distributable HTTP server for your HTMLs.

It is just a small single file with no configuration but gives your HTMLs the power of `http:` protocol instead of `file:` protocol.

Place a [**goserve** binary](bin/goserve-mac) into your directory and run to show your HTMLs under the directory.

BEFORE (without **goserve**):

```diff
 example/
 ├── index.html
 └── stub.json
```

Run `open index.html`, you will see:

![](doc/before.png)

AFTER (with **goserve** on macOS):

```diff
 example/
+├── goserve-mac
 ├── index.html
 └── stub.json
```

Run `./goserve-mac & open http://localhost:8080`, you will see:

![](doc/after.png)

(After that, do not forget to find the pid with `ps | grep goserve` and shutdown the server with `kill <found_pid>`.)
