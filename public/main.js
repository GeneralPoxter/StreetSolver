let script = document.createElement('script');
script.src = 'https://maps.googleapis.com/maps/api/js?key=AIzaSyBIwzALxUPNbatRBj3Xi1Uhp0fFzwWNBkE&callback=initMap';
script.async = true;

let map;
let panorama;
let marker;
let target

window.initMap = async function() {
	const sv = new google.maps.StreetViewService();
	sv.getPanorama(
		{
			location: await fetch('/getLoc').then((res) => res.text()).then((res) => JSON.parse(res)),
			source: google.maps.StreetViewPreference.OUTDOOR,
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

		google.maps.event.addListener(marker, 'dragend', function() {
			updateMarker();
		});

		guessingMap.addListener('click', (mapsMouseEvent) => {
			marker.setPosition(mapsMouseEvent.latLng);
			updateMarker();
		});

		panorama = new google.maps.StreetViewPanorama(document.getElementById('streetView'));
		panorama.setOptions({
			zoom: 0,
			addressControl: false,
			showRoadLabels: false
		});

		panorama.setPano(data.location.pano);
		panorama.setPov({
			heading: 270,
			pitch: 0
		});
		panorama.setVisible(true);

		target = {
			lat: data.location.latLng.lat(),
			lng: data.location.latLng.lng()
		}
	} else {
		console.log('Street View not found.');
		window.initMap();
	}
}

function updateMarker() {
	document.getElementById('marker').innerHTML =
		'Marker At: ' + marker.getPosition().lng() + ', ' + marker.getPosition().lat();
}

async function getScore() {
	let params = new URLSearchParams({
		targetLat: target.lat,
		targetLng: target.lng,
		markerLat: marker.getPosition().lat(),
		markerLng: marker.getPosition().lng()
	}).toString();
	console.log(params);
	const score = await fetch('/getScore?' + params).then((res) => res.text());
	document.getElementById('marker').innerHTML = 'Score: ' + score + ' / 5000';
}

document.head.appendChild(script);
