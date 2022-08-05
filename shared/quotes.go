package shared

import (
	"encoding/json"
	"errors"
	"io/fs"
	"os"
)

var Quotes = map[string][]string{}

func LoadQuotes() (err error) {
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

func SaveQuotes() (err error) {
	f, err := os.Create("quotes.json")
	if err != nil {
		return
	}
	defer f.Close()
	e := json.NewEncoder(f)
	err = e.Encode(Quotes)
	return
}

func AddQuote(who, quote string) uint {
	quotelist, ok := Quotes[who]
	if !ok {
		quotelist = []string{}
	}
	quotelist = append(quotelist, quote)
	Quotes[who] = quotelist
	return uint(len(quotelist))
}
