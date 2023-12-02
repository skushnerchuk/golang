package main

import (
	"fmt"

	"golang.org/x/example/hello/reverse"
)

func main() {
	reversed := reverse.String("Hello, OTUS!")
	fmt.Print(reversed)
}
