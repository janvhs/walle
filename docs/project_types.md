# Projects in different languages

- [Projects in different languages](#projects-in-different-languages)
  - [Supported languages](#supported-languages)
  - [Possible identifier and Target](#possible-identifier-and-target)
  - [Types of identifiers](#types-of-identifiers)
  - [Language-specific commands](#language-specific-commands)

## Supported languages

- [x] Javascript/Typescript
- [x] Rust
- [ ] Swift
  - [x] SwiftPM
  - [ ] Xcode

## Possible identifier and Target

| Language    | Identifier                       | Target             |
| ----------- | -------------------------------- | ------------------ |
| Javascript  | package.json                     | node_modules       |
| Rust        | Cargo.toml                       | target             |
| Swift (PM)  | Package.swift                    | .build             |
| Python      | \_\_pycache\_\_/\*.pyc oder *.py | \_\_pycache\_\_    |
| Python Venv | \*/pyvenv.cfg                    | dir(\*/pyvenv.cfg) |
| PHP         | composer.json                    | vendor             |

## Types of identifiers

- Exact file name e.g. `Cargo.toml`
- File in specific directory e.g. `__pycache__/*.pyc`
- File in target directory e.g. `*/pyvenv.cfg`
- File with specific extension `*.py`

## Language-specific commands

- `go clean -cache` for Go (any directory). Cleans global build cache
- `pipenv --rm` for Python in project root. Removes virtual environment
- `cargo clean` for Rust in project root. Removes local target directory