package cmplx

import (
	"math/rand"

	lua "github.com/yuin/gopher-lua"
)

func lib(mod map[string]lua.LGFunction) func(L *lua.LState) int {
	return func(L *lua.LState) int {
		L.Push(L.SetFuncs(L.NewTable(), mod))
		return 1
	}
}

func Say(L *lua.LState) int {
	discord.ChannelMessageSend(channel, L.ToString(1))
	return 0
}

func Typing(L *lua.LState) int {
	discord.ChannelTyping(channel)
	return 0
}

func RandI(L *lua.LState) int {
	max := L.ToInt(1)
	L.Push(lua.LNumber(rand.Intn(max)))
	return 1
}

func RandF(L *lua.LState) int {
	L.Push(lua.LNumber(rand.Float64()))
	return 1
}
