package learnermodel_test

import (
	"fmt"
	"math"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"testing"

	lm "github.com/altnovatrail/learner-model"
)

// The golden vectors of reading: the logs of the golden files of the fold, testdata/fold,
// folded and read at three times, with what Read and Readout.Task give written out as text, every
// float as its bits, testdata/read. The test takes the evidence from the fold's file, so that the
// two files pin one log. LEARNERMODEL_UPDATE=1 writes the files anew, and they change only with the
// version of learner-model.

// anchorsLog is a learner of grade 3 in olympiad maths and school maths, which teaches the anchor
// skills olympiad maths relies on: division is mastered over pages of school maths.
func anchorsLog() []lm.Evidence {
	const (
		// Three answers of olympiad maths an hour apart, the second of them wrong.
		answers     = 3
		hour        = 3600
		wrongAnswer = 1
		// Then six pages of school maths a day apart, of four items each, all right but one item
		// of the fifth page.
		pages          = 6
		pageItems      = 4
		pageDifficulty = 4
		oneWrongPage   = 4
	)
	at := int64(readEpoch)
	var log []lm.Evidence
	for i := range answers {
		at += hour
		log = append(log, choiceAnswer(choiceSpec{skill: topic("a"), place: 2, difficulty: 3, correct: i != wrongAnswer, start: ladderBand(2), at: at}))
	}
	for i := range pages {
		at += day
		right := pageItems
		if i == oneWrongPage {
			right--
		}
		e := page(pageSpec{items: pageItems, right: right, primary: []string{schoolSkill("division")}, at: at})
		e.Task.ID, e.Task.Band, e.Task.Difficulty = school+"/task/page", lm.BandRef{Place: 1, ID: "grade-3"}, pageDifficulty
		e.StartBand = lm.BandRef{Place: 1, ID: "grade-3"}
		log = append(log, e)
	}
	return log
}

func init() {
	goldenJournals = append(goldenJournals, goldenJournal{name: "anchors",
		note:   "a learner of grade 3 in olympiad maths and school maths, where division with remainder is mastered",
		graphs: anchorGraphs, log: anchorsLog, kind: 1})
}

// goldenReadPause is the pause after a golden log's last evidence of its last reading, past the
// pause after which a module's level forgets.
const goldenReadPause = 60 * day

// goldenLog reads the evidence of a golden file of the fold.
func goldenLog(t *testing.T, name string) []lm.Evidence {
	t.Helper()
	raw, err := os.ReadFile(filepath.Join("testdata", "fold", name+".txt"))
	if err != nil {
		t.Fatal(err)
	}
	var log []lm.Evidence
	for _, line := range strings.Split(strings.TrimSuffix(string(raw), "\n"), "\n") {
		if hex, ok := strings.CutPrefix(line, "evidence "); ok {
			var b []byte
			if _, err := fmt.Sscanf(hex, "%x", &b); err != nil {
				t.Fatal(err)
			}
			e, err := lm.DecodeEvidence(b)
			if err != nil {
				t.Fatal(err)
			}
			log = append(log, e)
		}
	}
	return log
}

func bits(v float64) string { return fmt.Sprintf("%016x(%.6g)", math.Float64bits(v), v) }

func band(b lm.BandRef) string { return fmt.Sprintf("%d/%s", b.Place, b.ID) }

func review(r lm.ReviewDate) string {
	if r.Already {
		return "already"
	}
	return fmt.Sprint(r.At)
}

// writeRead writes out what a reading gives, and Task of each of the candidates.
func writeRead(t *testing.T, out *strings.Builder, r lm.Readout, candidates []lm.Task) {
	t.Helper()
	for _, m := range r.Modules {
		level := "none"
		if m.HasLevelBand {
			level = band(m.LevelBand)
		}
		fmt.Fprintf(out, "module %s evidence=%v start=%s level=%s trial-over=%v last-in-trial=%v",
			m.Module, m.HasEvidence, band(m.StartBand), level, m.TrialOver, m.LastInTrial)
		for _, trap := range m.Traps {
			fmt.Fprintf(out, " trap=%s:%d:%d", trap.ID, trap.Count, trap.Last)
		}
		fmt.Fprintln(out)
		for _, sk := range m.Skills {
			fmt.Fprintf(out, "  skill %s count=%d last=%d estimate=%s uncertainty=%s retention=%s",
				sk.Skill, sk.EvidenceCount, sk.LastPractice, bits(sk.Estimate), bits(sk.Uncertainty), bits(sk.Retention))
			for _, trap := range sk.Traps {
				fmt.Fprintf(out, " trap=%s:%d:%d", trap.ID, trap.Count, trap.Last)
			}
			if sk.Mastered {
				fmt.Fprintf(out, " mastered=%s cautious=%s review=%s", band(sk.MasteryBand), bits(sk.CautiousProbability), review(sk.Review))
			}
			fmt.Fprintln(out)
		}
	}
	for _, a := range r.Anchors {
		fmt.Fprintf(out, "anchor %s status=%d reason=%d\n", a.Anchor, a.Status, a.Reason)
		for _, tr := range a.Teachers {
			fmt.Fprintf(out, "  teacher %s highest=%s", tr.Skill, band(tr.HighestBand))
			if tr.Practised {
				fmt.Fprintf(out, " estimate=%s retention=%s mastered=%v", bits(tr.Estimate), bits(tr.Retention), tr.Mastered)
				if tr.Mastered {
					fmt.Fprintf(out, " band=%s review=%s", band(tr.MasteryBand), review(tr.Review))
				}
			}
			fmt.Fprintln(out)
		}
	}
	for _, task := range candidates {
		got, err := r.Task(task)
		if err != nil {
			t.Fatal(err)
		}
		fmt.Fprintf(out, "task %s\n", task.ID)
		for _, p := range got.Items {
			if p.HasVerdict {
				fmt.Fprintf(out, "  item %d p=%s\n", p.Place, bits(p.Probability))
			} else {
				fmt.Fprintf(out, "  item %d no verdict\n", p.Place)
			}
		}
	}
}

// readGolden gives the text of the golden file of reading a log: the state at its end read at its
// last practice time and a day later with all the modules of the graphs the engine's and a grade,
// and 60 days later with none and no grade; the candidates are the last task of each module with a
// current version.
func readGolden(t *testing.T, note string, log []lm.Evidence, g lm.Graphs) string {
	t.Helper()
	s, err := lm.Recompute(log, g)
	if err != nil {
		t.Fatal(err)
	}
	const moduleParts = 2 // a module's ID is the first two parts of a task's, <publisher>/<name>
	var last int64
	var candidates []lm.Task
	for i := len(log) - 1; i >= 0; i-- {
		e := log[i]
		last = max(last, e.PracticeTime)
		module := strings.Join(strings.SplitN(e.Task.ID, "/", moduleParts+1)[:moduleParts], "/")
		if !e.LevelSetByAdult && g[module].Current && !slices.ContainsFunc(candidates, func(c lm.Task) bool { return strings.HasPrefix(c.ID, module+"/") }) {
			candidates = append(candidates, e.Task)
		}
	}
	var engine []string
	for _, id := range sortedKeys(g) {
		if g[id].Current {
			engine = append(engine, id)
		}
	}
	var out strings.Builder
	fmt.Fprintf(&out, "# %s\n", note)
	reads := []struct {
		at     int64
		grade  lm.Grade
		engine []string
	}{
		{last, lm.Grade{Value: 3, Known: true}, engine},
		{last + day, lm.Grade{Value: 3, Known: true}, engine},
		{last + goldenReadPause, lm.Grade{}, nil},
	}
	for _, rd := range reads {
		r, err := lm.Read(s, g, rd.at, rd.grade, rd.engine)
		if err != nil {
			t.Fatal(err)
		}
		fmt.Fprintf(&out, "read at=%d grade=%d/%v engine=%s\n", rd.at, rd.grade.Value, rd.grade.Known, strings.Join(rd.engine, ","))
		writeRead(t, &out, r, candidates)
	}
	return out.String()
}

func TestGoldenRead(t *testing.T) {
	const dirMode, fileMode = 0o755, 0o644 // the golden files are readable by all, writable by their owner
	for _, j := range goldenJournals {
		t.Run(j.name, func(t *testing.T) {
			path := filepath.Join("testdata", "read", j.name+".txt")
			got := readGolden(t, j.note, goldenLog(t, j.name), j.graphs(t))
			if os.Getenv(updateEnv) == "1" {
				if err := os.MkdirAll(filepath.Dir(path), dirMode); err != nil {
					t.Fatal(err)
				}
				if err := os.WriteFile(path, []byte(got), fileMode); err != nil {
					t.Fatal(err)
				}
			}
			want, err := os.ReadFile(path)
			if err != nil {
				t.Fatal(err)
			}
			if got != string(want) {
				gl, wl := strings.Split(got, "\n"), strings.Split(string(want), "\n")
				for i := range min(len(gl), len(wl)) {
					if gl[i] != wl[i] {
						t.Fatalf("%s:%d: reading gives\n%s\nnot\n%s", path, i+1, gl[i], wl[i])
					}
				}
				t.Fatalf("%s: reading gives %d lines, the file has %d", path, len(gl), len(wl))
			}
			if j.name == "anchors" && !strings.Contains(got, "anchor "+division+" status=1") {
				t.Errorf("the anchors log never holds %s", division)
			}
		})
	}
}
