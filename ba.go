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
	bits []Bit
	n    int // no. of bits
}

// New creates a new BitArray of `n` bits.
func New(n int) BitArray {
	return BitArray{
		bits: make([]Bit, int(math.Ceil(float64(n)/64))),
		n:    n,
	}
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

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func swapBlocks(a, b []uint64, bi, m uint64) {
	for ; bi < m; bi++ {
		a[bi], b[bi] = b[bi], a[bi]
	}
}

// SwapRange swaps the bits between bit arrays `a`and `b` starting at position `k`. If the
// sizes are different, then min(a.Size(), b.Size()) number of bits are swapped.
func SwapRange(a, b *BitArray, k int) {
	n := b.n
	if a.n < n {
		n = a.n
	}

	bi, si := biandsi(k)
	m := uint64(min(len(a.bits), len(b.bits)))
	// if the index is a multiple of 64 then we'll just swap wholesale the bit banks
	if si == 0 {
		swapBlocks(a.bits, b.bits, bi, m)
		return
	}

	// index is arbitrary, i.e. si is arbitrary, meaning we need to get to the next boundary
	// so let's start by swapping till the next boundary
	nb := n - k // number of bits
	anum, bnum := &a.bits[bi], &b.bits[bi]
	for ; si < 64; si++ {
		if nb != 0 {
			v := uint64(^(1 << si))
			abit := (*anum >> si) & 1          // bit from a
			bbit := (*bnum >> si) & 1          // bit from b
			*anum = (*anum & v) | (bbit << si) // write bbit into a
			*bnum = (*bnum & v) | (abit << si) // write abit into b
			nb--
		} else {
			return
		}
	}

	// we've swapped bits until we hit the start of the next boundary
	// and now we'll resort to swapping blocks
	bi++
	swapBlocks(a.bits, b.bits, bi, m)
}

// TODO: Add GetRange, SetRange methods

// Size returns the no. of bits stored.
func (ba *BitArray) Size() int { return ba.n }

func biandsi(k int) (uint64, uint64) {
	i := uint64(k)
	return i / 64, i % 64
}

func (ba *BitArray) set1(bi, si uint64) { ba.bits[bi] |= 1 << si }
func (ba *BitArray) set0(bi, si uint64) { ba.bits[bi] &= ^(1 << si) }

// Set sets the bit at position k.
func (ba *BitArray) Set(k int) { ba.set1(biandsi(k)) }

// SetAll sets all the bits.
func (ba *BitArray) SetAll() {
	for i := range ba.bits {
		ba.bits[i] = math.MaxUint64
	}
}

// Clr clears the bit at position k.
func (ba *BitArray) Clr(k int) { ba.set0(biandsi(k)) }

// ClrAll clears all the bits.
func (ba *BitArray) ClrAll() {
	for i := range ba.bits {
		ba.bits[i] = 0
	}
}

func (ba *BitArray) chk(bi, si uint64) bool { return (ba.bits[bi]>>si)&1 > 0 }

// ChkSet returns the value of the bit at position k before setting it.
func (ba *BitArray) ChkSet(k int) (b bool) {
	bi, si := biandsi(k)
	b = ba.chk(bi, si)
	if !b {
		ba.set1(bi, si)
	}
	return

}

// ChkClr returns the value of the bit at position k before clearing it.
func (ba *BitArray) ChkClr(k int) (b bool) {
	bi, si := biandsi(k)
	b = ba.chk(bi, si)
	if b {
		ba.set0(bi, si)
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
		n += bits.OnesCount64(uint64(b))
	}
	return
}

// Chk returns the value of the bit at position k.
func (ba *BitArray) Chk(k int) bool {
	bi, si := biandsi(k)
	return (ba.bits[bi]>>si)&1 > 0
}

// Put sets the value of the bit at position k to v.
func (ba *BitArray) Put(k int, v Bit) {
	bi, si := biandsi(k)
	ba.bits[bi] = (ba.bits[bi] & ^(1 << si)) | (v << si)
}

// Swap swaps the value of bit at position k with v. On return, v contains the old value.
func (ba *BitArray) Swap(k int, v *Bit) {
	bi, si := biandsi(k)
	ob := (ba.bits[bi] >> si) & 1
	ba.bits[bi] = (ba.bits[bi] & ^(1 << si)) | (*v << si)
	*v = ob
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
