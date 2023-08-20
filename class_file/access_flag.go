package class_file

type AccessFlag uint16

const (
	PublicFlag       AccessFlag = 0x0001
	PrivateFlag      AccessFlag = 0x0002
	ProtectedFlag    AccessFlag = 0x0004
	StaticFlag       AccessFlag = 0x0008
	FinalFlag        AccessFlag = 0x0010
	SynchronizedFlag AccessFlag = 0x0020
	SuperFlag        AccessFlag = 0x0020
	BridgeFlag       AccessFlag = 0x0040
	VolatileFlag     AccessFlag = 0x0040
	VarArgsFlag      AccessFlag = 0x0080
	TransientFlag    AccessFlag = 0x0080
	NativeFlag       AccessFlag = 0x0100
	InterfaceFlag    AccessFlag = 0x0200
	AbstractFlag     AccessFlag = 0x0400
	StrictFlag       AccessFlag = 0x0800
	SyntheticFlag    AccessFlag = 0x1000
	AnnotationFlag   AccessFlag = 0x2000
	EnumFlag         AccessFlag = 0x4000
)

func (flg AccessFlag) Contain(flags AccessFlag) bool {
	return (flg & flags) == flags
}
