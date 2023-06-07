package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
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

func toFloat(number string) float64 {
	floated, _ := strconv.ParseFloat(number, 64)
	return floated
}

func processLine(line string) string {
	values := strings.Split(line, ",")
	x1, y1, z1, x2, y2, z2 := toFloat(values[0]), toFloat(values[1]), toFloat(values[2]), toFloat(values[3]), toFloat(values[4]), toFloat(values[5])

	// drop X bc it should be on parallel to the normal and we don't represent this axis as it's depth (doesn't matter in isometric 3d)
	_, y1Prime, z1Prime := projectOnPlane(x1, y1, z1)
	_, y2Prime, z2Prime := projectOnPlane(x2, y2, z2)

	svgLine := fmt.Sprintf("<line x1=\"%f\" y1=\"%f\" x2=\"%f\" y2=\"%f\" stroke=\"black\" stroke-width=\"5\"/>", y1Prime+100, z1Prime+100, y2Prime+100, z2Prime+100)

	return svgLine
}

func projectOnPlane(x float64, y float64, z float64) (float64, float64, float64) {
	// multipliy by the matrix of the rotation on z and y
	//			|       1        1  -sqrt(2) |
	//     0.5	| -sqrt(2)  sqrt(2)       0  |
	//			|       1        1   sqrt(2) |

	projectedX := 0.5 * (float64(x) - 1.414*float64(y) + float64(z))
	projectedY := 0.5 * (float64(x) + 1.414*float64(y) + float64(z))
	projectedZ := 0.5 * (-1.414*float64(x) + 1.414*float64(z))
	return projectedX, projectedY, projectedZ

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
