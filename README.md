# bitarray
Package bitarray provides a convenient and fast `BitArray` data structure that compactly stores bits in Go.

# Installation
Just do `go get github.com/c2akula/bitarray` to install the package.

# Usage
The package provides a single type `BitArray`. It exposes an api to work with booleans but internally stores them as bits in uint64s.
## Creating a BitArray
```go
ba := bitarray.New(65) // creates a bitarray containing 65 bits
```
## Basic Operations
```go
ba.Set(5) // sets the bit at position 5
b := ba.Chk(5) // gets the bit at position 5
ba.Tgl(5) // toggles the bit at position 5
ba.Clr(5) // clears the bit at position 5
ba.Cnt() // returns the number of set bits
```
## Other Operations
```go
oldValue := ba.ChkSet(5) // returns the value at position 5 before setting it
oldValue := ba.ChkClr(5) // returns the value at position 5 before clearing it

v := bit.One
ba.Swap(5, &v) // swaps the value at position 5 with v

b1 := bitarray.FromStr("110110010")
b2 := bitarray.FromStr("0011011")
bitarray.SwapRange(&b1, &b2, 3) // swap bits starting at position 3
fmt.Println("b1: ", &b1, "b2: ", &b2) // b1 = "110101110", b2 = "0011100" 


ba.SetAll() // sets all the bits
ba.ClrAll() // clears all the bits
```

## Tests and Benchmarks
Tests and benchmarks can be found in ba_test.go.
```
goos: windows
goarch: amd64
pkg: github.com/c2akula/bitarray
BenchmarkBitArray/chk-8                                   20200155                311 ns/op            0 B/op          0 allocs/op
BenchmarkBitArray/put-8                                 1000000000               1.41 ns/op            0 B/op          0 allocs/op
BenchmarkBitArray/set-8                                 1000000000               1.31 ns/op            0 B/op          0 allocs/op
BenchmarkBitArray/tgl-8                                 1000000000               1.42 ns/op            0 B/op          0 allocs/op
BenchmarkBitArray/swap-range,_best-case-8                632504266               9.50 ns/op            0 B/op          0 allocs/op
BenchmarkBitArray/swap-range,_worst-case-8                33005805                185 ns/op            0 B/op          0 allocs/op
BenchmarkBitArray/copy-8                                1000000000               5.00 ns/op            0 B/op          0 allocs/op
BenchmarkBitArray_Cnt/cnt_-_chk-8                         13548182                442 ns/op            0 B/op          0 allocs/op
BenchmarkBitArray_Cnt/cnt_-_bits-8                       988570477               6.06 ns/op            0 B/op          0 allocs/op
```

# Issues
If you come across any bugs, please use the issue tracker. I will try to get back to you as soon as possible.
