package learnermodel

// StateKind is the kind of state this version writes. It reads every kind up to it. Kind 2 adds
// a module's FailuresInARow and TrialEndedByAdult, the state of the guard and of the adult's correction,
// which kind 1 does not record.
const StateKind = 2

// guardKind is the kind of state that adds FailuresInARow and TrialEndedByAdult of a module, fields
// 9 and 10: the writer writes them from this kind on, the reader requires them there and refuses
// them below.
const guardKind = 2

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

// SkillState is the state of one skill of the module whose assessability is verdict.
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
	// still needs to know. Evidence that moves nothing leaves it as it is.
	LastInTrial bool
	// FailuresInARow is the count of the guard: the pieces of evidence in a row, past the
	// trial series, on which the model expected the learner to know at least half of the answers,
	// above the guessing floor, and fewer than half of the items with an answer were solved at the
	// end of the attempt. The third lowers the level by a band; the count goes on until a piece of
	// evidence that is no such failure.
	FailuresInARow int32
	// TrialEndedByAdult tells whether an adult's correction of the level ended the trial series
	// before its five answers: the series is then over whatever the weight of answers.
	TrialEndedByAdult bool
}

// State is the state of one learner: the fold of the evidence log.
type State struct {
	Kind    int32
	Modules []ModuleState // ordered by module ID byte by byte
}
