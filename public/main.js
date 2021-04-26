let script = document.createElement('script');
script.src = 'https://maps.googleapis.com/maps/api/js?key=AIzaSyBIwzALxUPNbatRBj3Xi1Uhp0fFzwWNBkE&callback=initMap';
script.async = true;

let map;
let panorama;
let marker;
let markerCorrect;
let target;
let line;

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
		
		markerCorrect = new google.maps.Marker({
			map: guessingMap,
			position: { lat: 0, lng: 0 },
			icon: {
				path: google.maps.SymbolPath.CIRCLE,
				scale: 5,
			},
			draggable: false,
			title: 'Actual Location'
		});
		markerCorrect.setVisible(false);

		line = new google.maps.Polyline({
			icons: [
				{
					icon: google.maps.SymbolPath.FORWARD_CLOSED_ARROW,
					offset: "100%",
				},
			],
			map: guessingMap,
		});
		line.setVisible(false);

		google.maps.event.addListener(marker, 'dragend', function() {
			updateMarker();
		});

		guessingMap.addListener('click', (mapsMouseEvent) => {
			if (guessButton.innerHTML=="Guess"){
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
	info.innerHTML = 'Marker At: ' + marker.getPosition().lng() + ', ' + marker.getPosition().lat();
}

async function getScore() {
	if (guessButton.innerHTML == 'Guess') {
		let params = new URLSearchParams({
			targetLat: target.lat,
			targetLng: target.lng,
			markerLat: marker.getPosition().lat(),
			markerLng: marker.getPosition().lng()
		}).toString();
		const score = await fetch('/getScore?' + params).then((res) => res.text());
		info.innerHTML = 'Score: ' + score + ' / 5000';
		guessButton.innerHTML = 'Next';
		markerCorrect.setPosition({lat: target.lat, lng: target.lng});
		markerCorrect.setVisible(true)
		marker.setDraggable(false);
		line.setVisible(true)
		line.setPath(
			[
				{lat: marker.getPosition().lat(), lng: marker.getPosition().lng()},
				{lat: target.lat, lng: target.lng},
			]
		)
	} else {
		guessButton.innerHTML = 'Guess';
		info.innerHTML = 'Marker At: 0,0';
		marker.setDraggable(true);
		markerCorrect.setVisible(false);
		line.setVisible(false);
		initMap();
	}
}

document.head.appendChild(script);
