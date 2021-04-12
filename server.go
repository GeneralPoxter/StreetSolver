package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math"
	"math/rand"
	"net/http"
	"os"
	"time"
)

type Loc struct {
	Lat float64 `json:"lat"`
	Lng float64 `json:"lng"`
}

var location Loc
var MarkerLocation Loc

func main() {
	fileServer := http.FileServer(http.Dir("./public"))
	http.Handle("/", fileServer)
	http.HandleFunc("/getLoc", getLoc)
	http.HandleFunc("/getScore", getScore)
	http.HandleFunc("/sendScore", sendScore)

	if err := http.ListenAndServe(getPort(), nil); err != nil {
		log.Fatal(err)
	}
}

// Port function for Heroku deployment
func getPort() string {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
		fmt.Println("INFO: No PORT environment variable detected, defaulting to " + port)
	}
	fmt.Println("Starting server at port " + port)
	return ":" + port
}

func getLoc(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/getLoc" {
		http.Error(w, "404 not found.", http.StatusNotFound)
		return
	}

	if r.Method != "GET" {
		http.Error(w, "Method is not supported.", http.StatusNotFound)
		return
	}

	rand.Seed(time.Now().UnixNano())
	loc := Loc{randRange(-80, 80), randRange(-180, 180)}
	location = loc
	fmt.Println(loc)
	js, err := json.Marshal(loc)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}

func getScore(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/getScore" {
		http.Error(w, "404 not found.", http.StatusNotFound)
		return
	}

	if r.Method != "GET" {
		http.Error(w, "Method is not supported.", http.StatusNotFound)
		return
	}

	//var markerLoc Loc
	//markerLoc.Lat = 0
	//markerLoc.Lng = 0
	js, err := json.Marshal(calculateScore(MarkerLocation, location))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}

func sendScore(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/sendScore" {
		http.Error(w, "404 not found.", http.StatusNotFound)
		return
	}
	if r.Method != "POST" {
		http.Error(w, "Method is not supported.", http.StatusNotFound)
		return
	}
	decoder := json.NewDecoder(r.Body)
	var t Loc
	err := decoder.Decode(&t)
	if err != nil {
		panic(err)
	}
	log.Println(t)
	MarkerLocation.Lat = t.Lat
	MarkerLocation.Lng = t.Lng
}

func randRange(min, max float64) float64 {
	return rand.Float64()*(max-min) + min
}

func calculateScore(loc1, loc2 Loc) int {
	latRad1 := loc1.Lat * math.Pi / 180
	latRad2 := loc2.Lat * math.Pi / 180
	deltaLat := (loc2.Lat - loc1.Lat) * math.Pi / 180
	deltaLng := (loc2.Lng - loc1.Lng) * math.Pi / 180
	a := math.Pow(math.Sin(deltaLat/2), 2) + math.Cos(latRad1)*math.Cos(latRad2)*math.Pow(math.Sin(deltaLng/2), 2)
	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))

	distance := 6371e3 * c / 1000
	score := int(math.Round(5000 * math.Pow(math.E, (-distance/2000))))

	return score
}
