package zip

import (
	"github.com/murakmii/gj/vm"
	"hash/crc32"
)

func CRC32UpdateBytes(thread *vm.Thread, args []interface{}) error {
	crc := uint32(args[0].(int))
	b := args[1].(*vm.Instance).AsArray()
	off := args[2].(int)
	size := args[3].(int)

	bytes := make([]byte, size)
	for i := range bytes {
		bytes[i] = byte(b[off+i].(int))
	}

	thread.CurrentFrame().PushOperand(int(crc32.Update(crc, crc32.IEEETable, bytes)))
	return nil
}
