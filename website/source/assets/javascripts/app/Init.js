(function(Sidebar){

// Quick and dirty IE detection
var isIE = (function(){
	if (window.navigator.userAgent.match('Trident')) {
		return true;
	} else {
		return false;
	}
})();

// isIE = true;

var Init = {

	start: function(){
		var id = document.body.id.toLowerCase();

		if (this.Pages[id]) {
			this.Pages[id]();
		}
	},

	initializeSidebar: function(){
		new Sidebar();
	},

	initializeWaypoints: function(){
		$('#header').waypoint(function(event, direction) {
		    $(this.element).addClass('showit');
		}, {
		    offset: function() {
		    	return '25%';
		    }
		});

		$('#hero').waypoint(function(event, direction) {
		    $(this.element).addClass('showit');
		}, {
		    offset: function() {
		    	return '25%';
		    }
		});			
	},	

	Pages: {
		'page-home': function(){
			Init.initializeSidebar();
			Init.initializeWaypoints();
		}
	}

};

Init.start();

})(window.Sidebar);