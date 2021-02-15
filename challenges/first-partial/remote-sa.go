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

func getOrientation(p, q, r Point) int {
	val := (float64(q.Y - p.Y) * (r.X - q.X)) - (float64(q.X - p.X) * (r.Y - q.Y)) 
    if val > 0 {
		// Clockwise orientation 
        return 1
	} else if val < 0 {
		// Counterclockwise orientation 
        return 2
	} else {
		// Colinear orientation 
        return 0
	}
}

func onSegment(p, q, r Point) bool {
	if ( (q.X <= math.Max(p.X, r.X)) && (q.X >= math.Min(p.X, r.X)) && 
	(q.Y <= math.Max(p.Y, r.Y)) && (q.Y >= math.Min(p.Y, r.Y))){
		return true
	}
	return false
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
	collisions := false
	if len(vertices) > 3{
		for i := 0; i < len(vertices)-3; i++ {
			p1 := vertices[i]
			p2 := vertices[i+1]
			q1 := vertices[i+2]
			q2 := vertices[i+3]

			o1 := getOrientation(p1, q1, p2) 
			o2 := getOrientation(p1, q1, q2) 
			o3 := getOrientation(p2, q2, p1) 
			o4 := getOrientation(p2, q2, q1) 

			// General case 
			if ((o1 != o2) && (o3 != o4)) {
				collisions = true
			} 
				
			// Special Cases 
			// p1 , q1 and p2 are colinear and p2 lies on segment p1q1 
			if ((o1 == 0) && onSegment(p1, p2, q1)){
				collisions = true
				break
			} 
		
			// p1 , q1 and q2 are colinear and q2 lies on segment p1q1 
			if ((o2 == 0) && onSegment(p1, q2, q1)) {
				collisions = true
				break
			} 
		
			// p2 , q2 and p1 are colinear and p1 lies on segment p2q2 
			if ((o3 == 0) && onSegment(p2, p1, q2)) {
				collisions = true
				break
			}
		
			// p2 , q2 and q1 are colinear and q1 lies on segment p2q2 
			if ((o4 == 0) && onSegment(p2, q1, q2)) {
				collisions = true
				break
			}

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
	} else if collisions{
		response += fmt.Sprintf("ERROR - Your shape has collisions.")
	} else {
		response += fmt.Sprintf(" - Vertices        : %v\n", vertices)
		response += fmt.Sprintf(" - Perimeter       : %v\n", perimeter)
		response += fmt.Sprintf(" - Area            : %v\n", area)
	}


	// Send response to client
	fmt.Fprintf(w, response)
}
