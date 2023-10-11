package class_file

import "github.com/murakmii/gj/util"

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

func readAttributes(r *util.BinReader, cp *ConstantPool) []interface{} {
	attrs := make([]interface{}, r.ReadUint16())
	for i := 0; i < len(attrs); i++ {
		attrs[i] = readAttribute(r, cp)
	}
	return attrs
}

func readAttribute(r *util.BinReader, cp *ConstantPool) interface{} {
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

	//case enclosingMethodAttr:
	case exceptionsAttr:
		attr := ExceptionsAttr(make([]uint16, r.ReadUint16()))
		for i := 0; i < len(attr); i++ {
			attr[i] = r.ReadUint16()
		}
		return attr

	//case innerClassesAttr:
	//case lineNumberTableAttr:
	//case localVariableTableAttr:
	//case localVariableTypeTableAttr:
	//case methodParametersAttr:
	//case runtimeInvisibleAnnotationsAttr:
	//case runtimeInvisibleParameterAnnotationsAttr:
	//case runtimeInvisibleTypeAnnotationsAttr:
	//case runtimeVisibleAnnotationsAttr:
	//case runtimeVisibleParameterAnnotationsAttr:
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

func readCodeAttr(r *util.BinReader, cp *ConstantPool) interface{} {
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

func (attr *CodeAttr) MaxLocals() uint16 {
	return attr.maxLocals
}

func (attr *CodeAttr) Code() []byte {
	return attr.code
}

func (attr *CodeAttr) ExceptionTable() []*ExceptionTable {
	return attr.exceptionTables
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
