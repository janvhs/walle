build:
    go build .

install:
    @just build
    mv ./walle ~/.local/bin/

install-sys:
    @just build
    sudo mv ./walle /usr/local/bin/
