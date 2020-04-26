# bitarray
Package bitarray provides a convenient and fast BitArray data structure that compactly stores bits in Go.

# Installation
Just do `go get github.com/c2akula/bitarray` to install the package.

# Usage
The package provides one type BitArray. It exposes an api to work with booleans but internally stores them in uint64s.
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

ba.SetAll() // sets all the bits
ba.ClrAll() // clears all the bits
```

## Tests and Benchmarks
Tests and benchmarks can be found in ba_test.go.
```
goos: windows
goarch: amd64
pkg: github.com/c2akula/bitarray
BenchmarkNew
BenchmarkNew-8                  37570326                27.5 ns/op            32B/op          1 allocs/op
BenchmarkBitArray
BenchmarkBitArray/chk
BenchmarkBitArray/chk-8         191312590                5.96 ns/op            0B/op          0 allocs/op
BenchmarkBitArray/put
BenchmarkBitArray/put-8         174875466                6.78 ns/op            0B/op          0 allocs/op
BenchmarkBitArray/set
BenchmarkBitArray/set-8         208177489                5.79 ns/op            0B/op          0 allocs/op
BenchmarkBitArray/tgl
BenchmarkBitArray/tgl-8         205501374                5.82 ns/op            0B/op          0 allocs/op
BenchmarkBitArray_Cnt/cnt_-_bits-8              213713914                5.62 ns/op            0B/op          0 allocs/op
```

# Issues
If you come across any bugs, please use the issue tracker. I will try to get back to you as soon as possible.
