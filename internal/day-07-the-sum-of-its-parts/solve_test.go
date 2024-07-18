package day07thesumofitsparts_test

import (
	"testing"

	. "github.com/ewoutquax/advent-of-code-2018/internal/day-07-the-sum-of-its-parts"
	"github.com/stretchr/testify/assert"
)

func TestParseInput(t *testing.T) {
	var universe = ParseInput(testInput())

	assert.IsType(t, universe, Universe{})
	assert.Len(t, universe.Orders, 7)
	assert.Equal(t, Item("C"), universe.Orders[0].Before)
	assert.Equal(t, Item("A"), universe.Orders[0].After)
	assert.Equal(t, Item("F"), universe.Orders[len(universe.Orders)-1].Before)
	assert.Equal(t, Item("E"), universe.Orders[len(universe.Orders)-1].After)
}

func TestBuildOrder(t *testing.T) {
	universe := ParseInput(testInput())
	order := BuildOrder(universe)

	assert.Equal(t, "CABDFE", order)
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
