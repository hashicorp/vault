Demo.DemoStepRoute = Ember.Route.extend({
  model: function(params) {
    return this.store.find('step', params.id);
  },

  afterModel: function(model) {
    var clock = Ember.Clock.create({
      defaultPollInterval: 5000,
      pollImmediately: false,
      onPoll: function() {
        var socket = this.controllerFor('demo').get('socket');
        socket.send(JSON.stringify({type: "ping"}));
      }.bind(this)
    });

    this.set('clock', clock);
  },

  activate: function() {
    this.get('clock').startPolling();
  },

  deactivate: function() {
    var clock = this.get('clock');
    if(clock.get('isPolling')) {
      clock.stopPolling();
    }
  },
});
