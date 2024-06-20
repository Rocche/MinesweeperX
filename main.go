package main

import (
	"fmt"
	"net/http"
)

var game Game

func InitController() {
	gameController := GameController{
		Game: *NewGame(9, 9, 10),
	}
	http.HandleFunc("/", gameController.Home)
	http.HandleFunc("/game", gameController.StartGame)
	http.HandleFunc("/click", gameController.ClickCell)
	http.HandleFunc("/flag", gameController.FlagCell)
	http.HandleFunc("/chord", gameController.ChordCell)
	http.HandleFunc("/count", gameController.MinesCounter)
	http.HandleFunc("/status", gameController.GameStatus)
	http.HandleFunc("/instructions", gameController.Instructions)
	http.HandleFunc("/click-instructions", gameController.InstructionsClicked)
	http.HandleFunc("/difficulty/easy", gameController.DifficultyEasy)
	http.HandleFunc("/difficulty/medium", gameController.DifficultyMedium)
	http.HandleFunc("/difficulty/hard", gameController.DifficultyHard)
}

func main() {
	InitController()
	fmt.Println("Listening on port 3000")
	http.ListenAndServe(":3000", nil)
}
