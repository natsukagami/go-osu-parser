package parser

import (
	"math"
	"sort"
	"strconv"
	"strings"
)

var curveTypes = map[string]string{
	"C": "catmull",
	"B": "bezier",
	"L": "linear",
	"P": "pass-through",
}

// HitObject represents an osu! hit object.
type HitObject struct {
	ObjectName  string    `json:"objectName"` // "slider", "spinner", "circle"
	StartTime   int       `json:"startTime"`
	EndTime     int       `json:"endTime"`     // Spinner only
	RepeatCount int       `json:"repeatCount"` // Sliders only
	PixelLength float64   `json:"pixelLength"` // Sliders only
	Points      []Point   `json:"points"`      // Sliders only
	Duration    int       `json:"duration"`    // Sliders only
	CurveType   string    `json:"curveType"`   // Sliders only, "catmull", "bezier", "linear", "pass-through"
	EndPosition Point     `json:"endPosition"` // Sliders only
	Edges       []Edge    `json:"edges"`
	NewCombo    bool      `json:"newCombo"`
	SoundTypes  []string  `json:"soundTypes"` // contains "whistle", "finish", "clap", "normal"
	Position    Point     `json:"position"`
	Additions   *Addition `json:"additions"`
}

type hitObjectSorter []HitObject

func (b hitObjectSorter) Len() int           { return len(b) }
func (b hitObjectSorter) Less(i, j int) bool { return b[i].StartTime < b[j].StartTime }
func (b hitObjectSorter) Swap(i, j int)      { b[i], b[j] = b[j], b[i] }

func sortHitObjects(b []HitObject) {
	sort.Sort(hitObjectSorter(b))
}

// Edge represents a slider edge.
type Edge struct {
	SoundTypes []string  `json:"soundTypes"`
	Additions  *Addition `json:"addictions"`
}

// Gets sound types
func parseSoundType(t int) []string {
	a := make([]string, 0)
	/**
	 * sound type is a bitwise flag enum
	 * 0 : normal
	 * 2 : whistle
	 * 4 : finish
	 * 8 : clap
	 */
	if (t & 2) > 0 {
		a = append(a, "whistle")
	}
	if (t & 4) > 0 {
		a = append(a, "finish")
	}
	if (t & 8) > 0 {
		a = append(a, "clap")
	}
	if len(a) == 0 {
		a = append(a, "normal")
	}
	return a
}

func (b *Beatmap) parseHitObject(line string) (err error) {
	h := HitObject{}
	members := strings.Split(line, ",")
	var (
		soundType, objectType int
	)
	if soundType, err = strconv.Atoi(members[4]); err != nil {
		return
	}
	if objectType, err = strconv.Atoi(members[3]); err != nil {
		return
	}
	if h.StartTime, err = strconv.Atoi(members[2]); err != nil {
		return
	}
	h.NewCombo = (objectType & 4) > 0
	h.SoundTypes = make([]string, 0)
	h.Edges = make([]Edge, 0)
	if h.Position, err = parsePoint(members[0], members[1]); err != nil {
		return
	}
	h.SoundTypes = parseSoundType(soundType)
	/**
	 * object type is a bitwise flag enum
	 * 1: circle
	 * 2: slider
	 * 8: spinner
	 */
	if (objectType & 1) > 0 {
		// Circle
		h.ObjectName = "circle"
		b.NbCircles++
		if len(members) > 5 {
			if h.Additions, err = parseAddition(members[5]); err != nil {
				return
			}
		}
	} else if (objectType & 8) > 0 {
		h.ObjectName = "spinner"
		b.NbSpinners++
		if h.EndTime, err = strconv.Atoi(members[5]); err != nil {
			return
		}
		if len(members) > 6 {
			if h.Additions, err = parseAddition(members[6]); err != nil {
				return
			}
		}
	} else if (objectType & 2) > 0 {
		h.ObjectName = "slider"
		b.NbSliders++
		if h.RepeatCount, err = strconv.Atoi(members[6]); err != nil {
			return
		}
		if h.PixelLength, err = strconv.ParseFloat(members[7], 64); err != nil {
			return
		}
		if len(members) > 10 {
			if h.Additions, err = parseAddition(members[10]); err != nil {
				return
			}
		}
		h.Points = []Point{h.Position}
		/**
		 * Calculate slider duration
		 */
		timing := b.getTimingPoint(h.StartTime)
		if timing != nil {
			pxPerBeat := b.SliderMultiplier * 100 * timing.Velocity
			beatsNumber := h.PixelLength * float64(h.RepeatCount) / pxPerBeat
			h.Duration = int(math.Ceil(beatsNumber * timing.BeatLength))
			h.EndTime = h.StartTime + h.Duration
		}
		/**
		 * Parse slider points
		 */
		points := strings.Split(members[5], "|")
		if len(points) > 0 {
			if typ, ok := curveTypes[points[0]]; ok {
				h.CurveType = typ
			} else {
				h.CurveType = "unknown"
			}
			for i := 1; i < len(points); i++ {
				coords := strings.Split(points[i], ":")
				var x Point
				if x, err = parsePoint(coords[0], coords[1]); err != nil {
					return
				}
				h.Points = append(h.Points, x)
			}
		}
		edgeSounds := make([]string, h.RepeatCount+1)
		edgeAdditions := make([]string, h.RepeatCount+1)
		if len(members) > 8 && len(members[8]) > 0 {
			edgeSounds = strings.Split(members[8], "|")
		}
		if len(members) > 9 && len(members[9]) > 0 {
			edgeAdditions = strings.Split(members[9], "|")
		}
		/**
		 * Get soundTypes and additions for each slider edge
		 */
		for j := 0; j < h.RepeatCount+1; j++ {
			edge := Edge{}
			if edge.Additions, err = parseAddition(edgeAdditions[j]); err != nil {
				return
			}
			if len(edgeSounds[j]) > 0 {
				var sound int
				if sound, err = strconv.Atoi(edgeSounds[j]); err != nil {
					return
				}
				edge.SoundTypes = parseSoundType(sound)
			} else {
				edge.SoundTypes = []string{"normal"}
			}
			h.Edges = append(h.Edges, edge)
		}
		// get coordinates of the slider endpoint
		endPoint := getSliderEndPoint(h.CurveType, float64(h.PixelLength), h.Points)
		if endPoint.X != 0 && endPoint.Y != 0 {
			h.EndPosition = Point{math.Trunc(endPoint.X + 0.5), math.Trunc(endPoint.Y + 0.5)}
		} else {
			// If endPosition could not be calculated, approximate it by setting it to the last point
			h.EndPosition = h.Points[len(h.Points)-1]
		}
	} else {
		h.ObjectName = "unknown"
	}
	b.HitObjects = append(b.HitObjects, h)
	return
}
