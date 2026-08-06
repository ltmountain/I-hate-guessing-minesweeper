package main

import (
	"fmt"
	"math/rand/v2"
	"slices"
)

type Vec2 struct {
	x int
	y int
}

func SetMine(startPos Vec2, boardSize Vec2, mineCount int) []int {
	tileCount := boardSize.x * boardSize.y
	array := make([]int, 0, tileCount)
	dontTouch := CanSearchPos(startPos, boardSize, 1)
	dontTouch = append(dontTouch, BoardVec2ToArrayVec2(startPos, boardSize.x))
	for i := 0; i < mineCount; i++ {
		if !slices.Contains(dontTouch, i) {
			array = append(array, 1)
		} else {
			mineCount += 1
			array = append(array, 0)
		}
	}
	for i := mineCount; i < tileCount; i++ {
		array = append(array, 0)
	}
	FisherYates(array, dontTouch)

	return array
}

func CanSearchPos(pos Vec2, boardSize Vec2, size int) []int {
	init := Vec2{x: pos.x - size, y: pos.y - size}
	if pos.x-size < 0 {
		init.x = 0
	}
	if pos.y-size < 0 {
		init.y = 0
	}
	var array []int
	for y := init.y; y < boardSize.y && y <= pos.y+size; y++ {
		for x := init.x; x < boardSize.x && x <= pos.x+size; x++ {
			if x != pos.x || y != pos.y {
				array = append(array, BoardVec2ToArrayVec2(Vec2{x: x, y: y}, boardSize.x))
			}
		}
	}
	return array
}

func Open(vec2 Vec2, gameBoard []int, mineBoard []int, boardSize Vec2) {
	mineCount := 0
	canSearchPos := CanSearchPos(vec2, boardSize, 1)
	for _, v := range canSearchPos {
		if mineBoard[v] == 1 {
			mineCount += 1
		}
	}
	gameBoard[BoardVec2ToArrayVec2(vec2, boardSize.x)] = mineCount
	if mineCount == 0 {
		for _, v := range canSearchPos {
			if gameBoard[v] != 0 {
				Open(ArrayVec2ToBoardVec2(v, boardSize.x), gameBoard, mineBoard, boardSize)
			}
		}
	}
}

func FisherYates(array []int, dontTouch []int) {
	for i := len(array) - 1; i > 0; i-- {
		j := rand.N(i + 1)
		if !slices.Contains(dontTouch, i) && !slices.Contains(dontTouch, j) {
			array[j], array[i] = array[i], array[j]
		}
	}
}

func BoardVec2ToArrayVec2(vec2 Vec2, xsize int) int {
	return vec2.y*xsize + vec2.x
}

func ArrayVec2ToBoardVec2(len int, xsize int) Vec2 {
	return Vec2{x: len % xsize, y: len / xsize}
}

func ArrayToBoard(array []int, size Vec2) [][]int {
	board := make([][]int, 0, size.y)
	for y := 0; y < size.y; y++ {
		board = append(board, make([]int, 0, size.x))
		for x := 0; x < size.x; x++ {
			board[y] = append(board[y], array[BoardVec2ToArrayVec2(Vec2{x: x, y: y}, size.x)])
		}
	}
	return board
}

func Render2D(board [][]int) {
	for _, v := range board {
		fmt.Println(v)
	}
}

/*
-2 = mine
-1 = closed
0 = opend
0 < nearby mine count
*/
func main() {
	var boardSize Vec2 = Vec2{x: 25, y: 25}
	var mineCount int = 100

	var startPos Vec2 = Vec2{x: rand.N(boardSize.x), y: rand.N(boardSize.y)}
	startPos = Vec2{x: 1, y: 1}
	var array []int = SetMine(startPos, boardSize, mineCount)
	var gameBoard []int = slices.Repeat([]int{-1}, boardSize.x*boardSize.y)
	Open(startPos, gameBoard, array, boardSize)

	Render2D(ArrayToBoard(array, boardSize))
	Render2D(ArrayToBoard(gameBoard, boardSize))
}

//TODO 렌더링만 2차원 나머지는 1차원으로 통일
