export CGO_ENABLED := "0"

build:
    go build .

env:
    go env

install:
    @just build
    mv ./walle ~/.local/bin/

install-sys:
    @just build
    sudo mv ./walle /usr/local/bin/
