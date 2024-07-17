package day06chronalcoordinates_test

import (
	"testing"

	. "github.com/ewoutquax/advent-of-code-2018/internal/day-06-chronal-coordinates"
	"github.com/stretchr/testify/assert"
)

func TestParseInput(t *testing.T) {
	universe := ParseInput(testInput())

	assert := assert.New(t)
	assert.IsType(universe, Universe{})
	assert.Len(universe.Areas, 6)

	assert.Equal(1, universe.Areas[0].X)
	assert.Equal(1, universe.Areas[0].Y)
	assert.Equal(0, universe.Areas[0].Size)
	assert.Equal(8, universe.Areas[5].X)
	assert.Equal(9, universe.Areas[5].Y)
	assert.Equal(0, universe.Areas[5].Size)

	assert.Equal(1, universe.MinX)
	assert.Equal(1, universe.MinY)
	assert.Equal(8, universe.MaxX)
	assert.Equal(9, universe.MaxY)
}

func TestMaxFiniteSize(t *testing.T) {
	universe := ParseInput(testInput())
	maxSize := MaxFiniteSize(universe)

	assert.Equal(t, 17, maxSize)
}

func TestSizeOfRegionWithDistanceToAllBelowThreshold(t *testing.T) {
	universe := ParseInput(testInput())

	var size int = SizeOfRegionWithDistanceToAllBelowThreshold(universe, 32)
	assert.Equal(t, 16, size)
}

func testInput() []string {
	return []string{
		"1, 1",
		"1, 6",
		"8, 3",
		"3, 4",
		"5, 5",
		"8, 9",
	}
}
