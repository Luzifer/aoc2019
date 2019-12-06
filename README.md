# Luzifer / aoc2019

This repository contains my solutions for [Advent of Code 2019](https://adventofcode.com/2019/).

To execute them just use the tests for the respective day through `make dayXX` (some days require extra code files, the `make` instructions contains them):

```console
# make day01
go test -cover -v \
	day01.go day01_test.go \
	helpers.go intcode.go
=== RUN   TestCalculateDay1_Examples
--- PASS: TestCalculateDay1_Examples (0.00s)
=== RUN   TestCalculateDay1_Part1
--- PASS: TestCalculateDay1_Part1 (0.00s)
    day01_test.go:34: Solution Day 1 Part 1: 3126794
=== RUN   TestCalculateDay1_Part2
--- PASS: TestCalculateDay1_Part2 (0.00s)
    day01_test.go:43: Solution Day 1 Part 2: 4687331
PASS
coverage: 23.3% of statements
ok  	command-line-arguments	(cached)	coverage: 23.3% of statements
```
