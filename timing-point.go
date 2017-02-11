package parser

import (
	"math"
	"sort"
	"strconv"
	"strings"
)

// TimingPoint is an osu! timing point.
type TimingPoint struct {
	Offset            int     `json:"offset"`
	BeatLength        float64 `json:"beatLength"`
	Velocity          float64 `json:"velocity"`
	Bpm               float64 `json:"bpm"`
	TimingSignature   int     `json:"timingSignature"`
	SampleSetID       int     `json:"sampleSetID"`
	CustomSampleIndex int     `json:"customSampleIndex"`
	SampleVolume      int     `json:"sampleVolume"`
	TimingChange      bool    `json:"timingChange"`
	KiaiTimeActive    bool    `json:"kiaiTimeActive"`
}

type timingPointSorter []TimingPoint

func (b timingPointSorter) Len() int           { return len(b) }
func (b timingPointSorter) Less(i, j int) bool { return b[i].Offset < b[j].Offset }
func (b timingPointSorter) Swap(i, j int)      { b[i], b[j] = b[j], b[i] }

func sortTimingPoints(b []TimingPoint) {
	sort.Sort(timingPointSorter(b))
}

// Gets the timing point affecting a specific offset.
func (b Beatmap) getTimingPoint(offset int) *TimingPoint {
	for i := len(b.TimingPoints) - 1; i >= 0; i-- {
		if b.TimingPoints[i].Offset <= offset {
			return &b.TimingPoints[i]
		}
	}
	return nil
}

// Parse a timing line
func (b *Beatmap) parseTimingPoint(line string) (err error) {
	members := strings.Split(line, ",")
	p := TimingPoint{}
	if p.Offset, err = strconv.Atoi(members[0]); err != nil {
		return
	}
	if p.BeatLength, err = strconv.ParseFloat(members[1], 64); err != nil {
		p.BeatLength = 0
	}
	p.Velocity = 1
	if p.TimingSignature, err = strconv.Atoi(members[2]); err != nil {
		return
	}
	if p.SampleSetID, err = strconv.Atoi(members[3]); err != nil {
		return
	}
	if p.CustomSampleIndex, err = strconv.Atoi(members[4]); err != nil {
		return
	}
	if p.SampleVolume, err = strconv.Atoi(members[5]); err != nil {
		return
	}
	x, err := strconv.Atoi(members[6])
	if err != nil {
		return
	}
	p.TimingChange = (x == 1)
	x, err = strconv.Atoi(members[7])
	if err != nil {
		return
	}
	p.KiaiTimeActive = (x == 1)
	if p.BeatLength != 0 {
		if p.BeatLength > 0 {
			p.Bpm = math.Trunc(60000.0/p.BeatLength + 0.5)
			b.BpmMin = math.Min(b.BpmMin, p.Bpm)
			b.BpmMax = math.Max(b.BpmMax, p.Bpm)
		} else {
			p.Velocity = math.Abs(100 / p.BeatLength)
		}
	}
	b.TimingPoints = append(b.TimingPoints, p)
	return
}
