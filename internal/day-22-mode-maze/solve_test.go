package day22modemaze_test

import (
	"fmt"
	"testing"

	. "github.com/ewoutquax/advent-of-code-2018/internal/day-22-mode-maze"
	"github.com/stretchr/testify/assert"
)

func TestGeologicalIndex(t *testing.T) {
	targetLocation := Location{10, 10}
	caveDepth := CaveDepth(510)

	testCases := map[Location]int{
		{0, 0}: 0,
		{1, 0}: 16807,
		{0, 1}: 48271,
		{1, 1}: 145722555,
	}

	for inputLocation, expectedResult := range testCases {
		actualResult := inputLocation.GeologicalIndex(targetLocation, caveDepth)
		assert.Equal(t, expectedResult, actualResult)
	}
}

func TestErosionLevel(t *testing.T) {
	targetLocation := Location{10, 10}
	caveDepth := CaveDepth(510)

	testCases := map[Location]int{
		{0, 0}: 510,
		{1, 0}: 17317,
		{0, 1}: 8415,
		{1, 1}: 1805,
	}

	for inputLocation, expectedResult := range testCases {
		actualResult := inputLocation.ErosionLevel(targetLocation, caveDepth)
		assert.Equal(t, expectedResult, actualResult)
	}
}

func TestErosionLevelType(t *testing.T) {
	targetLocation := Location{10, 10}
	caveDepth := CaveDepth(510)

	testCases := map[Location]ErosionLevelType{
		{0, 0}: ErosionLevelTypeRocky,
		{1, 0}: ErosionLevelTypeWet,
		{0, 1}: ErosionLevelTypeRocky,
		{1, 1}: ErosionLevelTypeNarrow,
	}

	for inputLocation, expectedResult := range testCases {
		actualResult := inputLocation.ErosionLevelType(targetLocation, caveDepth)
		assert.Equal(t, expectedResult, actualResult)
	}
}

func TestCalculateRisk(t *testing.T) {
	targetLocation := Location{10, 10}
	caveDepth := CaveDepth(510)

	assert.Equal(t, 114, CalculateRisk(targetLocation, caveDepth))
}

func TestDrawCave(t *testing.T) {
	targetLocation := Location{10, 10}
	caveDepth := CaveDepth(510)

	// cachedErosionLevels = make(map[Location]int, targetLocation.X*targetLocation.Y)

	for y := 0; y < 16; y++ {
		for x := 0; x < 16; x++ {
			location := Location{
				X: x,
				Y: y,
			}
			locationType := location.ErosionLevelType(targetLocation, caveDepth)

			stringType := map[ErosionLevelType]string{
				ErosionLevelTypeRocky:  ".",
				ErosionLevelTypeWet:    "=",
				ErosionLevelTypeNarrow: "|",
			}[locationType]

			fmt.Printf("%s", stringType)
		}
		fmt.Println("")
	}

	assert.True(t, false)
}

func TestParseInput(t *testing.T) {
	input := ParseInput(testInput())

	assert.IsType(t, Input{}, input)
	assert.Equal(t, CaveDepth(4002), input.CaveDepth)
	assert.Equal(t, 5, input.TargetLocation.X)
	assert.Equal(t, 746, input.TargetLocation.Y)
}

func TestPathVisitedKey(t *testing.T) {
	path := Path{
		Location:         Location{42, 1337},
		CurrentEquipment: EquipmentTorch,
		NrMinutes:        0,
	}

	assert.Equal(t, VisitedPathKey("[42, 1337]/5"), path.VisitedKey())
}

func TestPreventDuplicateVisitedPaths(t *testing.T) {
	InitVisitedPaths()

	path := Path{
		Location:         Location{1337, 42},
		CurrentEquipment: EquipmentTorch,
		NrMinutes:        18,
	}

	assert.False(t, path.IsVisited())

	path.SetVisited()
	assert.True(t, path.IsVisited())

	// Changing gear with see this as a new entry
	path.CurrentEquipment = EquipmentClimbingGear
	assert.False(t, path.IsVisited())

	// Lowering the NrMinutes should see the new path as a better alternative
	path.SetVisited()
	assert.True(t, path.IsVisited())
	path.NrMinutes -= 1
	assert.False(t, path.IsVisited())
}

func TestFastestTime(t *testing.T) {
	targetLocation := Location{10, 10}
	caveDepth := CaveDepth(510)

	assert.Equal(t, 45, FastestTime(targetLocation, caveDepth))
}

func testInput() []string {
	return []string{
		"depth: 4002",
		"target: 5,746",
	}
}

func solutionExample2() []string {
	return []string{
		"[0, 0] / tool Torch: 0",
		"[0, 1] / tool Torch: 1",
		"[1, 1] / tool Torch: 2",
		"[1, 1] / tool None: 9",
		"[2, 1] / tool None: 10",
		"[3, 1] / tool None: 11",
		"[4, 1] / tool None: 12",
		"[4, 1] / tool Climbing gear: 19",
		"[4, 2] / tool Climbing gear: 20",
		"[4, 3] / tool Climbing gear: 21",
		"[4, 4] / tool Climbing gear: 22",
		"[4, 5] / tool Climbing gear: 23",
		"[4, 6] / tool Climbing gear: 24",
		"[4, 7] / tool Climbing gear: 25",
		"[4, 8] / tool Climbing gear: 26",
		"[5, 8] / tool Climbing gear: 27",
		"[5, 9] / tool Climbing gear: 28",
		"[5, 10] / tool Climbing gear: 29",
		"[5, 11] / tool Climbing gear: 30",
		"[6, 11] / tool Climbing gear: 31",
		"[6, 12] / tool Climbing gear: 32",
		"[7, 12] / tool Climbing gear: 33",
		"[8, 12] / tool Climbing gear: 34",
		"[9, 12] / tool Climbing gear: 35",
		"[10, 12] / tool Climbing gear: 36",
		"[10, 11] / tool Climbing gear: 37",
		"[10, 10] / tool Climbing gear: 38",
		"[10, 10] / tool Torch: 45",
	}
}
