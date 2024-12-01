package day23experimentalemergencyteleportation

import "github.com/ewoutquax/advent-of-code-2018/pkg/register"

const (
	Day string = "23"
)

func init() {
	register.Day(Day+"a", solvePart1)
	register.Day(Day+"b", solvePart2)
}
