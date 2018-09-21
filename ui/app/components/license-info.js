import Ember from 'ember';
import { allFeatures } from 'vault/helpers/all-features';

export default Ember.Component.extend({
  expirationTime: '',
  startTime: '',
  licenseId: '',
  features: null,
  licenseText: '',
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

  actions: {
    saveLicense() {
      debugger;
      let model = this.get('model');
      model.store.createRecord('license', { text: this.get('licenseText') });
      model.save().then(() => {
        this.setProperties(this.get('model'));
        this.set('licenseText', '');
      });
    },
  },
});
