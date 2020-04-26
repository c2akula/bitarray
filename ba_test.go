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
}

func BenchmarkNew(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		New(256)
	}
}

func BenchmarkBitArray(b *testing.B) {
	b.Run("chk", func(b *testing.B) {
		b.ReportAllocs()
		b.StopTimer()
		ba := New(257)
		b.StartTimer()
		for i := 0; i < b.N; i++ {
			ba.Chk(127)
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

}

func BenchmarkBitArray_Cnt(b *testing.B) {
	b.Run("cnt - chk", func(b *testing.B) {
		b.ReportAllocs()
		b.StopTimer()
		ba := New(257)
		randomize(&ba, rand.New(rand.NewSource(time.Now().UnixNano())))
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
		randomize(&ba, rand.New(rand.NewSource(time.Now().UnixNano())))
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
