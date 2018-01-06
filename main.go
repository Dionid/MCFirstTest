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

func (this *TreeSt) Fill(rootPath string, printFiles bool) error {

	splitedRootPath := strings.Split(rootPath, string(os.PathSeparator))

	// If user ended rootPath with `rootPath/` than we need to remove last empty element from array
	if splitedRootPath[len(splitedRootPath)-1] == "" {
		splitedRootPath = splitedRootPath[:len(splitedRootPath)-1]
	}

	// This done because `filepath.Walk` works with paths that doesn't have `./` in the beginning
	// Ex: `./testdata` => `testdata`
	if splitedRootPath[0] == "." {
		rootPath = string(rootPath[2:])
		splitedRootPath = splitedRootPath[1:]
	}

	//fmt.Println(rootPath)
	//fmt.Println(splitedRootPath)

	return filepath.Walk(rootPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// We don't process rootPath, because we don't need it by task
		if path == rootPath {
			return nil
		}

		//fmt.Println(rootPath)
		//fmt.Println(path)

		splitedPath := strings.Split(path, string(os.PathSeparator))[len(splitedRootPath):]
		splitedPathL := len(splitedPath)

		//fmt.Println(splitedPath)

		// This done not to show all hidden files
		// TODO: add flag for hidden files
		for _, v := range splitedPath {
			if string(v[0]) == "." {
				return nil
			}
		}

		TreeEl := TreeElSt{}
		TreeEl.Name = info.Name()

		if info.IsDir() {
			TreeEl.IsDir = true
			TreeEl.Inner = TreeSt{}
		} else {
			// If file check for print files
			if !printFiles {
				// If we don't need to print files than continue
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

	if tL == 0 {
		out.Write([]byte("Directory is empty"))
		out.Write(NewLineSymbol)
	}

	cL := 0
	for _, v := range keys {
		cL++
		this.DisplayEl(out, *(*this)[v], cL == tL, []byte{})
	}
}


func dirTree(out io.Writer, path string, printFiles bool) error {
	Tree := TreeSt{}

	if err := Tree.Fill(path, printFiles); err != nil {
		out.Write([]byte(err.Error()))
		out.Write(NewLineSymbol)
		return nil
	}
	Tree.Display(out)

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
