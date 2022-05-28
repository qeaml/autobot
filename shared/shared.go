package shared

import (
	"os"
	"strings"
	"time"

	"github.com/qeaml/ggpc"
	lua "github.com/yuin/gopher-lua"
)

var Error error = nil
var Value lua.LValue = lua.LNil
var Whitespacer = strings.NewReplacer("\t", " ", "\n", " ")
var Start time.Time

var (
	UsageBuiltin *ggpc.Cache[string, uint64]
	UsageSimple  *ggpc.Cache[string, uint64]
	UsageComplex *ggpc.Cache[string, uint64]
)

func Cache[K comparable, V any](fn string) (*ggpc.Cache[K, V], *os.File, error) {
	var f *os.File
	_, err := os.Stat(fn)
	if os.IsNotExist(err) {
		f, err = os.Create(fn)
		if err != nil {
			return nil, nil, err
		} else {
			f.WriteString("{}\n")
			f.Seek(0, 0)
		}
	} else if err != nil {
		return nil, nil, err
	} else {
		f, err = os.OpenFile(fn, os.O_RDWR, 0)
		if err != nil {
			return nil, nil, err
		}
	}
	cache, err := ggpc.LoadStored[K, V](f)
	return cache, f, err
}

type cfg struct {
	Token      string
	Admin      string
	WeatherKey string `yaml:"owm-key"`
	Prefix     string `yaml:"nprefix"`
	CCPrefix   string `yaml:"cprefix"`
}

var Config cfg
