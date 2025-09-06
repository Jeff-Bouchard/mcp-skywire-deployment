# Dmsgget

`dmsgget` is a utility exec which can download from HTTP servers hosted over the `dmsg` network (similar to a simplified `wget` over `dmsg`).

```
$ dmsgget --help

  Skycoin dmsgget v0.1.0, wget over dmsg.
  Usage: dmsgget [OPTION]... [URL]
  
    -O FILE
          write documents to FILE (default ".")
    -U AGENT
          identify as AGENT (default "dmsgget/v0.1.0")
    -dmsg-disc URL
          dmsg discovery URL (default "http://dmsgd.skywire.skycoin.com")
    -dmsg-sessions NUMBER
          connect to NUMBER of dmsg servers (default 1)
    -h    
    -help
          print this help
    -t NUMBER
          set number of retries to NUMBER (0 unlimits) (default 1)
    -w SECONDS
          wait SECONDS between retrievals
```

### Example usage

In this example, we will use the `dmsg` network where the `dmsg.Discovery` address is `http://dmsgd.skywire.skycoin.com`. However, any `dmsg.Discovery` would work.

First, lets create a folder where we will host files to serve over `dmsg` and create a `hello.txt` file within.

```shell script
// Create serving folder.
$ mkdir /tmp/dmsghttp -p

// Create file.
$ echo 'Hello World!' > /tmp/dmsghttp/hello.txt
```

Next, let's serve this over `http` via `dmsg` as transport. We have an example exec for this located within `/example/dmsgget/dmsg-example-http-server`.

```shell script
# Generate public/private key pair
$ go run ./examples/dmsgget/gen-keys/gen-keys.go
#   PK: 038dde2d050803db59e2ad19e5a6db0f58f8419709fc65041c48b0cb209bb7a851
#   SK: e5740e093bd472c2730b0a58944a5dee220d415de62acf45d1c559f56eea2b2d

# Run dmsg http server.
#   (replace 'e5740e093bd472c2730b0a58944a5dee220d415de62acf45d1c559f56eea2b2d' with the SK returned from above command)
$ go run ./examples/dmsgget/dmsg-example-http-server/dmsg-example-http-server.go --dir /tmp/dmsghttp --sk e5740e093bd472c2730b0a58944a5dee220d415de62acf45d1c559f56eea2b2d
```

Now we can use `dsmgget` to download the hosted file. Open a new terminal and run the following.

```shell script
# Replace '038dde2d050803db59e2ad19e5a6db0f58f8419709fc65041c48b0cb209bb7a851' with the generated PK.
$ dmsgget dmsg://038dde2d050803db59e2ad19e5a6db0f58f8419709fc65041c48b0cb209bb7a851:80/hello.txt

# Check downloaded file.
$ cat hello.txt
#   Hello World!
```

