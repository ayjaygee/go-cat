package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"
)

const (
	COUNT_LINES = 1 << iota
	EXCLUDE_BLANKS
	DEFAULT_MODE = 0
)

type myWriter struct {
	dest     io.Writer
	lineNum  int
	mode     int
	prevByte byte
}

func (out *myWriter) Write(bytes []byte) (int, error) {
	if out.mode&COUNT_LINES == 0 {
		return out.dest.Write(bytes)
	}
	for i, b := range bytes {
		var p []byte
		if out.prefixRequired(b) {
			p = []byte(fmt.Sprintf("%6v  ", out.lineNum))
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

func (out *myWriter) prefixRequired(b byte) bool {
	return out.mode&COUNT_LINES != 0 &&
		out.prevByte == '\n' &&
		(out.mode&EXCLUDE_BLANKS == 0 || b != '\n')
}

func main() {
	args := os.Args[1:]
	fileNames := []string{}
	mode := DEFAULT_MODE
	for _, arg := range args {
		if strings.HasPrefix(arg, "-") {
			switch arg {
			case "-":
				fileNames = append(fileNames, arg)
			case "-n":
				mode = mode | COUNT_LINES
			case "-b":
				mode = mode | COUNT_LINES | EXCLUDE_BLANKS
			default:
				log.Fatalf("unknown flag %s", arg)
			}
		} else {
			fileNames = append(fileNames, arg)
		}
	}
	if len(fileNames) == 0 {
		fileNames = []string{"-"}
	}
	out := newWriter(os.Stdout, mode)
	newSigHandler()
	err := iterateAndOutput(fileNames, out)
	if err != nil {
		log.Fatal(err)
	}
}

func newSigHandler() {
	sigs := make(chan os.Signal)
	signal.Notify(sigs, syscall.SIGPIPE)
	go func() {
		sig := <-sigs
		os.Exit(0)
		fmt.Println(sig)
	}()
}

func newWriter(dest io.Writer, mode int) *myWriter {
	return &myWriter{
		dest:     dest,
		lineNum:  1,
		mode:     mode,
		prevByte: '\n',
	}
}

func iterateAndOutput(fileNames []string, out io.Writer) error {
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
