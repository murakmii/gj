package class_file

import "github.com/murakmii/gojiai/support"

type (
	// Code属性を表す構造体
	CodeAttr struct {
		maxStack  uint16
		maxLocals uint16
		code      []byte
	}
)

const (
	// Code属性を示す文字列
	codeAttr = "Code"
)

// 各種attributesフィールドを読み取る。
// 属性の種別を判断するため定数プールにアクセスする必要があるので、引数で定数プールを表す構造体ConstantPoolを取る。
func readAttributes(r *support.ByteSeq, cp *ConstantPool) []interface{} {
	attrs := make([]interface{}, r.ReadUint16()) // attributes_countフィールドの読み取り
	for i := 0; i < len(attrs); i++ {
		attrs[i] = readAttribute(r, cp)
	}
	return attrs
}

func readAttribute(r *support.ByteSeq, cp *ConstantPool) interface{} {
	// attribute_name_indexフィールドから定数プール中のインデックスを読み取り、
	// 既に実装したUtf8メソッドで属性の名前を取得。
	name := cp.Utf8(r.ReadUint16())

	// attribute_lengthフィールドの読み取り
	size := r.ReadUint32()

	switch name {
	case codeAttr:
		return readCodeAttr(r, cp)
	default:
		// Code属性以外の属性については、詳細を知らなくてもサイズが分かっているため単にスキップ可能
		r.Skip(int(size))
		return nil
	}
}

// Code属性の読み取り
func readCodeAttr(r *support.ByteSeq, cp *ConstantPool) interface{} {
	// max_stack、max_locals、codeフィールドを読み取る
	attr := &CodeAttr{
		maxStack:  r.ReadUint16(),
		maxLocals: r.ReadUint16(),
		code:      r.ReadBytes(int(r.ReadUint32())),
	}

	// exception_tableフィールドをスキップ。
	// (exception_table_lengthフィールドが示す分、u2フィールド4つ=8バイトをスキップ)
	r.Skip(int(r.ReadUint16() * 8))

	// 属性をスキップ(読み取るが戻り値を捨てる)
	readAttributes(r, cp)

	return attr
}

func (ca *CodeAttr) MaxLocals() uint16 {
	return ca.maxLocals
}

func (ca *CodeAttr) Code() []byte {
	return ca.code
}
