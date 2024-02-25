package main

import (
	"flag"
	"fmt"
)

func main() {
	dir := flag.String("d", ".", "directory to work on")
	flag.Parse()

	fmt.Println("pineapage")
	fmt.Println(string(*dir))
}
