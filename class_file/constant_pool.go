package class_file

import (
	"fmt"
	"github.com/murakmii/gj/util"
	"math"
	"strings"
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

func readCP(r *util.BinReader) *ConstantPool {
	cpCount := r.ReadUint16()
	cp := &ConstantPool{cpInfo: make([]interface{}, cpCount)}

	// cp.cpInfo[0] won't be used(cp_info entries indexed from 1)
	for i := uint16(1); i < cpCount; i++ {
		switch r.ReadByte() {
		case utf8Tag:
			s := string(r.ReadBytes(int(r.ReadUint16())))
			cp.cpInfo[i] = &s

		case intTag:
			cp.cpInfo[i] = int(r.ReadUint32())

		case floatTag:
			cp.cpInfo[i] = math.Float32frombits(r.ReadUint32())

		case longTag:
			cp.cpInfo[i] = int64(r.ReadUint64())
			i++ // long occupies 2 entries

		case doubleTag:
			cp.cpInfo[i] = math.Float64frombits(r.ReadUint64())
			i++ // double occupies 2 entries

		case classTag, strTag:
			cp.cpInfo[i] = r.ReadUint16()

		case fieldRefTag, methodRefTag, ifMethodRefTag:
			cp.cpInfo[i] = &ReferenceCpInfo{class: r.ReadUint16(), nameAndType: r.ReadUint16()}

		case nameAndTypeTag:
			cp.cpInfo[i] = &NameAndTypeCpInfo{name: r.ReadUint16(), desc: r.ReadUint16()}

		case methodHandleTag:
			cp.cpInfo[i] = &MethodHandleCpInfo{kind: r.ReadByte(), index: r.ReadUint16()}

		case methodTypeTag:
			cp.cpInfo[i] = r.ReadUint16()

		case invokeDynTag:
			cp.cpInfo[i] = &InvokeDynamicCpInfo{
				bootstrapMethodAttr: r.ReadUint16(),
				nameAndType:         r.ReadUint16(),
			}
		}
	}

	return cp
}

func (cp *ConstantPool) ClassInfo(index uint16) *string {
	classInfo, ok := cp.cpInfo[index].(uint16)
	if !ok {
		return nil
	}
	return cp.Utf8(classInfo)
}

func (cp *ConstantPool) Utf8(index uint16) *string {
	s, ok := cp.cpInfo[index].(*string)
	if !ok {
		return nil
	}
	return s
}

func (cp *ConstantPool) String() string {
	sb := &strings.Builder{}
	sb.WriteString(fmt.Sprintf("Entries: %d\n", len(cp.cpInfo)-1))

	for i := 1; i < len(cp.cpInfo); i++ {
		sb.WriteString(fmt.Sprintf("[%4d] ", i))

		cpInfo := cp.cpInfo[i]
		switch ci := cpInfo.(type) {
		case *string:
			sb.WriteString(fmt.Sprintf("UTF-8: '%s'", *ci))
		case int, float32:
			sb.WriteString(fmt.Sprintf("%T: %v", ci, ci))
		case int64, float64:
			sb.WriteString(fmt.Sprintf("%T: %v", ci, ci))
			i++ // long/double occupies 2 entries
		case *NameAndTypeCpInfo:
			sb.WriteString(fmt.Sprintf("NameAndType: Name=%d, Type=%d", ci.name, ci.desc))
		case uint16:
			sb.WriteString(fmt.Sprintf("Class/Str/MethodType: %d", ci))
		case *ReferenceCpInfo:
			sb.WriteString(fmt.Sprintf("Field/Method/InterfaceMethodRef: Class=%d, NameAndType=%d", ci.class, ci.nameAndType))
		case *MethodHandleCpInfo:
			sb.WriteString(fmt.Sprintf("MethodHandle: Kind=%d, Index=%d", ci.kind, ci.index))
		case *InvokeDynamicCpInfo:
			sb.WriteString(fmt.Sprintf("InvokeDynamic: Bootstrap=%d, NameAndType=%d", ci.bootstrapMethodAttr, ci.nameAndType))
		}

		sb.WriteString("\n")
	}

	return sb.String()
}
