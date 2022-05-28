package cmplx

import (
	"encoding/json"

	"github.com/qeaml/autobot/shared"

	"github.com/wI2L/jettison"
	lua "github.com/yuin/gopher-lua"
)

var Json = lib(map[string]lua.LGFunction{
	"encode": func(L *lua.LState) int {
		vals := table2map(L, L.ToTable(1))
		json, err := jettison.Marshal(vals)
		if err != nil {
			shared.Error = err
			return 0
		}
		L.Push(lua.LString(json))
		return 1
	},
	"decode": func(L *lua.LState) int {
		vals := map[string]any{}
		err := json.Unmarshal([]byte(L.ToString(1)), &vals)
		if err != nil {
			shared.Error = err
			return 0
		}
		L.Push(go2lua(L, vals))
		return 1
	},
})
