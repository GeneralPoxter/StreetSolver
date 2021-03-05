var script = document.createElement('script');
script.src = 'https://maps.googleapis.com/maps/api/js?key=AIzaSyBIwzALxUPNbatRBj3Xi1Uhp0fFzwWNBkE&callback=initMap';
script.async = true;

window.initMap = function() {
	console.log('Map loaded');
	map = new google.maps.Map(document.getElementById('map'), {
		center: { lat: 39.018, lng: -77.013 },
		zoom: 8
	});
};

document.head.appendChild(script);
