import Component from '@ember/component';
import { computed } from '@ember/object';

export default Component.extend({
  title: 'Vault Enterprise',
  featureName: computed('title', function() {
    let title = this.get('title');
    return title === 'Vault Enterprise' ? 'This' : title;
  }),
  minimumEdition: 'Vault Enterprise',
});
