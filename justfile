export CGO_ENABLED := "0"

build:
    go build -ldflags="-s -w" .

build-version VERSION:
    go build -ldflags="-s -w -X main.Version="{{ VERSION }}"" .

env:
    go env

install:
    @just build
    mv ./walle ~/.local/bin/

install-sys:
    @just build
    sudo mv ./walle /usr/local/bin/
