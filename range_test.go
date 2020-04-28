package bitarray

import (
	"math/rand"
	"strings"
	"testing"
	"time"
)

func TestSwapRange(t *testing.T) {
	t.Run("swap-range, same-length < 64 bits", func(t *testing.T) {
		const (
			sa = "100110011"
			sb = "000001110"
			ea = "000000011"
			eb = "100111110"
		)

		a := FromStr(sa)
		b := FromStr(sb)

		SwapRange(a.Range(0, 5), b.Range(0, 5))

		if a.String() != ea {
			t.Fatalf("Test `a` failed. got = %s, exp = %s\n", a.String(), ea)
		}

		if b.String() != eb {
			t.Fatalf("Test `b` failed. got = %s, exp = %s\n", b.String(), ea)
		}
	})

	t.Run("swap-range, min-length < 64 bits", func(t *testing.T) {
		const (
			sa = "100110011"
			sb = "000001110"
			ea = "100011011"
			eb = "000011010"
		)

		a := FromStr(sa)
		b := FromStr(sb)

		SwapRange(a.Range(2, 4), b.Range(3, 5))

		if a.String() != ea {
			t.Fatalf("Test `a` failed. got = %s, exp = %s\n", a.String(), ea)
		}

		if b.String() != eb {
			t.Fatalf("Test `b` failed. got = %s, exp = %s\n", b.String(), ea)
		}
	})

	t.Run("swap-range, same-length > 64 bits", func(t *testing.T) {
		const (
			sa = "011011101111011001101011110101101111001100010001000000011110111000101000"
			sb = "110000001000110111010111010111111100011111111001001011101011000111110000"
			ea = "110000001000110111010111010111111100011111111001001011101011000111110000"
			eb = "011011101111011001101011110101101111001100010001000000011110111000101000"
			k  = 0
			nb = 69
		)

		a := FromStr(sa)
		b := FromStr(sb)

		SwapRange(a.Range(k, nb), b.Range(k, nb))

		if a.String() != ea {
			t.Fatalf("Test `a` failed. got = %s, exp = %s\n", a.String(), ea)
		}

		if b.String() != eb {
			t.Fatalf("Test `b` failed. got = %s, exp = %s\n", b.String(), ea)
		}
	})

	t.Run("swap-range, aligned+unaligned, same-length > 64 bits", func(t *testing.T) {
		const (
			sa = "011011101111011001101011110101101111001100010001000000011110111000101000"
			sb = "110000001000110111010111010111111100011111111001001011101011000111110000"
			ea = "011110000001000110111010111010111111100011111111001001011101011000111110"
			eb = "011101111011001101011110101101111001100010001000000011110111000101000000"
			nb = 69
		)

		a := FromStr(sa)
		b := FromStr(sb)

		SwapRange(a.Range(3, nb), b.Range(0, nb))

		if a.String() != ea {
			t.Fatalf("Test `a` failed. got = %s, exp = %s\n", a.String(), ea)
		}

		if b.String() != eb {
			t.Fatalf("Test `b` failed. got = %s, exp = %s\n", b.String(), ea)
		}
	})

	t.Run("swap-range, aligned+unaligned, min-length > 64 bits", func(t *testing.T) {
		const (
			sa = "011011101111011001101011110101101111001100010001000000011110111000101000"
			sb = "110000001000110111010111010111111100011111111001001011101011000111110000"
			ea = "011011000000100011011101011101011111110001111111100100101110101100011111"
			eb = "111011110110011010111101011011110011000100010000000111101110001010000000"
		)

		a := FromStr(sa)
		b := FromStr(sb)

		SwapRange(a.Range(6, 66), b.Range(2, b.n-2))

		if a.String() != ea {
			t.Fatalf("Test `a` failed. got = %s, exp = %s\n", a.String(), ea)
		}

		if b.String() != eb {
			t.Fatalf("Test `b` failed. got = %s, exp = %s\n", b.String(), ea)
		}
	})

}

func BenchmarkSwapRange(b *testing.B) {
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	b.Run("swap-range, worst-case", func(b *testing.B) {
		b.ReportAllocs()
		b.StopTimer()
		b1 := New(512)
		randomize(&b1, rng)
		b2 := New(512)
		randomize(&b2, rng)

		b1r, b2r := b1.Range(47, 445), b2.Range(67, b2.n-67)
		b.StartTimer()

		for i := 0; i < b.N; i++ {
			SwapRange(b1r, b2r)
		}
	})

	b.Run("swap-range, start position same, mixed > 64 bits", func(b *testing.B) {
		b.ReportAllocs()
		b.StopTimer()
		b1 := New(512)
		b2 := New(b1.n)
		randomize(&b1, rng)
		randomize(&b2, rng)
		b1r := b1.Range(47, 445)
		b2r := b2.Range(47, 437)
		b.StartTimer()

		for i := 0; i < b.N; i++ {
			SwapRange(b1r, b2r)
		}
	})

	b.Run("swap-range, start position same == 0, mixed > 64 bits", func(b *testing.B) {
		b.ReportAllocs()
		b.StopTimer()
		b1 := New(512)
		b2 := New(b1.n)
		randomize(&b1, rng)
		randomize(&b2, rng)
		b1r := b1.Range(0, 445)
		b2r := b2.Range(0, 437)
		b.StartTimer()

		for i := 0; i < b.N; i++ {
			SwapRange(b1r, b2r)
		}
	})

	b.Run("swap-range, start position same == 0, aligned > 64 bits", func(b *testing.B) {
		b.ReportAllocs()
		b.StopTimer()
		b1 := New(512)
		b2 := New(b1.n)
		randomize(&b1, rng)
		randomize(&b2, rng)
		b1r := b1.Range(0, 448)
		b2r := b2.Range(0, 448)
		b.StartTimer()

		for i := 0; i < b.N; i++ {
			SwapRange(b1r, b2r)
		}
	})

}

func TestCopyRange(t *testing.T) {
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	t.Run("non-aligned copy < 64 bits", func(t *testing.T) {
		const ss = "011100000"
		const ds = "101000110"
		src := FromStr(ss)
		dst := FromStr(ds)
		CopyRange(dst.Range(0, 5), src.Range(0, 5))
		ed := ss[:5] + ds[5:]
		if dst.String() != ed {
			t.Fatalf("Test CopyRange failed. got = %s, exp = %s\n", dst.String(), ed)
		}
	})

	t.Run("aligned copy, nbits % 64 == 0", func(t *testing.T) {
		src := New(256)
		randomize(&src, rng)
		dst := New(512)
		es := src.String() + strings.Repeat("0", src.n)
		CopyRange(dst.Range(0, dst.n), src.Range(0, src.n))
		if dst.String() != es {
			t.Fatalf("Test CopyRange failed. got = %s\nexp=%s\n", dst.String(), es)
		}
	})

	t.Run("aligned+unaligned copy, nbits > 64", func(t *testing.T) {
		src := New(256)
		randomize(&src, rng)
		dst := New(512)
		es := src.String()[:49] + strings.Repeat("0", dst.n-49)
		CopyRange(dst.Range(0, dst.n), src.Range(0, 49))
		if dst.String() != es {
			t.Fatalf("Test CopyRange failed. got = %s\nexp=%s\n", dst.String(), es)
		}
	})

	t.Run("unaligned copy, nbits > 64", func(t *testing.T) {
		const db, dn = 49, 128 // no. of bits from dst
		const sb, sn = 23, 230 // no. of bits from src

		src := New(256)
		randomize(&src, rng)
		dst := New(256)
		randomize(&dst, rng)
		//                  49  +                   128  +                  79
		es := dst.String()[:db] + src.String()[sb:sb+dn] + dst.String()[db+dn:]
		CopyRange(dst.Range(db, dn), src.Range(sb, sn))
		if dst.String() != es {
			t.Fatalf("Test CopyRange failed. got = %s\nexp=%s\n", dst.String(), es)
		}
	})
}

func BenchmarkCopyRange(b *testing.B) {
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))

	b.Run("worst case - unaligned copy", func(b *testing.B) {
		b.ReportAllocs()
		b.StopTimer()
		src := New(512)
		randomize(&src, rng)
		dst := New(512)
		randomize(&dst, rng)
		dr, sr := dst.Range(127, 385), src.Range(113, 385)
		b.StartTimer()

		for i := 0; i < b.N; i++ {
			CopyRange(dr, sr)
		}
	})

	b.Run("best case - aligned copy", func(b *testing.B) {
		b.ReportAllocs()
		b.StopTimer()
		src := New(512)
		randomize(&src, rng)
		dst := New(512)
		randomize(&dst, rng)
		dr, sr := dst.Range(128, 384), src.Range(64, 384)
		b.StartTimer()

		for i := 0; i < b.N; i++ {
			CopyRange(dr, sr)
		}
	})

}
