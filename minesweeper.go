package main

import (
	"fmt"
	"math/rand/v2"
)

type Vec2 struct {
	x int
	y int
}

type Board struct {
	board [][]int
	size  Vec2
}

func CanSearchPos(pos Vec2, boardSize Vec2, size int) []Vec2 {
	init := Vec2{x: pos.x - size, y: pos.y - size}
	if pos.x-size < 0 {
		init.x = 0
	}
	if pos.y-size < 0 {
		init.y = 0
	}
	var array []Vec2
	for y := init.y; y < boardSize.y && y <= pos.y+size; y++ {
		for x := init.x; x < boardSize.x && x <= pos.x+size; x++ {
			if x != pos.x || y != pos.y {
				array = append(array, Vec2{x: x, y: y})
			}
		}
	}
	return array
}

func FisherYates(startPos Vec2, array []int) {
	for i := len(array) - 1; i > 0; i-- {
		j := rand.N(i + 1)
		array[j], array[i] = array[i], array[j]
	}
}

func SetMine(startPos Vec2, boardSize Vec2, mineCount int) []int {
	tileCount := boardSize.x * boardSize.y
	array := make([]int, 0, tileCount)
	for i := 0; i < tileCount; i++ {
		//지뢰 먼저 배치
		if i < mineCount {
			array = append(array, -2)
		} else {
			array = append(array, -1)
		}
	}
	FisherYates(startPos, array)

	return array
}

func BoardVec2ToArrayVec2(vec2 Vec2, xsize int) int {
	return vec2.y*xsize + vec2.x
}

func Open(pos Vec2, board Board) {

}

func ArrayToBoard(array []int, size Vec2) Board {
	board := Board{board: make([][]int, 0, size.y), size: size}
	for y := 0; y < size.y; y++ {
		board.board = append(board.board, make([]int, 0, size.x))
		for x := 0; x < size.x; x++ {
			board.board[y] = append(board.board[y], array[BoardVec2ToArrayVec2(Vec2{x: x, y: y}, size.x)])
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
	var array []int = SetMine(startPos, boardSize, mineCount)
	var board Board = ArrayToBoard(array, boardSize)

	Render2D(board.board)
	fmt.Println(CanSearchPos(Vec2{x: 2, y: 2}, Vec2{x: 3, y: 3}, 1))
}

//TODO 시작 지점 인근 지뢰 매설 금지
//TODO 렌더링만 2차원 나머지는 1차원으로 통일
