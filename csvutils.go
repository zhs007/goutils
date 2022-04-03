package goutils

import (
	"encoding/csv"
	"io"
	"os"

	"go.uber.org/zap"
)

type FuncIsCSVHeadRow func(i int, row []string) bool
type FuncProcCSVRow func(i int, row []string, mapHeader map[int]string) error

func LoadCSVFile(fn string, funcIsHeadRow FuncIsCSVHeadRow, funcProcCSVRow FuncProcCSVRow) error {
	csvFile, err := os.Open(fn)
	if err != nil {
		Error("LoadCSVFile:Open",
			zap.Error(err))

		return err
	}
	defer csvFile.Close()

	header := make(map[int]string)
	i := 0
	csvReader := csv.NewReader(csvFile)
	for {
		row, err := csvReader.Read()
		if err == io.EOF {
			break
		} else if err != nil {
			Error("LoadCSVFile:Read",
				zap.Int("i", i),
				zap.Error(err))

			return err
		}

		if funcIsHeadRow(i, row) {
			for col, v := range row {
				header[col] = v
			}
		} else {
			err = funcProcCSVRow(i, row, header)
			if err != nil {
				Error("LoadCSVFile:funcProcCSVRow",
					zap.Int("i", i),
					zap.Error(err))

				return err
			}
		}

		i++
	}

	return nil
}
