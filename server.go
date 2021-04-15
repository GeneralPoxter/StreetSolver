package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math"
	"math/rand"
	"net/http"
	"os"
	"strconv"
	"time"
)

type Loc struct {
	Lat float64 `json:"lat"`
	Lng float64 `json:"lng"`
}

type Polygon []Loc

var target Loc
var marker Loc

//maryland
var poly = []Loc{{39.72108607946068, -79.47666224600735},
	{39.72322905639499, -75.78820613062912},
	{38.46024225250412, -75.69355492964675},
	{38.45107088128103, -75.04934375153289},
	{38.02838839691087, -75.24220525055418},
	{38.40467596253332, -77.04430644149704},
	{38.90229143401203, -76.90148419419953},
	{39.00054556187369, -77.04155988235236},
	{39.22404469392229, -77.45286054814271},
	{39.699283697533666, -77.93533222975897},
	{39.69542029957551, -78.18507807210698},
	{39.64867115193669, -78.76619476351625},
	{39.20622885088143, -79.48645056628894}}

func main() {
	fileServer := http.FileServer(http.Dir("./public"))
	http.Handle("/", fileServer)
	http.HandleFunc("/getLoc", getLoc)
	http.HandleFunc("/getScore", getScore)

	if err := http.ListenAndServe(getPort(), nil); err != nil {
		log.Fatal(err)
	}
}

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

	boundaryCheck := true
	for boundaryCheck {
		target = Loc{randRange(38, 40), randRange(-78, -75)}
		if isLocInPoly(poly, target) {
			boundaryCheck = false
		}
	}

	js, err := json.Marshal(target)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}

func randRange(min, max float64) float64 {
	rand.Seed(time.Now().UnixNano())
	return rand.Float64()*(max-min) + min
}

func isLocInPoly(poly Polygon, location Loc) bool {
	var inside bool
	i, j := 0, len(poly)-1
	for ; i < len(poly); j, i = i, i+1 {
		latI, lngI := poly[i].Lat, poly[i].Lng
		latJ, lngJ := poly[j].Lat, poly[j].Lng

		intersect := ((lngI > location.Lng) != (lngJ > location.Lng)) && (location.Lat < (latJ-latI)*(location.Lng-lngI)/(lngJ-lngI)+latI)
		if intersect {
			inside = !inside
		}
	}

	return inside
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

	query := r.URL.Query()
	values := make(map[string]float64)
	for key, _ := range query {
		s, err := strconv.ParseFloat(query.Get(key), 64)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		values[key] = s
	}

	target = Loc{Lat: values["targetLat"], Lng: values["targetLng"]}
	marker = Loc{Lat: values["markerLat"], Lng: values["markerLng"]}

	score := strconv.Itoa(calculateScore(marker, target))
	fmt.Fprint(w, score)
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
