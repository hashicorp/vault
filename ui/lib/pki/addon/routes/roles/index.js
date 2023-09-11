/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import Route from '@ember/routing/route';
import { inject as service } from '@ember/service';
import { withConfig } from 'pki/decorators/check-issuers';
import { hash } from 'rsvp';
import { getCliMessage } from 'pki/routes/overview';
@withConfig()
export default class PkiRolesIndexRoute extends Route {
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
      hasConfig: this.shouldPromptConfig,
      roles: this.fetchRoles(),
      parentModel: this.modelFor('roles'),
    });
  }

  setupController(controller, resolvedModel) {
    super.setupController(controller, resolvedModel);
    const roles = resolvedModel.roles;

    if (roles?.length) controller.notConfiguredMessage = getCliMessage('roles');
    else controller.notConfiguredMessage = getCliMessage();
  }
}
