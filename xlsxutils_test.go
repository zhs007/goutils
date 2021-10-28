package goutils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_Pos2Cell(t *testing.T) {
	assert.Equal(t, Pos2Cell(0, 0), "A1")
	assert.Equal(t, Pos2Cell(1, 0), "B1")
	assert.Equal(t, Pos2Cell(0, 1), "A2")
	assert.Equal(t, Pos2Cell(1, 1), "B2")
	assert.Equal(t, Pos2Cell(25, 0), "Z1")
	assert.Equal(t, Pos2Cell(26, 0), "AA1")
	assert.Equal(t, Pos2Cell(0, 10), "A11")
	assert.Equal(t, Pos2Cell(51, 0), "AZ1")
	assert.Equal(t, Pos2Cell(52, 0), "BA1")

	t.Logf("Test_Pos2Cell OK")
}
