package main

import (
	"math/rand"
)

type Status int

const (
	CLOSE Status = iota
	OPEN
	FLAG
)

type GameStatus int

const (
	RUNNING GameStatus = iota
	WIN
	LOSE
)

type Coord struct {
	Row int
	Col int
}

type Cell struct {
	Row     int
	Col     int
	Status  Status
	Content int
}

type Game struct {
	Rows       int
	Cols       int
	Cells      [][]Cell
	Bombs      int
	GameStatus GameStatus
}

func (g *Game) updateGameStatus() {
	status := WIN
	for r := 0; r < g.Rows; r++ {
		for c := 0; c < g.Rows; c++ {
			cell := g.Cells[r][c]
			// bomb and open: lose
			if cell.Content == 9 && cell.Status == OPEN {
				g.GameStatus = LOSE
				return
			}
			// not bomb and close: running
			if cell.Content != 9 && cell.Status == CLOSE {
				status = RUNNING
			}
		}
	}
	g.GameStatus = status
}

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

func (g *Game) Flag(row, col int) {
	if row < 0 || row >= g.Rows || col < 0 || col >= g.Cols {
		return
	}
	switch g.Cells[row][col].Status {
	case CLOSE:
		g.Cells[row][col].Status = FLAG
	case FLAG:
		g.Cells[row][col].Status = CLOSE
	default:
		return
	}
}

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

func (g *Game) GetRemainingBombsCount() int {
	count := g.Bombs
	for _, cellsRow := range g.Cells {
		for _, cell := range cellsRow {
			if cell.Status == FLAG {
				count--
			}
		}
	}
	return count
}

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

func NewGame(rows, cols, nBombs int) *Game {
	game := &Game{
		Rows:       rows,
		Cols:       cols,
		Bombs:      nBombs,
		Cells:      [][]Cell{},
		GameStatus: RUNNING,
	}
	bombs := generateBombs(rows, cols, nBombs)
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
				content = getNumberOfNearBombs(r, c, bombs)
			}
			row = append(row, Cell{r, c, CLOSE, content})
		}
		game.Cells = append(game.Cells, row)
	}
	return game
}

func generateBombs(rows, cols, nBombs int) []Coord {
	bombsMap := make(map[Coord]struct{})
	for i := 0; i < nBombs; i++ {
		r := rand.Intn(rows)
		c := rand.Intn(cols)
		coord := Coord{r, c}
		_, ok := bombsMap[coord]
		if !ok {
			bombsMap[coord] = struct{}{}
		} else {
			// repeat the cycle because the coord is already taken
			i--
		}
	}
	coords := make([]Coord, 0)
	for coord := range bombsMap {
		coords = append(coords, coord)
	}
	return coords
}

func getNumberOfNearBombs(row, col int, bombs []Coord) int {
	nBombs := 0
	for _, bomb := range bombs {
		if ((row-bomb.Row <= 1) && (row-bomb.Row >= -1)) && ((col-bomb.Col <= 1) && (col-bomb.Col >= -1)) {
			nBombs++
		}
	}
	return nBombs
}
