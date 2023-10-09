package vm

type (
	Instruction           func(frame *Frame) (CurrentFrameOperation, error)
	CurrentFrameOperation uint8
)

const (
	NoFrameOp CurrentFrameOperation = iota
	ReturnFromFrame
	ThrowFromFrame
)

var InstructionSet [255]Instruction

func init() {
	InstructionSet[0x02] = instrIconst(-1)
	InstructionSet[0x03] = instrIconst(0)
	InstructionSet[0x04] = instrIconst(1)
	InstructionSet[0x05] = instrIconst(2)
	InstructionSet[0x06] = instrIconst(3)
	InstructionSet[0x07] = instrIconst(4)
	InstructionSet[0x08] = instrIconst(5)

	InstructionSet[0x59] = instrDup(1, 0)
	InstructionSet[0x5A] = instrDup(1, 1)
	InstructionSet[0x5B] = instrDup(1, 2)
	InstructionSet[0x5C] = instrDup(2, 0)
	InstructionSet[0x5D] = instrDup(2, 1)
	InstructionSet[0x5E] = instrDup(2, 2)

	InstructionSet[0xAC] = instrReturn
	InstructionSet[0xAD] = instrReturn
	InstructionSet[0xAE] = instrReturn
	InstructionSet[0xAF] = instrReturn
	InstructionSet[0xB0] = instrReturn
	InstructionSet[0xB1] = instrReturnVoid

	InstructionSet[0xB8] = func(frame *Frame) (CurrentFrameOperation, error) { // invokestatic
		return NoFrameOp, nil
	}
}

func instrIconst(n int) Instruction {
	return func(frame *Frame) (CurrentFrameOperation, error) {
		frame.PushOperand(n)
		return NoFrameOp, nil
	}
}

func instrDup(n, x int) Instruction {
	return func(frame *Frame) (CurrentFrameOperation, error) {
		top := make([]interface{}, n+x)
		for i := len(top) - 1; i >= 0; i-- {
			top[i] = frame.PopOperand()
		}

		for i := 0; i < n; i++ {
			frame.PushOperand(top[x+i])
		}

		for _, v := range top {
			frame.PushOperand(v)
		}

		return NoFrameOp, nil
	}
}

func instrReturn(frame *Frame) (CurrentFrameOperation, error) {
	ret := frame.PopOperand()
	frame.ClearOperand()
	frame.PushOperand(ret)
	return ReturnFromFrame, nil
}

func instrReturnVoid(frame *Frame) (CurrentFrameOperation, error) {
	frame.ClearOperand()
	return ReturnFromFrame, nil
}
