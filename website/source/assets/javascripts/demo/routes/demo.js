Demo.DemoRoute = Ember.Route.extend({
  activate: function() {
    // connect to the websocket once we enter the application route
    // var socket = window.io.connect('http://localhost:8080');
    var socket = new WebSocket("ws://vault-demo-server.herokuapp.com/socket");

    this.controllerFor('application').set('socket', socket);

    socket.onmessage = function(message) {
      var data = JSON.parse(message.data),
          controller = this.controllerFor('demo');

      // Add the item
      if (data.stdout) {
        controller.appendLog(data.stdout);
      }

      if (data.stderr) {
        controller.appendLog(data.stderr);
      }

      controller.set('isLoading', false);
    }.bind(this);
  }
});
