let script = document.createElement('script');
script.src = 'https://maps.googleapis.com/maps/api/js?key=AIzaSyBIwzALxUPNbatRBj3Xi1Uhp0fFzwWNBkE&callback=initMap';
script.async = true;

let map;
let panorama;

var markerLat, markerLong;
var randomLoc;

window.initMap = function() {
	const sv = new google.maps.StreetViewService();
	/** For US-only:
	const randomLoc = { lat: rangeRandom(30, 50), lng: rangeRandom(-125, -65) };
	*/
	randomLoc = { lat: rangeRandom(-80, 80), lng: rangeRandom(-180, 180) };
	sv.getPanorama(
		{
			location: randomLoc,
			radius: 100000,
			source: google.maps.StreetViewSource.OUTDOOR
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

		var marker = new google.maps.Marker({
			map: guessingMap,
			position: { lat: 0, lng: 0 },
			draggable: true,
			title: 'Guessing Marker'
		});

		google.maps.event.addListener(marker, 'dragend', function() {
			updateMarker();
		});

		guessingMap.addListener('click', (mapsMouseEvent) => {
			marker.setPosition(mapsMouseEvent.latLng);
			updateMarker();
		});

		function updateMarker() {
			markerLat = marker.getPosition().lat();
			markerLong = marker.getPosition().lng();
			document.getElementById('marker').innerHTML = 'Marker At: ' + markerLat + ',' + markerLong;
			console.log(markerLat + ',' + markerLong);
		}

		document.getElementById('guess').click(function() {
			console.log(calcScore());
		});

		panorama = new google.maps.StreetViewPanorama(document.getElementById('streetView'));
		panorama.setOptions({
			addressControl: false,
			showRoadLabels: false
		});

		panorama.setPano(data.location.pano);
		panorama.setPov({
			heading: 270,
			pitch: 0
		});
		panorama.setVisible(true);
	} else {
		console.log('Street View not found.');
		window.initMap();
	}
}

function rangeRandom(min, max) {
	return Math.random() * (max - min) + min;
}

function calcScore() {
	const latRad1 = markerLat * Math.PI / 180;
	const latRad2 = randomLoc.lat * Math.PI / 180;
	const deltaLat = (randomLoc.lat - markerLat) * Math.PI / 180;
	const deltaLng = (randomLoc.lng - markerLong) * Math.PI / 180;
	const a =
		Math.pow(Math.sin(deltaLat / 2), 2) + Math.cos(latRad1) * Math.cos(latRad2) * Math.pow(Math.sin(deltaLng), 2);
	const c = 2 * Math.atan2(Math.sqrt(a), Marth.sqrt(1 - a));

	const distance = 6371e3 * c; //in meters

	var score = 5000 * Math.pow(e, -distance / 2000);
	return score;
}

document.head.appendChild(script);
