package learnermodel

// EvidenceKind is the kind of evidence this version writes. It reads every kind up to it, and a
// later one too, skipping what it does not know (docs/architecture.md, section 3.5).
const EvidenceKind = 1

// OutcomeKind is how an attempt ended (concepts 11 of spec). A value this version does not know
// is kept as its number and weighs nothing.
type OutcomeKind int32

const (
	OutcomeUnspecified OutcomeKind = iota
	OutcomeFirstSubmission
	OutcomeAfterHint
	OutcomeWithSketch
	OutcomeSolutionShown
	OutcomeLeftUnfinished
	OutcomeUnchecked
)

// Outcome is the outcome of an attempt, with the rung of the hint after which it ended, when
// HasRung: from 1, for OutcomeAfterHint.
type Outcome struct {
	Kind    OutcomeKind
	Rung    int32
	HasRung bool
}

// Verdict is an item's verdict: right or not yet. A value this version does not know is kept as
// its number, and its item reads as one with no answer.
type Verdict int32

const (
	VerdictUnspecified Verdict = iota
	VerdictCorrect
	VerdictNotYet
)

// UnverdictedReason is why a submission gave no verdict (concepts 9.1 of spec).
type UnverdictedReason int32

const (
	UnverdictedUnspecified UnverdictedReason = iota
	UnverdictedUnreadableBlank
	UnverdictedUnreadableOther
	UnverdictedUnchecked
)

// ItemAnswer is what happened to one item in the attempt, by its place.
type ItemAnswer struct {
	Place     int32
	Submitted bool
	// Unverdicted is why each submission before the first verdict gave none, in their order.
	Unverdicted []UnverdictedReason
	// FirstVerdict is the item's first verdict, when HasFirstVerdict, and WrongBeforeFirst the
	// number of wrong answers in the attempt before the submission that gave it.
	FirstVerdict     Verdict
	HasFirstVerdict  bool
	WrongBeforeFirst int32
	// FinalVerdict is the item's verdict at the end, when HasFinalVerdict.
	FinalVerdict    Verdict
	HasFinalVerdict bool
	// WrongBeforeFirstCorrect is the number of wrong answers in the attempt before the item's
	// first correct verdict, when HasFirstCorrect.
	WrongBeforeFirstCorrect int32
	HasFirstCorrect         bool
}

// FoundTrap is a trap the program found in the attempt by a wrong value.
type FoundTrap struct {
	ID         string
	Skills     []string // its skills in the trap catalog of the attempt's version, by full ID
	ItemPlaces []int32  // where it was found, ascending
}

// Evidence is the evidence of one finished attempt (docs/architecture.md, section 3.5): its
// observations, not the model's conclusions.
type Evidence struct {
	Kind      int32
	Task      Task
	StartBand BandRef
	Outcome   Outcome
	// PracticeTime is the time of the attempt's last action, in whole seconds of Unix time, UTC.
	PracticeTime        int64
	TechniqueOnStep     bool
	TechniqueOnRequest  bool
	ExplainAgainCount   int32
	SolutionSeenBefore  bool
	HasSolutionSeenFact bool // false for a template, whose instances differ
	// Items holds what happened to each item; the slots and the item's own skills are in Task.
	Items []ItemAnswer
	Traps []FoundTrap // ordered by ID byte by byte
}
