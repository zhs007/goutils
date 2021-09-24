package goutils

import (
	"bytes"

	"github.com/buger/jsonparser"
	"go.uber.org/zap"
)

// GetJsonString - number to string
func GetJsonString(data []byte, keys ...string) (val string, err error) {
	v, t, _, e := jsonparser.Get(data, keys...)

	if e != nil {
		if e != jsonparser.KeyPathNotFoundError {
			return "", e
		}

		return "", nil
	}

	if t == jsonparser.Null {
		return "", nil
	}

	if t == jsonparser.Number {
		return string(v), nil
		// if strings.Contains(string(v), ".") {
		// 	nf, err := jsonparser.ParseFloat(v)
		// 	if err != nil {
		// 		return "", err
		// 	}

		// 	str := strconv.FormatFloat(nf, 'E', -1, 64)

		// 	return str, nil
		// }

		// iv, err := jsonparser.ParseInt(v)
		// if err != nil {
		// 	return "", err
		// }

		// str := strconv.FormatInt(iv, 10)

		// return str, nil
	}

	if t != jsonparser.String {
		return "", ErrInvalidJsonString
	}

	// If no escapes return raw content
	if bytes.IndexByte(v, '\\') == -1 {
		return string(v), nil
	}

	return jsonparser.ParseString(v)
}

func GetJsonInt(data []byte, keys ...string) (val int64, err error) {
	v, t, _, e := jsonparser.Get(data, keys...)

	if e != nil {
		if e != jsonparser.KeyPathNotFoundError {
			return 0, e
		}

		return 0, nil
	}

	if t == jsonparser.Null {
		return 0, nil
	}

	if t == jsonparser.String {
		if len(v) == 0 {
			return 0, nil
		}

		// If no escapes return raw content
		if bytes.IndexByte(v, '\\') == -1 {
			n, err := String2Int64(string(v))
			if err != nil {
				return 0, err
			}

			return n, nil
		}

		s, err := jsonparser.ParseString(v)
		if err != nil {
			return 0, err
		}

		n, err := String2Int64(s)
		if err != nil {
			return 0, err
		}

		return n, nil
	}

	if t != jsonparser.Number {
		return 0, ErrInvalidJsonInt
	}

	return String2Int64(string(v))
}

func GetJsonFloat(data []byte, keys ...string) (val float64, err error) {
	v, t, _, e := jsonparser.Get(data, keys...)

	if e != nil {
		if e != jsonparser.KeyPathNotFoundError {
			return 0, e
		}

		return 0, nil
	}

	if t == jsonparser.Null {
		return 0, nil
	}

	if t == jsonparser.String {
		if len(v) == 0 {
			return 0, nil
		}

		// If no escapes return raw content
		if bytes.IndexByte(v, '\\') == -1 {
			n, err := String2Float64(string(v))
			if err != nil {
				return 0, err
			}

			return n, nil
		}

		s, err := jsonparser.ParseString(v)
		if err != nil {
			return 0, err
		}

		n, err := String2Float64(s)
		if err != nil {
			return 0, err
		}

		return n, nil
	}

	if t != jsonparser.Number {
		return 0, ErrInvalidJsonInt
	}

	return String2Float64(string(v))
}

func GetJsonArrayEachInt(value1 []byte, dataType1 jsonparser.ValueType, offset1 int, err1 error) (int64, error) {
	if err1 != nil {
		if err1 != jsonparser.KeyPathNotFoundError {
			return 0, err1
		}

		return 0, nil
	}

	if dataType1 == jsonparser.Null {
		return 0, nil
	}

	if dataType1 == jsonparser.String {
		if len(value1) == 0 {
			return 0, nil
		}

		// If no escapes return raw content
		if bytes.IndexByte(value1, '\\') == -1 {
			n, err := String2Int64(string(value1))
			if err != nil {
				return 0, err
			}

			return n, nil
		}

		s, err := jsonparser.ParseString(value1)
		if err != nil {
			return 0, err
		}

		n, err := String2Int64(s)
		if err != nil {
			return 0, err
		}

		return n, nil
	}

	if dataType1 != jsonparser.Number {
		return 0, ErrInvalidJsonInt
	}

	return String2Int64(string(value1))
}

func GetJsonIntArr(data []byte, key string) ([]int, error) {
	arr := []int{}

	offset, err := jsonparser.ArrayEach(data, func(value1 []byte, dataType1 jsonparser.ValueType, offset1 int, err1 error) {
		if err1 != nil {
			if err1 != jsonparser.KeyPathNotFoundError {
				Error("GetJsonIntArr:ArrayEach:func",
					zap.String("key", key),
					zap.Int("offset", offset1),
					zap.Error(err1))

				return
			}

			return
		}

		cv, err2 := GetJsonArrayEachInt(value1, dataType1, offset1, err1)
		if err2 != nil {
			Error("GetJsonIntArr:ArrayEach:func:GetJsonArrayEachInt",
				zap.Int("offset", offset1),
				zap.Error(err2))

			return
		}

		arr = append(arr, int(cv))
	}, key)
	if err != nil {
		if err != jsonparser.KeyPathNotFoundError {
			Error("GetJsonIntArr:ArrayEach",
				zap.Int("offset", offset),
				zap.Error(err))

			return nil, err
		}

		return nil, nil
	}

	return arr, nil
}

func GetJsonInt64Arr(data []byte, key string) ([]int64, error) {
	arr := []int64{}

	offset, err := jsonparser.ArrayEach(data, func(value1 []byte, dataType1 jsonparser.ValueType, offset1 int, err1 error) {
		if err1 != nil {
			if err1 != jsonparser.KeyPathNotFoundError {
				Error("GetJsonInt64Arr:ArrayEach:func",
					zap.String("key", key),
					zap.Int("offset", offset1),
					zap.Error(err1))

				return
			}

			return
		}

		cv, err2 := GetJsonArrayEachInt(value1, dataType1, offset1, err1)
		if err2 != nil {
			Error("GetJsonInt64Arr:ArrayEach:func:GetJsonArrayEachInt",
				zap.Int("offset", offset1),
				zap.Error(err2))

			return
		}

		arr = append(arr, cv)
	}, key)
	if err != nil {
		if err != jsonparser.KeyPathNotFoundError {
			Error("GetJsonInt64Arr:ArrayEach",
				zap.Int("offset", offset),
				zap.Error(err))

			return nil, err
		}

		return nil, nil
	}

	return arr, nil
}
