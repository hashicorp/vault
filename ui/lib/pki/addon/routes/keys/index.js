/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
 */

import Route from '@ember/routing/route';
import { inject as service } from '@ember/service';
import { withConfig } from 'pki/decorators/check-config';
import { hash } from 'rsvp';
import { PKI_DEFAULT_EMPTY_STATE_MSG } from 'pki/routes/overview';

@withConfig()
export default class PkiKeysIndexRoute extends Route {
  @service secretMountPath;
  @service store;

  model() {
    return hash({
      hasConfig: this.shouldPromptConfig,
      parentModel: this.modelFor('keys'),
      keyModels: this.store.query('pki/key', { backend: this.secretMountPath.currentPath }).catch((err) => {
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
      { label: 'secrets', route: 'secrets', linkExternal: true },
      { label: this.secretMountPath.currentPath, route: 'overview' },
      { label: 'keys', route: 'keys.index' },
    ];
    controller.message = PKI_DEFAULT_EMPTY_STATE_MSG;
  }
}
