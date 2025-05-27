/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */
import { service } from '@ember/service';
import { alias } from '@ember/object/computed';
import Controller, { inject as controller } from '@ember/controller';
import { task } from 'ember-concurrency';
import { sanitizePath } from 'core/utils/sanitize-path';

export default Controller.extend({
  vaultController: controller('vault'),
  clusterController: controller('vault.cluster'),
  flashMessages: service(),
  namespaceService: service('namespace'),
  flagsService: service('flags'),
  version: service(),
  auth: service(),
  router: service(),
  customMessages: service(),
  queryParams: [{ authMount: 'with', oidcProvider: 'o' }],
  namespaceQueryParam: alias('clusterController.namespaceQueryParam'),
  redirectTo: alias('vaultController.redirectTo'),
  hvdManagedNamespaceRoot: alias('flagsService.hvdManagedNamespaceRoot'),
  shouldRefocusNamespaceInput: false,

  // Query params
  authMount: '',
  oidcProvider: '',
  unwrapTokenError: '',

  fullNamespaceFromInput(value) {
    const strippedNs = sanitizePath(value);
    if (this.hvdManagedNamespaceRoot) {
      return `${this.hvdManagedNamespaceRoot}/${strippedNs}`;
    }
    return strippedNs;
  },

  updateNamespace: task(function* (value) {
    const ns = this.fullNamespaceFromInput(value);
    this.namespaceService.setNamespace(ns, true);
    yield this.customMessages.fetchMessages();
    this.set('namespaceQueryParam', ns);
    // if user is inputting a namespace, maintain input focus as the param updates
    this.set('shouldRefocusNamespaceInput', true);
  }).restartable(),

  actions: {
    // TODO CMB move to auth service?
    authSuccess({ isRoot, namespace }) {
      let transition;
      this.version.fetchVersion();
      if (this.redirectTo) {
        // here we don't need the namespace because it will be encoded in redirectTo
        transition = this.router.transitionTo(this.redirectTo);
        // reset the value on the controller because it's bound here
        this.set('redirectTo', '');
      } else {
        transition = this.router.transitionTo('vault.cluster', { queryParams: { namespace } });
      }
      transition.followRedirects().then(() => {
        if (this.version.isEnterprise) {
          this.customMessages.fetchMessages();
        }

        if (isRoot) {
          this.auth.set('isRootToken', true);
          this.flashMessages.warning(
            'You have logged in with a root token. As a security precaution, this root token will not be stored by your browser and you will need to re-authenticate after the window is closed or refreshed.'
          );
        }
      });
    },
    backToLogin() {
      // reset error
      this.set('unwrapTokenError', '');
      // reset query params and go back to auth route
      this.router.replaceWith('vault.cluster.auth', { queryParams: { wrapped_token: null } });
    },
  },
});
