package learnermodel_test

import (
	"bytes"
	"errors"
	"fmt"
	"math"
	"math/rand/v2"
	"os"
	"path/filepath"
	"slices"
	"strconv"
	"strings"
	"testing"

	lm "github.com/altnovatrail/learner-model"
)

// The golden bytes of the codec, testdata/wire/*.binpb, are one for every build, and their JSON
// beside them is what buf convert gives of them by the schema: just ci-proto holds the two to
// each other, both ways. LEARNERMODEL_UPDATE=1 writes the bytes anew; the JSON, ci-proto's
// update. Every value here is synthetic.

const updateEnv = "LEARNERMODEL_UPDATE"

const module = "example/workbook"

func skill(name string) string { return module + "/skill/" + name }

func goldenEvidence() lm.Evidence {
	return lm.Evidence{
		Kind: lm.EvidenceKind,
		Task: lm.Task{
			ID:                  module + "/task/page-7",
			Band:                lm.BandRef{Place: 2, ID: "grade-3"},
			Difficulty:          3,
			Primary:             []string{skill("subtraction"), skill("addition")},
			Supporting:          []string{"core.counting"},
			KeyProvenance:       lm.ProvenanceAuthor,
			SkillTagsProvenance: lm.ProvenanceAuthor,
			Items: []lm.TaskItem{
				{Place: 2, OwnSkills: true, Primary: []string{skill("comparison")},
					Slots: []lm.Slot{{Check: lm.CheckKindChoice, Options: 3, HasOptions: true}}},
				{Place: 1, Slots: []lm.Slot{{Check: lm.CheckKindNumber}, {Check: lm.CheckKindFraction}}},
			},
		},
		StartBand:           lm.BandRef{Place: 1, ID: "grade-2"},
		Outcome:             lm.Outcome{Kind: lm.OutcomeAfterHint, Rung: 1, HasRung: true},
		PracticeTime:        1762041600,
		TechniqueOnRequest:  true,
		ExplainAgainCount:   2,
		SolutionSeenBefore:  false,
		HasSolutionSeenFact: true,
		Items: []lm.ItemAnswer{
			{Place: 1, Submitted: true, Unverdicted: []lm.UnverdictedReason{lm.UnverdictedUnreadableBlank},
				FirstVerdict: lm.VerdictNotYet, HasFirstVerdict: true,
				FinalVerdict: lm.VerdictCorrect, HasFinalVerdict: true,
				WrongBeforeFirstCorrect: 1, HasFirstCorrect: true},
			{Place: 2, Submitted: true, FirstVerdict: lm.VerdictCorrect, HasFirstVerdict: true, WrongBeforeFirst: 1,
				FinalVerdict: lm.VerdictCorrect, HasFinalVerdict: true, WrongBeforeFirstCorrect: 1, HasFirstCorrect: true},
		},
		Traps: []lm.FoundTrap{{ID: module + "/trap/borrow-forgotten", Skills: []string{skill("subtraction")}, ItemPlaces: []int32{1}}},
	}
}

func goldenState() lm.State {
	return lm.State{
		Kind: lm.StateKind,
		Modules: []lm.ModuleState{{
			Module: module, Level: 0.5, AnswerWeight: 22.5, StartBand: lm.BandRef{Place: 1, ID: "grade-2"},
			Trial: []lm.TrialAnswer{{BandPlace: 1, Difficulty: 3, Slots: []lm.Slot{{Check: lm.CheckKindNumber}},
				Verdict: lm.VerdictCorrect, Weight: 0.25}},
			Traps: []lm.TrapCount{{ID: module + "/trap/digits-swapped", Count: 1, Last: 1762041600}},
			Skills: []lm.SkillState{
				{Skill: skill("subtraction"), Correction: math.Copysign(0, -1), AnswerWeight: 10.5, EvidenceCount: 4,
					LastPractice: 1762041600, Stability: 1.81, FailuresInARow: 1,
					Traps: []lm.TrapCount{{ID: module + "/trap/borrow-forgotten", Count: 2, Last: 1762041600}}},
				{Skill: skill("addition"), Correction: 0.58, AnswerWeight: 5, EvidenceCount: 2, LastPractice: -5,
					Stability: 2.5, MasteryBand: 2, Mastered: true},
			},
		}},
	}
}

// canonical gives a value as the decoder returns the bytes of the encoder: lists in their
// canonical order, and nil for empty.
func canonicalEvidence(e lm.Evidence) lm.Evidence {
	b, err := lm.EncodeEvidence(e)
	if err != nil {
		panic(err)
	}
	d, err := lm.DecodeEvidence(b)
	if err != nil {
		panic(err)
	}
	return d
}

func golden(t *testing.T, name string, got []byte) {
	t.Helper()
	path := filepath.Join("testdata", "wire", name)
	if os.Getenv(updateEnv) == "1" {
		if err := os.WriteFile(path, got, 0o644); err != nil {
			t.Fatal(err)
		}
	}
	want, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(got, want) {
		t.Errorf("%s: the encoder gives %x, the file has %x", path, got, want)
	}
}

func TestGoldenEvidence(t *testing.T) {
	b, err := lm.EncodeEvidence(goldenEvidence())
	if err != nil {
		t.Fatal(err)
	}
	golden(t, "evidence.binpb", b)
	e, err := lm.DecodeEvidence(b)
	if err != nil {
		t.Fatal(err)
	}
	again, err := lm.EncodeEvidence(e)
	if err != nil || !bytes.Equal(again, b) {
		t.Errorf("encode(decode(b)) = %x, %v; not b", again, err)
	}
	if e.Task.Items[0].Place != 1 || e.Task.Primary[0] != skill("addition") {
		t.Errorf("the lists are not in their canonical order: %+v", e.Task)
	}
}

func TestGoldenState(t *testing.T) {
	b, err := lm.EncodeState(goldenState())
	if err != nil {
		t.Fatal(err)
	}
	golden(t, "state.binpb", b)
	s, err := lm.DecodeState(b)
	if err != nil {
		t.Fatal(err)
	}
	again, err := lm.EncodeState(s)
	if err != nil || !bytes.Equal(again, b) {
		t.Errorf("encode(decode(b)) = %x, %v; not b", again, err)
	}
	sk := s.Modules[0].Skills
	if sk[0].Skill != skill("addition") || !sk[0].Mastered || sk[1].Mastered {
		t.Errorf("skills: %+v", sk)
	}
	if !math.Signbit(sk[1].Correction) {
		t.Errorf("−0 did not survive: %v", sk[1].Correction)
	}
}

// The round trip on random values: the decoder gives back what the encoder was given, in its
// canonical order, and the bytes of the two are the same.
func TestRoundTrip(t *testing.T) {
	r := rand.New(rand.NewPCG(5, 0))
	for i := range 2000 {
		e := randomEvidence(r)
		b, err := lm.EncodeEvidence(e)
		if err != nil {
			t.Fatalf("%d: %v", i, err)
		}
		d, err := lm.DecodeEvidence(b)
		if err != nil {
			t.Fatalf("%d: %v", i, err)
		}
		if !equalEvidence(d, normalEvidence(e)) {
			t.Fatalf("%d: decoded %+v, not %+v", i, d, e)
		}
		s := randomState(r)
		sb, err := lm.EncodeState(s)
		if err != nil {
			t.Fatalf("%d: %v", i, err)
		}
		ds, err := lm.DecodeState(sb)
		if err != nil {
			t.Fatalf("%d: %v", i, err)
		}
		if got, want := fmt.Sprintf("%#v", ds), fmt.Sprintf("%#v", normalState(s)); got != want {
			t.Fatalf("%d: decoded %s, not %s", i, got, want)
		}
		again, _ := lm.EncodeState(ds)
		if !bytes.Equal(again, sb) {
			t.Fatalf("%d: the state did not come back to its bytes", i)
		}
	}
}

func randomID(r *rand.Rand, kind string) string {
	return module + "/" + kind + "/s" + strconv.Itoa(int(r.Uint32N(1000)))
}

func randomIDs(r *rand.Rand, kind string) []string {
	var ids []string
	for range r.IntN(4) {
		id := randomID(r, kind)
		if !slices.Contains(ids, id) {
			ids = append(ids, id)
		}
	}
	return ids
}

func randomBand(r *rand.Rand) lm.BandRef {
	return lm.BandRef{Place: r.Int32N(6) + 1, ID: "band-" + strconv.Itoa(r.IntN(6))}
}

func randomSlots(r *rand.Rand) []lm.Slot {
	slots := make([]lm.Slot, r.IntN(3)+1)
	for i := range slots {
		slots[i].Check = lm.CheckKind(r.Int32N(6))
		if slots[i].Check == lm.CheckKindChoice {
			slots[i].Options, slots[i].HasOptions = r.Int32N(5)+2, true
		}
	}
	return slots
}

func randomEvidence(r *rand.Rand) lm.Evidence {
	e := lm.Evidence{
		Kind: lm.EvidenceKind,
		Task: lm.Task{
			ID: randomID(r, "task"), Template: r.IntN(2) == 0, Band: randomBand(r), Difficulty: r.Int32N(5) + 1,
			Primary: randomIDs(r, "skill"), Supporting: randomIDs(r, "skill"),
			KeyProvenance: lm.Provenance(r.Int32N(5)), SkillTagsProvenance: lm.Provenance(r.Int32N(5)),
		},
		StartBand:          randomBand(r),
		Outcome:            lm.Outcome{Kind: lm.OutcomeKind(r.Int32N(8))},
		PracticeTime:       r.Int64() - math.MaxInt64/2,
		TechniqueOnStep:    r.IntN(2) == 0,
		TechniqueOnRequest: r.IntN(2) == 0,
		ExplainAgainCount:  r.Int32N(4),
	}
	if e.Outcome.Kind == lm.OutcomeAfterHint {
		e.Outcome.Rung, e.Outcome.HasRung = r.Int32N(3)+1, true
	}
	if !e.Task.Template {
		e.SolutionSeenBefore, e.HasSolutionSeenFact = r.IntN(2) == 0, true
	}
	for _, place := range r.Perm(r.IntN(12)) {
		item := lm.TaskItem{Place: int32(place) + 1, Slots: randomSlots(r)}
		if r.IntN(3) == 0 {
			item.OwnSkills, item.Primary, item.Supporting = true, randomIDs(r, "skill"), randomIDs(r, "skill")
		}
		e.Task.Items = append(e.Task.Items, item)
		a := lm.ItemAnswer{Place: item.Place, Submitted: r.IntN(4) > 0}
		for range r.IntN(3) {
			a.Unverdicted = append(a.Unverdicted, lm.UnverdictedReason(r.Int32N(4)))
		}
		if r.IntN(4) > 0 {
			a.FirstVerdict, a.HasFirstVerdict, a.WrongBeforeFirst = lm.Verdict(r.Int32N(4)), true, r.Int32N(4)
			a.FinalVerdict, a.HasFinalVerdict = lm.Verdict(r.Int32N(3)), true
		}
		if r.IntN(2) == 0 {
			a.WrongBeforeFirstCorrect, a.HasFirstCorrect = r.Int32N(4), true
		}
		e.Items = append(e.Items, a)
	}
	for _, id := range randomIDs(r, "trap") {
		trap := lm.FoundTrap{ID: id, Skills: randomIDs(r, "skill")}
		for _, p := range r.Perm(3)[:r.IntN(3)] {
			trap.ItemPlaces = append(trap.ItemPlaces, int32(p)+1)
		}
		e.Traps = append(e.Traps, trap)
	}
	return e
}

func randomTraps(r *rand.Rand) []lm.TrapCount {
	var traps []lm.TrapCount
	for _, id := range randomIDs(r, "trap") {
		traps = append(traps, lm.TrapCount{ID: id, Count: r.Int64N(100) + 1, Last: r.Int64()})
	}
	return traps
}

// randomFloat draws any finite float64, its bits uniform: −0, subnormals and the largest too.
func randomFloat(r *rand.Rand) float64 {
	for {
		if v := math.Float64frombits(r.Uint64()); !math.IsNaN(v) && !math.IsInf(v, 0) {
			return v
		}
	}
}

func randomState(r *rand.Rand) lm.State {
	s := lm.State{Kind: lm.StateKind}
	for _, id := range randomIDs(r, "module") {
		m := lm.ModuleState{Module: id, Level: randomFloat(r), AnswerWeight: randomFloat(r), StartBand: randomBand(r),
			Traps: randomTraps(r)}
		for range r.IntN(4) {
			m.Trial = append(m.Trial, lm.TrialAnswer{BandPlace: r.Int32N(6) + 1, Difficulty: r.Int32N(5) + 1,
				Slots: randomSlots(r), Verdict: lm.Verdict(r.Int32N(4)), Weight: randomFloat(r)})
		}
		for _, sid := range randomIDs(r, "skill") {
			sk := lm.SkillState{Skill: sid, Correction: randomFloat(r), AnswerWeight: randomFloat(r),
				EvidenceCount: r.Int64(), LastPractice: r.Int64() - math.MaxInt64/2, Stability: randomFloat(r),
				Traps: randomTraps(r), FailuresInARow: r.Int32N(3)}
			if r.IntN(2) == 0 {
				sk.MasteryBand, sk.Mastered = r.Int32N(6)+1, true
			}
			m.Skills = append(m.Skills, sk)
		}
		s.Modules = append(s.Modules, m)
	}
	return s
}

// equalEvidence compares two pieces of evidence field by field, as Go prints them: %#v prints
// every field, a float as the shortest text that reads back to its bits, −0 too.
func equalEvidence(a, b lm.Evidence) bool {
	return fmt.Sprintf("%#v", a) == fmt.Sprintf("%#v", b)
}

// sortedNil sorts a copy of s, and gives nil for an empty one, as the decoder does.
func sortedNil[T int32 | string](s []T) []T {
	if len(s) == 0 {
		return nil
	}
	c := slices.Clone(s)
	slices.Sort(c)
	return c
}

func orNil[T any](s []T) []T {
	if len(s) == 0 {
		return nil
	}
	return s
}

// normalEvidence puts e as the decoder gives it, by the rules of the schema written out here and
// not by the encoder: lists of IDs and places in order, items by place, traps by ID, nil for empty.
func normalEvidence(e lm.Evidence) lm.Evidence {
	t := &e.Task
	t.Primary, t.Supporting = sortedNil(t.Primary), sortedNil(t.Supporting)
	t.Items = orNil(slices.Clone(t.Items))
	slices.SortFunc(t.Items, func(a, b lm.TaskItem) int { return int(a.Place - b.Place) })
	for i := range t.Items {
		it := &t.Items[i]
		it.Primary, it.Supporting, it.Slots = sortedNil(it.Primary), sortedNil(it.Supporting), orNil(it.Slots)
	}
	e.Items = orNil(slices.Clone(e.Items))
	slices.SortFunc(e.Items, func(a, b lm.ItemAnswer) int { return int(a.Place - b.Place) })
	for i := range e.Items {
		e.Items[i].Unverdicted = orNil(e.Items[i].Unverdicted)
	}
	e.Traps = orNil(slices.Clone(e.Traps))
	slices.SortFunc(e.Traps, func(a, b lm.FoundTrap) int { return strings.Compare(a.ID, b.ID) })
	for i := range e.Traps {
		e.Traps[i].Skills, e.Traps[i].ItemPlaces = sortedNil(e.Traps[i].Skills), sortedNil(e.Traps[i].ItemPlaces)
	}
	return e
}

func normalTraps(traps []lm.TrapCount) []lm.TrapCount {
	traps = orNil(slices.Clone(traps))
	slices.SortFunc(traps, func(a, b lm.TrapCount) int { return strings.Compare(a.ID, b.ID) })
	return traps
}

// normalState puts s as the decoder gives it: modules, skills and traps by ID, nil for empty.
func normalState(s lm.State) lm.State {
	s.Modules = orNil(slices.Clone(s.Modules))
	slices.SortFunc(s.Modules, func(a, b lm.ModuleState) int { return strings.Compare(a.Module, b.Module) })
	for i := range s.Modules {
		m := &s.Modules[i]
		m.Traps, m.Trial = normalTraps(m.Traps), orNil(slices.Clone(m.Trial))
		for j := range m.Trial {
			m.Trial[j].Slots = orNil(m.Trial[j].Slots)
		}
		m.Skills = orNil(slices.Clone(m.Skills))
		slices.SortFunc(m.Skills, func(a, b lm.SkillState) int { return strings.Compare(a.Skill, b.Skill) })
		for j := range m.Skills {
			m.Skills[j].Traps = normalTraps(m.Skills[j].Traps)
		}
	}
	return s
}

// Bytes cut inside a field are refused; a cut between two fields of the top level may read, as
// protobuf has it, but only as a value without the fields after the cut.
func TestTruncated(t *testing.T) {
	e, _ := lm.EncodeEvidence(goldenEvidence())
	s, _ := lm.EncodeState(goldenState())
	for _, c := range []struct {
		name string
		b    []byte
		read func([]byte) error
	}{
		{"evidence", e, func(b []byte) error { _, err := lm.DecodeEvidence(b); return err }},
		{"state", s, func(b []byte) error { _, err := lm.DecodeState(b); return err }},
	} {
		bounds := topLevelBounds(t, c.b)
		for n := range len(c.b) {
			err := c.read(c.b[:n])
			if !slices.Contains(bounds, n) && !errors.Is(err, lm.ErrCorrupt) {
				t.Errorf("%s cut at %d of %d: %v, not ErrCorrupt", c.name, n, len(c.b), err)
			}
		}
	}
}

// topLevelBounds gives the offsets where a top-level field of b ends.
func topLevelBounds(t *testing.T, b []byte) []int {
	t.Helper()
	var bounds []int
	for off := 0; off < len(b); {
		tag, n := uvarint(b[off:])
		off += n
		switch tag & 7 {
		case 0:
			_, n = uvarint(b[off:])
			off += n
		case 1:
			off += 8
		case 2:
			l, n := uvarint(b[off:])
			off += n + int(l)
		}
		bounds = append(bounds, off)
	}
	return bounds
}

func uvarint(b []byte) (uint64, int) {
	var v uint64
	for i, c := range b {
		v |= uint64(c&0x7f) << (7 * i)
		if c < 0x80 {
			return v, i + 1
		}
	}
	return 0, len(b)
}

func TestCorrupt(t *testing.T) {
	nan := math.Float64bits(math.Float64frombits(0x7FF8000000000001))
	fixed := func(v uint64) []byte {
		b := make([]byte, 8)
		for i := range b {
			b[i] = byte(v >> (8 * i))
		}
		return b
	}
	state := func(module []byte) []byte {
		return slices.Concat([]byte{0x08, 0x01, 0x12, byte(len(module))}, module)
	}
	for _, c := range []struct {
		name string
		b    []byte
		err  error
	}{
		{"field number 0", []byte{0x00, 0x01}, lm.ErrCorrupt},
		{"a group", []byte{0x0b}, lm.ErrCorrupt},
		{"wire type 7", []byte{0x0f}, lm.ErrCorrupt},
		{"a varint past 10 bytes", slices.Concat([]byte{0x08}, bytes.Repeat([]byte{0xff}, 10), []byte{0x01}), lm.ErrCorrupt},
		{"a length past the end", []byte{0x12, 0x05, 0x00}, lm.ErrCorrupt},
		{"the kind as a string", []byte{0x0a, 0x00}, lm.ErrCorrupt},
		{"the kind twice", []byte{0x08, 0x01, 0x08, 0x01}, lm.ErrCorrupt},
		{"an int32 past its range", slices.Concat([]byte{0x08}, []byte{0x80, 0x80, 0x80, 0x80, 0x10}), lm.ErrCorrupt},
		{"no kind", []byte{}, lm.ErrCorrupt},
		{"a module without its level", state([]byte{0x0a, 0x01, 'm'}), lm.ErrCorrupt},
		{"a level that is NaN", state(slices.Concat([]byte{0x0a, 0x01, 'm', 0x11}, fixed(nan))), lm.ErrNotFinite},
		{"a level that is +Inf", state(slices.Concat([]byte{0x0a, 0x01, 'm', 0x11}, fixed(0x7FF0000000000000))), lm.ErrNotFinite},
		{"a double cut short", state([]byte{0x0a, 0x01, 'm', 0x11, 0x00}), lm.ErrCorrupt},
		{"a fixed32 cut short", []byte{0x08, 0x01, 0xfd, 0x06, 0x00}, lm.ErrCorrupt},
	} {
		_, err := lm.DecodeState(c.b)
		if !errors.Is(err, c.err) {
			t.Errorf("%s: %v, not %v", c.name, err, c.err)
		}
	}
	for _, c := range []struct {
		name string
		b    []byte
	}{
		{"a bool of 2", []byte{0x18, 0x02}},
		{"an unpacked enum list", slices.Concat(evidenceTop(), []byte{0x92, 0x01, 0x04, 0x08, 0x01, 0x38, 0x01})},
		{"a packed list cut short", slices.Concat(evidenceTop(), []byte{0x92, 0x01, 0x04, 0x08, 0x01, 0x3a, 0x05})},
		{"a packed value past int32", slices.Concat(evidenceTop(), []byte{0x92, 0x01, 0x09, 0x08, 0x01, 0x3a, 0x05, 0x80, 0x80, 0x80, 0x80, 0x10})},
	} {
		if _, err := lm.DecodeEvidence(c.b); !errors.Is(err, lm.ErrCorrupt) {
			t.Errorf("%s: %v, not ErrCorrupt", c.name, err)
		}
	}
}

// evidenceTop gives the required top-level fields of the golden evidence, without its items and
// traps, so that a test can append an item of its own.
func evidenceTop() []byte {
	e := goldenEvidence()
	e.Task.Items, e.Items, e.Traps = nil, nil, nil
	b, err := lm.EncodeEvidence(e)
	if err != nil {
		panic(err)
	}
	return b
}

// A field this version does not know, of any wire type, is skipped, and an enum value it does
// not know is kept as its number.
func TestUnknown(t *testing.T) {
	b, _ := lm.EncodeEvidence(goldenEvidence())
	want := normalEvidence(goldenEvidence())
	for _, extra := range [][]byte{
		{0xb8, 0x06, 0x2a},                   // 103, a varint
		{0xb9, 0x06, 1, 2, 3, 4, 5, 6, 7, 8}, // 103, a fixed64
		{0xba, 0x06, 0x02, 'h', 'i'},         // 103, length-delimited
		{0xbd, 0x06, 1, 2, 3, 4},             // 103, a fixed32
	} {
		e, err := lm.DecodeEvidence(slices.Concat(b, extra))
		if err != nil || !equalEvidence(e, want) {
			t.Errorf("an unknown field %x: %v", extra, err)
		}
	}
	e := goldenEvidence()
	e.Outcome.Kind = 42
	e.Task.Items[0].Slots[0].Check = 17
	e.Items[0].FirstVerdict = 9
	d := canonicalEvidence(e)
	if d.Outcome.Kind != 42 || d.Task.Items[1].Slots[0].Check != 17 || d.Items[0].FirstVerdict != 9 {
		t.Errorf("unknown enum values not kept: %+v", d)
	}
}

func TestEncodeRefuses(t *testing.T) {
	for name, e := range map[string]func(*lm.Evidence){
		"a primary skill twice":    func(e *lm.Evidence) { e.Task.Primary = append(e.Task.Primary, e.Task.Primary[0]) },
		"a supporting skill twice": func(e *lm.Evidence) { e.Task.Supporting = append(e.Task.Supporting, e.Task.Supporting[0]) },
		"an item skill twice": func(e *lm.Evidence) {
			e.Task.Items[0].Primary = append(e.Task.Items[0].Primary, e.Task.Items[0].Primary[0])
		},
		"an item supporting skill twice": func(e *lm.Evidence) {
			e.Task.Items[0].Supporting = []string{"core.a", "core.a"}
		},
		"an item place twice":     func(e *lm.Evidence) { e.Task.Items[1].Place = 2; e.Items[0].Place = 2 },
		"places that differ":      func(e *lm.Evidence) { e.Items[0].Place = 9 },
		"a trap twice":            func(e *lm.Evidence) { e.Traps = append(e.Traps, e.Traps[0]) },
		"a trap skill twice":      func(e *lm.Evidence) { e.Traps[0].Skills = []string{"x", "x"} },
		"a trap item place twice": func(e *lm.Evidence) { e.Traps[0].ItemPlaces = []int32{1, 1} },
	} {
		ev := goldenEvidence()
		e(&ev)
		if _, err := lm.EncodeEvidence(ev); !errors.Is(err, lm.ErrInvalid) {
			t.Errorf("%s: %v, not ErrInvalid", name, err)
		}
	}
	nan := math.Float64frombits(0x7FF8000000000001)
	for name, c := range map[string]struct {
		change func(*lm.State)
		err    error
	}{
		"a NaN level":         {func(s *lm.State) { s.Modules[0].Level = nan }, lm.ErrNotFinite},
		"an infinite weight":  {func(s *lm.State) { s.Modules[0].Trial[0].Weight = math.Inf(-1) }, lm.ErrNotFinite},
		"a NaN stability":     {func(s *lm.State) { s.Modules[0].Skills[0].Stability = nan }, lm.ErrNotFinite},
		"a module twice":      {func(s *lm.State) { s.Modules = append(s.Modules, s.Modules[0]) }, lm.ErrInvalid},
		"a skill twice":       {func(s *lm.State) { m := &s.Modules[0]; m.Skills = append(m.Skills, m.Skills[0]) }, lm.ErrInvalid},
		"a module trap twice": {func(s *lm.State) { m := &s.Modules[0]; m.Traps = append(m.Traps, m.Traps[0]) }, lm.ErrInvalid},
		"a skill trap twice": {func(s *lm.State) {
			sk := &s.Modules[0].Skills[0]
			sk.Traps = append(sk.Traps, sk.Traps[0])
		}, lm.ErrInvalid},
	} {
		s := goldenState()
		c.change(&s)
		if _, err := lm.EncodeState(s); !errors.Is(err, c.err) {
			t.Errorf("%s: %v, not %v", name, err, c.err)
		}
	}
	if err := (func() error { _, err := lm.EncodeState(goldenState()); return err })(); err != nil {
		t.Fatal(err)
	}
	if _, err := lm.DecodeEvidence([]byte{0x08}); err == nil || !strings.Contains(err.Error(), "evidence") {
		t.Errorf("an error that does not say what it read: %v", err)
	}
}

// Every known field of every message, at every depth of the golden bytes, is refused written
// twice and refused of another wire type: the bytes of each message are rebuilt with the field
// repeated, or with one more field of its number and a wrong type, inside their unchanged parents.
func TestEveryFieldStrict(t *testing.T) {
	e, _ := lm.EncodeEvidence(goldenEvidence())
	s, _ := lm.EncodeState(goldenState())
	// The repeated fields whose elements keep the order they came in, by the path of their numbers
	// from the top: an element of them may come twice. Every other repeated field is in canonical
	// order, each element once, and an element twice is refused like a field twice.
	evidenceRepeated := []string{"18.5", "18.7"}
	stateRepeated := []string{"2.5", "2.5.3"}
	for _, c := range []struct {
		name     string
		b        []byte
		read     func([]byte) error
		repeated []string
	}{
		{"evidence", e, func(b []byte) error { _, err := lm.DecodeEvidence(b); return err }, evidenceRepeated},
		{"state", s, func(b []byte) error { _, err := lm.DecodeState(b); return err }, stateRepeated},
	} {
		for _, m := range mutations(c.b, "", nil) {
			err := c.read(m.bytes)
			wantErr := m.wrongType || !slices.Contains(c.repeated, m.path)
			if wantErr && !errors.Is(err, lm.ErrCorrupt) && !errors.Is(err, lm.ErrNotFinite) {
				t.Errorf("%s, field %s, %s: %v, not refused", c.name, m.path, m.what, err)
			}
			if !wantErr && err != nil {
				t.Errorf("%s, field %s, %s: %v", c.name, m.path, m.what, err)
			}
		}
	}
}

type mutation struct {
	path, what string
	wrongType  bool
	bytes      []byte
}

// field is one field of a message's bytes: its number, wire type, the bytes of its tag and value.
type rawField struct {
	number, wireType int
	all, value       []byte
}

func splitFields(b []byte) []rawField {
	var fields []rawField
	for off := 0; off < len(b); {
		start := off
		tag, n := uvarint(b[off:])
		off += n
		f := rawField{number: int(tag >> 3), wireType: int(tag & 7)}
		valueStart := off
		switch f.wireType {
		case 0:
			_, n = uvarint(b[off:])
			off += n
		case 1:
			off += 8
		case 2:
			l, n := uvarint(b[off:])
			valueStart = off + n
			off += n + int(l)
		}
		f.all, f.value = b[start:off], b[valueStart:off]
		fields = append(fields, f)
	}
	return fields
}

func appendTag(b []byte, number, wireType int) []byte {
	v := uint64(number)<<3 | uint64(wireType)
	for v >= 0x80 {
		b = append(b, byte(v)|0x80)
		v >>= 7
	}
	return append(b, byte(v))
}

// mutations gives, for each field of b at any depth, b with the field written twice and b with a
// field of its number and another wire type appended to its message; wrap rebuilds the parents.
func mutations(b []byte, prefix string, wrap func([]byte) []byte) []mutation {
	if wrap == nil {
		wrap = func(b []byte) []byte { return b }
	}
	var out []mutation
	fields := splitFields(b)
	for i, f := range fields {
		path := strconv.Itoa(f.number)
		if prefix != "" {
			path = prefix + "." + path
		}
		out = append(out, mutation{path, "twice", false, wrap(slices.Concat(b, f.all))})
		for _, wt := range []int{0, 1, 2} {
			if wt == f.wireType {
				continue
			}
			bad := appendTag(slices.Clone(b), f.number, wt)
			switch wt {
			case 0:
				bad = append(bad, 0x01)
			case 1:
				bad = append(bad, 0, 0, 0, 0, 0, 0, 0, 0)
			case 2:
				bad = append(bad, 0x00)
			}
			out = append(out, mutation{path, "wire type " + strconv.Itoa(wt), true, wrap(bad)})
		}
		if f.wireType == 2 && looksLikeMessage(f.value) {
			head := b[:len(b)-len(slices.Concat(restAll(fields[i:])...))]
			tail := b[len(head)+len(f.all):]
			number := f.number
			inner := func(v []byte) []byte {
				nb := appendTag(slices.Clone(head), number, 2)
				nb = appendLen(nb, len(v))
				nb = append(nb, v...)
				return wrap(append(nb, tail...))
			}
			out = append(out, mutations(f.value, path, inner)...)
		}
	}
	return out
}

func restAll(fields []rawField) [][]byte {
	var all [][]byte
	for _, f := range fields {
		all = append(all, f.all)
	}
	return all
}

func appendLen(b []byte, n int) []byte {
	v := uint64(n)
	for v >= 0x80 {
		b = append(b, byte(v)|0x80)
		v >>= 7
	}
	return append(b, byte(v))
}

// looksLikeMessage tells a nested message from a string or a packed list: its bytes split into
// fields with numbers from 1 to 20 and known wire types, exactly to their end.
func looksLikeMessage(b []byte) bool {
	if len(b) == 0 {
		return false
	}
	for off := 0; off < len(b); {
		tag, n := uvarint(b[off:])
		if n == 0 || tag>>3 == 0 || tag>>3 > 20 {
			return false
		}
		off += n
		switch tag & 7 {
		case 0:
			_, n = uvarint(b[off:])
			off += n
		case 1:
			off += 8
		case 2:
			l, n := uvarint(b[off:])
			off += n + int(l)
		default:
			return false
		}
		if off > len(b) {
			return false
		}
	}
	return true
}

// What the writer refuses beyond duplicates, and what the reader refuses as no bytes the writer
// would have given.
func TestStrictValues(t *testing.T) {
	for name, change := range map[string]func(*lm.Evidence){
		"kind 0":                     func(e *lm.Evidence) { e.Kind = 0 },
		"a kind past this version's": func(e *lm.Evidence) { e.Kind = lm.EvidenceKind + 1 },
		"a task ID of no UTF-8":      func(e *lm.Evidence) { e.Task.ID = "x\xff" },
		"a skill of no UTF-8":        func(e *lm.Evidence) { e.Task.Primary = []string{"\xc3"} },
		"a rung without its flag":    func(e *lm.Evidence) { e.Outcome.HasRung = false },
		"a fact without its flag":    func(e *lm.Evidence) { e.SolutionSeenBefore, e.HasSolutionSeenFact = true, false },
		"options without their flag": func(e *lm.Evidence) { e.Task.Items[0].Slots[0].HasOptions = false },
		"a first verdict without its flag": func(e *lm.Evidence) {
			e.Items[1].HasFirstVerdict = false
		},
		"wrong answers without a first verdict": func(e *lm.Evidence) {
			e.Items[0].FirstVerdict, e.Items[0].HasFirstVerdict, e.Items[0].WrongBeforeFirst = 0, false, 3
		},
		"a final verdict without its flag": func(e *lm.Evidence) { e.Items[0].HasFinalVerdict = false },
		"a first correct without its flag": func(e *lm.Evidence) { e.Items[0].HasFirstCorrect = false },
	} {
		e := goldenEvidence()
		change(&e)
		if _, err := lm.EncodeEvidence(e); !errors.Is(err, lm.ErrInvalid) {
			t.Errorf("evidence with %s: %v, not ErrInvalid", name, err)
		}
	}
	for name, change := range map[string]func(*lm.State){
		"kind 0":                        func(s *lm.State) { s.Kind = 0 },
		"a kind past this version's":    func(s *lm.State) { s.Kind = lm.StateKind + 1 },
		"a module of no UTF-8":          func(s *lm.State) { s.Modules[0].Module = "\xff" },
		"a mastery band without a flag": func(s *lm.State) { s.Modules[0].Skills[1].Mastered = false },
	} {
		s := goldenState()
		change(&s)
		if _, err := lm.EncodeState(s); !errors.Is(err, lm.ErrInvalid) {
			t.Errorf("a state with %s: %v, not ErrInvalid", name, err)
		}
	}
	for name, b := range map[string][]byte{
		"skills out of order":  slices.Concat(evidenceTop(), []byte{0x3a, 0x01, 'a'}),
		"a skill of no UTF-8":  slices.Concat(evidenceTop(), []byte{0x3a, 0x01, 0xff}),
		"field 9 without 8":    slices.Concat(evidenceTop(), []byte{0x92, 0x01, 0x08, 0x08, 0x01, 0x10, 0x00, 0x30, 0x00, 0x48, 0x01}),
		"kind 0":               slices.Concat([]byte{0x08, 0x00}, evidenceTop()[2:]),
		"items out of order":   slices.Concat(evidenceTop(), item(2), item(1)),
		"an item place twice":  slices.Concat(evidenceTop(), item(1), item(1)),
		"two traps, unordered": slices.Concat(evidenceTop(), []byte{0x9a, 0x01, 0x03, 0x0a, 0x01, 'b', 0x9a, 0x01, 0x03, 0x0a, 0x01, 'a'}),
	} {
		if _, err := lm.DecodeEvidence(b); !errors.Is(err, lm.ErrCorrupt) {
			t.Errorf("evidence with %s: %v, not ErrCorrupt", name, err)
		}
	}
	if _, err := lm.DecodeEvidence(slices.Concat(evidenceTop(), item(1), item(2))); err != nil {
		t.Errorf("two items in order: %v", err)
	}
	if _, err := lm.DecodeState([]byte{0x08, 0x00}); !errors.Is(err, lm.ErrCorrupt) {
		t.Errorf("a state of kind 0: %v, not ErrCorrupt", err)
	}
	if _, err := lm.DecodeState([]byte{0x08, 0x05}); err != nil {
		t.Errorf("a state of a later kind: %v", err)
	}
}

// item gives the bytes of field 18 with the least an item carries: its place, no own skills, not
// submitted.
func item(place byte) []byte {
	return []byte{0x92, 0x01, 0x06, 0x08, place, 0x10, 0x00, 0x30, 0x00}
}
