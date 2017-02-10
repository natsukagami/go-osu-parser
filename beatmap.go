package parser

// Beatmap is the returned struct, representing
// an osu! beatmap.
type Beatmap struct {

	// Metadata
	Artist        string
	ArtistUnicode string
	Title         string
	TitleUnicode  string
	AudioFilename string
	Creator       string
	Source        string
	Tags          []string
	// Game Metadata
	Version           string
	BeatmapID         int
	BeatmapSetID      int
	FileFormat        string
	Mode              int
	AudioLeadIn       int
	SampleSet         string
	BgFilename        string
	Countdown         int
	BeatDivisor       int
	StackLeniency     float64
	DistanceSpacing   int
	GridSize          int
	LetterboxInBreaks bool
	PreviewTime       int
	CircleSize        int
	HPDrainRate       int
	OverallDifficulty int
	ApproachRate      int
	// Beatmap information
	NbCircles        int
	NbSliders        int
	NbSpinners       int
	TotalTime        int
	DrainingTime     int
	MaxCombo         int
	BpmMin           float64
	BpmMax           float64
	SliderMultiplier float64
	SliderTickRate   int
	TimingPoints     []TimingPoint
	HitObjects       []HitObject
	BreakTimes       []BreakTime
	OtherAttributes  map[string]string
}

func newBeatmap() *Beatmap {
	b := Beatmap{}
	b.TimingPoints = make([]TimingPoint, 0)
	b.HitObjects = make([]HitObject, 0)
	b.BreakTimes = make([]BreakTime, 0)
	b.OtherAttributes = make(map[string]string)
	return &b
}
