package main

import (
	"fmt"
)

var HOME_POSITIONS = [16]int{32, 33, 47, 48, 41, 42, 56, 57, 176, 177, 191, 192, 167, 168, 182, 183}
var WIN_POSITIONS = [4]int{111, 97, 113, 127}
var START_POSITIONS = [4]int{91, 23, 133, 201}
var START_TO_END_POSITIONS = generateStartToEndPositions()

var loopPositions = []int{90, 91, 92, 93, 94, 95, 96, 81, 66, 51, 36, 21, 6, 7, 8, 23, 38, 53, 68, 83, 98, 99, 100, 101, 102, 103, 104, 119, 134, 133, 132, 131, 130, 129, 126, 143, 158, 173, 188, 203, 218, 217, 216, 201, 186, 171, 156, 141, 128, 125, 124, 123, 122, 121, 120, 105}
var startToColored = [4][2]int{{91, 105}, {23, 7}, {133, 119}, {201, 217}}
var coloredPaths = [4][]int{
	{106, 107, 108, 109, 110, 111},
	{22, 37, 52, 67, 82, 97},
	{118, 117, 116, 115, 114, 113},
	{202, 187, 172, 157, 142, 127},
}


func generateStartToEndPositions() [16][]int {
 startToEndPositions := [16][]int{}

	for playerNumber := 0; playerNumber < 4; playerNumber++ {
        for i:= 0; i<4; i++{
            pawnNumber := playerNumber*4 + i
            homePath := []int{HOME_POSITIONS[pawnNumber]}
		toColoredPath := reslice(startToColored[playerNumber][0], startToColored[playerNumber][1])
		coloredPath := coloredPaths[playerNumber]
		fullPath := append(append(homePath,toColoredPath...),  coloredPath...)
		startToEndPositions[pawnNumber] = fullPath
        }

	}

    return startToEndPositions

}
func findIndex(slice []int, value int) int {
	for i, item := range slice {
		if item == value {
			return i
		}
	}
	return -1
}

func GetPlayerNumber(pawnNumber int) int {
	return pawnNumber / 4
}

func reslice(startPos int, endPos int) []int {
	i1 := findIndex(loopPositions, startPos)
	i2 := findIndex(loopPositions, endPos)

	if i1 == -1 || i2 == -1 {
		fmt.Println("Invalid start or end position", startPos, endPos)
		return nil
	}

	resliced := []int{}
	for i := i1; i < 1000; i++ {
		rotatedIndex := i % len(loopPositions)
		resliced = append(resliced, loopPositions[rotatedIndex])
		if rotatedIndex == i2 {
			break
		}
	}

	return resliced
}


	//if pawn is on other player's pawn, move it to home

func ErrorIfInvalidMove(pawnNumber int, position int, dice int) error{
    //If pawn is home, cant move unless dice is 6
    if position == HOME_POSITIONS[pawnNumber] && dice != 6{
        return fmt.Errorf("invalid move, pawn is in home and dice is not 6")
    }

    //If pawn is at last position, dont move it
    if position == WIN_POSITIONS[GetPlayerNumber(pawnNumber)]{
        return fmt.Errorf("invalid move, pawn is at last position")
    }
    return nil
}
	
func getMovementPath(pawnNumber, position int, dice int)([]int, error){
    startIndex := findIndex(START_TO_END_POSITIONS[pawnNumber], position)
    if startIndex == -1 {
        return nil, fmt.Errorf("invalid position %d", position)
    }

    path := []int{}
    for i := 0; i <= dice; i++ {
        curIndex := startIndex + i
        if curIndex >= len(START_TO_END_POSITIONS[pawnNumber]){
            break
        }
        path = append(path, START_TO_END_POSITIONS[pawnNumber][curIndex])
    }
    return path,nil
}

func GetBackToHomePath(pawnNumber, position int)([]int, error){
    startIndex := findIndex(START_TO_END_POSITIONS[pawnNumber], position)

    if startIndex == -1 {
        return nil, fmt.Errorf("invalid position %d", position)
    }

    path := []int{}

    for i := startIndex; i >=0; i-- {
        path = append(path, START_TO_END_POSITIONS[pawnNumber][i])
    }
    return path,nil
}






func main() {
    // fmt.Println(GetMovementPath(0, 92, 2))
    // fmt.Println(GetMovementPath(0, 32, 2))
    // fmt.Println(GetMovementPath(0, 105, 2))
    // fmt.Println(GetMovementPath(0, 110, 20))
    // fmt.Println(GetBackToHomePath(0, 110))
    // fmt.Println(GetBackToHomePath(0, 163))
    // fmt.Println(GetBackToHomePath(11, 133))

}
