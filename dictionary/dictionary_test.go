package dictionary

import (
	"math/rand"
	"testing"
)

func newDictionary(t testing.TB, bits []bool) *Dictionary {
	t.Helper()
	d := New(len(bits))
	for i, b := range bits {
		d.SetBit(i, b)
	}
	d.CreateIndex()
	return d
}

func randBits(t testing.TB, size int) []bool {
	t.Helper()
	b := make([]bool, size)
	for i := range b {
		b[i] = rand.Int()%2 == 0
	}
	return b
}

func oneBits(t testing.TB, size int) []bool {
	t.Helper()
	b := make([]bool, size)
	for i := range b {
		b[i] = true
	}
	return b
}

func zeroBits(t testing.TB, size int) []bool {
	t.Helper()
	return make([]bool, size)
}

func TestNew(t *testing.T) {
	size := 1000000
	tests := map[string][]bool{
		"all zero bits": zeroBits(t, size),
		"all one bits":  oneBits(t, size),
		"random bits":   randBits(t, size),
	}
	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			s := len(test)
			d := New(s)
			l := s / bitsSize
			if s%bitsSize != 0 {
				l++
			}
			if ds := bitsSize * len(d.bits); ds < s {
				t.Fatalf("bitsSize * len(d.bits) = %d; is less than len(test) = %d", ds, s)
			}
			if got, want := len(d.bits), l; got != want {
				t.Fatalf("len(d.bits) = %d; want %d", got, want)
			}
		})
	}
}

func assertSetBit(t *testing.T, d *Dictionary, wants []bool) {
	t.Helper()
	gots := make([]bool, len(wants))
	for i, b := range d.bits {
		for j := 0; j < bitsSize; j++ {
			if b&(1<<j) > 0 {
				gots[i*bitsSize+j] = true
			}
		}
	}
	for i := range wants {
		if got, want := gots[i], wants[i]; got != want {
			t.Fatalf("pos = %d; got = %t; want %t", i, got, want)
		}
	}
}

func TestSetBit(t *testing.T) {
	size := 1000000
	tests := map[string][]bool{
		"all zero bits": zeroBits(t, size),
		"all one bits":  oneBits(t, size),
		"random bits":   randBits(t, size),
	}
	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			d := New(len(test))
			for i, b := range test {
				d.SetBit(i, b)
			}
			assertSetBit(t, d, test)
		})
	}
}

func TestSetBitZero(t *testing.T) {
	size := 1000000
	tests := map[string][]bool{
		"all zero bits": zeroBits(t, size),
		"all one bits":  oneBits(t, size),
		"random bits":   randBits(t, size),
	}
	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			d := New(len(test))
			for i := range test {
				d.SetBit(i, true)
			}
			for i, b := range test {
				d.SetBit(i, b)
			}
			assertSetBit(t, d, test)
		})
	}
}

func TestBit(t *testing.T) {
	size := 1000000
	tests := map[string][]bool{
		"all zero bits": zeroBits(t, size),
		"all one bits":  oneBits(t, size),
		"random bits":   randBits(t, size),
	}
	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			d := New(len(test))
			for i, b := range test {
				d.SetBit(i, b)
			}
			for i := range test {
				if got, want := d.Bit(i), test[i]; got != want {
					t.Fatalf("d.Bit(%d) = %t; want = %t", i, got, want)
				}
			}
		})
	}
}

func TestLen(t *testing.T) {
	tests := map[string]struct {
		input []bool
		want  int
	}{
		"less": {
			input: zeroBits(t, bitsSize-1),
			want:  bitsSize,
		},
		"equal": {
			input: zeroBits(t, bitsSize),
			want:  bitsSize,
		},
		"more": {
			input: zeroBits(t, bitsSize+1),
			want:  bitsSize * 2,
		},
		"large number": {
			input: zeroBits(t, bitsSize*1234+5),
			want:  bitsSize * 1235,
		},
	}
	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			d := New(len(test.input))
			if got := d.Len(); got != test.want {
				t.Fatalf("d.Len() = %d; want = %d", got, test.want)
			}
		})
	}
}

func TestRank(t *testing.T) {
	size := 1000000
	tests := map[string][]bool{
		"all zero bits": zeroBits(t, size),
		"all one bits":  oneBits(t, size),
		"random bits":   randBits(t, size),
	}
	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			d := newDictionary(t, test)
			want := 0
			for i, b := range test {
				if b {
					want++
				}
				if got := d.Rank(i); got != want {
					t.Fatalf("d.Rank(%d) = %d; want %d", i, got, want)
				}
			}
		})
	}
}

func TestRank0(t *testing.T) {
	size := 1000000
	tests := map[string][]bool{
		"all zero bits": zeroBits(t, size),
		"all one bits":  oneBits(t, size),
		"random bits":   randBits(t, size),
	}
	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			d := newDictionary(t, test)
			want := 0
			for i, b := range test {
				if !b {
					want++
				}
				if got := d.Rank0(i); got != want {
					t.Fatalf("d.Rank0(%d) = %d; want %d", i, got, want)
				}
			}
		})
	}
}

func TestSelect(t *testing.T) {
	size := 1000000
	tests := map[string][]bool{
		"all zero bits": zeroBits(t, size),
		"all one bits":  oneBits(t, size),
		"random bits":   randBits(t, size),
	}
	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			d := newDictionary(t, test)
			rank := 0
			td := make(map[int]int)
			for i, b := range test {
				if b {
					rank++
					td[rank] = i
				}
			}
			td[len(test)+100] = d.Len()
			for rank, want := range td {
				if got := d.Select(rank); got != want {
					t.Fatalf("d.Select(%d) = %d; want %d", rank, got, want)
				}
			}
		})
	}
}

func TestSelect0(t *testing.T) {
	size := 1000000
	tests := map[string][]bool{
		"all zero bits": zeroBits(t, size),
		"all one bits":  oneBits(t, size),
		"random bits":   randBits(t, size),
	}
	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			d := newDictionary(t, test)
			rank := 0
			td := make(map[int]int)
			for i, b := range test {
				if !b {
					rank++
					td[rank] = i
				}
			}
			td[len(test)+100] = d.Len()
			for rank, want := range td {
				if got := d.Select0(rank); got != want {
					t.Fatalf("d.Select0(%d) = %d; want %d", rank, got, want)
				}
			}
		})
	}
}

func BenchmarkCreateIndex(b *testing.B) {
	size := 10000000
	r := randBits(b, size)
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			d := New(size)
			for i, b := range r {
				d.SetBit(i, b)
			}
			d.CreateIndex()
		}
	})
}

func BenchmarkRank(b *testing.B) {
	size := 10000000
	d := newDictionary(b, randBits(b, size))
	r := rand.Intn(size)
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			d.Rank(r)
		}
	})
}

func BenchmarkSelect(b *testing.B) {
	size := 10000000
	d := newDictionary(b, randBits(b, size))
	r := rand.Intn(d.Rank(d.Len() - 1))
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			d.Select(r)
		}
	})
}

func BenchmarkSelect0(b *testing.B) {
	size := 10000000
	d := newDictionary(b, randBits(b, size))
	r := rand.Intn(d.Rank(d.Len() - 1))
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			d.Select0(r)
		}
	})
}

func FuzzRank(f *testing.F) {
	for _, size := range []uint32{1, 10, 100, 1000, 10000, 100000, 1000000, 10000000} {
		f.Add(size)
	}
	f.Fuzz(func(t *testing.T, size uint32) {
		bs := randBits(t, int(size))
		d := newDictionary(t, bs)
		want := 0
		for i, b := range bs {
			if b {
				want++
			}
			if got := d.Rank(i); got != want {
				t.Fatalf("size = %d; d.Rank(%d) = %d; want %d", size, i, got, want)
			}
		}
	})
}

func FuzzSelect(f *testing.F) {
	for _, size := range []uint32{1, 10, 100, 1000, 10000, 100000, 1000000, 10000000} {
		f.Add(size)
	}
	f.Fuzz(func(t *testing.T, size uint32) {
		bs := randBits(t, int(size))
		d := newDictionary(t, bs)
		rank := 0
		td := make(map[int]int)
		for i, b := range bs {
			if b {
				rank++
				td[rank] = i
			}
		}
		td[len(bs)+100] = d.Len()
		for rank, want := range td {
			if got := d.Select(rank); got != want {
				t.Fatalf("size = %d; d.Select(%d) = %d; want %d", size, rank, got, want)
			}
		}
	})
}

func FuzzSelect0(f *testing.F) {
	for _, size := range []uint32{1, 10, 100, 1000, 10000, 100000, 1000000, 10000000} {
		f.Add(size)
	}
	f.Fuzz(func(t *testing.T, size uint32) {
		bs := randBits(t, int(size))
		d := newDictionary(t, bs)
		rank := 0
		td := make(map[int]int)
		for i, b := range bs {
			if !b {
				rank++
				td[rank] = i
			}
		}
		td[len(bs)+100] = d.Len()
		for rank, want := range td {
			if got := d.Select0(rank); got != want {
				t.Fatalf("size = %d; d.Select0(%d) = %d; want %d", size, rank, got, want)
			}
		}
	})
}
