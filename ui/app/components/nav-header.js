import Component from '@ember/component';
import { inject as service } from '@ember/service';
import { computed } from '@ember/object';

export default Component.extend({
  router: service(),
  'data-test-navheader': true,
  classNameBindings: 'consoleFullscreen:panel-fullscreen',
  tagName: 'header',
  navDrawerOpen: false,
  consoleFullscreen: false,
  hideLinks: computed('router.currentRouteName', function() {
    let currentRoute = this.router.currentRouteName;
    if ('vault.cluster.oidc-provider' === currentRoute) {
      return true;
    }
    return false;
  }),
  actions: {
    toggleNavDrawer(isOpen) {
      if (isOpen !== undefined) {
        return this.set('navDrawerOpen', isOpen);
      }
      this.toggleProperty('navDrawerOpen');
    },
  },
});
