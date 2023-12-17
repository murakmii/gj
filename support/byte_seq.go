package support

import (
	"encoding/binary"
	"io"
)

type ByteSeq struct {
	src    []byte
	bytes  []byte
	offset int
}

func NewByteSeq(src io.Reader) (*ByteSeq, error) {
	var err error
	r := &ByteSeq{}

	r.bytes, err = io.ReadAll(src)
	if err != nil {
		return nil, err
	}

	r.src = r.bytes
	return r, nil
}

func (r *ByteSeq) Skip(n int) {
	r.bytes = r.bytes[n:]
	r.offset += n
}

func (r *ByteSeq) SkipToAlign(align int) {
	mod := r.offset % align
	skip := 0
	if mod > 0 {
		skip = align - mod
	}
	r.Skip(skip)
}

func (r *ByteSeq) Seek(pos int) {
	r.bytes = r.src[pos:]
	r.offset = pos
}

func (r *ByteSeq) Pos() int {
	return r.offset
}

func (r *ByteSeq) ReadByte() uint8 {
	b := r.bytes[0]
	r.bytes = r.bytes[1:]
	r.offset += 1
	return b
}

func (r *ByteSeq) ReadBytes(n int) []byte {
	bytes := r.bytes[0:n]
	r.bytes = r.bytes[n:]
	r.offset += n
	return bytes
}

func (r *ByteSeq) ReadUint16() uint16 {
	i := binary.BigEndian.Uint16(r.bytes)
	r.bytes = r.bytes[2:]
	r.offset += 2
	return i
}

func (r *ByteSeq) ReadUint32() uint32 {
	i := binary.BigEndian.Uint32(r.bytes)
	r.bytes = r.bytes[4:]
	r.offset += 4
	return i
}

func (r *ByteSeq) ReadUint64() uint64 {
	i := binary.BigEndian.Uint64(r.bytes)
	r.bytes = r.bytes[8:]
	r.offset += 8
	return i
}

func (r *ByteSeq) Remain() int {
	return len(r.bytes)
}
