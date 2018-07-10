import Ember from 'ember';
const { computed } = Ember;

export default Ember.Component.extend({
  content: null,
  list: computed('content', function() {
    return this.get('content').keys;
  }),
});
