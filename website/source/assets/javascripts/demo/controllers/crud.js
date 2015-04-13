Demo.DemoCrudController = Ember.ObjectController.extend({
  needs: ['demo', 'application'],
  isLoading: Ember.computed.alias('controllers.demo.isLoading'),
  currentText: Ember.computed.alias('controllers.demo.currentText'),
  logs: Ember.computed.alias('controllers.demo.logs'),
  logPrefix: Ember.computed.alias('controllers.demo.logPrefix'),
  currentMarker: Ember.computed.alias('controllers.demo.currentMarker'),
  notCleared: Ember.computed.alias('controllers.demo.notCleared'),
  socket: Ember.computed.alias('controllers.application.socket'),

  sendCommand: function() {
      this.set('isLoading', true);

      var demoController = this.get('controllers.demo');
      var command = this.getWithDefault('currentText', '');
      var log = this.get('log');

      this.set('currentText', '');
      demoController.logCommand(command);
      demoController.appendLog(command, true);

      switch(command) {
        case "clear":
          this.set('logs', "");
          this.set('notCleared', false);
          break;
        default:
          var data = JSON.stringify({type: "cli", data: {command: command}});
          console.log("Sending: ", data);
          this.get('socket').send(data);
      }
  },

  actions: {
    submitText: function() {
      this.set('controllers.demo.isLoading', true);

      // Send the actual request (fake for now)
      this.sendCommand();
    }
  }
});
