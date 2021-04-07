package fp_test

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/ghostec/fp"
)

func TestMap(t *testing.T) {
	input := []int{13, 14, 15}

	strs := fp.
		Map(input, func(val, idx int) int { return val*100 + idx }).
		Map(func(val int) string { return fmt.Sprintf("%d, got it?", val) }).
		Get().([]string)

	if !reflect.DeepEqual(strs, []string{
		"1300, got it?",
		"1401, got it?",
		"1502, got it?",
	}) {
		t.Errorf("unexpected output. Got %v", strs)
	}
}
