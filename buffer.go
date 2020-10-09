package snd

type buffer struct {
	data     []byte
	position int
}

func (b *buffer) rewind(n int) {
	b.position -= n
}

func (b *buffer) skip(n int) {
	b.position += n
}

func (b *buffer) put(data []byte) {
	copy(b.data[b.position:], data)
	b.position += len(data)
}

func (b *buffer) p(v byte) {
	b.data[b.position] = v
	b.position++
}

func (b *buffer) p8(v int) {
	b.data[b.position] = byte(v)
	b.position++
}

func (b *buffer) u8() (v byte) {
	v = b.data[b.position]
	b.position++
	return
}

func (b *buffer) s8() int8 {
	return int8(b.u8())
}

func (b *buffer) p16(v int) {
	b.p8(v >> 8)
	b.p8(v)
}

func (b *buffer) p16le(v int) {
	b.p8(v)
	b.p8(v >> 8)
}

func (b *buffer) u16() uint16 {
	return (uint16(b.u8()) << 8) | uint16(b.u8())
}

func (b *buffer) s16() int16 {
	return int16(b.u16())
}

func (b *buffer) p24(v int) {
	b.p8(v >> 16)
	b.p8(v >> 8)
	b.p8(v)
}

func (b *buffer) u24() uint32 {
	return (uint32(b.u8()) << 16) | (uint32(b.u8()) << 8) | uint32(b.u8())
}

func (b *buffer) s24() int32 {
	v := int32(b.u24())
	if v >= 0x1000000 {
		v -= 0x800000
	}
	return v
}

func (b *buffer) p32(v int) {
	b.p8(v >> 24)
	b.p8(v >> 16)
	b.p8(v >> 8)
	b.p8(v)
}

func (b *buffer) p32le(v int) {
	b.p8(v)
	b.p8(v >> 8)
	b.p8(v >> 16)
	b.p8(v >> 24)
}

func (b *buffer) u32() uint32 {
	return (uint32(b.u8()) << 24) | (uint32(b.u8()) << 16) | (uint32(b.u8()) << 8) | uint32(b.u8())
}

func (b *buffer) i32() int32 {
	return int32(b.u32())
}

func (b *buffer) p64(v int64) {
	b.p32(int(v >> 32))
	b.p32(int(v))
}

func (b *buffer) u64() uint64 {
	var v uint64
	v |= uint64(b.u32()) << 32
	v |= uint64(b.u32())
	return v
}

func (b *buffer) s64() int64 {
	return int64(b.u64())
}

func (b *buffer) smart() int {
	if b.data[b.position] < 128 {
		return int(b.u8()) - 64
	} else {
		return int(b.u16()) - 49152
	}
}

func (b *buffer) usmart() int {
	if b.data[b.position] < 128 {
		return int(b.u8())
	} else {
		return int(b.u16()) - 32768
	}
}
