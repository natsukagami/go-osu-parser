package parser

import (
	"strconv"
	"strings"
)

// Addition represents additional attributes of a hit object.
type Addition struct {
	Sample            string `json:"sample"`           // "normal", "soft", "drum"
	AdditionalSample  string `json:"additionalSample"` // "normal", "soft", "drum"
	CustomSampleIndex int    `json:"customSampleIndex"`
	HitsoundVolume    int    `json:"hitSoundVolume"`
	Hitsound          string `json:"hitSound"`
}

// Parse additional members
func parseAddition(str string) (p *Addition, err error) {
	a := Addition{}
	if len(str) == 0 {
		return nil, nil
	}
	adds := strings.Split(str, ":")
	if len(adds) > 0 && len(adds[0]) > 0 && adds[0] != "0" {
		switch adds[0] {
		case "1":
			a.Sample = "normal"

			p = &a
		case "2":
			a.Sample = "soft"

			p = &a
		case "3":
			a.Sample = "drum"

			p = &a
		}
	}
	if len(adds) > 1 && len(adds[1]) > 0 && adds[1] != "0" {
		switch adds[1] {
		case "1":
			a.AdditionalSample = "normal"

			p = &a
		case "2":
			a.AdditionalSample = "soft"

			p = &a
		case "3":
			a.AdditionalSample = "drum"

			p = &a
		}
	}
	if len(adds) > 2 && len(adds[2]) > 0 && adds[2] != "0" {
		a.CustomSampleIndex, err = strconv.Atoi(adds[2])
		if err != nil {
			return
		}
		p = &a
	}
	if len(adds) > 3 && len(adds[3]) > 0 && adds[3] != "0" {
		a.HitsoundVolume, err = strconv.Atoi(adds[3])
		if err != nil {
			return
		}
		p = &a
	}
	if len(adds) > 4 && len(adds[4]) > 0 {
		a.Hitsound = adds[4]

		p = &a
	}
	return
}
