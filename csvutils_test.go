package goutils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type csvdata struct {
	GameMod  string
	Tag      string
	Symbol   string
	TotalBet int
	X        []int
	TotalWin int
}

func newcsvdata() *csvdata {
	return &csvdata{
		X: make([]int, 5),
	}
}

func Test_LoadCSVFile(t *testing.T) {
	lst := []*csvdata{}

	LoadCSVFile("./unittestdata/test.csv", func(i int, row []string) bool {
		return i == 0
	}, func(i int, row []string, mapHeader map[int]string) error {
		curdata := newcsvdata()

		for j, str := range row {
			if mapHeader[j] == "gamemod" {
				curdata.GameMod = str
			} else if mapHeader[j] == "tag" {
				curdata.Tag = str
			} else if mapHeader[j] == "symbol" {
				curdata.Symbol = str
			} else if mapHeader[j] == "totalbet" {
				i64, err := String2Int64(str)
				assert.NoError(t, err)

				curdata.TotalBet = int(i64)
			} else if mapHeader[j] == "X1" {
				i64, _ := String2Int64(str)
				// assert.NoError(t, err)

				curdata.X[0] = int(i64)
			} else if mapHeader[j] == "X2" {
				i64, _ := String2Int64(str)
				// assert.NoError(t, err)

				curdata.X[1] = int(i64)
			} else if mapHeader[j] == "X3" {
				i64, _ := String2Int64(str)
				// assert.NoError(t, err)

				curdata.X[2] = int(i64)
			} else if mapHeader[j] == "X4" {
				i64, _ := String2Int64(str)
				// assert.NoError(t, err)

				curdata.X[3] = int(i64)
			} else if mapHeader[j] == "X5" {
				i64, _ := String2Int64(str)
				// assert.NoError(t, err)

				curdata.X[4] = int(i64)
			} else if mapHeader[j] == "totalwin" {
				i64, err := String2Int64(str)
				assert.NoError(t, err)

				curdata.TotalWin = int(i64)
			}
		}

		lst = append(lst, curdata)
		return nil
	})

	assert.Equal(t, len(lst), 19)

	assert.Equal(t, lst[1].GameMod, "bg")
	assert.Equal(t, lst[1].Tag, "")
	assert.Equal(t, lst[1].Symbol, "1")
	assert.Equal(t, lst[1].TotalBet, 300)
	assert.Equal(t, lst[1].X[0], 0)
	assert.Equal(t, lst[1].X[1], 0)
	assert.Equal(t, lst[1].X[2], 100)
	assert.Equal(t, lst[1].X[3], 100)
	assert.Equal(t, lst[1].X[4], 0)
	assert.Equal(t, lst[1].TotalWin, 200)

	t.Logf("Test_LoadCSVFile OK")
}
