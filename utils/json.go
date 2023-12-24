package utils

import (
	"bytes"
	jsoniter "github.com/json-iterator/go"
	errors "github.com/wuntsong-org/wterrors"
	"io"
)

func FormatJson(obj interface{}) (str string, err errors.WTError) {
	data, err := JsonMarshal(obj)
	if err != nil {
		return
	}
	str = string(data)
	return
}

func ParseJson(str string, t interface{}) errors.WTError {
	return JsonUnmarshal([]byte(str), t)
}

func JsonMarshal(v interface{}) ([]byte, errors.WTError) {
	buf := bytes.NewBuffer(make([]byte, 0))
	e := jsoniter.NewEncoder(buf)
	err := e.Encode(v)
	if err != nil {
		return nil, errors.WarpQuick(err)
	}

	res, err := io.ReadAll(buf)
	if err != nil {
		return nil, errors.WarpQuick(err)
	}

	return res, nil
}

func JsonUnmarshal(data []byte, v interface{}) errors.WTError {
	d := jsoniter.NewDecoder(bytes.NewBuffer(data))
	d.UseNumber()
	err := d.Decode(v)
	if err != nil {
		return errors.WarpQuick(err)
	}

	return nil
}

func IsJsonString(data string) bool {
	tmp := make(map[string]interface{}, 10)
	return JsonUnmarshal([]byte(data), &tmp) == nil
}
