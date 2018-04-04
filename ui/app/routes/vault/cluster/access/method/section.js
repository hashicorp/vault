import Ember from 'ember';
import DS from 'ember-data';

export default Ember.Route.extend({
  model(params) {
    const { section_name: section } = params;
    if (section !== 'configuration') {
      const error = new DS.AdapterError();
      Ember.set(error, 'httpStatus', 404);
      throw error;
    }
    return this.modelFor('vault.cluster.access.method');
  },

  setupController(controller) {
    const { section_name: section } = this.paramsFor(this.routeName);
    this._super(...arguments);
    controller.set('section', section);
  },
});
