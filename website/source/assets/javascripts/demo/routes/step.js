Demo.DemoStepRoute = Ember.Route.extend({
  model: function(params) {
    return this.store.find('step', params.id);
  }
  // socket: function() {
  //   return this.controllerFor('application').get('socket');
  // }.property(),

  // activate: function() {
  //   var data = JSON.stringify({type: "cli", data: {command: "vault init -key-shares=1 -key-threshold=1"}});
  //   var socket = this.get('socket');

  //   socket.onopen = function() {
  //     console.log("ws open");
  //     socket.send(data);
  //   };
  // },
});
