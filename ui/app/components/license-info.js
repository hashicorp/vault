import Ember from 'ember';
import { allFeatures } from 'vault/helpers/all-features';

export default Ember.Component.extend({
  expirationTime: '',
  startTime: '',
  licenseId: '',
  features: null,
  text: '',
  showForm: false,
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
  saveModel() {},
  actions: {
    createModel(text) {
      this.get('saveModel')(text);
    },
    toggleForm() {
      this.toggleProperty('showForm');
    },
  },
});
