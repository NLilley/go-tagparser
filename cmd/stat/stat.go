package main

import (
	"bufio"
	"flag"
	"fmt"
	parser "go-tagparser/pkg"
	"io"
	"log"
	"os"
)

var inputPath string

func setupCommandLine() {
	flag.Usage = func() {
		w := flag.CommandLine.Output()
		fmt.Fprintf(w, "\nstat - Simple Tag Document Statistics\n")
		fmt.Fprintf(w, "-------------------------------------\n")
		fmt.Fprintf(w, "This program will quickly parse a tag document, outputting helpful statics to standard out\n")
		fmt.Fprintf(w, "By default, the document will be read from standard in. If you provide the -i flag however, it will read directly from the file instead\n")
		fmt.Fprintf(w, "If there is an error in your input, this error will be logged instead.\n\n")
		fmt.Fprintf(w, "Arguments:\n")
		flag.PrintDefaults()
	}

	flag.StringVar(&inputPath, "i", "", "Path to input tag document. If not provided, the document will be read from stdin instead")
	flag.Parse()
}

func getDocumentBytes() []byte {
	var reader *bufio.Reader
	if len(inputPath) == 0 {
		reader = bufio.NewReader(os.Stdin)
	} else {
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

	log.Println("Input:")
	print(stringInput)
	print("\n\n")

	stats := parser.CalculateStats(&result.Root)
	renderedStats := stats.Render()

	log.Println(renderedStats)
}
