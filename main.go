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
	TabSymbol = []byte("     ")
	HorLineSymbol = []byte("────")
	VertLineSymbol = []byte("│")
)

type TreeElSt struct {
	Name string
	SizeInB []byte
	Inner TreeSt
}

func (this *TreeElSt) Fill() {

}

type TreeSt map[string]*TreeElSt


func (this *TreeSt) FillLvl() {

}

func (this *TreeSt) Display(out io.Writer, what TreeSt, tabC int) {
	//for _, v := range what {
	//	for i := 0; i < tabC; i++ {
	//		out.Write(TabSymbol)
	//	}
	//	if v.Inner != nil {
	//		fmt.Print(v.Name)
	//		fmt.Print(" ")
	//		fmt.Println(v.Inner)
	//		this.Display(out, v.Inner, tabC + 1)
	//	} else {
	//		fmt.Println(v.Name)
	//	}
	//}

	for _, v := range what {
		out.Write(VertLineSymbol)
		for i := 0; i < tabC; i++ {
			out.Write(TabSymbol)
		}
		out.Write(VertLineSymbol)
		out.Write(HorLineSymbol)
		if v.Inner != nil {
			out.Write([]byte(v.Name))
			out.Write(NewLineSymbol)
			this.Display(out, v.Inner, tabC+1)
		} else {
			out.Write([]byte(v.Name + " "))
			//out.Write(v.SizeInB)
			out.Write(NewLineSymbol)
		}
	}
}


func dirTree(out io.Writer, path string, printFiles bool) error {

	Tree := TreeSt{}

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
			b := make([]byte, 8)
			binary.LittleEndian.PutUint64(b, uint64(info.Size()))
			TreeEl.SizeInB = b
		}

		var CurTree TreeSt

		if splitedPathL == 1 {
			CurTree = Tree
		} else {
			tmpTree := Tree
			for i := 0; i < splitedPathL-1; i++ {
				tmpTree = tmpTree[splitedPath[i]].Inner
			}
			CurTree = tmpTree
		}

		CurTree[TreeEl.Name] = &TreeEl

		return nil
	})

	Tree.Display(out, Tree, 0)

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
