package main

import (
	"fmt"
	"log"
	"os"
)

func main() {
	fileName := os.Args[1]
	bytes, err := os.ReadFile(fileName)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%s", bytes)
}
