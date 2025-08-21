package day13minecartmadness

import (
	"container/heap"
	"fmt"
	"image"
	"strings"

	"github.com/ewoutquax/advent-of-code-2018/pkg/register"
	"github.com/ewoutquax/advent-of-code-2018/pkg/utils"
)

type (
	Turn      int
	Direction uint
	Location  = image.Point
	Cart      struct {
		Location Location
		NrMoves  int
		NextTurn Turn
		Direction
	}
	Rail struct {
		Location
		char       string
		Neighbours map[Direction]*Rail
	}
	Track struct {
		Carts map[Location]Cart
		Rails map[Location]*Rail
	}
	CartHeap []Cart
)

func (h CartHeap) Len() int { return len(h) }
func (h CartHeap) Less(i, j int) bool {
	return (h[i].NrMoves) < h[j].NrMoves ||
		(h[i].NrMoves == h[j].NrMoves && h[i].Location.Y < h[j].Location.Y) ||
		(h[i].NrMoves == h[j].NrMoves && h[i].Location.Y == h[j].Location.Y && h[i].Location.X < h[j].Location.X)
}
func (h CartHeap) Swap(i, j int) { h[i], h[j] = h[j], h[i] }
func (h *CartHeap) Push(x any)   { *h = append(*h, x.(Cart)) }
func (h *CartHeap) Pop() any {
	old := *h
	n := len(old)
	x := old[n-1]
	*h = old[0 : n-1]
	return x
}

const (
	DirectionUp Direction = iota
	DirectionRight
	DirectionDown
	DirectionLeft
)
const (
	Day string = "13"

	TurnLeft    Turn = -1
	TurnForward Turn = 0
	TurnRight   Turn = 1
)

func (c *Cart) Move(track Track) {
	c.Location = track.Rails[c.Location].Neighbours[c.Direction].Location
	c.NrMoves++

	// Do we need to turn, because we're in a corner?
	if _, ok := track.Rails[c.Location].Neighbours[c.Direction]; !ok {
		_, okDirectionUp := track.Rails[c.Location].Neighbours[DirectionUp]
		_, okDirectionRight := track.Rails[c.Location].Neighbours[DirectionRight]
		_, okDirectionDown := track.Rails[c.Location].Neighbours[DirectionDown]
		_, okDirectionLeft := track.Rails[c.Location].Neighbours[DirectionLeft]

		switch true {
		case okDirectionUp && (c.Direction == DirectionLeft || c.Direction == DirectionRight):
			c.Direction = DirectionUp
		case okDirectionRight && (c.Direction == DirectionUp || c.Direction == DirectionDown):
			c.Direction = DirectionRight
		case okDirectionDown && (c.Direction == DirectionLeft || c.Direction == DirectionRight):
			c.Direction = DirectionDown
		case okDirectionLeft && (c.Direction == DirectionUp || c.Direction == DirectionDown):
			c.Direction = DirectionLeft
		default:
			panic("No valid case found")
		}
	}

	// Do we need to turn, because we're on an intersection?
	if len(track.Rails[c.Location].Neighbours) == 4 {
		c.Direction = (c.Direction + Direction(c.NextTurn) + 4) % 4
		c.NextTurn = getNextTurn(c.NextTurn)
	}
}

func FindLastCrashLocation(track Track) string {
	var cartHeap = make(CartHeap, 0, len(track.Carts))                // Priority queue that will pop the cart to move
	var cartLocations = make(map[Location]struct{}, len(track.Carts)) // Indexed list of all locations with carts, to detect crashes
	var maxNrMoves int                                                // All carts should do the same number of moves

	// Load the carts into more usable structs
	for _, cart := range track.Carts {
		cartHeap = append(cartHeap, cart)
		cartLocations[cart.Location] = struct{}{}
	}

	heap.Init(&cartHeap)
	for len(cartHeap) > 1 {
		currentCart := heap.Pop(&cartHeap).(Cart)

		if _, ok := cartLocations[currentCart.Location]; ok {
			// There is a cart on this location, which will not happen after a crash
			delete(cartLocations, currentCart.Location)
			currentCart.Move(track)

			if _, ok := cartLocations[currentCart.Location]; ok {
				// There is already a cart on this location: we found a crash!
				delete(cartLocations, currentCart.Location)
			} else {
				heap.Push(&cartHeap, currentCart)
				cartLocations[currentCart.Location] = struct{}{}
			}
		}

		if maxNrMoves < currentCart.NrMoves {
			maxNrMoves = currentCart.NrMoves
		}
	}

	lastCart := cartHeap[0]
	if lastCart.NrMoves < maxNrMoves {
		lastCart.Move(track)
	}

	return fmt.Sprintf("%d,%d", lastCart.Location.X, lastCart.Location.Y)
}

func FindFirstCrashLocation(track Track) string {
	var cartHeap = make(CartHeap, 0, len(track.Carts))                // Priority queue that will pop the cart to move
	var cartLocations = make(map[Location]struct{}, len(track.Carts)) // Indexed list of all locations with carts, to detect crashes
	var crashLocation image.Point                                     // The solution for part-1

	// Load the carts into more usable structs
	for _, cart := range track.Carts {
		cartHeap = append(cartHeap, cart)
		cartLocations[cart.Location] = struct{}{}
	}

	heap.Init(&cartHeap)
	for len(cartHeap) == len(track.Carts) {
		currentCart := heap.Pop(&cartHeap).(Cart)

		delete(cartLocations, currentCart.Location)
		currentCart.Move(track)

		if _, ok := cartLocations[currentCart.Location]; ok {
			// There is already a cart on this location: we found a crash!
			crashLocation = currentCart.Location
		} else {
			heap.Push(&cartHeap, currentCart)
			cartLocations[currentCart.Location] = struct{}{}
		}
	}

	return fmt.Sprintf("%d,%d", crashLocation.X, crashLocation.Y)
}

func ParseInput(lines []string) Track {
	var rails = make(map[Location]*Rail)

	// First, create all the rails
	for y, line := range lines {
		for x, char := range strings.Split(line, "") {
			if char != " " {
				rails[Location(image.Pt(x, y))] = &Rail{
					Location:   Location(image.Pt(x, y)),
					Neighbours: make(map[Direction]*Rail, 0),
					char:       char,
				}
			}
		}
	}

	// Now, link all the rail structs
	for y, line := range lines {
		for x, char := range strings.Split(line, "") {
			currentRail := rails[Location(image.Pt(x, y))]

			switch char {
			case "-", ">", "<":
				currentRail.Neighbours[DirectionLeft] = rails[currentRail.Location.Add(toVector(DirectionLeft))]
				currentRail.Neighbours[DirectionRight] = rails[currentRail.Location.Add(toVector(DirectionRight))]
			case "|", "^", "v":
				currentRail.Neighbours[DirectionUp] = rails[currentRail.Location.Add(toVector(DirectionUp))]
				currentRail.Neighbours[DirectionDown] = rails[currentRail.Location.Add(toVector(DirectionDown))]
			case "+":
				currentRail.Neighbours[DirectionUp] = rails[currentRail.Location.Add(toVector(DirectionUp))]
				currentRail.Neighbours[DirectionRight] = rails[currentRail.Location.Add(toVector(DirectionRight))]
				currentRail.Neighbours[DirectionDown] = rails[currentRail.Location.Add(toVector(DirectionDown))]
				currentRail.Neighbours[DirectionLeft] = rails[currentRail.Location.Add(toVector(DirectionLeft))]
			case "/":
				locLeft := currentRail.Location.Add(image.Pt(-1, 0))
				if railsLeft, ok := rails[locLeft]; ok && (railsLeft.char == "-" || railsLeft.char == "+" || railsLeft.char == "<" || railsLeft.char == ">") {
					currentRail.Neighbours[DirectionLeft] = rails[currentRail.Location.Add(toVector(DirectionLeft))]
					currentRail.Neighbours[DirectionUp] = rails[currentRail.Location.Add(toVector(DirectionUp))]
				} else {
					currentRail.Neighbours[DirectionRight] = rails[currentRail.Location.Add(toVector(DirectionRight))]
					currentRail.Neighbours[DirectionDown] = rails[currentRail.Location.Add(toVector(DirectionDown))]
				}
			case `\`:
				locLeft := currentRail.Location.Add(image.Pt(-1, 0))
				if railsLeft, ok := rails[locLeft]; ok && (railsLeft.char == "-" || railsLeft.char == "+" || railsLeft.char == "<" || railsLeft.char == ">") {
					currentRail.Neighbours[DirectionLeft] = rails[currentRail.Location.Add(toVector(DirectionLeft))]
					currentRail.Neighbours[DirectionDown] = rails[currentRail.Location.Add(toVector(DirectionDown))]
				} else {
					currentRail.Neighbours[DirectionRight] = rails[currentRail.Location.Add(toVector(DirectionRight))]
					currentRail.Neighbours[DirectionUp] = rails[currentRail.Location.Add(toVector(DirectionUp))]
				}
			case " ":
			default:
				panic(fmt.Sprintf("No valid case found: '%v'", char))
			}
		}
	}

	// Lastly, build the carts
	var carts = make(map[Location]Cart)
	for y, line := range lines {
		for x, char := range strings.Split(line, "") {
			switch char {
			case "^":
				carts[Location(image.Pt(x, y))] = Cart{
					Location:  Location(image.Pt(x, y)),
					NrMoves:   0,
					Direction: DirectionUp,
					NextTurn:  TurnLeft,
				}
			case ">":
				carts[Location(image.Pt(x, y))] = Cart{
					Location:  Location(image.Pt(x, y)),
					NrMoves:   0,
					Direction: DirectionRight,
					NextTurn:  TurnLeft,
				}
			case "v":
				carts[Location(image.Pt(x, y))] = Cart{
					Location:  Location(image.Pt(x, y)),
					NrMoves:   0,
					Direction: DirectionDown,
					NextTurn:  TurnLeft,
				}
			case "<":
				carts[Location(image.Pt(x, y))] = Cart{
					Location:  Location(image.Pt(x, y)),
					NrMoves:   0,
					Direction: DirectionLeft,
					NextTurn:  TurnLeft,
				}
			default:
			}
		}
	}

	return Track{
		Carts: carts,
		Rails: rails,
	}
}

func getNextTurn(t Turn) Turn {
	return map[Turn]Turn{
		TurnLeft:    TurnForward,
		TurnForward: TurnRight,
		TurnRight:   TurnLeft,
	}[t]
}

func toVector(d Direction) image.Point {
	return map[Direction]image.Point{
		DirectionUp:    image.Pt(0, -1),
		DirectionRight: image.Pt(1, 0),
		DirectionDown:  image.Pt(0, 1),
		DirectionLeft:  image.Pt(-1, 0),
	}[d]
}

func solvePart1(inputFile string) {
	lines := utils.ReadFileAsLines(inputFile)
	track := ParseInput(lines)

	fmt.Printf("Result of day-%s / part-1: %s\n", Day, FindFirstCrashLocation(track))
}

func solvePart2(inputFile string) {
	lines := utils.ReadFileAsLines(inputFile)
	track := ParseInput(lines)

	fmt.Printf("Result of day-%s / part-2: %s\n", Day, FindLastCrashLocation(track))
}

func init() {
	register.Day(Day+"a", solvePart1)
	register.Day(Day+"b", solvePart2)
}
