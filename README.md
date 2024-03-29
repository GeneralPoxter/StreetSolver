<h1 align="center">
    <p>StreetSolver</p>
    <img src="public/img/favicon.png" width="200"/>
</h1>

<p align="center">Geographic discovery game based off of GeoGuessr, but completely free.</p>

## Usage
### Setup
The game was built for Go v^1.16. To locally serve the game, run:
```sh
go run server.go
```

Environment variables `PORT` and `API_KEY` can be modified in the `.env` file. By default, the game is served on localhost port 8080 and uses Google's public Maps JavaScript API Key. The rationale behind this is to avoid API costs out-of-the-box, and because this open source key only works on localhosts and URLs permitted by Google. This key is changed occassionally, so expired keys can be locally updated with the current one found [here](https://github.com/googlemaps/js-samples/blob/08d6e630e8baa89d9fef856d9596258b9550293f/dist/samples/add-map/index.html#L58).

### Game
For each round, use the Google StreetView interface to navigate the surroundings of the starting position.
Guess where the starting position is by clicking or dragging the marker on the world map.
To return to the starting position, use the "Return" button, and to submit the guess, use "Guess".
Each round is scored out of 5000 based on proximity, for a 5-round game totaling 25000 points.
Select an option in the dropdown menu to restart the game in that play region (default region: United States).  
Game round progress is tracked even after the user quits the frontend, as long as the server itself remains running.

## Developers
Made for Guided Research B: Future of Programming Languages:
* Natanel Ha
* Jason Liu
* Henry Ren
