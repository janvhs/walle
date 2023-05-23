# walle aka WALL·E

Keep your computer tidy with this little helper 🤖

```sh
walle ~/Projects
```

## What can *walle* do for me?

*walle* can crawl your hard drive and look for directories that
contain build artefacts, such as rust's *target* directory or
npm's *node_modules* directory.

These are usually very large and contain files that are only used
during development and can therefore be safely deleted.

## How do I install *walle*?

These instructions are for Linux and Mac.

> **Note**
> Please make sure that `~/.local/bin` is on your `$PATH`.

> **Warning**
> On Windows, adjust these to suit your needs.

### Via [just](https://just.systems)

```sh
git clone https://github.com/bode-fun/walle.git
cd walle
just install
```

### Manual

```sh
git clone https://github.com/bode-fun/walle.git
cd walle
go build .
mv walle ~/.local/bin/
```

## Which programming languages are currently supported?

The languages currently supported by *walle* are listed below.

However, adding a new Language is as easy as adding a configuration
in [main.go](main.go)

> **Note**
> If your language is missing, please add it in a pull request
instead of opening an issue.

- JavaScript
  - npm
- Python
  - Virtual environments aka venv
  - \_\_pycache\_\_
- Rust
  - Cargo
- php
  - Composer
- Swift
  - Swift Package Manager

## TODOs

For a list of open TODOs, please take  a look at the `// TODO:` comments in the source ✌️

Please contribute, if you have any ideas for cool features or
just want to improve something.

## One more thing

Software is more than just bits and bytes. It should be elegant,
easy to read, correct, maintainable and working on it
should teach you something.

One can only truly fulfil this goal, if one truly thinks about,
understands and knows about the source code one writes.

In my opinion, this is only possible by writing the source code
yourself, paraphrasing source code from external sources or,
to some extent, generating source code via a compiler.

Therefore, this software is **100%** handmade.

Built with 🫶 and 💅 by [Jan Fooken](https://github.com/bode-fun).
Licensed under [GPLv3](LICENSE)
