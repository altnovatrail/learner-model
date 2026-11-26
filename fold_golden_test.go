package learnermodel_test

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"testing"

	lm "github.com/altnovatrail/learner-model"
)

// The golden vectors of the fold: synthetic logs, each piece of evidence with the bytes of
// the state after it, one file for every build, testdata/fold. A line is "evidence <hex>" or "state
// <hex>", the bytes of the codec; a line opening with # is a note. The test reads the evidence from
// the file, appends it and holds the state's bytes to the next line, so that the file pins what
// goes in as well as what comes out. LEARNERMODEL_UPDATE=1 writes the files anew from the logs
// below, and they change only with the version of learner-model.

// The times of the golden logs: foldLogStart is when each log begins, a synthetic time in seconds,
// and foldAnswerGap the hour between two answers of standalone's tasks; a page of the workbook
// comes a day after the one before it.
const (
	foldLogStart  = int64(epoch)
	foldAnswerGap = 3600
)

// The answers of the log break: before the break, past 20, and after it; n = 20 + foldBreakAnswersAfter
// at the end.
const (
	foldBreakAnswersBefore = 24
	foldBreakAnswersAfter  = 4
)

// ownKinds is the kind of a golden journal whose log names the kind of each piece itself.
const ownKinds = 0

// goldenJournal is a log of a golden file with the graphs it is folded against, and the kind its
// evidence is written in: 1 for the logs that kind 1 wrote, kept as api keeps them, and ownKinds
// for a log that names the kind of each piece itself.
type goldenJournal struct {
	name   string
	note   string
	graphs func(*testing.T) lm.Graphs
	log    func() []lm.Evidence
	kind   int32
}

func mergeGraphs(gs ...lm.Graphs) lm.Graphs {
	out := lm.Graphs{}
	for _, g := range gs {
		for k, v := range g {
			out[k] = v
		}
	}
	return out
}

func foldGraphs(t *testing.T) lm.Graphs {
	return mergeGraphs(
		olympiadGraphs(t, map[string][]int32{topic("a"): {1, 2}, topic("b"): {2, 3}}),
		workbookGraphs(t),
		lm.Graphs{"example/waiting": {Current: false}},
	)
}

// olympiadLog is a learner of grade 3 in olympiad maths: the trial series, a mastery, a slip that
// holds it, a second that takes it away, and after a break of 60 days a mastery again.
func olympiadLog() []lm.Evidence {
	const (
		gradeBand   = 2        // the band 3-4, a learner of grade 3's
		masteryRun  = 6        // the right answers in a row that master topic a
		summerBreak = 60 * day // retention falls, and the next right answers raise τ much
		answersBack = 4        // the right answers after the break
	)
	start := ladderBand(gradeBand)
	var log []lm.Evidence
	at := foldLogStart
	add := func(c choiceSpec) {
		at += foldAnswerGap
		c.start, c.at = start, at
		log = append(log, choiceAnswer(c))
	}
	add(choiceSpec{skill: topic("a"), place: 2, difficulty: 3, correct: true})
	add(choiceSpec{skill: topic("b"), place: 2, difficulty: 4, correct: false})
	add(choiceSpec{skill: topic("a"), place: 2, difficulty: 4, correct: true})
	add(choiceSpec{skill: topic("b"), place: 3, difficulty: 1, correct: true})
	add(choiceSpec{skill: topic("a"), place: 1, difficulty: 5, correct: true})
	for range masteryRun {
		add(choiceSpec{skill: topic("a"), place: 2, difficulty: 3, correct: true})
	}
	add(choiceSpec{skill: topic("a"), place: 2, difficulty: 4, correct: false})
	log[len(log)-1].Traps = []lm.FoundTrap{{ID: olympiad + "/trap/order", Skills: []string{topic("a")}, ItemPlaces: []int32{1}}}
	add(choiceSpec{skill: topic("a"), place: 2, difficulty: 5, correct: false})
	at += summerBreak
	for range answersBack {
		add(choiceSpec{skill: topic("a"), place: 2, difficulty: 2, correct: true})
	}
	add(choiceSpec{skill: topic("b"), place: 3, difficulty: 2, correct: false})
	return log
}

// workbookLog is pages of the workbook: two primary skills on every item, J = 2, an item that
// names its own skill, one with a rubric, items with no answer and a failed page; a page after a
// hint, one with the solution shown, one with the technique on request, and traps found.
func workbookLog() []lm.Evidence {
	// The pages with something of their own, by their index in the log.
	const (
		afterHintPage     = 2
		solutionShownPage = 5
		techniquePage     = 7
	)
	both := []string{skill("addition"), skill("subtraction")}
	var log []lm.Evidence
	for i, p := range []pageSpec{
		{items: 10, right: 8}, {items: 10, right: 9}, {items: 6, right: 2}, {items: 12, right: 12},
		{items: 10, right: 7}, {items: 10, right: 3}, {items: 10, right: 3}, {items: 20, right: 20},
	} {
		p.primary, p.at = both, foldLogStart+int64(day*i)
		e := page(p)
		e.Task.Items[0].OwnSkills, e.Task.Items[0].Primary = true, []string{skill("addition")}
		e.Task.Items[1].OwnSkills, e.Task.Items[1].Primary = true, []string{skill("essay")}
		e.Items[p.items-1].HasFirstVerdict, e.Items[p.items-1].FirstVerdict, e.Items[p.items-1].Submitted = false, 0, false
		switch i {
		case afterHintPage: // the items first answered after the rung are with help
			e.Outcome = lm.Outcome{Kind: lm.OutcomeAfterHint, Rung: 1, HasRung: true}
			// The first wrong answer, the item after the right ones, came before them.
			for k := p.right + 1; k < p.items-1; k++ {
				e.Items[k].WrongBeforeFirst = 1
			}
		case solutionShownPage: // the item never submitted is not solved
			e.Outcome = lm.Outcome{Kind: lm.OutcomeSolutionShown}
		case techniquePage: // a page all right that masters nothing
			e.TechniqueOnRequest = true
		}
		if p.right < p.items-1 {
			e.Traps = []lm.FoundTrap{
				{ID: workbook + "/trap/borrow", Skills: []string{skill("subtraction")}, ItemPlaces: []int32{int32(p.items - 1)}},
				{ID: workbook + "/trap/loose", ItemPlaces: []int32{int32(p.right + 1)}},
			}
		}
		log = append(log, e)
	}
	return log
}

// mixedLog interleaves the two, with evidence of a module that has no current version, evidence
// with no answer and a task whose skill tags the AI gave, which weighs nothing.
func mixedLog() []lm.Evidence {
	o, w := olympiadLog(), workbookLog()
	waiting := choiceAnswer(choiceSpec{skill: topic("a"), place: 1, difficulty: 3, correct: true, start: ladderBand(1), at: foldLogStart})
	waiting.Task.ID = "example/waiting/task/x"
	ai := page(pageSpec{items: 4, right: 4, primary: []string{skill("addition")}, at: foldLogStart + 2})
	ai.Task.SkillTagsProvenance = lm.ProvenanceAI // a task from paper the AI tagged: it weighs nothing
	none := page(pageSpec{items: 3, right: 3, primary: []string{skill("addition")}, at: foldLogStart + 1})
	for i := range none.Items {
		none.Items[i].HasFirstVerdict, none.Items[i].FirstVerdict = false, 0
	}
	log := []lm.Evidence{waiting, none}
	for i := range max(len(o), len(w)) {
		if i < len(w) {
			log = append(log, w[i])
		}
		if i < len(o) {
			log = append(log, o[i])
		}
	}
	return append(log, ai, waiting)
}

// transferGraphs is the workbook at its third version: version 2 split addition into carry and
// place, and version 3 merged carry into counting and renamed subtraction, retired tally with no
// replacement and still teaches place in no task. The graph names each retired skill's
// replacements in this version, so addition goes to counting and place in one step.
func transferGraphs(t *testing.T) lm.Graphs {
	return versionGraphs(t, []string{skill("counting"), skill("place"), skill("subtraction2")}, []string{skill("place")}, nil,
		lm.RetiredSkill{ID: skill("addition"), ReplacedBy: []string{skill("counting"), skill("place")}},
		lm.RetiredSkill{ID: skill("carry"), ReplacedBy: []string{skill("counting")}},
		lm.RetiredSkill{ID: skill("subtraction"), ReplacedBy: []string{skill("subtraction2")}},
		lm.RetiredSkill{ID: skill("tally")},
	)
}

// transferLog3 is a log of the workbook across its three versions, recomputed by the third: pages
// of version 1, of version 2 with the parts of addition side by side, and of version 3, with a
// skill retired without a replacement and one the graph does not know.
func transferLog3() []lm.Evidence {
	s := skill
	return transferLog(
		[]string{s("addition")}, []string{s("addition"), s("subtraction")}, []string{s("subtraction")},
		[]string{s("carry"), s("place")}, []string{s("carry")}, []string{s("tally")},
		[]string{s("counting"), s("subtraction2")}, []string{s("unknown"), s("counting")}, []string{s("place")},
	)
}

var goldenJournals = []goldenJournal{
	{name: "olympiad", note: "a learner of grade 3 in olympiad maths: the trial series, mastery held, lost and set again after a break of 60 days", graphs: foldGraphs, log: olympiadLog, kind: lm.FirstKind},
	{name: "workbook", note: "pages of the workbook with two primary skills, J = 2, a hint, a shown solution, help and traps", graphs: foldGraphs, log: workbookLog, kind: lm.FirstKind},
	{name: "mixed", note: "the two interleaved, with a module without a current version, evidence with no answer and a task tagged by AI", graphs: foldGraphs, log: mixedLog, kind: lm.FirstKind},
	{name: "transfer", note: "the workbook across three versions, by the third: a split, a merge, a rename, a skill retired with [] and one unknown", graphs: transferGraphs, log: transferLog3, kind: lm.FirstKind},
	{name: "format04", note: "pages of kind 1, then of kind 2 in format 0.4: checks of the kind time, a trap of the 12-hour clock, and the leading question on request, which masters nothing", graphs: foldGraphs, log: format04Log, kind: ownKinds},
	{name: "break", note: "a learner of grade 3 in olympiad maths past 20 answers, back after a break of 30 days: the level forgets, and n falls back to 20", graphs: foldGraphs, log: breakLog, kind: ownKinds},
}

// breakLog is a learner of grade 3 in olympiad maths with 24 answers over two skills, who comes back
// after a break of 30 days: the pause of the module brings n back to 20 before the first answer
// after it, and the four answers after the break take it to 24, not 28.
func breakLog() []lm.Evidence {
	const (
		gradeBand     = 2        // the band 3-4, a learner of grade 3's
		lowDifficulty = 2        // the answers before it run through the difficulties 2, 3 and 4
		difficulties  = 3        // in turn
		wrongEvery    = 4        // and every fourth of them is wrong
		breakPause    = 30 * day // the pause of the module, past 14 days
	)
	topics := []string{"a", "b"}
	start := ladderBand(gradeBand)
	var log []lm.Evidence
	at := foldLogStart
	add := func(c choiceSpec) {
		at += foldAnswerGap
		c.start, c.at = start, at
		log = append(log, choiceAnswer(c))
	}
	for i := range foldBreakAnswersBefore {
		add(choiceSpec{skill: topic(topics[i%len(topics)]), place: gradeBand, difficulty: int32(lowDifficulty + i%difficulties),
			correct: i%wrongEvery != wrongEvery-1})
	}
	at += breakPause
	// The foldBreakAnswersAfter answers after the break.
	add(choiceSpec{skill: topic("a"), place: 2, difficulty: 3, correct: false})
	add(choiceSpec{skill: topic("a"), place: 2, difficulty: 3, correct: true})
	add(choiceSpec{skill: topic("b"), place: 2, difficulty: 4, correct: true})
	add(choiceSpec{skill: topic("b"), place: 3, difficulty: 2, correct: true})
	return log
}

// format04Log is pages of the workbook of kind 1, as api kept them before format 0.4, then pages of
// kind 2: one of times of day with the trap of the 12-hour clock, one with a number and a time in
// each item, and hard pages of times all right: two, then one after the leading question, which
// masters nothing, and the same page without it, which masters addition.
func format04Log() []lm.Evidence {
	const (
		kindBefore04 = 1  // the kind of evidence api kept before format 0.4
		kind04       = 2  // the kind that records the leading question on request
		pageItems    = 10 // the items of the pages before the hard ones
		timesDay     = 3  // the days of the log the pages of times, of both slots and of the hard
		bothDay      = 4  // pages fall on
		firstHardDay = 5
	)
	add := []string{skill("addition")}
	at := foldLogStart
	var log []lm.Evidence
	next := func(e lm.Evidence, kind int32) {
		e.Kind = kind
		log = append(log, e)
	}
	for i, right := range []int{8, 9, 10} {
		next(page(pageSpec{items: pageItems, right: right, primary: add, at: at + int64(day*i)}), kindBefore04)
	}
	times := page(pageSpec{items: pageItems, right: 8, primary: add, at: at + timesDay*day})
	for i := range times.Task.Items {
		times.Task.Items[i].Slots = []lm.Slot{{Check: lm.CheckKindTime}}
	}
	times.Traps = []lm.FoundTrap{{ID: workbook + "/trap/twelve-hours-off", Skills: add, ItemPlaces: []int32{9, 10}}}
	next(times, kind04)
	both := page(pageSpec{items: pageItems, right: 9, primary: add, at: at + bothDay*day})
	for i := range both.Task.Items {
		both.Task.Items[i].Slots = []lm.Slot{{Check: lm.CheckKindNumber}, {Check: lm.CheckKindTime}}
	}
	next(both, kind04)
	for i, leading := range []bool{false, false, true, false} {
		e := hard(page(pageSpec{items: 12, right: 12, primary: add, at: at + int64(day*(firstHardDay+i))}))
		for k := range e.Task.Items {
			e.Task.Items[k].Slots = []lm.Slot{{Check: lm.CheckKindTime}}
		}
		e.LeadingQuestionOnRequest = leading
		next(e, kind04)
	}
	return log
}

func TestGoldenFold(t *testing.T) {
	for _, j := range goldenJournals {
		t.Run(j.name, func(t *testing.T) {
			g := j.graphs(t)
			path := filepath.Join("testdata", "fold", j.name+".txt")
			if os.Getenv(updateEnv) == "1" {
				writeGoldenFold(t, path, j, g)
			}
			raw, err := os.ReadFile(path)
			if err != nil {
				t.Fatal(err)
			}
			var s lm.State
			var log []lm.Evidence
			var mastered, lost bool
			for n, line := range strings.Split(strings.TrimSuffix(string(raw), "\n"), "\n") {
				kind, hex, _ := strings.Cut(line, " ")
				var b []byte
				if kind != "#" {
					if _, err := fmt.Sscanf(hex, "%x", &b); err != nil {
						t.Fatalf("%s:%d: %v", path, n+1, err)
					}
				}
				switch kind {
				case "#":
				case "evidence":
					e, err := lm.DecodeEvidence(b)
					if err != nil {
						t.Fatalf("%s:%d: %v", path, n+1, err)
					}
					if s, e, err = lm.Append(s, e, g); err != nil {
						t.Fatalf("%s:%d: %v", path, n+1, err)
					}
					log = append(log, e)
				case "state":
					if got := encodeState(t, s); !bytes.Equal(got, b) {
						t.Fatalf("%s:%d: the state after %d pieces of evidence is\n%x\nnot\n%x", path, n+1, len(log), got, b)
					}
					for _, m := range s.Modules {
						for _, sk := range m.Skills {
							lost = lost || (mastered && !sk.Mastered && sk.FailuresInARow >= lm.FailuresToLose)
							mastered = mastered || sk.Mastered
						}
					}
				default:
					t.Fatalf("%s:%d: a line of kind %q", path, n+1, kind)
				}
			}
			again, err := lm.Recompute(log, g)
			if err != nil || !sameState(t, again, s) {
				t.Errorf("Recompute of the file's log gives another state: %v", err)
			}
			if j.name == "olympiad" && !(mastered && lost) {
				t.Errorf("the log never masters a skill and loses it: mastered %v, lost %v", mastered, lost)
			}
			if j.name == "format04" {
				// The page after the leading question moves the skill as the page without it does,
				// but only the page without it masters.
				before, err := lm.Recompute(log[:len(log)-1], g)
				if err != nil {
					t.Fatal(err)
				}
				sk := skillOf(t, s, skill("addition"))
				if !log[len(log)-2].LeadingQuestionOnRequest || skillOf(t, before, skill("addition")).Mastered || !sk.Mastered {
					t.Errorf("addition is mastered after the leading question %v, at the end %v; want only at the end",
						skillOf(t, before, skill("addition")).Mastered, sk.Mastered)
				}
				// The control: the same log without the leading question masters at that page, and moves
				// the skill to the same bits.
				control := slices.Clone(log[:len(log)-1])
				control[len(control)-1].LeadingQuestionOnRequest = false
				without, err := lm.Recompute(control, g)
				if err != nil {
					t.Fatal(err)
				}
				a, b := skillOf(t, before, skill("addition")), skillOf(t, without, skill("addition"))
				if !b.Mastered || a.Correction != b.Correction || a.AnswerWeight != b.AnswerWeight || a.Stability != b.Stability ||
					a.EvidenceCount != b.EvidenceCount {
					t.Errorf("without the leading question the skill is %+v; with it %+v: want it mastered, and moved alike", b, a)
				}
				if len(sk.Traps) != 1 || sk.Traps[0].ID != workbook+"/trap/twelve-hours-off" {
					t.Errorf("addition has the traps %v; want the trap of the 12-hour clock", sk.Traps)
				}
			}
			if j.name == "break" {
				// The answers before the break and after it, and n = 20 + the answers after it.
				if n := s.Modules[0].AnswerWeight; len(log) != foldBreakAnswersBefore+foldBreakAnswersAfter || n != lm.LevelWeight+foldBreakAnswersAfter {
					t.Errorf("after the break n is %v over %d pieces of evidence, want 24 over 28", n, len(log))
				}
			}
			if j.name == "transfer" {
				var skills []string
				for _, sk := range s.Modules[0].Skills {
					skills = append(skills, strings.TrimPrefix(sk.Skill, workbook+"/skill/"))
				}
				if fmt.Sprint(skills) != "[counting place subtraction2]" || !skillOf(t, s, skill("counting")).Mastered || skillOf(t, s, skill("place")).Mastered {
					t.Errorf("the transfer gives skills %v; want counting mastered, place not, and subtraction2", skills)
				}
			}
		})
	}
}

func writeGoldenFold(t *testing.T, path string, j goldenJournal, g lm.Graphs) {
	t.Helper()
	// The folder and the files are readable by all and writable by their owner.
	const (
		foldDirMode  = 0o755
		foldFileMode = 0o644
	)
	var out strings.Builder
	fmt.Fprintf(&out, "# %s\n", j.note)
	var s lm.State
	for _, e := range j.log() {
		if j.kind != ownKinds {
			e.Kind = j.kind
		}
		b, err := lm.EncodeEvidence(e)
		if err != nil {
			t.Fatal(err)
		}
		if s, _, err = lm.Append(s, e, g); err != nil {
			t.Fatal(err)
		}
		fmt.Fprintf(&out, "evidence %x\nstate %x\n", b, encodeState(t, s))
	}
	if err := os.MkdirAll(filepath.Dir(path), foldDirMode); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, []byte(out.String()), foldFileMode); err != nil {
		t.Fatal(err)
	}
}
