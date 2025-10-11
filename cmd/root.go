/*
Copyright © 2025 John, Sing Dao, Siu <john.sd.siu@gmail.com>

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.
*/

package cmd

import (
	"errors"
	"os"

	"github.com/J-Siu/go-helper/v2/errs"
	"github.com/J-Siu/go-helper/v2/ezlog"
	"github.com/J-Siu/go-png2ico/v2/global"
	"github.com/J-Siu/go-png2ico/v2/p2i"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:     "go-png2ico <PNG file> <PNG file> ... <ICO file>",
	Version: p2i.Version,
	Short:   "Build ICO file from PNGs",
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		// --- debug setup

		ezlog.SetLogLevel(ezlog.ERR)
		if global.Flag.Debug {
			ezlog.SetLogLevel(ezlog.DEBUG)
		}
		ezlog.Debug().N("Version").Mn(p2i.Version).Nn("Flag").M(&global.Flag).Out()

		// --- check number for filename, minimum 2

		if len(args) < 2 {
			errs.Queue("", errors.New("Input/Output file missing"))
		}

		// --- Pre-checking

		var icoFile string
		if errs.IsEmpty() {
			icoFile = args[len(args)-1]
			// Make sure icoFile is *NOT* PNG
			png := new(p2i.PNG).New().Read(icoFile)
			if png.Err == nil && png.IsPNG() {
				errs.Queue("", errors.New(png.File+": is PNG"))
			} else {
				// Clear the errs queue, as the icoFile may not exist yet and generated a read error
				errs.Clear()
			}
		}

	},
	Run: func(cmd *cobra.Command, args []string) {
		var icoFile string
		if errs.IsEmpty() {
			icoFile = args[len(args)-1]
		}
		ico := new(p2i.ICO).New(icoFile)
		if errs.IsEmpty() {
			// Add PNGs into ico struct
			pngc := len(args) - 1
			for i := range pngc {
				if ico.Err != nil {
					break
				}
				ico.AddPngFile(args[i])
				if global.Flag.Verbose {
					ezlog.Log().N("Add").M(args[i]).Out()
				}
			}
			errs.Queue("", ico.Err)
		}
		if errs.IsEmpty() {
			ico.Write()
			errs.Queue("", ico.Err)
		}
		if errs.IsEmpty() && global.Flag.Verbose {
			ezlog.Log().N("ICO").M(icoFile).Out()
		}
	},
	PersistentPostRun: func(cmd *cobra.Command, args []string) {
		if errs.NotEmpty() {
			ezlog.Err().L().M(errs.Errs).Out()
			cmd.Usage()
		}
	},
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().BoolVarP(&global.Flag.Debug, "debug", "d", false, "Enable debug")
	rootCmd.PersistentFlags().BoolVarP(&global.Flag.Verbose, "verbose", "v", false, "Verbose")
}
