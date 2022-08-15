package main

import (
	"fmt"
	"math/rand"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/qeaml/autobot/commands"
	"github.com/qeaml/autobot/commands/cmplx"
	"github.com/qeaml/autobot/model"
	"github.com/qeaml/autobot/shared"

	"github.com/bwmarrin/discordgo"
	"gopkg.in/yaml.v3"
)

func main() {
	shared.PushLog("Bot started")

	shared.PushLog("Loading state")

	shared.PushLog("Loading internal commands")
	if err := commands.LoadInternal(); err != nil {
		log("error loading internal commands: %s", err.Error())
		return
	}

	shared.SwapLog("Loading simple commands")
	if err := commands.LoadSimple(); err != nil {
		log("error loading simple commands: %s", err.Error())
		return
	}

	shared.SwapLog("Loading complex commands")
	if err := cmplx.Load(); err != nil {
		log("error loading complex commands: %s", err.Error())
		return
	}

	shared.SwapLog("Loading usages")
	ubCh, ubf, err := shared.Cache[string, uint64]("usages.builtin.json")
	if err != nil {
		log("error opening builtin usages: %s", err.Error())
	}
	defer ubf.Close()
	shared.UsageBuiltin = ubCh
	usCh, usf, err := shared.Cache[string, uint64]("usages.simple.json")
	if err != nil {
		log("error opening simple usages: %s", err.Error())
	}
	defer usf.Close()
	shared.UsageSimple = usCh
	ucCh, ucf, err := shared.Cache[string, uint64]("usages.complex.json")
	if err != nil {
		log("error opening complex usages: %s", err.Error())
	}
	defer ucf.Close()
	shared.UsageComplex = ucCh

	shared.SwapLog("Loading model")
	err = model.Load()
	if err != nil {
		log("error loading model: " + err.Error())
		return
	}

	shared.SwapLog("Loading quotes")
	err = shared.LoadQuotes()
	if err != nil {
		log("error loading quotes: " + err.Error())
		return
	}

	shared.SwapLog("Loading WTF")
	err = shared.LoadWTF()
	if err != nil {
		log("error loading WTF: " + err.Error())
		return
	}

	shared.PopLog()

	log("Reading config...")
	cf, err := os.OpenFile("config.yml", os.O_RDONLY|os.O_CREATE, os.ModePerm)
	if err != nil {
		println("error opening config: " + err.Error())
		return
	}

	dec := yaml.NewDecoder(cf)
	err = dec.Decode(&shared.Config)
	if err != nil {
		println("error parsing cofnig: %s", err.Error())
		return
	}
	cf.Close()

	shared.SwapLog("Connecting to Discord...")
	discord, err := discordgo.New("Bot " + shared.Config.Token)
	if err != nil {
		log("error creating Discord session: %s", err.Error())
		return
	}
	shared.Start = time.Now()
	discord.Identify.Intents = discordgo.IntentsAllWithoutPrivileged
	if err = discord.Open(); err != nil {
		log("error connecting to Discord: %s", err.Error())
		return
	}
	discord.AddHandler(onMessage)
	discord.AddHandler(onEditMessage)

	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)

	shared.SwapLog("Bot running. Press Ctrl+C to stop.")
	<-sc
	stop := time.Now()

	shared.SwapLog("Disconnecting discord...")
	if err = discord.Close(); err != nil {
		fmt.Printf("\nerror: %s\n", err.Error())
	}
	shared.SwapLog("Saving WTF...")
	if err = shared.SaveWTF(); err != nil {
		fmt.Printf("\nerror: %s\n", err.Error())
	}
	shared.SwapLog("Saving quotes...")
	if err = shared.SaveQuotes(); err != nil {
		fmt.Printf("\nerror: %s\n", err.Error())
	}
	shared.SwapLog("Saving model...")
	if err = model.Memory.Save(); err != nil {
		fmt.Printf("\nerror: %s\n", err.Error())
	} else {
		model.Source.Close()
	}
	shared.SwapLog("Saving builtin usage...")
	if err = shared.UsageBuiltin.Save(); err != nil {
		fmt.Printf("\nerror: %s\n", err.Error())
	}
	shared.SwapLog("Saving SCC...")
	if err = commands.SaveSimple(); err != nil {
		fmt.Printf("\nerror : %s\n", err.Error())
	}
	shared.SwapLog("Saving SCC usage...")
	if err = shared.UsageSimple.Save(); err != nil {
		fmt.Printf("\nerror : %s\n", err.Error())
	}
	shared.SwapLog("Saving CCC usage...")
	if err = shared.UsageComplex.Save(); err != nil {
		fmt.Printf("\nerror: %s\n", err.Error())
	}
	shared.PopLog()
	shared.SwapLog("Shut down complete. Bot ran for %s.\n", stop.Sub(shared.Start).Round(time.Second))
}

func onMessage(discord *discordgo.Session, evt *discordgo.MessageCreate) {
	handleMessage(discord, evt.Message)
}

func onEditMessage(discord *discordgo.Session, evt *discordgo.MessageUpdate) {
	handleMessage(discord, evt.Message)
}

func handleMessage(discord *discordgo.Session, msg *discordgo.Message) {
	if msg == nil {
		return
	}

	if msg.Author == nil {
		return
	}

	if msg.Author.Bot || msg.Content == "" {
		return
	}

	shared.PushLog("Handling message...")

	sane := shared.Whitespacer.Replace(msg.Content)
	args := strings.Split(sane, " ")
	cmd := args[0][1:]
	if args[0][0] == shared.Config.Prefix[0] {
		c, ok := commands.Internal[cmd]
		if ok {
			shared.PushLog("Internal command: %s", cmd)
			c(discord, msg, args)
			shared.PopLog()
		}
	} else if args[0][0] == shared.Config.CCPrefix[0] {
		c, ok := commands.Simple[cmd]
		if ok {
			shared.PushLog("Simple command: %s", cmd)
			c(discord, msg, args)
			shared.PopLog()
			return
		}
		c, ok = commands.Complex[cmd]
		if ok {
			shared.PushLog("Complex command: %s", cmd)
			c(discord, msg, args)
			shared.PopLog()
		}
	} else if strings.HasPrefix(msg.Content, discord.State.User.Mention()) {
		shared.PushLog("Mention")
		txt := msg.Author.Mention() + model.Generate("", 5+rand.Intn(10))
		discord.ChannelMessageSend(msg.ChannelID, txt)
		shared.PopLog()
	} else {
		shared.PushLog("Training text generator")
		model.TrainString(msg)
		shared.PopLog()
		if rand.Intn(11) == 5 {
			shared.PushLog("Replying")
			discord.ChannelMessageSendReply(
				msg.ChannelID,
				model.Generate("", 5+rand.Intn(10)),
				msg.Reference())
			shared.PopLog()
		}
		lower := strings.ToLower(msg.Content)
		if strings.Contains(lower, "wtf") || strings.Contains(lower, "what the fuck") {
			shared.WTF++
			discord.ChannelMessageSend(msg.ChannelID,
				fmt.Sprintf("WTF moments: %d", shared.WTF))
		}
	}

	shared.PopLog()
}

func log(format string, args ...any) {
	fmt.Printf("\x1b[2K\x1b[0G"+format, args...)
}
