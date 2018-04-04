import Ember from 'ember';

export default Ember.Component.extend({
  tagName: 'span',
  classNames: 'badge edition-badge',
  abbreviation: Ember.computed('edition', function() {
    const edition = this.get('edition');
    if (edition == 'Enterprise') {
      return 'Ent';
    } else {
      return edition;
    }
  }),
  attributeBindings: ['edition:aria-label'],
});
