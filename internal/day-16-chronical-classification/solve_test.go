package day16chronicalclassification_test

import (
	"fmt"
	"testing"

	. "github.com/ewoutquax/advent-of-code-2018/internal/day-16-chronical-classification"
	"github.com/stretchr/testify/assert"
)

func TestSetRegisters(t *testing.T) {
	assert := assert.New(t)
	registers := SetRegisters(testInput()[0])

	assert.IsType(Registers{}, registers)
	assert.Equal(3, registers[0])
	assert.Equal(2, registers[1])
	assert.Equal(1, registers[2])
	assert.Equal(1, registers[3])
}

func TestParseInstruction(t *testing.T) {
	assert := assert.New(t)

	line := testInput()[1]
	instruction := ParseInstruction(line)

	assert.IsType(Instruction{}, instruction)
	assert.Equal(2, instruction.InputA)
	assert.Equal(1, instruction.InputB)
	assert.Equal(2, instruction.Output)
}

func TestRegistersToS(t *testing.T) {
	registers := Registers{
		0: 9,
		1: 8,
		2: 7,
		3: 6,
	}

	assert.Equal(t, "[9, 8, 7, 6]", registers.ToS())
}

func TestValidOpcodes(t *testing.T) {
	assert := assert.New(t)
	validOpcodes := ValidOpcodes(testInput())

	assert.Len(validOpcodes, 3)
	assert.Contains(validOpcodes, Mulr)
	assert.Contains(validOpcodes, Addi)
	assert.Contains(validOpcodes, Seti)
}

func TestBorr(t *testing.T) {
	answer1 := 3 | 4
	fmt.Printf("answer1: %v\n", answer1)

	answer2 := 3 & 4
	fmt.Printf("answer2: %v\n", answer2)

	blockBori := []string{
		"Before: [3, 3, 1, 0]",
		"9 1 4 3",
		"After:  [3, 3, 1, 7]",
	}

	blockBorr := []string{
		"Before: [3, 3, 4, 0]",
		"9 1 2 3",
		"After:  [3, 3, 4, 7]",
	}

	blockBani := []string{
		"Before: [3, 3, 1, 9]",
		"9 1 4 3",
		"After:  [3, 3, 1, 0]",
	}

	blockBanr := []string{
		"Before: [3, 3, 4, 9]",
		"9 1 2 3",
		"After:  [3, 3, 4, 0]",
	}

	assert.True(t, IsOpcodeValidForBlock(Bori, blockBori))
	assert.True(t, IsOpcodeValidForBlock(Borr, blockBorr))
	assert.True(t, IsOpcodeValidForBlock(Bani, blockBani))
	assert.True(t, IsOpcodeValidForBlock(Banr, blockBanr))
}

func testInput() []string {
	return []string{
		"Before: [3, 2, 1, 1]",
		"9 2 1 2",
		"After:  [3, 2, 2, 1]",
	}
}
