package model

import (
	"math/rand"
	"os"
	"sort"
	"strings"

	"github.com/qeaml/autobot/shared"

	"github.com/bwmarrin/discordgo"
	"github.com/qeaml/ggpc"
)

var Memory *ggpc.Cache[string, map[string]uint64]
var Source *os.File

func Load() error {
	mCh, mf, err := shared.Cache[string, map[string]uint64]("model.json")
	if err != nil {
		return err
	}
	Memory = mCh
	Source = mf
	return nil
}

func Generate(starter string, count int) string {
	smdl := Memory.Solidify()
	starters := []string{}
	for st := range smdl {
		starters = append(starters, st)
	}
	lastW := starter
	out := lastW + " "
	for i := 0; i < count; i++ {
		wmap, ok := smdl[lastW]
		if !ok {
			if rand.Intn(16) == 3 {
				wmap = smdl[starters[rand.Intn(len(starters))]]
			} else {
				return out
			}
		}
		chcs := []string{}
		for chc := range wmap {
			chcs = append(chcs, chc)
		}
		sort.Slice(chcs, func(e, f int) bool {
			return wmap[chcs[e]] < wmap[chcs[f]]
		})
		max := 1 + len(chcs)/(1+rand.Intn(4))
		if max > len(chcs) {
			max = len(chcs)
		}
		lastW = chcs[rand.Intn(max)]
		out += lastW + " "
	}
	return out
}

func TrainString(msg *discordgo.Message) {
	recons := strings.TrimSpace(shared.Whitespacer.Replace(msg.Content))
	for strings.Contains(recons, "  ") {
		recons = strings.ReplaceAll(recons, "  ", " ")
	}
	words := strings.Split(" "+recons, " ")
	if len(words) < 2 {
		return
	}
	Train(words)
}

func Train(words []string) {
	var st, nd string
	for i := 1; i < len(words); i++ {
		st = words[i-1]
		nd = words[i]
		wmap := Memory.GetOrDefault(st, map[string]uint64{nd: 0})
		wmap[nd]++
		Memory.Set(st, wmap)
	}
}
