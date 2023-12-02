package main

import (
	"fmt"

	"golang.org/x/example/hello/reverse"
)

func main() {
	v := "Hello, OTUS!"
	fmt.Print(reverse.String(v))
}
