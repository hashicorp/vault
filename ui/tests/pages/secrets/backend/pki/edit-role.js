/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
 */

import { Base } from '../create';
import { settled } from '@ember/test-helpers';
import { clickable, visitable, create, fillable } from 'ember-cli-page-object';

export default create({
  ...Base,
  visitEdit: visitable('/vault/secrets/:backend/edit/:id'),
  visitEditRoot: visitable('/vault/secrets/:backend/edit'),
  toggleDomain: clickable('[data-test-toggle-group="Domain Handling"]'),
  toggleOptions: clickable('[data-test-toggle-group="Options"]'),
  name: fillable('[data-test-input="name"]'),
  allowAnyName: clickable('[data-test-input="allowAnyName"]'),
  allowedDomains: fillable('[data-test-input="allowedDomains"] .input'),
  save: clickable('[data-test-role-create]'),

  async createRole(name, allowedDomains) {
    await this.toggleDomain();
    await settled();
    await this.toggleOptions();
    await settled();
    await this.name(name);
    await settled();
    await this.allowAnyName();
    await settled();
    await this.allowedDomains(allowedDomains);
    await settled();
    await this.save();
    await settled();
  },
});
