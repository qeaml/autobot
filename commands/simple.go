package commands

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/qeaml/autobot/shared"

	"github.com/bwmarrin/discordgo"
)

func LoadSimple() error {
	f, err := os.Open("cc.json")
	if err != nil {
		return err
	}
	defer f.Close()
	dec := json.NewDecoder(f)
	simpleSrc = make(map[string]string)
	err = dec.Decode(&simpleSrc)
	if err != nil {
		return err
	}
	for cmd, txt := range simpleSrc {
		Simple[cmd] = str2cmd(cmd, txt)
	}
	return nil
}

func SaveSimple() error {
	f, err := os.Create("cc.json")
	if err != nil {
		return err
	}
	defer f.Close()
	enc := json.NewEncoder(f)
	return enc.Encode(simpleSrc)
}

var simpleSrc map[string]string

func str2cmd(cmd, txt string) Cmd {
	return func(sh *discordgo.Session, msg *discordgo.Message, args []string) error {
		slots := strings.Count(txt, "%s")
		fmtArgs := []any{}
		if len(args) == 1 {
			for i := 0; i < slots; i++ {
				fmtArgs = append(fmtArgs, "")
			}
		} else {
			for _, a := range args[1:slots] {
				fmtArgs = append(fmtArgs, a)
			}
			fmtArgs = append(fmtArgs, strings.Join(args[slots:], " "))
		}
		old := shared.UsageSimple.GetOrDefault(cmd, 0)
		shared.UsageSimple.Set(cmd, old+1)
		sh.ChannelMessageSend(msg.ChannelID, fmt.Sprintf(txt, fmtArgs...))
		return nil
	}
}

func Custom(sh *discordgo.Session, msg *discordgo.Message, args []string) error {
	if msg.Author.ID == shared.Config.Admin {
		txt := strings.Join(args[2:], " ")
		simpleSrc[args[1]] = txt
		Simple[args[1]] = str2cmd(args[1], txt)
		sh.ChannelMessageSend(msg.ChannelID, fmt.Sprintf(
			"Command **%s** updated.", args[1]))
	} else {
		sh.ChannelMessageSend(msg.ChannelID, "No permissions?")
	}
	return nil
}
