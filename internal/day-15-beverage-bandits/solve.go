package day15beveragebandits

import (
	"fmt"
	"slices"
	"strings"

	"github.com/ewoutquax/advent-of-code-2018/pkg/register"
)

const Day string = "15"

type Direction uint

const CharGoblin string = "G"
const CharElve string = "E"

type TypeCreature uint

const (
	TypeCreatureGoblin TypeCreature = iota + 1
	TypeCreatureElve
)

type Creature struct {
	Type TypeCreature
	*Position
	Health int
}

func (c Creature) isDifferentType(other Creature) bool {
	return c.Type != other.Type
}

func (c Creature) isAlive() bool {
	panic("Not yet implemented")
}

const (
	DirectionUp Direction = iota + 1
	DirectionLeft
	DirectionRight
	DirectionDown
)

type Location struct {
	X int
	Y int
}

func (l Location) toI() int { return l.Y*1000 + l.X }
func (l Location) inDirection(d Direction) Location {
	var vector [2]int

	switch d {
	case DirectionUp:
		vector = [2]int{0, -1}
	case DirectionLeft:
		vector = [2]int{-1, 0}
	case DirectionRight:
		vector = [2]int{1, 0}
	case DirectionDown:
		vector = [2]int{0, 1}
	default:
		panic("No valid direction found")
	}

	return Location{
		X: l.X + vector[0],
		Y: l.Y + vector[1],
	}
}

// type Goblin struct{ *Position }
// type Elve struct{ *Position }
//
// func (g Goblin) position() *Position             { return g.Position }
// func (g Goblin) isDifferentType(c Creature) bool { return fmt.Sprintf("%T", c) != fmt.Sprintf("%T", g) }
// func (g Goblin) isAlive() bool                   { return true }
// func (e Elve) position() *Position               { return e.Position }
// func (e Elve) isDifferentType(c Creature) bool   { return fmt.Sprintf("%T", c) != fmt.Sprintf("%T", e) }
// func (e Elve) isAlive() bool                     { return true }

type Position struct {
	Location
	*Creature
	LinkedPositions []*Position
}

func (p Position) toI() int    { return p.Location.toI() }
func (p Position) ToS() string { return fmt.Sprintf("[%d,%d]", p.Location.X, p.Location.Y) }

type Path struct {
	Positions []*Position
}

func (p Path) ToS() string {
	var out = make([]string, 0, len(p.Positions))

	for _, pos := range p.Positions {
		out = append(out, pos.ToS())
	}

	return strings.Join(out, "; ")
}

type Universe struct {
	Positions []*Position
	Creatures []Creature
}

func (c *Creature) MoveTo(newPos *Position) {
	c.Position.Creature = nil
	newPos.Creature = c
	c.Position = newPos
}

func FindCreatureToAttack(c Creature) (bool, Creature) {
	var isFound bool = false
	var FoundCreature Creature

	for _, nextPos := range c.Position.LinkedPositions {
		if nextPos.Creature != nil {
			if !isFound || isFound && FoundCreature.Health > nextPos.Creature.Health {
			}
		}
	}

	panic("Not yet implemented")
}

func FindPathToNearestEnemy(creature Creature) Path {
	var doContinue bool = true
	var maxSteps int = 0
	var visitedLocations = make(map[int]bool)
	var foundPath Path
	var foundPaths = make(map[int][]Path)
	foundPaths[0] = []Path{
		{Positions: []*Position{
			creature.Position,
		}},
	}
	visitedLocations[creature.Position.toI()] = true

	var ctrPanic int = 0

	for doContinue {
		if ctrPanic > 100 {
			panic("Counter reached")
		}
		ctrPanic++
		foundPaths[maxSteps+1] = make([]Path, 0)

		for _, currPath := range foundPaths[maxSteps] {
			fmt.Printf("currPath: %v\n", currPath.ToS())

			currPos := currPath.Positions[len(currPath.Positions)-1]
			for _, linkedPosition := range currPos.LinkedPositions {
				_, exists := visitedLocations[linkedPosition.toI()]
				if !exists && linkedPosition.Creature == nil {
					fmt.Printf("we move from %s to %s\n", currPos.ToS(), linkedPosition.ToS())

					visitedLocations[linkedPosition.toI()] = true
					foundPaths[maxSteps+1] = append(foundPaths[maxSteps+1], Path{
						Positions: append(currPath.Positions, linkedPosition),
					})
				} else {
					if !exists && linkedPosition.Creature != nil && linkedPosition.Creature.isDifferentType(creature) {
						if len(foundPath.Positions) == 0 {
							foundPath = currPath
						}

						doContinue = false

						fmt.Printf("Creature found at %s\n", linkedPosition.ToS())
						fmt.Printf("foundPath: %v\n", foundPath.ToS())
						break
					}
				}
			}
		}

		maxSteps += 1
		fmt.Printf("All positions for current nr step exhausted; Increasing to: %v\n", maxSteps)
	}

	fmt.Printf("returning foundPath: %v\n", foundPath.ToS())
	return foundPath
}

func PlayOrderOfCreatures(u *Universe) []Creature {
	type CreatureWithOrderNr struct {
		orderNr int
		Creature
	}

	var creaturesWithOrderNr []CreatureWithOrderNr = make([]CreatureWithOrderNr, 0, len(u.Creatures))

	for _, creature := range u.Creatures {
		creaturesWithOrderNr = append(creaturesWithOrderNr, CreatureWithOrderNr{
			orderNr:  creature.Position.Location.toI(),
			Creature: creature,
		})
	}

	slices.SortFunc(creaturesWithOrderNr, func(i, j CreatureWithOrderNr) int {
		if i.orderNr < j.orderNr {
			return -1
		}
		return 1
	})

	var out = make([]Creature, 0, len(u.Creatures))

	for _, creaturesWithOrderNr := range creaturesWithOrderNr {
		out = append(out, creaturesWithOrderNr.Creature)
	}

	return out
}

func ParseInput(lines []string) Universe {
	var universe = Universe{
		Positions: make([]*Position, 0),
		Creatures: make([]Creature, 0),
	}

	for y, line := range lines {
		for x, char := range strings.Split(line, "") {
			if char != "#" {
				loc := Location{X: x, Y: y}

				pos := &Position{
					Location:        loc,
					Creature:        nil,
					LinkedPositions: make([]*Position, 0),
				}
				universe.Positions = append(universe.Positions, pos)

				switch char {
				case CharGoblin:
					goblin := Creature{
						Type:     TypeCreatureGoblin,
						Position: pos,
						Health:   200,
					}
					universe.Creatures = append(universe.Creatures, goblin)
					pos.Creature = &goblin
				case CharElve:
					elve := Creature{
						Type:     TypeCreatureElve,
						Position: pos,
						Health:   200,
					}
					universe.Creatures = append(universe.Creatures, elve)
					pos.Creature = &elve
				default:
					continue
				}
			}
		}
	}

	linkPositions(&universe)

	return universe
}

func linkPositions(u *Universe) {
	var indexedPositions = make(map[Location]*Position, len(u.Positions))

	for _, position := range u.Positions {
		indexedPositions[position.Location] = position
	}

	for curLoc, curPos := range indexedPositions {

		for _, d := range directionsInReadingOrder() {
			if newPos, exists := indexedPositions[curLoc.inDirection(d)]; exists {
				curPos.LinkedPositions = append(curPos.LinkedPositions, newPos)
			}
		}
	}
}

func directionsInReadingOrder() []Direction {
	return []Direction{
		DirectionUp, DirectionLeft, DirectionRight, DirectionDown,
	}
}

func solvePart1(inputFile string) {
	// lines := utils.ReadFileAsLines(inputFile)
	var count int = 0
	fmt.Printf("Result of day-%s / part-1: %d\n", Day, count)
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
