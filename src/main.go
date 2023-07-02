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

const rotY float64 = 135
const rotZ float64 = 135
const HEIGHT = 400
const WIDTH = 400


func main() {
	file, err := os.Open("test.3dsvg")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	var lines []string
	

	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	
	svgLines := getSvgData(lines)
	write(svgLines)

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

}


func getSvgData(lines []string) []string {

	var svgLines []string
	svgHead := fmt.Sprintf(`<svg version="1.1" width="%d" height="%d" xmlns="http://www.w3.org/2000/svg">`, WIDTH, HEIGHT)
	
	svgLines = append(svgLines, svgHead)

	rotMatrix := getRotMatrix(rotY, rotZ)
	for _, fileline := range lines {
		values := strings.Split(fileline, ":")
		processed := ""
		if (values[0] == "rect") {
			processed = processRect(values[1], rotMatrix)
		} else{
			processed = processLine(values[1], rotMatrix)
		}


		svgLines = append(svgLines, processed)
	}
	svgLines = append(svgLines, "</svg>")
	return svgLines

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

	svgLine := fmt.Sprintf("<line x1=\"%f\" y1=\"%f\" x2=\"%f\" y2=\"%f\" stroke=\"green\"/>", y1Prime + (WIDTH/2), z1Prime+(HEIGHT/2), y2Prime +(WIDTH/2), z2Prime+(HEIGHT/2))

	return svgLine
}


func processRect(rect string, rotMatrix mat.Dense) string {
	values := strings.Split(rect, ",")
	x1, y1, z1, x2, y2, z2, x3, y3, z3 := toFloat(values[0]), toFloat(values[1]), toFloat(values[2]), toFloat(values[3]), toFloat(values[4]), toFloat(values[5]),  toFloat(values[6]), toFloat(values[7]), toFloat(values[8])

	// we consider the imagninary plane (our svg) has a base with  y' and z' axes parallel to it. so by doing a rotation it's as if we're projecting on this plane
	// thus we drop X bc it should be on on the normal and we don't represent this axis as it's depth (doesn't matter in isometric 3d)
	_, y1Prime, z1Prime := projectOnPlane(x1, y1, z1, rotMatrix)
	_, y2Prime, z2Prime := projectOnPlane(x2, y2, z2, rotMatrix)
	_, y3Prime, z3Prime := projectOnPlane(x3, y3, z3, rotMatrix)
	y4Prime := (y2Prime - y1Prime) + y3Prime
	z4Prime := (z2Prime - z1Prime) + z3Prime


	// we go 1 2 4 3 because we ploygon draw like a pen so we stay on the outside and have a rectangle	
	svgLine := fmt.Sprintf("<polygon points=\"%f,%f %f,%f %f,%f %f,%f\" />", y1Prime + (WIDTH/2), z1Prime+(HEIGHT/2), y2Prime +(WIDTH/2), z2Prime+(HEIGHT/2),  y4Prime +(WIDTH/2), z4Prime+(HEIGHT/2), y3Prime +(WIDTH/2), z3Prime+(HEIGHT/2))

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
