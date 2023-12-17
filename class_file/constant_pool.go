package class_file

import (
	"github.com/murakmii/gojiai/support"
)

type (
	ConstantPool struct {
		cpInfo []interface{}
	}

	// CONSTANT_Methodref
	MethodRefCpInfo struct {
		class       uint16
		nameAndType uint16
	}

	// CONSTANT_Class
	ClassCpInfo uint16

	// CONSTANT_NameAndType
	NameAndTypeCpInfo struct {
		name uint16
		desc uint16
	}
)

const (
	utf8Tag        uint8 = 1
	classTag       uint8 = 7
	methodRefTag   uint8 = 10
	nameAndTypeTag uint8 = 12
)

func readCP(r *support.ByteSeq) *ConstantPool {
	cpCount := r.ReadUint16()

	// 仕様上、constant_poolフィールドの長さはconstant_pool_count-1となっているため1減じる必要があるが、
	// インデックス番号が1から始まることもあり、スライスの長さはconstant_pool_countとし、
	// スライスの2番目の要素から定数プールエントリーを格納していく。
	cp := &ConstantPool{cpInfo: make([]interface{}, cpCount)}
	for i := uint16(1); i < cpCount; i++ {
		switch r.ReadByte() {
		case utf8Tag:
			// CONSTANT_Utf8用の構造体は定義せず、単にstringとする
			cp.cpInfo[i] = string(r.ReadBytes(int(r.ReadUint16())))

		case classTag:
			cp.cpInfo[i] = ClassCpInfo(r.ReadUint16())

		case methodRefTag:
			cp.cpInfo[i] = &MethodRefCpInfo{class: r.ReadUint16(), nameAndType: r.ReadUint16()}

		case nameAndTypeTag:
			cp.cpInfo[i] = &NameAndTypeCpInfo{name: r.ReadUint16(), desc: r.ReadUint16()}
		}
	}

	return cp
}

func (cp *ConstantPool) Utf8(index uint16) string {
	s, ok := cp.cpInfo[index].(string)
	if !ok {
		return ""
	}
	return s
}
