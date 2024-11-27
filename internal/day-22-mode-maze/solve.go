package day22modemaze

import (
	"container/heap"
	"fmt"
	"math"
	"strconv"
	"strings"

	"github.com/ewoutquax/advent-of-code-2018/pkg/utils"
)

type (
	ErosionLevelType uint
	Equipment        uint
	CaveDepth        int
)

func (e Equipment) toS() string {
	return map[Equipment]string{
		EquipmentClimbingGear: "Climbing gear",
		EquipmentTorch:        "Torch",
		EquipmentNone:         "None",
	}[e]
}

type Location struct {
	X int
	Y int
}

func (l Location) toS() string {
	return fmt.Sprintf("[%d, %d]", l.X, l.Y)
}

const (
	ErosionLevelTypeRocky  ErosionLevelType = 0
	ErosionLevelTypeWet    ErosionLevelType = 1
	ErosionLevelTypeNarrow ErosionLevelType = 2

	EquipmentClimbingGear Equipment = iota + 1
	EquipmentTorch
	EquipmentNone
)

var (
	cachedErosionLevels map[Location]int
	visitedPaths        map[VisitedPathKey]int
)

func (l Location) GeologicalIndex(tl Location, cd CaveDepth) int {
	switch true {
	case l.X == 0 && l.Y == 0:
		return 0
	case l.X == tl.X && l.Y == tl.Y:
		return 0
	case l.Y == 0:
		return l.X * 16807
	case l.X == 0:
		return l.Y * 48271
	default:
		return Location{l.X - 1, l.Y}.ErosionLevel(tl, cd) * Location{l.X, l.Y - 1}.ErosionLevel(tl, cd)
	}
}

func (l Location) ErosionLevel(tl Location, cd CaveDepth) int {
	if cachedErosionLevels == nil {
		cachedErosionLevels = make(map[Location]int, tl.X*tl.Y)
	}

	if level, exists := cachedErosionLevels[l]; exists {
		return level
	}

	newLevel := (l.GeologicalIndex(tl, cd) + int(cd)) % 20183
	cachedErosionLevels[l] = newLevel
	return newLevel
}

func (l Location) ErosionLevelType(tl Location, cd CaveDepth) ErosionLevelType {
	myType := l.ErosionLevel(tl, cd) % 3

	return ErosionLevelType(myType)
}

type Input struct {
	CaveDepth
	TargetLocation Location
}

func FastestTime(tl Location, cd CaveDepth) int {
	InitVisitedPaths()
	var minNrMinutes int = math.MaxInt

	var initPath = Path{
		Location:         Location{0, 0},
		CurrentEquipment: EquipmentTorch,
		NrMinutes:        0,
	}

	pathHeap := PathHeap{initPath}

	visitedPaths[initPath.VisitedKey()] = 0
	heap.Init(&pathHeap)

	for pathHeap.Len() > 0 {
		currentPath := heap.Pop(&pathHeap).(Path)

		if currentPath.Location.X == tl.X && currentPath.Location.Y == tl.Y && currentPath.CurrentEquipment == EquipmentTorch {
			minNrMinutes = currentPath.NrMinutes
			fmt.Printf("Found path to end with torch: %v\n", minNrMinutes)
		}

		nextPaths := currentPath.getNextPaths(tl, cd)
		for _, nextPath := range nextPaths {
			if nextPath.NrMinutes < minNrMinutes && !(nextPath.IsVisited()) {
				nextPath.SetVisited()

				heap.Push(&pathHeap, nextPath)
			}
		}
	}

	return minNrMinutes
}

func ParseInput(lines []string) Input {
	partsDepth := strings.Split(lines[0], ": ")
	partsLocation := strings.Split(lines[1], ": ")
	partsLocations := strings.Split(partsLocation[1], ",")

	return Input{
		CaveDepth: CaveDepth(convAtoI(partsDepth[1])),
		TargetLocation: Location{
			X: convAtoI(partsLocations[0]),
			Y: convAtoI(partsLocations[1]),
		},
	}
}

func allowedEquipments(t ErosionLevelType) []Equipment {
	type Equipments []Equipment

	return map[ErosionLevelType]Equipments{
		ErosionLevelTypeRocky:  {EquipmentClimbingGear, EquipmentTorch},
		ErosionLevelTypeWet:    {EquipmentClimbingGear, EquipmentNone},
		ErosionLevelTypeNarrow: {EquipmentTorch, EquipmentNone},
	}[t]
}

func convAtoI(s string) int {
	nr, err := strconv.Atoi(s)
	if err != nil {
		panic(err)
	}
	return nr
}

func CalculateRisk(tl Location, cd CaveDepth) int {
	var sum int = 0

	for y := 0; y <= tl.Y; y++ {
		for x := 0; x <= tl.X; x++ {
			location := Location{x, y}
			sum += int(location.ErosionLevelType(tl, cd))
		}
	}

	return sum
}

func allTypes() []ErosionLevelType {
	return []ErosionLevelType{
		ErosionLevelTypeRocky,
		ErosionLevelTypeWet,
		ErosionLevelTypeNarrow,
	}
}

func solvePart1(inputFile string) {
	lines := utils.ReadFileAsLines(inputFile)

	input := ParseInput(lines)

	fmt.Printf("Result of day-%s / part-1: %d\n", Day, CalculateRisk(
		input.TargetLocation,
		input.CaveDepth,
	))
}

func solvePart2(inputFile string) {
	lines := utils.ReadFileAsLines(inputFile)
	input := ParseInput(lines)

	fmt.Printf("Result of day-%s / part-2: %d\n", Day, FastestTime(input.TargetLocation, input.CaveDepth))

	/**
	  Answers
	  -------
	  1027: Too low
	*/
}
