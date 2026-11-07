package learnermodel_test

import (
	"errors"
	"fmt"
	"testing"

	lm "github.com/altnovatrail/learner-model"
)

func graphInput() ([]lm.Band, []lm.Skill, []lm.RetiredSkill, []lm.Task) {
	bands := []lm.Band{
		{Place: 1, ID: "grade-2", Grades: []int32{2}},
		{Place: 2, ID: "grade-3", Grades: []int32{3}},
		{Place: 3, ID: "grade-4", Grades: []int32{4}},
	}
	skills := []lm.Skill{
		{ID: skill("subtraction"), Assessability: lm.AssessabilityVerdict, Prerequisites: []string{skill("addition")},
			Tier: 2, BandsTaught: []int32{9}},
		{ID: skill("addition"), Assessability: lm.AssessabilityVerdict, Teaches: []string{"core.b", "core.a"}, Tier: 1},
		{ID: skill("comparison"), Assessability: lm.AssessabilityVerdict, Tier: 1},
		{ID: skill("essay"), Assessability: lm.AssessabilityRubric, Tier: 1},
	}
	retired := []lm.RetiredSkill{
		{ID: skill("minus"), ReplacedBy: []string{skill("subtraction")}},
		{ID: skill("counting"), ReplacedBy: nil},
	}
	tasks := []lm.Task{
		{ID: "t1", Band: lm.BandRef{Place: 2, ID: "grade-3"}, Difficulty: 3, Primary: []string{skill("subtraction")}},
		{ID: "t2", Band: lm.BandRef{Place: 1, ID: "grade-2"}, Difficulty: 1, Primary: []string{skill("subtraction"), skill("addition")}},
		{ID: "t3", Band: lm.BandRef{Place: 3, ID: "grade-4"}, Difficulty: 5, Supporting: []string{skill("comparison")},
			Items: []lm.TaskItem{
				{Place: 1, OwnSkills: true, Primary: []string{skill("comparison")}},
				{Place: 2, Primary: []string{skill("essay")}}, // not its own skills: they do not count
			}},
		{ID: "t4", Band: lm.BandRef{Place: 2, ID: "grade-3"}, Difficulty: 2, Primary: []string{skill("subtraction")}},
		// Every item names its own skills: the task's primary skill is trained by none of them.
		{ID: "t5", Band: lm.BandRef{Place: 3, ID: "grade-4"}, Difficulty: 1, Primary: []string{skill("addition")},
			Items: []lm.TaskItem{{Place: 1, OwnSkills: true, Primary: []string{skill("comparison")}}}},
	}
	return bands, skills, retired, tasks
}

func TestNewGraph(t *testing.T) {
	bands, skills, retired, tasks := graphInput()
	g, err := lm.NewGraph(module, bands, skills, retired, tasks)
	if err != nil {
		t.Fatal(err)
	}
	got := map[string]string{}
	for _, s := range g.Skills {
		got[s.ID] = fmt.Sprint(s.BandsTaught, s.Teaches, s.Prerequisites, s.Tier)
	}
	want := map[string]string{
		skill("addition"):    "[1] [core.a core.b] [] 1",
		skill("comparison"):  "[3] [] [] 1",
		skill("essay"):       "[] [] [] 1",
		skill("subtraction"): "[1 2] [] [" + skill("addition") + "] 2",
	}
	if fmt.Sprint(got) != fmt.Sprint(want) {
		t.Errorf("skills: %v, want %v", got, want)
	}
	if g.Skills[0].ID != skill("addition") || g.Retired[0].ID != skill("counting") || g.Retired[0].ReplacedBy != nil {
		t.Errorf("not in canonical order: %+v", g)
	}
	if g.Module != module || len(g.Bands) != 3 || g.Bands[1].ID != "grade-3" {
		t.Errorf("bands: %+v", g.Bands)
	}
	skills[0].Teaches = append(skills[0].Teaches, "core.changed")
	bands[0].Grades[0] = 7
	if len(g.Skills[3].Teaches) != 0 || g.Bands[0].Grades[0] != 2 {
		t.Errorf("the graph shares its input")
	}
}

func TestNewGraphRefuses(t *testing.T) {
	for name, change := range map[string]func([]lm.Band, []lm.Skill, []lm.RetiredSkill, []lm.Task) ([]lm.Band, []lm.Skill, []lm.RetiredSkill, []lm.Task){
		"a band out of place": func(b []lm.Band, s []lm.Skill, r []lm.RetiredSkill, t []lm.Task) ([]lm.Band, []lm.Skill, []lm.RetiredSkill, []lm.Task) {
			b[1].Place = 5
			return b, s, r, t
		},
		"a band twice": func(b []lm.Band, s []lm.Skill, r []lm.RetiredSkill, t []lm.Task) ([]lm.Band, []lm.Skill, []lm.RetiredSkill, []lm.Task) {
			b[1].ID = b[0].ID
			return b, s, r, t
		},
		"a skill twice": func(b []lm.Band, s []lm.Skill, r []lm.RetiredSkill, t []lm.Task) ([]lm.Band, []lm.Skill, []lm.RetiredSkill, []lm.Task) {
			return b, append(s, s[0]), r, t
		},
		"a task band past the bands": func(b []lm.Band, s []lm.Skill, r []lm.RetiredSkill, t []lm.Task) ([]lm.Band, []lm.Skill, []lm.RetiredSkill, []lm.Task) {
			t[0].Band.Place = 4
			return b, s, r, t
		},
		"a task band of 0": func(b []lm.Band, s []lm.Skill, r []lm.RetiredSkill, t []lm.Task) ([]lm.Band, []lm.Skill, []lm.RetiredSkill, []lm.Task) {
			t[0].Band.Place = 0
			return b, s, r, t
		},
		"a difficulty of 6": func(b []lm.Band, s []lm.Skill, r []lm.RetiredSkill, t []lm.Task) ([]lm.Band, []lm.Skill, []lm.RetiredSkill, []lm.Task) {
			t[0].Difficulty = 6
			return b, s, r, t
		},
		"an unknown primary skill": func(b []lm.Band, s []lm.Skill, r []lm.RetiredSkill, t []lm.Task) ([]lm.Band, []lm.Skill, []lm.RetiredSkill, []lm.Task) {
			t[2].Items[0].Primary = []string{"nope"}
			return b, s, r, t
		},
		"a retired ID that is a skill": func(b []lm.Band, s []lm.Skill, r []lm.RetiredSkill, t []lm.Task) ([]lm.Band, []lm.Skill, []lm.RetiredSkill, []lm.Task) {
			r[0].ID = s[0].ID
			return b, s, r, t
		},
		"a retired ID twice": func(b []lm.Band, s []lm.Skill, r []lm.RetiredSkill, t []lm.Task) ([]lm.Band, []lm.Skill, []lm.RetiredSkill, []lm.Task) {
			return b, s, append(r, r[0]), t
		},
		"an unknown replacement": func(b []lm.Band, s []lm.Skill, r []lm.RetiredSkill, t []lm.Task) ([]lm.Band, []lm.Skill, []lm.RetiredSkill, []lm.Task) {
			r[1].ReplacedBy = []string{"nope"}
			return b, s, r, t
		},
		"a replacement twice": func(b []lm.Band, s []lm.Skill, r []lm.RetiredSkill, t []lm.Task) ([]lm.Band, []lm.Skill, []lm.RetiredSkill, []lm.Task) {
			r[0].ReplacedBy = append(r[0].ReplacedBy, r[0].ReplacedBy[0])
			return b, s, r, t
		},
		"an anchor taught twice": func(b []lm.Band, s []lm.Skill, r []lm.RetiredSkill, t []lm.Task) ([]lm.Band, []lm.Skill, []lm.RetiredSkill, []lm.Task) {
			s[1].Teaches = []string{"core.a", "core.a"}
			return b, s, r, t
		},
		"a prerequisite twice": func(b []lm.Band, s []lm.Skill, r []lm.RetiredSkill, t []lm.Task) ([]lm.Band, []lm.Skill, []lm.RetiredSkill, []lm.Task) {
			s[0].Prerequisites = append(s[0].Prerequisites, s[0].Prerequisites[0])
			return b, s, r, t
		},
		"a grade twice": func(b []lm.Band, s []lm.Skill, r []lm.RetiredSkill, t []lm.Task) ([]lm.Band, []lm.Skill, []lm.RetiredSkill, []lm.Task) {
			b[0].Grades = []int32{2, 2}
			return b, s, r, t
		},
		"a task band of another ID": func(b []lm.Band, s []lm.Skill, r []lm.RetiredSkill, t []lm.Task) ([]lm.Band, []lm.Skill, []lm.RetiredSkill, []lm.Task) {
			t[0].Band.ID = "grade-5"
			return b, s, r, t
		},
		"an unknown task skill no item takes": func(b []lm.Band, s []lm.Skill, r []lm.RetiredSkill, t []lm.Task) ([]lm.Band, []lm.Skill, []lm.RetiredSkill, []lm.Task) {
			t[4].Primary = []string{"nope"}
			return b, s, r, t
		},
	} {
		b, s, r, tasks := change(graphInput())
		_, err := lm.NewGraph(module, b, s, r, tasks)
		if !errors.Is(err, lm.ErrInvalid) {
			t.Errorf("%s: %v, not ErrInvalid", name, err)
		}
	}
}
