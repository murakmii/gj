package class_file

import "github.com/murakmii/gojiai/support"

type (
	CodeAttr struct {
		maxStack        uint16
		maxLocals       uint16
		code            []byte
		exceptionTables []*ExceptionTable
		attributes      []interface{}
	}

	ConstantValueAttr uint16

	DeprecatedAttr struct{}
	ExceptionsAttr []uint16
	SignatureAttr  uint16
	SourceFileAttr uint16
	SyntheticAttr  struct{}

	ExceptionTable struct {
		startPC   uint16
		endPC     uint16
		handlerPC uint16
		catchType uint16
	}

	RuntimeVisibleAnnotationsAttr struct {
		rawBytes []byte
	}

	RuntimeVisibleParameterAnnotationsAttr struct {
		rawBytes []byte
	}

	EnclosingMethodAttr struct {
		class  uint16
		method uint16
	}

	InnerClassesAttr []*InnerClassInfo

	InnerClassInfo struct {
		class      uint16
		outer      uint16
		name       uint16
		accessFlag AccessFlag
	}

	LineNumberTableAttr map[uint16]uint16
)

const (
	annotationDefaultAttr                    = "AnnotationDefault"
	bootstrapMethodsAttr                     = "BootstrapMethods"
	codeAttr                                 = "Code"
	constantValueAttr                        = "ConstantValue"
	deprecatedAttr                           = "Deprecated"
	enclosingMethodAttr                      = "EnclosingMethod"
	exceptionsAttr                           = "Exceptions"
	innerClassesAttr                         = "InnerClasses"
	lineNumberTableAttr                      = "LineNumberTable"
	localVariableTableAttr                   = "LocalVariableTable"
	localVariableTypeTableAttr               = "LocalVariableTypeTable"
	methodParametersAttr                     = "MethodParameters"
	runtimeInvisibleAnnotationsAttr          = "RuntimeInvisibleAnnotations"
	runtimeInvisibleParameterAnnotationsAttr = "RuntimeInvisibleParameterAnnotations"
	runtimeInvisibleTypeAnnotationsAttr      = "RuntimeInvisibleTypeAnnotations"
	runtimeVisibleAnnotationsAttr            = "RuntimeVisibleAnnotations"
	runtimeVisibleParameterAnnotationsAttr   = "RuntimeVisibleParameterAnnotations"
	runtimeVisibleTypeAnnotationsAttr        = "RuntimeVisibleTypeAnnotations"
	signatureAttr                            = "Signature"
	sourceDebugExtensionAttr                 = "SourceDebugExtension"
	sourceFileAttr                           = "SourceFile"
	stackMapTableAttr                        = "StackMapTable"
	syntheticAttr                            = "Synthetic"
)

func readAttributes(r *support.ByteSeq, cp *ConstantPool) []interface{} {
	attrs := make([]interface{}, r.ReadUint16())
	for i := 0; i < len(attrs); i++ {
		attrs[i] = readAttribute(r, cp)
	}
	return attrs
}

func readAttribute(r *support.ByteSeq, cp *ConstantPool) interface{} {
	name := cp.Utf8(r.ReadUint16())
	if name == nil {
		return nil
	}

	size := r.ReadUint32()

	switch *name {
	//case annotationDefaultAttr:
	//case bootstrapMethodsAttr:
	case codeAttr:
		return readCodeAttr(r, cp)

	case constantValueAttr:
		return ConstantValueAttr(r.ReadUint16())

	case deprecatedAttr:
		return DeprecatedAttr{}

	case enclosingMethodAttr:
		return &EnclosingMethodAttr{class: r.ReadUint16(), method: r.ReadUint16()}

	case exceptionsAttr:
		attr := ExceptionsAttr(make([]uint16, r.ReadUint16()))
		for i := 0; i < len(attr); i++ {
			attr[i] = r.ReadUint16()
		}
		return attr

	case innerClassesAttr:
		attr := InnerClassesAttr(make([]*InnerClassInfo, r.ReadUint16()))
		for i := range attr {
			attr[i] = &InnerClassInfo{
				class:      r.ReadUint16(),
				outer:      r.ReadUint16(),
				name:       r.ReadUint16(),
				accessFlag: AccessFlag(r.ReadUint16()),
			}
		}
		return attr

	case lineNumberTableAttr:
		n := r.ReadUint16()
		attr := make(LineNumberTableAttr, n)
		for i := uint16(0); i < n; i++ {
			pc := r.ReadUint16()
			attr[pc] = r.ReadUint16()
		}
		return attr

	//case localVariableTableAttr:
	//case localVariableTypeTableAttr:
	//case methodParametersAttr:
	//case runtimeInvisibleAnnotationsAttr:
	//case runtimeInvisibleParameterAnnotationsAttr:
	//case runtimeInvisibleTypeAnnotationsAttr:
	case runtimeVisibleAnnotationsAttr:
		return &RuntimeVisibleAnnotationsAttr{rawBytes: r.ReadBytes(int(size))}

	case runtimeVisibleParameterAnnotationsAttr:
		return &RuntimeVisibleParameterAnnotationsAttr{rawBytes: r.ReadBytes(int(size))}

	//case runtimeVisibleTypeAnnotationsAttr:
	case signatureAttr:
		return SignatureAttr(r.ReadUint16())

	case sourceDebugExtensionAttr:
		r.Skip(int(size)) // ignore
		return nil

	case sourceFileAttr:
		return SourceFileAttr(r.ReadUint16())

	//case stackMapTableAttr:
	case syntheticAttr:
		return SyntheticAttr{}

	default:
		r.Skip(int(size))
		return nil
	}
}

func readCodeAttr(r *support.ByteSeq, cp *ConstantPool) interface{} {
	attr := &CodeAttr{
		maxStack:        r.ReadUint16(),
		maxLocals:       r.ReadUint16(),
		code:            r.ReadBytes(int(r.ReadUint32())),
		exceptionTables: make([]*ExceptionTable, r.ReadUint16()),
	}

	for i := 0; i < len(attr.exceptionTables); i++ {
		attr.exceptionTables[i] = &ExceptionTable{
			startPC:   r.ReadUint16(),
			endPC:     r.ReadUint16(),
			handlerPC: r.ReadUint16(),
			catchType: r.ReadUint16(),
		}
	}

	attr.attributes = readAttributes(r, cp)
	return attr
}

func (ca *CodeAttr) MaxLocals() uint16 {
	return ca.maxLocals
}

func (ca *CodeAttr) Code() []byte {
	return ca.code
}

func (ca *CodeAttr) OverrideCode(code []byte) {
	ca.code = code
}

func (ca *CodeAttr) ExceptionTable() []*ExceptionTable {
	return ca.exceptionTables
}

func (ca *CodeAttr) LineNumberTable() LineNumberTableAttr {
	for _, attr := range ca.attributes {
		if table, ok := attr.(LineNumberTableAttr); ok {
			return table
		}
	}
	return nil
}

func (e *ExceptionTable) HandlerStart() uint16 {
	return e.startPC
}

func (e *ExceptionTable) HandlerEnd() uint16 {
	return e.endPC
}

func (e *ExceptionTable) HandlerPC() uint16 {
	return e.handlerPC
}

func (e *ExceptionTable) CatchType() uint16 {
	return e.catchType
}

func (anno *RuntimeVisibleAnnotationsAttr) RawBytes() []byte {
	return anno.rawBytes
}

func (anno *RuntimeVisibleParameterAnnotationsAttr) RawBytes() []byte {
	return anno.rawBytes
}

func (enc *EnclosingMethodAttr) Class() uint16 {
	return enc.class
}

func (enc *EnclosingMethodAttr) Method() uint16 {
	return enc.method
}

func (inner *InnerClassInfo) Class() uint16 {
	return inner.class
}

func (inner *InnerClassInfo) Outer() uint16 {
	return inner.outer
}

func (inner *InnerClassInfo) Name() uint16 {
	return inner.name
}
