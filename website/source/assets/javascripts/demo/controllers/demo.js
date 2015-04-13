Demo.DemoController = Ember.ObjectController.extend({
  currentText: "vault help",
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
  },

  logCommand: function(command) {
    var commandLog = this.get('commandLog');

    commandLog.push(command);

    this.set('commandLog', commandLog);
  },

  actions: {
    close: function() {
      this.transitionTo('index');
    },
  }
});
