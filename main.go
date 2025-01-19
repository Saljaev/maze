package main

import (
	"bufio"
	"errors"
	"fmt"
	"math"
	"os"
	"strconv"
	"strings"
)

var (
	ErrInvalidMazeSize   = errors.New("invalid size of maze")
	ErrInvalidCellCount  = errors.New("invalid cell count")
	ErrInvalidPointInput = errors.New("invalid start/end point input")
	ErrUnreachablePoint  = errors.New("start/end is unreachable point")
	ErrNoPath            = errors.New("path no found")
)

type Point struct {
	x, y int
}

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)

	scanner.Scan()
	mazeSize := strings.Split(scanner.Text(), " ")
	if len(mazeSize) != 2 {
		fmt.Fprintln(os.Stderr, ErrInvalidMazeSize)
		return
	}

	rows := atoi(mazeSize[0])
	columns := atoi(mazeSize[1])

	maze := make([][]int, rows)
	for i := 0; i < rows; i++ {
		scanner.Scan()
		mazeCells := strings.Fields(scanner.Text())
		if len(mazeCells) != columns {
			fmt.Fprintln(os.Stderr, ErrInvalidCellCount)
			return
		}
		maze[i] = make([]int, columns)
		for j := 0; j < columns; j++ {
			maze[i][j] = atoi(mazeCells[j])
		}
	}

	scanner.Scan()
	points := strings.Fields(scanner.Text())
	if len(points) != 4 {
		fmt.Fprintln(os.Stderr, ErrInvalidPointInput)
		return
	}

	startPoint := Point{atoi(points[0]), atoi(points[1])}
	endPoint := Point{atoi(points[2]), atoi(points[3])}

	if !isPointValid(startPoint, rows, columns, maze) || !isPointValid(endPoint, rows, columns, maze) {
		fmt.Fprintln(os.Stderr, ErrUnreachablePoint)
		return
	}

	shortestPath, err := findShortestPath(maze, startPoint, endPoint)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return
	} else {
		writePath(writer, shortestPath)
	}
}

func atoi(s string) int {
	result, _ := strconv.Atoi(s)
	return result
}

func isPointValid(p Point, rows, columns int, maze [][]int) bool {
	return p.x >= 0 && p.x < rows && p.y >= 0 && p.y < columns && maze[p.x][p.y] != 0
}

func findShortestPath(maze [][]int, start, end Point) ([]Point, error) {
	rows := len(maze)
	columns := len(maze[0])

	directions := []Point{
		{0, 1},
		{0, -1},
		{1, 0},
		{-1, 0},
	}

	dist := make([][]int, rows)
	prev := make([][]Point, rows)
	for i := 0; i < rows; i++ {
		dist[i] = make([]int, columns)
		prev[i] = make([]Point, columns)
		for j := 0; j < columns; j++ {
			dist[i][j] = int(math.Inf(1))
			prev[i][j] = Point{-1, -1}
		}
	}

	dist[start.x][start.y] = 0

	type Queue struct {
		point Point
		cost  int
	}

	queue := []Queue{{start, 0}}

	for len(queue) > 0 {
		curr := queue[0]
		queue = queue[1:]

		if curr.point == end {
			break
		}

		for _, d := range directions {
			neighbor := Point{curr.point.x + d.x, curr.point.y + d.y}
			if isPointValid(neighbor, rows, columns, maze) {
				newCost := curr.cost + maze[neighbor.x][neighbor.y]
				if newCost < dist[neighbor.x][neighbor.y] {
					dist[neighbor.x][neighbor.y] = newCost
					prev[neighbor.x][neighbor.y] = curr.point
					queue = append(queue, Queue{neighbor, newCost})
				}
			}
		}
	}

	var path []Point
	for at := end; at != (Point{-1, -1}); at = prev[at.x][at.y] {
		path = append([]Point{at}, path...)
	}

	if len(path) > 0 && path[0] == start {
		return path, nil
	}

	return nil, ErrNoPath
}

func writePath(writer *bufio.Writer, points []Point) {
	for i := range len(points) {
		pointCoordinate := strconv.Itoa(points[i].x) + " " + strconv.Itoa(points[i].y)
		writer.WriteString(pointCoordinate)
		writer.WriteRune('\n')
	}
	writer.WriteRune('.')
	writer.Flush()
}
