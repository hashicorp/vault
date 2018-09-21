import Ember from 'ember';
import { allFeatures } from 'vault/helpers/all-features';

export default Ember.Component.extend({
  expirationTime: null,
  startTime: null,
  licenseId: null,
  features: null,
  isTemporary: Ember.computed('licenseId', function() {
    return this.get('licenseId') === 'temporary';
  }),
  featuresInfo: Ember.computed('features', function() {
    let info = [];
    allFeatures().forEach(feature => {
      let active = this.get('features').includes(feature) ? true : false;
      info.push({ name: feature, active: active });
    });
    return info;
  }),

  init() {
    this._super(...arguments);
    this.setProperties(this.get('model'));
  },

  actions: {
    saveLicense(text) {
      let model = this.get('model');
      model = model.createRecord(text);
      this.set('model', model);
      this.get('model').save().then(() => {
        this.setProperties(this.get('model'));
      });
    },
  },
});
