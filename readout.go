package learnermodel

// The output of reading the state at a time t₀ (docs/architecture.md, section 3.7, the table of
// operations). These are the types alone: Read and Readout.Task compute them from T20 on.

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
	// For a mastered skill: its mastery band, the cautious probability with retention at t₀ and
	// the review date. Zero while the skill is not mastered.
	Mastered            bool
	MasteryBand         int32
	CautiousProbability float64
	Review              ReviewDate
}

// ModuleReadout is what reading gives of a module.
type ModuleReadout struct {
	Module string
	// HasEvidence is false for a module with no evidence yet, whose StartBand is then the starting
	// point by the learner's grade.
	HasEvidence bool
	StartBand   BandRef
	Skills      []SkillReadout // ordered by skill ID byte by byte
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
// evidence; for AnchorUnknown, no one teaches it or only skills without a verdict do.
type AnchorReason int32

const (
	AnchorReasonNone AnchorReason = iota
	AnchorReasonNotMastered
	AnchorReasonNoEvidence
	AnchorReasonNoTeacher
	AnchorReasonNoVerdictTeacher
)

// TeacherReadout is a skill that teaches an anchor skill, as reading gives it with the anchor.
type TeacherReadout struct {
	Skill       string
	Estimate    float64
	Mastered    bool
	MasteryBand int32
	HighestBand int32 // the highest band the skill is taught at
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

// Readout is what reading the state at t₀ gives.
type Readout struct {
	Modules []ModuleReadout // ordered by module ID byte by byte
	Anchors []AnchorReadout // ordered by anchor ID byte by byte
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
