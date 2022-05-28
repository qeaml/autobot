package cmplx

import (
	"net/http"
	"strings"

	"github.com/qeaml/autobot/shared"

	lua "github.com/yuin/gopher-lua"
)

var Http = lib(map[string]lua.LGFunction{
	"get": func(L *lua.LState) int {
		url := L.ToString(1)
		resp, err := http.Get(url)
		if err != nil {
			shared.Error = err
			return 0
		}
		table2response(L, resp)
		return 1
	},
	"post": func(L *lua.LState) int {
		url := L.ToString(1)
		ctype := L.ToString(2)
		body := L.ToString(3)
		resp, err := http.Post(url, ctype, strings.NewReader(body))
		if err != nil {
			shared.Error = err
			return 0
		}
		table2response(L, resp)
		return 1
	},
})
