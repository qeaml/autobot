package internal

import (
	"fmt"
	"io"
	"math/rand"
	"os"
	"os/exec"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/qeaml/autobot/model"
	"github.com/qeaml/autobot/shared"

	"github.com/bwmarrin/discordgo"
)

func Generate(sh *discordgo.Session, msg *discordgo.Message, args []string) error {
	var starter string
	if len(args) >= 2 {
		starter = args[len(args)-1]
	} else {
		starter = ""
	}
	passage := model.Generate(starter, 5+rand.Intn(10))
	var prefix string
	if len(args) > 2 {
		prefix = strings.Join(args[1:len(args)-1], " ") + " "
		model.Train(args[1:])
	} else {
		prefix = ""
	}
	_, err := sh.ChannelMessageSend(msg.ChannelID, prefix+passage)
	return err
}

func Model(sh *discordgo.Session, msg *discordgo.Message, args []string) error {
	smdl := model.Memory.Solidify()
	links := len(smdl)
	for _, wmap := range smdl {
		links += len(wmap)
	}
	_, err := sh.ChannelMessageSend(msg.ChannelID, fmt.Sprintf(
		"I know **%d** starter words. I know of **%d** word pairs.",
		len(smdl), links))
	return err
}

func Error(sh *discordgo.Session, msg *discordgo.Message, args []string) error {
	if msg.Author.ID == shared.Config.Admin {
		sh.ChannelMessageSend(msg.ChannelID, shared.Error.Error())
	}
	return nil
}

func RunC(sh *discordgo.Session, msg *discordgo.Message, args []string) error {
	if len(args) < 3 {
		sh.ChannelMessageSend(msg.ChannelID,
			"Usage: !runc <compiler> <code...>")
		return nil
	}

	sh.ChannelTyping(msg.ChannelID)

	compiler := args[1]

	altArgs := strings.SplitN(msg.Content, " ", 3)
	code := altArgs[2]

	code = strings.TrimPrefix(code, "```c")
	code = strings.TrimPrefix(code, "```")
	code = strings.TrimSuffix(code, "```")

	shared.PushLog("Saving source to file...")
	f, err := os.Create("run.c")
	if err != nil {
		return err
	}
	_, err = f.WriteString(code)
	if err != nil {
		return err
	}
	err = f.Close()
	if err != nil {
		return err
	}

	shared.SwapLog("Compiling...")
	var (
		out     string
		badCode bool
	)
	switch compiler {
	case "gcc":
		out, badCode, err = compileGCC()
	case "clang":
		out, badCode, err = compileClang()
	case "cl":
		out, badCode, err = compileCl(false)
	case "cl-checked":
		out, badCode, err = compileCl(true)
	default:
		sh.ChannelMessageSend(msg.ChannelID, fmt.Sprintf(
			"Uknown compiler: %s", compiler))
		shared.PopLog()
		return nil
	}
	if badCode {
		sh.ChannelMessageSend(msg.ChannelID, fmt.Sprintf(
			"Could not compile:\n```c\n%s\n```",
			out))
		shared.PopLog()
		return nil
	}
	if err != nil {
		return err
	}

	shared.SwapLog("Running...")
	runOut := strings.Builder{}
	run := exec.Command(".\\run.exe")
	run.Stdout = &runOut
	run.Stderr = &runOut
	err = run.Run()
	if err != nil {
		if _, ok := err.(*exec.ExitError); ok {
			sh.ChannelMessageSend(msg.ChannelID, fmt.Sprintf(
				"An error occured trying to run your program:\n```c\n%s\n```",
				runOut.String()))
			shared.PopLog()
			return nil
		}
		return err
	}
	shared.SwapLog("Sending results...")
	result := runOut.String()
	if len(result) == 0 {
		sh.ChannelMessageSend(msg.ChannelID, "**Result empty.**")
	} else {
		sh.ChannelMessageSend(msg.ChannelID, runOut.String())
	}
	shared.PopLog()
	return nil
}

func RunPy(sh *discordgo.Session, msg *discordgo.Message, args []string) error {
	if len(args) < 2 {
		sh.ChannelMessageSend(msg.ChannelID, "Provide code to run")
		return nil
	}
	sh.ChannelTyping(msg.ChannelID)
	src := strings.TrimSpace(msg.Content[len(args[0])+1:])
	if strings.HasPrefix(src, "```python") {
		src = strings.TrimSpace(src[9:])
	}
	if strings.HasPrefix(src, "```py") {
		src = strings.TrimSpace(src[5:])
	}
	if strings.HasPrefix(src, "```") {
		src = strings.TrimSpace(src[3:])
	}
	if strings.HasSuffix(src, "```") {
		src = strings.TrimSpace(src[:len(src)-3])
	}
	shared.PushLog("Writing source code to file...")
	pyFile, err := os.Create("run.py")
	if err != nil {
		return err
	}
	_, err = pyFile.WriteString(src)
	if err != nil {
		return err
	}
	err = pyFile.Close()
	if err != nil {
		return err
	}
	shared.SwapLog("Running...")
	runOut := strings.Builder{}
	run := exec.Command("py", "run.py")
	run.Stdout = &runOut
	run.Stderr = &runOut
	err = run.Run()
	if err != nil {
		if ex, ok := err.(*exec.ExitError); ok {
			txt := fmt.Sprintf(
				"**Execution failed. (code %x)**\n```py\n%s\n```", ex.ExitCode(), runOut.String())
			sh.ChannelMessageSend(msg.ChannelID, txt)
			shared.PopLog()
			return nil
		}
		return err
	}
	shared.SwapLog("Sending results...")
	result := runOut.String()
	if len(result) == 0 {
		sh.ChannelMessageSend(msg.ChannelID, "**Result empty.**")
	} else {
		sh.ChannelMessageSend(msg.ChannelID, runOut.String())
	}
	shared.PopLog()
	return nil
}

func Doc(sh *discordgo.Session, msg *discordgo.Message, args []string) error {
	query := strings.Join(args[1:], " ")
	subpath := strings.ReplaceAll(query, ".", "/")
	path := "docs"
	if len(subpath) > 0 {
		path += "/" + subpath
	}

	dfi, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			path += ".md"
			dfi, err = os.Stat(path)
			if err != nil {
				return err
			}
		} else {
			return err
		}
	}
	if dfi.IsDir() {
		ovw := "No overview"
		entries := []string{}
		subpages := []string{}
		files, err := os.ReadDir(path)
		if err != nil {
			return err
		}
		for _, d := range files {
			if err != nil {
				return err
			}
			if d.Name() == "overview" {
				f, err := os.Open(path + "/overview")
				if err != nil {
					return err
				}
				defer f.Close()
				ovwRaw, err := io.ReadAll(f)
				if err != nil {
					return err
				}
				ovw = strings.TrimSpace(string(ovwRaw))
			} else if d.IsDir() {
				if d.Name() != dfi.Name() {
					subpages = append(subpages, d.Name())
				}
			} else {
				entries = append(entries, strings.TrimSuffix(d.Name(), ".md"))
			}
		}
		if err != nil {
			return err
		}
		txt := "**" + dfi.Name() + "**\n*" + ovw + "*\n"
		if len(entries) != 0 {
			txt += "**Entries:**\n"
			for _, e := range entries {
				txt += "> `" + e + "`\n"
			}
		}
		if len(subpages) != 0 {
			txt += "**Subpages**:\n"
			for _, p := range subpages {
				txt += "> `" + p + "`\n"
			}
		}
		sh.ChannelMessageSend(msg.ChannelID, txt)
	} else {
		f, err := os.Open(path)
		if err != nil {
			return err
		}
		defer f.Close()
		txt, err := io.ReadAll(f)
		if err != nil {
			return err
		}
		sh.ChannelMessageSend(msg.ChannelID, string(txt))
	}
	return nil
}

func Ping(sh *discordgo.Session, msg *discordgo.Message, args []string) error {
	a := time.Now()
	pmsg, err := sh.ChannelMessageSend(msg.ChannelID, "Ping!")
	if err != nil {
		return err
	}
	_, err = sh.ChannelMessageEdit(
		pmsg.ChannelID, pmsg.ID,
		"Pong! "+time.Since(a).Round(time.Millisecond).String())
	return err
}

func Usage(sh *discordgo.Session, msg *discordgo.Message, args []string) error {
	bcu := sortUsageMap(shared.UsageBuiltin.Solidify(), 5)
	scu := sortUsageMap(shared.UsageSimple.Solidify(), 5)
	ccu := sortUsageMap(shared.UsageComplex.Solidify(), 5)

	txt := "**Most used builtin commands:**\n"
	for cmd, amt := range bcu {
		txt += "> `" + shared.Config.Prefix + cmd + "` - " + fmt.Sprint(amt) + "\n"
	}
	txt += "**Most used simple custom commands:**\n"
	for cmd, amt := range scu {
		txt += "> `" + shared.Config.CCPrefix + cmd + "` - " + fmt.Sprint(amt) + "\n"
	}
	txt += "**Most used complex custom commands:**\n"
	for cmd, amt := range ccu {
		txt += "> `" + shared.Config.CCPrefix + cmd + "` - " + fmt.Sprint(amt) + "\n"
	}

	_, err := sh.ChannelMessageSend(msg.ChannelID, txt)
	return err
}

func sortUsageMap(m map[string]uint64, l int) map[string]uint64 {
	keys := []string{}
	for k := range m {
		keys = append(keys, k)
	}
	sort.Slice(keys, func(i, j int) bool {
		return m[keys[i]] > m[keys[j]]
	})
	newMap := map[string]uint64{}
	i := 0
	for _, k := range keys {
		if i < l {
			newMap[k] = m[k]
		} else {
			break
		}
	}
	return newMap
}

func Uptime(sh *discordgo.Session, msg *discordgo.Message, args []string) error {
	sh.ChannelMessageSend(msg.ChannelID, "I've been runing for "+
		time.Since(shared.Start).Round(time.Second).String())
	return nil
}

func Quote(sh *discordgo.Session, msg *discordgo.Message, args []string) error {
	if len(shared.Quotes) == 0 {
		sh.ChannelMessageSend(msg.ChannelID,
			"There are no shared.")
		return nil
	}

	who := ""
	num := 0
	if len(args) >= 2 {
		who = args[1]
	}
	if len(args) >= 3 {
		num64, err := strconv.ParseUint(args[2], 10, 32)
		if err != nil || num64 == 0 {
			sh.ChannelMessageSend(msg.ChannelID, "Provide a valid quote number (1+)")
			return nil
		}
		num = int(num64 - 1)
	}

	if who == "" {
		persons := make([]string, len(shared.Quotes))
		i := 0
		for p := range shared.Quotes {
			persons[i] = p
			i++
		}
		who = persons[rand.Intn(len(persons))]
	}

	qlist, ok := shared.Quotes[who]
	if !ok {
		sh.ChannelMessageSend(msg.ChannelID, "Could not find quote")
		return nil
	}

	if num == 0 {
		num = rand.Intn(len(qlist))
	}

	sh.ChannelMessageSend(msg.ChannelID,
		fmt.Sprintf("Quote **#%d** from **%s**:\n*“%s”*", num+1, who, qlist[num]))
	return nil
}

func AddQuote(sh *discordgo.Session, msg *discordgo.Message, args []string) error {
	if msg.Author.ID != shared.Config.Admin {
		return nil
	}
	if msg.ReferencedMessage == nil {
		return nil
	}
	if len(args) < 2 {
		return nil
	}
	num := shared.AddQuote(args[1], msg.ReferencedMessage.ContentWithMentionsReplaced())
	sh.ChannelMessageSend(msg.ChannelID,
		fmt.Sprintf("Added quote %d", num))
	return nil
}

func WTF(sh *discordgo.Session, msg *discordgo.Message, args []string) error {
	if len(args) >= 2 && args[1] == "+" {
		shared.WTF++
	}
	sh.ChannelMessageSend(msg.ChannelID,
		fmt.Sprintf("WTF moments: **%d**", shared.WTF))
	return nil
}

func compileGCC() (out string, badCode bool, err error) {
	gccOut := strings.Builder{}
	gcc := exec.Command("gcc", "-Wall", "-Wpedantic", "-o", "run.exe", "run.c")
	gcc.Stdout = &gccOut
	gcc.Stderr = &gccOut
	err = gcc.Run()
	if _, ok := err.(*exec.ExitError); ok {
		out = gccOut.String()
		badCode = true
		err = nil
		return
	}
	return
}

func compileClang() (out string, badCode bool, err error) {
	clangOut := strings.Builder{}
	clang := exec.Command("clang", "-Wall", "-Wpedantic", "-o", "run.exe", "run.c")
	clang.Stdout = &clangOut
	clang.Stderr = &clangOut
	err = clang.Run()
	if _, ok := err.(*exec.ExitError); ok {
		out = clangOut.String()
		badCode = true
		err = nil
		return
	}
	return
}

func compileCl(sanitize bool) (out string, badCode bool, err error) {
	args := []string{"/W4", "/WX", "/Ferun.exe"}
	if sanitize {
		args = append(args, "/fsanitize=address", "/RTC1")
	}
	args = append(args, "run.c")
	clOut := strings.Builder{}
	cl := exec.Command("cl.exe", args...)
	cl.Stderr = &clOut
	err = cl.Run()
	if _, ok := err.(*exec.ExitError); ok {
		out = clOut.String()
		badCode = true
		err = nil
		return
	}
	return
}
