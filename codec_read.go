package learnermodel

import (
	"errors"
	"math"
	"strconv"
	"unicode/utf8"

	"github.com/altnovatrail/learner-model/internal/wire"
)

// The readers of the codec. A message reads its fields into a value, each by one helper below,
// and marks each non-repeated field it has read, so that a field written twice and a required
// field missing are both refused. The wire type is checked in one place, read, before any value is
// read or any flag set.

type fields uint64

func (f *fields) mark(field int) error {
	bit := fields(1) << field
	if *f&bit != 0 {
		return errors.New("field " + strconv.Itoa(field) + " twice")
	}
	*f |= bit
	return nil
}

func (f fields) has(field int) bool { return f&(fields(1)<<field) != 0 }

// require refuses a message that lacks one of the fields, which every kind from 1 writes.
func (f fields) require(what string, numbers ...int) error {
	for _, n := range numbers {
		if !f.has(n) {
			return errors.New(what + " without field " + strconv.Itoa(n))
		}
	}
	return nil
}

// field is a field of a message as it is read: its number and wire type, and the reader at its
// value.
type field struct {
	number, wireType int
	r                *wire.Reader
}

// reading is what a message's visit gives for a known field: the wire type the field must have,
// and how to read its value. A field visit does not know it gives as nil.
type reading struct {
	wireType int
	read     func(*wire.Reader) error
}

// readMessage reads the fields of b in turn. visit names, for each field number it knows, the
// wire type and the reading; a field of another wire type is refused, and a field visit does not
// know is skipped by its wire type.
func readMessage(b []byte, visit func(number int, seen *fields) *reading) (fields, error) {
	var seen fields
	r := wire.NewReader(b)
	for !r.Done() {
		number, wireType, err := r.Next()
		if err != nil {
			return seen, err
		}
		rd := visit(number, &seen)
		switch {
		case rd == nil:
			err = r.Skip(wireType)
		case rd.wireType != wireType:
			err = errors.New("field " + strconv.Itoa(number) + " of wire type " + strconv.Itoa(wireType))
		default:
			err = rd.read(r)
		}
		if err != nil {
			return seen, err
		}
	}
	return seen, nil
}

// one reads a non-repeated scalar field into dst, marking it in seen, and sets the flags of
// presence given, after it has read.
func one[T any](number int, seen *fields, wireType int, read func(*wire.Reader) (T, error), dst *T, flags ...*bool) *reading {
	return &reading{wireType, func(r *wire.Reader) error {
		if err := seen.mark(number); err != nil {
			return err
		}
		v, err := read(r)
		if err != nil {
			return err
		}
		*dst = v
		for _, flag := range flags {
			*flag = true
		}
		return nil
	}}
}

func int32Field(number int, seen *fields, dst *int32, flags ...*bool) *reading {
	return one(number, seen, wire.Varint, (*wire.Reader).Int32, dst, flags...)
}

func int64Field(number int, seen *fields, dst *int64) *reading {
	return one(number, seen, wire.Varint, (*wire.Reader).Int64, dst)
}

func boolField(number int, seen *fields, dst *bool, flags ...*bool) *reading {
	return one(number, seen, wire.Varint, (*wire.Reader).Bool, dst, flags...)
}

func readDouble(r *wire.Reader) (float64, error) {
	v, err := r.Double()
	if err == nil && (math.IsNaN(v) || math.IsInf(v, 0)) {
		err = &detailError{ErrNotFinite, "a double that is not finite"}
	}
	return v, err
}

func doubleField(number int, seen *fields, dst *float64) *reading {
	return one(number, seen, wire.Fixed64, readDouble, dst)
}

func readString(r *wire.Reader) (string, error) {
	b, err := r.Bytes()
	if err == nil && !utf8.Valid(b) {
		err = errors.New("a string that is no UTF-8")
	}
	return string(b), err
}

func stringField(number int, seen *fields, dst *string) *reading {
	return one(number, seen, wire.Bytes, readString, dst)
}

// enumField reads an enum field into a Go enum type, keeping a value this version does not know.
func enumField[T ~int32](number int, seen *fields, dst *T, flags ...*bool) *reading {
	read := func(r *wire.Reader) (T, error) {
		v, err := r.Int32()
		return T(v), err
	}
	return one(number, seen, wire.Varint, read, dst, flags...)
}

// appendString reads one element of a repeated string field.
func appendString(dst *[]string) *reading {
	return &reading{wire.Bytes, func(r *wire.Reader) error {
		s, err := readString(r)
		*dst = append(*dst, s)
		return err
	}}
}

// message reads a message field into read; with seen, a non-repeated one, which it marks.
func message(number int, seen *fields, read func([]byte) error) *reading {
	return &reading{wire.Bytes, func(r *wire.Reader) error {
		if seen != nil {
			if err := seen.mark(number); err != nil {
				return err
			}
		}
		b, err := r.Bytes()
		if err != nil {
			return err
		}
		return read(b)
	}}
}

// element reads one element of a repeated message field and appends it to dst.
func element[T any](dst *[]T, read func([]byte) (T, error)) *reading {
	return message(0, nil, func(b []byte) error {
		v, err := read(b)
		*dst = append(*dst, v)
		return err
	})
}

// packed reads a packed list of int32 values, enums or places; an unpacked element is refused, as
// the writer never writes one.
func packed[T ~int32](number int, dst *[]T) *reading {
	return &reading{wire.Bytes, func(r *wire.Reader) error {
		values, err := r.PackedInt64()
		if err != nil {
			return err
		}
		for _, v := range values {
			if v < math.MinInt32 || v > math.MaxInt32 {
				return errors.New("field " + strconv.Itoa(number) + " holds " + strconv.FormatInt(v, 10))
			}
			*dst = append(*dst, T(v))
		}
		return nil
	}}
}

func readBand(b []byte) (BandRef, error) {
	var band BandRef
	seen, err := readMessage(b, func(n int, seen *fields) *reading {
		switch n {
		case 1:
			return int32Field(n, seen, &band.Place)
		case 2:
			return stringField(n, seen, &band.ID)
		}
		return nil
	})
	if err != nil {
		return BandRef{}, err
	}
	return band, seen.require("a band", 1, 2)
}

func bandField(number int, seen *fields, dst *BandRef) *reading {
	return message(number, seen, func(b []byte) (err error) {
		*dst, err = readBand(b)
		return err
	})
}

func readSlot(b []byte) (Slot, error) {
	var s Slot
	seen, err := readMessage(b, func(n int, seen *fields) *reading {
		switch n {
		case 1:
			return enumField(n, seen, &s.Check)
		case 2:
			return int32Field(n, seen, &s.Options, &s.HasOptions)
		}
		return nil
	})
	if err != nil {
		return Slot{}, err
	}
	return s, seen.require("an answer slot", 1)
}

func readTrapCount(b []byte) (TrapCount, error) {
	var t TrapCount
	seen, err := readMessage(b, func(n int, seen *fields) *reading {
		switch n {
		case 1:
			return stringField(n, seen, &t.ID)
		case 2:
			return int64Field(n, seen, &t.Count)
		case 3:
			return int64Field(n, seen, &t.Last)
		}
		return nil
	})
	if err != nil {
		return TrapCount{}, err
	}
	return t, seen.require("a trap count", 1, 2, 3)
}

func readEvidence(b []byte) (Evidence, error) {
	var e Evidence
	t := &e.Task
	seen, err := readMessage(b, func(n int, seen *fields) *reading {
		switch n {
		case 1:
			return int32Field(n, seen, &e.Kind)
		case 2:
			return stringField(n, seen, &t.ID)
		case 3:
			return boolField(n, seen, &t.Template)
		case 4:
			return bandField(n, seen, &t.Band)
		case 5:
			return int32Field(n, seen, &t.Difficulty)
		case 6:
			return bandField(n, seen, &e.StartBand)
		case 7:
			return appendString(&t.Primary)
		case 8:
			return appendString(&t.Supporting)
		case 9:
			return enumField(n, seen, &t.KeyProvenance)
		case 10:
			return enumField(n, seen, &t.SkillTagsProvenance)
		case 11:
			return enumField(n, seen, &e.Outcome.Kind)
		case 12:
			return int32Field(n, seen, &e.Outcome.Rung, &e.Outcome.HasRung)
		case 13:
			return int64Field(n, seen, &e.PracticeTime)
		case 14:
			return boolField(n, seen, &e.TechniqueOnStep)
		case 15:
			return boolField(n, seen, &e.TechniqueOnRequest)
		case 16:
			return int32Field(n, seen, &e.ExplainAgainCount)
		case 17:
			return boolField(n, seen, &e.SolutionSeenBefore, &e.HasSolutionSeenFact)
		case 18:
			return message(n, nil, func(b []byte) error {
				item, answer, err := readItem(b)
				t.Items = append(t.Items, item)
				e.Items = append(e.Items, answer)
				return err
			})
		case 19:
			return element(&e.Traps, readFoundTrap)
		case 20:
			return boolField(n, seen, &e.LeadingQuestionOnRequest)
		}
		return nil
	})
	if err != nil {
		return Evidence{}, err
	}
	if err := seen.require("evidence", 1, 2, 3, 4, 5, 6, 9, 10, 11, 13, 14, 15, 16); err != nil {
		return e, err
	}
	// Kind 2 adds field 20, which every kind from it writes and kind 1 does not. A kind below 1 is
	// refused for its kind, by canonicalEvidence.
	if e.Kind >= leadingQuestionKind {
		return e, seen.require("evidence of kind "+itoa(e.Kind), 20)
	}
	if e.Kind >= 1 && seen.has(20) {
		return e, errors.New("evidence of kind " + itoa(e.Kind) + " with field 20")
	}
	return e, nil
}

func readItem(b []byte) (TaskItem, ItemAnswer, error) {
	var item TaskItem
	var a ItemAnswer
	seen, err := readMessage(b, func(n int, seen *fields) *reading {
		switch n {
		case 1:
			return int32Field(n, seen, &item.Place)
		case 2:
			return boolField(n, seen, &item.OwnSkills)
		case 3:
			return appendString(&item.Primary)
		case 4:
			return appendString(&item.Supporting)
		case 5:
			return element(&item.Slots, readSlot)
		case 6:
			return boolField(n, seen, &a.Submitted)
		case 7:
			return packed(n, &a.Unverdicted)
		case 8:
			return enumField(n, seen, &a.FirstVerdict, &a.HasFirstVerdict)
		case 9:
			return int32Field(n, seen, &a.WrongBeforeFirst)
		case 10:
			return enumField(n, seen, &a.FinalVerdict, &a.HasFinalVerdict)
		case 11:
			return int32Field(n, seen, &a.WrongBeforeFirstCorrect, &a.HasFirstCorrect)
		}
		return nil
	})
	a.Place = item.Place
	if err != nil {
		return item, a, err
	}
	// The writer writes the wrong answers before the first verdict with it, and only with it.
	if seen.has(8) != seen.has(9) {
		return item, a, errors.New("an item with only one of fields 8 and 9")
	}
	return item, a, seen.require("an item", 1, 2, 6)
}

func readFoundTrap(b []byte) (FoundTrap, error) {
	var t FoundTrap
	seen, err := readMessage(b, func(n int, seen *fields) *reading {
		switch n {
		case 1:
			return stringField(n, seen, &t.ID)
		case 2:
			return appendString(&t.Skills)
		case 3:
			return packed(n, &t.ItemPlaces)
		}
		return nil
	})
	if err != nil {
		return FoundTrap{}, err
	}
	return t, seen.require("a found trap", 1)
}

func readState(b []byte) (State, error) {
	var s State
	seen, err := readMessage(b, func(n int, seen *fields) *reading {
		switch n {
		case 1:
			return int32Field(n, seen, &s.Kind)
		case 2:
			return element(&s.Modules, readModule)
		}
		return nil
	})
	if err != nil {
		return State{}, err
	}
	return s, seen.require("a state", 1)
}

func readModule(b []byte) (ModuleState, error) {
	var m ModuleState
	seen, err := readMessage(b, func(n int, seen *fields) *reading {
		switch n {
		case 1:
			return stringField(n, seen, &m.Module)
		case 2:
			return doubleField(n, seen, &m.Level)
		case 3:
			return doubleField(n, seen, &m.AnswerWeight)
		case 4:
			return bandField(n, seen, &m.StartBand)
		case 5:
			return element(&m.Trial, readTrialAnswer)
		case 6:
			return element(&m.Traps, readTrapCount)
		case 7:
			return element(&m.Skills, readSkill)
		case 8:
			return boolField(n, seen, &m.LastInTrial)
		}
		return nil
	})
	if err != nil {
		return ModuleState{}, err
	}
	return m, seen.require("a module", 1, 2, 3, 4, 8)
}

func readTrialAnswer(b []byte) (TrialAnswer, error) {
	var a TrialAnswer
	seen, err := readMessage(b, func(n int, seen *fields) *reading {
		switch n {
		case 1:
			return int32Field(n, seen, &a.BandPlace)
		case 2:
			return int32Field(n, seen, &a.Difficulty)
		case 3:
			return element(&a.Slots, readSlot)
		case 4:
			return enumField(n, seen, &a.Verdict)
		case 5:
			return doubleField(n, seen, &a.Weight)
		}
		return nil
	})
	if err != nil {
		return TrialAnswer{}, err
	}
	return a, seen.require("a trial answer", 1, 2, 4, 5)
}

func readSkill(b []byte) (SkillState, error) {
	var sk SkillState
	seen, err := readMessage(b, func(n int, seen *fields) *reading {
		switch n {
		case 1:
			return stringField(n, seen, &sk.Skill)
		case 2:
			return doubleField(n, seen, &sk.Correction)
		case 3:
			return doubleField(n, seen, &sk.AnswerWeight)
		case 4:
			return int64Field(n, seen, &sk.EvidenceCount)
		case 5:
			return int64Field(n, seen, &sk.LastPractice)
		case 6:
			return doubleField(n, seen, &sk.Stability)
		case 7:
			return element(&sk.Traps, readTrapCount)
		case 8:
			return int32Field(n, seen, &sk.MasteryBand, &sk.Mastered)
		case 9:
			return int32Field(n, seen, &sk.FailuresInARow)
		}
		return nil
	})
	if err != nil {
		return SkillState{}, err
	}
	return sk, seen.require("a skill", 1, 2, 3, 4, 5, 6, 9)
}
