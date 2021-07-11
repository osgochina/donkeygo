package dhash_test

import (
	"github.com/osgochina/donkeygo/encoding/dhash"
	"testing"
)

var (
	str = []byte("This is the test string for hash.")
)

func BenchmarkBKDRHash(b *testing.B) {
	for i := 0; i < b.N; i++ {
		dhash.BKDRHash(str)
	}
}

func BenchmarkBKDRHash64(b *testing.B) {
	for i := 0; i < b.N; i++ {
		dhash.BKDRHash64(str)
	}
}

func BenchmarkSDBMHash(b *testing.B) {
	for i := 0; i < b.N; i++ {
		dhash.SDBMHash(str)
	}
}

func BenchmarkSDBMHash64(b *testing.B) {
	for i := 0; i < b.N; i++ {
		dhash.SDBMHash64(str)
	}
}

func BenchmarkRSHash(b *testing.B) {
	for i := 0; i < b.N; i++ {
		dhash.RSHash(str)
	}
}

func BenchmarkSRSHash64(b *testing.B) {
	for i := 0; i < b.N; i++ {
		dhash.RSHash64(str)
	}
}

func BenchmarkJSHash(b *testing.B) {
	for i := 0; i < b.N; i++ {
		dhash.JSHash(str)
	}
}

func BenchmarkJSHash64(b *testing.B) {
	for i := 0; i < b.N; i++ {
		dhash.JSHash64(str)
	}
}

func BenchmarkPJWHash(b *testing.B) {
	for i := 0; i < b.N; i++ {
		dhash.PJWHash(str)
	}
}

func BenchmarkPJWHash64(b *testing.B) {
	for i := 0; i < b.N; i++ {
		dhash.PJWHash64(str)
	}
}

func BenchmarkELFHash(b *testing.B) {
	for i := 0; i < b.N; i++ {
		dhash.ELFHash(str)
	}
}

func BenchmarkELFHash64(b *testing.B) {
	for i := 0; i < b.N; i++ {
		dhash.ELFHash64(str)
	}
}

func BenchmarkDJBHash(b *testing.B) {
	for i := 0; i < b.N; i++ {
		dhash.DJBHash(str)
	}
}

func BenchmarkDJBHash64(b *testing.B) {
	for i := 0; i < b.N; i++ {
		dhash.DJBHash64(str)
	}
}

func BenchmarkAPHash(b *testing.B) {
	for i := 0; i < b.N; i++ {
		dhash.APHash(str)
	}
}

func BenchmarkAPHash64(b *testing.B) {
	for i := 0; i < b.N; i++ {
		dhash.APHash64(str)
	}
}
