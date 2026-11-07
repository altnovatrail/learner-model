package learnermodel

import (
	"slices"
	"strconv"
	"strings"
)

// NewGraph builds the graph of a module version from its bands, its skills, its retired skills
// and the descriptions of its tasks, so that api and sim build one graph of one module alike
// (R10). It derives each skill's BandsTaught — the BandsTaught given is ignored — and orders
// every list canonically. It copies what it is given: the graph shares nothing with its input.
//
// A skill is taught at the band of a task where it is a primary skill of an item: of the item's
// own skills, or of the task's for an item without its own, or for a task without items.
//
// It refuses, with ErrInvalid: bands out of their places 1, 2, …; a duplicate band, skill or
// retired ID, and an ID or grade twice in one list; a retired ID that is a skill; a replacement or a
// primary skill of a task or an item the module does not have; a task band outside the bands or
// whose ID is not the ID of the band at its place; and a difficulty outside 1 to 5.
func NewGraph(module string, bands []Band, skills []Skill, retired []RetiredSkill, tasks []Task) (Graph, error) {
	g := Graph{Module: module}
	bandIDs := map[string]bool{}
	for i, b := range bands {
		if int(b.Place) != i+1 {
			return Graph{}, invalid("band " + b.ID + " at place " + strconv.Itoa(int(b.Place)) + ", not " + strconv.Itoa(i+1))
		}
		if bandIDs[b.ID] {
			return Graph{}, invalid("band " + b.ID + " twice")
		}
		bandIDs[b.ID] = true
		if err := uniqueSorted("a grade of band "+b.ID, sortedCopy(b.Grades)); err != nil {
			return Graph{}, err
		}
		g.Bands = append(g.Bands, Band{Place: b.Place, ID: b.ID, Grades: sortedCopy(b.Grades)})
	}

	taught := map[string][]int32{}
	for _, s := range skills {
		if _, ok := taught[s.ID]; ok {
			return Graph{}, invalid("skill " + s.ID + " twice")
		}
		taught[s.ID] = nil
	}
	for _, t := range tasks {
		if t.Band.Place < 1 || int(t.Band.Place) > len(bands) {
			return Graph{}, invalid("task " + t.ID + " in band place " + strconv.Itoa(int(t.Band.Place)))
		}
		if t.Band.ID != bands[t.Band.Place-1].ID {
			return Graph{}, invalid("task " + t.ID + " in band " + t.Band.ID + " at the place of band " + bands[t.Band.Place-1].ID)
		}
		if t.Difficulty < 1 || t.Difficulty > 5 {
			return Graph{}, invalid("task " + t.ID + " of difficulty " + strconv.Itoa(int(t.Difficulty)))
		}
		// The task's own primary skills count when an item takes them, or when the task has no
		// items to say otherwise; an item's own skills take the place of the task's for it.
		var primary []string
		takesTask := len(t.Items) == 0
		for _, item := range t.Items {
			if item.OwnSkills {
				primary = append(primary, item.Primary...)
			} else {
				takesTask = true
			}
		}
		for _, id := range slices.Concat(primary, t.Primary) {
			if _, ok := taught[id]; !ok {
				return Graph{}, invalid("task " + t.ID + " trains skill " + id + ", which the module does not have")
			}
		}
		if takesTask {
			primary = append(primary, t.Primary...)
		}
		for _, id := range primary {
			if places := taught[id]; !slices.Contains(places, t.Band.Place) {
				taught[id] = append(places, t.Band.Place)
			}
		}
	}

	for _, s := range skills {
		if err := uniqueSorted("an anchor skill that skill "+s.ID+" teaches", sortedCopy(s.Teaches)); err != nil {
			return Graph{}, err
		}
		if err := uniqueSorted("a prerequisite of skill "+s.ID, sortedCopy(s.Prerequisites)); err != nil {
			return Graph{}, err
		}
		g.Skills = append(g.Skills, Skill{
			ID:            s.ID,
			Assessability: s.Assessability,
			Teaches:       sortedCopy(s.Teaches),
			BandsTaught:   sortedCopy(taught[s.ID]),
			Prerequisites: sortedCopy(s.Prerequisites),
			Tier:          s.Tier,
		})
	}
	slices.SortFunc(g.Skills, func(a, b Skill) int { return strings.Compare(a.ID, b.ID) })

	seen := map[string]bool{}
	for _, r := range retired {
		if _, isSkill := taught[r.ID]; isSkill || seen[r.ID] {
			return Graph{}, invalid("retired skill " + r.ID + " twice or still a skill")
		}
		seen[r.ID] = true
		if err := uniqueSorted("a replacement of retired skill "+r.ID, sortedCopy(r.ReplacedBy)); err != nil {
			return Graph{}, err
		}
		for _, id := range r.ReplacedBy {
			if _, ok := taught[id]; !ok {
				return Graph{}, invalid("retired skill " + r.ID + " replaced by " + id + ", which the module does not have")
			}
		}
		g.Retired = append(g.Retired, RetiredSkill{ID: r.ID, ReplacedBy: sortedCopy(r.ReplacedBy)})
	}
	slices.SortFunc(g.Retired, func(a, b RetiredSkill) int { return strings.Compare(a.ID, b.ID) })
	return g, nil
}

// sortedCopy gives a sorted copy of s, nil for an empty one: strings byte by byte, numbers
// ascending.
func sortedCopy[T int32 | string](s []T) []T {
	if len(s) == 0 {
		return nil
	}
	c := slices.Clone(s)
	slices.Sort(c)
	return c
}
