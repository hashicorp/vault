Demo.ApplicationController = Ember.ObjectController.extend({
  needs: ['demo'],

  init: function() {
    this._super.apply(this, arguments);

    // connect to the websocket once we enter the application route
    // var socket = window.io.connect('http://localhost:8080');
    var socket = new WebSocket("ws://vault-demo-server.herokuapp.com/socket");

    this.set('socket', socket);

    socket.onmessage = function(message) {
      var data = JSON.parse(message.data),
          controller = this.get('controllers.demo');

      // Add the item
      if (data.stdout !== "") {
        console.log("stdout:", data.stout);
        controller.appendLog(data.stdout, false);
      }

      if (data.stderr !== "") {
        console.log("stderr:", data.stderr);
        controller.appendLog(data.stderr, false);
      }

      controller.set('isLoading', false);
    }.bind(this);
  }
});
