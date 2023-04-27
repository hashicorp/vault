/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
 */

import Component from '@ember/component';
import { inject as service } from '@ember/service';
import { computed } from '@ember/object';

export default Component.extend({
  router: service(),
  currentCluster: service(),
  'data-test-navheader': true,
  attributeBindings: ['data-test-navheader'],
  classNameBindings: 'consoleFullscreen:panel-fullscreen',
  tagName: 'header',
  navDrawerOpen: false,
  consoleFullscreen: false,
  hideLinks: computed('router.currentRouteName', function () {
    const currentRoute = this.router.currentRouteName;
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
