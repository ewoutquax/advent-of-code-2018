package day23experimentalemergencyteleportation_test

import (
	"testing"

	. "github.com/ewoutquax/advent-of-code-2018/internal/day-23-experimental-emergency-teleportation"
	"github.com/stretchr/testify/assert"
)

func TestParseInput(t *testing.T) {
	assert := assert.New(t)

	bots := ParseInput(testInput())

	assert.IsType(Bot{}, bots[0])
	assert.Len(bots, 9)

	firstBot := bots[0]
	lastBot := bots[len(bots)-1]

	assert.IsType(Bot{}, firstBot)
	assert.Equal(0, firstBot.X)
	assert.Equal(0, firstBot.Y)
	assert.Equal(0, firstBot.Z)
	assert.Equal(Radius(4), firstBot.Radius)

	assert.IsType(Bot{}, firstBot)
	assert.Equal(1, lastBot.X)
	assert.Equal(3, lastBot.Y)
	assert.Equal(1, lastBot.Z)
	assert.Equal(Radius(1), lastBot.Radius)
}

func TestDistanceBetweenBots(t *testing.T) {
	bots := ParseInput(testInput())

	testCases := map[Bot]int{
		bots[0]: 0,
		bots[1]: 1,
		bots[2]: 4,
		bots[3]: 2,
		bots[4]: 5,
		bots[5]: 3,
		bots[6]: 3,
		bots[7]: 4,
		bots[8]: 5,
	}

	sourceBot := bots[0]
	for targetBot, expectedDistance := range testCases {
		assert.Equal(t, expectedDistance, sourceBot.Distance(targetBot))
	}
}

func TestCountWithinRangeOfStrongest(t *testing.T) {
	bots := ParseInput(testInput())

	assert.Equal(t, 7, CountWithinRangeOfStrongest(bots))
}

func TestCalculateBotWithFullOverlap(t *testing.T) {
	assert := assert.New(t)

	lines := []string{
		"Pos=<0,0,0>, r=4",
		"Pos=<4,0,0>, r=4",
		"Pos=<0,0,0>, r=0",
		"Pos=<4,0,0>, r=0",
	}

	var overlapBot Bot
	var err error

	bots := ParseInput(lines)

	// No overlap exists
	overlapBot, err = GenerateBotWithFullOverlap(bots[2], bots[3])
	assert.NotNil(err)
	assert.ErrorContains(err, "No overlap")
	assert.IsType(Bot{}, overlapBot)

	// Overlap exists with bot1 at the edge
	overlapBot, err = GenerateBotWithFullOverlap(bots[2], bots[1])
	assert.Nil(err)
	assert.IsType(Bot{}, overlapBot)
	assert.Equal(0, overlapBot.X)
	assert.Equal(0, overlapBot.Y)
	assert.Equal(0, overlapBot.Z)
	assert.Equal(Radius(0), overlapBot.Radius)

	// Overlap exists with bot2 at the edge
	overlapBot, err = GenerateBotWithFullOverlap(bots[3], bots[0])
	assert.Nil(err)
	assert.IsType(Bot{}, overlapBot)
	assert.Equal(4, overlapBot.X)
	assert.Equal(0, overlapBot.Y)
	assert.Equal(0, overlapBot.Z)
	assert.Equal(Radius(0), overlapBot.Radius)

	// Overlap exists with bot in the middle
	overlapBot, err = GenerateBotWithFullOverlap(bots[0], bots[1])
	assert.Nil(err)
	assert.IsType(Bot{}, overlapBot)
	assert.Equal(2, overlapBot.X)
	assert.Equal(0, overlapBot.Y)
	assert.Equal(0, overlapBot.Z)
	assert.Equal(Radius(2), overlapBot.Radius)
}

func TestFindLocationWithMostCoverage(t *testing.T) {
	bots := ParseInput(testInput2())

	loc := FindLocationWithMostCoverage(bots)
	assert.IsType(t, Location{}, loc)
	assert.Equal(t, 12, loc.X)
	assert.Equal(t, 12, loc.Y)
	assert.Equal(t, 12, loc.Z)
}

func testInput() []string {
	return []string{
		"Pos=<0,0,0>, r=4",
		"Pos=<1,0,0>, r=1",
		"Pos=<4,0,0>, r=3",
		"Pos=<0,2,0>, r=1",
		"Pos=<0,5,0>, r=3",
		"Pos=<0,0,3>, r=1",
		"Pos=<1,1,1>, r=1",
		"Pos=<1,1,2>, r=1",
		"Pos=<1,3,1>, r=1",
	}
}

func testInput2() []string {
	return []string{
		"pos=<10,12,12>, r=2",
		"pos=<12,14,12>, r=2",
		"pos=<16,12,12>, r=4",
		"pos=<14,14,14>, r=6",
		"pos=<50,50,50>, r=200",
		"pos=<10,10,10>, r=5",
	}
}
