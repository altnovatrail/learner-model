package learnermodel

import (
	"cmp"
	"errors"
	"math"
	"slices"
	"strings"
	"unicode/utf8"
)

// The checks of the codec: before writing, a value the writer can give canonically; after
// reading, a value the writer would have given (codec.go).

func checkEvidence(e Evidence) error {
	if e.Kind < 1 || e.Kind > EvidenceKind {
		return invalid("evidence of kind " + itoa(e.Kind))
	}
	return checkEvidenceContent(e)
}

// checkEvidenceContent is checkEvidence but for the kind: Append takes evidence of a later kind
// too.
func checkEvidenceContent(e Evidence) error {
	t := e.Task
	if err := utf8Strings(t.ID, t.Band.ID, e.StartBand.ID); err != nil {
		return err
	}
	if err := idList("primary skill", t.Primary); err != nil {
		return err
	}
	if err := idList("supporting skill", t.Supporting); err != nil {
		return err
	}
	if !e.Outcome.HasRung && e.Outcome.Rung != 0 {
		return invalid("a rung without HasRung")
	}
	if !e.HasSolutionSeenFact && e.SolutionSeenBefore {
		return invalid("SolutionSeenBefore without HasSolutionSeenFact")
	}
	taskPlaces := make([]int32, 0, len(t.Items))
	for _, item := range t.Items {
		if err := idList("primary skill of an item", item.Primary); err != nil {
			return err
		}
		if err := idList("supporting skill of an item", item.Supporting); err != nil {
			return err
		}
		if err := checkSlots(item.Slots); err != nil {
			return err
		}
		taskPlaces = append(taskPlaces, item.Place)
	}
	answerPlaces := make([]int32, 0, len(e.Items))
	for _, a := range e.Items {
		if err := checkAnswer(a); err != nil {
			return err
		}
		answerPlaces = append(answerPlaces, a.Place)
	}
	slices.Sort(taskPlaces)
	slices.Sort(answerPlaces)
	if err := uniqueSorted("item place", taskPlaces); err != nil {
		return err
	}
	if !slices.Equal(taskPlaces, answerPlaces) {
		return invalid("the items of the task and the answers name different places")
	}
	traps := make([]string, 0, len(e.Traps))
	for _, trap := range e.Traps {
		traps = append(traps, trap.ID)
		if err := idList("skill of trap "+trap.ID, trap.Skills); err != nil {
			return err
		}
		if err := uniqueSorted("item place of trap "+trap.ID, sortedCopy(trap.ItemPlaces)); err != nil {
			return err
		}
	}
	return idList("trap", traps)
}

func checkAnswer(a ItemAnswer) error {
	switch {
	case !a.HasFirstVerdict && (a.FirstVerdict != 0 || a.WrongBeforeFirst != 0):
		return invalid("a first verdict without HasFirstVerdict at item " + itoa(a.Place))
	case !a.HasFinalVerdict && a.FinalVerdict != 0:
		return invalid("a final verdict without HasFinalVerdict at item " + itoa(a.Place))
	case !a.HasFirstCorrect && a.WrongBeforeFirstCorrect != 0:
		return invalid("WrongBeforeFirstCorrect without HasFirstCorrect at item " + itoa(a.Place))
	}
	return nil
}

func checkSlots(slots []Slot) error {
	for _, s := range slots {
		if !s.HasOptions && s.Options != 0 {
			return invalid("options without HasOptions")
		}
	}
	return nil
}

func checkState(s State) error {
	if s.Kind < 1 || s.Kind > StateKind {
		return invalid("a state of kind " + itoa(s.Kind))
	}
	modules := make([]string, 0, len(s.Modules))
	for _, m := range s.Modules {
		modules = append(modules, m.Module)
		if err := utf8Strings(m.StartBand.ID); err != nil {
			return err
		}
		if err := finite("the level of "+m.Module, m.Level, m.AnswerWeight); err != nil {
			return err
		}
		for _, a := range m.Trial {
			if err := finite("a trial answer of "+m.Module, a.Weight); err != nil {
				return err
			}
			if err := checkSlots(a.Slots); err != nil {
				return err
			}
		}
		if err := trapCounts(m.Module, m.Traps); err != nil {
			return err
		}
		skills := make([]string, 0, len(m.Skills))
		for _, sk := range m.Skills {
			skills = append(skills, sk.Skill)
			if err := finite("skill "+sk.Skill, sk.Correction, sk.AnswerWeight, sk.Stability); err != nil {
				return err
			}
			if !sk.Mastered && sk.MasteryBand != 0 {
				return invalid("a mastery band without Mastered at skill " + sk.Skill)
			}
			if err := trapCounts(sk.Skill, sk.Traps); err != nil {
				return err
			}
		}
		if err := idList("skill of "+m.Module, skills); err != nil {
			return err
		}
	}
	return idList("module", modules)
}

func trapCounts(owner string, traps []TrapCount) error {
	ids := make([]string, 0, len(traps))
	for _, t := range traps {
		ids = append(ids, t.ID)
	}
	return idList("trap of "+owner, ids)
}

func finite(what string, values ...float64) error {
	for _, v := range values {
		if math.IsNaN(v) || math.IsInf(v, 0) {
			return &detailError{ErrNotFinite, what}
		}
	}
	return nil
}

// idList refuses a list of IDs with one twice or one that is no UTF-8.
func idList(what string, ids []string) error {
	if err := utf8Strings(ids...); err != nil {
		return err
	}
	return uniqueSorted(what, sortedCopy(ids))
}

func utf8Strings(ss ...string) error {
	for _, s := range ss {
		if !utf8.ValidString(s) {
			return invalid("a string that is no UTF-8: " + shownBytes(s))
		}
	}
	return nil
}

// shownBytes shows a string byte by byte, as \xNN past ASCII, without strconv.Quote, which
// reads Unicode tables (rule 6).
func shownBytes(s string) string {
	const hex = "0123456789abcdef"
	var b strings.Builder
	for i := range len(s) {
		c := s[i]
		if c >= 0x20 && c < 0x7f {
			b.WriteByte(c)
		} else {
			b.WriteString(`\x`)
			b.WriteByte(hex[c>>4])
			b.WriteByte(hex[c&15])
		}
	}
	return b.String()
}

func uniqueSorted[T int32 | string](what string, sorted []T) error {
	for i := 1; i < len(sorted); i++ {
		if sorted[i] == sorted[i-1] {
			return invalid(what + " twice")
		}
	}
	return nil
}

// After reading: the elements of each repeated field the writer orders are in that order, each
// once. The writer gives them so, and anything else is no bytes of it.

func canonicalEvidence(e Evidence) error {
	t := e.Task
	if e.Kind < 1 {
		return errors.New("evidence of kind " + itoa(e.Kind))
	}
	if err := ascending("primary_skills", t.Primary, strings.Compare); err != nil {
		return err
	}
	if err := ascending("supporting_skills", t.Supporting, strings.Compare); err != nil {
		return err
	}
	if err := ascending("items", t.Items, func(a, b TaskItem) int { return cmp.Compare(a.Place, b.Place) }); err != nil {
		return err
	}
	for _, item := range t.Items {
		if err := ascending("primary_skills of an item", item.Primary, strings.Compare); err != nil {
			return err
		}
		if err := ascending("supporting_skills of an item", item.Supporting, strings.Compare); err != nil {
			return err
		}
	}
	if err := ascending("traps", e.Traps, func(a, b FoundTrap) int { return strings.Compare(a.ID, b.ID) }); err != nil {
		return err
	}
	for _, trap := range e.Traps {
		if err := ascending("skills of a trap", trap.Skills, strings.Compare); err != nil {
			return err
		}
		if err := ascending("item_places of a trap", trap.ItemPlaces, cmp.Compare[int32]); err != nil {
			return err
		}
	}
	return nil
}

func canonicalState(s State) error {
	if s.Kind < 1 {
		return errors.New("a state of kind " + itoa(s.Kind))
	}
	byID := func(a, b TrapCount) int { return strings.Compare(a.ID, b.ID) }
	if err := ascending("modules", s.Modules, func(a, b ModuleState) int { return strings.Compare(a.Module, b.Module) }); err != nil {
		return err
	}
	for _, m := range s.Modules {
		if err := ascending("traps of a module", m.Traps, byID); err != nil {
			return err
		}
		if err := ascending("skills", m.Skills, func(a, b SkillState) int { return strings.Compare(a.Skill, b.Skill) }); err != nil {
			return err
		}
		for _, sk := range m.Skills {
			if err := ascending("traps of a skill", sk.Traps, byID); err != nil {
				return err
			}
		}
	}
	return nil
}

// ascending refuses a list whose elements are not in strictly ascending order.
func ascending[T any](what string, s []T, compare func(a, b T) int) error {
	for i := 1; i < len(s); i++ {
		if compare(s[i-1], s[i]) >= 0 {
			return errors.New(what + " out of order or twice")
		}
	}
	return nil
}
