/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import RoleEdit from './role-edit';
import { action } from '@ember/object';
import { tracked } from '@glimmer/tracking';

export default class TotpEdit extends RoleEdit {
  @tracked hasGenerated = false;
  successCallback;

  init() {
    super.init(...arguments);
    this.set('backendType', 'totp');
  }

  persist(method, successCallback) {
    const model = this.model;
    return model[method]().then(() => {
      if (!model.isError) {
        if (model.backend === 'totp' && model.generate) {
          this.hasGenerated = true;
          this.successCallback = successCallback;
        } else {
          successCallback(model);
        }
      }
    });
  }

  @action
  reset() {
    this.model.unloadRecord();
    this.successCallback(null);
  }
}
