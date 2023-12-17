package support

import (
	"encoding/binary"
	"io"
)

// バイト列に対する頻出操作をメソッドとしてまとめた構造体
type ByteSeq struct {
	seq    []byte
	offset int
}

func NewByteSeq(src io.Reader) (*ByteSeq, error) {
	r := &ByteSeq{}

	var err error
	r.seq, err = io.ReadAll(src)
	if err != nil {
		return nil, err
	}

	return r, nil
}

// n バイトスキップ
func (r *ByteSeq) Skip(n int) {
	r.offset += n
}

// pos で指定される絶対位置へ現在位置を更新
func (r *ByteSeq) Seek(pos int) {
	r.offset = pos
}

// 現在位置を取得
func (r *ByteSeq) Pos() int {
	return r.offset
}

// 1バイト読み取って現在位置を更新
func (r *ByteSeq) ReadByte() uint8 {
	b := r.seq[r.offset]
	r.offset += 1
	return b
}

// n バイト読み取って現在位置を更新
func (r *ByteSeq) ReadBytes(n int) []byte {
	bytes := r.seq[r.offset : r.offset+n]
	r.offset += n
	return bytes
}

// 非負の2バイト整数を読み取って現在位置を更新
func (r *ByteSeq) ReadUint16() uint16 {
	i := binary.BigEndian.Uint16(r.seq[r.offset:])
	r.offset += 2
	return i
}

// 非負の4バイト整数を読み取って現在位置を更新
func (r *ByteSeq) ReadUint32() uint32 {
	i := binary.BigEndian.Uint32(r.seq[r.offset:])
	r.offset += 4
	return i
}
