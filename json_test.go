package goutils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_GetJsonString(t *testing.T) {
	s, err := GetJsonString([]byte(`{"abc":"123"}`), "abc")
	assert.NoError(t, err)
	assert.Equal(t, s, "123")

	s, err = GetJsonString([]byte(`{"abc":123}`), "abc")
	assert.NoError(t, err)
	assert.Equal(t, s, "123")

	s, err = GetJsonString([]byte(`{"abc":123.456}`), "abc")
	assert.NoError(t, err)
	assert.Equal(t, s, "123.456")

	s, err = GetJsonString([]byte(`{"abc":null}`), "abc")
	assert.NoError(t, err)
	assert.Equal(t, s, "")

	t.Logf("Test_GetJsonString OK")
}

func Test_GetJsonInt(t *testing.T) {
	i64, err := GetJsonInt([]byte(`{"abc":"123"}`), "abc")
	assert.NoError(t, err)
	assert.Equal(t, i64, int64(123))

	i64, err = GetJsonInt([]byte(`{"abc":123}`), "abc")
	assert.NoError(t, err)
	assert.Equal(t, i64, int64(123))

	i64, err = GetJsonInt([]byte(`{"abc":123.456}`), "abc")
	assert.NoError(t, err)
	assert.Equal(t, i64, int64(123))

	i64, err = GetJsonInt([]byte(`{"abc":"123.456"}`), "abc")
	assert.NoError(t, err)
	assert.Equal(t, i64, int64(123))

	i64, err = GetJsonInt([]byte(`{"abc":""}`), "abc")
	assert.NoError(t, err)
	assert.Equal(t, i64, int64(0))

	i64, err = GetJsonInt([]byte(`{"abc":null}`), "abc")
	assert.NoError(t, err)
	assert.Equal(t, i64, int64(0))

	t.Logf("Test_GetJsonInt OK")
}

func Test_GetJsonFloat(t *testing.T) {
	f64, err := GetJsonFloat([]byte(`{"abc":"123"}`), "abc")
	assert.NoError(t, err)
	assert.Equal(t, f64, float64(123))

	f64, err = GetJsonFloat([]byte(`{"abc":123}`), "abc")
	assert.NoError(t, err)
	assert.Equal(t, f64, float64(123))

	f64, err = GetJsonFloat([]byte(`{"abc":123.456}`), "abc")
	assert.NoError(t, err)
	assert.Equal(t, f64, float64(123.456))

	f64, err = GetJsonFloat([]byte(`{"abc":"123.456"}`), "abc")
	assert.NoError(t, err)
	assert.Equal(t, f64, float64(123.456))

	f64, err = GetJsonFloat([]byte(`{"abc":""}`), "abc")
	assert.NoError(t, err)
	assert.Equal(t, f64, float64(0))

	f64, err = GetJsonFloat([]byte(`{"abc":null}`), "abc")
	assert.NoError(t, err)
	assert.Equal(t, f64, float64(0))

	t.Logf("Test_GetJsonFloat OK")
}

func Test_GetJsonIntArr(t *testing.T) {
	arr1, err := GetJsonIntArr([]byte(`{"abc":[1,2,3]}`), "abc")
	assert.NoError(t, err)
	assert.Equal(t, len(arr1), 3)
	assert.Equal(t, arr1[0], 1)
	assert.Equal(t, arr1[1], 2)
	assert.Equal(t, arr1[2], 3)

	arr2, err := GetJsonIntArr([]byte(`{"abc":[1,"2",3.8]}`), "abc")
	assert.NoError(t, err)
	assert.Equal(t, len(arr2), 3)
	assert.Equal(t, arr2[0], 1)
	assert.Equal(t, arr2[1], 2)
	assert.Equal(t, arr2[2], 3)

	arr3, err := GetJsonIntArr([]byte(`{"abc":[1,"2",3.8]}`), "ab")
	assert.NoError(t, err)
	assert.Nil(t, arr3)

	t.Logf("Test_GetJsonIntArr OK")
}

func Test_GetJsonInt64Arr(t *testing.T) {
	arr1, err := GetJsonInt64Arr([]byte(`{"abc":[1,2,3]}`), "abc")
	assert.NoError(t, err)
	assert.Equal(t, len(arr1), 3)
	assert.Equal(t, arr1[0], int64(1))
	assert.Equal(t, arr1[1], int64(2))
	assert.Equal(t, arr1[2], int64(3))

	arr2, err := GetJsonInt64Arr([]byte(`{"abc":[1,"2",3.8]}`), "abc")
	assert.NoError(t, err)
	assert.Equal(t, len(arr2), 3)
	assert.Equal(t, arr2[0], int64(1))
	assert.Equal(t, arr2[1], int64(2))
	assert.Equal(t, arr2[2], int64(3))

	arr3, err := GetJsonInt64Arr([]byte(`{"abc":[1,"2",3.8]}`), "ab")
	assert.NoError(t, err)
	assert.Nil(t, arr3)

	t.Logf("Test_GetJsonInt64Arr OK")
}

func Test_GetJsonIntArr2(t *testing.T) {
	arr1, err := GetJsonIntArr2([]byte(`{"abc":[[1,2,3],[4,5]]}`), "abc")
	assert.NoError(t, err)
	assert.Equal(t, len(arr1), 2)
	assert.Equal(t, len(arr1[0]), 3)
	assert.Equal(t, len(arr1[1]), 2)
	assert.Equal(t, arr1[0][0], 1)
	assert.Equal(t, arr1[0][1], 2)
	assert.Equal(t, arr1[0][2], 3)
	assert.Equal(t, arr1[1][0], 4)
	assert.Equal(t, arr1[1][1], 5)

	arr2, err := GetJsonIntArr2([]byte(`{"abc":[1,"2",3.8]}`), "abc")
	assert.NoError(t, err)
	assert.Nil(t, arr2)

	arr3, err := GetJsonIntArr2([]byte(`{"abc":[1,"2",3.8]}`), "ab")
	assert.NoError(t, err)
	assert.Nil(t, arr3)

	arr5, err := GetJsonIntArr2([]byte(`{"abc":[[1,2,3.5],["4",5]]}`), "abc")
	assert.NoError(t, err)
	assert.Equal(t, len(arr5), 2)
	assert.Equal(t, len(arr5[0]), 3)
	assert.Equal(t, len(arr5[1]), 2)
	assert.Equal(t, arr5[0][0], 1)
	assert.Equal(t, arr5[0][1], 2)
	assert.Equal(t, arr5[0][2], 3)
	assert.Equal(t, arr5[1][0], 4)
	assert.Equal(t, arr5[1][1], 5)

	t.Logf("Test_GetJsonIntArr2 OK")
}

func Test_GetJsonInt64Arr2(t *testing.T) {
	arr1, err := GetJsonInt64Arr2([]byte(`{"abc":[[1,2,3],[4,5]]}`), "abc")
	assert.NoError(t, err)
	assert.Equal(t, len(arr1), 2)
	assert.Equal(t, len(arr1[0]), 3)
	assert.Equal(t, len(arr1[1]), 2)
	assert.Equal(t, arr1[0][0], int64(1))
	assert.Equal(t, arr1[0][1], int64(2))
	assert.Equal(t, arr1[0][2], int64(3))
	assert.Equal(t, arr1[1][0], int64(4))
	assert.Equal(t, arr1[1][1], int64(5))

	arr2, err := GetJsonInt64Arr2([]byte(`{"abc":[1,"2",3.8]}`), "abc")
	assert.NoError(t, err)
	assert.Nil(t, arr2)

	arr3, err := GetJsonInt64Arr2([]byte(`{"abc":[1,"2",3.8]}`), "ab")
	assert.NoError(t, err)
	assert.Nil(t, arr3)

	arr5, err := GetJsonInt64Arr2([]byte(`{"abc":[[1,2,3.5],["4",5]]}`), "abc")
	assert.NoError(t, err)
	assert.Equal(t, len(arr5), 2)
	assert.Equal(t, len(arr5[0]), 3)
	assert.Equal(t, len(arr5[1]), 2)
	assert.Equal(t, arr5[0][0], int64(1))
	assert.Equal(t, arr5[0][1], int64(2))
	assert.Equal(t, arr5[0][2], int64(3))
	assert.Equal(t, arr5[1][0], int64(4))
	assert.Equal(t, arr5[1][1], int64(5))

	t.Logf("Test_GetJsonInt64Arr2 OK")
}
