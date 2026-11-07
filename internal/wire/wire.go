// Package wire writes and reads the binary form of protobuf, the part of it the schema of
// learner-model uses, on Go's standard library alone (R10 in docs/decisions.md;
// docs/architecture.md, section 3.7): varints, 64-bit fixed numbers in little-endian, and
// length-delimited strings, messages and packed lists. What the fields mean, in what order they
// go and which are required is the caller's.
package wire

import (
	"encoding/binary"
	"errors"
	"math"
	"strconv"
)

// The wire types the schema uses; 3 and 4, the groups, and 5, the 32-bit numbers, it does not.
const (
	Varint  = 0
	Fixed64 = 1
	Bytes   = 2
)

// ErrCorrupt is bytes that break the binary form: cut short, a varint past 10 bytes, a length
// past the end, a field number of 0 or a wire type the form does not have.
var ErrCorrupt = errors.New("corrupt bytes")

// Writer appends fields to a buffer. A caller writes them in ascending order of their numbers.
type Writer struct {
	buf []byte
}

// Bytes is what has been written.
func (w *Writer) Bytes() []byte { return w.buf }

func (w *Writer) tag(field, wireType int) {
	w.buf = binary.AppendUvarint(w.buf, uint64(field)<<3|uint64(wireType))
}

// Uint64 writes a varint field.
func (w *Writer) Uint64(field int, v uint64) {
	w.tag(field, Varint)
	w.buf = binary.AppendUvarint(w.buf, v)
}

// Int64 writes an int64 or int32 field: a negative number is ten bytes, as protobuf has it.
func (w *Writer) Int64(field int, v int64) { w.Uint64(field, uint64(v)) }

// Bool writes a bool field.
func (w *Writer) Bool(field int, v bool) {
	var b uint64
	if v {
		b = 1
	}
	w.Uint64(field, b)
}

// Double writes a double field: its eight IEEE bytes in little-endian, as they are, −0 too. The
// caller refuses NaN and the infinities first.
func (w *Writer) Double(field int, v float64) {
	w.tag(field, Fixed64)
	w.buf = binary.LittleEndian.AppendUint64(w.buf, math.Float64bits(v))
}

// String writes a string field.
func (w *Writer) String(field int, s string) {
	w.tag(field, Bytes)
	w.buf = binary.AppendUvarint(w.buf, uint64(len(s)))
	w.buf = append(w.buf, s...)
}

// Message writes a message field: the bytes body writes into a writer of its own.
func (w *Writer) Message(field int, body func(*Writer)) {
	var inner Writer
	body(&inner)
	w.tag(field, Bytes)
	w.buf = binary.AppendUvarint(w.buf, uint64(len(inner.buf)))
	w.buf = append(w.buf, inner.buf...)
}

// PackedInt64 writes a packed list of int32, int64 or enum values; an empty list, nothing.
func (w *Writer) PackedInt64(field int, values []int64) {
	if len(values) == 0 {
		return
	}
	var inner []byte
	for _, v := range values {
		inner = binary.AppendUvarint(inner, uint64(v))
	}
	w.tag(field, Bytes)
	w.buf = binary.AppendUvarint(w.buf, uint64(len(inner)))
	w.buf = append(w.buf, inner...)
}

// Reader reads the fields of one message in turn.
type Reader struct {
	buf []byte
}

// NewReader reads b.
func NewReader(b []byte) *Reader { return &Reader{buf: b} }

// Done tells whether every field has been read.
func (r *Reader) Done() bool { return len(r.buf) == 0 }

func corrupt(detail string) error { return &corruptError{detail} }

type corruptError struct{ detail string }

func (e *corruptError) Error() string { return ErrCorrupt.Error() + ": " + e.detail }

func (e *corruptError) Unwrap() error { return ErrCorrupt }

func (r *Reader) uvarint() (uint64, error) {
	v, n := binary.Uvarint(r.buf)
	if n <= 0 {
		return 0, corrupt("a varint cut short or past 10 bytes")
	}
	r.buf = r.buf[n:]
	return v, nil
}

// Next reads the tag of the next field: its number and wire type.
func (r *Reader) Next() (field, wireType int, err error) {
	t, err := r.uvarint()
	if err != nil {
		return 0, 0, err
	}
	if t>>3 == 0 || t>>3 > math.MaxInt32 {
		return 0, 0, corrupt("field number " + strconv.FormatUint(t>>3, 10))
	}
	wireType = int(t & 7)
	switch wireType {
	case Varint, Fixed64, Bytes, 5:
	default:
		return 0, 0, corrupt("wire type " + strconv.Itoa(wireType))
	}
	return int(t >> 3), wireType, nil
}

// Int64 reads the value of an int64 field.
func (r *Reader) Int64() (int64, error) {
	v, err := r.uvarint()
	return int64(v), err
}

// Int32 reads the value of an int32 or enum field, refusing one outside the int32 range.
func (r *Reader) Int32() (int32, error) {
	v, err := r.Int64()
	if err != nil {
		return 0, err
	}
	if v < math.MinInt32 || v > math.MaxInt32 {
		return 0, corrupt("an int32 of " + strconv.FormatInt(v, 10))
	}
	return int32(v), nil
}

// Bool reads the value of a bool field, refusing a varint other than 0 and 1.
func (r *Reader) Bool() (bool, error) {
	v, err := r.uvarint()
	if err != nil {
		return false, err
	}
	if v > 1 {
		return false, corrupt("a bool of " + strconv.FormatUint(v, 10))
	}
	return v == 1, nil
}

// Double reads the value of a double field.
func (r *Reader) Double() (float64, error) {
	if len(r.buf) < 8 {
		return 0, corrupt("a double cut short")
	}
	v := math.Float64frombits(binary.LittleEndian.Uint64(r.buf))
	r.buf = r.buf[8:]
	return v, nil
}

// Bytes reads the value of a length-delimited field: a string, a message or a packed list.
func (r *Reader) Bytes() ([]byte, error) {
	n, err := r.uvarint()
	if err != nil {
		return nil, err
	}
	if n > uint64(len(r.buf)) {
		return nil, corrupt("a length of " + strconv.FormatUint(n, 10) + " past the end")
	}
	b := r.buf[:n]
	r.buf = r.buf[n:]
	return b, nil
}

// PackedInt64 reads the value of a packed list of varints.
func (r *Reader) PackedInt64() ([]int64, error) {
	b, err := r.Bytes()
	if err != nil {
		return nil, err
	}
	inner := Reader{buf: b}
	var values []int64
	for !inner.Done() {
		v, err := inner.Int64()
		if err != nil {
			return nil, err
		}
		values = append(values, v)
	}
	return values, nil
}

// Skip passes over the value of a field this version does not know, by its wire type.
func (r *Reader) Skip(wireType int) error {
	var err error
	switch wireType {
	case Varint:
		_, err = r.uvarint()
	case Fixed64:
		_, err = r.Double()
	case Bytes:
		_, err = r.Bytes()
	case 5:
		if len(r.buf) < 4 {
			return corrupt("a fixed32 cut short")
		}
		r.buf = r.buf[4:]
	default:
		err = corrupt("wire type " + strconv.Itoa(wireType))
	}
	return err
}
