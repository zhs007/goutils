package goutils

func LoadVersion(fn string) (string, error) {
	data, err := ioutil.ReadFile(fn)
	if err != nil {
		return "", err
	}

	return string(data), nil
}
