package class_file

import (
	"encoding/binary"
	"io"
)

type reader struct {
	bytes []byte
}

func open(cfReader io.Reader) (*reader, error) {
	var err error
	r := &reader{}

	r.bytes, err = io.ReadAll(cfReader)
	if err != nil {
		return nil, err
	}

	return r, nil
}

func (r *reader) skip(n int) {
	r.bytes = r.bytes[n:]
}

func (r *reader) readByte() uint8 {
	b := r.bytes[0]
	r.bytes = r.bytes[1:]
	return b
}

func (r *reader) readBytes(n int) []byte {
	bytes := r.bytes[0:n]
	r.bytes = r.bytes[n:]
	return bytes
}

func (r *reader) readUint16() uint16 {
	i := binary.BigEndian.Uint16(r.bytes)
	r.bytes = r.bytes[2:]
	return i
}

func (r *reader) readUint32() uint32 {
	i := binary.BigEndian.Uint32(r.bytes)
	r.bytes = r.bytes[4:]
	return i
}

func (r *reader) readUint64() uint64 {
	i := binary.BigEndian.Uint64(r.bytes)
	r.bytes = r.bytes[8:]
	return i
}
