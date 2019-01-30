import { computed } from '@ember/object';
import Component from '@ember/component';

export default Component.extend({
  tagName: 'span',
  classNames: 'tag is-outlined edition-badge',
  attributeBindings: ['edition:aria-label'],
  icon: computed('edition', function() {
    const edition = this.get('edition');
    const entEditions = ['Enterprise', 'Premium', 'Pro'];

    if (entEditions.includes(edition)) {
      return 'edition-enterprise';
    } else {
      return 'edition-oss';
    }
  }),
});
