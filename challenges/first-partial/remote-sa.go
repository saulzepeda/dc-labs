package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"math"
)

type Point struct {
	X, Y float64
}

func main() {
	http.HandleFunc("/", handler)
	log.Fatal(http.ListenAndServe("localhost:8000", nil))
}

func distance(pointA, pointB Point) float64{
	return math.Sqrt(math.Pow(pointB.X - pointA.X, 2) + math.Pow(pointB.Y - pointA.Y, 2))
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
	sum1 := 0.0
	sum2 := 0.0
	for i := 0; i < len(points)-1; i++ {
		sum1 += points[i].X * points[i+1].Y
		sum2 += points[i].Y * points[i+1].X
	}

	sum1 += points[len(points)-1].X * points[0].Y
	sum2 += points[len(points)-1].Y * points[0].X


	return math.Abs(sum1 - sum2) / 2
}

// getPerimeter gets the perimeter from a given array of connected points
func getPerimeter(points []Point) float64 {
	sum := 0.0
	for i := 0; i < len(points)-1; i++ {
		sum += distance(points[i], points[i+1])
	}
	sum += distance(points[len(points)-1], points[0])
	return sum
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
	
	//Check if there are collisions
	collisions = false
	if len(vertices) > 3{
		for i := 0; i < len(points)-2; i++ {
			p1 = points[i]
			p2 = points[i+1]
			
			q1 = points[i+2]
			q2 = points[i+3]
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

	if len(vertices) < 3{
		response += fmt.Sprintf("ERROR - Your shape is not compliying with the minimum number of vertices.")
	} else {
		response += fmt.Sprintf(" - Vertices        : %v\n", vertices)
		response += fmt.Sprintf(" - Perimeter       : %v\n", perimeter)
		response += fmt.Sprintf(" - Area            : %v\n", area)
	}


	// Send response to client
	fmt.Fprintf(w, response)
}
