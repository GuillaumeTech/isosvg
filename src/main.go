package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
	"unicode"
)

func main() {
	file, err := os.Open("test.3dsvg")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	var lines []string
	var svgLines []string

	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	svgLines = append(svgLines, `<svg version="1.1" width="300" height="200" xmlns="http://www.w3.org/2000/svg">`)
	for _, line := range lines {
		processed := processLine(line)
		svgLines = append(svgLines, processed)
	}
	svgLines = append(svgLines, "</svg>")

	write(svgLines)

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

}

func processLine(line string) string {

	const svgLine = `<line x1="50" x2="50" y1="110" y2="150" stroke="black" stroke-width="5"/>`

	return svgLine
}

func write(lines []string) {

	file, err := os.Create("result.svg")

	if err != nil {
		log.Fatalf("failed creating file: %s", err)
	}

	datawriter := bufio.NewWriter(file)

	for _, line := range lines {
		_, _ = datawriter.WriteString(line + "\n")
	}

	datawriter.Flush()
	file.Close()
}
