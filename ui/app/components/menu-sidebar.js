/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import Component from '@ember/component';

export default Component.extend({
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
