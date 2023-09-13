package class_file

import (
	"io"
	"os"
)

type (
	Class struct {
		cp         *ConstantPool
		accessFlag AccessFlag
		this       uint16
		super      uint16
		interfaces []uint16
		fields     []*FieldInfo
		methods    []*MethodInfo
		attributes []interface{}
	}

	FieldInfo struct {
		accessFlag AccessFlag
		name       uint16
		desc       uint16
		attributes []interface{}
	}

	MethodInfo struct {
		*FieldInfo
	}
)

const magicNumber = 0xCAFEBABE

func ReadClassFile(cfReader io.Reader) (*Class, error) {
	return readClassFile(cfReader)
}

func OpenClassFile(path string) (*Class, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	return readClassFile(f)
}

func readClassFile(cfReader io.Reader) (*Class, error) {
	r, err := open(cfReader)
	if err != nil {
		return nil, err
	}

	if r.readUint32() != magicNumber {
		return nil, nil // TODO: error
	}

	r.skip(4) // skip major/minor versions

	class := &Class{}

	class.cp = readCP(r)
	class.accessFlag = AccessFlag(r.readUint16())
	class.this = r.readUint16()
	class.super = r.readUint16()
	class.interfaces = readInterfaces(r)
	class.fields = readFields(r, class.cp)
	class.methods = readMethods(r, class.cp)
	class.attributes = readAttributes(r, class.cp)

	return class, nil
}

func readInterfaces(r *reader) []uint16 {
	ifCount := r.readUint16()
	interfaces := make([]uint16, ifCount)

	for i := uint16(0); i < ifCount; i++ {
		interfaces[i] = r.readUint16()
	}

	return interfaces
}

func readFields(r *reader, cp *ConstantPool) []*FieldInfo {
	fCount := r.readUint16()
	fields := make([]*FieldInfo, fCount)

	for i := uint16(0); i < fCount; i++ {
		fields[i] = &FieldInfo{
			accessFlag: AccessFlag(r.readUint16()),
			name:       r.readUint16(),
			desc:       r.readUint16(),
			attributes: readAttributes(r, cp),
		}
	}

	return fields
}

func readMethods(r *reader, cp *ConstantPool) []*MethodInfo {
	fields := readFields(r, cp)
	methods := make([]*MethodInfo, len(fields))

	for i := 0; i < len(methods); i++ {
		methods[i] = &MethodInfo{fields[i]}
	}

	return methods
}
