package goutils

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"strings"
	"text/template"

	"go.uber.org/zap"
)

type VersionObj struct {
	Major int
	Minor int
	Patch int
}

func (vobj *VersionObj) ToString() string {
	return fmt.Sprintf("v%v.%v.%v", vobj.Major, vobj.Minor, vobj.Patch)
}

func (vobj *VersionObj) IncPatch() {
	vobj.Patch++
}

func LoadVersion(fn string) (string, error) {
	data, err := ioutil.ReadFile(fn)
	if err != nil {
		return "", err
	}

	return string(data), nil
}

func ParseVersion(str string) (*VersionObj, error) {
	if str[0:1] != "v" && str[0:1] != "V" {
		Warn("ParseVersion",
			zap.String("str", str))

		return nil, ErrInvalidVersion
	}

	arr := strings.Split(str[1:], ".")
	if len(arr) != 3 {
		Warn("ParseVersion:Split",
			zap.String("str", str))

		return nil, ErrInvalidVersion
	}

	v0, err := String2Int64(arr[0])
	if err != nil {
		Warn("ParseVersion:String2Int64:arr0",
			zap.String("str", str))

		return nil, ErrInvalidVersion
	}

	v1, err := String2Int64(arr[1])
	if err != nil {
		Warn("ParseVersion:String2Int64:arr1",
			zap.String("str", str))

		return nil, ErrInvalidVersion
	}

	v2, err := String2Int64(arr[2])
	if err != nil {
		Warn("ParseVersion:String2Int64:arr2",
			zap.String("str", str))

		return nil, ErrInvalidVersion
	}

	return &VersionObj{
		Major: int(v0),
		Minor: int(v1),
		Patch: int(v2),
	}, nil
}

func BuildVersionFile(fn string, tmpfn string, vobj *VersionObj) error {
	data, err := ioutil.ReadFile(tmpfn)
	if err != nil {
		Warn("BuildVersionFile:ReadFile",
			zap.String("tmpfn", tmpfn))

		return err
	}

	buf := new(bytes.Buffer)

	t, err := template.New("buildversion").Parse(string(data))
	if err != nil {
		Warn("BuildVersionFile:template.New",
			zap.String("tmpfn", tmpfn),
			zap.Error(err))

		return err
	}

	err = t.Execute(buf, vobj)
	if err != nil {
		Warn("BuildVersionFile:Execute",
			zap.String("tmpfn", tmpfn),
			zap.String("version", vobj.ToString()),
			zap.Error(err))

		return err
	}

	err = ioutil.WriteFile(fn, buf.Bytes(), 0644)
	if err != nil {
		Warn("BuildVersionFile:WriteFile",
			zap.String("fn", fn),
			zap.Error(err))

		return err
	}

	return nil
}
