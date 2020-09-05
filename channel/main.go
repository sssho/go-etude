package main

import "fmt"

func listWords() chan string {
	result := make(chan string)
	words := []string{"hoge", "fuga", "piyo"}

	go func() {
		for _, w := range words {
			result <- w
		}
		close(result)
	}()

	return result
}

func main() {
	words := listWords()

	for w := range words {
		fmt.Println(w)
	}
	return
}
