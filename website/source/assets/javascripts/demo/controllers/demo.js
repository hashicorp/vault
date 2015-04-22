Demo.DemoController = Ember.ObjectController.extend({
  actions: {
    close: function() {
      this.transitionTo('index');
    },
  }
});
