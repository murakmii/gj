package class_file

import (
	"fmt"
	"github.com/murakmii/gojiai/support"
	"io"
	"strings"
)

type (
	ClassFile struct {
		cp      *ConstantPool
		methods []*MethodInfo
	}

	MethodInfo struct {
		accessFlag uint16
		name       string
		desc       string
		attributes []interface{}
	}
)

const magicNumber = 0xCAFEBABE

func ReadClassFile(cfReader io.Reader) (*ClassFile, error) {
	r, err := support.NewByteSeq(cfReader)
	if err != nil {
		return nil, err
	}

	if r.ReadUint32() != magicNumber {
		return nil, fmt.Errorf("not java class file")
	}

	r.Skip(4)

	class := &ClassFile{cp: readCP(r)}

	r.Skip(6)
	r.Skip(int(r.ReadUint16() * 2))

	readMethodInfo(r, class.cp)
	class.methods = readMethodInfo(r, class.cp)

	return class, nil
}

func readMethodInfo(r *support.ByteSeq, cp *ConstantPool) []*MethodInfo {
	count := r.ReadUint16()
	methods := make([]*MethodInfo, count)

	for i := uint16(0); i < count; i++ {
		methods[i] = &MethodInfo{
			accessFlag: r.ReadUint16(),
			name:       cp.Utf8(r.ReadUint16()),
			desc:       cp.Utf8(r.ReadUint16()),
			attributes: readAttributes(r, cp),
		}
	}

	return methods
}

func (c *ClassFile) FindMethod(name, desc string) *MethodInfo {
	for _, method := range c.methods {
		if method.name == name && method.desc == desc {
			return method
		}
	}
	return nil
}

func (m *MethodInfo) Code() *CodeAttr {
	for _, attr := range m.attributes {
		if code, ok := attr.(*CodeAttr); ok {
			return code
		}
	}
	return nil
}

// デコードしたメソッド一覧を返す
func (c *ClassFile) Methods() []*MethodInfo {
	return c.methods
}

// メソッドの文字列表現を返す(名前とシグネチャ)
func (m *MethodInfo) String() string {
	return m.name + m.desc
}

func (m *MethodInfo) NumArgs() int {
	return strings.Index(m.desc[1:], ")")
}
