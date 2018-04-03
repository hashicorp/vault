import Ember from 'ember';

export default Ember.Component.extend({
  classNames: ['column', 'is-sidebar'],
  classNameBindings: ['isActive:is-active'],
  isActive: false,
  actions: {
    openMenu() {
      this.set('isActive', true);
    },
    closeMenu() {
      this.set('isActive', false);
    },
  },
});
