package main

import (
	"fmt"
	"os"

	"github.com/practigo/gomp4"
)

func run(filename string) error {
	f, err := gomp4.Open(filename)
	if err != nil {
		return err
	}

	return gomp4.View(f)
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("missing argument, provide an mp4 file!")
		os.Exit(1)
	}

	if err := run(os.Args[1]); err != nil {
		panic(err)
	}
}
