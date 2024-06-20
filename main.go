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
}

func main() {
	InitController()
	fmt.Println("Listening on port 3000")
	http.ListenAndServe(":3000", nil)
}
