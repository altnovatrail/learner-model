package learnermodel

import (
	"cmp"
	"slices"
	"strings"
)

// Reading the state at a time t₀ (docs/architecture.md, sections 3.1–3.4 and 3.7): what the
// engine, api and the skill map read of the model. Reading writes nothing into the cache: what
// depends on the time, the learner's grade and the modules of the engine is worked out here, from
// the state as the fold left it (R08 of api; the amendment of 2026-11-02 to R80 of spec).

// ReviewDate is when a mastered skill is due for review: At, in whole seconds of Unix time, UTC,
// or Already, when it is due whatever the time.
type ReviewDate struct {
	At      int64
	Already bool
}

// SkillReadout is what reading gives of a skill with practice.
type SkillReadout struct {
	Skill         string
	EvidenceCount int64
	LastPractice  int64
	Estimate      float64 // θ + δ
	Uncertainty   float64 // v
	Retention     float64 // R at t₀
	Traps         []TrapCount
	// For a mastered skill: its mastery band — the ID from the graph of the module's current
	// version, empty without one —, the cautious probability with retention at t₀, R · q, and
	// the review date. Zero while the skill is not mastered.
	Mastered            bool
	MasteryBand         BandRef
	CautiousProbability float64
	Review              ReviewDate
}

// ModuleReadout is what reading gives of a module.
type ModuleReadout struct {
	Module string
	// HasEvidence is false for a module of the engine with no evidence yet, whose StartBand is
	// then the starting point by the learner's grade, or the zero BandRef when the module's
	// current version has no bands.
	HasEvidence bool
	StartBand   BandRef
	// TrialOver tells whether the module's trial series is over, and LastInTrial whether the last
	// evidence that moved the module was still in it, the evidence that ended it too (R104 of
	// spec). Both are false before the module's first evidence.
	TrialOver   bool
	LastInTrial bool
	// Traps are the traps the program found that have no skill with a state, counted at the module
	// (docs/architecture.md, section 3.3).
	Traps  []TrapCount
	Skills []SkillReadout // ordered by skill ID byte by byte
}

// AnchorStatus is whether a learner has an anchor skill (docs/architecture.md, section 3.4).
type AnchorStatus int32

const (
	AnchorUnspecified AnchorStatus = iota
	AnchorHas
	AnchorHasNot
	AnchorUnknown
)

// AnchorReason is why an anchor skill has its status: for AnchorHasNot, not mastered or no
// evidence; for AnchorUnknown, no one teaches it or only skills that may not be mastered do —
// without a verdict, or with no band they are taught at, as in a module with no tasks —, which
// AnchorReasonNoVerdictTeacher names.
type AnchorReason int32

const (
	AnchorReasonNone AnchorReason = iota
	AnchorReasonNotMastered
	AnchorReasonNoEvidence
	AnchorReasonNoTeacher
	AnchorReasonNoVerdictTeacher
)

// TeacherReadout is a skill that teaches an anchor skill, as reading gives it with the anchor:
// the highest band it is taught at — Place 0 when there is none —, and, when Practised, its
// estimate, mastery, retention at t₀ and review date, as its SkillReadout has them.
type TeacherReadout struct {
	Skill       string
	HighestBand BandRef
	Practised   bool
	Estimate    float64
	Mastered    bool
	MasteryBand BandRef
	Retention   float64
	Review      ReviewDate
}

// AnchorReadout is the status of an anchor skill, with its reason and the skills that teach it.
type AnchorReadout struct {
	Anchor   string
	Status   AnchorStatus
	Reason   AnchorReason
	Teachers []TeacherReadout // ordered by skill ID byte by byte
}

// Readout is what reading the state at t₀ gives. It keeps what Task needs: a copy of the state,
// the graphs, t₀ and the grade.
type Readout struct {
	Modules []ModuleReadout // ordered by module ID byte by byte
	Anchors []AnchorReadout // ordered by anchor ID byte by byte

	state   State
	graphs  Graphs
	checked []checkedModule // the modules whose current graphs Read checked, by ID
	at      int64
	grade   Grade
}

// checkedModule is a module whose current graph Read checked, with what Task reads of it the same
// for every candidate: where it starts when it has no state, and R at t₀ of each skill of the
// graph with a verdict, by ID, as the step works it out — of the skill's state, or 1 for a skill
// with none. Task only reads them, so that it may be called while another call runs.
type checkedModule struct {
	id         string
	start      BandRef
	hasStart   bool
	retentions []skillRetention
}

// ItemProbability is the probability of a right answer to an item of a candidate task at t₀:
// Probability, when HasVerdict; an item without a verdict has none.
type ItemProbability struct {
	Place       int32
	Probability float64
	HasVerdict  bool
}

// TaskReadout is what Readout.Task gives of a candidate task: the P of each item, by place.
type TaskReadout struct {
	Items []ItemProbability
}

// anchorPrefix opens the ID of an anchor skill of the shared taxonomy, `core.<name>`.
const anchorPrefix = "core."

// Read reads a state — the cache api keeps — at the time at, in whole seconds of Unix time, UTC,
// for a learner of a grade, against the graphs of the current versions of the modules and the
// modules of the engine: those connected and fit for the learner, whose graphs the engine gets.
// g gives an entry to every module of the state and of the engine; the Readout keeps g for Task,
// and the caller changes neither while it reads.
//
// Modules are the modules of the state and of the engine. A module with a state gives its start
// band, the two flags of the trial series and its skills; a module of the engine with no state
// and a current version gives the starting point by the grade, BandByGrade, or none when that
// version has no bands, as a module with no tasks may have none; one with neither is left out.
// Anchors are the anchor skills a graph of a module of the engine names as a prerequisite or a
// skill it teaches, each with its status (docs/architecture.md, section 3.4): held, when a skill
// that teaches it is mastered in any of the modules with a current version; otherwise not held,
// when a skill of the engine's modules that may be mastered teaches it — a skill with a verdict
// and a band it is taught at —, not mastered when one of them has answers, no evidence
// otherwise; and unknown otherwise: no one teaches it, or only skills that may not be mastered
// do. A module with no current version gives no graph, so its mastery is not seen.
//
// It refuses, with ErrInvalid, a state Append would refuse, a module of the engine twice or not
// named by the full ID of a module, a module with no entry in g, a graph Append would refuse, and
// bands, other than none, that BandByGrade would refuse.
func Read(s State, g Graphs, at int64, grade Grade, engineModules []string) (Readout, error) {
	if err := checkInputState(s); err != nil {
		return Readout{}, err
	}
	engine := sortedCopy(engineModules)
	if err := idList("module of the engine", engine); err != nil {
		return Readout{}, err
	}
	for _, id := range engine {
		if !isModuleID(id) {
			return Readout{}, invalid("module of the engine " + shownBytes(id) + " is no full ID of a module")
		}
	}
	r := Readout{state: cloneState(s), graphs: g, at: at, grade: grade}
	mods, err := r.readModules(engine)
	if err != nil {
		return Readout{}, err
	}
	for _, m := range mods {
		mr, ok, err := m.readout(at, grade)
		if err != nil {
			return Readout{}, err
		}
		if ok {
			r.Modules = append(r.Modules, mr)
		}
		// A module with no state and a current version starts where its reading said: the same
		// call, which cannot refuse here, as readout has just made it.
		if c, read := r.checkedModule(m.id); read && m.state == nil {
			c.start, c.hasStart, _ = startingPoint(grade, m.graph.Graph.Bands)
		}
	}
	r.Anchors = anchors(mods, at)
	return r, nil
}

// Task gives the probability of a right answer to each item of a candidate task at the time of
// reading, by place: the P the step would take of evidence at that time (docs/architecture.md,
// sections 3.1 and 3.2), by the same function — with the retention of the item's skills at t₀,
// and without it, R̄ = 1, while the module's trial series lasts. A module with no state starts at
// the starting point by the grade. An item with no slot, or with one of a kind of check this
// version does not know, has no verdict and so no P. The provenance of the task's key and skill
// tags does not change P: evidence that weighs nothing by it moves nothing, but the chance of a
// right answer is the same.
//
// It refuses, with ErrInvalid, a task the codec would refuse, one whose ID names no module, one
// not of a difficulty from 1 to 5 or at a band place below 1, a choice without its number of
// options or with fewer than two, a module with no entry in the graphs or with no current
// version, a graph Append would refuse, bands BandByGrade would, and a module with no state whose
// current version has no bands: a task of it would have none to be in.
//
// It copies nothing it need not: tags and items already in their order are read as they are, and
// R at t₀ and the start of a module Read has read come from Read (docs/architecture.md, section
// 3.7; R29).
func (r Readout) Task(t Task) (TaskReadout, error) {
	if err := checkTaskTags(t); err != nil {
		return TaskReadout{}, err
	}
	items := t.Items
	if !slices.IsSortedFunc(items, byPlace) {
		items = slices.Clone(items)
		slices.SortFunc(items, byPlace)
	}
	for i := 1; i < len(items); i++ {
		if items[i].Place == items[i-1].Place {
			return TaskReadout{}, invalid("item place twice")
		}
	}
	module, ok := moduleOf(t.ID)
	if !ok {
		return TaskReadout{}, invalid("task " + shownBytes(t.ID) + " names no module")
	}
	if err := checkTaskPoint("task ", t.ID, t.Band.Place, t.Difficulty); err != nil {
		return TaskReadout{}, err
	}
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
	c, read := r.checkedModule(module)
	if !read {
		if err := checkGraph(module, mg.Graph); err != nil {
			return TaskReadout{}, err
		}
	}
	var ms ModuleState
	if i, ok := r.state.moduleIndex(module); ok {
		ms = r.state.Modules[i]
	} else {
		var start BandRef
		var ok bool
		if read {
			start, ok = c.start, c.hasStart
		} else {
			var err error
			if start, ok, err = startingPoint(r.grade, mg.Graph.Bands); err != nil {
				return TaskReadout{}, err
			}
		}
		if !ok {
			return TaskReadout{}, invalid("task " + t.ID + " of module " + module + ", whose current version has no bands")
		}
		ms = newModule(module, start)
	}
	trial := ms.inTrial()
	// Every skill transfer gives is one of the graph with a verdict, so that R of each is found
	// and none is inserted; were one, the clipped slice would be copied, never written.
	var rs []skillRetention
	if read {
		rs = slices.Clip(c.retentions)
	}
	out := TaskReadout{Items: make([]ItemProbability, 0, len(items))}
	for _, it := range items {
		ip := ItemProbability{Place: it.Place}
		if len(it.Slots) > 0 && knownChecks(it.Slots) {
			ip.HasVerdict = true
			ip.Probability = ms.chanceOf(t, it, &mg.Graph, &rs, nil, r.at, trial).p
		}
		out.Items = append(out.Items, ip)
	}
	return out, nil
}

// byPlace orders the items of a task by their place.
func byPlace(a, b TaskItem) int { return cmp.Compare(a.Place, b.Place) }

// checkedModule finds a module whose current graph Read checked, by ID.
func (r *Readout) checkedModule(id string) (*checkedModule, bool) {
	i, ok := slices.BinarySearchFunc(r.checked, id, func(c checkedModule, id string) int { return strings.Compare(c.id, id) })
	if !ok {
		return nil, false
	}
	return &r.checked[i], true
}

// moduleRead is a module reading reads: its entry in the graphs, its state, if it has one, and
// whether it is a module of the engine.
type moduleRead struct {
	id     string
	graph  ModuleGraph
	state  *ModuleState
	engine bool
}

// currentGraph is the graph of the module's current version, or nil without one.
func (m moduleRead) currentGraph() *Graph {
	if !m.graph.Current {
		return nil
	}
	return &m.graph.Graph
}

// readModules gives the modules of the state and of the engine, ordered by ID, each once, with
// their entries in the graphs, each graph checked.
func (r *Readout) readModules(engine []string) ([]moduleRead, error) {
	ids := slices.Clone(engine)
	for _, m := range r.state.Modules {
		ids = append(ids, m.Module)
	}
	slices.Sort(ids)
	ids = slices.Compact(ids)
	mods := make([]moduleRead, 0, len(ids))
	for _, id := range ids {
		mg, ok := r.graphs[id]
		if !ok {
			return nil, invalid("module " + id + " has no entry in the graphs")
		}
		m := moduleRead{id: id, graph: mg}
		if i, ok := r.state.moduleIndex(id); ok {
			m.state = &r.state.Modules[i]
		}
		if mg.Current {
			if err := checkGraph(id, mg.Graph); err != nil {
				return nil, err
			}
			ms := m.state
			if ms == nil {
				ms = &ModuleState{}
			}
			r.checked = append(r.checked, checkedModule{id: id, retentions: ms.retentions(&mg.Graph, r.at)})
		}
		_, m.engine = slices.BinarySearch(engine, id)
		mods = append(mods, m)
	}
	return mods, nil
}

// retentions gives R at a time of each skill of a graph with a verdict, by ID, by retention, as the
// step works it out: of the skill's state, or 1 for a skill with none.
func (ms *ModuleState) retentions(g *Graph, at int64) []skillRetention {
	rs := make([]skillRetention, 0, len(g.Skills))
	for _, s := range g.Skills {
		if s.Assessability == AssessabilityVerdict {
			ms.retention(&rs, s.ID, at)
		}
	}
	return rs
}

// readout gives what reading gives of a module, and false for a module with neither a state nor
// a current version.
func (m moduleRead) readout(at int64, grade Grade) (ModuleReadout, bool, error) {
	g := m.currentGraph()
	if m.state == nil {
		if g == nil {
			return ModuleReadout{}, false, nil
		}
		// With no place to start at, the engine still gets the module's reading, for its lessons
		// and anchor skills.
		start, _, err := startingPoint(grade, g.Bands)
		if err != nil {
			return ModuleReadout{}, false, err
		}
		return ModuleReadout{Module: m.id, StartBand: start}, true, nil
	}
	ms := m.state
	out := ModuleReadout{
		Module:      m.id,
		HasEvidence: true,
		StartBand:   ms.StartBand,
		TrialOver:   !ms.inTrial(),
		LastInTrial: ms.LastInTrial,
		Traps:       slices.Clone(ms.Traps),
		Skills:      make([]SkillReadout, 0, len(ms.Skills)),
	}
	for _, sk := range ms.Skills {
		out.Skills = append(out.Skills, ms.skillReadout(sk, g, at))
	}
	return out, true, nil
}

// skillReadout is what reading at t₀ gives of a skill with a state (docs/architecture.md, section
// 3.2): θ + δ, v and R at t₀, and for a mastered skill R · q at t₀ and the review date, by the
// functions of model.go. g names the mastery band; it is nil without a current version.
func (ms *ModuleState) skillReadout(sk SkillState, g *Graph, at int64) SkillReadout {
	out := SkillReadout{
		Skill:         sk.Skill,
		EvidenceCount: sk.EvidenceCount,
		LastPractice:  sk.LastPractice,
		Estimate:      ms.Level + sk.Correction,
		Uncertainty:   uncertainty(ms.AnswerWeight, sk.AnswerWeight),
		Retention:     retentionAt(sk, at),
		Traps:         slices.Clone(sk.Traps),
	}
	if sk.Mastered {
		q := masteryChance(ms.Level, ms.AnswerWeight, sk)
		out.Mastered = true
		out.MasteryBand = bandAt(g, sk.MasteryBand)
		out.CautiousProbability = float64(out.Retention * q)
		out.Review = reviewDate(q, sk.Stability, sk.LastPractice)
	}
	return out
}

// bandAt names the band at a place by the graph: its ID, when the graph has a band there.
func bandAt(g *Graph, place int32) BandRef {
	b := BandRef{Place: place}
	if g != nil && place >= 1 && int(place) <= len(g.Bands) {
		b.ID = g.Bands[place-1].ID
	}
	return b
}

// anchors gives the status of each anchor skill a graph of a module of the engine names, ordered
// by ID.
func anchors(mods []moduleRead, at int64) []AnchorReadout {
	var ids []string
	for _, m := range mods {
		g := m.currentGraph()
		if !m.engine || g == nil {
			continue
		}
		for _, sk := range g.Skills {
			for _, id := range slices.Concat(sk.Prerequisites, sk.Teaches) {
				if strings.HasPrefix(id, anchorPrefix) {
					ids = append(ids, id)
				}
			}
		}
	}
	slices.Sort(ids)
	ids = slices.Compact(ids)
	out := make([]AnchorReadout, 0, len(ids))
	for _, id := range ids {
		out = append(out, anchorOf(id, mods, at))
	}
	return out
}

// anchorOf gives the status of an anchor skill from the skills that teach it in the modules with
// a current version (docs/architecture.md, section 3.4).
func anchorOf(anchor string, mods []moduleRead, at int64) AnchorReadout {
	out := AnchorReadout{Anchor: anchor}
	var held, engineTeaches, masterable, practised bool
	for _, m := range mods {
		g := m.currentGraph()
		if g == nil {
			continue
		}
		for _, sk := range g.Skills {
			if !slices.Contains(sk.Teaches, anchor) {
				continue
			}
			t := TeacherReadout{Skill: sk.ID}
			if n := len(sk.BandsTaught); n > 0 {
				t.HighestBand = bandAt(g, sk.BandsTaught[n-1])
			}
			if m.state != nil {
				if i, ok := m.state.skillIndex(sk.ID); ok {
					sr := m.state.skillReadout(m.state.Skills[i], g, at)
					t.Practised = true
					t.Estimate, t.Mastered, t.MasteryBand = sr.Estimate, sr.Mastered, sr.MasteryBand
					t.Retention, t.Review = sr.Retention, sr.Review
				}
			}
			held = held || t.Mastered
			if m.engine {
				engineTeaches = true
				// Only a skill with a verdict and a band it is taught at may be mastered.
				if sk.Assessability == AssessabilityVerdict && len(sk.BandsTaught) > 0 {
					masterable = true
					practised = practised || t.Practised
				}
			}
			out.Teachers = append(out.Teachers, t)
		}
	}
	// Full IDs of skills of different modules need not sort as their modules do.
	slices.SortFunc(out.Teachers, func(a, b TeacherReadout) int { return strings.Compare(a.Skill, b.Skill) })
	switch {
	case held:
		out.Status = AnchorHas
	case masterable && practised:
		out.Status, out.Reason = AnchorHasNot, AnchorReasonNotMastered
	case masterable:
		out.Status, out.Reason = AnchorHasNot, AnchorReasonNoEvidence
	case engineTeaches:
		out.Status, out.Reason = AnchorUnknown, AnchorReasonNoVerdictTeacher
	default:
		out.Status, out.Reason = AnchorUnknown, AnchorReasonNoTeacher
	}
	return out
}

// isModuleID tells whether an ID is the full ID of a module, `<publisher>/<name>`.
func isModuleID(id string) bool {
	publisher, name, ok := strings.Cut(id, "/")
	return ok && publisher != "" && name != "" && !strings.Contains(name, "/")
}

// startingPoint gives where a module with no state starts for a learner of a grade, BandByGrade,
// and false, with the zero BandRef, when the module's current version has no bands: a module with
// no tasks may have none (the format's section 11), and then no place to start at either.
func startingPoint(grade Grade, bands []Band) (BandRef, bool, error) {
	if len(bands) == 0 {
		return BandRef{}, false, nil
	}
	start, err := BandByGrade(grade, bands)
	if err != nil {
		return BandRef{}, false, err
	}
	return start, true, nil
}
