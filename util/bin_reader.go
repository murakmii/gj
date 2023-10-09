package util

import (
	"encoding/binary"
	"io"
)

type BinReader struct {
	bytes  []byte
	offset int
}

func NewBinReader(src io.Reader) (*BinReader, error) {
	var err error
	r := &BinReader{}

	r.bytes, err = io.ReadAll(src)
	if err != nil {
		return nil, err
	}

	return r, nil
}

func (r *BinReader) Skip(n int) {
	r.bytes = r.bytes[n:]
	r.offset += n
}

func (r *BinReader) SkipToAlign(align int) {
	r.Skip(r.offset % align)
}

func (r *BinReader) ReadByte() uint8 {
	b := r.bytes[0]
	r.bytes = r.bytes[1:]
	r.offset += 1
	return b
}

func (r *BinReader) ReadBytes(n int) []byte {
	bytes := r.bytes[0:n]
	r.bytes = r.bytes[n:]
	r.offset += n
	return bytes
}

func (r *BinReader) ReadUint16() uint16 {
	i := binary.BigEndian.Uint16(r.bytes)
	r.bytes = r.bytes[2:]
	r.offset += 2
	return i
}

func (r *BinReader) ReadUint32() uint32 {
	i := binary.BigEndian.Uint32(r.bytes)
	r.bytes = r.bytes[4:]
	r.offset += 4
	return i
}

func (r *BinReader) ReadUint64() uint64 {
	i := binary.BigEndian.Uint64(r.bytes)
	r.bytes = r.bytes[8:]
	r.offset += 8
	return i
}

func (r *BinReader) Remain() int {
	return len(r.bytes)
}
