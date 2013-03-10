//stollen from bitbucket.org/StephaneBunel/xxhash-go
package xxhash

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"hash/adler32"
	"hash/crc32"
	"hash/fnv"
	"testing"
)

var (
	blob1       = []byte("Lorem ipsum dolor sit amet, consectetuer adipiscing elit, ")
	blob2       = []byte("sed diam nonummy nibh euismod tincidunt ut laoreet dolore magna aliquam erat volutpat.")
	VeryBigFile = "a-very-big-file"
)

// MurmurHash 3
// mmh3.Hash32 stollen from https://github.com/reusee/mmh3
func mmh3Hash32(key []byte) uint32 {
	length := len(key)
	if length == 0 {
		return 0
	}
	var c1, c2 uint32 = 0xcc9e2d51, 0x1b873593
	nblocks := length / 4
	var h, k uint32
	buf := bytes.NewBuffer(key)
	for i := 0; i < nblocks; i++ {
		binary.Read(buf, binary.LittleEndian, &k)
		k *= c1
		k = (k << 15) | (k >> (32 - 15))
		k *= c2
		h ^= k
		h = (h << 13) | (h >> (32 - 13))
		h = (h * 5) + 0xe6546b64
	}
	k = 0
	tailIndex := nblocks * 4
	switch length & 3 {
	case 3:
		k ^= uint32(key[tailIndex+2]) << 16
		fallthrough
	case 2:
		k ^= uint32(key[tailIndex+1]) << 8
		fallthrough
	case 1:
		k ^= uint32(key[tailIndex])
		k *= c1
		k = (k << 13) | (k >> (32 - 15))
		k *= c2
		h ^= k
	}
	h ^= uint32(length)
	h ^= h >> 16
	h *= 0x85ebca6b
	h ^= h >> 13
	h *= 0xc2b2ae35
	h ^= h >> 16
	return h
}

func Test_Checksum32(t *testing.T) {
	h32 := Checksum32(blob1)
	fmt.Printf("Checksum32(\"%v\") = 0x%08x\n", string(blob1), h32)
	if h32 != 0x1130e7d4 {
		t.Fail()
	}
}

func Test_Checksum32Seed(t *testing.T) {
	h32 := Checksum32Seed(blob1, 1471)
	fmt.Printf("Checksum32Seed(\"%v\", 1471) = 0x%08x\n", string(blob1), h32)
	if h32 != 0xba59a258 {
		t.Fail()
	}
}

func Test_New32(t *testing.T) {
	var digest = New(0)
	digest.Write(blob1)
	digest.Write(blob2)
	h32 := digest.Sum32()
	fmt.Printf("Sum32 = 0x%08x\n", h32)
	if h32 != 0x0d44373a {
		t.Fail()
	}
}

func Test_New32Seed(t *testing.T) {
	var digest = New(1471)
	digest.Write(blob1)
	digest.Write(blob2)
	h32 := digest.Sum32()
	fmt.Printf("Sum32 = 0x%08x\n", h32)
	if h32 != 0x3265e220 {
		t.Fail()
	}
}

func Test_Reset(t *testing.T) {
	var digest = New(0)
	digest.Write(blob2)
	digest.Reset()
	digest.Write(blob1)
	h32 := digest.Sum32()
	fmt.Printf("Sum32 = 0x%08x\n", h32)
	if h32 != 0x1130e7d4 {
		t.Fail()
	}
}

func Benchmark_xxhash32(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Checksum32(blob1)
	}
}

func Benchmark_CRC32IEEE(b *testing.B) {
	for i := 0; i < b.N; i++ {
		crc32.ChecksumIEEE(blob1)
	}
}

func Benchmark_Adler32(b *testing.B) {
	for i := 0; i < b.N; i++ {
		adler32.Checksum(blob1)
	}
}

func Benchmark_Fnv32(b *testing.B) {
	h := fnv.New32()
	for i := 0; i < b.N; i++ {
		h.Sum(blob1)
	}
}

func Benchmark_MurmurHash3Hash32(b *testing.B) {
	for i := 0; i < b.N; i++ {
		mmh3Hash32(blob1)
	}
}
