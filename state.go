package learnermodel

// StateKind is the kind of state this version writes.
const StateKind = 1

// TrapCount is a trap the program found, with the number of pieces of evidence it was found in
// and the practice time of the latest.
type TrapCount struct {
	ID    string
	Count int64
	Last  int64
}

// TrialAnswer is one answer of the trial series, as its likelihood needs it.
type TrialAnswer struct {
	BandPlace  int32
	Difficulty int32
	Slots      []Slot
	Verdict    Verdict
	Weight     float64 // u, the answer's weight in its series
}

// SkillState is the state of one skill of the module whose assessability is verdict
// (docs/architecture.md, sections 3.1–3.3).
type SkillState struct {
	Skill         string
	Correction    float64 // δ
	AnswerWeight  float64 // m
	EvidenceCount int64
	LastPractice  int64   // whole seconds of Unix time, UTC
	Stability     float64 // τ, in days
	Traps         []TrapCount
	// MasteryBand is the place of the mastery band, when Mastered.
	MasteryBand    int32
	Mastered       bool
	FailuresInARow int32
}

// ModuleState is the state of a learner in one module.
type ModuleState struct {
	Module       string
	Level        float64 // θ
	AnswerWeight float64 // n
	StartBand    BandRef
	Trial        []TrialAnswer // in the order of the log, while the trial series lasts
	Traps        []TrapCount   // the traps with no skill that has a state
	Skills       []SkillState  // ordered by skill ID byte by byte
	// LastInTrial tells whether the last piece of evidence that moved the module was still in the
	// trial series, the one that ended it too: that one clears the series' answers, and the engine
	// still needs to know (R104 of spec). Evidence that moves nothing leaves it as it is.
	LastInTrial bool
}

// State is the state of one learner: the fold of the evidence log.
type State struct {
	Kind    int32
	Modules []ModuleState // ordered by module ID byte by byte
}
