package cmplx

import (
	"io"
	"net/http"
	"reflect"

	lua "github.com/yuin/gopher-lua"
)

func go2luaName(goName string) string {
	return string([]byte{goName[0] | 0b100000}) + goName[1:]
}

func go2lua(L *lua.LState, gv any) (lv lua.LValue) {
	if i, ok := gv.(int); ok {
		lv = lua.LNumber(i)
	}
	if i, ok := gv.(int8); ok {
		lv = lua.LNumber(i)
	}
	if i, ok := gv.(int16); ok {
		lv = lua.LNumber(i)
	}
	if i, ok := gv.(int32); ok {
		lv = lua.LNumber(i)
	}
	if i, ok := gv.(int64); ok {
		lv = lua.LNumber(i)
	}

	if i, ok := gv.(uint); ok {
		lv = lua.LNumber(i)
	}
	if i, ok := gv.(uint8); ok {
		lv = lua.LNumber(i)
	}
	if i, ok := gv.(uint16); ok {
		lv = lua.LNumber(i)
	}
	if i, ok := gv.(uint32); ok {
		lv = lua.LNumber(i)
	}
	if i, ok := gv.(uint64); ok {
		lv = lua.LNumber(i)
	}

	if f, ok := gv.(float32); ok {
		lv = lua.LNumber(f)
	}
	if f, ok := gv.(float64); ok {
		lv = lua.LNumber(f)
	}

	if s, ok := gv.(string); ok {
		lv = lua.LString(s)
	}

	rv := reflect.ValueOf(gv)

	if rv.Kind() == reflect.Map {
		tbl := L.NewTable()
		it := rv.MapRange()
		for it.Next() {
			if it.Key().Kind() != reflect.String {
				continue
			}
			tbl.RawSetString(it.Key().String(), go2lua(L, it.Value().Interface()))
		}
		lv = tbl
	}

	if rv.Kind() == reflect.Slice {
		tbl := L.NewTable()
		for i := 0; i < rv.Len(); i++ {
			tbl.Append(go2lua(L, rv.Index(i).Interface()))
		}
		lv = tbl
	}

	if rv.Kind() == reflect.Struct {
		tbl := L.NewTable()
		for i := 0; i < rv.NumField(); i++ {
			k := go2luaName(rv.Type().Field(i).Name)
			v := go2lua(L, rv.Field(i).Interface())
			tbl.RawSetString(k, v)
		}
		lv = tbl
	}

	if lv == nil {
		panic("unsupported Go value of type " + rv.Type().Name())
	}

	return
}

func lua2go(L *lua.LState, lv lua.LValue) (gv any) {
	switch lv.Type() {
	case lua.LTNil:
		gv = nil
	case lua.LTBool:
		gv = lua.LVAsBool(lv)
	case lua.LTNumber:
		gv = lua.LVAsNumber(lv)
	case lua.LTString:
		gv = lua.LVAsString(lv)
	case lua.LTTable:
		tbl := lv.(*lua.LTable)
		if tbl.MaxN() == 0 {
			gv = table2map(L, tbl)
		} else {
			gv = table2arr(L, tbl)
		}
	default:
		panic("unsupporrted Lua value of type " + lv.Type().String())
	}
	return
}

func table2map(L *lua.LState, tbl *lua.LTable) map[string]any {
	vals := map[string]any{}
	tbl.ForEach(func(k, v lua.LValue) {
		if k.Type() != lua.LTString {
			return
		}
		vals[lua.LVAsString(k)] = lua2go(L, v)
	})
	return vals
}

func table2arr(L *lua.LState, tbl *lua.LTable) []any {
	vals := []any{}
	tbl.ForEach(func(k, v lua.LValue) {
		if k.Type() != lua.LTNumber {
			return
		}
		vals = append(vals, lua2go(L, v))
	})
	return vals
}

func map2table(L *lua.LState, vals map[string]any) *lua.LTable {
	tbl := L.NewTable()
	for k, v := range vals {
		tbl.RawSetString(k, go2lua(L, v))
	}
	return tbl
}

func arr2table(L *lua.LState, a []any) *lua.LTable {
	tbl := L.NewTable()
	for i, v := range a {
		tbl.RawSetInt(i+1, go2lua(L, v))
	}
	return tbl
}

func table2response(L *lua.LState, resp *http.Response) {
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return
	}
	resp.Body.Close()
	respt := L.NewTable()
	respt.RawSetString("status", lua.LNumber(resp.StatusCode))
	respt.RawSetString("text", lua.LString(body))
	L.Push(respt)
}
