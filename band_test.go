package learnermodel_test

import (
	"bytes"
	"errors"
	"slices"
	"strconv"
	"testing"

	lm "github.com/altnovatrail/learner-model"
)

func itoa(i int) string { return strconv.Itoa(i) }

func sortedKeys[V any](m map[string]V) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	slices.Sort(keys)
	return keys
}

// sameState tells whether two states encode to the same bytes.
func sameState(t *testing.T, a, b lm.State) bool {
	t.Helper()
	x, err := lm.EncodeState(a)
	if err != nil {
		t.Fatal(err)
	}
	y, err := lm.EncodeState(b)
	if err != nil {
		t.Fatal(err)
	}
	return bytes.Equal(x, y)
}

// The start band by grade (docs/architecture.md, section 3.1).
func TestBandByGrade(t *testing.T) {
	gap := []lm.Band{
		{Place: 1, ID: "early", Grades: []int32{2}},
		{Place: 2, ID: "none"},
		{Place: 3, ID: "late", Grades: []int32{4}},
	}
	noGrades := []lm.Band{{Place: 1, ID: "a"}, {Place: 2, ID: "b"}}
	known := func(g int32) lm.Grade { return lm.Grade{Value: g, Known: true} }
	cases := []struct {
		name  string
		grade lm.Grade
		bands []lm.Band
		want  string
	}{
		{"a band names the grade", known(3), ladder(), "3-4"},
		{"grade 7 goes to the nearest band", known(7), ladder(), "5-6"},
		{"grade 0 goes to the nearest band", known(0), ladder(), "1-2"},
		{"no grade starts at the first band", lm.Grade{}, ladder(), "1-2"},
		{"grade 0 is not no grade", known(0), gap, "early"},
		{"as near to two bands, the lower", known(3), gap, "early"},
		{"a band without grades is never nearest", known(4), gap, "late"},
		{"a module without grades starts at the first band", known(5), noGrades, "a"},
	}
	for _, c := range cases {
		got, err := lm.BandByGrade(c.grade, c.bands)
		if err != nil || got.ID != c.want || got.Place != c.bands[slices.IndexFunc(c.bands, func(b lm.Band) bool { return b.ID == c.want })].Place {
			t.Errorf("%s: %+v, %v; want %s", c.name, got, err, c.want)
		}
	}
	for _, bands := range [][]lm.Band{nil, {{Place: 2, ID: "x"}}, {{Place: 1, ID: "x"}, {Place: 2, ID: "x"}}} {
		if _, err := lm.BandByGrade(known(3), bands); !errors.Is(err, lm.ErrInvalid) {
			t.Errorf("bands %+v: %v, want ErrInvalid", bands, err)
		}
	}
}
