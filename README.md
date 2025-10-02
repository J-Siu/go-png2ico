# go-png2ico [![Paypal donate](https://www.paypalobjects.com/en_US/i/btn/btn_donate_LG.gif)](https://www.paypal.com/donate/?business=HZF49NM9D35SJ&no_recurring=0&currency_code=CAD)

Command line tool to create ICO(favicon) from PNG images.

### Table Of Content
<!-- TOC -->

- [Module](#module)
  - [Install](#install)
  - [Import](#import)
  - [ICO Usage](#ico-usage)
- [Command Line](#command-line)
  - [What It Does](#what-it-does)
  - [What It Does Not](#what-it-does-not)
  - [Limitation](#limitation)
- [Binary](#binary)
- [Compile](#compile)
- [Usage](#usage)
- [Repository](#repository)
- [Contributors](#contributors)
- [Change Log](#change-log)
- [License](#license)

<!-- /TOC -->
<!--more-->

### Module

#### Install

```sh
go get github.com/J-Siu/go-png2ico/v2
```

#### Import

```go
import "github.com/J-Siu/go-png2ico/v2/p2i"
```

#### ICO Usage

```go
ico := new(p2i.ICO).New(icoFile)
ico.AddPngFile(pngFile) // Can be repeated
ico.WriteAll()
```

Full example in [root.go](/cmd/root.go)

### Command Line

#### What It Does

- Create ICO file from PNG files
- ICO use PNG format for storage
- Minimum overhead(16byte) per PNG added
- PNG header check for input files
- PNG header check for output file to avoid mistake

#### What It Does Not

- Change PNG to BMP inside ICO
- Check file extension
- Transform PNG

#### Limitation

- ICO file always created from scratch
- No append nor replace within existing ICO file
- PNG to ICO only, other format/conversion not supported

### Binary

https://github.com/J-Siu/go-png2ico/releases

### Compile

```sh
go get github.com/J-Siu/go-png2ico
cd $GOPATH/src/github.com/J-Siu/go-png2ico
go install
```

### Usage

```sh
go-png2ico -h
```

```sh
Build ICO file from PNGs

Usage:
  go-png2ico <PNG file> <PNG file> ... <ICO file> [flags]

Flags:
  -d, --debug     Enable debug
  -h, --help      help for go-png2ico
  -v, --verbose   Verbose
      --version   version for go-png2ico
```

### Repository

- [go-png2ico](https://github.com/J-Siu/go-png2ico)

### Contributors

- [John Sing Dao Siu](https://github.com/J-Siu)

### Change Log

- 1.0
  - Initial Commit
- 1.0.1
  - Fix
    - debug log msg
    - error check
    - png detection
- 1.0.2
  - Use mod
- 1.0.3
  - Use github.com/J-Siu/go-helper
- 1.0.4
  - Use Go 1.16
- v1.0.5
  - Fix `goreleaser`
- v1.0.6
  - Update to Go 1.20 and dependency
- v1.0.7
  - Support PNG with height and width above 256
- v1.0.8
  - Fix Github workflows
- v2.0.0
  - expose ICO and PNG as package p2i
  - move to go-helper/v2
  - use cobra for cli flag
- v2.0.1
  - resolve conflict

### License

The MIT License

Copyright (c) 2025 John Siu

Permission is hereby granted, free of charge, to any person obtaining a copy of this software and associated documentation files (the "Software"), to deal in the Software without restriction, including without limitation the rights to use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of the Software, and to permit persons to whom the Software is furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
