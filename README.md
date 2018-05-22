# Libretro Cheat Fix

Parses and fixes multiline libretro cheat files.

#### Example:
The following
```
cheat0_desc = "Infinite Lives"
cheat0_code = "AAAAAA+BBBBBB"
cheat0_enable = true
```

Becomes
```
cheat0_desc = "Infinite Lives"
cheat0_code = "AAAAAA"
cheat0_enable = true

cheat1_desc = "Infinite Lives"
cheat1_code = "AAAAAA"
cheat1_enable = true
```

## Getting Started

[Download](https://golang.org/dl/) and install golang.
Clone into your $GOPATH/src directory.
Run `go build` in the libretro-chtfix directory.
Symlink or add $GOPATH/src/libretro-chtfix to your $PATH variable.
Run `libretro-chtfix -in=/path/to/cheat/file.cht`

CLI Options:
**in:** A single file to parse and fix.
**out:** Optional output file name. Defaults to the file specified using the `in` paramater with *_fixed* added to the end of the file name.
**dir:** Will scan and parse all cheat files in a directory.

__**You must supply either `in` or `dir`**__

### Prerequisites

Golang > 1.8

## License

This project is licensed under the MIT License - see the [LICENSE.md](LICENSE.md) file for details
