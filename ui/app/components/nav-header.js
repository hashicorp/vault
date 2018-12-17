import Component from '@ember/component';
import { inject as service } from '@ember/service';

export default Component.extend({
  permissions: service(),
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
