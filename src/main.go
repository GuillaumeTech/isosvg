package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
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
	values := strings.Split(line, ",")
	x1, y1, x2, y2 := values[0], values[1], values[2], values[3]
	svgLine := fmt.Sprintf("<line x1=\"%s\" y1=\"%s\" x2=\"%s\" y2=\"%s\" stroke=\"black\" stroke-width=\"5\"/>", x1, y1, x2, y2)

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
