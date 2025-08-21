package day13minecartmadness_test

import (
	"fmt"
	"image"
	"testing"

	. "github.com/ewoutquax/advent-of-code-2018/internal/day-13-mine-cart-madness"
	"github.com/stretchr/testify/assert"
)

func TestParseInput(t *testing.T) {
	assert := assert.New(t)

	track := ParseInput(testInput())

	assert.IsType(Track{}, track)
	assert.Len(track.Rails, 48)
	assert.Len(track.Carts, 2)
	assert.Contains(track.Carts, Location(image.Pt(2, 0)))
	assert.Contains(track.Carts, Location(image.Pt(9, 3)))

	// Cart
	cart := track.Carts[Location(image.Pt(2, 0))]
	assert.Equal("day13minecartmadness.Cart", fmt.Sprintf("%T", cart))
	assert.Equal(2, cart.Location.X)
	assert.Equal(0, cart.Location.Y)
	assert.Equal(DirectionRight, cart.Direction)
	assert.Equal(0, cart.NrMoves)
}

func TestCartMoveCorner(t *testing.T) {
	track := ParseInput(testInput())
	cart := track.Carts[Location(image.Pt(2, 0))]

	// Move one location right
	cart.Move(track)
	assert.Equal(t, 3, cart.Location.X)
	assert.Equal(t, 0, cart.Location.Y)
	assert.Equal(t, 1, cart.NrMoves)
	assert.Equal(t, DirectionRight, cart.Direction)
	assert.Equal(t, TurnLeft, cart.NextTurn)

	// Move one location right, and make a forced turn
	cart.Move(track)
	assert.Equal(t, 4, cart.Location.X)
	assert.Equal(t, 0, cart.Location.Y)
	assert.Equal(t, 2, cart.NrMoves)
	assert.Equal(t, DirectionDown, cart.Direction)
	assert.Equal(t, TurnLeft, cart.NextTurn)
}

func TestCartMoveIntersection(t *testing.T) {
	track := ParseInput(testInput())
	cart := track.Carts[Location(image.Pt(9, 3))]

	// Move one location down, and make the preferred turn
	cart.Move(track)
	assert.Equal(t, 9, cart.Location.X)
	assert.Equal(t, 4, cart.Location.Y)
	assert.Equal(t, 1, cart.NrMoves)
	assert.Equal(t, DirectionRight, cart.Direction)
	assert.Equal(t, TurnForward, cart.NextTurn)
}

func TestFindFirstCrashLocation(t *testing.T) {
	track := ParseInput(testInput())

	assert.Equal(t, "7,3", FindFirstCrashLocation(track))
}

func TestFindLastCrashLocation(t *testing.T) {
	track := ParseInput(testInput2())

	assert.Equal(t, "6,4", FindLastCrashLocation(track))
}

func testInput() []string {
	return []string{
		`/->-\`,
		`|   |  /----\`,
		`| /-+--+-\  |`,
		`| | |  | v  |`,
		`\-+-/  \-+--/`,
		`  \------/`,
	}
}

func testInput2() []string {
	return []string{
		`/>-<\`,
		`|   |`,
		`| /<+-\`,
		`| | | v`,
		`\>+</ |`,
		`  |   ^`,
		`  \<->/`,
	}
}
