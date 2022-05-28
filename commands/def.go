package commands

import (
	"github.com/bwmarrin/discordgo"
)

type Cmd func(*discordgo.Session, *discordgo.Message, []string) error

var (
	Internal map[string]Cmd = make(map[string]Cmd)
	Simple   map[string]Cmd = make(map[string]Cmd)
	Complex  map[string]Cmd = make(map[string]Cmd)
)
