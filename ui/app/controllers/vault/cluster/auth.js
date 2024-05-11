/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */
import { service } from '@ember/service';
import { alias } from '@ember/object/computed';
import Controller, { inject as controller } from '@ember/controller';
import { task, timeout } from 'ember-concurrency';
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
  queryParams: [{ authMethod: 'with', oidcProvider: 'o' }],
  namespaceQueryParam: alias('clusterController.namespaceQueryParam'),
  wrappedToken: alias('vaultController.wrappedToken'),
  redirectTo: alias('vaultController.redirectTo'),
  hvdManagedNamespaceRoot: alias('flagsService.hvdManagedNamespaceRoot'),
  authMethod: '',
  oidcProvider: '',

  get namespaceInput() {
    const namespaceQP = this.clusterController.namespaceQueryParam;
    if (this.hvdManagedNamespaceRoot) {
      // When managed, the user isn't allowed to edit the prefix `admin/` for their nested namespace
      const split = namespaceQP.split('/');
      if (split.length > 1) {
        split.shift();
        return `/${split.join('/')}`;
      }
      return '';
    }
    return namespaceQP;
  },

  fullNamespaceFromInput(value) {
    const strippedNs = sanitizePath(value);
    if (this.hvdManagedNamespaceRoot) {
      return `${this.hvdManagedNamespaceRoot}/${strippedNs}`;
    }
    return strippedNs;
  },

  updateNamespace: task(function* (value) {
    // debounce
    yield timeout(500);
    const ns = this.fullNamespaceFromInput(value);
    this.namespaceService.setNamespace(ns, true);
    this.customMessages.fetchMessages(ns);
    this.set('namespaceQueryParam', ns);
  }).restartable(),

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
        this.customMessages.fetchMessages(namespace);
      }

      if (isRoot) {
        this.auth.set('isRootToken', true);
        this.flashMessages.warning(
          'You have logged in with a root token. As a security precaution, this root token will not be stored by your browser and you will need to re-authenticate after the window is closed or refreshed.'
        );
      }
    });
  },

  actions: {
    onAuthResponse(authResponse, backend, data) {
      const { mfa_requirement } = authResponse;
      // if an mfa requirement exists further action is required
      if (mfa_requirement) {
        this.set('mfaAuthData', { mfa_requirement, backend, data });
      } else {
        this.authSuccess(authResponse);
      }
    },
    onMfaSuccess(authResponse) {
      this.authSuccess(authResponse);
    },
    onMfaErrorDismiss() {
      this.setProperties({
        mfaAuthData: null,
        mfaErrors: null,
      });
    },
    cancelAuthentication() {
      this.set('cancelAuth', true);
      this.set('waitingForOktaNumberChallenge', false);
    },
  },
});
