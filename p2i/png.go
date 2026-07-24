/*
The MIT License

Copyright (c) 2025 John Siu

Permission is hereby granted, free of charge, to any person obtaining a copy of this software and associated documentation files (the "Software"), to deal in the Software without restriction, including without limitation the rights to use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of the Software, and to permit persons to whom the Software is furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
*/

package p2i

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"os"
)

// PNG structure
//
// Must use New() to initialize
type PNG struct {
	Debug         bool
	Err           error
	Initialized   bool
	MyType        string
	OnErrContinue bool

	Buf    []byte `json:"Buf"`
	Depth  uint8  `json:"Depth"` // bit/pixel
	File   string `json:"File"`  // filename
	Width  uint32 `json:"Width"`
	Height uint32 `json:"Height"`
	Size   uint32 `json:"Size"`
	isPNG  bool
}

func (t *PNG) New(debug bool) *PNG {
	t.MyType = "PNG"
	t.Initialized = true
	t.Debug = debug
	return t
}

// Valid after `Check`. Else will be `false`
func (t *PNG) IsPNG() bool {
	return t.isPNG
}

// read PNG file and populate other fields
//
// use IsPNG() after Read()
func (t *PNG) Read(file string) *PNG {
	prefix := t.MyType + ".Read"
	t.File = file
	t.Buf, t.Err = os.ReadFile(t.File)
	if t.Err == nil {
		if t.Debug {
			fmt.Println(prefix+":byte", len(t.Buf))
		}
		if t.checkMagic() {
			t.info()
		}
	} else {
		t.Err = errors.New(prefix + ":" + file + ":" + t.Err.Error())
	}
	return t
}

// Verify if Buf is PNG, update and return isPNG
func (t *PNG) checkMagic() bool {
	prefix := t.MyType + ".checkMagic"

	t.isPNG = false

	/*
		25byte PNG header - BigEndian
		00:	89 50 4e 47 0d 0a 1a 0a // 8byte - magic number <- CHECK 1
		IHDR chunk
		08:	xx xx xx xx // 4byte - chunk length
		12:	49 48 44 52 // 4byte - chunk type(IHDR) <- CHECK 2
		16:	xx xx xx xx // 4byte - width
		20:	xx xx xx xx // 4byte - height
		24:	xx          // 1byte - bit depth (bit/pixel)
	*/

	// CHECK 1: 8byte header[0:8] - magic number
	magic := []byte{0x89, 0x50, 0x4e, 0x47, 0x0d, 0x0a, 0x1a, 0x0a}
	t.isPNG = bytes.Equal(magic[:], t.Buf[:8])
	if t.isPNG {
		if t.Debug {
			fmt.Println(prefix + ":Found PNG magic")
		}
		// CHECK 2: 4byte header[12:16] - chunk type IHDR
		t.isPNG = t.isPNG && bytes.Equal([]byte("IHDR"), t.Buf[12:16])
	}
	return t.isPNG
}

// Only update info fields if `isPNG` == true
func (t *PNG) info() *PNG {
	prefix := t.MyType + ".getInfo"
	if t.isPNG {
		// 4byte header[16:20] - width
		t.Width = binary.BigEndian.Uint32(t.Buf[16:20])
		// 4byte header[20:24] - height
		t.Height = binary.BigEndian.Uint32(t.Buf[20:24])
		// 1byte header[25] - color depth
		t.Depth = uint8(t.Buf[24])
		// file size
		stat, _ := os.Stat(t.File)
		t.Size = uint32(stat.Size())
		if t.Debug {
			PrintStruct(prefix+":png", t)
		}
	}
	return t
}

func (t *PNG) CheckErrInit(prefix string) (pass bool) {
	pass = true
	if !t.OnErrContinue {
		// check error first
		if t.Err != nil {
			pass = false
		} else if !t.Initialized {
			errMsg := "not initialized"
			if prefix != "" {
				errMsg = prefix + ": " + errMsg
			} else if t.MyType != "" {
				errMsg = t.MyType + ": " + errMsg
			}
			t.Err = errors.New(errMsg)
			pass = false
		}
	}
	return pass
}
