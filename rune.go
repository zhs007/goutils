package goutils

const RuneStr = "0"

// Rune2Int
func Rune2Int(r rune) int {
	return int(byte(r) - RuneStr[0])
}
