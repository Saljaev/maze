package main

import (
	"reflect"
	"testing"
)

func Test_findShortestPath(t *testing.T) {
	type args struct {
		maze  [][]int
		start Point
		end   Point
	}
	type want struct {
		path []Point
		err  error
	}

	tests := []struct {
		name string
		args args
		want want
	}{
		{
			"Valid maze and valid path",
			args{
				maze:  [][]int{{1, 2, 0}, {2, 0, 1}, {9, 1, 0}},
				start: Point{0, 0},
				end:   Point{2, 1},
			},
			want{
				[]Point{{0, 0}, {1, 0}, {2, 0}, {2, 1}},
				nil,
			},
		},
		{
			"No path",
			args{
				maze:  [][]int{{1, 2, 0}, {0, 0, 1}, {9, 1, 0}},
				start: Point{0, 0},
				end:   Point{2, 1},
			},
			want{
				path: nil,
				err:  ErrNoPath,
			},
		},
		{
			"Maze with 2 path",
			args{
				maze:  [][]int{{1, 5, 1}, {1, 0, 1}, {1, 9, 1}},
				start: Point{0, 0},
				end:   Point{2, 2},
			},
			want{
				path: []Point{{0, 0}, {0, 1}, {0, 2}, {1, 2}, {2, 2}},
				err:  nil,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got, err := findShortestPath(tt.args.maze, tt.args.start, tt.args.end); !reflect.DeepEqual(got, tt.want.path) || !reflect.DeepEqual(err, tt.want.err) {
				t.Errorf("findShortestPath() = %v %v, want %v %v", got, err, tt.want.path, tt.want.err)
			}
		})
	}
}

func Test_isPointValid(t *testing.T) {
	type labyrinth struct {
		rows    int
		columns int
		maze    [][]int
	}

	type args struct {
		p    Point
		maze labyrinth
	}

	maze := labyrinth{
		rows:    3,
		columns: 3,
		maze:    [][]int{{1, 2, 0}, {2, 0, 1}, {9, 0, 1}},
	}

	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			"Valid points",
			args{
				p:    Point{0, 0},
				maze: maze,
			},
			true,
		},
		{
			"Point unreachable by bound",
			args{
				p:    Point{1, 1},
				maze: maze,
			},
			false,
		},
		{
			"Point is out of bound (x < 0)",
			args{
				p:    Point{-1, 0},
				maze: maze,
			},
			false,
		},
		{
			"Point is out of bound (x > rows)",
			args{
				p:    Point{maze.rows + 1, 0},
				maze: maze,
			},
			false,
		},
		{
			"Point is out of bound (y < 0)",
			args{
				p:    Point{0, -1},
				maze: maze,
			},
			false,
		},
		{
			"Point is out of bound (y > columns)",
			args{
				p:    Point{0, maze.columns + 1},
				maze: maze,
			},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := isPointValid(tt.args.p, tt.args.maze.rows, tt.args.maze.columns, tt.args.maze.maze); got != tt.want {
				t.Errorf("isPointValid() = %v, want %v", got, tt.want)
			}
		})
	}
}
