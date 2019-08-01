import Component from '@ember/component';
import { computed } from '@ember/object';
import layout from '../templates/components/upgrade-page';

export default Component.extend({
  layout,
  title: 'Vault Enterprise',
  featureName: computed('title', function() {
    let title = this.get('title');
    return title === 'Vault Enterprise' ? 'this feature' : title;
  }),
  minimumEdition: 'Vault Enterprise',
});
