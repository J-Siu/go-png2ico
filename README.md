# go-png2ico [![Paypal donate](https://www.paypalobjects.com/en_US/i/btn/btn_donate_LG.gif)](https://www.paypal.com/donate/?business=HZF49NM9D35SJ&no_recurring=0&currency_code=CAD)

Command line tool to create ICO(favicon) from PNG images.

### Table Of Content
<!-- TOC -->

- [Table Of Content](#table-of-content)
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
### What It Does

- Create ICO file from PNG files
- ICO use PNG format for storage
- Minimum overhead(16byte) per PNG added
- PNG header check for input files
- PNG header check for output file to avoid mistake

### What It Does Not

- Change PNG to BMP inside ICO
- Check file extension
- Transform PNG

### Limitation

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
go-png2ico
```

```sh
go-png2ico MIT License Copyright (c) 2023 John Siu
Support: https://github.com/J-Siu/go-png2ico/issues
Usage: go-png2ico <PNG file> <PNG file> ... <ICO file>
Debug: export _DEBUG=true
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

### License

The MIT License

Copyright (c) 2023 John Siu

Permission is hereby granted, free of charge, to any person obtaining a copy of this software and associated documentation files (the "Software"), to deal in the Software without restriction, including without limitation the rights to use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of the Software, and to permit persons to whom the Software is furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.