package main

import (
	"fmt"
	"log"
	"math"
	"net/http"
	"strconv"
	"strings"
)

type Point struct {
	X, Y float64
}

func main() {
	http.HandleFunc("/", handler)
	log.Fatal(http.ListenAndServe("localhost:8000", nil))
}

//generatePoints array
func generatePoints(s string) ([]Point, error) {

	points := []Point{}

	s = strings.Replace(s, "(", "", -1)
	s = strings.Replace(s, ")", "", -1)
	vals := strings.Split(s, ",")
	if len(vals) < 2 {
		return []Point{}, fmt.Errorf("Point [%v] was not well defined", s)
	}

	var x, y float64

	for idx, val := range vals {

		if idx%2 == 0 {
			x, _ = strconv.ParseFloat(val, 64)
		} else {
			y, _ = strconv.ParseFloat(val, 64)
			points = append(points, Point{x, y})
		}
	}
	return points, nil
}

// getArea gets the area inside from a given shape
func getArea(points []Point) float64 {
	// shoelace formula
	area := 0.0
	j := len(points) - 1
	for i := 0; i < len(points); i++ {
		area += (points[j].X + points[i].X) * (points[j].Y - points[i].Y)
		j = i
	}
	return math.Abs(area / 2.0)
}

// getPerimeter gets the perimeter from a given array of connected points
func getPerimeter(points []Point) float64 {
	// calculate segments distance,and adding them to get perimeter
	perimeter := 0.0
	for i := 0; i < len(points)-1; i++ {
		dx := (points[i+1].X - points[i].X)
		dy := (points[i+1].Y - points[i].Y)
		perimeter += math.Sqrt(math.Pow(dx, 2) + math.Pow(dy, 2))
	}
	dx := points[0].X - points[len(points)-1].X
	dy := points[0].Y - points[len(points)-1].Y
	perimeter += math.Sqrt(math.Pow(dx, 2) + math.Pow(dy, 2))
	return perimeter
}

// handler handles the web request and reponds it
func handler(w http.ResponseWriter, r *http.Request) {

	var vertices []Point
	for k, v := range r.URL.Query() {
		if k == "vertices" {
			points, err := generatePoints(v[0])
			if err != nil {
				fmt.Fprintf(w, fmt.Sprintf("error: %v", err))
				return
			}
			vertices = points
			break
		}
	}

	// Results gathering
	area := getArea(vertices)
	perimeter := getPerimeter(vertices)

	// Logging in the server side
	log.Printf("Received vertices array: %v", vertices)

	// Response construction
	response := fmt.Sprintf("Welcome to the Remote Shapes Analyzer\n")
	response += fmt.Sprintf(" - Your figure has : [%v] vertices\n", len(vertices))
	//if there are 2 vertex, we can't calculate an area, so we skip that
	if len(vertices) <= 2 {
		response += fmt.Sprintf("ERROR - Your shape is not compliying with the minimum number of vertices.\n")
	} else {
		response += fmt.Sprintf(" - Vertices        : %v\n", vertices)
		response += fmt.Sprintf(" - Perimeter       : %v\n", perimeter)
		response += fmt.Sprintf(" - Area            : %v\n", area)
	}
	// Send response to client
	fmt.Fprintf(w, response)
}
