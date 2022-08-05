package shared

import (
	"errors"
	"fmt"
	"io"
	"io/fs"
	"os"
	"strconv"
)

var WTF uint

func LoadWTF() (err error) {
	f, err := os.Open("wtf")
	if errors.Is(err, fs.ErrNotExist) {
		f, err = os.Create("wtf")
		if err != nil {
			return
		}
		_, err = f.WriteString("0")
		if err != nil {
			return
		}
		_, err = f.Seek(0, 0)
		if err != nil {
			return
		}
	}
	defer f.Close()
	data, err := io.ReadAll(f)
	if err != nil {
		return
	}
	num, err := strconv.ParseUint(string(data), 10, 32)
	if err != nil {
		return
	}
	WTF = uint(num)
	return
}

func SaveWTF() (err error) {
	f, err := os.Create("wtf")
	if err != nil {
		return
	}
	defer f.Close()
	_, err = fmt.Fprintf(f, "%d", WTF)
	return
}
