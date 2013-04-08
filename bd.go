/*
 * Copyright (c) 2013, Chris Anderson
 * All rights reserved.
 * Redistribution and use in source and binary forms, with or without
 * modification, are permitted provided that the following conditions are met:
 *
 * 1. Redistributions of source code must retain the above copyright notice, this
 * list of conditions and the following disclaimer.
 *
 * THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS "AS IS" AND
 * ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE IMPLIED
 * WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE ARE
 * DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT OWNER OR CONTRIBUTORS BE LIABLE FOR
 * ANY DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES
 * (INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES;
 * LOSS OF USE, DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND
 * ON ANY THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT
 * (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE OF THIS
 * SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.
 */
package main

import (
	"fmt"
	"io/ioutil"
	"os"
)

type FileArr [][]uint8
var maxFile int
var maxIndex int

func isGraph(char uint8) bool {
	return (char > 0x20 && char <= 0x7E)
}

/* Mirrors output of the hd invocation of the linux hexdump */
func hexdump(f []uint8) {
	for offset := 0; offset < len(f); offset += 16 {
		fmt.Printf("%08x  ", offset)
		for i := 0; i < 16; i++ {
			fmt.Printf("%s ", outputHex(offset + i, f))

			if i == 7 { fmt.Print(" ") }
		}
		fmt.Print(" |")
		for  i := 0; i < 16; i++ {
			if offset + i < len(f) {
				tmp := f[offset+i]
				if isGraph(tmp) {
					fmt.Printf("%c", tmp)
				} else {
					fmt.Print(".")
				}
			} else {
				fmt.Print(".")
			}
		}
		fmt.Print("|\n")
	}
}

/* Determines whether if a 16 byte row has a mismatch in any supplied file */
func lineDiffers(offset int, base []uint8, remote []uint8) [16]bool {
	var diffMask [16]bool
	var target uint8

	for i := 0; i < 16; i++ {
		// No reason to check beyond the largest of the files
		if offset+i >= maxFile {
			return diffMask
		}

		// If the lengths don't match then it's an obvious diff, otherwise check target
		if offset+i < len(base) {
			target = base[offset+i]
			if offset+i >= len(remote) {
				diffMask[i] = true
			} else if remote[offset+i] != target {
				diffMask[i] = true
			}
		} else {
			diffMask[i] = true
		}
	}

	return diffMask
}

func outputHex(offset int, file []uint8) string {
	if offset >= len(file) {
		return "  "
	}

	return fmt.Sprintf("%02x", file[offset])
}

func printHexLine(offset int, fileData FileArr) {
	for f := range fileData {
		diffMask := lineDiffers(offset, fileData[0], fileData[f])
		fmt.Printf("%08x  ", offset)
		for i := 0; i < 16; i++ {
			tmp := outputHex(offset+i, fileData[f])
			if diffMask[i] {
				fmt.Printf("\x1b[31;1m%s\x1b[0m ", tmp)
			} else {
				fmt.Printf("%s ", tmp)
			}

			if i == 7 {
				fmt.Print(" ")
			}
		}
		fmt.Print("  ")
	}
	fmt.Print("\n")
}

func getMaxes(f FileArr) {
	for i := range f {
		if len(f[i]) > maxFile {
			maxFile = len(f[i])
			maxIndex = i
		}
	}
}

func main() {
	var err error
	var paths = os.Args[1:]

	if len(paths) < 1 {
		println("usage:", os.Args[0], "<files>")
		os.Exit(1)
	}

	fileData := make(FileArr, len(paths))
	for i := 0; i < len(paths); i++ {
		fileData[i], err = ioutil.ReadFile(paths[i])
		if err != nil {
			println("couldn't open", paths[i], ":", err.Error())
			os.Exit(1)
		}
	}

	/* If only one file is provided then act like hd */
	if len(paths) == 1 {
		hexdump(fileData[0])
	} else {
		getMaxes(fileData)
		for offset := 0; offset < maxFile; offset += 16 {
			printHexLine(offset, fileData)
		}
	}
}

/* vim: set noexpandtab:ts=4:sw=4:sts:4 */
