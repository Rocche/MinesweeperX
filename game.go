package main

import (
	"math/rand"
)

// Status is the status of a cell
type Status int

const (
	CLOSE Status = iota
	OPEN
	FLAG
)

// GameStatus is the status of the game
type GameStatus int

const (
	RUNNING GameStatus = iota
	WIN
	LOSE
)

// Coord is a struct representing the coordinates in the game grid
type Coord struct {
	Row int
	Col int
}

// Cell is a struct defining a cell in the game grid
type Cell struct {
	Row     int
	Col     int
	Status  Status
	Content int
}

// Game is the struct containing all the info needed to render the game
type Game struct {
	Rows           int
	Cols           int
	Cells          [][]Cell
	Mines          int
	RemainingMines int
	GameStatus     GameStatus
}

// updateGameStatus checks all cells in order to decide the game status
func (g *Game) updateGameStatus() {
	status := WIN
	for r := 0; r < g.Rows; r++ {
		for c := 0; c < g.Rows; c++ {
			cell := g.Cells[r][c]
			// mine and open: lose
			if cell.Content == 9 && cell.Status == OPEN {
				g.GameStatus = LOSE
				return
			}
			// not mine and close: running
			if cell.Content != 9 && cell.Status == CLOSE {
				status = RUNNING
			}
		}
	}
	g.GameStatus = status
}

// reveal sets the cell status to [OPEN] and if set reveals the
// adjacent cells in case the cell is empty in a recursive manner
func (g *Game) reveal(row, col int, propagateIfEmpty bool) {
	if row < 0 || row >= g.Rows || col < 0 || col >= g.Cols {
		return
	}
	cell := &g.Cells[row][col]
	if cell.Status != CLOSE {
		return
	}
	cell.Status = OPEN
	if !propagateIfEmpty || cell.Content != 0 {
		return
	}
	g.reveal(row-1, col-1, propagateIfEmpty)
	g.reveal(row-1, col, propagateIfEmpty)
	g.reveal(row-1, col+1, propagateIfEmpty)
	g.reveal(row, col-1, propagateIfEmpty)
	g.reveal(row, col+1, propagateIfEmpty)
	g.reveal(row+1, col-1, propagateIfEmpty)
	g.reveal(row+1, col, propagateIfEmpty)
	g.reveal(row+1, col+1, propagateIfEmpty)
}

// OpenCell reveals the clicked cell
func (g *Game) OpenCell(row, col int, propagateIfEmpty bool) {
	if row < 0 || row >= g.Rows || col < 0 || col >= g.Cols {
		return
	}
	if g.Cells[row][col].Status != CLOSE {
		return
	}
	// check if bomb
	if g.Cells[row][col].Content == 9 {
		// reveal all bombs
		for r := 0; r < g.Rows; r++ {
			for c := 0; c < g.Cols; c++ {
				if g.Cells[r][c].Content == 9 {
					g.Cells[r][c].Status = OPEN
				}
			}
		}
		g.updateGameStatus()
		return
	}
	g.reveal(row, col, propagateIfEmpty)
	g.updateGameStatus()
}

// Flag flags the given cell
func (g *Game) Flag(row, col int) {
	if row < 0 || row >= g.Rows || col < 0 || col >= g.Cols {
		return
	}
	switch g.Cells[row][col].Status {
	case CLOSE:
		g.Cells[row][col].Status = FLAG
		g.RemainingMines--
	case FLAG:
		g.Cells[row][col].Status = CLOSE
		g.RemainingMines++
	default:
		return
	}
}

// Chord reveals the adjacent cells to an already open one if
// the content of the cells is equal ato the adjacent flagged cells
func (g *Game) Chord(row, col int) {
	if row < 0 || row >= g.Rows || col < 0 || col >= g.Cols {
		return
	}
	// check if cell is actually opened
	if g.Cells[row][col].Status != OPEN {
		return
	}
	// count the number of adjacent flagged cells
	adjacentFlaggedCells := g.getNumberOfAdjacentFlaggedCells(row, col)
	// can open adjacent cells only if flagged cells are the same
	// as the cell content
	if g.Cells[row][col].Content != adjacentFlaggedCells {
		return
	}
	for r := row - 1; r <= row+1; r++ {
		for c := col - 1; c <= col+1; c++ {
			if r >= 0 && r < g.Rows && c >= 0 && c < g.Cols && g.Cells[r][c].Status != FLAG {
				g.OpenCell(r, c, false)
			}
		}
	}
	g.updateGameStatus()
}

// GetRemainingMinesCount returns the number of mines remaining,
// which is the initial number of mines minus the flagged cells
func (g *Game) GetRemainingMinesCount() int {
	return g.RemainingMines
}

// getNumberOfAdjacentFlaggedCells gets how many cells adjacent to the
// given one are flagged
func (g *Game) getNumberOfAdjacentFlaggedCells(row, col int) int {
	count := 0
	for r := row - 1; r <= row+1; r++ {
		for c := col - 1; c <= col+1; c++ {
			if r < 0 || r >= g.Rows || c < 0 || c >= g.Cols {
				continue
			}
			if g.Cells[r][c].Status == FLAG {
				count++
			}
		}
	}
	return count
}

// NewGame initializes a new [Game] with given rows, columns and number of mines
func NewGame(rows, cols, nMines int) *Game {
	game := &Game{
		Rows:           rows,
		Cols:           cols,
		Mines:          nMines,
		RemainingMines: nMines,
		Cells:          [][]Cell{},
		GameStatus:     RUNNING,
	}
	bombs := generateMines(rows, cols, nMines)
	for r := 0; r < rows; r++ {
		row := make([]Cell, 0)
		for c := 0; c < cols; c++ {
			content := 0
			isBomb := false
			for _, bomb := range bombs {
				if r == bomb.Row && c == bomb.Col {
					isBomb = true
					break
				}
			}
			if isBomb {
				content = 9
			} else {
				content = getNumberOfNearMines(r, c, bombs)
			}
			row = append(row, Cell{r, c, CLOSE, content})
		}
		game.Cells = append(game.Cells, row)
	}
	return game
}

// generateMines returns a slice of [Coord] representing where the mines
// should be placed
func generateMines(rows, cols, nMines int) []Coord {
	// first we populate a map of [Coord] in order
	// to check that the random generated coordinates are unique
	minesMap := make(map[Coord]struct{})
	for i := 0; i < nMines; i++ {
		r := rand.Intn(rows)
		c := rand.Intn(cols)
		coord := Coord{r, c}
		_, ok := minesMap[coord]
		if !ok {
			minesMap[coord] = struct{}{}
		} else {
			// repeat the cycle because the coord is already taken
			i--
		}
	}
	// finally we transform the map into a slice
	coords := make([]Coord, 0)
	for coord := range minesMap {
		coords = append(coords, coord)
	}
	return coords
}

// getNumberOfNearMines gets the number of near mines to a given cell
func getNumberOfNearMines(row, col int, mines []Coord) int {
	nMines := 0
	for _, mine := range mines {
		if ((row-mine.Row <= 1) && (row-mine.Row >= -1)) && ((col-mine.Col <= 1) && (col-mine.Col >= -1)) {
			nMines++
		}
	}
	return nMines
}
