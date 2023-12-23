package utils

import (
	"bytes"
	"encoding/json"
	jsoniter "github.com/json-iterator/go"
	"io"
)

func FormatJson(obj interface{}) (str string, err error) {
	data, err := json.Marshal(obj)
	if err != nil {
		return
	}
	str = string(data)
	return
}

func ParseJson(str string, t interface{}) error {
	return json.Unmarshal([]byte(str), t)
}

func JsonMarshal(v interface{}) ([]byte, error) {
	buf := bytes.NewBuffer(make([]byte, 0))
	e := jsoniter.NewEncoder(buf)
	err := e.Encode(v)
	if err != nil {
		return nil, err
	}

	res, err := io.ReadAll(buf)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func JsonUnmarshal(data []byte, v interface{}) error {
	d := jsoniter.NewDecoder(bytes.NewBuffer(data))
	d.UseNumber()
	err := d.Decode(v)
	if err != nil {
		return err
	}

	return nil
}

func IsJsonString(data string) bool {
	tmp := make(map[string]interface{}, 10)
	return JsonUnmarshal([]byte(data), &tmp) == nil
}
