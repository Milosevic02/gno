package tmhash

import "testing"

func benchamrkHashCommon(b *testing.B, size int) {
	b.Helper()

	data := []byte("this is a tx")

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		Sum(data)
	}
}

func BenchmarkHash_10(b *testing.B) {
	benchamrkHashCommon(b, 10)
}

func BenchmarkHash_100(b *testing.B) {
	benchamrkHashCommon(b, 100)
}

func BenchmarkHash_1000(b *testing.B) {
	benchamrkHashCommon(b, 1000)
}

func BenchmarkHash_10000(b *testing.B) {
	benchamrkHashCommon(b, 10000)
}
