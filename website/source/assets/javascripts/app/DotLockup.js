(function(){

DotLockup = Base.extend({

	$keyWrap: null,
	$keys: null,

	constructor: function(){
		var _this = this;

		_this.$keyWrap = $('.keys');
		_this.$keys = $('.keys span');

		//(3000)

		_this.addEventListeners();
		_this.animateFull()
			.then(_this.animateOff.bind(this))
			.then(_this.animateFull.bind(this))
			.then(_this.animatePress.bind(this))
			.then(_this.resetKeys.bind(this));
	},

	addEventListeners: function(){
		var _this = this;
	},

	animateFull: function(uberDelay){
		var _this = this,
			uberDelay = uberDelay || 0,
			deferred = $.Deferred();

		setTimeout( function(){
			_this.updateEachKeyClass('full', 'off', 1000, 150, deferred.resolve);
		}, uberDelay)

		return deferred;
	},

	animateOff: function(){
		var deferred = $.Deferred();

		this.updateEachKeyClass('full off', '', 1000, 150, deferred.resolve, true);
		
		return deferred;
	},

	animatePress: function(){
		var _this = this,
			deferred = $.Deferred(),
			len = _this.$keys.length,
			presses = _this.randomNumbersIn(len),
			delay = 250,
			interval = 600;

		for(var i=0; i < len; i++){
			(function(index){
				setTimeout(function(){
					_this.$keys.eq(presses[index]).addClass('press');
					if(index == len -1 ){
						deferred.resolve();
					}				
				}, delay)

				delay += interval;
			}(i))
		}

		return deferred;
	},

	resetKeys: function(){
		var _this = this,
			len = _this.$keys.length,
			delay = 2500,
			interval = 250;

		setTimeout(function(){
			_this.$keys.removeClass('full press');	
		}, delay)
		/*for(var i=0; i < len; i++){
			(function(index){
				setTimeout(function(){
					_this.$keys.eq(index).removeClass('full press');	
				}, delay)

				delay += interval;
			}(i))
		}*/		
	},

	updateEachKeyClass: function(clsAdd, clsRemove, delay, interval, resolve, reverse){
		var delay = delay;
		this.$keys.each(function(index){
			var span = this;
			var finishIndex = (reverse) ? 0 : 9; // final timeout at 0 or 9 depending on if class removal is reversed on the span list
			setTimeout( function(){ 
				$(span).removeClass(clsRemove).addClass(clsAdd);
				if(index == finishIndex ){
					resolve();
				}
			}, delay);

			if(reverse){
				delay -= interval;
			}else{
				delay += interval;
			}
		})		

	},

	randomNumbersIn: function(len){
		var arr = [];
		while(arr.length < len){
		  	var randomnumber=Math.floor(Math.random()*len)
		  	var found=false;
		  	for(var i=0;i<arr.length;i++){
				if(arr[i]==randomnumber){found=true;break}
		  	}
		  	if(!found)arr[arr.length]=randomnumber;
		}
		return arr;
	}

});

window.DotLockup = DotLockup;

})();
