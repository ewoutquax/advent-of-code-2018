package day17reservoirresearch

import (
	"fmt"
	"slices"
	"strings"

	"github.com/ewoutquax/advent-of-code-2018/pkg/register"
	"github.com/ewoutquax/advent-of-code-2018/pkg/utils"
)

const Day string = "17"

type ItemType string

const (
	TypeClay         string = "clay"
	TypeFallingWater string = "falling"
	TypeStillWater   string = "still"
	TypeSlidingWater string = "sliding"
)

type Location struct {
	X int
	Y int
}

func (l Location) toI() int {
	return (l.Y * 2000) + l.X
}

type Clay struct {
	Type ItemType
	Location
}

type StillWater struct {
	Type ItemType
	Location
}

type FallingWater struct {
	Type ItemType
	Location
	SlidingChildren []*SlidingWater
}

func (w *FallingWater) allSiblingsSettled() bool {
	for _, s := range w.SlidingChildren {
		if !s.IsSettled {
			return false
		}
	}
	return true
}

type SlidingWater struct {
	Type ItemType
	Location
	Parent        *FallingWater
	IsSlidingLeft bool
	IsSettled     bool
	IsFalling     bool
}

type Universe struct {
	ActiveWaterItems []interface{}
	Items            map[int]interface{}
	WaterDrops       int
	MaxY             int // The highest measured Y value
}

func (u *Universe) AddInitialWater() {
	u.ActiveWaterItems = append(u.ActiveWaterItems, &FallingWater{
		Type:     ItemType(TypeFallingWater),
		Location: Location{500, 0},
	})
	u.WaterDrops = 0
}

func (u *Universe) Simulate() {
	for _, item := range u.ActiveWaterItems {
		switch item.(type) {
		case *FallingWater:
			u.simulateFallingWater(item.(*FallingWater))
		case *SlidingWater:
			u.SimulateSlidingWater(item.(*SlidingWater))
		}
	}
}

func (u *Universe) simulateFallingWater(item *FallingWater) {
	var newLoc Location = item.Location
	newLoc.Y += 1

	if newLoc.Y > u.MaxY {
		// The stream is going out-of-reach; delete the active water
		u.ActiveWaterItems = slices.DeleteFunc(u.ActiveWaterItems, func(waterItem interface{}) bool {
			return waterItem == item
		})
	} else if u.Items[newLoc.toI()] == nil {
		item.Location = newLoc
		u.WaterDrops += 1
	} else {
		u.Items[item.Location.toI()] = &StillWater{
			Type:     ItemType(TypeStillWater),
			Location: item.Location,
		}

		// Delete the now-ended stream from the active water items
		u.ActiveWaterItems = slices.DeleteFunc(u.ActiveWaterItems, func(waterItem interface{}) bool {
			return waterItem == item
		})

		// Convert the falling stream into one or two sliding streams, and link them all together
		for _, sliding := range convertFallingIntoSliding(item, u) {
			u.ActiveWaterItems = append(u.ActiveWaterItems, sliding)
		}
	}
}

func (u *Universe) SimulateSlidingWater(item *SlidingWater) {
	newLoc := item.Location
	if item.IsSlidingLeft {
		newLoc.X -= 1
	} else {
		newLoc.X += 1
	}

	if u.Items[newLoc.toI()] == nil {
		item.Location = newLoc
		u.WaterDrops += 1
		u.Items[item.Location.toI()] = &StillWater{
			Type:     ItemType(TypeStillWater),
			Location: item.Location,
		}

		locBelow := item.Location
		locBelow.Y += 1
		if u.Items[locBelow.toI()] == nil {
			u.ActiveWaterItems = slices.DeleteFunc(u.ActiveWaterItems, func(waterItem interface{}) bool {
				return waterItem == item
			})
			u.ActiveWaterItems = append(u.ActiveWaterItems, convertSlidingIntoFalling(item, u))
		}

	} else {
		item.IsSettled = true
		u.ActiveWaterItems = slices.DeleteFunc(u.ActiveWaterItems, func(waterItem interface{}) bool {
			return waterItem == item
		})

		if item.Parent.allSiblingsSettled() {
			for _, sliding := range convertFallingIntoSliding(item.Parent, u) {
				u.ActiveWaterItems = append(u.ActiveWaterItems, sliding)
			}
		}
	}
}

func convertSlidingIntoFalling(item *SlidingWater, u *Universe) *FallingWater {
	newLoc := item.Location
	newLoc.Y += 1

	u.WaterDrops += 1
	return &FallingWater{
		Type:     ItemType(TypeFallingWater),
		Location: newLoc,
	}
}

func convertFallingIntoSliding(item *FallingWater, u *Universe) []*SlidingWater {
	newSliding := make([]*SlidingWater, 0, 2)

	// Move the falling stream one position up, for reuse
	newLoc := item.Location
	newLoc.Y -= 1
	item.Location = newLoc

	// Try a sliding-stream to the left
	locLeft := Location{item.Location.X - 1, item.Location.Y + 1}
	fmt.Printf("convertFallingIntoSliding: locLeft: %v\n", locLeft)
	if elem, exists := u.Items[locLeft.toI()]; !exists {
		fmt.Println("Adding sliding water to the left")
		newSliding = append(newSliding, &SlidingWater{
			Type:          ItemType(TypeSlidingWater),
			Location:      locLeft,
			Parent:        item,
			IsSlidingLeft: true,
		})
		u.Items[locLeft.toI()] = StillWater{
			Type:     ItemType(TypeStillWater),
			Location: locLeft,
		}
	} else {
		fmt.Printf("blocking element for sliding left: %v\n", elem)
	}

	// Try a sliding-stream to the right
	locRight := Location{item.Location.X + 1, item.Location.Y + 1}
	fmt.Printf("convertFallingIntoSliding: locRight: %v\n", locRight)
	if elem, exists := u.Items[locRight.toI()]; !exists {
		fmt.Println("Adding sliding water to the right")
		newSliding = append(newSliding, &SlidingWater{
			Type:          ItemType(TypeSlidingWater),
			Location:      locRight,
			Parent:        item,
			IsSlidingLeft: false,
		})
		u.Items[locRight.toI()] = StillWater{
			Type:     ItemType(TypeStillWater),
			Location: locRight,
		}
	} else {
		fmt.Printf("blocking element for sliding right: %v\n", elem)
	}

	item.SlidingChildren = newSliding
	u.WaterDrops += len(newSliding)

	fmt.Printf("newSliding: %v\n", newSliding)

	return newSliding
}

func SimulateFlow(u *Universe) {
	for len(u.ActiveWaterItems) > 0 {
		u.Simulate()
	}
}

func ParseInput(lines []string) Universe {
	var universe = Universe{
		ActiveWaterItems: make([]interface{}, 0),
		WaterDrops:       0,
		Items:            make(map[int]interface{}),
		MaxY:             0,
	}

	for _, line := range lines {
		var xPart, yPart string
		var xFrom, xTo, yFrom, yTo int

		parts := strings.Split(line, ", ")
		if strings.Index(parts[0], "x") == 0 {
			xPart = parts[0]
			yPart = parts[1]
		} else {
			xPart = parts[1]
			yPart = parts[0]
		}

		xFrom, xTo = extractNumbers(strings.Split(xPart, "=")[1])
		yFrom, yTo = extractNumbers(strings.Split(yPart, "=")[1])

		universe.MaxY = max(universe.MaxY, yTo)

		for x := xFrom; x <= xTo; x++ {
			for y := yFrom; y <= yTo; y++ {
				clay := Clay{
					Type:     ItemType(TypeClay),
					Location: Location{x, y},
				}
				universe.Items[clay.toI()] = &clay
			}
		}
	}

	return universe
}

func extractNumbers(input string) (int, int) {
	if strings.Contains(input, "..") {
		parts := strings.Split(input, "..")
		return utils.ConvStrToI(parts[0]), utils.ConvStrToI(parts[1])
	} else {
		return utils.ConvStrToI(input), utils.ConvStrToI(input)
	}
}

func max(a, b int) int {
	if a < b {
		return b
	} else {
		return a
	}
}

func solvePart1(inputFile string) {
	lines := utils.ReadFileAsLines(inputFile)
	universe := ParseInput(lines)
	universe.AddInitialWater()

	SimulateFlow(&universe)

	// Too low:
	// --------
	// 2032

	fmt.Printf("Result of day-%s / part-1: %d\n", Day, universe.WaterDrops)
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
