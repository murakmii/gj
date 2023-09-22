package class_file

import (
	"math"
)

type (
	ConstantPool struct {
		cpInfo []interface{}
	}

	ReferenceCpInfo struct {
		class       uint16
		nameAndType uint16
	}

	NameAndTypeCpInfo struct {
		name uint16
		desc uint16
	}

	MethodHandleCpInfo struct {
		kind  uint8
		index uint16
	}

	InvokeDynamicCpInfo struct {
		bootstrapMethodAttr uint16
		nameAndType         uint16
	}
)

const (
	utf8Tag         uint8 = 1
	intTag          uint8 = 3
	floatTag        uint8 = 4
	longTag         uint8 = 5
	doubleTag       uint8 = 6
	classTag        uint8 = 7
	strTag          uint8 = 8
	fieldRefTag     uint8 = 9
	methodRefTag    uint8 = 10
	ifMethodRefTag  uint8 = 11
	nameAndTypeTag  uint8 = 12
	methodHandleTag uint8 = 15
	methodTypeTag   uint8 = 16
	invokeDynTag    uint8 = 18
)

func readCP(r *reader) *ConstantPool {
	cpCount := r.readUint16()
	cp := &ConstantPool{cpInfo: make([]interface{}, cpCount)}

	// cp.cpInfo[0] won't be used(cp_info entries indexed from 1)
	for i := uint16(1); i < cpCount; i++ {
		switch r.readByte() {
		case utf8Tag:
			s := string(r.readBytes(int(r.readUint16())))
			cp.cpInfo[i] = &s

		case intTag:
			cp.cpInfo[i] = int(r.readUint32())

		case floatTag:
			cp.cpInfo[i] = math.Float32frombits(r.readUint32())

		case longTag:
			cp.cpInfo[i] = int64(r.readUint64())

		case doubleTag:
			cp.cpInfo[i] = math.Float64frombits(r.readUint64())

		case classTag, strTag:
			cp.cpInfo[i] = r.readUint16()

		case fieldRefTag, methodRefTag, ifMethodRefTag:
			cp.cpInfo[i] = &ReferenceCpInfo{class: r.readUint16(), nameAndType: r.readUint16()}

		case nameAndTypeTag:
			cp.cpInfo[i] = &NameAndTypeCpInfo{name: r.readUint16(), desc: r.readUint16()}

		case methodHandleTag:
			cp.cpInfo[i] = &MethodHandleCpInfo{kind: r.readByte(), index: r.readUint16()}

		case methodTypeTag:
			cp.cpInfo[i] = r.readUint16()

		case invokeDynTag:
			cp.cpInfo[i] = &InvokeDynamicCpInfo{
				bootstrapMethodAttr: r.readUint16(),
				nameAndType:         r.readUint16(),
			}
		}
	}

	return cp
}

func (cp *ConstantPool) Size() int {
	return len(cp.cpInfo)
}

func (cp *ConstantPool) Utf8(index uint16) *string {
	s, ok := cp.cpInfo[index].(*string)
	if !ok {
		return nil
	}
	return s
}
