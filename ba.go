package bitarray

import (
	"bytes"
	"math"
	"math/bits"

	"github.com/bmkessler/fastdiv"
)

var c64 fastdiv.Uint64

func init() {
	c64 = fastdiv.NewUint64(64)
}

type Bit = uint64

const (
	Zero = Bit(iota)
	One
)

// BitArray is an array data structure that compactly stores bits.
// Bits externally represented as `bool` are stored internally as `uint64`s.
// The total number of bits stored is set at creation and is immutable.
type BitArray struct {
	bits []uint64
	n    int // no. of bits
}

// New creates a new BitArray of `n` bits.
func New(n int) BitArray {
	return BitArray{
		bits: make([]uint64, int(math.Ceil(float64(n)/64))),
		n:    n,
	}
}

// Size returns the no. of bits stored.
func (ba *BitArray) Size() int { return ba.n }

func biandsi(k int) (uint64, uint64) { return c64.DivMod(uint64(k)) }

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
		n += bits.OnesCount64(b)
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
	bs := bytes.Repeat([]byte("0"), ba.n)
	for i := range bs {
		if ba.Chk(i) {
			bs[i] = '1'
		}
	}
	return string(bs)
}
