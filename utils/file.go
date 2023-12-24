package utils

import (
	errors "github.com/wuntsong-org/wterrors"
	"os"
	"strings"
)

func WriteString(path string, content string, append bool) errors.WTError {
	flag := os.O_RDWR | os.O_CREATE
	if append {
		flag = flag | os.O_APPEND
	}
	file, err := os.OpenFile(path, flag, 0644)
	if err != nil {
		return errors.WarpQuick(err)
	}
	defer Close(file)

	_, err = file.WriteString(content)
	return errors.WarpQuick(err)
}

func AppendLine(path string, content string) errors.WTError {
	file, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		return errors.WarpQuick(err)
	}
	defer Close(file)

	content = strings.Join([]string{content, "\n"}, "")
	_, err = file.WriteString(content)
	return errors.WarpQuick(err)
}

func PathExists(path string) (bool, errors.WTError) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, errors.WarpQuick(err)
}
