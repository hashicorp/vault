import Component from '@ember/component';
import { allFeatures } from 'vault/helpers/all-features';
import { computed } from '@ember/object';

export default Component.extend({
  expirationTime: '',
  startTime: '',
  licenseId: '',
  features: null,
  model: null,
  text: '',
  showForm: false,
  isTemporary: computed('licenseId', function() {
    return this.licenseId === 'temporary';
  }),
  featuresInfo: computed('model', 'features', function() {
    return allFeatures().map(feature => {
      let active = this.features.includes(feature);
      if (active && feature === 'Performance Standby') {
        let count = this.model.performanceStandbyCount;
        return {
          name: feature,
          active: count ? active : false,
          count,
        };
      }
      return { name: feature, active };
    });
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
