package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math"
	"math/rand"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

type Loc struct {
	Lat float64 `json:"lat"`
	Lng float64 `json:"lng"`
}

type Polygon []Loc

type GameData struct {
	Target     Loc    `json:"target"`
	Score      int    `json:"score"`
	TotalScore int    `json:"totalScore"`
	Round      int    `json:"round"`
	HighScore  int    `json:"highScore"`
	Region     string `json:"region"`
}

var game GameData

var distanceFactor float64

var defaultRegion = "United States"

var radii = map[string]int{
	"Maryland":      50,
	"United States": 100,
	"World":         1000,
}

var borders = map[string][]Loc{
	"Maryland": {
		{38, -78},
		{40, -75},
	},
	"United States": {
		{25, -125},
		{50, -60},
	},
	"World": {
		{-80, -180},
		{80, 180},
	},
}

var polys = map[string][]Loc{
	"Maryland": {
		{39.72108607946068, -79.47666224600735},
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
		{39.20622885088143, -79.48645056628894},
	},
	"United States": {
		{49.00501224324328, -123.33937857614913},
		{48.999349663558, -95.15417471469046},
		{46.45274297417163, -84.48370346783305},
		{41.532978006779395, -82.74786331314247},
		{45.866805771566554, -67.92181360743781},
		{30.50989864735576, -81.48443541811491},
		{25.403251595341164, -80.39678902598565},
		{30.486232716490345, -83.76959156032217},
		{29.027031987602395, -96.5357054240675},
		{32.1473984353643, -106.31353783259205},
		{32.90510566423088, -119.0363490137464},
	},
	"World": {},
}

func main() {
	fileServer := http.FileServer(http.Dir("./public"))
	http.Handle("/", fileServer)
	http.HandleFunc("/getData", getData)
	http.HandleFunc("/getKey", getKey)
	http.HandleFunc("/getLoc", getLoc)
	http.HandleFunc("/getRadius", getRadius)
	http.HandleFunc("/receiveTarget", receiveTarget)
	http.HandleFunc("/getRoundData", getRoundData)
	http.HandleFunc("/getHighScore", getHighScore)
	http.HandleFunc("/restart", restart)

	game = GameData{Loc{0, 0}, 0, 0, 0, 0, defaultRegion}

	if err := http.ListenAndServe(getPort(), nil); err != nil {
		log.Fatal(err)
	}
}

func getEnv(key string) string {
	if err := godotenv.Load(".env"); err != nil {
		fmt.Println("INFO: .env file not found")
		return ""
	}

	return os.Getenv(key)
}

func getPort() string {
	port := getEnv("PORT")
	if port == "" {
		port = "8080"
		fmt.Println("INFO: No PORT environment variable detected, defaulting to " + port)
	}
	fmt.Println("Starting server at port " + port)
	return ":" + port
}

func getData(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/getKey" {
		http.Error(w, "404 not found.", http.StatusNotFound)
		return
	}

	if r.Method != "GET" {
		http.Error(w, "Method is not supported.", http.StatusNotFound)
		return
	}
}

func getKey(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/getKey" {
		http.Error(w, "404 not found.", http.StatusNotFound)
		return
	}

	if r.Method != "GET" {
		http.Error(w, "Method is not supported.", http.StatusNotFound)
		return
	}

	fmt.Fprintf(w, getEnv("API_KEY"))
}

func parseQuery(query url.Values) (map[string]float64, error) {
	values := make(map[string]float64)
	for key := range query {
		s, err := strconv.ParseFloat(query.Get(key), 64)
		if err != nil {
			return nil, err
		}
		values[key] = s
	}
	return values, nil
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

	var candidate Loc
	if (game.Target == Loc{0, 0}) {
		boundaryCheck := true
		for boundaryCheck {
			candidate = Loc{
				randRange(borders[game.Region][0].Lat, borders[game.Region][1].Lat),
				randRange(borders[game.Region][0].Lng, borders[game.Region][1].Lng),
			}

			if game.Region == "World" || isLocInPoly(polys[game.Region], candidate) {
				boundaryCheck = false
			}
		}
	} else {
		candidate = game.Target
	}

	js, err := json.Marshal(candidate)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}

func getRadius(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/getRadius" {
		http.Error(w, "404 not found.", http.StatusNotFound)
		return
	}

	if r.Method != "GET" {
		http.Error(w, "Method is not supported.", http.StatusNotFound)
		return
	}

	js, err := json.Marshal(radii[game.Region])
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

func receiveTarget(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/receiveTarget" {
		http.Error(w, "404 not found.", http.StatusNotFound)
		return
	}

	if r.Method != "GET" {
		http.Error(w, "Method is not supported.", http.StatusNotFound)
		return
	}

	values, err := parseQuery(r.URL.Query())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	game.Target = Loc{values["targetLat"], values["targetLng"]}
}

func getRoundData(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/getRoundData" {
		http.Error(w, "404 not found.", http.StatusNotFound)
		return
	}

	if r.Method != "GET" {
		http.Error(w, "Method is not supported.", http.StatusNotFound)
		return
	}

	if (game.Target == Loc{0, 0}) {
		http.Error(w, "Next round's location not determined.", http.StatusInternalServerError)
		return
	}

	values, err := parseQuery(r.URL.Query())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	if game.Region == "World" {
		distanceFactor = math.Pow((2000 / 14), 2)
	} else {
		updateDistanceFactor(polys[game.Region])
	}

	marker := Loc{values["markerLat"], values["markerLng"]}
	scoreInt := calculateScore(marker, game.Target)
	game.Score = scoreInt
	game.TotalScore += scoreInt
	game.Round++
	if game.Round == 5 && game.HighScore < game.TotalScore {
		game.HighScore = game.TotalScore
	}

	js, err := json.Marshal(game)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(js)

	game.Target = Loc{0, 0}
	if game.Round == 5 {
		game = GameData{Loc{0, 0}, 0, 0, 0, game.HighScore, defaultRegion}
	}
}

func updateDistanceFactor(poly Polygon) {
	latMin, latMax := poly[0].Lat, poly[0].Lng
	lngMin, lngMax := poly[0].Lng, poly[0].Lng
	for i := 1; i < len(poly); i++ {
		if poly[i].Lat < latMin {
			latMin = poly[i].Lat
		}
		if poly[i].Lat > latMax {
			latMax = poly[i].Lat
		}
		if poly[i].Lng < lngMin {
			lngMin = poly[i].Lng
		}
		if poly[i].Lng > lngMax {
			lngMax = poly[i].Lng
		}
	}
	distanceFactor = (lngMax - lngMin) * (latMax - latMin)
	fmt.Println(distanceFactor)
}

func getHighScore(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/getHighScore" {
		http.Error(w, "404 not found.", http.StatusNotFound)
		return
	}

	if r.Method != "GET" {
		http.Error(w, "Method is not supported.", http.StatusNotFound)
		return
	}

	js, err := json.Marshal(game.HighScore)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}

func calculateScore(loc1, loc2 Loc) int {
	latRad1 := loc1.Lat * math.Pi / 180
	latRad2 := loc2.Lat * math.Pi / 180
	deltaLat := (loc2.Lat - loc1.Lat) * math.Pi / 180
	deltaLng := (loc2.Lng - loc1.Lng) * math.Pi / 180
	a := math.Pow(math.Sin(deltaLat/2), 2) + math.Cos(latRad1)*math.Cos(latRad2)*math.Pow(math.Sin(deltaLng/2), 2)
	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))

	distance := 6371e3 * c / 1000
	score := int(math.Round(5000 * math.Pow(math.E, (-distance/(float64(14)*math.Sqrt(distanceFactor))))))

	return score
}

func restart(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/restart" {
		http.Error(w, "404 not found.", http.StatusNotFound)
		return
	}

	if r.Method != "GET" {
		http.Error(w, "Method is not supported.", http.StatusNotFound)
		return
	}

	region := r.URL.Query().Get("region")
	if _, exist := polys[region]; exist {
		game = GameData{Loc{0, 0}, 0, 0, 0, 0, region}
		fmt.Fprintf(w, "OK")
		return
	}

	fmt.Fprintf(w, "Region not found")
}
