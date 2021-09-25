package goutils

import (
	"bytes"

	"github.com/buger/jsonparser"
	"go.uber.org/zap"
)

func HasJsonKey(data []byte, keys ...string) bool {
	_, _, _, e := jsonparser.Get(data, keys...)
	if e != nil {
		return e != jsonparser.KeyPathNotFoundError
	}

	return true
}

// GetJsonString - number to string
func GetJsonString(data []byte, keys ...string) (string, bool, error) {
	v, t, _, e := jsonparser.Get(data, keys...)

	if e != nil {
		if e != jsonparser.KeyPathNotFoundError {
			return "", false, e
		}

		return "", false, nil
	}

	if t == jsonparser.Null {
		return "", false, nil
	}

	if t == jsonparser.Number {
		return string(v), true, nil
	}

	if t != jsonparser.String {
		return "", false, ErrInvalidJsonString
	}

	// If no escapes return raw content
	if bytes.IndexByte(v, '\\') == -1 {
		return string(v), true, nil
	}

	str, err := jsonparser.ParseString(v)
	if err != nil {
		return "", false, err
	}

	return str, true, nil
}

func GetJsonInt(data []byte, keys ...string) (int64, bool, error) {
	v, t, _, e := jsonparser.Get(data, keys...)

	if e != nil {
		if e != jsonparser.KeyPathNotFoundError {
			return 0, false, e
		}

		return 0, false, nil
	}

	if t == jsonparser.Null {
		return 0, false, nil
	}

	if t == jsonparser.String {
		if len(v) == 0 {
			return 0, false, nil
		}

		// If no escapes return raw content
		if bytes.IndexByte(v, '\\') == -1 {
			n, err := String2Int64(string(v))
			if err != nil {
				return 0, false, err
			}

			return n, true, nil
		}

		s, err := jsonparser.ParseString(v)
		if err != nil {
			return 0, false, err
		}

		n, err := String2Int64(s)
		if err != nil {
			return 0, false, err
		}

		return n, true, nil
	}

	if t != jsonparser.Number {
		return 0, false, ErrInvalidJsonInt
	}

	i64, err := String2Int64(string(v))
	if err != nil {
		return 0, false, err
	}

	return i64, true, nil
}

func GetJsonFloat(data []byte, keys ...string) (float64, bool, error) {
	v, t, _, e := jsonparser.Get(data, keys...)

	if e != nil {
		if e != jsonparser.KeyPathNotFoundError {
			return 0, false, e
		}

		return 0, false, nil
	}

	if t == jsonparser.Null {
		return 0, false, nil
	}

	if t == jsonparser.String {
		if len(v) == 0 {
			return 0, false, nil
		}

		// If no escapes return raw content
		if bytes.IndexByte(v, '\\') == -1 {
			n, err := String2Float64(string(v))
			if err != nil {
				return 0, false, err
			}

			return n, true, nil
		}

		s, err := jsonparser.ParseString(v)
		if err != nil {
			return 0, false, err
		}

		n, err := String2Float64(s)
		if err != nil {
			return 0, false, err
		}

		return n, true, nil
	}

	if t != jsonparser.Number {
		return 0, false, ErrInvalidJsonInt
	}

	f64, err := String2Float64(string(v))
	if err != nil {
		return 0, false, err
	}

	return f64, true, nil
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

func GetJsonIntArr(data []byte, keys ...string) ([]int, error) {
	arr := []int{}

	offset, err := jsonparser.ArrayEach(data, func(value1 []byte, dataType1 jsonparser.ValueType, offset1 int, err1 error) {
		if err1 != nil {
			if err1 != jsonparser.KeyPathNotFoundError {
				Error("GetJsonIntArr:ArrayEach:func",
					zap.Strings("keys", keys),
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
	}, keys...)
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

func GetJsonInt64Arr(data []byte, keys ...string) ([]int64, error) {
	arr := []int64{}

	offset, err := jsonparser.ArrayEach(data, func(value1 []byte, dataType1 jsonparser.ValueType, offset1 int, err1 error) {
		if err1 != nil {
			if err1 != jsonparser.KeyPathNotFoundError {
				Error("GetJsonInt64Arr:ArrayEach:func",
					zap.Strings("keys", keys),
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
	}, keys...)
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

func GetJsonObjectArr(data []byte, cb func(value []byte, dataType jsonparser.ValueType, offset int, err error), keys ...string) error {
	offset, err := jsonparser.ArrayEach(data, func(value1 []byte, dataType1 jsonparser.ValueType, offset1 int, err1 error) {
		if err1 != nil {
			if err1 != jsonparser.KeyPathNotFoundError {
				Error("GetJsonObjectArr:ArrayEach:func",
					zap.Strings("keys", keys),
					zap.Int("offset", offset1),
					zap.Error(err1))

				return
			}

			return
		}

		cb(value1, dataType1, offset1, err1)
	}, keys...)
	if err != nil {
		if err != jsonparser.KeyPathNotFoundError {
			Error("GetJsonObjectArr:ArrayEach",
				zap.Int("offset", offset),
				zap.Error(err))

			return err
		}
	}

	return nil
}

func GetJsonIntArr2(data []byte, keys ...string) ([][]int, error) {
	arr := [][]int{}

	offset, err := jsonparser.ArrayEach(data, func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
		if err != nil {
			if err != jsonparser.KeyPathNotFoundError {
				Error("GetJsonIntArr2:ArrayEach:func",
					zap.Int("offset", offset),
					zap.Error(err))

				return
			}
		}

		if dataType == jsonparser.Array {
			arr0 := []int{}

			offset3, err3 := jsonparser.ArrayEach(value, func(value2 []byte, dataType2 jsonparser.ValueType, offset2 int, err2 error) {
				if err2 != nil {
					if err != jsonparser.KeyPathNotFoundError {
						Error("GetJsonIntArr2:ArrayEach:func2",
							zap.Int("offset", offset2),
							zap.Error(err2))

						return
					}
				}

				cv, err5 := GetJsonArrayEachInt(value2, dataType2, offset2, err2)
				if err5 != nil {
					Error("GetJsonIntArr2:ArrayEach:func2:GetJsonArrayEachInt",
						zap.Int("offset", offset2),
						zap.Error(err2))

					return
				}

				arr0 = append(arr0, int(cv))
			})
			if err3 != nil {
				Error("GetJsonIntArr2:ArrayEach:func:ArrayEach",
					zap.Int("offset", offset3),
					zap.Error(err3))

				return
			}

			arr = append(arr, arr0)

			return
		}

		Error("GetJsonIntArr2:ArrayEach:func:dataType",
			zap.Int("offset", offset),
			zap.String("dataType", dataType.String()))
	}, keys...)
	if err != nil {
		if err != jsonparser.KeyPathNotFoundError {
			Error("GetJsonIntArr2:ArrayEach",
				zap.Strings("keys", keys),
				zap.Int("offset", offset),
				zap.Error(err))

			return nil, err
		}
	}

	if len(arr) > 0 {
		return arr, nil
	}

	return nil, nil
}

func GetJsonInt64Arr2(data []byte, keys ...string) ([][]int64, error) {
	arr := [][]int64{}

	offset, err := jsonparser.ArrayEach(data, func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
		if err != nil {
			if err != jsonparser.KeyPathNotFoundError {
				Error("GetJsonInt64Arr2:ArrayEach:func",
					zap.Int("offset", offset),
					zap.Error(err))

				return
			}
		}

		if dataType == jsonparser.Array {
			arr0 := []int64{}

			offset3, err3 := jsonparser.ArrayEach(value, func(value2 []byte, dataType2 jsonparser.ValueType, offset2 int, err2 error) {
				if err2 != nil {
					if err != jsonparser.KeyPathNotFoundError {
						Error("GetJsonInt64Arr2:ArrayEach:func2",
							zap.Int("offset", offset2),
							zap.Error(err2))

						return
					}
				}

				cv, err5 := GetJsonArrayEachInt(value2, dataType2, offset2, err2)
				if err5 != nil {
					Error("GetJsonInt64Arr2:ArrayEach:func2:GetJsonArrayEachInt",
						zap.Int("offset", offset2),
						zap.Error(err2))

					return
				}

				arr0 = append(arr0, cv)
			})
			if err3 != nil {
				Error("GetJsonInt64Arr2:ArrayEach:func:ArrayEach",
					zap.Int("offset", offset3),
					zap.Error(err3))

				return
			}

			arr = append(arr, arr0)

			return
		}

		Error("GetJsonInt64Arr2:ArrayEach:func:dataType",
			zap.Int("offset", offset),
			zap.String("dataType", dataType.String()))
	}, keys...)
	if err != nil {
		if err != jsonparser.KeyPathNotFoundError {
			Error("GetJsonInt64Arr2:ArrayEach",
				zap.Strings("keys", keys),
				zap.Int("offset", offset),
				zap.Error(err))

			return nil, err
		}
	}

	if len(arr) > 0 {
		return arr, nil
	}

	return nil, nil
}

func GetJsonIntArr3(data []byte, keys ...string) ([][][]int, error) {
	arr := [][][]int{}

	offset, err := jsonparser.ArrayEach(data, func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
		if err != nil {
			if err != jsonparser.KeyPathNotFoundError {
				Error("GetJsonIntArr3:ArrayEach:func",
					zap.Int("offset", offset),
					zap.Error(err))

				return
			}
		}

		if dataType == jsonparser.Array {
			arr0 := [][]int{}

			offset3, err3 := jsonparser.ArrayEach(value, func(value2 []byte, dataType2 jsonparser.ValueType, offset2 int, err2 error) {
				if err2 != nil {
					if err != jsonparser.KeyPathNotFoundError {
						Error("GetJsonIntArr3:ArrayEach:func2",
							zap.Int("offset", offset2),
							zap.Error(err2))

						return
					}
				}

				if dataType2 == jsonparser.Array {
					arr1 := []int{}

					offset6, err6 := jsonparser.ArrayEach(value2, func(value5 []byte, dataType5 jsonparser.ValueType, offset5 int, err5 error) {
						if err5 != nil {
							if err != jsonparser.KeyPathNotFoundError {
								Error("GetJsonIntArr3:ArrayEach:func3",
									zap.Int("offset", offset5),
									zap.Error(err5))

								return
							}
						}

						cv, err7 := GetJsonArrayEachInt(value5, dataType5, offset5, err5)
						if err7 != nil {
							Error("GetJsonIntArr3:ArrayEach:func3:GetJsonArrayEachInt",
								zap.Int("offset", offset5),
								zap.Error(err5))

							return
						}

						arr1 = append(arr1, int(cv))
					})
					if err6 != nil {
						Error("GetJsonIntArr3:ArrayEach:func2:ArrayEach",
							zap.Int("offset", offset6),
							zap.Error(err6))

						return
					}

					arr0 = append(arr0, arr1)

					return
				}

				Error("GetJsonIntArr3:ArrayEach:func2:dataType",
					zap.Int("offset", offset2),
					zap.String("dataType", dataType2.String()))
			})
			if err3 != nil {
				Error("GetJsonIntArr3:ArrayEach:func:ArrayEach",
					zap.Int("offset", offset3),
					zap.Error(err3))

				return
			}

			arr = append(arr, arr0)

			return
		}

		Error("GetJsonIntArr3:ArrayEach:func:dataType",
			zap.Int("offset", offset),
			zap.String("dataType", dataType.String()))
	}, keys...)
	if err != nil {
		if err != jsonparser.KeyPathNotFoundError {
			Error("GetJsonIntArr3:ArrayEach",
				zap.Strings("keys", keys),
				zap.Int("offset", offset),
				zap.Error(err))

			return nil, err
		}
	}

	if len(arr) > 0 {
		return arr, nil
	}

	return nil, nil
}

func GetJsonInt64Arr3(data []byte, keys ...string) ([][][]int64, error) {
	arr := [][][]int64{}

	offset, err := jsonparser.ArrayEach(data, func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
		if err != nil {
			if err != jsonparser.KeyPathNotFoundError {
				Error("GetJsonIntArr3:ArrayEach:func",
					zap.Int("offset", offset),
					zap.Error(err))

				return
			}
		}

		if dataType == jsonparser.Array {
			arr0 := [][]int64{}

			offset3, err3 := jsonparser.ArrayEach(value, func(value2 []byte, dataType2 jsonparser.ValueType, offset2 int, err2 error) {
				if err2 != nil {
					if err != jsonparser.KeyPathNotFoundError {
						Error("GetJsonIntArr3:ArrayEach:func2",
							zap.Int("offset", offset2),
							zap.Error(err2))

						return
					}
				}

				if dataType2 == jsonparser.Array {
					arr1 := []int64{}

					offset6, err6 := jsonparser.ArrayEach(value2, func(value5 []byte, dataType5 jsonparser.ValueType, offset5 int, err5 error) {
						if err5 != nil {
							if err != jsonparser.KeyPathNotFoundError {
								Error("GetJsonIntArr3:ArrayEach:func3",
									zap.Int("offset", offset5),
									zap.Error(err5))

								return
							}
						}

						cv, err7 := GetJsonArrayEachInt(value5, dataType5, offset5, err5)
						if err7 != nil {
							Error("GetJsonIntArr3:ArrayEach:func3:GetJsonArrayEachInt",
								zap.Int("offset", offset5),
								zap.Error(err5))

							return
						}

						arr1 = append(arr1, cv)
					})
					if err6 != nil {
						Error("GetJsonIntArr3:ArrayEach:func2:ArrayEach",
							zap.Int("offset", offset6),
							zap.Error(err6))

						return
					}

					arr0 = append(arr0, arr1)

					return
				}

				Error("GetJsonIntArr3:ArrayEach:func2:dataType",
					zap.Int("offset", offset2),
					zap.String("dataType", dataType2.String()))
			})
			if err3 != nil {
				Error("GetJsonIntArr3:ArrayEach:func:ArrayEach",
					zap.Int("offset", offset3),
					zap.Error(err3))

				return
			}

			arr = append(arr, arr0)

			return
		}

		Error("GetJsonIntArr3:ArrayEach:func:dataType",
			zap.Int("offset", offset),
			zap.String("dataType", dataType.String()))
	}, keys...)
	if err != nil {
		if err != jsonparser.KeyPathNotFoundError {
			Error("GetJsonIntArr3:ArrayEach",
				zap.Strings("keys", keys),
				zap.Int("offset", offset),
				zap.Error(err))

			return nil, err
		}
	}

	if len(arr) > 0 {
		return arr, nil
	}

	return nil, nil
}
