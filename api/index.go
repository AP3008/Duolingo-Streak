package api

import (
	_ "embed"
	"fmt"
	"net/http"
	"strings"

	"duolingo-api/duolingo"
)

func Handler(w http.ResponseWriter, r *http.Request) {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /user/{username}", handleUserStreak)
	mux.HandleFunc("GET /user/svg/{username}", handleSVG)

	mux.ServeHTTP(w, r)
}

//go:embed duolingo-bird.svg
var birdSVG string

func handleUserStreak(w http.ResponseWriter, r *http.Request) {
	username := r.PathValue("username")

	streak, err := duolingo.GetStreak(username)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-type", "text/plain")
	fmt.Fprint(w, streak)
}

func handleSVG(w http.ResponseWriter, r *http.Request) {
	username := r.PathValue("username")

	streak, err := duolingo.GetStreak(username)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	cleanBird := strings.Replace(birdSVG, `<?xml version="1.0" encoding="UTF-8"?>`, "", 1)

	svg := fmt.Sprintf(`
	<svg width="300" height="120" viewBox="0 0 300 120" fill="none" xmlns="http://www.w3.org/2000/svg">
		<style>
			.streak-text { font: 800 35px 'Segoe UI', Arial, sans-serif; fill: #FF9600; }
			.days-label { font: 600 15px 'Segoe UI', Arial, sans-serif; fill: #c9d1d9; }
			.circle-border { stroke: #FF9600; stroke-width: 4; fill: none; }
			@keyframes slideIn { from { transform: translateX(-10px); opacity: 0; } to { transform: translateX(0); opacity: 1; } }
			.animate { animation: slideIn 0.5s ease-out forwards; }
		</style>

		<g transform="translate(0, 10) scale(0.6)" class="animate">
			%s
		</g>

		<g transform="translate(230, 60)" class="animate">
			<circle r="45" class="circle-border" />
			<text x="0" y="5" text-anchor="middle" class="streak-text">%d</text>
			<text x="0" y="25" text-anchor="middle" class="days-label">DAYS</text>
		</g>
	</svg>`, cleanBird, streak)

	w.Header().Set("Content-Type", "image/svg+xml")
	w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
	fmt.Fprint(w, svg)
}
