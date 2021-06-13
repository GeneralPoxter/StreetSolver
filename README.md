# StreetSolver

Geographic discovery game based off of GeoGuessr, but completely free. Made for Guided Research B: Future of Programming Languages.

## Usage
### Setup
The game was built for Go v^1.16. To locally serve the game, run:
```sh
go run server.go
```

Environment variables `PORT` and `API_KEY` can be modified in the `.env` file. By default, the game is served on localhost port 8080 and uses Google's public Maps JavaScript API Key. The rationale behind this is to avoid API costs out-of-the-box, and because this open source key only works on localhosts and URLs permitted by Google. This key is changed occassionally, so expired keys can be locally updated with the current one found [here](https://github.com/googlemaps/js-samples/blob/08d6e630e8baa89d9fef856d9596258b9550293f/dist/samples/add-map/index.html#L58).

### Game
For each round, use the Google StreetView interface to navigate the surroundings of the starting position.
Select the play region with the "Region" dropdown menu followed by a game restart with the "Restart" button. 
Guess where the starting position is by clicking or dragging the marker on the world map.
To return to the starting position, use the "Return" button, and to submit the guess, use "Guess".
Each round is scored out of 5000 based on proximity, for a 5-round game totaling 25000 points.  
Game round progress is tracked even after the user quits the frontend, as long as the server itself remains running.

## Developers
* Natanel Ha
* Jason Liu
* Henry Ren
