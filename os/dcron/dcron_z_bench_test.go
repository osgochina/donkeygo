package dcron_test

import (
	"github.com/osgochina/donkeygo/os/dcron"
	"testing"
)

func Benchmark_Add(b *testing.B) {
	for i := 0; i < b.N; i++ {
		dcron.Add("1 1 1 1 1 1", func() {

		})
	}
}
