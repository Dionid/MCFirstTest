package main

import (
	"io"
	"os"
	"path/filepath"
	"strings"
	"strconv"
	"sort"
)

var (
	NewLineSymbol = []byte("\n")
	TabSymbol = []byte("	")
	HorLineSymbol = []byte("───")
	VertLineSymbol = []byte("│")
	PointVertLineSymbol = []byte("├")
	EndVertLineSymbol = []byte("└")
)

type TreeElSt struct {
	Name string
	IsDir bool
	SizeInB int64
	Inner TreeSt
}

type TreeSt map[string]*TreeElSt

func (this *TreeSt) Fill(rootPath string, printFiles bool) {
	filepath.Walk(rootPath, func(path string, info os.FileInfo, err error) error {
		if path == rootPath || string(path[0]) == "." {
			return nil
		}

		splitedPath := strings.Split(path, string(os.PathSeparator))
		splitedPathL := len(splitedPath)

		if splitedPath[0] == rootPath {
			splitedPath = splitedPath[1:]
			splitedPathL -= 1
		}

		TreeEl := TreeElSt{}
		TreeEl.Name = info.Name()

		if info.IsDir() {
			TreeEl.IsDir = true
			TreeEl.Inner = TreeSt{}
		} else {
			if !printFiles {
				return nil
			}
			TreeEl.IsDir = false
			TreeEl.SizeInB = info.Size()
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

func (this *TreeSt) DisplayEl(out io.Writer, data TreeElSt, end bool, prefix []byte) {
	out.Write(prefix)
	if end {
		out.Write(EndVertLineSymbol)
	} else {
		out.Write(PointVertLineSymbol)
	}
	out.Write(HorLineSymbol)
	out.Write([]byte(data.Name))
	if !data.IsDir {
		out.Write([]byte(" ("))
		if data.SizeInB == 0 {
			out.Write([]byte("empty"))
		}  else {
			// TODO: REWORK
			out.Write([]byte(strconv.Itoa(int(data.SizeInB)) + "b"))
		}
		out.Write([]byte(")"))
	}
	out.Write(NewLineSymbol)
	if data.Inner != nil {
		tL := len(data.Inner)
		cL := 0
		if !end {
			prefix = append(prefix, VertLineSymbol...)
		}
		prefix = append(prefix, TabSymbol...)

		var keys []string
		for k := range data.Inner {
			keys = append(keys, k)
		}
		sort.Strings(keys)

		for _, v := range keys {
			cL++
			this.DisplayEl(out, *data.Inner[v], cL == tL, prefix)
		}
	}
}

func (this *TreeSt) Display(out io.Writer) {

	var keys []string
	for k := range *this {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	tL := len(*this)
	cL := 0
	for _, v := range keys {
		cL++
		this.DisplayEl(out, *(*this)[v], cL == tL, []byte{})
	}
}


func dirTree(out io.Writer, path string, printFiles bool) error {
	Tree := TreeSt{}

	Tree.Fill(path, printFiles)
	Tree.Display(out)

	return nil
}

//func oldDirTree(out io.Writer, path string, printFiles bool) error {
//
//	var (
//		curPrefix []byte
//		curFolder = ""
//	)
//
//	filepath.Walk(path, func(path string, info os.FileInfo, err error) error {
//		if string(path[0]) == "." {
//			return nil
//		}
//
//		out.Write(curPrefix)
//
//		if end {
//			out.Write(EndVertLineSymbol)
//		} else {
//			out.Write(PointVertLineSymbol)
//		}
//		out.Write(HorLineSymbol)
//		out.Write([]byte(info.Name()))
//
//		splitedPath := strings.Split(path, string(os.PathSeparator))
//		splitedPathL := len(splitedPath)
//
//
//		if info.IsDir() {
//
//		} else {
//
//		}
//
//		if curFolder == "" {
//
//		}
//
//
//		return nil
//	})
//
//	return nil
//}

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
