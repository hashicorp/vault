/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
 */

import { module, test } from 'qunit';
import { setupTest } from 'vault/tests/helpers';
import { setupMirage } from 'ember-cli-mirage/test-support';
import { allowAllCapabilitiesStub } from 'vault/tests/helpers/stubs';

module('Unit | Adapter | pki/tidy', function (hooks) {
  setupTest(hooks);
  setupMirage(hooks);

  hooks.beforeEach(function () {
    this.store = this.owner.lookup('service:store');
    this.secretMountPath = this.owner.lookup('service:secret-mount-path');
    this.backend = 'pki-test';
    this.secretMountPath.currentPath = this.backend;
    this.server.post('/sys/capabilities-self', allowAllCapabilitiesStub());
  });

  test('it exists', function (assert) {
    const adapter = this.owner.lookup('adapter:pki/tidy');
    assert.ok(adapter);
  });

  test('it calls the correct endpoint when tidyType = manual-tidy', async function (assert) {
    assert.expect(1);

    this.server.post(`${this.backend}/tidy`, () => {
      assert.ok(true, 'request made to correct endpoint on create');
      return {};
    });
    this.payload = {
      tidy_cert_store: true,
      tidy_revocation_queue: false,
      safetyBuffer: '120h',
      backend: this.backend,
    };
    await this.store
      .createRecord('pki/tidy', this.payload)
      .save({ adapterOptions: { tidyType: 'manual-tidy' } });
  });

  test('it calls the correct endpoint when tidyType = auto-tidy', async function (assert) {
    assert.expect(1);
    this.server.post(`${this.backend}/config/auto-tidy`, () => {
      assert.ok(true, 'request made to correct endpoint on create');
      return {};
    });
    this.payload = {
      enabled: true,
      interval_duration: '72h',
      backend: this.backend,
    };
    await this.store
      .createRecord('pki/tidy', this.payload)
      .save({ adapterOptions: { tidyType: 'auto-tidy' } });
  });
});
