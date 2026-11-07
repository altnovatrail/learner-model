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
// holds it, a second that takes it away, and a mastery again.
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
	add(topic("a"), 2, 5, false)
	for range 4 {
		add(topic("a"), 2, 2, true)
	}
	add(topic("b"), 3, 2, false)
	return log
}

// workbookLog is pages of the workbook: two primary skills on every item, J = 2, an item that
// names its own skill, one with a rubric, items with no answer and a failed page.
func workbookLog() []lm.Evidence {
	both := []string{skill("addition"), skill("subtraction")}
	var log []lm.Evidence
	for i, p := range []struct{ items, right int }{{10, 8}, {10, 9}, {6, 2}, {12, 12}, {10, 7}, {10, 3}, {10, 3}, {20, 20}} {
		e := page(p.items, p.right, both, 1762041600+int64(86400*i))
		e.Task.Items[0].OwnSkills, e.Task.Items[0].Primary = true, []string{skill("addition")}
		e.Task.Items[1].OwnSkills, e.Task.Items[1].Primary = true, []string{skill("essay")}
		e.Items[p.items-1].HasFirstVerdict, e.Items[p.items-1].FirstVerdict, e.Items[p.items-1].Submitted = false, 0, false
		log = append(log, e)
	}
	return log
}

// mixedLog interleaves the two, with evidence of a module that has no current version and
// evidence with no answer.
func mixedLog() []lm.Evidence {
	o, w := olympiadLog(), workbookLog()
	waiting := choiceAnswer(topic("a"), 1, 3, true, ladderBand(1), 1762041600)
	waiting.Task.ID = "example/waiting/task/x"
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
	return append(log, waiting)
}

var goldenJournals = []goldenJournal{
	{"olympiad", "a learner of grade 3 in olympiad maths: the trial series, mastery held, lost and set again", foldGraphs, olympiadLog},
	{"workbook", "pages of the workbook with two primary skills, J = 2", foldGraphs, workbookLog},
	{"mixed", "the two interleaved, with a module without a current version and evidence with no answer", foldGraphs, mixedLog},
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
