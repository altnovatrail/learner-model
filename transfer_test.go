package learnermodel_test

import (
	"bytes"
	"fmt"
	"slices"
	"testing"

	lm "github.com/altnovatrail/learner-model"
)

// The transfer by retired_skills (docs/architecture.md, section 3.5): the evidence of a retired
// skill counts for each skill that replaced it, in full, by the graph of the module's current
// version, while the log stays as it is.

// versionGraphs gives the graphs of a version of the workbook: skills with a verdict, each taught
// at band 1 but those in untaught; essay and the skills in rubric, with a rubric; and the retired
// skills given.
func versionGraphs(t *testing.T, verdict, untaught, rubric []string, retired ...lm.RetiredSkill) lm.Graphs {
	t.Helper()
	bands := []lm.Band{{Place: 1, ID: "grade-2", Grades: []int32{2}}, {Place: 2, ID: "grade-3", Grades: []int32{3}}}
	var skills []lm.Skill
	for _, id := range append([]string{skill("essay")}, rubric...) {
		skills = append(skills, lm.Skill{ID: id, Assessability: lm.AssessabilityRubric})
	}
	var tasks []lm.Task
	for _, id := range verdict {
		skills = append(skills, lm.Skill{ID: id, Assessability: lm.AssessabilityVerdict})
		if !slices.Contains(untaught, id) {
			tasks = append(tasks, lm.Task{ID: workbook + "/task/g", Band: lm.BandRef{Place: 1, ID: "grade-2"}, Difficulty: 3, Primary: []string{id}})
		}
	}
	g, err := lm.NewGraph(workbook, bands, skills, retired, tasks)
	if err != nil {
		t.Fatal(err)
	}
	return lm.Graphs{workbook: {Graph: g, Current: true}}
}

// transferLog is a log of the workbook, a page after a page, each naming the primary skills
// given. The first two pages are all right, 8 and 13 items, so that the trial series ends and the
// skills of both are mastered; the later ones have a wrong answer or two, come at growing pauses,
// so that forgetting shows, and every second one has a trap at the page's skills.
func transferLog(pages ...[]string) []lm.Evidence {
	var log []lm.Evidence
	at := int64(1762041600)
	for i, primary := range pages {
		items, right := 8, 8-i%3
		if i == 1 {
			items, right = 13, 13
		}
		e := hard(page(items, right, primary, at))
		if i%2 == 0 {
			e.Traps = []lm.FoundTrap{{ID: workbook + "/trap/borrow", Skills: primary, ItemPlaces: []int32{1}}}
		}
		log = append(log, e)
		at += int64(i+1) * 86400
	}
	return log
}

// renamed is the log as it would be had its pages named the skills of the current version: each
// primary skill and each skill of a trap goes to the skills names gives it, each once, and a skill
// names leaves out stays as it is.
func renamed(log []lm.Evidence, names map[string][]string) []lm.Evidence {
	rename := func(ids []string) []string {
		var out []string
		for _, id := range ids {
			to, ok := names[id]
			if !ok {
				to = []string{id}
			}
			for _, n := range to {
				if !slices.Contains(out, n) {
					out = append(out, n)
				}
			}
		}
		return out
	}
	var out []lm.Evidence
	for _, e := range log {
		e.Task.Primary = rename(e.Task.Primary)
		e.Traps = slices.Clone(e.Traps)
		for i := range e.Traps {
			e.Traps[i].Skills = rename(e.Traps[i].Skills)
		}
		out = append(out, e)
	}
	return out
}

func recompute(t *testing.T, log []lm.Evidence, g lm.Graphs) lm.State {
	t.Helper()
	s, err := lm.Recompute(log, g)
	if err != nil {
		t.Fatal(err)
	}
	return s
}

// The vectors of the transfer: recomputed by the graph of the current version, a log that names
// the skills of earlier versions gives the bytes of the log as if it had named their replacements
// — a renamed skill, skills merged in different items and in one, where the merged skill gets
// 2/J of the step and counts the item once, a chain of retirements, which the graph of the current
// version names in one step, and skills that go to none: retired with `[]` or with a replacement
// with a rubric, a skill that now has a rubric, and one the graph does not know, whose items move
// only θ and n and whose traps are the module's.
func TestTransfer(t *testing.T) {
	add, sub, count := skill("addition"), skill("subtraction"), skill("counting")
	addition2, carry, borrow, unknown := skill("addition2"), skill("carry"), skill("borrow"), skill("unknown")
	merged := versionGraphs(t, []string{count}, nil, nil,
		lm.RetiredSkill{ID: add, ReplacedBy: []string{count}}, lm.RetiredSkill{ID: sub, ReplacedBy: []string{count}})
	for _, c := range []struct {
		name   string
		g      lm.Graphs
		log    []lm.Evidence
		names  map[string][]string
		skills []string // the skills the state has
	}{
		{"a renamed skill", versionGraphs(t, []string{addition2, sub}, nil, nil, lm.RetiredSkill{ID: add, ReplacedBy: []string{addition2}}),
			transferLog([]string{add}, []string{add, sub}, []string{sub}, []string{add}, []string{add, sub}),
			map[string][]string{add: {addition2}}, []string{addition2, sub}},
		{"skills merged in different items and in one", merged,
			transferLog([]string{add}, []string{sub}, []string{add, sub}, []string{sub}, []string{add, sub}),
			map[string][]string{add: {count}, sub: {count}}, []string{count}},
		// Version 2 renamed addition to carry, and version 3 carry to counting: the graph of version
		// 3 names counting for both.
		{"a chain of retirements", versionGraphs(t, []string{count}, nil, nil,
			lm.RetiredSkill{ID: add, ReplacedBy: []string{count}}, lm.RetiredSkill{ID: carry, ReplacedBy: []string{count}}),
			transferLog([]string{add}, []string{add}, []string{carry}, []string{carry}, []string{count}),
			map[string][]string{add: {count}, carry: {count}}, []string{count}},
		{"skills that go to none", versionGraphs(t, []string{sub}, nil, []string{add},
			lm.RetiredSkill{ID: carry}, lm.RetiredSkill{ID: borrow, ReplacedBy: []string{skill("essay")}}),
			transferLog([]string{add}, []string{carry, sub}, []string{borrow}, []string{unknown}, []string{sub, add}, []string{carry}),
			map[string][]string{add: nil, carry: nil, borrow: nil, unknown: nil}, []string{sub}},
	} {
		got := recompute(t, c.log, c.g)
		want := recompute(t, renamed(c.log, c.names), c.g)
		if !bytes.Equal(encodeState(t, got), encodeState(t, want)) {
			t.Errorf("%s: the log gives\n%+v\nand the log of the current skills\n%+v", c.name, got, want)
		}
		var skills []string
		for _, sk := range got.Modules[0].Skills {
			skills = append(skills, sk.Skill)
		}
		if !slices.Equal(skills, c.skills) || !(got.Modules[0].AnswerWeight > 5) {
			t.Errorf("%s: skills %v, n %v; want skills %v and n past the trial series", c.name, skills, got.Modules[0].AnswerWeight, c.skills)
		}
	}
}

// A split (docs/architecture.md, section 3.5): each part gets the whole history of the skill it
// came from, beside another primary skill, and θ, n and the other skill move as before the split,
// bit for bit, since the mean of two equal corrections is the correction. The mastery comes anew,
// by the bands each part is taught at now: a part no task trains is mastered never.
func TestTransferSplit(t *testing.T) {
	add, sub, carry, place := skill("addition"), skill("subtraction"), skill("carry"), skill("place")
	log := transferLog([]string{add, sub}, []string{add, sub}, []string{add}, []string{sub}, []string{add, sub}, []string{add})
	before := recompute(t, log, versionGraphs(t, []string{add, sub}, nil, nil))
	after := recompute(t, log, versionGraphs(t, []string{carry, place, sub}, []string{place}, nil,
		lm.RetiredSkill{ID: add, ReplacedBy: []string{carry, place}}))
	m0, m1 := before.Modules[0], after.Modules[0]
	if m0.Level != m1.Level || m0.AnswerWeight != m1.AnswerWeight || !same(m0.Traps, m1.Traps) {
		t.Errorf("the split moved the module: θ %v → %v, n %v → %v", m0.Level, m1.Level, m0.AnswerWeight, m1.AnswerWeight)
	}
	a := skillOf(t, before, add)
	if !a.Mastered || len(a.Traps) == 0 {
		t.Fatalf("the log masters no addition and finds no trap of it: %+v", a)
	}
	for _, part := range []string{carry, place} {
		got, want := skillOf(t, after, part), a
		got.Skill = add
		if part == place {
			want.Mastered, want.MasteryBand = false, 0
		}
		if !same(got, want) {
			t.Errorf("part %s: %+v, want %+v", part, got, want)
		}
	}
	if got, want := skillOf(t, after, sub), skillOf(t, before, sub); !same(got, want) {
		t.Errorf("subtraction beside the split: %+v, want %+v", got, want)
	}
}

// The thresholds of the sums of weights on a log with the transfer: a page of 8 items of a skill
// since retired and one of 13 of its replacement make m = 5 within ε, and master it.
func TestWeightThresholdsTransferred(t *testing.T) {
	add, count := skill("addition"), skill("counting")
	g := versionGraphs(t, []string{count}, nil, nil, lm.RetiredSkill{ID: add, ReplacedBy: []string{count}})
	s := mustAppend(t, lm.State{}, hard(page(8, 8, []string{add}, 1)), g)
	s = mustAppend(t, s, hard(page(13, 13, []string{count}, 2)), g)
	sk := skillOf(t, s, count)
	if len(s.Modules[0].Skills) != 1 || !(sk.AnswerWeight > 5-1e-9 && sk.AnswerWeight < 5+1e-9) || !sk.Mastered || sk.MasteryBand != 1 {
		t.Errorf("counting after the pages: m %v, mastered %v at %d; want m = 5, mastered at band 1", sk.AnswerWeight, sk.Mastered, sk.MasteryBand)
	}
}

// Anchor skills among the supporting skills go as they are, and the model reads none of them
// (docs/architecture.md, section 3.5): a log whose pages name some gives the bytes of the log
// without them.
func TestAnchorSkills(t *testing.T) {
	g := workbookGraphs(t)
	log := transferLog([]string{skill("addition")}, []string{skill("addition"), skill("subtraction")}, []string{skill("subtraction")})
	var anchored []lm.Evidence
	for _, e := range log {
		e.Task.Supporting = []string{"core.counting", "core.place-value"}
		anchored = append(anchored, e)
	}
	if !bytes.Equal(encodeState(t, recompute(t, anchored, g)), encodeState(t, recompute(t, log, g))) {
		t.Errorf("anchor skills among the supporting skills moved the state")
	}
}

// Evidence of a kind later than this version's, as its bytes come from the log
// (docs/architecture.md, section 3.5): a field this version does not know is skipped, an item of
// an unknown kind of check or an unknown verdict gives no answer, and an unknown outcome reads as
// "unchecked". The state is the one the same evidence gives without what this version does not
// know.
func TestReadKinds(t *testing.T) {
	g := workbookGraphs(t)
	e := page(4, 2, []string{skill("addition")}, 200)
	e.Outcome = lm.Outcome{Kind: 42}
	e.Task.Items[1].Slots[0].Check = 17
	e.Items[2].FirstVerdict = 9
	b, err := lm.EncodeEvidence(e)
	if err != nil {
		t.Fatal(err)
	}
	// The kind is field 1, written first: 0x08 and its value.
	if b[0] != 0x08 || b[1] != byte(lm.EvidenceKind) {
		t.Fatalf("the evidence opens with %x, not its kind", b[:2])
	}
	b[1] = byte(lm.EvidenceKind + 1)
	b = append(b, 0xba, 0x06, 0x02, 'h', 'i') // field 103, length-delimited
	later, err := lm.DecodeEvidence(b)
	if err != nil || later.Kind != lm.EvidenceKind+1 {
		t.Fatalf("evidence of a later kind: kind %d, %v", later.Kind, err)
	}
	got, _, err := lm.Append(ready(t), later, g)
	if err != nil {
		t.Fatal(err)
	}
	known := page(4, 2, []string{skill("addition")}, 200)
	known.Outcome = lm.Outcome{Kind: lm.OutcomeUnchecked}
	known.Task.Items = []lm.TaskItem{known.Task.Items[0], known.Task.Items[3]}
	known.Items = []lm.ItemAnswer{known.Items[0], known.Items[3]}
	if want := mustAppend(t, ready(t), known, g); !bytes.Equal(encodeState(t, got), encodeState(t, want)) {
		t.Errorf("evidence of a later kind gives\n%+v\nwant\n%+v", got, want)
	}
}

// same compares two values by what %+v prints of them, which prints every float exactly.
func same(a, b any) bool { return fmt.Sprintf("%+v", a) == fmt.Sprintf("%+v", b) }
