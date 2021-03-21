let script = document.createElement('script');
script.src = 'https://maps.googleapis.com/maps/api/js?key=AIzaSyBIwzALxUPNbatRBj3Xi1Uhp0fFzwWNBkE&callback=initMap';
script.async = true;

let map;
let panorama;

var markerLat, markerLong;

window.initMap = function() {
	const sv = new google.maps.StreetViewService();
	const randomLoc = { lat: rangeRandom(-80, 80), lng: rangeRandom(-180, 180) };
	/** For US-only:
	const randomLoc = { lat: rangeRandom(30, 50), lng: rangeRandom(-125, -65) };
	*/

	sv.getPanorama({
			location: randomLoc,
			radius: 100000,
			source: google.maps.StreetViewSource.OUTDOOR
		},
		function(data, status) {
			if (status === 'OK') {
				guessingMap = new google.maps.Map(document.getElementById('guessingMap'), {
					center: {
						lat: 0,
						lng: 0
					},
					zoom: 1,
					streetViewControl: false,
					mapTypeId: google.maps.MapTypeId.ROADMAP
				});

				var marker = new google.maps.Marker({
					map: guessingMap,
					position: {lat: 0, lng: 0},
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

				function updateMarker(){
					markerLat = marker.getPosition().lat();
					markerLong = marker.getPosition().lng();
					document.getElementById("marker").innerHTML = "Marker At: "+markerLat+","+markerLong;
					console.log(markerLat+","+markerLong);
				}

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
	);
};

function rangeRandom(min, max) {
	return Math.random() * (max - min) + min;
}

/*** 
function() {

	streetview = new google.maps.StreetViewPanorama(document.getElementById('streetView'), {
		position: { 
			lat: 48.8584,
			lng: 2.296
		}
	});

};
*/

document.head.appendChild(script);
