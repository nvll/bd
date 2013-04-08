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

func compareLine(offset int, fileData FileArr) {
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
