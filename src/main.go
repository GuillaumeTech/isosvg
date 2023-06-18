package main

import (
	"bufio"
	"fmt"
	"log"
	"math"
	"os"
	"strconv"
	"strings"

	"gonum.org/v1/gonum/mat"
)

const rotY float64 = 45.0
const rotZ float64 = 45.0

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

	rotMatrix := getRotMatrix(rotY, rotZ)
	for _, line := range lines {
		processed := processLine(line, rotMatrix)
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

func processLine(line string, rotMatrix mat.Dense) string {
	values := strings.Split(line, ",")
	x1, y1, z1, x2, y2, z2 := toFloat(values[0]), toFloat(values[1]), toFloat(values[2]), toFloat(values[3]), toFloat(values[4]), toFloat(values[5])

	// we consider the imagninary plane (our svg) has a base with  y' and z' axes parallel to it. so by doing a rotation it's as if we're projecting on this plane
	// thus we drop X bc it should be on on the normal and we don't represent this axis as it's depth (doesn't matter in isometric 3d)
	_, y1Prime, z1Prime := projectOnPlane(x1, y1, z1, rotMatrix)
	_, y2Prime, z2Prime := projectOnPlane(x2, y2, z2, rotMatrix)

	svgLine := fmt.Sprintf("<line x1=\"%f\" y1=\"%f\" x2=\"%f\" y2=\"%f\" stroke=\"black\" stroke-width=\"5\"/>", y1Prime+100, z1Prime+100, y2Prime+100, z2Prime+100)

	return svgLine
}

func getRotMatrix(rotationY float64, rotationZ float64) mat.Dense {
	var resultZY mat.Dense
	var resultInverse mat.Dense

	rotationYradians := rotationY * (math.Pi / 180)
	rotationZradians := rotationZ * (math.Pi / 180)

	cosY := math.Cos(rotationYradians)
	cosZ := math.Cos(rotationZradians)
	sinY := math.Sin(rotationYradians)
	sinZ := math.Sin(rotationZradians)

	y := mat.NewDense(3, 3, []float64{
		cosY, 0, sinY,
		0, 1, 0,
		-sinY, 0, cosY,
	})
	z := mat.NewDense(3, 3, []float64{
		cosZ, -sinZ, 0,
		sinZ, cosZ, 0,
		0, 0, 1,
	})
	resultZY.Mul(y, z)
	resultInverse.Inverse(&resultZY)
	return resultInverse
}

func projectOnPlane(x float64, y float64, z float64, rotMatrix mat.Dense) (float64, float64, float64) {
	// multipliy by the matrix of the rotation
	var result mat.Dense
	cords := mat.NewDense(1, 3, []float64{
		x, y, z,
	})
	result.Mul(cords, &rotMatrix)

	return result.At(0, 0), result.At(0, 1), result.At(0, 2)

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
