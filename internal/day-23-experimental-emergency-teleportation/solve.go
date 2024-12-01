package day23experimentalemergencyteleportation

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/ewoutquax/advent-of-code-2018/pkg/utils"
)

type Radius int

type Location struct {
	X int
	Y int
	Z int
}

type Bot struct {
	Location
	Radius
}

func (b Bot) Distance(target Bot) int {
	return abs(target.X-b.X) +
		abs(target.Y-b.Y) +
		abs(target.Z-b.Z)
}

func (b Bot) toS() string {
	return fmt.Sprintf("Pos=<%d,%d,%d>, r=%d", b.X, b.Y, b.Z, b.Radius)
}

func FindLocationWithMostCoverage(bots []Bot) Location {
	var bestBot Bot
	var maxCount int = 0

	var bestBots []Bot = make([]Bot, 0, len(bots))

	for _, sourceBot := range bots {
		count := 1
		fmt.Println("\nNew sourcebot")
		newBot := sourceBot

		for _, targetBot := range bots {
			tempBot, err := GenerateBotWithFullOverlap(newBot, targetBot)
			if err == nil {
				// fmt.Printf("newBot: '%s' and targetBot '%s' have full overlap at '%s'\n", newBot.toS(), targetBot.toS(), tempBot.toS())

				count++
				newBot = tempBot
				// } else {
				// fmt.Printf("newBot: '%s' and targetBot '%s' have NO overlap\n", newBot.toS(), targetBot.toS())
			}

		}

		if maxCount < count {
			fmt.Printf("new bestLocation with bot: %v\n", newBot.toS())
			bestBot = newBot
			maxCount = count
			bestBots = append(bestBots, newBot)
		}
	}

	fmt.Println("FindLocationWithMostCoverage: bestBots (after round 1)")
	fmt.Println("------------------------------------------------------")
	fmt.Printf("len(bestBots): %v\n", len(bestBots))

	var bestestBots []Bot = make([]Bot, 0, len(bestBots))
	for _, sourceBot := range bestBots {
		count := 1
		fmt.Println("\nNew sourceBestBot")
		newBot := sourceBot

		for _, targetBot := range bots {
			tempBot, err := GenerateBotWithFullOverlap(newBot, targetBot)
			if err == nil {
				count++
				newBot = tempBot
			}

		}

		if maxCount < count {
			bestBot = newBot
			maxCount = count
			bestestBots = append(bestestBots, newBot)
		}
	}

	fmt.Println("FindLocationWithMostCoverage: bestestBots (after round 2)")
	fmt.Println("---------------------------------------------------------")
	for idx, b := range bestestBots {
		fmt.Printf("%d: %s\n", idx, b.toS())
	}

	fmt.Printf("overall bestBot: %s\n", bestBot.toS())

	return bestBot.Location
}

func CountWithinRangeOfStrongest(bots []Bot) int {
	var strongestBot Bot = bots[0]

	for _, bot := range bots {
		if strongestBot.Radius < bot.Radius {
			strongestBot = bot
		}
	}

	var count int = 0
	for _, bot := range bots {
		if strongestBot.Distance(bot) <= int(strongestBot.Radius) {
			count++
		}
	}

	return count
}

func GenerateBotWithFullOverlap(b1, b2 Bot) (Bot, error) {
	if b1.Distance(b2) > int(b1.Radius+b2.Radius) {
		return Bot{}, errors.New("No overlap")
	}

	var newBotMinX, newBotMaxX, newBotMinY, newBotMaxY, newBotMinZ, newBotMaxZ int

	newBotMinX = maximum(b1.X-int(b1.Radius), b2.X-int(b2.Radius))
	newBotMaxX = minimum(b1.X+int(b1.Radius), b2.X+int(b2.Radius))
	newBotMinY = maximum(b1.Y-int(b1.Radius), b2.Y-int(b2.Radius))
	newBotMaxY = minimum(b1.Y+int(b1.Radius), b2.Y+int(b2.Radius))
	newBotMinZ = maximum(b1.Z-int(b1.Radius), b2.Z-int(b2.Radius))
	newBotMaxZ = minimum(b1.Z+int(b1.Radius), b2.Z+int(b2.Radius))

	var newBot Bot = Bot{
		Location: Location{
			X: (newBotMinX + newBotMaxX) / 2,
			Y: (newBotMinY + newBotMaxY) / 2,
			Z: (newBotMinZ + newBotMaxZ) / 2,
		},
		Radius: 0,
	}

	newBot.Radius = Radius(minimum(
		(newBotMaxX-newBotMinX)/2,
		(newBotMaxY-newBotMinY)/2,
		(newBotMaxZ-newBotMinZ)/2,
	))

	return newBot, nil
}

func maximum(ints ...int) int {
	if len(ints) == 1 {
		return ints[0]
	}

	i2 := maximum(ints[1:]...)
	if ints[0] > i2 {
		return ints[0]
	}

	return i2
}

func minimum(ints ...int) int {
	if len(ints) == 1 {
		return ints[0]
	}

	i2 := minimum(ints[1:]...)
	if ints[0] < i2 {
		return ints[0]
	}

	return i2
}

func ParseInput(lines []string) []Bot {
	var bots []Bot = make([]Bot, 0, len(lines))

	for _, line := range lines {
		bots = append(bots, parseLine(line))
	}

	return bots
}

func parseLine(line string) Bot {
	parts := strings.Split(line, ", ")

	locationParts := strings.TrimRight(
		strings.Split(parts[0], "<")[1],
		">",
	)

	locations := strings.Split(locationParts, ",")
	radius := strings.Split(parts[1], "=")[1]

	return Bot{
		Location: Location{
			X: convAtoi(locations[0]),
			Y: convAtoi(locations[1]),
			Z: convAtoi(locations[2]),
		},
		Radius: Radius(convAtoi(radius)),
	}
}

func convAtoi(s string) int {
	nr, err := strconv.Atoi(s)
	if err != nil {
		panic(err)
	}

	return nr
}

func abs(i int) int {
	if i < 0 {
		return 0 - i
	}

	return i
}

func solvePart1(inputFile string) {
	lines := utils.ReadFileAsLines(inputFile)
	bots := ParseInput(lines)

	fmt.Printf("Result of day-%s / part-1: %d\n", Day, CountWithinRangeOfStrongest(bots))
}

func solvePart2(inputFile string) {
	lines := utils.ReadFileAsLines(inputFile)
	bots := ParseInput(lines)

	location := FindLocationWithMostCoverage(bots)
	distance := abs(location.X) +
		abs(location.Y) +
		abs(location.Z)

		/**
		  Answers:
		  --------
		  120339115: too high
		*/

	fmt.Printf("Result of day-%s / part-2: %d\n", Day, distance)
}
