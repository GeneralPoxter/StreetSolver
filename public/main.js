let script = document.createElement('script');
script.src = 'https://maps.googleapis.com/maps/api/js?key=AIzaSyBIwzALxUPNbatRBj3Xi1Uhp0fFzwWNBkE&callback=initMap';
script.async = true;

let map;
let panorama;
let markerLat = 0,
	markerLng = 0;

window.initMap = async function() {
	markerLat = 0;
	markerLng = 0;
	const sv = new google.maps.StreetViewService();
	/** For US-only:
	const randomLoc = { lat: rangeRandom(30, 50), lng: rangeRandom(-125, -65) };
	*/
	sv.getPanorama(
		{
			location: await fetch('/getLoc').then((res) => res.text()).then((res) => JSON.parse(res)),
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

		let marker = new google.maps.Marker({
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
			markerLng = marker.getPosition().lng();
			document.getElementById('marker').innerHTML = 'Marker At: ' + markerLat + ',' + markerLng;
			console.log(markerLat + ',' + markerLng);
		}

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
	} else {
		console.log('Street View not found.');
		window.initMap();
	}
}

document.head.appendChild(script);
