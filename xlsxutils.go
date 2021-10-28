package goutils

import "fmt"

const cellName string = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"

// Pos2Cell - (0, 0) -> A1, (1, 0) -> B1, (0, 1) -> A2
func Pos2Cell(x, y int) string {
	if x < 0 || y < 0 {
		return ""
	}

	if x < len(cellName) {
		return fmt.Sprintf("%c%v", cellName[x], y+1)
	}

	t := int(x / len(cellName))
	if t <= len(cellName) {
		t1 := x % len(cellName)

		return fmt.Sprintf("%c%c%v", cellName[t-1], cellName[t1], y+1)
	}

	return ""
}
