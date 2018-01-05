package main

import (
	"io"
	"os"
	"path/filepath"
	"strings"
)

var (
	NewLineSymbol = []byte("\n")
	TabSymbol = []byte("───")
)

type TreeElSt struct {
	Name string
	SizeInB int
	Inner *TreeSt
}

func (this *TreeElSt) Fill() {

}

type TreeSt []TreeElSt


func (this *TreeSt) FillLvl() {

}

func (this *TreeSt) Display(out io.Writer) {
	//for k, v := range *this {
	//
	//	out.Write(NewLineSymbol)
	//}
}


func dirTree(out io.Writer, path string, printFiles bool) error {
	//info, err := os.Lstat(path)
	//
	//if err != nil {
	//	return nil
	//}
	//
	//out.Write([]byte(info.Name()))
	//out.Write(NewLineSymbol)

	prevSplitedL := 0

	filepath.Walk(path, func(path string, info os.FileInfo, err error) error {
		if path == "." {
			return nil
		}

		out.Write([]byte("│"))

		splitedPath := strings.Split(path, "/")
		splitedPathL := len(splitedPath)

		if splitedPathL > 1 {
			for i := 0; i < splitedPathL-1; i++ {
				out.Write([]byte("    "))
			}
			if splitedPathL < prevSplitedL {
				out.Write([]byte("└"))
			} else {
				out.Write([]byte("│"))
			}
			out.Write([]byte(TabSymbol))
		} else {
			out.Write([]byte(TabSymbol))
		}

		prevSplitedL = splitedPathL

		out.Write([]byte(info.Name() + " "))

		if !info.IsDir() {
			if printFiles {
				//fmt.Print(info.Size())
				//b := make([]byte, 8)
				//binary.LittleEndian.PutUint64(b, uint64(info.Size()))
				//out.Write(b)
				//fmt.Print(uint64(info.Size()))
				//out.Write([]byte("b"))
			} else {
				return nil
			}
		}

		out.Write(NewLineSymbol)

		return nil
	})

	return nil
}


func main() {
	out := os.Stdout
	if !(len(os.Args) == 2 || len(os.Args) == 3) {
		panic("usage go run main.go . [-f]")
	}
	path := os.Args[1]
	printFiles := len(os.Args) == 3 && os.Args[2] == "-f"
	err := dirTree(out, path, printFiles)
	if err != nil {
		panic(err.Error())
	}
}
