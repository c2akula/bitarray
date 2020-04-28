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

ba.SetAll() // sets all the bits
ba.ClrAll() // clears all the bits
```

## Range Operations
There are two procedures `CopyRange` and `SwapRange` to help work with a range of bits. A`Range` represents
a span over a certain number of bits starting at a specific position.
### CopyRange
```go
b1 := FromStr("11001001")
b2 := FromStr("00101110")
// copy 3 bits starting at position 3 from `b2` into `b1` starting at position 1 
bitarray.CopyRange(b1.Range(1, 3), b2.Range(3, 3))
fmt.Println("b1: ", &b1) // b1 = "10111001"
```
`CopyRange` copies number of bits equal to that of the smaller range.

### SwapRange
```go
b1 := FromStr("11001001")
b2 := FromStr("00101110")
// swap 3 bits starting at position 3 from `b2` with `b1` starting at position 1
bitarray.SwapRange(b1.Range(1, 3), b2.Range(3, 3))
fmt.Println("b1: ", b1, "b2: ", b2) // b1 = "10111001" b2 = "00110010" 
```
`SwapRange` swaps number of bits equal to that of the smaller range.

## Tests and Benchmarks
Tests and benchmarks can be found in ba_test.go.
```
goos: windows
goarch: amd64
pkg: github.com/c2akula/bitarray

BenchmarkNew-8                                            40094756              28.1 ns/op            32 B/op          1 allocs/op
BenchmarkBitArray/chk-8                                    3957813               297 ns/op             0 B/op          0 allocs/op
BenchmarkBitArray/put-8                                  853075514              1.41 ns/op             0 B/op          0 allocs/op
BenchmarkBitArray/set-8                                  846415746              1.43 ns/op             0 B/op          0 allocs/op
BenchmarkBitArray/tgl-8                                  847111244              1.45 ns/op             0 B/op          0 allocs/op
BenchmarkBitArray/swap-range,_best-case-8                 80213366              14.9 ns/op             0 B/op          0 allocs/op
BenchmarkBitArray/swap-range,_worst-case-8                 6011475               184 ns/op             0 B/op          0 allocs/op
BenchmarkBitArray/copy-8                                 251309268              4.76 ns/op             0 B/op          0 allocs/op
BenchmarkBitArray_Cnt/cnt_-_bits-8                       191233108              6.28 ns/op             0 B/op          0 allocs/op
BenchmarkCopyRange/worst_case_-_unaligned_copy-8            705885              1703 ns/op             0 B/op          0 allocs/op
BenchmarkCopyRange/best_case_-_aligned_copy-8            122504306              9.54 ns/op             0 B/op          0 allocs/op

```

# Issues
If you come across any bugs, please use the issue tracker. I will try to get back to you as soon as possible.
