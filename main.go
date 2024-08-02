package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"strings"
)

type MyWriter struct {
	dest          io.Writer
	lineNum       int
	excludeBlanks bool
	prevByte      byte
}

func (out *MyWriter) Write(bytes []byte) (int, error) {
	for i, b := range bytes {
		var p []byte
		if out.prefixRequired(b) {
			p = out.prefixGen()
			out.lineNum++
		}
		p = append(p, b)
		_, err := out.dest.Write(p)
		if err != nil {
			return i, err
		}
		out.prevByte = b
	}
	return len(bytes), nil
}

func (out *MyWriter) prefixGen() []byte {
	return []byte(fmt.Sprintf("%6v  ", out.lineNum))
}

func (out *MyWriter) prefixRequired(r byte) bool {
	return out.prevByte == '\n' && (!out.excludeBlanks || r != '\n')
}

func NewLineNumberer(dest io.Writer, excludeBlanks bool) io.Writer {
	return &MyWriter{dest, 1, excludeBlanks, '\n'}
}

func main() {
	args := os.Args[1:]
	flags := []string{}
	fileNames := []string{}
	var out io.Writer = os.Stdout
	for _, arg := range args {
		if strings.HasPrefix(arg, "-") {
			flags = append(flags, arg)
		} else {
			fileNames = append(fileNames, arg)
		}
	}
	for _, flag := range flags {
		switch flag {
		case "-":
			fileNames = append(fileNames, flag)
		case "-n":
			out = NewLineNumberer(os.Stdout, false)
		case "-b":
			out = NewLineNumberer(os.Stdout, true)
		default:
			log.Fatalf("unknown flag %s", flag)
		}
	}
	if len(fileNames) == 0 {
		fileNames = []string{"-"}
	}
	err := IterateAndOutput(fileNames, out)
	if err != nil {
		log.Fatal(err)
	}
}

func IterateAndOutput(fileNames []string, out io.Writer) error {
	for _, fileName := range fileNames {
		switch fileName {
		case "-":
			_, err := os.Stdin.WriteTo(out)
			if err != nil {
				return err
			}
		default:
			file, openErr := os.Open(fileName)
			if openErr != nil {
				return openErr
			}
			defer file.Close()

			_, err := file.WriteTo(out)
			if err != nil {
				return err
			}
		}
	}
	return nil
}
