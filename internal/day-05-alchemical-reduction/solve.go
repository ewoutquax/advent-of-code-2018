package day05alchemicalreduction

import (
	"fmt"
	"strings"

	"github.com/ewoutquax/advent-of-code-2018/pkg/register"
	"github.com/ewoutquax/advent-of-code-2018/pkg/utils"
)

const Day string = "05"
const ALPHABET = "abcdefghijklmnopqrstuvwxyz"

func ShortestPolymerLengthWithExtraction(polymer string) int {
	return len(ShortestPolymerWithExtraction(polymer))
}

func PolymerLengthAfterTrigger(polymer string) int {
	return len(TriggerPolymer(polymer))
}

func ShortestPolymerWithExtraction(polymer string) string {
	var shortestPolymer string = polymer

	for _, extractable := range strings.Split(ALPHABET, "") {
		result := triggerPolymerWithExtraction(polymer, extractable)
		if len(shortestPolymer) > len(result) {
			shortestPolymer = result
		}
	}

	return shortestPolymer
}

func triggerPolymerWithExtraction(polymer, extractable string) string {
	polymer = strings.Replace(polymer, extractable, "", -1)
	polymer = strings.Replace(polymer, strings.ToUpper(extractable), "", -1)

	return TriggerPolymer(polymer)
}

func TriggerPolymer(polymer string) string {
	var idx int = 0
	var doAnotherRun bool = true

	for doAnotherRun {
		idx = 0
		doAnotherRun = false
		for idx < len(polymer)-1 {
			if DoesReact(polymer[idx : idx+2]) {
				polymer = polymer[:idx] + polymer[idx+2:]
				doAnotherRun = true
			} else {
				idx++
			}
		}
	}

	return polymer
}

func DoesReact(input string) bool {
	parts := strings.Split(input, "")

	return parts[0] != parts[1] && strings.ToUpper(parts[0]) == strings.ToUpper(parts[1])
}

func solvePart1(inputFile string) {
	line := utils.ReadFileAsLine(inputFile)
	polymer := TriggerPolymer(line)
	fmt.Printf("Result of day-%s / part-1: %d\n", Day, len(polymer))
}

func solvePart2(inputFile string) {
	line := utils.ReadFileAsLine(inputFile)
	polymer := ShortestPolymerWithExtraction(line)
	fmt.Printf("Result of day-%s / part-2: %d\n", Day, len(polymer))
}

func init() {
	register.Day(Day+"a", solvePart1)
	register.Day(Day+"b", solvePart2)
}
