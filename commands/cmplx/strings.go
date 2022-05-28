package cmplx

import (
	"crypto/sha256"
	"strings"

	lua "github.com/yuin/gopher-lua"
)

var Strings = lib(map[string]lua.LGFunction{
	"join": func(L *lua.LState) int {
		elemt := L.ToTable(1)
		sep := L.ToString(2)
		out := ""
		elemt.ForEach(func(k, v lua.LValue) {
			if k.Type() == lua.LTNumber {
				out += lua.LVAsString(v) + sep
			}
		})
		L.Push(lua.LString(out[:len(out)-len(sep)]))
		return 1
	},
	"startswith": func(L *lua.LState) int {
		a := L.ToString(1)
		b := L.ToString(2)
		L.Push(lua.LBool(strings.HasPrefix(a, b)))
		return 1
	},
	"endswith": func(L *lua.LState) int {
		a := L.ToString(1)
		b := L.ToString(2)
		L.Push(lua.LBool(strings.HasSuffix(a, b)))
		return 1
	},
	"trimspace": func(L *lua.LState) int {
		str := L.ToString(1)
		L.Push(lua.LString(strings.TrimSpace(str)))
		return 1
	},
	"hash": func(L *lua.LState) int {
		input := L.ToString(1)
		hashRaw := sha256.Sum256([]byte(input))
		var hash uint32
		for _, v := range hashRaw {
			hash += uint32(v)
		}
		L.Push(lua.LNumber(hash))
		return 1
	},
	"lower": func(L *lua.LState) int {
		L.Push(lua.LString(strings.ToLower(L.ToString(1))))
		return 1
	},
	"upper": func(L *lua.LState) int {
		L.Push(lua.LString(strings.ToUpper(L.ToString(1))))
		return 1
	},
	"cap": func(L *lua.LState) int {
		s := L.ToString(1)
		max := int(L.ToNumber(2))
		if len(s) > max {
			L.Push(lua.LString(s[:max]))
		} else {
			L.Push(lua.LString(s))
		}
		return 1
	},
})
