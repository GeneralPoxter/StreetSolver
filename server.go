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

	geojson "github.com/paulmach/go.geojson"
)

type Loc struct {
	Lat float64 `json:"lat"`
	Lng float64 `json:"lng"`
}

type Polygon []Loc

type tagger struct {
	Tags []string `json:"tags"`
}

type GameData struct {
	Target     Loc `json:"target"`
	Score      int `json:"score"`
	TotalScore int `json:"totalScore"`
	Round      int `json:"round"`
}

var game GameData

var distanceFactor = 2000

// Maryland
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
	http.HandleFunc("/receiveTarget", receiveTarget)
	http.HandleFunc("/getRoundData", getRoundData)

	g, _ := NewGeoJSON(Loc{0, 0}, []string{"foo", "bar"})
	fmt.Println(string(g))

	game = GameData{Loc{0, 0}, 0, 0, 0}

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
			/* // World map
			candidate = Loc{randRange(-80, 80), randRange(-180, 180)}
			// TODO: Make a Poly that excludes oceans
			boundaryCheck = false */

			candidate = Loc{randRange(38, 40), randRange(-78, -75)}
			if isLocInPoly(poly, candidate) {
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

func NewGeoJSON(loc Loc, tags []string) ([]byte, error) {
	featureCollection := geojson.NewFeatureCollection()
	feature := geojson.NewPointFeature([]float64{loc.Lng, loc.Lat})
	featureCollection.AddFeature(feature)
	return featureCollection.MarshalJSON()
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

	//Adjusted score for smaller area of MD
	updateDistanceFactor(poly)

	marker := Loc{values["markerLat"], values["markerLng"]}
	scoreInt := calculateScore(marker, game.Target)
	game.Score = scoreInt
	game.TotalScore += scoreInt
	game.Round++

	js, err := json.Marshal(game)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(js)

	game.Target = Loc{0, 0}
	if game.Round == 5 {
		game = GameData{Loc{0, 0}, 0, 0, 0}
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
	distanceFactor = int(2 * (lngMax - lngMin) * (latMax - latMin))
}

func calculateScore(loc1, loc2 Loc) int {
	latRad1 := loc1.Lat * math.Pi / 180
	latRad2 := loc2.Lat * math.Pi / 180
	deltaLat := (loc2.Lat - loc1.Lat) * math.Pi / 180
	deltaLng := (loc2.Lng - loc1.Lng) * math.Pi / 180
	a := math.Pow(math.Sin(deltaLat/2), 2) + math.Cos(latRad1)*math.Cos(latRad2)*math.Pow(math.Sin(deltaLng/2), 2)
	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))

	distance := 6371e3 * c / 1000
	score := int(math.Round(5000 * math.Pow(math.E, (-distance/float64(distanceFactor)))))

	return score
}
