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
    return this.get('licenseId') === 'temporary';
  }),
  featuresInfo: computed('features', function() {
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
