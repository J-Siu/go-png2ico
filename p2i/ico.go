/*
The MIT License

Copyright (c) 2025 John Siu

Permission is hereby granted, free of charge, to any person obtaining a copy of this software and associated documentation files (the "Software"), to deal in the Software without restriction, including without limitation the rights to use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of the Software, and to permit persons to whom the Software is furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
*/

package p2i

import (
	"encoding/binary"
	"encoding/hex"
	"errors"
	"os"

	"github.com/J-Siu/go-helper/v2/basestruct"
	"github.com/J-Siu/go-helper/v2/errs"
	"github.com/J-Siu/go-helper/v2/ezlog"
)

const (
	lenIconDir      uint32 = 6
	lenIconDirEntry uint32 = 16
)

// ICO structure
//
// Must use New() to initialize
type ICO struct {
	basestruct.Base

	File       string `json:"File"`
	fileHandle *os.File
	pngCount   uint16
	pngs       []*PNG
}

func (t *ICO) New(file string) *ICO {
	t.MyType = "ICO"
	t.File = file
	t.Initialized = true
	return t
}

func (t *ICO) PngCount() uint16 {
	return t.pngCount
}

func (t *ICO) AddPng(png *PNG) *ICO {
	prefix := t.MyType + ".AddPng"
	if !t.CheckErrInit(prefix) {
		return t
	}
	ezlog.Trace().N(prefix).N("png").M(png).Out()
	if png.Err == nil {
		if png.IsPNG() {
			t.pngs = append(t.pngs, png)
			t.pngCount++
		} else {
			t.Err = errors.New(png.File + " not PNG")
		}
	}
	errs.Queue(prefix, t.Err)
	return t
}

func (t *ICO) AddPngFile(file string) *ICO {
	prefix := t.MyType + ".AddPngFile"
	if !t.CheckErrInit(prefix) {
		return t
	}
	return t.AddPng(new(PNG).New().Read(file))
}

func (t *ICO) Write() *ICO {
	prefix := t.MyType + ".WriteAll"
	if !t.CheckErrInit(prefix) {
		return t
	}
	// Write ICONDIR
	if t.open().Err == nil {
		t.writeByte(t.iconDir(t.pngCount))
	}
	// Write all ICONDIRENTRY
	for index := range t.pngs {
		if t.Err != nil {
			break
		}
		t.writeByte(t.iconDirEntry(index))
	}
	// Write all PNGs
	for _, png := range t.pngs {
		if t.Err != nil {
			break
		}
		t.writeByte(&png.Buf)
	}
	errs.Queue(prefix, t.Err)

	return t
}

// open ICO file handle
func (t *ICO) open() *ICO {
	prefix := t.MyType + ".open"
	if !t.CheckErrInit(prefix) {
		return t
	}
	ezlog.Debug().N(prefix).M(t.File).Out()
	t.fileHandle, t.Err = os.Create(t.File)
	return t
}

// writeByte ICO
func (t *ICO) writeByte(b *[]byte) *ICO {
	prefix := t.MyType + ".write"
	if !t.CheckErrInit(prefix) {
		return t
	}
	var n int
	n, t.Err = t.fileHandle.Write(*b)
	ezlog.Debug().N(prefix).N("byte").M(n).Out()
	return t
}

// return iconDir byte array
func (t *ICO) iconDir(num uint16) *[]byte {
	prefix := t.MyType + ".icondir"
	/*
		6byte ICONDIR - LittleEndian
		00:   00 00 // 2byte, must be 0
		02:   01 00 // 2byte, 1 for ICO
		04:   xx xx // 2byte, img number
	*/
	b := []byte{0, 0, 1, 0, 0, 0}
	binary.LittleEndian.PutUint16(b[4:6], num)
	ezlog.Debug().N(prefix).M(hex.EncodeToString(b)).Out()
	return &b
}

// return iconDirEntry byte array
func (t *ICO) iconDirEntry(pngIndex int) *[]byte {
	prefix := t.MyType + ".iconDirEntry"
	var (
		b                  []byte = make([]byte, 16)
		existingPngSize    uint32                                        // Sum of all PNGs' size before pngIndex
		lenIconDirEntryAll uint32 = lenIconDirEntry * uint32(t.pngCount) // Always base on final number of PNGs
		offset             uint32
	)
	for index := range pngIndex {
		existingPngSize += t.pngs[index].Size
	}
	offset = lenIconDir + lenIconDirEntryAll + existingPngSize
	/*
		16byte ICONDIRENTRY - LittleEndian
		00:   xx    // 1byte, width
		01:   xx    // 1byte, height
		02:   00    // 1byte, color palette number, 0 for PNG
		03:   00    // 1byte, reserved, always 0
		04:   00 00 // 2byte, color planes, 0 for PNG
		06:   xx xx // 2byte, color depth
		08:   xx xx xx xx // 4byte, image size
		12:   xx xx xx xx // 4byte, image offset
	*/
	png := t.pngs[pngIndex]
	copy(b[0:6], []byte{png.Width, png.Height, 0, 0, 0, 0})
	binary.LittleEndian.PutUint16(b[6:8], png.Depth)
	binary.LittleEndian.PutUint32(b[8:12], png.Size)
	binary.LittleEndian.PutUint32(b[12:16], offset)

	ezlog.Debug().N(prefix).N("byte").M(hex.EncodeToString(b)).N("PNG").M(png.File).Out()

	return &b
}
