package parser

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"testing"
)

func runTest(num int) func(t *testing.T) {
	filename := fmt.Sprintf("testfiles/v%d.osu", num)
	return func(t *testing.T) {
		b, err := ParseFile(filename)
		if err != nil {
			t.Error(err)
		}
		if bytes, err := json.MarshalIndent(b, "", "\t"); err != nil {
			t.Error(err)
		} else if err := ioutil.WriteFile(fmt.Sprintf("testfiles/v%0d_out.json", num), bytes, 0644); err != nil {
			t.Error(err)
		}
	}
}

func TestAll(t *testing.T) {
	for i := 7; i <= 13; i++ {
		t.Run(fmt.Sprintf("Test v%d", i), runTest(i))
	}
}
