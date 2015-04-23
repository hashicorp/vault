Demo.DemoStepRoute = Ember.Route.extend({
  model: function(params) {
    return this.store.find('step', params.id);
  }
});
