/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { inject as service } from '@ember/service';
import { alias } from '@ember/object/computed';
import Component from '@ember/component';
import { computed } from '@ember/object';

export default Component.extend({
  namespaceService: service('namespace'),
  currentNamespace: alias('namespaceService.path'),

  tagName: '',
  //public api
  targetNamespace: null,
  showLastSegment: false,
  // set to true if targetNamespace is passed in unmodified
  // otherwise, this assumes it is parsed as in namespace-picker
  unparsed: false,

  normalizedNamespace: computed('targetNamespace', 'unparsed', function () {
    const ns = this.targetNamespace || '';
    return this.unparsed ? ns : ns.replace(/\.+/g, '/').replace(/☃/g, '.');
  }),

  namespaceDisplay: computed('normalizedNamespace', 'showLastSegment', function () {
    const ns = this.normalizedNamespace;
    if (!ns) return 'root';
    const showLastSegment = this.showLastSegment;
    const parts = ns?.split('/');
    return showLastSegment ? parts[parts.length - 1] : ns;
  }),

  isCurrentNamespace: computed('targetNamespace', 'currentNamespace', function () {
    return this.currentNamespace === this.targetNamespace;
  }),

  get namespaceLink() {
    const origin =
      window.location.protocol +
      '//' +
      window.location.hostname +
      (window.location.port ? ':' + window.location.port : '');

    if (!this.normalizedNamespace) return `${origin}/ui/vault/dashboard`;
    // The full URL/origin is required so that the page is reloaded.
    return `${origin}/ui/vault/dashboard?namespace=${encodeURIComponent(this.normalizedNamespace)}`;
  },
});
