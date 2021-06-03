import { equal } from '@ember/object/computed';
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
  featuresInfo: computed('features', 'model.performanceStandbyCount', function() {
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
});
