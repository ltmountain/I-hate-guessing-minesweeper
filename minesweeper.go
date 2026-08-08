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
	array := make([]int, 0, size*size)
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
	if mineBoard[BoardVec2ToArrayVec2(vec2, boardSize.x)] == 1 {
		panic("boom")
	}
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

func CanSetFlags(gameBoard []int, boardSize Vec2) []int {
	flags := make([]int, 0, boardSize.x*boardSize.y)
	for i, v := range gameBoard {
		if v > 0 {
			closes := 0
			nearbyTiles := CanSearchPos(ArrayVec2ToBoardVec2(i, boardSize.x), boardSize, 1)
			for _, tileIndex := range nearbyTiles {
				if gameBoard[tileIndex] == -1 || gameBoard[tileIndex] == -2 {
					closes += 1
				}
			}
			if closes == v {
				for _, tileIndex := range nearbyTiles {
					if gameBoard[tileIndex] == -1 {
						flags = append(flags, tileIndex)
					}
				}
			}
		}
	}
	return flags
}

func CanOpenTiles(gameBoard []int, boardSize Vec2) []int {
	opens := make([]int, 0, boardSize.x*boardSize.y)
	for i, v := range gameBoard {
		if v > 0 {
			flags := 0
			nearbyTiles := CanSearchPos(ArrayVec2ToBoardVec2(i, boardSize.x), boardSize, 1)
			for _, tileIndex := range nearbyTiles {
				if gameBoard[tileIndex] == -2 {
					flags += 1
				}
			}
			if flags == v {
				for _, tileIndex := range nearbyTiles {
					if gameBoard[tileIndex] == -1 {
						opens = append(opens, tileIndex)
					}
				}
			}
		}
	}
	return opens
}

func Deduction(gameBoard []int, boardSize Vec2) ([]int, []int) {
	flags := make([]int, 0, boardSize.x*boardSize.y)
	opens := make([]int, 0, boardSize.x*boardSize.y)
	for i, v := range gameBoard {
		if v > 0 {
			nearbyTilesOut := CanSearchPos(ArrayVec2ToBoardVec2(i, boardSize.x), boardSize, 2)
			for _, tileIndexOut := range nearbyTilesOut {
				if gameBoard[tileIndexOut] > 0 {
					nearbyTilesIn := CanSearchPos(ArrayVec2ToBoardVec2(tileIndexOut, boardSize.x), boardSize, 1)
					nearbyClosedTilesIn := 0
					countflagTilesOut := 0
					for _, tileIndexIn := range nearbyTilesIn {
						switch gameBoard[tileIndexIn] {
						case -1:
							nearbyClosedTilesIn += 1
						case -2:
							countflagTilesOut += 1
						}
					}
					addedTempFlags := make([]int, 0, 8)
					for _, tileIndexIn := range nearbyTilesIn {
						if gameBoard[tileIndexIn] == -1 {
							gameBoard[tileIndexIn] = -3
							addedTempFlags = append(addedTempFlags, tileIndexIn)
						}
					}

					nearbyTiles := CanSearchPos(ArrayVec2ToBoardVec2(i, boardSize.x), boardSize, 1)
					countTempFlags := make([]int, 0, 8)
					nearbyClosedTiles := make([]int, 0, 8)
					countflagTiles := 0
					for _, tileIndex := range nearbyTiles {
						switch gameBoard[tileIndex] {
						case -3:
							countTempFlags = append(countTempFlags, tileIndex)
						case -1:
							nearbyClosedTiles = append(nearbyClosedTiles, tileIndex)
						case -2:
							countflagTiles += 1
						}
					}

					if slices.Equal(addedTempFlags, countTempFlags) && len(nearbyClosedTiles) > 0 && v-countflagTiles == gameBoard[tileIndexOut]-countflagTilesOut {
						for _, nearbyClosedTile := range nearbyClosedTiles {
							opens = append(opens, nearbyClosedTile)
						}
					}
					if len(countTempFlags) > 0 && len(nearbyClosedTiles) == (v-countflagTiles)-(gameBoard[tileIndexOut]-countflagTilesOut) {
						for _, nearbyClosedTile := range nearbyClosedTiles {
							flags = append(flags, nearbyClosedTile)
						}
					}
					for _, tileIndexIn := range nearbyTilesIn {
						if gameBoard[tileIndexIn] == -3 {
							gameBoard[tileIndexIn] = -1
						}
					}
				}
			}
		}
	}
	return flags, opens
}

func Solve(gameBoard []int, mineBoard []int, boardSize Vec2) {
	for {
		done := true
		for _, tileIndex := range CanSetFlags(gameBoard, boardSize) {
			if mineBoard[tileIndex] != 1 {
				panic("Something is wrong")
			}
			gameBoard[tileIndex] = -2
			done = false
		}
		for _, tileIndex := range CanOpenTiles(gameBoard, boardSize) {
			Open(ArrayVec2ToBoardVec2(tileIndex, boardSize.x), gameBoard, mineBoard, boardSize)
			done = false
		}
		flags, opens := Deduction(gameBoard, boardSize)
		for _, tileIndex := range flags {
			if mineBoard[tileIndex] != 1 {
				panic("Something is wrong")
			}
			gameBoard[tileIndex] = -2
			done = false
		}
		for _, tileIndex := range opens {
			Open(ArrayVec2ToBoardVec2(tileIndex, boardSize.x), gameBoard, mineBoard, boardSize)
			done = false
		}
		if done {
			break
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
	for _, y := range board {
		for _, x := range y {
			if x < 0 {
				fmt.Print(x)
			} else {
				fmt.Print(" ")
				fmt.Print(x)
			}
		}
		fmt.Print("\n")
	}
}

/*
-3 = temp flag
-2 = mine / flag
-1 = closed
0 = opend
0 < nearby mine count
*/
func main() {
	fmt.Println("start")
	var boardSize Vec2 = Vec2{x: 25, y: 25}
	var mineCount int = 100

	var startPos Vec2 = Vec2{x: rand.N(boardSize.x), y: rand.N(boardSize.y)}
	startPos = Vec2{x: 1, y: 1}
	var mineBoard []int = SetMine(startPos, boardSize, mineCount)
	var gameBoard []int = slices.Repeat([]int{-1}, boardSize.x*boardSize.y)
	Open(startPos, gameBoard, mineBoard, boardSize)

	Solve(gameBoard, mineBoard, boardSize)
	Render2D(ArrayToBoard(gameBoard, boardSize))
	for _, v := range gameBoard {
		if v == -1 {
			fmt.Println("fail")
			break
		}
	}
}
