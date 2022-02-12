package goutils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_LoadVersion(t *testing.T) {
	str, err := LoadVersion("./VERSION")
	assert.NoError(t, err)
	assert.Equal(t, str, Version)

	t.Logf("Test_LoadVersion OK")
}

func Test_ParseVersion(t *testing.T) {
	_, err0 := ParseVersion("1.2.03")
	assert.Error(t, err0)

	_, err1 := ParseVersion("v1.2.")
	assert.Error(t, err1)

	vobj2, err2 := ParseVersion("v1.20.03")
	assert.NoError(t, err2)
	assert.Equal(t, vobj2.Major, 1)
	assert.Equal(t, vobj2.Minor, 20)
	assert.Equal(t, vobj2.Patch, 3)
	assert.Equal(t, vobj2.ToString(), "v1.20.3")

	vobj2.IncPatch()
	assert.Equal(t, vobj2.Major, 1)
	assert.Equal(t, vobj2.Minor, 20)
	assert.Equal(t, vobj2.Patch, 4)
	assert.Equal(t, vobj2.ToString(), "v1.20.4")

	t.Logf("Test_LoadVersion OK")
}

func Test_BuildVersionFile(t *testing.T) {
	vobj0, err0 := ParseVersion("v1.20.03")
	assert.NoError(t, err0)

	err := BuildVersionFile("./unittestdata/version.txt", "./unittestdata/version.temp", vobj0)
	assert.NoError(t, err)

	isok := IsSameFile("./unittestdata/version.txt", "./unittestdata/version.cmp")
	assert.Equal(t, isok, true)

	t.Logf("Test_BuildVersionFile OK")
}
