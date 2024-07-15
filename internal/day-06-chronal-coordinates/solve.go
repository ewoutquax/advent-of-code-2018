package day06chronalcoordinates

import (
	"fmt"
	"math"
	"regexp"

	"github.com/ewoutquax/advent-of-code-2018/pkg/register"
	"github.com/ewoutquax/advent-of-code-2018/pkg/utils"
)

type Location struct {
	X int
	Y int
}

func (l Location) ManhattanDistance(other Location) int {
	return utils.Abs(l.X-other.X) + utils.Abs(l.Y-other.Y)
}
func (l Location) toS() string {
	return fmt.Sprintf("%d,%d", l.X, l.Y)
}

type Area struct {
	Location
	Size       int
	IsInfinite bool
}

type Universe struct {
	Areas []*Area
	MinX  int
	MinY  int
	MaxX  int
	MaxY  int
}

const Day string = "06"
const INFINITE int = math.MaxInt

func init() {
	register.Day(Day+"a", solvePart1)
	// register.Day(Day+"b", solvePart2)
}

func solvePart1(inputFile string) {
	lines := utils.ReadFileAsLines(inputFile)
	universe := ParseInput(lines)
	var count int = MaxFiniteSize(universe)

	fmt.Printf("Result of day-%s / part-1: %d\n", Day, count)
}

func solvePart2(inputFile string) {
	fmt.Printf("Result of day-%s / part-2: %d\n", Day, 0)
}

func MaxFiniteSize(universe Universe) int {
	var maxSize int = 0

	for y := universe.MinY - 1; y <= universe.MaxY+1; y++ {
		for x := universe.MinX - 1; x <= universe.MaxX+1; x++ {
			currentLocation := Location{X: x, Y: y}
			var closestArea *Area
			minDistance := INFINITE
			var unique bool = true

			for _, area := range universe.Areas {
				distance := currentLocation.ManhattanDistance(area.Location)
				if minDistance == distance {
					unique = false
				}
				if minDistance > distance {
					unique = true
					minDistance = distance
					closestArea = area
				}
			}

			if unique {
				closestArea.Size++
				if currentLocation.X < universe.MinX ||
					currentLocation.X > universe.MaxX ||
					currentLocation.Y < universe.MinY ||
					currentLocation.Y > universe.MaxY {
					closestArea.IsInfinite = true
				}
				if maxSize < closestArea.Size && !closestArea.IsInfinite {
					maxSize = closestArea.Size
				}
			}
		}
	}

	return maxSize
}

func ParseInput(lines []string) Universe {
	var areas []*Area = make([]*Area, 0, len(lines))
	var minX, minY int = INFINITE, INFINITE
	var maxX, maxY int = 0, 0

	for _, line := range lines {
		ex := regexp.MustCompile(`(-?\d+)+`)
		matches := ex.FindAllString(line, -1)

		area := Area{
			Location: Location{
				X: utils.ConvStrToI(matches[0]),
				Y: utils.ConvStrToI(matches[1]),
			},
			IsInfinite: false,
			Size:       0,
		}
		areas = append(areas, &area)

		if minX > area.X {
			minX = area.X
		}
		if minY > area.Y {
			minY = area.Y
		}
		if maxX < area.X {
			maxX = area.X
		}
		if maxY < area.Y {
			maxY = area.Y
		}
	}

	return Universe{
		Areas: areas,
		MinX:  minX,
		MinY:  minY,
		MaxX:  maxX,
		MaxY:  maxY,
	}
}
