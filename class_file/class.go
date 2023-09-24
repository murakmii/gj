package class_file

import (
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
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

func (c *Class) ConstantPool() *ConstantPool {
	return c.cp
}

func (c *Class) String() string {
	sb := &strings.Builder{}
	sb.WriteString("# Class file\n\n")

	sb.WriteString("## Constant pool\n\n")
	sb.WriteString(c.ConstantPool().String())
	sb.WriteByte('\n')

	sb.WriteString("## This:Super\n\n")
	sb.WriteString(fmt.Sprintf("%d:%d\n\n", c.this, c.super))

	if len(c.interfaces) > 0 {
		sb.WriteString("## Implement interfaces\n\n")
		for _, i := range c.interfaces {
			sb.WriteString(strconv.Itoa(int(i)))
			sb.WriteByte(' ')
		}
		sb.WriteByte('\n')
	}

	if len(c.fields) > 0 {
		sb.WriteString("## Fields\n\n")
		for _, f := range c.fields {
			sb.WriteString(fmt.Sprintf("* Name: %d, Desc: %d\n", f.name, f.desc))
		}
		sb.WriteByte('\n')
	}

	if len(c.methods) > 0 {
		sb.WriteString("## Methods\n\n")
		for _, m := range c.methods {
			sb.WriteString(fmt.Sprintf("* Name: %d, Desc: %d\n", m.name, m.desc))

			for _, attr := range m.attributes {
				code, ok := attr.(*CodeAttr)
				if ok {
					dumpCode(sb, code)
					sb.WriteByte('\n')
				}
			}
		}
	}

	return sb.String()
}

func dumpCode(sb *strings.Builder, code *CodeAttr) {
	r := &reader{bytes: code.code}
	for r.remain() > 0 {
		sb.WriteString("   ")

		opCode := r.readByte()
		switch opCode {
		case 0x32:
			sb.WriteString("aaload")
		case 0x53:
			sb.WriteString("aastore")
		case 0x01:
			sb.WriteString("aconst_null")
		case 0x19:
			sb.WriteString(fmt.Sprintf("aload local:%d", r.readByte()))
		case 0x2A, 0x2B, 0x2C, 0x2D:
			sb.WriteString(fmt.Sprintf("aload_%d", opCode-0x2A))
		case 0xBD:
			sb.WriteString(fmt.Sprintf("anewarray cp:%d", r.readUint16()))
		case 0xB0:
			sb.WriteString("areturn")
		case 0xBE:
			sb.WriteString("arraylength")
		case 0x3A:
			sb.WriteString(fmt.Sprintf("astore local:%d", r.readByte()))
		case 0x4B, 0x4C, 0x4D, 0x4E:
			sb.WriteString(fmt.Sprintf("astore_%d", opCode-0x4B))
		case 0xBF:
			sb.WriteString("athrow")

		case 0x33:
			sb.WriteString("baload")
		case 0x54:
			sb.WriteString("bastore")
		case 0x10:
			sb.WriteString(fmt.Sprintf("bipush %d", r.readByte()))

		case 0x34:
			sb.WriteString("caload")
		case 0x55:
			sb.WriteString("castore")

		case 0xC0:
			sb.WriteString(fmt.Sprintf("checkcast cp:%d", r.readUint16()))

		case 0x90:
			sb.WriteString("d2f")
		case 0x8E:
			sb.WriteString("d2i")
		case 0x8F:
			sb.WriteString("d2l")
		case 0x63:
			sb.WriteString("dadd")
		case 0x31:
			sb.WriteString("daload")
		case 0x52:
			sb.WriteString("dastore")
		case 0x98:
			sb.WriteString("dcmpg")
		case 0x97:
			sb.WriteString("dcmpl")
		case 0x0E:
			sb.WriteString("dconst_0")
		case 0x0F:
			sb.WriteString("dconst_1")
		case 0x6F:
			sb.WriteString("ddiv")
		case 0x18:
			sb.WriteString(fmt.Sprintf("dload local:%d", r.readByte()))
		case 0x26, 0x27, 0x28, 0x29:
			sb.WriteString(fmt.Sprintf("dload_%d", opCode-0x26))
		case 0x6B:
			sb.WriteString("dmul")
		case 0x77:
			sb.WriteString("dneg")
		case 0x73:
			sb.WriteString("drem")
		case 0xAF:
			sb.WriteString("dreturn")
		case 0x39:
			sb.WriteString(fmt.Sprintf("dstore local:%d", r.readByte()))
		case 0x47, 0x48, 0x49, 0x4A:
			sb.WriteString(fmt.Sprintf("dstore_%d", opCode-0x47))
		case 0x67:
			sb.WriteString("dsub")

		case 0x59:
			sb.WriteString("dup")
		case 0x5A:
			sb.WriteString("dup_x1")
		case 0x5B:
			sb.WriteString("dup_x2")
		case 0x5C:
			sb.WriteString("dup2")
		case 0x5D:
			sb.WriteString("dup2_x1")
		case 0x5E:
			sb.WriteString("dup2_x2")

		case 0x8D:
			sb.WriteString("f2d")
		case 0x8B:
			sb.WriteString("f2i")
		case 0x8C:
			sb.WriteString("f2l")
		case 0x62:
			sb.WriteString("fadd")
		case 0x30:
			sb.WriteString("faload")
		case 0x51:
			sb.WriteString("fastore")
		case 0x96:
			sb.WriteString("fcmpg")
		case 0x95:
			sb.WriteString("fcmpl")
		case 0x0B, 0x0C, 0x0D:
			sb.WriteString(fmt.Sprintf("fconst_%d", opCode-0x0B))
		case 0x6E:
			sb.WriteString("fdiv")
		case 0x17:
			sb.WriteString(fmt.Sprintf("fload local:%d", r.readByte()))
		case 0x22, 0x23, 0x24, 0x25:
			sb.WriteString(fmt.Sprintf("fload_%d", opCode-0x22))
		case 0x6A:
			sb.WriteString("fmul")
		case 0x76:
			sb.WriteString("fneg")
		case 0x72:
			sb.WriteString("frem")
		case 0xAE:
			sb.WriteString("freturn")
		case 0x38:
			sb.WriteString(fmt.Sprintf("fstore local:%d", r.readByte()))
		case 0x43, 0x44, 0x45, 0x46:
			sb.WriteString(fmt.Sprintf("fstore_%d", opCode-0x43))
		case 0x66:
			sb.WriteString("fsub")

		case 0xB4:
			sb.WriteString(fmt.Sprintf("getfield cp:%d", r.readUint16()))
		case 0xB2:
			sb.WriteString(fmt.Sprintf("getstatic cp:%d", r.readUint16()))

		case 0xA7:
			sb.WriteString(fmt.Sprintf("goto pc:%d", r.readUint16()))
		case 0xC8:
			sb.WriteString(fmt.Sprintf("goto_w pc:%d", r.readUint32()))

		case 0x91:
			sb.WriteString("i2b")
		case 0x92:
			sb.WriteString("i2c")
		case 0x87:
			sb.WriteString("i2d")
		case 0x86:
			sb.WriteString("i2f")
		case 0x85:
			sb.WriteString("i2l")
		case 0x93:
			sb.WriteString("i2s")
		case 0x60:
			sb.WriteString("iadd")
		case 0x2E:
			sb.WriteString("iaload")
		case 0x7E:
			sb.WriteString("iand")
		case 0x4F:
			sb.WriteString("iastore")
		case 0x02:
			sb.WriteString("iconst_m1")
		case 0x03, 0x04, 0x05, 0x06, 0x07, 0x08:
			sb.WriteString(fmt.Sprintf("iconst_%d", opCode-0x03))
		case 0x6C:
			sb.WriteString("idiv")

		case 0xA5:
			sb.WriteString(fmt.Sprintf("if_acmpeq pc:%d", r.readUint16()))
		case 0xA6:
			sb.WriteString(fmt.Sprintf("if_acmpne pc:%d", r.readUint16()))
		case 0x9F:
			sb.WriteString(fmt.Sprintf("if_icmpeq pc:%d", r.readUint16()))
		case 0xA0:
			sb.WriteString(fmt.Sprintf("if_acmpne pc:%d", r.readUint16()))
		case 0xA1:
			sb.WriteString(fmt.Sprintf("if_acmplt pc:%d", r.readUint16()))
		case 0xA2:
			sb.WriteString(fmt.Sprintf("if_acmpge pc:%d", r.readUint16()))
		case 0xA3:
			sb.WriteString(fmt.Sprintf("if_acmpgt pc:%d", r.readUint16()))
		case 0xA4:
			sb.WriteString(fmt.Sprintf("if_acmple pc:%d", r.readUint16()))

		case 0x99:
			sb.WriteString(fmt.Sprintf("ifeq pc:%d", r.readUint16()))
		case 0x9A:
			sb.WriteString(fmt.Sprintf("ifne pc:%d", r.readUint16()))
		case 0x9B:
			sb.WriteString(fmt.Sprintf("iflt pc:%d", r.readUint16()))
		case 0x9C:
			sb.WriteString(fmt.Sprintf("ifge pc:%d", r.readUint16()))
		case 0x9D:
			sb.WriteString(fmt.Sprintf("ifgt pc:%d", r.readUint16()))
		case 0x9E:
			sb.WriteString(fmt.Sprintf("ifle pc:%d", r.readUint16()))

		case 0xC7:
			sb.WriteString(fmt.Sprintf("ifnonnull pc:%d", r.readUint16()))
		case 0xC6:
			sb.WriteString(fmt.Sprintf("ifnull pc:%d", r.readUint16()))

		case 0x84:
			sb.WriteString(fmt.Sprintf("iinc local:%d, incr:%d", r.readByte(), r.readByte()))
		case 0x15:
			sb.WriteString(fmt.Sprintf("iload local:%d", r.readByte()))
		case 0x1A, 0x1B, 0x1C, 0x1D:
			sb.WriteString(fmt.Sprintf("iload_%d", opCode-0x1A))
		case 0x68:
			sb.WriteString("imul")
		case 0x74:
			sb.WriteString("ineg")

		case 0xC1:
			sb.WriteString(fmt.Sprintf("instanceof cp:%d", r.readUint16()))

		case 0xBA:
			sb.WriteString(fmt.Sprintf("invokedynamic cp:%d", r.readUint16()))
			r.skip(2)
		case 0xB9:
			sb.WriteString(fmt.Sprintf("invokeinterface cp:%d, count:%d", r.readUint16(), r.readByte()))
			r.skip(1)
		case 0xB7:
			sb.WriteString(fmt.Sprintf("invokespecial cp:%d", r.readUint16()))
		case 0xB8:
			sb.WriteString(fmt.Sprintf("invokestatic cp:%d", r.readUint16()))
		case 0xB6:
			sb.WriteString(fmt.Sprintf("invokevirtual cp:%d", r.readUint16()))

		case 0x80:
			sb.WriteString("ior")
		case 0x70:
			sb.WriteString("irem")
		case 0xAC:
			sb.WriteString("ireturn")
		case 0x78:
			sb.WriteString("ishl")
		case 0x7A:
			sb.WriteString("ishr")
		case 0x36:
			sb.WriteString(fmt.Sprintf("istore local:%d", r.readByte()))
		case 0x3B, 0x3C, 0x3D, 0x3E:
			sb.WriteString(fmt.Sprintf("istore_%d", opCode-0x3B))
		case 0x64:
			sb.WriteString("isub")
		case 0x7C:
			sb.WriteString("iushr")
		case 0x82:
			sb.WriteString("ixor")

		case 0xA8:
			sb.WriteString(fmt.Sprintf("jsr pc:%d", r.readUint16()))
		case 0xC9:
			sb.WriteString(fmt.Sprintf("jsr_w pc:%d", r.readUint32()))

		case 0x8A:
			sb.WriteString("l2d")
		case 0x89:
			sb.WriteString("l2f")
		case 0x88:
			sb.WriteString("l2i")
		case 0x61:
			sb.WriteString("ladd")
		case 0x2F:
			sb.WriteString("laload")
		case 0x7F:
			sb.WriteString("land")
		case 0x50:
			sb.WriteString("lastore")
		case 0x94:
			sb.WriteString("lcmp")
		case 0x09:
			sb.WriteString("lconst_0")
		case 0x0A:
			sb.WriteString("lconst_1")

		case 0x12:
			sb.WriteString(fmt.Sprintf("ldc cp:%d", r.readByte()))
		case 0x13:
			sb.WriteString(fmt.Sprintf("ldc_w cp:%d", r.readUint16()))
		case 0x14:
			sb.WriteString(fmt.Sprintf("ldc2_w cp:%d", r.readUint16()))

		case 0x6D:
			sb.WriteString("ldiv")
		case 0x16:
			sb.WriteString(fmt.Sprintf("lload local:%d+1", r.readByte()))
		case 0x1E, 0x1F, 0x20, 0x21:
			sb.WriteString(fmt.Sprintf("lload_%d", opCode-0x1E))
		case 0x69:
			sb.WriteString("lmul")
		case 0x75:
			sb.WriteString("lneg")

		case 0xAB:
			r.skipToAlign(4)
			sb.WriteString(fmt.Sprintf("lookupswitch default:%d+1", int(r.readUint32())))
			for i := r.readUint32(); i > 0; i-- {
				sb.WriteString(fmt.Sprintf("    * %d -> %d", int(r.readUint32()), int(r.readUint32())))
				if i != 1 {
					sb.WriteByte('\n')
				}
			}

		case 0x81:
			sb.WriteString("lor")
		case 0x71:
			sb.WriteString("lrem")
		case 0xAD:
			sb.WriteString("lreturn")
		case 0x79:
			sb.WriteString("lshl")
		case 0x7B:
			sb.WriteString("lshr")
		case 0x37:
			sb.WriteString(fmt.Sprintf("lstore local:%d+1", r.readByte()))
		case 0x3F, 0x40, 0x41, 0x42:
			sb.WriteString(fmt.Sprintf("lstore_%d", opCode-0x3F))
		case 0x65:
			sb.WriteString("lsub")
		case 0x7D:
			sb.WriteString("lushr")
		case 0x83:
			sb.WriteString("lxor")

		case 0xC2:
			sb.WriteString("monitorenter")
		case 0xC3:
			sb.WriteString("monitorexit")

		case 0xC5:
			sb.WriteString(fmt.Sprintf("multianewarray cp:%d, dimensions:%d", r.readUint16(), r.readByte()))

		case 0xBB:
			sb.WriteString(fmt.Sprintf("new cp:%d", r.readUint16()))
		case 0xBC:
			sb.WriteString(fmt.Sprintf("newarray type:%d", r.readByte()))

		case 0x00:
			sb.WriteString("nop")
		case 0x57:
			sb.WriteString("pop")
		case 0x58:
			sb.WriteString("pop2")

		case 0xB5:
			sb.WriteString(fmt.Sprintf("putfield cp:%d", r.readUint16()))
		case 0xB3:
			sb.WriteString(fmt.Sprintf("putstatic cp:%d", r.readUint16()))

		case 0xA9:
			sb.WriteString(fmt.Sprintf("ret local:%d", r.readByte()))
		case 0xB1:
			sb.WriteString("return")

		case 0x35:
			sb.WriteString("saload")
		case 0x56:
			sb.WriteString("sastore")
		case 0x11:
			sb.WriteString(fmt.Sprintf("sipush value:%d", int(r.readUint16())))

		case 0x5F:
			sb.WriteString("swap")

		case 0xAA:
			r.skipToAlign(4)
			sb.WriteString(fmt.Sprintf("tableswitch default:%d+1", int(r.readUint32())))
			low := int(r.readUint32())
			high := int(r.readUint32())
			for i := 0; i < high-low+1; i++ {
				sb.WriteString(fmt.Sprintf("    * %d -> %d", low+i, int(r.readUint32())))
				if i != 1 {
					sb.WriteByte('\n')
				}
			}

		case 0xC4:
			op := r.readByte()
			sb.WriteString(fmt.Sprintf("wide op:%d, local:%d", op, r.readUint16()))

			if op == 0x84 {
				sb.WriteString(fmt.Sprintf(", incr:%d", r.readUint16()))
			}
		}

		sb.WriteByte('\n')
	}
}
