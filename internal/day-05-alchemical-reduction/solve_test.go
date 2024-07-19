package day05alchemicalreduction_test

import (
	"testing"

	. "github.com/ewoutquax/advent-of-code-2018/internal/day-05-alchemical-reduction"
	"github.com/stretchr/testify/assert"
)

func TestReacts(t *testing.T) {
	testCases := map[string]bool{
		"aa": false,
		"aA": true,
		"Aa": true,
		"AA": false,
		"aB": false,
	}

	for input, expectedOutput := range testCases {
		output := DoesReact(input)
		assert.Equal(t, expectedOutput, output, input)

	}
}
func TestTriggerPolymers(t *testing.T) {
	result := TriggerPolymer("dabAcCaCBAcCcaDA")
	assert.Equal(t, "dabCBAcaDA", result)
}

func TestPolymerLengthAfterTrigger(t *testing.T) {
	length := PolymerLengthAfterTrigger("dabAcCaCBAcCcaDA")
	assert.Equal(t, 10, length)
}

func TestShortestPolymerWithExtraction(t *testing.T) {
	shortest := ShortestPolymerWithExtraction("dabAcCaCBAcCcaDA")
	assert.Equal(t, "daDA", shortest)
}

func TestShortestPolymerLengthWithExtraction(t *testing.T) {
	length := ShortestPolymerLengthWithExtraction("dabAcCaCBAcCcaDA")
	assert.Equal(t, 4, length)
}
