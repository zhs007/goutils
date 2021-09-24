package goutils

import (
	"testing"

	"github.com/buger/jsonparser"
	"github.com/stretchr/testify/assert"
)

func Test_HasJsonKey(t *testing.T) {
	isok := HasJsonKey([]byte(`{"abc":"123"}`), "abc")
	assert.Equal(t, isok, true)

	isok = HasJsonKey([]byte(`{"abc":123}`), "ab")
	assert.Equal(t, isok, false)

	isok = HasJsonKey([]byte(`{"abc":null}`), "abc")
	assert.Equal(t, isok, true)

	isok = HasJsonKey([]byte(`{"abc":null}`), "ab")
	assert.Equal(t, isok, false)

	t.Logf("Test_HasJsonKey OK")
}

func Test_GetJsonString(t *testing.T) {
	s, isok, err := GetJsonString([]byte(`{"abc":"123"}`), "abc")
	assert.NoError(t, err)
	assert.Equal(t, isok, true)
	assert.Equal(t, s, "123")

	s, isok, err = GetJsonString([]byte(`{"abc":123}`), "abc")
	assert.NoError(t, err)
	assert.Equal(t, isok, true)
	assert.Equal(t, s, "123")

	s, isok, err = GetJsonString([]byte(`{"abc":123.456}`), "abc")
	assert.NoError(t, err)
	assert.Equal(t, isok, true)
	assert.Equal(t, s, "123.456")

	s, isok, err = GetJsonString([]byte(`{"abc":null}`), "abc")
	assert.NoError(t, err)
	assert.Equal(t, isok, false)
	assert.Equal(t, s, "")

	t.Logf("Test_GetJsonString OK")
}

func Test_GetJsonInt(t *testing.T) {
	i64, isok, err := GetJsonInt([]byte(`{"abc":"123"}`), "abc")
	assert.NoError(t, err)
	assert.Equal(t, isok, true)
	assert.Equal(t, i64, int64(123))

	i64, isok, err = GetJsonInt([]byte(`{"abc":123}`), "abc")
	assert.NoError(t, err)
	assert.Equal(t, isok, true)
	assert.Equal(t, i64, int64(123))

	i64, isok, err = GetJsonInt([]byte(`{"abc":123.456}`), "abc")
	assert.NoError(t, err)
	assert.Equal(t, isok, true)
	assert.Equal(t, i64, int64(123))

	i64, isok, err = GetJsonInt([]byte(`{"abc":"123.456"}`), "abc")
	assert.NoError(t, err)
	assert.Equal(t, isok, true)
	assert.Equal(t, i64, int64(123))

	i64, isok, err = GetJsonInt([]byte(`{"abc":""}`), "abc")
	assert.NoError(t, err)
	assert.Equal(t, isok, false)
	assert.Equal(t, i64, int64(0))

	i64, isok, err = GetJsonInt([]byte(`{"abc":null}`), "abc")
	assert.NoError(t, err)
	assert.Equal(t, isok, false)
	assert.Equal(t, i64, int64(0))

	t.Logf("Test_GetJsonInt OK")
}

func Test_GetJsonFloat(t *testing.T) {
	f64, isok, err := GetJsonFloat([]byte(`{"abc":"123"}`), "abc")
	assert.NoError(t, err)
	assert.Equal(t, isok, true)
	assert.Equal(t, f64, float64(123))

	f64, isok, err = GetJsonFloat([]byte(`{"abc":123}`), "abc")
	assert.NoError(t, err)
	assert.Equal(t, isok, true)
	assert.Equal(t, f64, float64(123))

	f64, isok, err = GetJsonFloat([]byte(`{"abc":123.456}`), "abc")
	assert.NoError(t, err)
	assert.Equal(t, isok, true)
	assert.Equal(t, f64, float64(123.456))

	f64, isok, err = GetJsonFloat([]byte(`{"abc":"123.456"}`), "abc")
	assert.NoError(t, err)
	assert.Equal(t, isok, true)
	assert.Equal(t, f64, float64(123.456))

	f64, isok, err = GetJsonFloat([]byte(`{"abc":""}`), "abc")
	assert.NoError(t, err)
	assert.Equal(t, isok, false)
	assert.Equal(t, f64, float64(0))

	f64, isok, err = GetJsonFloat([]byte(`{"abc":null}`), "abc")
	assert.NoError(t, err)
	assert.Equal(t, isok, false)
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

func Test_GetJsonObjectArr(t *testing.T) {
	arr1 := []int64{}
	err := GetJsonObjectArr([]byte(`{"abc":[{"a":1},{"a":2},{"a":3}]}`), "abc", func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
		if dataType == jsonparser.Object {
			a, isok, err := GetJsonInt(value, "a")
			assert.NoError(t, err)

			if isok {
				arr1 = append(arr1, a)
			}
		}
	})
	assert.NoError(t, err)
	assert.Equal(t, len(arr1), 3)
	assert.Equal(t, arr1[0], int64(1))
	assert.Equal(t, arr1[1], int64(2))
	assert.Equal(t, arr1[2], int64(3))

	arr2 := []int64{}
	err = GetJsonObjectArr([]byte(`{"abc":[{"a":1},{"a":2},{"a":3}]}`), "ab", func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
		if dataType == jsonparser.Object {
			a, isok, err := GetJsonInt(value, "a")
			assert.NoError(t, err)

			if isok {
				arr2 = append(arr2, a)
			}
		}
	})
	assert.NoError(t, err)
	assert.Nil(t, err)
	assert.Equal(t, len(arr2), 0)

	arr3 := []int64{}
	err = GetJsonObjectArr([]byte(`{"abc":[{"a":1},{"ab":2},{"a":3}]}`), "abc", func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
		if dataType == jsonparser.Object {
			a, isok, err := GetJsonInt(value, "a")
			assert.NoError(t, err)

			if isok {
				arr3 = append(arr3, a)
			}
		}
	})
	assert.NoError(t, err)
	assert.Nil(t, err)
	assert.Equal(t, len(arr3), 2)
	assert.Equal(t, arr3[0], int64(1))
	assert.Equal(t, arr3[1], int64(3))

	t.Logf("Test_GetJsonIntArr OK")
}

func Test_GetJsonIntArr3(t *testing.T) {
	arr1, err := GetJsonIntArr3([]byte(`{"abc":[[[1,2,3],[4,5]],[[6,7],[8]]]}`), "abc")
	assert.NoError(t, err)
	assert.Equal(t, len(arr1), 2)
	assert.Equal(t, len(arr1[0]), 2)
	assert.Equal(t, len(arr1[1]), 2)
	assert.Equal(t, len(arr1[0][0]), 3)
	assert.Equal(t, len(arr1[0][1]), 2)
	assert.Equal(t, len(arr1[1][0]), 2)
	assert.Equal(t, len(arr1[1][1]), 1)
	assert.Equal(t, arr1[0][0][0], 1)
	assert.Equal(t, arr1[0][0][1], 2)
	assert.Equal(t, arr1[0][0][2], 3)
	assert.Equal(t, arr1[0][1][0], 4)
	assert.Equal(t, arr1[0][1][1], 5)
	assert.Equal(t, arr1[1][0][0], 6)
	assert.Equal(t, arr1[1][0][1], 7)
	assert.Equal(t, arr1[1][1][0], 8)

	t.Logf("Test_GetJsonIntArr2 OK")
}

func Test_GetJsonInt64Arr3(t *testing.T) {
	arr1, err := GetJsonInt64Arr3([]byte(`{"abc":[[[1,2,3],[4,5]],[[6,7],[8]]]}`), "abc")
	assert.NoError(t, err)
	assert.Equal(t, len(arr1), 2)
	assert.Equal(t, len(arr1[0]), 2)
	assert.Equal(t, len(arr1[1]), 2)
	assert.Equal(t, len(arr1[0][0]), 3)
	assert.Equal(t, len(arr1[0][1]), 2)
	assert.Equal(t, len(arr1[1][0]), 2)
	assert.Equal(t, len(arr1[1][1]), 1)
	assert.Equal(t, arr1[0][0][0], int64(1))
	assert.Equal(t, arr1[0][0][1], int64(2))
	assert.Equal(t, arr1[0][0][2], int64(3))
	assert.Equal(t, arr1[0][1][0], int64(4))
	assert.Equal(t, arr1[0][1][1], int64(5))
	assert.Equal(t, arr1[1][0][0], int64(6))
	assert.Equal(t, arr1[1][0][1], int64(7))
	assert.Equal(t, arr1[1][1][0], int64(8))

	t.Logf("Test_GetJsonInt64Arr3 OK")
}
