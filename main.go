package main

import (
	"bufio"
	"fmt"
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
	go newSigHandler()
	iterateAndOutput(fileNames, mode)
}

func newSigHandler() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGPIPE)
	_ = <-c
	os.Exit(0)
}

func iterateAndOutput(fileNames []string, mode int) {
	lineCount := 1
	for _, fileName := range fileNames {
		var s *bufio.Scanner
		if fileName == "-" {
			s = bufio.NewScanner(os.Stdin)
		} else {
			file, openErr := os.Open(fileName)
			if openErr != nil {
				log.Fatal(openErr)
			}
			s = bufio.NewScanner(file)
		}
		for s.Scan() {
			text := s.Text()
			if mode&COUNT_LINES != 0 && (mode&EXCLUDE_BLANKS == 0 || len(text) > 0) {
				text = fmt.Sprintf("%6d  %s", lineCount, text)
				lineCount++
			}
			text += "\n"
			os.Stdout.Write([]byte(text))
			//fmt.Print(text, "\n")
		}
	}
}
