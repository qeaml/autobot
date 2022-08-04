package quotes

import (
	"encoding/json"
	"errors"
	"io/fs"
	"math/rand"
	"os"
)

var Quotes = map[string][]string{}

func Load() (err error) {
	f, err := os.Open("quotes.json")
	if errors.Is(err, fs.ErrNotExist) {
		f, err = os.Create("quotes.json")
		if err != nil {
			return
		}
		_, err = f.WriteString("{}\n")
		if err != nil {
			return
		}
		_, err = f.Seek(0, 0)
		if err != nil {
			return
		}
	}
	if err != nil {
		return
	}
	defer f.Close()
	d := json.NewDecoder(f)
	err = d.Decode(&Quotes)
	return
}

func Save() (err error) {
	f, err := os.Create("quotes.json")
	if err != nil {
		return
	}
	defer f.Close()
	e := json.NewEncoder(f)
	err = e.Encode(Quotes)
	return
}

func Add(who, quote string) uint {
	quotelist, ok := Quotes[who]
	if !ok {
		quotelist = []string{}
	}
	quotelist = append(quotelist, quote)
	Quotes[who] = quotelist
	return uint(len(quotelist))
}

func Get(who string, num int) (q string, ok bool) {
	quotelist, ok := Quotes[who]
	if !ok {
		return
	}
	i := num - 1
	if i >= len(quotelist) || i <= 0 {
		ok = false
		return
	}
	q = quotelist[i]
	return
}

func GetRandom(who string) (q string, ok bool) {
	quotelist, ok := Quotes[who]
	if !ok {
		return
	}
	q = quotelist[rand.Intn(len(quotelist))]
	return
}
