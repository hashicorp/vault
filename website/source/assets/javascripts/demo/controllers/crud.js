Demo.DemoCrudController = Ember.ObjectController.extend({
  needs: ['demo'],
  isLoading: Ember.computed.alias('controllers.demo.isLoading'),
  currentText: Ember.computed.alias('controllers.demo.currentText'),
  currentLog: Ember.computed.alias('controllers.demo.currentLog'),
  logPrefix: Ember.computed.alias('controllers.demo.logPrefix'),
  currentMarker: Ember.computed.alias('controllers.demo.currentMarker'),
  notCleared: Ember.computed.alias('controllers.demo.notCleared'),

  sendCommand: function() {
    // Request
    Ember.run.later(this, function() {
      var command = this.getWithDefault('currentText', '');
      var currentLogs = this.get('currentLog').toArray();

      // Add the last log item
      currentLogs.push(command);

      // Clean the state
      this.set('currentText', '');

      // Push the new logs
      this.set('currentLog', currentLogs);

      switch(command) {
        case "clear":
          this.set('currentLog', []);
          this.set('notCleared', false);
          break;
        default:
          console.log("Submitting: ", command);
      }

      this.set('isLoading', false);
    }, 1000);
  },

  actions: {
    submitText: function() {
      this.set('isLoading', true);

      // Send the actual request (fake for now)
      this.sendCommand();
    }
  }
});
