package learnermodel

// The types the library speaks with api, engine and sim (R10 in docs/decisions.md;
// docs/architecture.md, section 3.7). They are values: a field that may be absent by its meaning
// is a value with a flag beside it, never a pointer, so that a state is copied by assignment and
// appending shares nothing with the state it started from.

// Grade is a learner's school grade, or no grade: Known is false for a learner without one, and
// Value is then 0. "No grade" is a case of its own, not grade 0.
type Grade struct {
	Value int32
	Known bool
}

// BandRef names a band of a module's difficulty scale by its place in the `bands` of the module
// version, from 1, and by its ID in the manifest. The model reads the place; the ID is kept.
type BandRef struct {
	Place int32
	ID    string
}

// Band is a band of a module's difficulty scale, in the order of the manifest's `bands`.
type Band struct {
	Place  int32
	ID     string
	Grades []int32 // the grades the band names, ascending; empty in a module without grades
}

// Assessability is how the platform can know a learner's command of a skill. Only a skill with
// AssessabilityVerdict has a state.
type Assessability int32

const (
	AssessabilityUnspecified Assessability = iota
	AssessabilityVerdict
	AssessabilityRubric
	AssessabilityNone
)

// Skill is a skill of a module's graph, by its full ID.
type Skill struct {
	ID            string
	Assessability Assessability
	Teaches       []string // the anchor skills it teaches, `core.<name>`, ordered byte by byte
	// BandsTaught is the places of the bands where the skill is taught: the bands of the tasks
	// where it is a primary skill of the task or of an item, ascending. NewGraph derives it.
	BandsTaught   []int32
	Prerequisites []string // skills of the module by full ID and anchor skills; the engine reads them
	Tier          int32    // the skill's tier, `skill_tiers` of the server part; the engine reads it
}

// RetiredSkill is an entry of `retired_skills`: a skill of an earlier version of the module and
// the skills of this one that replaced it. An empty ReplacedBy is `[]`: no skill took over.
type RetiredSkill struct {
	ID         string
	ReplacedBy []string
}

// Graph is the graph of one version of a module, as the model and the engine read it.
type Graph struct {
	Module  string
	Bands   []Band         // by place, from 1
	Skills  []Skill        // ordered by ID byte by byte
	Retired []RetiredSkill // ordered by ID byte by byte
}

// ModuleGraph is the entry of a module in Graphs: the graph of its current version, or no
// current version, when Current is false and Graph is empty (docs/architecture.md, section 3.5).
type ModuleGraph struct {
	Graph   Graph
	Current bool
}

// Graphs gives each module the graph of its current version, by the module's full ID. The model
// looks modules up in it and never ranges over it, so its order cannot matter.
type Graphs map[string]ModuleGraph

// CheckKind is how the answer in a slot is checked: the kinds that give a verdict. A value this
// version does not know is kept as its number, and its item reads as one with no answer.
type CheckKind int32

const (
	CheckKindUnspecified CheckKind = iota
	CheckKindNumber
	CheckKindFraction
	CheckKindChoice
	CheckKindPredicate
)

// Slot is an answer slot of an item, as the model needs it for the guessing of the item.
type Slot struct {
	Check CheckKind
	// Options is the number of options of a choice, and is meaningful only when HasOptions.
	Options    int32
	HasOptions bool
}

// Provenance is where a key or the skill tags of a task came from (concepts 9.4 of spec).
type Provenance int32

const (
	ProvenanceUnspecified Provenance = iota
	ProvenanceAuthor
	ProvenanceSolver
	ProvenanceAdult
	ProvenanceAI
)

// TaskItem is an item of a task, by its place, from 1.
type TaskItem struct {
	Place int32
	// OwnSkills tells whether the item names skills of its own, Primary and Supporting, which take
	// the place of the task's for it.
	OwnSkills  bool
	Primary    []string
	Supporting []string
	Slots      []Slot // in their order
}

// Task describes a task, the same for evidence and for a candidate task of the engine.
type Task struct {
	ID         string // the full ID of the task; of a template, of the template
	Template   bool
	Band       BandRef
	Difficulty int32 // from 1 to 5
	Primary    []string
	Supporting []string // full IDs or anchor skills, `core.<name>`
	// KeyProvenance and SkillTagsProvenance are where the key and the skill tags came from.
	KeyProvenance       Provenance
	SkillTagsProvenance Provenance
	Items               []TaskItem // by place
}
