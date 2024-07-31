package main

import (
	"io"
	"log"
	"os"
)

func main() {
	fileName := "-"
	if len(os.Args) > 1 {
		fileName = os.Args[1]
	}
	if fileName == "-" {
		if err := TransferBufferData(os.Stdin, os.Stdout); err != nil {
			log.Fatal(err)
		}
	} else {
		file, err := os.Open(fileName)
		if err != nil {
			log.Fatal(err)
		}
		if err = TransferBufferData(file, os.Stdout); err != nil {
			log.Fatal(err)
		}
	}
}

func TransferBufferData(in *os.File, out io.Writer) error {
	if _, err := in.WriteTo(out); err != nil {
		return err
	}
	return nil
}
