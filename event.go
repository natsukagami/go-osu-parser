package parser

import (
	"regexp"
	"sort"
	"strconv"
	"strings"
)

// BreakTime represents a beatmap's breaktime
type BreakTime struct {
	StartTime int
	EndTime   int
}

type breakTimeSorter []BreakTime

func (b breakTimeSorter) Len() int           { return len(b) }
func (b breakTimeSorter) Less(i, j int) bool { return b[i].StartTime < b[j].StartTime }
func (b breakTimeSorter) Swap(i, j int)      { b[i], b[j] = b[j], b[i] }

func sortBreakTimes(b []BreakTime) {
	sort.Sort(breakTimeSorter(b))
}

// Event is either a background image, or a break time?

func (b *Beatmap) parseEvent(line string) (err error) {
	members := strings.Split(line, ",")
	if members[0] == "0" && members[1] == "0" && members[2] != "" {
		bgName := strings.Trim(members[2], " ")
		if bgName[0] == '"' && bgName[len(bgName)-1] == '"' {
			b.BgFilename = bgName[1 : len(bgName)-1]
		} else {
			b.BgFilename = bgName
		}
	} else if members[0] == "2" {
		r := regexp.MustCompile("^[0-9]+$")
		if r.MatchString(members[2]) && r.MatchString(members[1]) {
			bt := BreakTime{}
			if bt.StartTime, err = strconv.Atoi(members[1]); err != nil {
				return
			}
			if bt.EndTime, err = strconv.Atoi(members[2]); err != nil {
				return
			}
			b.BreakTimes = append(b.BreakTimes, bt)
		}
	}
	return
}
