package main

import (
	"fmt"
	"os"
	"path/filepath"
)

func listDirAndFiles() chan string {
	result := make(chan string)

	go func() {
		err := filepath.Walk(".", func(path string, info os.FileInfo, err error) error {
			if err != nil {
				result <- fmt.Sprintf("prevent panic: %q: %v", path, err)
				return nil
			}
			if info.IsDir() {
				result <- fmt.Sprintf("dir: %q", path)
			}
			result <- fmt.Sprintf("visited file or dir: %q", path)
			return nil
		})

		if err != nil {
			close(result)
		}
		close(result)
	}()

	return result
}

func main() {
	list := listDirAndFiles()

	for l := range list {
		fmt.Println(l)
	}
	return
}
