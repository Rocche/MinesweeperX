package main

import (
	"fmt"
	"html/template"
	"net/http"
	"strconv"
)

// GameController is the component responsible of managing the
// client's requests and serving the components based on the Game status
type GameController struct {
	Game Game
}

// Home serves the initial view
func (g *GameController) Home(w http.ResponseWriter, r *http.Request) {
	tmplFile := "templates/index.tmpl"
	tmpl, err := template.ParseFiles(tmplFile)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
	err = tmpl.Execute(w, nil)
}

// StartGame serves the game view with the selected difficulty
func (g *GameController) StartGame(w http.ResponseWriter, r *http.Request) {
	tmplFile := "templates/game.tmpl"
	gridFile := "templates/game-grid.tmpl"
	counterFile := "templates/bombs-counter.tmpl"
	statusFile := "templates/game-status.tmpl"
	instructionsFile := "templates/instructions.tmpl"
	tmpl, err := template.ParseFiles(tmplFile, gridFile, counterFile, instructionsFile, statusFile)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
	err = tmpl.Execute(w, g.Game)
}

// DifficultyEasy serves the difficulty easy component
func (g *GameController) DifficultyEasy(w http.ResponseWriter, r *http.Request) {
	tmplFile := "templates/difficulty-easy.tmpl"
	logoFile := "templates/logo.tmpl"
	btnFile := "templates/start-button.tmpl"
	tmpl, err := template.ParseFiles(tmplFile, logoFile, btnFile)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
	err = tmpl.Execute(w, nil)
	g.Game = *NewGame(9, 9, 10)
}

// DifficultyMedium serves the difficulty medium component
func (g *GameController) DifficultyMedium(w http.ResponseWriter, r *http.Request) {
	tmplFile := "templates/difficulty-medium.tmpl"
	logoFile := "templates/logo.tmpl"
	btnFile := "templates/start-button.tmpl"
	tmpl, err := template.ParseFiles(tmplFile, logoFile, btnFile)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
	err = tmpl.Execute(w, nil)
	g.Game = *NewGame(16, 16, 40)
}

// DifficultyHard serves the difficulty hard component
func (g *GameController) DifficultyHard(w http.ResponseWriter, r *http.Request) {
	tmplFile := "templates/difficulty-hard.tmpl"
	logoFile := "templates/logo.tmpl"
	btnFile := "templates/start-button.tmpl"
	tmpl, err := template.ParseFiles(tmplFile, logoFile, btnFile)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
	err = tmpl.Execute(w, nil)
	g.Game = *NewGame(16, 30, 99)
}

// ClickCell updates the game after a cell click and serves the new
// game grid
func (g *GameController) ClickCell(w http.ResponseWriter, r *http.Request) {
	row, err := strconv.Atoi(r.URL.Query().Get("row"))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, "Cannot parse row")
	}
	col, err := strconv.Atoi(r.URL.Query().Get("col"))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, "Cannot parse col")
	}
	tmplFile := "templates/game-grid.tmpl"
	tmpl, err := template.ParseFiles(tmplFile)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	g.Game.OpenCell(row, col, true)
	if g.Game.GameStatus != RUNNING {
		w.Header().Add("HX-Trigger", "gameover")
	}
	err = tmpl.ExecuteTemplate(w, "gameGrid", g.Game)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
}

// FlagCell updates the game after a cell right click and serves the new
// game grid
func (g *GameController) FlagCell(w http.ResponseWriter, r *http.Request) {
	row, err := strconv.Atoi(r.URL.Query().Get("row"))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, "Cannot parse row")
	}
	col, err := strconv.Atoi(r.URL.Query().Get("col"))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, "Cannot parse col")
	}
	tmplFile := "templates/game-grid.tmpl"
	tmpl, err := template.ParseFiles(tmplFile, tmplFile)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
	g.Game.Flag(row, col)
	w.Header().Add("HX-Trigger", "flagged") // needed to trigger bombs counter event
	if g.Game.GameStatus != RUNNING {
		w.Header().Add("HX-Trigger", "gameover")
	}
	err = tmpl.ExecuteTemplate(w, "gameGrid", g.Game)
}

// ChordCell updates the game after a cell ctrl+click and serves the new
// game grid
func (g *GameController) ChordCell(w http.ResponseWriter, r *http.Request) {
	row, err := strconv.Atoi(r.URL.Query().Get("row"))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, "Cannot parse row")
	}
	col, err := strconv.Atoi(r.URL.Query().Get("col"))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, "Cannot parse col")
	}
	tmplFile := "templates/game-grid.tmpl"
	tmpl, err := template.ParseFiles(tmplFile)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
	g.Game.Chord(row, col)
	if g.Game.GameStatus != RUNNING {
		w.Header().Add("HX-Trigger", "gameover")
	}
	err = tmpl.ExecuteTemplate(w, "gameGrid", g.Game)
}

// MinesCounter serves the mine counter component with the actual
// mines count
func (g *GameController) MinesCounter(w http.ResponseWriter, r *http.Request) {
	counterFile := "templates/bombs-counter.tmpl"
	tmpl, err := template.ParseFiles(counterFile)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
	count := g.Game.GetRemainingBombsCount()
	err = tmpl.ExecuteTemplate(w, "bombsCounter", struct{ Bombs int }{count})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
}

// GameStatus serves the victory condition of the game
func (g *GameController) GameStatus(w http.ResponseWriter, r *http.Request) {
	counterFile := "templates/game-status.tmpl"
	tmpl, err := template.ParseFiles(counterFile)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
	err = tmpl.ExecuteTemplate(w, "gameStatus", g.Game)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
}

// InstructionsClicked sends a HTMX HX-Trigger header in order to let
// the view fetch the instructions component
func (g *GameController) InstructionsClicked(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("HX-Trigger", "instructions")
}

// Instructions serves the instructions view after the instructions HX-Trigger
func (g *GameController) Instructions(w http.ResponseWriter, r *http.Request) {
	counterFile := "templates/instructions.tmpl"
	tmpl, err := template.ParseFiles(counterFile)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
	err = tmpl.ExecuteTemplate(w, "instructions", g.Game)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
}
