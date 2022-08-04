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
	"github.com/qeaml/autobot/quotes"
	"github.com/qeaml/autobot/shared"

	"github.com/bwmarrin/discordgo"
	"gopkg.in/yaml.v3"
)

func main() {
	log("Loading internal commands")
	if err := commands.LoadInternal(); err != nil {
		log("error loading internal commands: %s", err.Error())
		return
	}
	log("Loading simple commands")
	if err := commands.LoadSimple(); err != nil {
		log("error loading simple commands: %s", err.Error())
		return
	}
	log("Loading complex commands")
	if err := cmplx.Load(); err != nil {
		log("error loading complex commands: %s", err.Error())
		return
	}

	log("Loading usages")
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

	log("Loading model")
	err = model.Load()
	if err != nil {
		log("error loading model: " + err.Error())
		return
	}

	log("Loading quotes")
	err = quotes.Load()
	if err != nil {
		log("error loading quotes: " + err.Error())
		return
	}

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

	log("Connecting to Discord...")
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

	log("Bot running. Press Ctrl+C to stop.")
	<-sc
	stop := time.Now()

	log("Disconnecting discord...")
	if err = discord.Close(); err != nil {
		fmt.Printf("\nerror: %s\n", err.Error())
	}
	log("Saving quotes...")
	if err = quotes.Save(); err != nil {
		fmt.Printf("\nerror: %s\n", err.Error())
	}
	log("Saving model...")
	if err = model.Memory.Save(); err != nil {
		fmt.Printf("\nerror: %s\n", err.Error())
	} else {
		model.Source.Close()
	}
	log("Saving builtin usage...")
	if err = shared.UsageBuiltin.Save(); err != nil {
		fmt.Printf("\nerror: %s\n", err.Error())
	}
	log("Saving SCC...")
	if err = commands.SaveSimple(); err != nil {
		fmt.Printf("\nerror : %s\n", err.Error())
	}
	log("Saving SCC usage...")
	if err = shared.UsageSimple.Save(); err != nil {
		fmt.Printf("\nerror : %s\n", err.Error())
	}
	log("Saving CCC usage...")
	if err = shared.UsageComplex.Save(); err != nil {
		fmt.Printf("\nerror: %s\n", err.Error())
	}
	log("Shut down complete. Bot ran for %s.\n", stop.Sub(shared.Start).Round(time.Second))
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

	sane := shared.Whitespacer.Replace(msg.Content)
	args := strings.Split(sane, " ")
	cmd := args[0][1:]
	if args[0][0] == shared.Config.Prefix[0] {
		c, ok := commands.Internal[cmd]
		if ok {
			c(discord, msg, args)
		}
	} else if args[0][0] == shared.Config.CCPrefix[0] {
		c, ok := commands.Simple[cmd]
		if ok {
			c(discord, msg, args)
			return
		}
		c, ok = commands.Complex[cmd]
		if ok {
			c(discord, msg, args)
		}
	} else if strings.HasPrefix(msg.Content, discord.State.User.Mention()) {
		txt := msg.Author.Mention() + model.Generate("", 5+rand.Intn(10))
		discord.ChannelMessageSend(msg.ChannelID, txt)
	} else {
		model.TrainString(msg)
		if rand.Intn(11) == 5 {
			discord.ChannelMessageSendReply(
				msg.ChannelID,
				model.Generate("", 5+rand.Intn(10)),
				msg.Reference())
		}
	}
}

func log(format string, args ...any) {
	fmt.Printf("\x1b[2K\x1b[0G"+format, args...)
}
