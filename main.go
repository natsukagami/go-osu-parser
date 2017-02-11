// Package parser provides an .osu beatmap parser.
// The parser output is made for maximum compability with the Node.js
// implementation which can be found at https://github.com/nojhamster/osu-parser.
// The output JSON has basically the same attributes with the Node.js implementation,
// however there are key differences:
//  - Empty additions are presented as null instead of an empty struct.
//  - Unlisted properties (in Beatmap object) are stored in a struct under "OtherAttributes",
//  some notable properties are "ComboX" and some other color properties.
package parser

import (
	"io/ioutil"
	"os"
	"strings"
)

// ParseError represents a parser error.
type ParseError string

func (p ParseError) Error() string { return string(p) }

// ParseFile parses a file given a filepath.
func ParseFile(file string) (b Beatmap, e error) {
	if stat, err := os.Stat(file); err != nil || stat.IsDir() {
		return b, ParseError("Invalid file")
	}
	var (
		bytes []byte
		err   error
	)
	if bytes, err = ioutil.ReadFile(file); err != nil {
		return b, err
	}
	return ParseBytes(bytes)
}

// ParseBytes parses a file given its contents as a byte array.
func ParseBytes(bytes []byte) (b Beatmap, err error) {
	return ParseString(string(bytes))
}

// ParseString parses a file given its contents as a string.
func ParseString(str string) (b Beatmap, err error) {
	p := newBeatmapParser()
	for _, line := range strings.Split(str, "\n") {
		p.ReadLine(line)
	}
	B, err := p.BuildBeatmap()
	if err != nil {
		return b, err
	}
	return *B, nil
}
