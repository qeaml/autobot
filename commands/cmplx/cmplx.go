package cmplx

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/qeaml/autobot/commands"
	"github.com/qeaml/autobot/shared"

	"github.com/bwmarrin/discordgo"
	lua "github.com/yuin/gopher-lua"
)

var state *lua.LState
var discord *discordgo.Session
var channel string
var raw map[string]string = map[string]string{}

func Load() error {
	state = lua.NewState(lua.Options{
		SkipOpenLibs:        true,
		IncludeGoStackTrace: true,
	})
	// base Lua API
	state.Push(state.NewFunction(lua.OpenPackage))
	state.Push(lua.LString("package"))
	state.Call(1, 0)
	state.Push(state.NewFunction(lua.OpenBase))
	state.Push(lua.LString("base"))
	state.Call(1, 0)
	// global functions
	state.SetGlobal("say", state.NewFunction(Say))
	state.SetGlobal("typing", state.NewFunction(Typing))
	state.SetGlobal("rand_i", state.NewFunction(RandI))
	state.SetGlobal("rand_f", state.NewFunction(RandF))
	// modules
	state.PreloadModule("strings", Strings)
	state.PreloadModule("http", Http)
	state.PreloadModule("json", Json)

	f, err := os.Open("ccc.json")
	if err != nil {
		return err
	}
	dec := json.NewDecoder(f)
	err = dec.Decode(&raw)
	if err != nil {
		return err
	}
	for cmd, src := range raw {
		cf, err := Compile(cmd, src)
		if err != nil {
			return err
		} else {
			commands.Complex[cmd] = cf
		}
	}
	commands.Internal["ccc"] = ComplexCustom
	return nil
}

func Compile(cmd, src string) (commands.Cmd, error) {
	decl := fmt.Sprintf("function __ccc__%s(msg, args)\n%s\nend", cmd, src)
	err := state.DoString(decl)
	if err != nil {
		return nil, err
	}
	return func(sh *discordgo.Session, msg *discordgo.Message, args []string) error {
		channel = msg.ChannelID
		discord = sh
		return state.CallByParam(lua.P{
			Fn:      state.GetGlobal("__ccc__" + cmd),
			NRet:    0,
			Protect: true,
		}, go2lua(state, messageSimplify(msg)), go2lua(state, args[1:]))
	}, nil
}

type simpleUser struct {
	Name   string
	Id     string
	Avatar string
}

func userSimplify(org *discordgo.User) simpleUser {
	return simpleUser{
		Name:   org.Username,
		Id:     org.ID,
		Avatar: org.AvatarURL(""),
	}
}

type simpleMessage struct {
	Content   string
	Id        string
	ChannelID string
	Author    simpleUser
}

func messageSimplify(msg *discordgo.Message) simpleMessage {
	return simpleMessage{
		Content:   msg.Content,
		Id:        msg.ID,
		ChannelID: msg.ChannelID,
		Author:    userSimplify(msg.Author),
	}
}

func ComplexCustom(sh *discordgo.Session, msg *discordgo.Message, args []string) error {
	if msg.Author.ID == shared.Config.Admin {
		src := strings.TrimSpace(msg.Content[len(args[0])+len(args[1])+2:])
		if strings.HasPrefix(src, "```lua") {
			src = strings.TrimSpace(src[6:])
		}
		if strings.HasPrefix(src, "```") {
			src = strings.TrimSpace(src[3:])
		}
		if strings.HasSuffix(src, "```") {
			src = strings.TrimSpace(src[:len(src)-3])
		}
		cf, err := Compile(args[1], src)
		if err != nil {
			sh.ChannelMessageSend(msg.ChannelID, "***Compilation error:***\n```"+err.Error()+"```")
		} else {
			commands.Complex[args[1]] = cf
			sh.ChannelMessageSend(msg.ChannelID, fmt.Sprintf(
				"Command **%s** updated.", args[1]))
			raw[args[1]] = src
			f, _ := os.Create("ccc.json")
			enc := json.NewEncoder(f)
			enc.Encode(raw)
		}
	} else {
		sh.ChannelMessageSend(msg.ChannelID, "No permissions?")
	}
	return nil
}
