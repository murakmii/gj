package zip

import (
	"github.com/murakmii/gojiai/vm"
	"hash/crc32"
)

func init() {
	class := "java/util/zip/CRC32"

	vm.NativeMethods.Register(class, "updateBytes", "(I[BII)I", func(thread *vm.Thread, args []interface{}) error {
		crc := uint32(args[0].(int32))
		b := args[1].(*vm.Instance).AsArray()
		off := int(args[2].(int32))
		size := args[3].(int32)

		bytes := make([]byte, size)
		for i := range bytes {
			bytes[i] = byte(b[off+i].(int32))
		}

		thread.CurrentFrame().PushOperand(int32(crc32.Update(crc, crc32.IEEETable, bytes)))
		return nil
	})
}
