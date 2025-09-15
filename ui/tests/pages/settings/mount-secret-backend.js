/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { create, visitable, fillable, clickable } from 'ember-cli-page-object';
import { visit, click, fillIn } from '@ember/test-helpers';
import { GENERAL } from 'vault/tests/helpers/general-selectors';

export default create({
  visit: visitable('/vault/secrets/mounts'),
  version: fillable('[data-test-input="options.version"]'),
  setMaxVersion: fillable('[data-test-input="kv_config.max_versions"]'),
  maxTTLVal: fillable('[data-test-ttl-value="Max Lease TTL"]'),
  maxTTLUnit: fillable('[data-test-ttl-unit="Max Lease TTL"] [data-test-select="ttl-unit"]'),
  enableEngine: clickable('[data-test-enable-engine]'),
  secretList: clickable('[data-test-sidebar-nav-link="Secrets Engines"]'),
  defaultTTLVal: fillable('input[data-test-ttl-value="Default Lease TTL"]'),
  defaultTTLUnit: fillable('[data-test-ttl-unit="Default Lease TTL"] [data-test-select="ttl-unit"]'),
  enable: async function (type, path) {
    // Navigate to the secrets engines catalog
    await visit('/vault/secrets/mounts');
    // Click the engine type card to proceed to configuration
    await click(GENERAL.cardContainer(type));
    // Fill in the path if provided
    if (path) {
      await fillIn(GENERAL.inputByAttr('path'), path);
    }
    // Submit the form to mount the engine
    await click(GENERAL.submitButton);
  },
});
