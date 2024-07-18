package day07thesumofitsparts

import (
	"fmt"
	"regexp"
	"sort"
	"strings"

	"github.com/ewoutquax/advent-of-code-2018/pkg/register"
	"github.com/ewoutquax/advent-of-code-2018/pkg/utils"
)

const Day string = "07"

type Item struct {
	Name               string
	BlockedBy          []*Item
	RemainingBuildTime int
	IsInProgress       bool
}

func (item *Item) CanBeStarted() bool {
	if item.IsInProgress || item.RemainingBuildTime == 0 {
		return false
	}

	for _, blockingItem := range item.BlockedBy {
		if blockingItem.RemainingBuildTime > 0 {
			return false
		}
	}

	return true
}

type Order struct {
	Before *Item
	After  *Item
}

type Universe struct {
	Orders []Order
	Items  map[string]*Item
}

func (u Universe) AllItemsBuild() bool {
	for _, item := range u.Items {
		if item.RemainingBuildTime > 0 {
			return false
		}
	}

	return true
}

func BuildMetrics(universe Universe, nrWorkers int) (string, int) {
	var order []string = make([]string, 0, len(universe.Items))
	var elapsedTime int = 0
	var nrAvailableWorkers int = nrWorkers

	for !universe.AllItemsBuild() {
		startableItems := getStartableItems(universe)
		for len(startableItems) > 0 && nrAvailableWorkers > 0 {
			universe.Items[startableItems[0]].IsInProgress = true
			order = append(order, startableItems[0])

			nrAvailableWorkers--
			startableItems = getStartableItems(universe)
		}

		elapsedTime++

		for _, item := range universe.Items {
			if item.IsInProgress {
				item.RemainingBuildTime--
				if item.RemainingBuildTime == 0 {
					item.IsInProgress = false
					nrAvailableWorkers++
				}
			}
		}
	}

	return strings.Join(order, ""), elapsedTime
}

func getStartableItems(universe Universe) []string {
	items := make([]string, 0, len(universe.Items))
	for _, item := range universe.Items {
		if item.CanBeStarted() {
			items = append(items, item.Name)
		}
	}
	sort.Strings(items)

	return items
}

func ParseInput(lines []string, offsetBuildTime int) Universe {
	var u Universe = Universe{
		Orders: make([]Order, 0, len(lines)),
		Items:  make(map[string]*Item, len(lines)*2),
	}

	for _, line := range lines {
		ex := regexp.MustCompile(`Step ([A-Z]) must be finished before step ([A-Z]) can begin.`)
		match := ex.FindStringSubmatch(line)

		itemBefore, existsBefore := u.Items[match[1]]
		if !existsBefore {
			itemBefore = &Item{
				Name:               match[1],
				RemainingBuildTime: offsetBuildTime + 1 + int(match[1][0]-'A'),
			}
		}

		itemAfter, existsAfter := u.Items[match[2]]
		if !existsAfter {
			itemAfter = &Item{
				Name:               match[2],
				RemainingBuildTime: offsetBuildTime + 1 + int(match[2][0]-'A'),
			}
		}

		u.Items[itemBefore.Name] = itemBefore
		u.Items[itemAfter.Name] = itemAfter
		u.Items[itemAfter.Name].BlockedBy = append(u.Items[itemAfter.Name].BlockedBy, itemBefore)

		u.Orders = append(u.Orders, Order{
			Before: u.Items[itemBefore.Name],
			After:  u.Items[itemAfter.Name],
		})
	}

	return u
}

func solvePart1(inputFile string) {
	lines := utils.ReadFileAsLines(inputFile)
	universe := ParseInput(lines, 0)
	order, _ := BuildMetrics(universe, 1)

	fmt.Printf("Result of day-%s / part-1: %s\n", Day, order)
}

func solvePart2(inputFile string) {
	lines := utils.ReadFileAsLines(inputFile)
	universe := ParseInput(lines, 60)
	_, elapsedTime := BuildMetrics(universe, 5)

	fmt.Printf("Result of day-%s / part-2: %d\n", Day, elapsedTime)
}

func init() {
	register.Day(Day+"a", solvePart1)
	register.Day(Day+"b", solvePart2)
}
