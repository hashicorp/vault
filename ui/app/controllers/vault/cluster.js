/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

/* eslint-disable ember/no-observers */
import { service } from '@ember/service';
import { alias } from '@ember/object/computed';
import Controller from '@ember/controller';
import { observer } from '@ember/object';
export default Controller.extend({
  auth: service(),
  store: service(),
  media: service(),
  router: service(),
  permissions: service(),
  namespaceService: service('namespace'),
  flashMessages: service(),
  customMessages: service(),

  vaultVersion: service('version'),
  console: service(),

  queryParams: [
    {
      namespaceQueryParam: {
        scope: 'controller',
        as: 'namespace',
      },
    },
  ],

  namespaceQueryParam: '',

  onQPChange: observer('namespaceQueryParam', function () {
    this.namespaceService.setNamespace(this.namespaceQueryParam);
  }),

  consoleOpen: alias('console.isOpen'),
  activeCluster: alias('auth.activeCluster'),

  permissionBanner: alias('permissions.permissionsBanner'),

  actions: {
    toggleConsole() {
      this.toggleProperty('consoleOpen');
    },
  },
});
