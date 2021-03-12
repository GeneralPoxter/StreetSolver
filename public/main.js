var script = document.createElement('script');
script.src = 'https://maps.googleapis.com/maps/api/js?key=AIzaSyBIwzALxUPNbatRBj3Xi1Uhp0fFzwWNBkE&callback=initMap';
script.async = true;

window.initMap = function() {
	
	streetview = new google.maps.StreetViewPanorama(document.getElementById('streetView'),{
		position: { 
			lat: 48.8584,
			lng: 2.296
		},
	})
	console.log('Streetview Map Loaded')

	guessingMap = new google.maps.Map(document.getElementById('guessingMap'), {
		center: { 
			lat: 0, 
			lng: 0
		},
		zoom: 1,
		streetViewControl: false,
		mapTypeId: google.maps.MapTypeId.ROADMAP
	});
	console.log('Guessing Map loaded');
};

document.head.appendChild(script);
