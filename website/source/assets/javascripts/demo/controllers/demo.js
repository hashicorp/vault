Demo.DemoController = Ember.ObjectController.extend({
  isLoading: false,
  logs: "",

  init: function() {
    this._super.apply(this, arguments);

    // connect to the websocket once we enter the application route
    // var socket = window.io.connect('http://localhost:8080');
    var socket = new WebSocket("wss://vault-demo-server.herokuapp.com/socket");

    // Set socket on application controller
    this.set('socket', socket);

    socket.onmessage = function(message) {
      var data = JSON.parse(message.data),
          controller = this;

      // ignore pongs
      if (data.pong) {
        return
      }

      // Add the item
      if (data.stdout !== "") {
        controller.appendLog(data.stdout, false);
      }

      if (data.stderr !== "") {
        controller.appendLog(data.stderr, false);
      }

      controller.set('isLoading', false);
    }.bind(this);
  },

  appendLog: function(data, prefix) {
    var newline;

    if (prefix) {
      data = '$ ' + data;
    } else {
      newline = '';
    }

    newline = '\n';

    this.set('logs', this.get('logs')+data+newline);

    Ember.run.later(function() {
      var element = $('.demo-terminal');
      // Scroll to the bottom of the element
      element.scrollTop(element[0].scrollHeight);

      element.find('input.shell')[0].focus();
    }, 5);
  },
});
