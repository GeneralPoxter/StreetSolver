# StreetSolver

Geographic discovery game based off of GeoGuessr, but completely free. Made for Guided Research B: Future of Programming Languages.

## Usage
### Setup
The game was built for Go v^1.16. To locally serve the game, run:
```sh
go run server.go
```
Port can be specified with the environment variable PORT.

By default, the game uses Google's public Maps JavaScript API Key, which only works for localhosts and specific URLs. This key is changed occassionally, so expired keys can be replaced with the one found [here](https://github.com/googlemaps/js-samples/blob/08d6e630e8baa89d9fef856d9596258b9550293f/dist/samples/add-map/index.html#L58). To replace the API key with a personal API key so the game can be hosted on services like Heroku, replace the value of `API_KEY` in `src/config.js` with the desired key.

### Game
Navigate to the localhost port to play on the StreetSolver frontend.  
For each round, use the Google StreetView interface to navigate the surroundings of the starting position. Then, guess where the starting position is by clicking or dragging the marker on the world map. To return to the starting position, use the "Return" button, and to conclude the round, click "Guess". Each round is scored out of 5000 based on proximity, for a 5-round game totaling 25000 points.  
Game round progress is tracked even after the user quits the frontend, as long as the server itself remains running.

## Developers
* Natanel Ha
* Jason Liu
* Henry Ren
