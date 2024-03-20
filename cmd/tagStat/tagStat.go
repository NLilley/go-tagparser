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
		fmt.Fprintln(w, "\ntagStat - Simple Tag Statistics")
		fmt.Fprintln(w, "-------------------------------------")
		fmt.Fprintln(w, "This program will quickly parse a tag document, outputting helpful statistics to standard out")
		fmt.Fprintln(w, "You must specify either -i to provide an input file, or -stdin to read your input from standard in.")
		fmt.Fprintln(w, "Invalid input will result in an error log and program termination.")
		fmt.Fprintln(w, "\nArguments:")
		flag.PrintDefaults()

		fmt.Println("\nExamples:")
		fmt.Println("\ttagStat -i ./test.html")
		fmt.Println("\techo \"<div>Hello, World!</div>\" | tagStat -stdin")
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

	fmt.Println("Input:")
	fmt.Println(stringInput)
	fmt.Print("\n\n")

	stats := parser.CalculateStats(&result.Root)
	renderedStats := stats.Render()

	log.Println(renderedStats)
}
