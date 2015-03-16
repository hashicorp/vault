Demo.DemoCrudController = Ember.ObjectController.extend({
  needs: ['demo'],
  currentText: Ember.computed.alias('controllers.demo.currentText'),
  currentLog: Ember.computed.alias('controllers.demo.currentLog'),
  logPrefix: Ember.computed.alias('controllers.demo.logPrefix'),
  currentMarker: Ember.computed.alias('controllers.demo.currentMarker'),
  notCleared: Ember.computed.alias('controllers.demo.notCleared'),

  actions: {
    submitText: function() {
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
    }
  }
});
