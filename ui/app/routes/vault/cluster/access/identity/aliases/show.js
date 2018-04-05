import Ember from 'ember';
import DS from 'ember-data';
import { TABS } from 'vault/helpers/tabs-for-identity-show';

export default Ember.Route.extend({
  model(params) {
    let { section } = params;
    let itemType = this.modelFor('vault.cluster.access.identity') + '-alias';
    let tabs = TABS[itemType];
    let modelType = `identity/${itemType}`;
    if (!tabs.includes(section)) {
      const error = new DS.AdapterError();
      Ember.set(error, 'httpStatus', 404);
      throw error;
    }
    // TODO peekRecord here to see if we have the record already
    return Ember.RSVP.hash({
      model: this.store.findRecord(modelType, params.item_alias_id),
      section,
    });
  },

  setupController(controller, resolvedModel) {
    let { model, section } = resolvedModel;
    controller.setProperties({
      model,
      section,
    });
  },
});
