package bitarray

// A Range represents a span over a certain number of bits in a BitArray starting
// at specific position.
type Range struct {
	*BitArray
	b, n int
}

// Range creates a Range object representing `n` bits starting at `b`.
func (ba *BitArray) Range(b, n int) Range {
	if (n-b < 0) || (n-b > ba.n) {
		panic("index out of bounds")
	}
	return Range{ba, b, n}
}

// CopyRange copies bits from `src` into `dst` specified by the ranges.
// The procedure copies number of bits equal to the the minimum of the two ranges.
// It is undefined behavior to copy overlapping ranges.
func CopyRange(dst, src Range) {
	nb := min(dst.n, src.n) // no. of bits to get

	dbi, dsi := biandsi(dst.b)
	sbi, ssi := biandsi(src.b)

	if dsi == 0 && ssi == 0 {
		m, n := biandsi(nb) // nb/64, nb%64
		copy(dst.bits[dbi:dbi+m], src.bits[sbi:sbi+m])
		if n != 0 {
			// copy remaining < 64 bits, one bit at a time
			copyu64bits(n, &dst.bits[dbi+m], src.bits[sbi+m])
		}
		return
	}

	if dst.b == src.b {
		bi, si := biandsi(dst.b)

		// wind the first block till the next boundary
		for ; si < 64; si++ {
			if nb != 0 {
				setbit(&dst.bits[bi], si, getbit(src.bits[bi], si))
			} else {
				return
			}
		}
		// we've wound up to the next boundary
		// so we'll copy blocks up until the last-but-one block
		bi++
		m, n := biandsi(nb) // nb/64, nb%64
		copy(dst.bits[bi:bi+m], src.bits[bi:bi+m])

		if n != 0 {
			// copy the remaining bits from the last block
			copyu64bits(n, &dst.bits[bi+m], src.bits[bi+m])
		}
		return
	}

	unalignedCopy(nb, dst.bits, dbi, dsi, src.bits, sbi, ssi)
}

func copyu64bits(n uint64, d *uint64, s uint64) {
	for i := uint64(0); n != 0; i++ {
		setbit(d, i, getbit(s, i))
		n--
	}
}

// SwapRange swaps the bits of two ranges `a` and `b`.
// The procedure copies number of bits equal to the the minimum of the two ranges.
// It is undefined behavior to swap overlapping ranges.
func SwapRange(a, b Range) {
	nb := min(a.n, b.n) // no. of bits to swap

	abi, asi := biandsi(a.b)
	bbi, bsi := biandsi(b.b)

	if asi == 0 && bsi == 0 {
		m, n := biandsi(nb) // nb/64, nb%64
		// 	swap the aligned blocks first
		for m != 0 {
			a.bits[abi], b.bits[bbi] = b.bits[bbi], a.bits[abi]
			abi++
			bbi++
			m--
		}

		if n != 0 {
			// 	swap the bits of the last block, bit-by-bit
			swapu64bits(n, &a.bits[abi], &b.bits[bbi])
		}
		return
	}

	// starting point is same
	if a.b == b.b {
		bi, si := biandsi(a.b)
		// process the first block
		anum, bnum := &a.bits[bi], &b.bits[bi]
		for ; si < 64; si++ {
			if nb != 0 {
				abit := getbit(*anum, si)
				bbit := getbit(*bnum, si)
				setbit(anum, si, bbit)
				setbit(bnum, si, abit)
				nb--
			} else {
				return
			}
		}

		// 	we've wound up to the beginning of the next boundary
		// process from the second till the last-but-one block
		bi++
		m, n := biandsi(nb) // nb/64, nb%64
		// swap until the last block
		for m != 0 {
			a.bits[bi], b.bits[bi] = b.bits[bi], a.bits[bi]
			bi++
			m--
		}

		if n != 0 {
			// swap the remaining bits, if any, in the last block bitwise
			swapu64bits(n, &a.bits[bi], &b.bits[bi])
		}
		return
	}

	// 	general case
	unalignedSwap(nb, a.bits, abi, asi, b.bits, bbi, bsi)
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func getbit(u uint64, i uint64) Bit     { return (u >> i) & 1 }
func setbit(u *uint64, i uint64, b Bit) { *u = (*u & ^(1 << i)) | (b << i) }

func alignedcopy(nb int, dst []Bit, di uint64, src []Bit, si uint64) (int, uint64, uint64) {
	m := uint64(nb / 64) // no. of blocks to copy
	copy(dst[di:di+m], src[si:si+m])
	// for nb != 0 {
	// 	dst[di] = src[si]
	// 	di++
	// 	si++
	// 	nb -= 64
	// }
	return nb, di + m, si + m
}

func unalignedCopy(nb int, d []Bit, dbi, dsi uint64, s []Bit, sbi, ssi uint64) {
	for nb != 0 {
		// dst[dbi][dsi] <- src[sbi][ssi]
		setbit(&d[dbi], dsi, getbit(s[sbi], ssi))

		// move to the next src bit
		ssi++
		if ssi >= 64 {
			sbi++
			ssi = 0
		}

		// move to the next dst bit
		dsi++
		if dsi >= 64 {
			dbi++
			dsi = 0
		}

		// one less bit to process
		nb--
	}
}

func swapu64bits(n uint64, a, b *uint64) {
	for i := uint64(0); n != 0; i++ {
		abit := getbit(*a, i)
		bbit := getbit(*b, i)
		setbit(a, i, bbit)
		setbit(b, i, abit)
		n--
	}
}

func unalignedSwap(nb int, a []Bit, abi, asi uint64, b []Bit, bbi, bsi uint64) {
	for nb != 0 {
		anum, bnum := &a[abi], &b[bbi]
		abit := getbit(*anum, asi)
		bbit := getbit(*bnum, bsi)
		setbit(anum, asi, bbit)
		setbit(bnum, bsi, abit)

		asi++
		if asi >= 64 {
			abi++
			asi = 0
		}

		bsi++
		if bsi >= 64 {
			bbi++
			bsi = 0
		}

		nb--
	}
}
