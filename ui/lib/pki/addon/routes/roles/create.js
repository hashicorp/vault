/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import Route from '@ember/routing/route';
import { service } from '@ember/service';
import { withConfirmLeave } from 'core/decorators/confirm-leave';
import { hash } from 'rsvp';

@withConfirmLeave('model.role', ['model.issuers'])
export default class PkiRolesCreateRoute extends Route {
  @service store;
  @service secretMountPath;

  model() {
    const backend = this.secretMountPath.currentPath;
    return hash({
      role: this.store.createRecord('pki/role', { backend }),
      issuers: this.store.query('pki/issuer', { backend }).catch((err) => {
        if (err.httpStatus === 404) {
          return [];
        } else {
          throw err;
        }
      }),
    });
  }

  setupController(controller, resolvedModel) {
    super.setupController(controller, resolvedModel);
    controller.breadcrumbs = [
      { label: 'Secrets', route: 'secrets', linkExternal: true },
      { label: this.secretMountPath.currentPath, route: 'overview', model: this.secretMountPath.currentPath },
      { label: 'Roles', route: 'roles.index', model: this.secretMountPath.currentPath },
      { label: 'Create' },
    ];
  }

  willTransition() {
    // after upgrading to Ember Data 5.3.2 we saw duplicate records in the store after creating and saving a new role
    // it's unclear why this ghost record is persisting, manually unloading refreshes the store
    this.store.unloadAll('pki/role');
  }
}
