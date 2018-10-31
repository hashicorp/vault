import Component from '@ember/component';

export default Component.extend({
  'data-test-navheader': true,
  tagName: 'header',
  navDrawerOpen: false,
  actions: {
    toggleNavDrawer(isOpen) {
      if (isOpen !== undefined) {
        return this.set('navDrawerOpen', isOpen);
      }
      this.toggleProperty('navDrawerOpen');
    },
  },
});
