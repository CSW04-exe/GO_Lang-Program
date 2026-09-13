# GO_Lang-Program

## Purpose

This is my assignment for CS424-01 (Fall 2025) at UAH. The course asked me to pick up a language I hadn't used for real work before and build something non-trivial in it, so I chose Go. I wanted a project with actual substance rather than a toy exercise, so I built a small basketball statistics reporter: it reads player data from a file, computes a couple of standard advanced shooting metrics, and prints a formatted report. The point wasn't the basketball math — it was forcing myself to learn Go's type system, its approach to structs and methods, its error-handling style, and its standard library (file I/O, scanning, sorting, string parsing) well enough to use them correctly in one sitting.

## Problem and Approach

The assignment needed a program that:

- Reads a list of players and their raw box-score totals (games, points, field goals made/attempted, three-pointers made/attempted, free throws made/attempted) from a text file whose name is supplied at runtime.
- Computes two advanced shooting metrics for each player:
  - **eFG% (effective field goal percentage)**: `(FGMade + 0.5 * TPMade) / FGAtt`
  - **TS% (true shooting percentage)**: `Points / (2 * (FGAtt + 0.44 * FTAtt))`
- Handles rows with missing or malformed data instead of crashing on them.
- Sorts the players alphabetically (last name, then first name) before printing.
- Prints a clean, aligned report with a header, total points scored, and one row per player.

My approach in Go was to model each player as a `Player` struct holding the raw integer totals plus a `BadData` flag, attach `eFG()` and `ts()` as methods on that struct so the metric math lives next to the data it operates on, and keep the pipeline as a straight sequence of small, single-purpose functions: read the file into a slice of `Player`, sort the slice in place, then walk it to print the report. Rows that don't parse (too few fields, or a non-numeric stat) are still kept in the slice with `BadData: true` so they show up in the output as flagged rather than being silently dropped.

## Structure and Methodologies

**Language constructs used:**

- **Structs and methods** — `Player` is a struct with named fields for identity (`First`, `Last`), the eight raw counting stats, and a `BadData` bool. `eFG()` and `ts()` are value-receiver methods on `Player`, which was my first real exposure to Go's receiver-based method syntax instead of a class.
- **Slices** — players are collected into a `[]Player` that grows as the file is read (`append`), gets sorted in place, and is iterated twice (once to sum points, once to print).
- **Multiple return values and explicit error handling** — `readPlayersFromFile` returns `([]Player, error)` instead of throwing an exception, and `strconv.Atoi` is checked with the `v, err :=` pattern, which is different from the try/catch style I was used to.
- **Closures** — `sortPlayers` uses `sort.Slice` with an inline comparator closure that does a case-insensitive compare on last name, falling back to first name.
- **String formatting/parsing** — `strings.Fields` tokenizes each line, `strconv.Atoi` converts tokens to ints, and `fmt.Printf`/`fmt.Sprintf` with width specifiers (`%-*s`, `%6.1f`) build the aligned table.
- **Defer** — `defer f.Close()` releases the file handle regardless of how the function returns.
- **Named constants** — `nameColWidth` controls the column width used by both the header and the row formatting so they can't drift apart.

**Standard library packages imported:** `bufio` (buffered file/stdin reading), `fmt` (formatted I/O), `os` (file open, stdin, exit codes), `sort` (in-place slice sorting), `strconv` (string-to-int conversion), `strings` (tokenizing, trimming, case-folding, repeat).

No third-party packages and no goroutines/channels — the program is single-threaded and sequential, since the assignment's focus was Go's core language features and standard library rather than its concurrency model.

## Process

Step by step, from `main()`:

1. `main` opens a buffered reader on `os.Stdin` and prompts `Enter input filename:`, reading a line and trimming whitespace to get the filename.
2. It calls `readPlayersFromFile(fn)`, which opens the file, and scans it line by line with a `bufio.Scanner`, skipping blank lines. Each non-blank line is handed to `parsePlayer`.
3. `parsePlayer` splits the line into whitespace-separated fields. If there are fewer than 10 fields (first name, last name, and 8 stats), or if any of the 8 stat fields fails to convert to an integer, the function returns a `Player` with only `First`/`Last` set (when available) and `BadData: true`. Otherwise it returns a fully populated `Player`.
4. Once every line has been parsed into the `[]Player` slice, `main` checks for a file-level error (e.g., the file didn't exist) and exits with a message on stderr if one occurred.
5. `main` calls `sortPlayers`, which sorts the slice in place by last name, then first name, both case-insensitively.
6. `main` calls `printReport`, which:
   - Prints a header line with the total player count and the sum of all valid players' points (`totalPoints`, which skips `BadData` rows).
   - Prints a column header (`PLAYER NAME`, `eFG%`, `TS%`) and a dashed rule sized to match.
   - Loops over the sorted players: a `BadData` player prints as `Last, First` with a `*missing input data*` marker instead of numbers; a good player prints `Last, First` followed by `eFG()` and `ts()` (each computed on the fly and shown as a percentage with one decimal place).
7. The program returns after the report is printed — there's no loop back to prompt again.

## Outcome

Running the program produces a readable, aligned console report, for example:

```
BASKETBALL TEAM REPORT --- 5 PLAYERS FOUND IN FILE
TOTAL POINTS SCORED: 412

PLAYER NAME                    eFG%    TS%
------------------------------------------
Adams, Jordan                  54.2   58.1
Brooks, Taylor          *missing input data*
...
```

That output is directly checkable against the math by hand, which is what let me confirm the eFG%/TS% formulas and the bad-data handling were actually correct rather than just "didn't crash."

**How to run:**

```
go run WardCS424GOLang.go
```

Then, at the prompt, enter the path to a text file where each non-blank line has the form:
`FirstName LastName Games Points FGMade FGAtt TPMade TPAtt FTMade FTAtt`

**What I got out of it:** this was my first substantial Go program, and it made the differences from the languages I already knew concrete rather than abstract — explicit multi-value error returns instead of exceptions, methods attached to a type via a receiver instead of defined inside a class body, and a statically-typed struct as the natural unit of "a record" instead of a dictionary or object literal. None of those ideas were unfamiliar in the abstract, but writing them myself — and hitting the compiler's opinions about unused imports, zero values, and `:=` versus `=` — is what actually made them stick. I came away with a working feel for reading a Go standard-library doc page and using it correctly the same day, which is really the skill this assignment was testing: that I can pick up a new language's idioms quickly enough to ship something real in it under a deadline, not just follow a tutorial.
