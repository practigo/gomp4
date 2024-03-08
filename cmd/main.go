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

	for i, b := range f.Boxes {
		fmt.Println(i, b.Repr())
	}

	return nil
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("missing argument, provide an mp4 file!")
		return
	}

	if err := run(os.Args[1]); err != nil {
		panic(err)
	}
}