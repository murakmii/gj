package vm

import "sync"

type (
	NativeMemAllocator struct {
		allocated map[int64][]byte
		total     int64
		lock      *sync.Mutex
	}
)

func CreateNativeMemAllocator() *NativeMemAllocator {
	return &NativeMemAllocator{
		allocated: make(map[int64][]byte),
		total:     0,
		lock:      &sync.Mutex{},
	}
}

func (allocator *NativeMemAllocator) Alloc(size int64) int64 {
	allocator.lock.Lock()
	defer allocator.lock.Unlock()

	addr := allocator.total
	allocator.total += size
	allocator.allocated[addr] = make([]byte, size)

	return addr
}

func (allocator *NativeMemAllocator) Ref(addr int64) []byte {
	startAddr, block := allocator.findBlock(addr)
	if startAddr == -1 {
		return nil
	}
	return block[addr-startAddr:]
}

func (allocator *NativeMemAllocator) Free(addr int64) {
	allocator.lock.Lock()
	defer allocator.lock.Unlock()

	delete(allocator.allocated, addr)
}

func (allocator *NativeMemAllocator) findBlock(addr int64) (int64, []byte) {
	for startAddr, block := range allocator.allocated {
		endAddr := startAddr + int64(len(block))
		if startAddr <= addr && addr < endAddr {
			return startAddr, block
		}
	}
	return -1, nil
}
