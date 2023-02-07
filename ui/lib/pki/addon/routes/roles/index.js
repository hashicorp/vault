/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
 */

import PkiOverviewRoute from '../overview';
import { inject as service } from '@ember/service';
import { hash } from 'rsvp';

export default class PkiRolesIndexRoute extends PkiOverviewRoute {
  @service store;
  @service secretMountPath;

  async fetchRoles() {
    try {
      return await this.store.query('pki/role', { backend: this.secretMountPath.currentPath });
    } catch (e) {
      if (e.httpStatus === 404) {
        return { parentModel: this.modelFor('roles') };
      } else {
        throw e;
      }
    }
  }

  model() {
    return hash({
      hasConfig: this.hasConfig(),
      roles: this.fetchRoles(),
      parentModel: this.modelFor('roles'),
    });
  }
}
