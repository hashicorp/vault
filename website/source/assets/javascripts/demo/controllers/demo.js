Demo.DemoController = Ember.ObjectController.extend({
  currentText: "vault help",
  currentLog: [],
  logPrefix: "$ ",
  cursor: 0,
  notCleared: true,

  setFromHistory: function() {
    var index = this.get('currentLog.length') + this.get('cursor');

    this.set('currentText', this.get('currentLog')[index]);
  }.observes('cursor')
});
