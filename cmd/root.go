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
	"fmt"
	"os"

	"github.com/J-Siu/go-png2ico/v2/global"
	"github.com/J-Siu/go-png2ico/v2/p2i"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:     "go-png2ico <PNG file> <PNG file> ... <ICO file>",
	Version: p2i.Version,
	Short:   "Build ICO file from PNGs",
	PersistentPreRunE: func(cmd *cobra.Command, args []string) (err error) {
		// --- debug setup

		if global.Flag.Debug {
			fmt.Println("Version:", p2i.Version)
			p2i.PrintStruct("Flag", global.Flag)
		}
		// --- check number for filename, minimum 2

		if len(args) < 2 {
			err = errors.New("Input/Output file missing")
		}
		return
	},
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		var (
			ico     *p2i.ICO
			icoFile string
			png     *p2i.PNG
			pngNum  = len(args) - 1
		)

		icoFile = args[len(args)-1]
		// check output is NOT PNG
		png = new(p2i.PNG).New(global.Flag.Debug).Read(icoFile)
		if png.Err == nil && png.IsPNG() {
			err = errors.New(png.File + ": is PNG")
		} else {
			ico = new(p2i.ICO).New(icoFile, global.Flag.Debug)
			for i := range pngNum {
				if ico.Err != nil {
					break
				}
				ico.PngAddFile(args[i])
				if global.Flag.Verbose {
					fmt.Println("Add:", args[i])
				}
			}
			if ico.Err == nil {
				ico.Write()

				fmt.Println("ICO:", icoFile)
			}
			err = ico.Err
		}
		return err
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
