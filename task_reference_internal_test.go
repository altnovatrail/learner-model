package learnermodel

import (
	"cmp"
	"slices"
	"strings"
	"testing"
)

// The reference of Readout.Task (T40 in RUN.md): Task as v0.1.5 had it, which copied the tags and
// the items of each candidate, split its ID, and worked out R of its skills at t₀ anew. Its text
// and that of its checks are v0.1.5's word for word, but for their names and for two lines: which
// modules Read has checked comes from the field that holds them now, and P comes of chanceOf as
// the step calls it — R worked out anew into an empty list, and the shares gathered —, the very
// count v0.1.5 made on reading. So TestTaskIsReference holds Task's checks and its R and start,
// which Read now works out once, to v0.1.5's, and its P to the step's, bit for bit.

// TaskReference is Readout.Task of v0.1.5, for the tests outside the package.
func (r Readout) TaskReference(t Task) (TaskReadout, error) {
	if _, err := checkTaskContentReference(t); err != nil {
		return TaskReadout{}, err
	}
	module, ok := moduleOfReference(t.ID)
	if !ok {
		return TaskReadout{}, invalid("task " + shownBytes(t.ID) + " names no module")
	}
	if err := checkTaskPointReference("task "+t.ID, t.Band.Place, t.Difficulty); err != nil {
		return TaskReadout{}, err
	}
	items := slices.Clone(t.Items)
	slices.SortFunc(items, func(a, b TaskItem) int { return cmp.Compare(a.Place, b.Place) })
	for _, it := range items {
		if err := checkChoices(it.Slots); err != nil {
			return TaskReadout{}, err
		}
	}
	mg, ok := r.graphs[module]
	if !ok {
		return TaskReadout{}, invalid("module " + module + " has no entry in the graphs")
	}
	if !mg.Current {
		return TaskReadout{}, invalid("module " + module + " of task " + t.ID + " has no current version")
	}
	// Read has checked the graphs of its modules; a candidate of another module is checked here.
	if _, checked := r.checkedModule(module); !checked {
		if err := checkGraph(module, mg.Graph); err != nil {
			return TaskReadout{}, err
		}
	}
	var ms ModuleState
	if i, ok := r.state.moduleIndex(module); ok {
		ms = r.state.Modules[i]
	} else {
		start, ok, err := startingPoint(r.grade, mg.Graph.Bands)
		if err != nil {
			return TaskReadout{}, err
		}
		if !ok {
			return TaskReadout{}, invalid("task " + t.ID + " of module " + module + ", whose current version has no bands")
		}
		ms = newModule(module, start)
	}
	trial := ms.inTrial()
	var rs []skillRetention
	out := TaskReadout{Items: make([]ItemProbability, 0, len(items))}
	for _, it := range items {
		ip := ItemProbability{Place: it.Place}
		if len(it.Slots) > 0 && knownChecks(it.Slots) {
			ip.HasVerdict = true
			// P as the step takes it: R worked out anew, and the shares gathered.
			var shares []share
			ip.Probability = ms.chanceOf(t, it, &mg.Graph, &rs, &shares, r.at, trial).p
		}
		out.Items = append(out.Items, ip)
	}
	return out, nil
}

func checkTaskContentReference(t Task) ([]int32, error) {
	if err := utf8Strings(t.ID, t.Band.ID); err != nil {
		return nil, err
	}
	if err := idListReference("primary skill", t.Primary); err != nil {
		return nil, err
	}
	if err := idListReference("supporting skill", t.Supporting); err != nil {
		return nil, err
	}
	places := make([]int32, 0, len(t.Items))
	for _, item := range t.Items {
		if err := idListReference("primary skill of an item", item.Primary); err != nil {
			return nil, err
		}
		if err := idListReference("supporting skill of an item", item.Supporting); err != nil {
			return nil, err
		}
		if err := checkSlots(item.Slots); err != nil {
			return nil, err
		}
		places = append(places, item.Place)
	}
	slices.Sort(places)
	return places, uniqueSorted("item place", places)
}

func idListReference(what string, ids []string) error {
	if err := utf8Strings(ids...); err != nil {
		return err
	}
	return uniqueSorted(what, sortedCopy(ids))
}

func moduleOfReference(id string) (string, bool) {
	parts := strings.SplitN(id, "/", 4)
	if len(parts) < 4 || slices.Contains(parts, "") {
		return "", false
	}
	return parts[0] + "/" + parts[1], true
}

func checkTaskPointReference(what string, place, difficulty int32) error {
	if place < 1 {
		return invalid(what + " in band place " + itoa(place))
	}
	if difficulty < 1 || difficulty > 5 {
		return invalid(what + " of difficulty " + itoa(difficulty))
	}
	return nil
}

// moduleOf gives what the split of v0.1.5 gave, on IDs of every shape it tells apart.
func TestModuleOfIsReference(t *testing.T) {
	ids := []string{
		"", "/", "a", "a/b", "a/b/c", "a/b/c/", "a/b/c/d", "a/b/c/d/e", "a/b/c//", "a/b//d", "a//c/d", "/b/c/d",
		"a/b/c/d/", "ab/cd/task/x", "example/m0/task/t1", "example/m0/task/t/1", "example/m0/task", "///", "////",
		"a/b/c/\xff", "\xff/b/c/d",
	}
	for _, id := range ids {
		got, ok := moduleOf(id)
		want, wantOK := moduleOfReference(id)
		if got != want || ok != wantOK {
			t.Errorf("moduleOf(%q) = %q, %v; v0.1.5 gave %q, %v", id, got, ok, want, wantOK)
		}
	}
}
