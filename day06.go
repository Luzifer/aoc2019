package aoc2019

import (
	"bufio"
	"io"
	"os"
	"strings"

	"github.com/pkg/errors"
)

type day06Planet struct {
	Name   string
	Orbits *day06Planet
}

func (d *day06Planet) getPathLengthToPlanet(tp *day06Planet) int {
	if d == tp {
		return 0
	}

	return 1 + d.Orbits.getPathLengthToPlanet(tp)
}

type day06OrbitMap map[string]*day06Planet

func (d day06OrbitMap) getPathLengthToCOMByName(start string) int {
	return d[start].getPathLengthToPlanet(d["COM"])
}

func (d day06OrbitMap) getCommonPlanet(a, b *day06Planet) *day06Planet {
	var (
		pathAToCOM, pathBToCOM []*day06Planet
		ptr                    *day06Planet
	)

	// Get path A
	ptr = a
	for {
		pathAToCOM = append(pathAToCOM, ptr)
		if ptr.Orbits == nil {
			break
		}
		ptr = ptr.Orbits
	}

	// Get path B
	ptr = b
	for {
		pathBToCOM = append(pathBToCOM, ptr)
		if ptr.Orbits == nil {
			break
		}
		ptr = ptr.Orbits
	}

	// Find first common planet
	for _, ptrA := range pathAToCOM {
		for _, ptrB := range pathBToCOM {
			if ptrA == ptrB {
				return ptrA
			}
		}
	}

	return nil
}

func day06ParseOrbitMap(localMap io.Reader) (day06OrbitMap, error) {
	m := day06OrbitMap{}

	scanner := bufio.NewScanner(localMap)
	for scanner.Scan() {
		if err := scanner.Err(); err != nil {
			return nil, errors.Wrap(err, "Scanner failed to scan")
		}

		parts := strings.Split(scanner.Text(), ")")
		for _, p := range parts {
			if _, ok := m[p]; !ok {
				m[p] = &day06Planet{Name: p}
			}
		}

		m[parts[1]].Orbits = m[parts[0]]
	}

	return m, nil
}

func solveDay6Part1(inFile string) (int, error) {
	f, err := os.Open(inFile)
	if err != nil {
		return 0, errors.Wrap(err, "Unable to open input file")
	}
	defer f.Close()

	oMap, err := day06ParseOrbitMap(f)
	if err != nil {
		return 0, errors.Wrap(err, "Unable to parse orbit map")
	}

	var c int
	for _, p := range oMap {
		c += p.getPathLengthToPlanet(oMap["COM"])
	}

	return c, nil
}

func solveDay6Part2(inFile string) (int, error) {
	f, err := os.Open(inFile)
	if err != nil {
		return 0, errors.Wrap(err, "Unable to open input file")
	}
	defer f.Close()

	oMap, err := day06ParseOrbitMap(f)
	if err != nil {
		return 0, errors.Wrap(err, "Unable to parse orbit map")
	}

	cp := oMap.getCommonPlanet(oMap["YOU"], oMap["SAN"])
	dist := oMap["YOU"].getPathLengthToPlanet(cp) + oMap["SAN"].getPathLengthToPlanet(cp)

	// Distance is 2 too high as we're not a planet but a ship cycling a planet: Close enough
	return dist - 2, nil
}
