package parser

import "math"

// Compute the total time and the draining time of the beatmap.
func (b *Beatmap) computeDuration() {
	totalBreakTime := 0
	for _, bt := range b.BreakTimes {
		totalBreakTime += bt.EndTime - bt.StartTime
	}
	if len(b.HitObjects) > 0 {
		var (
			firstObj = b.HitObjects[0]
			lastObj  = b.HitObjects[len(b.HitObjects)-1]
		)
		b.TotalTime = int(math.Trunc(float64(lastObj.StartTime) / 1000))
		b.DrainingTime = int(math.Trunc(float64(lastObj.StartTime-firstObj.StartTime) / 1000))
	}
}

// Browse objects and compute max combo.
func (b *Beatmap) computeMaxCombo() {
	if len(b.TimingPoints) == 0 {
		return
	}
	maxCombo := 0
	var (
		sMul  = b.SliderMultiplier
		sTick = b.SliderTickRate
	)
	var (
		curTp TimingPoint
		nxOff = 0.0
		i     = 0
	)
	for _, h := range b.HitObjects {
		if float64(h.StartTime) >= nxOff {
			curTp = b.TimingPoints[i]
			i++
			nxOff = math.MaxFloat64
			if len(b.TimingPoints) > i {
				nxOff = b.TimingPoints[1].Offset
			}
		}
		switch h.ObjectName {
		case "spinner":
		case "circle":
			maxCombo++
		case "slider":
			var (
				osupxPerBeat = sMul * 100 * curTp.Velocity
				tickLength   = osupxPerBeat / float64(sTick)
				tickPerSide  = int(math.Ceil((math.Floor(float64(h.PixelLength)/tickLength*100) / 100) - 1))
			)
			maxCombo += (len(h.Edges)-1)*(tickPerSide+1) + 1 // 1 combo for each tick and endpoint
		}
	}
	b.MaxCombo = maxCombo
}
