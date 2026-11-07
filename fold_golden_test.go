package learnermodel_test

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	lm "github.com/altnovatrail/learner-model"
)

// The golden vectors of the fold (docs/architecture.md, section 3.6): synthetic logs, each piece
// of evidence with the bytes of the state after it, one file for every build, testdata/fold. A
// line is "evidence <hex>" or "state <hex>", the bytes of the codec; a line opening with # is a
// note. The test reads the evidence from the file, appends it and holds the state's bytes to the
// next line, so that the file pins what goes in as well as what comes out. LEARNERMODEL_UPDATE=1
// writes the files anew from the logs below, and they change only with the version of
// learner-model.

// goldenJournal is a log of a golden file with the graphs it is folded against.
type goldenJournal struct {
	name   string
	note   string
	graphs func(*testing.T) lm.Graphs
	log    func() []lm.Evidence
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
	start := ladderBand(2)
	var log []lm.Evidence
	at := int64(1762041600)
	add := func(skillID string, place, d int32, correct bool) {
		at += 3600
		log = append(log, choiceAnswer(skillID, place, d, correct, start, at))
	}
	add(topic("a"), 2, 3, true)
	add(topic("b"), 2, 4, false)
	add(topic("a"), 2, 4, true)
	add(topic("b"), 3, 1, true)
	add(topic("a"), 1, 5, true)
	for range 6 {
		add(topic("a"), 2, 3, true)
	}
	add(topic("a"), 2, 4, false)
	log[len(log)-1].Traps = []lm.FoundTrap{{ID: olympiad + "/trap/order", Skills: []string{topic("a")}, ItemPlaces: []int32{1}}}
	add(topic("a"), 2, 5, false)
	at += 60 * 86400 // a summer break: retention falls, and the next right answers raise τ much
	for range 4 {
		add(topic("a"), 2, 2, true)
	}
	add(topic("b"), 3, 2, false)
	return log
}

// workbookLog is pages of the workbook: two primary skills on every item, J = 2, an item that
// names its own skill, one with a rubric, items with no answer and a failed page; a page after a
// hint, one with the solution shown, one with the technique on request, and traps found.
func workbookLog() []lm.Evidence {
	both := []string{skill("addition"), skill("subtraction")}
	var log []lm.Evidence
	for i, p := range []struct{ items, right int }{{10, 8}, {10, 9}, {6, 2}, {12, 12}, {10, 7}, {10, 3}, {10, 3}, {20, 20}} {
		e := page(p.items, p.right, both, 1762041600+int64(86400*i))
		e.Task.Items[0].OwnSkills, e.Task.Items[0].Primary = true, []string{skill("addition")}
		e.Task.Items[1].OwnSkills, e.Task.Items[1].Primary = true, []string{skill("essay")}
		e.Items[p.items-1].HasFirstVerdict, e.Items[p.items-1].FirstVerdict, e.Items[p.items-1].Submitted = false, 0, false
		switch i {
		case 2: // a page after a hint: the items first answered after the rung are with help
			e.Outcome = lm.Outcome{Kind: lm.OutcomeAfterHint, Rung: 1, HasRung: true}
			for k := 3; k < p.items-1; k++ { // the first wrong answer, item 3, came before them
				e.Items[k].WrongBeforeFirst = 1
			}
		case 5: // the solution shown: the item never submitted is not solved
			e.Outcome = lm.Outcome{Kind: lm.OutcomeSolutionShown}
		case 7: // the technique on request: a page all right that masters nothing
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
	waiting := choiceAnswer(topic("a"), 1, 3, true, ladderBand(1), 1762041600)
	waiting.Task.ID = "example/waiting/task/x"
	ai := page(4, 4, []string{skill("addition")}, 1762041602)
	ai.Task.SkillTagsProvenance = lm.ProvenanceAI // a task from paper the AI tagged: it weighs nothing
	none := page(3, 3, []string{skill("addition")}, 1762041601)
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
	{"olympiad", "a learner of grade 3 in olympiad maths: the trial series, mastery held, lost and set again after a break of 60 days", foldGraphs, olympiadLog},
	{"workbook", "pages of the workbook with two primary skills, J = 2, a hint, a shown solution, help and traps", foldGraphs, workbookLog},
	{"mixed", "the two interleaved, with a module without a current version, evidence with no answer and a task tagged by AI", foldGraphs, mixedLog},
	{"transfer", "the workbook across three versions, by the third: a split, a merge, a rename, a skill retired with [] and one unknown", transferGraphs, transferLog3},
}

func TestGoldenFold(t *testing.T) {
	for _, j := range goldenJournals {
		t.Run(j.name, func(t *testing.T) {
			g := j.graphs(t)
			path := filepath.Join("testdata", "fold", j.name+".txt")
			if os.Getenv(updateEnv) == "1" {
				writeGoldenFold(t, path, j.note, j.log(), g)
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
							lost = lost || (mastered && !sk.Mastered && sk.FailuresInARow >= 2)
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

func writeGoldenFold(t *testing.T, path, note string, log []lm.Evidence, g lm.Graphs) {
	t.Helper()
	var out strings.Builder
	fmt.Fprintf(&out, "# %s\n", note)
	var s lm.State
	for _, e := range log {
		b, err := lm.EncodeEvidence(e)
		if err != nil {
			t.Fatal(err)
		}
		if s, _, err = lm.Append(s, e, g); err != nil {
			t.Fatal(err)
		}
		fmt.Fprintf(&out, "evidence %x\nstate %x\n", b, encodeState(t, s))
	}
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, []byte(out.String()), 0o644); err != nil {
		t.Fatal(err)
	}
}
