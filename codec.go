package learnermodel

import (
	"cmp"
	"errors"
	"slices"
	"strings"

	"github.com/altnovatrail/learner-model/internal/wire"
)

// The codec of the evidence and the state: protobuf's binary form by the schema in
// proto/novatrail/learnermodel/v1/, written by this code alone, on the standard library without the
// protobuf runtime. One value gives one run of bytes: the fields go in ascending order of their
// numbers, a scalar with a value is written even when it is 0, lists of scalars are packed, and the
// elements of a repeated field go in their canonical order, which the writer sets itself.
//
// The reader is strict, as only this writer writes: it skips a field it does not know by its wire
// type and keeps an enum value it does not know as its number, and it refuses bytes cut short, a
// field of a wrong wire type or written twice, a field its kind writes that is missing, a field a
// later kind adds in an earlier one — field 20 of evidence in kind 1, field 21 in kinds 1 and 2,
// fields 9 and 10 of a module in a state of kind 1 —, a value whose flag of presence is missing, a
// string that is no UTF-8, elements of a repeated field out of their canonical order or twice,
// and NaN and the infinities. What it accepts, the writer gives back as the same bytes, but for
// the fields it does not know.

// EncodeEvidence gives the bytes of e. It refuses, with ErrInvalid: a kind outside 1 to
// EvidenceKind; a string that is no UTF-8; a duplicate place, skill or trap; items of the task and
// answers whose places do not match one for one; a value set without its flag of presence — a
// rung without HasRung, say; a leading question on request in evidence of kind 1, and a
// correction in kinds 1 and 2, which do not record them; and a correction with skills, items,
// traps or the other facts of an attempt, or whose task is not that of LevelCorrection.
func EncodeEvidence(e Evidence) ([]byte, error) {
	if err := checkEvidence(e); err != nil {
		return nil, err
	}
	var w wire.Writer
	writeEvidence(&w, e)
	return w.Bytes(), nil
}

// DecodeEvidence reads evidence from b, of any kind from 1: a later kind too.
func DecodeEvidence(b []byte) (Evidence, error) {
	e, err := readEvidence(b)
	if err == nil {
		err = canonicalEvidence(e)
	}
	if err != nil {
		return Evidence{}, wrapRead("evidence", err)
	}
	return e, nil
}

// EncodeState gives the bytes of s. It refuses, with ErrNotFinite, a NaN or an infinity, and,
// with ErrInvalid, a kind outside 1 to StateKind, a string that is no UTF-8, a duplicate module,
// skill or trap, a mastery band without Mastered, a count of the guard below 0, and, in a
// state of kind 1, which does not record them, a count of the guard or a trial series ended by an
// adult.
func EncodeState(s State) ([]byte, error) {
	if err := checkState(s); err != nil {
		return nil, err
	}
	var w wire.Writer
	writeState(&w, s)
	return w.Bytes(), nil
}

// DecodeState reads a state from b.
func DecodeState(b []byte) (State, error) {
	s, err := readState(b)
	if err == nil {
		err = canonicalState(s)
	}
	if err != nil {
		return State{}, wrapRead("state", err)
	}
	return s, nil
}

func wrapRead(what string, err error) error {
	if errors.Is(err, ErrNotFinite) {
		return err
	}
	return &detailError{ErrCorrupt, what + ": " + err.Error()}
}

// The writers, field by field, in ascending order of their numbers.

func writeBand(w *wire.Writer, field int, b BandRef) {
	w.Message(field, func(w *wire.Writer) {
		w.Int64(fieldBandPlace, int64(b.Place))
		w.String(fieldBandID, b.ID)
	})
}

func writeSlots(w *wire.Writer, field int, slots []Slot) {
	for _, s := range slots {
		w.Message(field, func(w *wire.Writer) {
			w.Int64(fieldSlotCheckKind, int64(s.Check))
			if s.HasOptions {
				w.Int64(fieldSlotOptionCount, int64(s.Options))
			}
		})
	}
}

func writeStrings(w *wire.Writer, field int, s []string) {
	for _, v := range sortedCopy(s) {
		w.String(field, v)
	}
}

func writeTrapCounts(w *wire.Writer, field int, traps []TrapCount) {
	traps = slices.Clone(traps)
	slices.SortFunc(traps, func(a, b TrapCount) int { return strings.Compare(a.ID, b.ID) })
	for _, t := range traps {
		w.Message(field, func(w *wire.Writer) {
			w.String(fieldTrapCountID, t.ID)
			w.Int64(fieldTrapCountEvidenceCount, t.Count)
			w.Int64(fieldTrapCountLastFoundTime, t.Last)
		})
	}
}

func writeEvidence(w *wire.Writer, e Evidence) {
	t := e.Task
	w.Int64(fieldEvidenceKind, int64(e.Kind))
	w.String(fieldEvidenceTaskID, t.ID)
	w.Bool(fieldEvidenceTemplate, t.Template)
	writeBand(w, fieldEvidenceTaskBand, t.Band)
	w.Int64(fieldEvidenceDifficulty, int64(t.Difficulty))
	writeBand(w, fieldEvidenceStartBand, e.StartBand)
	writeStrings(w, fieldEvidencePrimarySkills, t.Primary)
	writeStrings(w, fieldEvidenceSupportingSkills, t.Supporting)
	w.Int64(fieldEvidenceKeyProvenance, int64(t.KeyProvenance))
	w.Int64(fieldEvidenceSkillTagsProvenance, int64(t.SkillTagsProvenance))
	w.Int64(fieldEvidenceOutcome, int64(e.Outcome.Kind))
	if e.Outcome.HasRung {
		w.Int64(fieldEvidenceHintRung, int64(e.Outcome.Rung))
	}
	w.Int64(fieldEvidencePracticeTime, e.PracticeTime)
	w.Bool(fieldEvidenceTechniqueOnStep, e.TechniqueOnStep)
	w.Bool(fieldEvidenceTechniqueOnRequest, e.TechniqueOnRequest)
	w.Int64(fieldEvidenceExplainAgainCount, int64(e.ExplainAgainCount))
	if e.HasSolutionSeenFact {
		w.Bool(fieldEvidenceSolutionSeenBefore, e.SolutionSeenBefore)
	}
	items := slices.Clone(t.Items)
	slices.SortFunc(items, func(a, b TaskItem) int { return cmp.Compare(a.Place, b.Place) })
	for _, item := range items {
		i := slices.IndexFunc(e.Items, func(a ItemAnswer) bool { return a.Place == item.Place })
		writeItem(w, item, e.Items[i])
	}
	traps := slices.Clone(e.Traps)
	slices.SortFunc(traps, func(a, b FoundTrap) int { return strings.Compare(a.ID, b.ID) })
	for _, trap := range traps {
		w.Message(fieldEvidenceTraps, func(w *wire.Writer) {
			w.String(fieldFoundTrapID, trap.ID)
			writeStrings(w, fieldFoundTrapSkills, trap.Skills)
			w.PackedInt64(fieldFoundTrapItemPlaces, widen(sortedCopy(trap.ItemPlaces)))
		})
	}
	// Kind 2 adds field 20, and kind 1 does not record it: its bytes stay as they were written.
	if e.Kind >= leadingQuestionKind {
		w.Bool(fieldEvidenceLeadingQuestionOnRequest, e.LeadingQuestionOnRequest)
	}
	// Kind 3 adds field 21 the same way.
	if e.Kind >= correctionKind {
		w.Bool(fieldEvidenceLevelSetByAdult, e.LevelSetByAdult)
	}
}

func writeItem(w *wire.Writer, item TaskItem, a ItemAnswer) {
	w.Message(fieldEvidenceItems, func(w *wire.Writer) {
		w.Int64(fieldItemPlace, int64(item.Place))
		w.Bool(fieldItemOwnSkills, item.OwnSkills)
		writeStrings(w, fieldItemPrimarySkills, item.Primary)
		writeStrings(w, fieldItemSupportingSkills, item.Supporting)
		writeSlots(w, fieldItemAnswerSlots, item.Slots)
		w.Bool(fieldItemSubmitted, a.Submitted)
		w.PackedInt64(fieldItemUnverdictedReasons, widen(a.Unverdicted))
		if a.HasFirstVerdict {
			w.Int64(fieldItemFirstVerdict, int64(a.FirstVerdict))
			w.Int64(fieldItemWrongBeforeFirstVerdict, int64(a.WrongBeforeFirst))
		}
		if a.HasFinalVerdict {
			w.Int64(fieldItemFinalVerdict, int64(a.FinalVerdict))
		}
		if a.HasFirstCorrect {
			w.Int64(fieldItemWrongBeforeFirstCorrect, int64(a.WrongBeforeFirstCorrect))
		}
	})
}

func widen[T ~int32](values []T) []int64 {
	out := make([]int64, len(values))
	for i, v := range values {
		out[i] = int64(v)
	}
	return out
}

func writeState(w *wire.Writer, s State) {
	w.Int64(fieldStateKind, int64(s.Kind))
	modules := slices.Clone(s.Modules)
	slices.SortFunc(modules, func(a, b ModuleState) int { return strings.Compare(a.Module, b.Module) })
	for _, m := range modules {
		w.Message(fieldStateModules, func(w *wire.Writer) {
			w.String(fieldModuleID, m.Module)
			w.Double(fieldModuleLevel, m.Level)
			w.Double(fieldModuleAnswerWeight, m.AnswerWeight)
			writeBand(w, fieldModuleStartBand, m.StartBand)
			for _, a := range m.Trial {
				w.Message(fieldModuleTrialAnswers, func(w *wire.Writer) {
					w.Int64(fieldTrialAnswerBandPlace, int64(a.BandPlace))
					w.Int64(fieldTrialAnswerDifficulty, int64(a.Difficulty))
					writeSlots(w, fieldTrialAnswerAnswerSlots, a.Slots)
					w.Int64(fieldTrialAnswerVerdict, int64(a.Verdict))
					w.Double(fieldTrialAnswerWeight, a.Weight)
				})
			}
			writeTrapCounts(w, fieldModuleTraps, m.Traps)
			skills := slices.Clone(m.Skills)
			slices.SortFunc(skills, func(a, b SkillState) int { return strings.Compare(a.Skill, b.Skill) })
			for _, sk := range skills {
				w.Message(fieldModuleSkills, func(w *wire.Writer) { writeSkill(w, sk) })
			}
			w.Bool(fieldModuleLastEvidenceInTrial, m.LastInTrial)
			// Kind 2 adds fields 9 and 10, and kind 1 does not record them.
			if s.Kind >= guardKind {
				w.Int64(fieldModuleFailuresInARow, int64(m.FailuresInARow))
				w.Bool(fieldModuleTrialEndedByAdult, m.TrialEndedByAdult)
			}
		})
	}
}

func writeSkill(w *wire.Writer, sk SkillState) {
	w.String(fieldSkillID, sk.Skill)
	w.Double(fieldSkillCorrection, sk.Correction)
	w.Double(fieldSkillAnswerWeight, sk.AnswerWeight)
	w.Int64(fieldSkillEvidenceCount, sk.EvidenceCount)
	w.Int64(fieldSkillLastPracticeTime, sk.LastPractice)
	w.Double(fieldSkillStability, sk.Stability)
	writeTrapCounts(w, fieldSkillTraps, sk.Traps)
	if sk.Mastered {
		w.Int64(fieldSkillMasteryBandPlace, int64(sk.MasteryBand))
	}
	w.Int64(fieldSkillFailuresInARow, int64(sk.FailuresInARow))
}
