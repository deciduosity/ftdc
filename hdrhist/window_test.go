package hdrhist_test

import (
	"testing"

	"github.com/deciduosity/ftdc/hdrhist"
	"github.com/stretchr/testify/assert"
)

func TestWindowedHistogram(t *testing.T) {
	w := hdrhist.NewWindowed(2, 1, 1000, 3)

	for i := 0; i < 100; i++ {
		assert.NoError(t, w.Current.RecordValue(int64(i)))
	}
	w.Rotate()

	for i := 100; i < 200; i++ {
		assert.NoError(t, w.Current.RecordValue(int64(i)))
	}
	w.Rotate()

	for i := 200; i < 300; i++ {
		assert.NoError(t, w.Current.RecordValue(int64(i)))
	}

	if v, want := w.Merge().ValueAtQuantile(50), int64(199); v != want {
		t.Errorf("Median was %v, but expected %v", v, want)
	}
}

func BenchmarkWindowedHistogramRecordAndRotate(b *testing.B) {
	w := hdrhist.NewWindowed(3, 1, 10000000, 3)
	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		if err := w.Current.RecordValue(100); err != nil {
			b.Fatal(err)
		}

		if i%100000 == 1 {
			w.Rotate()
		}
	}
}

func BenchmarkWindowedHistogramMerge(b *testing.B) {
	w := hdrhist.NewWindowed(3, 1, 10000000, 3)
	for i := 0; i < 10000000; i++ {
		if err := w.Current.RecordValue(100); err != nil {
			b.Fatal(err)
		}

		if i%100000 == 1 {
			w.Rotate()
		}
	}
	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		w.Merge()
	}
}
