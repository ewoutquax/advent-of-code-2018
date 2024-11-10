package day16chronicalclassification

import (
	"fmt"
	"slices"
	"sort"
	"strconv"
	"strings"

	"github.com/ewoutquax/advent-of-code-2018/pkg/register"
	"github.com/ewoutquax/advent-of-code-2018/pkg/utils"
)

type Opcode uint

const (
	Addr Opcode = iota + 1
	Addi
	Muli
	Mulr
	Banr
	Bani
	Borr
	Bori
	Setr
	Seti
	Gtir
	Gtri
	Gtrr
	Eqir
	Eqri
	Eqrr

	MIN_VALID_OPCODES int    = 3
	Day               string = "16"
)

type MappingOpcode map[int]Opcode

type Instruction struct {
	Opcode
	InputA int
	InputB int
	Output int
}

type Registers map[int]int

func (r Registers) ToS() string {
	return fmt.Sprintf("[%d, %d, %d, %d]", r[0], r[1], r[2], r[3])
}

func (i Instruction) Exec(registers Registers) {
	valueTesting := 0

	switch i.Opcode {
	case Addr:
		registers[i.Output] = registers[i.InputA] + registers[i.InputB]
	case Addi:
		registers[i.Output] = registers[i.InputA] + i.InputB
	case Mulr:
		registers[i.Output] = registers[i.InputA] * registers[i.InputB]
	case Muli:
		registers[i.Output] = registers[i.InputA] * i.InputB
	case Banr:
		registers[i.Output] = registers[i.InputA] & registers[i.InputB]
	case Bani:
		registers[i.Output] = registers[i.InputA] & i.InputB
	case Borr:
		registers[i.Output] = registers[i.InputA] | registers[i.InputB]
	case Bori:
		registers[i.Output] = registers[i.InputA] | i.InputB
	case Setr:
		registers[i.Output] = registers[i.InputA]
	case Seti:
		registers[i.Output] = i.InputA
	case Gtir:
		if i.InputA > registers[i.InputB] {
			valueTesting = 1
		}
		registers[i.Output] = valueTesting
	case Gtri:
		if registers[i.InputA] > i.InputB {
			valueTesting = 1
		}
		registers[i.Output] = valueTesting
	case Gtrr:
		if registers[i.InputA] > registers[i.InputB] {
			valueTesting = 1
		}
		registers[i.Output] = valueTesting
	case Eqir:
		if i.InputA == registers[i.InputB] {
			valueTesting = 1
		}
		registers[i.Output] = valueTesting
	case Eqri:
		if registers[i.InputA] == i.InputB {
			valueTesting = 1
		}
		registers[i.Output] = valueTesting
	case Eqrr:
		if registers[i.InputA] == registers[i.InputB] {
			valueTesting = 1
		}
		registers[i.Output] = valueTesting
	default:
		panic("No valid case found")
	}
}

func ValidOpcodes(lines []string) []Opcode {
	var validOpcodes []Opcode = make([]Opcode, 0, len(allOpcodes()))

	for _, currentOpcode := range allOpcodes() {
		if IsOpcodeValidForBlock(currentOpcode, lines) {
			validOpcodes = append(validOpcodes, currentOpcode)
		}
	}

	return validOpcodes
}

func IsOpcodeValidForBlock(opcode Opcode, lines []string) bool {
	var registers Registers = SetRegisters(lines[0])

	instruction := ParseInstruction(lines[1])
	instruction.Opcode = opcode
	instruction.Exec(registers)

	return fmt.Sprintf("After:  %s", registers.ToS()) == lines[2]
}

func SetRegisters(line string) Registers {
	var registers Registers = make(Registers, 4)

	suffix := strings.Split(line[:len(line)-1], "[")[1]
	parts := strings.Split(suffix, ", ")

	registers[0] = convAtoI(parts[0])
	registers[1] = convAtoI(parts[1])
	registers[2] = convAtoI(parts[2])
	registers[3] = convAtoI(parts[3])

	return registers
}

func ParseInstruction(line string) Instruction {
	parts := strings.Split(line, " ")

	return Instruction{
		Opcode: Opcode(convAtoI(parts[0])),
		InputA: convAtoI(parts[1]),
		InputB: convAtoI(parts[2]),
		Output: convAtoI(parts[3]),
	}
}

func convAtoI(s string) int {
	nr, err := strconv.Atoi(s)
	if err != nil {
		panic(err)
	}
	return nr
}

func allOpcodes() []Opcode {
	return []Opcode{
		Addr, Addi,
		Muli, Mulr,
		Banr, Bani,
		Borr, Bori,
		Setr, Seti,
		Gtir, Gtri, Gtrr,
		Eqir, Eqri, Eqrr,
	}
}

func solvePart1(inputFile string) {
	blocks := utils.ReadFileAsBlocks(inputFile)

	var count int = 0

	for idx := 0; idx < len(blocks)-2; idx++ {
		validOpcodes := ValidOpcodes(blocks[idx])
		if len(validOpcodes) >= MIN_VALID_OPCODES {
			count++
		}
	}

	fmt.Printf("Result of day-%s / part-1: %d\n", Day, count)
}

func solvePart2(inputFile string) {
	blocks := utils.ReadFileAsBlocks(inputFile)

	mappedOpcodes := mapOpcodes(blocks)

	var registers Registers = Registers{
		0: 0,
		1: 0,
		2: 0,
		3: 0,
	}

	for _, rawInstruction := range blocks[len(blocks)-1] {
		instruction := ParseInstruction(rawInstruction)
		instruction.Opcode = mappedOpcodes[int(instruction.Opcode)]

		instruction.Exec(registers)
	}

	fmt.Printf("Result of day-%s / part-2: %d\n", Day, registers[0])
}

func mapOpcodes(blocks [][]string) MappingOpcode {
	type OpcodeToInts map[Opcode][]int

	var opcodeToInts OpcodeToInts = make(OpcodeToInts)

	// initialize
	for idx := 0; idx <= 15; idx++ {
		opcodeToInts[Opcode(idx+1)] = make([]int, 0, 16)
	}

	// find valid Opcode by opcode-id
	for idx := 0; idx < len(blocks)-2; idx++ {
		instruction := ParseInstruction(blocks[idx][1])

		opcodes := ValidOpcodes(blocks[idx])

		// store possible opcode-ids by Opcode
		for _, opcode := range opcodes {
			opcodeToInts[opcode] = append(opcodeToInts[opcode], int(instruction.Opcode))
		}
	}

	// remove duplicate opcode-ids per opcode, and sort them for convinience
	for opcode, ints := range opcodeToInts {
		uniqInts := uniq(ints)
		sort.Ints(uniqInts)
		opcodeToInts[opcode] = uniqInts
	}

	var mappedOpcodes map[int]Opcode = make(map[int]Opcode, 16)

	// When an opcode is linked to just 1 opcode-id, then copy the link to the output and remove the opcode-id from all opcodes
	for len(mappedOpcodes) < 16 {
		for opcode, ints := range opcodeToInts {
			if len(ints) == 1 {
				mappedOpcodes[ints[0]] = opcode

				for subOpcode, subInts := range opcodeToInts {
					var cleanedInts []int = make([]int, 0)
					if len(subInts) > 0 {
						cleanedInts = slices.DeleteFunc(subInts, func(subInt int) bool {
							return subInt == ints[0]
						})
					}

					opcodeToInts[subOpcode] = cleanedInts
				}
			}
		}
	}

	return mappedOpcodes
}

func uniq(ints []int) []int {
	var uniqInts []int = make([]int, 0, 15)
	var unique map[int]bool = make(map[int]bool, 15)

	for _, currentInt := range ints {
		unique[currentInt] = true
	}

	for currentInt := range unique {
		uniqInts = append(uniqInts, currentInt)
	}

	return uniqInts
}

func init() {
	register.Day(Day+"b", solvePart2)
	register.Day(Day+"a", solvePart1)
}
