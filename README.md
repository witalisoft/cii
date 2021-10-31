# cii - container image info

cii is a tool to provide basic information about the image without pulling it locally. The information contains the number of layers, available platforms with sizes (compressed versions) and detailed layer history. History option is similar to what `docker history` offers but extended with shell formatter if applicable.  The tool is based on [crane](https://github.com/google/go-containerregistry/tree/main/cmd/crane).

## Build from source
Project is using Go modules, just build:
```console
$ go build
```

## Usage
```console
$ ./cii -h
Usage of ./cii:
  -image string
        image name
  -no-color
        disable color output
  -no-format
        don't try to format shell scripts
  -platform string
        specify platform of which you want to get layers history in the form os/arch (default "linux/amd64")
  -version
        print the version and exit
```

## Examples
```console
$ ./cii -image postgres:latest
Available platforms:

         platform           size                                                                         digest
      linux/amd64       130.6MiB        sha256:047d261a65cd42d974cfa3e4ce2462e32b9856a1b4616bfcd04aeac8a9e35482
     linux/arm/v5       124.2MiB        sha256:40d981463bd01563bc5af77e2cf179776316483a28f4a588b1cbbcfb72726586
     linux/arm/v7       119.2MiB        sha256:6be1c085a40cfc9907583f3c82d0cad939883e37317667b30d777b31d7d8dba9
   linux/arm64/v8       125.5MiB        sha256:092709339ae4cfddc3459406cfb36fd81b23da8907bac56235d93e6ed161c092
        linux/386       132.6MiB        sha256:19a83c9999859536d94ffc4979867023de714ba796227b8310068e21399cd019
   linux/mips64le       125.6MiB        sha256:1f8c4bc490b6b0fe3a0ca5846ceef79b82f3258549e8a423dd0d293ddcfb8f31
    linux/ppc64le       139.1MiB        sha256:22b98acd86bfeced06064b1c49839b9e0c060358a7a66b2e3f1c3d98a485c211
      linux/s390x       134.3MiB        sha256:ff97afc7a1a86a843c7d22191d2f4c9b6e8f19ddffd7f194da5cda30296e7cc1


Data layers: 13
Empty layers: 12
Last pushed: 4 days


Layers history for platform: linux/amd64

layer: 1
size: 29.9MiB
empty_layer: false
created: 2 weeks
created_by: /bin/sh -c #(nop) ADD file:16dc2c6d1932194edec28d730b004fd6deca3d0f0e1a07bc5b8b6e8a1662f7af in / 

layer: 2
size: 0B
empty_layer: true
created: 2 weeks
created_by: /bin/sh -c #(nop)  CMD ["bash"]

layer: 3
size: 4.21MiB
empty_layer: false
created: 2 weeks
created_by: |
        /bin/sh -c set -ex
        if ! command -v gpg >/dev/null; then
                apt-get update
                apt-get install -y --no-install-recommends gnupg dirmngr
                rm -rf /var/lib/apt/lists/*
        fi
...
```
