/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { create, visitable, fillable, clickable } from 'ember-cli-page-object';
import { settled } from '@ember/test-helpers';
import { mountBackend } from 'vault/tests/helpers/components/mount-backend-form-helpers';

export default create({
  visit: visitable('/vault/settings/mount-secret-backend'),
  version: fillable('[data-test-input="version"]'),
  setMaxVersion: fillable('[data-test-input="maxVersions"]'),
  enableMaxTtl: clickable('[data-test-toggle-input="Max Lease TTL"]'),
  maxTTLVal: fillable('[data-test-ttl-value="Max Lease TTL"]'),
  maxTTLUnit: fillable('[data-test-ttl-unit="Max Lease TTL"] [data-test-select="ttl-unit"]'),
  enableDefaultTtl: clickable('[data-test-toggle-input="Default Lease TTL"]'),
  enableEngine: clickable('[data-test-enable-engine]'),
  secretList: clickable('[data-test-sidebar-nav-link="Secrets Engines"]'),
  defaultTTLVal: fillable('input[data-test-ttl-value="Default Lease TTL"]'),
  defaultTTLUnit: fillable('[data-test-ttl-unit="Default Lease TTL"] [data-test-select="ttl-unit"]'),
  enable: async function (type, path) {
    await this.visit();
    await settled();
    await mountBackend(type, path);
    await settled();
  },
});
