# Go Language Program — Basketball Stats Reporter

**Type:** Individual project
**Contributor:** Carter Ward
**Course:** CS 424-01 (Programming Languages), Fall 2025
**Completed:** 09/28/2025

## Purpose
This is my assignment for CS 424-01, which asked me to pick up a language I hadn't used for real work before and build something non-trivial in it. I chose Go and wanted a project with actual substance rather than a toy exercise, so I built a small basketball statistics reporter. The goal wasn't the basketball math — it was learning Go's structs/methods, error-handling style, and standard library well enough to use correctly in one sitting.

## Problem and Approach
The program reads a list of players and their raw box-score totals (games, points, field goals, three-pointers, free throws) from a text file named at runtime, computes eFG% (effective field goal percentage) and TS% (true shooting percentage) for each player, handles rows with missing or malformed data without crashing, sorts players alphabetically by last name, then first, and prints an aligned report with a header and total points. Each player is modeled as a `Player` struct holding the raw stats plus a `BadData` flag, with `eFG()` and `ts()` as methods on it. The pipeline is a straight sequence: read the file into a slice of `Player`, sort in place, then print. Rows that fail to parse are kept with `BadData: true` and flagged in the output rather than dropped.

## Structure and Methodologies
- `Player` struct with identity fields, eight raw counting stats, and a `BadData` bool; `eFG()`/`ts()` as value-receiver methods
- Players collected into a growing `[]Player` slice, sorted in place with `sort.Slice` and a case-insensitive comparator closure
- Multiple return values and explicit `err` checks (e.g. `readPlayersFromFile` returns `([]Player, error)`) instead of exceptions
- `strings.Fields` tokenizes lines, `strconv.Atoi` parses ints, `fmt.Printf`/`Sprintf` with width specifiers build the aligned table
- `defer f.Close()` for file cleanup; a named constant controls shared column width
- Stdlib only: `bufio`, `fmt`, `os`, `sort`, `strconv`, `strings` — no third-party packages, no concurrency

## Process
1. Prompt for and read an input filename from stdin.
2. Open the file and scan it line by line, parsing each into a `Player`.
3. Flag lines with too few fields or non-numeric stats as `BadData` instead of failing.
4. Exit with an error message if the file itself can't be opened.
5. Sort the parsed players by last name, then first name.
6. Print a header (count, total points), then one row per player, flagging bad-data rows.

## Outcome
Running the program produces a readable, aligned console report, e.g. `Adams, Jordan   54.2   58.1` alongside a flagged `*missing input data*` row for bad input. That output was directly checkable by hand, which let me confirm the eFG%/TS% math and bad-data handling were actually correct. This was my first substantial Go program, and it made the language's differences concrete: multi-value error returns instead of exceptions, receiver-based methods instead of classes, and structs as the natural "record" type. I came away with a working feel for reading Go stdlib docs and using them correctly the same day.

**How to run:**
```
go run WardCS424GOLang.go
```
At the prompt, enter a file path where each line is: `FirstName LastName Games Points FGMade FGAtt TPMade TPAtt FTMade FTAtt`
