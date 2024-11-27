package day22modemaze

import "fmt"

type Path struct {
	Location
	CurrentEquipment Equipment
	NrMinutes        int
	distance         int
}

func (p Path) IsVisited() bool {
	/**
	fmt.Printf("\nlen(visitedPaths): %v\n", len(visitedPaths))
	fmt.Println("------------------")
	for key, value := range visitedPaths {
		fmt.Printf("'%s': %v\n", key, value)
	}
	*/

	visitedNrMinutes, exists := visitedPaths[p.VisitedKey()]
	return !(!exists || exists && visitedNrMinutes > p.NrMinutes)
}

func (p Path) SetVisited() {
	visitedPaths[p.VisitedKey()] = p.NrMinutes
}

type VisitedPathKey string

func (p Path) VisitedKey() VisitedPathKey {
	return VisitedPathKey(fmt.Sprintf("%s/%d", p.Location.toS(), p.CurrentEquipment))
}

// An IntHeap is a min-heap of ints.
type PathHeap []Path

func (h PathHeap) Len() int { return len(h) }

// func (h PathHeap) Less(i, j int) bool { return h[i].NrMinutes < h[j].NrMinutes }
func (h PathHeap) Less(i, j int) bool { return h[i].distance < h[j].distance }
func (h PathHeap) Swap(i, j int)      { h[i], h[j] = h[j], h[i] }
func (h *PathHeap) Push(x any)        { *h = append(*h, x.(Path)) }

func (h *PathHeap) Pop() any {
	old := *h
	n := len(old)
	x := old[n-1]
	*h = old[0 : n-1]
	return x
}

func (p Path) toS() string {
	return fmt.Sprintf(
		"%s / tool %s: %d",
		p.Location.toS(),
		p.CurrentEquipment.toS(),
		p.NrMinutes,
	)
}

func (p Path) getNextPaths(targetLocation Location, caveDepth CaveDepth) []Path {
	var out []Path = make([]Path, 0, 8)

	for _, vector := range vectors() {
		newLocation := Location{
			X: p.Location.X + vector[0],
			Y: p.Location.Y + vector[1],
		}

		if newLocation.X >= 0 && newLocation.Y >= 0 {
			equipments := allowedEquipments(newLocation.ErosionLevelType(targetLocation, caveDepth))
			for _, equipment := range equipments {
				newPath := Path{
					Location: newLocation, CurrentEquipment: equipment,
					NrMinutes: p.NrMinutes + 1,
				}
				if equipment != p.CurrentEquipment {
					newPath.NrMinutes += 7
				}
				newPath.distance = abs(targetLocation.X-newLocation.X) + abs(targetLocation.Y-newLocation.Y) + newPath.NrMinutes
				out = append(out, newPath)
			}
		}
	}

	return out
}

func InitVisitedPaths() {
	visitedPaths = make(map[VisitedPathKey]int)
}

func abs(i int) int {
	if i < 0 {
		return 0 - i
	}
	return i
}

func vectors() [4][2]int {
	return [4][2]int{
		{-1, 0},
		{1, 0},
		{0, -1},
		{0, 1},
	}
}
