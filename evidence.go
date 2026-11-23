package learnermodel

// EvidenceKind is the kind of evidence this version writes. It reads every kind up to it, and a
// later one too, skipping what it does not know. Kind 2 adds LeadingQuestionOnRequest, which
// kind 1 does not record, and kind 3 LevelSetByAdult, which kinds 1 and 2 do not.
const EvidenceKind = 3

// leadingQuestionKind is the kind that adds LeadingQuestionOnRequest, field 20: the writer writes
// it from this kind on, the reader requires it there and refuses it below.
const leadingQuestionKind = 2

// correctionKind is the kind that adds LevelSetByAdult, field 21, as leadingQuestionKind adds
// field 20.
const correctionKind = 3

// OutcomeKind is how an attempt ended. A value this version does not know is kept as its number and
// weighs nothing.
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

// UnverdictedReason is why a submission gave no verdict.
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

// Evidence is the evidence of one finished attempt: its observations, not the model's conclusions,
// so that a later version recomputes the state from it.
type Evidence struct {
	Kind      int32
	Task      Task
	StartBand BandRef
	Outcome   Outcome
	// PracticeTime is the time of the attempt's last action, in whole seconds of Unix time, UTC.
	PracticeTime       int64
	TechniqueOnStep    bool
	TechniqueOnRequest bool
	// LeadingQuestionOnRequest is whether a hint the learner asked for gave the task's leading
	// question, which a hint gives only before the first wrong answer. Kind 1 does not record it, and
	// it is false there.
	LeadingQuestionOnRequest bool
	ExplainAgainCount        int32
	SolutionSeenBefore       bool
	HasSolutionSeenFact      bool // false for a template, whose instances differ
	// Items holds what happened to each item; the slots and the item's own skills are in Task.
	Items []ItemAnswer
	Traps []FoundTrap // ordered by ID byte by byte
	// LevelSetByAdult tells that the evidence is no attempt but an adult's correction of the
	// module's level, as LevelCorrection gives it. Kinds 1 and 2 do not record it, and it is
	// false there.
	LevelSetByAdult bool
}

// levelKind is the kind of entry the ID of a correction's task names in place of a task's: no
// kind of the format's, so that it names no entry of a module.
const levelKind = "/level/"

// LevelCorrection gives the evidence of an adult's correction of the level of a module, by its
// full ID: the band the adult set, by its place and ID in the module's current
// version, at a time in whole seconds of Unix time, UTC, with the start band api gives any evidence
// of the module, by BandByGrade. Append sets the module's level in that band in one step, or keeps
// it when it stands there already, ends the trial series and starts the guard's count anew, so
// that a repeat changes nothing; a band above the version's highest is the highest. Its task is
// no task: its ID is `<module>/level/<band ID>`, its band the band set and its difficulty 3, the
// band's middle, where the level goes, and the adult its provenance; it has no skills and no
// items, so that a version before kind 3, which does not know LevelSetByAdult, reads it as
// evidence with no answer, which moves nothing.
func LevelCorrection(module string, band, start BandRef, at int64) Evidence {
	return Evidence{
		Kind: EvidenceKind,
		Task: Task{
			ID: module + levelKind + band.ID, Band: band, Difficulty: 3,
			KeyProvenance: ProvenanceAdult, SkillTagsProvenance: ProvenanceAdult,
		},
		StartBand:       start,
		PracticeTime:    at,
		LevelSetByAdult: true,
	}
}
