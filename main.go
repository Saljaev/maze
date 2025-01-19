package main

import (
	"bufio"
	"container/heap"
	"errors"
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
	errWriter := bufio.NewWriter(os.Stderr)

	scanner.Scan()
	mazeSize := strings.Split(scanner.Text(), " ")
	if len(mazeSize) != 2 {
		writeError(errWriter, ErrInvalidMazeSize)
		return
	}

	rows := atoi(mazeSize[0])
	columns := atoi(mazeSize[1])

	maze := make([][]int, rows)
	for i := 0; i < rows; i++ {
		scanner.Scan()
		mazeCells := strings.Fields(scanner.Text())
		if len(mazeCells) != columns {
			writeError(errWriter, ErrInvalidCellCount)
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
		writeError(errWriter, ErrInvalidPointInput)
		return
	}

	startPoint := Point{atoi(points[0]), atoi(points[1])}
	endPoint := Point{atoi(points[2]), atoi(points[3])}

	if !isPointValid(startPoint, rows, columns, maze) || !isPointValid(endPoint, rows, columns, maze) {
		writeError(errWriter, ErrUnreachablePoint)
		return
	}

	shortestPath, err := findShortestPath(maze, startPoint, endPoint)
	if err != nil {
		writeError(errWriter, err)
		return
	} else {
		writePath(writer, shortestPath)
	}
}

// atoi переводит string в int без возвращения ошибки.
func atoi(s string) int {
	result, _ := strconv.Atoi(s)
	return result
}

// isPointValid проверяет выходят ли координаты Point за границы лабиринта,
// является ли эта точка доступной (проверка на значение 0).
func isPointValid(p Point, rows, columns int, maze [][]int) bool {
	return p.x >= 0 && p.x < rows && p.y >= 0 && p.y < columns && maze[p.x][p.y] != 0
}

type PriorityQueueItem struct {
	point    Point
	priority int
	index    int
}

type PriorityQueue []*PriorityQueueItem

// Len возвращает длину очереди.
func (pq *PriorityQueue) Len() int {
	return len(*pq)
}

// Less сравнивает приоритетность i и j элемента.
func (pq *PriorityQueue) Less(i, j int) bool {
	return (*pq)[i].priority < (*pq)[j].priority
}

// Swap меняет местами элементы i и j.
func (pq *PriorityQueue) Swap(i, j int) {
	(*pq)[i], (*pq)[j] = (*pq)[j], (*pq)[i]
	(*pq)[i].index = i
	(*pq)[j].index = j
}

// Push добавляет x в очередь.
func (pq *PriorityQueue) Push(x interface{}) {
	item := x.(*PriorityQueueItem)
	item.index = len(*pq)
	*pq = append(*pq, item)
}

// Pop достает верхний элемент с очереди.
func (pq *PriorityQueue) Pop() interface{} {
	old := *pq
	n := len(old)
	item := old[n-1]
	item.index = -1
	*pq = old[0 : n-1]
	return item
}

// Update обновляет приоритет элемента в очереди.
func (pq *PriorityQueue) Update(item *PriorityQueueItem, priority int) {
	item.priority = priority
	heap.Fix(pq, item.index)
}

// findShortestPath ищет путь в лабиринте maze от точки start до точки end,
// использует алгоритм Дейкстры вместе с двоичной кучей.
// При отсутствии пути возвращает ErrNoPath.
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

	priorityQueue := &PriorityQueue{}
	heap.Init(priorityQueue)
	heap.Push(priorityQueue, &PriorityQueueItem{point: start, priority: 0})

	for priorityQueue.Len() > 0 {
		currItem := heap.Pop(priorityQueue).(*PriorityQueueItem)
		curr := currItem.point

		if curr == end {
			break
		}

		for _, d := range directions {
			neighbor := Point{curr.x + d.x, curr.y + d.y}
			if isPointValid(neighbor, rows, columns, maze) {
				newDist := dist[curr.x][curr.y] + maze[neighbor.x][neighbor.y]
				if newDist < dist[neighbor.x][neighbor.y] {
					dist[neighbor.x][neighbor.y] = newDist
					prev[neighbor.x][neighbor.y] = curr
					heap.Push(priorityQueue, &PriorityQueueItem{point: neighbor, priority: newDist})
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

func writeError(writer *bufio.Writer, err error) {
	writer.Write([]byte(err.Error()))
	writer.Flush()
}
