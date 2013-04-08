package main

import(
	"fmt"
	"os"
	"io/ioutil"
)

type FileArr [][]uint8

func byteDiffers(offset int, files FileArr) bool {
	char := files[0][offset]
	for f := range files {
		if len(files[f]) <= offset || files[f][offset] != char {
			return true
		}
	}

	return false
}

func lineDiffers(offset int, fileData FileArr) (bool, [16]bool) {
	var differs bool
	var diffMask [16]bool
	var target uint8

	for i := 0; i < 16; i++ {
		if offset + i >= len(fileData[0]) {
			return differs, diffMask
		}

		for f := range fileData {
			if f == 0 {
				target = fileData[f][offset + i]
			} else {
				if fileData[f][offset + i] != target || offset + i > len(fileData[f]) {
					diffMask[i] = true
					differs = true
					break
				}
			}
		}
	}

	return differs, diffMask
}


func outputHex(offset int, file []uint8) string {
	if offset >= len(file) {
		return "  "
	}

	return fmt.Sprintf("%02x", file[offset])
}


func compareLine(offset int, fileData FileArr) {
	differs, diffMask := lineDiffers(offset, fileData)
	_,_ = differs, diffMask
	fmt.Printf("%08x  ", offset)
	for f := range fileData {
		for i := 0; i < 16; i++ {
			tmp := outputHex(offset + i, fileData[f])
			if (diffMask[i]) {
				fmt.Printf("\x1b[31;1m%s\x1b[0m ", tmp)
			} else {
				fmt.Printf("%s ", tmp)
			}

			if i == 7 {
				fmt.Print(" ")
			}
		}
		fmt.Print("\t")
	}
	fmt.Print("\n")
}

func maxSize(f FileArr) int {
	var max int

	for i := range f {
		if len(f[i]) > max {
			max = len(f[i])
		}
	}

	return max
}

func isGraph(char uint8) bool {
	return (char >= 0x20 && char <= 0x7E)
}

func main() {
	var err error
	var paths = os.Args[1:]

	if len(paths) < 2 {
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

	
	maxfile := maxSize(fileData)
	for i := 0; i < maxfile; i += 16 {
		compareLine(i, fileData)
	}

}

/* vim: set noexpandtab:ts=4:sw=4:sts:4 */
