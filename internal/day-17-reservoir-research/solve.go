package day17reservoirresearch

import (
	"fmt"
	"image/color"
	"log"
	"reflect"
	"slices"
	"strings"

	"github.com/ewoutquax/advent-of-code-2018/pkg/register"
	"github.com/ewoutquax/advent-of-code-2018/pkg/utils"
	"github.com/hajimehoshi/ebiten/v2"
)

const Day string = "17"

type ItemType string

const (
	TypeClay         string = "clay"
	TypeFallingWater string = "falling"
	TypeStillWater   string = "still"
	TypeSlidingWater string = "sliding"
)

var (
	white = color.RGBA{
		R: 255,
		G: 255,
		B: 255,
		A: 255,
	}
	blue = color.RGBA{
		R: 30,
		G: 30,
		B: 255,
		A: 255,
	}
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
	Parent          *FallingWater
	SlidingChildren []*SlidingWater
	IsFalling       bool
}

func (w *FallingWater) hasActiveSlidingChildren() bool {
	for _, child := range w.SlidingChildren {
		if !child.IsSettled && child.Child == nil {
			return true
		}
	}

	return false
}

func (w *FallingWater) hasAllSlidingChildrenSettled() bool {
	for _, child := range w.SlidingChildren {
		if !child.IsSettled {
			return false
		}
	}

	return true
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
	Child         *FallingWater
	IsSlidingLeft bool
	IsSettled     bool
}

type Universe struct {
	CurrentFallingWater *FallingWater
	Items               map[int]interface{}
	WaterDrops          int
	MaxY                int // The highest measured Y value
}

func (u *Universe) AddInitialWater() {
	u.CurrentFallingWater = &FallingWater{
		Type:      ItemType(TypeFallingWater),
		Location:  Location{500, 0},
		Parent:    nil,
		IsFalling: true,
	}
	u.WaterDrops = 0
}

func (u *Universe) Simulate() {
	switch true {
	case u.CurrentFallingWater.IsFalling:
		u.CurrentFallingWater.simulateFalling(u)
	case u.CurrentFallingWater.hasActiveSlidingChildren():
		u.CurrentFallingWater.simulateSliding(u)
	case u.CurrentFallingWater.hasAllSlidingChildrenSettled():
		// Raise the level of the falling water, and spawn new sliding
		panic("Not yet implemented")
	default:
		// dunno... Stop the stream and set its parent as active?
		panic("What to do now")
	}
}

func (w *FallingWater) simulateFalling(u *Universe) {
	var newLoc Location = w.Location
	newLoc.Y += 1

	if newLoc.Y > u.MaxY {
		// The stream is going out-of-reach; delete the active water
		u.CurrentFallingWater = w.Parent
	} else if itemBelow, exists := u.Items[newLoc.toI()]; !exists {
		switch itemBelow.(type) {
		case *Clay:
			// The falling stream hit clay; stop falling, and create 2 sliding streams
			w.IsFalling = false
			u.Items[w.Location.toI()] = &StillWater{
				Type:     ItemType(TypeStillWater),
				Location: w.Location,
			}

			// Convert the falling stream into one or two sliding streams, and link them all together
			w.SlidingChildren = convertFallingIntoSliding(w, u)
		case *StillWater:
			// existing water was hit; this stream should end, and we continue with the parent stream
			w.IsFalling = false
			u.CurrentFallingWater = w.Parent
		default:
			panic("What type of item did we fall on?")
		}
	} else {
		// The stream is still falling
		w.Location = newLoc
		u.WaterDrops += 1
	}
}

func (w *FallingWater) simulateSliding(u *Universe) {
	for _, child := range w.SlidingChildren {
		newLoc := w.Location
		if child.IsSlidingLeft {
			newLoc.X -= 1
		} else {
			newLoc.X += 1
		}

		if u.Items[newLoc.toI()] == nil {
			child.Location = newLoc
			u.WaterDrops += 1
			u.Items[w.Location.toI()] = &StillWater{
				Type:     ItemType(TypeStillWater),
				Location: child.Location,
			}

			locBelow := child.Location
			locBelow.Y += 1
			if u.Items[locBelow.toI()] == nil {
				// The stream has sliden over an edge, and becomes a the new current falling stream
				u.CurrentFallingWater = convertSlidingIntoFalling(child, u)
			}
		} else {
			child.IsSettled = true
		}
	}
}

type Game struct {
	universe *Universe
}

func (g *Game) Update() error {
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	fmt.Println("Start drawing")
	for _, item := range g.universe.Items {
		switch item.(type) {
		case *Clay:
			screen.Set(item.(*Clay).Location.X, item.(*Clay).Location.Y, white)

		case *StillWater:
			screen.Set(item.(*StillWater).Location.X, item.(*StillWater).Location.Y, blue)
		case StillWater:
			screen.Set(item.(StillWater).Location.X, item.(StillWater).Location.Y, blue)
		default:
			fmt.Printf("item: %v\n", item)
			fmt.Printf("reflect.TypeOf(item): %v\n", reflect.TypeOf(item))
			panic("What is this type?")
		}
	}
	fmt.Println("Done drawing")
}

func (g *Game) Layout(_, _ int) (w, h int) {
	return 1000, 2000
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
	if _, exists := u.Items[locLeft.toI()]; !exists {
		newSliding = append(newSliding, &SlidingWater{
			Type:          ItemType(TypeSlidingWater),
			Location:      locLeft,
			Parent:        item,
			IsSlidingLeft: true,
			IsSettled:     false,
		})
		u.Items[locLeft.toI()] = StillWater{
			Type:     ItemType(TypeStillWater),
			Location: locLeft,
		}
	}

	// Try a sliding-stream to the right
	locRight := Location{item.Location.X + 1, item.Location.Y + 1}
	if _, exists := u.Items[locRight.toI()]; !exists {
		newSliding = append(newSliding, &SlidingWater{
			Type:          ItemType(TypeSlidingWater),
			Location:      locRight,
			Parent:        item,
			IsSlidingLeft: false,
			IsSettled:     false,
		})
		u.Items[locRight.toI()] = StillWater{
			Type:     ItemType(TypeStillWater),
			Location: locRight,
		}
	}

	item.SlidingChildren = newSliding
	u.WaterDrops += len(newSliding)

	fmt.Printf("newSliding: %v\n", newSliding)

	return newSliding
}

func SimulateFlow(u *Universe) {
	for u.CurrentFallingWater != nil {
		u.Simulate()
	}
}

func ParseInput(lines []string) Universe {
	var universe = Universe{
		CurrentFallingWater: nil,
		WaterDrops:          0,
		Items:               make(map[int]interface{}),
		MaxY:                0,
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

	ebiten.SetWindowSize(1000*2, 2000*2)
	ebiten.SetWindowTitle("Reservoir Research")
	if err := ebiten.RunGame(&Game{universe: &universe}); err != nil {
		log.Fatal(err)
	}

	go SimulateFlow(&universe)

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
