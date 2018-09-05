import Ember from 'ember';

const { Component, computed } = Ember;

export default Component.extend({
  tagName: '',

  text: computed(function() {
    return 'home';
  }),

  computedClasses: computed('classNames', function() {
    return this.get('classNames').join(' ');
  }),
});
