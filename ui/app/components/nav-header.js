import Component from '@ember/component';
export default Component.extend({
  'data-test-navheader': true,
  classNameBindings: 'consoleFullscreen:panel-fullscreen',
  tagName: 'header',
  navDrawerOpen: false,
  consoleFullscreen: false,
  actions: {
    toggleNavDrawer(isOpen) {
      if (isOpen !== undefined) {
        return this.set('navDrawerOpen', isOpen);
      }
      this.toggleProperty('navDrawerOpen');
    },
  },
});
