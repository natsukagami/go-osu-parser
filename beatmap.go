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
	FileFormat        int `json:"fileFormat"`
	Mode              int
	AudioLeadIn       int
	SampleSet         string
	BgFilename        string `json:"bgFilename"`
	Countdown         int
	BeatDivisor       int
	StackLeniency     float64
	DistanceSpacing   int
	GridSize          int
	LetterboxInBreaks bool
	PreviewTime       int
	CircleSize        float64
	HPDrainRate       float64
	OverallDifficulty float64
	ApproachRate      float64
	// Beatmap information
	NbCircles        int     `json:"nbCircles"`
	NbSliders        int     `json:"nbSliders"`
	NbSpinners       int     `json:"nbSpinners"`
	TotalTime        int     `json:"totalTime"`
	DrainingTime     int     `json:"drainingTime"`
	MaxCombo         int     `json:"maxCombo"`
	BpmMin           float64 `json:"bpmMin"`
	BpmMax           float64 `json:"bpmMax"`
	SliderMultiplier float64
	SliderTickRate   int
	TimingPoints     []TimingPoint `json:"timingPoints"`
	HitObjects       []HitObject   `json:"hitObjects"`
	BreakTimes       []BreakTime   `json:"breakTimes"`
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
