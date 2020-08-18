/*
The MIT License

Copyright (c) 2020 John Siu

Permission is hereby granted, free of charge, to any person obtaining a copy of this software and associated documentation files (the "Software"), to deal in the Software without restriction, including without limitation the rights to use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of the Software, and to permit persons to whom the Software is furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
*/

package main

import (
	"bytes"
	"encoding/binary"
	"encoding/hex"
	"errors"
	"fmt"
	"os"
	"strconv"

	"github.com/J-Siu/go-png2ico/helper"
)

// ICO structire
type ICO struct {
	file string
	fh   *os.File
}

// PNG structure
type PNG struct {
	file   string // filename
	fh     *os.File
	height uint8
	width  uint8
	depth  uint16 // bit/pixel
	size   uint32
	offset uint32
	isPNG  bool
}

func log(msg ...interface{}) {
	if helper.Debug {
		fmt.Println("log:", msg)
	}
}

func errCheck(e error) {
	if e != nil {
		fmt.Println(e)
		os.Exit(1)
	}
}

// Open : open PNG file
func (png *PNG) Open(file string) error {
	log("PNG:Open:", file)

	var e error
	var fh *os.File
	var n int

	fh, e = os.Open(file)
	if e != nil {
		return e
	}

	png.file = file
	png.fh = fh
	png.isPNG = false

	/* PNG header check(25byte)
	89 50 4e 47 0d 0a 1a 0a // 8byte - magic number
	// IHDR chunk
	xx xx xx xx // 4byte - chunk length
	49 48 44 52 // 4byte - chunk type(IHDR)
	xx xx xx xx // 4byte - width
	xx xx xx xx // 4byte - height
	xx          // 1byte - bit death (bit/pixel)
	*/
	headerLen := 25
	header := make([]byte, headerLen)
	n, e = png.fh.Read(header)
	if e != nil {
		return e
	}
	log("PNG:Open:Header:", hex.EncodeToString(header), "(", n, ")")

	// 8byte header[0:8] - magic number
	magic := []byte{0x89, 0x50, 0x4e, 0x47, 0x0d, 0x0a, 0x1a, 0x0a}
	if bytes.Equal(magic[:], header[:8]) {
		log("Found PNG magic.")
	} else {
		return errors.New("Not PNG")
	}
	// 4byte header[8:12] - chunk length - skipped
	// 4byte header[12:16] - chunk type IHDR
	chunk := []byte{0x49, 0x48, 0x44, 0x52}
	if bytes.Equal(chunk[:], header[12:16]) {
		log("Found IHDR chunk.")
	} else {
		return errors.New("PNG no IHDR chunk")
	}

	// It is PNG
	png.isPNG = true

	// 4byte header[16:20] - width
	width := binary.BigEndian.Uint32(header[16:20])

	// 4byte header[20:24] - height
	height := binary.BigEndian.Uint32(header[20:24])

	if width <= 256 && height <= 256 {
		log("PNG:Open:", width, "x", height)
		// ICO format use 0 for 256px
		if width == 256 {
			width = 0
		}
		if height == 256 {
			height = 0
		}
	} else {
		return errors.New(png.file + "(" + strconv.FormatUint(uint64(width), 10) + "x" + strconv.FormatUint(uint64(height), 10) + "): Width and height cannot be larger than 256.")
	}
	png.width = uint8(width)
	png.height = uint8(height)

	// 1byte header[25] - color depth
	png.depth = uint16(uint8(header[24]))
	log("PNG:Open:Bit/pixel:", png.depth)

	stat, _ := os.Stat(file)
	png.size = uint32(stat.Size())
	log("PNG:Open:Size:", png.size)

	// Pass all check, create PNG struct
	log("PNG:Open:PNG:", *png)

	return nil
}

// Read : read PNG file
func (png *PNG) Read() *[]byte {
	log("PNG:Read:", png.file)

	var e error
	var n int
	var n64 int64

	n64, e = png.fh.Seek(0, 0)
	errCheck(e)
	log("PNG:Read:FH reset:", n64)

	b := make([]byte, png.size)
	n, e = png.fh.Read(b)
	errCheck(e)

	log("PNG:Read:byte", n)

	return &b
}

// Open : open ICO filehandle
func (ico *ICO) Open(file string) {
	var e error
	log("ICO:Open:" + file)
	ico.fh, e = os.Create(file)
	errCheck(e)
}

// Write : write ICO
func (ico *ICO) Write(b *[]byte) {
	var e error
	var n int
	n, e = ico.fh.Write(*b)
	errCheck(e)
	log("ICO:Write:byte:", n)
}

// ICONDIR - return ICONDIR byte array
func (ico *ICO) ICONDIR(num uint16) *[]byte {
	/*
		6byte ICONDIR structure
		00:   00 00 // 2byte, must be 0
		02:   01 00 // 2byte, 1 for ICO
		04:   xx xx // 2byte, img number
	*/

	b := []byte{0, 0, 1, 0, 0, 0}
	binary.LittleEndian.PutUint16(b[4:6], num)
	log("ICO:ICONDIR:", hex.EncodeToString(b))
	return &b
}

// ICONDIRENTRY - return ICONDIRENTRY byte array
func (png *PNG) ICONDIRENTRY() *[]byte {
	log("ICONDIRENTRY:", *png)
	/*
		16byte ICONDIRENTRY structure
		00:   xx    // 1byte, width
		01:   xx    // 1byte, height
		02:   00    // 1byte, color palette number, 0 for PNG
		03:   00    // 1byte, reserved, always 0
		04:   00 00 // 2byte, color planes, 0 for PNG
		06:   xx xx // 2byte, color depth
		08:   xx xx xx xx // 4byte, image size
		12:   xx xx xx xx // 4byte, image offset
	*/

	b := make([]byte, 16)

	copy(b[0:6], []byte{png.width, png.height, 0, 0, 0, 0})
	binary.LittleEndian.PutUint16(b[6:8], png.depth)
	binary.LittleEndian.PutUint32(b[8:12], png.size)
	binary.LittleEndian.PutUint32(b[12:16], png.offset)
	log("PNG:ICONDIRENTRY:", hex.EncodeToString(b))

	return &b
}

func usage() {
	fmt.Println("go-png2ico MIT License	Copyright (c) 2020 John Siu")
	fmt.Println("Support: https://github.com/J-Siu/go-png2ico/issues")
	fmt.Println("Usage: go-png2ico <PNG file> <PNG file> ... <ICO file>")
	fmt.Println("Debug: export _DEBUG=true")
}

func main() {
	var e error

	//Debug
	if os.Getenv("_DEBUG") == "true" {
		helper.Debug = true
	}

	// ARGs
	args := os.Args[1:]
	argc := len(args)
	switch argc {
	case 0:
		usage()
		os.Exit(0)
	case 1:
		e = errors.New("Input/Output file missing")
	}
	errCheck(e)

	fileout := args[argc-1]

	// Make sure destination file is not PNG
	png := new(PNG)
	if png.Open(fileout) == nil || png.isPNG {
		e = errors.New("Output file (" + png.file + ") is a PNG file.")
	} else {
		e = nil
	}
	errCheck(e)

	// Get and calculate all PNGs info
	pngs := []*PNG{}
	pngc := argc - 1
	var pngTotalSize uint32 = 0
	var LenICONDIR uint32 = 6
	var LenICONDIRENTRY uint32 = 16
	var LenAllICONDIRENTRY uint32 = LenICONDIRENTRY * uint32(pngc)
	for i := 0; i < pngc; i++ {
		png := new(PNG)
		e = png.Open(args[i])
		errCheck(e)
		// offset = len(ICONDIR) + len(all ICONDIRENTRY) + len(all PNG before current one)
		png.offset = LenICONDIR + LenAllICONDIRENTRY + pngTotalSize
		pngs = append(pngs, png)
		pngTotalSize += png.size
	}

	// Open ICON
	ico := new(ICO)
	ico.Open(fileout)
	ico.Write(ico.ICONDIR(uint16(pngc)))
	// Write ICONDIRENTRY
	for i := 0; i < pngc; i++ {
		ico.Write(pngs[i].ICONDIRENTRY())
	}
	// Copy PNG
	for i := 0; i < pngc; i++ {
		ico.Write(pngs[i].Read())
	}
}
