package main

import (
	"fmt"
	"html/template"
	"net/http"
	"strconv"
)

var game Game

type GameController struct {
	Game Game
}

func (g *GameController) serve(w http.ResponseWriter, r *http.Request) {
	tmplFile := "templates/index.tmpl"
	tmpl, err := template.ParseFiles(tmplFile)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
	err = tmpl.Execute(w, nil)
}

func (g *GameController) game(w http.ResponseWriter, r *http.Request) {
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

func (g *GameController) easy(w http.ResponseWriter, r *http.Request) {
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
func (g *GameController) medium(w http.ResponseWriter, r *http.Request) {
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
func (g *GameController) hard(w http.ResponseWriter, r *http.Request) {
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

func (g *GameController) reset(w http.ResponseWriter, r *http.Request) {
	tmplFile := "templates/game.tmpl"
	tmpl, err := template.ParseFiles(tmplFile)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
	g.Game = *NewGame(g.Game.Rows, g.Game.Cols, g.Game.Bombs)
	err = tmpl.Execute(w, g.Game)
}

func (g *GameController) click(w http.ResponseWriter, r *http.Request) {
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

func (g *GameController) flag(w http.ResponseWriter, r *http.Request) {
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

func (g *GameController) chord(w http.ResponseWriter, r *http.Request) {
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

func (g *GameController) count(w http.ResponseWriter, r *http.Request) {
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

func (g *GameController) status(w http.ResponseWriter, r *http.Request) {
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

func (g *GameController) clickInstructions(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("HX-Trigger", "instructions")
}

func (g *GameController) instructions(w http.ResponseWriter, r *http.Request) {
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

func main() {
	gameController := GameController{
		Game: *NewGame(9, 9, 10),
	}
	http.HandleFunc("/", gameController.serve)
	http.HandleFunc("/game", gameController.game)
	http.HandleFunc("/click", gameController.click)
	http.HandleFunc("/flag", gameController.flag)
	http.HandleFunc("/chord", gameController.chord)
	http.HandleFunc("/reset", gameController.reset)
	http.HandleFunc("/count", gameController.count)
	http.HandleFunc("/status", gameController.status)
	http.HandleFunc("/instructions", gameController.instructions)
	http.HandleFunc("/click-instructions", gameController.clickInstructions)
	http.HandleFunc("/difficulty/easy", gameController.easy)
	http.HandleFunc("/difficulty/medium", gameController.medium)
	http.HandleFunc("/difficulty/hard", gameController.hard)
	fmt.Println("Listening on port 3000")
	http.ListenAndServe(":3000", nil)
}
