import Route from '@ember/routing/route';

export default Route.extend({
  model() {
    return this.store.findAll('mfa-method').then((data) => {
      return data;
    });
  },
  setupController(controller, model) {
    this._super(...arguments);
    controller.setProperties({
      model: model,
    });
  },
  actions: {
    willTransition(transition) {
      window.scrollTo(0, 0);
      if (transition.targetName !== this.routeName) {
        this.store.clearAllDatasets();
      }
      return true;
    },
    reload() {
      this.store.clearAllDatasets();
      this.refresh();
    },
  },
});
