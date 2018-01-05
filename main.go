package main

import (
	"io"
	"os"
	"path/filepath"
	"strings"
	"fmt"
	"encoding/binary"
)

var (
	NewLineSymbol = []byte("\n")
	TabSymbol = []byte("    ")
	HorLineSymbol = []byte("───")
	VertLineSymbol = []byte("│")
	PointVertLineSymbol = []byte("├")
	EndVertLineSymbol = []byte("└")
)

type TreeElSt struct {
	Name string
	SizeInB *int64
	Inner TreeSt
}

type TreeSt map[string]*TreeElSt


func (this *TreeSt) Fill(path string) {
	filepath.Walk(path, func(path string, info os.FileInfo, err error) error {
		if string(path[0]) == "." {
			return nil
		}

		splitedPath := strings.Split(path, "/")
		splitedPathL := len(splitedPath)

		TreeEl := TreeElSt{}
		TreeEl.Name = info.Name()

		if info.IsDir() {
			TreeEl.Inner = TreeSt{}
		} else {
			//b := make([]byte, 8)
			//binary.LittleEndian.PutUint64(b, uint64(info.Size()))
			size := info.Size()
			TreeEl.SizeInB = &size
		}

		var CurTree TreeSt

		if splitedPathL == 1 {
			CurTree = *this
		} else {
			tmpTree := *this
			for i := 0; i < splitedPathL-1; i++ {
				tmpTree = tmpTree[splitedPath[i]].Inner
			}
			CurTree = tmpTree
		}

		CurTree[TreeEl.Name] = &TreeEl

		return nil
	})
}


func (this *TreeSt) DisplayEl(out io.Writer, data TreeElSt, end bool, tabC int, vertLPosArr []int) {
	if end {
		out.Write(EndVertLineSymbol)
	} else {
		out.Write(PointVertLineSymbol)
	}
	out.Write(HorLineSymbol)
	out.Write([]byte(data.Name))
	if data.SizeInB != nil {
		out.Write([]byte(" ("))
		if *data.SizeInB == 0 {
			out.Write([]byte("empty"))
		}  else {
			b := make([]byte, 8)
			binary.LittleEndian.PutUint64(b, uint64(*data.SizeInB))
			fmt.Print(uint64(*data.SizeInB))
			//out.Write(b)
			out.Write([]byte("b"))
		}
		out.Write([]byte(")"))
	}
	out.Write(NewLineSymbol)
	if data.Inner != nil {
		if !end {
			vertLPosArr = append(vertLPosArr, tabC-1)
		}
		tL := len(data.Inner)
		cL := 0
		for _, v := range data.Inner {
			cL++
			for i := 0; i < tabC; i++ {
				for _, vlpV := range vertLPosArr {
					if i == vlpV {
						out.Write(VertLineSymbol)
					}
				}
				out.Write(TabSymbol)
			}
			this.DisplayEl(out, *v, cL == tL, tabC+1, vertLPosArr)
		}
	}
}

func (this *TreeSt) Display(out io.Writer) {
	tL := len(*this)
	cL := 0
	for _, v := range *this {
		cL++
		this.DisplayEl(out, *v, cL == tL, 1, []int{})
	}
}


func dirTree(out io.Writer, path string, printFiles bool) error {

	Tree := TreeSt{}

	Tree.Fill(path)
	Tree.Display(out)

	fmt.Println("Done")

	return nil
}

func oldDirTree(out io.Writer, path string, printFiles bool) error {
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

		splitedPath := strings.Split(path, string(os.PathSeparator))
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
