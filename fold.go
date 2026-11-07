package learnermodel

import (
	"cmp"
	"slices"
	"strconv"
	"strings"
)

// The state is a fold of the evidence log (docs/architecture.md, section 3.5): Recompute folds the
// whole log from the empty state, and Append is one step of the fold over the cache, the same
// step. So the state an appended piece of evidence gives is the state recomputed from the log, bit
// for bit, as long as the cache is kept without loss, which the codec does.

// Append appends a piece of evidence to a state — the cache api keeps — against the graphs of the
// current versions of the modules, and gives the new state and the evidence with the start band
// the step took: the module's, when it has a state, and the evidence's own otherwise. api writes
// both in one transaction. The zero State is the empty one; the state Append gives is of
// StateKind and shares no memory with s, while the evidence it gives shares e's lists.
//
// The step reads only the graph of the evidence's module. A module whose entry says it has no
// current version moves nothing: its evidence waits in the log. Evidence with no answer moves
// nothing either, and gives a module no state.
//
// It refuses, with ErrInvalid:
//   - a state of a kind outside 1 to StateKind, unless it is the zero State, one not in its
//     canonical order, one with a skill of a stability below τ₀, and one EncodeState would
//     refuse;
//   - evidence of a kind below 1, evidence EncodeEvidence would refuse but for a later kind, and
//     evidence whose task ID names no module, whose task is not of a difficulty from 1 to 5, or
//     whose task band or start band is at a place below 1;
//   - a choice without its number of options or with fewer than two, and an item with an answer
//     and no slot;
//   - a module with no entry in g, and a graph whose module is not its entry's or whose lists are
//     not in the order NewGraph gives.
func Append(s State, e Evidence, g Graphs) (State, Evidence, error) {
	if err := checkInputState(s); err != nil {
		return State{}, Evidence{}, err
	}
	in, err := checkInput(e, g, nil)
	if err != nil {
		return State{}, Evidence{}, err
	}
	out := cloneState(s)
	e = out.append(e, in)
	return out, e, nil
}

// Recompute folds the evidence log, in the order of its numbers, from the empty state, against
// the graphs of the current versions of the modules: Append of each piece of evidence in turn. Its
// errors are Append's, with the index of the evidence in log, from 0. It checks the graph of each
// module once.
func Recompute(log []Evidence, g Graphs) (State, error) {
	s := State{Kind: StateKind}
	checked := map[string]bool{}
	for i, e := range log {
		in, err := checkInput(e, g, checked)
		if err != nil {
			return State{}, &detailError{err, "at index " + strconv.Itoa(i) + " of the log"}
		}
		s.append(e, in)
	}
	return s, nil
}

// input is evidence as the step reads it, once checked: its module, its items paired with their
// answers by place, and the module's entry in the graphs.
type input struct {
	module string
	items  []item
	graph  ModuleGraph
}

// append is the step of the fold on s, which it changes in place, and gives the evidence with the
// start band it took.
func (s *State) append(e Evidence, in input) Evidence {
	s.Kind = StateKind
	i, ok := slices.BinarySearchFunc(s.Modules, in.module, func(m ModuleState, id string) int { return strings.Compare(m.Module, id) })
	if ok {
		e.StartBand = s.Modules[i].StartBand
	}
	if !in.graph.Current {
		return e
	}
	if ok {
		s.Modules[i].step(e, in.items, &in.graph.Graph)
		return e
	}
	// θ needs no start of its own: a new module is in the trial series, whose θ comes from θ₀.
	ms := ModuleState{Module: in.module, StartBand: e.StartBand}
	if ms.step(e, in.items, &in.graph.Graph) {
		s.Modules = slices.Insert(s.Modules, i, ms)
	}
	return e
}

// moduleOf gives the module of a full ID, `<publisher>/<name>/<kind>/<ID>`: its first two parts.
func moduleOf(id string) (string, bool) {
	parts := strings.SplitN(id, "/", 4)
	if len(parts) < 4 || slices.Contains(parts, "") {
		return "", false
	}
	return parts[0] + "/" + parts[1], true
}

// checkInputState checks a state as the step needs it, beside what EncodeState checks: every skill
// has a stability of τ₀ or more, and a trial answer is one the step writes — a point on the scale,
// slots of known kinds, a verdict right or not yet and a positive weight — and the trial series is
// over once n ≥ 5.
func checkInputState(s State) error {
	if s.Kind == 0 && len(s.Modules) == 0 {
		return nil
	}
	if err := checkState(s); err != nil {
		return err
	}
	if err := canonicalState(s); err != nil {
		return invalid("a state " + err.Error())
	}
	for _, m := range s.Modules {
		if m.StartBand.Place < 1 {
			return invalid("module " + m.Module + " starts at band place " + itoa(m.StartBand.Place))
		}
		if len(m.Trial) > 0 && m.AnswerWeight >= trialWeight-weightTolerance {
			return invalid("module " + m.Module + " keeps trial answers past the trial series")
		}
		// A skill has a state only once it has been practised, and its stability is then τ₀ or more.
		for _, sk := range m.Skills {
			if !(sk.Stability >= stabilityFloor) {
				return invalid("skill " + sk.Skill + " of a stability below τ₀")
			}
		}
		for _, a := range m.Trial {
			what := "a trial answer of " + m.Module
			if err := checkTaskPoint(what, a.BandPlace, a.Difficulty); err != nil {
				return err
			}
			if err := checkChoices(a.Slots); err != nil {
				return err
			}
			if !isAnswer(TaskItem{Slots: a.Slots}, ItemAnswer{FirstVerdict: a.Verdict, HasFirstVerdict: true}) || len(a.Slots) == 0 {
				return invalid(what + " with no answer the step reads")
			}
			if !(a.Weight > 0) {
				return invalid(what + " of a weight that is not positive")
			}
		}
	}
	return nil
}

// checkInput checks evidence as the step needs it, and gives it as the step reads it. A module in
// checked has had its graph checked already, and checkInput marks the modules it checks there.
func checkInput(e Evidence, g Graphs, checked map[string]bool) (input, error) {
	t := e.Task
	if e.Kind < 1 {
		return input{}, invalid("evidence of kind " + itoa(e.Kind))
	}
	if err := checkEvidenceContent(e); err != nil {
		return input{}, err
	}
	module, ok := moduleOf(t.ID)
	if !ok {
		return input{}, invalid("task " + shownBytes(t.ID) + " names no module")
	}
	if err := checkTaskPoint("task "+t.ID, t.Band.Place, t.Difficulty); err != nil {
		return input{}, err
	}
	if e.StartBand.Place < 1 {
		return input{}, invalid("evidence of task " + t.ID + " starts at band place " + itoa(e.StartBand.Place))
	}
	items := pairItems(e)
	for _, it := range items {
		if err := checkChoices(it.task.Slots); err != nil {
			return input{}, err
		}
		if len(it.task.Slots) == 0 && isAnswer(it.task, it.answer) {
			return input{}, invalid("item " + itoa(it.task.Place) + " of task " + t.ID + " has an answer and no slot")
		}
	}
	mg, ok := g[module]
	if !ok {
		return input{}, invalid("module " + module + " has no entry in the graphs")
	}
	if mg.Current && !checked[module] {
		if err := checkGraph(module, mg.Graph); err != nil {
			return input{}, err
		}
		if checked != nil {
			checked[module] = true
		}
	}
	return input{module: module, items: items, graph: mg}, nil
}

func checkTaskPoint(what string, place, difficulty int32) error {
	if place < 1 {
		return invalid(what + " in band place " + itoa(place))
	}
	if difficulty < 1 || difficulty > 5 {
		return invalid(what + " of difficulty " + itoa(difficulty))
	}
	return nil
}

func checkChoices(slots []Slot) error {
	for _, s := range slots {
		if s.Check == CheckKindChoice && (!s.HasOptions || s.Options < 2) {
			return invalid("a choice of " + itoa(s.Options) + " options")
		}
	}
	return nil
}

// checkGraph checks what the step relies on in a graph: its module, and skills and the bands
// each is taught at in their order, as NewGraph gives them.
func checkGraph(module string, g Graph) error {
	if g.Module != module {
		return invalid("the graph of module " + g.Module + " in the entry of " + module)
	}
	if err := ascending("skills of the graph of "+module, g.Skills, func(a, b Skill) int { return strings.Compare(a.ID, b.ID) }); err != nil {
		return invalid(err.Error())
	}
	for _, s := range g.Skills {
		if err := ascending("bands skill "+s.ID+" is taught at", s.BandsTaught, cmp.Compare[int32]); err != nil {
			return invalid(err.Error())
		}
	}
	return nil
}

// cloneState copies a state deeply, so that the copy shares no memory with it.
func cloneState(s State) State {
	out := State{Kind: s.Kind}
	for _, m := range s.Modules {
		m.Trial = slices.Clone(m.Trial)
		for i := range m.Trial {
			m.Trial[i].Slots = slices.Clone(m.Trial[i].Slots)
		}
		m.Traps = slices.Clone(m.Traps)
		m.Skills = slices.Clone(m.Skills)
		for i := range m.Skills {
			m.Skills[i].Traps = slices.Clone(m.Skills[i].Traps)
		}
		out.Modules = append(out.Modules, m)
	}
	return out
}
