import Ember from 'ember';
import allFeatures from 'vault/helpers/all-features';

export default Ember.Component.extend({
  expirationTime: null,
  startTime: null,
  licenseId: null,
  features: null,
  featuresInfo: Ember.computed('features', function() {
    let info = {};
    for (let feature in allFeatures()) {
      info[feature] = this.get('features').includes(feature) ? true : false;
      debugger;
    }
    return info;
  }),

  init() {
    this._super(...arguments);
    this.setProperties(this.get('model'));
    debugger;
  },
});
