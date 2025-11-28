// -----------------------------------------------------------------------------
// Author: Carter Ward
// Course: CS424-01 (Fall 2025)
// Date  : 9/28/2025
// Environment: VS Code with Go 1.25
//
// Purpose:
//   This program reads basketball player statistics from a text file provided
//   by the user, stores them in a structured list, sorts the players by last
//   name (and first name if needed), calculates eFG% and TS% for each player,
//   and displays the results in a simple, organized report. Players with
//   incomplete data are clearly flagged.
// -----------------------------------------------------------------------------

package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
	"strconv"
	"strings"
)

// ---------------------------
// Data model
// ---------------------------

// Player holds the raw totals read from the file.
// BadData marks a row with missing numeric fields.
type Player struct {
	First, Last   string
	Games, Points int
	FGMade, FGAtt int
	TPMade, TPAtt int
	FTMade, FTAtt int
	BadData       bool
}

// eFG: (FGmade + 0.5*TPmade) / FGattempts. Returns a fraction [0,1].
func (p Player) eFG() float64 {
	if p.FGAtt == 0 {
		return 0
	}
	return (float64(p.FGMade) + 0.5*float64(p.TPMade)) / float64(p.FGAtt)
}

// TS: Points / (2 * (FGattempts + 0.44*FTattempts)). Returns a fraction [0,1].
func (p Player) ts() float64 {
	den := 2.0 * (float64(p.FGAtt) + 0.44*float64(p.FTAtt))
	if den == 0 {
		return 0
	}
	return float64(p.Points) / den
}

// ---------------------------
// Parsing / file I/O
// ---------------------------

// parsePlayer reads one line into a Player.
// Expects: first last + 8 integers. Marks BadData if fields are missing.
func parsePlayer(line string) Player {
	fields := strings.Fields(line)

	if len(fields) < 10 {
		var first, last string
		if len(fields) >= 1 {
			first = fields[0]
		}
		if len(fields) >= 2 {
			last = fields[1]
		}
		return Player{First: first, Last: last, BadData: true}
	}

	nums := make([]int, 0, 8)
	for _, s := range fields[2:10] {
		v, err := strconv.Atoi(s)
		if err != nil {
			return Player{First: fields[0], Last: fields[1], BadData: true}
		}
		nums = append(nums, v)
	}

	return Player{
		First:  fields[0],
		Last:   fields[1],
		Games:  nums[0],
		Points: nums[1],
		FGMade: nums[2],
		FGAtt:  nums[3],
		TPMade: nums[4],
		TPAtt:  nums[5],
		FTMade: nums[6],
		FTAtt:  nums[7],
	}
}

// readPlayersFromFile loads all players from filename.
// Ignores blank lines; returns any scan error.
func readPlayersFromFile(filename string) ([]Player, error) {
	f, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	var players []Player
	sc := bufio.NewScanner(f)
	for sc.Scan() {
		txt := strings.TrimSpace(sc.Text())
		if txt == "" {
			continue
		}
		players = append(players, parsePlayer(txt))
	}
	if err := sc.Err(); err != nil {
		return nil, err
	}
	return players, nil
}

// ---------------------------
// Core logic
// ---------------------------

// sortPlayers: last name, then first; case-insensitive.
func sortPlayers(players []Player) {
	sort.Slice(players, func(i, j int) bool {
		li := strings.ToLower(players[i].Last)
		lj := strings.ToLower(players[j].Last)
		if li != lj {
			return li < lj
		}
		fi := strings.ToLower(players[i].First)
		fj := strings.ToLower(players[j].First)
		return fi < fj
	})
}

// totalPoints: sum of Points for valid rows.
func totalPoints(players []Player) int {
	sum := 0
	for _, p := range players {
		if !p.BadData {
			sum += p.Points
		}
	}
	return sum
}

// ---------------------------
// Reporting
// ---------------------------

const nameColWidth = 28 // enough for long last names

// printReport prints the header and aligned table.
func printReport(players []Player) {
	fmt.Printf("BASKETBALL TEAM REPORT --- %d PLAYERS FOUND IN FILE\n", len(players))
	fmt.Printf("TOTAL POINTS SCORED: %d\n\n", totalPoints(players))

	fmt.Printf("%-*s %6s %6s\n", nameColWidth, "PLAYER NAME", "eFG%", "TS%")
	fmt.Println(strings.Repeat("-", nameColWidth+1+6+1+6))

	for _, p := range players {
		name := fmt.Sprintf("%s, %s", p.Last, p.First)
		if p.BadData {
			fmt.Printf("%-*s %s\n", nameColWidth, name, "*missing input data*")
			continue
		}
		fmt.Printf("%-*s %6.1f %6.1f\n", nameColWidth, name, p.eFG()*100, p.ts()*100)
	}
}

// ---------------------------
// main
// ---------------------------

func main() {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Enter input filename: ")
	fn, _ := reader.ReadString('\n')
	fn = strings.TrimSpace(fn)

	players, err := readPlayersFromFile(fn)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	sortPlayers(players)
	printReport(players)
}
