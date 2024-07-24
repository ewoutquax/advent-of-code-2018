package day17reservoirresearch_test

import (
	"testing"

	. "github.com/ewoutquax/advent-of-code-2018/internal/day-17-reservoir-research"
	"github.com/stretchr/testify/assert"
)

func TestParseInput(t *testing.T) {
	universe := ParseInput(testInput())

	assert := assert.New(t)
	assert.IsType(universe, Universe{})
	assert.Len(universe.Items, 34)
	assert.Equal(0, universe.WaterDrops)
	assert.Equal(13, universe.MaxY)
}

func TestAddInitialWater(t *testing.T) {
	universe := ParseInput(testInput())
	universe.AddInitialWater()

	assert.Len(t, universe.ActiveWaterItems, 1)
	item := universe.ActiveWaterItems[0].(*FallingWater)
	assert.Equal(t, ItemType("falling"), item.Type)
	assert.Equal(t, 0, universe.WaterDrops)
}

func TestSimulate6Moves(t *testing.T) {
	universe := ParseInput(testInput())
	universe.AddInitialWater()

	for move := 6; move > 0; move-- {
		universe.Simulate()
	}

	assert.Len(t, universe.ActiveWaterItems, 1)
	assert.Equal(t, universe.ActiveWaterItems[0].(*FallingWater).Location.X, 500)
	assert.Equal(t, universe.ActiveWaterItems[0].(*FallingWater).Location.Y, 6)
	assert.Equal(t, 6, universe.WaterDrops)
}

func TestSimulate7Moves(t *testing.T) {
	// The falling stream converts into settled water, and spawn just a left-sliding stream
	universe := ParseInput(testInput())
	universe.AddInitialWater()

	for move := 7; move > 0; move-- {
		universe.Simulate()
	}

	assert := assert.New(t)
	assert.Len(universe.ActiveWaterItems, 1)
	assert.Len(universe.Items, 36)
	assert.IsType(universe.ActiveWaterItems[0], &SlidingWater{})
	assert.Equal(universe.ActiveWaterItems[0].(*SlidingWater).Location.X, 499)
	assert.Equal(universe.ActiveWaterItems[0].(*SlidingWater).Location.Y, 6)
	assert.Equal(7, universe.WaterDrops)
}

func TestSimulate10Moves(t *testing.T) {
	// The sliding stream will move to the end of the block, and not have hit the wall yet
	universe := ParseInput(testInput())
	universe.AddInitialWater()

	for move := 10; move > 0; move-- {
		universe.Simulate()
	}

	assert := assert.New(t)
	assert.Len(universe.ActiveWaterItems, 1)
	assert.Len(universe.Items, 39)
	assert.IsType(universe.ActiveWaterItems[0], &SlidingWater{})
	assert.Equal(universe.ActiveWaterItems[0].(*SlidingWater).Location.X, 496)
	assert.Equal(universe.ActiveWaterItems[0].(*SlidingWater).Location.Y, 6)
	assert.Equal(10, universe.WaterDrops)
}

func TestSimulate11Moves(t *testing.T) {
	// The sliding stream has reached a wall, and will raise one level
	universe := ParseInput(testInput())
	universe.AddInitialWater()

	for move := 11; move > 0; move-- {
		universe.Simulate()
	}

	assert := assert.New(t)
	assert.Len(universe.ActiveWaterItems, 1)
	assert.Len(universe.Items, 40)
	assert.IsType(universe.ActiveWaterItems[0], &SlidingWater{})
	assert.Equal(universe.ActiveWaterItems[0].(*SlidingWater).Location.X, 499)
	assert.Equal(universe.ActiveWaterItems[0].(*SlidingWater).Location.Y, 5)
	assert.Equal(11, universe.WaterDrops)
}

func TestSimulate17Moves(t *testing.T) {
	// The bucket has been filled to the rim, and will start splilling in the next move
	universe := ParseInput(testInput())
	universe.AddInitialWater()

	for move := 17; move > 0; move-- {
		universe.Simulate()
	}

	assert := assert.New(t)
	assert.Len(universe.ActiveWaterItems, 2)
	assert.Len(universe.Items, 47) // also a stream going to the right now
	assert.IsType(universe.ActiveWaterItems[0], &SlidingWater{})
	assert.Equal(universe.ActiveWaterItems[0].(*SlidingWater).Location.X, 499)
	assert.Equal(universe.ActiveWaterItems[0].(*SlidingWater).Location.Y, 2)
	assert.IsType(universe.ActiveWaterItems[1], &SlidingWater{})
	assert.Equal(universe.ActiveWaterItems[1].(*SlidingWater).Location.X, 501)
	assert.Equal(universe.ActiveWaterItems[1].(*SlidingWater).Location.Y, 2)
	assert.Equal(18, universe.WaterDrops)
}

func TestSimulate18Moves(t *testing.T) {
	// The left-sliding is halted by a wall, and the right-sliding goes over the edge, becomes a falling and drop one position
	universe := ParseInput(testInput())
	universe.AddInitialWater()

	for move := 18; move > 0; move-- {
		universe.Simulate()
	}

	assert := assert.New(t)
	assert.Len(universe.ActiveWaterItems, 1)
	assert.Len(universe.Items, 48) // The right-stream did leave a still-water, before becoming a falling
	assert.IsType(universe.ActiveWaterItems[0], &FallingWater{})
	assert.Equal(universe.ActiveWaterItems[0].(*FallingWater).Location.X, 502)
	assert.Equal(universe.ActiveWaterItems[0].(*FallingWater).Location.Y, 3)
	assert.Equal(20, universe.WaterDrops) // The right-stream has become a falling, and already dropped one
}

func TestSimulate36Moves(t *testing.T) {
	// The falling stream landed in the bucket below, and filled it till the rim
	universe := ParseInput(testInput())
	universe.AddInitialWater()

	for move := 36; move > 0; move-- {
		universe.Simulate()
	}

	assert := assert.New(t)
	assert.Len(universe.ActiveWaterItems, 1) // The stream to the right has settled 2 moves ago
	assert.Len(universe.Items, 61)           // Only the 15 spaces in the bucket are filled with still water, minus the 3 (or 2) that used to be the falling
	assert.IsType(universe.ActiveWaterItems[0], &SlidingWater{})
	assert.Equal(499, universe.ActiveWaterItems[0].(*SlidingWater).Location.X)
	assert.Equal(10, universe.ActiveWaterItems[0].(*SlidingWater).Location.Y)
	assert.Equal(41, universe.WaterDrops)
}

func TestSimulate37Moves(t *testing.T) {
	universe := ParseInput(testInput())
	universe.AddInitialWater()

	for move := 37; move > 0; move-- {
		universe.Simulate()
	}

	assert := assert.New(t)
	assert.Len(universe.Items, 63)
	assert.Len(universe.ActiveWaterItems, 2)
	assert.IsType(universe.ActiveWaterItems[0], &SlidingWater{})
	assert.Equal(501, universe.ActiveWaterItems[0].(*SlidingWater).Location.X)
	assert.Equal(9, universe.ActiveWaterItems[0].(*SlidingWater).Location.Y)
	assert.IsType(universe.ActiveWaterItems[1], &SlidingWater{})
	assert.Equal(503, universe.ActiveWaterItems[1].(*SlidingWater).Location.X)
	assert.Equal(9, universe.ActiveWaterItems[1].(*SlidingWater).Location.Y)
	assert.Equal(43, universe.WaterDrops)
}

func TestSimulateFlow(t *testing.T) {
	// Both buckets have overflown, and the falling streams have fallen out of reach
	universe := ParseInput(testInput())
	universe.AddInitialWater()

	SimulateFlow(&universe)

	assert := assert.New(t)
	assert.Len(universe.ActiveWaterItems, 0) // The stream to the right has settled 2 moves ago
	assert.Equal(57, universe.WaterDrops)
}

func testInput() []string {
	return []string{
		"x=495, y=2..7",
		"y=7, x=495..501",
		"x=501, y=3..7",
		"x=498, y=2..4",
		"x=506, y=1..2",
		"x=498, y=10..13",
		"x=504, y=10..13",
		"y=13, x=498..504",
	}
}
