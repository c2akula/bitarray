package bitarray

import (
	"fmt"
	"math/rand"
	"strings"
	"testing"
	"time"
)

func randomize(ba *BitArray, rng *rand.Rand) {
	for i := 0; i < ba.n; i++ {
		if rng.Float32() < 0.5 {
			ba.Tgl(i)
		}
	}
}

func TestNew(t *testing.T) {
	tests := []struct {
		n, nbits int
	}{
		{4, 1},
		{64, 1},
		{65, 2},
		{256, 4},
		{257, 5},
	}

	for i, tt := range tests {
		ba := New(tt.n)
		if tt.nbits != len(ba.bits) {
			t.Fatalf("test %d failed. got = %d, exp = %d\n", i, len(ba.bits), tt.nbits)
		}
	}
}

func TestBitArray(t *testing.T) {
	t.Run("swap", func(t *testing.T) {
		ba := New(257)
		v := One
		for k := 0; k < ba.n; k++ {
			ba.Swap(k, &v)
			if v != Zero && !ba.Chk(k) {
				t.Fatalf("Test at bit %d failed. got = %d, exp = %d\n", k, One, Zero)
			}
		}
	})

	t.Run("chkset", func(t *testing.T) {
		ba := New(257)
		ba.ClrAll()
		for k := 0; k < ba.n; k++ {
			if ba.ChkSet(k) != false {
				t.Fatalf("Test at bit %d failed. got = %t, exp = %t\n", k, true, false)
			}
		}
	})

	t.Run("chkclr", func(t *testing.T) {
		ba := New(257)
		ba.SetAll()

		for k := 0; k < ba.n; k++ {
			if ba.ChkClr(k) != true {
				t.Fatalf("Test at bit %d failed. got = %t, exp = %t\n", k, false, true)
			}
		}
	})

	t.Run("chk", func(t *testing.T) {
		eb := "10000100010111111000011100001100111000000111011011110110001101000111100000010001001000110001001101110001010001100010010010101010001100010101101101011101100001111011110100100011100001101001001001110101000010000000011010000010111111111111111100000001110100001"
		ba := New(len(eb))
		for k, b := range eb {
			if b == '1' {
				ba.Put(k, One)
			}
		}

		// test
		for k, b := range eb {
			if (b == '1') != ba.Chk(k) {
				t.Fatalf("Test failed. got = %s\nexp = %s\n", ba.String(), eb)
			}
		}

	})

	t.Run("put", func(t *testing.T) {
		eb := "10000100010111111000011100001100111000000111011011110110001101000111100000010001001000110001001101110001010001100010010010101010001100010101101101011101100001111011110100100011100001101001001001110101000010000000011010000010111111111111111100000001110100001"
		ba := New(len(eb))
		for k, b := range eb {
			if b == '1' {
				ba.Put(k, One)
			}
		}
		if eb != ba.String() {
			t.Fatalf("Test failed. got = %s\nexp = %s\n", ba.String(), eb)
		}
	})

	t.Run("set and get", func(t *testing.T) {
		eb := "11010011000010000010001110000010001000111011111010001000001000110011111110000001010111111010000001101110010011000010011100100001110100111110010100101100101000101001110100001101100010001010001001100011100001001111011011001100101011001111100010100000001001111"
		ba := New(len(eb))
		for k, b := range eb {
			if b == '1' {
				ba.Set(k)
			}
		}
		if eb != ba.String() {
			t.Fatalf("Test failed. got = %s\nexp = %s\n", ba.String(), eb)
		}
	})

	t.Run("count", func(t *testing.T) {
		eb := "00011100111011110101011110110010010100101010000110001010100011011010111001110110011111001101100000101011010011001110110000111010101100010110001101101011111100110000110101100000000100011110100111000101100001110010010011001001000001111010101110000001001101001"
		ba := New(len(eb))
		for k, b := range eb {
			if b == '1' {
				ba.Set(k)
			}
		}
		slow_count := func(ba *BitArray) (n int) {
			for i := 0; i < ba.n; i++ {
				if ba.Chk(i) {
					n++
				}
			}
			return
		}

		ec := slow_count(&ba)
		gc := ba.Cnt()
		if gc != ec {
			t.Fatalf("Test failed. got = %d, exp = %d\n", gc, ec)
		}

		fmt.Println("gc: ", gc)
	})

	t.Run("copy", func(t *testing.T) {
		const sb = "00011100111011110101011110110010010100101010000110001010100011011010111001110110011111001101100000101011010011001110110000111010101100010110001101101011111100110000110101100000000100011110100111000101100001110010010011001001000001111010101110000001001101001"
		src := New(len(sb))
		for i, b := range sb {
			if b == '1' {
				src.Set(i)
			}
		}

		dst := New(len(sb))
		Copy(&dst, &src)
		if dst.String() != src.String() {
			t.Fatalf("Test failed. got = %s\nexp = %s\n", dst.String(), sb)
		}
	})

	t.Run("fromstring", func(t *testing.T) {
		const sb = "00011100111011110101011110110010010100101010000110001010100011011010111001110110011111001101100000101011010011001110110000111010101100010110001101101011111100110000110101100000000100011110100111000101100001110010010011001001000001111010101110000001001101001"
		ba := FromStr(sb)
		if ba.String() != sb {
			t.Fatalf("Test failed. got = %s\nexp = %s\n", ba.String(), sb)
		}
	})

	t.Run("fromuint64", func(t *testing.T) {
		const u = uint64(0xfefefefefefefefe)
		ba := FromUint64(u)
		var sb strings.Builder
		for i := 0; i < ba.n; i++ {
			sb.WriteByte(chkbit(i, u))
		}
		fmt.Println("sb", sb.String(), sb.Len())
		if ba.String() != sb.String() {
			t.Fatalf("Test failed. got = %s\nexp = %s\n", ba.String(), sb.String())
		}
	})

	// TODO: test with
	t.Run("get-slice", func(t *testing.T) {
		slice := func(ba *BitArray, b, n int) (oa BitArray) {
			if (b+n < 0) || (b+n > ba.n) {
				panic("range specified is outside bounds")
			}

			if b != n {
				oa = New(n)
				for i := 0; i < n; i++ {
					if ba.Chk(b) {
						oa.Set(i)
					}
					b++
				}
			}
			return
		}

		ba := FromStr("010111101")
		oa := slice(&ba, 1, 5)
		const eb = "10111"
		if gb := oa.String(); gb != eb {
			t.Fatalf("Test failed. got = %s, exp: %s\n", gb, eb)
		}
	})

	t.Run("swap-range, same-length > 64bits", func(t *testing.T) {
		const sa = "11010011000010000010001110000010001000111011111010001000001000110011111110000001010111111010000001101110010011000010011100100001110100111110010100101100101000101001110100001101100010001010001001100011100001001111011011001100101011001111100010100000001001111"
		const sb = "00011100111011110101011110110010010100101010000110001010100011011010111001110110011111001101100000101011010011001110110000111010101100010110001101101011111100110000110101100000000100011110100111000101100001110010010011001001000001111010101110000001001101001"

		// length same
		a := FromStr(sa)
		b := FromStr(sb)
		SwapRange(&a, &b, 3)
		const ea = "11011100111011110101011110110010010100101010000110001010100011011010111001110110011111001101100000101011010011001110110000111010101100010110001101101011111100110000110101100000000100011110100111000101100001110010010011001001000001111010101110000001001101001"
		const eb = "00010011000010000010001110000010001000111011111010001000001000110011111110000001010111111010000001101110010011000010011100100001110100111110010100101100101000101001110100001101100010001010001001100011100001001111011011001100101011001111100010100000001001111"

		if a.String() != ea {
			t.Fatalf("Test `a` failed. got = %s, exp = %s\n", a.String(), ea)
		}

		if b.String() != eb {
			t.Fatalf("Test `b` failed. got = %s, exp = %s\n", b.String(), eb)
		}
	})

	t.Run("swap-range, same-length < 64bits", func(t *testing.T) {

		// length same
		a := FromStr("1110011101")
		b := FromStr("1110110010")
		SwapRange(&a, &b, 4)
		const ea = "1110110010"
		const eb = "1110011101"

		if a.String() != ea {
			t.Fatalf("Test `a` failed. got = %s, exp = %s\n", a.String(), ea)
		}

		if b.String() != eb {
			t.Fatalf("Test `b` failed. got = %s, exp = %s\n", b.String(), eb)
		}
	})

	t.Run("swap-range, min-length < 64bits", func(t *testing.T) {
		// minimum length swap
		a := FromStr("1110011")
		b := FromStr("1110110010")
		// swap_slice(&a, &b, 3)
		SwapRange(&a, &b, 3)
		const ea = "1110110"
		const eb = "1110011010"

		if a.String() != ea {
			t.Fatalf("Test `a` failed. got = %s, exp = %s\n", a.String(), ea)
		}

		if b.String() != eb {
			t.Fatalf("Test `b` failed. got = %s, exp = %s\n", b.String(), eb)
		}

	})

	t.Run("swap-range, min-length > 64bits", func(t *testing.T) {

		sa := "00011100111011110101011110110010010100101010000110001010100011011010111001110110011111001101100000101011010011001110110000111010101100010110001101101011111100110000110101100000000100011110100111000101100001110010010011001001000001111010101110000001001101001"
		sb := strings.Repeat("0", 127)

		a := FromStr(sa)
		b := FromStr(sb)
		//      7       +   120  + 130
		ea := "0001110" + sb[7:] + sa[127:]
		eb := sb[:7] + sa[7:127]

		SwapRange(&a, &b, 7)

		if a.String() != ea {
			t.Fatalf("Test `a` failed. got = %s, exp = %s\n", a.String(), ea)
		}

		if b.String() != eb {
			t.Fatalf("Test `b` failed. got = %s, exp = %s\n", b.String(), eb)
		}

	})
}

func BenchmarkNew(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		New(256)
	}
}

func BenchmarkBitArray(b *testing.B) {
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))

	b.Run("chk", func(b *testing.B) {
		b.ReportAllocs()
		b.StopTimer()
		ba := New(257)
		b.StartTimer()
		for i := 0; i < b.N; i++ {
			for k := 0; k < ba.n; k++ {
				ba.Chk(k)
			}
		}
	})

	b.Run("put", func(b *testing.B) {
		b.ReportAllocs()
		b.StopTimer()
		ba := New(257)
		b.StartTimer()
		for i := 0; i < b.N; i++ {
			ba.Put(127, One)
		}
	})

	b.Run("set", func(b *testing.B) {
		b.ReportAllocs()
		b.StopTimer()
		ba := New(257)
		b.StartTimer()
		for i := 0; i < b.N; i++ {
			ba.Set(127)
		}
	})

	b.Run("tgl", func(b *testing.B) {
		b.ReportAllocs()
		b.StopTimer()
		ba := New(257)
		b.StartTimer()
		for i := 0; i < b.N; i++ {
			ba.Tgl(127)
		}
	})

	b.Run("swap-range, best-case", func(b *testing.B) {
		b.ReportAllocs()
		b.StopTimer()
		b1 := New(257)
		b2 := New(b1.n)
		randomize(&b1, rng)
		randomize(&b2, rng)
		b.StartTimer()

		for i := 0; i < b.N; i++ {
			SwapRange(&b1, &b2, 0)
		}
	})

	// the worst-case is when we are at position 1 in a multiple 64-bit blocks situation
	// we have to swap bit-by-bit till we reach the next block start
	b.Run("swap-range, worst-case", func(b *testing.B) {
		b.ReportAllocs()
		b.StopTimer()
		b1 := New(257)
		b2 := New(b1.n)
		randomize(&b1, rng)
		randomize(&b2, rng)
		b.StartTimer()

		for i := 0; i < b.N; i++ {
			SwapRange(&b1, &b2, 1)
		}
	})

	b.Run("copy", func(b *testing.B) {
		b.ReportAllocs()
		b.StopTimer()
		b1 := New(257)
		b2 := New(b1.n)
		randomize(&b1, rng)
		randomize(&b2, rng)
		b.StartTimer()

		for i := 0; i < b.N; i++ {
			Copy(&b2, &b1)
		}
	})

}

func BenchmarkBitArray_Cnt(b *testing.B) {
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	b.Run("cnt - chk", func(b *testing.B) {
		b.ReportAllocs()
		b.StopTimer()
		ba := New(257)
		randomize(&ba, rng)
		cnt := func(ba *BitArray) (n int) {
			for k := 0; k < ba.n; k++ {
				if ba.Chk(k) {
					n++
				}
			}
			return
		}
		b.StartTimer()
		for i := 0; i < b.N; i++ {
			cnt(&ba)
		}
	})

	b.Run("cnt - bits", func(b *testing.B) {
		b.ReportAllocs()
		b.StopTimer()
		ba := New(257)
		randomize(&ba, rng)
		b.StartTimer()
		for i := 0; i < b.N; i++ {
			ba.Cnt()
		}
	})

}

func BenchmarkBiandSi(b *testing.B) {
	for i := 0; i < b.N; i++ {
		biandsi(i)
	}
}

func chkbit(k int, u uint64) byte { return '0' + byte((u>>k)&1) }
