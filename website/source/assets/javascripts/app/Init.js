(function(Sidebar, DotLockup){

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

    //always initialize sidebar
		Init.initializeSidebar();
	},

	initializeSidebar: function(){
		new Sidebar();
	},

	initializeDotLockup: function(){
		new DotLockup();
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
			Init.initializeDotLockup();
			Init.initializeWaypoints();
		}
	}

};

Init.start();

})(window.Sidebar, window.DotLockup);
