let script = document.createElement('script');
script.src = 'https://maps.googleapis.com/maps/api/js?key=AIzaSyBIwzALxUPNbatRBj3Xi1Uhp0fFzwWNBkE&callback=initMap';
script.async = true;

let map;
let panorama;
let marker;
let markerCorrect;
let target;
let line;
let ready = false;

let info = document.getElementById('info');
let guessButton = document.getElementById('guess');
let returnButton = document.getElementById('return');

window.initMap = async function() {
	const sv = new google.maps.StreetViewService();
	sv.getPanorama(
		{
			location: await fetch('/getLoc').then((res) => res.text()).then((res) => JSON.parse(res)),
			source: google.maps.StreetViewPreference.OUTDOOR
			// preference: google.maps.StreetViewPreference.BEST
		},
		initGame
	);
};

function initGame(data, status) {
	if (status === 'OK') {
		guessingMap = new google.maps.Map(document.getElementById('guessingMap'), {
			center: {
				lat: 0,
				lng: 0
			},
			zoom: 1,
			disableDefaultUI: true,
			fullscreenControl: true,
			zoomControl: true,
			zoomControlOptions: {
				style: google.maps.ZoomControlStyle.LARGE
			}
		});

		marker = new google.maps.Marker({
			map: guessingMap,
			position: { lat: 0, lng: 0 },
			draggable: true,
			title: 'Guessing Marker'
		});
		updateMarker();

		markerCorrect = new google.maps.Marker({
			map: guessingMap,
			position: { lat: 0, lng: 0 },
			icon: {
				path: google.maps.SymbolPath.CIRCLE,
				scale: 5
			},
			draggable: false,
			title: 'Actual Location'
		});
		markerCorrect.setVisible(false);

		line = new google.maps.Polyline({
			icons: [
				{
					icon: google.maps.SymbolPath.FORWARD_CLOSED_ARROW,
					offset: '100%'
				}
			],
			map: guessingMap
		});
		line.setVisible(false);

		google.maps.event.addListener(marker, 'dragend', function() {
			updateMarker();
		});

		guessingMap.addListener('click', (mapsMouseEvent) => {
			if (guessButton.innerHTML == 'Guess') {
				marker.setPosition(mapsMouseEvent.latLng);
				updateMarker();
			}
		});

		panorama = new google.maps.StreetViewPanorama(document.getElementById('streetView'));
		panorama.setOptions({
			zoom: 0,
			addressControl: false,
			showRoadLabels: false
		});

		goTo(data.location.pano);

		returnButton.onclick = function() {
			goTo(data.location.pano);
		};

		guessButton.onclick = getScore;

		target = {
			lat: data.location.latLng.lat(),
			lng: data.location.latLng.lng()
		};

		ready = true;
	} else {
		console.log('Street View not found.');
		window.initMap();
	}
}

function goTo(loc) {
	panorama.setPano(loc);
	panorama.setPov({
		heading: 270,
		pitch: 0
	});
	panorama.setVisible(true);
}

function updateMarker() {
	info.innerHTML = `Marker: {${marker.getPosition().lng()}, ${marker.getPosition().lat()}}`;
}

async function getScore() {
	if (guessButton.innerHTML == 'Guess') {
		if (!ready) {
			return;
		}
		let params = new URLSearchParams({
			targetLat: target.lat,
			targetLng: target.lng,
			markerLat: marker.getPosition().lat(),
			markerLng: marker.getPosition().lng()
		}).toString();
		const score = await fetch('/getScore?' + params).then((res) => res.text());
		const round = await fetch('/getRound').then((res) => res.text());
		if (round < 5) {
			resetUI('Next', `Score: ${score} / 5000`);
		} else {
			resetUI('Finish', `Score: ${score} / 5000`);
		}
		markerCorrect.setVisible(true);
		marker.setDraggable(false);
		markerCorrect.setPosition({ lat: target.lat, lng: target.lng });
		line.setVisible(true);
		line.setPath([
			{ lat: marker.getPosition().lat(), lng: marker.getPosition().lng() },
			{ lat: target.lat, lng: target.lng }
		]);
	} else if (guessButton.innerHTML == 'Finish') {
		const finalScore = await fetch('/getTotalScore').then((res) => res.text());
		resetUI('Restart', `Final Score: ${finalScore} / 25000`);
	} else {
		resetUI('Guess', 'Loading street view...');
		ready = false;
		initMap();
	}
}

function resetUI(guessText, infoText) {
	guessButton.innerHTML = guessText;
	info.innerHTML = infoText;
	marker.setDraggable(true);
	markerCorrect.setVisible(false);
	line.setVisible(false);
}

document.head.appendChild(script);
