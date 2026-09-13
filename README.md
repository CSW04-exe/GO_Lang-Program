# GO_Lang-Program

A command-line Go program that reads basketball player statistics from a
text file, computes advanced shooting metrics for each player, and prints
a clean, sorted performance report to the terminal. Written by Carter Ward
for CS424-01 (Fall 2025).

## 1. Purpose

This project exists to answer a simple sports-analytics question with real
code: given a list of players and their raw box-score totals (games,
points, field goals, three-pointers, free throws), how do you turn that
into the advanced shooting metrics coaches and analysts actually care
about? Raw points-per-game numbers can be misleading — a player who
shoots a lot of low-percentage shots can look productive on the surface
while actually hurting their team's efficiency. Effective Field Goal
Percentage (eFG%) and True Shooting Percentage (TS%) correct for that by
accounting for the extra value of three-pointers and the free-throw line.

Beyond the basketball angle, the program was built as a course assignment
for CS424 to demonstrate practical Go skills: reading and parsing
untrusted text input, structuring data with custom types, handling
errors and missing data gracefully, and producing formatted, readable
output — the same fundamentals that show up in almost any real-world
data-processing tool.

## 2. Problem and approach

The assignment was to build a program that ingests a roster of players
from a file (format and contents not guaranteed to be clean), computes
per-player advanced stats, and reports them in an organized way — while
being resilient to bad or incomplete rows rather than crashing on them.

The approach taken:

- **Read the raw file line by line** using `bufio.Scanner`, trimming
  whitespace and skipping blank lines so the parser only has to deal with
  real data rows.
- **Parse each line defensively.** Each row is expected to contain a
  first name, last name, and eight integer stats (games, points, field
  goals made/attempted, three-pointers made/attempted, free throws
  made/attempted). If a row has too few fields, or a field that isn't a
  valid integer, the player is flagged with a `BadData` boolean instead
  of aborting the whole program — so one malformed line in the input file
  doesn't take down the report for everyone else.
- **Compute the metrics as methods on the data type** rather than as
  free-floating functions operating on raw numbers, so the math for eFG%
  and TS% lives right next to the data it operates on and is protected
  from divide-by-zero (a player with zero field-goal attempts, for
  example, safely returns 0 instead of crashing).
- **Sort players alphabetically** by last name (falling back to first
  name for ties), case-insensitively, so the report reads like a roster
  rather than an unordered dump.
- **Report bad rows visibly** instead of silently dropping them, printing
  `*missing input data*` in place of stats so the reader knows exactly
  which entries in the source file need to be fixed.

## 3. Structure and methodologies

The program is a single, self-contained Go file (`WardCS424GOLang.go`)
with no external dependencies and no `go.mod` module file — everything it
needs comes from the Go standard library, keeping the project easy to
build and run with nothing more than a Go toolchain (developed and tested
against Go 1.25 in VS Code):

- **`bufio`** — for buffered, line-oriented reading of both the input
  file (`bufio.Scanner`) and interactive stdin input (`bufio.NewReader`)
  for the filename prompt.
- **`os`** — for opening the input file and writing errors/exit codes to
  stderr.
- **`sort`** — specifically `sort.Slice`, to order the parsed players by
  last/first name with a custom comparator.
- **`strconv`** — to convert the text stat fields into integers, with
  error handling for malformed numeric input.
- **`strings`** — for whitespace trimming, case-insensitive comparisons,
  splitting each line into fields, and building the repeated-dash table
  divider.
- **`fmt`** — for all formatted console output, including fixed-width
  columns via `%-*s` and one-decimal percentages via `%6.1f`.

Data-wise, the whole program revolves around one struct, `Player`, which
holds the raw counting stats plus a `BadData` flag, and two methods,
`eFG()` and `ts()`, that derive the advanced percentages from those raw
totals on demand rather than storing precomputed values. The code is
organized into clearly commented sections — data model, parsing/file I/O,
core logic (sorting and totals), and reporting — which keeps each concern
isolated even though it's all in one file.

## 4. Process

Based on the file's own header comments and structure, the build followed
a fairly standard "define the data, then the pipeline" process typical of
a course assignment with a tight, well-specified deliverable:

1. **Model the data first.** The `Player` struct was defined up front with
   exactly the fields the assignment's stat formulas require (games,
   points, field goal/three-point/free-throw made-and-attempted pairs),
   plus a `BadData` flag anticipated from the start rather than bolted on
   later — a sign the "what happens with messy input" question was
   considered early rather than as an afterthought.
2. **Implement the math as methods.** `eFG()` and `ts()` were written
   directly against the standard formulas (eFG% = (FGM + 0.5×3PM) / FGA;
   TS% = PTS / (2×(FGA + 0.44×FTA))), each guarded against a zero
   denominator so a player with no attempts doesn't crash the report.
3. **Build the parser and file reader.** `parsePlayer` turns one line of
   text into a `Player`, and `readPlayersFromFile` wraps it in a
   `bufio.Scanner` loop, skipping blanks and propagating real I/O errors
   up to `main`. Validation (field count, integer parsing) was folded
   into the same function so bad rows are caught at the earliest possible
   point.
4. **Add sorting and aggregation.** `sortPlayers` (case-insensitive,
   last-then-first) and `totalPoints` (summing only valid rows) were
   layered on top of the already-parsed slice of players — logic kept
   separate from both parsing and display.
5. **Build the report last.** `printReport` formats a header, a total
   points line, and an aligned table, using constant column widths and
   `fmt`'s width/precision verbs so the output lines up regardless of
   name length, printing a clear placeholder for any flagged bad-data
   row instead of a formatting error or a zeroed-out line.
6. **Wire it together in `main`.** The final step was the thin driver:
   prompt for a filename on stdin, read and parse the file, exit cleanly
   with a stderr message on I/O failure, then sort and print. Interactive
   input was chosen over a command-line argument, keeping the program
   simple to run for anyone testing it without needing to remember flags.

The heavy inline commenting throughout (including a full header block
documenting author, course, date, and environment) suggests the process
was documentation-as-you-go rather than a final documentation pass,
which is consistent with how the file reads today: each section explains
its own intent before the code that implements it.

## 5. Outcome

The finished program successfully takes an arbitrary player-stats file —
including one with incomplete or malformed rows — and produces a
correctly sorted, aligned, human-readable report with accurate eFG% and
TS% calculations, without crashing on bad input. That resilience (a
malformed line degrades gracefully into a flagged row instead of an
unhandled exception) is the clearest concrete proof the program works as
intended, since it's the exact failure mode the design specifically
guards against.

Working through this assignment reinforced several practical Go and
software-design skills:

- **Defensive parsing of untrusted input** — treating every line from a
  file as something that might be malformed, and deciding up front what
  "gracefully handle it" actually means in code (a flag and a fallback
  message, not a crash).
- **Struct-and-methods design** — attaching the domain math (`eFG()`,
  `ts()`) directly to the data type it describes, which keeps the
  business logic close to the data and easy to test or extend later
  (e.g., adding another advanced stat is just one more method).
- **Idiomatic use of the Go standard library** for I/O and text
  processing (`bufio.Scanner`, `strings.Fields`, `strconv.Atoi`,
  `sort.Slice`) without reaching for any third-party dependencies —
  proof that the standard library alone is enough for a real, useful
  command-line tool.
- **Formatted output as a first-class concern** — using `fmt`'s width and
  precision verbs to produce an aligned table, rather than treating
  console output as an afterthought.

More broadly, the project demonstrates the ability to take a
loosely-specified, real-world-flavored problem (messy sports data) and
turn it into a small, well-organized program with clear separation
between data modeling, parsing, computation, and presentation — a
structure that would scale reasonably well if the assignment grew to
support more stats, more input formats, or a larger roster.
