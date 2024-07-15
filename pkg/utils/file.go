package utils

import (
	"os"
	"strings"
)

func ReadFileAsNumbers(baseDir string) (numbers []int) {
	var lines []string = ReadFileAsLines(baseDir)

	for _, string := range lines {
		numbers = append(numbers, ConvStrToI(string))
	}

	return
}

func ReadFileAsBlocks(baseDir string) (blocks [][]string) {
	var block_inputs []string = strings.Split(readFile(baseDir), "\n\n")

	for _, block_input := range block_inputs {
		blocks = append(blocks, strings.Split(block_input, "\n"))
	}
	return
}

func ReadFileAsLines(inputFile string) []string {
	return strings.Split(readFile(inputFile), "\n")
}

func ReadFileAsLine(inputFile string) string {
	return readFile(inputFile)
}

func readFile(inputFile string) string {
	raw, err := os.ReadFile(inputFile)
	check(err)

	return strings.TrimSuffix(string(raw), "\n")
}
