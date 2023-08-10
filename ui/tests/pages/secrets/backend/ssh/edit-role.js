/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
 */

import { Base } from '../create';
import { clickable, visitable, create, fillable } from 'ember-cli-page-object';

export default create({
  ...Base,
  visitEdit: visitable('/vault/secrets/:backend/edit/:id'),
  visitEditRoot: visitable('/vault/secrets/:backend/edit'),
  keyType: fillable('[data-test-input="keyType"]'),
  defaultUser: fillable('[data-test-input="defaultUser"]'),
  toggleMore: clickable('[data-test-toggle-group="Options"]'),
  name: fillable('[data-test-input="name"]'),
  CIDR: fillable('[data-test-input="cidrList"]'),
  save: clickable('[data-test-role-ssh-create]'),

  async createOTPRole(name) {
    await this.name(name);
    await this.toggleMore().keyType('otp').defaultUser('admin').CIDR('0.0.0.0/0').save();
  },
});
