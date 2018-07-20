package parser

import (
	"regexp"
	"strconv"
	"strings"
)

var (
	sectionReg = regexp.MustCompile("^\\[([a-zA-Z0-9]+)\\]$")
	keyValReg  = regexp.MustCompile("^([a-zA-Z0-9]+)[ ]*:[ ]*(.+)$")
)

type beatmapParser struct {
	*Beatmap
	TimingLines    []string
	HitObjectLines []string
	EventLines     []string
	OsuSection     string
}

func (b *beatmapParser) ReadLine(line string) (err error) {
	line = strings.Trim(line, " \r\n")
	if len(line) == 0 {
		return
	}
	if match := sectionReg.FindStringSubmatch(line); match != nil {
		b.OsuSection = strings.ToLower(match[1])
		return
	}
	switch b.OsuSection {
	case "timingpoints":
		b.TimingLines = append(b.TimingLines, line)
	case "hitobjects":
		b.HitObjectLines = append(b.HitObjectLines, line)
	case "events":
		b.EventLines = append(b.EventLines, line)
	default:
		if b.OsuSection == "" {
			fmtRegex := regexp.MustCompile("^osu file format (v[0-9]+)$")
			if match := fmtRegex.FindStringSubmatch(line); match != nil {
				b.FileFormat = match[1]
				return
			}
		}
		// Apart from events, timingpoints and hitobjects sections, lines are "key: value"
		if match := keyValReg.FindStringSubmatch(line); match != nil {
			switch match[1] {
			case "SliderMultiplier":
				if b.SliderMultiplier, err = strconv.ParseFloat(match[2], 64); err != nil {
					return
				}
			case "SliderTickRate":
				if b.SliderTickRate, err = strconv.Atoi(match[2]); err != nil {
					return
				}
			case "Artist":
				b.Artist = match[2]
			case "ArtistUnicode":
				b.ArtistUnicode = match[2]
			case "Title":
				b.Title = match[2]
			case "TitleUnicode":
				b.TitleUnicode = match[2]
			case "AudioFilename":
				b.AudioFilename = match[2]
			case "Creator":
				b.Creator = match[2]
			case "Source":
				b.Source = match[2]
			case "Version":
				b.Version = match[2]
			case "BeatmapID":
				if b.BeatmapID, err = strconv.Atoi(match[2]); err != nil {
					return
				}
			case "BeatmapSetID":
				if b.BeatmapSetID, err = strconv.Atoi(match[2]); err != nil {
					return
				}
			case "FileFormat":
				b.FileFormat = match[2]
			case "Mode":
				if b.Mode, err = strconv.Atoi(match[2]); err != nil {
					return
				}
			case "AudioLeadIn":
				if b.AudioLeadIn, err = strconv.Atoi(match[2]); err != nil {
					return
				}
			case "SampleSet":
				b.SampleSet = match[2]
			case "Countdown":
				if b.Countdown, err = strconv.Atoi(match[2]); err != nil {
					return
				}
			case "BeatDivisor":
				if b.BeatDivisor, err = strconv.Atoi(match[2]); err != nil {
					return
				}
			case "StackLeniency":
				if b.StackLeniency, err = strconv.ParseFloat(match[2], 64); err != nil {
					return
				}
			case "DistanceSpacing":
				if b.DistanceSpacing, err = strconv.Atoi(match[2]); err != nil {
					return
				}
			case "GridSize":
				if b.GridSize, err = strconv.Atoi(match[2]); err != nil {
					return
				}
			case "LetterboxInBreaks":
				b.LetterboxInBreaks = (match[2] == "1")
			case "PreviewTime":
				if b.PreviewTime, err = strconv.Atoi(match[2]); err != nil {
					return
				}
			case "CircleSize":
				if b.CircleSize, err = strconv.ParseFloat(match[2], 64); err != nil {
					return
				}
			case "HPDrainRate":
				if b.HPDrainRate, err = strconv.ParseFloat(match[2], 64); err != nil {
					return
				}
			case "OverallDifficulty":
				if b.OverallDifficulty, err = strconv.ParseFloat(match[2], 64); err != nil {
					return
				}
			case "ApproachRate":
				if b.ApproachRate, err = strconv.ParseFloat(match[2], 64); err != nil {
					return
				}
			default:
				b.OtherAttributes[match[1]] = match[2]
			}
		}
	}
	return
}

func (b *beatmapParser) BuildBeatmap() (*Beatmap, error) {
	if tags, ok := b.OtherAttributes["Tags"]; ok {
		b.Tags = strings.Split(tags, " ")
		delete(b.OtherAttributes, "Tags")
	}
	var err error
	for _, line := range b.EventLines {
		if err = b.parseEvent(line); err != nil {
			return nil, err
		}
	}
	sortBreakTimes(b.BreakTimes)
	for _, line := range b.TimingLines {
		if err = b.parseTimingPoint(line); err != nil {
			return nil, err
		}
	}
	sortTimingPoints(b.TimingPoints)
	for i := 1; i < len(b.TimingPoints); i++ {
		if b.TimingPoints[i].Bpm == 0 {
			b.TimingPoints[i].Bpm = b.TimingPoints[i-1].Bpm
			b.TimingPoints[i].BeatLength = b.TimingPoints[i-1].BeatLength
		}
	}
	for _, line := range b.HitObjectLines {
		if err = b.parseHitObject(line); err != nil {
			return nil, err
		}
	}
	sortHitObjects(b.HitObjects)
	b.computeMaxCombo()
	b.computeDuration()
	return b.Beatmap, nil
}

func newBeatmapParser() beatmapParser {
	b := beatmapParser{}
	b.Beatmap = newBeatmap()
	b.EventLines = make([]string, 0)
	b.HitObjectLines = make([]string, 0)
	b.TimingLines = make([]string, 0)
	return b
}
