package day07thesumofitsparts_test

import (
	"testing"

	. "github.com/ewoutquax/advent-of-code-2018/internal/day-07-the-sum-of-its-parts"
	"github.com/stretchr/testify/assert"
)

func TestParseInput(t *testing.T) {
	var universe = ParseInput(testInput(), 0)

	assert := assert.New(t)

	assert.Len(universe.Items, 6)

	assert.IsType(universe, Universe{})
	assert.Len(universe.Orders, 7)
	assert.Equal("C", universe.Orders[0].Before.Name)
	assert.Equal("A", universe.Orders[0].After.Name)
	assert.Equal("F", universe.Orders[len(universe.Orders)-1].Before.Name)
	assert.Equal("E", universe.Orders[len(universe.Orders)-1].After.Name)

	assert.Equal(3, universe.Orders[0].Before.RemainingBuildTime)
	assert.Equal(1, universe.Orders[0].After.RemainingBuildTime)
	assert.Equal(6, universe.Orders[len(universe.Orders)-1].Before.RemainingBuildTime)
	assert.Equal(5, universe.Orders[len(universe.Orders)-1].After.RemainingBuildTime)

	assert.Len(universe.Items["C"].BlockedBy, 0)
	assert.Len(universe.Items["A"].BlockedBy, 1)
	assert.Equal("C", universe.Items["A"].BlockedBy[0].Name)
	assert.False(universe.Items["C"].IsInProgress)
	assert.False(universe.Items["A"].IsInProgress)
}

func TestItemCanBeStarted(t *testing.T) {
	universe := ParseInput(testInput(), 0)

	assert.True(t, universe.Items["C"].CanBeStarted())
}

func TestItemCanNotBeStartedWhenBlocked(t *testing.T) {
	universe := ParseInput(testInput(), 0)

	assert.False(t, universe.Items["A"].CanBeStarted())
}

func TestItemCanBeStartedWhenBlockedByCompletedItems(t *testing.T) {
	universe := ParseInput(testInput(), 0)
	universe.Items["C"].RemainingBuildTime = 0

	assert.True(t, universe.Items["A"].CanBeStarted())
}

func TestItemCanNotBeStartedWhenAlreadyInProgress(t *testing.T) {
	universe := ParseInput(testInput(), 0)
	universe.Items["C"].IsInProgress = true

	assert.False(t, universe.Items["C"].CanBeStarted())
}

func TestItemCanNotBeStartedWhenDone(t *testing.T) {
	universe := ParseInput(testInput(), 0)
	universe.Items["C"].IsInProgress = false
	universe.Items["C"].RemainingBuildTime = 0

	assert.False(t, universe.Items["C"].CanBeStarted())
}

func TestBuildMetrics(t *testing.T) {
	universe := ParseInput(testInput(), 0)
	order, _ := BuildMetrics(universe, 1)

	assert.Equal(t, "CABDFE", order)
}

func TestCalculateBuildTimeWithWorkers(t *testing.T) {
	universe := ParseInput(testInput(), 0)
	_, elapsed := BuildMetrics(universe, 2)

	assert.Equal(t, 15, elapsed)
}

func testInput() []string {
	return []string{
		"Step C must be finished before step A can begin.",
		"Step C must be finished before step F can begin.",
		"Step A must be finished before step B can begin.",
		"Step A must be finished before step D can begin.",
		"Step B must be finished before step E can begin.",
		"Step D must be finished before step E can begin.",
		"Step F must be finished before step E can begin.",
	}
}
