package bitpack

type BitSet64 uint64

type BitPosition = uint8

type (
	Packed8  [1]byte // 8 бит
	Packed16 [2]byte // 16 бит
	Packed24 [3]byte // 24 бита
	Packed32 [4]byte // 32 бита
	Packed40 [5]byte // 40 бит
	Packed48 [6]byte // 48 бит
	Packed56 [7]byte // 56 бит
	Packed64 [8]byte // 64 бита
)

type UnsignedInteger interface {
	~uint8 | ~uint16 | ~uint32 | ~uint64
}

type SignedInteger interface {
	~int8 | ~int16 | ~int32 | ~int64
}
