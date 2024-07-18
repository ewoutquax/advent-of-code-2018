package day07thesumofitsparts

import (
	"fmt"
	"regexp"
	"sort"

	"github.com/ewoutquax/advent-of-code-2018/pkg/register"
	"github.com/ewoutquax/advent-of-code-2018/pkg/utils"
)

const Day string = "07"

type Item string

type Order struct {
	Before Item
	After  Item
}

type Universe struct {
	Orders []Order
}

func BuildOrder(universe Universe) string {
	var tmpSorting map[Item]bool = make(map[Item]bool, len(universe.Orders)*2)

	for _, order := range universe.Orders {
		tmpSorting[order.Before] = true
		tmpSorting[order.After] = true
	}

	// Copy the unique items to a list
	var sorting = make([]Item, 0, len(tmpSorting))
	for item := range tmpSorting {
		sorting = append(sorting, item)
	}

	// Sort the list in alphabetical order
	sort.Slice(sorting, func(i, j int) bool {
		return sorting[i] < sorting[j]
	})

	// Loop through the orders, and swap the positions of the items
	for _, order := range universe.Orders {
		idxBefore := indexOf(sorting, order.Before)
		idxAfter := indexOf(sorting, order.After)

		if idxBefore > idxAfter {
			start := sorting[:idxAfter]
			middle := sorting[idxAfter+1 : idxBefore]
			end := sorting[idxBefore+1:]
			before := sorting[idxBefore]
			after := sorting[idxAfter]

			sorting = append(start, before)
			sorting = append(sorting, middle...)
			sorting = append(sorting, after)
			sorting = append(sorting, end...)
		}
	}

	// Copy the sorted items to a string
	var out string = ""
	for _, item := range sorting {
		out += string(item)
	}

	return out
}

func ParseInput(lines []string) Universe {
	var u Universe

	for _, line := range lines {
		ex := regexp.MustCompile(`Step ([A-Z]) must be finished before step ([A-Z]) can begin.`)
		match := ex.FindStringSubmatch(line)

		u.Orders = append(u.Orders, Order{
			Before: Item(match[1]),
			After:  Item(match[2]),
		})
	}

	return u
}

func indexOf(haystack []Item, needle Item) int {
	for idx, item := range haystack {
		if item == needle {
			return idx
		}
	}
	return -1
}

func solvePart1(inputFile string) {
	lines := utils.ReadFileAsLines(inputFile)
	universe := ParseInput(lines)

	fmt.Printf("Result of day-%s / part-1: %s\n", Day, BuildOrder(universe))
}

func solvePart2(inputFile string) {
	// lines := utils.ReadFileAsLines(inputFile)

	var count int = 0
	fmt.Printf("Result of day-%s / part-2: %d\n", Day, count)
}

func init() {
	register.Day(Day+"a", solvePart1)
	// register.Day(Day+"b", solvePart2)
}
