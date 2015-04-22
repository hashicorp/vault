Demo.DemoStepController = Ember.ObjectController.extend({
  needs: ['application'],
  socket: Ember.computed.alias('controllers.application.socket'),

  currentText: "",
  commandLog: [],
  logs: "",
  cursor: 0,
  notCleared: true,
  isLoading: false,

  setFromHistory: function() {
    var index = this.get('commandLog.length') + this.get('cursor');
    var previousMessage = this.get('commandLog')[index];

    this.set('currentText', previousMessage);
  }.observes('cursor'),

  appendLog: function(data, prefix) {
    if (prefix) {
      data = '$ ' + data;
    }

    this.set('logs', this.get('logs')+'\n'+data);

    Ember.run.later(function() {
      var element = $('.demo-overlay');
      // Scroll to the bottom of the element
      element.scrollTop(element[0].scrollHeight);
    }, 5);
  },

  logCommand: function(command) {
    var commandLog = this.get('commandLog');

    commandLog.push(command);

    this.set('commandLog', commandLog);
  },

  actions: {
    submitText: function() {
      // Send the actual request (fake for now)
      this.sendCommand();
    }
  },

  sendCommand: function() {
      var demoController = this.get('controllers.demo');
      var command = this.getWithDefault('currentText', '');
      var log = this.get('log');

      this.set('currentText', '');
      demoController.logCommand(command);
      demoController.appendLog(command, true);

      switch(command) {
        case "":
          break;
        case "clear":
          this.set('logs', "");
          this.set('notCleared', false);
          break;
        default:
          this.set('isLoading', true);
          var data = JSON.stringify({type: "cli", data: {command: command}});
          this.get('socket').send(data);
      }
  },
});
