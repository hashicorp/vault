Demo.DemoStepController = Ember.ObjectController.extend({
  needs: ['demo'],
  socket: Ember.computed.alias('controllers.demo.socket'),
  logs: Ember.computed.alias('controllers.demo.logs'),
  isLoading: Ember.computed.alias('controllers.demo.isLoading'),

  currentText: "",
  commandLog: [],
  cursor: 0,
  notCleared: true,
  fullscreen: false,

  renderedLogs: function() {
    return this.get('logs');
  }.property('logs.length'),

  setFromHistory: function() {
    var index = this.get('commandLog.length') + this.get('cursor');
    var previousMessage = this.get('commandLog')[index];

    this.set('currentText', previousMessage);
  }.observes('cursor'),

  logCommand: function(command) {
    var commandLog = this.get('commandLog');

    commandLog.push(command);

    this.set('commandLog', commandLog);
  },

  actions: {
    submitText: function() {
      // Send the actual request (fake for now)
      this.sendCommand();
    },

    close: function() {
      this.transitionTo('index');
    },

    next: function() {
      var nextStepNumber = parseInt(this.get('model.id'), 10) + 1;
      this.transitionTo('demo.step', nextStepNumber);
    },

    previous: function() {
      var prevStepNumber = parseInt(this.get('model.id'), 10) - 1;
      this.transitionTo('demo.step', prevStepNumber);
    },
  },

  sendCommand: function() {
      var command = this.getWithDefault('currentText', '');
      var log = this.get('log');

      this.set('currentText', '');
      this.logCommand(command);
      this.get('controllers.demo').appendLog(command, true);

      switch(command) {
        case "":
          break;
        case "next":
        case "forward":
          this.set('notCleared', true);
          this.send('next');
          break;
        case "previous":
        case "back":
        case "prev":
          this.set('notCleared', true);
          this.send('previous');
          break;
        case "quit":
        case "exit":
          this.send('close');
          break;
        case "clear":
          this.set('logs', "");
          this.set('notCleared', false);
          break;
        case "fu":
        case "fullscreen":
          this.set('fullscreen', true);
          break;
        case "help":
          this.get('controllers.demo').appendLog('You can use `vault path-help <command>` ' +
            'to learn more about specific Vault commands, or `next` ' +
            'and `previous` to navigate. Or, `fu` to go fullscreen.', false);
          break;
        default:
          this.set('isLoading', true);
          var data = JSON.stringify({type: "cli", data: {command: command}});
          this.get('socket').send(data);
      }
  },
});
