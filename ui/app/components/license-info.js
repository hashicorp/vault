import Component from '@ember/component';
import { allFeatures } from 'vault/helpers/all-features';
import { computed } from '@ember/object';

export default Component.extend({
  expirationTime: '',
  startTime: '',
  licenseId: '',
  features: null,
  text: '',
  showForm: false,
  isTemporary: computed('licenseId', function() {
    return this.licenseId === 'temporary';
  }),
  featuresInfo: computed('features', function() {
    let info = [];
    allFeatures().forEach(feature => {
      let active = this.features.includes(feature) ? true : false;
      info.push({ name: feature, active: active });
    });
    return info;
  }),
  saveModel() {},
  actions: {
    saveModel(text) {
      this.saveModel(text);
    },
    toggleForm() {
      this.toggleProperty('showForm');
    },
  },
});
