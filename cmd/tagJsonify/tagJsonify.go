package main

import (
	"bufio"
	"flag"
	"fmt"
	parser "go-tagparser/pkg/tagparser"
	"io"
	"log"
	"os"
)

var inputPath string
var useStdIn bool

func setupCommandLine() {
	flag.Usage = func() {
		w := flag.CommandLine.Output()
		fmt.Fprintln(w, "\ntagJsonify - Simple Tag Document Jsonification")
		fmt.Fprintln(w, "-------------------------------------")
		fmt.Fprintln(w, "This program will convert your json input into a simple json file.")
		fmt.Fprintln(w, "You must specify either -i to provide an input file, or -stdin to read your input from standard in.")
		fmt.Fprintln(w, "Output will be written to standard out.")
		fmt.Fprintln(w, "Invalid input will result in an error log and program termination.")
		fmt.Fprintln(w, "\nArguments:")
		flag.PrintDefaults()

		fmt.Fprintln(w, "\nExamples:")
		fmt.Fprintln(w, "\ttagJsonify -i ./test.html > test.json")
		fmt.Fprintln(w, "\techo \"<div>Hello, World!</div>\" | tagJsonify -stdin")
	}

	flag.StringVar(&inputPath, "i", "", "Read Tag Document from a file. You must provide the path to the tag document.")
	flag.BoolVar(&useStdIn, "stdin", false, "Read Tag Document from stdin.")
	flag.Parse()
}

func getDocumentBytes() (input []byte) {
	var reader *bufio.Reader
	switch {
	case len(inputPath) > 0:
		file, err := os.Open(inputPath)
		if err != nil {
			log.Panic(err)
		}

		cleanUp := func() {
			err := file.Close()
			if err != nil {
				log.Panic(err)
			}
		}

		defer cleanUp()

		reader = bufio.NewReader(file)
	case useStdIn:
		reader = bufio.NewReader(os.Stdin)
	default:
		{
			flag.Usage()
			os.Exit(0)
			return
		}
	}

	input, err := io.ReadAll(reader)
	if err != nil {
		log.Panic(err)
	}

	return input
}

func main() {
	setupCommandLine()
	input := getDocumentBytes()

	stringInput := string(input)
	result, error := parser.Parse([]rune(stringInput))
	if error != nil {
		log.Fatalf("Error occurred parsing input - %v", error)
	}

	json := result.Root.ToJson()
	os.Stdout.WriteString(json)
}
