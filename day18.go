package aoc2019

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"math"
	"sort"
	"strings"
	"sync"

	"github.com/pkg/errors"
)

type day18TileType int

const (
	day18TileTypeWall day18TileType = iota
	day18TileTypeFloor
	day18TileTypeDoor
	day18TileTypeKey
	day18TileTypeActor
)

type day18Tile struct {
	X, Y      int
	Type      day18TileType
	KeyLockID string

	distFromActor int64
}

func (d *day18Tile) updateBFSDistance(grid day18Grid, dist int64, wg *sync.WaitGroup) {
	defer wg.Done()
	if d.distFromActor < dist {
		// Do not update tiles which dist has been set already
		return
	}

	// Set distance
	d.distFromActor = dist

	// Seek adjacent tiles and update them too
	for _, t := range []*day18Tile{
		grid[grid.key(d.X-1, d.Y)],
		grid[grid.key(d.X+1, d.Y)],
		grid[grid.key(d.X, d.Y-1)],
		grid[grid.key(d.X, d.Y+1)],
	} {
		if t.Type == day18TileTypeWall || t.Type == day18TileTypeDoor {
			// Doors and walls are inaccessible objects and therefore
			// immediately stop distance progressing
			continue
		}

		wg.Add(1)
		go t.updateBFSDistance(grid, dist+1, wg)
	}
}

type day18Grid map[string]*day18Tile

func (d day18Grid) accessibleKeys() []string {
	var keys []string

	for _, key := range d.getTilesByType(day18TileTypeKey) {
		if key.distFromActor < math.MaxInt64 {
			keys = append(keys, key.KeyLockID)
		}
	}

	sort.Strings(keys)
	return keys
}

func (d day18Grid) bounds() (minX, minY, maxX, maxY int) {
	minX, minY = math.MaxInt64, math.MaxInt64

	for _, t := range d {
		if t.X < minX {
			minX = t.X
		}
		if t.Y < minY {
			minY = t.Y
		}
		if t.X > maxX {
			maxX = t.X
		}
		if t.Y > maxY {
			maxY = t.Y
		}
	}

	return
}

func (d day18Grid) clone() day18Grid {
	var out = make(day18Grid)

	for key, tile := range d {
		var tmpTile = &day18Tile{}
		*tmpTile = *tile
		out[key] = tmpTile
	}

	return out
}

func (d day18Grid) getKeyByID(id string) *day18Tile {
	for _, key := range d.getTilesByType(day18TileTypeKey) {
		if key.KeyLockID == id {
			return key
		}
	}
	return nil
}

func (d day18Grid) getPathDistances(path string, shortest *int64) (map[string]int64, error) {
	var (
		dist int64
		// Always work on a clone and leave the original grid intact
		grid = d.clone()
		out  = make(map[string]int64)
	)

	grid.updateBFSDistance()

	// Apply movements
	for _, c := range path {
		key := grid.getKeyByID(string(c))
		moveDist, err := grid.move(key.X, key.Y)
		if err != nil {
			return nil, errors.Wrap(err, "Unable to move")
		}
		dist += moveDist
		grid.updateBFSDistance()

		if dist >= *shortest {
			return nil, nil
		}
	}

	remainingKeys := grid.getTilesByType(day18TileTypeKey)

	if len(remainingKeys) == 0 {
		fmt.Printf("DBG %s = %d\n", path, dist)
		out[path] = dist
		if dist < *shortest {
			*shortest = dist
		}
		return out, nil
	}

	sort.Slice(remainingKeys, func(i, j int) bool { return remainingKeys[i].distFromActor < remainingKeys[j].distFromActor })

	for _, k := range remainingKeys {
		if k.distFromActor == math.MaxInt64 {
			continue
		}

		recDists, err := d.getPathDistances(path+k.KeyLockID, shortest)
		if err != nil {
			return nil, err
		}

		for key, dist := range recDists {
			out[key] = dist
		}
	}

	return out, nil
}

func (d day18Grid) getTilesByType(t day18TileType) []*day18Tile {
	var out []*day18Tile

	for _, tile := range d {
		if tile.Type == t {
			out = append(out, tile)
		}
	}

	return out
}

func (day18Grid) key(x, y int) string { return fmt.Sprintf("%d:%d", x, y) }

func (d day18Grid) move(x, y int) (int64, error) {
	var (
		origin = d.getTilesByType(day18TileTypeActor)[0] // There MUST be only one actor
		target = d[d.key(x, y)]
	)

	if target.distFromActor == math.MaxInt64 {
		return 0, errors.New("Distance not initialized / move to unreachable destination")
	}

	if target.Type == day18TileTypeKey {
		for _, door := range d.getTilesByType(day18TileTypeDoor) {
			if door.KeyLockID == strings.ToUpper(target.KeyLockID) {
				// Unlock door, transform to floor
				door.Type = day18TileTypeFloor
				break
			}
		}
	}

	// Move the actor
	target.Type = day18TileTypeActor
	origin.Type = day18TileTypeFloor

	return target.distFromActor, nil
}

func (d day18Grid) print() {
	var minX, minY, maxX, maxY = d.bounds()

	for y := minY; y <= maxY; y++ {
		for x := minX; x <= maxX; x++ {
			var t = d[d.key(x, y)]
			switch t.Type {
			case day18TileTypeActor:
				fmt.Print("@")
			case day18TileTypeDoor, day18TileTypeKey:
				fmt.Print(t.KeyLockID)
			case day18TileTypeFloor:
				fmt.Print(".")
			case day18TileTypeWall:
				fmt.Print("\u2588")
			}
		}
		fmt.Println()
	}
}

func (d day18Grid) updateBFSDistance() {
	// Reset distances
	for _, t := range d {
		t.distFromActor = math.MaxInt64
	}

	wg := new(sync.WaitGroup)
	wg.Add(1)
	d.getTilesByType(day18TileTypeActor)[0].updateBFSDistance(d, 0, wg)
	wg.Wait()
}

func day18ReadGrid(raw []byte) day18Grid {
	var (
		out  = make(day18Grid)
		x, y int
	)

	for _, c := range raw {
		var t = &day18Tile{X: x, Y: y, KeyLockID: string(c)}

		switch {
		case c == '\n':
			x, y = 0, y+1
			continue
		case c == '#':
			t.Type = day18TileTypeWall
		case c == '.':
			t.Type = day18TileTypeFloor
		case c == '@':
			t.Type = day18TileTypeActor
		case 'a' <= c && c <= 'z':
			t.Type = day18TileTypeKey
		case 'A' <= c && c <= 'Z':
			t.Type = day18TileTypeDoor
		}

		out[out.key(x, y)] = t
		x++
	}

	return out
}

func solveDay18Part1(inFile string) (int64, error) {
	raw, err := ioutil.ReadFile(inFile)
	if err != nil {
		return 0, errors.Wrap(err, "Unable to read file")
	}

	grid := day18ReadGrid(bytes.TrimSpace(raw))

	grid.print()

	// Calculate max heuristic distance by using always using shortest
	// single distance which will not be the shortest total distance
	// but set an upper limit for the heuristic
	cgrid := grid.clone()
	cgrid.updateBFSDistance()
	var minDist int64
	for len(cgrid.accessibleKeys()) > 0 {
		keys := cgrid.getTilesByType(day18TileTypeKey)
		sort.Slice(keys, func(i, j int) bool { return keys[i].distFromActor < keys[j].distFromActor })
		d, err := cgrid.move(keys[0].X, keys[0].Y)
		if err != nil {
			return 0, errors.Wrap(err, "Failed to move")
		}
		minDist += d
		cgrid.updateBFSDistance()
	}

	log.Printf("Shortest distance for heuristic: %d", minDist)

	dists, err := grid.getPathDistances("", &minDist)
	if err != nil {
		return 0, errors.Wrap(err, "Unable to determine distances")
	}

	fmt.Printf("%#v\n", dists)

	return minDist, nil
}

func solveDay18Part2(inFile string) (int64, error) { return 0, errors.New("Not implemented") }
