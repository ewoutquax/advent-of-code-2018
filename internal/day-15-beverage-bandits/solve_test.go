package day15beveragebandits_test

import (
	"fmt"
	"testing"

	. "github.com/ewoutquax/advent-of-code-2018/internal/day-15-beverage-bandits"
	"github.com/stretchr/testify/assert"
)

func TestParseInput(t *testing.T) {
	universe := ParseInput(testInputParse())

	assert := assert.New(t)
	assert.IsType(Universe{}, universe)

	assert.Len(universe.Positions, 15)                           // Test number of parsed locations
	assert.Len(universe.Creatures, 7)                            // Test number of creatures
	assert.Equal(TypeCreatureGoblin, universe.Creatures[0].Type) // Test type of creatures
	assert.IsType(TypeCreatureElve, universe.Creatures[len(universe.Creatures)-1].Type)
	assert.Equal(2, universe.Creatures[0].Location.X) // Test location of creatures
	assert.Equal(1, universe.Creatures[0].Location.Y)
	assert.Equal(4, universe.Creatures[len(universe.Creatures)-1].Location.X)
	assert.Equal(3, universe.Creatures[len(universe.Creatures)-1].Location.Y)

	// Test positions are linked to creatures
	posTopLeft := universe.Positions[1]
	assert.Equal(2, posTopLeft.Location.X)
	assert.Equal(1, posTopLeft.Location.Y)
	assert.IsType(TypeCreatureGoblin, posTopLeft.Creature.Type)
}

func TestLinkPositions(t *testing.T) {
	universe := ParseInput(testInputParse())

	assert := assert.New(t)

	posTopLeft := universe.Positions[0]
	assert.Equal(1, posTopLeft.Location.X)
	assert.Equal(1, posTopLeft.Location.Y)
	assert.Len(posTopLeft.LinkedPositions, 2)

	posBotRight := universe.Positions[len(universe.Positions)-1]
	assert.Equal(5, posBotRight.Location.X)
	assert.Equal(3, posBotRight.Location.Y)
	assert.Len(posBotRight.LinkedPositions, 2)
}

func TestPlayOrderOfCreatures(t *testing.T) {
	universe := ParseInput(testInputParse())

	creaturesInPlayOrder := PlayOrderOfCreatures(&universe)

	assert := assert.New(t)
	assert.Len(creaturesInPlayOrder, 7)
	assert.Equal(TypeCreatureGoblin, creaturesInPlayOrder[0].Type) // assert type of first and last creature
	assert.Equal(TypeCreatureElve, creaturesInPlayOrder[len(creaturesInPlayOrder)-1].Type)
	assert.Equal(2, creaturesInPlayOrder[0].Location.X) // assert location of first and last creature
	assert.Equal(1, creaturesInPlayOrder[0].Location.Y)
	assert.Equal(4, creaturesInPlayOrder[len(creaturesInPlayOrder)-1].Location.X)
	assert.Equal(3, creaturesInPlayOrder[len(creaturesInPlayOrder)-1].Location.Y)
}

func TestFindPaths(t *testing.T) {
	universe := ParseInput(testInputMove())
	elve := universe.Creatures[0]

	var path Path = elve.FindPathToNearestEnemy()
	fmt.Printf("pathToNearestEnemy: %v\n", path.ToS())

	assert := assert.New(t)

	assert.Len(path.Positions, 3)

	assert.Equal(2, path.Positions[1].Location.X)
	assert.Equal(1, path.Positions[1].Location.Y)
	assert.Equal(3, path.Positions[2].Location.X)
	assert.Equal(1, path.Positions[2].Location.Y)
}

func TestFindPathsWhenNoPathToEnemy(t *testing.T) {
	universe := ParseInput(testInputMoveNoPath())
	elve := universe.Creatures[0]

	var path Path = elve.FindPathToNearestEnemy()

	assert.Len(t, path.Positions, 0)
}

func TestFindPathsWhenNoMoveRequired(t *testing.T) {
	universe := ParseInput(testInputMoveNoMoveRequired())
	elve := universe.Creatures[0]

	var path Path = elve.FindPathToNearestEnemy()

	assert.Len(t, path.Positions, 1)
}

func TestMove(t *testing.T) {
	universe := ParseInput(testInputMove())

	elve := universe.Creatures[0]
	origPos := elve.Position
	destPos := universe.Positions[1]
	fmt.Printf("destPos: %v\n", destPos.ToS())
	elve.MoveTo(destPos)

	assert := assert.New(t)

	assert.Nil(origPos.Creature)                          // The original position should no longer contain a creature
	assert.Equal(TypeCreatureElve, destPos.Creature.Type) // The destination position should now contain a creature
	assert.Equal(2, elve.Location.X)                      // The elve is now at the new position
	assert.Equal(1, elve.Location.Y)
}

func TestFindCreatureToAttack(t *testing.T) {
	assert := assert.New(t)

	universe := ParseInput(testInputAttack())
	universe.Creatures[0].Health = 2  // first found creature, but of same type
	universe.Creatures[2].Health = 10 // first found creature
	universe.Creatures[3].Health = 5  // second found creature, with least healt
	universe.Creatures[4].Health = 2  // Diagonal creature, not viable

	elve1 := universe.Creatures[0]
	assert.Equal(1, elve1.Position.Location.X)
	assert.Equal(1, elve1.Position.Location.Y)
	isFound, _ := elve1.FindCreatureToAttack()
	assert.False(isFound)

	elve2 := universe.Creatures[1]
	assert.Equal(1, elve2.Position.Location.X) // Did we select the correct creature/
	assert.Equal(2, elve2.Position.Location.Y)

	isFound, creatureToAttack := elve2.FindCreatureToAttack()
	assert.True(isFound)
	assert.Equal(TypeCreatureGoblin, creatureToAttack.Type)
	assert.Equal(1, creatureToAttack.Position.Location.X) // Are we attacking the intended creature
	assert.Equal(3, creatureToAttack.Position.Location.Y)
}

func TestDoAttack(t *testing.T) {
	universe := ParseInput(testInputAttack())

	assert := assert.New(t)

	universe.Creatures[0].Health = 2  // first found creature, but of same type
	universe.Creatures[2].Health = 10 // first found creature
	universe.Creatures[3].Health = 5  // second found creature, with least healt
	universe.Creatures[4].Health = 2  // Diagonal creature, not viable

	elve := universe.Creatures[1]
	_, creatureToAttack := elve.FindCreatureToAttack()
	posWithAttacker := creatureToAttack.Position

	elve.DoAttack(creatureToAttack)
	assert.Equal(200, elve.Health)
	assert.True(creatureToAttack.IsAlive())
	assert.Equal(2, creatureToAttack.Health)

	// attack 'm again
	elve.DoAttack(creatureToAttack)
	assert.Equal(200, elve.Health)
	assert.False(creatureToAttack.IsAlive())
	assert.Nil(posWithAttacker.Creature)
}

func TestIsExtraRoundRequiredTrue(t *testing.T) {
	universe := ParseInput(testInputAttack())

	assert.True(t, universe.IsExtraRoundRequired())
}

func TestIsExtraRoundRequiredFalse(t *testing.T) {
	universe := ParseInput(testInputAttack())

	// Kill both elves
	universe.Creatures[0].Health = -2
	universe.Creatures[1].Health = -2

	assert.False(t, universe.IsExtraRoundRequired())
}

func TestExample1(t *testing.T) {
	universe := ParseInput(testInputExample1())

	score := PlayGame(&universe)

	assert.Equal(t, 2, score)
}

func testInputParse() []string {
	return []string{
		"#######",
		"#.G.E.#",
		"#E.G.E#",
		"#.G.E.#",
		"#######",
	}
}

func testInputMove() []string {
	return []string{
		"#######",
		"#E..G.#",
		"#.E.#.#",
		"#.G.#G#",
		"#######",
	}
}

func testInputMoveNoPath() []string {
	return []string{
		"#######",
		"#EE.G.#",
		"###.#.#",
		"#.G.#G#",
		"#######",
	}
}

func testInputMoveNoMoveRequired() []string {
	return []string{
		"#######",
		"#.....#",
		"#.EG..#",
		"#.....#",
		"#######",
	}
}

func testInputAttack() []string {
	return []string{
		"######",
		"#E.. #",
		"#EG. #",
		"#GG. #",
		"######",
	}
}

func testInputExample1() []string {
	return []string{
		"#######",
		"#.G...#",
		"#...EG#",
		"#.#.#G#",
		"#..G#E#",
		"#.....#",
		"#######",
	}
}
