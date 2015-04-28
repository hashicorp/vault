Ember.Clock = Ember.Object.extend({
  second: null,
  checks: 0,
  isPolling: false,
  pollImmediately: true,
  pollInterval: null,
  tickInterval: null,
  defaultPollInterval: 10000,
  defaultTickInterval: 1000,

  init: function () {
      this.scheduleTick();
      this.pollImmediately && this.startPolling();
  },

  computedPollInterval: function() {
    return this.pollInterval || this.defaultPollInterval; // Time between polls (in ms)
  }.property().readOnly(),

  computedTickInterval: function() {
    return this.tickInterval ||  this.defaultTickInterval; // Time between ticks (in ms)
  }.property().readOnly(),

  // Schedules the function `f` to be executed every `interval` time.
  schedulePoll: function(f) {
    return Ember.run.later(this, function() {
      f.apply(this);
      this.incrementProperty('checks')
      this.set('timer', this.schedulePoll(f));
    }, this.get('computedPollInterval'));
  },

  scheduleTick: function () {
    return Ember.run.later(this, function () {
          this.tick();
          this.scheduleTick();
      }, this.get('computedTickInterval'));
  },

  // Stops the polling
  stopPolling: function() {
    this.set('isPolling', false);
    Ember.run.cancel(this.get('timer'));
  },

  // Starts the polling, i.e. executes the `onPoll` function every interval.
  startPolling: function() {
    this.set('isPolling', true);
    this.set('timer', this.schedulePoll(this.get('onPoll')));
  },

  // Moves the clock
  tick: function () {
      var now = new Date();
      this.setProperties({
          second: now.getSeconds()
      });
  },

  // Override this while making a new Clock object.
  onPoll: function(){
  },
});
