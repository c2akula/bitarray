package bitarray

import (
	"math"
	"math/bits"
)

type Bit = uint64

const (
	Zero = Bit(iota)
	One
)

// BitArray is an array data structure that compactly stores bits.
// Bits externally represented as `bool` are stored internally as `uint64`s.
// The total number of bits stored is set at creation and is immutable.
type BitArray struct {
	// buf is a backing array that bits writes into by default when the no. of bits requested to allocate is
	// < 512. Only if more is asked, we'll skip buf and allocate directly into bits
	buf  [8]Bit
	bits []Bit
	n    int // no. of bits
}

// New creates a new BitArray of `n` bits. If n <= 512, no allocation is done.
func New(n int) (ba BitArray) {
	nblk := nbitsToNblks(n)
	ba.n = n
	ba.bits = ba.buf[:]
	if nblk <= 8 {
		ba.bits = ba.buf[:nblk]
		return
	}
	ba.bits = append(ba.bits, make([]Bit, nblk-len(ba.buf))...)
	return
}

// Copy copies src into dst.
func Copy(dst, src *BitArray) {
	if dst != src && src != nil {
		if src.n != dst.n {
			panic("size of bit arrays must be the same for copy")
		}

		if dst.n == 0 {
			// nothing to do here, since the source `oa` has nothing to copy from
			return
		}

		copy(dst.bits, src.bits)
	}
}

// FromStr creates a BitArray from a bit string
func FromStr(bs string) BitArray {
	ba := New(len(bs))
	for i, b := range bs {
		if b == '1' {
			ba.Set(i)
		}
	}
	return ba
}

// FromUint64 creates a BitArray from the bit representation of u.
func FromUint64(u uint64) BitArray {
	ba := New(64)
	ba.bits[0] = u
	return ba
}

// Size returns the no. of bits stored.
func (ba *BitArray) Size() int { return ba.n }

// Set sets the bit at position k.
func (ba *BitArray) Set(k int) { bi, si := biandsi(k); set(&ba.bits[bi], si) }

// SetAll sets all the bits.
func (ba *BitArray) SetAll() {
	for i := range ba.bits {
		ba.bits[i] = math.MaxUint64
	}
}

// Clr clears the bit at position k.
func (ba *BitArray) Clr(k int) { bi, si := biandsi(k); clr(&ba.bits[bi], si) }

// ClrAll clears all the bits.
func (ba *BitArray) ClrAll() {
	for i := range ba.bits {
		ba.bits[i] = 0
	}
}

// ChkSet returns the value of the bit at position k before setting it.
func (ba *BitArray) ChkSet(k int) (b bool) {
	bi, si := biandsi(k)
	u := &ba.bits[bi]
	b = chk(*u, si) != 0
	if !b {
		set(u, si)
	}
	return
}

// ChkClr returns the value of the bit at position k before clearing it.
func (ba *BitArray) ChkClr(k int) (b bool) {
	bi, si := biandsi(k)
	u := &ba.bits[bi]
	b = chk(*u, si) != 0
	if b {
		clr(u, si)
	}
	return
}

// Tgl toggles the bit at position k.
func (ba *BitArray) Tgl(k int) {
	bi, si := biandsi(k)
	ba.bits[bi] ^= 1 << si
}

// Cnt returns the number of set bits.
func (ba *BitArray) Cnt() (n int) {
	for _, b := range ba.bits {
		n += bits.OnesCount64(b)
	}
	return
}

// Chk returns the value of the bit at position k.
func (ba *BitArray) Chk(k int) bool {
	bi, si := biandsi(k)
	return chk(ba.bits[bi], si) != 0
}

// Put sets the value of the bit at position k to v.
func (ba *BitArray) Put(k int, v Bit) {
	bi, si := biandsi(k)
	put(&ba.bits[bi], si, v)
}

// Swap swaps the value of bit at position k with v. On return, v contains the old value.
func (ba *BitArray) Swap(k int, b *Bit) {
	bi, si := biandsi(k)
	t := &ba.bits[bi]
	ob := chk(*t, si)
	if ob == *b {
		return
	}
	put(t, si, ob)
	*b = ob
}

func (ba *BitArray) String() string {
	sb := make([]byte, ba.n)
	for i := range sb {
		sb[i] = '0'
		if ba.Chk(i) {
			sb[i] = '1'
		}
	}
	return string(sb)
}

func nbitsToNblks(n int) int { return int(math.Ceil(float64(n) / 64)) }

func set(u *uint64, si uint64)        { *u |= 1 << si }
func clr(u *uint64, si uint64)        { *u &= ^(1 << si) }
func chk(u uint64, si uint64) Bit     { return (u >> si) & 1 }
func put(u *uint64, si uint64, b Bit) { *u = (*u & ^(1 << si)) | (b << si) }

func biandsi(k int) (uint64, uint64) {
	i := uint64(k)
	return i / 64, i % 64
}
