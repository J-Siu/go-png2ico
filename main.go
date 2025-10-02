/*
The MIT License

Copyright (c) 2025 John Siu

Permission is hereby granted, free of charge, to any person obtaining a copy of this software and associated documentation files (the "Software"), to deal in the Software without restriction, including without limitation the rights to use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of the Software, and to permit persons to whom the Software is furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
*/

package main

import (
	"errors"
	"fmt"
	"os"

	"github.com/J-Siu/go-helper/v2/errs"
	"github.com/J-Siu/go-helper/v2/ezlog"
	"github.com/J-Siu/go-png2ico/p2i"
)

const version = "v1.1.0"

func usage() {
	fmt.Println("go-png2ico version " + version)
	fmt.Println("License : MIT License Copyright (c) 2025 John Siu")
	fmt.Println("Support : https://github.com/J-Siu/go-png2ico/issues")
	fmt.Println("Debug   : export _DEBUG=true")
	fmt.Println("Usage   : go-png2ico <PNG file> <PNG file> ... <ICO file>")
}

func main() {
	prefix := "main"
	// ARGs
	args := os.Args[1:]
	argc := len(args)
	switch argc {
	case 0:
		usage()
		os.Exit(0)
	case 1:
		errs.Queue(prefix, errors.New("Input/Output file missing"))
	}

	fileout := args[argc-1]

	if errs.IsEmpty() {
		// Make sure destination file is *not* PNG
		png := new(p2i.PNG).New().Read(fileout)
		if png.IsPNG() {
			errs.Queue(prefix, errors.New(png.File+": is PNG"))
		}
	}
	ezlog.SetLogLevel(ezlog.DEBUG)
	ico := new(p2i.ICO).New(fileout)

	pngc := argc - 1
	for i := 0; i < pngc; i++ {
		ico.AddPngFile(args[i])
	}

	ico.WriteAll()

	ezlog.Debug().Nn("ico").Mn(ico).N("png count").M(ico.PngCount()).Out()

	if !errs.IsEmpty() {
		ezlog.Err().Ln().M(errs.Errs).Out()
	}
}
