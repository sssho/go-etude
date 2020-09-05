package main

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"
)

type fileTime struct {
	Path    string
	ModTime time.Time
}

func findFile(root string, ext string) []fileTime {
	var files []fileTime

	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return nil
		}
		if info.IsDir() {
			return nil
		}

		if ext != "" {
			if strings.HasSuffix(path, ext) {
				files = append(files, fileTime{path, info.ModTime()})
			}
		} else {
			files = append(files, fileTime{path, info.ModTime()})
		}
		return nil
	})
	if err != nil {
		return nil
	}
	return files
}

func main() {
	files := findFile(".", "")

	// sort by file modification time, decending order
	sort.Slice(files, func(i, j int) bool {
		return files[i].ModTime.After(files[j].ModTime)
	})
	for _, f := range files {
		fmt.Printf("%v\n", f.Path)
	}
	return
}
