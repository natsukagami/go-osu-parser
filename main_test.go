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
		} else if ans, err := ioutil.ReadFile(fmt.Sprintf("testfiles/v%d_out.json", num)); err != nil {
			t.Error(err)
		} else if string(bytes) != string(ans) {
			t.Errorf("Test v%d: Output are not the same", num)
		}
	}
}

func TestAll(t *testing.T) {
	for i := 7; i <= 13; i++ {
		t.Run(fmt.Sprintf("Test v%d", i), runTest(i))
	}
}
