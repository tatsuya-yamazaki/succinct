// Package dictionary provides a succinct indexable dictionary with space-efficient support for rank and select operations.
//
// This package implements a compressed bit vector data structure that allows efficient querying of bit positions.
// It supports the following operations:
//
//   - Rank: Returns the number of 1-bits up to a given position.
//   - Rank0: Returns the number of 0-bits up to a given position.
//   - Select: Finds the smallest position of the 1-bit with the specified rank.
//   - Select0: Finds the smallest position of the 0-bit with the specified rank.
//
// The data structure is designed to minimize memory usage while providing fast access to bit-level information.
// It is useful for applications where memory constraints are critical, and quick bit manipulation is required.
package dictionary

import (
	"math/bits"
)

const (
	// bitsSize represents the number of bits in each element of the underlying bit array (1 byte = 8 bits).
	bitsSize = 8

	// bitsPerRankIndexLarge specifies the number of bit units that one large rank index entry spans.
	bitsPerRankIndexLarge = 8191
)

// Dictionary represents a succinct bit vector with rank and select operations.
// It stores the bits as a slice of bytes and maintains rank indexes for efficient queries.
type Dictionary struct {
	// Array of bits stored in bytes.
	bits []uint8
	// Rank index structure for efficient rank operations.
	rank rankIndex
}

// New creates a new Dictionary of the specified size.
// It initializes the bit vector and prepares it for bit manipulations.
func New(size int) *Dictionary {
	d := new(Dictionary)
	l := size / bitsSize
	if size%bitsSize > 0 {
		l++
	}
	d.bits = make([]uint8, l)
	return d
}

// Len returns the total number of bits in the bit vector.
func (d *Dictionary) Len() int {
	return len(d.bits) * bitsSize
}

// bitsIndex returns the index in the byte array (bits) that corresponds to the given bit position.
func (d *Dictionary) bitsIndex(pos int) int {
	return pos / bitsSize
}

// bitPos returns the position of a specific bit within a byte in the bit array.
func (d *Dictionary) bitPos(pos int) int {
	return pos % bitsSize
}

// SetBit sets the bit at the given position to either 1 (true) or 0 (false).
// The flag parameter determines whether to set or clear the bit.
func (d *Dictionary) SetBit(pos int, flag bool) {
	bi := d.bitsIndex(pos)
	var b uint8 = 1 << d.bitPos(pos)
	if flag {
		d.bits[bi] |= b
	} else {
		d.bits[bi] &^= b
	}
}

// Bit returns true if the bit at the given position is 1, and false otherwise.
func (d *Dictionary) Bit(pos int) bool {
	return d.bits[d.bitsIndex(pos)]&(1<<d.bitPos(pos)) > 0
}

// CreateIndex builds the index for efficient rank and select operations.
func (d *Dictionary) CreateIndex() {
	d.rank = newRankIndex(d.Len())
	c := 0
	for i, b := range d.bits {
		c += oneBitsCount(b, bitsSize-1)
		d.rank.update(i, c)
	}
}

// Rank returns the number of 1-bits in the bit vector up to the given position.
// It efficiently calculates the number of set bits from the start to the specified position.
func (d *Dictionary) Rank(pos int) int {
	bi := d.bitsIndex(pos)
	r := d.rank.rank(bi)
	r += oneBitsCount(d.bits[bi], d.bitPos(pos))
	return r
}

// Rank0 returns the number of 0-bits in the bit vector up to the given position.
// It efficiently calculates the number of set bits from the start to the specified position.
func (d *Dictionary) Rank0(pos int) int {
	return pos - d.Rank(pos) + 1
}

// Select returns the smallest position of the 1-bit with the specified rank in the bit vector.
// It efficiently finds the first occurrence of the specified number of set bits.
func (d *Dictionary) Select(rank int) (pos int) {
	l, r := 0, d.Len()
	for l != r {
		m := (l + r) / 2
		if d.Rank(m) < rank {
			l = m + 1
		} else {
			r = m
		}
	}
	return l
}

// Select0 returns the smallest position of the 0-bit with the specified rank in the bit vector.
// It efficiently finds the first occurrence of the specified number of set bits.
func (d *Dictionary) Select0(rank int) (pos int) {
	l, r := 0, d.Len()
	for l != r {
		m := (l + r) / 2
		if d.Rank0(m) < rank {
			l = m + 1
		} else {
			r = m
		}
	}
	return l
}

// rankIndex represents the structure used to store rank information for efficient rank query.
type rankIndex struct {
	// Small rank index storing cumulative counts within smaller intervals.
	small []uint16
	// Large rank index storing cumulative counts at larger intervals (e.g., every bitsPerRankIndexLarge bits).
	large []int
}

// newRankIndex creates a new rankIndex of the specified size.
func newRankIndex(size int) rankIndex {
	sl := size/bitsSize + 1
	if size%bitsSize != 0 {
		sl++
	}

	ls := bitsSize * bitsPerRankIndexLarge
	ll := size/ls + 1
	if size%ls != 0 {
		ll++
	}

	return rankIndex{
		small: make([]uint16, sl),
		large: make([]int, ll),
	}
}

// largeIndex returns the index in the large rank index that corresponds to the given bitsIndex.
func (r *rankIndex) largeIndex(bitsIndex int) int {
	return bitsIndex / bitsPerRankIndexLarge
}

// update updates the rank index with the cumulative count of 1-bits at the given bitsIndex.
func (r *rankIndex) update(bitsIndex, onesCount int) {
	li := r.largeIndex(bitsIndex)
	if bitsIndex%bitsPerRankIndexLarge == bitsPerRankIndexLarge-1 {
		r.large[li+1] = onesCount
		return
	}
	r.small[bitsIndex+1] = uint16(onesCount - r.large[li])
}

// rank returns the cumulative number of 1-bits up to the given bitsIndex.
func (r *rankIndex) rank(bitsIndex int) int {
	return int(r.large[r.largeIndex(bitsIndex)] + int(r.small[bitsIndex]))
}

// oneBitsCount returns the number of 1-bits in the byte x, up to the given bit position pos.
func oneBitsCount(x uint8, pos int) int {
	return bits.OnesCount8(x & uint8((1<<(pos+1))-1))
}
