package learnermodel

import (
	"cmp"
	"errors"
	"slices"
	"strings"

	"github.com/altnovatrail/learner-model/internal/wire"
)

// The codec of the evidence and the state: protobuf's binary form by the schema in
// proto/novatrail/learnermodel/v1/, written by this code alone (R10; docs/architecture.md, section
// 3.7). One value gives one run of bytes: the fields go in ascending order of their numbers, a
// scalar with a value is written even when it is 0, lists of scalars are packed, and the elements
// of a repeated field go in their canonical order, which the writer sets itself.
//
// The reader is strict, as only this writer writes: it skips a field it does not know by its wire
// type and keeps an enum value it does not know as its number, and it refuses bytes cut short, a
// field of a wrong wire type or written twice, a field every kind writes that is missing, a value
// whose flag of presence is missing, a string that is no UTF-8, elements of a repeated field out
// of their canonical order or twice, and NaN and the infinities. What it accepts, the writer gives
// back as the same bytes, but for the fields it does not know.

// EncodeEvidence gives the bytes of e. It refuses, with ErrInvalid: a kind outside 1 to
// EvidenceKind; a string that is no UTF-8; a duplicate place, skill or trap; items of the task and
// answers whose places do not match one for one; and a value set without its flag of presence —
// a rung without HasRung, say.
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
// skill or trap, and a mastery band without Mastered.
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
		w.Int64(1, int64(b.Place))
		w.String(2, b.ID)
	})
}

func writeSlots(w *wire.Writer, field int, slots []Slot) {
	for _, s := range slots {
		w.Message(field, func(w *wire.Writer) {
			w.Int64(1, int64(s.Check))
			if s.HasOptions {
				w.Int64(2, int64(s.Options))
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
			w.String(1, t.ID)
			w.Int64(2, t.Count)
			w.Int64(3, t.Last)
		})
	}
}

func writeEvidence(w *wire.Writer, e Evidence) {
	t := e.Task
	w.Int64(1, int64(e.Kind))
	w.String(2, t.ID)
	w.Bool(3, t.Template)
	writeBand(w, 4, t.Band)
	w.Int64(5, int64(t.Difficulty))
	writeBand(w, 6, e.StartBand)
	writeStrings(w, 7, t.Primary)
	writeStrings(w, 8, t.Supporting)
	w.Int64(9, int64(t.KeyProvenance))
	w.Int64(10, int64(t.SkillTagsProvenance))
	w.Int64(11, int64(e.Outcome.Kind))
	if e.Outcome.HasRung {
		w.Int64(12, int64(e.Outcome.Rung))
	}
	w.Int64(13, e.PracticeTime)
	w.Bool(14, e.TechniqueOnStep)
	w.Bool(15, e.TechniqueOnRequest)
	w.Int64(16, int64(e.ExplainAgainCount))
	if e.HasSolutionSeenFact {
		w.Bool(17, e.SolutionSeenBefore)
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
		w.Message(19, func(w *wire.Writer) {
			w.String(1, trap.ID)
			writeStrings(w, 2, trap.Skills)
			w.PackedInt64(3, widen(sortedCopy(trap.ItemPlaces)))
		})
	}
}

func writeItem(w *wire.Writer, item TaskItem, a ItemAnswer) {
	w.Message(18, func(w *wire.Writer) {
		w.Int64(1, int64(item.Place))
		w.Bool(2, item.OwnSkills)
		writeStrings(w, 3, item.Primary)
		writeStrings(w, 4, item.Supporting)
		writeSlots(w, 5, item.Slots)
		w.Bool(6, a.Submitted)
		w.PackedInt64(7, widen(a.Unverdicted))
		if a.HasFirstVerdict {
			w.Int64(8, int64(a.FirstVerdict))
			w.Int64(9, int64(a.WrongBeforeFirst))
		}
		if a.HasFinalVerdict {
			w.Int64(10, int64(a.FinalVerdict))
		}
		if a.HasFirstCorrect {
			w.Int64(11, int64(a.WrongBeforeFirstCorrect))
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
	w.Int64(1, int64(s.Kind))
	modules := slices.Clone(s.Modules)
	slices.SortFunc(modules, func(a, b ModuleState) int { return strings.Compare(a.Module, b.Module) })
	for _, m := range modules {
		w.Message(2, func(w *wire.Writer) {
			w.String(1, m.Module)
			w.Double(2, m.Level)
			w.Double(3, m.AnswerWeight)
			writeBand(w, 4, m.StartBand)
			for _, a := range m.Trial {
				w.Message(5, func(w *wire.Writer) {
					w.Int64(1, int64(a.BandPlace))
					w.Int64(2, int64(a.Difficulty))
					writeSlots(w, 3, a.Slots)
					w.Int64(4, int64(a.Verdict))
					w.Double(5, a.Weight)
				})
			}
			writeTrapCounts(w, 6, m.Traps)
			skills := slices.Clone(m.Skills)
			slices.SortFunc(skills, func(a, b SkillState) int { return strings.Compare(a.Skill, b.Skill) })
			for _, sk := range skills {
				w.Message(7, func(w *wire.Writer) { writeSkill(w, sk) })
			}
		})
	}
}

func writeSkill(w *wire.Writer, sk SkillState) {
	w.String(1, sk.Skill)
	w.Double(2, sk.Correction)
	w.Double(3, sk.AnswerWeight)
	w.Int64(4, sk.EvidenceCount)
	w.Int64(5, sk.LastPractice)
	w.Double(6, sk.Stability)
	writeTrapCounts(w, 7, sk.Traps)
	if sk.Mastered {
		w.Int64(8, int64(sk.MasteryBand))
	}
	w.Int64(9, int64(sk.FailuresInARow))
}
