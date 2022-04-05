package goutils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_AppendString(t *testing.T) {

	in := [][]string{
		{"abc", "efg", "hijklmn"},
		{"abc", "", "hijklmn"},
	}

	out := []string{
		"abcefghijklmn",
		"abchijklmn",
	}

	for i, v := range in {
		ret := AppendString(v...)
		if ret != out[i] {
			t.Fatalf("Test_AppendString AppendString \"%s\" != \"%s\" [ %+v ]",
				ret, out[i], in[i])
		}
	}

	t.Logf("Test_AppendString OK")
}

func Test_String2Int64(t *testing.T) {
	i64, err := String2Int64("123")
	assert.NoError(t, err)
	assert.Equal(t, i64, int64(123))

	i64, err = String2Int64("123.456")
	assert.NoError(t, err)
	assert.Equal(t, i64, int64(123))

	_, err = String2Int64("123.456,")
	assert.Error(t, err)

	t.Logf("Test_String2Int64 OK")
}

func Test_String2Float64(t *testing.T) {
	f64, err := String2Float64("123")
	assert.NoError(t, err)
	assert.Equal(t, f64, float64(123))

	f64, err = String2Float64("123.456")
	assert.NoError(t, err)
	assert.Equal(t, f64, 123.456)

	_, err = String2Float64("123.456,")
	assert.Error(t, err)

	t.Logf("Test_String2Float64 OK")
}

func Test_Int2StringWithArr(t *testing.T) {
	str := Int2StringWithArr(0, "abc")
	assert.Equal(t, str, "a")

	str = Int2StringWithArr(1, "abc")
	assert.Equal(t, str, "b")

	str = Int2StringWithArr(2, "abc")
	assert.Equal(t, str, "c")

	str = Int2StringWithArr(3, "abc")
	assert.Equal(t, str, "aa")

	str = Int2StringWithArr(4, "abc")
	assert.Equal(t, str, "ab")

	str = Int2StringWithArr(5, "abc")
	assert.Equal(t, str, "ac")

	str = Int2StringWithArr(6, "abc")
	assert.Equal(t, str, "ba")

	str = Int2StringWithArr(7, "abc")
	assert.Equal(t, str, "bb")

	str = Int2StringWithArr(8, "abc")
	assert.Equal(t, str, "bc")

	str = Int2StringWithArr(9, "abc")
	assert.Equal(t, str, "ca")

	str = Int2StringWithArr(10, "abc")
	assert.Equal(t, str, "cb")

	str = Int2StringWithArr(11, "abc")
	assert.Equal(t, str, "cc")

	str = Int2StringWithArr(12, "abc")
	assert.Equal(t, str, "aaa")

	t.Logf("Test_Int2StringWithArr OK")
}
