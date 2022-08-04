package commands

import (
	"strings"

	"github.com/qeaml/autobot/commands/internal"
	"github.com/qeaml/autobot/shared"

	"github.com/bwmarrin/discordgo"
)

func LoadInternal() error {
	Internal = map[string]Cmd{
		"generate": internal.Generate,
		"model":    internal.Model,
		"cc":       Custom,
		"ccc":      nil,
		"error":    internal.Error,
		"commands": Commands,
		"runc":     internal.RunC,
		"runpy":    internal.RunPy,
		"doc":      internal.Doc,
		"ping":     internal.Ping,
		"usage":    internal.Usage,
		"uptime":   internal.Uptime,
		"quote":    internal.Quote,
		"addquote": internal.AddQuote,
	}
	return nil
}

var internalDesc = map[string]string{
	"generate": "Generate text with an optional starter word",
	"model":    "Check the size of the text generation model",
	"cc":       "Modify simple custom commands",
	"ccc":      "Modify complex custom commands",
	"error":    "Check what error has last occured, and when",
	"commands": "View a list of all commands",
	"runc":     "Compile and run a C program",
	"runpy":    "Run a Python script",
	"doc":      "View the complex custom command documentation",
	"ping":     "Check the bot's ping to Discord",
	"uptime":   "Check the bot's uptime",
	"quote":    "View a quote",
	"addquote": "Add a quote (admin only)",
}

func Commands(sh *discordgo.Session, msg *discordgo.Message, args []string) error {
	txt := strings.Builder{}
	txt.WriteString("**Built-in commands:**\n")
	for cmd := range Internal {
		txt.WriteString("> `")
		txt.WriteString(shared.Config.Prefix)
		txt.WriteString(cmd)
		txt.WriteString("` - ")
		txt.WriteString(internalDesc[cmd])
		txt.WriteString("\n")
	}
	txt.WriteString("**Simple custom commands:**\n> ")
	for cmd := range Simple {
		txt.WriteString("`")
		txt.WriteString(shared.Config.CCPrefix)
		txt.WriteString(cmd)
		txt.WriteString("`, ")
	}
	txt.WriteString("\n**Complex custom commands:**\n> ")
	for cmd := range Complex {
		txt.WriteString("`")
		txt.WriteString(shared.Config.CCPrefix)
		txt.WriteString(cmd)
		txt.WriteString("`, ")
	}
	sh.ChannelMessageSend(msg.ChannelID, txt.String())
	return nil
}
