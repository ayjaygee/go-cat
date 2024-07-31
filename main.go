package main

import (
	"io"
	"log"
	"os"
)

func main() {
	err := IterateAndOutput(os.Args[1:], os.Stdout)
	if err != nil {
		log.Fatal(err)
	}
}

func IterateAndOutput(fileNames []string, out io.Writer) error {
	if len(fileNames) == 0 {
		fileNames = []string{"-"}
	}
	for _, fileName := range fileNames {
		switch fileName {
		case "-":
			err := OutputBufferData(os.Stdin, out)
			if err != nil {
				return err
			}
		default:
			file, openErr := os.Open(fileName)
			if openErr != nil {
				return openErr
			}
			defer file.Close()

			err := OutputBufferData(file, out)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func OutputBufferData(in *os.File, out io.Writer) error {
	if _, err := in.WriteTo(out); err != nil {
		return err
	}
	return nil
}
