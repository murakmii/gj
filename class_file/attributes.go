package class_file

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

func readAttributes(r *reader, cp *ConstantPool) []interface{} {
	attrs := make([]interface{}, r.readUint16())
	for i := 0; i < len(attrs); i++ {
		attrs[i] = readAttribute(r, cp)
	}
	return attrs
}

func readAttribute(r *reader, cp *ConstantPool) interface{} {
	name := cp.Utf8(r.readUint16())
	if name == nil {
		return nil
	}

	size := r.readUint32()

	switch *name {
	//case annotationDefaultAttr:
	//case bootstrapMethodsAttr:
	case codeAttr:
		return readCodeAttr(r, cp)

	case constantValueAttr:
		return ConstantValueAttr(r.readUint16())

	case deprecatedAttr:
		return DeprecatedAttr{}

	//case enclosingMethodAttr:
	case exceptionsAttr:
		attr := ExceptionsAttr(make([]uint16, r.readUint16()))
		for i := 0; i < len(attr); i++ {
			attr[i] = r.readUint16()
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
		return SignatureAttr(r.readUint16())

	case sourceDebugExtensionAttr:
		r.skip(int(size)) // ignore
		return nil

	case sourceFileAttr:
		return SourceFileAttr(r.readUint16())

	//case stackMapTableAttr:
	case syntheticAttr:
		return SyntheticAttr{}

	default:
		r.skip(int(size))
		return nil
	}
}

func readCodeAttr(r *reader, cp *ConstantPool) interface{} {
	attr := &CodeAttr{
		maxStack:        r.readUint16(),
		maxLocals:       r.readUint16(),
		code:            r.readBytes(int(r.readUint32())),
		exceptionTables: make([]*ExceptionTable, r.readUint16()),
	}

	for i := 0; i < len(attr.exceptionTables); i++ {
		attr.exceptionTables[i] = &ExceptionTable{
			startPC:   r.readUint16(),
			endPC:     r.readUint16(),
			handlerPC: r.readUint16(),
			catchType: r.readUint16(),
		}
	}

	attr.attributes = readAttributes(r, cp)
	return attr
}
